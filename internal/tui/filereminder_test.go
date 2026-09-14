package tui

import (
	"strconv"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The board mentions a pile of finished tickets and does nothing about it. The
// threshold is the whole behaviour: below it the board says nothing, above it
// one line, and never a file moved.
func TestTheBoardMentionsAPileAndOnlyMentionsIt(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	terminal := lanes.Terminal()
	if terminal == nil {
		t.Fatal("the built-in lanes have no terminal lane")
	}

	m := &Model{lanes: lanes}
	fill := func(n int) {
		m.tickets = nil
		for i := 0; i < n; i++ {
			m.tickets = append(m.tickets, &ticket.Ticket{ID: "01T", Status: terminal.ID})
		}
	}

	fill(fileReminder - 1)
	if got := m.readyToFile(); got != 0 {
		t.Errorf("a normal day tripped the reminder: %d", got)
	}
	fill(fileReminder)
	if got := m.readyToFile(); got != fileReminder {
		t.Errorf("at the threshold the board said %d", got)
	}
	fill(fileReminder + 7)
	if got := m.readyToFile(); got != fileReminder+7 {
		t.Errorf("above the threshold the board said %d", got)
	}

	// A ticket that is only on a ref is not this board's to file.
	m.tickets = nil
	for i := 0; i < fileReminder+3; i++ {
		m.tickets = append(m.tickets, &ticket.Ticket{ID: "01R", Status: terminal.ID, ReadOnly: true})
	}
	if got := m.readyToFile(); got != 0 {
		t.Errorf("ref-only tickets counted towards filing: %d", got)
	}
}

// The reminder has to reach the screen, and it has to name the cut. A counter
// nobody renders is a counter nobody reads, and a line that says only how many
// leaves the reader with no way to act on it — so this test looks at the
// rendered status bar rather than at readyToFile.
func TestTheRenderedHintNamesTheFilingCommand(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	terminal := lanes.Terminal()
	if terminal == nil {
		t.Fatal("the built-in lanes have no terminal lane")
	}

	m := &Model{lanes: lanes, width: 200}
	fill := func(n int) {
		m.tickets = nil
		for i := 0; i < n; i++ {
			m.tickets = append(m.tickets, &ticket.Ticket{ID: "01T", Status: terminal.ID})
		}
	}

	fill(fileReminder - 1)
	if bar := m.statusBar(); strings.Contains(bar, "to file") || strings.Contains(bar, fileCommand) {
		t.Errorf("below the threshold the board mentioned filing: %q", bar)
	}

	fill(fileReminder)
	bar := m.statusBar()
	if !strings.Contains(bar, "to file") {
		t.Errorf("at the threshold the board said nothing about filing: %q", bar)
	}
	if !strings.Contains(bar, strconv.Itoa(fileReminder)) {
		t.Errorf("the line does not say how many are waiting: %q", bar)
	}
	if !strings.Contains(bar, fileCommand) {
		t.Errorf("the line does not name %q, so a reader cannot act on it: %q", fileCommand, bar)
	}
}
