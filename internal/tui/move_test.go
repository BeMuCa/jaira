package tui

import (
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// applyMove threads the acting identity into the gate, the same as the CLI's
// move command does. Without it, the ownership check would see an empty actor
// and refuse every board move of an assigned ticket, including the owner's
// own — which would make dragging a ticket across the board in the TUI broken
// for anything but unassigned scrap.
func TestApplyMoveAllowsOwnerToMoveOwnTicket(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	m := newTestModel(t, 150, 32)

	now := time.Now()
	tk, err := m.store.Create(map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     "Owned by berk",
		ticket.FieldStatus:    "todo",
		ticket.FieldReady:     "true",
		ticket.FieldAssignee:  "berk",
		ticket.FieldGoal:      "prove the owner can move their own ticket",
		ticket.FieldDoD:       "the move succeeds",
		ticket.FieldContext:   "ownership gate regression test",
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}, map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.selectByID(tk.ID)
	if got := m.selected(); got == nil || got.ID != tk.ID {
		t.Fatalf("did not select the fixture ticket, got %v", got)
	}
	m.openMove()
	for i, l := range m.lanes.Lanes {
		if l.ID == "in-progress" {
			m.moveTarget = i
		}
	}

	m.applyMove()

	if m.isErr {
		t.Fatalf("owner's own move was refused: %s", m.message)
	}
	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "in-progress" {
		t.Errorf("ticket did not move: status = %q", got.Status)
	}
}

// The same move, attempted by someone other than the assignee, is refused —
// proving the gate is actually wired into the board's drag-move, not only
// into the CLI.
func TestApplyMoveRefusesNonOwner(t *testing.T) {
	t.Setenv("JAIRA_USER", "someone-else")
	m := newTestModel(t, 150, 32)

	now := time.Now()
	tk, err := m.store.Create(map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     "Owned by berk",
		ticket.FieldStatus:    "todo",
		ticket.FieldReady:     "true",
		ticket.FieldAssignee:  "berk",
		ticket.FieldGoal:      "prove a non-owner is refused",
		ticket.FieldDoD:       "the move is refused",
		ticket.FieldContext:   "ownership gate regression test",
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}, map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.selectByID(tk.ID)
	m.openMove()
	for i, l := range m.lanes.Lanes {
		if l.ID == "in-progress" {
			m.moveTarget = i
		}
	}

	m.applyMove()

	if !m.isErr {
		t.Fatal("expected the move to be refused for a non-owner")
	}
	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "todo" {
		t.Errorf("ticket moved despite refusal: status = %q", got.Status)
	}
}
