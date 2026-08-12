package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/project"
	"github.com/berk/jaira/core/session"
	"github.com/berk/jaira/core/ticket"
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

	// Chosen is the board the user picked; the caller opens it.
	Chosen string
	// Quit is set when the user asked to leave rather than pick.
	Quit bool

	width, height int
}

// NewHome builds the launcher from the registry plus anything discovered below
// the working directory.
func NewHome(extraRoots []string) (*Home, error) {
	lanes, err := lane.Load()
	if err != nil {
		return nil, err
	}
	h := &Home{lanes: lanes}

	seen := map[string]bool{}
	add := func(root string) {
		abs, err := filepath.Abs(root)
		if err != nil || seen[abs] {
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
	case tea.KeyPressMsg:
		switch msg.String() {
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
	// Everything below is truncated against the pane, which is the list width
	// when the icon is beside it and the whole terminal when it is not.
	paneW := h.width
	if h.width >= iconWidth+34 {
		paneW = h.width - iconWidth - 2
	}
	paneW = max(12, paneW)

	var right strings.Builder
	right.WriteString(styLaneTitle.Render("jaira") + "\n")
	right.WriteString(styMeta.Render("your boards") + "\n\n")

	if len(h.entries) == 0 {
		right.WriteString(styWarn.Render(truncate("No boards yet.", paneW)) + "\n\n")
		right.WriteString(styMeta.Render(truncate("Run 'jaira init' in a repository, or", paneW)) + "\n")
		right.WriteString(styMeta.Render(truncate("'jaira projects add <path>' to register one.", paneW)) + "\n")
	}

	for i, e := range h.entries {
		lead := "  "
		name := truncate(e.Name, max(1, paneW-2))
		if i == h.idx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		right.WriteString(lead + name + "\n")

		var bits []string
		if e.Err != nil {
			bits = append(bits, styErr.Render("unreadable"))
		} else {
			bits = append(bits, styMeta.Render(fmt.Sprintf("%d open / %d", e.Open, e.Total)))
		}
		// Two different colours because they ask for two different things: work
		// is happening without you, versus it has stopped and is waiting on you.
		if e.Agents > 0 {
			bits = append(bits, styOK.Render(fmt.Sprintf("● %d agent(s)", e.Agents)))
		}
		if e.Waiting > 0 {
			bits = append(bits, styReview.Render(fmt.Sprintf("◆ %d awaiting you", e.Waiting)))
		}
		right.WriteString("    " + truncate(strings.Join(bits, "  "), max(1, paneW-4)) + "\n")
	}

	right.WriteString("\n" + styMeta.Render(truncate("enter open · jk move · q quit", paneW)))

	// The icon sits beside the list when there is room for both, and is dropped
	// entirely when there is not — a launcher that cannot show its list because
	// of decoration would be a bad trade.
	body := right.String()
	if h.width >= iconWidth+34 {
		return lipgloss.JoinHorizontal(lipgloss.Top, iconArt, "  ", body)
	}
	return body
}
