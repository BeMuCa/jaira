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
