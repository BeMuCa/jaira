package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
)

// The whole point of the two mix levels: on a palette-only terminal a quiet
// tint does not exist, so the glow must be mixed far enough to survive the
// snap. If any tag lands on a grey, that tag's card says nothing about itself.
func TestGlowKeepsItsColourOnAPaletteTerminal(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.trueColor = false

	// The colours .jaira/tags hands out, across the palette.
	for _, tag := range []int{170, 45, 111, 33, 63, 214, 141, 78} {
		_, params := m.glowFor(tag)
		idx := strings.TrimPrefix(params, "5;")
		if idx == params {
			t.Errorf("tag %d: glow is %q, want a palette colour on a palette terminal", tag, params)
			continue
		}
		var n int
		for _, c := range idx {
			n = n*10 + int(c-'0')
		}
		if r, g, b := paletteRGB(n); r == g && g == b {
			t.Errorf("tag %d: glow snapped to grey %d, so the tag's colour is gone", tag, n)
		}
	}
}

// A 24-bit terminal gets the mix itself, with no palette in the way — that is
// the only way the quiet tint is reachable at all.
func TestGlowIsTrueColourWhenTheTerminalCan(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.Update(tea.ColorProfileMsg{Profile: colorprofile.TrueColor})
	if !m.trueColor {
		t.Fatal("board did not take the terminal's 24-bit colour profile")
	}
	if _, params := m.glowFor(45); !strings.HasPrefix(params, "2;") {
		t.Errorf("glow is %q, want 24-bit parameters", params)
	}
}

// Whatever the tag and whatever the terminal, the glow has to stay lighter than
// the lane around it — it is the selection first and a tag colour second.
func TestGlowStaysAboveTheLaneShades(t *testing.T) {
	for _, trueColor := range []bool{false, true} {
		m := newTestModel(t, 120, 40)
		m.trueColor = trueColor
		for _, tag := range []int{170, 45, 111, 33, 63, 214, 141, 78, 16} {
			r, g, b := paletteRGB(tag)
			mix := glowMix256
			if trueColor {
				mix = glowMixTrue
			}
			for _, c := range []int{r, g, b} {
				got := glowBase + int(mix*float64(c-glowBase))
				if got < 16 {
					t.Errorf("tag %d (trueColor=%v): a channel mixes to %d, below the lane's own shades",
						tag, trueColor, got)
				}
			}
		}
	}
}

// A ticket with no coloured tag has nothing to glow with, so it keeps the
// neutral fill rather than picking an arbitrary hue.
func TestGlowFallsBackToNeutralWithoutATag(t *testing.T) {
	m := newTestModel(t, 120, 40)
	if _, params := m.selectionFill(0, false); params != "5;"+selBgDark {
		t.Errorf("untagged selection is %q, want the neutral fill %q", params, "5;"+selBgDark)
	}
}

// c is the toggle, and it changes only the hue — the fill itself stays, so the
// cursor is no harder to find with the glow off.
func TestGlowTogglesWithC(t *testing.T) {
	m := newTestModel(t, 120, 40)
	if !m.glow {
		t.Fatal("glow is off on a fresh board; it ships on")
	}

	m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	if m.glow {
		t.Error("c did not turn the glow off")
	}
	if _, params := m.selectionFill(45, true); params != "5;"+selBgDark {
		t.Errorf("glow off still tints the fill: %q", params)
	}

	m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	if !m.glow {
		t.Error("c did not turn the glow back on")
	}
	if _, params := m.selectionFill(45, true); params == "5;"+selBgDark {
		t.Error("glow on left the fill neutral")
	}
}

// The hint names what the next press does, the way the thin-empty hint beside
// it does — a toggle labelled with its current state reads backwards.
func TestGlowHintNamesTheNextPress(t *testing.T) {
	m := newTestModel(t, 120, 40)
	if out := m.View().Content; !strings.Contains(out, "c plain") {
		t.Error("with the glow on, the hint does not offer to turn it off")
	}
	m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	if out := m.View().Content; !strings.Contains(out, "c glow") {
		t.Error("with the glow off, the hint does not offer to turn it on")
	}
}
