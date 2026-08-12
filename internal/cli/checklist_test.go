package cli

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// dodTestStore builds a fresh store with one ticket carrying a two-item
// Definition of Done, and returns its directory and handle.
func dodTestStore(t *testing.T) (dir, id string) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  "t",
		ticket.FieldStatus: "backlog",
	}, nil, "## Definition of Done\n\n- [ ] one\n- [ ] two\n")
	if err != nil {
		t.Fatal(err)
	}
	return dir, tk.ID
}

// runDoD runs the real 'dod' command against a store, the same path a user or
// agent invokes it through, and returns its stdout.
func runDoD(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "dod"}, args...))
	err := root.Execute()
	return out.String(), err
}

// A single call ticks the item and records the evidence for it.
func TestDodProofTicksAndRecordsInOneCall(t *testing.T) {
	dir, id := dodTestStore(t)
	out, err := runDoD(t, dir, id, "1", "--done", "--proof", "internal/x.go:12; TestX")
	if err != nil {
		t.Fatalf("dod --done --proof: %v\n%s", err, out)
	}
	if !strings.Contains(out, "[x] one") {
		t.Errorf("item was not marked done:\n%s", out)
	}
	if !strings.Contains(out, "proof: internal/x.go:12; TestX") {
		t.Errorf("proof was not recorded or shown:\n%s", out)
	}

	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(tk.DoDItems) != 2 || tk.DoDItems[0].Proof != "internal/x.go:12; TestX" {
		t.Errorf("proof did not survive on disk: %#v", tk.DoDItems)
	}
}

// --proof with no state flag records the proof and leaves the marker alone.
func TestDodProofAloneLeavesMarkerAlone(t *testing.T) {
	dir, id := dodTestStore(t)
	if _, err := runDoD(t, dir, id, "1", "--proof", "still open, but here is where it will land"); err != nil {
		t.Fatal(err)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if tk.DoDItems[0].State != ticket.StateTodo {
		t.Errorf("state changed to %v; --proof alone must not change the marker", tk.DoDItems[0].State)
	}
	if tk.DoDItems[0].Proof != "still open, but here is where it will land" {
		t.Errorf("proof = %q", tk.DoDItems[0].Proof)
	}
}

// Setting proof twice on the same item replaces the line rather than
// stacking a second one.
func TestDodProofTwiceReplaces(t *testing.T) {
	dir, id := dodTestStore(t)
	if _, err := runDoD(t, dir, id, "1", "--proof", "first"); err != nil {
		t.Fatal(err)
	}
	if _, err := runDoD(t, dir, id, "1", "--proof", "second"); err != nil {
		t.Fatal(err)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(tk.DoDItems) != 2 {
		t.Fatalf("got %d items, want 2 (a stray proof line must not be read as a new one)", len(tk.DoDItems))
	}
	if tk.DoDItems[0].Proof != "second" {
		t.Errorf("proof = %q, want %q", tk.DoDItems[0].Proof, "second")
	}
}

// --proof cannot be combined with --add or --option: those do not address a
// single existing item.
func TestDodProofRejectsAddAndOption(t *testing.T) {
	dir, id := dodTestStore(t)
	if _, err := runDoD(t, dir, id, "--add", "three", "--proof", "x"); err == nil {
		t.Error("expected a usage error combining --proof with --add")
	}
	if _, err := runDoD(t, dir, id, "1", "--option", "planning", "--proof", "x"); err == nil {
		t.Error("expected a usage error combining --proof with --option")
	}
}
