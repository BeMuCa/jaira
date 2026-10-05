package tui

import (
	"image/color"
	"slices"
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// screen is a rendered frame read back cell by cell, so a test can ask what
// glyph sits where and in which colour — the two things a line is made of.
type screen struct {
	cv   *lipgloss.Canvas
	w, h int
}

func readScreen(m *Model) screen {
	out := m.render()
	h := lipgloss.Height(out)
	return screen{cv: lipgloss.NewCanvas(m.width, h).Compose(lipgloss.NewLayer(out)), w: m.width, h: h}
}

func (s screen) glyph(x, y int) string {
	if c := s.cv.CellAt(x, y); c != nil {
		return c.Content
	}
	return ""
}

func (s screen) fg(x, y int) color.Color {
	if c := s.cv.CellAt(x, y); c != nil {
		return c.Style.Fg
	}
	return nil
}

// find is the cell the first occurrence of text starts at.
func (s screen) find(t *testing.T, text string) (int, int) {
	t.Helper()
	want := strings.Split(text, "")
	for y := 0; y < s.h; y++ {
		for x := 0; x+len(want) <= s.w; x++ {
			hit := true
			for i, g := range want {
				if s.glyph(x+i, y) != g {
					hit = false
					break
				}
			}
			if hit {
				return x, y
			}
		}
	}
	t.Fatalf("%q is not on screen", text)
	return 0, 0
}

// border walks from x along row y in direction d to the nearest lane border,
// which a line running along it may already have drawn over.
func (s screen) border(t *testing.T, x, y, d int) int {
	t.Helper()
	for ; x >= 0 && x < s.w; x += d {
		if g := s.glyph(x, y); g == "│" || slices.Contains(lineGlyphs, g) {
			return x
		}
	}
	t.Fatalf("no lane border from column %d on row %d", x, y)
	return 0
}

func sameColour(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// expect fails unless the cell holds glyph in colour.
func (s screen) expect(t *testing.T, x, y int, glyph string, colour color.Color) {
	t.Helper()
	if g := s.glyph(x, y); g != glyph || !sameColour(s.fg(x, y), colour) {
		t.Errorf("cell (%d,%d) = %q in %v, want %q in %v", x, y, g, s.fg(x, y), glyph, colour)
	}
}

func titled(t *testing.T, m *Model, title string) string {
	t.Helper()
	for _, tk := range m.tickets {
		if tk.Title == title {
			return tk.ID
		}
	}
	t.Fatalf("no ticket titled %q", title)
	return ""
}

// addTicket files one more ticket into the seeded board. updated decides
// where it sorts in its lane: newer is higher.
func addTicket(t *testing.T, m *Model, title, status string, updated time.Time, fields map[string]string, lists map[string][]string) string {
	t.Helper()
	f := map[string]string{
		ticket.FieldID:        ticket.NewID(time.Now()),
		ticket.FieldTitle:     title,
		ticket.FieldStatus:    status,
		ticket.FieldUpdatedAt: ticket.FormatTime(updated),
	}
	for k, v := range fields {
		f[k] = v
	}
	tk, err := m.store.Create(f, lists, "")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond) // ids are time-ordered; keep them distinct
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	return tk.ID
}

var lineGlyphs = []string{"━", "┃", "┏", "┓", "┗", "┛"}

// The seeded "Refactor auth middleware" waits on "Fix session cookie…",
// two lanes to its right and level with it: the line runs straight from the
// cursor card's right border to the blocker's left border, in red, across
// whatever lies between.
func TestLinkLineRunsStraightToABlockerOnTheSameRow(t *testing.T) {
	m := newTestModel(t, 200, 40)
	m.selectByID(titled(t, m, "Refactor auth middleware"))
	s := readScreen(m)

	sx, sy := s.find(t, "Refactor")
	bx, by := s.find(t, "Fix sess")
	if sy != by {
		t.Fatalf("setup: the two cards are on rows %d and %d, want one row", sy, by)
	}
	from, to := s.border(t, sx, sy, 1), s.border(t, bx, by, -1)
	for x := from; x <= to; x++ {
		s.expect(t, x, sy+1, "━", colErr)
	}
	// The card's own title row is untouched: the line takes the middle row.
	if g := s.glyph(from, sy); g != "│" {
		t.Errorf("the title row's border became %q", g)
	}
}

// A target lower in its lane is reached by turning down in the border gap
// just before its column and entering it on its middle row.
func TestLinkLineBendsDownToALowerCard(t *testing.T) {
	m := newTestModel(t, 200, 40)
	addTicket(t, m, "Newer card on top", "in-progress", time.Now().Add(time.Hour), nil, nil)
	m.selectByID(titled(t, m, "Refactor auth middleware"))
	s := readScreen(m)

	sx, sy := s.find(t, "Refactor")
	bx, by := s.find(t, "Fix sess")
	if by != sy+cardSlots {
		t.Fatalf("setup: blocker on row %d, want one card below the cursor's row %d", by, sy)
	}
	from, to := s.border(t, sx, sy, 1), s.border(t, bx, by, -1)
	bend := to - 1
	for x := from; x < bend; x++ {
		s.expect(t, x, sy+1, "━", colErr)
	}
	s.expect(t, bend, sy+1, "┓", colErr)
	for y := sy + 2; y < by+1; y++ {
		s.expect(t, bend, y, "┃", colErr)
	}
	s.expect(t, bend, by+1, "┗", colErr)
	s.expect(t, to, by+1, "━", colErr)
}

// A link to a lane on the left runs leftwards, in the colour of its family.
func TestLinkLineRunsLeftToWhatItFollows(t *testing.T) {
	m := newTestModel(t, 200, 40)
	prior := titled(t, m, "Investigate flaky logout test")
	id := addTicket(t, m, "Follow-up card", "in-progress", time.Now().Add(time.Hour),
		map[string]string{ticket.FieldFollows: prior}, nil)
	m.selectByID(id)
	s := readScreen(m)

	sx, sy := s.find(t, "Follow-up")
	fx, fy := s.find(t, "Investig")
	if sy != fy {
		t.Fatalf("setup: the two cards are on rows %d and %d, want one row", sy, fy)
	}
	from, to := s.border(t, sx, sy, -1), s.border(t, fx, fy, 1)
	for x := to; x <= from; x++ {
		s.expect(t, x, sy+1, "━", colOK)
	}
}

// A link inside one lane runs along that lane's right border.
func TestLinkLineInOneLaneRunsAlongItsBorder(t *testing.T) {
	m := newTestModel(t, 200, 40)
	parent := titled(t, m, "Refactor auth middleware")
	child := addTicket(t, m, "Child of the refactor", "todo", time.Now().Add(-time.Hour),
		map[string]string{ticket.FieldParent: parent}, nil)
	m.selectByID(child)
	s := readScreen(m)

	cx, cy := s.find(t, "Child of")
	_, py := s.find(t, "Refactor")
	if py != cy-cardSlots {
		t.Fatalf("setup: parent on row %d, want the card above the child's row %d", py, cy)
	}
	edge := s.border(t, cx, cy, 1)
	blue := lineFamilies[1].colour
	s.expect(t, edge, cy+1, "┛", blue)
	for y := cy; y > py+1; y-- {
		s.expect(t, edge, y, "┃", blue)
	}
	s.expect(t, edge, py+1, "┓", blue)
}

// A card with no links draws no line and no legend.
func TestUnlinkedCardDrawsNothing(t *testing.T) {
	m := newTestModel(t, 200, 40)
	m.selectByID(titled(t, m, "Rate limit the login endpoint"))
	out := stripANSI(m.render())

	for _, g := range lineGlyphs {
		if strings.Contains(out, g) {
			t.Errorf("an unlinked card drew %q", g)
		}
	}
	for _, word := range []string{"blocker", "off screen"} {
		if strings.Contains(out, word) {
			t.Errorf("an unlinked card put %q in the status bar", word)
		}
	}
}

// A link whose other end has no card on screen draws nothing, and the bar
// counts it — while the legend names only the families that were drawn.
func TestLinkWithoutACardOnScreenIsCounted(t *testing.T) {
	m := newTestModel(t, 200, 40)
	nowhere := ticket.NewID(time.Unix(0, 0))
	id := addTicket(t, m, "Card with a dangling blocker", "backlog", time.Now().Add(time.Hour), nil,
		map[string][]string{
			ticket.FieldBlockedBy: {nowhere},
			ticket.FieldRelated:   {titled(t, m, "Fix session cookie dropped on 302")},
		})
	m.selectByID(id)
	out := stripANSI(m.render())

	if !strings.Contains(out, "1 link off screen · L") {
		t.Errorf("the bar does not count the undrawn link:\n%s", out)
	}
	if !strings.Contains(out, "related") {
		t.Error("the legend does not name the related line that was drawn")
	}
	if strings.Contains(out, "blocker") {
		t.Error("the legend names the blocker family although no blocker line was drawn")
	}
}

// Lines belong to the board alone: an open ticket, the link window over the
// board, and the move picker all show it without them.
func TestLinkLinesOnlyOnTheBoard(t *testing.T) {
	for _, k := range []string{"enter", "L", "m"} {
		m := newTestModel(t, 200, 40)
		m.selectByID(titled(t, m, "Refactor auth middleware"))
		if !strings.Contains(stripANSI(m.render()), "━") {
			t.Fatal("setup: the board itself draws no line")
		}
		m.key(key(k))
		if out := stripANSI(m.render()); strings.Contains(out, "━") {
			t.Errorf("after %q (mode %v) the screen still draws a line", k, m.mode)
		}
	}
}

// Where two lines share cells the stronger family is the one left showing: a
// blocker and a related ticket on the same row both leave along the cursor
// card's middle row, and up to the blocker that run is red.
func TestStrongerLineWinsWhereLinesOverlap(t *testing.T) {
	m := newTestModel(t, 200, 40)
	id := addTicket(t, m, "Card with two links", "todo", time.Now().Add(time.Hour), nil,
		map[string][]string{
			ticket.FieldBlockedBy: {titled(t, m, "Fix session cookie dropped on 302")},
			ticket.FieldRelated:   {titled(t, m, "Decide on cookie SameSite policy")},
		})
	m.selectByID(id)
	s := readScreen(m)

	sx, sy := s.find(t, "Card with")
	bx, by := s.find(t, "Fix sess")
	rx, ry := s.find(t, "Decide o")
	if sy != by || sy != ry {
		t.Fatalf("setup: cards on rows %d, %d, %d, want one row", sy, by, ry)
	}
	from, blocker := s.border(t, sx, sy, 1), s.border(t, bx, by, -1)
	for x := from; x <= blocker; x++ {
		s.expect(t, x, sy+1, "━", colErr)
	}
	for x := s.border(t, bx, by, 1); x <= s.border(t, rx, ry, -1); x++ {
		s.expect(t, x, sy+1, "━", colDim)
	}
}

// A blocker named by its handle — what 'jaira create --blocked-by SEFFWC'
// writes — is the same link as one named by its full id.
func TestLinkLineFindsABlockerNamedByItsHandle(t *testing.T) {
	m := newTestModel(t, 200, 40)
	blocker := titled(t, m, "Fix session cookie dropped on 302")
	id := addTicket(t, m, "Card naming a handle", "todo", time.Now().Add(time.Hour), nil,
		map[string][]string{ticket.FieldBlockedBy: {ticket.Handle(blocker)}})
	m.selectByID(id)
	s := readScreen(m)

	sx, sy := s.find(t, "Card nam")
	bx, by := s.find(t, "Fix sess")
	if sy != by {
		t.Fatalf("setup: the two cards are on rows %d and %d, want one row", sy, by)
	}
	for x := s.border(t, sx, sy, 1); x <= s.border(t, bx, by, -1); x++ {
		s.expect(t, x, sy+1, "━", colErr)
	}
	if out := stripANSI(m.render()); strings.Contains(out, "off screen") {
		t.Error("the handle was counted as a link with no card")
	}
}
