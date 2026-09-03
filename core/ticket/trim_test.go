package ticket

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// trimStore builds a store whose done lane holds n tickets, updated a minute
// apart, and returns them oldest first.
func trimStore(t *testing.T, n int) (*Store, []*Ticket) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	var out []*Ticket
	for i := 0; i < n; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			FieldID:        NewID(stamp),
			FieldTitle:     fmt.Sprintf("t%02d", i),
			FieldStatus:    "done",
			FieldCreatedAt: FormatTime(stamp),
			FieldUpdatedAt: FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, tk)
	}
	return s, out
}

func TestOverflowReturnsExactlyTheOldestBeyondTheCap(t *testing.T) {
	s, ts := trimStore(t, 12)
	over, err := s.Overflow("done", 10)
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 2 {
		t.Fatalf("overflow of 12 over 10 = %d tickets, want 2", len(over))
	}
	if over[0].ID != ts[0].ID || over[1].ID != ts[1].ID {
		t.Errorf("overflow = %s, %s — want the two oldest %s, %s",
			over[0].Title, over[1].Title, ts[0].Title, ts[1].Title)
	}
}

func TestOverflowAtOrUnderTheCapIsEmpty(t *testing.T) {
	s, _ := trimStore(t, 10)
	over, err := s.Overflow("done", 10)
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 0 {
		t.Errorf("exactly 10 over a cap of 10 overflows %d tickets, want none", len(over))
	}
	// keep <= 0 means no cap at all.
	over, err = s.Overflow("done", 0)
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 0 {
		t.Errorf("a cap of 0 (unlimited) overflows %d tickets, want none", len(over))
	}
}

func TestTrimLaneFilesTheOverflowAndRestoreBringsItBack(t *testing.T) {
	s, ts := trimStore(t, 11)
	trimmed, err := s.TrimLane("done", 10, "bc-20260903")
	if err != nil {
		t.Fatalf("TrimLane: %v", err)
	}
	if len(trimmed) != 1 || trimmed[0].ID != ts[0].ID {
		t.Fatalf("trimmed %v, want exactly the oldest %s", trimmed, ts[0].Title)
	}
	if _, err := os.Stat(trimmed[0].Path); err != nil {
		t.Errorf("trimmed file not in the logbook at %q: %v", trimmed[0].Path, err)
	}
	if _, err := os.Stat(ts[0].Path); !os.IsNotExist(err) {
		t.Errorf("oldest ticket still on the board at %q", ts[0].Path)
	}
	left, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(left) != 10 {
		t.Errorf("%d tickets remain on the board, want 10", len(left))
	}
	// The way back exists: restore returns the ticket to the board.
	if _, err := s.Restore(filepath.Base(trimmed[0].Path)); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	left, err = s.List()
	if err != nil {
		t.Fatalf("List after restore: %v", err)
	}
	if len(left) != 11 {
		t.Errorf("%d tickets after restore, want 11", len(left))
	}
}
