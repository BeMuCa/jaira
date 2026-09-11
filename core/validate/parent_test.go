package validate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// parented builds a ticket that is part of another.
func parented(id, parent string) *ticket.Ticket {
	t := ok1(id)
	t.Parent = parent
	return t
}

// ok1 is a ticket with nothing wrong with it beyond what a test adds.
func ok1(id string) *ticket.Ticket {
	return &ticket.Ticket{
		ID:       id,
		Title:    "t",
		Status:   "backlog",
		Goal:     "g",
		Context:  "c",
		DoD:      "d",
		Assignee: "a",
		Path:     "/tmp/t.md",
	}
}

const (
	idA = "01KZTT3XZ2YQBX93TTSR7BVRC1"
	idB = "01KZTT3XZ2YQBX93TTSR7BVRC2"
	idC = "01KZTT3XZ2YQBX93TTSR7BVRC3"
)

// The parent chain is followed recursively wherever children are rendered, so
// a ring in it hangs the view rather than merely misdrawing it.
func TestParentCycleIsAnError(t *testing.T) {
	a := parented(idA, idB)
	b := parented(idB, idA)

	ps := Tickets([]*ticket.Ticket{a, b}, lanes(t), nil)
	if !has(ps, CodeParentCycle) {
		t.Errorf("a ↔ b is a parent ring and must be reported, got %v", ps)
	}
}

func TestLongerParentCycleIsAnError(t *testing.T) {
	ps := Tickets([]*ticket.Ticket{
		parented(idA, idB), parented(idB, idC), parented(idC, idA),
	}, lanes(t), nil)
	if !has(ps, CodeParentCycle) {
		t.Errorf("a → b → c → a must be reported, got %v", ps)
	}
}

func TestSelfParentIsAnError(t *testing.T) {
	ps := Tickets([]*ticket.Ticket{parented(idA, idA)}, lanes(t), nil)
	if !has(ps, CodeSelfParent) {
		t.Errorf("a ticket cannot be its own parent, got %v", ps)
	}
}

// A parent that finished and was filed away still exists — known says so, and
// nothing is wrong with the link.
func TestFiledAwayParentIsNotDangling(t *testing.T) {
	known := func(id string) bool { return id == idB }

	ps := Tickets([]*ticket.Ticket{parented(idA, idB)}, lanes(t), known)
	if has(ps, CodeDanglingParent) {
		t.Errorf("a parent in the logbook is not dangling, got %v", ps)
	}
}

func TestParentThatExistsNowhereIsDangling(t *testing.T) {
	ps := Tickets([]*ticket.Ticket{parented(idA, idB)}, lanes(t), nil)
	if !has(ps, CodeDanglingParent) {
		t.Errorf("a parent nothing knows must be reported, got %v", ps)
	}
}

// The same rule the gate now applies: a dependency that finished and left the
// board has cleared, and calling that an error made finishing work look like
// breaking it.
func TestFiledAwayBlockerIsNotDangling(t *testing.T) {
	a := ok1(idA)
	a.BlockedBy = []string{idB}
	known := func(id string) bool { return id == idB }

	ps := Tickets([]*ticket.Ticket{a}, lanes(t), known)
	if has(ps, CodeDanglingDep) {
		t.Errorf("a blocker in the logbook is not dangling, got %v", ps)
	}
}

// A chain that simply ends is not a ring.
func TestPlainParentChainIsFine(t *testing.T) {
	ps := Tickets([]*ticket.Ticket{
		parented(idA, idB), parented(idB, idC), ok1(idC),
	}, lanes(t), nil)
	if has(ps, CodeParentCycle) || has(ps, CodeDanglingParent) {
		t.Errorf("a → b → c is a normal chain, got %v", ps)
	}
}

// A child explaining itself names its epic — that is what the context is
// for. Telling it to put the epic in blocked-by would contradict the rule
// that a parent is not a gate, and would make the warning fire on every
// well-written child.
func TestAParentNamedInTheContextIsNotAnUndeclaredDependency(t *testing.T) {
	epic := ok1(idB)
	child := parented(idA, idB)
	child.Context = "part of " + handleOf(idB) + ", which sets the direction"

	ps := Tickets([]*ticket.Ticket{epic, child}, lanes(t), nil)
	if has(ps, CodeUndeclaredDep) {
		t.Errorf("naming the parent is not a hidden dependency, got %v", ps)
	}
}

// Same for a ticket this one merely relates to.
func TestARelatedTicketNamedInTheContextIsNotADependency(t *testing.T) {
	other := ok1(idB)
	a := ok1(idA)
	a.Related = []string{idB}
	a.Context = "see also " + handleOf(idB)

	ps := Tickets([]*ticket.Ticket{other, a}, lanes(t), nil)
	if has(ps, CodeUndeclaredDep) {
		t.Errorf("naming a related ticket is not a hidden dependency, got %v", ps)
	}
}

// A handle is six characters, so two different ids can end with the same
// one. A parent written as a full id that exists nowhere is dangling — it
// must not be resolved to whichever ticket shares that tail, or validate
// reports a ring that does not exist.
func TestAHandleCollisionDoesNotInventACycle(t *testing.T) {
	// idC and ghost share their last six characters.
	ghost := "01ZZZZZZZZZZZZZZZZZZZ" + idC[len(idC)-5:]
	a := parented(idA, ghost)
	c := parented(idC, idA)

	ps := Tickets([]*ticket.Ticket{a, c}, lanes(t), nil)
	if has(ps, CodeParentCycle) {
		t.Errorf("a full id that names nothing is dangling, not a ring: %v", ps)
	}
	if !has(ps, CodeDanglingParent) {
		t.Errorf("and it must be reported as dangling, got %v", ps)
	}
}
