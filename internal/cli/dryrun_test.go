package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// moveTestStore builds a store with one under-specified backlog ticket, which is
// what a gate refuses.
func moveTestStore(t *testing.T) (dir, id string) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))
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
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return dir, tk.ID
}

// readTicketFile returns the raw bytes of the one ticket in the store, so a test
// can prove nothing was written rather than trusting the command's word.
func readTicketFile(t *testing.T, dir string) []byte {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, ".jaira", "tickets", "*.md"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("expected one ticket file, got %v (%v)", matches, err)
	}
	raw, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// The point of the flag: ask the gate and change nothing. Without it the only
// way to find out was to try the move — which is how one of the user's tickets
// was moved by accident.
func TestMoveDryRunWritesNothingWhenRefused(t *testing.T) {
	dir, id := moveTestStore(t)
	before := readTicketFile(t, dir)

	out, err := runCLI(t, dir, "move", id, "--to", "in-progress", "--dry-run")
	if err == nil {
		t.Fatalf("an under-specified ticket was not refused:\n%s", out)
	}
	if !strings.Contains(err.Error(), "nothing was written") {
		t.Errorf("the refusal does not say the run was dry: %v", err)
	}
	if got := readTicketFile(t, dir); string(got) != string(before) {
		t.Errorf("the ticket file changed during a dry run:\n%s", got)
	}
}

// The staged fields are the trap: a real move writes them before the gate runs,
// so a dry run that reused that path would leave --what and an assignee behind
// on a move it then refused.
func TestMoveDryRunDoesNotWriteStagedFields(t *testing.T) {
	dir, id := moveTestStore(t)
	before := readTicketFile(t, dir)

	out, err := runCLI(t, dir, "move", id, "--to", "review", "--dry-run",
		"--what", "streamed the writer", "--why", "it buffered", "--resolves", "the export completes")
	if err == nil {
		t.Fatalf("expected a refusal:\n%s", out)
	}
	got := string(readTicketFile(t, dir))
	if strings.Contains(got, "streamed the writer") {
		t.Error("an outcome staged for a dry run was written to the ticket")
	}
	if got != string(before) {
		t.Errorf("the ticket file changed during a dry run:\n%s", got)
	}
}

// A move that would be allowed reports so, exits zero, and still writes nothing.
func TestMoveDryRunOnAnAllowedMove(t *testing.T) {
	dir, id := moveTestStore(t)
	if out, err := runCLI(t, dir, "set", id,
		"goal=make it work", "context=came up while debugging", "definition-of-done=a test covers it"); err != nil {
		t.Fatalf("setup: %v\n%s", err, out)
	}
	before := readTicketFile(t, dir)

	out, err := runCLI(t, dir, "move", id, "--to", "todo", "--dry-run")
	if err != nil {
		t.Fatalf("an allowed move was refused in a dry run: %v\n%s", err, out)
	}
	if !strings.Contains(out, "would be allowed") {
		t.Errorf("the dry run does not say the move would pass:\n%s", out)
	}
	if got := readTicketFile(t, dir); string(got) != string(before) {
		t.Error("an allowed dry run still wrote to the ticket")
	}

	// And the real move afterwards actually moves it, so the dry run did not
	// consume the move.
	if out, err := runCLI(t, dir, "move", id, "--to", "todo"); err != nil {
		t.Fatalf("the real move failed after a dry run: %v\n%s", err, out)
	}
	if got := string(readTicketFile(t, dir)); !strings.Contains(got, "status: todo") {
		t.Errorf("the ticket did not move:\n%s", got)
	}
}
