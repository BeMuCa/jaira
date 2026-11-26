package tui

import (
	"strings"
	"testing"
)

// In the ten-lane fixture review (6) and signoff (7) are the only two empty
// lanes side by side; brainstorm (1), pre-process (3) and blocked (9) are
// empty on their own.

// thinColumnRows is the strip's text, one entry per row, border and padding
// stripped; a blank row (the space in a name) falls away.
func thinColumnRows(s string) []string {
	var rows []string
	for _, l := range strings.Split(stripANSI(s), "\n") {
		l = strings.Trim(l, "│ ")
		if l != "" && !strings.ContainsAny(l, "┌└") {
			rows = append(rows, l)
		}
	}
	return rows
}

// Two neighbouring empty lanes are one strip under z, no wider than one thin
// lane; the board draws one border fewer than it has lanes.
func TestZFoldsNeighbouringEmptyLanesIntoOneStrip(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	win := m.boardFit(240)
	if got := win.fold[6]; got != 8 {
		t.Fatalf("fold[review] = %d, want 8 (review and signoff in one strip); fold = %v", got, win.fold)
	}
	out := stripANSI(m.render())
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "┌") {
			if got, want := strings.Count(l, "┌"), len(m.cols)-1; got != want {
				t.Errorf("%d boxes on the board, want %d (ten lanes, two of them in one strip):\n%s", got, want, l)
			}
			break
		}
	}
	strip := m.renderThinColumn(6, 8, 30)
	if w := len([]rune(strings.Split(stripANSI(strip), "\n")[0])); w != thinColWidth {
		t.Errorf("folded strip is %d cells wide, want %d like a single thin lane", w, thinColWidth)
	}
}

// The strip says what it holds: both names down the column with a '─' row
// between them, and only their initials when the names do not fit.
func TestFoldedStripNamesItsLanes(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	tall := strings.Join(thinColumnRows(m.renderThinColumn(6, 8, 30)), "")
	if want := "Review─HumanReview"; tall != want {
		t.Errorf("tall strip reads %q, want %q", tall, want)
	}
	short := strings.Join(thinColumnRows(m.renderThinColumn(6, 8, 5)), "")
	if want := "RH"; short != want {
		t.Errorf("short strip reads %q, want the initials %q", short, want)
	}
}

// An empty lane between two full ones keeps a strip of its own, and without z
// nothing folds at all.
func TestSingleEmptyLaneStaysItsOwnStrip(t *testing.T) {
	m := newTestModel(t, 240, 40)
	if win := m.boardFit(240); len(win.fold) != 0 {
		t.Errorf("lanes folded without z: %v", win.fold)
	}
	m.key(key("z"))
	win := m.boardFit(240)
	for _, i := range []int{1, 3, 9} {
		if !win.thin[i] {
			t.Errorf("lane %q is not thin", m.cols[i].lane.ID)
		}
		if _, ok := win.fold[i]; ok {
			t.Errorf("lane %q, alone between full lanes, was folded", m.cols[i].lane.ID)
		}
	}
	if len(win.fold) != 1 {
		t.Errorf("fold = %v, want only review+signoff", win.fold)
	}
}

// The lane holding the cursor never folds: h and l over empty lanes must still
// show where the cursor is.
func TestCursorLaneIsNotFolded(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	m.laneIdx = 6
	if win := m.boardFit(240); len(win.fold) != 0 {
		t.Errorf("fold = %v with the cursor on review, want review and signoff drawn apart", win.fold)
	}
}

// The width a fold saves goes to the lanes with tickets: five full columns
// share the terminal less four strips, not five.
func TestFoldGivesItsWidthToFullLanes(t *testing.T) {
	const w = 204
	m := newTestModel(t, w, 40)
	m.key(key("z"))
	win := m.boardFit(w)
	if want := (w - 4*thinColWidth) / 5; win.colW != want {
		t.Errorf("colW = %d, want %d", win.colW, want)
	}
	if unfolded := (w - 5*thinColWidth) / 5; win.colW <= unfolded {
		t.Errorf("colW = %d, no wider than the %d of five separate strips", win.colW, unfolded)
	}
}

// An unknown lane gives its own initial in a short strip, not the "?" it is
// prefixed with in full — two of them would otherwise read the same.
func TestFoldedStripInitialOfUnknownLane(t *testing.T) {
	m := newTestModel(t, 240, 40)
	m.key(key("z"))
	m.cols[6].lane.Unknown = true
	m.cols[7].lane.Unknown = true
	if got := strings.Join(thinColumnRows(m.renderThinColumn(6, 8, 5)), ""); got != "RH" {
		t.Errorf("short strip of two unknown lanes reads %q, want \"RH\"", got)
	}
}
