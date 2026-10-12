package lane

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestJairaHomeHidesTheRealDefaultBoard: a machine whose ~/.jaira carries a
// default board and a lane catalogue of its own must not change what a board
// set up under a relocated JAIRA_HOME gets. Sixteen tests went red on one
// machine because its ~/.jaira/default-board.md left out backlog, while
// JAIRA_HOME — set by the test package — moved everything except these two.
func TestJairaHomeHidesTheRealDefaultBoard(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("JAIRA_DEFAULT_BOARD", "")
	t.Setenv("JAIRA_LANES_DIR", "")
	writeLane(t, filepath.Join(home, ".jaira"), "default-board.md",
		"---\nlanes: [todo, in-progress, done]\n---\n")
	writeLane(t, filepath.Join(home, ".jaira", "lanes"), "fallback.md",
		"---\nname: Fallback\nafter: human\nprecedence: 41\n---\n")

	jairaHome := t.TempDir()
	t.Setenv("JAIRA_HOME", jairaHome)
	if got, want := DefaultBoardPath(), filepath.Join(jairaHome, "default-board.md"); got != want {
		t.Errorf("DefaultBoardPath() = %q, want %q", got, want)
	}
	if got, want := UserLanesDir(), filepath.Join(jairaHome, "lanes"); got != want {
		t.Errorf("UserLanesDir() = %q, want %q", got, want)
	}

	root := t.TempDir()
	if err := os.MkdirAll(ProjectLanesDir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := set.IDs(); !reflect.DeepEqual(got, builtinIDList(t)) {
		t.Errorf("board = %v, want the built-ins — the real ~/.jaira leaked in", got)
	}
}
