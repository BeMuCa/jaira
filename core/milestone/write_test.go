package milestone

import (
	"errors"
	"os"
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
	names []string
	last  []byte
	err   error
}

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
	if _, err := os.Stat(testPath(root, "a")); err != nil {
		t.Error("not saved")
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
	if c, _ := Load(root, "c"); strings.Join(c.After, ",") != "a,x" {
		t.Errorf("c must now come after a, got %v", c.After)
	}
	if strings.Join(rec.names, ",") != "b,c" {
		t.Errorf("recorded %v, want b and c only", rec.names)
	}
}

// A rename is the file under its new name, every successor naming the new
// name, and a mark left under the old one for the other clones — in one
// locked write. This disk keeps no file under the old name.
func TestRenameMovesTheFileRewritesSuccessorsAndMarksTheOldName(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\ngoal: g\nafter: []\n---\n\n# a\n\n- "+idA+"\n")
	write(t, root, "b", "---\ncolor: 46\nafter: [x, a]\n---\n")
	write(t, root, "c", "---\ncolor: 47\n---\n")
	lk, rec := &fakeLock{}, &fakeRec{}
	ms, err := Rename(root, lk, rec, "a", "z", func(string) error { return nil })
	if err != nil || ms.Name != "z" {
		t.Fatalf("rename: %v", err)
	}
	if lk.taken != 1 {
		t.Errorf("the rename and the successors must share one lock, taken %d times", lk.taken)
	}
	if _, err := os.Stat(testPath(root, "a")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the old file is still on disk: %v", err)
	}
	z, err := Load(root, "z")
	if err != nil {
		t.Fatal(err)
	}
	if z.Goal != "g" || !z.Has(idA) || !z.Placed || z.Filed() {
		t.Errorf("z must carry a's goal, ticket and place: %s", z.Bytes())
	}
	if body := string(z.Bytes()); !strings.Contains(body, "\n# z\n") || strings.Contains(body, "# a") {
		t.Errorf("the heading keeps the old name:\n%s", body)
	}
	if b, _ := Load(root, "b"); strings.Join(b.After, ",") != "x,z" {
		t.Errorf("b must name z where it named a, got %v", b.After)
	}
	if strings.Join(rec.names, ",") != "z,b,a" {
		t.Errorf("recorded %v, want z, b and the mark on a", rec.names)
	}
	mark := FromBytes("a", rec.last)
	if !mark.Filed() || mark.RenamedTo != "z" || mark.HoldsName() {
		t.Errorf("the old ref must carry a rename mark that frees the name:\n%s", rec.last)
	}
}

// The new name has to be free — on the board, filed here, or as the caller's
// taken says — while the mark an earlier rename left is free, so a rename can
// be taken back.
func TestRenameRefusesATakenNameButNotTheMarkOfARename(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\n---\n")
	write(t, root, "b", "---\ncolor: 46\n---\n")
	write(t, root, "f", "---\ncolor: 47\nstatus: filed\n---\n")
	write(t, root, "old", "---\ncolor: 48\nstatus: filed\nrenamed-to: a\n---\n")
	rec := &fakeRec{}
	free := func(string) error { return nil }

	var ee *ExistsError
	if _, err := Rename(root, &fakeLock{}, rec, "a", "b", free); !errors.As(err, &ee) || ee.Name != "b" {
		t.Errorf("onto a milestone on the board: %v", err)
	}
	var fe *FiledError
	if _, err := Rename(root, &fakeLock{}, rec, "a", "f", free); !errors.As(err, &fe) || fe.Name != "f" || fe.RenamedTo != "" {
		t.Errorf("onto a filed one: %v", err)
	}
	boom := errors.New("filed on a ref")
	if _, err := Rename(root, &fakeLock{}, rec, "a", "n", func(string) error { return boom }); !errors.Is(err, boom) {
		t.Errorf("taken must be asked: %v", err)
	}
	if _, err := Rename(root, &fakeLock{}, rec, "old", "n", free); !errors.As(err, &fe) || fe.RenamedTo != "a" {
		t.Errorf("a mark is not a milestone to rename: %v", err)
	}
	if _, err := Rename(root, &fakeLock{}, rec, "a", "a", free); err != nil || len(rec.names) != 0 {
		t.Fatalf("refusals and a rename to the same name write nothing: err=%v recorded=%v", err, rec.names)
	}

	if _, err := Rename(root, &fakeLock{}, rec, "a", "old", free); err != nil {
		t.Fatalf("renaming back onto a rename's mark: %v", err)
	}
	if old, err := Load(root, "old"); err != nil || old.Filed() || old.RenamedTo != "" {
		t.Errorf("old must be on the board again, unmarked: %v\n%s", err, old.Bytes())
	}
}
