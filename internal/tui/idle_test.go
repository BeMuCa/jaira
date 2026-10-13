package tui

import (
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Not parallel: ticket.Counts is process-wide.

func titleOf(m *Model, id string) string {
	for _, t := range m.tickets {
		if t.ID == id {
			return t.Title
		}
	}
	return ""
}

// A tick on a board where nothing changed decodes no ticket and parses no
// frontmatter; a change made behind the board's back shows on the next tick,
// and costs exactly one decode.
func TestIdleTickDecodesNothingAndChangeAppearsNextTick(t *testing.T) {
	m := newTestModel(t, 150, 32)
	defer m.Close()

	dec0, par0 := ticket.Counts()
	for i := 0; i < 3; i++ {
		m.Update(tickMsg(time.Now()))
	}
	dec1, par1 := ticket.Counts()
	if dec1 != dec0 || par1 != par0 {
		t.Fatalf("idle ticks did work: %d decodes, %d header parses", dec1-dec0, par1-par0)
	}

	id := m.tickets[0].ID
	if _, err := m.store.Mutate(id, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldTitle, "changed behind the board")
	}); err != nil {
		t.Fatal(err)
	}
	dec2, _ := ticket.Counts()
	m.Update(tickMsg(time.Now()))
	dec3, _ := ticket.Counts()
	if got := titleOf(m, id); got != "changed behind the board" {
		t.Fatalf("the change did not show on the next tick: title %q", got)
	}
	if dec3-dec2 != 1 {
		t.Fatalf("one changed ticket should cost one decode, got %d", dec3-dec2)
	}
}
