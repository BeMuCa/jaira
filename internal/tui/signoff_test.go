package tui

import (
	"strings"
	"testing"

	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/ticket"
)

// The review lane stops for a person, and this is that person's screen: what was
// wrong, what the agent did, why, and whether it holds. Those four questions are
// the whole reason the step exists, so they are laid out in that order rather
// than being left among the other twenty fields.
func TestSignOffViewShowsTheFourQuestions(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Rate limit login", Status: "signoff",
		Goal:    "stop credential stuffing",
		Context: "came up while reading the auth logs",
		Outcome: ticket.Outcome{
			What:     "added a token bucket per IP",
			Why:      "the endpoint had no limit at all",
			Resolves: "429 now returned above 100/min",
		},
		ReviewVerdict: "the diff matches the criteria; no defects found",
	}
	m.detail = tk
	out := stripANSI(m.renderSignOff())

	for _, want := range []string{
		"stop credential stuffing",         // the issue
		"added a token bucket per IP",      // what the agent did
		"the endpoint had no limit at all", // why
		"429 now returned above 100/min",   // whether it solved it
		"the diff matches the criteria",    // the reviewer's own verdict
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sign-off view is missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "accept") || !strings.Contains(out, "follow-up") {
		t.Errorf("the two actions are not offered:\n%s", out)
	}
}

// The implementer's account and the reviewer's judgement are different claims
// and must both survive. Storing the verdict in outcome-resolves would destroy
// exactly the pair this screen exists to show.
func TestVerdictIsSeparateFromTheOutcome(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.detail = &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "signoff",
		Outcome:       ticket.Outcome{What: "w", Why: "y", Resolves: "the implementer says it works"},
		ReviewVerdict: "the reviewer disagrees",
	}
	out := stripANSI(m.renderSignOff())
	if !strings.Contains(out, "the implementer says it works") {
		t.Error("the implementer's account was lost")
	}
	if !strings.Contains(out, "the reviewer disagrees") {
		t.Error("the reviewer's verdict was lost")
	}
}

// An agent must not be able to release a ticket from the lane that stops for a
// human; a person pressing a key in the board must. Review is now the model's
// own step and it may leave that freely — sign-off is the checkpoint.
func TestHumanExitGate(t *testing.T) {
	lanes, err := lane.Load()
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "signoff",
		Goal: "g", Context: "c", Assignee: "a",
		ReviewVerdict: "the diff holds",
		Outcome:       ticket.Outcome{What: "w", Why: "y", Resolves: "r"},
	}
	env := gate.Env{Lanes: lanes}

	blocked := gate.CheckAdvance(env, tk, gate.Request{To: "done"})
	if !hasCode(blocked, gate.CodeNeedsHuman) {
		t.Errorf("an agent was allowed out of the sign-off lane: %v", blocked)
	}

	allowed := gate.CheckAdvance(env, tk, gate.Request{To: "done", Interactive: true})
	if hasCode(allowed, gate.CodeNeedsHuman) {
		t.Errorf("a person was blocked from signing off: %v", allowed)
	}
}

func hasCode(vs gate.Violations, code string) bool {
	for _, v := range vs {
		if v.Code == code {
			return true
		}
	}
	return false
}
