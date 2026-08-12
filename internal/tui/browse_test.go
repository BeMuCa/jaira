package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A board created from the browse screen's 'i' key must be gitignored and
// announced in the agent instruction files, exactly like one created by
// 'jaira init' — this is the regression guard for that bug.
func TestBrowseKeyInitPreparesTheBoard(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	tmp := t.TempDir()

	b := newBrowser(tmp)
	added, done := b.key("i")
	if !done {
		t.Fatalf("done = false, want true")
	}
	if len(added) != 1 || added[0] != tmp {
		t.Fatalf("added = %v, want [%s]", added, tmp)
	}

	if fi, err := os.Stat(filepath.Join(tmp, ".jaira", "tickets")); err != nil || !fi.IsDir() {
		t.Errorf(".jaira/tickets was not created: %v", err)
	}

	gi, err := os.ReadFile(filepath.Join(tmp, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if !strings.Contains(string(gi), "/.jaira/") {
		t.Errorf(".gitignore = %q, want it to contain /.jaira/", gi)
	}

	claude, err := os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("reading CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(claude), "<!-- jaira:start -->") {
		t.Errorf("CLAUDE.md does not contain the jaira marker")
	}
}
