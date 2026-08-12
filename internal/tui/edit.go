package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/ticket"
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
		// A goal or a piece of context is often a few sentences, so enter inserts
		// a line rather than ending the edit. Saving is deliberate.
		m.editBuf += "\n"
		return m, nil
	case "ctrl+s", "ctrl+d":
		m.commitEdit()
		m.mode = modeDetail
		m.editBuf = ""
		return m, nil
	case "esc":
		m.mode = modeDetail
		m.editBuf = ""
		return m, nil
	case "tab":
		m.commitEdit()
		m.editIdx = (m.editIdx + 1) % len(editableFields)
		m.editBuf = editableFields[m.editIdx].get(m.detail)
		return m, nil
	case "shift+tab":
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

// renderEdit draws the field editor.
//
// The field being edited gets the whole width and as many rows as it needs. The
// earlier version gave every field one truncated line, which made writing a
// paragraph of context in a maximised terminal worse than writing it in the CLI
// — the space was there and went unused.
func (m *Model) renderEdit() string {
	if m.detail == nil {
		return m.renderBoard()
	}
	w := max(20, m.width-2)
	var b strings.Builder

	fmt.Fprintf(&b, "%s  %s\n", styHandle.Render(ticket.Handle(m.detail.ID)),
		styLaneTitle.Render(truncate(m.detail.Title, max(1, w-10))))
	b.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	// Rows left for the editing box, once the other fields and the footer have
	// taken theirs. One line each for the unfocused fields.
	overhead := len(editableFields) + 6
	boxRows := max(3, m.height-overhead)

	for i, f := range editableFields {
		label := styMeta.Render(fmt.Sprintf("%-9s", f.name))
		if i != m.editIdx {
			val := strings.ReplaceAll(f.get(m.detail), "\n", " ⏎ ")
			if strings.TrimSpace(val) == "" {
				val = styMeta.Render("—")
			}
			fmt.Fprintf(&b, "  %s %s\n", label, truncate(val, max(1, w-13)))
			continue
		}

		fmt.Fprintf(&b, "%s %s\n", styDoing.Render("→ "), styDoing.Render(f.name))
		// The buffer is wrapped to the pane rather than truncated: this is where
		// the text is being written, so all of it has to be visible.
		lines := strings.Split(m.editBuf+"▌", "\n")
		var shown []string
		for _, ln := range lines {
			shown = append(shown, strings.Split(wrap(ln, max(10, w-4), 0), "\n")...)
		}
		// Keep the end of what is being typed on screen.
		if len(shown) > boxRows {
			shown = shown[len(shown)-boxRows:]
		}
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(colAccent).
			Width(max(10, w-2)).Height(min(boxRows, max(1, len(shown))))
		b.WriteString(box.Render(strings.Join(shown, "\n")) + "\n")
	}

	b.WriteString(styMeta.Render(truncate(
		"enter newline · ctrl+s save · tab next field · esc cancel", w)))
	return b.String()
}
