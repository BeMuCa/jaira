package cli

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Accepting a ticket into done from the worktree of a branch whose name is
// not the ticket's id writes every commit on that branch that names the
// ticket in its subject — the handle alone, as people write it. Before
// 8FE2EG the gate found them and the move wrote nothing, and the later
// stamp, run from another checkout, wrote an empty list.
func TestMoveToDoneInAWorktreeWritesTheBranchCommits(t *testing.T) {
	if !gitrepo.Available() {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	main := filepath.Join(root, "repo")
	wt := filepath.Join(root, "wt")
	t.Setenv("JAIRA_USER", "Alexander Sacharov")
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	git := func(dir string, args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("git", "init", "-q", "-b", "master", main).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	git(main, "config", "user.email", "fixture@example.com")
	git(main, "config", "user.name", "Fixture")

	s, err := ticket.At(main)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	id := ticket.NewID(time.Now())
	if _, err := s.Create(map[string]string{
		ticket.FieldID: id, ticket.FieldTitle: "t", ticket.FieldStatus: "signoff",
	}, nil, ""); err != nil {
		t.Fatal(err)
	}
	git(main, "add", "-f", ".")
	git(main, "commit", "-q", "-m", "chore: board")

	git(main, "worktree", "add", "-q", "-b", "feat/another-name", wt)
	h := ticket.Handle(id)
	var want []string
	for _, subj := range []string{"fix(" + h + "): first", "feat(" + h + "): second"} {
		git(wt, "commit", "-q", "--allow-empty", "-m", subj)
		want = append(want, git(wt, "rev-parse", "HEAD"))
	}

	if out, err := runCLI(t, wt, "move", h, "--to", "done", "--force"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	ws, err := ticket.At(wt)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ws.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	for _, sha := range want {
		found := false
		for _, c := range got.Commits {
			if c == sha {
				found = true
			}
		}
		if !found {
			t.Errorf("commits = %v, want it to hold %s from feat/another-name", got.Commits, sha)
		}
	}
}
