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
	rows    []linkRow
	cursor  int
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

// linkKindOrder fixes the order the groups appear in, strongest obligation
// first, so the window does not reshuffle between two openings of the same
// ticket.
var linkKindOrder = []link.Kind{
	link.KindBlockedBy, link.KindBlocks,
	link.KindParent, link.KindChild,
	link.KindRelated,
	link.KindFollows, link.KindFollowedBy,
}

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
	ix := link.Build(m.store, m.lanes, m.tickets)
	entries := ix.Relations(t.ID)
	v := &linkView{subject: t.ID}
	for _, k := range linkKindOrder {
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
		m.selectByID(e.Ref.ID)
		// Following a link out of an open ticket opens the ticket it leads
		// to: the reader was reading, not navigating the board.
		if from == modeDetail {
			m.openDetail()
		}
	}
}

// renderLinks draws the link window.
func (m *Model) renderLinks() string {
	v := m.links
	if v == nil {
		return m.renderBoard()
	}
	var b strings.Builder
	title := "Links"
	if t := m.byID(v.subject); t != nil {
		title = "Links · " + t.Title
	}
	b.WriteString(styLaneTitle.Render(title) + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 72))) + "\n\n")
	if len(v.rows) == 0 {
		b.WriteString(styMeta.Render("Nothing is linked to this ticket yet.") + "\n")
		b.WriteString("\n" + styMeta.Render("esc  back") + "\n")
		return b.String()
	}
	for i, row := range v.rows {
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
	b.WriteString("\n" + styMeta.Render("j/k  move    enter  jump    esc  back") + "\n")
	return b.String()
}

// byID finds a listed ticket by id.
func (m *Model) byID(id string) *ticket.Ticket {
	for _, t := range m.tickets {
		if t.ID == id {
			return t
		}
	}
	return nil
}
