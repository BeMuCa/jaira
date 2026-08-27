package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
)

// A key the terminal is too narrow to show is a key the reader does not know
// exists — hints wrap onto more lines instead of dropping the tail.
func TestWrapHintsKeepsEveryItem(t *testing.T) {
	items := []string{"enter open", "v compact", "z thin empty", "n new", "m move", "S settings", "/ filter", "? help", "q quit"}
	lines := wrapHints(items, 24)
	if len(lines) < 2 {
		t.Fatalf("width 24 must force wrapping, got %d line(s): %v", len(lines), lines)
	}
	joined := strings.Join(lines, "\n")
	for _, it := range items {
		if !strings.Contains(joined, it) {
			t.Errorf("wrapped hints lost %q:\n%s", it, joined)
		}
	}
	for _, l := range lines {
		if len([]rune(l)) > 24 {
			t.Errorf("line wider than the terminal: %q", l)
		}
	}
}

// The board on a narrow terminal must still name every key, and the whole
// render must still fit the window — the wrapped bar eats into the columns'
// height rather than pushing the board past the bottom.
func TestNarrowBoardShowsAllKeysAndFits(t *testing.T) {
	m := newTestModel(t, 30, 24)
	out := stripANSI(m.render())
	for _, key := range []string{"enter open", "v compact", "z thin empty", "n new", "m move", "S settings", "/ filter", "? help", "q quit"} {
		if !strings.Contains(out, key) {
			t.Errorf("narrow board lost the %q hint", key)
		}
	}
	if got := len(strings.Split(out, "\n")); got > 24 {
		t.Errorf("board rendered %d lines into a 24-line terminal", got)
	}
}

// The open ticket's footer wraps the same way, and the pane still fits.
func TestNarrowDetailShowsAllKeysAndFits(t *testing.T) {
	m := newTestModel(t, 34, 20)
	m.detail = longTicket()
	out := stripANSI(m.renderDetail())
	for _, key := range []string{"e fields", "E editor", "y copy id", "m move", "esc back"} {
		if !strings.Contains(out, key) {
			t.Errorf("narrow detail footer lost %q:\n%s", key, out)
		}
	}
	if got := len(strings.Split(out, "\n")); got > 20 {
		t.Errorf("detail rendered %d lines into a 20-line terminal", got)
	}
}

// Every catalogue lane that is not on the board shows as its own dimmed
// column, and enter installs it — no detour through the '+' column.
func TestAvailableLaneShowsAndAddsOnEnter(t *testing.T) {
	s, catalogue := laneTestStore(t)
	// A project lane directory makes itself authoritative, so the catalogue
	// lane below is present on the machine but not on this board.
	builtins, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lane.Export(builtins[0], lane.ProjectLanesDir(s.Root), false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(catalogue, "extra.md"), []byte(`---
id: extra
name: Extra
after: human
precedence: 41
---
`), 0o644); err != nil {
		t.Fatal(err)
	}
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)

	if len(ls.available) != 1 || ls.available[0].ID != "extra" {
		t.Fatalf("available = %v, want exactly [extra]", ls.available)
	}
	if out := stripANSI(ls.render(120, 40)); !strings.Contains(out, "not on board") {
		t.Errorf("the missing lane is not visible on the screen:\n%s", out)
	}

	ls.idx = len(ls.lanes) // the one available lane
	ls.key("enter")
	if indexOfLane(ls.lanes, "extra") < 0 {
		t.Fatalf("enter did not install the available lane: %v", ls.msg)
	}
	if len(ls.available) != 0 {
		t.Errorf("extra is installed but still listed as available: %v", ls.available)
	}
}

// E edits the lane under the cursor. A built-in has no file, so the edit
// writes a catalogue override copy first and opens that.
func TestEditBuiltinLaneWritesTheOverrideCopy(t *testing.T) {
	s, catalogue := laneTestStore(t)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	ls.idx = indexOfLane(ls.lanes, "review")
	if ls.idx < 0 {
		t.Fatal("no review lane on the default board")
	}

	t.Setenv("EDITOR", "true")
	if cmd := ls.editLane(); cmd == nil {
		t.Fatalf("editLane returned no command: %v", ls.msg)
	}
	copy := filepath.Join(catalogue, "review.md")
	if _, err := os.Stat(copy); err != nil {
		t.Fatalf("no override copy was written at %s: %v", copy, err)
	}

	// A second E must reuse the copy rather than refusing to overwrite it.
	if cmd := ls.editLane(); cmd == nil {
		t.Fatalf("second editLane refused: %v", ls.msg)
	}
}

// The lane settings screen wraps its columns onto further rows on a narrow
// terminal. It used to clip them at the right edge with no indicator, which
// read as review, human review and done simply not existing.
func TestLaneScreenShowsEveryLaneOnANarrowTerminal(t *testing.T) {
	s, _ := laneTestStore(t)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)

	out := stripANSI(ls.render(40, 40))
	for _, l := range set.Lanes {
		if !strings.Contains(out, truncateForTest(l.Name, laneColWidth-2)) {
			t.Errorf("lane %q is not visible at width 40:\n%s", l.Name, out)
		}
	}
}

func truncateForTest(s string, w int) string { return truncate(s, w) }
