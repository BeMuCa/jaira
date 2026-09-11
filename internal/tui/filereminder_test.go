package tui

import (
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
