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

// rebaseRoad answers a lost race for a road by merging it (road.Merge), as
// rebaseEpic does for an epic.
func (y *Syncer) rebaseRoad(r outbox.Result) outbox.Result {
	return y.rebaseRef(r, refKind{
		pending: y.Box.PendingRoad,
		queue:   y.Box.QueueRoad,
		drop:    y.Box.DropRoad,
		record:  y.RecordRoad,
		readAt:  y.Repo.ReadRoadAt,
		read:    y.Repo.ReadRoad,
		write:   y.Repo.WriteRoad,
		merge:   road.Merge,
	})
}
