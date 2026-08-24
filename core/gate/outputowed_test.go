package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// OutputOwed is what the board asks in order to show that a lane has not been
// run on a ticket, and it must be the same question the gate asks when it
// refuses the move out of that lane. Two implementations of "the lane owes
// this field" would drift, and the drift would show as a card that looks worked
// and a move that is refused.
func TestOutputOwedAgreesWithTheRefusal(t *testing.T) {
	env := testEnv(t)
	l, ok := env.Lanes.Get("review")
	if !ok {
		t.Fatal("no review lane")
	}

	tk := ticketWith("")
	tk.ReviewSummary = "streamed the writer instead of buffering the whole file"
	// verdict, gaps and check deliberately left empty.

	owed := OutputOwed(l, tk)
	if len(owed) == 0 {
		t.Fatal("review declares four outputs and three are empty; OutputOwed reported none")
	}
	for _, f := range owed {
		if !hasFieldViolation(CheckAdvance(env, tk, Request{To: "signoff"}), CodeMissingProduces, f) {
			t.Errorf("OutputOwed reported %q but the gate does not refuse on it", f)
		}
	}
	for _, f := range owed {
		if f == ticket.FieldReviewSummary {
			t.Error("review-summary is filled and was still reported as owed")
		}
	}
}

// Nothing owed is the signal that the lane has been run, so a fully produced
// contract must come back empty rather than "almost empty".
func TestOutputOwedEmptyWhenTheLaneHasRun(t *testing.T) {
	env := testEnv(t)
	l, _ := env.Lanes.Get("review")

	tk := ticketWith("")
	tk.ReviewSummary = "s"
	tk.ReviewGaps = "none"
	tk.ReviewVerdict = "v"
	tk.ReviewCheck = "run it and watch the log"

	if owed := OutputOwed(l, tk); len(owed) != 0 {
		t.Errorf("a produced review still owes %v", owed)
	}
}

// A lane the ticket opted out of owes nothing — the same reason the gate does
// not refuse the move: the step does not apply to this ticket. Without this a
// skipped lane would mark every card passing through it as unworked forever.
func TestOutputOwedIgnoresASkippedLane(t *testing.T) {
	l := &lane.Lane{
		ID:             "critique",
		RequiresOption: "critique",
		OutputProduces: []string{ticket.FieldReviewSummary},
	}
	tk := ticketWith("")
	if owed := OutputOwed(l, tk); len(owed) != 0 {
		t.Errorf("a lane this ticket opted out of owes %v", owed)
	}
}

// Called from the renderers, which run against whatever the board holds, so the
// nil cases must not be a crash. A lane declaring an output that this package
// does not model once crashed on exactly this path.
func TestOutputOwedIsNilSafe(t *testing.T) {
	if owed := OutputOwed(nil, ticketWith("")); owed != nil {
		t.Errorf("nil lane owed %v", owed)
	}
	if owed := OutputOwed(&lane.Lane{ID: "x", OutputProduces: []string{"whatever"}}, nil); owed != nil {
		t.Errorf("nil ticket owed %v", owed)
	}
}
