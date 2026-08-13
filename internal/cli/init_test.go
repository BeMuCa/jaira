package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// runInit drives the real 'init' command tree against dir.
func runInit(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "init"}, args...))
	err := root.Execute()
	return out.String(), err
}

// runCreate drives the real 'create' command tree against dir.
func runCreate(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "create"}, args...))
	err := root.Execute()
	return out.String(), err
}

// initTestHome isolates a test's state directory and catalogue from the real
// user's, the same way other cli tests do.
func initTestHome(t *testing.T) {
	t.Helper()
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
}

// TestInitWithNoDefaultBoardWritesNoLaneFiles is the criterion's load-bearing
// half at the CLI boundary: a project whose owner changed nothing carries no
// .jaira/lanes/ directory at all.
func TestInitWithNoDefaultBoardWritesNoLaneFiles(t *testing.T) {
	initTestHome(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "does-not-exist.md"))
	dir := t.TempDir()

	out, err := runInit(t, dir)
	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if _, statErr := os.Stat(filepath.Join(dir, ".jaira", "lanes")); !os.IsNotExist(statErr) {
		t.Errorf("expected no .jaira/lanes directory, stat err = %v", statErr)
	}
	if !strings.Contains(out, "built-in") {
		t.Errorf("init output = %q, want it to say the built-in lanes are in use", out)
	}
}

// TestInitWithDefaultBoardDroppingALaneMaterialisesTheSelection asserts a
// default board that drops a lane creates one file per selected lane, and no
// others.
func TestInitWithDefaultBoardDroppingALaneMaterialisesTheSelection(t *testing.T) {
	initTestHome(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(board, []byte("---\nlanes: [backlog, todo, done]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", board)
	dir := t.TempDir()

	out, err := runInit(t, dir)
	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	projDir := lane.ProjectLanesDir(dir)
	for _, want := range []string{"backlog.md", "todo.md", "done.md"} {
		if _, statErr := os.Stat(filepath.Join(projDir, want)); statErr != nil {
			t.Errorf("expected %s to exist: %v", want, statErr)
		}
	}
	if _, statErr := os.Stat(filepath.Join(projDir, "review.md")); !os.IsNotExist(statErr) {
		t.Error("a lane not in the selection must not be materialised")
	}
	if !strings.Contains(out, "Wrote 3 lane file") {
		t.Errorf("init output = %q, want it to report the count written", out)
	}
}

// TestInitTwiceDoesNotReapplyDefaultBoardOverProjectChoices asserts init
// stays safe to run twice: once a project has its own .jaira/lanes/, a
// second init must not re-materialise (and potentially discard a hand edit).
func TestInitTwiceDoesNotReapplyDefaultBoardOverProjectChoices(t *testing.T) {
	initTestHome(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(board, []byte("---\nlanes: [backlog, todo, done]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", board)
	dir := t.TempDir()

	if _, err := runInit(t, dir); err != nil {
		t.Fatal(err)
	}
	projDir := lane.ProjectLanesDir(dir)
	// Simulate a hand edit after the first init.
	handEdited := []byte("---\nid: todo\nname: Hand-edited Todo\nafter: backlog\nprecedence: 20\n---\n")
	if err := os.WriteFile(filepath.Join(projDir, "todo.md"), handEdited, 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runInit(t, dir)
	if err != nil {
		t.Fatalf("second init: %v\n%s", err, out)
	}
	if !strings.Contains(out, "already scopes its own lanes") {
		t.Errorf("second init output = %q, want it to say the project already has its own lanes", out)
	}
	got, err := os.ReadFile(filepath.Join(projDir, "todo.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(handEdited) {
		t.Error("a second init must not overwrite a hand-edited project lane")
	}
}

// TestCreateAfterInitTicksDefaultBoardOptions asserts a ticket created after
// init in a root with a default board carries the board's pre-ticked
// options — the default board changes what a new ticket starts with, not
// only what jaira init writes to disk.
func TestCreateAfterInitTicksDefaultBoardOptions(t *testing.T) {
	initTestHome(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(board, []byte("---\noptions: [brainstorm]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", board)
	dir := t.TempDir()

	if out, err := runInit(t, dir); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCreate(t, dir, "A ticket"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}

	s, err := ticket.Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	tickets, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(tickets) != 1 {
		t.Fatalf("expected exactly one ticket, got %d", len(tickets))
	}
	full, err := s.Load(tickets[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(full.Body, "- [x] brainstorm") {
		t.Errorf("ticket body = %q, want brainstorm ticked", full.Body)
	}
}
