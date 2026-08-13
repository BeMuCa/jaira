package gate

import (
	"os"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// scrapTicket is a ticket with only a title: nothing that the promotion gate
// asks for is filled in.
func scrapTicket(status string) *ticket.Ticket {
	return &ticket.Ticket{
		ID:     "01KZTT3XZ2YQBX93TTSR7BVRCX",
		Title:  "a one-line scrap",
		Status: status,
	}
}

// installCustomLane points JAIRA_LANES_DIR at a temp directory containing one
// custom lane, the way a real installation would, and returns an Env built
// from it. t.Setenv restores the previous value when the test ends.
func installCustomLane(t *testing.T, filename, body string) Env {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/"+filename, []byte(body), 0o644); err != nil {
		t.Fatalf("write custom lane: %v", err)
	}
	t.Setenv("JAIRA_LANES_DIR", dir)
	return testEnv(t)
}

// This is today's behaviour and must be unchanged: a ticket may not leave the
// backlog for todo without goal, definition-of-done, context and assignee.
func TestBacklogToTodoMissingFieldsRefused(t *testing.T) {
	tk := scrapTicket("backlog")
	vs := CheckAdvance(testEnv(t), tk, Request{To: "todo"})
	for _, f := range []string{ticket.FieldGoal, ticket.FieldDoD, ticket.FieldContext, ticket.FieldAssignee} {
		if !hasFieldViolation(vs, CodeMissingField, f) {
			t.Errorf("expected a %s violation naming %s, got %v", CodeMissingField, f, vs)
		}
	}
}

func TestBacklogToTodoAllFieldsPresentAllowed(t *testing.T) {
	tk := scrapTicket("backlog")
	tk.Goal, tk.DoD, tk.Context, tk.Assignee = "g", "d", "c", "a"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "todo"})
	for _, v := range vs {
		if v.Code == CodeMissingField {
			t.Errorf("unexpected missing-field violation: %v", v)
		}
	}
}

// Skipping todo must not skip the gate: going straight from backlog to
// in-progress still demands the promotion fields.
func TestBacklogToInProgressSkippingTodoMissingFieldsRefused(t *testing.T) {
	tk := scrapTicket("backlog")
	vs := CheckAdvance(testEnv(t), tk, Request{To: "in-progress"})
	if !hasFieldViolation(vs, CodeMissingField, ticket.FieldGoal) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingField, ticket.FieldGoal, vs)
	}
}

// A ticket with only a title can sit in, and move into, a lane placed before
// the specified zone that carries no requires-specified.
func TestScrapTicketMovesIntoLaneBeforeSpecifiedZone(t *testing.T) {
	env := installCustomLane(t, "triage.md", "---\nid: triage\nname: Triage\nafter: backlog\n---\n")
	tk := scrapTicket("backlog")
	vs := CheckAdvance(env, tk, Request{To: "triage"})
	if len(vs) != 0 {
		t.Errorf("expected the move into a pre-zone lane to be allowed, got %v", vs)
	}
}

// The same ticket is refused the moment it tries to cross from that lane into
// todo, the start of the specified zone.
func TestScrapTicketRefusedEnteringSpecifiedZone(t *testing.T) {
	env := installCustomLane(t, "triage.md", "---\nid: triage\nname: Triage\nafter: backlog\n---\n")
	tk := scrapTicket("triage")
	vs := CheckAdvance(env, tk, Request{To: "todo"})
	if !hasFieldViolation(vs, CodeMissingField, ticket.FieldGoal) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingField, ticket.FieldGoal, vs)
	}
}

// A move between two lanes both already inside the specified zone must not
// re-run the promotion check, even for a ticket that (however it got there)
// is missing its promotion fields.
func TestMoveWithinSpecifiedZoneDoesNotRerunPromotionCheck(t *testing.T) {
	tk := scrapTicket("todo")
	vs := CheckAdvance(testEnv(t), tk, Request{To: "in-progress"})
	for _, v := range vs {
		if v.Code == CodeMissingField {
			t.Errorf("promotion check re-ran on a move inside the specified zone: %v", v)
		}
	}
}

// With a lane set where no lane declares requires-specified, the gate falls
// back to today's rule rather than letting an unspecified ticket through.
func TestFallbackWithNoDeclaredBoundaryStillRefuses(t *testing.T) {
	env := testEnv(t)
	for _, l := range env.Lanes.Lanes {
		l.RequiresSpecified = false
	}
	tk := scrapTicket("backlog")
	vs := CheckAdvance(env, tk, Request{To: "todo"})
	if !hasFieldViolation(vs, CodeMissingField, ticket.FieldGoal) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingField, ticket.FieldGoal, vs)
	}
}

// A ticket belongs to its assignee: a move by anyone else, in an ordinary
// lane, is refused, and the refusal names the owner.
func TestOwnershipRefusesWriteByOtherThanAssignee(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.Assignee = "anna"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: "bob"})
	found := findViolation(vs, CodeNotOwner)
	if found == nil {
		t.Fatalf("expected a %s violation, got %v", CodeNotOwner, vs)
	}
	if !strings.Contains(found.Message, "anna") {
		t.Errorf("message %q does not name the owner", found.Message)
	}
}

// The human checkpoint lanes are exempt: reviewing and signing off someone
// else's work is the entire point of them.
func TestOwnershipAllowsLeavingRequiresHumanExitLane(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "signoff"
	tk.Assignee = "anna"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "done", Actor: "bob"})
	if found := findViolation(vs, CodeNotOwner); found != nil {
		t.Errorf("expected no ownership violation leaving a human-exit lane, got %v", found)
	}
}

func TestOwnershipAllowsLeavingRequiresQuestionLane(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "human"
	tk.Assignee = "anna"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: "bob"})
	if found := findViolation(vs, CodeNotOwner); found != nil {
		t.Errorf("expected no ownership violation leaving a requires-question lane, got %v", found)
	}
}

// Taking a ticket over is always allowed, so a ticket is never frozen by an
// absent owner.
func TestOwnershipAllowsHandOver(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.Assignee = "anna"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: "bob", NewAssignee: "bob"})
	if found := findViolation(vs, CodeNotOwner); found != nil {
		t.Errorf("expected a hand-over to be allowed, got %v", found)
	}
}

// A ticket with no assignee belongs to nobody and is writable by anyone.
func TestOwnershipAllowsWhenNoAssignee(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.Assignee = ""
	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: "bob"})
	if found := findViolation(vs, CodeNotOwner); found != nil {
		t.Errorf("expected an unassigned ticket to be writable by anyone, got %v", found)
	}
}

// The comparison is case-insensitive and trimmed: the same person spelled
// differently is still the owner.
func TestOwnershipAllowsSameActorDifferentCase(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.Assignee = "Anna"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: " anna "})
	if found := findViolation(vs, CodeNotOwner); found != nil {
		t.Errorf("expected the assignee to be able to write their own ticket, got %v", found)
	}
}

func findViolation(vs Violations, code string) *Violation {
	for i := range vs {
		if vs[i].Code == code {
			return &vs[i]
		}
	}
	return nil
}
