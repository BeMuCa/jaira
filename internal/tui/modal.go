package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// A modal is a window drawn on top of the board rather than instead of it.
//
// Every dialog here used to replace the whole screen, which cost the reader
// the one thing they were reasoning about: a link window that hides the board
// makes you remember which card you were on and which lane it sat in. Keeping
// the board visible behind the box means the answer and the question are on
// screen together.
//
// It is one helper rather than per-dialog drawing because the alternative is
// each window inventing its own geometry, and a box that lands somewhere else
// for every dialog reads as a different program each time.

// styModal frames a modal: a rounded border in the accent colour so the box
// reads as being in front, and padding so the board behind never touches the
// text.
var styModal = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colAccent).
	Padding(0, 1)

// modal centres content in a bordered box over whatever is behind it. The
// content is clipped to the window rather than overflowing it: a box taller
// than the terminal would push the board out of the frame and defeat the
// point.
func (m *Model) modal(content string) string {
	return m.modalOver(content, m.returnTo)
}

// modalOver draws content over a named view rather than over whatever the
// last dialog happened to be returning to. A modal opened from a modal — a
// refusal raised inside the link window — needs to say what it covers, or
// the two dialogs fight over one field and the wrong screen shows through.
func (m *Model) modalOver(content string, under mode) string {
	behind := m.underlying(under)
	if m.width <= 0 || m.height <= 0 {
		return behind
	}
	// Leave a margin so the board is visible around all four sides — that
	// margin is what says "this is on top of something", not "this is the
	// screen now".
	maxW := max(20, m.width-8)
	maxH := max(6, m.height-6)
	// The trailing newline every view ends with counts as a line once the
	// content is split, and one line over the budget is what the clip then
	// eats — the hint line at the bottom, which is the one part of a window
	// a reader cannot work without.
	box := styModal.MaxWidth(maxW).Render(clipHeight(strings.TrimRight(content, "\n"), maxH-2))

	x := max(0, (m.width-lipgloss.Width(box))/2)
	y := max(0, (m.height-lipgloss.Height(box))/2)
	// A Layer only carries its position for a Compositor to read; drawing one
	// straight onto a Canvas ignores x and y and lands everything at the
	// origin, which renders as the box alone on an empty screen.
	return lipgloss.NewCanvas(m.width, m.height).Compose(
		lipgloss.NewCompositor(
			lipgloss.NewLayer(fitCanvas(behind, m.width, m.height)),
			lipgloss.NewLayer(box).X(x).Y(y).Z(1),
		),
	).Render()
}

// underlying renders the view a modal sits on, so closing the modal changes
// nothing but the box going away.
func (m *Model) underlying(md mode) string {
	switch md {
	case modeDetail:
		return m.renderDetail()
	case modeLaneFocus:
		return m.renderLaneFocus()
	case modePipeline:
		return m.renderPipeline()
	case modeLinks:
		// A dialog raised from the link window keeps the window on screen
		// under it, rather than dropping the reader back to the board.
		if m.links != nil {
			return m.modalOver(m.renderLinks(), m.links.from)
		}
	}
	return m.renderBoard()
}

// clipHeight keeps the first n lines and says nothing about the rest — a
// modal that scrolled its own content would be a second scrolling model to
// keep in step with the pane it covers.
func clipHeight(s string, n int) string {
	if n <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= n {
		return s
	}
	kept := append([]string{}, lines[:n-1]...)
	return strings.Join(append(kept, styMeta.Render("…")), "\n")
}

// fitCanvas pads a rendered view out to the full window, so the layer under
// the modal covers the whole canvas instead of leaving holes where the view
// happened to end early.
func fitCanvas(s string, w, h int) string {
	lines := strings.Split(s, "\n")
	for len(lines) < h {
		lines = append(lines, "")
	}
	for i, l := range lines {
		if pad := w - lipgloss.Width(l); pad > 0 {
			lines[i] = l + strings.Repeat(" ", pad)
		}
	}
	return strings.Join(lines[:h], "\n")
}
