package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
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

// z hides the five lanes holding no tickets and keeps the five that do, and
// says how many went and how to bring them back.
func TestZHidesEmptyLanesAndNotices(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	out := stripANSI(m.render())

	for _, name := range emptyLaneNames {
		if strings.Contains(out, name) {
			t.Errorf("empty lane %q still rendered after z:\n%s", name, out)
		}
	}
	for _, name := range filledLaneNames {
		if !strings.Contains(out, name) {
			t.Errorf("filled lane %q disappeared after z:\n%s", name, out)
		}
	}
	if !strings.Contains(out, "5 lane(s) hidden") {
		t.Errorf("notice does not name the count 5 and the word 'hidden':\n%s", out)
	}
	if !strings.Contains(out, "z to show") {
		t.Errorf("notice does not name z as the way back:\n%s", out)
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
	if strings.Contains(out, "hidden") {
		t.Errorf("notice survived the second z:\n%s", out)
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

// With the toggle on, l must never land the cursor on a hidden lane — it
// steps to the next lane that holds tickets.
func TestLWithToggleOnSkipsHiddenLanes(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.laneIdx = laneIndex(t, m, "backlog")
	m.key(key("z"))

	m.moveLane(1)
	if got := m.cols[m.laneIdx].lane.ID; got != "todo" {
		t.Errorf("l from backlog with empty lanes hidden landed on %q, want todo (skipping brainstorm)", got)
	}
}

// The focused lane renders even when it holds nothing, because the cursor
// must never point at something invisible — and the notice counts only the
// lanes that actually went, not the one kept for the cursor. The cursor lands
// on the empty lane the way another shell's move would (a direct index
// assignment, bypassing moveLane's own skip), the case boardFit's keep-rule
// exists for.
func TestFocusedEmptyLaneStaysVisibleAndIsNotCountedHidden(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z")) // toggle on first, while the cursor sits on non-empty backlog
	m.laneIdx = laneIndex(t, m, "blocked")
	out := stripANSI(m.render())

	if !strings.Contains(out, "Blocked") {
		t.Errorf("focused empty lane disappeared:\n%s", out)
	}
	if !strings.Contains(out, "4 lane(s) hidden") {
		t.Errorf("notice should count 4 hidden lanes (blocked kept for the cursor):\n%s", out)
	}
}

// Turning the toggle on while parked on an empty lane moves the cursor to the
// nearest lane that holds tickets, searching right first — otherwise the
// toggle would leave the cursor stuck on the one lane it was pressed to get
// rid of, kept on screen forever by boardFit's own keep-rule.
func TestTogglingOnMovesCursorOffEmptyLaneSearchingRightFirst(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.laneIdx = laneIndex(t, m, "pre-process") // empty; in-progress (filled) is immediately right
	m.key(key("z"))

	got := m.cols[m.laneIdx].lane.ID
	if got != "in-progress" {
		t.Errorf("cursor landed on %q after toggling on from an empty lane, want in-progress (search right)", got)
	}
}

// When nothing sits to the right, the search continues left.
func TestTogglingOnMovesCursorOffEmptyLaneSearchingLeftAsFallback(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.laneIdx = laneIndex(t, m, "blocked") // empty, and the last lane: nothing to its right
	m.key(key("z"))

	got := m.cols[m.laneIdx].lane.ID
	if got != "done" {
		t.Errorf("cursor landed on %q after toggling on from the last (empty) lane, want done (search left fallback)", got)
	}
}

// The toggle is a board preference: laneFocusKey's h/l must step to the
// immediately adjacent lane regardless of hideEmpty, since that view never
// shows the toggle at all.
func TestLaneFocusIgnoresTheToggle(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.mode = modeLaneFocus
	m.laneIdx = laneIndex(t, m, "backlog")
	m.hideEmpty = true

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
	if !strings.Contains(out, "z hide empty") {
		t.Errorf("hint bar does not offer 'z hide empty':\n%s", out)
	}

	m.key(key("z"))
	out = stripANSI(m.render())
	if !strings.Contains(out, "z show empty") {
		t.Errorf("hint bar does not offer 'z show empty' once lanes are hidden:\n%s", out)
	}
}

// renderHelp names z at any width.
func TestHelpNamesZ(t *testing.T) {
	for _, w := range []int{40, 100, 200} {
		m := newTestModel(t, w, 40)
		m.mode = modeHelp
		out := stripANSI(m.render())
		if !strings.Contains(out, "z") {
			t.Errorf("help at width %d does not mention z:\n%s", w, out)
		}
	}
}

// scrolloff: the margin shrinks with the window rather than demanding room
// that does not exist. Free of any Model, so a later reader arguing about the
// margin can settle it from this table without opening the renderer.
func TestScrolloffTable(t *testing.T) {
	cases := []struct {
		perScreen int
		want      int
	}{
		{5, 2}, {6, 2}, {20, 2},
		{3, 1}, {4, 1},
		{1, 0}, {2, 0},
	}
	for _, c := range cases {
		if got := scrolloff(c.perScreen); got != c.want {
			t.Errorf("scrolloff(%d) = %d, want %d", c.perScreen, got, c.want)
		}
	}
}

// windowStart: the statement of what the scroll rule is. total 10, perScreen
// 5 unless said otherwise, so scrolloff is 2 and the rule is exact centring.
func TestWindowStartTable(t *testing.T) {
	const total, perScreen = 10, 5
	cases := []struct {
		name             string
		total, perScreen int
		start, cursor    int
		want             int
	}{
		{"cursor 5 from start 0 centres", total, perScreen, 0, 5, 3},
		{"cursor 5 from start 3 (already correct) stays", total, perScreen, 3, 5, 3},
		{"cursor 5 from start 9 (out of range) centres", total, perScreen, 9, 5, 3},
		{"cursor 0: nothing further left to scroll to", total, perScreen, 4, 0, 0},
		{"cursor 1: same left edge", total, perScreen, 4, 1, 0},
		{"cursor 9: flush with the right edge, no padding", total, perScreen, 4, 9, 5},
		{"cursor 8: same right edge", total, perScreen, 4, 8, 5},
		{"a start already satisfying the margin is returned untouched", total, perScreen, 2, 4, 2},
		{"negative start is pulled back inside", total, perScreen, -3, 0, 0},
		{"start past total-perScreen is pulled back inside", total, perScreen, 9, 9, 5},
		{"everything fits: perScreen == total, no window", 10, 10, 7, 5, 0},
		{"everything fits: perScreen > total, no window", 10, 12, 7, 5, 0},
		// perScreen 7 (scrolloff still 2): cursors 2, 3 and 4 all leave a
		// start of 0 alone — the steadiness the whole change is for.
		{"steady band: cursor 2 leaves start 0", 10, 7, 0, 2, 0},
		{"steady band: cursor 3 leaves start 0", 10, 7, 0, 3, 0},
		{"steady band: cursor 4 leaves start 0", 10, 7, 0, 4, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := windowStart(c.total, c.perScreen, c.start, c.cursor); got != c.want {
				t.Errorf("windowStart(%d, %d, %d, %d) = %d, want %d",
					c.total, c.perScreen, c.start, c.cursor, got, c.want)
			}
		})
	}
}

// At 240 columns (everything fits) the render carries no window indicator at
// all, whether or not empty lanes are hidden — a wide terminal is untouched.
func TestWideBoardHasNoWindowIndicator(t *testing.T) {
	m := newTestModel(t, 240, 40)
	out := stripANSI(m.render())
	if strings.Contains(out, "‹") || strings.Contains(out, "›") {
		t.Errorf("240-column board shows a window indicator it should not need:\n%s", out)
	}
}

// At 125x40 (the user's own terminal) perScreen is 5, and the indicator names
// the window: positions 1-5 out of 10, plus the h/l keys.
func TestNarrowBoardNamesTheWindow(t *testing.T) {
	m := newTestModel(t, 125, 40)
	out := stripANSI(m.render())
	if !strings.Contains(out, "lanes 1–5 of 10") {
		t.Errorf("125-column board does not name its window (positions 1-5 of 10):\n%s", out)
	}
	if !strings.Contains(out, "h") || !strings.Contains(out, "l") {
		t.Errorf("125-column board's window indicator does not name the h/l keys:\n%s", out)
	}
}

// Walking l from the leftmost lane at 125x40 keeps the columns still while
// the cursor sits inside its two-lane margin, and advances the window by
// exactly one lane at a time once the margin is exhausted.
func TestWalkingLKeepsTheTwoLaneMargin(t *testing.T) {
	m := newTestModel(t, 125, 40)
	m.render() // establish the window at start=0 for the initial cursor

	starts := []int{}
	for i := 0; i < 5; i++ {
		m.key(key("l"))
		m.render()
		starts = append(starts, m.laneStart)
	}
	// The first two presses move no columns at all: the cursor is still
	// inside its margin against the left end.
	if starts[0] != 0 || starts[1] != 0 {
		t.Fatalf("first two l presses moved the window: %v, want [0 0 ...]", starts)
	}
	// The third, fourth and fifth each advance the window by exactly one.
	want := []int{0, 0, 1, 2, 3}
	if starts[0] != want[0] || starts[1] != want[1] || starts[2] != want[2] || starts[3] != want[3] || starts[4] != want[4] {
		t.Fatalf("window starts = %v, want %v", starts, want)
	}

	if want := laneIndex(t, m, "human"); m.laneIdx != want {
		t.Fatalf("cursor landed at index %d, want human at %d (lane 6 of 10)", m.laneIdx, want)
	}
	out := stripANSI(m.render())
	if !strings.Contains(out, "lanes 4–8 of 10") {
		t.Errorf("window should sit at positions 4-8 (two either side of lane 6):\n%s", out)
	}
	// The lanes that scrolled off are absent from the render.
	if strings.Contains(out, "Backlog") {
		t.Errorf("Backlog should have scrolled off by now:\n%s", out)
	}
}

// Walking to the very last lane clamps the window at the right edge with no
// padding; walking back to the first mirrors it at the left edge.
func TestWalkingToEitherEndClampsWithoutPadding(t *testing.T) {
	m := newTestModel(t, 125, 40)
	m.render()
	for i := 0; i < 9; i++ {
		m.key(key("l"))
	}
	out := stripANSI(m.render())
	if !strings.Contains(out, "lanes 6–10 of 10") {
		t.Errorf("window should clamp at positions 6-10 at the last lane:\n%s", out)
	}
	if want := laneIndex(t, m, "blocked"); m.laneIdx != want {
		t.Fatalf("cursor should sit on the last lane (blocked): laneIdx=%d, want %d", m.laneIdx, want)
	}
	if strings.Contains(out, "Backlog") {
		t.Errorf("no lane before the window should render at the right edge:\n%s", out)
	}

	for i := 0; i < 9; i++ {
		m.key(key("h"))
	}
	out = stripANSI(m.render())
	if !strings.Contains(out, "lanes 1–5 of 10") {
		t.Errorf("window should clamp back at positions 1-5 at the first lane:\n%s", out)
	}
	if m.laneIdx != laneIndex(t, m, "backlog") {
		t.Fatalf("cursor should be back on backlog: laneIdx=%d", m.laneIdx)
	}
}

// At 170x40 (perScreen 7, scrolloff 2), moving the cursor through the
// interior leaves m.laneStart alone and renders the same lane names in the
// same order — the steadiness the whole change is for, checked on both the
// field and the actual render so neither could pass while the other fails.
func TestSteadyInteriorAt170Columns(t *testing.T) {
	m := newTestModel(t, 170, 40)
	m.laneIdx = laneIndex(t, m, "todo") // index 2
	first := stripANSI(m.render())
	if m.laneStart != 0 {
		t.Fatalf("laneStart = %d before any move, want 0", m.laneStart)
	}

	m.key(key("l")) // -> pre-process, index 3
	second := stripANSI(m.render())
	if m.laneStart != 0 {
		t.Fatalf("laneStart = %d after moving to lane 4, want 0", m.laneStart)
	}

	m.key(key("l")) // -> in-progress, index 4
	third := stripANSI(m.render())
	if m.laneStart != 0 {
		t.Fatalf("laneStart = %d after moving to lane 5, want 0", m.laneStart)
	}

	for _, name := range []string{"Backlog", "Brainstorm", "Todo", "Pre-process", "Implementing", "HITL", "Review"} {
		for _, out := range []string{first, second, third} {
			if !strings.Contains(out, name) {
				t.Errorf("lane %q missing from a steady-interior render:\n%s", name, out)
			}
		}
	}
}

// At 125x40 with the toggle on, the five lanes that hold tickets all fit —
// no window part in the notice at all, though the hidden-lanes half (Task 1)
// still names what the toggle removed. This is the user's actual terminal:
// the empty-lane toggle alone solves it.
func TestNarrowBoardWithToggleOnHasNoWindow(t *testing.T) {
	m := newTestModel(t, 125, 40)
	m.key(key("z"))
	out := stripANSI(m.render())

	if strings.Contains(out, "‹") || strings.Contains(out, "›") {
		t.Errorf("125-column board with empty lanes hidden should need no window:\n%s", out)
	}
	if !strings.Contains(out, "5 lane(s) hidden") {
		t.Errorf("hidden-lanes notice should still name what the toggle removed:\n%s", out)
	}
	for _, name := range filledLaneNames {
		if !strings.Contains(out, name) {
			t.Errorf("lane %q missing even though all five should fit:\n%s", name, out)
		}
	}
}

// At exactly minColWidth a card still shows its whole six-cell handle with no
// ellipsis and at least one intact flag word — the floor means what it
// claims to mean.
func TestMinColWidthKeepsHandleAndAFlagIntact(t *testing.T) {
	m := newTestModel(t, 240, 40)
	idx := laneIndex(t, m, "backlog")
	tk := m.cols[idx].tickets[0]
	handle := ticket.Handle(tk.ID)

	out := stripANSI(m.renderColumn(idx, minColWidth, 20))
	if !strings.Contains(out, handle) {
		t.Errorf("handle %q not intact at minColWidth:\n%s", handle, out)
	}
	if !strings.Contains(out, "spec") {
		t.Errorf("flag word 'spec' not intact at minColWidth:\n%s", out)
	}
}

// A resize re-clamps the stored window on the next frame with no resize
// handler of its own: growing past every lane drops the indicator, and
// shrinking back re-applies it with the cursor's margin honoured. A stored
// start left out of range by the shrink must not slice past win.cols.
func TestResizeReclampsInBothDirectionsWithNoHandler(t *testing.T) {
	m := newTestModel(t, 125, 40)
	m.render()
	for i := 0; i < 8; i++ {
		m.key(key("l"))
	}
	m.render()
	if m.laneStart == 0 {
		t.Fatal("expected a non-zero window before resizing, so growing/shrinking is a real exercise")
	}

	m.width = 240
	out := stripANSI(m.render())
	if strings.Contains(out, "‹") {
		t.Errorf("240-column render after growing still shows a window indicator:\n%s", out)
	}
	for _, name := range append(append([]string{}, emptyLaneNames...), filledLaneNames...) {
		if !strings.Contains(out, name) {
			t.Errorf("lane %q missing after growing to 240:\n%s", name, out)
		}
	}

	m.width = 125
	out = stripANSI(m.render()) // must not panic slicing win.cols with a stale start
	if !strings.Contains(out, "‹") {
		t.Errorf("125-column render after shrinking back lost its window indicator:\n%s", out)
	}
	if m.laneStart < 0 || m.laneStart > 10-5 {
		t.Errorf("laneStart %d out of [0, 5] after resizing back down", m.laneStart)
	}
}

// A ticket moved by another shell, one the user was not looking at, must not
// scroll the window — only the ticket the cursor follows can do that.
func TestExternalMoveOfUnrelatedTicketDoesNotScrollTheWindow(t *testing.T) {
	m := newTestModel(t, 125, 40)
	for i := 0; i < 8; i++ {
		m.key(key("l"))
	}
	m.render()
	before := m.laneStart

	other := backlogTicket(t, m, "Rate limit the login endpoint")
	touchOnDisk(t, m, other.ID, "backlog") // bumps updated_at, reorders within backlog
	m.render()

	if m.laneStart != before {
		t.Errorf("window moved for a ticket the user was not looking at: %d -> %d", before, m.laneStart)
	}
}

// A ticket moved by another shell that IS the selected one scrolls the
// window to include it, because following the selected ticket is the
// board's own stated behaviour and a cursor pointing off-screen would be
// worse.
func TestExternalMoveOfSelectedTicketScrollsToInclude(t *testing.T) {
	m := newTestModel(t, 125, 40)
	moved := focusBacklog(t, m, "Investigate flaky logout test")
	m.render()
	if m.laneStart != 0 {
		t.Fatalf("laneStart = %d before the move, want 0", m.laneStart)
	}

	touchOnDisk(t, m, moved.ID, "blocked")
	out := stripANSI(m.render())

	if want := laneIndex(t, m, "blocked"); m.laneIdx != want {
		t.Fatalf("cursor did not follow the moved ticket: laneIdx=%d, want blocked at %d", m.laneIdx, want)
	}
	if m.laneStart == 0 {
		t.Errorf("window did not scroll to include the moved, selected ticket's lane")
	}
	if !strings.Contains(out, "Blocked") {
		t.Errorf("blocked lane is not on screen after following the selected ticket:\n%s", out)
	}
}
