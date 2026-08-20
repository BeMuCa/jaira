package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// An empty review-summary or review-gaps must refuse the move out of review,
// the same way any other lane output is enforced — otherwise a reviewer can
// leave the diff unjudged and the ticket moves on regardless.
func TestReviewSummaryRequiredToLeaveReview(t *testing.T) {
	tk := ticketWith("")
	tk.ReviewVerdict = "the diff matches the criteria"
	tk.ReviewGaps = "none"
	// ReviewSummary deliberately left empty.

	vs := CheckAdvance(testEnv(t), tk, Request{To: "signoff"})
	if !hasFieldViolation(vs, CodeMissingProduces, ticket.FieldReviewSummary) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingProduces, ticket.FieldReviewSummary, vs)
	}
}

func TestReviewGapsRequiredToLeaveReview(t *testing.T) {
	tk := ticketWith("")
	tk.ReviewVerdict = "the diff matches the criteria"
	tk.ReviewSummary = "streamed the writer instead of buffering the whole file"
	// ReviewGaps deliberately left empty.

	vs := CheckAdvance(testEnv(t), tk, Request{To: "signoff"})
	if !hasFieldViolation(vs, CodeMissingProduces, ticket.FieldReviewGaps) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingProduces, ticket.FieldReviewGaps, vs)
	}
}

// "none" is the reviewer saying they looked and found nothing; it must count
// as filled, not be mistaken for the empty string.
func TestReviewGapsNoneCountsAsFilled(t *testing.T) {
	tk := ticketWith("")
	tk.ReviewVerdict = "the diff matches the criteria"
	tk.ReviewSummary = "streamed the writer instead of buffering the whole file"
	tk.ReviewGaps = "none"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "signoff"})
	if hasFieldViolation(vs, CodeMissingProduces, ticket.FieldReviewGaps) {
		t.Errorf("review-gaps: none was treated as empty: %v", vs)
	}
}

// A review with no check cannot leave either: the check is the one field the
// reader acts on, so it is enforced like the rest of the lane's output.
func TestReviewCheckRequiredToLeaveReview(t *testing.T) {
	tk := ticketWith("")
	tk.ReviewVerdict = "the diff matches the criteria"
	tk.ReviewSummary = "streamed the writer instead of buffering the whole file"
	tk.ReviewGaps = "none"
	// ReviewCheck deliberately left empty.

	vs := CheckAdvance(testEnv(t), tk, Request{To: "signoff"})
	if !hasFieldViolation(vs, CodeMissingProduces, ticket.FieldReviewCheck) {
		t.Errorf("expected a %s violation naming %s, got %v", CodeMissingProduces, ticket.FieldReviewCheck, vs)
	}
}

// Every field filled must produce no missing-produces violation at all.
func TestReviewAllFieldsFilledDoesNotBlock(t *testing.T) {
	tk := ticketWith("")
	tk.ReviewVerdict = "the diff matches the criteria"
	tk.ReviewSummary = "streamed the writer instead of buffering the whole file"
	tk.ReviewGaps = "none"
	tk.ReviewCheck = "1. go run ./cmd/app  2. open /export  3. the download starts at once"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "signoff"})
	for _, v := range vs {
		if v.Code == CodeMissingProduces {
			t.Errorf("unexpected missing-produces violation: %v", v)
		}
	}
}

// The review lane's contract is the source of the enforcement above; if it
// stops declaring every field, the tests above would stop meaning anything.
func TestReviewLaneDeclaresEveryReviewField(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	review, ok := lanes.Get("review")
	if !ok {
		t.Fatal("no review lane installed")
	}
	for _, f := range []string{
		ticket.FieldReviewSummary, ticket.FieldReviewGaps,
		ticket.FieldReviewVerdict, ticket.FieldReviewCheck,
	} {
		found := false
		for _, p := range review.OutputProduces {
			if p == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("review lane output-produces %v does not list %s", review.OutputProduces, f)
		}
	}
}

func hasFieldViolation(vs Violations, code, field string) bool {
	for _, v := range vs {
		if v.Code == code && v.Field == field {
			return true
		}
	}
	return false
}
