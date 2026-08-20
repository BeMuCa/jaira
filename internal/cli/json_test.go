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
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	body := `## Plan

- [x] write the spec
  proof: docs/spec.md
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
	// Proof rides alongside text and state in every payload that carries a
	// checklist, so an agent reading --json output can see the evidence
	// without a second command.
	if plan[0]["proof"] != "docs/spec.md" {
		t.Errorf("plan_items[0][\"proof\"] = %#v, want %q", plan[0]["proof"], "docs/spec.md")
	}
	if plan[1]["proof"] != "" {
		t.Errorf("plan_items[1][\"proof\"] = %#v, want empty", plan[1]["proof"])
	}

	dod, ok := out["dod_items"].([]map[string]any)
	if !ok {
		t.Fatalf("dod_items missing or wrong type: %#v", out["dod_items"])
	}
	if len(dod) != 1 || dod[0]["state"] != "todo" || dod[0]["done"] != false {
		t.Errorf("dod_items = %#v", dod)
	}
	if dod[0]["proof"] != "" {
		t.Errorf("dod_items[0][\"proof\"] = %#v, want empty", dod[0]["proof"])
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
	lanes, err := lane.Load("")
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

// The review fields existed only in the TUI: an agent asked to hand a reader the
// review could not read one. review-check especially, since it is the only field
// the reader acts on.
func TestTicketJSONCarriesTheReview(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "signoff",
		ReviewSummary: "streamed the writer",
		ReviewGaps:    "none",
		ReviewVerdict: "matches the criteria",
		ReviewCheck:   "1. go run ./cmd/app  2. open /export  3. the download starts at once",
	}

	out := ticketJSON(tk, lanes)

	rev, ok := out["review"].(map[string]string)
	if !ok {
		t.Fatalf("review missing or wrong type: %#v", out["review"])
	}
	for field, want := range map[string]string{
		"summary": "streamed the writer",
		"gaps":    "none",
		"verdict": "matches the criteria",
		"check":   "1. go run ./cmd/app  2. open /export  3. the download starts at once",
	} {
		if rev[field] != want {
			t.Errorf("review[%q] = %q, want %q", field, rev[field], want)
		}
	}
}

// The route stops being something every caller derives for itself.
func TestTicketJSONNamesTheNextLane(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ status, want string }{
		{"todo", "in-progress"},
		{"review", "signoff"},
		{"done", ""},
	} {
		tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: c.status}
		if got := ticketJSON(tk, lanes)["next_lane"]; got != c.want {
			t.Errorf("from %s: next_lane = %v, want %q", c.status, got, c.want)
		}
	}
}
