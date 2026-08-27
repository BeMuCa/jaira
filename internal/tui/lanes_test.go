package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

// TestLaneScreenNavigationClamps asserts h/l move the selected column and
// never run off either end — the last column being the '+' one.
func TestLaneScreenNavigationClamps(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)

	if done := ls.key("h"); done {
		t.Fatal("h at the first column must not finish the screen")
	}
	if ls.idx != 0 {
		t.Errorf("idx = %d, want 0 (clamped at the first column)", ls.idx)
	}
	for range ls.lanes {
		ls.key("l")
	}
	ls.key("l") // one more than there are lanes: lands on '+', then clamps
	if ls.idx != len(ls.lanes) {
		t.Errorf("idx = %d, want %d (clamped at the '+' column)", ls.idx, len(ls.lanes))
	}
	if !ls.isPlusColumn() {
		t.Error("idx at len(lanes) must report as the '+' column")
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

// TestLaneScreenEveryBoardLaneIsAFileOfThisBoard covers the rule the screen
// now sits on: a board is its lane directory, so every lane listed has a file
// there, and an unchanged copy of a shipped lane is not marked as differing.
func TestLaneScreenEveryBoardLaneIsAFileOfThisBoard(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	dir := lane.ProjectLanesDir(m.store.Root)
	if len(ls.lanes) == 0 {
		t.Fatal("no lanes on the board")
	}
	for _, l := range ls.lanes {
		if !strings.HasPrefix(l.Source, dir) {
			t.Errorf("lane %q comes from %s, want a file under %s", l.ID, l.Source, dir)
		}
		if l.Overrides != "" {
			t.Errorf("unchanged copy %q is marked as differing from the shipped lane", l.ID)
		}
		if l.Builtin {
			t.Errorf("lane %q is marked Builtin on a board; nothing is injected any more", l.ID)
		}
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

	path := filepath.Join(lane.ProjectLanesDir(m.store.Root), "new-lane.md")
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
	first, err := os.ReadFile(filepath.Join(lane.ProjectLanesDir(m.store.Root), "new-lane.md"))
	if err != nil {
		t.Fatalf("expected new-lane.md to exist: %v", err)
	}

	ls.key("n")
	second, err := os.ReadFile(filepath.Join(lane.ProjectLanesDir(m.store.Root), "new-lane-2.md"))
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
	if _, err := os.Stat(filepath.Join(lane.ProjectLanesDir(ls.store.Root), "new-lane.md")); !os.IsNotExist(err) {
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

// TestLaneScreenRendersBoardWithPlusColumn asserts the screen shows the
// project's lane names and a '+' column, not a vertical list. The board only
// fits a few columns at once (like the main board it mirrors, see
// renderBoard), so this scrolls all the way across, collecting every frame,
// rather than expecting one screen to hold every lane at once.
func TestLaneScreenRendersBoardWithPlusColumn(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	var seen strings.Builder
	seen.WriteString(ls.render(150, 32))
	for range ls.lanes {
		ls.key("l")
		seen.WriteString(ls.render(150, 32))
	}
	out := seen.String()
	for _, l := range ls.lanes {
		if !strings.Contains(out, l.Name) {
			t.Errorf("rendered board missing lane name %q", l.Name)
		}
	}
	if !strings.Contains(out, "+") {
		t.Error("rendered board missing the '+' column")
	}
}

// TestLaneScreenLMovesSelectedLaneAndWritesOrderFile asserts 'L' moves the
// selected lane one step right and persists it to the order file, through
// the same lane.MoveLane the CLI's 'lanes move' calls.
func TestLaneScreenLMovesSelectedLaneAndWritesOrderFile(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	firstID, secondID := ls.lanes[0].ID, ls.lanes[1].ID

	ls.key("L")
	if ls.isErr {
		t.Fatalf("L produced an error: %s", ls.msg)
	}
	if ls.lanes[0].ID != secondID || ls.lanes[1].ID != firstID {
		t.Fatalf("after L, lanes = %v, want %s then %s", ls.lanes[0:2], secondID, firstID)
	}
	ids, err := lane.LoadOrder(m.store.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		t.Fatal("L did not write an order file")
	}
	if ls.idx != 1 {
		t.Errorf("idx = %d, want 1 (the moved lane stays selected at its new position)", ls.idx)
	}
}

// TestLaneScreenXRemovesSelectedLane asserts 'x' opens a yes/no confirmation
// (default no, so a bare enter removes nothing) and only removes the
// selected lane from the project, through lane.Remove, once 'right'/'enter'
// selects yes. Uses a fresh, ticketless store (laneTestStore) rather than
// newTestModel's, which seeds sample tickets into backlog and would
// otherwise make the very first lane refuse removal.
func TestLaneScreenXRemovesSelectedLane(t *testing.T) {
	s, _ := laneTestStore(t)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	before := len(ls.lanes)
	id := ls.lanes[0].ID

	// A bare enter (default no) must remove nothing — the quick-enter
	// protection the confirmation exists for.
	ls.key("x")
	if ls.confirm == nil {
		t.Fatal("x must open a confirmation before removing anything")
	}
	if ls.confirm.yes {
		t.Error("confirmation must default to no")
	}
	ls.key("enter")
	if len(ls.lanes) != before {
		t.Fatalf("lanes = %d after x+enter (default no), want %d (nothing removed)", len(ls.lanes), before)
	}
	if ls.confirm != nil {
		t.Error("enter must clear the confirmation")
	}

	ls.key("x")
	ls.key("right")
	if !ls.confirm.yes {
		t.Error("right must select yes")
	}
	ls.key("enter")
	if ls.isErr {
		t.Fatalf("x produced an error: %s", ls.msg)
	}
	if len(ls.lanes) != before-1 {
		t.Fatalf("lanes = %d after x,right,enter, want %d", len(ls.lanes), before-1)
	}
	if _, ok := ls.set.Get(id); ok {
		t.Errorf("removed lane %q still present", id)
	}
}

// TestLaneScreenXRefusesWhenTicketPresent asserts 'x', right, enter is
// refused, naming the ticket, when one currently sits in the selected lane —
// the same refusal the CLI's 'lanes remove' gives.
func TestLaneScreenXRefusesWhenTicketPresent(t *testing.T) {
	m := newTestModel(t, 150, 32)
	ls := newLaneScreen(m.store, m.lanes)
	target := ls.lanes[0]

	tk, err := m.store.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  "occupies " + target.ID,
		ticket.FieldStatus: target.ID,
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	ls.key("x")
	ls.key("right")
	ls.key("enter")
	if !ls.isErr {
		t.Fatal("x on an occupied lane must refuse, not remove it")
	}
	if !strings.Contains(ls.msg, ticket.Handle(tk.ID)) {
		t.Errorf("refusal %q does not name the occupying ticket", ls.msg)
	}
	if _, ok := ls.set.Get(target.ID); !ok {
		t.Error("a refused removal must leave the lane in place")
	}
}

// TestLaneScreenXDeletesAvailableCatalogueLane asserts x, right, enter on a
// not-installed (available) catalogue lane deletes its file from disk and it
// disappears from ls.available after the reload x triggers.
func TestLaneScreenXDeletesAvailableCatalogueLane(t *testing.T) {
	s, catalogue := laneTestStore(t)
	// Materialise the project directory first: with no project directory of
	// its own, D-03's fallback merges the whole catalogue into the board
	// itself, and my-lane would never show up as merely available — see
	// TestLaneScreenPlusColumnOpensCatalogueListingOnlyMissingLanes.
	builtins, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range builtins {
		if _, err := lane.Export(l, lane.ProjectLanesDir(s.Root), false); err != nil {
			t.Fatal(err)
		}
	}
	writeLaneFile(t, catalogue, "new-lane.md", `---
id: my-lane
name: My Lane
after: human
precedence: 41
---
`)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	if len(ls.available) != 1 || ls.available[0].ID != "my-lane" {
		t.Fatalf("setup: available = %v, want exactly [my-lane]", ls.available)
	}
	ls.idx = len(ls.lanes) // the available lane's column

	path := ls.available[0].Source
	ls.key("x")
	if ls.confirm == nil || ls.confirm.path != path {
		t.Fatalf("x on an available lane must open a confirmation naming its file, got confirm = %+v", ls.confirm)
	}
	ls.key("right")
	ls.key("enter")
	if ls.isErr {
		t.Fatalf("deleting the catalogue lane produced an error: %s", ls.msg)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected %s to be deleted, stat err = %v", path, err)
	}
	for _, l := range ls.available {
		if l.ID == "my-lane" {
			t.Error("deleted lane must not still show up in available after reload")
		}
	}
}

// TestLaneScreenXOnAvailableDefaultNoLeavesFileOnDisk asserts a bare enter
// (default no) after x on an available catalogue lane leaves its file in
// place.
func TestLaneScreenXOnAvailableDefaultNoLeavesFileOnDisk(t *testing.T) {
	s, catalogue := laneTestStore(t)
	builtins, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range builtins {
		if _, err := lane.Export(l, lane.ProjectLanesDir(s.Root), false); err != nil {
			t.Fatal(err)
		}
	}
	writeLaneFile(t, catalogue, "new-lane.md", `---
id: my-lane
name: My Lane
after: human
precedence: 41
---
`)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	if len(ls.available) != 1 || ls.available[0].ID != "my-lane" {
		t.Fatalf("setup: available = %v, want exactly [my-lane]", ls.available)
	}
	ls.idx = len(ls.lanes)
	path := ls.available[0].Source

	ls.key("x")
	ls.key("enter")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("default no must leave %s in place, stat err = %v", path, err)
	}
}

// TestLaneScreenXOnAvailableBuiltinRefusesWithoutConfirmation asserts x on an
// available built-in lane (one removed from the board, so it reappears as
// available with Builtin true) shows the "nothing to delete" error and opens
// no confirmation, since a built-in has no file to delete.
func TestLaneScreenXOnAvailableBuiltinRefusesWithoutConfirmation(t *testing.T) {
	s, _ := laneTestStore(t)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	builtinID := ls.lanes[0].ID

	ls.key("x")
	ls.key("right")
	ls.key("enter")
	if ls.isErr {
		t.Fatalf("removing the builtin from the board failed: %s", ls.msg)
	}

	var found *lane.Lane
	for _, l := range ls.available {
		if l.ID == builtinID {
			found = l
		}
	}
	if found == nil || !found.Builtin {
		t.Fatalf("expected %q to reappear as an available builtin, available = %v", builtinID, ls.available)
	}
	ls.idx = indexOfLane(ls.available, builtinID) + len(ls.lanes)

	ls.key("x")
	if !ls.isErr {
		t.Fatal("x on an available built-in must refuse, not open a confirmation")
	}
	if !strings.Contains(ls.msg, builtinID) {
		t.Errorf("refusal %q does not name the built-in lane", ls.msg)
	}
	if ls.confirm != nil {
		t.Error("x on a built-in must not open a confirmation")
	}
}

// TestLaneScreenConfirmEscCancelsWithoutFinishing asserts esc while the
// confirmation is open clears it without finishing the screen, and the
// render before esc shows "no" before "yes".
func TestLaneScreenConfirmEscCancelsWithoutFinishing(t *testing.T) {
	s, _ := laneTestStore(t)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	before := len(ls.lanes)

	ls.key("x")
	out := ls.render(150, 32)
	if strings.Index(out, "no") > strings.Index(out, "yes") || !strings.Contains(out, "no") || !strings.Contains(out, "yes") {
		t.Errorf("confirmation render must show \"no\" before \"yes\":\n%s", out)
	}

	if done := ls.key("esc"); done {
		t.Fatal("esc while confirming must not finish the screen")
	}
	if ls.confirm != nil {
		t.Error("esc must clear the confirmation")
	}
	if len(ls.lanes) != before {
		t.Error("esc while confirming must not remove anything")
	}
}

// TestLaneScreenPlusColumnOpensCatalogueListingOnlyMissingLanes asserts enter
// on the '+' column lists exactly the lanes not already in this project.
//
// The project is materialised (every built-in exported as its own lane
// directory's copy) before the catalogue gains "extra": with no project
// directory of its own, D-03's fallback already shows the whole catalogue,
// so nothing would be missing from it to demonstrate the filter.
func TestLaneScreenPlusColumnOpensCatalogueListingOnlyMissingLanes(t *testing.T) {
	s, catalogue := laneTestStore(t)
	builtins, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range builtins {
		if _, err := lane.Export(l, lane.ProjectLanesDir(s.Root), false); err != nil {
			t.Fatal(err)
		}
	}
	writeLaneFile(t, catalogue, "extra.md", `---
id: extra
name: Extra
after: human
precedence: 41
---
`)
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(s, set)
	// The not-installed lane shows as its own dimmed column before '+'.
	if len(ls.available) != 1 || ls.available[0].ID != "extra" {
		t.Fatalf("available = %v, want exactly [extra]", ls.available)
	}
	ls.idx = len(ls.lanes) + len(ls.available) // the '+' column

	ls.key("enter")
	if !ls.catalogueOpen {
		t.Fatal("enter on the '+' column must open the catalogue")
	}
	if len(ls.catalogue) != 1 || ls.catalogue[0].ID != "extra" {
		t.Fatalf("catalogue = %v, want exactly [extra]", ls.catalogue)
	}

	ls.key("enter")
	if ls.isErr {
		t.Fatalf("adding from the catalogue failed: %s", ls.msg)
	}
	if ls.catalogueOpen {
		t.Error("choosing a lane must close the catalogue")
	}
	if _, ok := ls.set.Get("extra"); !ok {
		t.Error("chosen lane was not added to the project")
	}
}
