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
			Agentic: true, ModelTier: "strong", RejectsTo: []string{"in-progress"}, Produces: []string{"review-summary"}},
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

// A lane file edited by hand changes the pipeline without going through any
// command that would regenerate the note — which is how a board ends up handing
// an agent a route that no longer exists.
func TestNoteIsCurrentSpotsAHandEditedBoard(t *testing.T) {
	root := t.TempDir()
	facts := demoFacts()
	if _, err := AnnounceInAgentFiles(root, facts); err != nil {
		t.Fatal(err)
	}
	if current, stale := NoteIsCurrent(root, facts); !current {
		t.Errorf("a note just written reports as stale: %v", stale)
	}

	// The board gains a lane without the note being rewritten.
	changed := append(facts, LaneFact{ID: "optimize", Agentic: true, Description: "Strips what the change does not need."})
	current, stale := NoteIsCurrent(root, changed)
	if current {
		t.Fatal("a board that gained a lane still reports its note as current")
	}
	if len(stale) != 2 {
		t.Errorf("stale = %v, want both agent files", stale)
	}

	// And regenerating settles it.
	if _, err := AnnounceInAgentFiles(root, changed); err != nil {
		t.Fatal(err)
	}
	if current, stale := NoteIsCurrent(root, changed); !current {
		t.Errorf("still stale after regenerating: %v", stale)
	}
}

// A repository with no agent file is not out of date; nothing should nag about
// a file nobody wrote.
func TestNoteIsCurrentIgnoresMissingFiles(t *testing.T) {
	if current, stale := NoteIsCurrent(t.TempDir(), demoFacts()); !current {
		t.Errorf("an empty repository reports a stale note: %v", stale)
	}
}

// A lane may hand work back to more than one place — a flaw goes back to be
// implemented, a decision goes to a person — and the note has to name both, or
// the second edge exists only in the lane file nobody generating this read.
func TestLaneSectionNamesEveryBackEdge(t *testing.T) {
	facts := demoFacts()
	for i := range facts {
		if facts[i].ID == "critique" {
			facts[i].RejectsTo = []string{"in-progress", "signoff"}
		}
	}
	got := laneSection(facts)
	if !strings.Contains(got, "critique sends work back to in-progress or signoff") {
		t.Errorf("both declared back edges are not named:\n%s", got)
	}
}
