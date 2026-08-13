package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	Err error
}

// Home is the launcher: the icon, and every board with what it is doing.
type Home struct {
	entries []HomeEntry
	idx     int
	lanes   *lane.Set

	// browse is the directory picker, non-nil while adding a board.
	browse *browser

	// board is the default board screen, non-nil while it is open. Home is
	// the only per-user surface there is, which is why this hangs off it
	// rather than off any one project's Model.
	board *defaultBoardScreen

	// msg reports what an action did. Registering boards silently looked like
	// nothing had happened whenever they were already in the list.
	msg string

	// Chosen is the board the user picked; the caller opens it.
	Chosen string
	// Quit is set when the user asked to leave rather than pick.
	Quit bool

	// startDir is where the directory picker opens, normally where jaira was run.
	startDir string

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
	h := &Home{lanes: lanes}

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

	if h.msg != "" {
		b.WriteString("\n" + centre.Render(styOK.Render(truncate(h.msg, h.width))) + "\n")
	}

	b.WriteString("\n" + centre.Render(styMeta.Render(truncate(
		"enter open · jk move · a add a board · d default board · r refresh · q quit", h.width))))
	return b.String()
}

// header renders the wordmark and icon side by side, centred, or nothing at all
// when the terminal is too small to carry them without crowding the list.
func (h *Home) header() string {
	if h.width < wordmarkWidth+iconWidth+6 || h.height < 20 {
		if h.width >= wordmarkWidth+2 && h.height >= 14 {
			return lipgloss.NewStyle().Width(h.width).Align(lipgloss.Center).
				Render(styLaneTitle.Render(wordmark))
		}
		return ""
	}
	block := lipgloss.JoinHorizontal(lipgloss.Center,
		styLaneTitle.Render(wordmark), "    ", iconArt)
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
