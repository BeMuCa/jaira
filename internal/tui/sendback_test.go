package tui

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// signoffModel is a board with one ticket in signoff, open on its screen.
func signoffModel(t *testing.T) (*Model, *ticket.Ticket) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))
	t.Setenv("JAIRA_USER", "berk")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	stamp := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	tk, err := s.Create(map[string]string{
		ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: "half done",
		ticket.FieldStatus: "signoff", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
		ticket.FieldUpdatedAt: ticket.FormatTime(stamp), ticket.FieldAssignee: "berk",
		ticket.FieldGoal: "g", ticket.FieldContext: "c", ticket.FieldDoD: "it works",
		"outcome-what": "w", "outcome-why": "y", "outcome-resolves": "r",
	}, nil, "# x\n\n## Definition of Done\n\n- [ ] it works\n")
	if err != nil {
		t.Fatal(err)
	}
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	for li, c := range m.cols {
		for ci, x := range c.tickets {
			if x.ID == tk.ID {
				m.laneIdx, m.cardIdx = li, ci
			}
		}
	}
	m.openDetail()
	if m.detail == nil || m.detail.ID != tk.ID {
		t.Fatal("the signoff ticket did not open")
	}
	return m, tk
}

func pickLane(t *testing.T, m *Model, id string) {
	t.Helper()
	m.Update(tea.KeyPressMsg{Code: 'm', Text: "m"})
	if m.mode != modeMove {
		t.Fatalf("m opened mode %v, want modeMove", m.mode)
	}
	m.moveTarget = m.lanes.Index(id)
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
}

// A person sends work back out of signoff with m, says what is still open,
// and the ticket lands in the earlier lane with that reason as a note —
// no --force, no follow-up ticket.
func TestSendBackFromSignoffWritesTheReason(t *testing.T) {
	m, tk := signoffModel(t)
	if out := stripANSI(m.render()); !strings.Contains(out, "m send back") {
		t.Errorf("the signoff footer does not name the send-back:\n%s", out)
	}
	pickLane(t, m, "in-progress")
	if m.mode != modeSendBack {
		t.Fatalf("a send-back out of signoff opened mode %v, want the reason prompt", m.mode)
	}
	// No reason, no move: the prompt stays open.
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.mode != modeSendBack {
		t.Fatalf("enter with no reason left mode %v", m.mode)
	}
	typeInto(m, "DoD 2 fehlt noch")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "in-progress" {
		t.Fatalf("the ticket sits in %q, want in-progress (message: %q)", got.Status, m.message)
	}
	if !strings.Contains(got.Body, "sent back to in-progress from signoff: DoD 2 fehlt noch") {
		t.Errorf("the send-back left no note with its reason:\n%s", got.Body)
	}
	if got.BackBy != ticket.BackByPerson {
		t.Errorf("back-by = %q, want a person's", got.BackBy)
	}
}

// esc on the reason prompt moves nothing.
func TestSendBackCanBeCancelled(t *testing.T) {
	m, tk := signoffModel(t)
	pickLane(t, m, "in-progress")
	typeInto(m, "egal")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.mode != modeDetail {
		t.Errorf("esc left mode %v, want the ticket still open", m.mode)
	}
	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "signoff" {
		t.Errorf("esc moved the ticket to %q", got.Status)
	}
}

// Forward stays the accept key's alone: m into done is refused, and asks
// for no reason, since it is no send-back.
func TestMoveDialogCannotAcceptOutOfSignoff(t *testing.T) {
	m, tk := signoffModel(t)
	pickLane(t, m, "done")
	if m.mode == modeSendBack {
		t.Fatal("a move forward asked for a send-back reason")
	}
	got, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "signoff" {
		t.Errorf("the move dialog accepted the ticket into %q", got.Status)
	}
	if !strings.Contains(m.message, "human checkpoint") {
		t.Errorf("the refusal does not say why: %q", m.message)
	}
}
