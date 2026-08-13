package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The board and the compact view both used to show nothing until 'p' was
// pressed once, because m.projects was assigned in exactly one place: inside
// the 'p' handler. These tests build a model and read its list, and read the
// views that depend on it, without ever sending a key.

// TestProjectsLoadedAtBuildNotOnlyOnP covers the bug directly: a freshly built
// model must already know the recorded boards.
func TestProjectsLoadedAtBuildNotOnlyOnP(t *testing.T) {
	s := newTestStore(t) // sets JAIRA_HOME for us
	other := t.TempDir()
	otherStore, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := otherStore.Init(); err != nil {
		t.Fatal(err)
	}
	project.Remember(s.Root)
	project.Remember(other)

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.projects) != 2 {
		t.Fatalf("New() did not load the recorded boards: got %d, want 2", len(m.projects))
	}
}

// TestCompactViewTabsRenderWithoutP covers the visible symptom: the compact
// view's numbered tabs must appear on the very first render.
func TestCompactViewTabsRenderWithoutP(t *testing.T) {
	s := newTestStore(t)
	other := t.TempDir()
	otherStore, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := otherStore.Init(); err != nil {
		t.Fatal(err)
	}
	project.Remember(s.Root)
	project.Remember(other)

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 130, 34
	m.mode = modePipeline

	out := stripANSI(m.renderPipeline())
	for i, p := range m.projects {
		want := fmt.Sprintf("%d %s", i+1, p.Name)
		if !strings.Contains(out, want) {
			t.Errorf("compact view missing tab %q before p was ever pressed:\n%s", want, out)
		}
	}
}

// TestBoardViewShowsBoardsWithoutP covers the board view side of the same
// bug: the numbered line above the columns must also need no 'p' first.
func TestBoardViewShowsBoardsWithoutP(t *testing.T) {
	s := newTestStore(t)
	other := t.TempDir()
	otherStore, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := otherStore.Init(); err != nil {
		t.Fatal(err)
	}
	project.Remember(s.Root)
	project.Remember(other)

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32

	out := stripANSI(m.render())
	for i, p := range m.projects {
		want := fmt.Sprintf("%d %s", i+1, p.Name)
		if !strings.Contains(out, want) {
			t.Errorf("board view missing tab %q before p was ever pressed:\n%s", want, out)
		}
	}
}

// TestBoardViewNumberSwitchesProject covers the other half of the bug report:
// the board view had no 1-9 binding at all.
func TestBoardViewNumberSwitchesProject(t *testing.T) {
	m := newTestModel(t, 130, 34)
	m.mode = modeBoard
	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: "/somewhere/else", Name: "other"},
	}

	m.key(tea.KeyPressMsg{Code: '2', Text: "2"})
	if m.SwitchTo != "/somewhere/else" {
		t.Errorf("SwitchTo = %q, want /somewhere/else", m.SwitchTo)
	}

	// Choosing the board already open is a no-op, not a pointless reload.
	m2 := newTestModel(t, 130, 34)
	m2.mode = modeBoard
	m2.projects = []project.Project{{Root: m2.store.Root, Name: "current"}}
	m2.key(tea.KeyPressMsg{Code: '1', Text: "1"})
	if m2.SwitchTo != "" {
		t.Errorf("SwitchTo was set to %q for the board already open", m2.SwitchTo)
	}
	if m2.mode != modeBoard {
		t.Errorf("choosing the current board left board mode: %v", m2.mode)
	}
}

// TestBoardViewNumberBeyondListIsIgnored covers "a number beyond the list is
// ignored, not a panic" — the point of the test is that this does not panic.
func TestBoardViewNumberBeyondListIsIgnored(t *testing.T) {
	m := newTestModel(t, 130, 34)
	m.mode = modeBoard
	m.projects = []project.Project{{Root: m.store.Root, Name: "current"}}

	m.key(tea.KeyPressMsg{Code: '9', Text: "9"})
	if m.SwitchTo != "" {
		t.Errorf("SwitchTo was set for a number with no board behind it: %q", m.SwitchTo)
	}
	if m.mode != modeBoard {
		t.Errorf("an unused number left board mode: %v", m.mode)
	}
}
