package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/session"
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
	other := newTestStore(t)
	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: other.Root, Name: "other"},
	}

	m.key(tea.KeyPressMsg{Code: '2', Text: "2"})
	if m.store.Root != other.Root {
		t.Errorf("store root = %q, want the other board %q", m.store.Root, other.Root)
	}
	if m.mode != modeBoard {
		t.Errorf("switching boards left board mode: %v", m.mode)
	}

	// Choosing the board already open is a no-op, not a pointless reload.
	m2 := newTestModel(t, 130, 34)
	m2.mode = modeBoard
	m2.projects = []project.Project{{Root: m2.store.Root, Name: "current"}}
	before := m2.store
	m2.key(tea.KeyPressMsg{Code: '1', Text: "1"})
	if m2.store != before {
		t.Error("choosing the current board swapped the store")
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

	before := m.store
	m.key(tea.KeyPressMsg{Code: '9', Text: "9"})
	if m.store != before {
		t.Error("a number with no board behind it swapped the store")
	}
	if m.mode != modeBoard {
		t.Errorf("an unused number left board mode: %v", m.mode)
	}
}

// A board that has a running agent is worth telling apart from one that does
// not, at a glance — these tests cover boardHasLiveSession directly, since it
// is what decides that.

func TestBoardHasLiveSession(t *testing.T) {
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	other := t.TempDir()
	st, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}

	if boardHasLiveSession(other) {
		t.Fatal("a board with no session file must not read as live")
	}

	fresh := session.Session{ID: "s1", Focus: "doing work", UpdatedAt: ticket.FormatTime(time.Now())}
	if err := session.Save(st, fresh); err != nil {
		t.Fatal(err)
	}
	if !boardHasLiveSession(other) {
		t.Error("a fresh session must mark the board live")
	}

	stale := session.Session{ID: "s1", Focus: "doing work", UpdatedAt: ticket.FormatTime(time.Now().Add(-2 * session.StaleAfter))}
	if err := session.Save(st, stale); err != nil {
		t.Fatal(err)
	}
	if boardHasLiveSession(other) {
		t.Error("a stale session must not mark the board live — a crashed agent must not look busy")
	}
}

// TestBoardHasLiveSessionUnreadableBoard covers a board that is gone, or was
// never there: reading it must never fail, only report not-live.
func TestBoardHasLiveSessionUnreadableBoard(t *testing.T) {
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	gone := filepath.Join(t.TempDir(), "does-not-exist", "nested", "deeper")
	if boardHasLiveSession(gone) {
		t.Error("a board directory that does not exist cannot be live")
	}
}

// TestRefreshLiveBoardsToleratesAGoneBoard covers the same requirement one
// level up: refreshing the whole switcher must not fail render just because
// one of several recorded boards can no longer be read.
func TestRefreshLiveBoardsToleratesAGoneBoard(t *testing.T) {
	m := newTestModel(t, 130, 34)
	live := t.TempDir()
	liveStore, err := ticket.At(live)
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Save(liveStore, session.Session{
		ID: "s1", Focus: "working", UpdatedAt: ticket.FormatTime(time.Now()),
	}); err != nil {
		t.Fatal(err)
	}
	gone := filepath.Join(t.TempDir(), "does-not-exist")

	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: live, Name: "live"},
		{Root: gone, Name: "gone"},
	}
	m.refreshLiveBoards() // must not panic

	if !m.liveBoards[live] {
		t.Error("a board with a fresh session must be marked live")
	}
	if m.liveBoards[gone] {
		t.Error("a board that cannot be read must render unmarked, not marked")
	}
}

// TestProjectTabsMarkLiveBoards covers the rendering side: the marker is a
// glyph, not colour alone, and only the live board carries it.
func TestProjectTabsMarkLiveBoards(t *testing.T) {
	m := newTestModel(t, 130, 34)
	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: "/somewhere/live", Name: "live"},
		{Root: "/somewhere/stale", Name: "stale"},
	}
	m.liveBoards = map[string]bool{"/somewhere/live": true}

	tabs := m.projectTabs()
	if len(tabs) != 3 {
		t.Fatalf("got %d tabs, want 3", len(tabs))
	}
	if strings.Contains(stripANSI(tabs[0]), "◆") {
		t.Errorf("current board incorrectly marked live: %q", stripANSI(tabs[0]))
	}
	if !strings.Contains(stripANSI(tabs[1]), "◆") {
		t.Errorf("live board not marked: %q", stripANSI(tabs[1]))
	}
	if strings.Contains(stripANSI(tabs[2]), "◆") {
		t.Errorf("stale board incorrectly marked live: %q", stripANSI(tabs[2]))
	}
}
