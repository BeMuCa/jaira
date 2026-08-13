package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareMakesANewBoardPrivateAndAnnounced(t *testing.T) {
	root := t.TempDir()

	p := Prepare(root)
	if !p.Ignored {
		t.Errorf("Ignored = false, want true on a fresh board")
	}
	if p.IgnoreErr != nil {
		t.Errorf("IgnoreErr = %v, want nil", p.IgnoreErr)
	}
	if p.NoteErr != nil {
		t.Errorf("NoteErr = %v, want nil", p.NoteErr)
	}
	if len(p.Notes) != 2 {
		t.Errorf("Notes = %v, want two entries", p.Notes)
	}

	gi, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if !strings.Contains(string(gi), IgnoreLine) {
		t.Errorf(".gitignore = %q, want it to contain %q", gi, IgnoreLine)
	}

	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		if !strings.Contains(string(b), "<!-- jaira:start -->") {
			t.Errorf("%s does not contain the jaira marker", name)
		}
	}

	// Preparing an already-prepared board is a no-op: nothing new is ignored,
	// nothing new is announced, and the gitignore entry does not duplicate.
	p2 := Prepare(root)
	if p2.Ignored {
		t.Errorf("Ignored = true on a second Prepare, want false")
	}
	if len(p2.Notes) != 0 {
		t.Errorf("Notes = %v on a second Prepare, want none", p2.Notes)
	}

	gi2, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if n := strings.Count(string(gi2), IgnoreLine); n != 1 {
		t.Errorf(".gitignore contains %s %d times, want exactly 1: %q", IgnoreLine, n, gi2)
	}
}

// TestAddLanesIgnoreOnFreshBoard asserts that on a shared board — one where
// RemoveIgnore has already stopped ignoring the whole /.jaira/ tree —
// AddLanesIgnore still ignores /.jaira/lanes/ on its own.
func TestAddLanesIgnoreOnFreshBoard(t *testing.T) {
	root := t.TempDir()

	if _, err := RemoveIgnore(root); err != nil {
		t.Fatalf("RemoveIgnore: %v", err)
	}
	changed, err := AddLanesIgnore(root)
	if err != nil {
		t.Fatalf("AddLanesIgnore: %v", err)
	}
	if !changed {
		t.Errorf("AddLanesIgnore changed = false on a fresh board, want true")
	}

	gi, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	for _, line := range strings.Split(string(gi), "\n") {
		if strings.TrimSpace(line) == IgnoreLine {
			t.Errorf(".gitignore = %q, must not ignore the whole board", gi)
		}
	}
	if !strings.Contains(string(gi), LanesIgnoreLine) {
		t.Errorf(".gitignore = %q, want it to contain %q", gi, LanesIgnoreLine)
	}
}

// TestAddLanesIgnoreIsIdempotent asserts a second call reports changed=false
// and does not duplicate the line.
func TestAddLanesIgnoreIsIdempotent(t *testing.T) {
	root := t.TempDir()

	if _, err := AddLanesIgnore(root); err != nil {
		t.Fatalf("AddLanesIgnore (first): %v", err)
	}
	changed, err := AddLanesIgnore(root)
	if err != nil {
		t.Fatalf("AddLanesIgnore (second): %v", err)
	}
	if changed {
		t.Errorf("AddLanesIgnore changed = true on the second call, want false")
	}
	gi, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if n := strings.Count(string(gi), LanesIgnoreLine); n != 1 {
		t.Errorf(".gitignore contains %s %d times, want exactly 1: %q", LanesIgnoreLine, n, gi)
	}
}

// TestRemoveLanesIgnoreDropsTheLine asserts the share --undo path removes the
// lanes-only line, since /.jaira/ already covers it once the board itself is
// private again and a leftover line would be a puzzle for the next reader.
func TestRemoveLanesIgnoreDropsTheLine(t *testing.T) {
	root := t.TempDir()

	if _, err := AddLanesIgnore(root); err != nil {
		t.Fatalf("AddLanesIgnore: %v", err)
	}
	changed, err := RemoveLanesIgnore(root)
	if err != nil {
		t.Fatalf("RemoveLanesIgnore: %v", err)
	}
	if !changed {
		t.Errorf("RemoveLanesIgnore changed = false, want true")
	}
	gi, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if strings.Contains(string(gi), LanesIgnoreLine) {
		t.Errorf(".gitignore = %q, want the lanes-only line gone", gi)
	}
}

// TestIgnoredNotFooledByLanesLine asserts Ignored(root) reports true only for
// the whole-board entry. isShared and bindDriverIfShared both branch on this,
// so a false positive would silently disable the merge driver on a shared
// board that merely keeps its lanes private.
func TestIgnoredNotFooledByLanesLine(t *testing.T) {
	root := t.TempDir()

	if _, err := AddLanesIgnore(root); err != nil {
		t.Fatalf("AddLanesIgnore: %v", err)
	}
	if Ignored(root) {
		t.Errorf("Ignored(root) = true with only the lanes-only line present, want false")
	}

	if _, err := AddIgnore(root); err != nil {
		t.Fatalf("AddIgnore: %v", err)
	}
	if !Ignored(root) {
		t.Errorf("Ignored(root) = false once the whole-board line is present, want true")
	}
}
