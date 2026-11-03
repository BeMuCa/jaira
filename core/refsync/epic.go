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

// rebaseEpic answers a lost race for an epic by merging it (epic.Merge), since
// an epic has no local file to fall back on.
func (y *Syncer) rebaseEpic(r outbox.Result) outbox.Result {
	return y.rebaseRef(r, refKind{
		pending: y.Box.PendingEpic,
		queue:   y.Box.QueueEpic,
		drop:    y.Box.DropEpic,
		record:  y.RecordEpic,
		readAt:  y.Repo.ReadEpicAt,
		read:    y.Repo.ReadEpic,
		write:   y.Repo.WriteEpic,
		merge:   epic.Merge,
	})
}
