package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
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

// A doorway lane files the just-landed ticket and everything still sitting in
// the lane straight into the logbook, and settleMessage names every file with
// its restore path. The holds (cap) branch is pinned at the core and CLI
// layers.
//
// The doorway is switched on here rather than taken from the shipped lanes:
// jaira ships done without one, because filing is a cut somebody makes with
// 'jaira logbook --all' and not something finishing a ticket does to
// everybody else's work.
func TestSettleLaneFilesADoorwayLane(t *testing.T) {
	s, ts := holdsStore(t, 11)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	l, ok := m.lanes.Get("done")
	if !ok {
		t.Fatal(`"done" lane not found`)
	}
	l.LogbookOnEntry = true
	trimmed, filed, err := lane.Settle(s, l, m.settleFolder(), "", m.settlePrepare())
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	msg := settleMessage(trimmed, filed, l.Holds, "done")
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

// settleLane (now core/move.Move's own settle step) skipped a lane that
// resolves to nothing; the same guard is m.lanes.Get returning false.
func TestSettleLaneIgnoresAnUnknownLane(t *testing.T) {
	s, _ := holdsStore(t, 1)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.lanes.Get("nosuch"); ok {
		t.Fatal(`"nosuch" unexpectedly resolves to a lane`)
	}
}
