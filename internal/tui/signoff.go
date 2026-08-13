package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/ticket"
)

// renderSignOff is the screen a person sees when a ticket is waiting on them.
//
// It answers these questions in order — what was wrong, what the agent did,
// why, whether it actually holds, what the reviewer says it does, what the
// reviewer found missing, and the reviewer's verdict — because that is the
// judgement being asked for. The implementer's account and the reviewer's
// reading of it are shown side by side rather than merged: they are different
// claims, and when they disagree that disagreement is the most useful thing on
// the screen.
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
	section("What the reviewer says it does", t.ReviewSummary)
	section("What the reviewer found missing", t.ReviewGaps)
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
			// A ticked box without its evidence is a claim nobody can check, and
			// this is exactly the screen where that judgement is made.
			if it.Proof != "" {
				fmt.Fprintf(&b, "      %s\n", styMeta.Render(truncate("proof: "+it.Proof, max(1, w-6))))
			}
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
		To: next.ID, Actor: identity.Current(m.store.Root), Interactive: true,
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

// followUpContext builds the new ticket's "why", replacing the previous body
// prose so the reason for the ticket lives only in context, not split between
// a frontmatter field and a body heading. It is naturally multi-line, which is
// the end-to-end proof that a block scalar is readable in the file and in a
// diff, not just in a unit test.
func followUpContext(src *ticket.Ticket) string {
	parts := []string{"Raised from the review of " + ticket.Handle(src.ID) + "."}
	if why := firstNonEmpty(src.Context, src.Goal); strings.TrimSpace(why) != "" {
		parts = append(parts, why)
	}
	if strings.TrimSpace(src.ReviewVerdict) != "" {
		parts = append(parts, "The reviewer said:\n\n> "+strings.ReplaceAll(src.ReviewVerdict, "\n", "\n> "))
	}
	return strings.Join(parts, "\n\n")
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
	me := identity.Current(m.store.Root)
	def := m.lanes.Default()

	fields := map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     "Follow-up: " + src.Title,
		ticket.FieldStatus:    def.ID,
		ticket.FieldReady:     "false",
		ticket.FieldCreator:   me,
		ticket.FieldAssignee:  firstNonEmpty(src.Assignee, me),
		ticket.FieldContext:   followUpContext(src),
		ticket.FieldFollows:   src.ID,
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	lists := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}

	body := "# Follow-up: " + src.Title + "\n\n" +
		"## Definition of Done\n\n- [ ] <What must be true that is not true yet>\n\n## Notes\n\n"

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
