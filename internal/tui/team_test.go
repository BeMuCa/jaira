package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// After a pull the board can hold teammates' claims. Whose ticket a card is
// must be scannable: someone else's carries an @-marked name, your own stays
// the quiet dimmed name it always was.
func TestForeignAssigneeIsMarkedOnTheCard(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.me = "berk"

	mine := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "mine", Assignee: "berk"}
	theirs := &ticket.Ticket{ID: "01KZZR4CBGDM5T35SZDR72PQYG", Title: "theirs", Assignee: "sam"}

	own := stripANSI(m.renderCard(mine, 30, false))
	if strings.Contains(own, "@") {
		t.Errorf("own ticket carries the foreign marker:\n%s", own)
	}
	foreign := stripANSI(m.renderCard(theirs, 30, false))
	if !strings.Contains(foreign, "@sam") {
		t.Errorf("someone else's ticket is not marked @sam:\n%s", foreign)
	}
}

// The TUI's move claims an unassigned ticket the same way the CLI's does, and
// the claim satisfies the promotion gate within the same move.
func TestApplyMoveClaimsTheUnassignedTicket(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	s := newTestStore(t)
	f := map[string]string{
		ticket.FieldID: ticket.NewID(time.Now()), ticket.FieldTitle: "unclaimed",
		ticket.FieldStatus: "backlog", ticket.FieldCreator: "berk",
		ticket.FieldGoal: "g", ticket.FieldContext: "c", ticket.FieldDoD: "d",
	}
	tk, err := s.Create(f, map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}, "")
	if err != nil {
		t.Fatal(err)
	}
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 120, 40
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.selectByID(tk.ID)
	for i, l := range m.lanes.Lanes {
		if l.ID == "todo" {
			m.moveTarget = i
		}
	}
	m.applyMove()

	got, err := s.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "todo" {
		t.Fatalf("status = %q, want todo (move refused?)", got.Status)
	}
	if got.Assignee != "berk" {
		t.Errorf("assignee = %q, want berk: the pull must claim", got.Assignee)
	}
}

// The remaining body sections render like fields — label left, content right —
// and a heading with nothing under it is skipped entirely.
func TestBodySectionsRenderAsColumns(t *testing.T) {
	m := newTestModel(t, 100, 40)
	m.detail = &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo",
		Body: "# t\n\n## Options\n\n- [x] brainstorm\n- [ ] planning\n\n## Progress\n\n- **2026-08-18 10:00 · berk** — tried X, it fails on Y\n\n## Notes\n\n",
	}

	out := stripANSI(m.renderDetail())
	if strings.Contains(out, "##") {
		t.Errorf("raw markdown headings leaked into the pane:\n%s", out)
	}
	if !strings.Contains(out, "options") || !strings.Contains(out, "[x] brainstorm") {
		t.Errorf("options section lost its label or items:\n%s", out)
	}
	if !strings.Contains(out, "progress") || !strings.Contains(out, "tried X, it fails on Y") {
		t.Errorf("progress section lost its label or entry:\n%s", out)
	}
	if strings.Contains(out, "notes") {
		t.Errorf("an empty section rendered its heading anyway:\n%s", out)
	}
}

// The board's n key captures the same way the CLI does: unassigned, creator
// recorded — the claim happens at the pull, not at the capture.
func TestTUICreateLeavesTheTicketUnassigned(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	m := newTestModel(t, 120, 40)
	m.createTicket("captured from the board")

	ts, err := m.store.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range ts {
		if tk.Title == "captured from the board" {
			if tk.Assignee != "" {
				t.Errorf("assignee = %q, want nobody", tk.Assignee)
			}
			if tk.Creator != "berk" {
				t.Errorf("creator = %q, want berk", tk.Creator)
			}
			return
		}
	}
	t.Fatal("the created ticket is not in the store")
}
