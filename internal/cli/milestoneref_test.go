package cli

import (
	"os"
	"os/exec"
	"path/filepath"
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
	filedMS, err := milestone.Load(ada, "round-one")
	if err != nil {
		t.Fatal(err)
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

// The tree that filed it stands in a state no other tree is in: the copy is in
// its own logbook AND the mark is on the ref it wrote. It must be told about
// the copy lying here, not sent to "the tree that filed it" — which is the
// tree it is standing in. Asking the ref first said exactly that, so this is
// the test that holds the order down. It runs on a board with refs, because a
// board without them cannot get the order wrong.
func TestTheFilingTreeIsPointedAtItsOwnLogbook(t *testing.T) {
	ada, _ := twoBoards(t)

	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if out, err := runCLI(t, ada, "create", "something to group", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create ticket: %v\n%s", err, out)
	}
	ticketID := handleFromList(t, ada, "something to group")

	doors := map[string][]string{
		"add": {"milestone", "add", "round-one", ticketID},
		// Filing the same milestone twice: the copy is in this tree's
		// logbook, so this door owes the reader the restore that runs here
		// like the other two do.
		"logbook": {"logbook", "round-one"},
	}
	for door, args := range doors {
		out, err := runCLI(t, ada, args...)
		if err == nil {
			t.Fatalf("'%s' went through a milestone this tree filed:\n%s", door, out)
		}
		for _, want := range []string{
			"into the logbook",
			"jaira restore round-one-",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("the %s refusal does not mention %q: %v", door, want, err)
			}
		}
		// The ref carries the mark here as well, and naming it would send the
		// reader away from the copy that is lying in this very tree.
		if strings.Contains(err.Error(), gitref.MilestoneRefName("round-one")) {
			t.Errorf("the %s refusal points at the ref although the copy is here: %v", door, err)
		}
		if strings.Contains(err.Error(), "the tree that filed it") {
			t.Errorf("the %s refusal sends the filing tree to itself: %v", door, err)
		}
	}
}

// Fetch is the fifth writer of a milestone file, after create, add/rm, logbook
// and restore, and the only one that is not a command somebody typed:
// maybeFetch spawns a detached 'jaira fetch --json' after every command, so it
// runs beside the next 'jaira milestone add'. Unlocked, one of the two writes
// to the same file is lost.
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
	// And it has written nothing while waiting: a fetch that blocked only
	// after putting the file down would pass the wait above and still lose a
	// concurrent add's write.
	if _, err := os.Stat(milestonePath(grace, "round-one")); !os.IsNotExist(err) {
		unlock()
		t.Fatalf("the fetch wrote the milestone file while the lock was held (stat: %v)", err)
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
	if _, err := milestone.Load(grace, "round-one"); err != nil {
		t.Fatalf("the fetch did not write the milestone once it had the lock: %v", err)
	}
}

// An outbox report is read by whoever lost the race, and it has to name the
// thing they wrote. Put through ticket.Handle, the milestone 'next-release'
// comes out as 'elease' — six characters that name nothing, that no command
// takes back, and that an agent reading --json finds under the key "ticket".
//
// The second line matters as much: a rejected ticket keeps its local file and
// shows both sides on the board, but a rejected milestone is overwritten by
// the next fetch (refsync.IncomingMilestones), so the edit has to be made
// again on the fetched file. Telling that user their file is unchanged sends
// them away believing work is safe that is about to be discarded.
func TestARejectedMilestoneIsNamedAndSaysTheFetchWillReplaceIt(t *testing.T) {
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

	if !strings.Contains(stderr, "jaira: next-release was not sent") {
		t.Errorf("the rejection does not name the milestone:\n%s", stderr)
	}
	if !strings.Contains(stderr, "the next fetch replaces your file") {
		t.Errorf("the rejection does not say the fetch will replace the file:\n%s", stderr)
	}
	if strings.Contains(stderr, "your file is unchanged") {
		t.Errorf("the rejection tells a milestone writer a ticket's story:\n%s", stderr)
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
	one, err := milestone.Load(grace, "round-one")
	if err != nil {
		t.Fatalf("grace has to have the milestone before it is renamed: %v", err)
	}

	onGrace := func(want string) {
		t.Helper()
		all, err := milestone.LoadAll(grace)
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
	first, err := milestone.Load(ada, "a")
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
	all, err := milestone.LoadAll(grace)
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
	if b, _ := milestone.Load(grace, "b"); strings.Join(b.After, ",") != first.ID {
		t.Errorf("grace's b comes after %v, want %s", b.After, first.ID)
	}
}

// A new milestone lives in <name>-<ID>.md and travels on a ref named after its
// ID alone, carrying that file: the name is nowhere in the ref's name.
func TestAMilestoneIsFiledUnderItsIDAndTravelsOnTheIDRef(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "backend"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(ada, "backend")
	if err != nil {
		t.Fatal(err)
	}
	if !milestone.IsID(ms.ID) || filepath.Base(ms.Path(ada)) != "backend-"+ms.ID+".md" {
		t.Fatalf("the file is %s", ms.Path(ada))
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
	a, _ := milestone.Load(ada, "a")
	b, _ := milestone.Load(ada, "b")
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
	if got, _ := milestone.Load(ada, "b"); strings.Join(got.After, ",") != a.ID || got.FileName() != b.FileName() {
		t.Errorf("b was touched: after %v", got.After)
	}
}

// The first write on a board that still has files from before IDs moves them
// all, under the lock, whatever it was asked to write.
func TestTheFirstWriteMovesEveryOldFileToItsID(t *testing.T) {
	dir := emptyStore(t)
	dirMS := milestone.Dir(dir)
	if err := os.MkdirAll(dirMS, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"q4.md": "---\ncolor: 45\nafter: []\n---\n",
		"q5.md": "---\ncolor: 46\nafter: [q4]\n---\n",
	} {
		if err := os.WriteFile(filepath.Join(dirMS, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if out, err := runCLI(t, dir, "milestone", "ls"); err != nil || !strings.Contains(out, "q4") || !strings.Contains(out, "q5") {
		t.Fatalf("reading must work on old files: %v\n%s", err, out)
	}
	if entries, _ := os.ReadDir(dirMS); len(entries) != 2 || entries[0].Name() != "q4.md" {
		t.Fatalf("reading changed the directory: %v", entries)
	}
	if out, err := runCLI(t, dir, "milestone", "create", "q6"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	var names []string
	entries, _ := os.ReadDir(dirMS)
	for _, e := range entries {
		names = append(names, e.Name())
	}
	q6, _ := milestone.Load(dir, "q6")
	want := "q4-" + milestone.LegacyID("q4") + ".md,q5-" + milestone.LegacyID("q5") + ".md,q6-" + q6.ID + ".md"
	if strings.Join(names, ",") != want {
		t.Errorf("directory %v, want %s", names, want)
	}
	if q5, _ := milestone.Load(dir, "q5"); strings.Join(q5.After, ",") != milestone.LegacyID("q4") {
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

// The move is a write nobody asked for by name, so it is told — once, on the
// first command that writes, and not again.
func TestTheMoveToIDNamesIsToldOnTheFirstWriteOnly(t *testing.T) {
	dir := emptyStore(t)
	oldMilestoneFiles(t, dir, map[string]string{"q4.md": "---\ncolor: 45\n---\n", "q5.md": "---\ncolor: 46\n---\n"})

	_, stderr := captureStdio(t, func() {
		if out, err := runCLI(t, dir, "milestone", "create", "q6"); err != nil {
			t.Fatalf("create: %v\n%s", err, out)
		}
	})
	for _, want := range []string{"moved 2 milestone file(s) to ID names", "q4.md -> q4-" + milestone.LegacyID("q4") + ".md", "q5.md -> q5-" + milestone.LegacyID("q5") + ".md"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("the message lacks %q:\n%s", want, stderr)
		}
	}
	_, stderr = captureStdio(t, func() {
		if out, err := runCLI(t, dir, "milestone", "create", "q7"); err != nil {
			t.Fatalf("second create: %v\n%s", err, out)
		}
	})
	if strings.Contains(stderr, "moved") {
		t.Errorf("a second write told the move again:\n%s", stderr)
	}
}

// Filing and restoring speak file names: the logbook keeps <name>-<ID>.md and
// restore takes it back to the board and to its ref; an OLD <name>.md found in
// a logbook lands as <name>-<LegacyID>.md, on its ID ref, with the move told.
func TestLogbookAndRestoreSpeakIDFileNamesAndMigrateAnOldOne(t *testing.T) {
	ada, _ := twoBoards(t)
	done := mkDoneTicket(t, ada, "finished")
	if out, err := runAndSend(t, ada, "milestone", "create", "round-one", done); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, _ := milestone.Load(ada, "round-one")
	if out, err := runAndSend(t, ada, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	} else if !strings.Contains(out, "jaira restore "+ms.FileName()) {
		t.Errorf("logbook names no restorable file %s:\n%s", ms.FileName(), out)
	}
	if out, err := runAndSend(t, ada, "restore", ms.FileName()); err != nil {
		t.Fatalf("restore: %v\n%s", err, out)
	}
	if back, err := milestone.Load(ada, "round-one"); err != nil || back.ID != ms.ID || back.Filed() {
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
	got, err := milestone.Load(ada, "old-one")
	if err != nil || got.ID != id || got.FileName() != "old-one-"+id+".md" || got.Filed() {
		t.Fatalf("restored old file: %v %+v", err, got)
	}
	if !strings.Contains(stderr, "old-one.md -> old-one-"+id+".md") {
		t.Errorf("the move was not told:\n%s", stderr)
	}
	file, _, _, err := (&gitref.Repo{Dir: ada, Remote: "origin"}).ReadMilestone(id)
	if err != nil || file != "old-one-"+id+".md" {
		t.Errorf("the ID ref carries %q (%v)", file, err)
	}
}
