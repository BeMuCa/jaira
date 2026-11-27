package ticket_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

func mkTickets(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ticket.DirName, ticket.TicketsSubdir), 0o755); err != nil {
		t.Fatal(err)
	}
}

func realPath(t *testing.T, p string) string {
	t.Helper()
	r, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// TestDiscoverInLinkedWorktree covers which board a command run inside a
// linked worktree opens, one case per shape the worktree and its main checkout
// can have.
func TestDiscoverInLinkedWorktree(t *testing.T) {
	cases := []struct {
		name      string
		mainBoard string // where the main checkout's board sits, relative; "-" none
		treeBoard string // where the worktree's own board sits, relative; "-" none
		from      string // where in the worktree the command runs
		want      string // "main", "tree" or "none", plus the relative board root
		wantRel   string
	}{
		{"no board in the worktree: the main checkout's", "", "-", "", "main", ""},
		{"from a subdirectory of the worktree", "", "-", "a/b", "main", ""},
		{"board in a subdirectory of the main checkout", "sub", "-", "sub/deeper", "main", "sub"},
		{"the worktree's own board wins", "", "", "", "tree", ""},
		{"the worktree's own board wins from below", "", "", "a", "tree", ""},
		{"neither has a board", "-", "-", "", "none", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			main, tree := clonePair(t)
			if c.mainBoard != "-" {
				mkTickets(t, filepath.Join(main, c.mainBoard))
			}
			if c.treeBoard != "-" {
				mkTickets(t, filepath.Join(tree, c.treeBoard))
			}
			from := filepath.Join(tree, c.from)
			if err := os.MkdirAll(from, 0o755); err != nil {
				t.Fatal(err)
			}
			s, err := ticket.Discover(from)
			if c.want == "none" {
				if !errors.Is(err, ticket.ErrNoStore) {
					t.Fatalf("got %v, %v; want ErrNoStore", s, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			base := main
			if c.want == "tree" {
				base = tree
			}
			if got, want := realPath(t, s.Root), realPath(t, filepath.Join(base, c.wantRel)); got != want {
				t.Fatalf("Root = %q, want %q", got, want)
			}
		})
	}
}
