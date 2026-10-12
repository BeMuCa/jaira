package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

// staleOnDisk sets up the K1W2JR accident: a ticket pulled here, moved on
// since on its ref, and a file on disk that still holds the earlier state — as
// a branch does that missed the lanes which commit nothing. It returns the
// handle, the file's path and the ref's SHA once the ref is ahead.
func staleOnDisk(t *testing.T) (clone, h, path, sha string) {
	t.Helper()
	t.Setenv("JAIRA_USER", "ada")
	clone = cloneWithRemotes(t, "origin")
	h = ticketIn(t, clone, "stale file")
	if out, err := runCLI(t, clone, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(h)
	if err != nil {
		t.Fatal(err)
	}
	path = tk.Path
	old, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, clone, "move", h, "--to", "todo", "--force"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "note", h, "the review found nothing"); err != nil {
		t.Fatalf("note: %v\n%s", err, out)
	}
	// What runs after every command in production; runCLI stops short of it.
	flushRefs()
	// The earlier bytes go back, stamped well before the ref's last write so
	// the test does not hinge on two writes landing in different seconds.
	var lines []string
	for _, line := range strings.Split(string(old), "\n") {
		if strings.HasPrefix(line, "updated-at:") {
			line = "updated-at: 2026-09-10T08:00:00Z"
		}
		lines = append(lines, line)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	sha, err = (&gitref.Repo{Dir: clone}).SHA(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	return clone, h, path, sha
}

// The accident itself: move read the stale file, --force overrode the gate's
// refusals that were all symptoms of it, and the old state went onto the ref.
// Now the move is refused before the gate, with both states named, and the ref
// is where it was.
func TestMoveRefusesAFileOlderThanItsRefEvenWithForce(t *testing.T) {
	clone, h, path, sha := staleOnDisk(t)
	before, _ := os.ReadFile(path)

	out, err := runCLI(t, clone, "move", h, "--to", "done", "--force")
	if !errors.Is(err, ticket.ErrStale) {
		t.Fatalf("move --force on a stale file was not refused as stale: %v\n%s", err, out)
	}
	msg := err.Error()
	for _, want := range []string{"file: backlog, updated 2026-09-10 08:00 UTC", "ref:  todo", "jaira pull " + h} {
		if !strings.Contains(msg, want) {
			t.Errorf("the refusal does not say %q:\n%s", want, msg)
		}
	}
	id := ticket.IDFromFilename(filepath.Base(path))
	if got, _ := (&gitref.Repo{Dir: clone}).SHA(id); got != sha {
		t.Errorf("the ref moved: %s, was %s", got, sha)
	}
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Error("the refused move changed the file")
	}
}

// Every write reads the file and puts it on the ref, so a note is refused the
// same way — the move is only where it was noticed.
func TestANoteOnAStaleFileIsRefusedToo(t *testing.T) {
	clone, h, _, _ := staleOnDisk(t)
	if _, err := runCLI(t, clone, "note", h, "lost"); !errors.Is(err, ticket.ErrStale) {
		t.Fatalf("a note on a stale file was not refused as stale: %v", err)
	}
}

// pull was a no-op as soon as any file existed, so it could not repair the
// state the refusal points at. It brings the file up to the ref now, and the
// move that was refused goes through afterwards.
func TestPullBringsAStaleFileUpToItsRef(t *testing.T) {
	clone, h, path, sha := staleOnDisk(t)

	out, err := runCLI(t, clone, "pull", h)
	if err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	for _, want := range []string{"was behind its ref", "was: backlog", "now: todo"} {
		if !strings.Contains(out, want) {
			t.Errorf("pull does not say %q:\n%s", want, out)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ref, _, err := (&gitref.Repo{Dir: clone}).Read(ticket.IDFromFilename(filepath.Base(path)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(ref) {
		t.Errorf("the file is not the ref's bytes after pull:\n%s\n--- ref ---\n%s", got, ref)
	}
	if now, _ := (&gitref.Repo{Dir: clone}).SHA(ticket.IDFromFilename(filepath.Base(path))); now != sha {
		t.Error("bringing a file up to its ref wrote the ref")
	}
	if out, err := runCLI(t, clone, "note", h, "after the pull"); err != nil {
		t.Fatalf("a note after the pull was refused: %v\n%s", err, out)
	}
}
