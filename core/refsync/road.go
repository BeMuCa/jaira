package refsync

import (
	"errors"
	"fmt"
	"sort"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/road"
)

// The Syncer is where roads are kept, as it is where epics are: they live on
// refs only, so the ref and the outbox together are the whole store. These
// three methods are road.Refs.

// RoadNames lists every road this clone knows: the local refs, and the roads
// written here that have not reached the remote yet.
func (y *Syncer) RoadNames() ([]string, error) {
	if y == nil || y.Usable() != nil {
		return nil, road.ErrNoRefs
	}
	names, err := y.Repo.ListRoads()
	if err != nil {
		return nil, err
	}
	entries, err := y.Box.List()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, n := range names {
		seen[n] = true
	}
	for _, e := range entries {
		if e.Kind == outbox.KindRoad && !seen[e.ID] {
			seen[e.ID] = true
			names = append(names, e.ID)
		}
	}
	sort.Strings(names)
	return names, nil
}

// RoadContent reads one road. An unsent write wins over the ref, for the same
// reason an unsent epic does: it is newer, and it is the only copy.
func (y *Syncer) RoadContent(name string) ([]byte, error) {
	if y == nil || y.Usable() != nil {
		return nil, road.ErrNoRefs
	}
	if e, ok := y.Box.PendingRoad(name); ok {
		return []byte(e.Content), nil
	}
	content, _, err := y.Repo.ReadRoad(name)
	if errors.Is(err, gitref.ErrNoRef) {
		return nil, fmt.Errorf("%w: %s", road.ErrNotFound, name)
	}
	return content, err
}

// RecordRoad queues a road's whole content for its ref, leased against the SHA
// the remote last confirmed, the same rule an epic write follows.
func (y *Syncer) RecordRoad(name string, content []byte) error {
	if y == nil || y.Usable() != nil {
		return road.ErrNoRefs
	}
	lease, err := y.Repo.RoadSHA(name)
	if err != nil && !errors.Is(err, gitref.ErrNoRef) {
		return err
	}
	if err := y.Box.QueueRoad(name, content, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// rebaseRoad answers a lost race for a road by merging, as rebaseEpic does for
// an epic: the queued entry is the only copy of the edit, so it is merged
// (road.Merge keeps both sides' member changes), sent again, and replaces the
// entry only once it is sent or requeued.
func (y *Syncer) rebaseRoad(r outbox.Result) outbox.Result {
	e, ok := y.Box.PendingRoad(r.ID)
	if !ok {
		return r
	}
	name, ours := r.ID, []byte(e.Content)
	var base []byte
	if e.Lease != "" {
		b, err := y.Repo.ReadRoadAt(name, e.Lease)
		if err != nil {
			// The common state is unreadable, so there is nothing to merge
			// against; the edit stays queued rather than guessed at.
			r.Outcome, r.Err = outbox.Failed, err
			return r
		}
		base = b
	}
	for range rebaseAttempts {
		theirs, lease, err := y.Repo.ReadRoad(name)
		if err != nil && !errors.Is(err, gitref.ErrNoRef) {
			r.Err = err
			break
		}
		merged := road.Merge(name, base, ours, theirs)
		_, err = y.Repo.WriteRoad(name, merged, lease)
		switch {
		case err == nil:
			r.Outcome, r.Err = outbox.Sent, y.Box.DropRoad(name)
			return r
		case errors.Is(err, gitref.ErrRaceLost):
			// Somebody wrote it again between our read and our push: their
			// state is now the one to merge onto, ours stays what it was.
			base, ours = theirs, merged
			if ferr := y.Repo.Fetch(); ferr != nil {
				err = ferr
			} else {
				continue
			}
		}
		r.Err = err
		break
	}
	// Not sent: queue what we have, leased on the ref it was merged onto, so
	// the next command tries again.
	if err := y.Box.DropRoad(name); err != nil {
		r.Err = err
	} else if err := y.RecordRoad(name, ours); err != nil {
		_ = y.Box.QueueRoad(name, []byte(e.Content), e.Lease, e.By)
		r.Err = err
	}
	r.Outcome = outbox.Unsent
	return r
}
