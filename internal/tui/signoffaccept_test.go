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
// lands a ticket in done — the cap must fire there exactly as it does on a
// move. This test came out of the review that found accept() skipping it.
func TestAcceptEnforcesTheCap(t *testing.T) {
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
	// Ten already in done, a minute apart, oldest first.
	var done []*ticket.Ticket
	for i := 0; i < 10; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: fmt.Sprintf("t%02d", i),
			ticket.FieldStatus: "done", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		done = append(done, tk)
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

	after, err := s.Load(sg.ID)
	if err != nil {
		t.Fatalf("reload after accept: %v", err)
	}
	if after.Status != "done" {
		t.Fatalf("accept() was refused, ticket still in %q (message: %q)", after.Status, m.message)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, tk := range all {
		if tk.Status == "done" {
			n++
		}
	}
	if n != 10 {
		t.Errorf("done holds %d after accept, want 10 — the cap did not fire on the accept path", n)
	}
	if _, err := os.Stat(done[0].Path); !os.IsNotExist(err) {
		t.Errorf("oldest ticket %s is still on the board", ticket.Handle(done[0].ID))
	}
	if !strings.Contains(m.message, "logbook") || !strings.Contains(m.message, ticket.Handle(done[0].ID)) {
		t.Errorf("accept message does not say who left for the logbook: %q", m.message)
	}
}
