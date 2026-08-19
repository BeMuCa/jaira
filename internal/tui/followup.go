package tui

import (
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/ticket"
)

// splitMinWidth is the narrowest terminal that gets two panes. Below it the
// follow-up takes the screen: the predecessor is context, and context is the
// first thing a narrow terminal has to give up.
const splitMinWidth = 80

// splitMinHeight is the shortest terminal that gets two panes. The field editor
// needs its own rows; below this there is no pane left to put beside it.
const splitMinHeight = 20

// draft is a follow-up that has not been written yet. The field editor writes
// each field straight to disk (edit.go), which is right for a ticket that exists
// and wrong for one the user may still think better of — so a draft is edited in
// memory and becomes a file only when it is saved.
type draft struct {
	fields map[string]string
	lists  map[string][]string
	body   string
}

// follow is the split view: the ticket a follow-up is being written for on the
// left, the follow-up itself on the right. It exists so the reason for the
// follow-up is still on screen while it is written — which is the whole reason
// this board exists, applied to the moment a ticket is born.
type follow struct {
	src       *ticket.Ticket
	srcScroll int
	// draft is non-nil while the follow-up is being written. Once it is saved the
	// right pane is a real ticket in m.detail and reads like any other.
	draft *draft
	// focusLeft decides which pane the plain scroll keys move. The action keys
	// always belong to the ticket being worked, which is the right one.
	focusLeft bool
}

// startFollowUp opens a follow-up of the open ticket beside it. Nothing is
// written: the ticket appears on the board when it is saved, so a follow-up
// thought better of leaves no trace. Called again from a saved follow-up it
// chains — that ticket becomes the one on the left.
func (m *Model) startFollowUp() {
	src := m.detail
	if src == nil {
		return
	}
	fields, lists, body := m.followUpFields(src, "Follows on from "+ticket.Handle(src.ID)+".")
	m.follow = &follow{
		src: src,
		// The predecessor keeps the position it was being read at.
		srcScroll: m.detailScroll,
		draft:     &draft{fields: fields, lists: lists, body: body},
	}
	m.detail = nil
	m.detailScroll = 0
	m.mode = modeEdit
	m.editIdx = 0
	m.editBuf = m.editValue(0)
}

// saveFollowUp writes the draft. Readiness is recomputed straight after, the way
// the on-disk editor recomputes it after every field, so a saved follow-up is
// not left claiming a readiness nobody derived.
func (m *Model) saveFollowUp() {
	f := m.follow
	if f == nil || f.draft == nil {
		return
	}
	d := f.draft
	t, err := m.store.Create(d.fields, d.lists, d.body)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	if _, err := m.store.Mutate(t.ID, func(x *ticket.Ticket) error {
		return ticket.SetReady(x.Doc(), gate.Ready(x))
	}); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	full, err := m.store.Load(t.ID)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	f.draft = nil
	m.detail = full
	m.detailScroll = 0
	m.editBuf = ""
	m.mode = modeDetail
}

// closeFollowUp leaves the split, whether the follow-up was written or not, and
// puts the ticket it was for back on screen alone — that is the ticket the user
// came from, so it is the one esc owes them.
func (m *Model) closeFollowUp() {
	f := m.follow
	if f == nil {
		return
	}
	m.follow = nil
	m.detail = f.src
	m.detailScroll = f.srcScroll
	m.editBuf = ""
	m.mode = modeDetail
}

// scrollFocused moves whichever pane tab has focused, or the open ticket when
// there is no split.
func (m *Model) scrollFocused(d int) {
	if m.follow != nil && m.follow.draft == nil && m.follow.focusLeft {
		m.follow.srcScroll += d
		return
	}
	m.detailScroll += d
}

// leaveDetail closes the open ticket, and the split with it — walking to the
// neighbouring ticket leaves the whole follow-up screen behind.
func (m *Model) leaveDetail() {
	m.follow = nil
	m.mode = m.detailFrom
	m.detail = nil
}

// scrollSrc moves the left pane, which stays readable in both states — while the
// follow-up is being written the editor owns tab, so the predecessor needs a
// scroll of its own.
func (m *Model) scrollSrc(d int) {
	if m.follow != nil {
		m.follow.srcScroll += d
	}
}

// editValue reads the field the editor is on from whichever subject it is
// editing: the draft when there is one, the open ticket otherwise.
func (m *Model) editValue(i int) string {
	if d := m.editDraft(); d != nil {
		return d.fields[editableFields[i].field]
	}
	if m.detail == nil {
		return ""
	}
	return editableFields[i].get(m.detail)
}

func (m *Model) editDraft() *draft {
	if m.follow == nil {
		return nil
	}
	return m.follow.draft
}

// followUpFields builds the frontmatter, lists and body a follow-up of src
// carries. The sign-off path and the split writer share it so the two cannot
// drift apart. lead is the first sentence of the context: only the sign-off path
// may claim a review happened.
func (m *Model) followUpFields(src *ticket.Ticket, lead string) (map[string]string, map[string][]string, string) {
	now := time.Now()
	me := identity.Current(m.store.Root)
	def := m.lanes.Default()

	fields := map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     "Follow-up: " + src.Title,
		ticket.FieldStatus:    def.ID,
		ticket.FieldReady:     "false",
		ticket.FieldCreator:   me,
		ticket.FieldAssignee:  firstNonEmpty(src.Assignee, me),
		ticket.FieldContext:   followUpContext(src, lead),
		ticket.FieldFollows:   src.ID,
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	lists := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}
	body := "# Follow-up: " + src.Title + "\n\n" +
		"## Definition of Done\n\n- [ ] <What must be true that is not true yet>\n\n## Progress\n\n"
	return fields, lists, body
}

// renderSplit draws the follow-up beside the ticket it follows.
func (m *Model) renderSplit() string {
	f := m.follow
	if f == nil {
		return m.renderBoard()
	}
	if m.width < splitMinWidth || m.height < splitMinHeight {
		// One readable column beats two unreadable ones, and a pane too short to
		// hold the editor is not a pane.
		if f.draft != nil {
			return m.renderEdit()
		}
		return m.renderDetail()
	}

	items := strings.Split(m.splitHints(), " · ")
	footer := wrapHints(items, max(1, m.width))
	// A pane costs a row top and bottom and a column each side for its box, and
	// the boxes are what say which pane is live.
	rows := max(3, m.height-len(footer)-3)
	inner := max(10, (m.width-1)/2-2)

	left, srcScroll := clipPane(m.detailBody(f.src, inner), inner, rows, f.srcScroll)
	f.srcScroll = srcScroll

	var right string
	if f.draft != nil {
		right, _ = clipPane(m.editBody(inner, rows), inner, rows, 0)
	} else {
		right, m.detailScroll = clipPane(m.detailBody(m.detail, inner), inner, rows, m.detailScroll)
	}

	// While the follow-up is being written the editor holds the keyboard, so the
	// right pane is live whatever focusLeft says.
	leftLive := f.draft == nil && f.focusLeft
	var b strings.Builder
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		splitPane(left, leftLive), " ", splitPane(right, !leftLive)))
	for _, l := range footer {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return clampBlock(b.String(), m.width, m.height)
}

// splitPane boxes one pane. The live pane's border carries the accent colour:
// which half the keys are talking to must never be a guess.
//
// Neither Width nor Height is set: clipPane already handed over an exact
// rectangle, and a sized lipgloss style measures its frame including the border
// — so setting Width here re-wraps every full-width row onto a second line and
// grows the pane past the terminal. The border hugs the rectangle instead.
func splitPane(content string, live bool) string {
	border := colFaint
	if live {
		border = colAccent
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(border).
		Render(content)
}

func (m *Model) splitHints() string {
	if m.editDraft() != nil {
		return "ctrl+s save · tab next field · shift+up/down scroll left · esc discard"
	}
	return m.detailHints(m.detail) + " · tab other pane · esc close"
}
