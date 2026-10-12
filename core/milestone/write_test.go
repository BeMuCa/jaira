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
	ms, err := Create(root, lk, rec, "a", func(existing []*Milestone) (*Milestone, error) {
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
	_, err := Create(root, &fakeLock{}, rec, "a", func([]*Milestone) (*Milestone, error) { return nil, boom })
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
