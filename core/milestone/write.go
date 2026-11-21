package milestone

import (
	"bytes"
	"fmt"
	"slices"
)

// LockName keys the store lock every milestone write is made under. One name
// for all of them rather than one per milestone: creating a milestone reads
// every existing one to avoid their colours, so two creations at once have to
// be ordered even though they write different refs.
const LockName = "milestones"

// Locker is the one thing of the ticket store a write needs. An interface so
// this package does not import the store's locking and a test can count.
type Locker interface {
	Lock(name string) (unlock func(), err error)
}

// FiledError is Edit's refusal of a milestone that is marked filed. Writing
// it would put the change on its ref for every clone that fetches.
type FiledError struct {
	Name string
	ID   string
}

func (e *FiledError) Error() string {
	return fmt.Sprintf("milestone %q has been filed; 'jaira restore %s-%s' brings it back", e.Name, e.Name, e.ID)
}

// ExistsError is CheckFree's answer for a name a milestone on the board
// already has.
type ExistsError struct {
	Name string
	ID   string // the ID of the one that has it; of two, the first by ID
}

func (e *ExistsError) Error() string { return fmt.Sprintf("milestone %q already exists", e.Name) }

// Edit is the one write path for a milestone that exists: lock, read, let fn
// change it, record it on its ref. The CLI and the board both write through it
// so they cannot drift — what a person does in the layer is exactly what an
// agent does with the command.
//
// fn may refuse by returning an error; nothing is written then. The milestone
// is recorded only when fn changed its bytes, because a re-add of a ticket
// already there is not a write and must not queue a ref. A missing milestone
// comes back as the Load error (errors.Is os.ErrNotExist) so the caller can
// say why in its own words, and a filed one as a *FiledError before fn runs.
//
// When fn takes a placed milestone out of the layer, Edit also closes the
// chain behind it under the same lock: every milestone that came after it now
// comes after what came before it (ReattachSuccessors). Here and not in the
// callers, so the layer's x and 'jaira milestone set --no-roadmap' cannot leave
// two different pictures behind.
func Edit(lock Locker, refs Refs, name string, fn func(*Milestone) error) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()
	return editHeld(refs, name, false, fn)
}

// EditFiled is Edit for the one writer that may change a filed milestone:
// filing and restoring it, which set and take off the very line that makes it
// filed.
func EditFiled(lock Locker, refs Refs, name string, fn func(*Milestone) error) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()
	return editHeld(refs, name, true, fn)
}

func editHeld(refs Refs, name string, filedToo bool, fn func(*Milestone) error) (*Milestone, error) {
	ms, err := Load(refs, name)
	if err != nil {
		return nil, err
	}
	if ms.Filed() && !filedToo {
		return ms, &FiledError{Name: ms.Name, ID: ms.ID}
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
		// milestone is being written anyway.
		all, err := LoadAll(refs)
		if err != nil {
			return ms, err
		}
		if ids := SpellAsIDs(all, ms.After); !slices.Equal(ids, ms.After) {
			ms.SetAfter(ids)
		}
	}
	var succ []*Milestone
	if wasPlaced && !ms.Placed {
		// Read before the write, while this milestone's own after line is
		// still on its ref: it is what its successors move up to.
		all, err := LoadAll(refs)
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
	return ms, recordAll(refs, append([]*Milestone{ms}, succ...))
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
// line: a rename is one write on the one ref, carrying the document under its
// new file name. Nobody who follows the milestone is touched.
//
// to has to be a name no milestone on the board carries (CheckFree). The
// milestone refuses the way Edit does when it is filed, and renaming it to
// its own name writes nothing.
func Rename(lock Locker, refs Refs, from, to string) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()

	ms, err := Load(refs, from)
	if err != nil {
		return nil, err
	}
	if ms.Filed() {
		return ms, &FiledError{Name: ms.Name, ID: ms.ID}
	}
	if to == ms.Name {
		return ms, nil
	}
	if err := CheckFree(refs, to); err != nil {
		return ms, err
	}
	renamed := FromBytes(to+"-"+ms.ID+".md", ms.Bytes())
	renamed.retitle(ms.Name, to)
	return renamed, record(refs, renamed)
}

// CheckFree asks the refs whether a name can be given out: a milestone on the
// board under it is an *ExistsError, and nothing else holds a name. A filed one
// does not — its ID tells it from a new one — and neither does a milestone on
// a ref nobody has fetched here yet: the names of two people's milestones are
// told apart by ID, so a double is something to see and rename, not something
// to prevent across machines. Every door that hands out a name asks here:
// Rename, and the CLI's and the layer's create.
func CheckFree(refs Refs, name string) error {
	all, err := LoadAll(refs)
	if err != nil {
		return err
	}
	if i := slices.IndexFunc(all, func(m *Milestone) bool { return m.Name == name }); i >= 0 {
		return &ExistsError{Name: name, ID: all[i].ID}
	}
	return nil
}

// retitle rewrites the "# <name>" heading New puts under the frontmatter, so
// a renamed milestone does not open with its old name. Only that exact line:
// any other is somebody's prose.
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
func Create(lock Locker, refs Refs, build func(existing []*Milestone) (*Milestone, error)) (*Milestone, error) {
	unlock, err := lock.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()

	existing, err := LoadAll(refs)
	if err != nil {
		return nil, err
	}
	ms, err := build(existing)
	if err != nil {
		return nil, err
	}
	// The ID was rolled before the lock; a milestone here, filed ones
	// included, may have it, and the name tells nobody.
	taken, err := scan(refs)
	if err != nil {
		return nil, err
	}
	for slices.ContainsFunc(taken, func(m *Milestone) bool { return m.ID == ms.ID }) || refs.HasMilestoneRef(ms.ID) {
		ms.ID = NewID()
	}
	return ms, record(refs, ms)
}

// recordAll records every milestone, stopping at the first that fails.
func recordAll(refs Refs, ms []*Milestone) error {
	for _, m := range ms {
		if err := record(refs, m); err != nil {
			return err
		}
	}
	return nil
}

func record(refs Refs, ms *Milestone) error {
	if err := refs.RecordMilestone(ms.ID, ms.FileName(), ms.Bytes()); err != nil {
		return fmt.Errorf("milestone %q: %w", ms.Name, err)
	}
	return nil
}
