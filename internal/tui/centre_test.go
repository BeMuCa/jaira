package tui

import (
	"strings"
	"testing"
)

// His words: the selected lane should be in the middle, so the lanes after it
// stay visible while the window allows. Anchoring it to the right edge meant
// what comes next — where the work goes — was never on screen.
func TestFocusedLaneIsCentred(t *testing.T) {
	m := newTestModel(t, 90, 32)
	if len(m.cols) < 5 {
		t.Skipf("board has only %d lanes; nothing to scroll", len(m.cols))
	}
	const perScreen = 3

	mid := len(m.cols) / 2
	m.laneIdx = mid
	start, end := laneWindow(m.laneIdx, len(m.cols), perScreen)
	if end-start != perScreen {
		t.Fatalf("window %d..%d is not a full row of %d", start, end, perScreen)
	}
	if start >= mid {
		t.Errorf("window %d..%d shows nothing before the focused lane %d", start, end, mid)
	}
	if end <= mid+1 {
		t.Errorf("window %d..%d shows nothing after the focused lane %d — the old edge-anchored behaviour", start, end, mid)
	}
}

// At either end the row stays full rather than padded with blanks, and the
// focused lane is always inside it — a clamp that scrolled past the lane you are
// on would be worse than the edge-anchoring it replaced.
func TestTheWindowStaysFullAndHoldsTheFocusedLane(t *testing.T) {
	m := newTestModel(t, 90, 32)
	if len(m.cols) < 5 {
		t.Skipf("board has only %d lanes", len(m.cols))
	}
	for _, perScreen := range []int{1, 2, 3, 4} {
		for at := range m.cols {
			m.laneIdx = at
			start, end := laneWindow(m.laneIdx, len(m.cols), perScreen)
			if end-start != min(perScreen, len(m.cols)) {
				t.Errorf("perScreen=%d lane=%d: window %d..%d is not full", perScreen, at, start, end)
			}
			if at < start || at >= end {
				t.Errorf("perScreen=%d lane=%d: the focused lane is outside the window %d..%d", perScreen, at, start, end)
			}
		}
	}
}

// More lanes than fit is the only case where any of this matters, so the whole
// board is rendered at a width that forces it.
func TestCentringSurvivesARealRender(t *testing.T) {
	m := newTestModel(t, 90, 32)
	idx := -1
	for i, c := range m.cols {
		if c.lane.ID == "in-progress" {
			idx = i
		}
	}
	if idx < 0 {
		t.Skip("no in-progress lane on this board")
	}
	m.laneIdx, m.cardIdx = idx, 0
	out := stripANSI(m.render())
	if out == "" {
		t.Fatal("empty render")
	}
	if !strings.Contains(out, m.cols[idx].lane.Name) {
		t.Errorf("the focused lane is not on screen:\n%s", out)
	}
	// The lane after it — where this lane's work goes next.
	if next := m.cols[idx+1]; !strings.Contains(out, next.lane.Name) {
		t.Errorf("the lane after the focused one is not on screen:\n%s", out)
	}
	t.Logf("\n%s", out)
}
