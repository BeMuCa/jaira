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

// TestDiffLeavesTicketFilesOut pins that the commit diff a lane reads carries
// the code and not the ticket file committed beside it: the ticket's patch
// holds the implementer's outcome, which the review lane must not read.
func TestDiffLeavesTicketFilesOut(t *testing.T) {
	dir, repo := fixtureRepo(t)
	fixtureCommit(t, dir, map[string]string{"base.txt": "base\n"}, "chore: base",
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	head := fixtureCommit(t, dir, map[string]string{
		"code.go":                                "package code\n",
		".jaira/tickets/AAA.md":                  "outcome-what: the author's account\n",
		".jaira/logbook/someone-2026-01-02/B.md": "outcome-what: a filed account\n",
		".jaira/archive/C.md":                    "outcome-what: a rejected account\n",
	}, "feat(AAA): code", time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC))

	diff, unavailable := repo.Diff([]string{head})
	if len(unavailable) != 0 {
		t.Fatalf("unavailable: %v", unavailable)
	}
	if !strings.Contains(diff, "package code") {
		t.Errorf("the code is missing from the diff:\n%s", diff)
	}
	if strings.Contains(diff, "account") || strings.Contains(diff, ".jaira/") {
		t.Errorf("the ticket file is in the diff:\n%s", diff)
	}
}

// TestDiffNeverHandsGitAnOption pins that a value in a ticket's hand-editable
// commits: field reaches git only as a commit id: "--output=<path>" used to make
// 'git show' truncate or create that file and exit 0.
func TestDiffNeverHandsGitAnOption(t *testing.T) {
	dir, repo := fixtureRepo(t)
	fixtureCommit(t, dir, map[string]string{"base.txt": "base\n"}, "chore: base",
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	victim := filepath.Join(t.TempDir(), "victim")
	if err := os.WriteFile(victim, []byte("keep me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	bad := "--output=" + victim

	_, unavailable := repo.Diff([]string{bad})
	if len(unavailable) != 1 || unavailable[0] != bad {
		t.Errorf("unavailable = %v, want the option reported as unavailable", unavailable)
	}
	if _, err := repo.Stat([]string{bad}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Commits([]string{bad}); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(victim); err != nil || string(b) != "keep me\n" {
		t.Errorf("git wrote over the file: %q, %v", b, err)
	}
}

// TestDiffIsEmptyForTicketOnlyCommits pins that a commit which touched only its
// ticket file adds nothing to Diff, not even a blank line: the caller tells
// "no code to judge" from an empty string, and "\n" used to hide that.
func TestDiffIsEmptyForTicketOnlyCommits(t *testing.T) {
	dir, repo := fixtureRepo(t)
	fixtureCommit(t, dir, map[string]string{"base.txt": "base\n"}, "chore: base",
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	sha := fixtureCommit(t, dir, map[string]string{".jaira/tickets/AAA.md": "note\n"},
		"chore(AAA): ticket only", time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC))

	diff, unavailable := repo.Diff([]string{sha})
	if diff != "" || len(unavailable) != 0 {
		t.Errorf("Diff = %q, unavailable %v; want empty for a ticket-only commit", diff, unavailable)
	}
}
