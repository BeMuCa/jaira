package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// The art is generated, so the thing worth testing is that it still matches what
// the layout code assumes: a fixed width, no stray blank margin, and no line that
// would push the pane sideways.
func TestIconArtMatchesItsDeclaredWidth(t *testing.T) {
	lines := strings.Split(strings.TrimRight(iconArt, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("the icon is empty")
	}
	for i, l := range lines {
		if w := lipgloss.Width(stripANSI(l)); w > iconWidth {
			t.Errorf("row %d is %d cols, wider than the declared %d: %q", i, w, iconWidth, stripANSI(l))
		}
	}
	if strings.TrimSpace(stripANSI(lines[0])) == "" {
		t.Error("the icon starts with a blank row")
	}
	if strings.TrimSpace(stripANSI(lines[len(lines)-1])) == "" {
		t.Error("the icon ends with a blank row")
	}
}

func TestIconArtIsDrawnWithBlockGlyphs(t *testing.T) {
	if !strings.ContainsAny(iconArt, "▀▄█") {
		t.Error("the icon contains no block glyphs, so it is not the rendered art")
	}
	if !strings.Contains(iconArt, "\x1b[") {
		t.Error("the icon carries no colour")
	}
}
