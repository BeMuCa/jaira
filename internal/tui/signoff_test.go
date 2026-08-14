package tui

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The sign-off lane stops for a person, and this is that person's screen: what
// was wrong, what the agent did, why, whether it holds, and the reviewer's own
// account of what shipped, what it found missing, and its verdict. Those
// questions are the whole reason the step exists, so they are laid out in that
// order rather than being left among the other twenty fields.
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
		ReviewSummary: "added a token-bucket rate limiter per client IP",
		ReviewGaps:    "none",
		ReviewVerdict: "the diff matches the criteria; no defects found",
	}
	m.detail = tk
	out := stripANSI(m.renderSignOff())

	for _, want := range []string{
		"stop credential stuffing",                        // the issue
		"added a token bucket per IP",                     // what the agent did
		"the endpoint had no limit at all",                // why
		"429 now returned above 100/min",                  // whether it solved it
		"What the reviewer says it does",                  // the reviewer's summary, labelled
		"added a token-bucket rate limiter per client IP", // the reviewer's own account
		"What the reviewer found missing",                 // the reviewer's gaps, labelled
		"the diff matches the criteria",                   // the reviewer's own verdict
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
	lanes, err := lane.Load("")
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

// The sign-off screen is exactly where a ticked box without its evidence is a
// claim nobody can check, so a proof must render under its item here too.
func TestSignOffRendersProofUnderDoDItem(t *testing.T) {
	m := newTestModel(t, 120, 40)
	body := "## Definition of Done\n\n- [x] 429 returned above 100/min\n  proof: internal/x.go:12; TestRateLimit\n"
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "signoff",
		Body: body,
	}
	tk.DoDItems = ticket.ParseDoDItems(body)
	m.detail = tk
	out := stripANSI(m.renderSignOff())
	if !strings.Contains(out, "proof: internal/x.go:12; TestRateLimit") {
		t.Errorf("proof did not render under its item on the sign-off screen:\n%s", out)
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

// A follow-up outlives the ticket it came from: the predecessor gets archived
// off the board, and then follows: points at a file that is no longer in
// tickets/. The commits are therefore written into the context prose as well,
// so "what was already done" is still answerable from the follow-up alone.
func TestFollowUpContextCarriesTheCommits(t *testing.T) {
	src := &ticket.Ticket{
		ID:      "01KZZR4CBGDM5T35SZDR72PQYG",
		Title:   "Rate limit login",
		Context: "came up while reading the auth logs",
		Commits: []string{"bc615031d54de4b369d66bfabe0adf8846adc409"},
	}

	got := followUpContext(src)

	if !strings.Contains(got, "72PQYG") {
		t.Errorf("context does not name the predecessor:\n%s", got)
	}
	if !strings.Contains(got, "bc615031d54de4b369d66bfabe0adf8846adc409") {
		t.Errorf("context does not name the commit the work shipped in:\n%s", got)
	}
}

// With no commits recorded there is nothing to say, and a dangling "shipped in
// ." sentence would be worse than silence.
func TestFollowUpContextWithoutCommits(t *testing.T) {
	src := &ticket.Ticket{ID: "01KZZR4CBGDM5T35SZDR72PQYG", Title: "t"}

	if got := followUpContext(src); strings.Contains(got, "shipped in") {
		t.Errorf("context invented a commit sentence with no commits:\n%s", got)
	}
}

// The board is where the link is read, so the detail pane has to show it. The
// handle is what every command takes, so it is the form shown.
func TestDetailRendersFollows(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.detail = &ticket.Ticket{
		ID:      "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:   "Follow-up: rate limit login",
		Status:  "backlog",
		Follows: "01KZZR4CBGDM5T35SZDR72PQYG",
	}

	out := stripANSI(m.renderDetail())

	if !strings.Contains(out, "follows") || !strings.Contains(out, "72PQYG") {
		t.Errorf("detail pane does not show the predecessor:\n%s", out)
	}
}
