package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// ticketTitled is the board's ticket with that title, failing the test without one.
func ticketTitled(t *testing.T, m *Model, title string) *ticket.Ticket {
	t.Helper()
	for _, tk := range m.tickets {
		if tk.Title == title {
			return tk
		}
	}
	t.Fatalf("no ticket %q", title)
	return nil
}

// claim writes a claim onto a ticket, as 'jaira claim' does, and reloads.
func claim(t *testing.T, m *Model, title, by string, at time.Time) *ticket.Ticket {
	t.Helper()
	tk := ticketTitled(t, m, title)
	if _, err := m.store.Mutate(tk.ID, func(tk *ticket.Ticket) error {
		if err := tk.Doc().SetScalar(ticket.FieldClaimedBy, by); err != nil {
			return err
		}
		return tk.Doc().SetRaw(ticket.FieldClaimedAt, ticket.FormatTime(at))
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	return ticketTitled(t, m, title)
}

// moveTo puts a ticket in a lane directly, without the gates, and reloads.
func moveTo(t *testing.T, m *Model, title, status string) {
	t.Helper()
	tk := ticketTitled(t, m, title)
	if _, err := m.store.Mutate(tk.ID, func(tk *ticket.Ticket) error {
		return tk.Doc().SetScalar(ticket.FieldStatus, status)
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
}

func snakeBGIn(s string, k snakeKind) bool {
	return strings.Contains(s, "\x1b[48;"+snakeParams(k)+"m")
}

// A card nobody is working on has no snake; claim it and an agent's snake runs
// round it, in the agent's colour, with the card's text untouched.
func TestAClaimedCardCarriesTheAgentSnake(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := ticketTitled(t, m, "Fix session cookie dropped on 302")
	before := m.renderCardBlock(tk, 30, false, false)
	if m.snakeOf(tk) != snakeNone || snakeBGIn(before, snakeAgent) {
		t.Fatalf("an unclaimed in-progress card has a snake:\n%q", before)
	}

	tk = claim(t, m, "Fix session cookie dropped on 302", "sess-1", time.Now())
	if got := m.snakeOf(tk); got != snakeAgent {
		t.Fatalf("snakeOf = %v, want the agent's", got)
	}
	after := m.renderCardBlock(tk, 30, false, false)
	if !snakeBGIn(after, snakeAgent) {
		t.Errorf("the claimed card is not painted in the agent's colour:\n%q", after)
	}
	if stripANSI(after) != stripANSI(before) {
		t.Errorf("the snake changed the card's text:\nbefore %q\nafter  %q", stripANSI(before), stripANSI(after))
	}
}

// The colour says who the ticket waits on: a person in a lane that takes a
// question or a sign-off, whoever holds it; a reviewer in the fresh lane.
func TestTheSnakeColourFollowsTheLane(t *testing.T) {
	m := newTestModel(t, 150, 32)
	if got := m.snakeOf(ticketTitled(t, m, "Decide on cookie SameSite policy")); got != snakeHuman {
		t.Errorf("a HITL ticket: snakeOf = %v, want a person's", got)
	}
	moveTo(t, m, "Decide on cookie SameSite policy", "signoff")
	if got := m.snakeOf(ticketTitled(t, m, "Decide on cookie SameSite policy")); got != snakeHuman {
		t.Errorf("a sign-off ticket: snakeOf = %v, want a person's", got)
	}
	moveTo(t, m, "Fix session cookie dropped on 302", "review")
	tk := claim(t, m, "Fix session cookie dropped on 302", "sess-1", time.Now())
	if got := m.snakeOf(tk); got != snakeReview {
		t.Errorf("a claimed ticket in review: snakeOf = %v, want a reviewer's", got)
	}
	moveTo(t, m, "Add session test harness", "done")
	tk = claim(t, m, "Add session test harness", "sess-1", time.Now())
	if got := m.snakeOf(tk); got != snakeNone {
		t.Errorf("a claimed ticket in done: snakeOf = %v, want none", got)
	}
}

// Past its lease a claim keeps its snake only while the agent's process is
// known to be alive on this machine.
func TestAnExpiredClaimKeepsTheSnakeOnlyWhileItsAgentLives(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := claim(t, m, "Fix session cookie dropped on 302", "sess-1", time.Now().Add(-2*ticket.ClaimTTL))
	if got := m.snakeOf(tk); got != snakeNone {
		t.Fatalf("an expired claim with nothing alive: snakeOf = %v, want none", got)
	}
	m.liveClaims = map[string]bool{tk.ID: true}
	if got := m.snakeOf(tk); got != snakeAgent {
		t.Errorf("an expired claim whose agent lives: snakeOf = %v, want the agent's", got)
	}
}

// The snake moves with the clock while the animation runs, stands still once
// c stops it, and is drawn either way.
func TestCStopsTheSnakeWithoutHidingIt(t *testing.T) {
	m := newTestModel(t, 150, 32)
	id := ticketTitled(t, m, "Decide on cookie SameSite policy").ID
	if !m.anim {
		t.Fatal("the animation is off on a new board")
	}
	first := m.snakeCells(id, 20, 3)
	m.animPhase = 5
	moved := m.snakeCells(id, 20, 3)
	if equalCells(first, moved) {
		t.Error("a later phase drew the snake in the same cells")
	}

	press(m, "c")
	if m.anim {
		t.Fatal("c did not stop the animation")
	}
	if m.animating() {
		t.Error("the frame loop would keep running with the animation off")
	}
	still := m.snakeCells(id, 20, 3)
	if len(still) == 0 {
		t.Error("a stopped snake is not drawn")
	}
	m.animPhase = 9
	if !equalCells(still, m.snakeCells(id, 20, 3)) {
		t.Error("a stopped snake still moves")
	}
	if !strings.Contains(stripANSI(m.render()), "c animate") {
		t.Error("the key line does not say c starts it again")
	}
	press(m, "c")
	if !m.anim || !m.animating() {
		t.Error("a second c did not start it again")
	}
}

// The frame loop runs only while there is something to animate, and ends by
// itself once there is not.
func TestTheFrameLoopEndsWhenNothingIsWorked(t *testing.T) {
	m := newTestModel(t, 150, 32)
	if m.animKick() == nil {
		t.Fatal("a board with a ticket waiting on a person started no frame")
	}
	if m.animKick() != nil {
		t.Error("a second loop was started beside the first")
	}
	moveTo(t, m, "Decide on cookie SameSite policy", "backlog")
	if m.frame() != nil {
		t.Error("the loop asked for another frame with nothing left to animate")
	}
	if m.animPending {
		t.Error("a finished loop still counts as pending, so none could start again")
	}
}

// The snake goes once round the edge: the top row, the right column, the
// bottom row and the left column, never the middle.
func TestTheSnakeStaysOnTheEdge(t *testing.T) {
	m := newTestModel(t, 150, 32)
	for phase := 0; phase < 60; phase++ {
		m.animPhase = phase
		for r, cols := range m.snakeCells("01ABC", 20, 3) {
			for c := range cols {
				if r == 1 && c != 0 && c != 19 {
					t.Fatalf("phase %d: the snake is inside the card at row %d column %d", phase, r, c)
				}
			}
		}
	}
}

func equalCells(a, b map[int]map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for r, cols := range a {
		if len(cols) != len(b[r]) {
			return false
		}
		for c := range cols {
			if !b[r][c] {
				return false
			}
		}
	}
	return true
}
