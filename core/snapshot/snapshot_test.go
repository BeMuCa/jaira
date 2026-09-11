package snapshot_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/snapshot"
)

// A bare repo plus a clone, because the snapshot is a branch on a remote and
// the whole claim is about what the remote ends up holding.
func board(t *testing.T) (*gitref.Repo, string, string) {
	t.Helper()
	for _, tool := range []string{"git"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not on PATH", tool)
		}
	}
	root := t.TempDir()
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "ada")
	run(t, root, "git", "clone", "--quiet", bare, clone)
	run(t, clone, "git", "config", "user.name", "ada")
	run(t, clone, "git", "config", "user.email", "ada@example.test")
	return &gitref.Repo{Dir: clone, Remote: "origin", AuthorName: "ada", AuthorEmail: "ada@example.test"}, clone, bare
}

func ticketFile(id, status string) []byte {
	return []byte("---\nid: " + id + "\ntitle: t " + id + "\nstatus: " + status + "\n---\n\n# t\n")
}

func TestASnapshotHoldsWhateverTheRefsHold(t *testing.T) {
	repo, clone, _ := board(t)
	for _, id := range []string{"01AAA", "01BBB"} {
		if _, err := repo.Write(id, ticketFile(id, "backlog"), ""); err != nil {
			t.Fatalf("write %s: %v", id, err)
		}
	}
	r := &snapshot.Runner{Repo: repo}

	first, err := r.Run()
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if first.Commit == "" || first.Tickets != 2 {
		t.Fatalf("first snapshot: %+v", first)
	}
	if len(first.Added) != 2 {
		t.Errorf("added %v, want both ids", first.Added)
	}
	// The files are under board/, deliberately not under .jaira/tickets/: at
	// the same path, merging this branch by accident would collide with every
	// working ticket at once.
	names := run(t, clone, "git", "ls-tree", "-r", "--name-only", snapshot.DefaultBranch)
	for _, want := range []string{"board/01AAA.md", "board/01BBB.md"} {
		if !strings.Contains(names, want) {
			t.Errorf("%s missing from the snapshot:\n%s", want, names)
		}
	}
	if strings.Contains(names, ".jaira/") {
		t.Errorf("the snapshot writes into the working ticket path:\n%s", names)
	}
	// The working tree is untouched: the run is plumbing only, because it
	// happens in the background of somebody else's command.
	if out := run(t, clone, "git", "status", "--short"); strings.TrimSpace(out) != "" {
		t.Errorf("the working tree was touched:\n%s", out)
	}

	// An unchanged set of refs writes nothing at all.
	again, err := r.Run()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !again.Unchanged || again.Commit != "" {
		t.Errorf("an unchanged board produced a second snapshot: %+v", again)
	}
	if n := strings.Count(run(t, clone, "git", "log", "--oneline", snapshot.DefaultBranch), "\n"); n != 1 {
		t.Errorf("the branch has %d commits, want 1", n)
	}
}

// Adding and removing come out of rebuilding the tree from the current refs,
// so a ticket whose ref is gone is simply not in the next snapshot — and stays
// readable in the previous ones.
func TestRemovingARefRemovesItFromTheNextSnapshotAndKeepsItInHistory(t *testing.T) {
	repo, clone, _ := board(t)
	for _, id := range []string{"01AAA", "01BBB"} {
		if _, err := repo.Write(id, ticketFile(id, "backlog"), ""); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	r := &snapshot.Runner{Repo: repo}
	if _, err := r.Run(); err != nil {
		t.Fatalf("first run: %v", err)
	}

	lease, err := repo.SHA("01AAA")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete("01AAA", lease); err != nil {
		t.Fatalf("delete ref: %v", err)
	}
	if _, err := repo.Write("01CCC", ticketFile("01CCC", "todo"), ""); err != nil {
		t.Fatalf("write new: %v", err)
	}

	second, err := r.Run()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if second.Commit == "" {
		t.Fatal("a changed board produced no snapshot")
	}
	if len(second.Added) != 1 || second.Added[0] != "01CCC" {
		t.Errorf("added %v", second.Added)
	}
	if len(second.Removed) != 1 || second.Removed[0] != "01AAA" {
		t.Errorf("removed %v", second.Removed)
	}
	names := run(t, clone, "git", "ls-tree", "-r", "--name-only", snapshot.DefaultBranch)
	if strings.Contains(names, "01AAA") {
		t.Errorf("the removed ticket is still in the newest snapshot:\n%s", names)
	}
	// But it is not lost: the previous snapshot still holds it, which is what
	// makes this a backup rather than a mirror.
	old := run(t, clone, "git", "ls-tree", "-r", "--name-only", snapshot.DefaultBranch+"~1")
	if !strings.Contains(old, "board/01AAA.md") {
		t.Errorf("the previous snapshot lost the ticket:\n%s", old)
	}
	// And git itself says what changed on the board.
	diff := run(t, clone, "git", "diff", "--name-status", snapshot.DefaultBranch+"~1", snapshot.DefaultBranch)
	if !strings.Contains(diff, "01AAA") || !strings.Contains(diff, "01CCC") {
		t.Errorf("the diff does not describe the change:\n%s", diff)
	}
}

// The first snapshot has no parent: the branch is its own history and shares
// none with the code.
func TestTheFirstSnapshotIsParentless(t *testing.T) {
	repo, clone, _ := board(t)
	if _, err := repo.Write("01AAA", ticketFile("01AAA", "backlog"), ""); err != nil {
		t.Fatal(err)
	}
	r := &snapshot.Runner{Repo: repo}
	if _, err := r.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	parents := strings.TrimSpace(run(t, clone, "git", "rev-list", "--parents", "-1", snapshot.DefaultBranch))
	if len(strings.Fields(parents)) != 1 {
		t.Errorf("the first snapshot has a parent: %q", parents)
	}
}

// Reaping: a ticket filed away in a landing branch has arrived, so its ref has
// done its job. One still on the board there has not.
func TestOnlyATicketFiledAwayInALandingBranchIsReaped(t *testing.T) {
	repo, clone, _ := board(t)

	// The landing branch: one ticket in the logbook, one still on the board.
	run(t, clone, "git", "checkout", "--quiet", "-b", "main")
	mkfile(t, clone, ".jaira/logbook/ada-20260910/01AAA-t.md", ticketFile("01AAA", "done"))
	mkfile(t, clone, ".jaira/tickets/01BBB-t.md", ticketFile("01BBB", "todo"))
	run(t, clone, "git", "add", ".jaira")
	run(t, clone, "git", "commit", "--quiet", "-m", "land 01AAA, keep 01BBB on the board")

	for _, id := range []string{"01AAA", "01BBB"} {
		if _, err := repo.Write(id, ticketFile(id, "done"), ""); err != nil {
			t.Fatalf("write %s: %v", id, err)
		}
	}

	r := &snapshot.Runner{Repo: repo, LandingBranches: []string{"main"}}
	res, err := r.Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(res.Reaped) != 1 || res.Reaped[0] != "01AAA" {
		t.Fatalf("reaped %v, want only the landed ticket", res.Reaped)
	}
	if _, err := repo.SHA("01AAA"); err == nil {
		t.Error("the landed ticket's ref survived")
	}
	if _, err := repo.SHA("01BBB"); err != nil {
		t.Errorf("a ticket still on the board there lost its ref: %v", err)
	}
	// And the reaped ticket is in the snapshot written moments before, which is
	// what makes deleting the ref safe rather than lossy.
	names := run(t, clone, "git", "ls-tree", "-r", "--name-only", snapshot.DefaultBranch)
	if !strings.Contains(names, "board/01AAA.md") {
		t.Errorf("the ref was deleted without the snapshot holding it:\n%s", names)
	}
}

// With no landing branch resolved, nothing is reaped. A ref left standing costs
// nothing; one removed by mistake takes away the visibility this is for.
func TestWithoutALandingBranchNothingIsReaped(t *testing.T) {
	repo, _, _ := board(t)
	if _, err := repo.Write("01AAA", ticketFile("01AAA", "done"), ""); err != nil {
		t.Fatal(err)
	}
	r := &snapshot.Runner{Repo: repo}
	res, err := r.Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(res.Reaped) != 0 {
		t.Errorf("reaped %v with no landing branch configured", res.Reaped)
	}
	if _, err := repo.SHA("01AAA"); err != nil {
		t.Errorf("the ref was removed anyway: %v", err)
	}
}

// Drop is the explicit way out for a branch that will never be merged.
func TestDropRemovesARefOnPurpose(t *testing.T) {
	repo, _, _ := board(t)
	if _, err := repo.Write("01AAA", ticketFile("01AAA", "done"), ""); err != nil {
		t.Fatal(err)
	}
	r := &snapshot.Runner{Repo: repo}
	if err := r.Drop("01AAA"); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if _, err := repo.SHA("01AAA"); err == nil {
		t.Error("drop left the ref in place")
	}
}

func mkfile(t *testing.T, root, rel string, content []byte) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=ada", "GIT_AUTHOR_EMAIL=ada@example.test",
		"GIT_COMMITTER_NAME=ada", "GIT_COMMITTER_EMAIL=ada@example.test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}
