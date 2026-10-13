package tui

import (
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Not parallel: ticket.Counts is process-wide.

func titleOf(m *Model, id string) string {
	for _, t := range m.tickets {
		if t.ID == id {
			return t.Title
		}
	}
	return ""
}

// A tick on a board where nothing changed decodes no ticket and parses no
// frontmatter; a change made behind the board's back shows on the next tick,
// and costs exactly one decode.
func TestIdleTickDecodesNothingAndChangeAppearsNextTick(t *testing.T) {
	m := newTestModel(t, 150, 32)
	defer m.Close()

	dec0, par0 := ticket.Counts()
	for i := 0; i < 3; i++ {
		m.Update(tickMsg(time.Now()))
	}
	dec1, par1 := ticket.Counts()
	if dec1 != dec0 || par1 != par0 {
		t.Fatalf("idle ticks did work: %d decodes, %d header parses", dec1-dec0, par1-par0)
	}

	id := m.tickets[0].ID
	if _, err := m.store.Mutate(id, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldTitle, "changed behind the board")
	}); err != nil {
		t.Fatal(err)
	}
	dec2, _ := ticket.Counts()
	m.Update(tickMsg(time.Now()))
	dec3, _ := ticket.Counts()
	if got := titleOf(m, id); got != "changed behind the board" {
		t.Fatalf("the change did not show on the next tick: title %q", got)
	}
	if dec3-dec2 != 1 {
		t.Fatalf("one changed ticket should cost one decode, got %d", dec3-dec2)
	}
}

// agedTicket pushes every ticket file an hour into the past, outside the racy
// window, and lets the board cache them.
func agedBoard(t *testing.T) *Model {
	t.Helper()
	m := newTestModel(t, 150, 32)
	old := time.Now().Add(-time.Hour)
	for _, tk := range m.tickets {
		if err := os.Chtimes(tk.Path, old, old); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	return m
}

// rewriteKeepingSizeAndTime changes a ticket's title without changing its file's
// size or modification time: nothing but the watcher can tell it happened.
func rewriteKeepingSizeAndTime(t *testing.T, path, from, to string) {
	t.Helper()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	out := strings.Replace(string(raw), from, to, 1)
	if len(out) != len(raw) || out == string(raw) {
		t.Fatalf("bad rewrite %q -> %q", from, to)
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, fi.ModTime(), fi.ModTime()); err != nil {
		t.Fatal(err)
	}
}

// A change message naming one file reads that file and nothing else, even when
// the file kept its size and modification time.
func TestChangeMsgRereadsOnlyTheNamedFile(t *testing.T) {
	m := agedBoard(t)
	defer m.Close()
	tk := m.tickets[0]
	rewriteKeepingSizeAndTime(t, tk.Path, tk.Title, strings.ToUpper(tk.Title))

	m.Update(tickMsg(time.Now()))
	if got := titleOf(m, tk.ID); got != tk.Title {
		t.Fatalf("a tick alone saw an invisible change: %q", got)
	}

	dec0, par0 := ticket.Counts()
	m.Update(changeMsg{paths: []string{tk.Path}})
	dec1, par1 := ticket.Counts()
	if got := titleOf(m, tk.ID); got != strings.ToUpper(tk.Title) {
		t.Fatalf("the named file was not reread: %q", got)
	}
	if dec1-dec0 != 1 || par1-par0 != 1 {
		t.Fatalf("want 1 decode and 1 parse, got %d and %d", dec1-dec0, par1-par0)
	}
}

// The real watcher reports the path of a file that changed.
func TestWatcherReportsChangedPath(t *testing.T) {
	m := newTestModel(t, 150, 32)
	defer m.Close()
	m.startWatching()
	if m.closer == nil {
		t.Skip("no file watcher on this platform")
	}
	path := m.tickets[0].Path
	raw, _ := os.ReadFile(path)
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	select {
	case paths := <-m.watch:
		if !slices.Contains(paths, path) {
			t.Fatalf("watcher named %v, want it to include %s", paths, path)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no change event within 5s")
	}
}
