package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// The ten-lane fixture behind these tests holds tickets in five lanes
// (backlog, todo, in-progress, human, done) and leaves five empty
// (brainstorm, pre-process, review, signoff, blocked) — the same shape as a
// real board with more custom lanes declared than are actually in use.
var emptyLaneNames = []string{"Brainstorm", "Pre-process", "Review", "Human Review", "Blocked"}
var filledLaneNames = []string{"Backlog", "Todo", "Implementing", "HITL", "Done"}

// At 240 columns all ten lanes fit (10*(22+2) = 240), so this is the baseline
// the toggle is measured against: nothing about a wide terminal changes.
func TestWideBoardShowsAllTenLanesByDefault(t *testing.T) {
	m := newTestModel(t, 240, 40)
	out := stripANSI(m.render())
	for _, name := range append(append([]string{}, emptyLaneNames...), filledLaneNames...) {
		if !strings.Contains(out, name) {
			t.Errorf("lane %q missing from the baseline render:\n%s", name, out)
		}
	}
}

// z keeps every lane on the board and draws the five holding no tickets thin:
// the process is still all there, and a thin lane is how it says "nothing in
// here". The lanes with tickets get the room the thin ones give up.
func TestZDrawsEmptyLanesThinAndKeepsThemAll(t *testing.T) {
	m := newTestModel(t, 240, 40)
	before := m.boardFit(240)
	m.key(key("z"))
	win := m.boardFit(240)

	if win.start != 0 || win.end != len(m.cols) {
		t.Errorf("window %d..%d after z, want every one of the %d lanes", win.start, win.end, len(m.cols))
	}
	for i, c := range m.cols {
		if got, want := win.thin[i], len(c.tickets) == 0; got != want {
			t.Errorf("lane %q: thin = %v, want %v (it holds %d tickets)", c.lane.ID, got, want, len(c.tickets))
		}
	}
	if win.colW <= before.colW {
		t.Errorf("full columns did not widen: %d before z, %d after", before.colW, win.colW)
	}

	out := stripANSI(m.render())
	for _, name := range filledLaneNames {
		if !strings.Contains(out, name) {
			t.Errorf("filled lane %q disappeared after z:\n%s", name, out)
		}
	}
	if strings.Contains(out, "hidden") {
		t.Errorf("nothing is hidden any more, yet the render says so:\n%s", out)
	}
}

// The reason to press z at all: at 200 columns only eight full lanes fit, so
// two are off-screen; drawn thin, all ten fit (5×24 + 5×6 = 150 cells).
func TestZMakesRoomForLanesThatWereOffScreen(t *testing.T) {
	m := newTestModel(t, 200, 40)
	win := m.boardFit(200)
	if win.end-win.start >= len(m.cols) {
		t.Skipf("all %d lanes already fit at 200 columns; the fixture changed", len(m.cols))
	}

	m.key(key("z"))
	win = m.boardFit(200)
	if win.start != 0 || win.end != len(m.cols) {
		t.Errorf("window %d..%d after z at 200 columns, want all %d lanes", win.start, win.end, len(m.cols))
	}
}

// A second z restores every lane and removes the notice entirely.
func TestSecondZRestoresEverything(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	m.key(key("z"))
	out := stripANSI(m.render())

	for _, name := range append(append([]string{}, emptyLaneNames...), filledLaneNames...) {
		if !strings.Contains(out, name) {
			t.Errorf("lane %q missing after the second z:\n%s", name, out)
		}
	}
	if win := m.boardFit(240); len(win.thin) != 0 {
		t.Errorf("lanes still drawn thin after the second z: %v", win.thin)
	}
}

// The toggle changes only the picture: ticket counts, per-lane counts and the
// set of lane IDs in m.cols must be identical before and after both presses.
func TestToggleIsDisplayOnly(t *testing.T) {
	m := newTestModel(t, 240, 40)
	before := len(m.tickets)
	beforeCols := colTicketCounts(m)
	beforeIDs := colLaneIDs(m)

	m.key(key("z"))
	if got := len(m.tickets); got != before {
		t.Errorf("len(m.tickets) changed after z: %d -> %d", before, got)
	}
	if got := colTicketCounts(m); !equalIntMaps(got, beforeCols) {
		t.Errorf("per-lane ticket counts changed after z: %v -> %v", beforeCols, got)
	}
	if got := colLaneIDs(m); !equalStringSets(got, beforeIDs) {
		t.Errorf("lane ID set changed after z: %v -> %v", beforeIDs, got)
	}

	m.key(key("z"))
	if got := len(m.tickets); got != before {
		t.Errorf("len(m.tickets) changed after the second z: %d -> %d", before, got)
	}
	if got := colTicketCounts(m); !equalIntMaps(got, beforeCols) {
		t.Errorf("per-lane ticket counts changed after the second z: %v -> %v", beforeCols, got)
	}
	if got := colLaneIDs(m); !equalStringSets(got, beforeIDs) {
		t.Errorf("lane ID set changed after the second z: %v -> %v", beforeIDs, got)
	}
}

func colTicketCounts(m *Model) map[string]int {
	out := map[string]int{}
	for _, c := range m.cols {
		out[c.lane.ID] = len(c.tickets)
	}
	return out
}

func colLaneIDs(m *Model) map[string]bool {
	out := map[string]bool{}
	for _, c := range m.cols {
		out[c.lane.ID] = true
	}
	return out
}

func equalIntMaps(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func equalStringSets(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// assertStepSkipsEmptyLanes checks that stepping d lanes from `from` landed on
// `to`, that `to` holds tickets, and that every lane passed over is empty.
func assertStepSkipsEmptyLanes(t *testing.T, m *Model, from, to, d int) {
	t.Helper()
	n := len(m.cols)
	for i := ((from+d)%n + n) % n; i != to; i = ((i+d)%n + n) % n {
		if len(m.cols[i].tickets) != 0 {
			t.Errorf("step from %q to %q jumped over %q, which holds tickets", m.cols[from].lane.ID, m.cols[to].lane.ID, m.cols[i].lane.ID)
		}
	}
	if len(m.cols[to].tickets) == 0 {
		t.Errorf("step from %q landed on %q, which holds no tickets", m.cols[from].lane.ID, m.cols[to].lane.ID)
	}
}

// With z on, the board's ←/→ skip the thin lanes: from backlog, right lands on
// the next lane that holds a ticket, not on the empty brainstorm beside it.
func TestArrowsWithToggleOnSkipEmptyLanes(t *testing.T) {
	m := newTestModel(t, 240, 40)
	start := laneIndex(t, m, "backlog")
	m.laneIdx = start
	m.key(key("z"))

	m.key(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.laneIdx == laneIndex(t, m, "brainstorm") {
		t.Errorf("right from backlog with empty lanes thin landed on brainstorm, which is empty")
	}
	assertStepSkipsEmptyLanes(t, m, start, m.laneIdx, 1)
}

// Without z the arrows still step onto every lane, empty or not.
func TestArrowsWithToggleOffStepOntoEmptyLanes(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.laneIdx = laneIndex(t, m, "backlog")

	m.key(tea.KeyPressMsg{Code: tea.KeyRight})
	if got := m.cols[m.laneIdx].lane.ID; got != "brainstorm" {
		t.Errorf("right from backlog with z off landed on %q, want brainstorm (the adjacent lane)", got)
	}
}

// From backlog, left with z on wraps round the trailing empty lanes to the
// last lane that holds tickets.
func TestLeftWithToggleOnWrapsToLastFilledLane(t *testing.T) {
	m := newTestModel(t, 240, 40)
	start := laneIndex(t, m, "backlog")
	m.laneIdx = start
	m.key(key("z"))

	m.key(tea.KeyPressMsg{Code: tea.KeyLeft})
	if got := m.cols[m.laneIdx].lane.ID; got != "done" {
		t.Errorf("left from backlog with z on landed on %q, want done (the last lane holding tickets)", got)
	}
	assertStepSkipsEmptyLanes(t, m, start, m.laneIdx, -1)
}

// When the cursor's lane is the only one holding tickets, ←/→ have nowhere
// else to skip to and leave the cursor where it is.
func TestArrowsWithToggleOnStayWhenOnlyLaneHasTickets(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.filter = "Refactor"
	m.rebuild()
	filled := -1
	count := 0
	for i, c := range m.cols {
		if len(c.tickets) > 0 {
			filled = i
			count++
		}
	}
	if count != 1 {
		t.Fatalf("filter Refactor left %d lanes with tickets, want exactly 1", count)
	}
	m.laneIdx = filled
	m.key(key("z"))

	m.key(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.laneIdx != filled {
		t.Errorf("right moved to lane %d with only lane %d holding tickets", m.laneIdx, filled)
	}
	m.key(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.laneIdx != filled {
		t.Errorf("left moved to lane %d with only lane %d holding tickets", m.laneIdx, filled)
	}
}

// When no lane holds tickets at all, there is nothing to skip to, so the
// arrows step one lane at a time as they always did.
func TestArrowsWithToggleOnAllEmptyStillMove(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.filter = "zzz-matches-nothing"
	m.rebuild()
	for _, c := range m.cols {
		if len(c.tickets) != 0 {
			t.Fatalf("filter matching nothing left lane %q with tickets", c.lane.ID)
		}
	}
	m.laneIdx = 3
	m.key(key("z"))

	m.key(tea.KeyPressMsg{Code: tea.KeyRight})
	if want := (3 + 1) % len(m.cols); m.laneIdx != want {
		t.Errorf("right with every lane empty landed on lane %d, want %d", m.laneIdx, want)
	}
}

// Standing on an empty lane when the toggle is on, the arrows still go to the
// next lane that holds tickets in that direction.
func TestArrowsWithToggleOnLeaveEmptyStartLane(t *testing.T) {
	m := newTestModel(t, 240, 40)
	start := laneIndex(t, m, "pre-process")
	m.laneIdx = start
	m.key(key("z"))

	m.key(tea.KeyPressMsg{Code: tea.KeyRight})
	assertStepSkipsEmptyLanes(t, m, start, m.laneIdx, 1)
}

// Toggling on while parked on an empty lane leaves the cursor where it is:
// the lane is drawn thin, not removed, so there is nothing to move away from.
func TestTogglingOnLeavesTheCursorOnAnEmptyLane(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.laneIdx = laneIndex(t, m, "pre-process")
	m.key(key("z"))

	if got := m.cols[m.laneIdx].lane.ID; got != "pre-process" {
		t.Errorf("cursor moved to %q on toggling from an empty lane, want it to stay on pre-process", got)
	}
	if win := m.boardFit(240); !win.thin[m.laneIdx] {
		t.Errorf("the focused empty lane is not drawn thin")
	}
}

// The toggle is a board preference: laneFocusKey's h/l step to the
// immediately adjacent lane regardless of thinEmpty, since that view never
// shows the toggle at all.
func TestLaneFocusIgnoresTheToggle(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.mode = modeLaneFocus
	m.laneIdx = laneIndex(t, m, "backlog")
	m.thinEmpty = true

	m.key(key("h"))
	if got := m.cols[m.laneIdx].lane.ID; got != "blocked" {
		// backlog is lane 0; h wraps to the last lane, blocked, regardless of
		// whether blocked holds tickets.
		t.Errorf("lane focus h landed on %q, want blocked (the immediately adjacent lane)", got)
	}
}

// The hint bar names the key in both states.
func TestHintBarNamesTheToggleInBothStates(t *testing.T) {
	m := newTestModel(t, 240, 40)
	out := stripANSI(m.render())
	if !strings.Contains(out, "z thin empty") {
		t.Errorf("hint bar does not offer 'z thin empty':\n%s", out)
	}

	m.key(key("z"))
	out = stripANSI(m.render())
	if !strings.Contains(out, "z widen empty") {
		t.Errorf("hint bar does not offer 'z widen empty' once lanes are thin:\n%s", out)
	}
}

// renderHelp names z at any width.
func TestHelpNamesZ(t *testing.T) {
	for _, w := range []int{40, 100, 200} {
		m := newTestModel(t, w, 50)
		m.mode = modeHelp
		out := stripANSI(m.render())
		if !strings.Contains(out, "z") {
			t.Errorf("help at width %d does not mention z:\n%s", w, out)
		}
	}
}

// Every column is budgeted at the width it is drawn at, so the row of
// columns fills the terminal except for the remainder of an integer division
// (fewer cells than there are full columns) and never runs past it.
func TestBoardRowFillsTerminalWidth(t *testing.T) {
	for _, thin := range []bool{false, true} {
		for _, w := range []int{100, 130, 150, 177, 200, 240, 300} {
			m := newTestModel(t, w, 40)
			if thin {
				m.key(key("z"))
			}
			win := m.boardFit(w)
			full := 0
			for i := win.start; i < win.end; i++ {
				if !win.thin[i] {
					full++
				}
			}
			row := 0
			for _, l := range strings.Split(stripANSI(m.render()), "\n") {
				if strings.Contains(l, "┌") {
					row = max(row, len([]rune(strings.TrimRight(l, " "))))
				}
			}
			if row > w {
				t.Errorf("thin=%v w=%d: column row is %d wide, wider than the terminal", thin, w, row)
			}
			if gap := w - row; full > 0 && gap >= full {
				t.Errorf("thin=%v w=%d: %d cells blank on the right, want fewer than the %d full columns", thin, w, gap, full)
			}
		}
	}
}
