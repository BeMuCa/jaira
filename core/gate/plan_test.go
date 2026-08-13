package gate

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

func testEnv(t *testing.T) Env {
	t.Helper()
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatalf("lane.Load: %v", err)
	}
	return Env{Lanes: lanes}
}

func ticketWith(body string) *ticket.Ticket {
	tk := &ticket.Ticket{
		ID:       "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:    "t",
		Status:   "review",
		Goal:     "g",
		Context:  "c",
		Assignee: "a",
		Body:     body,
		Outcome:  ticket.Outcome{What: "w", Why: "y", Resolves: "r"},
	}
	tk.DoDItems = ticket.ParseDoDItems(body)
	tk.PlanItems = ticket.ParsePlanItems(body)
	return tk
}

// The plan is the method and the definition of done is the criteria, but they are
// not independent: the criteria cannot have been met while the work that meets
// them is still in progress. A ticket accepted with an unfinished plan is the
// "I thought it was built, only the design was done" failure this board exists to
// prevent, so it is refused rather than merely noted.
func TestUnfinishedPlanBlocksTerminalLane(t *testing.T) {
	tk := ticketWith(`## Plan

- [x] write the spec
- [~] implement

## Definition of Done

- [x] the criterion
`)
	vs := CheckAdvance(testEnv(t), tk, Request{To: "done"})
	var found *Violation
	for i := range vs {
		if vs[i].Code == CodePlanIncomplete {
			found = &vs[i]
		}
	}
	if found == nil {
		t.Fatalf("expected a plan violation, got %v", vs)
	}
	if !strings.Contains(found.Message, "implement") {
		t.Errorf("message %q does not name the unfinished step", found.Message)
	}
}

// A plan is optional. Tickets that do not carry one must be unaffected.
func TestAbsentPlanDoesNotBlock(t *testing.T) {
	tk := ticketWith("## Definition of Done\n\n- [x] the criterion\n")
	for _, v := range CheckAdvance(testEnv(t), tk, Request{To: "done"}) {
		if v.Code == CodePlanIncomplete {
			t.Errorf("a ticket with no plan was blocked: %v", v)
		}
	}
}

func TestCompletePlanDoesNotBlock(t *testing.T) {
	tk := ticketWith(`## Plan

- [x] write the spec
- [x] implement

## Definition of Done

- [x] the criterion
`)
	for _, v := range CheckAdvance(testEnv(t), tk, Request{To: "done"}) {
		if v.Code == CodePlanIncomplete {
			t.Errorf("a complete plan was blocked: %v", v)
		}
	}
}

// The plan governs acceptance, not movement through the pipeline. An agent must
// still be able to move a ticket into review while its plan is in progress.
func TestUnfinishedPlanDoesNotBlockNonTerminalLanes(t *testing.T) {
	tk := ticketWith("## Plan\n\n- [~] implement\n\n## Definition of Done\n\n- [ ] c\n")
	tk.Status = "in-progress"
	for _, v := range CheckAdvance(testEnv(t), tk, Request{To: "review"}) {
		if v.Code == CodePlanIncomplete {
			t.Errorf("moving to a non-terminal lane was blocked: %v", v)
		}
	}
}

// A ticket sitting in a lane this installation does not know about has no
// contract to enforce, so moving it out is refused rather than guessed at.
func TestMoveOutOfUnrecognizedLaneRefused(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "some-uninstalled-lane"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "todo"})
	found := false
	for _, v := range vs {
		if v.Code == CodeUnknownLane {
			found = true
			if !strings.Contains(v.Message, "some-uninstalled-lane") {
				t.Errorf("message %q does not name the unrecognized lane", v.Message)
			}
		}
	}
	if !found {
		t.Fatalf("expected a %s violation, got %v", CodeUnknownLane, vs)
	}
}

// Moving a ticket into an id no lane claims is refused with CodeNoSuchLane.
func TestMoveIntoNoSuchLaneRefused(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "backlog"
	vs := CheckAdvance(testEnv(t), tk, Request{To: "not-a-real-lane"})
	found := false
	for _, v := range vs {
		if v.Code == CodeNoSuchLane {
			found = true
			if !strings.Contains(v.Message, "not-a-real-lane") {
				t.Errorf("message %q does not name the target lane", v.Message)
			}
		}
	}
	if !found {
		t.Fatalf("expected a %s violation, got %v", CodeNoSuchLane, vs)
	}
}
