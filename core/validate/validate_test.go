package validate

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

func lanes(t *testing.T) *lane.Set {
	t.Helper()
	s, err := lane.Load()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func tk(id, title, status string) *ticket.Ticket {
	return &ticket.Ticket{
		ID: id, Title: title, Status: status,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
		Path: "/tmp/" + title + ".md",
	}
}

func codes(ps []Problem) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Code)
	}
	return out
}

func has(ps []Problem, code string) bool {
	for _, p := range ps {
		if p.Code == code {
			return true
		}
	}
	return false
}

// The gates fire when a ticket moves. Nothing checks a ticket at rest, so a
// hand-edited or agent-mangled file sits there looking fine until someone tries
// to use it. That is what this exists to catch.
func TestValidReturnsNothing(t *testing.T) {
	good := tk(ticket.NewID(time.Now()), "fine", "todo")
	good.Goal, good.Context, good.Assignee = "g", "c", "berk"
	good.DoD = "a test covers it"
	if ps := Tickets([]*ticket.Ticket{good}, lanes(t)); len(ps) != 0 {
		t.Errorf("a valid ticket produced problems: %v", codes(ps))
	}
}

func TestIDMustBeAULID(t *testing.T) {
	bad := tk("not-a-ulid", "bad id", "todo")
	if !has(Tickets([]*ticket.Ticket{bad}, lanes(t)), CodeBadID) {
		t.Error("a malformed id was accepted")
	}
	missing := tk("", "no id", "todo")
	if !has(Tickets([]*ticket.Ticket{missing}, lanes(t)), CodeBadID) {
		t.Error("a missing id was accepted")
	}
}

// A ticket in a lane this installation does not have is already read-only, but
// nothing tells you it is there until you try to touch it.
func TestUnknownLaneIsReported(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "stranded", "some-lane-nobody-installed")
	ps := Tickets([]*ticket.Ticket{bad}, lanes(t))
	if !has(ps, CodeUnknownLane) {
		t.Errorf("an unknown lane was accepted: %v", codes(ps))
	}
}

func TestMissingTitleIsReported(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "", "todo")
	bad.Path = "/tmp/x.md"
	if !has(Tickets([]*ticket.Ticket{bad}, lanes(t)), CodeNoTitle) {
		t.Error("a ticket with no title was accepted")
	}
}

// A dependency pointing at a ticket that does not exist silently blocks work
// forever: the gate refuses the move and nothing explains why.
func TestDanglingDependencyIsReported(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "depends", "todo")
	a.BlockedBy = []string{"01ZZZZZZZZZZZZZZZZZZZZZZZZ"}
	ps := Tickets([]*ticket.Ticket{a}, lanes(t))
	if !has(ps, CodeDanglingDep) {
		t.Errorf("a dangling dependency was accepted: %v", codes(ps))
	}
}

func TestResolvedDependencyIsFine(t *testing.T) {
	b := tk(ticket.NewID(time.Now()), "dependency", "todo")
	a := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "depends", "todo")
	a.BlockedBy = []string{b.ID}
	if ps := Tickets([]*ticket.Ticket{a, b}, lanes(t)); has(ps, CodeDanglingDep) {
		t.Errorf("a resolvable dependency was reported: %v", codes(ps))
	}
}

// Two files claiming one id means one of them is invisible. The store already
// picks a winner deterministically; this is what tells you it happened.
func TestDuplicateIDIsReported(t *testing.T) {
	id := ticket.NewID(time.Now())
	a, b := tk(id, "first", "todo"), tk(id, "second", "todo")
	ps := Tickets([]*ticket.Ticket{a, b}, lanes(t))
	if !has(ps, CodeDuplicateID) {
		t.Errorf("a duplicate id was accepted: %v", codes(ps))
	}
}

// A ticket cannot depend on itself; the gate treats this as blocked forever.
func TestSelfDependencyIsReported(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "loop", "todo")
	a.BlockedBy = []string{a.ID}
	if !has(Tickets([]*ticket.Ticket{a}, lanes(t)), CodeSelfDep) {
		t.Error("a self-dependency was accepted")
	}
}

// Timestamps are what the merge driver resolves conflicts with, so an unset one
// is a real problem rather than cosmetic.
func TestMissingTimestampIsReported(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "no dates", "todo")
	a.CreatedAt, a.UpdatedAt = time.Time{}, time.Time{}
	ps := Tickets([]*ticket.Ticket{a}, lanes(t))
	if !has(ps, CodeBadTimestamp) {
		t.Errorf("missing timestamps were accepted: %v", codes(ps))
	}
}

// Severity has to be distinguishable, because an incomplete ticket in the
// backlog is normal and a corrupt one is not.
func TestIncompleteBacklogTicketIsOnlyAWarning(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "just captured", "backlog")
	ps := Tickets([]*ticket.Ticket{a}, lanes(t))
	for _, p := range ps {
		if p.Severity == SeverityError {
			t.Errorf("a normal backlog ticket produced an error: %+v", p)
		}
	}
	if !has(ps, CodeIncomplete) {
		t.Errorf("an unspecified ticket was not flagged at all: %v", codes(ps))
	}
	if HasErrors(ps) {
		t.Error("HasErrors reported true for warnings only")
	}
}

func TestHasErrorsDistinguishesSeverity(t *testing.T) {
	bad := tk("nope", "bad", "todo")
	if !HasErrors(Tickets([]*ticket.Ticket{bad}, lanes(t))) {
		t.Error("HasErrors did not report a real error")
	}
}

func TestProblemMessagesNameTheTicket(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "stranded", "nope-lane")
	for _, p := range Tickets([]*ticket.Ticket{bad}, lanes(t)) {
		if p.Handle == "" && p.Path == "" {
			t.Errorf("problem is not attributable to a ticket: %+v", p)
		}
		if strings.TrimSpace(p.Message) == "" {
			t.Errorf("problem has no message: %+v", p)
		}
	}
}
