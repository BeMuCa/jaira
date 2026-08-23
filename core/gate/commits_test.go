package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The commits requirement sits on the done lane, not on the implementing
// lane's way out: work may move through review before it is committed, but
// nothing is accepted that cannot be checked — the diff shown at sign-off and
// recalled months later is the diff of exactly these commits.
func TestCommitsRequiredToEnterDone(t *testing.T) {
	tk := ticketWith(`## Definition of Done

- [x] it works
`)
	tk.Status = "signoff"
	// Commits deliberately left empty.

	vs := CheckAdvance(testEnv(t), tk, Request{To: "done", Actor: "a", Interactive: true})
	if !hasFieldViolation(vs, CodeNeedsCommits, ticket.FieldCommits) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeNeedsCommits, ticket.FieldCommits, vs)
	}
}

func TestCommitsRecordedEnterDone(t *testing.T) {
	tk := ticketWith(`## Definition of Done

- [x] it works
`)
	tk.Status = "signoff"
	tk.Commits = []string{"a1b2c3d"}

	vs := CheckAdvance(testEnv(t), tk, Request{To: "done", Actor: "a", Interactive: true})
	if err := vs.Err(); err != nil {
		t.Errorf("commits recorded and DoD met, expected the move into done to pass, got: %v", err)
	}
}

// Leaving the implementing lane without commits is allowed — recording them
// there is encouraged by the prompt, required only at acceptance.
func TestNoCommitsStillLeavesImplementing(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "review", Actor: "a"})
	if hasFieldViolation(vs, CodeMissingProduces, ticket.FieldCommits) ||
		hasFieldViolation(vs, CodeNeedsCommits, ticket.FieldCommits) {
		t.Errorf("leaving in-progress without commits must not be refused, got %v", vs)
	}
}

// TestDeriveCommitsFallbackAcceptsAnEmptyTicket asserts a ticket that records
// no commits of its own is still accepted when git can find them: jaira works
// the commits out from git itself rather than requiring a person to type a
// sha by hand (D-01).
func TestDeriveCommitsFallbackAcceptsAnEmptyTicket(t *testing.T) {
	tk := ticketWith(`## Definition of Done

- [x] it works
`)
	tk.Status = "signoff"
	// Commits deliberately left empty; the derivation stands in for them.

	env := testEnv(t)
	env.DeriveCommits = func(*ticket.Ticket) []string { return []string{"a1b2c3d"} }

	vs := CheckAdvance(env, tk, Request{To: "done", Actor: "a", Interactive: true})
	if hasFieldViolation(vs, CodeNeedsCommits, ticket.FieldCommits) {
		t.Errorf("expected no %s violation when DeriveCommits finds a sha, got %v", CodeNeedsCommits, vs)
	}
}

// TestDeriveCommitsNilOrEmptyStillRefuses asserts the existing refusal is
// unchanged both when no derivation is on offer (nil, the zero value every
// existing gate test already builds) and when a derivation runs and finds
// nothing — neither is a silent pass (D-01).
func TestDeriveCommitsNilOrEmptyStillRefuses(t *testing.T) {
	cases := []struct {
		name   string
		derive func(*ticket.Ticket) []string
	}{
		{name: "nil derivation", derive: nil},
		{name: "empty derivation", derive: func(*ticket.Ticket) []string { return nil }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tk := ticketWith(`## Definition of Done

- [x] it works
`)
			tk.Status = "signoff"

			env := testEnv(t)
			env.DeriveCommits = c.derive

			vs := CheckAdvance(env, tk, Request{To: "done", Actor: "a", Interactive: true})
			if !hasFieldViolation(vs, CodeNeedsCommits, ticket.FieldCommits) {
				t.Errorf("expected a %s violation naming %s, got %v", CodeNeedsCommits, ticket.FieldCommits, vs)
			}
		})
	}
}

// TestDeriveCommitsNeverCalledWhenTicketAlreadyHasCommits asserts explicit
// always beats derived: a ticket with commits already recorded must not even
// invoke the derivation (D-01).
func TestDeriveCommitsNeverCalledWhenTicketAlreadyHasCommits(t *testing.T) {
	tk := ticketWith(`## Definition of Done

- [x] it works
`)
	tk.Status = "signoff"
	tk.Commits = []string{"a1b2c3d"}

	calls := 0
	env := testEnv(t)
	env.DeriveCommits = func(*ticket.Ticket) []string {
		calls++
		return []string{"unused"}
	}

	if err := CheckAdvance(env, tk, Request{To: "done", Actor: "a", Interactive: true}).Err(); err != nil {
		t.Errorf("commits recorded and DoD met, expected the move into done to pass, got: %v", err)
	}
	if calls != 0 {
		t.Errorf("DeriveCommits was called %d time(s) for a ticket that already records commits, want 0", calls)
	}
}
