package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// parkIn moves the fixture's checkpoint ticket into a lane and selects it.
func parkIn(t *testing.T, m *Model, laneID string) *ticket.Ticket {
	t.Helper()
	tk := m.cols[laneIndex(t, m, "human")].tickets[0]
	if _, err := m.store.Mutate(tk.ID, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldStatus, laneID)
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.laneIdx, m.cardIdx = laneIndex(t, m, laneID), 0
	return m.cols[m.laneIdx].tickets[0]
}

// The sign-off marker belongs to the lane a person has to get the ticket out of,
// which is the lane that declares requires-human-exit — not the lane whose id
// happens to be "review". The model's review lane is the agent's, and marking it
// as waiting on the reader is what sent finished work to the wrong column.
func TestSignOffMarkerFollowsTheHumanExitLane(t *testing.T) {
	for _, c := range []struct {
		lane   string
		marked bool
	}{
		{"signoff", true},
		{"review", false},
	} {
		t.Run(c.lane, func(t *testing.T) {
			m := newTestModel(t, 120, 34)
			tk := parkIn(t, m, c.lane)

			card := stripANSI(m.renderLaneCard(tk, 70, true))
			board := stripANSI(m.renderCard(tk, 30, true))
			m.openDetail()
			detail := stripANSI(m.detailBody(m.detail, 90))

			for name, out := range map[string]string{"board card": board, "lane card": card} {
				if got := strings.Contains(out, "sign off"); got != c.marked {
					t.Errorf("%s: sign-off flag = %v, want %v:\n%s", name, got, c.marked, out)
				}
			}
			if got := strings.Contains(detail, "waiting on your sign-off"); got != c.marked {
				t.Errorf("detail: sign-off row = %v, want %v:\n%s", got, c.marked, detail)
			}
		})
	}
}

// The question marker stays on the lane that requires a question, untouched.
func TestAsksMarkerStaysOnTheQuestionLane(t *testing.T) {
	m := newTestModel(t, 120, 34)
	tk := parkIn(t, m, "human")

	if out := stripANSI(m.renderCard(tk, 30, true)); !strings.Contains(out, "asks") {
		t.Errorf("the HITL card lost its asks flag:\n%s", out)
	}
	m.openDetail()
	if out := stripANSI(m.detailBody(m.detail, 90)); !strings.Contains(out, "waiting on your answer") {
		t.Errorf("the HITL detail row lost its question marker:\n%s", out)
	}
}
