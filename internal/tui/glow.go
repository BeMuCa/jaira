package tui

import (
	"image/color"
	"strconv"

	"charm.land/lipgloss/v2"
)

// The selected card's fill can carry the colour of the ticket's own tag rather
// than a neutral grey — the glow. It is a second reading of the same cell: the
// fill says where the cursor is, and its hue says what the ticket is about,
// without spending a column or a glyph on either.
//
// Two mix levels, not one, because the 256-colour palette cannot express a
// quiet tint. Below roughly 45% mix every tag colour snaps onto the palette's
// grey ramp — 238, 239, 240 — since it holds no low-saturation dark colours
// between that ramp and its colour cube. A terminal that can show 24-bit
// colour gets the quiet tint it can actually render; one that cannot gets a
// louder mix that still survives the snap, instead of a grey that silently
// drops the tag.
const (
	glowMixTrue = 0.15
	glowMix256  = 0.45

	// glowBase is the grey the tag colour is mixed up from, the value of
	// 256-colour index 237. Mixing down from the tag towards black instead let
	// a dark tag land on the neutral fill's own shade, and the selection
	// stopped being visible at all.
	glowBase = 58
)

// cubeLevels are the six values each channel takes in the 256-colour palette's
// 6×6×6 cube.
var cubeLevels = [6]int{0, 95, 135, 175, 215, 255}

// paletteRGB is the RGB a 256-colour index stands for. The sixteen system
// colours have no fixed definition — a terminal theme decides them — so they
// are reported as a mid grey rather than guessed at.
func paletteRGB(n int) (int, int, int) {
	switch {
	case n >= 16 && n <= 231:
		n -= 16
		return cubeLevels[n/36], cubeLevels[(n/6)%6], cubeLevels[n%6]
	case n >= 232 && n <= 255:
		v := 8 + 10*(n-232)
		return v, v, v
	default:
		return 128, 128, 128
	}
}

// nearestPalette is the 256-colour index closest to an RGB triple. It searches
// the cube and the grey ramp only, for the same reason paletteRGB will not
// report the system colours: their actual values are not knowable here.
func nearestPalette(r, g, b int) int {
	best, bestD := 16, 1<<30
	for i := 16; i <= 255; i++ {
		cr, cg, cb := paletteRGB(i)
		if d := (cr-r)*(cr-r) + (cg-g)*(cg-g) + (cb-b)*(cb-b); d < bestD {
			best, bestD = i, d
		}
	}
	return best
}

// glowFor is the fill for a selected card whose tag has colour tag, as both the
// colour lipgloss takes and the SGR parameters refill writes. On a 24-bit
// terminal that is the mixed colour itself; otherwise the nearest palette entry
// to a louder mix.
func (m *Model) glowFor(tag int) (color.Color, string) {
	r, g, b := paletteRGB(tag)
	k := glowMix256
	if m.trueColor {
		k = glowMixTrue
	}
	mix := func(c int) int { return glowBase + int(k*float64(c-glowBase)) }
	mr, mg, mb := mix(r), mix(g), mix(b)

	if m.trueColor {
		return color.RGBA{R: uint8(mr), G: uint8(mg), B: uint8(mb), A: 0xff},
			"2;" + strconv.Itoa(mr) + ";" + strconv.Itoa(mg) + ";" + strconv.Itoa(mb)
	}
	idx := strconv.Itoa(nearestPalette(mr, mg, mb))
	return lipgloss.Color(idx), "5;" + idx
}

// selectionFill is what the selected card is painted with: the tag's glow when
// the glow is on and the ticket has a coloured tag, and the neutral fill
// otherwise.
func (m *Model) selectionFill(tag int, tagged bool) (color.Color, string) {
	if m.glow && tagged {
		return m.glowFor(tag)
	}
	if m.darkBG {
		return colSelBgDark, "5;" + selBgDark
	}
	return colSelBgLight, "5;" + selBgLight
}
