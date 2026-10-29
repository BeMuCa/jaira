package tui

import (
	"image/color"
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

// selectionFill is what the selected card is painted with: one step off the
// terminal's background, whichever way that is.
func (m *Model) selectionFill() (color.Color, string) {
	if m.darkBG {
		return colSelBgDark, "5;" + selBgDark
	}
	return colSelBgLight, "5;" + selBgLight
}
