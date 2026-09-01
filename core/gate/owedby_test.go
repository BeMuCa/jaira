package gate

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// OwedBy is what a renderer asks in order to print a field nobody has filled
// in yet, so it has to agree with OutputOwed lane by lane: a screen naming a
// debt the gate would not refuse on, or missing one it would, is worse than
// no screen at all.
func TestOwedByAgreesWithOutputOwed(t *testing.T) {
	env := testEnv(t)
	tk := ticketWith("")
	// The outcome fields are filled by ticketWith; the four review fields are
	// not, which is the ticket that reached review unworked.
	owed := OwedBy(env.Lanes, tk)

	// Which lane owes a review field depends on the board — this catalogue
	// splits the four between critique, optimize and review — so what is
	// asserted is that every one of them is owed by a lane, not which.
	for _, f := range []string{
		ticket.FieldReviewSummary, ticket.FieldReviewGaps,
		ticket.FieldReviewVerdict, ticket.FieldReviewCheck,
	} {
		if owed[f] == "" {
			t.Errorf("%s is empty and no lane owes it", f)
		}
	}
	for _, f := range []string{ticket.FieldOutcomeWhat, ticket.FieldOutcomeWhy, ticket.FieldOutcomeResolves} {
		if l, ok := owed[f]; ok {
			t.Errorf("%s is filled in and was still reported as owed by %q", f, l)
		}
	}
	for f, id := range owed {
		l, ok := env.Lanes.Get(id)
		if !ok {
			t.Fatalf("%s is owed by %q, which is not an installed lane", f, id)
		}
		var found bool
		for _, o := range OutputOwed(l, tk) {
			found = found || o == f
		}
		if !found {
			t.Errorf("OwedBy attributes %s to %q, but that lane does not owe it", f, id)
		}
	}
}

// pre-process declares plan and brainstorm declares goal, and both are entered
// only by a ticket that opted in. A ticket that did not opt in owes them
// nothing — the same reason the gate lets it move straight past them.
func TestOwedByIgnoresLanesOffTheRoute(t *testing.T) {
	env := testEnv(t)
	tk := ticketWith("")
	tk.Goal = ""

	if l, ok := OwedBy(env.Lanes, tk)["plan"]; ok {
		t.Errorf("a ticket that did not opt into planning owes plan to %q", l)
	}
	if l, ok := OwedBy(env.Lanes, tk)[ticket.FieldGoal]; ok {
		t.Errorf("a ticket that did not opt into brainstorming owes goal to %q", l)
	}

	tk.Body = "## Options\n\n- [x] brainstorm\n"
	tk.Options = ticket.ParseOptions(tk.Body)
	if l := OwedBy(env.Lanes, tk)[ticket.FieldGoal]; l != "brainstorm" {
		t.Errorf("a ticket that opted into brainstorming owes goal to %q, want brainstorm", l)
	}
}

// Two lanes may declare the same field; the earlier one owes it, the same
// precedence a lane-load warning uses when it names a producer.
func TestOwedByAttributesAFieldToTheFirstLaneThatDeclaresIt(t *testing.T) {
	set := &lane.Set{Lanes: []*lane.Lane{
		{ID: "critique", OutputProduces: []string{ticket.FieldReviewSummary}},
		{ID: "review", OutputProduces: []string{ticket.FieldReviewSummary}},
	}}
	if l := OwedBy(set, &ticket.Ticket{ID: "x"})[ticket.FieldReviewSummary]; l != "critique" {
		t.Errorf("review-summary is owed by %q, want the first lane that declares it", l)
	}
}

// Called from the renderers, which draw whatever the board holds.
func TestOwedByIsNilSafe(t *testing.T) {
	if owed := OwedBy(nil, ticketWith("")); owed != nil {
		t.Errorf("a nil lane set owed %v", owed)
	}
	if owed := OwedBy(testEnv(t).Lanes, nil); owed != nil {
		t.Errorf("a nil ticket owed %v", owed)
	}
}

// DeclaredBy is the label a renderer puts on a filled value, so it must use
// the same first-declarer attribution OwedBy gives a debt.
func TestDeclaredByNamesTheFirstDeclarer(t *testing.T) {
	set := &lane.Set{Lanes: []*lane.Lane{
		{ID: "critique", OutputProduces: []string{ticket.FieldReviewSummary}},
		{ID: "review", OutputProduces: []string{ticket.FieldReviewSummary, ticket.FieldReviewCheck}},
	}}
	src := DeclaredBy(set)
	// Both declarers, in board order: picking one would assert a provenance
	// the frontmatter cannot back.
	if got := strings.Join(src[ticket.FieldReviewSummary], "/"); got != "critique/review" {
		t.Errorf("review-summary declared by %q, want critique/review", got)
	}
	if got := strings.Join(src[ticket.FieldReviewCheck], "/"); got != "review" {
		t.Errorf("review-check declared by %q, want review", got)
	}
	if DeclaredBy(nil) != nil {
		t.Error("a nil set declared something")
	}
}
