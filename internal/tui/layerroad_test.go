package tui

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/road"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newRemoteModel is newTestModel on a board with a remote, so a road can be
// written to its ref, and read back from the outbox.
func newRemoteModel(t *testing.T, w, h int) *Model {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(root, "no-lanes"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "ada")
	run(t, root, "git", "clone", "--quiet", bare, clone)
	run(t, clone, "git", "config", "user.name", "ada")
	run(t, clone, "git", "config", "user.email", "ada@example.test")
	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = w, h
	return m
}

// Lanes no road names are Road-1, Road-2 and so on, left to right.
func TestUnnamedSwimlanesAreNumbered(t *testing.T) {
	m := newTestModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())
	if !strings.Contains(out, "Road-1") || !strings.Contains(out, "Road-2") {
		t.Errorf("the headers must read Road-1 and Road-2:\n%s", out)
	}
}

// A lane whose road exists shows the road's title, and the numbering counts
// only the lanes no road names.
func TestARoadsTitleHeadsItsLane(t *testing.T) {
	m := newRemoteModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, nil)
	mkPlaced(t, m, "c", 2, nil)
	reloadOrFail(t, m)
	if _, err := road.Create(m.store, m.refSync, "auth", "Auth rollout", "", []string{lid("b")}); err != nil {
		t.Fatal(err)
	}
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())
	// Lanes a, b (named), c: a is Road-1 and c, the second unnamed one, Road-2.
	for _, want := range []string{"Auth rollout", "Road-1", "Road-2"} {
		if !strings.Contains(out, want) {
			t.Errorf("the headers must show %q, the named lane by its title:\n%s", want, out)
		}
	}
}

// Enter on an unnamed lane's header names it: a road is written holding every
// milestone drawn in that lane.
func TestEnterNamesAnUnnamedLaneWithAllItsMilestones(t *testing.T) {
	m := newRemoteModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyUp)
	m.key(key("enter"))
	if m.layer.stage != stageRoadName {
		t.Fatalf("enter on the header must open the name field, stage %v", m.layer.stage)
	}
	for range len("Road-1") {
		m.key(key("backspace"))
	}
	typeInto(m, "Auth rollout")
	m.key(key("enter"))

	roads, err := road.LoadAll(m.refSync)
	if err != nil || len(roads) != 1 {
		t.Fatalf("roads = %v, %v; want one road", roads, err)
	}
	r := roads[0]
	if r.Title != "Auth rollout" || r.Name != "auth-rollout" || !r.Has(lid("a")) || !r.Has(lid("b")) {
		t.Errorf("road = %q (%s) with %v; want Auth rollout with a and b", r.Title, r.Name, r.Members())
	}
	if m.layer.stage != stageLayer || m.layer.sel.kind != onRoad {
		t.Errorf("after saving, stage %v sel %+v; want the picture with the header still selected", m.layer.stage, m.layer.sel)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "Auth rollout") {
		t.Errorf("the header must show the new name:\n%s", out)
	}
}

// Enter on a header whose lane has a road renames the road, and keeps its
// milestones.
func TestEnterRenamesALaneWithARoad(t *testing.T) {
	m := newRemoteModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	reloadOrFail(t, m)
	if _, err := road.Create(m.store, m.refSync, "auth", "Auth", "", []string{lid("a")}); err != nil {
		t.Fatal(err)
	}
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyUp)
	m.key(key("enter"))
	if got := m.layer.form.vals[0]; got != "Auth" {
		t.Errorf("the field starts as the road's title, got %q", got)
	}
	for range len("Auth") {
		m.key(key("backspace"))
	}
	typeInto(m, "Login")
	m.key(key("enter"))

	roads, _ := road.LoadAll(m.refSync)
	if len(roads) != 1 || roads[0].Title != "Login" || roads[0].Name != "auth" || !roads[0].Has(lid("a")) {
		t.Errorf("roads = %+v; want auth titled Login, still listing a", roads)
	}
}

// Esc writes nothing: the field closes and the road is as it was.
func TestEscInTheRoadFieldWritesNothing(t *testing.T) {
	m := newRemoteModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyUp)
	m.key(key("enter"))
	typeInto(m, "Never saved")
	m.key(key("esc"))
	if m.layer.stage != stageLayer {
		t.Errorf("esc must close the field, stage %v", m.layer.stage)
	}
	if roads, _ := road.LoadAll(m.refSync); len(roads) != 0 {
		t.Errorf("esc wrote a road: %+v", roads)
	}
}

// A blank name is refused where it is typed, and nothing is written.
func TestABlankRoadNameIsRefused(t *testing.T) {
	m := newRemoteModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyUp)
	m.key(key("enter"))
	for range len("Road-1") {
		m.key(key("backspace"))
	}
	m.key(key("enter"))
	if m.layer.stage != stageRoadName || m.layer.err == "" {
		t.Errorf("a blank name must stay in the field with a reason, stage %v err %q", m.layer.stage, m.layer.err)
	}
	if roads, _ := road.LoadAll(m.refSync); len(roads) != 0 {
		t.Errorf("a blank name wrote a road: %+v", roads)
	}
}
