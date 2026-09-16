package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
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

// Filing a milestone must not take its ref down, and this is why: a ref is
// what every other clone reads. Taken down, the name is free again and two
// machines can plan two different milestones under one identity. Left up and
// marked, it says "filed" to everyone who fetches — the file is not written to
// their board, and the name stays taken.
func TestAFiledMilestoneStaysOffTheOtherCloneAndKeepsItsRef(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if _, err := milestone.Load(ada, "round-one"); err == nil {
		t.Error("the milestone is still on ada's board after being filed")
	}
	out, err := runCLI(t, ada, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if strings.Contains(out, "round-one") {
		t.Errorf("'milestone ls' still names a filed milestone:\n%s", out)
	}

	// The ref is up and says so. Read through gitref rather than the file,
	// because the ref is the only thing the other clone will ever see.
	repo := &gitref.Repo{Dir: ada, Remote: "origin"}
	content, _, err := repo.ReadMilestone("round-one")
	if err != nil {
		t.Fatalf("the ref of a filed milestone was taken down: %v", err)
	}
	if got := milestone.FromBytes("round-one", content); !got.Filed() {
		t.Errorf("the ref carries status %q, want %q", got.Status, milestone.StatusFiled)
	}

	// And grace, who never saw it on her board, does not get it written there.
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if _, err := milestone.Load(grace, "round-one"); err == nil {
		t.Error("a fetch wrote a filed milestone onto grace's board")
	}
	out, err = runCLI(t, grace, "milestone", "ls")
	if err != nil {
		t.Fatalf("grace's ls: %v\n%s", err, out)
	}
	if strings.Contains(out, "round-one") {
		t.Errorf("grace's 'milestone ls' names a filed milestone:\n%s", out)
	}
}

// The other half of the same claim, and the reason the ref carries a status at
// all: a clone that ALREADY has the file on its board has to learn that the
// milestone was filed. Skipping the marked ref would leave grace planning a
// round of work ada closed, with nothing in either board ever telling her.
func TestAClonePlanningTheMilestoneLearnsItWasFiled(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's first fetch: %v\n%s", err, out)
	}
	if _, err := milestone.Load(grace, "round-one"); err != nil {
		t.Fatalf("grace has to have the milestone before it is filed: %v", err)
	}

	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	out, err := runCLI(t, grace, "fetch")
	if err != nil {
		t.Fatalf("grace's second fetch: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one") {
		t.Errorf("the fetch does not say the milestone was filed:\n%s", out)
	}

	// The file is still hers — nothing deletes a file jaira only read — but it
	// carries the mark, and the mark is what takes it off the board.
	ms, err := milestone.Load(grace, "round-one")
	if err != nil {
		t.Fatalf("the fetch removed a file instead of marking it: %v", err)
	}
	if !ms.Filed() {
		t.Errorf("grace's file carries status %q, want %q", ms.Status, milestone.StatusFiled)
	}
	out, err = runCLI(t, grace, "milestone", "ls")
	if err != nil {
		t.Fatalf("grace's ls: %v\n%s", err, out)
	}
	if strings.Contains(out, "round-one") {
		t.Errorf("grace's 'milestone ls' still names a milestone that was filed:\n%s", out)
	}
}

// The door a tree walks into when the filing happened somewhere else: grace
// never had the file, so her disk holds nothing and her logbook holds nothing
// — the only thing that says "filed" is the ref ada left standing. The refusal
// has to point at ada's tree, because 'jaira restore' here answers that the
// file is not in the archive.
func TestAMilestoneFiledOnItsRefPointsAtTheTreeThatFiledIt(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's fetch: %v\n%s", err, out)
	}
	if _, err := milestone.Load(grace, "round-one"); err == nil {
		t.Fatal("grace has the file, so this is not the ref-only state under test")
	}

	if out, err := runCLI(t, grace, "create", "something to group", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create ticket: %v\n%s", err, out)
	}
	ticketID := handleFromList(t, grace, "something to group")
	doors := map[string][]string{
		"create": {"milestone", "create", "round-one"},
		"add":    {"milestone", "add", "round-one", ticketID},
	}
	for door, args := range doors {
		out, err := runCLI(t, grace, args...)
		if err == nil {
			t.Fatalf("'%s' went through a milestone filed on its ref:\n%s", door, out)
		}
		for _, want := range []string{
			gitref.MilestoneRefName("round-one"),
			milestone.StatusFiled,
			"jaira restore round-one.md",
			"the tree that filed it",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("the %s refusal does not mention %q: %v", door, want, err)
			}
		}
		// The restore it names runs in ada's tree. Claiming a logbook here
		// would be claiming a copy grace does not have.
		if strings.Contains(err.Error(), "into the logbook") {
			t.Errorf("the %s refusal claims a logbook copy grace never had: %v", door, err)
		}
	}
}
