package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Colours are chosen from the 256-colour cube rather than truecolour so the
// board stays legible in terminals that do not advertise 24-bit colour, and
// every foreground is paired with a shape or label so the UI never depends on
// colour alone to convey state.
var (
	colDim     = lipgloss.Color("244")
	colFaint   = lipgloss.Color("240")
	colAccent  = lipgloss.Color("39")
	colWarn    = lipgloss.Color("214")
	colErr     = lipgloss.Color("203")
	colOK      = lipgloss.Color("78")
	colAgentic = lipgloss.Color("141")

	styLaneTitle = lipgloss.NewStyle().Bold(true)
	styLaneCount = lipgloss.NewStyle().Foreground(colFaint)
	styHandle    = lipgloss.NewStyle().Foreground(colFaint)
	styMeta      = lipgloss.NewStyle().Foreground(colDim)
	stySelected  = lipgloss.NewStyle().Foreground(colAccent).Bold(true)
	styWarn      = lipgloss.NewStyle().Foreground(colWarn)
	styErr       = lipgloss.NewStyle().Foreground(colErr)
	styOK        = lipgloss.NewStyle().Foreground(colOK)
	styAgentic   = lipgloss.NewStyle().Foreground(colAgentic)
	styHelpKey   = lipgloss.NewStyle().Foreground(colAccent)
	styBar       = lipgloss.NewStyle().Foreground(colDim)
	styAsks      = lipgloss.NewStyle().Foreground(colWarn).Bold(true)
	styReview    = lipgloss.NewStyle().Foreground(colAccent).Bold(true)
	styDoing     = lipgloss.NewStyle().Foreground(colAgentic).Bold(true)
)

const minColWidth = 22

// View satisfies tea.Model.
func (m *Model) View() tea.View {
	v := tea.NewView(m.render())
	// In v2 the alternate screen is declared on the view rather than requested
	// with a command.
	v.AltScreen = true
	v.WindowTitle = "jaira"
	return v
}

func (m *Model) render() string {
	if m.width == 0 {
		return "loading…"
	}
	// A follow-up being written beside its predecessor owns both halves of the
	// screen, whether the right half is still the editor or already a ticket.
	if m.follow != nil && (m.mode == modeDetail || m.mode == modeEdit) {
		return m.renderSplit()
	}
	switch m.mode {
	case modeHelp:
		return m.renderHelp()
	case modeDetail:
		// A ticket parked at a human checkpoint opens to the screen that asks for
		// the decision, rather than to the generic field dump.
		if m.atHumanCheckpoint() {
			return m.renderSignOff()
		}
		return m.renderDetail()
	case modeEdit:
		return m.renderEdit()
	case modePipeline:
		return m.renderPipeline()
	case modeLaneFocus:
		return m.renderLaneFocus()
	case modeMessage:
		return m.renderMessage()
	case modeProjects:
		return m.renderProjects()
	case modeLanes:
		return m.laneScreen.render(m.width, m.height)
	case modeSettings:
		return m.settingsScreen.render(m.width, m.height)
	case modeDefaultBoard:
		return m.board.render(m.width, m.height)
	case modeMove:
		return m.renderBoard() // picker is drawn into the status bar
	case modeDelete:
		return m.renderDelete()
	}
	return m.renderBoard()
}

// boardWindow is what boardFit decides: which of m.cols render, in what
// slice, and how wide. Its doc lives on boardFit, since the two are one
// decision split only for return-type reasons.
type boardWindow struct {
	cols   []int // indices into m.cols, in board order, with hidden lanes already dropped
	start  int   // first rendered position within cols
	end    int   // one past the last rendered position within cols
	hidden int   // how many lanes were dropped because they hold no tickets
	colW   int
}

// boardFit is the single place that decides which lanes render and how wide.
// Empty lanes are dropped first, and the window applies to what remains —
// one function, so the order is a fact of the code rather than a convention
// two call sites have to remember.
func (m *Model) boardFit(width int) boardWindow {
	var cols []int
	hidden := 0
	for i, c := range m.cols {
		// The focused lane is kept even when it is empty: switching the
		// toggle on must not make the cursor point at something invisible.
		if !m.hideEmpty || len(c.tickets) > 0 || i == m.laneIdx {
			cols = append(cols, i)
			continue
		}
		hidden++
	}

	// The columns that do fit stretch to fill the row: a capped column width
	// left the right third of a wide terminal blank, which read as wasted
	// screen rather than as a decision.
	perScreen := max(1, width/(minColWidth+2))
	shown := min(len(cols), perScreen)
	// A bordered column occupies its content width plus two border cells.
	colW := max(minColWidth, width/shown-2)

	// cursor is m.laneIdx's position within cols, not within m.cols — the
	// window is measured over the lanes actually on the board. Absent only if
	// the keep-rule above somehow missed it, which it does not.
	cursor := 0
	for i, ci := range cols {
		if ci == m.laneIdx {
			cursor = i
			break
		}
	}

	// Centred on the cursor, over the lanes actually on the board. The window
	// is measured here rather than by renderBoard because hiding empty lanes
	// changes what "the lane after this one" means: with the toggle on, the
	// successors worth showing are the ones that hold something.
	start, end := laneWindow(cursor, len(cols), perScreen)

	return boardWindow{cols: cols, start: start, end: end, hidden: hidden, colW: colW}
}

func (m *Model) renderBoard() string {
	var b strings.Builder
	b.WriteString(m.header())
	b.WriteString("\n")
	tabs := m.boardProjectLine()
	if tabs != "" {
		b.WriteString(tabs + "\n")
	}
	if panel := m.renderSessions(); panel != "" {
		b.WriteString(panel)
	}

	if len(m.cols) == 0 {
		b.WriteString("\n  No lanes are installed.\n")
		return b.String()
	}

	win := m.boardFit(m.width)

	tabsLine := 0
	if tabs != "" {
		tabsLine = 1
	}
	// The status bar is rendered first because it may wrap onto several lines
	// on a narrow terminal, and the columns get whatever height remains.
	sb := m.statusBar(win)
	sbLines := strings.Count(sb, "\n") + 1
	bodyHeight := m.height - 4 - sbLines - tabsLine - sessionPanelHeight(m.sessions)
	if bodyHeight < 3 {
		bodyHeight = 3
	}

	rendered := make([]string, 0, win.end-win.start)
	for _, ci := range win.cols[win.start:win.end] {
		rendered = append(rendered, m.renderColumn(ci, win.colW, bodyHeight))
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rendered...))
	b.WriteString("\n")
	b.WriteString(sb)
	return b.String()
}

// renderSessions shows what each agent session is working on: the board's
// window onto memory that would otherwise vanish when a session ends.
func (m *Model) renderSessions() string {
	if len(m.sessions) == 0 {
		return ""
	}
	var b strings.Builder
	for _, x := range m.sessions {
		if x.Focus == "" && x.TicketID == "" {
			continue
		}
		marker, style := "●", styOK
		if x.Stale() {
			// Dimmed and labelled rather than removed: a crashed session's last
			// known focus is still the most useful thing to show.
			marker, style = "○", styMeta
		}
		line := style.Render(marker + " " + truncate(x.Focus, max(20, m.width/2)))
		if x.TicketID != "" {
			line += styHandle.Render(" " + ticket.Handle(x.TicketID))
		}
		if x.Model != "" {
			line += styAgentic.Render(" " + x.Model)
		}
		if x.Reasoning != "" {
			line += styMeta.Render("  — " + truncate(x.Reasoning, max(10, m.width/3)))
		}
		if x.Stale() {
			line += styMeta.Render(fmt.Sprintf("  (quiet %s)", roughAge(x.Age())))
		}
		b.WriteString(truncate(line, m.width) + "\n")
	}
	if b.Len() == 0 {
		return ""
	}
	return b.String()
}

func sessionPanelHeight(ss []session.Session) int {
	n := 0
	for _, x := range ss {
		if x.Focus != "" || x.TicketID != "" {
			n++
		}
	}
	return n
}

func roughAge(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "moments"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// timespan says when a ticket appeared and when it was last touched. The two
// answer different questions — how old the thought is, and whether anyone has
// been near it since — and a ticket where those diverge is exactly the one worth
// noticing.
func timespan(created, updated time.Time) string {
	if created.IsZero() {
		return ""
	}
	s := created.Local().Format("2 Jan 2006") + " · " + roughAge(time.Since(created)) + " old"
	// A minute of slack: every ticket is written a moment after it is created, and
	// reporting that as a separate event would put "touched moments ago" on every
	// card ever made.
	if !updated.IsZero() && updated.After(created.Add(time.Minute)) {
		s += ", touched " + roughAge(time.Since(updated)) + " ago"
	}
	return s
}

// laneWindow is the run of lanes drawn when perScreen of n of them fit, with
// idx — the focused one — inside it.
//
// The focused lane sits in the middle of the row, not at its right edge.
// Edge-anchoring meant the lanes *after* the focused one were never on screen —
// and what comes after a lane is where its work goes next, which is the question
// the board is being read to answer.
//
// Clamped at both ends so the row is always full: padding one side with blanks
// to centre a first or last lane would trade the same information away again,
// for symmetry nobody asked for.
func laneWindow(idx, n, perScreen int) (start, end int) {
	start = max(0, min(idx-perScreen/2, n-perScreen))
	return start, min(n, start+perScreen)
}

func (m *Model) renderColumn(idx, w, h int) string {
	col := m.cols[idx]
	focused := idx == m.laneIdx

	border := lipgloss.NormalBorder()
	style := lipgloss.NewStyle().Border(border).BorderForeground(colFaint).Width(w).Height(h)
	if focused {
		style = style.BorderForeground(colAccent)
	}

	// Lane headers carry no agentic marker. Every lane can be driven by an agent
	// — creating, moving and filling in tickets are all CLI operations — so
	// marking two of them said less than it implied. What is worth marking is on
	// the cards: a ticket that needs an answer, or one waiting to be signed off.
	title := col.lane.Name
	if col.lane.Unknown {
		title = "? " + title
	}
	head := styLaneTitle.Render(truncate(title, w-6)) + " " + styLaneCount.Render(fmt.Sprintf("%d", len(col.tickets)))

	var body strings.Builder
	body.WriteString(head + "\n")
	if col.lane.Unknown {
		body.WriteString(styWarn.Render(truncate("read-only", w-4)) + "\n")
	}
	body.WriteString(styBar.Render(strings.Repeat("─", max(1, w-4))) + "\n")

	// Keep the cursor visible within a lane taller than the pane.
	visible := max(1, (h-4)/3)
	first := m.scroll[col.lane.ID]
	if focused {
		if m.cardIdx < first {
			first = m.cardIdx
		}
		if m.cardIdx >= first+visible {
			first = m.cardIdx - visible + 1
		}
		m.scroll[col.lane.ID] = first
	}
	if first > len(col.tickets) {
		first = 0
	}

	for i := first; i < len(col.tickets) && i < first+visible; i++ {
		body.WriteString(m.renderCard(col.tickets[i], w-4, focused && i == m.cardIdx))
	}
	if rest := len(col.tickets) - (first + visible); rest > 0 {
		body.WriteString(styMeta.Render(fmt.Sprintf(" +%d more", rest)))
	}
	return style.Render(clampBlock(body.String(), w, h))
}

func (m *Model) renderCard(t *ticket.Ticket, w int, selected bool) string {
	marker := "  "
	title := truncate(t.Title, w-2)
	if selected {
		marker = stySelected.Render("▌ ")
		title = stySelected.Render(title)
	}

	// State is shown with a glyph plus a word, never colour alone.
	var flags []string
	env := m.gateEnv()
	if !gate.Ready(t) {
		flags = append(flags, styWarn.Render("○ spec"))
	}
	if len(t.BlockedBy) > 0 && !gate.Actionable(env, t) {
		flags = append(flags, styErr.Render("■ blocked"))
	}
	// Two lanes both mean "waiting on you", but they ask for different things —
	// answering a question, or signing off finished work — so they are labelled
	// and coloured apart rather than collapsed into one "needs attention".
	if l, ok := m.lanes.Get(t.Status); ok {
		switch {
		case l.RequiresQuestion:
			flags = append(flags, styAsks.Render("▲ asks"))
		case l.RequiresHumanExit:
			flags = append(flags, styReview.Render("◆ sign off"))
		}
		// A lane that still owes its declared output has not been run on this
		// ticket. Without this, nine tickets sitting in a critique lane
		// uncritiqued looked exactly like the two that had been through it, and
		// the only way to tell was to open each one and find the field empty.
		//
		// Only for a lane whose work an agent does: 'todo' and 'backlog' declare
		// no output and would never light up anyway, but a hand-written lane may,
		// and the claim being made here is specifically that the agent has not
		// run.
		if l.Agentic && len(gate.OutputOwed(l, t)) > 0 {
			flags = append(flags, styAgentic.Render("◇ unworked"))
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
		flags = append(flags, styAgentic.Render(truncate(t.ExecutedBy, 8)))
	}

	meta := styHandle.Render(ticket.Handle(t.ID))
	if t.Assignee != "" {
		// Someone else's ticket is marked, not merely named: after a pull the
		// board can hold teammates' claims, and whose lane a card sits in is
		// exactly what you scan for before touching it. @ plus colour, so the
		// state never rests on colour alone.
		if strings.EqualFold(t.Assignee, m.me) {
			meta += styMeta.Render(" " + truncate(t.Assignee, max(3, w-14)))
		} else {
			meta += styWarn.Render(" @" + truncate(t.Assignee, max(3, w-15)))
		}
	}

	out := marker + title + "\n"
	out += "  " + truncate(meta, w) + "\n"
	if len(flags) > 0 {
		out += "  " + truncate(strings.Join(flags, " "), w) + "\n"
	} else {
		out += "\n"
	}
	return out
}

// checklistProgress counts items that are no longer outstanding. An item in
// progress counts as unfinished and a superseded one counts as settled, matching
// the gate: the point of the counter is to show how much work is left, not how
// much has been started. Counting only ticks would leave a ticket that can be
// completed sitting at 4/5 for good.
func checklistProgress(items []ticket.DoDItem) (done, total int) {
	for _, it := range items {
		if it.Settled() {
			done++
		}
	}
	return done, len(items)
}

// renderChecklist prints one checklist with its state markers, arrowing the item
// being worked on. That arrow is the answer to "which step is the agent on",
// which is the question the board exists to make answerable at a glance.
func renderChecklist(b *strings.Builder, label string, items []ticket.DoDItem, width int) {
	if len(items) == 0 {
		return
	}
	done, total := checklistProgress(items)
	// The same two-column shape as every other field: label left, content
	// right, so the checklists read as fields of the ticket rather than as a
	// second document below it. The label line carries the progress count; the
	// items follow in the content column, keeping their state markers, the
	// arrow on the item being worked on, and the proof line under its item.
	//
	// A pane too narrow for the label column falls back to items at the left
	// edge — the column layout must never be the thing that pushes a line past
	// the pane, which is what the width guard below is for.
	itemCol := 13
	if width < 40 {
		itemCol = 0
	}
	pad := strings.Repeat(" ", itemCol)
	fmt.Fprintf(b, "\n%s %s\n", styMeta.Render(fmt.Sprintf("%-12s", label)),
		styLaneCount.Render(fmt.Sprintf("%d/%d", done, total)))
	for _, it := range items {
		lead, sty := "  ", styMeta
		switch it.State {
		case ticket.StateDoing:
			lead, sty = styDoing.Render("→ "), styDoing
		case ticket.StateDone:
			sty = styOK
		}
		fmt.Fprintf(b, "%s%s%s %s\n", pad, lead,
			sty.Render("["+it.State.Marker()+"]"), wrap(it.Text, max(1, width-itemCol-6), itemCol+6))
		if it.Proof != "" {
			fmt.Fprintf(b, "%s      %s\n", pad,
				styMeta.Render(wrap("proof: "+it.Proof, max(1, width-itemCol-6), itemCol+13)))
		}
	}
}

// renderBodySections renders what remains of the body — Options, Progress,
// and whatever headings a hand edit added — in the same two-column shape as
// every other field: the heading as the left label, its lines in the content
// column. Raw "## Heading" markdown between styled fields read as a second
// document glued below the ticket. A heading with nothing under it is skipped
// entirely; the seeded empty Progress section says nothing yet.
func renderBodySections(b *strings.Builder, body string, width int) {
	label := ""
	var content []string
	flush := func() {
		filled := false
		for _, l := range content {
			if strings.TrimSpace(l) != "" {
				filled = true
				break
			}
		}
		if !filled {
			label, content = "", nil
			return
		}
		b.WriteString("\n")
		first := true
		for _, l := range content {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			// A bullet's dash is layout, not content, once the line sits in a
			// column of its own.
			l = strings.TrimPrefix(l, "- ")
			lead := strings.Repeat(" ", 13)
			if first {
				lead = styMeta.Render(fmt.Sprintf("%-12s", strings.ToLower(label))) + " "
				first = false
			}
			b.WriteString(lead + wrap(l, max(10, width-14), 13) + "\n")
		}
		label, content = "", nil
	}
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			flush()
			label = strings.TrimSpace(strings.TrimLeft(trimmed, "# "))
			continue
		}
		content = append(content, line)
	}
	flush()
}

// dropLeadingTitle removes the body's opening "# <title>" heading. Every new
// ticket's body starts with one (see ticket.NewBody), but the pane already
// names the ticket in its header — rendered raw it reads as a mystery section
// between the checklists and the notes.
func dropLeadingTitle(rest string) string {
	lines := strings.Split(rest, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimLeft(strings.Join(lines[i+1:], "\n"), "\n")
		}
		break
	}
	return rest
}

// stripChecklistSections removes the sections rendered above, so the remaining
// body can be shown without printing the same checklists twice in two formats.
func stripChecklistSections(body string) string {
	var out []string
	skipping := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			heading := strings.ToLower(strings.TrimLeft(trimmed, "# "))
			skipping = false
			for _, h := range append(append([]string{}, ticket.PlanHeadings()...), ticket.DoDHeadings()...) {
				if strings.Contains(heading, h) {
					skipping = true
					break
				}
			}
			if skipping {
				continue
			}
		}
		if !skipping {
			out = append(out, line)
		}
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

// boardProjectLine is a single line of numbered board names above the
// columns, so 1-9 means the same board here as in the compact view. The board
// has far less spare room than the compact view's own banner, so a line that
// does not fit is dropped entirely rather than wrapped — a wrapped second line
// would misalign the bordered columns drawn under it.
func (m *Model) boardProjectLine() string {
	if len(m.projects) < 2 {
		return ""
	}
	line := strings.Join(m.projectTabs(), styBar.Render("  │  "))
	// Truncated rather than dropped: a narrow terminal used to hide the whole
	// line, so the boards — and which one you are in — silently disappeared
	// instead of merely being cut short.
	return truncate(line, m.width)
}

func (m *Model) header() string {
	// The board's name is dropped when the project line below already carries it,
	// marked and in first position: the same word twice on two consecutive lines
	// is noise. With a single recorded board that line does not render, so the
	// name stays here rather than disappearing entirely.
	var left string
	if m.boardProjectLine() == "" {
		name := "jaira"
		if root := m.store.Root; root != "" {
			parts := strings.Split(root, string(os.PathSeparator))
			name = parts[len(parts)-1]
		}
		left = styLaneTitle.Render(name)
	}
	if m.filter != "" {
		left += styMeta.Render(fmt.Sprintf("   filter: %q", m.filter))
	}
	total := 0
	for _, c := range m.cols {
		total += len(c.tickets)
	}
	right := styMeta.Render(fmt.Sprintf("%d tickets", total))
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	if lipgloss.Width(left)+lipgloss.Width(right)+1 > m.width {
		return truncate(left, m.width)
	}
	return left + strings.Repeat(" ", gap) + right
}

func (m *Model) statusBar(win boardWindow) string {
	if m.mode == modeMove {
		var parts []string
		for i, l := range m.lanes.Lanes {
			label := l.Name
			if i == m.moveTarget {
				label = stySelected.Render("[" + label + "]")
			} else {
				label = styMeta.Render(label)
			}
			parts = append(parts, label)
		}
		return truncate("move to: "+strings.Join(parts, " "), m.width)
	}
	if m.mode == modeFilter {
		return truncate(styHelpKey.Render("/")+m.input+stySelected.Render("▏"), m.width)
	}
	if m.mode == modeCreate {
		return truncate(styHelpKey.Render("new: ")+m.input+stySelected.Render("▏"), m.width)
	}

	// Hints are dropped from the right as the terminal narrows, rather than
	// letting the bar wrap and push the board off-screen.
	// "m move" rather than "m lane": the key moves the ticket, and calling it
	// after its destination read as if m selected a lane to look at.
	// A toggle names what the next press will do, which is the only way a
	// hint bar can describe one honestly.
	zHint := "z hide empty"
	if m.hideEmpty {
		zHint = "z show empty"
	}
	keys := []string{"enter open", "v compact", zHint, "n new", "m move", "S settings", "/ filter", "? help", "q quit"}
	prefix := ""
	if len(m.warnings) > 0 {
		prefix += styWarn.Render(fmt.Sprintf("⚠ %d ", len(m.warnings)))
	}
	// Wrapped, never dropped: a key the bar has no room for is a key the reader
	// does not know exists. renderBoard measures this bar and gives the columns
	// whatever height is left.
	lines := wrapHints(keys, max(1, m.width-lipgloss.Width(prefix)))
	for i, l := range lines {
		if i == 0 {
			lines[i] = prefix + styMeta.Render(l)
		} else {
			lines[i] = styMeta.Render(l)
		}
	}
	// The window/hidden-lanes notice, if either applies. Nothing is appended
	// when neither does, so a wide terminal with nothing hidden renders
	// byte-for-byte what it renders today. This sits above m.versionLine,
	// which stays last.
	if notice := boardNotice(win); notice != "" {
		lines = append(lines, styMeta.Render(truncate(notice, m.width)))
	}
	if m.versionLine != "" {
		lines = append(lines, truncate(m.versionLine, m.width))
	}
	return strings.Join(lines, "\n")
}

// boardNotice states what boardFit hid or scrolled off-screen. A reader who
// cannot see the blocked lane, or the ninth of ten lanes, must not be left to
// infer nothing is there — this is the one line that says so, and it is
// empty when nothing was hidden and nothing was scrolled.
func boardNotice(win boardWindow) string {
	if win.hidden == 0 {
		return ""
	}
	return fmt.Sprintf("%d lane(s) hidden — z to show", win.hidden)
}

// renderDetail draws the open ticket at the full width of the terminal.
func (m *Model) renderDetail() string {
	if m.detail == nil {
		return m.renderBoard()
	}
	return m.clipToWindow(m.detailBody(m.detail, max(20, m.width)), m.detailHints(m.detail)+" · esc back")
}

// renderDelete keeps the ticket on screen while its handle is typed back. What
// is about to be destroyed stays in front of the person destroying it.
func (m *Model) renderDelete() string {
	if m.detail == nil {
		return m.renderBoard()
	}
	handle := ticket.Handle(m.detail.ID)
	prompt := styWarn.Render("delete "+handle+" — type the handle: ") + m.input + stySelected.Render("▏")
	return m.clipToWindow(m.detailBody(m.detail, max(20, m.width)), prompt+" · esc cancel")
}

// detailHints names the actions an open ticket offers. Basic movement (arrows,
// jk, paging) is deliberately not listed: the footer names actions, the help
// screen teaches movement. b appears only when there is a blocker to jump to.
func (m *Model) detailHints(t *ticket.Ticket) string {
	hint := "e fields · E editor · y copy id · m move · n follow-up · X delete"
	if len(t.BlockedBy) > 0 {
		hint += " · b blocked-by"
	}
	return hint
}

// detailBody builds an open ticket's content at a given width, so the same
// rendering serves the full screen and one half of the split follow-up view.
func (m *Model) detailBody(t *ticket.Ticket, width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s\n", styHandle.Render(ticket.Handle(t.ID)), styLaneTitle.Render(t.Title))
	b.WriteString(styBar.Render(strings.Repeat("─", max(1, width))) + "\n")

	row := func(k, v string) {
		if strings.TrimSpace(v) == "" {
			return
		}
		fmt.Fprintf(&b, "%s %s\n", styMeta.Render(fmt.Sprintf("%-12s", k)), wrap(v, max(10, width-14), 13))
	}
	if m.copied {
		row("id", t.ID+styOK.Render("  copied"))
	} else {
		row("id", t.ID)
	}
	if l, ok := m.lanes.Get(t.Status); ok && l.RequiresQuestion {
		row("lane", t.Status+styAsks.Render("  ▲ waiting on your answer"))
	} else if ok && l.RequiresHumanExit {
		row("lane", t.Status+styReview.Render("  ◆ waiting on your sign-off"))
	} else {
		row("lane", t.Status)
	}
	row("assignee", t.Assignee)
	row("creator", t.Creator)
	// A ticket that sat in a lane for three weeks and one touched an hour ago look
	// identical without this. The board is read to answer "where does this stand",
	// and how long it has stood there is part of that answer — especially for a
	// thought captured in one session and picked up in another.
	row("when", timespan(t.CreatedAt, t.UpdatedAt))
	row("executed-by", t.ExecutedBy)
	row("tier", t.ModelTier)
	// The identity rows above stay tight; the prose fields below each get a
	// blank line, because two wrapped paragraphs with no gap between them read
	// as one — which field a sentence belongs to should not need re-reading.
	prose := func(k, v string) {
		if strings.TrimSpace(v) == "" {
			return
		}
		b.WriteString("\n")
		row(k, v)
	}
	prose("goal", t.Goal)
	prose("context", t.Context)
	// The checklist below carries the same label; showing the one-line scalar
	// too would print "done when" twice for every ticket that has both.
	if len(t.DoDItems) == 0 {
		prose("done when", t.DoD)
	}
	if len(t.BlockedBy) > 0 {
		var hs []string
		for _, d := range t.BlockedBy {
			hs = append(hs, ticket.Handle(d))
		}
		prose("blocked by", strings.Join(hs, ", "))
	}
	if t.Follows != "" {
		prose("follows", ticket.Handle(t.Follows))
	}
	// Shown only while the ticket is parked, same as the CLI: a stale reason
	// rendered on an active ticket reads as its current state.
	if l, ok := m.lanes.Get(t.Status); ok && l.RequiresBlockedReason {
		prose("waiting on", t.BlockedReason)
	}
	prose("question", t.Question)

	if t.Outcome.What != "" || t.Outcome.Why != "" || t.Outcome.Resolves != "" {
		b.WriteString("\n" + styLaneTitle.Render("Outcome") + "\n")
		row("what", t.Outcome.What)
		row("why", t.Outcome.Why)
		row("resolves", t.Outcome.Resolves)
	}
	if t.ReviewSummary != "" || t.ReviewGaps != "" || t.ReviewVerdict != "" || t.ReviewCheck != "" {
		b.WriteString("\n" + styLaneTitle.Render("Review") + "\n")
		row("summary", t.ReviewSummary)
		row("gaps", t.ReviewGaps)
		row("verdict", t.ReviewVerdict)
		row("check", t.ReviewCheck)
	}
	if len(t.Commits) > 0 {
		b.WriteString("\n" + styLaneTitle.Render("Commits") + "\n")
		if stat, err := (&gitStat{root: m.store.Root}).of(t.Commits); err == nil && stat != "" {
			b.WriteString(styMeta.Render(stat) + "\n")
		} else {
			row("commits", strings.Join(t.Commits, " "))
		}
	}
	if miss := missing(t); len(miss) > 0 {
		b.WriteString("\n" + styWarn.Render("Before this can start: "+strings.Join(miss, ", ")) + "\n")
	}
	renderChecklist(&b, "plan", t.PlanItems, width)
	renderChecklist(&b, "done when", t.DoDItems, width)

	renderBodySections(&b, dropLeadingTitle(stripChecklistSections(t.Body)), width)
	return b.String()
}

// clipPane clips one pane of the split to a fixed height, clamping the scroll it
// is handed and padding short content so two panes joined side by side keep
// their rows lined up. It returns the clamped scroll, because only the renderer
// knows how long the content came out. The last row is spent saying how much is
// below rather than showing one more line of it — nothing off-screen goes
// unnamed here either.
// The block comes out exactly width x height, for the reason clampBlock spells
// out: a sized lipgloss style wraps overlong lines onto extra rows and pads
// without clipping, so a pane that is not already the right shape silently
// changes the shape of the layout around it.
func clipPane(content string, width, height, scroll int) (string, int) {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	height = max(1, height)
	if scroll > len(lines)-height {
		scroll = len(lines) - height
	}
	if scroll < 0 {
		scroll = 0
	}
	out := make([]string, 0, height)
	if len(lines) > scroll+height {
		end := scroll + height - 1
		out = append(out, lines[scroll:end]...)
		out = append(out, styMeta.Render(fmt.Sprintf(" +%d more", len(lines)-end)))
	} else {
		out = append(out, lines[scroll:min(len(lines), scroll+height)]...)
	}
	for len(out) < height {
		out = append(out, "")
	}
	// Every row is cut and then padded to exactly width, so the block is a true
	// rectangle and the border drawn around it cannot be ragged.
	for i, l := range out {
		l = truncate(l, width)
		if pad := width - lipgloss.Width(l); pad > 0 {
			l += strings.Repeat(" ", pad)
		}
		out[i] = l
	}
	return strings.Join(out, "\n"), scroll
}

// clipToWindow clips an open ticket's content to the terminal, applying and
// clamping the scroll offset, and appends the footer hint — prefixed with how
// much is hidden below, so a short terminal never silently swallows the rest.
// The content of an open ticket has no upper bound in length; rendering it
// whole pushed the handle, the title and the goal off the top of the terminal
// with no key that could bring them back. The offset is clamped here rather
// than in the key handler because only the renderer knows how long the content
// came out.
func (m *Model) clipToWindow(content, hint string) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	// The footer wraps rather than truncates — a key the terminal is too narrow
	// to show is a key the reader does not know exists — so the space reserved
	// for it is however many lines it needs, plus the blank line above it.
	width := max(1, m.width)
	items := strings.Split(hint, " · ")
	footer := wrapHints(items, width)
	visible := max(1, m.height-1-len(footer))

	// The scroll range is defined against the plain footer: at the very bottom
	// nothing is hidden below, so nothing extra is in the footer there.
	if m.detailScroll > len(lines)-visible {
		m.detailScroll = len(lines) - visible
	}
	if m.detailScroll < 0 {
		m.detailScroll = 0
	}
	end := min(len(lines), m.detailScroll+visible)

	if end < len(lines) {
		// Something is hidden below, so the footer leads with how much. The
		// extra item can wrap the footer onto one more line, which shrinks the
		// window again — one recomputation settles it, and the count is written
		// against the window actually shown.
		for range 2 {
			more := fmt.Sprintf("+%d more", len(lines)-end)
			footer = wrapHints(append([]string{more}, items...), width)
			visible = max(1, m.height-1-len(footer))
			end = min(len(lines), m.detailScroll+visible)
		}
	}
	var b strings.Builder
	b.WriteString(strings.Join(lines[m.detailScroll:end], "\n") + "\n")
	for _, l := range footer {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return b.String()
}

// wrapHints lays key hints into as many lines as the width needs, breaking at
// the " · " separators, so a narrow terminal shows every key instead of
// silently dropping or truncating the tail.
func wrapHints(items []string, width int) []string {
	var lines []string
	cur := ""
	for _, it := range items {
		if it == "" {
			continue
		}
		cand := it
		if cur != "" {
			cand = cur + " · " + it
		}
		if cur != "" && lipgloss.Width(cand) > width {
			lines = append(lines, cur)
			cur = it
			continue
		}
		cur = cand
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

type gitStat struct{ root string }

func (g *gitStat) of(shas []string) (string, error) {
	args := append([]string{"-C", g.root, "show", "--stat=120", "--format=%h %s"}, shas[0])
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// colorizeDiff highlights a patch without pulling in a syntax-highlighting
// dependency: for a diff, the leading character carries all the meaning.

func (m *Model) renderMessage() string {
	style := styOK
	label := "note"
	if m.isErr {
		style, label = styErr, "refused"
	}
	var b strings.Builder
	b.WriteString(style.Render(strings.ToUpper(label)) + "\n\n")
	b.WriteString(m.message + "\n\n")
	switch {
	case m.pending == nil:
		b.WriteString(styMeta.Render("esc dismiss"))
	case m.pending.confirm:
		b.WriteString(styMeta.Render("y override · n cancel"))
	default:
		b.WriteString(styMeta.Render("f override · esc dismiss"))
	}
	return b.String()
}

func (m *Model) renderProjects() string {
	var b strings.Builder
	b.WriteString(styLaneTitle.Render("Switch board") + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 78))) + "\n\n")
	for i, p := range m.projects {
		marker := "  "
		name := p.Name
		if i == m.projIdx {
			marker = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		cur := ""
		if p.Root == m.store.Root {
			cur = styMeta.Render("  (current)")
		}
		b.WriteString(marker + name + cur + "\n")
		b.WriteString("    " + styMeta.Render(truncate(p.Root, m.width-6)) + "\n")
	}
	b.WriteString("\n" + styMeta.Render("enter switch · esc back"))
	return b.String()
}

func (m *Model) renderHelp() string {
	var b strings.Builder
	b.WriteString(styLaneTitle.Render("jaira") + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 78))) + "\n\n")

	sections := []struct {
		name string
		keys [][2]string
	}{
		{"Move around", [][2]string{
			{"h l ← →", "previous / next lane"},
			{"j k ↓ ↑", "previous / next card"},
			{"g G", "first / last card in lane"},
		}},
		{"Look at things", [][2]string{
			{"enter", "open the selected ticket"},
			{"↓ ↑", "scroll an open ticket; jk jump to the next/previous one"},
			{"b", "open the ticket this one is blocked by (follow the chain)"},
			{"/", "filter tickets as you type; key:value narrows to one field"},
			{"esc", "clear the filter"},
			{"y", "copy the full ticket id (detail pane)"},
			{"z", "hide lanes with no tickets (press again to bring them back)"},
		}},
		{"Change things", [][2]string{
			{"n", "create a ticket, then fill it in straight away"},
			{"n", "on an open ticket: write its follow-up beside it (esc discards)"},
			{"tab", "in that split: the other pane; shift+↓↑ scrolls the left one"},
			{"e", "edit fields in the detail pane (enter newline, ctrl+s save)"},
			{"E", "open the ticket body and checklists in $EDITOR"},
			{"a f", "accept, or raise a follow-up (on a ticket awaiting sign-off)"},
			{"m", "move the selected ticket to another lane"},
			{"f", "override a move the gate refused, after confirming with y"},
			{"x", "archive the selected ticket (restore brings it back)"},
			{"X", "on an open ticket: delete its file, after typing the handle back"},
			{"r", "reload from disk now"},
			{"p", "switch to another board"},
			{"S", "settings: lanes (read a prompt, use it, publish it) and the default board"},
		}},
		{"Compact view", [][2]string{
			{"v", "the whole flow at a glance, agents counted per step"},
			{"enter", "open the highlighted step: one lane, filling the screen"},
			{"h l ← →", "in that view, next / previous lane, without leaving it"},
			{"q esc v", "back to the compact view"},
		}},
	}
	for _, s := range sections {
		b.WriteString(styLaneTitle.Render(s.name) + "\n")
		for _, k := range s.keys {
			fmt.Fprintf(&b, "  %s  %s\n", styHelpKey.Render(fmt.Sprintf("%-9s", k[0])), styMeta.Render(k[1]))
		}
		b.WriteString("\n")
	}

	b.WriteString(styMeta.Render(
		"Gates are enforced here exactly as on the command line: a ticket needs a goal,\n" +
			"a definition of done, context and an assignee before it can leave the backlog.\n"))

	if len(m.warnings) > 0 {
		b.WriteString("\n" + styWarn.Render("Warnings") + "\n")
		for _, w := range m.warnings {
			b.WriteString("  " + styWarn.Render("⚠ ") + truncate(w, m.width-4) + "\n")
		}
	}
	b.WriteString("\n" + styMeta.Render("esc back"))
	return b.String()
}

func missing(t *ticket.Ticket) []string {
	var out []string
	if strings.TrimSpace(t.Goal) == "" {
		out = append(out, "goal")
	}
	// A checklist in the body counts, exactly as it does for the gate. Checking
	// only the scalar told a ticket carrying a full checklist that it still
	// needed a definition of done — the board contradicting the rules it renders.
	if !t.HasDoD() {
		out = append(out, "definition-of-done")
	}
	if strings.TrimSpace(t.Context) == "" {
		out = append(out, "context")
	}
	if strings.TrimSpace(t.Assignee) == "" {
		out = append(out, "assignee")
	}
	return out
}

// truncate cuts to a display width, counting grapheme width rather than bytes so
// wide characters and emoji do not break column alignment.
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= n {
		return s
	}
	r := []rune(s)
	for len(r) > 0 && lipgloss.Width(string(r)+"…") > n {
		r = r[:len(r)-1]
	}
	return string(r) + "…"
}

// clampBlock cuts a rendered block to a hard w x h budget. lipgloss's Width()
// wraps overlong lines onto extra rows and Height() pads without clipping, so
// a block handed to a sized style can silently grow the layout around it —
// which is exactly how a flag-heavy review lane once pushed the board past the
// bottom of the terminal. Clamping here means over-tall content is cut, never
// spilled.
func clampBlock(s string, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for i, l := range lines {
		lines[i] = truncate(l, w)
	}
	return strings.Join(lines, "\n")
}

func wrap(s string, width, indent int) string {
	if width <= 8 {
		return s
	}
	words := strings.Fields(s)
	var lines []string
	cur := ""
	for _, w := range words {
		if cur == "" {
			cur = w
			continue
		}
		if lipgloss.Width(cur+" "+w) > width {
			lines = append(lines, cur)
			cur = w
			continue
		}
		cur += " " + w
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return strings.Join(lines, "\n"+strings.Repeat(" ", indent))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
