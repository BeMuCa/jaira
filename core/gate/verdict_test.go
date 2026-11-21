package gate

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// reviewed is a ticket in review with every output filled, so the only thing
// a move out of it can trip on is the verdict.
func reviewed(verdict string) *ticket.Ticket {
	tk := ticketWith("")
	tk.ReviewSummary = "s"
	tk.ReviewGaps = "none"
	tk.ReviewCheck = "c"
	tk.ReviewVerdict = verdict
	return tk
}

func verdictViolations(vs Violations) Violations {
	var out Violations
	for _, v := range vs {
		if v.Code == CodeVerdictSendBack || v.Code == CodeVerdictUnknown {
			out = append(out, v)
		}
	}
	return out
}

// A review that sent the work back cannot see it land in front of a person:
// the move on is refused, and the refusal names the verdict and where the
// ticket goes instead.
func TestSendBackVerdictRefusesSignoff(t *testing.T) {
	vs := CheckAdvance(testEnv(t), reviewed("send-back: DoD 2 is unmet"), Request{To: "signoff"})
	got := verdictViolations(vs)
	if len(got) != 1 || got[0].Code != CodeVerdictSendBack || got[0].Field != ticket.FieldReviewVerdict {
		t.Fatalf("want one %s violation, got %v", CodeVerdictSendBack, vs)
	}
	for _, want := range []string{`"send-back"`, `"in-progress"`, `"signoff"`} {
		if !strings.Contains(got[0].Message, want) {
			t.Errorf("message does not name %s: %s", want, got[0].Message)
		}
	}
}

// The way the verdict points is free: back to the work lane, without --force,
// and parking, which accepts nothing either.
func TestSendBackVerdictLetsTheTicketGoBack(t *testing.T) {
	for _, to := range []string{"in-progress", "testing"} {
		if vs := CheckAdvance(testEnv(t), reviewed("send-back: DoD 2 is unmet"), Request{To: to, Actor: "a"}); len(vs) != 0 {
			t.Errorf("review → %s with a send-back verdict refused: %v", to, vs)
		}
	}
	tk := reviewed("send-back: waiting on the vendor")
	if got := verdictViolations(CheckAdvance(testEnv(t), tk, Request{To: "blocked", Reason: "vendor"})); len(got) != 0 {
		t.Errorf("parking a sent-back ticket refused on its verdict: %v", got)
	}
}

// A move made in the board is a person deciding; the rail is for agents.
func TestSendBackVerdictInteractiveIsAPerson(t *testing.T) {
	vs := CheckAdvance(testEnv(t), reviewed("send-back: x"), Request{To: "signoff", Interactive: true})
	if got := verdictViolations(vs); len(got) != 0 {
		t.Errorf("a person's move in the board was refused on the verdict: %v", got)
	}
}

// The verdict is read off its first word, over the whole class of ways it is
// written; a sentence that only mentions a word is no verdict.
func TestVerdictWordClasses(t *testing.T) {
	cases := []struct {
		verdict string
		want    string // "" no verdict violation, else its code
	}{
		{"accept: the diff matches", ""},
		{"Accept", ""},
		{"Accept.", ""},
		{"send back!", CodeVerdictSendBack},
		{"  ACCEPT  : fine", ""},
		{"accept\nsecond line: with a colon", ""},
		{"send-back: DoD 2", CodeVerdictSendBack},
		{"send back", CodeVerdictSendBack},
		{"Send Back: x", CodeVerdictSendBack},
		{"send  -  back: x", CodeVerdictSendBack},
		{"send-back\r\nmore", CodeVerdictSendBack},
		{"\tsend-back", CodeVerdictSendBack},
		{"the diff matches the criteria; no defects found", CodeVerdictUnknown},
		{"I would not accept this", CodeVerdictUnknown},
		{"zurückschicken: DoD 1", CodeVerdictUnknown},
		{": empty word", CodeVerdictUnknown},
		{"accepted: x", CodeVerdictUnknown},
	}
	for _, c := range cases {
		got := verdictViolations(CheckAdvance(testEnv(t), reviewed(c.verdict), Request{To: "signoff"}))
		switch {
		case c.want == "" && len(got) != 0:
			t.Errorf("%q: refused, want accepted: %v", c.verdict, got)
		case c.want != "" && (len(got) != 1 || got[0].Code != c.want):
			t.Errorf("%q: got %v, want %s", c.verdict, got, c.want)
		}
	}

	// Empty is the lane's output contract's job, not this rule's: one
	// refusal, missing_lane_output, and no second one on top.
	vs := CheckAdvance(testEnv(t), reviewed(""), Request{To: "signoff"})
	if got := verdictViolations(vs); len(got) != 0 || !hasFieldViolation(vs, CodeMissingProduces, ticket.FieldReviewVerdict) {
		t.Errorf("empty verdict: %v", vs)
	}
}

// A board whose review lane predates the words keeps today's behaviour: its
// verdict is free text, and nothing reads it.
func TestLaneWithoutVerdictWordsReadsNothing(t *testing.T) {
	env := testEnv(t) // loaded afresh: changing its review lane touches no other test
	l, ok := env.Lanes.Get("review")
	if !ok {
		t.Fatal("no review lane")
	}
	l.VerdictAccept, l.VerdictSendBack = nil, nil
	for _, v := range []string{"send-back: x", "the diff matches"} {
		if got := verdictViolations(CheckAdvance(env, reviewed(v), Request{To: "signoff"})); len(got) != 0 {
			t.Errorf("%q on a lane without words was read: %v", v, got)
		}
	}
}
