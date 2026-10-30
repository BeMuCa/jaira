package lane

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// activeRoot is a repository with today's board and two named ones.
func activeRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range []string{"video", "docs"} {
		dir := filepath.Join(root, ticket.DirName, "boards", name, "lanes")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("---\nid: a\nname: A\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// Without the file every board is active, today's first; switching one off
// writes the rest by name, and the last one cannot be switched off.
func TestActiveBoards(t *testing.T) {
	root := activeRoot(t)
	if got, want := ActiveBoards(root), []string{"", "docs", "video"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("no file: ActiveBoards = %q, want %q", got, want)
	}

	if err := SetBoardActive(root, "docs", false); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(ActiveBoardsPath(root))
	if string(b) != "default\nvideo\n" {
		t.Errorf("file = %q, want default and video, one per line", b)
	}
	if got, want := ActiveBoards(root), []string{"", "video"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ActiveBoards = %q, want %q", got, want)
	}

	if err := SetBoardActive(root, "default", false); err != nil {
		t.Fatal(err)
	}
	if err := SetBoardActive(root, "video", false); !errors.Is(err, ErrLastActiveBoard) {
		t.Errorf("switching off the last active board: err = %v, want ErrLastActiveBoard", err)
	}
	if got, want := ActiveBoards(root), []string{"video"}; !reflect.DeepEqual(got, want) {
		t.Errorf("after the refusal ActiveBoards = %q, want %q", got, want)
	}

	// A file naming only boards that are gone falls back to all of them.
	if err := os.WriteFile(ActiveBoardsPath(root), []byte("gone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ActiveBoards(root); len(got) != 3 {
		t.Errorf("a file naming no existing board: ActiveBoards = %q, want all three", got)
	}
}
