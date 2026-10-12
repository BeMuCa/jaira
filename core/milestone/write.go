package milestone

import (
	"bytes"
	"fmt"
)

// LockName keys the store lock every milestone write is made under. One name
// for the whole directory rather than one per file: creating a milestone
// reads every existing one to avoid their colours, so two creations at once
// have to be ordered even though they write different files.
const LockName = "milestones"

// Locker is the one thing of the ticket store a write needs. An interface so
// this package does not import the store's locking and a test can count.
type Locker interface {
	Lock(name string) (unlock func(), err error)
}

// Recorder queues a milestone file for its ref; refsync's Syncer is one. A
// nil *Syncer is a valid Recorder that does nothing, which is how a board with
// no remote stays a supported way to work.
type Recorder interface {
	RecordMilestone(name string, content []byte) error
}

// RecordError reports that the file was saved and only the ref could not be
// written. The milestone is returned beside it: the write happened, so a
// caller that wants the old best-effort behaviour downgrades this to a warning
// and carries on.
type RecordError struct {
	Name string
	Err  error
}

func (e *RecordError) Error() string { return fmt.Sprintf("milestone %q: %v", e.Name, e.Err) }
func (e *RecordError) Unwrap() error { return e.Err }

// Edit is the one write path for a milestone that exists: lock, read, let fn
// change it, save atomically, record the ref. The CLI and the board both
// write through it so they cannot drift — what a person does in the layer is
// exactly what an agent does with the command.
//
// fn may refuse by returning an error; nothing is written then. The file is
// saved only when fn changed its bytes, because a re-add of a ticket already
// there is not a write and must not queue a ref. A missing file comes back as
// the Load error (errors.Is os.ErrNotExist) so the caller can say why in its
// own words.
func Edit(root string, lock Locker, rec Recorder, name string, fn func(*Milestone) error) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()

	ms, err := Load(root, name)
	if err != nil {
		return nil, err
	}
	before := ms.Bytes()
	if err := fn(ms); err != nil {
		return ms, err
	}
	if bytes.Equal(before, ms.Bytes()) {
		return ms, nil
	}
	return ms, saveAndRecord(root, rec, ms)
}

// Create is Edit's counterpart for a milestone that does not exist yet. build
// runs under the lock with every milestone on the board, because choosing a
// colour and checking the name is free are both reads that must not race a
// second creation; it returns the milestone to write, or refuses with an error
// and nothing is written.
func Create(root string, lock Locker, rec Recorder, name string, build func(existing []*Milestone) (*Milestone, error)) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()

	existing, err := LoadAll(root)
	if err != nil {
		return nil, err
	}
	ms, err := build(existing)
	if err != nil {
		return nil, err
	}
	return ms, saveAndRecord(root, rec, ms)
}

func saveAndRecord(root string, rec Recorder, ms *Milestone) error {
	if err := ms.Save(root); err != nil {
		return err
	}
	if err := rec.RecordMilestone(ms.Name, ms.Bytes()); err != nil {
		return &RecordError{Name: ms.Name, Err: err}
	}
	return nil
}
