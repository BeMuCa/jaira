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
//
// A value that differs from the canonical form only in its whitespace is the
// same failure reached by a value that is not even a typo: CanonicalMode trims
// it, but nothing on the read path does, so " conversational " prints as a mode
// and is compared against exactly one word by the worker, which does not match.
func TestBadModeIsReportedWithTheRepair(t *testing.T) {
	for _, c := range []struct{ name, mode string }{
		{"hand-written mode", "chat"},
		{"untrimmed mode", " " + ticket.ModeConversational + " "},
	} {
		t.Run(c.name, func(t *testing.T) {
			bad := tk(ticket.NewID(time.Now()), c.name, "todo")
			bad.Goal, bad.Context, bad.Assignee, bad.DoD = "g", "c", "berk", "d"
			bad.Mode = c.mode

			ps := Tickets([]*ticket.Ticket{bad}, lanes(t), nil)
			var p Problem
			for _, cand := range ps {
				if cand.Code == CodeBadMode {
					p = cand
				}
			}
			if p.Code == "" {
				t.Fatalf("mode %q produced %v", c.mode, codes(ps))
			}
			if p.Severity != SeverityWarning {
				t.Errorf("severity = %q, want a warning: the ticket itself is intact", p.Severity)
			}
			if p.Field != ticket.FieldMode {
				t.Errorf("field = %q, want %q", p.Field, ticket.FieldMode)
			}
			if !strings.Contains(p.Message, c.mode) {
				t.Errorf("the message does not name the offending value: %q", p.Message)
			}
			if !strings.Contains(p.Message, "mode="+ticket.ModeConversational) {
				t.Errorf("the message does not name the repair: %q", p.Message)
			}
		})
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

// stop-after-fix reads as on only for "true": a hand-written "yes" reads as
// off without a word, so validate names it and the repair; what jaira reads
// as meant is silent.
func TestUnreadableStopAfterFixIsReported(t *testing.T) {
	for _, c := range []struct {
		value string
		warn  bool
	}{
		{"yes", true}, {"on", true}, {"1", true},
		{"", false}, {"true", false}, {"false", false}, {"True", false}, {" true ", false},
	} {
		tkt := tk(ticket.NewID(time.Now()), "saf", "todo")
		tkt.Goal, tkt.Context, tkt.Assignee, tkt.DoD = "g", "c", "berk", "d"
		tkt.StopAfterFixRaw = c.value
		ps := Tickets([]*ticket.Ticket{tkt}, lanes(t), nil)
		if has(ps, CodeBadStopAfterFix) != c.warn {
			t.Errorf("stop-after-fix %q: reported %v, want %v (%v)", c.value, !c.warn, c.warn, codes(ps))
		}
		for _, p := range ps {
			if p.Code == CodeBadStopAfterFix && (p.Severity != SeverityWarning || !strings.Contains(p.Message, c.value) || !strings.Contains(p.Message, "stop-after-fix=true")) {
				t.Errorf("stop-after-fix %q: %+v, want a warning naming the value and the repair", c.value, p)
			}
		}
	}
}
