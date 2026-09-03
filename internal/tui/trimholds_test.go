package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// holdsStore builds a store whose done lane holds n tickets, updated a minute
// apart, and returns them oldest first.
func holdsStore(t *testing.T, n int) (*ticket.Store, []*ticket.Ticket) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	var out []*ticket.Ticket
	for i := 0; i < n; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			ticket.FieldID:        ticket.NewID(stamp),
			ticket.FieldTitle:     fmt.Sprintf("t%02d", i),
			ticket.FieldStatus:    "done",
			ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, tk)
	}
	return s, out
}

// The TUI's move enforces the same cap the CLI's does, and says so the same
// way: which ticket left, and the restore command that brings it back.
func TestTrimHoldsFilesTheOldestAndSaysSo(t *testing.T) {
	s, ts := holdsStore(t, 11)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := m.trimHolds("done", "")
	if err != nil {
		t.Fatalf("trimHolds: %v", err)
	}
	oldest := ts[0]
	if !strings.Contains(msg, ticket.Handle(oldest.ID)) || !strings.Contains(msg, "logbook") {
		t.Errorf("trim message %q does not name %s and the logbook", msg, ticket.Handle(oldest.ID))
	}
	left, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(left) != 10 {
		t.Errorf("%d tickets remain on the board, want 10", len(left))
	}
}

func TestTrimHoldsAtTheCapMovesNothing(t *testing.T) {
	s, _ := holdsStore(t, 10)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := m.trimHolds("done", "")
	if err != nil {
		t.Fatalf("trimHolds: %v", err)
	}
	if msg != "" {
		t.Errorf("a lane at its cap reported a trim: %q", msg)
	}
	left, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(left) != 10 {
		t.Errorf("%d tickets remain on the board, want 10", len(left))
	}
}
