package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/project"
)

// page names where the program stands, in enough detail to tell two pages of
// one screen apart: the catalogue over the lane editor, a dialog over the
// roadmap, the split beside an open ticket.
func page(m *Model) string {
	p := fmt.Sprintf("mode %d", m.mode)
	switch m.mode {
	case modeDetail:
		p += fmt.Sprintf(" from %d follow %t", m.detailFrom, m.follow != nil)
	case modeLanes:
		if ls := m.laneScreen; ls != nil {
			p += fmt.Sprintf(" catalogue %t confirm %t", ls.catalogueOpen, ls.confirm != nil)
		}
	case modeLayer:
		if l := m.layer; l != nil {
			p += fmt.Sprintf(" stage %d", l.stage)
		}
	}
	return p
}

// quits reports whether a key's command ends the program.
func quits(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// pressAll sends keys one after another through the real dispatch.
func pressAll(m *Model, keys ...string) {
	for _, k := range keys {
		press(m, k)
	}
}

// withBoards gives the board list a second board, so p opens the list
// instead of saying there is nothing to switch to.
func withBoards(t *testing.T) {
	t.Helper()
	t.Setenv("JAIRA_HOME", t.TempDir())
	for _, name := range []string{"one", "two"} {
		project.Remember(makeBoard(t, t.TempDir(), name))
	}
}

// Every page of the board that has no text field, opened the way a person
// opens it. q has to land where esc lands — and, except on the compact view,
// that is a page and not the end of the program. In a text field q is a
// letter, so the filter, a new ticket, the delete prompt, the field editor,
// the follow-up and the roadmap's form are not on this list: esc is their way
// back.
func TestQGoesWhereEscGoesOnEveryPage(t *testing.T) {
	cases := []struct {
		name  string
		open  func(t *testing.T, m *Model)
		at    mode // the page open must reach, so a case cannot pass on a page it never opened
		quits bool
	}{
		{"compact view", func(t *testing.T, m *Model) { press(m, "v") }, modePipeline, true},
		{"opened step", func(t *testing.T, m *Model) { pressAll(m, "v", "enter") }, modeLaneFocus, false},
		{"ticket from an opened step", func(t *testing.T, m *Model) {
			pressAll(m, "v", "enter")
			focusBacklog(t, m, "Investigate flaky logout test")
			press(m, "enter")
		}, modeDetail, false},
		{"open ticket", func(t *testing.T, m *Model) { openFirstTicket(t, m) }, modeDetail, false},
		{"links", func(t *testing.T, m *Model) { press(m, "l") }, modeLinks, false},
		{"links over an open ticket", func(t *testing.T, m *Model) {
			openFirstTicket(t, m)
			press(m, "l")
		}, modeLinks, false},
		{"move picker", func(t *testing.T, m *Model) { press(m, "m") }, modeMove, false},
		{"move picker over an open ticket", func(t *testing.T, m *Model) {
			openFirstTicket(t, m)
			press(m, "m")
		}, modeMove, false},
		{"message", func(t *testing.T, m *Model) { press(m, "r") }, modeMessage, false},
		{"help", func(t *testing.T, m *Model) { press(m, "?") }, modeHelp, false},
		{"tag legend", func(t *testing.T, m *Model) { press(m, "t") }, modeLegend, false},
		{"people picker", func(t *testing.T, m *Model) { press(m, "u") }, modeUsers, false},
		{"milestone picker", func(t *testing.T, m *Model) {
			mkMilestone(t, m, "round", 45)
			reloadOrFail(t, m)
			press(m, "M")
		}, modeMilestones, false},
		{"boards", func(t *testing.T, m *Model) {
			withBoards(t)
			press(m, "p")
		}, modeProjects, false},
		{"remove a board", func(t *testing.T, m *Model) {
			withBoards(t)
			pressAll(m, "p", "x")
		}, modeDropBoard, false},
		{"settings", func(t *testing.T, m *Model) { press(m, "S") }, modeSettings, false},
		{"lane editor", func(t *testing.T, m *Model) { pressAll(m, "S", "enter") }, modeLanes, false},
		{"lane catalogue", func(t *testing.T, m *Model) {
			pressAll(m, "S", "enter")
			for range len(m.laneScreen.lanes) + len(m.laneScreen.available) {
				press(m, "l")
			}
			press(m, "enter")
		}, modeLanes, false},
		{"removing a lane", func(t *testing.T, m *Model) { pressAll(m, "S", "enter", "x") }, modeLanes, false},
		{"default board", func(t *testing.T, m *Model) { pressAll(m, "S", "j", "enter") }, modeDefaultBoard, false},
		{"message over settings", func(t *testing.T, m *Model) {
			press(m, "S")
			m.notify("something happened", false)
		}, modeMessage, false},
		{"roadmap", func(t *testing.T, m *Model) { press(m, "L") }, modeLayer, false},
		{"roadmap: new or connect", func(t *testing.T, m *Model) { pressAll(m, "L", "enter") }, modeLayer, false},
		{"roadmap: connect", func(t *testing.T, m *Model) { pressAll(m, "L", "enter", "c") }, modeLayer, false},
		{"roadmap: open milestone", func(t *testing.T, m *Model) {
			mkPlaced(t, m, "round", 0, nil, titled(t, m, "Refactor auth middleware"))
			reloadOrFail(t, m)
			pressAll(m, "L", "enter")
		}, modeLayer, false},
		{"roadmap: ticket from a milestone", func(t *testing.T, m *Model) {
			mkPlaced(t, m, "round", 0, nil, titled(t, m, "Refactor auth middleware"))
			reloadOrFail(t, m)
			pressAll(m, "L", "enter", "enter")
		}, modeDetail, false},
		{"roadmap: take off", func(t *testing.T, m *Model) {
			mkPlaced(t, m, "round", 0, nil)
			reloadOrFail(t, m)
			pressAll(m, "L", "x")
		}, modeLayer, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			land := func(k string) (string, bool) {
				m := newTestModel(t, 140, 40)
				c.open(t, m)
				if m.mode != c.at {
					t.Fatalf("the page did not open: mode = %v, want %v", m.mode, c.at)
				}
				before := page(m)
				quit := quits(press(m, k))
				if quit != c.quits {
					t.Errorf("%s: quits = %t, want %t", k, quit, c.quits)
				}
				if !quit && page(m) == before {
					t.Errorf("%s left the page as it was: %s", k, before)
				}
				return page(m), quit
			}
			qPage, qQuit := land("q")
			escPage, escQuit := land("esc")
			if qPage != escPage || qQuit != escQuit {
				t.Errorf("q went to %q (quit %t), esc to %q (quit %t)", qPage, qQuit, escPage, escQuit)
			}
		})
	}
}

// The board is the top page: q leaves it — for the launcher when jaira was
// started from one, which the CLI loop decides — and esc only clears a filter.
func TestQOnTheBoardLeavesIt(t *testing.T) {
	m := newTestModel(t, 140, 40)
	if !quits(press(m, "q")) {
		t.Error("q on the board did not end the board")
	}
	m = newTestModel(t, 140, 40)
	m.filter = "auth"
	m.rebuild()
	if quits(press(m, "esc")) || m.filter != "" || m.mode != modeBoard {
		t.Errorf("esc on the board must clear the filter and stay: filter %q, mode %v", m.filter, m.mode)
	}
}

// The launcher's own pages go back to the list with q exactly as with esc,
// and on the list both leave jaira.
func TestQGoesWhereEscGoesOnTheLauncher(t *testing.T) {
	where := func(h *Home) string {
		return fmt.Sprintf("browse %t results %t drop %t board %t quit %t",
			h.browse != nil, h.browse != nil && h.browse.results != nil, h.drop != nil, h.board != nil, h.Quit)
	}
	cases := []struct {
		name string
		keys []string
	}{
		{"list", nil},
		{"add a board", []string{"a"}},
		{"boards a scan found", []string{"a", "s"}},
		{"remove a board", []string{"x"}},
		{"default board", []string{"d"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			land := func(k string) string {
				root := t.TempDir()
				h := newHome(t, []string{makeBoard(t, root, "alpha")}, 120, 30)
				h.startDir = root
				for _, k := range c.keys {
					h.Update(key(k))
				}
				before := where(h)
				h.Update(key(k))
				if where(h) == before {
					t.Errorf("%s left the page as it was: %s", k, before)
				}
				return where(h)
			}
			if q, esc := land("q"), land("esc"); q != esc {
				t.Errorf("q went to %q, esc to %q", q, esc)
			}
		})
	}
}

// The help says what q is, since the same key quitting on one page and going
// back on the next is what this rule replaced.
func TestHelpNamesQAsOnePageBack(t *testing.T) {
	m := newTestModel(t, 140, 60)
	m.mode = modeHelp
	if out := stripANSI(m.render()); !strings.Contains(out, "one page back") {
		t.Errorf("the help does not call q one page back:\n%s", out)
	}
}
