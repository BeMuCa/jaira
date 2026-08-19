package tui

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// focusBacklog points the model at the backlog lane with the cursor on the
// named ticket, and returns that ticket.
func focusBacklog(t *testing.T, m *Model, title string) *ticket.Ticket {
	t.Helper()
	m.laneIdx = laneIndex(t, m, "backlog")
	for ci, tk := range m.cols[m.laneIdx].tickets {
		if tk.Title == title {
			m.cardIdx = ci
			return tk
		}
	}
	t.Fatalf("no backlog ticket titled %q in fixture", title)
	return nil
}

// backlogTicket finds a backlog ticket by title without moving the cursor.
func backlogTicket(t *testing.T, m *Model, title string) *ticket.Ticket {
	t.Helper()
	for _, tk := range m.cols[laneIndex(t, m, "backlog")].tickets {
		if tk.Title == title {
			return tk
		}
	}
	t.Fatalf("no backlog ticket titled %q in fixture", title)
	return nil
}

// touchOnDisk edits a ticket behind the TUI's back, the way the CLI or another
// session does, then lets the TUI notice. Passing the lane it is already in
// only bumps updated_at, which is enough to reorder its lane.
func touchOnDisk(t *testing.T, m *Model, id, status string) {
	t.Helper()
	if _, err := m.store.Mutate(id, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldStatus, status)
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
}

// The lane filling the screen must not be swapped out from under the user
// because the ticket they had selected was moved elsewhere.
func TestLaneFocusHoldsItsLaneWhenTheTicketLeaves(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modeLaneFocus
	moved := focusBacklog(t, m, "Investigate flaky logout test")

	touchOnDisk(t, m, moved.ID, "done")

	if want := laneIndex(t, m, "backlog"); m.laneIdx != want {
		t.Errorf("lane focus jumped: laneIdx = %d, want backlog at %d", m.laneIdx, want)
	}
	sel := m.selected()
	if sel == nil {
		t.Fatal("nothing selected after the reload")
	}
	if sel.ID == moved.ID {
		t.Error("cursor followed the ticket out of the lane")
	}
	if sel.Title != "Rate limit the login endpoint" {
		t.Errorf("cursor landed on %q, want the remaining backlog ticket", sel.Title)
	}
}

// The compact view's focused step is m.laneIdx too, so it holds the same way.
func TestCompactViewHoldsItsStepWhenTheTicketLeaves(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modePipeline
	moved := focusBacklog(t, m, "Investigate flaky logout test")

	touchOnDisk(t, m, moved.ID, "done")

	if want := laneIndex(t, m, "backlog"); m.laneIdx != want {
		t.Errorf("compact view jumped: laneIdx = %d, want backlog at %d", m.laneIdx, want)
	}
}

// Leaving an open ticket returns to the view it was opened from, so an open
// ticket must not defeat the hold either.
func TestOpenTicketFromLaneFocusStillHoldsTheLane(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modeLaneFocus
	moved := focusBacklog(t, m, "Investigate flaky logout test")
	m.openDetail()
	if m.detailFrom != modeLaneFocus {
		t.Fatalf("detailFrom = %v, want modeLaneFocus", m.detailFrom)
	}

	touchOnDisk(t, m, moved.ID, "done")

	if want := laneIndex(t, m, "backlog"); m.laneIdx != want {
		t.Errorf("open ticket jumped its lane: laneIdx = %d, want backlog at %d", m.laneIdx, want)
	}
}

// On the multi-column board the new lane is already on screen, so the cursor
// tracking your ticket is the point — that must not regress.
func TestBoardStillFollowsTheTicket(t *testing.T) {
	m := newTestModel(t, 170, 34)
	m.mode = modeBoard
	moved := focusBacklog(t, m, "Investigate flaky logout test")

	touchOnDisk(t, m, moved.ID, "done")

	if want := laneIndex(t, m, "done"); m.laneIdx != want {
		t.Errorf("board stopped following: laneIdx = %d, want done at %d", m.laneIdx, want)
	}
	if sel := m.selected(); sel == nil || sel.ID != moved.ID {
		t.Errorf("board lost the ticket: selected %v, want %s", idOf(sel), moved.ID)
	}
}

// Holding the lane must still follow the ticket inside it: cards are sorted by
// UpdatedAt, so any unrelated edit reorders them.
func TestHeldLaneStillFollowsTheTicketWithinIt(t *testing.T) {
	m := newTestModel(t, 120, 34)
	m.mode = modeLaneFocus
	kept := focusBacklog(t, m, "Rate limit the login endpoint")
	other := backlogTicket(t, m, "Investigate flaky logout test")

	touchOnDisk(t, m, other.ID, "backlog")

	if sel := m.selected(); sel == nil || sel.ID != kept.ID {
		t.Errorf("cursor drifted off its ticket: selected %v, want %s", idOf(sel), kept.ID)
	}
}
