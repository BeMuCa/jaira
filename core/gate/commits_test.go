package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A ticket that leaves the implementing lane with no commits recorded is a
// change nobody can check: review is handed a diff it cannot assemble, and at
// sign-off a person cannot see what they are accepting. The implementing lane
// declares commits as its output, so the existing lane-output gate refuses the
// move — the same way an empty outcome does.
func TestCommitsRequiredToLeaveImplementing(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	// Commits deliberately left empty.

	vs := CheckAdvance(testEnv(t), tk, Request{To: "review"})
	if !hasFieldViolation(vs, CodeMissingProduces, ticket.FieldCommits) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingProduces, ticket.FieldCommits, vs)
	}
}

func TestCommitsRecordedLeavesImplementing(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "in-progress"
	tk.Commits = []string{"a1b2c3d"}

	vs := CheckAdvance(testEnv(t), tk, Request{To: "review"})
	if hasFieldViolation(vs, CodeMissingProduces, ticket.FieldCommits) {
		t.Errorf("commits are recorded, expected no %s violation for %s, got %v", CodeMissingProduces, ticket.FieldCommits, vs)
	}
}
