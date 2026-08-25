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
	case "esc", "v", "q":
		// Back to the compact view it was opened from, never to the board — the
		// board is the sideways view this screen replaces. q joins esc/v here
		// (rather than quitting, as it does on the board) because lane focus is
		// one level deeper than the compact view, and q is always one back.
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
		cards := make([]string, len(col.tickets))
		heights := make([]int, len(col.tickets))
		for i, t := range col.tickets {
			cards[i] = m.renderLaneCard(t, w, i == m.cardIdx)
			// +1 for the blank separator line the loop below writes after each card.
			heights[i] = strings.Count(cards[i], "\n") + 1
		}

		// The same budget the footer arithmetic below reserves — 2 lines for the
		// gap and the footer itself — measured against what the header above has
		// already used.
		headerLines := strings.Count(b.String(), "\n")
		budget := max(1, m.height-headerLines-2)

		// Reuse the board's own scroll state for this lane: the board recomputes
		// and clamps it for the focused lane on every render anyway, and the two
		// views agreeing on where a long lane is scrolled to is the point.
		first := m.scroll[col.lane.ID]
		if m.cardIdx < first {
			first = m.cardIdx
		}
		if first > len(cards) {
			first = 0
		}
		last := laneFocusFit(heights, first, budget)
		// Keep the cursor visible: page forward until it lands inside the window
		// the budget allows — the same shape as the board's own scroll fix, just
		// walked one card at a time because cards are not a fixed height here.
		for m.cardIdx >= last && first < len(cards)-1 {
			first++
			last = laneFocusFit(heights, first, budget)
		}
		m.scroll[col.lane.ID] = first

		// Nothing is hidden without saying so, above or below.
		if first > 0 {
			b.WriteString(styMeta.Render(fmt.Sprintf(" +%d more", first)) + "\n")
		}
		for i := first; i < last; i++ {
			b.WriteString(cards[i])
			b.WriteString("\n")
		}
		if last < len(cards) {
			b.WriteString(styMeta.Render(fmt.Sprintf(" +%d more", len(cards)-last)) + "\n")
		}
	}

	footer := styMeta.Render("enter open · q/esc/v back")
	used := strings.Count(b.String(), "\n") + 1
	if gap := m.height - used - 2; gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	// Last line of defence: this view can no more spill past the terminal than
	// a board column can.
	return clampBlock(b.String()+"\n"+footer, m.width, m.height)
}

// laneFocusFit returns the index one past the last card, starting at first,
// that fits within budget lines — reserving a line for each off-screen
// indicator ("+N more" above and/or below) that ends up actually shown. That
// reservation is a fixed point rather than a single pass, because whether the
// "more below" indicator is needed can only be known after seeing how many
// cards fit, the same way clipToWindow settles its own footer. At least one
// card is always returned so a budget too small for even the first card still
// shows something; clampBlock is the backstop if that pushes past the window.
func laneFocusFit(heights []int, first, budget int) int {
	reserve := 0
	if first > 0 {
		reserve = 1
	}
	last := first
	for iter := 0; iter < 2; iter++ {
		room := budget - reserve
		last = first
		sum := 0
		for last < len(heights) {
			if last > first && sum+heights[last] > room {
				break
			}
			sum += heights[last]
			last++
		}
		want := 0
		if first > 0 {
			want++
		}
		if last < len(heights) {
			want++
		}
		if want == reserve {
			break
		}
		reserve = want
	}
	return last
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
	env := m.gateEnv()
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
		case l.RequiresHumanExit:
			flags = append(flags, styReview.Render("◆ sign off"))
		}
		if l.Agentic && len(gate.OutputOwed(l, t)) > 0 {
			flags = append(flags, styAgentic.Render("◇ unworked"))
		}
	}
	if !m.isMe(t.UpdatedBy) && strings.TrimSpace(t.UpdatedBy) != "" {
		flags = append(flags, styAsks.Render("✎ "+truncate(t.UpdatedBy, 10)))
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
