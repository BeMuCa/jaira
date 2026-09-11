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
	// A short window can be smaller than the first group, so the scroll has
	// to be settled before the first render rather than on the first
	// keypress — otherwise the window opens with its selection off screen.
	v.scrollInto(m.linkLines())
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
		v.scrollInto(m.linkLines())
	case "k", "up":
		if n := v.next(v.cursor, -1); n >= 0 {
			v.cursor = n
		}
		v.scrollInto(m.linkLines())
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

// linkLines is how many lines the list may use, indicator lines included.
//
// It has to agree with what modalOver leaves for content — the box keeps a
// margin off the window and a line of border each side — minus this window's
// fixed furniture: the title, the rule under it and the blank line after it,
// then the blank line and the hint line at the bottom. When this number is
// too large the modal clips the list itself, and a clip cannot scroll; what
// it eats first is the hint line that says which keys do anything.
//
// The floor is two lines, one entry, because a window that shows nothing is
// worse than one that shows a little and says so.
func (m *Model) linkLines() int {
	const fixed = 3 + 2
	return max(2, m.height-8-fixed)
}

// scrollInto moves the visible run so the cursor is inside it, and no
// further — a window that recentred on every keypress makes the whole list
// move under the reader.
func (v *linkView) scrollInto(lines int) {
	// A heading belongs to the entry under it, so scrolling to an entry
	// shows the heading that names it rather than starting mid-group.
	want := v.cursor
	if want > 0 && want < len(v.rows) && v.rows[want-1].heading != "" {
		want--
	}
	if want < v.top {
		v.top = want
	}
	// Every row prints two lines: an entry is the ticket and where it lives,
	// a heading is a blank line and the heading itself.
	if span := (v.cursor - v.top + 1) * 2; span > lines {
		v.top = max(0, v.cursor-lines/4)
	}
	v.top = max(0, min(v.top, max(0, len(v.rows)-1)))
}

// visible is the run of rows on screen, and how much is hidden above and
// below it.
//
// Three things want the same lines and they are not equally important: the
// selected row has to be on screen or the window is unusable, the list is
// what the reader came for, and the two lines that say how much is hidden
// are a courtesy. So the markers are offered the space first and give it up
// again as soon as keeping them would push the selection off the screen —
// and when nothing fits, the run starts at the selection itself.
func (v *linkView) visible(lines int) (out []linkRow, first, above, below int) {
	pack := func(top, budget int) ([]linkRow, int) {
		var rows []linkRow
		used := 0
		for i := top; i < len(v.rows); i++ {
			// Every row prints two lines: an entry is the ticket and where
			// it lives, a heading is a blank line and the heading itself.
			if used+2 > budget {
				return rows, len(v.rows) - i
			}
			used += 2
			rows = append(rows, v.rows[i])
		}
		return rows, 0
	}
	for _, reserve := range []int{2, 1, 0} {
		rows, left := pack(v.top, lines-reserve)
		if v.top+len(rows) <= v.cursor {
			continue // the selection would be off screen
		}
		above, below = 0, 0
		if v.top > 0 {
			above = v.top
		}
		if left > 0 {
			below = left
		}
		need := 0
		if above > 0 {
			need++
		}
		if below > 0 {
			need++
		}
		if need > reserve {
			continue // not enough room to say what is hidden
		}
		return rows, v.top, above, below
	}
	// Nothing fits around the selection: start at it and show what there is.
	rows, _ := pack(v.cursor, lines)
	return rows, v.cursor, 0, 0
}

// renderLinks draws the link window.
func (m *Model) renderLinks() string {
	v := m.links
	if v == nil {
		return m.renderBoard()
	}
	var b strings.Builder
	title := "Links"
	if v.title != "" {
		title = "Links · " + v.title
	}
	b.WriteString(styLaneTitle.Render(title) + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 72))) + "\n\n")
	if len(v.rows) == 0 {
		b.WriteString(styMeta.Render("Nothing is linked to this ticket yet.") + "\n")
		b.WriteString("\n" + styMeta.Render("esc  back") + "\n")
		return b.String()
	}
	rows, first, above, below := v.visible(m.linkLines())
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
	b.WriteString("\n" + styMeta.Render("j/k  move    enter  jump    esc  back") + "\n")
	return b.String()
}
