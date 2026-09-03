package tui

import (
	"fmt"
	"strings"

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
	w := max(20, m.width)
	var b strings.Builder

	fmt.Fprintf(&b, "%s  %s\n", styHandle.Render(ticket.Handle(t.ID)),
		styleLines(styLaneTitle, wrap(t.Title, max(10, w-10), 8)))
	b.WriteString(styReview.Render(truncate("◆ waiting on your sign-off", w)) + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	// The same shape as the detail pane — label column left, text right, a
	// blank line between fields — so the two screens read as one tool. The
	// order is still the four questions a sign-off answers: what was wrong,
	// what was done and why, does it solve it, and what the reviewer made of it.
	section := func(label, body string) {
		if strings.TrimSpace(body) == "" {
			return
		}
		b.WriteString("\n")
		fmt.Fprintf(&b, "%s %s\n", styMeta.Render(fmt.Sprintf("%-12s", label)),
			wrapField(body, min(w-14, 64), 13))
	}
	// An empty section used to vanish, which on this screen is the worst
	// possible answer: the person is being asked to accept work, and a lane
	// that never ran left no trace at all. A field an installed lane declares
	// it produces keeps its place and names the lane that owes it.
	owed := gate.OwedBy(m.lanes, t)
	srcs := gate.DeclaredBy(m.lanes)
	// 78 is the pane width that leaves a debt row the 64 columns of text
	// section() gives a value, beside the 13-column label — so the rows line
	// up with the sections around them.
	const paneWidth = 78
	declared := func(label, field, value string) {
		if strings.TrimSpace(value) != "" {
			section(label, sourced(srcs, field, value))
			return
		}
		if l, ok := owed[field]; ok {
			b.WriteString("\n")
			owedRow(&b, label, l, min(w, paneWidth))
		}
	}
	// Before the seven questions, not among them: a tag says which subject the
	// work belongs to, which is context for every answer below rather than one
	// of them.
	section("tags", strings.Join(t.Tags, " "))
	// A ticket with no goal but a context has a problem statement, so this row
	// shows it and the goal's debt goes unmentioned here — deliberately: the
	// question this screen asks is what was wrong, and it has an answer.
	declared("problem", ticket.FieldGoal, firstNonEmpty(t.Goal, t.Context))
	declared("what", ticket.FieldOutcomeWhat, t.Outcome.What)
	declared("why", ticket.FieldOutcomeWhy, t.Outcome.Why)
	declared("resolves", ticket.FieldOutcomeResolves, t.Outcome.Resolves)
	declared("summary", ticket.FieldReviewSummary, t.ReviewSummary)
	declared("gaps", ticket.FieldReviewGaps, t.ReviewGaps)
	declared("verdict", ticket.FieldReviewVerdict, t.ReviewVerdict)
	// Last, because it is the one section that asks the reader to do something
	// rather than to read: everything above is the account, this is the check.
	declared("check", ticket.FieldReviewCheck, t.ReviewCheck)
	// After the seven, never among them: those labels are the order this
	// judgement is made in, and a field a board's own lane declares cannot be
	// slotted into it without guessing where it belongs. Appended, it is still
	// on the screen — which is the point, because the gate refuses the move on
	// it and the person accepting the work would otherwise never see it.
	for _, f := range m.laneFields(t, owed) {
		b.WriteString("\n")
		if strings.TrimSpace(f.value) != "" {
			fieldRow(&b, f.label, sourced(srcs, f.field, f.value), min(w, paneWidth))
			continue
		}
		owedRow(&b, f.label, owed[f.field], min(w, paneWidth))
	}

	// What is being accepted: the same Commits block the detail pane shows.
	// This is the screen where a person judges shipped work, and an account
	// with no pointer to the change sends them back out to find it. A ticket
	// here normally has NO commits recorded — only the done lane demands
	// them, so the field fills at acceptance, after this screen — which is
	// why an empty field falls back to deriving the list from git, once per
	// ticket (memoised on the model; per render it would exec on every
	// keypress). The heading says so: a derived list is git's account, not
	// the ticket's.
	shas, derived := t.Commits, false
	if len(shas) == 0 {
		if m.derivedFor != t.ID {
			m.derivedFor, m.derivedShas = t.ID, m.gateEnv().DeriveCommits(t)
		}
		shas, derived = m.derivedShas, len(m.derivedShas) > 0
	}
	if len(shas) > 0 {
		b.WriteString("\n" + styLaneTitle.Render("Commits"))
		if derived {
			b.WriteString(styMeta.Render("  derived from git — recorded at acceptance"))
		}
		b.WriteString("\n")
		if stat, err := (&gitStat{root: m.store.Root}).of(shas); err == nil && stat != "" {
			b.WriteString(styMeta.Render(wrapLines(stat, max(10, min(w, paneWidth)))) + "\n")
		} else {
			fieldRow(&b, "commits", strings.Join(shas, " "), min(w, paneWidth))
		}
	}

	if done, total := checklistProgress(t.DoDItems); total > 0 {
		b.WriteString("\n" + styLaneTitle.Render("Definition of Done") +
			" " + styLaneCount.Render(fmt.Sprintf("%d/%d", done, total)) + "\n")
		for _, it := range t.DoDItems {
			sty := styMeta
			if it.Checked() {
				sty = styOK
			}
			fmt.Fprintf(&b, "  %s %s\n", sty.Render("["+it.State.Marker()+"]"),
				wrap(it.Text, max(1, w-6), 6))
			// A ticked box without its evidence is a claim nobody can check, and
			// this is exactly the screen where that judgement is made.
			if it.Proof != "" {
				fmt.Fprintf(&b, "      %s\n", styleLines(styMeta, wrap("proof: "+it.Proof, max(1, w-13), 13)))
			}
		}
	}

	// Clipped and scrolled the same way the detail pane is: a review with a
	// long verdict and a long checklist is exactly the ticket that outgrows a
	// small terminal, and this is the screen where it must be readable.
	return m.clipToWindow(b.String(),
		"a accept → done · f follow-up ticket · e edit · E editor · esc back")
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
	env := m.gateEnv()
	if vs := gate.CheckAdvance(env, m.detail, gate.Request{
		To: next.ID, Actor: identity.Current(m.store.Root),
		ActorAliases: identity.Aliases(m.store.Root), Interactive: true,
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
	// The accept key is the fourth way a ticket lands in a lane, and the
	// usual one for the capped terminal lane — it enforces the cap exactly
	// as the other three move write-sites do.
	trimMsg, trimErr := m.trimHolds(next.ID, id)
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.mode = modeBoard
	m.detail = nil
	m.selectByID(id)
	msg := fmt.Sprintf("Accepted %s into %s.", ticket.Handle(id), next.Name)
	if trimMsg != "" {
		msg += "\n\n" + trimMsg
	}
	if trimErr != nil {
		m.notify(msg+"\n\ntrimming the "+next.ID+" lane failed: "+trimErr.Error(), true)
		return
	}
	m.notify(msg, false)
}

// followUpContext builds the new ticket's "why", replacing the previous body
// prose so the reason for the ticket lives only in context, not split between
// a frontmatter field and a body heading. It is naturally multi-line, which is
// the end-to-end proof that a block scalar is readable in the file and in a
// diff, not just in a unit test.
// lead is the first sentence: only the sign-off path may claim a review happened.
func followUpContext(src *ticket.Ticket, lead string) string {
	parts := []string{lead}
	// The commits are written into the prose as well as being reachable through
	// follows:, because the context has to still answer "what was already done"
	// after the predecessor has been archived off the board.
	if len(src.Commits) > 0 {
		parts = append(parts, "That work shipped in "+strings.Join(src.Commits, ", ")+".")
	}
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
	fields, lists, body := m.followUpFields(src, "Raised from the review of "+ticket.Handle(src.ID)+".")

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
