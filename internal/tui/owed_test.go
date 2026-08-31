package tui

// A field an installed lane declares it produces is part of the ticket's shape
// whether or not the lane has run. Nine tickets once reached a review lane
// unworked and read as finished work, because an empty field renders as
// nothing at all. These tests pin the opposite: the row stays and names the
// lane that owes it.

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/ticket"
)

// unworkedInReview is a ticket that reached the review lane without a single
// lane before it producing anything — the state on Berk's own board that
// started this.
func unworkedInReview() *ticket.Ticket {
	return &ticket.Ticket{
		ID:       "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:    "Rate limit the login endpoint",
		Status:   "review",
		Assignee: "berk",
		Creator:  "berk",
		Goal:     "stop credential stuffing",
		Context:  "came up while reading the auth logs",
		// Set because every ticket on disk carries them and the when row is one
		// of the base fields these tests assert the position of.
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}
}

// rowIndex returns the line number carrying a label column, so the reading
// order can be asserted without depending on the text beside the label.
func rowIndex(t *testing.T, out, label string) int {
	t.Helper()
	for i, l := range strings.Split(out, "\n") {
		if f := strings.Fields(l); len(f) > 0 && f[0] == label {
			return i
		}
	}
	t.Fatalf("no %q row in:\n%s", label, out)
	return -1
}

// The seven questions a reviewer asks, in order, on a ticket where six of them
// have no answer yet: each one still gets a row, and the row says who owes it.
func TestDetailShowsOwedFieldsNamingTheLane(t *testing.T) {
	m := newTestModel(t, 120, 40)
	out := stripANSI(m.detailBody(unworkedInReview(), 120))

	// in-progress declares the three outcome fields, review the four review
	// fields; both are on this ticket's route, so both are in debt.
	for _, want := range []string{
		"what         — owed by in-progress",
		"why          — owed by in-progress",
		"resolves     — owed by in-progress",
		"summary      — owed by review",
		"gaps         — owed by review",
		"verdict      — owed by review",
		"check        — owed by review",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing owed row %q:\n%s", want, out)
		}
	}

	// problem, what, why, resolves, summary, gaps, check — the order the
	// sign-off decision is made in.
	order := []string{"goal", "what", "why", "resolves", "summary", "gaps", "check"}
	prev := -1
	for _, label := range order {
		at := rowIndex(t, out, label)
		if at <= prev {
			t.Errorf("%q is out of order in:\n%s", label, out)
		}
		prev = at
	}

	// The base fields lead, in every lane, whatever else is on the ticket.
	base := []string{"id", "lane", "assignee", "creator", "when", "goal", "context"}
	prev = -1
	for _, label := range base {
		at := rowIndex(t, out, label)
		if at <= prev {
			t.Errorf("base field %q is out of order in:\n%s", label, out)
		}
		prev = at
	}
	if rowIndex(t, out, "context") > rowIndex(t, out, "what") {
		t.Errorf("the base fields do not lead:\n%s", out)
	}
}

// A field that has been filled in shows its value. The debt row is a stand-in
// for a missing answer, never an addition to one.
func TestDetailShowsTheValueRatherThanTheDebt(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := unworkedInReview()
	tk.Outcome.What = "added a token bucket per IP"
	tk.ReviewSummary = "the limiter is per client IP"
	out := stripANSI(m.detailBody(tk, 120))

	for _, want := range []string{"added a token bucket per IP", "the limiter is per client IP"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing value %q:\n%s", want, out)
		}
	}
	for _, unwanted := range []string{"what         — owed", "summary      — owed"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("a filled field still shows a debt row %q:\n%s", unwanted, out)
		}
	}
	// The fields beside them are still owed — filling one in does not settle
	// the lane.
	if !strings.Contains(out, "why          — owed by in-progress") {
		t.Errorf("the rest of the lane's output stopped being owed:\n%s", out)
	}
}

// The goal is a declared field too — a brainstorm lane produces it — and it
// used to render through a helper that suppressed empty values, so a ticket
// that opted into brainstorming and never got one showed no goal row at all:
// the very defect this change exists to remove, still in place one row above
// the rows it fixed.
func TestDetailShowsTheGoalDebtWhenBrainstormIsOnTheRoute(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := unworkedInReview()
	tk.Goal = ""
	tk.Body = "## Options\n\n- [x] brainstorm\n"
	tk.Options = ticket.ParseOptions(tk.Body)

	out := stripANSI(m.detailBody(tk, 120))
	if !strings.Contains(out, "goal         — owed by brainstorm") {
		t.Errorf("an empty goal a lane owes is invisible:\n%s", out)
	}
	// And the row keeps its place: still a base field, still above the rest.
	if rowIndex(t, out, "goal") > rowIndex(t, out, "what") {
		t.Errorf("the goal debt row is out of order:\n%s", out)
	}
}

// plan is declared by pre-process and lives in the body as a checklist, not as
// a label-and-value row. Its debt is the empty checklist itself; printing a
// second, worse version of the same fact in the field rows is the thing being
// pinned against here.
func TestDetailShowsNoDebtRowForABodyChecklistField(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := unworkedInReview()
	tk.Body = "## Options\n\n- [x] planning\n"
	tk.Options = ticket.ParseOptions(tk.Body)

	out := stripANSI(m.detailBody(tk, 120))
	if strings.Contains(out, "owed by pre-process") {
		t.Errorf("a body-checklist field was rendered as a field row:\n%s", out)
	}
	// The lane does owe it, though — the renderer is choosing not to print a
	// row, and that must stay a choice rather than becoming a silent gap.
	if l := gate.OwedBy(m.lanes, tk)["plan"]; l != "pre-process" {
		t.Errorf("plan is owed by %q, want pre-process", l)
	}
}

// The rule is that a field an installed lane declares is shown even when
// empty, which is board-wide: a backlog ticket carries the same debts, because
// they are what the pipeline will ask of it. Pinned as a decision rather than
// left as an accident — the compact first glance is the folds cut's job, and
// whoever takes it should have to change this test on purpose.
func TestBacklogDetailCarriesTheSameDebts(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := unworkedInReview()
	tk.Status = "backlog"

	out := stripANSI(m.detailBody(tk, 120))
	for _, want := range []string{
		"what         — owed by in-progress",
		"summary      — owed by review",
		"check        — owed by review",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("a backlog ticket is missing %q:\n%s", want, out)
		}
	}
}

// The debt rows go through the same wrap as every other field: a pane cannot
// scroll sideways, so nothing may be wider than the width it was given.
func TestOwedRowsStayInsideTheWidth(t *testing.T) {
	for _, w := range []int{30, 60, 120} {
		m := newTestModel(t, w, 40)
		out := m.detailBody(unworkedInReview(), w)
		for i, l := range strings.Split(out, "\n") {
			if got := lipgloss.Width(l); got > w {
				t.Errorf("width %d: line %d is %d wide: %q", w, i, got, stripANSI(l))
			}
		}
	}
}

// The sign-off screen is where a person accepts the work, so it is the worst
// place for a lane that never ran to leave no trace.
func TestSignOffShowsOwedFieldsNamingTheLane(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := unworkedInReview()
	tk.Status = "signoff"
	m.detail = tk
	out := stripANSI(m.renderSignOff())

	for _, want := range []string{
		"stop credential stuffing", // the problem is still the problem
		"what         — owed by in-progress",
		"resolves     — owed by in-progress",
		"summary      — owed by review",
		"check        — owed by review",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sign-off view is missing %q:\n%s", want, out)
		}
	}
}
