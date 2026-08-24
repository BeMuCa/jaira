package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// typeInto feeds a string to the model one key at a time, the way a person does.
func typeInto(m *Model, s string) {
	for _, r := range s {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

func openFirstTicket(t *testing.T, m *Model) *ticket.Ticket {
	t.Helper()
	m.laneIdx, m.cardIdx = 0, 0
	m.openDetail()
	if m.detail == nil {
		t.Fatal("no ticket opened")
	}
	return m.detail
}

// The second step is the handle typed back, so a stray return key cannot delete
// anything. Everything else the board does is reversible; this is not.
func TestDeleteNeedsTheHandleTypedBack(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := openFirstTicket(t, m)
	before := len(m.tickets)

	m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})
	if m.mode != modeDelete {
		t.Fatalf("X on an open ticket opened mode %v, want modeDelete", m.mode)
	}
	// The ticket stays on screen while its handle is typed.
	if out := stripANSI(m.render()); !strings.Contains(out, ticket.Handle(tk.ID)) {
		t.Errorf("the ticket being deleted is not shown:\n%s", out)
	}

	typeInto(m, "yes")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if len(m.tickets) != before {
		t.Fatal("a wrong answer deleted a ticket")
	}
	if _, err := m.store.Load(tk.ID); err != nil {
		t.Errorf("the ticket is gone after a refused confirmation: %v", err)
	}
}

func TestDeleteRemovesTheTicketWhenTheHandleMatches(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := openFirstTicket(t, m)
	before := len(m.tickets)

	m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})
	typeInto(m, ticket.Handle(tk.ID))
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if _, err := m.store.Load(tk.ID); err == nil {
		t.Error("the ticket file survived a confirmed delete")
	}
	if len(m.tickets) != before-1 {
		t.Errorf("the board still holds %d tickets, want %d", len(m.tickets), before-1)
	}
	if m.detail != nil {
		t.Error("the deleted ticket is still open")
	}
}

// esc leaves everything alone, and x must still mean archive.
func TestDeleteCanBeCancelledAndXStillArchives(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := openFirstTicket(t, m)

	m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})
	typeInto(m, ticket.Handle(tk.ID))
	m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.mode != modeDetail {
		t.Errorf("esc left mode %v, want the ticket still open", m.mode)
	}
	if _, err := m.store.Load(tk.ID); err != nil {
		t.Errorf("esc deleted the ticket: %v", err)
	}

	// x on the board archives, which is a move, not a removal.
	m.mode = modeBoard
	m.detail = nil
	m.laneIdx, m.cardIdx = 0, 0
	sel := m.selected()
	m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if _, err := m.store.Load(sel.ID); err == nil {
		t.Error("x did not take the ticket off the board")
	}
	names, err := m.store.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range names {
		if o.ID == sel.ID {
			t.Error("x left the ticket on the board")
		}
	}
}
