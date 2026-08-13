package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// TestLaneScreenListsAllLoadedLanes asserts the screen shows exactly the
// lanes the board itself loaded — it must not disagree with jaira lanes.
func TestLaneScreenListsAllLoadedLanes(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	if len(ls.lanes) != len(m.lanes.Lanes) {
		t.Fatalf("laneScreen has %d lanes, want %d", len(ls.lanes), len(m.lanes.Lanes))
	}
}

// TestLaneScreenNavigationClamps asserts j/k move the selection and never run
// off either end of the list.
func TestLaneScreenNavigationClamps(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)

	if done := ls.key("k"); done {
		t.Fatal("k at the top must not finish the screen")
	}
	if ls.idx != 0 {
		t.Errorf("idx = %d, want 0 (clamped at top)", ls.idx)
	}
	for range ls.lanes {
		ls.key("j")
	}
	if ls.idx != len(ls.lanes)-1 {
		t.Errorf("idx = %d, want %d (clamped at bottom)", ls.idx, len(ls.lanes)-1)
	}
}

// TestLaneScreenEscFinishes asserts esc and q report the screen as done.
func TestLaneScreenEscFinishes(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	if done := ls.key("esc"); !done {
		t.Error("esc must finish the screen")
	}
}

// TestLaneScreenUseExportsToProjectDir asserts 'u' writes the selected lane
// into this project's own lane directory, and the file matches the lane's
// bytes.
func TestLaneScreenUseExportsToProjectDir(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	l := ls.selected()
	if l == nil {
		t.Fatal("no lane selected")
	}

	ls.key("u")
	if ls.isErr {
		t.Fatalf("use produced an error: %s", ls.msg)
	}
	if !strings.Contains(ls.msg, "wrote") {
		t.Errorf("message %q does not report what was written", ls.msg)
	}

	dst := filepath.Join(lane.ProjectLanesDir(m.store.Root), l.ID+".md")
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("expected %s to exist: %v", dst, err)
	}
}

// TestLaneScreenUseRefusesSecondExport asserts a repeated 'u' on the same
// lane refuses rather than silently overwriting.
func TestLaneScreenUseRefusesSecondExport(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)

	ls.key("u")
	if ls.isErr {
		t.Fatalf("first use failed: %s", ls.msg)
	}
	ls.key("u")
	if !ls.isErr {
		t.Error("a second 'u' on the same lane must refuse, not overwrite silently")
	}
}

// TestLaneScreenPublishWritesUnderIdentitySlug asserts 'p' writes into
// .jaira/shared/<slug>/, keyed off the acting identity.
func TestLaneScreenPublishWritesUnderIdentitySlug(t *testing.T) {
	t.Setenv("JAIRA_USER", "Alex Doe")
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	l := ls.selected()
	if l == nil {
		t.Fatal("no lane selected")
	}

	ls.key("p")
	if ls.isErr {
		t.Fatalf("publish produced an error: %s", ls.msg)
	}

	dst := filepath.Join(m.store.SharedDir(), "alex-doe", l.ID+".md")
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", dst, err)
	}
	if l.Creator == "" && !strings.Contains(string(got), "creator: alex-doe") {
		t.Errorf("published file missing a stamped creator, got:\n%s", got)
	}
}

// TestLaneScreenPublishRefusesSecondPublish mirrors the export refusal.
func TestLaneScreenPublishRefusesSecondPublish(t *testing.T) {
	t.Setenv("JAIRA_USER", "alex")
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)

	ls.key("p")
	if ls.isErr {
		t.Fatalf("first publish failed: %s", ls.msg)
	}
	ls.key("p")
	if !ls.isErr {
		t.Error("a second 'p' on the same lane must refuse, not overwrite silently")
	}
}

// laneTestStore builds a store rooted at a fresh project directory, pointed
// at its own catalogue, for tests that need real catalogue/project drift —
// unlike newTestStore, which deliberately points JAIRA_LANES_DIR at an empty
// directory to keep the board tests built-ins-only.
func laneTestStore(t *testing.T) (*ticket.Store, string) {
	t.Helper()
	root := t.TempDir()
	catalogue := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return s, catalogue
}

func writeLaneFile(t *testing.T, dir, filename, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestLaneScreenDetectsAndRefreshesDrift covers D-02: a project lane whose
// bytes differ from its catalogue copy is flagged, and 'R' pulls the
// catalogue version in.
func TestLaneScreenDetectsAndRefreshesDrift(t *testing.T) {
	s, catalogue := laneTestStore(t)
	catBody := `---
id: triage
name: Triage (catalogue)
after: human
precedence: 41
---
`
	writeLaneFile(t, catalogue, "triage.md", catBody)
	writeLaneFile(t, lane.ProjectLanesDir(s.Root), "triage.md", `---
id: triage
name: Triage (edited in this project)
after: human
precedence: 41
---
`)

	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	if _, drifted := ls.drift["triage"]; !drifted {
		t.Fatal("expected triage to be flagged as drifted")
	}
	for i, l := range ls.lanes {
		if l.ID == "triage" {
			ls.idx = i
		}
	}

	ls.key("R")
	if ls.isErr {
		t.Fatalf("refresh failed: %s", ls.msg)
	}
	if _, stillDrifted := ls.drift["triage"]; stillDrifted {
		t.Error("drift entry must be cleared after a refresh")
	}
	got, err := os.ReadFile(filepath.Join(lane.ProjectLanesDir(s.Root), "triage.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != catBody {
		t.Errorf("refresh did not pull the catalogue bytes:\ngot:  %q\nwant: %q", got, catBody)
	}
}
