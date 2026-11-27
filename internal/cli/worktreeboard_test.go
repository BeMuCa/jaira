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
