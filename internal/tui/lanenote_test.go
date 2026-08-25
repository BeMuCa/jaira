package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The settings screen changes the pipeline, and its reload runs after every
// change it makes — so that is where the note has to catch up. Without this the
// board could be reordered from the TUI and the agent file would keep
// describing the old route until someone happened to run 'jaira update'.
func TestLaneScreenReloadRefreshesTheAgentNote(t *testing.T) {
	m := newTestModel(t, 150, 32)
	root := m.store.Root

	// Start from a note that names no lanes at all, as a board written before
	// the block described itself would have.
	stale := "<!-- jaira:start -->\nold text that names no lanes\n<!-- jaira:end -->\n"
	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(stale), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	ls := newLaneScreen(m.store, m.lanes)
	if err := ls.reload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if ls.isErr {
		t.Fatalf("reload reported: %s", ls.msg)
	}

	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(got), "This board's lanes") {
			t.Errorf("%s was not brought up to date:\n%s", name, got)
		}
		if !strings.Contains(string(got), "Order: ") {
			t.Errorf("%s carries no route:\n%s", name, got)
		}
		if strings.Contains(string(got), "old text that names no lanes") {
			t.Errorf("%s kept the stale block:\n%s", name, got)
		}
	}
}
