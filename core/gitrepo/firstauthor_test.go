package gitrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// The board's creator is whoever committed its lanes first, not whoever
// touched them last; an uncommitted path has none.
func TestFirstAuthor(t *testing.T) {
	dir, repo := fixtureRepo(t)
	if got := repo.FirstAuthor("lanes"); got != "" {
		t.Errorf("empty repo: FirstAuthor = %q, want empty", got)
	}
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureCommit(t, dir, map[string]string{"lanes/a.md": "a"}, "first", t0)
	if out, err := exec.Command("git", "-C", dir, "config", "user.name", "Later").CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	fixtureCommit(t, dir, map[string]string{"lanes/b.md": "b"}, "second", t0.Add(time.Hour))
	if got := repo.FirstAuthor("lanes"); got != "Fixture" {
		t.Errorf("FirstAuthor = %q, want Fixture", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "new.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := repo.FirstAuthor("new.md"); got != "" {
		t.Errorf("uncommitted file: FirstAuthor = %q, want empty", got)
	}
}
