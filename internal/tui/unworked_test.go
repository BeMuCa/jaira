package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The fixture's in-progress ticket has commits but none of the three outcome
// fields the lane declares it produces, which is exactly the state that was
// invisible: work claimed, lane not run.
func inProgressTicket(t *testing.T, m *Model) *ticket.Ticket {
	t.Helper()
	for _, tk := range m.tickets {
		if tk.Status == "in-progress" {
			return tk
		}
	}
	t.Fatal("fixture has no in-progress ticket")
	return nil
}

func TestUnworkedTicketIsMarkedOnTheBoard(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := inProgressTicket(t, m)

	out := stripANSI(m.render())
	if !strings.Contains(out, "unworked") {
		t.Errorf("a ticket owing its lane's whole output is not marked:\n%s", out)
	}

	// Produce what the lane declares, and the mark must go: this is the
	// difference the board exists to show.
	if _, err := m.store.Mutate(tk.ID, func(x *ticket.Ticket) error {
		for f, v := range map[string]string{
			ticket.FieldOutcomeWhat:     "streamed the writer",
			ticket.FieldOutcomeWhy:      "the whole file was buffered",
			ticket.FieldOutcomeResolves: "yes, the leak is gone",
		} {
			if err := x.Doc().SetScalar(f, v); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if out := stripANSI(m.render()); strings.Contains(out, "unworked") {
		t.Errorf("a ticket whose lane has produced its output is still marked unworked:\n%s", out)
	}
}

// Lane focus carries the same flag block as the board on purpose — a ticket must
// not look worked on one screen and unworked on the other.
func TestUnworkedTicketIsMarkedInLaneFocus(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := inProgressTicket(t, m)
	for i, c := range m.cols {
		if c.lane.ID == "in-progress" {
			m.laneIdx = i
		}
	}
	m.mode = modeLaneFocus

	if out := stripANSI(m.renderLaneFocus()); !strings.Contains(out, "unworked") {
		t.Errorf("lane focus does not mark an unworked ticket:\n%s", out)
	}

	if _, err := m.store.Mutate(tk.ID, func(x *ticket.Ticket) error {
		for f, v := range map[string]string{
			ticket.FieldOutcomeWhat:     "streamed the writer",
			ticket.FieldOutcomeWhy:      "the whole file was buffered",
			ticket.FieldOutcomeResolves: "yes, the leak is gone",
		} {
			if err := x.Doc().SetScalar(f, v); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if out := stripANSI(m.renderLaneFocus()); strings.Contains(out, "unworked") {
		t.Errorf("lane focus still marks a worked ticket:\n%s", out)
	}
}

// The trap this had to avoid: 'todo' and 'backlog' produce nothing, so every
// card in them would light up if the flag were derived from an empty contract
// rather than an unsatisfied one.
func TestLanesThatProduceNothingAreNotMarkedUnworked(t *testing.T) {
	m := newTestModel(t, 150, 32)
	for _, id := range []string{"todo", "backlog"} {
		for i, c := range m.cols {
			if c.lane.ID != id {
				continue
			}
			if len(c.tickets) == 0 {
				t.Fatalf("fixture has no ticket in %s", id)
			}
			m.laneIdx = i
			m.mode = modeLaneFocus
			if out := stripANSI(m.renderLaneFocus()); strings.Contains(out, "unworked") {
				t.Errorf("lane %q declares no output and its cards are marked unworked:\n%s", id, out)
			}
		}
	}
}
