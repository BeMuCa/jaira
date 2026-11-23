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

// Rewording is the whole point of --text: the criterion changed, its state and
// its evidence did not. Ticking the stale item with the proof "obsolete,
// replaced by item 6" was the only route before, and it leaves a [x] beside work
// nobody ever did.
func TestDodTextRewordsWithoutTouchingTheState(t *testing.T) {
	dir, id := dodTestStore(t)
	if out, err := runDoD(t, dir, id, "1", "--done", "--proof", "internal/x.go:12; TestX"); err != nil {
		t.Fatalf("setup: %v\n%s", err, out)
	}
	out, err := runDoD(t, dir, id, "1", "--text", "one, said properly")
	if err != nil {
		t.Fatalf("dod --text: %v\n%s", err, out)
	}
	if !strings.Contains(out, "[x] one, said properly") {
		t.Errorf("reworded item lost its tick:\n%s", out)
	}
	if !strings.Contains(out, "proof: internal/x.go:12; TestX") {
		t.Errorf("reworded item lost its proof:\n%s", out)
	}
	if !strings.Contains(out, "[ ] two") {
		t.Errorf("the other item was disturbed:\n%s", out)
	}
}

// --superseded retires an item: not achieved, and no longer in the way.
func TestDodSupersededRetiresAnItem(t *testing.T) {
	dir, id := dodTestStore(t)
	out, err := runDoD(t, dir, id, "2", "--superseded")
	if err != nil {
		t.Fatalf("dod --superseded: %v\n%s", err, out)
	}
	if !strings.Contains(out, "[-] two") {
		t.Errorf("item was not marked superseded:\n%s", out)
	}

	// Ticking what is left completes the ticket, which is what distinguishes a
	// superseded item from an open one.
	if out, err := runDoD(t, dir, id, "1", "--done"); err != nil {
		t.Fatalf("dod --done: %v\n%s", err, out)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	complete, remaining := tk.DoDComplete()
	if !complete {
		t.Errorf("a superseded item still blocks completion; remaining: %v", remaining)
	}
	for _, it := range tk.DoDItems {
		if it.Text == "two" && it.Checked() {
			t.Error("a superseded item reports as done")
		}
	}
}

// Two states are a contradiction, and the message must name every state so the
// caller does not have to go looking for the fourth one.
func TestDodRefusesTwoStates(t *testing.T) {
	dir, id := dodTestStore(t)
	out, err := runDoD(t, dir, id, "1", "--done", "--superseded")
	if err == nil {
		t.Fatalf("two state flags were accepted:\n%s", out)
	}
	if !strings.Contains(err.Error(), "--superseded") {
		t.Errorf("the refusal does not mention the new state: %v", err)
	}
}
