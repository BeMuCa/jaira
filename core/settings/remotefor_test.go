package settings_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/settings"
)

// repoWith builds a real repository carrying exactly the named remotes. Real,
// because what is under test is how jaira reads a checkout, and a fake would
// only prove that the fake answers as assumed.
func repoWith(t *testing.T, remotes ...string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	dir := t.TempDir()
	git(t, dir, "init", "--quiet")
	for _, name := range remotes {
		git(t, dir, "remote", "add", name, "https://example.test/"+name+".git")
	}
	return dir
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// The bug this feature exists for: one machine-wide "remote": "upstream", set
// for the one checkout that has a fork as its origin, took every other board
// down with "no remote \"upstream\"". A repository with a single remote has
// nothing to get wrong, so it works without anybody configuring it.
func TestRemoteForUsesTheOnlyRemoteWhenTheSettingIsNotTrueHere(t *testing.T) {
	dir := repoWith(t, "origin")
	s := settings.Settings{Remote: "upstream"}
	if got := s.RemoteFor(dir); got != "origin" {
		t.Errorf("a repository with only origin resolved to %q", got)
	}
}

// And the other half, which is what stops the fix from being "fall back to
// origin everywhere": where the configured remote does exist, it wins, so a
// ticket ref in a fork still goes upstream and not into the fork.
func TestRemoteForKeepsTheConfiguredRemoteWhereItExists(t *testing.T) {
	dir := repoWith(t, "origin", "upstream")
	s := settings.Settings{Remote: "upstream"}
	if got := s.RemoteFor(dir); got != "upstream" {
		t.Errorf("a fork with an upstream resolved to %q, which is the fork", got)
	}
}

// Several remotes and none of them the one asked for is ambiguous. Picking one
// would be the fork case again, so the configured name is kept and the caller
// stops on it loudly.
func TestRemoteForDoesNotGuessBetweenSeveralRemotes(t *testing.T) {
	dir := repoWith(t, "origin", "fork")
	s := settings.Settings{Remote: "upstream"}
	if got := s.RemoteFor(dir); got != "upstream" {
		t.Errorf("jaira guessed %q instead of stopping", got)
	}
}

// A name set for this clone is a decision about this repository and never falls
// back — not to the machine-wide setting, not to the only remote there is.
func TestBoardRemoteWinsAndNeverFallsBack(t *testing.T) {
	dir := repoWith(t, "origin", "upstream")
	git(t, dir, "config", "--local", "jaira.remote", "origin")
	s := settings.Settings{Remote: "upstream"}
	if got := s.RemoteFor(dir); got != "origin" {
		t.Errorf("git config jaira.remote was ignored: %q", got)
	}

	gone := repoWith(t, "origin")
	git(t, gone, "config", "--local", "jaira.remote", "upstream")
	if got := (settings.Settings{}).RemoteFor(gone); got != "upstream" {
		t.Errorf("a board-set remote fell back to %q; it must fail loudly instead", got)
	}
}

// Every worktree of a clone shares .git/config, so the board remote is set once
// and every worktree of that checkout answers the same.
func TestBoardRemoteIsSharedByWorktrees(t *testing.T) {
	dir := repoWith(t, "origin", "upstream")
	git(t, dir, "config", "--local", "jaira.remote", "upstream")
	git(t, dir, "config", "user.email", "t@example.test")
	git(t, dir, "config", "user.name", "t")
	git(t, dir, "commit", "--quiet", "--allow-empty", "-m", "root")
	tree := filepath.Join(t.TempDir(), "wt")
	git(t, dir, "worktree", "add", "--quiet", "-b", "side", tree)

	if got := (settings.Settings{}).RemoteFor(tree); got != "upstream" {
		t.Errorf("the worktree resolved to %q instead of the clone's setting", got)
	}
}
