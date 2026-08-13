package lane

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadDefaultBoardAbsentReturnsZeroBoard asserts no file at all is the
// normal state, not a missing-configuration error.
func TestLoadDefaultBoardAbsentReturnsZeroBoard(t *testing.T) {
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "does-not-exist.md"))
	b, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Lanes) != 0 || len(b.Options) != 0 {
		t.Errorf("expected a zero board, got %+v", b)
	}
	if len(b.Warnings) != 0 {
		t.Errorf("an absent file must not warn, got: %v", b.Warnings)
	}
}

// TestLoadDefaultBoardParsesSelection asserts lanes: and options: parse in
// order.
func TestLoadDefaultBoardParsesSelection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(path, []byte("---\nlanes: [backlog, todo, done]\noptions: [brainstorm]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", path)

	b, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	wantLanes := []string{"backlog", "todo", "done"}
	if len(b.Lanes) != len(wantLanes) {
		t.Fatalf("Lanes = %v, want %v", b.Lanes, wantLanes)
	}
	for i, id := range wantLanes {
		if b.Lanes[i] != id {
			t.Errorf("Lanes[%d] = %q, want %q", i, b.Lanes[i], id)
		}
	}
	if len(b.Options) != 1 || b.Options[0] != "brainstorm" {
		t.Errorf("Options = %v, want [brainstorm]", b.Options)
	}
}

// TestLoadDefaultBoardUnparseableWarnsAndIsAbsent asserts a file that exists
// but does not parse is treated the same as no file, with a warning.
func TestLoadDefaultBoardUnparseableWarnsAndIsAbsent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(path, []byte("not frontmatter at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", path)

	b, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Lanes) != 0 {
		t.Errorf("an unparseable file must be treated as absent, got Lanes=%v", b.Lanes)
	}
	if len(b.Warnings) == 0 {
		t.Error("an unparseable file must warn")
	}
}

func builtinIDList(t *testing.T) []string {
	t.Helper()
	lanes, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]string, len(lanes))
	for i, l := range lanes {
		ids[i] = l.ID
	}
	return ids
}

// TestDiffersFalseWhenSelectionMatchesBuiltins asserts naming exactly the
// built-in ids, with every one resolving to the built-in, is not a
// difference.
func TestDiffersFalseWhenSelectionMatchesBuiltins(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	board := &DefaultBoard{Lanes: builtinIDList(t)}
	if Differs(board, set) {
		t.Error("a selection matching the built-ins exactly must not differ")
	}
}

// TestDiffersFalseWhenBoardIsEmpty asserts an absent or empty selection means
// the built-ins.
func TestDiffersFalseWhenBoardIsEmpty(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if Differs(&DefaultBoard{}, set) {
		t.Error("an empty board selection must not differ")
	}
}

// TestDiffersTrueWhenALaneIsDropped asserts a shorter selection differs.
func TestDiffersTrueWhenALaneIsDropped(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	all := builtinIDList(t)
	board := &DefaultBoard{Lanes: all[:len(all)-1]}
	if !Differs(board, set) {
		t.Error("dropping a lane must differ from the built-ins")
	}
}

// TestDiffersTrueWhenSelectedLaneIsAnOverride asserts a selection naming
// exactly the built-in ids, but where one resolves to a catalogue override,
// still differs.
func TestDiffersTrueWhenSelectedLaneIsAnOverride(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "review.md", `---
id: review
name: My Review
after: human
precedence: 50
agentic: true
model-tier: strong
---
A different prompt.
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	board := &DefaultBoard{Lanes: builtinIDList(t)}
	if !Differs(board, set) {
		t.Error("a selection resolving one id to an override must differ")
	}
}

// TestMaterialiseWritesNothingWhenBoardMatchesBuiltins is the criterion's
// load-bearing half: a repo whose owner changed nothing carries no lane
// files at all.
func TestMaterialiseWritesNothingWhenBoardMatchesBuiltins(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	board := &DefaultBoard{Lanes: builtinIDList(t)}
	written, err := Materialise(root, set, board)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 0 {
		t.Errorf("expected nothing written, got %v", written)
	}
	if _, err := os.Stat(ProjectLanesDir(root)); !os.IsNotExist(err) {
		t.Errorf("expected %s to not exist, stat err = %v", ProjectLanesDir(root), err)
	}
}

// TestMaterialiseWritesOnlyTheSelection asserts a differing board writes
// exactly the selected lanes, each parsing back to the same lane.
func TestMaterialiseWritesOnlyTheSelection(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	board := &DefaultBoard{Lanes: []string{"backlog", "todo", "done"}}
	written, err := Materialise(root, set, board)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 3 {
		t.Fatalf("written = %v, want 3 files", written)
	}
	for _, name := range []string{"backlog", "todo", "done"} {
		p := filepath.Join(ProjectLanesDir(root), name+".md")
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
		l, err := parse(b, p, false)
		if err != nil {
			t.Fatalf("materialised %s did not parse: %v", p, err)
		}
		if l.ID != name {
			t.Errorf("materialised file %s parsed to id %q", p, l.ID)
		}
	}
	if _, err := os.Stat(filepath.Join(ProjectLanesDir(root), "review.md")); !os.IsNotExist(err) {
		t.Error("a lane not in the selection must not be materialised")
	}
}

// TestMaterialiseWarnsOnUnknownLaneAndContinues asserts an id the default
// board names that is not installed is a warning, and the rest still write.
func TestMaterialiseWarnsOnUnknownLaneAndContinues(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	board := &DefaultBoard{Lanes: []string{"backlog", "nosuchlane", "done"}}
	written, err := Materialise(root, set, board)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 2 {
		t.Fatalf("written = %v, want 2 files (the unknown lane skipped)", written)
	}
	found := false
	for _, w := range board.Warnings {
		if containsWarning([]string{w}, "nosuchlane") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a warning naming nosuchlane, got: %v", board.Warnings)
	}
}

// TestResolveOptionsTicksNamedOptions asserts the board's options are ticked
// against what the set actually installs, in display order.
func TestResolveOptionsTicksNamedOptions(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	board := &DefaultBoard{Options: []string{"brainstorm"}}
	opts := ResolveOptions(set, board)
	if len(opts) == 0 {
		t.Fatal("expected at least one option from the built-ins")
	}
	for _, o := range opts {
		if o.Name == "brainstorm" && !o.Ticked {
			t.Error("brainstorm should be ticked per the board")
		}
		if o.Name == "planning" && o.Ticked {
			t.Error("planning should not be ticked")
		}
	}
}

// TestSaveDefaultBoardRoundTrips asserts a saved board's lanes and options
// come back unchanged through LoadDefaultBoard.
func TestSaveDefaultBoardRoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default-board.md")
	t.Setenv("JAIRA_DEFAULT_BOARD", path)

	b := &DefaultBoard{Path: path, Lanes: []string{"backlog", "todo"}, Options: []string{"brainstorm"}}
	if err := SaveDefaultBoard(b); err != nil {
		t.Fatal(err)
	}

	got, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Lanes) != 2 || got.Lanes[0] != "backlog" || got.Lanes[1] != "todo" {
		t.Errorf("Lanes = %v, want [backlog todo]", got.Lanes)
	}
	if len(got.Options) != 1 || got.Options[0] != "brainstorm" {
		t.Errorf("Options = %v, want [brainstorm]", got.Options)
	}
}

// TestSaveDefaultBoardPreservesExistingBody asserts a user's own prose
// survives a save.
func TestSaveDefaultBoardPreservesExistingBody(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default-board.md")
	t.Setenv("JAIRA_DEFAULT_BOARD", path)

	if err := SaveDefaultBoard(&DefaultBoard{Path: path, Body: "# My notes\n\nDo not touch."}); err != nil {
		t.Fatal(err)
	}
	got, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if got.Body != "# My notes\n\nDo not touch." {
		t.Errorf("Body = %q, want the original prose preserved", got.Body)
	}

	// A second save with the loaded board (round-tripping Body) must not lose it.
	if err := SaveDefaultBoard(got); err != nil {
		t.Fatal(err)
	}
	got2, err := LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if got2.Body != "# My notes\n\nDo not touch." {
		t.Errorf("Body after second save = %q, want it preserved", got2.Body)
	}
}

// TestResolveOptionsDropsUnknownBoardOption asserts an option the board names
// that no lane requires does not appear (and does not panic).
func TestResolveOptionsDropsUnknownBoardOption(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	board := &DefaultBoard{Options: []string{"not-a-real-option"}}
	opts := ResolveOptions(set, board)
	for _, o := range opts {
		if o.Name == "not-a-real-option" {
			t.Error("an option no lane requires must not appear")
		}
	}
}
