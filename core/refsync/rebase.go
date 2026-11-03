package refsync

import (
	"errors"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
)

// rebaseAttempts bounds how often one flush re-merges an epic or a road that
// keeps losing. Each loss means somebody else wrote it in the meantime, so a few
// rounds are plenty; beyond that the write goes back in the queue.
const rebaseAttempts = 3

// refKind is what an epic and a road differ in: their outbox queue, their refs,
// and the merge that keeps both sides' member changes.
type refKind struct {
	pending func(name string) (outbox.Entry, bool)
	queue   func(name string, content []byte, lease, by string) error
	drop    func(name string) error
	record  func(name string, content []byte) error
	readAt  func(name, sha string) ([]byte, error)
	read    func(name string) ([]byte, string, error)
	write   func(name string, content []byte, lease string) (string, error)
	merge   func(name string, base, ours, theirs []byte) []byte
}

// rebaseRef answers a lost race for an epic or a road by merging instead of
// dropping the edit. The queued entry is the only copy of the edit, so dropping
// it would lose a ticket somebody sorted in. Flush leaves it queued; the merge
// keeps both sides' member changes, is sent again at once, and replaces the
// entry only once it is sent or requeued.
func (y *Syncer) rebaseRef(r outbox.Result, k refKind) outbox.Result {
	e, ok := k.pending(r.ID)
	if !ok {
		return r
	}
	name, ours := r.ID, []byte(e.Content)
	var base []byte
	if e.Lease != "" {
		b, err := k.readAt(name, e.Lease)
		if err != nil {
			// The common state is unreadable, so there is nothing to merge
			// against; the edit stays queued rather than guessed at.
			r.Outcome, r.Err = outbox.Failed, err
			return r
		}
		base = b
	}
	for range rebaseAttempts {
		theirs, lease, err := k.read(name)
		if err != nil && !errors.Is(err, gitref.ErrNoRef) {
			r.Err = err
			break
		}
		merged := k.merge(name, base, ours, theirs)
		_, err = k.write(name, merged, lease)
		switch {
		case err == nil:
			r.Outcome, r.Err = outbox.Sent, k.drop(name)
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
	// the next command tries again. The box keeps the lease of an entry that
	// is there, so the refused one goes first, and comes back if the merge
	// cannot be queued.
	if err := k.drop(name); err != nil {
		r.Err = err
	} else if err := k.record(name, ours); err != nil {
		_ = k.queue(name, []byte(e.Content), e.Lease, e.By)
		r.Err = err
	}
	r.Outcome = outbox.Unsent
	return r
}
