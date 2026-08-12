package tui

import (
	"os"
	"strings"
	"testing"
)

// The fallback chain matters here: EDITOR is unset on plenty of machines,
// including the one this was written on, and "no editor configured" is a worse
// answer than "we found vi".
func TestEditorCommandFallback(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")

	t.Setenv("EDITOR", "nano")
	if got := editorCommand(); got[0] != "nano" {
		t.Errorf("EDITOR ignored: %v", got)
	}

	// VISUAL wins over EDITOR, matching git and every other tool that reads both.
	t.Setenv("VISUAL", "vim")
	if got := editorCommand(); got[0] != "vim" {
		t.Errorf("VISUAL did not take precedence: %v", got)
	}

	// An editor with arguments must not be treated as one long filename.
	t.Setenv("VISUAL", "code --wait")
	got := editorCommand()
	if got[0] != "code" || len(got) != 2 || got[1] != "--wait" {
		t.Errorf("editor arguments were not split: %v", got)
	}
}

func TestEditorFallsBackWhenUnset(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	got := editorCommand()
	if len(got) == 0 || got[0] == "" {
		t.Fatal("no fallback editor was chosen")
	}
}

// Editing must go through the store's write path, not write the file directly:
// the atomic write and the lock are what make concurrent sessions safe, and an
// editor that bypassed them would be the one writer that can corrupt a ticket.
func TestApplyExternalEditWritesThroughTheStore(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	id := m.detail.ID

	tmp, err := os.CreateTemp(t.TempDir(), "*.md")
	if err != nil {
		t.Fatal(err)
	}
	original := m.detail.Doc().Body()
	if _, err := tmp.WriteString(original + "\n## Plan\n\n- [~] rewritten by the editor\n"); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	if err := m.applyExternalEdit(id, tmp.Name()); err != nil {
		t.Fatalf("applyExternalEdit: %v", err)
	}

	reloaded, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reloaded.Body, "rewritten by the editor") {
		t.Errorf("the edited body was not written back: %q", reloaded.Body)
	}
	if len(reloaded.PlanItems) != 1 || reloaded.PlanItems[0].State.String() != "doing" {
		t.Errorf("the new checklist was not parsed: %+v", reloaded.PlanItems)
	}
	// Frontmatter must survive untouched — only the body is handed to the editor.
	if reloaded.Title != m.detail.Title {
		t.Errorf("title changed: %q -> %q", m.detail.Title, reloaded.Title)
	}
}

// An editor that exits without saving, or on an empty file, must not wipe the
// ticket body.
func TestEmptyExternalEditIsRefused(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	id := m.detail.ID
	before, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}

	empty, err := os.CreateTemp(t.TempDir(), "*.md")
	if err != nil {
		t.Fatal(err)
	}
	empty.Close()

	if err := m.applyExternalEdit(id, empty.Name()); err == nil {
		t.Error("an empty edit was accepted")
	}
	after, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if after.Body != before.Body {
		t.Error("the body changed despite the edit being refused")
	}
}
