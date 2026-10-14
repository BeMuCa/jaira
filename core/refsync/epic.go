package refsync

import (
	"errors"
	"fmt"
	"sort"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
)

// The Syncer is where epics are kept: they live on refs only, so the ref and
// the outbox together are the whole store. These three methods are
// epic.Refs.

// EpicNames lists every epic this clone knows: the local refs, and the epics
// written here that have not reached the remote yet.
func (y *Syncer) EpicNames() ([]string, error) {
	if y == nil || y.Usable() != nil {
		return nil, epic.ErrNoRefs
	}
	names, err := y.Repo.ListEpics()
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
		if e.Kind == outbox.KindEpic && !seen[e.ID] {
			seen[e.ID] = true
			names = append(names, e.ID)
		}
	}
	sort.Strings(names)
	return names, nil
}

// EpicContent reads one epic. An unsent write wins over the ref: it is newer,
// and with no file on disk it is the only copy of that edit.
func (y *Syncer) EpicContent(name string) ([]byte, error) {
	if y == nil || y.Usable() != nil {
		return nil, epic.ErrNoRefs
	}
	if e, ok := y.Box.PendingEpic(name); ok {
		return []byte(e.Content), nil
	}
	content, _, err := y.Repo.ReadEpic(name)
	if errors.Is(err, gitref.ErrNoRef) {
		return nil, fmt.Errorf("%w: %s", epic.ErrNotFound, name)
	}
	return content, err
}

// RecordEpic queues an epic's whole content for its ref, leased against the
// SHA the remote last confirmed — the same rule a ticket write follows.
func (y *Syncer) RecordEpic(name string, content []byte) error {
	if y == nil || y.Usable() != nil {
		return epic.ErrNoRefs
	}
	lease, err := y.Repo.EpicSHA(name)
	if err != nil && !errors.Is(err, gitref.ErrNoRef) {
		return err
	}
	if err := y.Box.QueueEpic(name, content, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// rebaseAttempts bounds how often one flush re-merges an epic that keeps
// losing. Each loss means somebody else wrote it in the meantime, so a few
// rounds are plenty; beyond that the write goes back in the queue.
const rebaseAttempts = 3

// rebaseEpic answers a lost race for an epic by merging instead of dropping the
// edit. A ticket's loser keeps its local file and the board shows both sides;
// an epic has no local file, so the queued entry was the only copy of the edit
// and dropping it would lose a ticket somebody sorted in. The merge keeps both
// sides' member changes (epic.Merge) and is sent again at once.
func (y *Syncer) rebaseEpic(r outbox.Result) outbox.Result {
	name, ours := r.ID, []byte(r.Entry.Content)
	var base []byte
	if r.Entry.Lease != "" {
		b, err := y.Repo.ReadEpicAt(name, r.Entry.Lease)
		if err != nil {
			// The common state is unreadable, so there is nothing to merge
			// against; keep the edit queued rather than guess.
			r.Outcome, r.Err = outbox.Failed, err
			_ = y.Box.QueueEpic(name, ours, r.Entry.Lease, r.Entry.By)
			return r
		}
		base = b
	}
	for range rebaseAttempts {
		theirs, lease, err := y.Repo.ReadEpic(name)
		if err != nil && !errors.Is(err, gitref.ErrNoRef) {
			r.Err = err
			break
		}
		merged := epic.Merge(name, base, ours, theirs)
		_, err = y.Repo.WriteEpic(name, merged, lease)
		switch {
		case err == nil:
			r.Outcome, r.Err = outbox.Sent, nil
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
	// Not sent: queue what we have so the next command tries again.
	if err := y.RecordEpic(name, ours); err != nil {
		r.Err = err
	}
	r.Outcome = outbox.Unsent
	return r
}
