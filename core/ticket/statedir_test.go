package ticket_test

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/snapshot"
	"github.com/BeMuCa/jaira/core/ticket"
)

// TestRepoStateDirIsSharedByWorktrees is the regression the whole change is
// for: a second working tree of one clone must land on the same state
// directory, because a background job's clock belongs to the board rather than
// to the checkout that happened to trigger it.
func TestRepoStateDirIsSharedByWorktrees(t *testing.T) {
	main, tree := clonePair(t)

	a := (&ticket.Store{Root: main}).RepoStateDir()
	b := (&ticket.Store{Root: tree}).RepoStateDir()
	if a != b {
		t.Fatalf("worktrees disagree on the repository state dir:\n  %s\n  %s", a, b)
	}
	if x, y := (&ticket.Store{Root: main}).StateDir(), (&ticket.Store{Root: tree}).StateDir(); x == y {
		t.Fatalf("worktrees share the per-tree state dir %q; sessions and locks must stay separate", x)
	}
}

// TestSnapshotClockIsSharedAcrossWorktrees checks the behaviour that mattered:
// after a run stamped in one worktree, the other must not think a snapshot is
// due.
func TestSnapshotClockIsSharedAcrossWorktrees(t *testing.T) {
	main, tree := clonePair(t)

	first := (&ticket.Store{Root: main}).RepoStateDir()
	if !snapshot.Due(first, time.Hour) {
		t.Fatal("a board that has never snapshotted should be due")
	}
	if err := snapshot.WriteStamp(first, snapshot.Stamp{RanAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if second := (&ticket.Store{Root: tree}).RepoStateDir(); snapshot.Due(second, time.Hour) {
		t.Fatal("a fresh worktree started its own clock: the snapshot is due again")
	}
}

// TestRepoStateDirSeparatesClones guards the other direction: two clones are two
// boards and must not share a clock.
func TestRepoStateDirSeparatesClones(t *testing.T) {
	one, _ := clonePair(t)
	two, _ := clonePair(t)
	if a, b := (&ticket.Store{Root: one}).RepoStateDir(), (&ticket.Store{Root: two}).RepoStateDir(); a == b {
		t.Fatalf("two clones share %q", a)
	}
}

// TestRepoStateDirOutsideGitFallsBack keeps a board that is not in a repository
// working: there is one working tree there, so the per-tree directory is the
// right answer rather than an error.
func TestRepoStateDirOutsideGitFallsBack(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s := &ticket.Store{Root: dir}
	if got, want := s.RepoStateDir(), s.StateDir(); got != want {
		t.Fatalf("outside a repository: got %q, want the per-tree dir %q", got, want)
	}
}

// clonePair builds a repository with one extra worktree and returns both roots.
func clonePair(t *testing.T) (main, tree string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))

	main = t.TempDir()
	git(t, main, "init")
	git(t, main, "config", "user.email", "test@example.com")
	git(t, main, "config", "user.name", "test")
	git(t, main, "commit", "--allow-empty", "-m", "root")

	tree = filepath.Join(t.TempDir(), "wt")
	git(t, main, "worktree", "add", "-b", "side", tree)
	return main, tree
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
