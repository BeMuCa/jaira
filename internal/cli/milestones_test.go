package cli

import (
	"os"
	"strings"
	"testing"

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
