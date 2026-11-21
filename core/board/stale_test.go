package board

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// A jaira before 689E0T wrote "/sub/.jaira/" for a board in " sub". The
// repair takes out that line and only that line: never one a person wrote,
// never one that belongs to a board in the directory it names, never one
// for a path that starts with no space.
func TestDropStaleEntryTakesOutOnlyItsOwnLine(t *testing.T) {
	own := func(nl, line string) string { return nl + privateComment + nl + line + nl }
	cases := []struct {
		name     string
		dir      string // the board, below the top
		other    string // a directory holding a board of its own, "" for none
		existing func(nl string) string
		crlf     bool
		dropped  string
	}{
		{name: "own line", dir: " sub", existing: func(nl string) string { return "*.log" + nl + own(nl, "/sub/.jaira/") }, dropped: "/sub/.jaira/"},
		{name: "own line, CRLF", dir: " sub", crlf: true, existing: func(nl string) string { return "*.log" + nl + own(nl, "/sub/.jaira/") }, dropped: "/sub/.jaira/"},
		{name: "own line, nested", dir: " a/b", existing: func(nl string) string { return own(nl, "/a/b/.jaira/") }, dropped: "/a/b/.jaira/"},
		{name: "own line, glob characters", dir: "  x[1]", existing: func(nl string) string { return own(nl, `/x\[1]/.jaira/`) }, dropped: `/x\[1]/.jaira/`},
		{name: "a person's line", dir: " sub", existing: func(nl string) string { return "/sub/.jaira/" + nl }},
		{name: "own and a person's copy", dir: " sub", existing: func(nl string) string { return "/sub/.jaira/" + nl + own(nl, "/sub/.jaira/") }},
		{name: "a board stands in sub", dir: " sub", other: "sub", existing: func(nl string) string { return own(nl, "/sub/.jaira/") }},
		{name: "no leading space", dir: "sub", existing: func(nl string) string { return own(nl, "/sub/.jaira/") }},
		{name: "no stale line", dir: " sub", existing: func(nl string) string { return "*.log" + nl }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root, exclude := newRepo(t)
			board := filepath.Join(root, filepath.FromSlash(c.dir))
			if err := os.MkdirAll(board, 0o755); err != nil {
				t.Fatal(err)
			}
			if c.other != "" {
				if err := os.MkdirAll(filepath.Join(root, c.other, ".jaira"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			nl := "\n"
			if c.crlf {
				nl = "\r\n"
			}
			write(t, exclude, c.existing(nl))
			if _, err := AddIgnore(board); err != nil {
				t.Fatal(err)
			}
			before, _ := os.ReadFile(exclude)

			_, stale := StaleEntry(board)
			dropped, err := DropStaleEntry(board)
			if err != nil {
				t.Fatal(err)
			}
			if dropped != c.dropped || stale != c.dropped {
				t.Fatalf("StaleEntry = %q, DropStaleEntry = %q; want %q", stale, dropped, c.dropped)
			}
			after, _ := os.ReadFile(exclude)
			if c.dropped == "" {
				if string(after) != string(before) {
					t.Errorf("exclude changed:\n before %q\n after  %q", before, after)
				}
				return
			}
			if count(string(after), c.dropped) != 0 || strings.Count(string(after), privateMarker) != 1 {
				t.Errorf("the stale line or its comment is still there:\n%q", after)
			}
			if !Ignored(board) {
				t.Errorf("the board's own entry went with the stale one:\n%q", after)
			}
			if c.crlf && strings.Contains(strings.ReplaceAll(string(after), "\r\n", ""), "\n") {
				t.Errorf("CRLF file got an LF line:\n%q", after)
			}
			if again, _ := DropStaleEntry(board); again != "" {
				t.Errorf("a second DropStaleEntry dropped %q", again)
			}
		})
	}
}

// Two jaira update runs repairing the same board at once take the line out
// once and leave the board's own entry standing.
func TestParallelDropStaleEntry(t *testing.T) {
	root, exclude := newRepo(t)
	board := filepath.Join(root, " sub")
	if err := os.MkdirAll(board, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, exclude, "\n"+privateComment+"\n/sub/.jaira/\n")
	if _, err := AddIgnore(board); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	n := 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, err := DropStaleEntry(board)
			if err != nil {
				t.Error(err)
			}
			if d != "" {
				mu.Lock()
				n++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	got, _ := os.ReadFile(exclude)
	if n != 1 || count(string(got), "/sub/.jaira/") != 0 || !Ignored(board) {
		t.Errorf("dropped %d times; exclude:\n%s", n, got)
	}
}
