package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

// HomeEntry is one project on the home screen, with enough state to answer the
// question the screen exists for: which of my boards needs me.
type HomeEntry struct {
	Root  string
	Name  string
	Open  int // tickets not in a terminal lane
	Total int

	// Agents is how many sessions are actively working this board, and Waiting
	// how many tickets sit at a human checkpoint. Those are the two states worth
	// colouring: one means work is happening without you, the other means it has
	// stopped for you.
	Agents  int
	Waiting int

	// Logged is this board's logbook entries per day over the last
	// logbookDays, oldest first, today last. Counted from the dated folder
	// names, so listing a board costs no ticket read.
	Logged []int

	Err error
}

// logbookDays is how many days the launcher's activity chart covers.
const logbookDays = 7

// Home is the launcher: the icon, and every board with what it is doing.
type Home struct {
	entries []HomeEntry

	// drop is the board-removal screen, non-nil while it is up. The launcher is
	// where a board most obviously wants removing — it is the list you are
	// looking at when you notice one should not be on it.
	drop  *dropBoard
	idx   int
	lanes *lane.Set

	// browse is the directory picker, non-nil while adding a board.
	browse *browser

	// board is the default board screen, non-nil while it is open. Home is
	// the only per-user surface there is, which is why this hangs off it
	// rather than off any one project's Model.
	board *defaultBoardScreen

	// msg reports what an action did. Registering boards silently looked like
	// nothing had happened whenever they were already in the list.
	msg string

	// stats is whether the activity chart is shown under the list: logbook
	// entries per day over the last week, every board added together. Off by
	// default and session-only, like every other view preference.
	stats bool

	// Chosen is the board the user picked; the caller opens it.
	Chosen string
	// Quit is set when the user asked to leave rather than pick.
	Quit bool

	// startDir is where the directory picker opens, normally where jaira was run.
	startDir string

	// versionLine is the persistent "which version, is there an update"
	// indicator, computed once at construction — see versionLine() in
	// updatecheck.go for why not on every render.
	versionLine string

	width, height int
}

// refresh rebuilds the list, including anything just added, and reports how many
// of them were not already known.
func (h *Home) refresh(extra []string) int {
	seen := map[string]bool{}
	var entries []HomeEntry
	add := func(root string) {
		abs := project.Canonical(root)
		if seen[abs] {
			return
		}
		seen[abs] = true
		entries = append(entries, describe(abs, h.lanes))
	}
	before := map[string]bool{}
	for _, e := range h.entries {
		before[project.Canonical(e.Root)] = true
	}
	for _, p := range project.Load() {
		add(p.Root)
	}
	for _, r := range extra {
		add(r)
	}
	h.entries = entries
	if h.idx >= len(h.entries) {
		h.idx = max(0, len(h.entries)-1)
	}
	added := 0
	for _, e := range entries {
		if !before[project.Canonical(e.Root)] {
			added++
		}
	}
	return added
}

// NewHome builds the launcher from the registry plus anything discovered below
// the working directory.
func NewHome(extraRoots []string) (*Home, error) {
	// The launcher spans every registered board, not one project, so there is
	// no single root to scope lanes to: it can only ever show the catalogue.
	lanes, err := lane.Load("")
	if err != nil {
		return nil, err
	}
	h := &Home{lanes: lanes, versionLine: versionLine()}

	seen := map[string]bool{}
	add := func(root string) {
		abs := project.Canonical(root)
		if seen[abs] {
			return
		}
		seen[abs] = true
		h.entries = append(h.entries, describe(abs, lanes))
	}
	for _, p := range project.Load() {
		add(p.Root)
	}
	for _, r := range extraRoots {
		add(r)
	}
	if wd, err := os.Getwd(); err == nil {
		h.startDir = wd
	}
	return h, nil
}

// describe reads one board's counts. A board that cannot be read is listed with
// its error rather than dropped: a project vanishing silently from the launcher
// is worse than one that says what is wrong with it.
func describe(root string, lanes *lane.Set) HomeEntry {
	e := HomeEntry{Root: root, Name: filepath.Base(root)}
	s, err := ticket.At(root)
	if err != nil {
		e.Err = err
		return e
	}
	tickets, err := s.List()
	if err != nil {
		if pe, ok := err.(*ticket.PartialError); ok {
			e.Err = pe
		} else {
			e.Err = err
			return e
		}
	}
	e.Total = len(tickets)
	e.Logged = s.LoggedPerDay(time.Now(), logbookDays)
	for _, t := range tickets {
		l, known := lanes.Get(t.Status)
		if !known || !l.Terminal {
			e.Open++
		}
		if known && l.RequiresHumanExit {
			e.Waiting++
		}
	}
	if ss, err := session.Load(s); err == nil {
		for _, x := range ss {
			if !x.Stale() {
				e.Agents++
			}
		}
	}
	return e
}

func (h *Home) Init() tea.Cmd { return nil }

func (h *Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width, h.height = msg.Width, msg.Height
	case boardEditorDoneMsg:
		if h.board != nil {
			if msg.err != nil {
				h.board.msg, h.board.isErr = msg.err.Error(), true
			} else {
				h.board.msg, h.board.isErr = "", false
			}
		}
		return h, nil

	case tea.KeyPressMsg:
		if h.drop != nil {
			done, removed := h.drop.key(msg.String())
			if done {
				name := h.drop.name
				h.drop = nil
				if removed {
					h.refresh(nil)
					if h.idx >= len(h.entries) {
						h.idx = max(0, len(h.entries)-1)
					}
					h.msg = "removed " + name
				}
			}
			return h, nil
		}
		if h.board != nil {
			done, cmd := h.board.key(msg.String())
			if done {
				h.board = nil
			}
			return h, cmd
		}
		if h.browse != nil {
			added, done := h.browse.key(msg.String())
			if len(added) > 0 {
				// Say what happened. Registering boards that were already known
				// changes nothing on screen, which reads as the key not working.
				n := h.refresh(added)
				switch {
				case n == 0 && len(added) == 1:
					h.msg = filepath.Base(added[0]) + " was already on the list"
				case n == 0:
					h.msg = fmt.Sprintf("all %d board(s) found were already on the list", len(added))
				case n == len(added):
					h.msg = fmt.Sprintf("added %d board(s)", n)
				default:
					h.msg = fmt.Sprintf("added %d of %d board(s); the rest were already listed", n, len(added))
				}
			}
			if done {
				h.browse = nil
			}
			return h, nil
		}
		switch msg.String() {
		case "a":
			h.msg = ""
			h.browse = newBrowser(h.startDir)
			return h, nil
		case "r":
			h.refresh(nil)
			h.msg = "reloaded"
			return h, nil
		case "s":
			h.stats = !h.stats
			return h, nil
		case "x":
			// Opens a screen rather than acting: removing a board is the one
			// key on this list that cannot be taken back.
			if h.idx < len(h.entries) {
				e := h.entries[h.idx]
				h.msg = ""
				h.drop = newDropBoard(e.Root, e.Name, false)
			}
			return h, nil
		case "d":
			h.msg = ""
			db, err := lane.LoadDefaultBoard()
			if err != nil {
				h.msg = err.Error()
				return h, nil
			}
			h.board = newDefaultBoardScreen(h.lanes, db)
			return h, nil
		case "q", "ctrl+c", "esc":
			h.Quit = true
			return h, tea.Quit
		case "j", "down":
			if h.idx < len(h.entries)-1 {
				h.idx++
			}
		case "k", "up":
			if h.idx > 0 {
				h.idx--
			}
		case "enter":
			if h.idx < len(h.entries) {
				h.Chosen = h.entries[h.idx].Root
				return h, tea.Quit
			}
		}
	}
	return h, nil
}

func (h *Home) View() tea.View {
	v := tea.NewView(h.render())
	v.AltScreen = true
	v.WindowTitle = "jaira"
	return v
}

func (h *Home) render() string {
	if h.width == 0 {
		return "loading…"
	}
	if h.drop != nil {
		return h.drop.render(h.width, h.height)
	}
	if h.board != nil {
		return h.board.render(h.width, h.height)
	}
	if h.browse != nil {
		return h.browse.render(h.width, h.height)
	}
	// Everything below is truncated against the pane, which is the list width
	// when the icon is beside it and the whole terminal when it is not.
	paneW := h.width
	if h.width >= iconWidth+34 {
		paneW = h.width - iconWidth - 2
	}
	paneW = max(12, paneW)

	var b strings.Builder

	// Header: the wordmark with the icon beside it, centred as a block. This is
	// the one screen that is looked at rather than worked in, so it is allowed
	// to spend rows on identity — but only when they fit.
	if head := h.header(); head != "" {
		b.WriteString(head)
		b.WriteString("\n")
		b.WriteString(styBar.Render(strings.Repeat("─", h.width)) + "\n\n")
	}

	centre := lipgloss.NewStyle().Width(h.width).Align(lipgloss.Center)
	b.WriteString(centre.Render(styLaneTitle.Render("Projects:")) + "\n\n")

	if len(h.entries) == 0 {
		b.WriteString(centre.Render(styWarn.Render(truncate("No boards yet.", h.width))) + "\n\n")
		b.WriteString(centre.Render(styMeta.Render(truncate("Press a to find or create one.", h.width))) + "\n")
	}

	for i, e := range h.entries {
		lead := "    "
		name := truncate(e.Name, max(1, h.width-30))
		if i == h.idx {
			lead = "  " + stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		var bits []string
		if e.Err != nil {
			bits = append(bits, styErr.Render("unreadable"))
		} else {
			bits = append(bits, styMeta.Render(fmt.Sprintf("%d open / %d", e.Open, e.Total)))
		}
		// Two colours because they ask for two different things: work happening
		// without you, versus work stopped and waiting on you.
		if e.Agents > 0 {
			bits = append(bits, styOK.Render(fmt.Sprintf("● %d agent(s)", e.Agents)))
		}
		if e.Waiting > 0 {
			bits = append(bits, styReview.Render(fmt.Sprintf("◆ %d awaiting you", e.Waiting)))
		}
		b.WriteString(centre.Render(truncate(lead+padTo(name, 28)+strings.Join(bits, "  "), h.width)) + "\n")
	}

	if h.stats {
		b.WriteString("\n" + h.renderStats(centre))
	}

	if h.msg != "" {
		b.WriteString("\n" + centre.Render(styOK.Render(truncate(h.msg, h.width))) + "\n")
	}

	b.WriteString("\n" + centre.Render(styMeta.Render(truncate(
		"enter open · s stats · a add a board · x remove a board · d default board · r refresh · q quit", h.width))))
	if h.versionLine != "" {
		b.WriteString("\n" + centre.Render(truncate(h.versionLine, h.width)))
	}
	return b.String()
}

// renderStats is the activity chart: logbook entries per day over the last
// logbookDays, every board added together, today on the right and marked. It
// counts what went into the logbook, not what was accepted — a ticket left in
// done is not in it — and the title says so.
func (h *Home) renderStats(centre lipgloss.Style) string {
	days := make([]int, logbookDays)
	for _, e := range h.entries {
		for i, n := range e.Logged {
			days[i] += n
		}
	}
	peak := 0
	for _, n := range days {
		peak = max(peak, n)
	}

	const rows = 4
	const colW = 4
	today := logbookDays - 1
	// Every row is exactly logbookDays*colW wide, so centring lands each
	// column in the same place on every line.
	tone := func(i int) lipgloss.Style {
		if i == today {
			return stySelected
		}
		return styOK
	}

	var b strings.Builder
	b.WriteString(centre.Render(styMeta.Render("logbook entries per day · last 7 days · all boards")) + "\n")

	// Counts above the bars, so a day is readable without judging bar heights.
	var line strings.Builder
	for i, n := range days {
		cell := ""
		if n > 0 {
			cell = fmt.Sprintf("%d", n)
		}
		line.WriteString(tone(i).Render(padTo(cell, colW)))
	}
	b.WriteString(centre.Render(line.String()) + "\n")

	// Bars scale to the busiest day; a day with anything at all gets at least
	// one row, so one ticket is never drawn as nothing.
	for r := rows; r >= 1; r-- {
		line.Reset()
		for i, n := range days {
			cell := "  "
			if n > 0 && (n*rows+peak-1)/peak >= r {
				cell = "██"
			}
			line.WriteString(tone(i).Render(padTo(cell, colW)))
		}
		b.WriteString(centre.Render(line.String()) + "\n")
	}

	// The day of the month under each bar; today is the marked one.
	line.Reset()
	now := time.Now()
	for i := range days {
		day := now.AddDate(0, 0, i-today).Format("2")
		line.WriteString(tone(i).Render(padTo(day, colW)))
	}
	b.WriteString(centre.Render(line.String()) + "\n")
	return b.String()
}

// header renders the wordmark and icon side by side, centred, or nothing at all
// when the terminal is too small to carry them without crowding the list.
func (h *Home) header() string {
	if h.width < wordmarkWidth+iconWidth+6 || h.height < 20 {
		if h.width >= wordmarkWidth+2 && h.height >= 14 {
			return lipgloss.NewStyle().Width(h.width).Align(lipgloss.Center).
				Render(styledWordmark())
		}
		return ""
	}
	block := lipgloss.JoinHorizontal(lipgloss.Center,
		styledWordmark(), "    ", iconArt)
	return lipgloss.NewStyle().Width(h.width).Align(lipgloss.Center).Render(block)
}

// padTo pads a styled string to a column count, measuring display width so the
// escape sequences in it are not counted.
func padTo(s string, n int) string {
	if w := lipgloss.Width(s); w < n {
		return s + strings.Repeat(" ", n-w)
	}
	return s + " "
}
