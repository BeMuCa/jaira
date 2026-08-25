package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

func boardOnDisk(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(root, "no-lanes"))
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(map[string]string{
		ticket.FieldID: ticket.NewID(time.Now()), ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, ""); err != nil {
		t.Fatal(err)
	}
	return root
}

// Removing a board is irreversible, so the cursor starts on No and a stray
// return key must cancel rather than delete.
func TestDropBoardStartsOnNoAndEnterCancels(t *testing.T) {
	root := boardOnDisk(t)
	d := newDropBoard(root, "demo", false)

	if d.yes {
		t.Error("the cursor starts on Yes")
	}
	if d.idx != len(d.parts) {
		t.Errorf("idx = %d, want the answer row %d", d.idx, len(d.parts))
	}
	out := stripANSI(d.render(90, 30))
	if !strings.Contains(out, "No") || !strings.Contains(out, "Yes") {
		t.Errorf("no yes/no choice on screen:\n%s", out)
	}

	done, removed := d.key("enter")
	if !done || removed {
		t.Errorf("enter on No: done=%v removed=%v, want done and nothing removed", done, removed)
	}
	if _, err := os.Stat(filepath.Join(root, ticket.DirName)); err != nil {
		t.Errorf("the board was removed by a plain enter: %v", err)
	}
}

// Everything is ticked by default and the whole .jaira goes, which is what
// takes the board off the home screen.
func TestDropBoardRemovesEverythingByDefault(t *testing.T) {
	root := boardOnDisk(t)
	d := newDropBoard(root, "demo", false)
	for _, p := range d.parts {
		if !p.on {
			t.Errorf("%q is not ticked by default", p.label)
		}
	}

	d.key("l") // move to Yes
	done, removed := d.key("enter")
	if !done || !removed {
		t.Fatalf("done=%v removed=%v, want both", done, removed)
	}
	if _, err := os.Stat(filepath.Join(root, ticket.DirName)); !os.IsNotExist(err) {
		t.Errorf(".jaira survived: %v", err)
	}
}

// Unticking leaves that part alone — the point of the checkboxes.
func TestDropBoardKeepsWhatIsUnticked(t *testing.T) {
	root := boardOnDisk(t)
	d := newDropBoard(root, "demo", false)

	// Untick the tickets and the board directory itself.
	for i, p := range d.parts {
		if p.label == "tickets" || strings.HasPrefix(p.label, "the rest of .jaira") {
			d.idx = i
			d.key(" ")
		}
	}
	d.idx = len(d.parts)
	d.key("l")
	if done, removed := d.key("enter"); !done || !removed {
		t.Fatalf("done=%v removed=%v", done, removed)
	}

	if _, err := os.Stat(filepath.Join(root, ticket.DirName)); err != nil {
		t.Errorf(".jaira was removed although it was unticked: %v", err)
	}
	files, err := filepath.Glob(filepath.Join(root, ticket.DirName, ticket.TicketsSubdir, "*.md"))
	if err != nil || len(files) != 1 {
		t.Errorf("tickets = %v (%v), want the one kept", files, err)
	}
}

// The board you are looking at cannot be pulled out from under you.
func TestDropBoardRefusesTheOpenBoard(t *testing.T) {
	root := boardOnDisk(t)
	d := newDropBoard(root, "demo", true)
	d.key("l")
	if d.yes {
		t.Error("Yes is reachable for the board currently open")
	}
	out := stripANSI(d.render(90, 30))
	if !strings.Contains(out, "Switch to another one first") {
		t.Errorf("the screen does not say why:\n%s", out)
	}
}

// The catalogue is a different thing in a different place, and the screen says
// so — otherwise "delete this board's lanes" reads as "delete my lanes".
func TestDropBoardSaysTheCatalogueIsSafe(t *testing.T) {
	out := stripANSI(newDropBoard(boardOnDisk(t), "demo", false).render(90, 30))
	if !strings.Contains(out, "~/.jaira/lanes are never touched") {
		t.Errorf("the screen does not say the catalogue is safe:\n%s", out)
	}
}

// The launcher is the list you are looking at when you notice a board should
// not be on it, and it is a different screen from the board's own switcher —
// which is how x came to exist on one and not the other.
func TestLauncherRemovesABoard(t *testing.T) {
	root := boardOnDisk(t)
	// The launcher shows the catalogue, since it spans every board rather than
	// one project. NewHome refuses to build without it, so a test that skips it
	// tests something production cannot reach.
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	h := &Home{width: 100, height: 30, lanes: lanes, entries: []HomeEntry{{Root: root, Name: "demo"}}}

	if out := stripANSI(h.render()); !strings.Contains(out, "x remove a board") {
		t.Errorf("the launcher does not offer it:\n%s", out)
	}

	h.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if h.drop == nil {
		t.Fatal("x did not open the removal screen")
	}
	if out := stripANSI(h.render()); !strings.Contains(out, "Remove board") {
		t.Errorf("the removal screen is not shown:\n%s", out)
	}

	// Still starts on No.
	h.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if _, err := os.Stat(filepath.Join(root, ticket.DirName)); err != nil {
		t.Fatalf("enter on No removed the board: %v", err)
	}

	h.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	h.Update(tea.KeyPressMsg{Code: 'l', Text: "l"})
	h.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if _, err := os.Stat(filepath.Join(root, ticket.DirName)); !os.IsNotExist(err) {
		t.Errorf("the board survived: %v", err)
	}
	if h.drop != nil {
		t.Error("the removal screen is still up")
	}
}
