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
