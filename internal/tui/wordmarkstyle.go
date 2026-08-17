package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// colSunset is the sunset pink the AI in jAIra wears. From the 256-colour cube
// like every other colour here, so it renders the same in terminals without
// truecolour.
var stySunset = lipgloss.NewStyle().Foreground(lipgloss.Color("211")).Bold(true)

// styledWordmark renders the block wordmark with the a and i — the AI in
// jAIra — in sunset pink. The five letters sit at fixed column spans in the
// generated art (see wordmark.go), separated by all-blank columns, so slicing
// per row is safe as long as the art itself is not edited.
func styledWordmark() string {
	const aiStart, aiEnd = 8, 22
	var out []string
	for _, row := range strings.Split(strings.TrimRight(wordmark, "\n"), "\n") {
		// Sliced as runes, not bytes: the block glyph is three bytes wide and
		// one column wide, and the spans are columns.
		r := []rune(row)
		for len(r) < wordmarkWidth {
			r = append(r, ' ')
		}
		out = append(out, styLaneTitle.Render(string(r[:aiStart]))+
			stySunset.Render(string(r[aiStart:aiEnd]))+
			styLaneTitle.Render(string(r[aiEnd:])))
	}
	return strings.Join(out, "\n")
}
