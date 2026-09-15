package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/milestone"
)

// runAndSend runs a command the way the binary does: the queued write is sent
// after the command, not during it, so a test that only calls runCLI leaves
// everything sitting in the outbox.
func runAndSend(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	out, err := runCLI(t, dir, args...)
	flushRefs()
	return out, err
}

// twoBoards builds one bare remote and two real clones with a board in each,
// standing in for two teammates. Nothing is faked: the claim under test is
// about git refs, and a fake would only prove the fake behaves as assumed.
func twoBoards(t *testing.T) (ada, grace string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_HOME", home)
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(root, "no-lanes"))

	bare := filepath.Join(root, "board.git")
	gitRun(t, root, "init", "--bare", "--quiet", bare)
	mk := func(name string) string {
		dir := filepath.Join(root, name)
		gitRun(t, root, "clone", "--quiet", bare, dir)
		gitRun(t, dir, "config", "user.name", name)
		gitRun(t, dir, "config", "user.email", name+"@example.test")
		if out, err := runCLI(t, dir, "init"); err != nil {
			t.Fatalf("init %s: %v\n%s", name, err, out)
		}
		return dir
	}
	return mk("ada"), mk("grace")
}

// The reason a milestone travels on a ref: it is the one file everybody plans
// from, so it has to arrive without anyone merging a branch. If it waited for
// a merge, the file everybody is supposed to read would be the file nobody
// has.
func TestAMilestoneReachesTheOtherCloneWithoutAMerge(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	// Nothing was committed, let alone merged: an unshared board gitignores
	// .jaira entirely, so the milestone file is in no commit anywhere and the
	// ref is the only way it can travel.
	if out := gitOut(t, ada, "for-each-ref", "refs/heads/"); strings.TrimSpace(out) != "" {
		t.Fatalf("the clone has branches, so this test would not prove the ref carried it:\n%s", out)
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	ms, err := milestone.Load(grace, "round-one")
	if err != nil {
		t.Fatalf("grace does not have the milestone after a fetch: %v", err)
	}
	if ms.Name != "round-one" {
		t.Errorf("grace's milestone is %q", ms.Name)
	}
	out, err := runCLI(t, grace, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one") {
		t.Errorf("'milestone ls' in grace's clone does not show it:\n%s", out)
	}
}

// Adding tickets to a milestone is an edit of the same file, and it has to
// reach the other side the same way — otherwise grouping twenty tickets is one
// edit that only one person can see.
func TestAddingToAMilestoneReachesTheOtherClone(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "create", "the thing", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create ticket: %v\n%s", err, out)
	}
	h := handleFromList(t, ada, "the thing")
	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create milestone: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "milestone", "add", "round-one", h); err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	ms, err := milestone.Load(grace, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if len(ms.Members()) != 1 {
		t.Fatalf("grace's copy holds %v, want the one ticket", ms.Members())
	}
	// And the ticket is filterable by it there, which is the thing the file
	// was written for.
	got := jsonCLI(t, grace, "list", "--milestone", "round-one")
	if rows, _ := got["tickets"].([]any); len(rows) != 1 {
		t.Errorf("grace's 'list --milestone round-one' returned %d tickets, want 1", len(rows))
	}
}

// handleFromList reads a ticket's handle back through the CLI. A board with a
// remote keeps an unworked ticket on its ref and not on disk, so reading the
// tickets directory would find nothing.
func handleFromList(t *testing.T, dir, title string) string {
	t.Helper()
	got := jsonCLI(t, dir, "list")
	for _, r := range got["tickets"].([]any) {
		row, _ := r.(map[string]any)
		if row["title"] == title {
			h, _ := row["handle"].(string)
			if h == "" {
				h, _ = row["id"].(string)
			}
			return h
		}
	}
	t.Fatalf("no ticket titled %q in 'jaira list' from %s", title, dir)
	return ""
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}
