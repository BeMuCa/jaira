package tui

import (
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
