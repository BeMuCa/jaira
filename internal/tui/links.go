package tui

import (
	"fmt"
	"strings"

	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
)

// linkView is the open link window: the ticket it was opened on and the
// connections found for it, already flattened into the order they render in.
//
// The entries are captured when the window opens rather than recomputed per
// keypress. Answering "who names me as their parent" means reading the
// logbook, and doing that on every cursor move would make the window feel
// like the board does not.
type linkView struct {
	subject string
	// title is the subject's own title, kept because the window is opened
	// from the ticket itself — looking it up again through the board would
	// also fail for a ticket the board has no card for.
	title  string
	rows   []linkRow
	cursor int
	// top is the first row on screen. An epic with many children is longer
	// than the box, and a window that clipped without scrolling could hold a
	// selection the reader cannot see — which reads as the cursor being
	// stuck.
	top int
	// from is the view this window was opened over, kept here rather than in
	// the shared returnTo field: a refusal raised inside this window uses
	// that field for itself, and the two would overwrite each other.
	from mode
}

// linkRow is one printed line: either a group heading or one linked ticket.
type linkRow struct {
	heading string
	entry   link.Entry
}

// selectable reports whether the cursor may rest on this row. A heading is
// not a destination.
func (r linkRow) selectable() bool { return r.heading == "" }

// openLinks builds the link window for the selected card.
func (m *Model) openLinks() {
	// The open ticket when there is one: in the detail pane the card the
	// board cursor sits on is not necessarily what the reader is looking at.
	t := m.detail
	if m.mode != modeDetail || t == nil {
		t = m.selected()
	}
	if t == nil {
		return
	}
	entries := m.linkIndex().Relations(t.ID)
	v := &linkView{subject: t.ID, title: t.Title}
	for _, k := range link.Order {
		var group []link.Entry
		for _, e := range entries {
			if e.Kind == k {
				group = append(group, e)
			}
		}
		if len(group) == 0 {
			continue
		}
		v.rows = append(v.rows, linkRow{heading: k.Title()})
		for _, e := range group {
			v.rows = append(v.rows, linkRow{entry: e})
		}
	}
	v.cursor = v.next(-1, 1)
	v.from = m.mode
	m.links = v
	m.mode = modeLinks
}

// next finds the nearest selectable row from i in direction d, staying put
// when there is none — so the cursor can never land on a heading and can
// never run off either end.
func (v *linkView) next(i, d int) int {
	for j := i + d; j >= 0 && j < len(v.rows); j += d {
		if v.rows[j].selectable() {
			return j
		}
	}
	if i >= 0 && i < len(v.rows) && v.rows[i].selectable() {
		return i
	}
	return -1
}

// current is the entry under the cursor, if the cursor is on one.
func (v *linkView) current() (link.Entry, bool) {
	if v == nil || v.cursor < 0 || v.cursor >= len(v.rows) {
		return link.Entry{}, false
	}
	row := v.rows[v.cursor]
	if !row.selectable() {
		return link.Entry{}, false
	}
	return row.entry, true
}

// keyLinks drives the link window.
func (m *Model) keyLinks(s string) {
	v := m.links
	if v == nil {
		m.mode = modeBoard
		return
	}
	switch s {
	case "esc", "q", "L":
		m.mode = v.from
		m.links = nil
	case "j", "down":
		if n := v.next(v.cursor, 1); n >= 0 {
			v.cursor = n
		}
	case "k", "up":
		if n := v.next(v.cursor, -1); n >= 0 {
			v.cursor = n
		}
	case "enter":
		e, ok := v.current()
		if !ok {
			return
		}
		// Only a ticket that is on this board can be jumped to: the board
		// has no card for one in the logbook or on a ref, and moving the
		// cursor to nothing would read as the jump having failed silently.
		if e.Ref.Place != link.PlaceBoard {
			m.notify(fmt.Sprintf("%s is %s.\n\nThe board has no card for it, so there is nowhere to jump to.",
				ticket.Handle(e.Ref.ID), e.Ref.Whereabouts()), false)
			return
		}
		from := v.from
		m.links = nil
		m.mode = modeBoard
		if !m.selectByID(e.Ref.ID) {
			// On the board but not on screen: a filter is hiding it. Saying
			// so beats leaving the cursor where it was, which reads as the
			// jump having silently failed.
			m.notify(fmt.Sprintf("%s is on the board but the filter %q hides it.\n\nClear the filter with esc, then try again.",
				ticket.Handle(e.Ref.ID), m.filter), false)
			return
		}
		// Following a link out of an open ticket opens the ticket it leads
		// to: the reader was reading, not navigating the board.
		if from == modeDetail {
			m.openDetail()
		}
	}
}

// The window's layout is one calculation, done in one place, because three
// rounds of review found the same defect three times: the lines were being
// budgeted in one function and spent in another, and the two disagreed.
//
// Everything here is stated in rendered lines. A row costs two — an entry is
// the ticket and where it lives, a heading is a blank line and the heading
// itself — and the only caller that may move the visible run is frame(),
// which writes back the top it decided on, so what the reader sees and what
// the next keypress reasons about are always the same thing.

// content is how many lines modalOver leaves for this window to fill. Going
// over it means the modal clips, and a clip cannot scroll — what it eats
// first is the hint line that says which keys do anything.
func (m *Model) content() int { return m.modalContent() }

// frame decides what is on screen: it moves top so the selection is inside
// the run, then reports the run and how much is hidden either side.
//
// The claims on the space are ranked, because they are not equally
// important. The hint line is always drawn — a window whose keys are a
// secret is unusable. The selected row is next: losing it is what makes the
// list feel stuck. The heading above the selection, the rest of the list,
// the title, and finally the two lines that say how much is hidden, take
// what is left in that order.
func (v *linkView) frame(lines int) (rows []linkRow, first, above, below int) {
	if len(v.rows) == 0 {
		return nil, 0, 0, 0
	}
	// Reserve nothing for the markers yet: whether they are drawn depends on
	// whether anything is hidden, which is what this is working out.
	fit := func(reserve int) (int, int) {
		capacity := max(1, (lines-reserve)/2)
		top := v.top
		if v.cursor < top {
			top = v.cursor
		}
		if v.cursor >= top+capacity {
			top = v.cursor - capacity + 1
		}
		// A heading names the entry under it, so keep the two together when
		// the extra row is affordable.
		if top > 0 && top == v.cursor && capacity >= 2 && v.rows[top-1].heading != "" {
			top--
		}
		return top, min(len(v.rows), top+capacity)
	}
	top, end := fit(0)
	hidden := 0
	if top > 0 {
		hidden++
	}
	if end < len(v.rows) {
		hidden++
	}
	if hidden > 0 {
		// Now that it is known how many marker lines would be drawn, redo the
		// fit paying for them. If that costs the selection its place, the
		// markers are the ones that go: the reader can live without being
		// told how much is left over.
		if t2, e2 := fit(hidden); v.cursor >= t2 && v.cursor < e2 {
			top, end = t2, e2
		}
	}
	v.top = top
	// Whatever the fit worked out, the markers are drawn only out of lines
	// that are genuinely spare. The row capacity has a floor of one, so a
	// very short window can have no spare line at all — and a marker drawn
	// there is a line over budget, which the modal takes off the bottom.
	spare := lines - (end-top)*2
	above, below = 0, 0
	if top > 0 && spare > 0 {
		above, spare = top, spare-1
	}
	if end < len(v.rows) && spare > 0 {
		below = len(v.rows) - end
	}
	return v.rows[top:end], top, above, below
}

// renderLinks draws the link window.
func (m *Model) renderLinks() string {
	v := m.links
	if v == nil {
		return m.renderBoard()
	}
	avail := m.content()
	hint := styMeta.Render("j/k  move    enter  jump    esc  back")

	var b strings.Builder
	// The title and the rule under it are the first thing to go when the
	// window is short: they say what the reader already knows, having just
	// pressed the key that opened it.
	if avail >= 9 {
		title := "Links"
		if v.title != "" {
			title = "Links · " + v.title
		}
		b.WriteString(styLaneTitle.Render(title) + "\n")
		b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 72))) + "\n")
		avail -= 2
	}
	if len(v.rows) == 0 {
		b.WriteString(styMeta.Render("Nothing is linked to this ticket yet.") + "\n")
		b.WriteString(styMeta.Render("esc  back"))
		return b.String()
	}

	rows, first, above, below := v.frame(avail - 1) // the hint line
	if above > 0 {
		b.WriteString(styMeta.Render(fmt.Sprintf("↑ %d more above", above)) + "\n")
	}
	for j, row := range rows {
		i := first + j
		if row.heading != "" {
			b.WriteString("\n" + styMeta.Render(row.heading) + "\n")
			continue
		}
		e := row.entry
		marker := "  "
		if i == v.cursor {
			marker = "▸ "
		}
		indent := strings.Repeat("  ", e.Ref.Depth)
		b.WriteString(marker + indent + e.Ref.Label() + "\n")
		b.WriteString(styMeta.Render(fmt.Sprintf("  %s  %s", indent, e.Ref.Whereabouts())) + "\n")
	}
	if below > 0 {
		b.WriteString(styMeta.Render(fmt.Sprintf("↓ %d more below", below)) + "\n")
	}
	b.WriteString(hint)
	return b.String()
}
