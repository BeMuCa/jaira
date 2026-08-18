package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

func laneFocusModel(t *testing.T, w, h int) *Model {
	t.Helper()
	m := newTestModel(t, w, h)
	m.mode = modeLaneFocus
	return m
}

// laneIndex finds the column index for a lane ID, since the fixture's lane
// order is defined by the built-in lane files, not by this test.
func laneIndex(t *testing.T, m *Model, id string) int {
	t.Helper()
	for i, c := range m.cols {
		if c.lane.ID == id {
			return i
		}
	}
	t.Fatalf("no lane %q in fixture", id)
	return -1
}

// The compact view's one drill-down goes deeper into the step picked, not
// sideways into the multi-column board.
func TestLaneFocusOpensFromPipelineNotBoard(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modePipeline
	m.laneIdx = laneIndex(t, m, "backlog")
	m.pipelineKey("enter")
	if m.mode != modeLaneFocus {
		t.Fatalf("enter opened mode %v, want modeLaneFocus", m.mode)
	}
}

func TestLaneFocusShowsOnlyItsOwnLane(t *testing.T) {
	m := laneFocusModel(t, 120, 34)
	m.laneIdx = laneIndex(t, m, "backlog")
	out := stripANSI(m.renderLaneFocus())
	if !strings.Contains(out, "Backlog") {
		t.Errorf("lane name missing:\n%s", out)
	}
	if !strings.Contains(out, "Investigate flaky logout test") {
		t.Errorf("backlog ticket missing:\n%s", out)
	}
	if !strings.Contains(out, "Rate limit the login endpoint") {
		t.Errorf("backlog ticket missing:\n%s", out)
	}
	if strings.Contains(out, "Refactor auth middleware") {
		t.Errorf("a todo-lane ticket leaked into the backlog view:\n%s", out)
	}
}

func TestLaneFocusChangesLaneWithoutLeaving(t *testing.T) {
	m := laneFocusModel(t, 120, 34)
	start := m.laneIdx
	m.laneFocusKey("l")
	if m.laneIdx == start {
		t.Fatal("l did not move lane")
	}
	if m.mode != modeLaneFocus {
		t.Errorf("l left the single-lane view (mode=%v)", m.mode)
	}
	m.laneFocusKey("h")
	if m.laneIdx != start {
		t.Error("h did not move back")
	}
	if m.mode != modeLaneFocus {
		t.Errorf("h left the single-lane view (mode=%v)", m.mode)
	}
}

func TestLaneFocusEscReturnsToPipeline(t *testing.T) {
	m := laneFocusModel(t, 120, 34)
	m.laneFocusKey("esc")
	if m.mode != modePipeline {
		t.Errorf("esc returned to mode %v, want modePipeline", m.mode)
	}
}

func TestLaneFocusEnterOpensDetail(t *testing.T) {
	m := laneFocusModel(t, 120, 34)
	m.laneIdx = laneIndex(t, m, "backlog")
	m.cardIdx = 0
	m.laneFocusKey("enter")
	if m.mode != modeDetail {
		t.Fatalf("enter opened mode %v, want modeDetail", m.mode)
	}
	if m.detailFrom != modeLaneFocus {
		t.Errorf("detailFrom = %v, want modeLaneFocus so esc comes back here", m.detailFrom)
	}
}

// The single-lane view must never render taller than the terminal, whatever
// the height — it used to render every ticket in the lane with no cap at all.
func TestLaneFocusFitsTheTerminal(t *testing.T) {
	s := flagHeavyReviewLane(t, 35)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.mode = modeLaneFocus
	m.laneIdx = laneIndex(t, m, "review")

	for _, h := range []int{40, 30, 24, 20, 16, 12} {
		m.width, m.height = 120, h
		got := len(strings.Split(stripANSI(m.renderLaneFocus()), "\n"))
		if got > h {
			t.Errorf("terminal h=%d: lane focus rendered %d lines", h, got)
		}
	}
}

// Walking the cursor through a lane too long for the pane must keep the
// selected ticket on screen at every step, the single-lane equivalent of
// TestSelectionStaysVisibleWhenScrollingDown.
func TestLaneFocusKeepsTheCursorVisible(t *testing.T) {
	s := flagHeavyReviewLane(t, 35)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.mode = modeLaneFocus
	m.laneIdx = laneIndex(t, m, "review")
	m.width, m.height = 120, 20

	col := m.cols[m.laneIdx]
	for i := 0; i < len(col.tickets); i++ {
		m.cardIdx = i
		out := stripANSI(m.renderLaneFocus())
		want := ticket.Handle(col.tickets[i].ID)
		if !strings.Contains(out, want) {
			t.Fatalf("selection at index %d (%s) is not visible in the render:\n%s", i, want, out)
		}
	}
}

// A lane cut short by the terminal must say so — nothing is hidden silently.
func TestLaneFocusSaysWhatIsOffScreen(t *testing.T) {
	s := flagHeavyReviewLane(t, 35)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.mode = modeLaneFocus
	m.laneIdx = laneIndex(t, m, "review")
	m.width, m.height = 120, 20
	m.cardIdx = len(m.cols[m.laneIdx].tickets) / 2

	out := stripANSI(m.renderLaneFocus())
	if !strings.Contains(out, "more") {
		t.Errorf("nothing told the reader there is more off-screen:\n%s", out)
	}
}

func TestLaneFocusEmptyLaneDoesNotPanic(t *testing.T) {
	m := laneFocusModel(t, 120, 34)
	idx := laneIndex(t, m, "blocked") // built-in lane with no tickets in the fixture
	m.laneIdx = idx
	out := stripANSI(m.renderLaneFocus())
	if !strings.Contains(out, "no tickets") {
		t.Errorf("empty lane did not say so:\n%s", out)
	}
	if !strings.Contains(out, "Blocked") {
		t.Errorf("empty lane's name missing:\n%s", out)
	}
	m.laneFocusKey("l")
	if m.laneIdx == idx {
		t.Error("l did not move on from the empty lane")
	}
}
