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
	s, err := lane.Load("")
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

// A handle typed into prose reads as a dependency to a human, but the gates
// only ever read blocked-by — this is the gap that makes the board lie about
// a ticket being actionable.
func TestUndeclaredDependencyMentionIsReported(t *testing.T) {
	dep := tk(ticket.NewID(time.Now()), "dependency", "todo")
	citing := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "citing", "todo")
	citing.Context = "waiting on the auth work, see " + ticket.Handle(dep.ID)
	ps := Tickets([]*ticket.Ticket{dep, citing}, lanes(t))
	if !has(ps, CodeUndeclaredDep) {
		t.Errorf("an undeclared handle mention was not reported: %v", codes(ps))
	}
	for _, p := range ps {
		if p.Code != CodeUndeclaredDep {
			continue
		}
		if p.Severity != SeverityWarning {
			t.Errorf("an undeclared dependency mention was not a warning: %+v", p)
		}
		if p.Field != ticket.FieldContext {
			t.Errorf("a mention found in context did not carry the context field: %+v", p)
		}
	}
	if HasErrors(ps) {
		t.Error("HasErrors reported true for a warning-only mention")
	}
}

// The note (ticket body) is scanned the same as context, and this pins the
// source label in the message so a later edit that drops the body source
// cannot pass silently.
func TestUndeclaredDependencyMentionInBodyIsReported(t *testing.T) {
	dep := tk(ticket.NewID(time.Now()), "dependency", "todo")
	citing := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "citing", "todo")
	citing.Body = "## Progress\n\n- waiting on " + ticket.Handle(dep.ID) + "\n"
	ps := Tickets([]*ticket.Ticket{dep, citing}, lanes(t))
	var found *Problem
	for i, p := range ps {
		if p.Code == CodeUndeclaredDep {
			found = &ps[i]
		}
	}
	if found == nil {
		t.Fatalf("an undeclared handle mention in the body was not reported: %v", codes(ps))
	}
	if !strings.HasPrefix(found.Message, "note names") {
		t.Errorf("a body-sourced mention did not carry the note source label: %q", found.Message)
	}
	if found.Field != "" {
		t.Errorf("a mention found in the body should carry no field (no constant covers it): %+v", found)
	}
}

// The same handle typed twice in one field must produce one warning, not one
// per occurrence.
func TestRepeatedMentionIsReportedOnce(t *testing.T) {
	dep := tk(ticket.NewID(time.Now()), "dependency", "todo")
	citing := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "citing", "todo")
	h := ticket.Handle(dep.ID)
	citing.Context = "see " + h + ", also see " + h + " again"
	ps := Tickets([]*ticket.Ticket{dep, citing}, lanes(t))
	n := 0
	for _, p := range ps {
		if p.Code == CodeUndeclaredDep {
			n++
		}
	}
	if n != 1 {
		t.Errorf("a handle repeated twice in one field produced %d warnings, want 1: %v", n, codes(ps))
	}
}

// Once the same handle is declared in blocked-by, the mention is exactly what
// it looks like — an explanation, not a hidden dependency — and must stop
// being reported.
func TestDeclaredDependencyMentionIsNotReported(t *testing.T) {
	dep := tk(ticket.NewID(time.Now()), "dependency", "todo")
	citing := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "citing", "todo")
	citing.Context = "waiting on the auth work, see " + ticket.Handle(dep.ID)
	citing.BlockedBy = []string{dep.ID}
	if ps := Tickets([]*ticket.Ticket{dep, citing}, lanes(t)); has(ps, CodeUndeclaredDep) {
		t.Errorf("a declared dependency's mention was reported: %v", codes(ps))
	}
}

// A done ticket blocks nothing, so naming it in prose without a blocked-by
// entry is not the gap this check exists to catch.
func TestUndeclaredMentionOfTerminalTicketIsNotReported(t *testing.T) {
	dep := tk(ticket.NewID(time.Now()), "dependency", "done")
	citing := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "citing", "todo")
	citing.Context = "see " + ticket.Handle(dep.ID)
	if ps := Tickets([]*ticket.Ticket{dep, citing}, lanes(t)); has(ps, CodeUndeclaredDep) {
		t.Errorf("a mention of a done ticket was reported: %v", codes(ps))
	}
}

// A token can share a handle's six-character shape without being one; only a
// token that also resolves to a ticket actually on the board may warn.
func TestShapeMatchThatResolvesToNothingIsNotReported(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "note-writer", "todo")
	a.Body = "spec says ABCDEF for the header, which is not a ticket"
	if ps := Tickets([]*ticket.Ticket{a}, lanes(t)); has(ps, CodeUndeclaredDep) {
		t.Errorf("a shape-matching token that resolves to no ticket was reported: %v", codes(ps))
	}
}

// A ticket naming its own handle is not a dependency on itself.
func TestOwnHandleMentionIsNotReported(t *testing.T) {
	a := tk(ticket.NewID(time.Now()), "self-referential", "todo")
	a.Context = "tracked as " + ticket.Handle(a.ID)
	if ps := Tickets([]*ticket.Ticket{a}, lanes(t)); has(ps, CodeUndeclaredDep) {
		t.Errorf("a ticket naming its own handle was reported: %v", codes(ps))
	}
}

// A follow-up naming its parent in context is declaring the relation
// follows: already exists for. The parent does not block the follow-up, so
// this must not warn — and must not warn for any handle equal to follows,
// while an unrelated undeclared handle alongside it still should.
func TestFollowsMentionIsNotReported(t *testing.T) {
	parent := tk(ticket.NewID(time.Now()), "parent", "todo")
	other := tk(ticket.NewID(time.Now().Add(time.Millisecond)), "other", "todo")
	followUp := tk(ticket.NewID(time.Now().Add(2*time.Millisecond)), "follow-up", "todo")
	followUp.Follows = parent.ID
	followUp.Context = "review of " + ticket.Handle(parent.ID)

	ps := Tickets([]*ticket.Ticket{parent, other, followUp}, lanes(t))
	if has(ps, CodeUndeclaredDep) {
		t.Errorf("a follow-up naming its parent was reported: %v", codes(ps))
	}

	followUp.Context += ", also see " + ticket.Handle(other.ID)
	ps = Tickets([]*ticket.Ticket{parent, other, followUp}, lanes(t))
	if !has(ps, CodeUndeclaredDep) {
		t.Errorf("an unrelated undeclared handle next to a follows mention was not reported: %v", codes(ps))
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
