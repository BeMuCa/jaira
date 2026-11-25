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
	// ours is always the text of a merge onto the ref at oursLease: that pair
	// is what goes back in the queue when it is not sent. Leasing it on any
	// newer state would let the next flush push it over a write it never saw.
	name, ours, oursLease := r.ID, []byte(e.Content), e.Lease
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
			// Only the entry read above goes: an edit queued during the push
			// was made on ours, not on the merge, and stays leased where it
			// is, to lose the next push and be merged in turn.
			r.Outcome, r.Err = outbox.Sent, y.Box.Settle(e)
			return r
		case errors.Is(err, gitref.ErrRaceLost):
			// Somebody wrote it again between our read and our push: their
			// state is now the one to merge onto, ours stays what it was.
			base, ours, oursLease = theirs, merged, lease
			if ferr := y.fetch(); ferr != nil {
				err = ferr
			} else {
				continue
			}
		}
		r.Err = err
		break
	}
	// Not sent: queue what we have, leased on the ref it was merged onto, so
	// the next command tries again and loses the race to every write it has
	// not merged yet — not on the ref as it stands now, which after the fetch
	// above may hold a write this text never saw. It replaces only the entry
	// read above: an edit queued since stays, and is merged on the next try.
	if err := y.Box.Replace(e, ours, oursLease, y.Actor); err != nil {
		r.Err = err
	} else {
		y.dirty = true
	}
	r.Outcome = outbox.Unsent
	return r
}
