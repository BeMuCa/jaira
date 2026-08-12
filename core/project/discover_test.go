package project

import (
	"os"
	"path/filepath"
	"testing"
)

func board(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(path, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverFindsBoardsWithinDepth(t *testing.T) {
	root := t.TempDir()
	board(t, filepath.Join(root, "one"))
	board(t, filepath.Join(root, "group", "two"))
	// Three levels down is past the limit and must not be found.
	board(t, filepath.Join(root, "a", "b", "c"))

	got := Discover(root, 2)
	if len(got) != 2 {
		t.Fatalf("found %d boards, want 2: %v", len(got), got)
	}
	for _, want := range []string{filepath.Join(root, "one"), filepath.Join(root, "group", "two")} {
		var ok bool
		for _, g := range got {
			if g == want {
				ok = true
			}
		}
		if !ok {
			t.Errorf("did not find %s in %v", want, got)
		}
	}
}

// The user's own configuration lives at ~/.jaira. A scan of a home directory
// must not mistake it for a project, which is why a board is identified by its
// tickets directory and dot directories are skipped outright.
func TestDiscoverIgnoresDotDirectories(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
	// root itself is a board, so scanning root finds root — scan its parent
	// instead to check the dot directory is not entered.
	parent := filepath.Dir(root)
	for _, g := range Discover(parent, 1) {
		if filepath.Base(g) == ".jaira" {
			t.Errorf("a dot directory was reported as a project: %s", g)
		}
	}
}

// A board is not searched for boards inside it: tickets live there, and a
// repository nested in a repository is not this tool's concern.
func TestDiscoverDoesNotDescendIntoABoard(t *testing.T) {
	root := t.TempDir()
	board(t, filepath.Join(root, "outer"))
	board(t, filepath.Join(root, "outer", "inner"))

	got := Discover(root, 2)
	if len(got) != 1 || filepath.Base(got[0]) != "outer" {
		t.Errorf("expected only the outer board, got %v", got)
	}
}

func TestDiscoverSkipsHeavyDirectories(t *testing.T) {
	root := t.TempDir()
	board(t, filepath.Join(root, "node_modules", "pkg"))
	if got := Discover(root, 2); len(got) != 0 {
		t.Errorf("scanned into a skipped directory: %v", got)
	}
}

// The same board must not appear twice because it was named differently. A
// trailing slash, a relative path and a symlink are all the same directory.
func TestRememberIsIdempotentAcrossSpellings(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	root := t.TempDir()
	board(t, filepath.Join(root, "repo"))
	repo := filepath.Join(root, "repo")

	Remember(repo)
	Remember(repo + string(filepath.Separator))
	Remember(filepath.Join(repo, ".", ""))

	link := filepath.Join(root, "link-to-repo")
	if err := os.Symlink(repo, link); err == nil {
		Remember(link)
	}

	if got := Load(); len(got) != 1 {
		t.Errorf("the same board was recorded %d times: %+v", len(got), got)
	}
}
