package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/BeMuCa/jaira/core/ticket"
)

// recentMove is how long a lane counts as having just received something.
//
// The arrow into it is lit for that window, which is the whole point of the
// compact view: you are watching several agents at once and want to see work
// move without reading counts.
const recentMove = 20 * time.Second

// pipelineStep is one lane, reduced to what the compact view shows.
type pipelineStep struct {
	id      string
	name    string
	tickets int
	agents  int
	waiting bool // stopped, needs a person
	asking  bool // stopped, needs an answer
	fresh   bool // something arrived here just now
}

func (m *Model) pipelineSteps() []pipelineStep {
	now := time.Now()
	byLane := map[string][]*ticket.Ticket{}
	for _, t := range m.tickets {
		byLane[t.Status] = append(byLane[t.Status], t)
	}
	// Sessions say which ticket each agent is on, which is how a step knows an
	// agent is working in it rather than merely having tickets parked there.
	agentOn := map[string]int{}
	for _, s := range m.sessions {
		if s.Stale() || s.TicketID == "" {
			continue
		}
		for _, t := range m.tickets {
			if t.ID == s.TicketID {
				agentOn[t.Status]++
			}
		}
	}

	var out []pipelineStep
	for _, l := range m.lanes.Lanes {
		ts := byLane[l.ID]
		st := pipelineStep{id: l.ID, name: l.Name, tickets: len(ts), agents: agentOn[l.ID]}
		st.waiting = l.RequiresHumanExit && len(ts) > 0
		st.asking = l.RequiresQuestion && len(ts) > 0
		for _, t := range ts {
			if now.Sub(t.UpdatedAt) < recentMove {
				st.fresh = true
			}
		}
		out = append(out, st)
	}
	return out
}

// stepStyle picks the colour. The two states worth distinguishing are the ones
// that ask something of you, and they ask different things: an answer, or a
// decision on finished work.
func stepStyle(s pipelineStep, focused bool) lipgloss.Style {
	switch {
	case s.asking:
		return styAsks
	case s.waiting:
		return styReview
	case s.agents > 0:
		return styOK
	case focused:
		return stySelected
	case s.tickets == 0:
		return styMeta
	}
	return styLaneTitle
}

// renderStep draws one step as a rounded cell: name, a dot per ticket, and the
// agent count when any are working in it.
func renderStep(s pipelineStep, focused bool) string {
	const inner = 13
	sty := stepStyle(s, focused)

	dots := ""
	switch {
	case s.tickets == 0:
		dots = styMeta.Render("·")
	case s.tickets <= 6:
		dots = sty.Render(strings.Repeat("●", s.tickets))
	default:
		// Beyond a handful the dots stop being countable, so the number carries it.
		dots = sty.Render("●●●●●") + styMeta.Render(fmt.Sprintf("+%d", s.tickets-5))
	}

	tail := ""
	if s.agents > 0 {
		tail = styAgentic.Render(fmt.Sprintf("◆%d", s.agents))
	}

	border := lipgloss.RoundedBorder()
	box := lipgloss.NewStyle().Border(border).Width(inner).Align(lipgloss.Center)
	if focused {
		box = box.BorderForeground(colAccent)
	} else {
		box = box.BorderForeground(colFaint)
	}
	// The name must never wrap: a two-line cell is taller than its neighbours and
	// the whole row stops lining up.
	body := sty.Render(truncate(s.name, inner-2)) + "\n" + dots
	if tail != "" {
		body += "  " + tail
	}
	return box.Render(body)
}

// projectTabs renders each recorded board as "N name", numbered for the 1-9
// switch, marking the open one and any board with a live session. Shared by
// the compact view's banner and the board view's thinner line above the
// columns, so the two agree on what a number means and on what "live" looks
// like without two copies to keep in step.
func (m *Model) projectTabs() []string {
	var tabs []string
	for i, p := range m.projects {
		if i >= 9 {
			break
		}
		sty := styMeta
		label := fmt.Sprintf("%d %s", i+1, p.Name)
		if p.Root == m.store.Root {
			// The board you are standing in is marked with a glyph, not only with
			// a colour: the same rule the rest of the board follows, and colour
			// alone is easy to miss on a line that also carries the live marker.
			sty = stySelected
			label = "▸ " + label
		}
		tab := sty.Render(label)
		// Liveness is marked with a glyph, not colour alone, matching the rule
		// the rest of the board follows for state.
		if m.liveBoards[p.Root] {
			tab += styAgentic.Render(" ◆")
		}
		tabs = append(tabs, tab)
	}
	return tabs
}

// switchToProject sets SwitchTo when n (1-based, as a person types it) names a
// recorded board other than the one already open. Both the compact view and
// the board view call this, so a number means the same board everywhere.
// Switching to the board already open, or to a number beyond the list, is a
// no-op rather than a reload.
func (m *Model) switchToProject(n int) bool {
	i := n - 1
	if i < 0 || i >= len(m.projects) {
		return false
	}
	if m.projects[i].Root == m.store.Root {
		return false
	}
	m.SwitchTo = m.projects[i].Root
	return true
}

// renderPipeline is the compact overview: the whole flow on one screen.
func (m *Model) renderPipeline() string {
	var b strings.Builder

	// Projects across the top, switchable by number. One keystroke per board is
	// the point — this view exists to be glanced at while several agents run.
	if len(m.projects) > 0 {
		tabs := m.projectTabs()
		b.WriteString(truncate(strings.Join(tabs, styBar.Render("  │  ")), m.width) + "\n")
		b.WriteString(styBar.Render(strings.Repeat("─", m.width)) + "\n\n")
	}

	steps := m.pipelineSteps()
	if len(steps) == 0 {
		return "No lanes are installed."
	}

	// Steps run in a serpentine: left to right, then down and back right to left,
	// so the path stays continuous instead of jumping from the end of one row to
	// the start of the next. A wrapped row reads as a break in the flow; a snake
	// reads as the flow continuing.
	const cellW, arrowW = 15, 5
	perRow := max(1, (m.width+arrowW)/(cellW+arrowW))

	var rows []string
	for start := 0; start < len(steps); start += perRow {
		end := min(len(steps), start+perRow)
		idx := make([]int, 0, end-start)
		for i := start; i < end; i++ {
			idx = append(idx, i)
		}
		reversed := (start/perRow)%2 == 1
		if reversed {
			for l, r := 0, len(idx)-1; l < r; l, r = l+1, r-1 {
				idx[l], idx[r] = idx[r], idx[l]
			}
		}

		var parts []string
		for n, i := range idx {
			if n > 0 {
				// The arrow points the way the row runs, and lights when the step it
				// leads into has just received work.
				glyph, lit := " ──▶ ", " ━━▶ "
				if reversed {
					glyph, lit = " ◀── ", " ◀━━ "
				}
				arrow := styBar.Render(glyph)
				target := i
				if reversed {
					target = idx[n]
				}
				if steps[target].fresh {
					arrow = styOK.Render(lit)
				}
				parts = append(parts, lipgloss.NewStyle().MarginTop(1).Render(arrow))
			}
			parts = append(parts, renderStep(steps[i], i == m.laneIdx))
		}
		row := lipgloss.JoinHorizontal(lipgloss.Center, parts...)
		rows = append(rows, row)

		// The connector between rows sits under the end the snake turns at.
		if end < len(steps) {
			pad := 0
			if !reversed {
				pad = max(0, lipgloss.Width(row)-cellW/2-1)
			} else {
				pad = cellW / 2
			}
			rows = append(rows, strings.Repeat(" ", pad)+styBar.Render("▼"))
		}
	}
	b.WriteString(strings.Join(rows, "\n") + "\n")

	if st := m.focusedStep(steps); st != nil {
		b.WriteString("\n" + styLaneTitle.Render(st.name) + " ")
		b.WriteString(styMeta.Render(fmt.Sprintf("· %d ticket(s)", st.tickets)))
		if st.agents > 0 {
			b.WriteString(styAgentic.Render(fmt.Sprintf(" · %d agent(s) working here", st.agents)))
		}
		if st.waiting {
			b.WriteString(styReview.Render(" · waiting on you"))
		}
		if st.asking {
			b.WriteString(styAsks.Render(" · needs an answer"))
		}
		b.WriteString("\n")
	}

	body := b.String()
	footer := styMeta.Render(truncate(
		"enter open step · 1-9 switch project · v full board · q quit", m.width))

	// Push the hints to the bottom. Without this they sat directly under the
	// diagram with the rest of the terminal blank beneath them, which reads as
	// the view having failed to fill the screen.
	used := strings.Count(body, "\n") + 1
	if gap := m.height - used - 2; gap > 0 {
		body += strings.Repeat("\n", gap)
	}
	return body + "\n" + footer
}

func (m *Model) focusedStep(steps []pipelineStep) *pipelineStep {
	if m.laneIdx < 0 || m.laneIdx >= len(steps) {
		return nil
	}
	return &steps[m.laneIdx]
}

// pipelineKey handles the compact view. Navigation is wasd as well as the arrow
// keys, because this is the view you sit in with one hand.
func (m *Model) pipelineKey(s string) bool {
	switch s {
	case "a", "left", "h":
		if m.laneIdx > 0 {
			m.laneIdx--
		}
	case "d", "right", "l":
		if m.laneIdx < len(m.lanes.Lanes)-1 {
			m.laneIdx++
		}
	case "w", "up", "k":
		m.laneIdx = 0
	case "s", "down", "j":
		m.laneIdx = len(m.lanes.Lanes) - 1
	case "enter":
		// Opening a step goes deeper into that one lane, filling the screen — not
		// sideways into the multi-column board, which is a different view of
		// everything rather than more of what you picked.
		m.mode = modeLaneFocus
	case "v":
		m.mode = modeBoard
	default:
		if len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
			return m.switchToProject(int(s[0] - '0'))
		}
		return false
	}
	return false
}
