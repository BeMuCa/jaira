package tui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// catalogueHome points the per-user directory at a fresh JAIRA_HOME and
// clears JAIRA_LANES_DIR, so the catalogue is the one a person actually has:
// $JAIRA_HOME/lanes. It returns that catalogue directory.
func catalogueHome(t *testing.T) string {
	t.Helper()
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("JAIRA_HOME", home)
	t.Setenv("JAIRA_LANES_DIR", "")
	return filepath.Join(home, "lanes")
}

// initBoard makes a fresh repository with a jaira board whose lane directory
// is set up, under the JAIRA_HOME catalogueHome chose.
func initBoard(t *testing.T) *ticket.Store {
	t.Helper()
	s, err := ticket.At(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if _, err := lane.Load(s.Root); err != nil {
		t.Fatal(err)
	}
	return s
}

// offeredUnderPlus reports whether a board at root offers id in its '+'
// column — the lane screen's not-installed list.
func offeredUnderPlus(t *testing.T, s *ticket.Store, id string) bool {
	t.Helper()
	set, err := lane.Load(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range newLaneScreen(s, set).available {
		if l.ID == id {
			return true
		}
	}
	return false
}

const deployLane = `---
id: deploy-check
name: Deploy check
after: human
precedence: 41
---

Check the deploy.
`

// TestNewLaneLandsInCatalogueAndOnAnotherBoard: a lane made with n in one
// repository is, once the editor closes, also in ~/.jaira/lanes, and a second
// repository's lane screen offers it under '+'.
func TestNewLaneLandsInCatalogueAndOnAnotherBoard(t *testing.T) {
	catalogue := catalogueHome(t)
	first := initBoard(t)
	second := initBoard(t)

	m, err := New(first)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 40
	m.laneScreen = newLaneScreen(m.store, m.lanes)

	// What n and the editor leave behind: a skeleton the person filled in.
	path, err := writeLaneSkeleton(lane.ProjectLanesDir(first.Root))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(deployLane), 0o644); err != nil {
		t.Fatal(err)
	}
	m.Update(newLaneDoneMsg{path: path})

	if _, err := os.Stat(filepath.Join(catalogue, "deploy-check.md")); err != nil {
		t.Fatalf("new lane not in the catalogue: %v", err)
	}
	if m.laneScreen.isErr || !strings.Contains(m.laneScreen.msg, "catalogue") {
		t.Errorf("lane screen message = %q (err %v), want it to name the catalogue", m.laneScreen.msg, m.laneScreen.isErr)
	}
	if indexOfLane(m.laneScreen.lanes, "deploy-check") < 0 {
		t.Error("new lane must still be on the board it was made on")
	}
	if !offeredUnderPlus(t, second, "deploy-check") {
		t.Error("second repository does not offer the new lane under '+'")
	}
}

// TestNewLaneUneditedSkeletonStaysOutOfCatalogue: closing the editor on the
// untouched skeleton puts its placeholder id nowhere but this board.
func TestNewLaneUneditedSkeletonStaysOutOfCatalogue(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	path, err := writeLaneSkeleton(lane.ProjectLanesDir(s.Root))
	if err != nil {
		t.Fatal(err)
	}
	if msg, _ := catalogueNewLane(path); msg != "" {
		t.Errorf("message = %q, want none", msg)
	}
	if matches, _ := filepath.Glob(filepath.Join(catalogue, "*.md")); len(matches) != 0 {
		t.Errorf("catalogue = %v, want empty", matches)
	}
}

// TestNewLaneNeverOverwritesCatalogue: an id the catalogue already holds is
// kept as it is, and the lane screen says the new lane stayed on this board.
func TestNewLaneNeverOverwritesCatalogue(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	existing := strings.Replace(deployLane, "Check the deploy.", "The catalogue's own copy.", 1)
	writeLaneFile(t, catalogue, "deploy-check.md", existing)
	path := filepath.Join(lane.ProjectLanesDir(s.Root), "new-lane.md")
	writeLaneFile(t, lane.ProjectLanesDir(s.Root), "new-lane.md", deployLane)

	msg, isErr := catalogueNewLane(path)
	if !isErr || !strings.Contains(msg, "this board only") {
		t.Errorf("message = %q (err %v), want an error saying it stayed on this board", msg, isErr)
	}
	b, _ := os.ReadFile(filepath.Join(catalogue, "deploy-check.md"))
	if string(b) != existing {
		t.Error("catalogue copy was overwritten")
	}
}

// TestEditedLaneIsNotCatalogued: E on an existing lane reuses newLaneDoneMsg
// without a path, and must not copy anything into the catalogue.
func TestEditedLaneIsNotCatalogued(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	writeLaneFile(t, lane.ProjectLanesDir(s.Root), "deploy-check.md", deployLane)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.laneScreen = newLaneScreen(m.store, m.lanes)
	m.Update(newLaneDoneMsg{})
	if matches, _ := filepath.Glob(filepath.Join(catalogue, "*.md")); len(matches) != 0 {
		t.Errorf("catalogue = %v, want empty after editing an existing lane", matches)
	}
}

// TestAdoptedSharedLaneIsOfferedOnEveryBoard: a colleague's lane under
// .jaira/shared/<colleague>/, adopted with a, is in the catalogue and offered
// under '+' on a board in another repository.
func TestAdoptedSharedLaneIsOfferedOnEveryBoard(t *testing.T) {
	catalogue := catalogueHome(t)
	first := initBoard(t)
	second := initBoard(t)
	writeLaneFile(t, filepath.Join(first.SharedDir(), "sam"), "deploy-check.md", deployLane)

	set, err := lane.Load(first.Root)
	if err != nil {
		t.Fatal(err)
	}
	ls := newLaneScreen(first, set)
	ls.key("tab")
	ls.key("a")
	if ls.isErr {
		t.Fatalf("adopt failed: %s", ls.msg)
	}
	if _, err := os.Stat(filepath.Join(catalogue, "deploy-check.md")); err != nil {
		t.Fatalf("adopted lane not in the catalogue: %v", err)
	}
	for name, s := range map[string]*ticket.Store{"first": first, "second": second} {
		if !offeredUnderPlus(t, s, "deploy-check") {
			t.Errorf("%s repository does not offer the adopted lane under '+'", name)
		}
	}
}

// TestNewLaneWithBuiltinIDStaysOutOfCatalogue: a board-local lane that reuses
// a built-in's id is not copied into the catalogue, where it would stand in
// for the built-in on every board.
func TestNewLaneWithBuiltinIDStaysOutOfCatalogue(t *testing.T) {
	catalogue := catalogueHome(t)
	dir := t.TempDir()
	writeLaneFile(t, dir, "new-lane.md", strings.Replace(deployLane, "id: deploy-check", "id: review", 1))
	path := filepath.Join(dir, "new-lane.md")
	msg, isErr := catalogueNewLane(path)
	if !isErr || !strings.Contains(msg, "built-in") {
		t.Errorf("message = %q (err %v), want a refusal naming the built-in", msg, isErr)
	}
	if _, err := os.Stat(filepath.Join(catalogue, "review.md")); !os.IsNotExist(err) {
		t.Errorf("catalogue holds review.md (stat err %v), want none", err)
	}
}

// TestEditOfCataloguedLaneWritesThrough: after n put a lane on the board and
// into the catalogue, an E edit of the board file reaches the catalogue too,
// so the lane screen reports no drift and R has nothing old to pull back.
func TestEditOfCataloguedLaneWritesThrough(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	boardDir := lane.ProjectLanesDir(s.Root)
	writeLaneFile(t, boardDir, "new-lane.md", deployLane)
	if _, isErr := catalogueNewLane(filepath.Join(boardDir, "new-lane.md")); isErr {
		t.Fatal("setup: catalogueNewLane failed")
	}

	edited := filepath.Join(boardDir, "new-lane.md")
	syncTo, syncFrom := catalogueTwin(edited, "deploy-check")
	if syncTo != filepath.Join(catalogue, "deploy-check.md") {
		t.Fatalf("catalogueTwin = %q, want the catalogue copy", syncTo)
	}
	changed := strings.Replace(deployLane, "Check the deploy.", "Check the deploy twice.", 1)
	writeLaneFile(t, boardDir, "new-lane.md", changed)

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.laneScreen = newLaneScreen(m.store, m.lanes)
	m.Update(newLaneDoneMsg{edited: edited, syncTo: syncTo, syncFrom: syncFrom})

	b, _ := os.ReadFile(syncTo)
	if string(b) != changed {
		t.Errorf("catalogue copy = %q, want the edit", b)
	}
	if _, drifted := m.laneScreen.drift["deploy-check"]; drifted {
		t.Error("lane screen reports drift after the write-through")
	}
}

// TestEditDoesNotWriteThroughAnAlreadyDifferentCatalogue: a board copy that
// had already parted from its catalogue copy is edited on the board alone.
func TestEditDoesNotWriteThroughAnAlreadyDifferentCatalogue(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	boardDir := lane.ProjectLanesDir(s.Root)
	writeLaneFile(t, boardDir, "deploy-check.md", deployLane)
	other := strings.Replace(deployLane, "Check the deploy.", "Another board's version.", 1)
	writeLaneFile(t, catalogue, "deploy-check.md", other)
	if to, _ := catalogueTwin(filepath.Join(boardDir, "deploy-check.md"), "deploy-check"); to != "" {
		t.Errorf("catalogueTwin = %q, want none for copies that already differ", to)
	}
}

// TestWriteThroughLeavesACatalogueChangedMeanwhile: when the catalogue copy
// changed while the editor was open, it is not overwritten.
func TestWriteThroughLeavesACatalogueChangedMeanwhile(t *testing.T) {
	catalogue := catalogueHome(t)
	board := t.TempDir()
	writeLaneFile(t, board, "deploy-check.md", strings.Replace(deployLane, "deploy.", "deploy, edited.", 1))
	meanwhile := strings.Replace(deployLane, "deploy.", "deploy, from elsewhere.", 1)
	writeLaneFile(t, catalogue, "deploy-check.md", meanwhile)

	msg, isErr := writeThrough(filepath.Join(board, "deploy-check.md"), filepath.Join(catalogue, "deploy-check.md"), []byte(deployLane))
	if !isErr || !strings.Contains(msg, "meanwhile") {
		t.Errorf("message = %q (err %v), want it to say the catalogue changed meanwhile", msg, isErr)
	}
	b, _ := os.ReadFile(filepath.Join(catalogue, "deploy-check.md"))
	if string(b) != meanwhile {
		t.Error("catalogue copy changed meanwhile was overwritten")
	}
}

// TestWriteThroughKeepsCatalogueWhenIDChanged: an E edit that renames the
// lane's id does not put a lane of another id into the old catalogue file.
func TestWriteThroughKeepsCatalogueWhenIDChanged(t *testing.T) {
	catalogue := catalogueHome(t)
	board := t.TempDir()
	writeLaneFile(t, catalogue, "deploy-check.md", deployLane)
	writeLaneFile(t, board, "deploy-check.md", strings.Replace(deployLane, "id: deploy-check", "id: deploy-gate", 1))

	msg, isErr := writeThrough(filepath.Join(board, "deploy-check.md"), filepath.Join(catalogue, "deploy-check.md"), []byte(deployLane))
	if !isErr || !strings.Contains(msg, "deploy-gate") {
		t.Errorf("message = %q (err %v), want it to name the new id", msg, isErr)
	}
	b, _ := os.ReadFile(filepath.Join(catalogue, "deploy-check.md"))
	if string(b) != deployLane {
		t.Error("catalogue file of the old id was overwritten")
	}
}

// TestEditOfBuiltinIDDoesNotWriteThrough: a board lane carrying a built-in's
// id whose catalogue override is byte-equal is edited on the board alone —
// the override is what every board without its own lane directory runs.
func TestEditOfBuiltinIDDoesNotWriteThrough(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	boardDir := lane.ProjectLanesDir(s.Root)
	review := strings.Replace(deployLane, "id: deploy-check", "id: review", 1)
	writeLaneFile(t, boardDir, "review.md", review)
	writeLaneFile(t, catalogue, "review.md", review)
	if to, _ := catalogueTwin(filepath.Join(boardDir, "review.md"), "review"); to != "" {
		t.Errorf("catalogueTwin = %q, want none for a built-in id", to)
	}
}

// TestNewLaneEditedSkeletonWithPlaceholderIDStaysOutOfCatalogue: a skeleton
// the editor changed (here only whitespace) but whose id is still the
// template's placeholder is not a named lane, and must not reach the
// catalogue, where every board would offer it under '+'.
func TestNewLaneEditedSkeletonWithPlaceholderIDStaysOutOfCatalogue(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	path, err := writeLaneSkeleton(lane.ProjectLanesDir(s.Root))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(lane.Template+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	msg, isErr := catalogueNewLane(path)
	if !isErr || !strings.Contains(msg, "placeholder") {
		t.Errorf("message = %q (err %v), want an error naming the placeholder id", msg, isErr)
	}
	if matches, _ := filepath.Glob(filepath.Join(catalogue, "*.md")); len(matches) != 0 {
		t.Errorf("catalogue = %v, want empty", matches)
	}
}

// TestEditorErrorSaysCatalogueWasSkipped: an editor that exits non-zero
// copies nothing, and the lane screen says what that means for the lane —
// for n, that it stayed on this board; for E, that the catalogue keeps the
// old copy.
func TestEditorErrorSaysCatalogueWasSkipped(t *testing.T) {
	catalogue := catalogueHome(t)
	s := initBoard(t)
	boardDir := lane.ProjectLanesDir(s.Root)
	writeLaneFile(t, boardDir, "new-lane.md", deployLane)
	path := filepath.Join(boardDir, "new-lane.md")

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.laneScreen = newLaneScreen(m.store, m.lanes)
	m.Update(newLaneDoneMsg{err: errors.New("exit status 1"), path: path})
	if !m.laneScreen.isErr || !strings.Contains(m.laneScreen.msg, "exit status 1") || !strings.Contains(m.laneScreen.msg, "this board only") {
		t.Errorf("n: message = %q (err %v), want the editor error and 'this board only'", m.laneScreen.msg, m.laneScreen.isErr)
	}
	if _, err := os.Stat(filepath.Join(catalogue, "deploy-check.md")); !os.IsNotExist(err) {
		t.Errorf("catalogue holds deploy-check.md after an editor error (stat err %v)", err)
	}

	syncTo := filepath.Join(catalogue, "deploy-check.md")
	m.Update(newLaneDoneMsg{err: errors.New("exit status 1"), edited: path, syncTo: syncTo, syncFrom: []byte(deployLane)})
	if !m.laneScreen.isErr || !strings.Contains(m.laneScreen.msg, "keeps the old copy") {
		t.Errorf("E: message = %q (err %v), want 'keeps the old copy'", m.laneScreen.msg, m.laneScreen.isErr)
	}
}
