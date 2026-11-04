package gitrepo

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// branchFixture: master has a.txt and b.txt; feat/abc123-x changes a.txt and
// adds a file under .jaira; master then moves on by changing b.txt, which the
// branch diff must not show.
func branchFixture(t *testing.T) (string, *Repo) {
	t.Helper()
	dir, repo := fixtureRepo(t)
	when := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureCommit(t, dir, map[string]string{"a.txt": "one\n", "b.txt": "one\n"}, "base", when)
	gitIn(t, dir, "branch", "-M", "master")
	gitIn(t, dir, "checkout", "-q", "-b", "feat/abc123-x")
	fixtureCommit(t, dir, map[string]string{
		"a.txt":               "one\ntwo\n",
		".jaira/tickets/x.md": "ticket\n",
	}, "feat(ABC123): add two", when.Add(time.Hour))
	gitIn(t, dir, "checkout", "-q", "master")
	fixtureCommit(t, dir, map[string]string{"b.txt": "changed on master\n"}, "master moves on", when.Add(2*time.Hour))
	return dir, repo
}

func TestBranchForAndBaseFindTheTicketBranch(t *testing.T) {
	_, repo := branchFixture(t)
	if got := repo.BranchFor("ABC123"); got != "feat/abc123-x" {
		t.Errorf("BranchFor(ABC123) = %q, want feat/abc123-x", got)
	}
	if got := repo.BranchFor("ZZZ999"); got != "" {
		t.Errorf("BranchFor(ZZZ999) = %q, want none", got)
	}
	if got := repo.Base(); got != "master" {
		t.Errorf("Base() = %q, want master", got)
	}
}

func TestBranchDiffIsAgainstMergeBaseAndSkipsJaira(t *testing.T) {
	_, repo := branchFixture(t)
	patch, err := repo.BranchDiff("master", "feat/abc123-x")
	if err != nil {
		t.Fatalf("BranchDiff: %v", err)
	}
	if !strings.Contains(patch, "a.txt") || !strings.Contains(patch, "+two") {
		t.Errorf("patch misses the branch change to a.txt:\n%s", patch)
	}
	if strings.Contains(patch, "b.txt") {
		t.Errorf("patch shows master's own change to b.txt:\n%s", patch)
	}
	if strings.Contains(patch, ".jaira") {
		t.Errorf("patch shows a file under .jaira:\n%s", patch)
	}
}

func TestWorktreeOfFindsTheCheckoutOfTheBranch(t *testing.T) {
	dir, repo := branchFixture(t)
	wt := filepath.Join(t.TempDir(), "wt")
	gitIn(t, dir, "worktree", "add", "-q", wt, "feat/abc123-x")
	got := repo.WorktreeOf("feat/abc123-x")
	if got == "" {
		t.Fatal("WorktreeOf found no worktree for feat/abc123-x")
	}
	gotReal, _ := filepath.EvalSymlinks(got)
	wantReal, _ := filepath.EvalSymlinks(wt)
	if gotReal != wantReal {
		t.Errorf("WorktreeOf = %q, want %q", gotReal, wantReal)
	}
	if got := repo.WorktreeOf("nope"); got != "" {
		t.Errorf("WorktreeOf(nope) = %q, want none", got)
	}
}

// When several branches carry the handle, the one with the newest commit wins,
// not the first by name: an old experiment must not shadow the current work.
func TestBranchForPicksTheNewestBranchCarryingTheHandle(t *testing.T) {
	dir, repo := fixtureRepo(t)
	when := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureCommit(t, dir, map[string]string{"a.txt": "one\n"}, "base", when)
	gitIn(t, dir, "branch", "-M", "master")
	gitIn(t, dir, "checkout", "-q", "-b", "exp/abc123-a")
	fixtureCommit(t, dir, map[string]string{"a.txt": "exp\n"}, "exp", when.Add(time.Hour))
	gitIn(t, dir, "checkout", "-q", "master")
	gitIn(t, dir, "checkout", "-q", "-b", "feat/abc123-b")
	fixtureCommit(t, dir, map[string]string{"a.txt": "feat\n"}, "feat", when.Add(2*time.Hour))
	gitIn(t, dir, "checkout", "-q", "master")
	if got := repo.BranchFor("ABC123"); got != "feat/abc123-b" {
		t.Errorf("BranchFor(ABC123) = %q, want the newer feat/abc123-b, not exp/abc123-a", got)
	}
}
