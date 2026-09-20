package gitrepo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestWorktreeDiffKeepsUntrackedPathsGitWouldQuote pins the half of
// WorktreeDiff that reads untracked files. git quotes a path that is not plain
// ASCII unless it is asked not to, and a quoted path is not a path any later
// command can open — so the file used to fall out of the patch while the
// payload went on claiming to carry the worktree. That is the fragment
// reported as the whole, one step further on, which is exactly what this
// package must not do.
func TestWorktreeDiffKeepsUntrackedPathsGitWouldQuote(t *testing.T) {
	dir, repo := fixtureRepo(t)
	fixtureCommit(t, dir, map[string]string{"base.txt": "base\n"}, "chore: base",
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))

	for _, name := range []string{"Änderung.txt", "neue datei.txt", "plain.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content of "+name+"\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	diff, err := repo.WorktreeDiff()
	if err != nil {
		t.Fatalf("WorktreeDiff: %v", err)
	}
	for _, name := range []string{"Änderung.txt", "neue datei.txt", "plain.txt"} {
		if !strings.Contains(diff, "content of "+name) {
			t.Errorf("untracked %q is missing from the worktree diff:\n%s", name, diff)
		}
	}
	if strings.Contains(diff, `\303\204`) {
		t.Errorf("the diff carries a quoted path instead of the file's contents:\n%s", diff)
	}
}

// TestWorktreeDiffIgnoresAnEmptyUntrackedFile guards the other side of the
// same check: --no-index produces no patch for a file with no bytes in it, and
// that is not the unreadable-path failure the empty patch otherwise means.
func TestWorktreeDiffIgnoresAnEmptyUntrackedFile(t *testing.T) {
	dir, repo := fixtureRepo(t)
	fixtureCommit(t, dir, map[string]string{"base.txt": "base\n"}, "chore: base",
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	if err := os.WriteFile(filepath.Join(dir, "empty.txt"), nil, 0o644); err != nil {
		t.Fatalf("write empty.txt: %v", err)
	}

	if _, err := repo.WorktreeDiff(); err != nil {
		t.Fatalf("an empty untracked file must not be an error: %v", err)
	}
}
