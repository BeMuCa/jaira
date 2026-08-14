package gate

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Parking a ticket is the cheapest move on the board and the easiest to forget.
// A blocked ticket with no recorded reason looks the same on day one and day
// ninety, so the blocked lane refuses to take one.
func TestBlockedLaneNeedsAReason(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked"})
	if !hasFieldViolation(vs, CodeNeedsReason, ticket.FieldBlockedReason) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeNeedsReason, ticket.FieldBlockedReason, vs)
	}
}

func TestBlockedLaneAcceptsAReasonOnTheMove(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked", Reason: "vendor API returns 500s"})
	if hasFieldViolation(vs, CodeNeedsReason, ticket.FieldBlockedReason) {
		t.Errorf("a reason was given, expected no %s violation, got %v", CodeNeedsReason, vs)
	}
}

func TestBlockedLaneAcceptsAReasonAlreadyOnTheTicket(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.BlockedReason = "waiting on the security review"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked"})
	if hasFieldViolation(vs, CodeNeedsReason, ticket.FieldBlockedReason) {
		t.Errorf("the ticket already carries a reason, expected no %s violation, got %v", CodeNeedsReason, vs)
	}
}

// blocked-by already answers "waiting on what?" when the blocker is another
// ticket, so demanding the same thing be typed out twice would be a toll, not
// a gate.
func TestBlockedLaneAcceptsBlockedByInstead(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.BlockedBy = []string{"01KZZR4CBGDM5T35SZDR72PQYG"}

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked"})
	if hasFieldViolation(vs, CodeNeedsReason, ticket.FieldBlockedReason) {
		t.Errorf("blocked-by is recorded, expected no %s violation, got %v", CodeNeedsReason, vs)
	}
}

// The refusal has to say what to do, not only what is wrong: an agent reads it
// and acts on it without a human translating.
func TestBlockedReasonRefusalNamesTheFix(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked"})
	if err := vs.Err(); err == nil || !strings.Contains(err.Error(), "--reason") {
		t.Errorf("the refusal does not name the flag that fixes it: %v", err)
	}
}

// The dependency gate guards starting work. Parking is the opposite: the ticket
// most entitled to sit in Blocked is exactly the one whose blockers are still
// open, so entry must not be refused on those blockers.
func TestUnresolvedBlockerDoesNotBarTheBlockedLane(t *testing.T) {
	open := ticketWith("")
	open.ID = "01KZZR4CBGDM5T35SZDR72PQYG"
	open.Status = "in-progress"

	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.BlockedBy = []string{open.ID}

	env := testEnv(t)
	env.All = []*ticket.Ticket{open, tk}

	vs := CheckAdvance(env, tk, Request{To: "blocked", Actor: "a"})
	if err := vs.Err(); err != nil {
		t.Errorf("a ticket with an open blocker must be parkable in blocked, got: %v", err)
	}
}

// A ticket stopped mid-work has not produced its lane's output — that is what
// stopped it. If the leaving lane's contract fired on the way to Blocked, an
// implementing ticket with no commits yet could never be parked at all.
func TestParkingIsExemptFromTheLeavingLanesContract(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	// No commits, no outcome: work stopped before producing anything.
	tk.Outcome = ticket.Outcome{}

	vs := CheckAdvance(testEnv(t), tk, Request{To: "blocked", Reason: "vendor API down", Actor: "a"})
	if err := vs.Err(); err != nil {
		t.Errorf("parking mid-work must not demand the lane's output, got: %v", err)
	}
}
