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

// The builtin done is a doorway: settleLane files the just-landed ticket and
// everything still sitting in the lane straight into the logbook — the lane
// is self-migrating, and the message names every file with its restore path.
// The holds (cap) branch of settleLane is pinned at the core and CLI layers.
func TestSettleLaneFilesTheDoorwayLane(t *testing.T) {
	s, ts := holdsStore(t, 11)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := m.settleLane("done", "")
	if err != nil {
		t.Fatalf("settleLane: %v", err)
	}
	if got := strings.Count(msg, "filed to the logbook"); got != 11 {
		t.Errorf("%d filing lines, want 11:\n%s", got, msg)
	}
	if !strings.Contains(msg, ticket.Handle(ts[0].ID)) {
		t.Errorf("message does not name the oldest resident:\n%s", msg)
	}
	left, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Errorf("%d tickets remain on the board, want 0 — done is a doorway", len(left))
	}
}

func TestSettleLaneIgnoresAnUnknownLane(t *testing.T) {
	s, _ := holdsStore(t, 1)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if msg, err := m.settleLane("nosuch", ""); msg != "" || err != nil {
		t.Errorf("settleLane on an unknown lane = %q, %v; want silence", msg, err)
	}
}
