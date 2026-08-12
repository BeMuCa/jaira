package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/ticket"
)

// editableFields are the scalar fields worth editing in place. The body, and
// with it both checklists, is left to $EDITOR: a one-line buffer is the wrong
// tool for prose, and the fields below are the ones that block a ticket from
// leaving the backlog.
var editableFields = []struct {
	name  string
	field string
	get   func(*ticket.Ticket) string
}{
	{"goal", ticket.FieldGoal, func(t *ticket.Ticket) string { return t.Goal }},
	{"context", ticket.FieldContext, func(t *ticket.Ticket) string { return t.Context }},
	{"assignee", ticket.FieldAssignee, func(t *ticket.Ticket) string { return t.Assignee }},
	{"title", ticket.FieldTitle, func(t *ticket.Ticket) string { return t.Title }},
	{"tier", ticket.FieldModelTier, func(t *ticket.Ticket) string { return t.ModelTier }},
	{"question", ticket.FieldQuestion, func(t *ticket.Ticket) string { return t.Question }},
}

// startEdit opens the field editor on the ticket in the detail pane.
func (m *Model) startEdit() {
	if m.detail == nil {
		return
	}
	m.mode = modeEdit
	m.editIdx = 0
	m.editBuf = editableFields[0].get(m.detail)
}

// commitEdit writes the field being edited and reloads, so what the pane shows
// afterwards is what is on disk rather than what was typed.
func (m *Model) commitEdit() {
	if m.detail == nil {
		return
	}
	f := editableFields[m.editIdx]
	if f.get(m.detail) == m.editBuf {
		return
	}
	id := m.detail.ID
	if _, err := m.store.Mutate(id, func(t *ticket.Ticket) error {
		return t.Doc().SetScalar(f.field, m.editBuf)
	}); err != nil {
		m.notify(err.Error(), true)
		return
	}
	// Readiness is derived from the gate, so it has to be recomputed after any
	// field that feeds it changes.
	if _, err := m.store.Mutate(id, func(t *ticket.Ticket) error {
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	}); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if t, err := m.store.Load(id); err == nil {
		m.detail = t
	}
}

// editKey handles a keypress while a field is being edited.
func (m *Model) editKey(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "enter":
		m.commitEdit()
		m.mode = modeDetail
		m.editBuf = ""
		return m, nil
	case "esc":
		m.mode = modeDetail
		m.editBuf = ""
		return m, nil
	case "tab", "down":
		m.commitEdit()
		m.editIdx = (m.editIdx + 1) % len(editableFields)
		m.editBuf = editableFields[m.editIdx].get(m.detail)
		return m, nil
	case "shift+tab", "up":
		m.commitEdit()
		m.editIdx = (m.editIdx - 1 + len(editableFields)) % len(editableFields)
		m.editBuf = editableFields[m.editIdx].get(m.detail)
		return m, nil
	case "backspace":
		// Delete a character, not a byte: the tickets this board is meant to
		// hold are not all ASCII.
		if r := []rune(m.editBuf); len(r) > 0 {
			m.editBuf = string(r[:len(r)-1])
		}
		return m, nil
	}
	// Text carries what the key actually produced, including multi-byte
	// characters. Gating on a one-byte string here is what dropped every umlaut
	// in the filter and create prompts.
	if k.Text != "" {
		m.editBuf += k.Text
	}
	return m, nil
}

// renderEdit draws the field editor over the detail pane.
func (m *Model) renderEdit() string {
	if m.detail == nil {
		return m.renderBoard()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s\n", styHandle.Render(ticket.Handle(m.detail.ID)),
		styLaneTitle.Render(truncate(m.detail.Title, max(1, min(m.width, 78)-10))))
	b.WriteString(styBar.Render(strings.Repeat("─", max(1, min(m.width, 78)))) + "\n")
	for i, f := range editableFields {
		val := f.get(m.detail)
		lead := "  "
		if i == m.editIdx {
			lead = styDoing.Render("→ ")
			val = m.editBuf + "▌"
		}
		fmt.Fprintf(&b, "%s%s %s\n", lead,
			styMeta.Render(fmt.Sprintf("%-9s", f.name)),
			truncate(val, max(1, min(m.width, 78)-13)))
	}
	b.WriteString("\n" + styMeta.Render(truncate(
		"tab next field · enter save · esc cancel", max(1, min(m.width, 78)))))
	return b.String()
}
