package lane

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const scriptLane = "---\nid: script\nname: Script\nagentic: true\nmodel-tier: strong\n---\n\n# Prompt\n\nWrite the script.\n"

// newCatalogueHome gives the test a JAIRA_HOME of its own holding one
// catalogue lane, script, that no repository has.
func newCatalogueHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	t.Setenv("JAIRA_LANES_DIR", "")
	if err := os.MkdirAll(filepath.Join(home, "lanes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "lanes", "script.md"), []byte(scriptLane), 0o644); err != nil {
		t.Fatal(err)
	}
	return home
}

func typeNamed(t *testing.T, name string) (BoardType, bool) {
	t.Helper()
	types, err := BoardTypes()
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(types, func(bt BoardType) bool { return bt.Name == name })
	if i < 0 {
		return BoardType{}, false
	}
	return types[i], true
}

// A board built from picked lanes — one from today's board, one from the
// catalogue — is saved into the catalogue once, and another repository then
// builds the same board from it as a type.
func TestBoardInTheCatalogueIsATypeInEveryRepository(t *testing.T) {
	newCatalogueHome(t)
	a := newBoardRoot(t)
	set, err := Load(a)
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{"script", "backlog", "done"}
	if _, err := CreateBoardWithLanes(a, "video", set, ids, "Video work"); err != nil {
		t.Fatal(err)
	}
	vset, err := LoadBoard(a, "video")
	if err != nil {
		t.Fatal(err)
	}
	if got := vset.IDs(); !slices.Equal(got, ids) {
		t.Errorf("video lanes = %v, want %v", got, ids)
	}

	path, saved, err := SaveToCatalogue(a, "video")
	if err != nil || !saved || path != filepath.Join(UserBoardsDir(), "video") {
		t.Fatalf("SaveToCatalogue = %q, %v, %v; want the catalogue's video, saved", path, saved, err)
	}
	if _, saved, _ := SaveToCatalogue(a, "video"); saved {
		t.Error("a second save replaced the catalogue board")
	}

	bt, ok := typeNamed(t, "video")
	if !ok || !bt.Catalogue || bt.Description != "Video work" || !slices.Equal(bt.Lanes, ids) {
		t.Fatalf("catalogue type video = %+v (found %v)", bt, ok)
	}
	b := newBoardRoot(t)
	if _, got, err := CreateBoard(b, "clips", "video", "Clips"); err != nil || !slices.Equal(got, ids) {
		t.Fatalf("create from the catalogue type in another repo: %v, %v", got, err)
	}
	raw, _ := os.ReadFile(filepath.Join(boardLanesDir(b, "clips"), "script.md"))
	if string(raw) != scriptLane {
		t.Errorf("catalogue lane did not arrive verbatim:\n%s", raw)
	}

	// A board named like a shipped type stays out of the catalogue, so the
	// shipped type is never hidden behind a copy.
	if _, _, err := CreateBoard(a, DefaultBoardType, DefaultBoardType, "x"); err != nil {
		t.Fatal(err)
	}
	if path, saved, err := SaveToCatalogue(a, DefaultBoardType); saved || path != "" || err != nil {
		t.Errorf("save of a board named like a shipped type: %q, %v, %v", path, saved, err)
	}
}

func TestCreateBoardWithLanesRefusesUnknownAndRepeatedLanes(t *testing.T) {
	newCatalogueHome(t)
	root := newBoardRoot(t)
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, ids := range [][]string{nil, {"backlog", "nope"}, {"backlog", "backlog"}} {
		if _, err := CreateBoardWithLanes(root, "video", set, ids, "x"); err == nil {
			t.Errorf("lanes %v were accepted", ids)
		}
	}
	if Boards(root) != nil {
		t.Errorf("a refused create left a board behind: %v", Boards(root))
	}
}

// Publish writes the board under shared/<who>/boards/, the listing finds
// it, and adopt puts it into the catalogue — neither overwrites without
// being told to.
func TestPublishAndAdoptBoard(t *testing.T) {
	newCatalogueHome(t)
	root := newBoardRoot(t)
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateBoardWithLanes(root, "video", set, []string{"script", "done"}, "Video work"); err != nil {
		t.Fatal(err)
	}
	shared := filepath.Join(root, ".jaira", "shared")
	dst, err := PublishBoard(root, "video", shared, "berk", false)
	if err != nil || dst != filepath.Join(shared, "berk", "boards", "video") {
		t.Fatalf("publish: %q, %v", dst, err)
	}
	if _, err := PublishBoard(root, "video", shared, "berk", false); !errors.Is(err, ErrBoardExists) {
		t.Errorf("second publish: %v, want ErrBoardExists", err)
	}
	if _, err := PublishBoard(root, "video", shared, "berk", true); err != nil {
		t.Errorf("publish with overwrite: %v", err)
	}
	if dst, err := PublishBoard(root, "", shared, "berk", false); err != nil || filepath.Base(dst) != DefaultBoardName {
		t.Errorf("publish today's board: %q, %v; want it named default", dst, err)
	}
	// One broken board is named, not fatal.
	if err := os.MkdirAll(filepath.Join(shared, "ana", "boards", "broken"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Shared lanes stay lanes only: a published board's lanes are not read as
	// shared lanes.
	if sl, _, _ := Shared(root); len(sl) != 0 {
		t.Errorf("a published board's lanes showed up as shared lanes: %v", sl)
	}

	boards, warnings, err := SharedBoards(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(boards) != 2 || boards[1].Name != "video" || boards[1].Folder != "berk" || boards[1].Description != "Video work" {
		t.Fatalf("shared boards = %+v", boards)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "broken") {
		t.Errorf("warnings = %v, want the broken board named", warnings)
	}

	newCatalogueHome(t) // a teammate's machine: nothing in the catalogue yet
	bt, path, err := AdoptBoard(boards[1].Path, false)
	if err != nil || path != filepath.Join(UserBoardsDir(), "video") || !slices.Equal(bt.Lanes, []string{"script", "done"}) {
		t.Fatalf("adopt: %+v, %q, %v", bt, path, err)
	}
	if got, ok := typeNamed(t, "video"); !ok || !got.Catalogue {
		t.Errorf("adopted board is not a catalogue type: %+v", got)
	}
	if _, _, err := AdoptBoard(boards[1].Path, false); !errors.Is(err, ErrBoardExists) {
		t.Errorf("second adopt: %v, want ErrBoardExists", err)
	}
	if _, _, err := AdoptBoard(boards[1].Path, true); err != nil {
		t.Errorf("adopt with overwrite: %v", err)
	}
}
