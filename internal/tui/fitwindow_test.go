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
