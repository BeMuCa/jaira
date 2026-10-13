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

func (r *fakeRec) RecordMilestone(name string, content []byte) error {
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
	if _, err := os.Stat(Path(root, "a")); err != nil {
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
	if _, err := os.Stat(Path(root, "a")); err == nil {
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
	b, _ := os.ReadFile(Path(root, "a"))
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
