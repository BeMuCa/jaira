package lane

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// max-rounds is read as a whole number above zero. Anything else written on
// the line is refused rather than read as "no limit": a typo must not switch
// off the brake it was written for.
func TestParseMaxRounds(t *testing.T) {
	cases := []struct {
		line    string
		want    int
		wantErr bool
	}{
		{"", 0, false},
		{"max-rounds: 3", 3, false},
		{"max-rounds: 1", 1, false},
		{`max-rounds: "2"`, 2, false},
		{"max-rounds: 0", 0, true},
		{"max-rounds: -1", 0, true},
		{"max-rounds: three", 0, true},
		{"max-rounds: 2.5", 0, true},
		{"max-rounds:", 0, true},
		{`max-rounds: ""`, 0, true},
		{"max-rounds: [3]", 0, true},
		{"max-rounds: 99999999999999999999", 0, true},
	}
	for _, c := range cases {
		t.Run(c.line, func(t *testing.T) {
			src := "---\nid: loop\nname: Loop\n" + c.line + "\n---\n"
			l, err := Parse([]byte(src), "loop.md")
			if c.wantErr {
				if err == nil {
					t.Fatalf("accepted, MaxRounds=%d", l.MaxRounds)
				}
				if !strings.Contains(err.Error(), "max-rounds") {
					t.Errorf("error does not name the field: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if l.MaxRounds != c.want {
				t.Errorf("MaxRounds = %d, want %d", l.MaxRounds, c.want)
			}
		})
	}
}

// The shipped loop lanes carry the limit the worker role already stops at.
func TestShippedLoopLanesStopAfterThree(t *testing.T) {
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"critique", "optimize", "testing"} {
		l, ok := set.Get(id)
		if !ok {
			t.Fatalf("no %s lane", id)
		}
		if l.MaxRounds != 3 {
			t.Errorf("%s: max-rounds %d, want 3", id, l.MaxRounds)
		}
	}
}

// A copy of a shipped lane that only drops max-rounds behaves differently and
// is marked as an override, like any other changed line.
func TestMaxRoundsCountsForOverride(t *testing.T) {
	raw := string(builtinBytes(t, "testing"))
	changed := strings.Replace(raw, "max-rounds: 3\n", "", 1)
	if changed == raw {
		t.Fatal("fixture did not drop max-rounds")
	}
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "testing.md", changed)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, _ := set.Get("testing")
	if l.MaxRounds != 0 || l.Overrides != "testing" {
		t.Errorf("MaxRounds=%d Overrides=%q, want 0 and testing", l.MaxRounds, l.Overrides)
	}
}

// Rounds left read back-count only when the run is this lane's own.
func TestRoundsLeft(t *testing.T) {
	l := &Lane{ID: "testing", MaxRounds: 3}
	cases := []struct {
		from, by string
		count    int
		want     int
	}{
		{"", "", 0, 3},
		{"testing", ticket.BackByLane, 1, 2},
		{"testing", ticket.BackByLane, 3, 0},
		{"testing", ticket.BackByLane, 5, 0},
		{"testing", ticket.BackByPerson, 2, 3},
		{"critique", ticket.BackByLane, 2, 3},
	}
	for _, c := range cases {
		tk := &ticket.Ticket{BackFrom: c.from, BackBy: c.by, BackCount: c.count}
		if got, limited := l.RoundsLeft(tk); !limited || got != c.want {
			t.Errorf("%+v: got %d (limited %v), want %d", c, got, limited, c.want)
		}
	}
	if _, limited := (&Lane{ID: "review"}).RoundsLeft(&ticket.Ticket{}); limited {
		t.Error("a lane without max-rounds reports a limit")
	}
}

// The person's lane is one that takes a question, never a sign-off lane.
func TestPersonLane(t *testing.T) {
	human := &Lane{ID: "human", RequiresQuestion: true}
	ask := &Lane{ID: "ask", RequiresQuestion: true}
	signoff := &Lane{ID: "signoff", RequiresHumanExit: true}
	loop := &Lane{ID: "testing", RejectsTo: []string{"in-progress", "ask"}}
	set := assemble([]*Lane{loop, human, ask, signoff}, nil, nil)
	if got := set.PersonLane(loop); got != "ask" {
		t.Errorf("rejects-to names ask, got %q", got)
	}
	if got := set.PersonLane(&Lane{ID: "critique"}); got != "human" {
		t.Errorf("no back edge to a question lane: got %q, want the board's first, human", got)
	}
	bare := assemble([]*Lane{loop, signoff}, nil, nil)
	if got := bare.PersonLane(loop); got != "" {
		t.Errorf("a board with only a sign-off lane offered %q", got)
	}
}

// An agentic lane that requires a question is not where a person answers.
func TestPersonLaneSkipsAgenticQuestionLanes(t *testing.T) {
	work := &Lane{ID: "in-progress", Agentic: true, RequiresQuestion: true}
	human := &Lane{ID: "human", RequiresQuestion: true}
	loop := &Lane{ID: "testing", RejectsTo: []string{"in-progress"}}
	set := assemble([]*Lane{work, loop, human}, nil, nil)
	if got := set.PersonLane(loop); got != "human" {
		t.Errorf("got %q, want human", got)
	}
}

// The third return from one lane goes to the teamlead first: with a parent
// session the ticket stays in its lane with the question set, and only
// without one does the shipped prompt still send it to the human lane.
func TestInProgressPromptHandsTheThirdReturnToTheTeamleadFirst(t *testing.T) {
	ls, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	var prompt string
	for _, l := range ls {
		if l.ID == "in-progress" {
			prompt = l.Prompt
		}
	}
	for _, want := range []string{
		"to the teamlead first, a person only",
		"leave\n  the ticket in this lane, set the question",
		`jaira set <id> question="testing sent this back three times`,
		"No parent session, or it cannot be reached",
		`jaira move <id> --to human --question "testing sent this back three times`,
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("in-progress prompt does not say %q", want)
		}
	}
}
