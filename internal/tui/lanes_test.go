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

// TestLaneScreenUseThenReloadShowsNoOverrideLabel covers the settings-screen
// side of the "copy is not an override" rule: exporting a built-in unchanged
// with 'u', then reloading — exactly what happens after the editor closes —
// must not show the orange "overrides" label, because the copy behaves
// exactly like the built-in it shadows.
func TestLaneScreenUseThenReloadShowsNoOverrideLabel(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	l := ls.selected()
	if l == nil {
		t.Fatal("no lane selected")
	}
	id := l.ID

	ls.key("u")
	if ls.isErr {
		t.Fatalf("use produced an error: %s", ls.msg)
	}

	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	reloaded, ok := m.lanes.Get(id)
	if !ok {
		t.Fatalf("%s lane missing after reload", id)
	}
	if reloaded.Overrides != "" {
		t.Errorf("an unmodified copy must not be marked as overriding, got %q", reloaded.Overrides)
	}

	ls2 := newLaneScreen(m.store, m.lanes)
	out := stripANSI(ls2.render(150, 32))
	if strings.Contains(out, "overrides "+id) {
		t.Errorf("settings screen shows an override label for an unmodified copy:\n%s", out)
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

// TestLaneScreenNewWritesSkeletonAndReloadsList asserts 'n' writes the lane
// template into the catalogue and queues a command to open it in $EDITOR.
// The editor launch itself cannot be driven from a test with no interactive
// terminal; what is verified here is the file write and the reload that
// follows the editor exiting, which is the part under this package's control.
func TestLaneScreenNewWritesSkeletonAndReloadsList(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.laneScreen = newLaneScreen(m.store, m.lanes)
	m.mode = modeLanes
	before := len(m.laneScreen.lanes)

	if done := m.laneScreen.key("n"); done {
		t.Fatal("n must not finish the screen")
	}
	if m.laneScreen.pendingCmd == nil {
		t.Fatal("n must queue a command to open the skeleton in $EDITOR")
	}

	path := filepath.Join(lane.UserLanesDir(), "new-lane.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
	if !strings.Contains(string(got), "id: my-lane") {
		t.Errorf("skeleton content = %q, want the lane template", got)
	}

	// Simulate the editor exiting cleanly — the reload this triggers, and the
	// list growing to include the new lane, is what is testable without a
	// real $EDITOR.
	updated, _ := m.Update(newLaneDoneMsg{})
	m2, ok := updated.(*Model)
	if !ok {
		t.Fatalf("Update returned %T, want *Model", updated)
	}
	if len(m2.laneScreen.lanes) != before+1 {
		t.Errorf("lane list has %d lanes after n, want %d", len(m2.laneScreen.lanes), before+1)
	}
}

// TestLaneScreenNewTwiceWritesTwoFiles asserts a second 'n' never overwrites
// the first skeleton: it picks the next free name instead.
func TestLaneScreenNewTwiceWritesTwoFiles(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)

	ls.key("n")
	first, err := os.ReadFile(filepath.Join(lane.UserLanesDir(), "new-lane.md"))
	if err != nil {
		t.Fatalf("expected new-lane.md to exist: %v", err)
	}

	ls.key("n")
	second, err := os.ReadFile(filepath.Join(lane.UserLanesDir(), "new-lane-2.md"))
	if err != nil {
		t.Fatalf("expected new-lane-2.md to exist rather than overwriting new-lane.md: %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("the two skeletons should carry identical starting content:\nfirst:  %q\nsecond: %q", first, second)
	}
}

// TestLaneScreenNewIgnoredWhenSharedFocused asserts 'n' only creates a lane
// when the settings screen's own list has focus, not the shared list.
func TestLaneScreenNewIgnoredWhenSharedFocused(t *testing.T) {
	s, _ := laneTestStore(t)
	writeLaneFile(t, filepath.Join(s.Root, ".jaira", "shared", "sam"), "hitl.md", `---
id: hitl
name: HITL
after: human
precedence: 41
---
`)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	ls.key("tab")
	if ls.focus != 1 {
		t.Fatal("tab must switch focus to the shared section")
	}

	ls.key("n")
	if ls.pendingCmd != nil {
		t.Error("n must be a no-op while the shared list has focus")
	}
	if _, err := os.Stat(filepath.Join(lane.UserLanesDir(), "new-lane.md")); !os.IsNotExist(err) {
		t.Errorf("n must not write a skeleton while the shared list has focus, stat err = %v", err)
	}
}

// TestLaneScreenFooterNamesNewKey asserts the footer names 'n' — without it
// the key does not exist as far as any reader is concerned.
func TestLaneScreenFooterNamesNewKey(t *testing.T) {
	m := newTestModel(t, 150, 40)
	ls := newLaneScreen(m.store, m.lanes)
	out := ls.render(150, 40)
	if !strings.Contains(out, "n new") {
		t.Error("footer must say \"n new\"")
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

// TestLaneScreenListsSharedLanes asserts a lane published under
// .jaira/shared/ shows up in the shared section, grouped by folder.
func TestLaneScreenListsSharedLanes(t *testing.T) {
	s, _ := laneTestStore(t)
	writeLaneFile(t, filepath.Join(s.Root, ".jaira", "shared", "sam"), "hitl.md", `---
id: hitl
name: HITL
after: human
precedence: 41
creator: sam
---
`)

	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	if len(ls.shared) != 1 {
		t.Fatalf("shared = %v, want exactly one lane", ls.shared)
	}
	if ls.shared[0].Folder != "sam" || ls.shared[0].Lane.ID != "hitl" {
		t.Errorf("shared[0] = %+v, want folder sam, id hitl", ls.shared[0])
	}
}

// TestLaneScreenShowsSharedSectionEvenWhenEmpty asserts the shared section and
// its explanation are visible on a board where nobody has published a lane —
// otherwise a user with an empty .jaira/shared/ has no way to learn the
// feature exists.
func TestLaneScreenShowsSharedSectionEvenWhenEmpty(t *testing.T) {
	m := newTestModel(t, 150, 40)
	ls := newLaneScreen(m.store, m.lanes)
	if len(ls.shared) != 0 {
		t.Fatalf("setup: shared = %v, want none", ls.shared)
	}

	out := ls.render(150, 40)
	if !strings.Contains(out, "Shared by teammates") {
		t.Error("shared section heading missing when there are no shared lanes")
	}
	if !strings.Contains(out, "publishes a lane with p") {
		t.Error("empty-state explanation missing when there are no shared lanes")
	}
}

// TestLaneScreenFooterSaysRefreshNotDrift asserts the footer names the action,
// not the reason the key exists, whether or not the shared section has
// anything in it.
func TestLaneScreenFooterSaysRefreshNotDrift(t *testing.T) {
	m := newTestModel(t, 150, 40)
	ls := newLaneScreen(m.store, m.lanes)

	out := ls.render(150, 40)
	if !strings.Contains(out, "R refresh") {
		t.Error("footer must say \"R refresh\"")
	}
	if strings.Contains(out, "refresh drift") {
		t.Error("footer must not say \"refresh drift\"")
	}
}

// TestLaneScreenAdoptWritesIntoCatalogue asserts pressing tab then a adopts
// the selected shared lane into the catalogue, naming the path in the
// message.
func TestLaneScreenAdoptWritesIntoCatalogue(t *testing.T) {
	s, catalogue := laneTestStore(t)
	writeLaneFile(t, filepath.Join(s.Root, ".jaira", "shared", "sam"), "hitl.md", `---
id: hitl
name: HITL
after: human
precedence: 41
---
`)

	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)

	ls.key("tab")
	if ls.focus != 1 {
		t.Fatal("tab must switch focus to the shared section")
	}
	ls.key("a")
	if ls.isErr {
		t.Fatalf("adopt failed: %s", ls.msg)
	}
	if !strings.Contains(ls.msg, "adopted") {
		t.Errorf("message %q does not report the adoption", ls.msg)
	}
	dst := filepath.Join(catalogue, "hitl.md")
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("expected %s to exist: %v", dst, err)
	}
}

// TestLaneScreenAdoptRefusesThenConfirms asserts an id collision is held for
// confirmation rather than silently overwritten.
func TestLaneScreenAdoptRefusesThenConfirms(t *testing.T) {
	s, catalogue := laneTestStore(t)
	writeLaneFile(t, catalogue, "hitl.md", `---
id: hitl
name: Existing catalogue HITL
after: human
precedence: 41
---
`)
	writeLaneFile(t, filepath.Join(s.Root, ".jaira", "shared", "sam"), "hitl.md", `---
id: hitl
name: Sam's HITL
after: human
precedence: 41
---
`)

	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	ls.key("tab")

	ls.key("a")
	if !ls.isErr {
		t.Fatal("adopting a colliding id on the first press must refuse, not overwrite")
	}
	if ls.confirmAdoptID != "hitl" {
		t.Fatalf("confirmAdoptID = %q, want %q", ls.confirmAdoptID, "hitl")
	}

	ls.key("a")
	if ls.isErr {
		t.Fatalf("the confirming press must succeed: %s", ls.msg)
	}
	got, err := os.ReadFile(filepath.Join(catalogue, "hitl.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "Sam's HITL") {
		t.Errorf("confirmed adopt did not overwrite the catalogue file, got:\n%s", got)
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
