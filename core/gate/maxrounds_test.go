package gate

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The shipped loop lanes stop after three send-backs in a row. The table
// covers who sent the ticket back last, how often, and who sends it now —
// every case the limit must refuse and every case it must leave alone.
func TestMaxRoundsRefusesOnlyALanesOwnRunPastTheLimit(t *testing.T) {
	cases := []struct {
		name       string
		lane       string
		backFrom   string
		backBy     string
		count      int
		to         string
		req        Request
		wantRefuse bool
	}{
		{"fourth send-back from testing", "testing", "testing", ticket.BackByLane, 3, "in-progress", Request{}, true},
		{"fourth from critique", "critique", "critique", ticket.BackByLane, 3, "in-progress", Request{}, true},
		{"fourth from optimize, to an earlier lane than in-progress", "optimize", "optimize", ticket.BackByLane, 3, "todo", Request{}, true},
		{"count above the limit, a hand-edited file", "testing", "testing", ticket.BackByLane, 7, "in-progress", Request{}, true},
		{"third send-back is still allowed", "testing", "testing", ticket.BackByLane, 2, "in-progress", Request{}, false},
		{"no send-back recorded yet", "testing", "", "", 0, "in-progress", Request{}, false},
		{"the last ones came from another lane", "testing", "critique", ticket.BackByLane, 5, "in-progress", Request{}, false},
		{"the last one was a person's", "testing", "testing", ticket.BackByPerson, 3, "in-progress", Request{}, false},
		{"a person sends it back from the board", "testing", "testing", ticket.BackByLane, 3, "in-progress", Request{Interactive: true}, false},
		{"a person sends it back at a terminal", "testing", "testing", ticket.BackByLane, 3, "in-progress", Request{ByPerson: true}, false},
		{"moving on is no send-back", "testing", "testing", ticket.BackByLane, 3, "review", Request{}, false},
		{"parking is no send-back", "testing", "testing", ticket.BackByLane, 3, "blocked", Request{Reason: "r"}, false},
		{"a lane without max-rounds", "review", "review", ticket.BackByLane, 9, "in-progress", Request{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tk := ticketWith("")
			tk.Status = c.lane
			tk.BackFrom, tk.BackBy, tk.BackCount = c.backFrom, c.backBy, c.count
			req := c.req
			req.To = c.to

			v := findViolation(CheckAdvance(testEnv(t), tk, req), CodeMaxRounds)
			if c.wantRefuse && v == nil {
				t.Fatalf("send-back from %s with %d in a row was allowed", c.lane, c.count)
			}
			if !c.wantRefuse && v != nil {
				t.Fatalf("refused: %s", v.Message)
			}
		})
	}
}

// The refusal names the lane a person answers in, so the agent knows where to
// take the ticket instead of guessing or reaching for --force.
func TestMaxRoundsRefusalNamesThePersonsLane(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "critique"
	tk.BackFrom, tk.BackBy, tk.BackCount = "critique", ticket.BackByLane, 3

	v := findViolation(CheckAdvance(testEnv(t), tk, Request{To: "in-progress"}), CodeMaxRounds)
	if v == nil {
		t.Fatal("not refused")
	}
	for _, want := range []string{`"human"`, "--question", "3 times"} {
		if !strings.Contains(v.Message, want) {
			t.Errorf("message lacks %s: %s", want, v.Message)
		}
	}
}

// The limit is the lane file's number, not a constant three.
func TestMaxRoundsFollowsTheLaneFile(t *testing.T) {
	env := testEnv(t)
	l, _ := env.Lanes.Get("testing")
	l.MaxRounds = 1

	tk := ticketWith("")
	tk.Status = "testing"
	tk.BackFrom, tk.BackBy, tk.BackCount = "testing", ticket.BackByLane, 1
	if findViolation(CheckAdvance(env, tk, Request{To: "in-progress"}), CodeMaxRounds) == nil {
		t.Error("max-rounds: 1 let a second send-back through")
	}
}

// The hand-off the refusal prescribes is never refused itself, wherever the
// board draws the question lane: earlier than the loop lane, moving there is a
// send-back too.
func TestMaxRoundsLetsTheHandOffToAQuestionLaneThrough(t *testing.T) {
	for _, earlier := range []string{"backlog", "todo", "in-progress"} {
		env := testEnv(t)
		q, ok := env.Lanes.Get(earlier)
		if !ok {
			t.Fatalf("no %s lane", earlier)
		}
		q.RequiresQuestion = true
		crit, _ := env.Lanes.Get("critique")
		crit.RejectsTo = []string{"in-progress", earlier}

		tk := ticketWith("")
		tk.Status = "critique"
		tk.BackFrom, tk.BackBy, tk.BackCount = "critique", ticket.BackByLane, 3
		vs := CheckAdvance(env, tk, Request{To: earlier, Question: "is the definition of done wrong?"})
		if v := findViolation(vs, CodeMaxRounds); v != nil {
			t.Errorf("hand-off to question lane %s refused: %s", earlier, v.Message)
		}
	}
}
