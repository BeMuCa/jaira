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
	if got := strings.Count(m.message, "filed to the logbook"); got != 3 {
		t.Errorf("%d filing lines in the accept message, want 3:\n%s", got, m.message)
	}
	if _, err := os.Stat(full.Path); !os.IsNotExist(err) {
		t.Errorf("the accepted ticket is still on the board at %q", full.Path)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Errorf("%d tickets remain on the board, want 0", len(all))
	}
	// The commits survived the filing.
	matches, err := filepath.Glob(filepath.Join(s.LogbookDir(), "*", filepath.Base(full.Path)))
	if err != nil || len(matches) != 1 {
		t.Fatalf("accepted ticket not found in the logbook: %v (%d matches)", err, len(matches))
	}
	filed, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(filed), "deadbeef") {
		t.Errorf("the filed ticket lost its commits:\n%s", filed)
	}
}
