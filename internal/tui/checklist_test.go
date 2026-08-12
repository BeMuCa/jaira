package tui

import (
	"os"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/BeMuCa/jaira/core/ticket"
)

func withChecklists(body string) *ticket.Ticket {
	t := &ticket.Ticket{
		ID:     "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:  "Rate limit the login endpoint",
		Status: "in-progress",
		Body:   body,
	}
	t.DoDItems = ticket.ParseDoDItems(body)
	t.PlanItems = ticket.ParsePlanItems(body)
	return t
}

const twoChecklists = `## Plan

- [x] write the spec
- [~] design the interface
- [ ] implement

## Definition of Done

- [x] 429 returned above 100/min
- [ ] documented in README
`

// The counts are the whole point: "designed but not built" has to be legible
// without opening the ticket, because the failure being designed against is
// believing a ticket is finished when only its first step is.
func TestCardShowsChecklistProgress(t *testing.T) {
	m := newTestModel(t, 150, 32)
	out := stripANSI(m.renderCard(withChecklists(twoChecklists), 40, false))
	if !strings.Contains(out, "Plan 1/3") {
		t.Errorf("card does not show plan progress:\n%s", out)
	}
	if !strings.Contains(out, "DoD 1/2") {
		t.Errorf("card does not show definition-of-done progress:\n%s", out)
	}
}

// A ticket with neither checklist must not gain empty counters.
func TestCardWithoutChecklistsIsUnchanged(t *testing.T) {
	m := newTestModel(t, 150, 32)
	out := stripANSI(m.renderCard(withChecklists("no checklists here\n"), 40, false))
	if strings.Contains(out, "Plan ") || strings.Contains(out, "DoD ") {
		t.Errorf("counters appeared on a ticket with no checklists:\n%s", out)
	}
}

// Two different states both mean "this is waiting on you", but they call for
// different actions — answering a question, versus signing off a review — so
// they must be distinguishable rather than sharing one colour.
func TestWaitingStatesAreDistinct(t *testing.T) {
	m := newTestModel(t, 150, 32)

	human := withChecklists(twoChecklists)
	human.Status = "human"
	human.Question = "should this be per-IP or per-account?"
	humanOut := stripANSI(m.renderCard(human, 40, false))

	review := withChecklists(twoChecklists)
	review.Status = "review"
	reviewOut := stripANSI(m.renderCard(review, 40, false))

	if !strings.Contains(humanOut, "asks") {
		t.Errorf("a ticket waiting on an answer is not labelled:\n%s", humanOut)
	}
	if !strings.Contains(reviewOut, "sign off") {
		t.Errorf("a ticket waiting on sign-off is not labelled:\n%s", reviewOut)
	}
	if humanOut == reviewOut {
		t.Error("the two waiting states render identically")
	}
}

// The detail pane must show the step being worked on, marked distinctly from the
// rest, since "which one is the agent on" is the question it exists to answer.
func TestDetailRendersChecklists(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.detail = withChecklists(twoChecklists)
	out := stripANSI(m.renderDetail())

	for _, want := range []string{
		"Plan", "write the spec", "design the interface", "implement",
		"Definition of Done", "documented in README",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("detail pane is missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "→") {
		t.Errorf("the in-progress step is not marked:\n%s", out)
	}
	// The raw markdown must not also be dumped underneath the rendered lists.
	if strings.Contains(out, "- [x] write the spec") {
		t.Errorf("raw checklist markdown was printed as well as the rendered list:\n%s", out)
	}
}

// A ticked box without its evidence is a claim nobody can check, so a proof
// must render under its item in the detail pane, for both checklists.
func TestDetailRendersProofUnderItsItem(t *testing.T) {
	m := newTestModel(t, 150, 32)
	body := `## Plan

- [x] write the spec
  proof: docs/spec.md

## Definition of Done

- [x] 429 returned above 100/min
  proof: internal/x.go:12; TestRateLimit
- [ ] documented in README
`
	m.detail = withChecklists(body)
	out := stripANSI(m.renderDetail())

	if !strings.Contains(out, "proof: docs/spec.md") {
		t.Errorf("plan item's proof did not render:\n%s", out)
	}
	if !strings.Contains(out, "proof: internal/x.go:12; TestRateLimit") {
		t.Errorf("definition-of-done item's proof did not render:\n%s", out)
	}
}

// A narrow terminal must not be able to push a checklist line past the edge of
// the pane. Measured with lipgloss.Width rather than len, because display width
// is what actually wraps a terminal.
//
// Note: other parts of the detail pane still overflow at narrow widths — that is
// a known, separate bug. This test guards the checklist rendering only.
func TestChecklistLinesDoNotOverflow(t *testing.T) {
	m := newTestModel(t, 90, 32)
	m.detail = withChecklists(twoChecklists)
	for _, w := range []int{20, 24, 40, 80, 120} {
		m.width = w
		for _, line := range strings.Split(stripANSI(m.renderDetail()), "\n") {
			isChecklistLine := strings.Contains(line, "[x]") ||
				strings.Contains(line, "[~]") ||
				strings.Contains(line, "[ ]") ||
				strings.HasPrefix(strings.TrimSpace(line), "Plan ") ||
				strings.HasPrefix(strings.TrimSpace(line), "Definition of Done")
			if isChecklistLine && lipgloss.Width(line) > w {
				t.Errorf("width %d: checklist line is %d cols: %q", w, lipgloss.Width(line), line)
			}
		}
	}
}

// The board must not contradict the gate. A ticket whose definition of done is a
// checklist in the body satisfies the promotion rule, so the detail pane must not
// claim it is still missing one.
func TestChecklistCountsAsADefinitionOfDone(t *testing.T) {
	tk := withChecklists(twoChecklists)
	tk.Goal, tk.Context, tk.Assignee = "g", "c", "berk"
	if got := missing(tk); len(got) != 0 {
		t.Errorf("detail pane reports missing fields %v for a ticket with a DoD checklist", got)
	}
}

// Archiving from the board must move the file, not delete it, and the ticket
// must leave the board.
func TestArchiveFromTheBoard(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	before := len(m.tickets)
	id := m.selected().ID

	m.archiveSelected()

	if len(m.tickets) != before-1 {
		t.Errorf("board still has %d tickets, want %d", len(m.tickets), before-1)
	}
	if _, err := m.store.Load(id); err == nil {
		t.Error("the ticket is still loadable from the board")
	}
	names, err := os.ReadDir(m.store.ArchiveDir())
	if err != nil || len(names) != 1 {
		t.Fatalf("archive holds %v (err %v), want one file", names, err)
	}
}
