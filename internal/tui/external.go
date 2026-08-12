package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/ticket"
)

// editorCommand resolves the editor to launch.
//
// VISUAL before EDITOR is the convention every tool that reads both follows, and
// git's own core.editor is consulted next because a user who configured one there
// has already stated a preference. The final fallbacks exist because "no editor
// configured" is a worse answer than a working one — EDITOR is unset on plenty of
// machines.
func editorCommand() []string {
	for _, v := range []string{os.Getenv("VISUAL"), os.Getenv("EDITOR")} {
		if f := strings.Fields(v); len(f) > 0 {
			return f
		}
	}
	if out, err := exec.Command("git", "config", "--get", "core.editor").Output(); err == nil {
		if f := strings.Fields(strings.TrimSpace(string(out))); len(f) > 0 {
			return f
		}
	}
	for _, candidate := range []string{"nano", "vim", "vi"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return []string{candidate}
		}
	}
	return []string{"vi"}
}

// openInEditor suspends the board, hands the ticket's body to $EDITOR, and
// writes the result back when the editor exits.
//
// Only the body is handed over. The frontmatter stays under the tool's control,
// because a hand-edited field would bypass the gates and the derived readiness
// flag — and because the fields are already editable in place with e.
func (m *Model) openInEditor() (tea.Model, tea.Cmd) {
	if m.detail == nil {
		return m, nil
	}
	id := m.detail.ID
	dir, err := os.MkdirTemp("", "jaira-edit-")
	if err != nil {
		m.notify(err.Error(), true)
		return m, nil
	}
	path := filepath.Join(dir, ticket.Handle(id)+".md")
	if err := os.WriteFile(path, []byte(m.detail.Doc().Body()), 0o600); err != nil {
		m.notify(err.Error(), true)
		return m, nil
	}

	argv := append(editorCommand(), path)
	cmd := exec.Command(argv[0], argv[1:]...)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorDoneMsg{id: id, path: path, dir: dir, err: err}
	})
}

type editorDoneMsg struct {
	id, path, dir string
	err           error
}

// applyExternalEdit writes an edited body back through the store.
//
// It goes through store.Mutate rather than writing the file directly: the atomic
// write and the lock are what make concurrent sessions safe, and an editor that
// bypassed them would be the one writer able to corrupt a ticket.
func (m *Model) applyExternalEdit(id, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// An editor quit without saving, or one that truncated the file, must not be
	// able to erase a ticket's contents.
	if strings.TrimSpace(string(b)) == "" {
		return fmt.Errorf("the edited body was empty, so nothing was written")
	}
	body := string(b)
	if _, err := m.store.Mutate(id, func(t *ticket.Ticket) error {
		t.Doc().SetBody(body)
		return nil
	}); err != nil {
		return err
	}
	// The definition of done lives in the body, and readiness depends on whether
	// there is one, so it has to be recomputed after an edit.
	if _, err := m.store.Mutate(id, func(t *ticket.Ticket) error {
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	}); err != nil {
		return err
	}
	if err := m.reload(); err != nil {
		return err
	}
	if t, err := m.store.Load(id); err == nil {
		m.detail = t
	}
	return nil
}
