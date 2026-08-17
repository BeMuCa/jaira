package tui

import (
	"fmt"
	"strings"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/ticket"
)

// laneFocusKey handles the single-lane view: one lane, alone, filling the
// screen. h/l change lane without leaving this view — unlike the board,
// changing lane here never switches mode, because staying put is the entire
// point of the screen.
func (m *Model) laneFocusKey(s string) {
	switch s {
	case "j", "down":
		m.moveCard(1)
	case "k", "up":
		m.moveCard(-1)
	case "g":
		m.cardIdx = 0
	case "G":
		if len(m.cols) > 0 {
			m.cardIdx = len(m.cols[m.laneIdx].tickets) - 1
			m.clampCursor()
		}
	case "h", "left":
		m.moveLane(-1)
	case "l", "right":
		m.moveLane(1)
	case "enter":
		m.openDetail()
	case "esc", "v":
		// Back to the compact view it was opened from, never to the board — the
		// board is the sideways view this screen replaces.
		m.mode = modePipeline
	}
}

// renderLaneFocus draws one lane at the terminal's full width. A board column
// is clamped to 22-34 columns, so its card shows only a title, the handle and
// a row of flags. At full width there is room for what those columns cannot
// fit — the goal, and how long the ticket has sat here — without opening the
// detail pane, which is the reason this view exists rather than just widening
// a board column.
func (m *Model) renderLaneFocus() string {
	if len(m.cols) == 0 {
		return "No lanes are installed."
	}
	col := m.cols[m.laneIdx]
	l := col.lane

	var b strings.Builder
	title := l.Name
	if l.Unknown {
		title = "? " + title
	}
	b.WriteString(styLaneTitle.Render(title))
	b.WriteString(" " + styLaneCount.Render(fmt.Sprintf("%d ticket(s)", len(col.tickets))))
	switch {
	case l.Unknown:
		b.WriteString("  " + styWarn.Render("read-only"))
	case l.Agentic:
		tier := l.ModelTier
		if tier == "" {
			tier = "default"
		}
		b.WriteString("  " + styAgentic.Render("agentic · "+tier))
	}
	b.WriteString("\n")
	b.WriteString(styBar.Render(strings.Repeat("─", max(1, m.width))) + "\n\n")

	if len(col.tickets) == 0 {
		b.WriteString(styMeta.Render("no tickets") + "\n")
	} else {
		w := max(20, m.width-2)
		for i, t := range col.tickets {
			b.WriteString(m.renderLaneCard(t, w, i == m.cardIdx))
			b.WriteString("\n")
		}
	}

	footer := styMeta.Render(truncate(
		"enter open · esc/v back · q quit", m.width))
	used := strings.Count(b.String(), "\n") + 1
	if gap := m.height - used - 2; gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	return b.String() + "\n" + footer
}

// renderLaneCard is renderCard widened out. The same status flags the board
// shows keep a ticket looking the same on both screens; the goal and the
// timespan are added because a full-width view has room for the two lines a
// board column never does.
func (m *Model) renderLaneCard(t *ticket.Ticket, w int, selected bool) string {
	marker := "  "
	title := truncate(t.Title, w-2)
	if selected {
		marker = stySelected.Render("▌ ")
		title = stySelected.Render(title)
	}

	meta := styHandle.Render(ticket.Handle(t.ID))
	if t.Assignee != "" {
		meta += styMeta.Render(" " + t.Assignee)
	}
	if when := timespan(t.CreatedAt, t.UpdatedAt); when != "" {
		meta += styMeta.Render("  " + when)
	}

	// Same computation as renderCard's flags, so a ticket does not look
	// different depending on which screen it is read from.
	var flags []string
	env := gate.Env{Lanes: m.lanes, All: m.tickets}
	if !gate.Ready(t) {
		flags = append(flags, styWarn.Render("○ spec"))
	}
	if len(t.BlockedBy) > 0 && !gate.Actionable(env, t) {
		flags = append(flags, styErr.Render("■ blocked"))
	}
	if l, ok := m.lanes.Get(t.Status); ok {
		switch {
		case l.RequiresQuestion:
			flags = append(flags, styAsks.Render("▲ asks"))
		case l.ID == "review":
			flags = append(flags, styReview.Render("◆ sign off"))
		}
	}
	if n, total := checklistProgress(t.PlanItems); total > 0 {
		flags = append(flags, styMeta.Render(fmt.Sprintf("Plan %d/%d", n, total)))
	}
	if n, total := checklistProgress(t.DoDItems); total > 0 {
		flags = append(flags, styMeta.Render(fmt.Sprintf("DoD %d/%d", n, total)))
	}
	if len(t.Commits) > 0 {
		flags = append(flags, styOK.Render(fmt.Sprintf("✓ %d", len(t.Commits))))
	}
	if t.ExecutedBy != "" {
		flags = append(flags, styAgentic.Render(t.ExecutedBy))
	}

	out := marker + title + "\n"
	out += "  " + truncate(meta, w) + "\n"
	if len(flags) > 0 {
		out += "  " + truncate(strings.Join(flags, "  "), w) + "\n"
	}
	if strings.TrimSpace(t.Goal) != "" {
		out += "  " + styMeta.Render(truncate("goal: "+t.Goal, w)) + "\n"
	}
	return out
}
