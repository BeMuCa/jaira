package tui

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/lane"
)

// newDefaultBoardTestScreen builds a screen against an isolated catalogue and
// default board file, so a test never touches the real developer's machine.
func newDefaultBoardTestScreen(t *testing.T) *defaultBoardScreen {
	t.Helper()
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))
	set, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	board, err := lane.LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	return newDefaultBoardScreen(set, board)
}

// TestDefaultBoardScreenDefaultsToBuiltinsChecked asserts a first visit, with
// no default board yet, shows every built-in already ticked — an absent
// selection means the built-ins, and the screen must not misrepresent that
// as "nothing chosen".
func TestDefaultBoardScreenDefaultsToBuiltinsChecked(t *testing.T) {
	d := newDefaultBoardTestScreen(t)
	for _, l := range d.set.Lanes {
		if !d.lanes[l.ID] {
			t.Errorf("built-in %q must start ticked, got unticked", l.ID)
		}
	}
}

// TestDefaultBoardScreenTogglesAndSaves asserts unticking a lane and ticking
// an option, then saving, is what LoadDefaultBoard reads back.
func TestDefaultBoardScreenTogglesAndSaves(t *testing.T) {
	d := newDefaultBoardTestScreen(t)

	// Untick the first lane.
	droppedID := d.set.Lanes[0].ID
	d.idx = 0
	d.focus = 0
	d.key(" ")
	if d.lanes[droppedID] {
		t.Fatal("space did not untick the selected lane")
	}

	// Tick the first option, if the built-ins install one.
	opts := d.set.Options()
	if len(opts) == 0 {
		t.Fatal("expected the built-ins to declare at least one option")
	}
	d.key("tab")
	if d.focus != 1 {
		t.Fatal("tab must switch focus to the options list")
	}
	d.idx = 0
	d.key(" ")
	if !d.options[opts[0]] {
		t.Fatal("space did not tick the selected option")
	}

	done, cmd := d.key("s")
	if done || cmd != nil {
		t.Fatalf("save must not finish the screen or return a command, got done=%v cmd=%v", done, cmd)
	}
	if d.isErr {
		t.Fatalf("save failed: %s", d.msg)
	}

	saved, err := lane.LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range saved.Lanes {
		if id == droppedID {
			t.Errorf("saved board still names %q, want it dropped", droppedID)
		}
	}
	found := false
	for _, o := range saved.Options {
		if o == opts[0] {
			found = true
		}
	}
	if !found {
		t.Errorf("saved board options = %v, want %q ticked", saved.Options, opts[0])
	}
}

// TestDefaultBoardScreenEscFinishes asserts esc and q report the screen done.
func TestDefaultBoardScreenEscFinishes(t *testing.T) {
	d := newDefaultBoardTestScreen(t)
	done, cmd := d.key("esc")
	if !done || cmd != nil {
		t.Errorf("esc must finish the screen with no command, got done=%v cmd=%v", done, cmd)
	}
}

// TestDefaultBoardScreenEditRefusesForBuiltin asserts 'e' on a built-in lane
// (which has no file on disk) refuses rather than trying to launch an editor
// on nothing.
func TestDefaultBoardScreenEditRefusesForBuiltin(t *testing.T) {
	d := newDefaultBoardTestScreen(t)
	d.idx, d.focus = 0, 0
	done, cmd := d.key("e")
	if done {
		t.Error("'e' must not finish the screen")
	}
	if cmd != nil {
		t.Error("'e' on a built-in must not launch an editor")
	}
	if !d.isErr {
		t.Error("'e' on a built-in must report a refusal")
	}
}

// TestHomeDKeyOpensDefaultBoardScreen asserts the 'd' keybind on the home
// screen opens the default board screen.
func TestHomeDKeyOpensDefaultBoardScreen(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))
	h := newHome(t, nil, 120, 30)

	h.Update(tea.KeyPressMsg{Code: 'd', Text: "d"})
	if h.board == nil {
		t.Fatal("'d' must open the default board screen")
	}
}
