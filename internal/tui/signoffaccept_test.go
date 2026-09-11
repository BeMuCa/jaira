package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The accept key is the fourth status write-site and the usual way a person
// lands a ticket in done — since done is a doorway, accepting files the ticket
// (and any residents) straight into the logbook, commits kept, and says so.
func TestAcceptFilesTheTicketIntoTheLogbook(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))
	t.Setenv("JAIRA_USER", "berk")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	// Two residents from before the doorway existed — they go along.
	for i := 0; i < 2; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		if _, err := s.Create(map[string]string{
			ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: fmt.Sprintf("t%02d", i),
			ticket.FieldStatus: "done", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, ""); err != nil {
			t.Fatal(err)
		}
	}
	// One at signoff, fully gated so accept() is not refused.
	stamp := base.Add(time.Hour)
	body := "# x\n\n## Definition of Done\n\n- [x] it works\n  proof: shipped\n"
	sg, err := s.Create(map[string]string{
		ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: "the accepted one",
		ticket.FieldStatus: "signoff", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
		ticket.FieldUpdatedAt: ticket.FormatTime(stamp), ticket.FieldAssignee: "berk",
		ticket.FieldGoal: "g", ticket.FieldContext: "c", ticket.FieldDoD: "it works",
		"outcome-what": "w", "outcome-why": "y", "outcome-resolves": "r",
		"review-summary": "s", "review-gaps": "none", "review-verdict": "v", "review-check": "c",
	}, map[string][]string{ticket.FieldCommits: {"deadbeef"}}, body)
	if err != nil {
		t.Fatal(err)
	}

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	full, err := s.Load(sg.ID)
	if err != nil {
		t.Fatal(err)
	}
	m.detail = full
	m.accept()

	if !strings.Contains(m.message, "Accepted") {
		t.Fatalf("accept() did not accept (message: %q)", m.message)
	}
	// Accepting files nothing. It says the work is accepted, which is not a
	// statement about anybody's bookkeeping: the ticket waits in the terminal
	// lane with the others until somebody cuts with 'jaira logbook --all'.
	if got := strings.Count(m.message, "filed to the logbook"); got != 0 {
		t.Errorf("accepting filed %d ticket(s):\n%s", got, m.message)
	}
	if _, err := os.Stat(full.Path); err != nil {
		t.Errorf("the accepted ticket left the board: %v", err)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) == 0 {
		t.Error("the board is empty after an accept; the finished tickets should still be there")
	}
	// And nothing was filed by the accept at all — the logbook is untouched
	// until somebody cuts.
	matches, err := filepath.Glob(filepath.Join(s.LogbookDir(), "*", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Errorf("accepting put %d ticket(s) in the logbook", len(matches))
	}
}
