package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/ticket"
)

// renderSignOff is the screen a person sees when a ticket is waiting on them.
//
// It answers four questions in order — what was wrong, what the agent did, why,
// and whether it actually holds — because that is the judgement being asked for.
// The implementer's account and the reviewer's verdict are shown side by side
// rather than merged: they are different claims, and when they disagree that
// disagreement is the most useful thing on the screen.
func (m *Model) renderSignOff() string {
	t := m.detail
	if t == nil {
		return m.renderBoard()
	}
	w := max(20, min(m.width, 78))
	var b strings.Builder

	fmt.Fprintf(&b, "%s  %s\n", styHandle.Render(ticket.Handle(t.ID)),
		styLaneTitle.Render(truncate(t.Title, max(1, w-10))))
	b.WriteString(styReview.Render(truncate("◆ waiting on your sign-off", w)) + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	section := func(label, body string) {
		if strings.TrimSpace(body) == "" {
			return
		}
		b.WriteString("\n" + styLaneTitle.Render(truncate(label, w)) + "\n")
		b.WriteString(wrap(body, max(10, w-2), 2) + "\n")
	}
	section("What was wrong", firstNonEmpty(t.Goal, t.Context))
	section("What was done", t.Outcome.What)
	section("Why", t.Outcome.Why)
	section("Does it solve it — the implementer's account", t.Outcome.Resolves)
	section("The reviewer's verdict", t.ReviewVerdict)

	if done, total := checklistProgress(t.DoDItems); total > 0 {
		b.WriteString("\n" + styLaneTitle.Render("Definition of Done") +
			" " + styLaneCount.Render(fmt.Sprintf("%d/%d", done, total)) + "\n")
		for _, it := range t.DoDItems {
			sty := styMeta
			if it.Checked() {
				sty = styOK
			}
			fmt.Fprintf(&b, "  %s %s\n", sty.Render("["+it.State.Marker()+"]"),
				truncate(it.Text, max(1, w-6)))
		}
	}

	b.WriteString("\n" + styMeta.Render(truncate(
		"a accept → done · f follow-up ticket · e edit · E body · esc back", w)))
	return b.String()
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// accept moves the ticket on from the human checkpoint. The move goes through
// the same gate the CLI uses, with Interactive set — which is the one thing an
// agent cannot claim, since it cannot press this key.
func (m *Model) accept() {
	if m.detail == nil {
		return
	}
	id := m.detail.ID
	next := m.lanes.Terminal()
	if next == nil {
		m.notify("no terminal lane is installed, so there is nowhere to accept this to", true)
		return
	}
	env := gate.Env{Lanes: m.lanes, All: m.tickets}
	if vs := gate.CheckAdvance(env, m.detail, gate.Request{
		To: next.ID, Actor: identity(m.store.Root), Interactive: true,
	}); len(vs) > 0 {
		m.notify("Cannot accept yet:\n\n"+vs.Err().Error(), true)
		return
	}
	if _, err := m.store.Mutate(id, func(t *ticket.Ticket) error {
		if err := t.Doc().SetScalar(ticket.FieldStatus, next.ID); err != nil {
			return err
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	}); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.mode = modeBoard
	m.detail = nil
	m.selectByID(id)
	m.notify(fmt.Sprintf("Accepted %s into %s.", ticket.Handle(id), next.Name), false)
}

// followUp creates a new backlog ticket carrying the reviewed ticket's context,
// and links back to it. Rejecting work without recording what is left undone is
// how the reason for a ticket gets lost, which is the failure this board exists
// to prevent.
func (m *Model) followUp() {
	if m.detail == nil {
		return
	}
	src := m.detail
	now := time.Now()
	me := identity(m.store.Root)
	def := m.lanes.Default()

	fields := map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     "Follow-up: " + src.Title,
		ticket.FieldStatus:    def.ID,
		ticket.FieldReady:     "false",
		ticket.FieldCreator:   me,
		ticket.FieldAssignee:  firstNonEmpty(src.Assignee, me),
		ticket.FieldContext:   firstNonEmpty(src.Context, src.Goal),
		ticket.FieldFollows:   src.ID,
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	lists := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}

	body := "# Follow-up: " + src.Title + "\n\n## Description\n\n" +
		"Raised from the review of " + ticket.Handle(src.ID) + ".\n\n"
	if strings.TrimSpace(src.ReviewVerdict) != "" {
		body += "The reviewer said:\n\n> " + strings.ReplaceAll(src.ReviewVerdict, "\n", "\n> ") + "\n\n"
	}
	body += "## Definition of Done\n\n- [ ] <What must be true that is not true yet>\n\n## Notes\n\n"

	t, err := m.store.Create(fields, lists, body)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if full, err := m.store.Load(t.ID); err == nil {
		m.detail = full
		m.startEdit()
	}
}

// atHumanCheckpoint reports whether the open ticket is parked in a lane that
// only a person can release it from.
func (m *Model) atHumanCheckpoint() bool {
	if m.detail == nil {
		return false
	}
	l, ok := m.lanes.Get(m.detail.Status)
	return ok && l.RequiresHumanExit
}
