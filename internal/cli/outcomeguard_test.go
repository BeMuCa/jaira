package cli

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The rule in one table: only the lane that produces the outcome may replace
// it; anyone may fill an empty field or repeat the value already there.
func TestGuardOutcome(t *testing.T) {
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	written := ticket.Outcome{What: "streamed the writer", Why: "it buffered", Resolves: "the export completes"}
	cases := []struct {
		name    string
		status  string
		have    ticket.Outcome
		what    string
		refused bool
	}{
		{"review replaces the implementer's account", "review", written, "review fields set", true},
		{"critique replaces it on the way back", "critique", written, "as before", true},
		{"in-progress rewrites its own output on a later round", "in-progress", written, "streamed it in chunks", false},
		{"review fills an empty field", "review", ticket.Outcome{}, "streamed the writer", false},
		{"review repeats the same value", "review", written, written.What, false},
		{"review passes no --what at all", "review", written, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tk := &ticket.Ticket{ID: "01M4P7CJ8STJNPHCVSRMQTHD8H", Status: c.status, Outcome: c.have}
			err := guardOutcome(tk, lanes, map[string]string{ticket.FieldOutcomeWhat: c.what})
			if got := err != nil; got != c.refused {
				t.Fatalf("refused = %v, want %v (%v)", got, c.refused, err)
			}
			var ce *codedError
			if c.refused && (!errors.As(err, &ce) || ce.code != ExitValidation || ce.reason != "outcome_not_yours") {
				t.Errorf("refusal is not a validation error with reason outcome_not_yours: %#v", err)
			}
		})
	}
}

// The same rule through the real command: the refused move writes nothing,
// and the way out it names works.
func TestMoveRefusesToReplaceAnotherLanesOutcome(t *testing.T) {
	dir, id := moveTestStore(t)
	// Filling the empty outcome on the way in is allowed from any lane.
	if _, err := runCLI(t, dir, "move", id, "--to", "review", "--force",
		"--what", "streamed the writer", "--why", "it buffered", "--resolves", "the export completes"); err != nil {
		t.Fatal(err)
	}
	before := readTicketFile(t, dir)

	out, err := runCLI(t, dir, "move", id, "--to", "signoff", "--force", "--resolves", "review fields set")
	if err == nil {
		t.Fatalf("review replaced the implementer's outcome:\n%s", out)
	}
	if !strings.Contains(err.Error(), "jaira set") {
		t.Errorf("the refusal does not name the explicit way to replace it: %v", err)
	}
	if got := readTicketFile(t, dir); string(got) != string(before) {
		t.Errorf("the refused move wrote the ticket:\n%s", got)
	}

	if _, err := runCLI(t, dir, "set", id, "outcome-resolves=replaced on purpose"); err != nil {
		t.Fatalf("an explicit replacement was refused: %v", err)
	}
	if !strings.Contains(string(readTicketFile(t, dir)), "replaced on purpose") {
		t.Error("jaira set did not replace the outcome")
	}
}
