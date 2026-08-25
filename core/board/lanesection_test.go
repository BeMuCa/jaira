package board

import (
	"strings"
	"testing"
)

func demoFacts() []LaneFact {
	return []LaneFact{
		{ID: "todo", Name: "Todo", Description: "Specified and ready to be picked up."},
		{ID: "in-progress", Name: "Implementing", Description: "Carrying out the plan.",
			Agentic: true, ModelTier: "cheap", Produces: []string{"outcome-what"}},
		{ID: "critique", Name: "Critique", Description: "Judges the approach. Sends work back.",
			Agentic: true, ModelTier: "strong", RejectsTo: "in-progress", Produces: []string{"review-summary"}},
		{ID: "signoff", Name: "Human Review", Description: "Waiting for a person.", HumanExit: true},
		{ID: "done", Name: "Done", Description: "Accepted.", Terminal: true},
	}
}

// The static half of the note teaches the commands and cannot teach the board:
// every board's lanes are different, and an agent that has to discover the route
// by running a command will often not run it. Twelve tickets sat in a critique
// lane with nothing ever produced for them.
func TestLaneSectionNamesTheRouteAndTheLoop(t *testing.T) {
	got := laneSection(demoFacts())

	if !strings.Contains(got, "todo → in-progress → critique → signoff → done") {
		t.Errorf("the route is not stated on one line:\n%s", got)
	}
	if !strings.Contains(got, "critique sends work back to in-progress") {
		t.Errorf("the declared back edge is not named — it is the part an ordered list cannot show:\n%s", got)
	}
	for _, want := range []string{"tier cheap", "tier strong", "must produce review-summary", "terminal"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
}

// Which lanes an agent may work at all is the fact it acts on, so it is stated
// per lane rather than left to be inferred from a model tier.
func TestLaneSectionSaysWhoseEachLaneIs(t *testing.T) {
	got := laneSection(demoFacts())
	for _, want := range []string{
		"`in-progress` — yours to work",
		"`signoff` — **a person's, not yours**",
		"`todo` — no agent step",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
}

// The board does not run anything, and the note has to say so plainly: this is
// where a reader would otherwise assume a lane's prompt fires by itself.
func TestLaneSectionSaysNothingRunsOnItsOwn(t *testing.T) {
	got := laneSection(demoFacts())
	if !strings.Contains(got, "Nothing moves a ticket for you") || !strings.Contains(got, "no daemon and no runner") {
		t.Errorf("the note does not say the board runs nothing:\n%s", got)
	}
	if !strings.Contains(got, "jaira next --per-lane") {
		t.Errorf("the note does not say how to find the lane with work waiting:\n%s", got)
	}
}

// A caller with no board loaded writes exactly the note that existed before
// boards could describe themselves.
func TestLaneSectionIsEmptyWithoutLanes(t *testing.T) {
	if got := laneSection(nil); got != "" {
		t.Errorf("laneSection(nil) = %q, want empty", got)
	}
	if got := laneSection([]LaneFact{}); got != "" {
		t.Errorf("laneSection([]) = %q, want empty", got)
	}
}

// A lane whose description runs to a paragraph must not push the rest of the
// section off the screen.
func TestLaneSectionKeepsDescriptionsToTheirFirstClaim(t *testing.T) {
	got := laneSection([]LaneFact{{ID: "x", Description: "First claim. Second sentence nobody needs here."}})
	if !strings.Contains(got, "First claim.") {
		t.Errorf("the opening claim was dropped:\n%s", got)
	}
	if strings.Contains(got, "Second sentence") {
		t.Errorf("the whole description was inlined:\n%s", got)
	}
}
