package milestone

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
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
	RecordMilestone(key, file string, content []byte) error
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

// FiledError is Edit's refusal of a milestone whose file still lies here but
// is marked filed — what a fetch of somebody else's filing leaves behind.
// Writing it would put the change on its ref for every clone that fetches.
// RenamedTo is set when the mark is the one a rename left: then there is
// nothing to restore, and the new name is the thing to say.
type FiledError struct{ Name, RenamedTo string }

func (e *FiledError) Error() string {
	if e.RenamedTo != "" {
		return fmt.Sprintf("milestone %q was renamed to %q", e.Name, e.RenamedTo)
	}
	return fmt.Sprintf("milestone %q has been filed; 'jaira restore %s.md', run where it was filed, brings it back", e.Name, e.Name)
}

// ExistsError is CheckFree's answer for a name a milestone on the board
// already has.
type ExistsError struct{ Name string }

func (e *ExistsError) Error() string { return fmt.Sprintf("milestone %q already exists", e.Name) }

// Edit is the one write path for a milestone that exists: lock, read, let fn
// change it, save atomically, record the ref. The CLI and the board both
// write through it so they cannot drift — what a person does in the layer is
// exactly what an agent does with the command.
//
// fn may refuse by returning an error; nothing is written then. The file is
// saved only when fn changed its bytes, because a re-add of a ticket already
// there is not a write and must not queue a ref. A missing file comes back as
// the Load error (errors.Is os.ErrNotExist) so the caller can say why in its
// own words, and a filed one as a *FiledError before fn runs.
//
// When fn takes a placed milestone out of the layer, Edit also closes the
// chain behind it under the same lock: every milestone that came after it now
// comes after what came before it (ReattachSuccessors). Here and not in the
// callers, so the layer's x and 'jaira milestone set --no-roadmap' cannot leave
// two different pictures behind.
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
	if ms.Filed() {
		return ms, &FiledError{Name: ms.Name, RenamedTo: ms.RenamedTo}
	}
	before, wasPlaced := ms.Bytes(), ms.Placed
	if err := fn(ms); err != nil {
		return ms, err
	}
	if bytes.Equal(before, ms.Bytes()) {
		return ms, nil
	}
	var succ []*Milestone
	if wasPlaced && !ms.Placed {
		// Read before the save, while this milestone's own after line is
		// still on disk: it is what its successors move up to.
		all, err := LoadAll(root)
		if err != nil {
			return ms, err
		}
		reattach := ReattachSuccessors(all, ms.ID)
		for _, m := range all { // LoadAll's order, by name
			if after, ok := reattach[m.ID]; ok {
				m.SetAfter(after)
				succ = append(succ, m)
			}
		}
	}
	return ms, saveAll(root, rec, append([]*Milestone{ms}, succ...))
}

// Rename gives the milestone from the name to. The name is the file name and
// the ref name, so a rename is three writes under the one lock: the file
// under its new name, every successor's after line naming the new name where
// it named the old one, and a mark under the old name — the old file marked
// filed, with renamed-to: <to>, put on the old ref. The old file goes from
// this disk; the mark is for every other clone. Taking the old ref down
// instead would reach none of them: a fetch writes a marked ref over the file
// a clone already has, and deletes nothing (refsync.IncomingMilestones).
//
// to has to be free: CheckFree asks the directory, and taken asks what the
// directory cannot show — a filing in this tree's logbook or on a ref, or a
// teammate's milestone on a ref no fetch has written here — refusing in the
// caller's words. from itself refuses the way Edit does, and
// renaming it to its own name writes nothing.
func Rename(root string, lock Locker, rec Recorder, from, to string, taken func(name string) error) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()

	ms, err := Load(root, from)
	if err != nil {
		return nil, err
	}
	if ms.Filed() {
		return ms, &FiledError{Name: from, RenamedTo: ms.RenamedTo}
	}
	if to == from {
		return ms, nil
	}
	if err := CheckFree(root, to); err != nil {
		return ms, err
	}
	if err := taken(to); err != nil {
		return ms, err
	}

	all, err := LoadAll(root)
	if err != nil {
		return ms, err
	}
	renamed := FromBytes(to, ms.Bytes())
	renamed.retitle(from, to)
	writes := []*Milestone{renamed}
	r := newResolver(all)
	for _, m := range all { // LoadAll's order, by name
		if m.ID == ms.ID {
			continue
		}
		after, changed := slices.Clone(m.After), false
		for i, tok := range after {
			if r.resolve(tok) != ms.ID {
				continue
			}
			// The renamed milestone's ID is the one the new name gives it,
			// until IDs stop following names.
			after[i], changed = to, true
			if IsID(tok) {
				after[i] = renamed.ID
			}
		}
		if changed {
			m.SetAfter(after)
			writes = append(writes, m)
		}
	}
	// The new file before the old one goes: a rename cut short leaves the
	// milestone twice, never nowhere.
	err = saveAll(root, rec, writes)
	if _, recordOnly := err.(*RecordError); err != nil && !recordOnly {
		return renamed, err
	}
	if rmErr := os.Remove(ms.Path(root)); rmErr != nil {
		return renamed, rmErr
	}
	ms.SetStatus(StatusFiled)
	ms.RenamedTo = to
	ms.setFront("renamed-to", to, false)
	if recErr := rec.RecordMilestone(from, from+".md", ms.Bytes()); recErr != nil && err == nil {
		err = &RecordError{Name: from, Err: recErr}
	}
	return renamed, err
}

// CheckFree asks the milestone directory whether a name can be given out: a
// milestone on the board under it is an *ExistsError, a filed one lying here
// a *FiledError. No file is free, and so is the mark an earlier rename left
// (HoldsName) — which is what lets a rename be taken back and a new milestone
// take the old name. What the directory cannot show, a filing in this tree's
// logbook or on a ref or a milestone on a ref not fetched yet, is the
// caller's to ask. Every door that hands out a
// name asks here: Rename, and the CLI's and the layer's create.
func CheckFree(root, name string) error {
	have, err := Load(root, name)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return err
	case !have.Filed():
		return &ExistsError{Name: name}
	case have.HoldsName():
		return &FiledError{Name: name}
	}
	return nil
}

// retitle rewrites the "# <name>" heading New puts under the frontmatter, so
// a renamed file does not open with its old name. Only that exact line: any
// other is somebody's prose.
func (m *Milestone) retitle(from, to string) {
	for i, line := range m.lines {
		if line == "# "+from {
			m.lines[i] = "# " + to
			return
		}
	}
}

// Create is Edit's counterpart for a milestone that does not exist yet. build
// runs under the lock with every milestone on the board, because choosing a
// colour and checking the name is free are both reads that must not race a
// second creation; it returns the milestone to write, or refuses with an error
// and nothing is written.
func Create(root string, lock Locker, rec Recorder, build func(existing []*Milestone) (*Milestone, error)) (*Milestone, error) {
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

// saveAll saves every milestone and queues each for its ref. Every file is
// saved before a ref failure is reported: the ref is best effort, the save is
// the write. A failed save stops it; otherwise the first ref failure comes
// back, as a *RecordError.
func saveAll(root string, rec Recorder, ms []*Milestone) error {
	var recErr error
	for _, m := range ms {
		switch err := saveAndRecord(root, rec, m).(type) {
		case nil:
		case *RecordError:
			if recErr == nil {
				recErr = err
			}
		default:
			return err
		}
	}
	return recErr
}

func saveAndRecord(root string, rec Recorder, ms *Milestone) error {
	if err := ms.Save(root); err != nil {
		return err
	}
	if err := rec.RecordMilestone(ms.Name, ms.Name+".md", ms.Bytes()); err != nil {
		return &RecordError{Name: ms.Name, Err: err}
	}
	return nil
}
