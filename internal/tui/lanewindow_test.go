package tui

import (
	"strings"
	"testing"
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
