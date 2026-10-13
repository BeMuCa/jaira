package milestone

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeLock struct{ held, taken int }

func (f *fakeLock) Lock(string) (func(), error) {
	f.held++
	f.taken++
	return func() { f.held-- }, nil
}

type fakeRec struct {
	have  map[string]bool // IDs that already have a ref
	names []string
	last  []byte
	err   error
}

func (r *fakeRec) HasMilestoneRef(id string) bool { return r.have[id] }

func (r *fakeRec) RecordMilestone(name, file string, content []byte) error {
	r.names = append(r.names, name)
	r.last = content
	return r.err
}

func TestCreateSavesRecordsAndHoldsTheLockOnlyWhileBuilding(t *testing.T) {
	root := t.TempDir()
	lk, rec := &fakeLock{}, &fakeRec{}
	ms, err := Create(root, lk, rec, func(existing []*Milestone) (*Milestone, error) {
		if lk.held != 1 {
			t.Error("build must run under the lock")
		}
		return New("a", 45, time.Now()), nil
	})
	if err != nil || ms.Name != "a" {
		t.Fatalf("create: %v", err)
	}
	if lk.held != 0 || lk.taken != 1 {
		t.Errorf("lock held=%d taken=%d", lk.held, lk.taken)
	}
	if _, err := os.Stat(ms.Path(root)); err != nil || filepath.Base(ms.Path(root)) != "a-"+ms.ID+".md" {
		t.Errorf("not saved as a-%s.md: %v", ms.ID, err)
	}
	if len(rec.names) != 1 || string(rec.last) != string(ms.Bytes()) {
		t.Errorf("recorded %v", rec.names)
	}
}

func TestCreateRefusalWritesNothing(t *testing.T) {
	root := t.TempDir()
	rec := &fakeRec{}
	boom := errors.New("no")
	_, err := Create(root, &fakeLock{}, rec, func([]*Milestone) (*Milestone, error) { return nil, boom })
	if !errors.Is(err, boom) || len(rec.names) != 0 {
		t.Errorf("err=%v recorded=%v", err, rec.names)
	}
	if _, err := os.Stat(testPath(root, "a")); err == nil {
		t.Error("file written despite refusal")
	}
}

func TestEditWritesOnlyWhenTheBytesChanged(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\n---\n- "+idA+"\n")
	rec := &fakeRec{}

	ms, err := Edit(root, &fakeLock{}, rec, "a", func(m *Milestone) error { m.Add(idA); return nil })
	if err != nil || len(rec.names) != 0 {
		t.Fatalf("no-op edit: err=%v recorded=%v", err, rec.names)
	}
	if len(ms.Members()) != 1 {
		t.Error("milestone not returned")
	}

	ms, err = Edit(root, &fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("g"); return nil })
	if err != nil || len(rec.names) != 1 {
		t.Fatalf("edit: err=%v recorded=%v", err, rec.names)
	}
	got, _ := Load(root, "a")
	if got.Goal != "g" || ms.Goal != "g" {
		t.Error("edit not saved")
	}
}

func TestEditPassesAMissingFileAndAFnRefusalThrough(t *testing.T) {
	root := t.TempDir()
	if _, err := Edit(root, &fakeLock{}, &fakeRec{}, "zz", func(*Milestone) error { return nil }); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}
	write(t, root, "a", "---\ncolor: 45\n---\n")
	boom := errors.New("refused")
	rec := &fakeRec{}
	_, err := Edit(root, &fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("x"); return boom })
	if !errors.Is(err, boom) || len(rec.names) != 0 {
		t.Errorf("err=%v recorded=%v", err, rec.names)
	}
	b, _ := os.ReadFile(testPath(root, "a"))
	if strings.Contains(string(b), "goal") {
		t.Error("a refused edit was saved")
	}
}

// A recorder that fails must not undo a save that happened: the caller gets
// the milestone and a RecordError it may downgrade to a warning.
func TestRecordFailureIsReportedBesideTheSavedMilestone(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\n---\n")
	rec := &fakeRec{err: errors.New("queue full")}
	ms, err := Edit(root, &fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("g"); return nil })
	var re *RecordError
	if !errors.As(err, &re) || ms == nil {
		t.Fatalf("err=%v ms=%v", err, ms)
	}
	if got, _ := Load(root, "a"); got.Goal != "g" {
		t.Error("save was lost")
	}
}

// A filed milestone's file can still lie here — a fetch of somebody else's
// filing leaves it — and writing it would put the change on its ref for every
// clone. Edit refuses it before fn runs, so no caller has to remember to.
func TestEditRefusesAFiledMilestone(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\nstatus: filed\n---\n")
	rec := &fakeRec{}
	_, err := Edit(root, &fakeLock{}, rec, "a", func(*Milestone) error {
		t.Error("fn ran on a filed milestone")
		return nil
	})
	var fe *FiledError
	if !errors.As(err, &fe) || fe.Name != "a" || len(rec.names) != 0 {
		t.Errorf("err=%v recorded=%v", err, rec.names)
	}
}

// Taking a milestone out of the layer closes the chain behind it in the same
// locked write, whoever calls Edit: what came after it now comes after what
// came before it, and a milestone that does not name it is not written.
func TestEditTakingAMilestoneOutClosesTheChain(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\nafter: []\n---\n")
	write(t, root, "b", "---\ncolor: 46\nafter: [a]\n---\n- "+idA+"\n")
	write(t, root, "c", "---\ncolor: 47\nafter: [b, x]\n---\n")
	write(t, root, "d", "---\ncolor: 48\nafter: [a]\n---\n")
	lk, rec := &fakeLock{}, &fakeRec{}
	if _, err := Edit(root, lk, rec, "b", func(m *Milestone) error { m.Unplace(); return nil }); err != nil {
		t.Fatal(err)
	}
	if lk.taken != 1 {
		t.Errorf("the removal and the re-attach must share one lock, taken %d times", lk.taken)
	}
	if b, _ := Load(root, "b"); b.Placed || !b.Has(idA) {
		t.Errorf("b must leave the layer and keep its ticket, placed %v", b.Placed)
	}
	if c, _ := Load(root, "c"); strings.Join(c.After, ",") != LegacyID("a")+",x" {
		t.Errorf("c must now come after a (by ID), got %v", c.After)
	}
	if strings.Join(rec.names, ",") != LegacyID("b")+","+LegacyID("c") {
		t.Errorf("recorded %v, want b and c only", rec.names)
	}
}

// A rename is the file under its new name on the same ID: one write on the
// one ref, no successor touched, nothing left behind under the old name.
func TestRenameChangesOnlyTheFileNameOnTheSameRef(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-AAAAAA.md", "---\ncolor: 45\ngoal: g\nafter: []\n---\n\n# a\n\n- "+idA+"\n")
	succ := "---\ncolor: 46\nafter: [x, AAAAAA]\n---\n"
	writeFile(t, root, "b-BBBBBB.md", succ)
	lk, rec := &fakeLock{}, &fakeRec{}
	ms, err := Rename(root, lk, rec, "a", "z")
	if err != nil || ms.Name != "z" || ms.ID != "AAAAAA" {
		t.Fatalf("rename: %v %+v", err, ms)
	}
	if lk.taken != 1 {
		t.Errorf("one lock, taken %d times", lk.taken)
	}
	if got := dirNames(t, root); got != "b-BBBBBB.md,z-AAAAAA.md" {
		t.Errorf("directory %s", got)
	}
	z, err := Load(root, "AAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	if z.Goal != "g" || !z.Has(idA) || !z.Placed || z.Filed() {
		t.Errorf("z must carry a's goal, ticket and place: %s", z.Bytes())
	}
	if body := string(z.Bytes()); !strings.Contains(body, "\n# z\n") || strings.Contains(body, "# a") {
		t.Errorf("the heading keeps the old name:\n%s", body)
	}
	if b, _ := os.ReadFile(filepath.Join(Dir(root), "b-BBBBBB.md")); string(b) != succ {
		t.Errorf("a successor was rewritten:\n%s", b)
	}
	if strings.Join(rec.names, ",") != "AAAAAA" {
		t.Errorf("recorded %v, want one write on the ID ref", rec.names)
	}
}

// The new name has to be free on this board; a filed milestone is not
// renamed, a rename to the name it already has writes nothing, and a filed
// milestone does not hold its name.
func TestRenameRefusals(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-AAAAAA.md", "---\ncolor: 45\n---\n")
	writeFile(t, root, "b-BBBBBB.md", "---\ncolor: 46\n---\n")
	writeFile(t, root, "f-FFFFFF.md", "---\ncolor: 47\nstatus: filed\n---\n")
	rec := &fakeRec{}

	var ee *ExistsError
	if _, err := Rename(root, &fakeLock{}, rec, "a", "b"); !errors.As(err, &ee) || ee.Name != "b" {
		t.Errorf("onto a milestone on the board: %v", err)
	}
	var fe *FiledError
	if _, err := Rename(root, &fakeLock{}, rec, "f", "n"); !errors.As(err, &fe) || fe.Name != "f" || fe.File != "f-FFFFFF.md" {
		t.Errorf("a filed milestone is not renamed: %v", err)
	}
	if _, err := Rename(root, &fakeLock{}, rec, "a", "a"); err != nil || len(rec.names) != 0 {
		t.Fatalf("a rename to the same name writes nothing: err=%v recorded=%v", err, rec.names)
	}
	if _, err := Rename(root, &fakeLock{}, rec, "nothing", "n"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("an unknown milestone: %v", err)
	}
	if _, err := Rename(root, &fakeLock{}, rec, "a", "f"); err != nil {
		t.Errorf("a filed milestone held its name: %v", err)
	}
	if got := dirNames(t, root); got != "b-BBBBBB.md,f-AAAAAA.md,f-FFFFFF.md" {
		t.Errorf("directory %s", got)
	}
}

// Two milestones with one name are told apart by their ID, and a rename of one
// leaves the other alone.
func TestRenameOfOneOfTwoSameNamedMilestones(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "dup-AAAAAA.md", "---\ncolor: 45\n---\n")
	writeFile(t, root, "dup-BBBBBB.md", "---\ncolor: 46\n---\n")
	if _, err := Rename(root, &fakeLock{}, &fakeRec{}, "dup", "z"); err == nil {
		t.Fatal("an ambiguous name was renamed")
	}
	if _, err := Rename(root, &fakeLock{}, &fakeRec{}, "BBBBBB", "z"); err != nil {
		t.Fatal(err)
	}
	if got := dirNames(t, root); got != "dup-AAAAAA.md,z-BBBBBB.md" {
		t.Errorf("directory %s", got)
	}
}

// A write on a board that still has files from before IDs moves them all
// first, under the lock; a new milestone is rolled a fresh ID when the one it
// drew is taken, filed milestones and refs counting.
func TestWritesMigrateFirstAndCreateAvoidsTakenIDs(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "q4", "---\ncolor: 45\n---\n")
	rec := &fakeRec{}
	ms, err := Create(root, &fakeLock{}, rec, func([]*Milestone) (*Milestone, error) {
		return New("q5", 46, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := dirNames(t, root); got != "q4-"+LegacyID("q4")+".md,q5-"+ms.ID+".md" {
		t.Errorf("directory %s", got)
	}
	if strings.Join(rec.names, ",") != LegacyID("q4")+","+ms.ID {
		t.Errorf("recorded %v", rec.names)
	}
	taken := New("t", 45, time.Now())
	rec = &fakeRec{have: map[string]bool{taken.ID: true}}
	got, err := Create(root, &fakeLock{}, rec, func([]*Milestone) (*Milestone, error) { return taken, nil })
	if err != nil || got.ID == "" || got.ID == rec.names[0] && rec.have[got.ID] {
		t.Errorf("a ref's ID was reused: %v %v", got, err)
	}
	if rec.have[got.ID] {
		t.Errorf("ID %s has a ref already", got.ID)
	}
}

// A name a hand wrote in after: becomes the ID it names the next time Edit
// writes the file anyway; a token that names nothing stays as written, and a
// write that has nothing to change leaves the hand's spelling alone.
func TestEditSpellsHandWrittenNamesAsIDs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-K3F9QZ.md", "---\ncolor: 45\nafter: []\n---\n")
	write(t, root, "b", "---\ncolor: 46\nafter: [a,ghost]\n---\n")
	lk, rec := &fakeLock{}, &fakeRec{}
	if _, err := Edit(root, lk, rec, "b", func(m *Milestone) error { m.Add(idA); return nil }); err != nil {
		t.Fatal(err)
	}
	if b, _ := Load(root, "b"); strings.Join(b.After, ",") != "K3F9QZ,ghost" {
		t.Errorf("after = %v, want [K3F9QZ ghost]", b.After)
	}
	before, _ := os.ReadFile(testPath(root, "b"))
	if _, err := Edit(root, lk, rec, "b", func(m *Milestone) error { m.Add(idB); return nil }); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(testPath(root, "b")); strings.Count(string(got), "after:") != 1 || !strings.Contains(string(got), "after: [K3F9QZ, ghost]") || len(got) <= len(before) {
		t.Errorf("second write:\n%s", got)
	}
}

func TestResolveRefs(t *testing.T) {
	all := []*Milestone{
		FromBytes("a-K3F9QZ.md", nil), FromBytes("dup-AAAAAA.md", nil), FromBytes("dup-BBBBBB.md", nil), FromBytes("old.md", nil),
	}
	got, err := ResolveRefs(all, []string{"a", "K3F9QZ", "k3f9qz", "dup-bbbbbb", "old"})
	if err != nil || strings.Join(got, ",") != "K3F9QZ,BBBBBB,"+LegacyID("old") {
		t.Fatalf("got %v, %v", got, err)
	}
	var pe *PredecessorError
	if _, err := ResolveRefs(all, []string{"ghost"}); !errors.As(err, &pe) || pe.Pred != "ghost" {
		t.Errorf("unknown: %v", err)
	}
	var amb *AmbiguousError
	if _, err := ResolveRefs(all, []string{"dup"}); !errors.As(err, &amb) || len(amb.Choices) != 2 {
		t.Errorf("ambiguous: %v", err)
	}
	if got, err := ResolveRefs(all, nil); err != nil || got == nil || len(got) != 0 {
		t.Errorf("empty: %#v, %v", got, err)
	}
}
