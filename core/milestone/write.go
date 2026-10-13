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
	// HasMilestoneRef reports whether the milestone with this ID already has
	// its ref in this tree.
	HasMilestoneRef(id string) bool
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
type FiledError struct {
	Name string
	File string // the file it lies in, which is what 'jaira restore' takes
}

func (e *FiledError) Error() string {
	return fmt.Sprintf("milestone %q has been filed; 'jaira restore %s', run where it was filed, brings it back", e.Name, e.File)
}

// ExistsError is CheckFree's answer for a name a milestone on the board
// already has.
type ExistsError struct {
	Name string
	File string // the file of the one that has it; of two, the first by ID
}

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
	unlock, err := lockAndMigrate(root, lock, rec)
	if err != nil {
		return nil, err
	}
	defer unlock()

	ms, err := Load(root, name)
	if err != nil {
		return nil, err
	}
	if ms.Filed() {
		return ms, &FiledError{Name: ms.Name, File: ms.FileName()}
	}
	before, wasPlaced := ms.Bytes(), ms.Placed
	if err := fn(ms); err != nil {
		return ms, err
	}
	if bytes.Equal(before, ms.Bytes()) {
		return ms, nil
	}
	if ms.Placed {
		// A name a hand wrote in after: becomes the ID it names, now that the
		// file is being written anyway.
		all, err := LoadAll(root)
		if err != nil {
			return ms, err
		}
		if ids := SpellAsIDs(all, ms.After); !slices.Equal(ids, ms.After) {
			ms.SetAfter(ids)
		}
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

// lockAndMigrate takes the milestone lock and, under it, moves any file
// written before IDs existed to its ID name (Migrate), so a write never meets
// the old format and no caller can forget to. A ref that could not be queued
// for a moved file is not a reason to stop: the move happened, and the next
// fetch or write queues it again.
func lockAndMigrate(root string, lock Locker, rec Recorder) (func(), error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	if _, err := Migrate(root, rec); err != nil {
		var re *RecordError
		if !errors.As(err, &re) {
			unlock()
			return nil, err
		}
	}
	return unlock, nil
}

// SpellAsIDs gives every token that names exactly one milestone as that
// milestone's ID, and leaves the others as they are.
func SpellAsIDs(all []*Milestone, tokens []string) []string {
	r := newResolver(all)
	out := make([]string, len(tokens))
	for i, tok := range tokens {
		out[i] = r.spell(tok)
	}
	return out
}

// Rename gives the milestone, found by whatever names it, the name to. The
// ID stays, so the milestone keeps its ref and its place in every after:
// line: a rename is the file under its new name, the old file gone, and one
// write on the one ref. Nothing else is written, and nobody who follows the
// milestone is touched.
//
// to has to be a name no milestone on the board carries (CheckFree). The
// milestone refuses the way Edit does when it is filed, and renaming it to
// its own name writes nothing.
func Rename(root string, lock Locker, rec Recorder, from, to string) (*Milestone, error) {
	unlock, err := lockAndMigrate(root, lock, rec)
	if err != nil {
		return nil, err
	}
	defer unlock()

	ms, err := Load(root, from)
	if err != nil {
		return nil, err
	}
	if ms.Filed() {
		return ms, &FiledError{Name: ms.Name, File: ms.FileName()}
	}
	if to == ms.Name {
		return ms, nil
	}
	if err := CheckFree(root, to); err != nil {
		return ms, err
	}
	renamed := FromBytes(to+"-"+ms.ID+".md", ms.Bytes())
	renamed.retitle(ms.Name, to)
	// The new file before the old one goes: a rename cut short leaves the
	// milestone twice, never nowhere.
	err = saveAndRecord(root, rec, renamed)
	if _, recordOnly := err.(*RecordError); err != nil && !recordOnly {
		return renamed, err
	}
	if rmErr := os.Remove(ms.Path(root)); rmErr != nil {
		return renamed, rmErr
	}
	return renamed, err
}

// CheckFree asks the milestone directory whether a name can be given out: a
// milestone on the board under it is an *ExistsError, and nothing else holds a
// name. A filed one does not — its ID tells it from a new one — and neither
// does a milestone on a ref nobody has fetched here yet: the names of two
// people's milestones are told apart by ID, so a double is something to see
// and rename, not something to prevent across machines. Every door that hands
// out a name asks here: Rename, and the CLI's and the layer's create.
func CheckFree(root, name string) error {
	all, err := LoadAll(root)
	if err != nil {
		return err
	}
	if i := slices.IndexFunc(all, func(m *Milestone) bool { return m.Name == name }); i >= 0 {
		return &ExistsError{Name: name, File: all[i].FileName()}
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
	unlock, err := lockAndMigrate(root, lock, rec)
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
	// The ID was rolled before the lock; a milestone on this disk, filed ones
	// included, or a ref in this tree may have it, and the name tells nobody.
	taken, err := scan(root)
	if err != nil {
		return nil, err
	}
	for slices.ContainsFunc(taken, func(m *Milestone) bool { return m.ID == ms.ID }) || rec.HasMilestoneRef(ms.ID) {
		ms.ID = NewID()
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
	if err := rec.RecordMilestone(ms.ID, ms.FileName(), ms.Bytes()); err != nil {
		return &RecordError{Name: ms.Name, Err: err}
	}
	return nil
}
