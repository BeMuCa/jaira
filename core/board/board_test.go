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
