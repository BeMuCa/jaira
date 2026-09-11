package ticket

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// One clone logs a finished ticket while another still has it on the board.
// Merging those branches keeps both paths, and the ticket then lives twice —
// closed in the logbook, open on the board. The old duplicate check could not
// see this, because it only compared ids inside tickets/.
func TestATicketOnTheBoardAndInTheLogbookIsReported(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	s, err := At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	id := NewID(time.Now())
	doc := "---\nid: " + id + "\ntitle: two copies\nstatus: done\n---\n\n# two copies\n"
	if err := WriteAtomic(filepath.Join(s.TicketsDir(), Filename(id, "two copies")), []byte(doc)); err != nil {
		t.Fatal(err)
	}

	// Clean board: no complaint yet.
	if _, err := s.List(); err != nil {
		t.Fatalf("a single copy was reported as a problem: %v", err)
	}

	// The same ticket, also filed away by somebody else's branch.
	folder := filepath.Join(s.LogbookDir(), "grace-20260910")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(filepath.Join(folder, Filename(id, "two copies")), []byte(doc)); err != nil {
		t.Fatal(err)
	}

	_, err = s.List()
	if err == nil {
		t.Fatal("a ticket that is both on the board and in the logbook was not reported")
	}
	var pe *PartialError
	if !errors.As(err, &pe) {
		t.Fatalf("unexpected error type: %v", err)
	}
	joined := strings.Join(pe.Problems, "\n")
	if !strings.Contains(joined, "filed away") || !strings.Contains(joined, Handle(id)) {
		t.Errorf("the report does not name the ticket and both places:\n%s", joined)
	}
	// And it says the person decides, because guessing would sometimes reopen
	// finished work.
	if !strings.Contains(joined, "only you can say which") {
		t.Errorf("the report decides for the user:\n%s", joined)
	}
}

// The archive counts the same way: a ticket cannot be abandoned and open at
// once.
func TestATicketOnTheBoardAndInTheArchiveIsReported(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	s, _ := At(root)
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	id := NewID(time.Now())
	doc := "---\nid: " + id + "\ntitle: abandoned twice\nstatus: backlog\n---\n\n# abandoned twice\n"
	if err := WriteAtomic(filepath.Join(s.TicketsDir(), Filename(id, "abandoned twice")), []byte(doc)); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(s.ArchiveDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(filepath.Join(s.ArchiveDir(), Filename(id, "abandoned twice")), []byte(doc)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.List(); err == nil {
		t.Fatal("a ticket on the board and in the archive was not reported")
	}
}
