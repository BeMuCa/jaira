package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

const checkSteps = "1. go run ./cmd/app  2. open /export  3. the download starts at once"

// The sign-off screen is where a person decides. Every other field there is an
// account of what happened; the check is the only one they can act on, so it has
// to be on that screen.
func TestSignOffScreenShowsTheCheck(t *testing.T) {
	m := newTestModel(t, 140, 40)
	m.detail = &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "Fix the export", Status: "signoff",
		Goal: "stop the export hanging", ReviewVerdict: "matches the criteria",
		ReviewCheck: checkSteps,
	}

	out := stripANSI(m.renderSignOff())

	if !strings.Contains(out, "try it") {
		t.Errorf("the sign-off screen has no 'try it' row:\n%s", out)
	}
	if !strings.Contains(out, "go run ./cmd/app") {
		t.Errorf("the check steps are not on the sign-off screen:\n%s", out)
	}
	// It comes first, before the problem and the verdict: the person tries it
	// before they read the account.
	i, j, k := strings.Index(out, "try it"), strings.Index(out, "stop the export hanging"), strings.Index(out, "matches the criteria")
	if i < 0 || j < 0 || k < 0 || i > j || i > k {
		t.Errorf("the 'try it' row is not above the problem and the verdict:\n%s", out)
	}
	if !strings.Contains(out, "\n             2. open /export\n") {
		t.Errorf("the sign-off screen did not put one step per line:\n%s", out)
	}
}

// On a ticket in a person's lane the detail pane opens with the same 'try it'
// row, above the goal, and the Review group no longer repeats the check.
func TestDetailPaneOpensWithTheCheckOnAHumanLane(t *testing.T) {
	m := newTestModel(t, 140, 40)
	tk := &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "Fix the export", Status: "signoff",
		Goal: "stop the export hanging", ReviewVerdict: "matches the criteria",
		ReviewCheck: checkSteps,
	}

	out := stripANSI(m.detailBody(tk, 100))

	i, j := strings.Index(out, "try it"), strings.Index(out, "stop the export hanging")
	if i < 0 || j < 0 || i > j {
		t.Errorf("the detail pane does not open with 'try it' above the goal:\n%s", out)
	}
	if strings.Count(out, "go run ./cmd/app") != 1 {
		t.Errorf("the check is shown more than once on a human lane:\n%s", out)
	}
}

// Off a person's lane the check stays where it was: in the Review group, after
// the verdict, and not under 'try it'.
func TestDetailPaneKeepsTheCheckInReviewOffAHumanLane(t *testing.T) {
	m := newTestModel(t, 140, 40)
	tk := &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "Fix the export", Status: "review",
		ReviewVerdict: "matches the criteria", ReviewCheck: checkSteps,
	}

	out := stripANSI(m.detailBody(tk, 100))

	if strings.Contains(out, "try it") {
		t.Errorf("a non-human lane shows a 'try it' row:\n%s", out)
	}
	i, j := strings.Index(out, "verdict"), strings.Index(out, "check")
	if i < 0 || j < 0 || j < i {
		t.Errorf("the check is not after the verdict in the Review group:\n%s", out)
	}
}

func TestOneStepPerLine(t *testing.T) {
	cases := map[string]string{
		"1. a  2. b  3. c":           "1. a\n2. b\n3. c",
		"1. a 2. b":                  "1. a\n2. b",
		"1. build\n2. look":          "1. build\n2. look", // already lines: unchanged
		"just prose, no steps":       "just prose, no steps",
		"2. starts at two  3. three": "2. starts at two  3. three", // not step 1: unchanged
		"1. only one step":           "1. only one step",
		"  1. padded  2. steps  ":    "1. padded\n2. steps",
	}
	for in, want := range cases {
		if got := oneStepPerLine(in); got != want {
			t.Errorf("oneStepPerLine(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDetailPaneShowsTheCheck(t *testing.T) {
	m := newTestModel(t, 140, 40)
	tk := &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "Fix the export", Status: "review",
		ReviewCheck: checkSteps,
	}

	out := stripANSI(m.detailBody(tk, 100))

	if !strings.Contains(out, "check") || !strings.Contains(out, "open /export") {
		t.Errorf("the open ticket does not show the check:\n%s", out)
	}
}

// A ticket with only a check and none of the other review fields must still show
// the Review block: the block's condition used to name three fields.
func TestReviewBlockAppearsForACheckAlone(t *testing.T) {
	m := newTestModel(t, 140, 40)
	tk := &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "t", Status: "review",
		ReviewCheck: checkSteps,
	}

	if out := stripANSI(m.detailBody(tk, 100)); !strings.Contains(out, "Review") {
		t.Errorf("no Review heading for a ticket carrying only a check:\n%s", out)
	}
}

// Steps written one per line stay one per line — a numbered check must read
// as a list, not as a flattened paragraph, on both screens that show it.
func TestAMultilineCheckKeepsItsLines(t *testing.T) {
	steps := "1. build it\n2. open the board\n3. squint at the cards"
	indent := "\n             " // a following line starts under the label column

	m := newTestModel(t, 140, 40)
	tk := &ticket.Ticket{
		ID: "01M0F84EW52Y5HGH84R74VEH6C", Title: "t", Status: "review",
		ReviewCheck: steps,
	}
	out := stripANSI(m.detailBody(tk, 100))
	if !strings.Contains(out, indent+"2. open the board") || !strings.Contains(out, indent+"3. squint at the cards") {
		t.Errorf("detail pane flattened the check's lines:\n%s", out)
	}

	tk.Status = "signoff"
	tk.ReviewVerdict = "fine"
	m.detail = tk
	sign := stripANSI(m.renderSignOff())
	if !strings.Contains(sign, indent+"2. open the board") {
		t.Errorf("sign-off screen flattened the check's lines:\n%s", sign)
	}
}

// The third label-column renderer: lane-declared fields go through fieldRow on
// both screens, and it has to keep the author's lines exactly like the others —
// the review that accepted wrapField found this one still flattening.
func TestFieldRowKeepsTheAuthorsLines(t *testing.T) {
	var b strings.Builder
	fieldRow(&b, "check", "1. build\n2. look", 80)
	out := stripANSI(b.String())
	if !strings.Contains(out, "\n             2. look") {
		t.Errorf("fieldRow flattened the lines:\n%s", out)
	}
}
