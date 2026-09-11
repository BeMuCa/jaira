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
	behind := m.behind()
	if m.width <= 0 || m.height <= 0 {
		return behind
	}
	// Leave a margin so the board is visible around all four sides — that
	// margin is what says "this is on top of something", not "this is the
	// screen now".
	maxW := max(20, m.width-8)
	maxH := max(6, m.height-6)
	box := styModal.MaxWidth(maxW).Render(clipHeight(content, maxH-2))

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

// behind is the view the modal is drawn over: the one the dialog will return
// to, so closing it changes nothing but the box going away.
func (m *Model) behind() string {
	switch m.returnTo {
	case modeDetail:
		return m.renderDetail()
	case modeLaneFocus:
		return m.renderLaneFocus()
	case modePipeline:
		return m.renderPipeline()
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
