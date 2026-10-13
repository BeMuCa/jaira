package tui

import (
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/project"
)

// press sends one key through the real dispatch, so a test cannot pass against
// a handler the keypress never reaches.
func press(m *Model, s string) tea.Cmd {
	k := tea.KeyPressMsg{Text: s}
	if r := []rune(s); len(r) == 1 {
		k.Code = r[0]
	} else {
		switch s {
		case "esc":
			k.Code = tea.KeyEscape
			k.Text = ""
		case "enter":
			k.Code = tea.KeyEnter
			k.Text = ""
		}
	}
	_, cmd := m.key(k)
	return cmd
}

// laneOrder finds a lane's index in m.lanes.Lanes, which is what the move
// picker's moveTarget indexes.
func laneOrder(t *testing.T, m *Model, id string) int {
	t.Helper()
	for i, l := range m.lanes.Lanes {
		if l.ID == id {
			return i
		}
	}
	t.Fatalf("no lane %q in fixture", id)
	return -1
}

// openTicketInLaneFocus puts an open ticket on screen with lane focus behind it.
func openTicketInLaneFocus(t *testing.T, m *Model) {
	t.Helper()
	m.mode = modeLaneFocus
	focusBacklog(t, m, "Investigate flaky logout test")
	m.openDetail()
	if m.mode != modeDetail || m.detail == nil {
		t.Fatal("detail did not open")
	}
}

// A dialog opened over an open ticket belongs to that ticket, not to the board.
func TestMovePickerReturnsToTheOpenTicket(t *testing.T) {
	m := newTestModel(t, 120, 34)
	openTicketInLaneFocus(t, m)

	press(m, "m")
	if m.mode != modeMove {
		t.Fatalf("m did not open the move picker: mode = %v", m.mode)
	}
	press(m, "esc")

	if m.mode != modeDetail {
		t.Errorf("esc left the move picker on mode %v, want modeDetail", m.mode)
	}
	if m.detail == nil {
		t.Error("the open ticket was closed by the move picker")
	}
}

func TestMessageDismissesToThePageItWasRaisedOn(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, m *Model)
		want  mode
	}{
		{"lane focus", func(t *testing.T, m *Model) { m.mode = modeLaneFocus }, modeLaneFocus},
		{"compact view", func(t *testing.T, m *Model) { m.mode = modePipeline }, modePipeline},
		{"open ticket", openTicketInLaneFocus, modeDetail},
		{"board", func(t *testing.T, m *Model) { m.mode = modeBoard }, modeBoard},
		{"settings", func(t *testing.T, m *Model) {
			m.settingsScreen = newSettingsScreen(m.store.Root)
			m.mode = modeSettings
		}, modeSettings},
		{"lane editor", func(t *testing.T, m *Model) {
			m.settingsScreen = newSettingsScreen(m.store.Root)
			m.laneScreen = newLaneScreen(m.store, m.lanes)
			m.mode = modeLanes
		}, modeLanes},
		{"default board", func(t *testing.T, m *Model) {
			db, err := lane.LoadDefaultBoard()
			if err != nil {
				t.Fatal(err)
			}
			m.settingsScreen = newSettingsScreen(m.store.Root)
			m.board = newDefaultBoardScreen(m.lanes, db)
			m.mode = modeDefaultBoard
		}, modeDefaultBoard},
		{"boards", func(t *testing.T, m *Model) {
			m.projects = []project.Project{{Root: m.store.Root, Name: "current"}}
			m.mode = modeProjects
		}, modeProjects},
		// The legend keeps its own way back in returnTo, which a message over
		// it would overwrite, so it still dismisses to the board.
		{"legend", func(t *testing.T, m *Model) { press(m, "t") }, modeBoard},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestModel(t, 120, 34)
			c.setup(t, m)
			m.notify("something happened", false)
			if m.mode != modeMessage {
				t.Fatalf("notify did not show a message: mode = %v", m.mode)
			}
			// The top row sits above the box, so it shows what the message
			// covers — which has to be the page it goes back to.
			under := strings.SplitN(stripANSI(m.render()), "\n", 2)[0]
			press(m, "esc")
			if m.mode != c.want {
				t.Errorf("dismissed to %v, want %v", m.mode, c.want)
			}
			if back := strings.SplitN(stripANSI(m.render()), "\n", 2)[0]; back != under {
				t.Errorf("the message covered %q, but dismissing it showed %q", under, back)
			}
		})
	}
}

// A board that will not open says so over the list it was chosen from, and
// the message goes back to that list rather than to a board nobody switched to.
func TestABoardThatWillNotOpenSaysSoOverTheList(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: filepath.Join(t.TempDir(), "gone"), Name: "gone"},
	}
	m.mode = modeProjects
	m.projIdx = 1
	before := m.store

	press(m, "enter")
	if m.mode != modeMessage {
		t.Fatalf("a failed switch showed no message: mode = %v", m.mode)
	}
	if m.store != before {
		t.Error("a failed switch swapped the store")
	}
	press(m, "q")
	if m.mode != modeProjects {
		t.Errorf("dismissed to %v, want modeProjects", m.mode)
	}
}

// A message raised by the move picker belongs to the page behind the picker.
func TestMessageFromTheMovePickerReturnsBehindIt(t *testing.T) {
	m := newTestModel(t, 120, 34)
	openTicketInLaneFocus(t, m)
	press(m, "m")

	m.notify("the gate said no", true)
	press(m, "esc")

	if m.mode != modeDetail {
		t.Errorf("dismissed to %v, want modeDetail — the page behind the picker", m.mode)
	}
}

// Moving a ticket from an open ticket leaves it open, showing the lane it just
// moved to rather than the one it left.
func TestMoveFromAnOpenTicketLeavesItOpenAndFresh(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modeBoard
	// The human lane is the one the fixture can legally leave: a checkpoint is
	// exempt from the ownership rail and its exit needs no outcome.
	focusTicket(t, m, "human", "Decide on cookie SameSite policy")
	m.openDetail()
	id := m.detail.ID

	press(m, "m")
	m.moveTarget = laneOrder(t, m, "todo")
	press(m, "enter")

	if m.mode != modeDetail {
		t.Fatalf("the move closed the ticket: mode = %v", m.mode)
	}
	if m.detail == nil {
		t.Fatal("no ticket open after the move")
	}
	if m.detail.ID != id {
		t.Errorf("a different ticket is open: %s, want %s", m.detail.ID, id)
	}
	if m.detail.Status != "todo" {
		t.Errorf("the open ticket still shows %q, want todo", m.detail.Status)
	}
}

// A reload underneath a dialog must not jump lanes either.
func TestOverlaysOverLaneFocusStillHoldTheLane(t *testing.T) {
	for _, c := range []struct {
		name  string
		setup func(m *Model)
	}{
		{"move picker", func(m *Model) { press(m, "m") }},
		{"message", func(m *Model) { m.notify("note", false) }},
	} {
		t.Run(c.name, func(t *testing.T) {
			m := newTestModel(t, 120, 34)
			openTicketInLaneFocus(t, m)
			moved := m.detail
			c.setup(m)

			touchOnDisk(t, m, moved.ID, "done")

			if want := laneIndex(t, m, "backlog"); m.laneIdx != want {
				t.Errorf("lane jumped under the overlay: laneIdx = %d, want %d", m.laneIdx, want)
			}
		})
	}
}
