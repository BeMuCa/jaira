package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A worktree cut before 'jaira share' - or of a board never shared - has no
// .jaira of its own. Commands run there work on the main checkout's board, as
// the cockpit's Accept does when it moves a ticket inside the ticket's worktree.
func TestWorktreeWithoutBoardUsesTheMainCheckouts(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "create", "from main", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	ts, err := s.List()
	if err != nil || len(ts) != 1 {
		t.Fatalf("List = %v, %v", ts, err)
	}
	id := ticket.Handle(ts[0].ID)

	wt := filepath.Join(t.TempDir(), "wt")
	gitRun(t, dir, "worktree", "add", "--quiet", "-b", "feat/x", wt)
	if _, err := os.Stat(filepath.Join(wt, ticket.DirName)); !os.IsNotExist(err) {
		t.Fatalf("the worktree has a .jaira already (%v): the test needs one without", err)
	}

	if out, err := runCLI(t, wt, "list"); err != nil || !strings.Contains(out, id) {
		t.Fatalf("list in the worktree: %v\n%s", err, out)
	}
	if out, err := runCLI(t, wt, "show", id); err != nil || !strings.Contains(out, "from main") {
		t.Fatalf("show in the worktree: %v\n%s", err, out)
	}
	if out, err := runCLI(t, wt, "move", id, "--to", "blocked", "--reason", "r"); err != nil {
		t.Fatalf("move in the worktree: %v\n%s", err, out)
	}
	got, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "blocked" {
		t.Fatalf("status on the main checkout's board = %q, want blocked", got.Status)
	}
	if _, err := os.Stat(filepath.Join(wt, ticket.DirName)); !os.IsNotExist(err) {
		t.Fatalf("the worktree grew a .jaira: %v", err)
	}
}

// The board is the main checkout's, but git questions - the branch, its
// commits, its diff - are the worktree's. With the main checkout on another
// ticket's branch, a ticket filed in the worktree follows the worktree's
// ticket, and the worktree's commit is the one a review lane is shown.
func TestWorktreeWithoutBoardAsksGitAboutItself(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	handle := func(title string) string {
		t.Helper()
		out, err := runCLI(t, dir, "create", title, "--goal", "g", "--context", "c", "--dod", "d")
		if err != nil {
			t.Fatalf("create: %v\n%s", err, out)
		}
		f := strings.Fields(out)
		if len(f) < 2 {
			t.Fatalf("create printed no handle:\n%s", out)
		}
		return f[1]
	}
	mainID, wtID := handle("main's"), handle("worktree's")
	gitRun(t, dir, "checkout", "--quiet", "-b", "feat/"+strings.ToLower(mainID))

	wt := filepath.Join(t.TempDir(), "wt")
	gitRun(t, dir, "worktree", "add", "--quiet", "-b", "feat/"+strings.ToLower(wtID), wt, "HEAD")

	out, err := runCLI(t, wt, "create", "found on the way", "--goal", "g", "--context", "c", "--dod", "d")
	if err != nil {
		t.Fatalf("create in the worktree: %v\n%s", err, out)
	}
	if !strings.Contains(out, wtID) || strings.Contains(out, mainID) {
		t.Fatalf("create in the worktree on feat/%s should follow %s, not %s:\n%s", strings.ToLower(wtID), wtID, mainID, out)
	}

	if err := os.WriteFile(filepath.Join(wt, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wt, "add", "work.txt")
	gitRun(t, wt, "commit", "--quiet", "-m", "fix("+wtID+"): work")
	out, err = runCLI(t, wt, "show", wtID, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	if !strings.Contains(out, "work.txt") {
		t.Fatalf("the review payload does not carry the worktree's commit:\n%s", out)
	}
}
