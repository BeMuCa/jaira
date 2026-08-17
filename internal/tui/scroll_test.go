package tui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// reviewLane fills a store with n tickets in the review lane, so a single column
// is longer than any pane can show.
func reviewLane(t *testing.T, n int) *ticket.Store {
	t.Helper()
	s := newTestStore(t)
	now := time.Now()
	for i := 0; i < n; i++ {
		f := map[string]string{
			ticket.FieldID:        ticket.NewID(now),
			ticket.FieldTitle:     "Review item " + string(rune('A'+i%26)),
			ticket.FieldStatus:    "review",
			ticket.FieldReady:     "true",
			ticket.FieldCreator:   "berk",
			ticket.FieldAssignee:  "berk",
			ticket.FieldGoal:      "g",
			ticket.FieldDoD:       "d",
			ticket.FieldContext:   "c",
			ticket.FieldCreatedAt: ticket.FormatTime(now),
			ticket.FieldUpdatedAt: ticket.FormatTime(now),
		}
		l := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: {"abc1234"}}
		if _, err := s.Create(f, l, ""); err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Millisecond)
	}
	return s
}

// Scrolling down a lane taller than the pane must bring the selected card into
// view. Without this the selection walks off the bottom and the board shows a
// cursor that is not there — you cannot see what you are about to act on.
func TestSelectionStaysVisibleWhenScrollingDown(t *testing.T) {
	s := reviewLane(t, 30)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	// A deliberately short terminal: the case is "I made the terminal smaller".
	m.width, m.height = 120, 20
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}

	// Focus the review column.
	found := -1
	for i, c := range m.cols {
		if c.lane.ID == "review" {
			found = i
		}
	}
	if found < 0 {
		t.Fatal("no review column on the board")
	}
	m.laneIdx = found
	col := m.cols[m.laneIdx]
	if len(col.tickets) < 30 {
		t.Fatalf("review column has %d tickets, want 30", len(col.tickets))
	}

	// Walk the selection all the way down, rendering at each step the way the
	// real board does, so the scroll state evolves exactly as it would in use.
	for i := 0; i < len(col.tickets); i++ {
		m.cardIdx = i
		out := stripANSI(m.render())
		want := ticket.Handle(col.tickets[i].ID)
		if !strings.Contains(out, want) {
			t.Fatalf("selection at index %d (%s) is not visible in the render:\n%s", i, want, out)
		}
	}
}

// longTicket is a ticket whose rendered detail is taller than any small
// terminal: several checklists, an outcome, a review and a body.
func longTicket() *ticket.Ticket {
	body := `## Plan

- [x] write the spec
- [x] design the interface
- [~] implement
- [ ] test

## Definition of Done

- [ ] documented in README
- [ ] covered by a test
- [ ] reviewed

## Notes

Tried the buffered writer first; it holds the whole export in memory.
Streaming per 5k rows works.
`
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Export survives a slow network",
		Status: "review", Goal: "the CSV export completes on a 3G connection",
		Context:  "reported while debugging the failed month-end export",
		Assignee: "berk", Creator: "berk", ExecutedBy: "claude-opus-5",
		Outcome: ticket.Outcome{
			What: "streamed the writer", Why: "it buffered the whole export",
			Resolves: "a 40MB export now completes",
		},
		ReviewSummary: "streamed the writer", ReviewGaps: "none",
		ReviewVerdict: "the diff matches the criteria",
		Body:          body,
	}
	tk.DoDItems = ticket.ParseDoDItems(body)
	tk.PlanItems = ticket.ParsePlanItems(body)
	return tk
}

// The open ticket must never render taller than the terminal. It used to render
// whole regardless of height, which pushed the handle, title and goal off the top
// with no key that could bring them back.
func TestDetailFitsTheTerminal(t *testing.T) {
	for _, h := range []int{40, 30, 24, 20, 16, 12, 8} {
		m := newTestModel(t, 100, h)
		m.detail = longTicket()
		got := len(strings.Split(stripANSI(m.renderDetail()), "\n"))
		if got > h {
			t.Errorf("terminal h=%d: detail rendered %d lines", h, got)
		}
	}
}

// The top of the ticket is what identifies it, so it must be what you see first.
func TestDetailStartsAtTheTop(t *testing.T) {
	m := newTestModel(t, 100, 16)
	m.detail = longTicket()

	out := stripANSI(m.renderDetail())
	if !strings.Contains(out, "Export survives a slow network") {
		t.Errorf("the title is not on screen at scroll 0:\n%s", out)
	}
	if !strings.Contains(out, "more ·") {
		t.Errorf("nothing told the reader there is more below:\n%s", out)
	}
}

// Content below the fold has to be reachable, and scrolling must stop at the end
// rather than running off into blank screens.
func TestDetailScrollsToTheEndAndStops(t *testing.T) {
	m := newTestModel(t, 100, 16)
	m.detail = longTicket()

	// The notes are the last thing in the body, so they prove the bottom is
	// reachable at all.
	var reached bool
	for i := 0; i < 40; i++ {
		if strings.Contains(stripANSI(m.renderDetail()), "Streaming per 5k rows works.") {
			reached = true
			break
		}
		m.detailScroll += 4
	}
	if !reached {
		t.Fatal("the end of the ticket is not reachable by scrolling")
	}

	// Scrolling far past the end must clamp, not blank the pane.
	m.detailScroll = 10_000
	out := stripANSI(m.renderDetail())
	if strings.TrimSpace(out) == "" {
		t.Error("scrolling past the end blanked the pane")
	}
	if !strings.Contains(out, "Streaming per 5k rows works.") {
		t.Errorf("clamped scroll does not show the last page:\n%s", out)
	}
}

// The arrow keys are the base movement vocabulary: pressed in an open ticket
// they must scroll it, one line per press, and j/k must keep jumping between
// tickets rather than scrolling. This drives the real key dispatch, not the
// renderer.
func TestArrowKeysScrollTheOpenTicket(t *testing.T) {
	m := newTestModel(t, 100, 16)
	m.detail = longTicket()
	m.mode = modeDetail
	m.detailFrom = modeBoard

	m.key(tea.KeyPressMsg{Code: tea.KeyDown})
	m.key(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.detailScroll != 2 {
		t.Errorf("two ↓ presses moved detailScroll to %d, want 2", m.detailScroll)
	}
	m.key(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.detailScroll != 1 {
		t.Errorf("↑ moved detailScroll to %d, want 1", m.detailScroll)
	}
	if m.mode != modeDetail || m.detail == nil {
		t.Error("arrow keys left the open ticket, they must only scroll it")
	}

	// j leaves the pane and selects the neighbouring card — unchanged.
	m.key(key("j"))
	if m.mode == modeDetail {
		t.Error("j did not leave the detail pane to jump to the next ticket")
	}
}
