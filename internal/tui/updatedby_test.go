package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Several people and several agent sessions write the same store. updated-at
// said a ticket moved under you and never said who, which is the one thing you
// want to know before touching it.
func TestSomeoneElsesChangeIsMarkedOnTheCard(t *testing.T) {
	m := newTestModel(t, 150, 32)
	// A ticket in a lane the board actually draws: the window shows six of ten
	// lanes at this width, and 'done' is not one of them.
	tk := inProgressTicket(t, m)

	// Through the real write path: a teammate's session is a store whose Actor
	// is theirs, and the field is written by Mutate rather than by a caller.
	theirs := *m.store
	theirs.Actor = "someone-else"
	if _, err := theirs.Mutate(tk.ID, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldGoal, "they changed the goal")
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	// The glyph, not the whole name: a card is as wide as its column, and a
	// long name is truncated there like every other flag.
	if out := stripANSI(m.render()); !strings.Contains(out, "✎ someo") {
		t.Errorf("a ticket last written by somebody else is not marked:\n%s", out)
	}

	// Your own change clears it: the marker is about the other person, and one
	// that fired on everything would read as decoration.
	if _, err := m.store.Mutate(tk.ID, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldGoal, "and then I did")
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if out := stripANSI(m.render()); strings.Contains(out, "✎") {
		t.Errorf("your own change is marked as somebody else's:\n%s", out)
	}
}

// One person is not one string. A marker that called your own change somebody
// else's would be worse than no marker at all.
func TestAnAliasOfYoursIsNotSomeoneElse(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.myAliases = []string{"BeMuCa", "berk@example.test"}
	if !m.isMe("berk@example.test") {
		t.Error("an alias is not recognised as me")
	}
	if !m.isMe("bemuca") {
		t.Error("matching is not case-insensitive")
	}
	if m.isMe("teammate") {
		t.Error("a teammate is recognised as me")
	}
	if m.isMe("") {
		t.Error("an empty name is recognised as me")
	}
}

// A ticket written by a version that did not record it reads as unknown, not as
// somebody: an empty field must not put a marker on every old ticket.
func TestATicketWithNoUpdatedByIsNotMarked(t *testing.T) {
	m := newTestModel(t, 150, 32)
	for _, tk := range m.tickets {
		if tk.UpdatedBy != "" {
			t.Fatalf("fixture ticket %s already records an author", tk.ID)
		}
	}
	if out := stripANSI(m.render()); strings.Contains(out, "✎") {
		t.Errorf("tickets with no updated-by are marked:\n%s", out)
	}
}
