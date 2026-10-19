package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

func TestRejectOnTheBoardWritesTheReason(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.laneIdx, m.cardIdx = 0, 0
	sel := m.selected()

	m.Update(tea.KeyPressMsg{Code: 'R', Text: "R"})
	if m.mode != modeReject {
		t.Fatalf("R on the board opened mode %v, want modeReject", m.mode)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "reject "+ticket.Handle(sel.ID)) {
		t.Errorf("the prompt does not name the ticket:\n%s", out)
	}
	// An empty reason is not a rejection: the prompt stays open.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.mode != modeReject {
		t.Fatalf("enter with no reason left mode %v", m.mode)
	}
	typeInto(m, "nicht mehr nötig")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if _, err := m.store.Load(sel.ID); err == nil {
		t.Fatal("the rejected ticket is still on the board")
	}
	raw, err := os.ReadFile(filepath.Join(m.store.ArchiveDir(), filepath.Base(sel.Path)))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"rejected-kind: wontfix", "rejected-why: nicht mehr nötig", "rejected-at: "} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("archived file lacks %q:\n%s", want, raw)
		}
	}
}

// esc leaves the ticket where it was, open if it was open.
func TestRejectCanBeCancelledFromAnOpenTicket(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := openFirstTicket(t, m)

	m.Update(tea.KeyPressMsg{Code: 'R', Text: "R"})
	typeInto(m, "egal")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.mode != modeDetail {
		t.Errorf("esc left mode %v, want the ticket still open", m.mode)
	}
	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatalf("esc took the ticket off the board: %v", err)
	}
	if got.Rejected.Kind != "" {
		t.Errorf("esc left a rejection behind: %+v", got.Rejected)
	}
}
