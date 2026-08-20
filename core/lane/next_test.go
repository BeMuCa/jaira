package lane

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// inLane builds a ticket sitting in a lane, with the named Options ticked.
func inLane(status string, options ...string) *ticket.Ticket {
	t := &ticket.Ticket{Status: status}
	if len(options) > 0 {
		body := "## Options\n\n"
		for _, o := range options {
			body += "- [x] " + o + "\n"
		}
		t.Body = body
		t.Options = ticket.ParseOptions(body)
	}
	return t
}

func builtinSet(t *testing.T) *Set {
	t.Helper()
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	s, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The route was only ever implicit in the column order. Next answers it once so
// four agents do not each derive it differently.
func TestNextWalksTheBoard(t *testing.T) {
	s := builtinSet(t)
	for _, c := range []struct{ from, want string }{
		{"todo", "in-progress"},
		{"in-progress", "review"}, // human is a stop, not a step
		{"review", "signoff"},
		{"signoff", "done"},
	} {
		got := s.Next(inLane(c.from))
		if got == nil {
			t.Errorf("from %s: no next lane, want %s", c.from, c.want)
			continue
		}
		if got.ID != c.want {
			t.Errorf("from %s: next is %s, want %s", c.from, got.ID, c.want)
		}
	}
}

// An optional step is on this ticket's path only if the ticket says so.
func TestNextSkipsStepsTheTicketDidNotChoose(t *testing.T) {
	s := builtinSet(t)

	if got := s.Next(inLane("backlog")); got == nil || got.ID != "todo" {
		t.Errorf("without options, next after backlog is %v, want todo", idOfLane(got))
	}
	if got := s.Next(inLane("backlog", "brainstorm")); got == nil || got.ID != "brainstorm" {
		t.Errorf("with brainstorm ticked, next after backlog is %v, want brainstorm", idOfLane(got))
	}
	if got := s.Next(inLane("todo", "planning")); got == nil || got.ID != "pre-process" {
		t.Errorf("with planning ticked, next after todo is %v, want pre-process", idOfLane(got))
	}
}

// Being blocked, or being asked a question, is something that happens to a
// ticket — never the next step in its work.
func TestNextIsNeverAParkingOrQuestionLane(t *testing.T) {
	s := builtinSet(t)
	for _, from := range []string{"backlog", "todo", "in-progress", "review", "signoff"} {
		got := s.Next(inLane(from))
		if got == nil {
			continue
		}
		if got.RequiresBlockedReason {
			t.Errorf("from %s: next is the parking lane %s", from, got.ID)
		}
		if got.RequiresQuestion {
			t.Errorf("from %s: next is the question lane %s", from, got.ID)
		}
	}
}

// A ticket waiting on an answer resumes where it stopped, and the board does not
// record where that was — so walking forward would send it past the step that
// asked. Found against a real board: HITL reported "review", skipping critique.
func TestNextIsEmptyForATicketWaitingOnAnAnswer(t *testing.T) {
	s := builtinSet(t)
	if got := s.Next(inLane("human")); got != nil {
		t.Errorf("a ticket in HITL claims to go on to %s", got.ID)
	}
}

// Done is done, and a lane this installation does not have cannot be reasoned
// about.
func TestNextEndsAtTerminalAndUnknown(t *testing.T) {
	s := builtinSet(t)
	if got := s.Next(inLane("done")); got != nil {
		t.Errorf("a terminal lane has a next lane: %s", got.ID)
	}
	if got := s.Next(inLane("no-such-lane")); got != nil {
		t.Errorf("an unknown lane has a next lane: %s", got.ID)
	}
	// Parked work has no "next step" either: it resumes where it stopped.
	if got := s.Next(inLane("blocked")); got != nil {
		t.Errorf("the parking lane has a next lane: %s", got.ID)
	}
}

func idOfLane(l *Lane) string {
	if l == nil {
		return "<nil>"
	}
	return l.ID
}
