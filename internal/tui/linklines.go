package tui

import (
	"fmt"
	"image/color"
	"slices"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
)

// lineFamilies are the colours a link line is drawn in, one per pair of
// kinds rather than one per kind: which way a pair points is the link
// window's to say, and seven colours on one board are more than a reader
// can tell apart. Strongest first — the order the status bar names them in,
// and the reverse of the order they are drawn in, so where two lines share a
// cell the stronger one is the one left showing.
var lineFamilies = []struct {
	name   string
	colour color.Color
}{
	{"blocker", colErr},
	{"parent", lipgloss.Color("33")},
	{"follows", colOK},
	{"related", colDim},
}

// lineFamily is the index into lineFamilies a kind is drawn with.
func lineFamily(k link.Kind) int {
	switch k {
	case link.KindBlockedBy, link.KindBlocks:
		return 0
	case link.KindParent, link.KindChild:
		return 1
	case link.KindFollows, link.KindFollowedBy:
		return 2
	}
	return 3
}

// lineLink is one ticket the cursor card is connected to, with the strongest
// relation between the two: a pair joined twice gets one line, not two drawn
// over each other.
type lineLink struct {
	id     string
	family int
}

// cursorLinks lists the direct links of t, read only from what the board
// already holds in memory. link.Index.Relations would also find the links
// that filed tickets make, but it reads the logbook to do so, and this runs
// on every frame the cursor rests on a card.
//
// The order is the drawing order: weakest family first, then by id so the
// picture does not shuffle between frames.
func (m *Model) cursorLinks(t *ticket.Ticket) []lineLink {
	best := map[string]int{}
	note := func(id string, k link.Kind) {
		if f, ok := best[id]; !ok || lineFamily(k) < f {
			best[id] = lineFamily(k)
		}
	}
	// What t names, kept until a ticket in memory answers to it: a name
	// nothing here answers to is still a link, just one with no card.
	named := map[string]link.Kind{}
	for _, r := range t.BlockedBy {
		named[r] = link.KindBlockedBy
	}
	for _, r := range t.Related {
		named[r] = link.KindRelated
	}
	named[t.Parent] = link.KindParent
	named[t.Follows] = link.KindFollows
	delete(named, "")

	scan := func(o *ticket.Ticket) {
		if o.ID == t.ID {
			return
		}
		for r, k := range named {
			if refersTo(r, o.ID) {
				note(o.ID, k)
				delete(named, r)
			}
		}
		if slices.ContainsFunc(o.BlockedBy, func(r string) bool { return refersTo(r, t.ID) }) {
			note(o.ID, link.KindBlocks)
		}
		if refersTo(o.Parent, t.ID) {
			note(o.ID, link.KindChild)
		}
		if slices.ContainsFunc(o.Related, func(r string) bool { return refersTo(r, t.ID) }) {
			note(o.ID, link.KindRelated)
		}
		if refersTo(o.Follows, t.ID) {
			note(o.ID, link.KindFollowedBy)
		}
	}
	for _, o := range m.tickets {
		scan(o)
	}
	for _, l := range m.logged {
		scan(l.Ticket)
	}
	for r, k := range named {
		if !refersTo(r, t.ID) {
			note(r, k)
		}
	}

	out := make([]lineLink, 0, len(best))
	for id, f := range best {
		out = append(out, lineLink{id: id, family: f})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].family != out[j].family {
			return out[i].family > out[j].family
		}
		return out[i].id < out[j].id
	})
	return out
}

// refersTo reports whether a link field's value names the ticket id. Not
// every field holds a full id: 'jaira create --blocked-by SEFFWC' stores the
// handle as typed, so a value shorter than an id is matched as its tail, the
// way core/link resolves one.
func refersTo(ref, id string) bool {
	if ref == "" {
		return false
	}
	if ref == id {
		return true
	}
	ref = ticket.NormalizeIDPrefix(ref)
	return len(ref) < len(id) && strings.HasSuffix(id, ref)
}

// linkLegend is what the status bar says while the cursor card has links:
// the colour of each family a line is drawn in, and how many links have no
// card on screen to draw to.
func linkLegend(drawn []lineLink, offscreen int) string {
	seen := make([]bool, len(lineFamilies))
	for _, l := range drawn {
		seen[l.family] = true
	}
	var parts []string
	for i, f := range lineFamilies {
		if seen[i] {
			parts = append(parts, lipgloss.NewStyle().Foreground(f.colour).Render("━")+styMeta.Render(" "+f.name))
		}
	}
	s := strings.Join(parts, styMeta.Render("  "))
	if offscreen > 0 {
		noun := "links"
		if offscreen == 1 {
			noun = "link"
		}
		if s != "" {
			s += styMeta.Render(" · ")
		}
		s += styMeta.Render(fmt.Sprintf("%d %s off screen · L", offscreen, noun))
	}
	if s == "" {
		return ""
	}
	return s + "   "
}

// cardSpot is where one card landed on screen: the two border cells of its
// column, and the row the card's middle line is on.
type cardSpot struct {
	left, right, row int
}

// lineTarget is one line to draw: the card it ends at and its colour.
type lineTarget struct {
	at     cardSpot
	colour color.Color
}

// drawLinkLines draws a line from one card to each target over the finished
// board. A line runs along the middle row of the card it starts at, straight
// across whatever cards lie between, turns in the border gap just before the
// target's column, and enters the target on its middle row. A target in the
// same lane is reached along that lane's right border.
//
// Only a cell's glyph and foreground change: its background stays, so a line
// crossing a card keeps the card's band around it. A cell holding part of a
// wide character is left alone — overwriting half of one shifts the rest of
// the row.
func drawLinkLines(out string, width int, from cardSpot, targets []lineTarget) string {
	height := lipgloss.Height(out)
	cv := lipgloss.NewCanvas(width, height).Compose(lipgloss.NewLayer(fitCanvas(out, width, height)))
	for _, tg := range targets {
		put := func(x, y int, glyph string) {
			cell := cv.CellAt(x, y)
			if cell == nil || cell.Width != 1 {
				return
			}
			c := *cell
			c.Content, c.Style.Fg = glyph, tg.colour
			cv.SetCell(x, y, &c)
		}
		to := tg.at
		switch {
		case to.left > from.right:
			bend := to.left - 1
			for x := from.right; x < bend; x++ {
				put(x, from.row, "━")
			}
			bendLine(put, bend, from.row, to.row, "┛", "┓", "┏", "┗")
			put(to.left, to.row, "━")
		case to.right < from.left:
			bend := to.right + 1
			for x := from.left; x > bend; x-- {
				put(x, from.row, "━")
			}
			bendLine(put, bend, from.row, to.row, "┗", "┏", "┓", "┛")
			put(to.right, to.row, "━")
		default:
			bendLine(put, from.right, from.row, to.row, "┛", "┓", "┓", "┛")
		}
	}
	return cv.Render()
}

// bendLine draws the vertical part of a line in column x, from row y0 to row
// y1, with the corner glyphs for where it turns off the horizontal run and
// where it turns into the target. Equal rows need no turn at all.
func bendLine(put func(x, y int, glyph string), x, y0, y1 int, startUp, startDown, endFromBelow, endFromAbove string) {
	switch {
	case y1 == y0:
		put(x, y0, "━")
	case y1 < y0:
		put(x, y0, startUp)
		for y := y0 - 1; y > y1; y-- {
			put(x, y, "┃")
		}
		put(x, y1, endFromBelow)
	default:
		put(x, y0, startDown)
		for y := y0 + 1; y < y1; y++ {
			put(x, y, "┃")
		}
		put(x, y1, endFromAbove)
	}
}
