package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// The selected card is the one thing on the board that has to be findable
// without reading anything, so its fill is asserted on the rendered frame
// rather than on the style that produced it.
func TestSelectedCardIsFilled(t *testing.T) {
	m := newTestModel(t, 120, 40)
	out := m.View().Content

	dark := fillMark(selBgDark)
	if !strings.Contains(out, dark) {
		t.Fatalf("selected card carries no fill; want %q somewhere in the frame", dark)
	}
	if n := strings.Count(out, dark); n < 3 {
		t.Errorf("fill reaches %d spans, so it stops partway down the card; want the whole box", n)
	}
}

// Only one card is selected, so only one card may be filled — a fill that
// leaked onto its neighbours would say nothing about where the cursor is.
func TestUnselectedCardsAreNotFilled(t *testing.T) {
	m := newTestModel(t, 120, 40)

	dark := fillMark(selBgDark)
	filled := 0
	for _, line := range strings.Split(m.View().Content, "\n") {
		if strings.Contains(line, dark) {
			filled++
		}
	}
	// One card is five rows tall, borders included.
	if filled > 5 {
		t.Errorf("fill covers %d rows, more than the one selected card", filled)
	}
}

// Stacked cards share one border row, and the row has to belong to the
// selected card: given to the card above, the selection is a box open at the
// top whose fill starts mid-card.
func TestSelectedCardKeepsItsTopEdge(t *testing.T) {
	m := newTestModel(t, 120, 40)

	// Onto the second card of the lane, the first position where the top
	// border is the shared one.
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.cardIdx != 1 {
		t.Skipf("lane holds no second card to select (cardIdx=%d)", m.cardIdx)
	}

	dark := fillMark(selBgDark)
	for _, line := range strings.Split(m.View().Content, "\n") {
		if !strings.Contains(line, dark) {
			continue
		}
		// The first filled row is the selected card's top border.
		if !strings.Contains(line, "┌") {
			t.Errorf("first filled row carries no top border, so the selection is open at the top: %q", line)
		}
		return
	}
	t.Error("selected card carries no fill at all")
}

// A light terminal has nothing lighter than its background to offer, so the
// fill goes the other way. Without this the fill on a light theme is a black
// block.
func TestFillFollowsTerminalBackground(t *testing.T) {
	m := newTestModel(t, 120, 40)

	if _, code := m.selBg(); code != selBgDark {
		t.Errorf("default fill is %q, want the dark-terminal one %q", code, selBgDark)
	}

	m.Update(tea.BackgroundColorMsg{Color: white{}})
	if _, code := m.selBg(); code != selBgLight {
		t.Errorf("fill on a light terminal is %q, want %q", code, selBgLight)
	}
	if strings.Contains(m.View().Content, fillMark(selBgDark)) {
		t.Error("board still paints the dark fill after the terminal reported a light background")
	}
}

// refill must not re-open the background after a line's last reset: that would
// paint the terminal row past the card's right edge.
func TestRefillLeavesTrailingResetAlone(t *testing.T) {
	in := "\x1b[38;5;244mmeta\x1b[m plain\x1b[m"
	got := refill(in, selBgDark)
	if !strings.HasSuffix(got, "\x1b[m") {
		t.Fatalf("refill(%q) = %q, want it to end on a bare reset", in, got)
	}
	if strings.Contains(got, "\x1b[m\x1b[48;5;"+selBgDark+"m\x1b[m") {
		t.Errorf("refill reopened the fill at the end of the line: %q", got)
	}
	if !strings.Contains(got, "meta\x1b[m\x1b[48;5;"+selBgDark+"m plain") {
		t.Errorf("refill did not reopen the fill mid-line: %q", got)
	}
}

// fillMark is the SGR parameter the fill sets, not a whole escape sequence:
// lipgloss merges the background into the same sequence as a foreground where a
// row has both, so on a card's border rows the fill never stands alone.
func fillMark(code string) string { return "48;5;" + code + "m" }

type white struct{}

func (white) RGBA() (r, g, b, a uint32) { return 0xffff, 0xffff, 0xffff, 0xffff }
