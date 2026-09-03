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

	if !strings.Contains(out, "check") {
		t.Errorf("the sign-off screen has no check row:\n%s", out)
	}
	if !strings.Contains(out, "go run ./cmd/app") {
		t.Errorf("the check steps are not on the sign-off screen:\n%s", out)
	}
	// It comes after the account, because you read first and then go and look.
	if i, j := strings.Index(out, "verdict"), strings.Index(out, "check"); i >= 0 && j >= 0 && j < i {
		t.Error("the check is printed above the verdict; it belongs after the account")
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
