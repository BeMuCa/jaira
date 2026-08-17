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
