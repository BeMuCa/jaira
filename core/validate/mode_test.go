package validate

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The two CLI write paths already refuse a mode outside the closed set, so this
// check exists for the ways that go past them: a hand-edited file, a bad merge,
// an agent writing something unexpected. A value outside the set is worse than
// no value, because it prints as a mode and is read as none.
func TestUnknownModeIsReportedWithTheRepair(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "hand-written mode", "todo")
	bad.Goal, bad.Context, bad.Assignee, bad.DoD = "g", "c", "berk", "d"
	bad.Mode = "chat"

	ps := Tickets([]*ticket.Ticket{bad}, lanes(t), nil)
	if !has(ps, CodeBadMode) {
		t.Fatalf("an unknown mode produced %v", codes(ps))
	}
	var p Problem
	for _, cand := range ps {
		if cand.Code == CodeBadMode {
			p = cand
		}
	}
	if p.Severity != SeverityWarning {
		t.Errorf("severity = %q, want a warning: the ticket itself is intact", p.Severity)
	}
	if p.Field != ticket.FieldMode {
		t.Errorf("field = %q, want %q", p.Field, ticket.FieldMode)
	}
	if !strings.Contains(p.Message, "chat") {
		t.Errorf("the message does not name the offending value: %q", p.Message)
	}
	if !strings.Contains(p.Message, "mode="+ticket.ModeConversational) {
		t.Errorf("the message does not name the repair: %q", p.Message)
	}
}

// A value that differs from the canonical form only in its whitespace is the
// same failure as an unknown one: CanonicalMode trims it, but nothing on the
// read path does, so " conversational " prints as a mode and is compared
// against exactly one word by the worker, which does not match.
func TestUntrimmedModeIsReportedWithTheRepair(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "untrimmed mode", "todo")
	bad.Goal, bad.Context, bad.Assignee, bad.DoD = "g", "c", "berk", "d"
	bad.Mode = " " + ticket.ModeConversational + " "

	ps := Tickets([]*ticket.Ticket{bad}, lanes(t), nil)
	if !has(ps, CodeBadMode) {
		t.Fatalf("an untrimmed mode produced %v", codes(ps))
	}
	var p Problem
	for _, cand := range ps {
		if cand.Code == CodeBadMode {
			p = cand
		}
	}
	if p.Severity != SeverityWarning {
		t.Errorf("severity = %q, want a warning: the ticket itself is intact", p.Severity)
	}
	if !strings.Contains(p.Message, "mode="+ticket.ModeConversational) {
		t.Errorf("the message does not name the repair: %q", p.Message)
	}
}

// Both values the field actually takes are silent: no mode at all is the
// ordinary autonomous run, and the one word a worker compares against is the
// point of the field.
func TestKnownModesAreNotReported(t *testing.T) {
	for _, mode := range []string{"", ticket.ModeConversational} {
		fine := tk(ticket.NewID(time.Now()), "fine mode", "todo")
		fine.Goal, fine.Context, fine.Assignee, fine.DoD = "g", "c", "berk", "d"
		fine.Mode = mode
		if ps := Tickets([]*ticket.Ticket{fine}, lanes(t), nil); has(ps, CodeBadMode) {
			t.Errorf("mode %q was reported: %v", mode, codes(ps))
		}
	}
}
