package milestone

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateSavesRecordsAndHoldsTheLockOnlyWhileBuilding(t *testing.T) {
	root := t.TempDir()
	lk, rec := &fakeLock{}, &fakeRec{dir: Dir(root)}
	ms, err := Create(lk, rec, func(existing []*Milestone) (*Milestone, error) {
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
	if ms.FileName() != "a-"+ms.ID+".md" || ms.Ref() != RefPrefix+ms.ID {
		t.Errorf("not recorded as a-%s.md on its ref: %s %s", ms.ID, ms.FileName(), ms.Ref())
	}
	if len(rec.names) != 1 || string(rec.last) != string(ms.Bytes()) {
		t.Errorf("recorded %v", rec.names)
	}
}

func TestCreateRefusalWritesNothing(t *testing.T) {
	root := t.TempDir()
	rec := &fakeRec{dir: Dir(root)}
	boom := errors.New("no")
	_, err := Create(&fakeLock{}, rec, func([]*Milestone) (*Milestone, error) { return nil, boom })
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
	rec := &fakeRec{dir: Dir(root)}

	ms, err := Edit(&fakeLock{}, rec, "a", func(m *Milestone) error { m.Add(idA); return nil })
	if err != nil || len(rec.names) != 0 {
		t.Fatalf("no-op edit: err=%v recorded=%v", err, rec.names)
	}
	if len(ms.Members()) != 1 {
		t.Error("milestone not returned")
	}

	ms, err = Edit(&fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("g"); return nil })
	if err != nil || len(rec.names) != 1 {
		t.Fatalf("edit: err=%v recorded=%v", err, rec.names)
	}
	got, _ := Load(refsAt(root), "a")
	if got.Goal != "g" || ms.Goal != "g" {
		t.Error("edit not saved")
	}
}

func TestEditPassesAMissingFileAndAFnRefusalThrough(t *testing.T) {
	root := t.TempDir()
	if _, err := Edit(&fakeLock{}, &fakeRec{dir: Dir(root)}, "zz", func(*Milestone) error { return nil }); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}
	write(t, root, "a", "---\ncolor: 45\n---\n")
	boom := errors.New("refused")
	rec := &fakeRec{dir: Dir(root)}
	_, err := Edit(&fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("x"); return boom })
	if !errors.Is(err, boom) || len(rec.names) != 0 {
		t.Errorf("err=%v recorded=%v", err, rec.names)
	}
	b, _ := os.ReadFile(testPath(root, "a"))
	if strings.Contains(string(b), "goal") {
		t.Error("a refused edit was saved")
	}
}

// The ref is a milestone's only copy, so a ref that cannot be written is a
// write that did not happen, and the caller is told so.
func TestARefThatCannotBeWrittenFailsTheWrite(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\n---\n")
	boom := errors.New("queue full")
	rec := &fakeRec{dir: Dir(root), err: boom}
	if _, err := Edit(&fakeLock{}, rec, "a", func(m *Milestone) error { m.SetGoal("g"); return nil }); !errors.Is(err, boom) {
		t.Fatalf("err=%v", err)
	}
	if got, _ := Load(refsAt(root), "a"); got.Goal == "g" {
		t.Error("a write that failed is visible")
	}
}

// Writing a filed milestone would put the change on its ref for every clone.
// Edit refuses it before fn runs, so no caller has to remember to; EditFiled,
// which filing and restoring use, may.
func TestEditRefusesAFiledMilestone(t *testing.T) {
	root := t.TempDir()
	write(t, root, "a", "---\ncolor: 45\nstatus: filed\n---\n")
	rec := &fakeRec{dir: Dir(root)}
	_, err := Edit(&fakeLock{}, rec, "a", func(*Milestone) error {
		t.Error("fn ran on a filed milestone")
		return nil
	})
	var fe *FiledError
	if !errors.As(err, &fe) || fe.Name != "a" || len(rec.names) != 0 {
		t.Errorf("err=%v recorded=%v", err, rec.names)
	}
	ms, err := EditFiled(&fakeLock{}, rec, "a", func(m *Milestone) error { m.SetStatus(""); return nil })
	if err != nil || ms.Filed() || len(rec.names) != 1 {
		t.Errorf("EditFiled: err=%v filed=%v recorded=%v", err, ms.Filed(), rec.names)
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
	lk, rec := &fakeLock{}, &fakeRec{dir: Dir(root)}
	if _, err := Edit(lk, rec, "b", func(m *Milestone) error { m.Unplace(); return nil }); err != nil {
		t.Fatal(err)
	}
	if lk.taken != 1 {
		t.Errorf("the removal and the re-attach must share one lock, taken %d times", lk.taken)
	}
	if b, _ := Load(refsAt(root), "b"); b.Placed || !b.Has(idA) {
		t.Errorf("b must leave the layer and keep its ticket, placed %v", b.Placed)
	}
	if c, _ := Load(refsAt(root), "c"); strings.Join(c.After, ",") != LegacyID("a")+",x" {
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
	lk, rec := &fakeLock{}, &fakeRec{dir: Dir(root)}
	ms, err := Rename(lk, rec, "a", "z")
	if err != nil || ms.Name != "z" || ms.ID != "AAAAAA" {
		t.Fatalf("rename: %v %+v", err, ms)
	}
	if lk.taken != 1 {
		t.Errorf("one lock, taken %d times", lk.taken)
	}
	if got := dirNames(t, root); got != "b-BBBBBB.md,z-AAAAAA.md" {
		t.Errorf("directory %s", got)
	}
	z, err := Load(refsAt(root), "AAAAAA")
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
	rec := &fakeRec{dir: Dir(root)}

	var ee *ExistsError
	if _, err := Rename(&fakeLock{}, rec, "a", "b"); !errors.As(err, &ee) || ee.Name != "b" || ee.ID != "BBBBBB" {
		t.Errorf("onto a milestone on the board: %v", err)
	}
	var fe *FiledError
	if _, err := Rename(&fakeLock{}, rec, "f", "n"); !errors.As(err, &fe) || fe.Name != "f" || fe.ID != "FFFFFF" {
		t.Errorf("a filed milestone is not renamed: %v", err)
	}
	if _, err := Rename(&fakeLock{}, rec, "a", "a"); err != nil || len(rec.names) != 0 {
		t.Fatalf("a rename to the same name writes nothing: err=%v recorded=%v", err, rec.names)
	}
	if _, err := Rename(&fakeLock{}, rec, "nothing", "n"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("an unknown milestone: %v", err)
	}
	if _, err := Rename(&fakeLock{}, rec, "a", "f"); err != nil {
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
	if _, err := Rename(&fakeLock{}, &fakeRec{dir: Dir(root)}, "dup", "z"); err == nil {
		t.Fatal("an ambiguous name was renamed")
	}
	if _, err := Rename(&fakeLock{}, &fakeRec{dir: Dir(root)}, "BBBBBB", "z"); err != nil {
		t.Fatal(err)
	}
	if got := dirNames(t, root); got != "dup-AAAAAA.md,z-BBBBBB.md" {
		t.Errorf("directory %s", got)
	}
}

// A new milestone is rolled a fresh ID when the one it drew is taken, filed
// milestones and refs counting.
func TestCreateAvoidsTakenIDs(t *testing.T) {
	root := t.TempDir()
	taken := New("t", 45, time.Now())
	rec := &fakeRec{dir: Dir(root), have: map[string]bool{taken.ID: true}}
	got, err := Create(&fakeLock{}, rec, func([]*Milestone) (*Milestone, error) { return taken, nil })
	if err != nil || got.ID == "" || rec.have[got.ID] {
		t.Errorf("a ref's ID was reused: %v %v", got, err)
	}
}

// A name a hand wrote in after: becomes the ID it names the next time Edit
// writes the file anyway; a token that names nothing stays as written, and a
// write that has nothing to change leaves the hand's spelling alone.
func TestEditSpellsHandWrittenNamesAsIDs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-K3F9QZ.md", "---\ncolor: 45\nafter: []\n---\n")
	write(t, root, "b", "---\ncolor: 46\nafter: [a,ghost]\n---\n")
	lk, rec := &fakeLock{}, &fakeRec{dir: Dir(root)}
	if _, err := Edit(lk, rec, "b", func(m *Milestone) error { m.Add(idA); return nil }); err != nil {
		t.Fatal(err)
	}
	if b, _ := Load(refsAt(root), "b"); strings.Join(b.After, ",") != "K3F9QZ,ghost" {
		t.Errorf("after = %v, want [K3F9QZ ghost]", b.After)
	}
	before, _ := os.ReadFile(testPath(root, "b"))
	if _, err := Edit(lk, rec, "b", func(m *Milestone) error { m.Add(idB); return nil }); err != nil {
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
	got, err := ResolveRefs(all, "ZZZZZZ", []string{"a", "K3F9QZ", "k3f9qz", "dup-bbbbbb", "old"})
	if err != nil || strings.Join(got, ",") != "K3F9QZ,BBBBBB,"+LegacyID("old") {
		t.Fatalf("got %v, %v", got, err)
	}
	var pe *PredecessorError
	if _, err := ResolveRefs(all, "ZZZZZZ", []string{"ghost"}); !errors.As(err, &pe) || pe.Pred != "ghost" {
		t.Errorf("unknown: %v", err)
	}
	var amb *AmbiguousError
	if _, err := ResolveRefs(all, "ZZZZZZ", []string{"dup"}); !errors.As(err, &amb) || len(amb.Choices) != 2 {
		t.Errorf("ambiguous: %v", err)
	}
	if got, err := ResolveRefs(all, "ZZZZZZ", nil); err != nil || got == nil || len(got) != 0 {
		t.Errorf("empty: %#v, %v", got, err)
	}
}
