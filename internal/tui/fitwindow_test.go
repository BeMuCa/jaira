package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Every screen, at a terminal small enough to hurt. A line wider than the
// terminal does not spill, it wraps — silently pushing everything below it down
// a line, so one long warning can shove a footer off the bottom of a screen
// that had measured its own height correctly.
func TestEveryScreenFitsTheTerminal(t *testing.T) {
	for _, size := range [][2]int{{40, 12}, {60, 20}, {80, 24}, {150, 32}} {
		w, h := size[0], size[1]
		m := newTestModel(t, w, h)
		m.laneIdx, m.cardIdx = 0, 0
		m.openDetail()
		open := m.detail

		screens := []struct {
			name  string
			setup func()
		}{
			{"board", func() { m.mode, m.detail = modeBoard, nil }},
			{"help", func() { m.mode = modeHelp }},
			{"detail", func() { m.mode, m.detail = modeDetail, open }},
			{"delete", func() { m.mode, m.detail, m.input = modeDelete, open, "" }},
			{"lane focus", func() { m.mode, m.detail = modeLaneFocus, nil }},
			{"compact", func() { m.mode, m.detail = modePipeline, nil }},
			{"projects", func() { m.mode, m.detail = modeProjects, nil }},
			{"move", func() { m.mode, m.detail = modeMove, nil }},
			{"filter", func() { m.mode, m.input = modeFilter, "a very long filter string that keeps going and going and going" }},
			{"create", func() { m.mode, m.input = modeCreate, strings.Repeat("title ", 40) }},
			{"sign-off", func() {
				tk := longTicket()
				tk.Status = "signoff"
				tk.DoDItems[0].Proof = strings.Repeat("a proof line that keeps going ", 6)
				m.mode, m.detail = modeDetail, tk
			}},
			{"edit", func() {
				m.mode, m.detail = modeEdit, open
				m.editIdx = 0
				m.editBuf = strings.Repeat("a very long line being typed into the field ", 4)
			}},
			{"settings", func() { m.mode, m.settingsScreen = modeSettings, newSettingsScreen() }},
			{"lanes", func() { m.mode, m.laneScreen = modeLanes, newLaneScreen(m.store, m.lanes) }},
			{"drop-board", func() {
				m.mode, m.drop = modeDropBoard, newDropBoard(m.store.Root, strings.Repeat("a very long board name ", 4), false)
			}},
			{"message", func() {
				m.mode, m.message, m.isErr = modeMessage, strings.Repeat("a very long message that keeps going and going ", 3), true
			}},
			// follow-up last: it is entered by m.follow being non-nil rather
			// than a mode of its own, and that state must not leak into the
			// screens tested above it.
			{"follow-up", func() {
				m.detail = open
				m.startFollowUp()
			}},
		}
		for _, sc := range screens {
			sc.setup()
			// A warning long enough to wrap on its own, since that is the real
			// case: the board carries lane warnings it did not choose the
			// length of.
			m.warnings = []string{strings.Repeat("a lane file this installation does not understand ", 6)}
			view := m.View()
			out := stripANSI(view.Content)

			lines := strings.Split(out, "\n")
			if len(lines) > h {
				t.Errorf("%dx%d %s: %d lines, want at most %d", w, h, sc.name, len(lines), h)
			}
			for i, l := range lines {
				if got := lipgloss.Width(l); got > w {
					t.Errorf("%dx%d %s: line %d is %d wide, want at most %d: %q", w, h, sc.name, i+1, got, w, l)
				}
			}
		}
	}
}

// Edit, settings and lanes are reached through Model, but Model.View() ends
// in clampBlock, which hard-truncates every line — so measuring them through
// m.View() (as the table above does) can never fail on a wrap bug, only on a
// crash. Browse and home are not reached through Model at all: browse only
// ever hangs off Home, and Home is a separate top-level program with no clamp
// of its own anywhere in its View(). For all five, the only way to measure
// the width a wrap bug would actually produce is to call the render function
// itself, the way wrap_test.go already does for sign-off and the message
// screen.
func TestScreensWithNoOuterClampFitTheirOwnWidth(t *testing.T) {
	long := strings.Repeat("a very long line that keeps going and going and going ", 4)

	for _, w := range []int{34, 100} {
		// edit: the field being typed wraps to the pane instead of truncating.
		{
			m := newTestModel(t, w, 40)
			m.detail = longTicket()
			m.editIdx = 0
			m.editBuf = long
			checkLineWidths(t, w, "edit", stripANSI(m.editBody(max(20, w-2), 40)))
		}

		// settings: each entry's description wraps under its name.
		{
			s := newSettingsScreen()
			checkLineWidths(t, w, "settings", stripANSI(s.render(w, 24)))
		}

		// lanes: promptOf's title line ("name (id)") is the one line on this
		// screen that is neither wrapped nor truncated.
		{
			m := newTestModel(t, w, 40)
			ls := newLaneScreen(m.store, m.lanes)
			found := false
			for i, l := range ls.lanes {
				if l.Prompt != "" {
					l.Name, ls.idx, found = long, i, true
					break
				}
			}
			if !found {
				t.Fatal("no built-in lane carries a prompt to attach the long name to")
			}
			checkLineWidths(t, w, "lanes", stripANSI(ls.render(w, 40)))
		}

		// browse: the current directory path wraps instead of truncating.
		{
			b := newBrowser(t.TempDir())
			b.dir = "/home/someone/git/organisation/team/very-long-repository-name/subdir/more"
			checkLineWidths(t, w, "browse", stripANSI(b.render(w, 30)))
		}

		// home: the launcher's own message wraps under the board list, and
		// nothing downstream of Home.render clamps it further.
		{
			t.Setenv("JAIRA_HOME", t.TempDir())
			h, err := NewHome(nil)
			if err != nil {
				t.Fatal(err)
			}
			h.width, h.height = w, 30
			h.msg = long
			checkLineWidths(t, w, "home", stripANSI(h.render()))
		}
	}
}

func checkLineWidths(t *testing.T, w int, name, out string) {
	t.Helper()
	for i, l := range strings.Split(out, "\n") {
		if n := len([]rune(l)); n > w {
			t.Errorf("width %d %s: line %d is %d wide, want at most %d: %q", w, name, i+1, n, w, l)
		}
	}
}

// The help is the screen this was found on: it is longer than a terminal and
// used to be printed whole, so its own way out scrolled off the bottom.
func TestHelpScrollsAndKeepsItsWayOut(t *testing.T) {
	m := newTestModel(t, 100, 24)
	m.mode = modeHelp
	m.detailScroll = 0

	top := stripANSI(m.View().Content)
	if !strings.Contains(top, "esc back") {
		t.Errorf("the help does not say how to leave:\n%s", top)
	}
	if !strings.Contains(top, "Move around") {
		t.Errorf("the help does not start at the top:\n%s", top)
	}

	// Scrolling reaches what the first screen could not show, and the footer
	// travels with it.
	m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	bottom := stripANSI(m.View().Content)
	if bottom == top {
		t.Error("the help does not scroll")
	}
	if !strings.Contains(bottom, "esc back") {
		t.Errorf("scrolling lost the way out:\n%s", bottom)
	}
	if !strings.Contains(bottom, "Gates are enforced") {
		t.Errorf("the closing note is unreachable:\n%s", bottom)
	}
}
