package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/ticket"
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
	ms, err := milestoneLoad(t, grace, "round-one")
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
	ms, err := milestoneLoad(t, grace, "round-one")
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
// what every other clone reads, and its only copy. Left up and marked, it says
// "filed" to everyone who fetches, and it stays off their boards.
func TestAFiledMilestoneStaysOffTheOtherCloneAndKeepsItsRef(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	filedMS, err := milestoneLoad(t, ada, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if ms, err := milestoneLoad(t, ada, "round-one"); err != nil || !ms.Filed() {
		t.Errorf("the milestone is not filed on ada's board: %v", err)
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
	file, content, _, err := repo.ReadMilestone(filedMS.ID)
	if err != nil {
		t.Fatalf("the ref of a filed milestone was taken down: %v", err)
	}
	if got := milestone.FromBytes(file, content); !got.Filed() {
		t.Errorf("the ref carries status %q, want %q", got.Status, milestone.StatusFiled)
	}

	// And grace, who never saw it on her board, does not get it written there.
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if all, _ := milestone.LoadAll(refsOf(t, grace)); len(all) != 0 {
		t.Error("a fetch put a filed milestone onto grace's board")
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
	if _, err := milestoneLoad(t, grace, "round-one"); err != nil {
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
	ms, err := milestoneLoad(t, grace, "round-one")
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

// A filed milestone is refused at every door that would change it, and the
// refusal names its ref, the mark and the way back — the same in every clone,
// because the ref is.
func TestAFiledMilestoneIsRefusedWithItsWayBack(t *testing.T) {
	ada, _ := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, _ := milestoneLoad(t, ada, "round-one")
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	ticketID := aTicket(t, ada)
	doors := map[string][]string{
		"add":     {"milestone", "add", "round-one", ticketID},
		"logbook": {"logbook", "round-one"},
	}
	for door, args := range doors {
		out, err := runCLI(t, ada, args...)
		if err == nil {
			t.Fatalf("'%s' went through a filed milestone:\n%s", door, out)
		}
		for _, want := range []string{ms.Ref(), milestone.StatusFiled, "jaira restore round-one-" + ms.ID} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("the %s refusal does not mention %q: %v", door, want, err)
			}
		}
	}
}

// Fetch writes milestones too — it moves old name refs and files of an older
// build onto ID refs — and it is not a command somebody typed: maybeFetch
// spawns a detached 'jaira fetch --json' after every command, so it runs
// beside the next 'jaira milestone add'. It takes the milestone lock.
func TestFetchWaitsForTheMilestoneLock(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}

	s, err := ticket.At(grace)
	if err != nil {
		t.Fatal(err)
	}
	unlock, err := s.Lock(milestoneLockName)
	if err != nil {
		t.Fatal(err)
	}

	fetched := make(chan error, 1)
	go func() {
		_, err := runCLI(t, grace, "fetch")
		fetched <- err
	}()

	select {
	case err := <-fetched:
		unlock()
		t.Fatalf("fetch did not wait for the milestone lock (returned %v)", err)
	case <-time.After(150 * time.Millisecond):
	}
	unlock()
	select {
	case err := <-fetched:
		if err != nil {
			t.Fatalf("fetch after unlock: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("fetch never completed after the lock was released")
	}
	if _, err := milestoneLoad(t, grace, "round-one"); err != nil {
		t.Fatalf("the milestone did not arrive: %v", err)
	}
}

// An outbox report is read by whoever lost the race, and it has to name the
// thing they wrote. Put through ticket.Handle, the milestone 'next-release'
// comes out as 'elease' — six characters that name nothing, that no command
// takes back, and that an agent reading --json finds under the key "ticket".
//
// The second line matters as much: a rejected ticket keeps its local file and
// shows both sides on the board, and a rejected milestone is merged into the
// version that won (refsync.Flush), so the edit is neither dropped nor left for
// the next fetch to overwrite. Telling that user their file is unchanged sends
// them away believing something other than what happened.
func TestARejectedMilestoneIsNamedAndSaysItWasMergedNotDropped(t *testing.T) {
	ada, grace := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "next-release"); err != nil {
		t.Fatalf("ada create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace fetch: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "milestone", "set", "next-release", "--goal", "ada's"); err != nil {
		t.Fatalf("ada set: %v\n%s", err, out)
	}

	// Grace has not fetched ada's edit, so nothing locally stops her; the
	// ref has moved and the write loses the race.
	var stderr string
	_, stderr = captureStdio(t, func() {
		if out, err := runAndSend(t, grace, "milestone", "set", "next-release", "--goal", "grace's"); err != nil {
			t.Fatalf("grace set: %v\n%s", err, out)
		}
	})

	if !strings.Contains(stderr, "jaira: next-release: someone else wrote it first; your edit was merged") {
		t.Errorf("the rejection does not name the milestone and say it was merged:\n%s", stderr)
	}
	if !strings.Contains(stderr, "both changed goal: theirs is kept") {
		t.Errorf("the clash is not told:\n%s", stderr)
	}
	if strings.Contains(stderr, "your file is unchanged") {
		t.Errorf("the rejection tells a milestone writer a ticket's story:\n%s", stderr)
	}
	// the edit is neither lost nor left for the next fetch to overwrite
	if _, raw := alphaOf(t, grace, "next-release"); !strings.Contains(raw, "grace's") {
		t.Errorf("grace's goal is gone from her file:\n%s", raw)
	}
	if _, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatal(err)
	}
	if _, raw := alphaOf(t, grace, "next-release"); !strings.Contains(raw, "grace's") {
		t.Errorf("the next fetch replaced her edit:\n%s", raw)
	}
}

// The key an agent looks the subject up under, and the word a person reads,
// both come from the kind. This pins the ticket half too: it is the half that
// was already right and that the milestone branch must not change.
func TestRefSubjectNamesEachKindTheWayItIsAddressed(t *testing.T) {
	id := "01M2GGKMFB1XXK7V8AFW0YGWXQ"
	if key, name := refSubject(outbox.KindTicket, id, ""); key != "ticket" || name != ticket.Handle(id) {
		t.Errorf("ticket subject is (%q, %q), want (%q, %q)", key, name, "ticket", ticket.Handle(id))
	}
	if key, name := refSubject(outbox.KindMilestone, "K3F9QZ", "next-release-K3F9QZ.md"); key != "milestone" || name != "next-release" {
		t.Errorf("milestone subject is (%q, %q), want (%q, %q)", key, name, "milestone", "next-release")
	}
	if key, name := refSubject(outbox.KindMilestone, "next-release", ""); key != "milestone" || name != "next-release" {
		t.Errorf("an entry an older build queued is (%q, %q), want (%q, %q)", key, name, "milestone", "next-release")
	}
	if a := rejectedAdvice(outbox.KindTicket); !strings.Contains(a, "your file is unchanged") {
		t.Errorf("the ticket advice changed: %q", a)
	}
	if a := rejectedAdvice(outbox.KindMilestone); strings.Contains(a, "your file is unchanged") {
		t.Errorf("the milestone advice is still the ticket's: %q", a)
	}
}

// A rename has to reach a clone that already has the milestone: the file
// arrives under its new name on the same ref and the old file, which has the
// same ID, goes. What came after it names the ID and is not touched. Renaming
// back reaches her the same way.
func TestARenamedMilestoneArrivesUnderItsNewNameOnly(t *testing.T) {
	ada, grace := twoBoards(t)
	for _, args := range [][]string{
		{"milestone", "create", "round-one", "--after", ""},
		{"milestone", "create", "round-two", "--after", "round-one"},
	} {
		if out, err := runAndSend(t, ada, args...); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's first fetch: %v\n%s", err, out)
	}
	one, err := milestoneLoad(t, grace, "round-one")
	if err != nil {
		t.Fatalf("grace has to have the milestone before it is renamed: %v", err)
	}

	onGrace := func(want string) {
		t.Helper()
		all, err := milestone.LoadAll(refsOf(t, grace))
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, ms := range all {
			names = append(names, ms.Name)
			if ms.Name == "round-two" && strings.Join(ms.After, ",") != one.ID {
				t.Errorf("grace's round-two comes after %v, want %s", ms.After, one.ID)
			}
		}
		if strings.Join(names, ",") != want {
			t.Errorf("grace's board has %v, want %s", names, want)
		}
	}

	if out, err := runAndSend(t, ada, "milestone", "rename", "round-one", "first-round"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	out, err := runCLI(t, grace, "fetch")
	if err != nil {
		t.Fatalf("grace's second fetch: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one → first-round") || strings.Contains(out, "filed elsewhere") {
		t.Errorf("the fetch must report the rename, not a filing:\n%s", out)
	}
	onGrace("first-round,round-two")
	if out, _ := runCLI(t, grace, "milestone", "ls"); strings.Contains(out, "round-one") {
		t.Errorf("grace's 'milestone ls' still names the old name:\n%s", out)
	}

	if out, err := runAndSend(t, ada, "milestone", "rename", "first-round", "round-one"); err != nil {
		t.Fatalf("renaming back: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's third fetch: %v\n%s", err, out)
	}
	onGrace("round-one,round-two")
}

// A rename is one write on the milestone's own ref, so renaming and renaming
// back before anything was sent queues the same ref twice, and what arrives
// elsewhere is the milestone under its old name again.
func TestARenameCanBeTakenBackBeforeItsMarkIsSent(t *testing.T) {
	ada, grace := twoBoards(t)
	for _, args := range [][]string{
		{"milestone", "create", "a", "--after", ""},
		{"milestone", "create", "b", "--after", "a"},
	} {
		if out, err := runAndSend(t, ada, args...); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	first, err := milestoneLoad(t, ada, "a")
	if err != nil {
		t.Fatal(err)
	}
	// runCLI sends nothing: both renames stay in the outbox.
	if out, err := runCLI(t, ada, "milestone", "rename", "a", "alpha"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	if out, err := runCLI(t, ada, "milestone", "rename", "alpha", "a"); err != nil {
		t.Fatalf("taking the rename back before its mark was sent: %v\n%s", err, out)
	}
	flushRefs()

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's fetch: %v\n%s", err, out)
	}
	all, err := milestone.LoadAll(refsOf(t, grace))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, ms := range all {
		names = append(names, ms.Name)
	}
	if strings.Join(names, ",") != "a,b" {
		t.Errorf("grace's board has %v, want a and b", names)
	}
	if b, _ := milestoneLoad(t, grace, "b"); strings.Join(b.After, ",") != first.ID {
		t.Errorf("grace's b comes after %v, want %s", b.After, first.ID)
	}
}

// A new milestone lives on a ref named after its ID alone, carrying
// <name>-<ID>.md: the name is nowhere in the ref's name.
func TestAMilestoneIsFiledUnderItsIDAndTravelsOnTheIDRef(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "backend"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestoneLoad(t, ada, "backend")
	if err != nil {
		t.Fatal(err)
	}
	if !milestone.IsID(ms.ID) || ms.FileName() != "backend-"+ms.ID+".md" {
		t.Fatalf("the file is %s", ms.FileName())
	}
	if refs := strings.Fields(gitOut(t, ada, "for-each-ref", "--format=%(refname)", gitref.MilestonePrefix)); len(refs) != 1 || refs[0] != gitref.MilestoneRefName(ms.ID) {
		t.Fatalf("milestone refs %v, want only %s", refs, gitref.MilestoneRefName(ms.ID))
	}
	repo := &gitref.Repo{Dir: ada, Remote: "origin"}
	file, content, _, err := repo.ReadMilestone(ms.ID)
	if err != nil || file != ms.FileName() || string(content) != string(ms.Bytes()) {
		t.Errorf("the ref carries %q (%v), want %q with the file's bytes", file, err, ms.FileName())
	}
}

// A rename changes the file name on the one ref: the ID, the ref and every
// milestone that follows it stay as they are, and exactly one ref is written.
func TestARenameWritesOnlyTheFileNameOnTheSameRef(t *testing.T) {
	ada, _ := twoBoards(t)
	for _, args := range [][]string{
		{"milestone", "create", "a", "--after", ""},
		{"milestone", "create", "b", "--after", "a"},
	} {
		if out, err := runAndSend(t, ada, args...); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	a, _ := milestoneLoad(t, ada, "a")
	b, _ := milestoneLoad(t, ada, "b")
	before := gitOut(t, ada, "for-each-ref", "--format=%(refname) %(objectname)", gitref.MilestonePrefix)

	if out, err := runAndSend(t, ada, "milestone", "rename", "a", "alpha"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	after := gitOut(t, ada, "for-each-ref", "--format=%(refname) %(objectname)", gitref.MilestonePrefix)
	changed := 0
	for _, line := range strings.Split(strings.TrimSpace(after), "\n") {
		if !strings.Contains(before, line) {
			changed++
			if !strings.HasPrefix(line, gitref.MilestoneRefName(a.ID)+" ") {
				t.Errorf("the rename moved %s", line)
			}
		}
	}
	if changed != 1 || strings.Count(after, "\n") != 2 {
		t.Errorf("one ref must change and none appear:\nbefore %safter %s", before, after)
	}
	file, _, _, err := (&gitref.Repo{Dir: ada, Remote: "origin"}).ReadMilestone(a.ID)
	if err != nil || file != "alpha-"+a.ID+".md" {
		t.Errorf("the ref carries %q (%v)", file, err)
	}
	if got, _ := milestoneLoad(t, ada, "b"); strings.Join(got.After, ",") != a.ID || got.FileName() != b.FileName() {
		t.Errorf("b was touched: after %v", got.After)
	}
}

// The first milestone command on a board that still has files of an older
// build — before IDs, or before milestones lived on refs only — takes them all
// onto their refs and deletes them, a read included: from then on nothing
// reads that directory.
func TestTheFirstCommandTakesEveryOldFileOntoItsRef(t *testing.T) {
	dir := milestoneStore(t)
	oldMilestoneFiles(t, dir, map[string]string{
		"q4.md":        "---\ncolor: 45\nafter: []\n---\n",
		"q5.md":        "---\ncolor: 46\nafter: [q4]\n---\n",
		"q6-K3F9QZ.md": "---\ncolor: 47\n---\n",
	})
	if out, err := runCLI(t, dir, "milestone", "ls"); err != nil || !strings.Contains(out, "q4") || !strings.Contains(out, "q6") {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if _, err := os.Stat(milestone.Dir(dir)); !os.IsNotExist(err) {
		t.Errorf("the old directory is still there: %v", err)
	}
	if got := milestoneDirNames(t, dir); got != "q4-"+milestone.LegacyID("q4")+".md,q5-"+milestone.LegacyID("q5")+".md,q6-K3F9QZ.md" {
		t.Errorf("refs carry %s", got)
	}
	if q5, _ := milestoneLoad(t, dir, "q5"); strings.Join(q5.After, ",") != milestone.LegacyID("q4") {
		t.Errorf("q5 comes after %v, want q4's ID", q5.After)
	}
}

func oldMilestoneFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	if err := os.MkdirAll(milestone.Dir(dir), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(milestone.Dir(dir), name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// Taking the files off the disk is a write nobody asked for by name, so it is
// told — once, and not again.
func TestTheMoveOntoRefsIsToldOnce(t *testing.T) {
	dir := milestoneStore(t)
	oldMilestoneFiles(t, dir, map[string]string{"q4.md": "---\ncolor: 45\n---\n", "q5.md": "---\ncolor: 46\n---\n"})

	_, stderr := captureStdio(t, func() {
		if out, err := runCLI(t, dir, "milestone", "create", "q6"); err != nil {
			t.Fatalf("create: %v\n%s", err, out)
		}
	})
	for _, want := range []string{"took 2 file(s) off the disk", "q4.md -> refs/jaira/milestones/" + milestone.LegacyID("q4"), "q5.md -> refs/jaira/milestones/" + milestone.LegacyID("q5")} {
		if !strings.Contains(stderr, want) {
			t.Errorf("the message lacks %q:\n%s", want, stderr)
		}
	}
	_, stderr = captureStdio(t, func() {
		if out, err := runCLI(t, dir, "milestone", "create", "q7"); err != nil {
			t.Fatalf("second create: %v\n%s", err, out)
		}
	})
	if strings.Contains(stderr, "took") {
		t.Errorf("a second write told the move again:\n%s", stderr)
	}
}

// Filing marks the ref and writes no file; restore takes the mark off by the
// <name>-<ID> the logbook printed. An OLD <name>.md an older build filed into
// a logbook goes onto its ID ref as <name>-<LegacyID>.md, with the move told.
func TestLogbookAndRestoreSpeakIDFileNamesAndMigrateAnOldOne(t *testing.T) {
	ada, _ := twoBoards(t)
	done := mkDoneTicket(t, ada, "finished")
	if out, err := runAndSend(t, ada, "milestone", "create", "round-one", done); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, _ := milestoneLoad(t, ada, "round-one")
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	} else if !strings.Contains(out, "jaira restore round-one-"+ms.ID) {
		t.Errorf("logbook names no way back round-one-%s:\n%s", ms.ID, out)
	}
	if logbook := filepath.Join(ada, ticket.DirName, ticket.LogbookSubdir); milestoneFilesUnder(t, logbook) != "" {
		t.Errorf("filing wrote a file: %s", milestoneFilesUnder(t, logbook))
	}
	if out, err := runAndSend(t, ada, "restore", "round-one-"+ms.ID); err != nil {
		t.Fatalf("restore: %v\n%s", err, out)
	}
	if back, err := milestoneLoad(t, ada, "round-one"); err != nil || back.ID != ms.ID || back.Filed() {
		t.Fatalf("restore: %v %+v", err, back)
	}

	// An old filed file in the logbook, as an old build left it.
	folder := filepath.Join(ada, ticket.DirName, ticket.LogbookSubdir, "ada-20260101", ticket.MilestonesSubdir)
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "old-one.md"), []byte("---\ncolor: 45\nstatus: filed\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, stderr := captureStdio(t, func() {
		if out, err := runAndSend(t, ada, "restore", "old-one.md"); err != nil {
			t.Fatalf("restore old: %v\n%s", err, out)
		}
	})
	id := milestone.LegacyID("old-one")
	got, err := milestoneLoad(t, ada, "old-one")
	if err != nil || got.ID != id || got.FileName() != "old-one-"+id+".md" || got.Filed() {
		t.Fatalf("restored old file: %v %+v", err, got)
	}
	if !strings.Contains(stderr, "old-one.md -> refs/jaira/milestones/"+id) {
		t.Errorf("the move was not told:\n%s", stderr)
	}
	file, _, _, err := (&gitref.Repo{Dir: ada, Remote: "origin"}).ReadMilestone(id)
	if err != nil || file != "old-one-"+id+".md" {
		t.Errorf("the ID ref carries %q (%v)", file, err)
	}
}

// aTicket files a ticket on a two-board clone and returns its handle.
func aTicket(t *testing.T, dir string) string {
	t.Helper()
	if out, err := runCLI(t, dir, "create", "work to group", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create ticket: %v\n%s", err, out)
	}
	return handleFromList(t, dir, "work to group")
}

// milestoneDirNames is the file names the board's milestones are carried
// under on their refs, sorted.
func milestoneDirNames(t *testing.T, dir string) string {
	t.Helper()
	files, err := refsOf(t, dir).MilestoneFiles()
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

// milestoneFilesUnder lists every .md file under a milestones/ folder below
// dir.
func milestoneFilesUnder(t *testing.T, dir string) string {
	t.Helper()
	var found []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && filepath.Base(filepath.Dir(p)) == ticket.MilestonesSubdir {
			found = append(found, p)
		}
		return nil
	})
	return strings.Join(found, ",")
}

// DoD 2: a rename and a create of the same name, made at once on two clones
// that have not fetched, leave both milestones standing, in either order: they
// are two refs with two IDs, and nothing needs to be settled between them.
func TestARenameAndACreateOfTheSameNameAtOnceLoseNothing(t *testing.T) {
	for _, renameFirst := range []bool{true, false} {
		ada, grace := twoBoards(t)
		if out, err := runAndSend(t, ada, "milestone", "create", "backend"); err != nil {
			t.Fatalf("create: %v\n%s", err, out)
		}
		if out, err := runCLI(t, grace, "fetch"); err != nil {
			t.Fatalf("fetch: %v\n%s", err, out)
		}
		backend, err := milestoneLoad(t, ada, "backend")
		if err != nil {
			t.Fatal(err)
		}
		rename := func() {
			if out, err := runAndSend(t, ada, "milestone", "rename", "backend", "frontend"); err != nil {
				t.Fatalf("rename: %v\n%s", err, out)
			}
		}
		create := func() {
			if out, err := runAndSend(t, grace, "milestone", "create", "frontend"); err != nil {
				t.Fatalf("create: %v\n%s", err, out)
			}
		}
		if renameFirst {
			rename()
			create()
		} else {
			create()
			rename()
		}
		for _, dir := range []string{ada, grace} {
			if out, err := runCLI(t, dir, "fetch"); err != nil {
				t.Fatalf("fetch: %v\n%s", err, out)
			}
		}

		var ids [2]string
		for i, dir := range []string{ada, grace} {
			j := jsonCLI(t, dir, "milestone", "ls")
			rows, _ := j["milestones"].([]any)
			if len(rows) != 2 {
				t.Fatalf("renameFirst=%v: %s has %d milestones, want both: %v", renameFirst, dir, len(rows), rows)
			}
			a, b := rows[0].(map[string]any), rows[1].(map[string]any)
			if a["name"] != "frontend" || b["name"] != "frontend" || a["id"] == b["id"] {
				t.Fatalf("renameFirst=%v: %s sees %v", renameFirst, dir, rows)
			}
			if a["id"].(string) > b["id"].(string) {
				a, b = b, a
			}
			ids[i] = a["id"].(string) + "," + b["id"].(string)
		}
		if ids[0] != ids[1] {
			t.Errorf("the two boards disagree on the IDs: %v", ids)
		}
		for _, dir := range []string{ada, grace} {
			if names := milestoneDirNames(t, dir); strings.Count(names, "frontend-") != 2 || strings.Contains(names, "backend") {
				t.Errorf("renameFirst=%v: %s holds %s", renameFirst, dir, names)
			}
		}
		file, _, _, err := (&gitref.Repo{Dir: ada, Remote: "origin"}).ReadMilestone(backend.ID)
		if err != nil || file != "frontend-"+backend.ID+".md" {
			t.Errorf("the renamed milestone's ref carries %q (%v)", file, err)
		}
		task := aTicket(t, ada)
		_, err = runCLI(t, ada, "milestone", "add", "frontend", task)
		if err == nil || exitCode(err) != ExitValidation || strings.Count(err.Error(), "frontend-") < 2 {
			t.Errorf("the name is ambiguous now and says both forms: %v", err)
		}
		if out, err := runCLI(t, ada, "milestone", "add", "frontend-"+backend.ID, task); err != nil {
			t.Errorf("the name-ID form is the way out: %v\n%s", err, out)
		}
	}
}

// Boards from before IDs, and a remote that has only name refs, move to IDs
// without anybody speaking to anybody: an old file that is a copy of what a
// name ref carries is deleted, the first fetch moves the name refs to ID refs,
// the old name refs stay on the remote untouched, the other clone ends with
// the same bytes and no double, and a name ref that turns up later is taken
// once and then ignored.
func TestOldBoardsAndOldRefsMoveToIDsWithoutAgreement(t *testing.T) {
	ada, grace := twoBoards(t)
	remote := &gitref.Repo{Dir: ada, Remote: "origin"}
	q4 := "---\ncolor: 45\nafter: []\n---\n\n# q4\n"
	q5 := "---\ncolor: 46\nafter: [q4]\n---\n\n# q5\n"
	for name, body := range map[string]string{"q4": q4, "q5": q5} {
		if _, err := remote.WriteMilestone(name, name+".md", []byte(body), ""); err != nil {
			t.Fatal(err)
		}
	}
	nameRefs := func() string {
		return gitOut(t, ada, "ls-remote", "origin", gitref.MilestonePrefix+"q*")
	}
	before := nameRefs()

	// grace has the old files on disk: copies of what the name refs carry
	oldMilestoneFiles(t, grace, map[string]string{"q4.md": q4, "q5.md": q5})
	j := jsonCLI(t, grace, "milestone", "ls")
	rows, _ := j["milestones"].([]any)
	if len(rows) != 2 || rows[0].(map[string]any)["id"] != milestone.LegacyID("q4") {
		t.Fatalf("grace's ls: %v", rows)
	}
	if _, err := os.Stat(milestone.Dir(grace)); !os.IsNotExist(err) {
		t.Fatalf("the old files are still there: %v", err)
	}

	// ada's fetch moves the name refs to ID refs and sends them
	if out, err := runAndSend(t, ada, "fetch"); err != nil {
		t.Fatalf("ada's fetch: %v\n%s", err, out)
	}
	id4, id5 := milestone.LegacyID("q4"), milestone.LegacyID("q5")
	if want := "q4-" + id4 + ".md,q5-" + id5 + ".md"; milestoneDirNames(t, ada) != want {
		t.Fatalf("ada holds %s, want %s", milestoneDirNames(t, ada), want)
	}
	for _, id := range []string{id4, id5} {
		if _, _, _, err := remote.ReadMilestone(id); err != nil {
			t.Errorf("the ID ref %s was not pushed: %v", id, err)
		}
	}
	if nameRefs() != before {
		t.Errorf("the old name refs changed:\nbefore %safter %s", before, nameRefs())
	}

	// grace fetches: same bytes as ada, no double
	if out, err := runAndSend(t, grace, "fetch"); err != nil {
		t.Fatalf("grace's fetch: %v\n%s", err, out)
	}
	if got := milestoneDirNames(t, grace); got != milestoneDirNames(t, ada) {
		t.Errorf("grace holds %s, ada %s", got, milestoneDirNames(t, ada))
	}
	for _, id := range []string{id4, id5} {
		a, _ := milestoneLoad(t, ada, id)
		g, _ := milestoneLoad(t, grace, id)
		if string(a.Bytes()) != string(g.Bytes()) {
			t.Errorf("%s differs between the clones:\n%s\n%s", id, a.Bytes(), g.Bytes())
		}
	}

	// an old name ref turning up later is taken once, then ignored
	if _, err := remote.WriteMilestone("q6", "q6.md", []byte("---\ncolor: 47\n---\n"), ""); err != nil {
		t.Fatal(err)
	}
	id6 := milestone.LegacyID("q6")
	if out, err := runAndSend(t, grace, "fetch", "--json"); err != nil || !strings.Contains(out, "milestones_migrated") || !strings.Contains(out, gitref.MilestoneRefName(id6)) {
		t.Fatalf("grace's third fetch: %v\n%s", err, out)
	}
	if got := milestoneDirNames(t, grace); !strings.Contains(got, "q6-"+id6+".md") {
		t.Errorf("q6 was not taken and moved: %s", got)
	}
	if out, err := runCLI(t, grace, "fetch", "--json"); err != nil || strings.Contains(out, gitref.MilestoneRefName(id6)) {
		t.Fatalf("grace's fourth fetch moved q6 again: %v\n%s", err, out)
	}
	if got := milestoneDirNames(t, grace); strings.Count(got, "q6-") != 1 {
		t.Errorf("q6 doubled: %s", got)
	}
}

// An old <name>.md in the logbook restored beside a live milestone of the same
// name must bring back that file's milestone, not report the live one.
func TestRestoreOldFileBesideALiveTwinUnfilesTheRestoredOne(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "old-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	folder := filepath.Join(ada, ticket.DirName, ticket.LogbookSubdir, "ada-20260101", ticket.MilestonesSubdir)
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "old-one.md"), []byte("---\ncolor: 45\nstatus: filed\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runAndSend(t, ada, "restore", "old-one.md"); err != nil {
		t.Fatalf("restore: %v\n%s", err, out)
	}
	got, err := milestoneLoad(t, ada, milestone.LegacyID("old-one"))
	if err != nil {
		t.Fatalf("load restored: %v", err)
	}
	if got.Filed() {
		t.Errorf("the restored milestone is still filed")
	}
}

// Two milestones sharing a name: create's refusal names a ref that is there,
// and the <name>-<ID> forms 'add' takes.
func TestCreateOnATwinnedNameNamesARealRef(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "dup"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	twin := "dup-" + milestone.LegacyID("dup") + ".md"
	oldMilestoneFiles(t, ada, map[string]string{twin: "---\ncolor: 45\n---\n"})
	out, err := runAndSend(t, ada, "milestone", "create", "dup")
	if err == nil {
		t.Fatalf("create on a taken name succeeded:\n%s", out)
	}
	// The hint must be a form 'add' takes, not the name two milestones share.
	if !strings.Contains(out+err.Error(), "milestone add <dup-") {
		t.Errorf("hint does not spell the <name>-<ID> forms:\n%s%v", out, err)
	}
	all, _ := milestone.LoadAll(refsOf(t, ada))
	found := false
	for _, m := range all {
		if strings.Contains(out+err.Error(), m.Ref()) {
			found = true
		}
	}
	if len(all) != 2 || !found {
		t.Errorf("refusal names no ref that exists (%d milestones):\n%s%v", len(all), out, err)
	}
}

func alphaOf(t *testing.T, dir, name string) (*milestone.Milestone, string) {
	t.Helper()
	ms, err := milestoneLoad(t, dir, name)
	if err != nil {
		t.Fatal(err)
	}
	return ms, string(ms.Bytes())
}
