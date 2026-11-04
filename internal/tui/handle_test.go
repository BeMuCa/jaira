package tui

import (
	"image/color"
	"testing"
)

// The ticket handle must read on the selected card and on both card shades in
// either terminal background. Every fill the board paints is grey, so the
// level of one channel is enough to compare them.
func TestHandleToneReadsOnEveryCardFill(t *testing.T) {
	for _, dark := range []bool{true, false} {
		name := "light"
		if dark {
			name = "dark"
		}
		m := newTestModel(t, 150, 32)
		m.darkBG = dark

		sel, _ := m.selectionFill()
		shadeA, _ := m.laneShade(false)
		shadeB, _ := m.laneShade(true)
		fills := map[string]color.Color{"selection": sel, "shade A": shadeA, "shade B": shadeB}

		handle := grey(m.handleStyle().GetForeground())
		for what, f := range fills {
			d := handle - grey(f)
			if d < 0 {
				d = -d
			}
			if d < 60 {
				t.Errorf("%s: handle level %d is only %d from %s (level %d), want at least 60",
					name, handle, d, what, grey(f))
			}
		}
	}
}

// grey is the 8-bit level of a colour the board paints grey.
func grey(c color.Color) int {
	r, _, _, _ := c.RGBA()
	return int(r >> 8)
}
