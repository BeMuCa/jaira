package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A failed ref listing must not turn into "every ticket left the board".
// fetch and validate both report departures through departedHere: it says the
// check could not run and names nothing, and leaves the rest of the command's
// report alone.
//
// The listing is broken after the syncer has found the board usable, because
// that is the case that matters: a broken packed-refs from the start makes the
// whole board read as having no remote, and Departed is never asked.
func TestAFailedRefListingReportsNoDepartures(t *testing.T) {
	ada, grace := twoBoards(t)
	id := createTicket(t, ada, "still on the board")
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "pull", ticket.Handle(id)); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	// refs is the syncer grace's last command opened; it has seen the ticket.
	if err := refs.Usable(); err != nil {
		t.Fatalf("grace's board does not use refs: %v", err)
	}
	var w strings.Builder
	if d := departedHere(&w); len(d) != 0 || w.Len() != 0 {
		t.Fatalf("nothing has left the board yet: %+v, %q", d, w.String())
	}

	packed := filepath.Join(grace, ".git", "packed-refs")
	before, _ := os.ReadFile(packed) // absent is fine
	if err := os.WriteFile(packed, append(append([]byte{}, before...), []byte("not a ref line\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	w.Reset()
	if d := departedHere(&w); len(d) != 0 {
		t.Errorf("a failed ref listing reported departures: %+v", d)
	}
	if !strings.Contains(w.String(), "could not check for tickets taken off the board elsewhere") {
		t.Errorf("the failed listing was not said: %q", w.String())
	}
}
