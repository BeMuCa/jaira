package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// refuseAMove sets up a move the gate will refuse and runs it, leaving the
// refusal on screen. It returns the ticket and the lane it tried to reach.
func refuseAMove(t *testing.T, m *Model, laneID, title, to string) (*ticket.Ticket, string) {
	t.Helper()
	tk := focusTicket(t, m, laneID, title)
	m.openDetail()
	press(m, "m")
	m.moveTarget = laneOrder(t, m, to)
	press(m, "enter")
	if m.mode != modeMessage || !m.isErr {
		t.Fatalf("the move was not refused: mode = %v, isErr = %v", m.mode, m.isErr)
	}
	return tk, to
}

// statusOnDisk reads the ticket's lane from the store, not from the model, so a
// test cannot pass on a cached copy.
func statusOnDisk(t *testing.T, m *Model, id string) string {
	t.Helper()
	full, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	return full.Status
}

// The refusal has to say how to override it, and must not have moved anything.
func TestRefusedMoveOffersForce(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk, _ := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")

	if !strings.Contains(m.message, "press f to override") {
		t.Errorf("the refusal does not offer f:\n%s", m.message)
	}
	if strings.Contains(m.message, "use the CLI") {
		t.Errorf("the refusal still sends the user to the CLI:\n%s", m.message)
	}
	if got := statusOnDisk(t, m, tk.ID); got != "backlog" {
		t.Errorf("the refused move happened anyway: status = %q", got)
	}
}

// f asks once more. It does not move anything on its own.
func TestForceAsksBeforeItMoves(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk, _ := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")

	press(m, "f")

	if m.mode != modeMessage {
		t.Fatalf("f left the message: mode = %v", m.mode)
	}
	if !strings.Contains(m.message, "anyway") {
		t.Errorf("the confirmation does not ask again:\n%s", m.message)
	}
	if got := statusOnDisk(t, m, tk.ID); got != "backlog" {
		t.Errorf("f moved the ticket without confirmation: status = %q", got)
	}
}

// f then y is the whole override: the ticket moves, and it says what it broke.
func TestForceThenYesMovesAndReportsWhatItOverrode(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk, to := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")

	press(m, "f")
	press(m, "y")

	if got := statusOnDisk(t, m, tk.ID); got != to {
		t.Errorf("the forced move did not land: status = %q, want %q", got, to)
	}
	if m.mode != modeMessage {
		t.Fatalf("the override said nothing: mode = %v", m.mode)
	}
	if !strings.Contains(strings.ToLower(m.message), "overrode") {
		t.Errorf("the override is not reported:\n%s", m.message)
	}
	// Dismissing lands back on the page the move was started from.
	press(m, "esc")
	if m.mode != modeDetail {
		t.Errorf("landed on %v after the override, want the open ticket", m.mode)
	}
}

// The ready flag is recomputed on a forced move exactly as on a clean one, so a
// forced ticket does not carry a stale "ready" into its new lane.
func TestForcedMoveRecomputesReady(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk, _ := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")

	press(m, "f")
	press(m, "y")

	full, err := m.store.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if full.Ready {
		t.Error("a ticket forced past its own refusal is marked ready")
	}
}

func TestForceCancelledLeavesTheTicketAlone(t *testing.T) {
	for _, key := range []string{"n", "esc", "q"} {
		t.Run(key, func(t *testing.T) {
			m := newTestModel(t, 120, 34)
			tk, _ := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")
			press(m, "f")

			press(m, key)

			if got := statusOnDisk(t, m, tk.ID); got != "backlog" {
				t.Errorf("%q moved the ticket: status = %q", key, got)
			}
			if m.mode != modeDetail {
				t.Errorf("%q landed on %v, want the open ticket", key, m.mode)
			}
		})
	}
}

// Dismissing a refusal drops the offer with it, so a later f cannot fire a move
// the user has walked away from.
func TestDismissedRefusalDropsTheOffer(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk, _ := refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")

	press(m, "esc")
	if m.mode != modeDetail {
		t.Fatalf("dismissed to %v, want the open ticket", m.mode)
	}
	m.notify("something else entirely", false)
	press(m, "f")
	press(m, "y")

	if got := statusOnDisk(t, m, tk.ID); got != "backlog" {
		t.Errorf("a stale offer fired: status = %q", got)
	}
}

// Force is the CLI's --force, which overrides any refusal, not only ownership.
func TestForceOverridesADependencyRefusalToo(t *testing.T) {
	m := newTestModel(t, 120, 34)
	// "Refactor auth middleware" is blocked by work still in flight.
	tk, to := refuseAMove(t, m, "todo", "Refactor auth middleware", "in-progress")
	if !strings.Contains(strings.ToLower(m.message), "block") {
		t.Fatalf("expected a dependency refusal, got:\n%s", m.message)
	}

	press(m, "f")
	press(m, "y")

	if got := statusOnDisk(t, m, tk.ID); got != to {
		t.Errorf("a blocked ticket could not be forced: status = %q, want %q", got, to)
	}
}

func TestMessageFooterNamesWhatIsOnOffer(t *testing.T) {
	m := newTestModel(t, 120, 34)

	m.notify("plain note", false)
	if got := stripANSI(m.renderMessage()); !strings.Contains(got, "esc dismiss") || strings.Contains(got, "override") {
		t.Errorf("plain message footer:\n%s", got)
	}

	m = newTestModel(t, 120, 34)
	refuseAMove(t, m, "backlog", "Investigate flaky logout test", "todo")
	if got := stripANSI(m.renderMessage()); !strings.Contains(got, "f override") {
		t.Errorf("refusal footer does not offer f:\n%s", got)
	}

	press(m, "f")
	if got := stripANSI(m.renderMessage()); !strings.Contains(got, "y override") || !strings.Contains(got, "n cancel") {
		t.Errorf("confirmation footer:\n%s", got)
	}
}
