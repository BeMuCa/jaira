package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// mkTicket captures one ticket and returns its handle, the way a person would
// name it on the command line.
func mkTicket(t *testing.T, dir, title string) string {
	t.Helper()
	if out, err := runCLI(t, dir, "create", title, "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create %q: %v\n%s", title, err, out)
	}
	return firstTitledHandle(t, dir, title)
}

// Grouping twenty tickets has to be one edit of one file — that is the whole
// reason a milestone is a file and not a field on each ticket.
func TestOneAddCallWritesOneFileForEveryTicket(t *testing.T) {
	dir := emptyStore(t)
	a := mkTicket(t, dir, "first")
	b := mkTicket(t, dir, "second")
	c := mkTicket(t, dir, "third")

	if out, err := runCLI(t, dir, "milestone", "create", "Round One"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "milestone", "add", "round-one", a, b, c)
	if err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}
	if !strings.Contains(out, "holds 3 ticket(s)") {
		t.Errorf("add did not report all three:\n%s", out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(ms.Members()); got != 3 {
		t.Errorf("milestone holds %d tickets, want 3", got)
	}
	// And the file is the only thing that changed: no ticket file learned
	// about the milestone, because that is the cost the file exists to avoid.
	tk := ticketByHandle(t, dir, a)
	body, err := os.ReadFile(tk.Path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "round-one") {
		t.Error("the ticket file mentions the milestone; the grouping must live only in the milestone file")
	}
}

// A name lands in one stored form and the person is told, or a board grows two
// names for one round of work — the failure 'jaira tags' exists to prevent,
// happening again one directory over.
func TestCreateNormalizesTheNameAndSaysSo(t *testing.T) {
	dir := emptyStore(t)
	out, err := runCLI(t, dir, "milestone", "create", "Round One")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if !strings.Contains(out, `Filed "Round One" as "round-one"`) {
		t.Errorf("create did not say what it filed the name as:\n%s", out)
	}
	if _, err := os.Stat(milestone.Path(dir, "round-one")); err != nil {
		t.Errorf("no file at %s: %v", milestone.Path(dir, "round-one"), err)
	}
	// Creating the same milestone twice is a mistake to report, not a silent
	// overwrite of a file somebody has been editing.
	if out, err := runCLI(t, dir, "milestone", "create", "round-one"); err == nil {
		t.Errorf("creating an existing milestone succeeded:\n%s", out)
	}
	// And adding to one that does not exist says how to make it, rather than
	// creating it behind your back.
	out, err = runCLI(t, dir, "milestone", "add", "round-nine", mkTicket(t, dir, "orphan"))
	if err == nil {
		t.Fatalf("add to a missing milestone succeeded:\n%s", out)
	}
	if !strings.Contains(err.Error(), "jaira milestone create round-nine") {
		t.Errorf("the refusal does not say how to make it: %v", err)
	}
}

// Nobody picks the colour, and the colour is in the milestone's own file —
// there is no second registry to keep in step with it.
func TestCreateAssignsAColourIntoTheMilestoneFile(t *testing.T) {
	dir := emptyStore(t)
	for _, name := range []string{"one", "two"} {
		if out, err := runCLI(t, dir, "milestone", "create", name); err != nil {
			t.Fatalf("create %s: %v\n%s", name, err, out)
		}
	}
	all, err := milestone.LoadAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("board has %d milestones, want 2", len(all))
	}
	for _, ms := range all {
		if ms.Colour <= 0 || ms.Colour > 255 {
			t.Errorf("milestone %q got colour %d, want an assigned ANSI-256 value", ms.Name, ms.Colour)
		}
	}
	if all[0].Colour == all[1].Colour {
		t.Errorf("both milestones got colour %d; a free palette must not repeat", all[0].Colour)
	}
	// An explicit colour is still available, for the one case where the random
	// choice is wrong.
	if out, err := runCLI(t, dir, "milestone", "create", "three", "--color", "40"); err != nil {
		t.Fatalf("create with --color: %v\n%s", err, out)
	}
	ms, err := milestone.Load(dir, "three")
	if err != nil {
		t.Fatal(err)
	}
	if ms.Colour != 40 {
		t.Errorf("--color 40 stored %d", ms.Colour)
	}
	if out, err := runCLI(t, dir, "milestone", "create", "four", "--color", "999"); err == nil {
		t.Errorf("--color 999 was accepted:\n%s", out)
	}
}

// --milestone is how an agent narrows the board to one round of work, and it
// has to answer exactly the same question the board's own gesture does.
func TestListFiltersToAMilestoneExactly(t *testing.T) {
	dir := emptyStore(t)
	in := mkTicket(t, dir, "inside the round")
	mkTicket(t, dir, "outside the round")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", in); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}

	got := jsonCLI(t, dir, "list", "--milestone", "round-one")
	rows, _ := got["tickets"].([]any)
	if len(rows) != 1 {
		t.Fatalf("--milestone round-one listed %d tickets, want 1", len(rows))
	}
	row, _ := rows[0].(map[string]any)
	if row["title"] != "inside the round" {
		t.Errorf("listed %v, want the ticket in the milestone", row["title"])
	}
	names, _ := row["milestones"].([]any)
	if len(names) != 1 || names[0] != "round-one" {
		t.Errorf("row's milestones = %v, want [round-one]", names)
	}

	// A prefix is not a match: a milestone is a name from a closed set.
	got = jsonCLI(t, dir, "list", "--milestone", "round")
	if rows, _ := got["tickets"].([]any); len(rows) != 0 {
		t.Errorf("--milestone round matched %d tickets by prefix", len(rows))
	}
	// And a ticket in no milestone reports an empty list rather than null, so
	// an agent can iterate the field without a nil check.
	got = jsonCLI(t, dir, "list")
	for _, r := range got["tickets"].([]any) {
		row, _ := r.(map[string]any)
		if row["milestones"] == nil {
			t.Errorf("ticket %v has a null milestones field", row["title"])
		}
	}
}

// rm takes a line out and leaves the milestone standing: an empty milestone is
// still a plan.
func TestRmDropsTheLineAndKeepsTheMilestone(t *testing.T) {
	dir := emptyStore(t)
	h := mkTicket(t, dir, "only one")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", h); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "milestone", "rm", "round-one", h); err != nil {
		t.Fatalf("rm: %v\n%s", err, out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatalf("rm removed the milestone itself: %v", err)
	}
	if len(ms.Members()) != 0 {
		t.Errorf("milestone still holds %v", ms.Members())
	}
	// Removing what is not there is a no-op that says so, not an error.
	out, err := runCLI(t, dir, "milestone", "rm", "round-one", h)
	if err != nil {
		t.Fatalf("second rm: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Nothing to do") {
		t.Errorf("a no-op rm did not say so:\n%s", out)
	}
	// And 'ls' names it exactly as before: a milestone goes away on command
	// and never on its own, or a freshly created one would vanish the moment
	// you took its last ticket back out.
	out, err = runCLI(t, dir, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one") {
		t.Errorf("ls stopped naming the emptied milestone:\n%s", out)
	}
}

// The rule hangs on emptying a milestone, not on being empty: creating one
// with no tickets is exactly how a plan starts, and that file has to stay put.
func TestCreateWithNoTicketsLeavesTheFileLyingThere(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatalf("an empty milestone did not survive its own creation: %v", err)
	}
	if len(ms.Members()) != 0 {
		t.Errorf("milestone holds %v, want nothing", ms.Members())
	}
	out, err := runCLI(t, dir, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one") || !strings.Contains(out, "0 ticket(s)") {
		t.Errorf("ls does not name the empty milestone:\n%s", out)
	}
}

// ls is read before creating, the way 'jaira tags' is read before tagging.
func TestLsNamesEveryMilestoneWithItsCount(t *testing.T) {
	dir := emptyStore(t)
	out, err := runCLI(t, dir, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls on an empty board: %v\n%s", err, out)
	}
	if !strings.Contains(out, "jaira milestone create") {
		t.Errorf("ls on an empty board does not say how to start one:\n%s", out)
	}
	h := mkTicket(t, dir, "one thing")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", h); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err = runCLI(t, dir, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "round-one") || !strings.Contains(out, "1 ticket(s)") {
		t.Errorf("ls does not show the milestone and its count:\n%s", out)
	}
	j := jsonCLI(t, dir, "milestone", "ls")
	arr, _ := j["milestones"].([]any)
	if len(arr) != 1 {
		t.Fatalf("ls --json listed %d milestones, want 1", len(arr))
	}
	row, _ := arr[0].(map[string]any)
	if row["name"] != "round-one" || row["count"].(float64) != 1 {
		t.Errorf("ls --json row = %v", row)
	}
}

// A line pointing at no ticket is a dead line that looks like a plan, so the
// whole call is refused rather than half-written.
func TestAddRefusesAnIDThatNamesNothing(t *testing.T) {
	dir := emptyStore(t)
	h := mkTicket(t, dir, "real one")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "milestone", "add", "round-one", h, "ZZZZZZ")
	if err == nil {
		t.Fatalf("add accepted an unknown id:\n%s", out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if len(ms.Members()) != 0 {
		t.Errorf("the refused call still wrote %v", ms.Members())
	}
}

func ticketByHandle(t *testing.T, dir, handle string) *ticket.Ticket {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(handle)
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

// mkDoneTicket puts a ticket straight into the terminal lane. It bypasses the
// pipeline deliberately: what is under test is filing a milestone, not how a
// ticket gets to be finished.
func mkDoneTicket(t *testing.T, dir, title string) string {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	id := ticket.NewID(time.Now())
	if _, err := s.Create(map[string]string{
		ticket.FieldID:     id,
		ticket.FieldTitle:  title,
		ticket.FieldStatus: "done",
	}, nil, ""); err != nil {
		t.Fatal(err)
	}
	return id
}

// A milestone is filed the way a ticket is, and comes back the way a ticket
// does — with its ticket list and its colour intact, because a group restored
// as a different colour is a group nobody recognises on the board.
func TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack(t *testing.T) {
	dir := emptyStore(t)
	done := mkDoneTicket(t, dir, "finished work")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", done); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	before, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}

	// Unfinished work in it refuses the whole thing: filing a group whose work
	// is still on the board takes the plan away and leaves the work.
	open := mkTicket(t, dir, "not finished")
	if out, err := runCLI(t, dir, "milestone", "add", "round-one", open); err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "logbook", "round-one")
	if err == nil {
		t.Fatalf("filing a milestone with unfinished work succeeded:\n%s", out)
	}
	if !strings.Contains(err.Error(), "not finished") && !strings.Contains(err.Error(), "has not reached") {
		t.Errorf("the refusal does not name what is unfinished: %v", err)
	}
	if out, err := runCLI(t, dir, "milestone", "rm", "round-one", open); err != nil {
		t.Fatalf("rm: %v\n%s", err, out)
	}

	if out, err := runCLI(t, dir, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if _, err := milestone.Load(dir, "round-one"); err == nil {
		t.Error("the milestone is still on the board after being filed")
	}
	if out, err := runCLI(t, dir, "milestone", "ls"); err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	} else if strings.Contains(out, "round-one") {
		t.Errorf("'milestone ls' still names a filed milestone:\n%s", out)
	}
	// No card carries its colour any more: the index is built from the files
	// on the board, and the file is not on it.
	rows, _ := jsonCLI(t, dir, "list", "--lane", "done")["tickets"].([]any)
	if len(rows) == 0 {
		t.Fatal("the filed milestone's ticket is not listed, so the colour check would prove nothing")
	}
	for _, r := range rows {
		row, _ := r.(map[string]any)
		if names, _ := row["milestones"].([]any); len(names) != 0 {
			t.Errorf("ticket %v still carries %v", row["title"], names)
		}
	}
	// And the logbook says where it went, or the file is there and the list
	// denies it.
	out, err = runCLI(t, dir, "logbook")
	if err != nil {
		t.Fatalf("logbook listing: %v\n%s", err, out)
	}
	if !strings.Contains(out, filepath.Join("milestones", "round-one.md")) {
		t.Errorf("the logbook listing does not name the filed milestone:\n%s", out)
	}

	// The name is taken until somebody gives it up.
	out, err = runCLI(t, dir, "milestone", "create", "round-one")
	if err == nil {
		t.Fatalf("a filed milestone's name was handed out again:\n%s", out)
	}
	if !strings.Contains(err.Error(), "jaira restore") {
		t.Errorf("the refusal does not say how to get it back: %v", err)
	}

	if out, err := runCLI(t, dir, "restore", "round-one.md"); err != nil {
		t.Fatalf("restore: %v\n%s", err, out)
	}
	after, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatalf("restore did not put the milestone back on the board: %v", err)
	}
	if after.Filed() {
		t.Error("the restored milestone still says it is filed")
	}
	if after.Colour != before.Colour {
		t.Errorf("colour came back as %d, want %d", after.Colour, before.Colour)
	}
	if got := after.Members(); len(got) != 1 || got[0] != done {
		t.Errorf("ticket list came back as %v, want [%s]", got, done)
	}
	if out, err := runCLI(t, dir, "milestone", "ls"); err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	} else if !strings.Contains(out, "round-one") {
		t.Errorf("'milestone ls' does not name the restored milestone:\n%s", out)
	}
}

// A marked file lying on the board is the normal outcome of pulling somebody
// else's filing, not a corner case: core/refsync IncomingMilestones writes a
// filed ref over a file this tree already has. 'add' and 'rm' are the third
// door into that file, after create and logbook, and the write would not stay
// local — recordMilestone puts it on the ref, so every clone that fetches sees
// a milestone its filer took off the board being edited.
func TestAddAndRmRefuseAFiledMilestoneOnDisk(t *testing.T) {
	dir := emptyStore(t)
	member := mkTicket(t, dir, "in the group")
	other := mkTicket(t, dir, "not in it yet")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", member); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	ms.SetStatus(milestone.StatusFiled)
	if err := ms.Save(dir); err != nil {
		t.Fatal(err)
	}

	for _, verb := range []string{"add", "rm"} {
		out, err := runCLI(t, dir, "milestone", verb, "round-one", other)
		if err == nil {
			t.Fatalf("'milestone %s' edited a filed milestone:\n%s", verb, out)
		}
		// The refusal has to name the file, its mark and where the way back
		// runs — the same three things create's refusal names, because the
		// reader is standing in front of the same file.
		for _, want := range []string{milestone.Path(dir, "round-one"), milestone.StatusFiled, "jaira restore round-one.md"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("the %s refusal does not mention %q: %v", verb, want, err)
			}
		}
	}

	// And nothing was written: the members are as they were and the mark stands,
	// so the ref this tree records is still the filed one.
	after, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if !after.Filed() {
		t.Error("the refused call took the filed mark off the file")
	}
	if got := after.Members(); len(got) != 1 || got[0] != ticketByHandle(t, dir, member).ID {
		t.Errorf("the refused call wrote the member list as %v", got)
	}
}

// Filed in this very tree, the file is in the logbook and Load answers the same
// ErrNotExist a name nobody ever used answers with. Pointing at 'create' there
// walks the reader into create's own refusal: two steps for one answer, and the
// first one points away from the file.
func TestAddOnAMilestoneFiledInThisTreePointsAtRestore(t *testing.T) {
	dir := emptyStore(t)
	done := mkDoneTicket(t, dir, "finished work")
	other := mkTicket(t, dir, "something else")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", done); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "milestone", "add", "round-one", other)
	if err == nil {
		t.Fatalf("'milestone add' accepted a milestone that is in the logbook:\n%s", out)
	}
	if !strings.Contains(err.Error(), "jaira restore round-one.md") {
		t.Errorf("the refusal does not point at the restore: %v", err)
	}
	if strings.Contains(err.Error(), "milestone create") {
		t.Errorf("the refusal still sends the reader to 'create', which refuses in turn: %v", err)
	}
	if !strings.Contains(err.Error(), "logbook") {
		t.Errorf("the refusal does not say where the file went: %v", err)
	}
}

// Three commands walk into a milestone that is filed, and the reader is
// standing in front of the same file at all three: its file is on disk and
// marked, and the only way back runs in the tree that filed it. The refusal
// used to be typed out at each door, so this measures that the three facts
// arrive whichever door was used — a door that drops one of them teaches the
// reader that it is a different problem.
func TestEveryDoorIntoAFiledMilestoneSaysTheSameThings(t *testing.T) {
	dir := emptyStore(t)
	member := mkTicket(t, dir, "in the group")
	other := mkTicket(t, dir, "not in it yet")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", member); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	ms.SetStatus(milestone.StatusFiled)
	if err := ms.Save(dir); err != nil {
		t.Fatal(err)
	}

	doors := map[string][]string{
		"create":  {"milestone", "create", "round-one"},
		"add":     {"milestone", "add", "round-one", other},
		"rm":      {"milestone", "rm", "round-one", other},
		"logbook": {"logbook", "round-one"},
	}
	for door, args := range doors {
		out, err := runCLI(t, dir, args...)
		if err == nil {
			t.Fatalf("'%s' went through a filed milestone:\n%s", door, out)
		}
		for _, want := range []string{
			milestone.Path(dir, "round-one"),
			milestone.StatusFiled,
			"jaira restore round-one.md",
			"the tree that filed it",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("the %s refusal does not mention %q: %v", door, want, err)
			}
		}
	}
}

// Restore is the fourth writer of a milestone file, after create, add/rm and
// logbook, and it has to queue behind them: it moves the file into
// .jaira/milestones/ and then read-modify-writes it to take the "filed" line
// off, which is exactly what a concurrent 'milestone add' does to the same
// file. Without the lock one of the two writes is lost.
func TestRestoreWaitsForTheMilestoneLock(t *testing.T) {
	dir := emptyStore(t)
	done := mkDoneTicket(t, dir, "finished work")
	if out, err := runCLI(t, dir, "milestone", "create", "round-one", done); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "logbook", "round-one"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}

	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	unlock, err := s.Lock(milestoneLockName)
	if err != nil {
		t.Fatal(err)
	}

	restored := make(chan error, 1)
	go func() {
		_, err := runCLI(t, dir, "restore", "round-one.md")
		restored <- err
	}()

	select {
	case err := <-restored:
		unlock()
		t.Fatalf("restore did not wait for the milestone lock (returned %v)", err)
	case <-time.After(150 * time.Millisecond):
	}
	// The file has not moved either: the lock is taken before s.Restore, not
	// between the move and the unfiling.
	if _, err := os.Stat(milestone.Path(dir, "round-one")); !os.IsNotExist(err) {
		unlock()
		t.Fatalf("the milestone file was moved back while the lock was held (stat: %v)", err)
	}

	unlock()
	select {
	case err := <-restored:
		if err != nil {
			t.Fatalf("restore after unlock: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("restore never completed after the lock was released")
	}
	ms, err := milestone.Load(dir, "round-one")
	if err != nil {
		t.Fatalf("restore did not put the milestone back: %v", err)
	}
	if ms.Filed() {
		t.Error("the restored milestone still says it is filed")
	}
}
