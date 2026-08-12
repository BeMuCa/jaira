package cli

import (
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// An agent driving the board through --json must be able to see the checklists.
// Without them it cannot tell which items exist, which are outstanding, or which
// number to pass to 'jaira dod' — and the gate's own refusal messages tell it to
// do exactly that. Rendering the checklists only in the human output would leave
// the machine interface unable to act on the machine-readable error.
func TestTicketJSONCarriesChecklists(t *testing.T) {
	lanes, err := lane.Load()
	if err != nil {
		t.Fatal(err)
	}
	body := `## Plan

- [x] write the spec
- [~] implement

## Definition of Done

- [ ] documented
`
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo", Body: body}
	tk.DoDItems = ticket.ParseDoDItems(body)
	tk.PlanItems = ticket.ParsePlanItems(body)

	out := ticketJSON(tk, lanes)

	plan, ok := out["plan_items"].([]map[string]any)
	if !ok {
		t.Fatalf("plan_items missing or wrong type: %#v", out["plan_items"])
	}
	if len(plan) != 2 {
		t.Fatalf("plan_items = %d, want 2", len(plan))
	}
	// The number is what 'jaira dod' takes, so it has to be in the payload rather
	// than left for the agent to infer from array position.
	if plan[1]["n"] != 2 || plan[1]["state"] != "doing" || plan[1]["text"] != "implement" {
		t.Errorf("plan_items[1] = %#v", plan[1])
	}
	if plan[0]["state"] != "done" || plan[0]["done"] != true {
		t.Errorf("plan_items[0] = %#v", plan[0])
	}

	dod, ok := out["dod_items"].([]map[string]any)
	if !ok {
		t.Fatalf("dod_items missing or wrong type: %#v", out["dod_items"])
	}
	if len(dod) != 1 || dod[0]["state"] != "todo" || dod[0]["done"] != false {
		t.Errorf("dod_items = %#v", dod)
	}

	if out["plan_complete"] != false {
		t.Errorf("plan_complete = %v, want false", out["plan_complete"])
	}
	if out["dod_complete"] != false {
		t.Errorf("dod_complete = %v, want false", out["dod_complete"])
	}
}

// A ticket with no checklists must emit empty arrays rather than null, so an
// agent can iterate without a nil check.
func TestTicketJSONChecklistsEmptyNotNull(t *testing.T) {
	lanes, err := lane.Load()
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo"}
	out := ticketJSON(tk, lanes)
	for _, k := range []string{"plan_items", "dod_items"} {
		v, ok := out[k].([]map[string]any)
		if !ok {
			t.Fatalf("%s is not an array: %#v", k, out[k])
		}
		if v == nil {
			t.Errorf("%s is nil; want an empty array", k)
		}
	}
	// A ticket with no plan is vacuously complete; one with no criteria is not.
	if out["plan_complete"] != true {
		t.Errorf("plan_complete = %v, want true for a ticket with no plan", out["plan_complete"])
	}
	if out["dod_complete"] != false {
		t.Errorf("dod_complete = %v, want false for a ticket with no criteria", out["dod_complete"])
	}
}
