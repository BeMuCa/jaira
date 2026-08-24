package ticket

import (
	"strings"
	"testing"
)

// [-] has to survive a parse and a write in both directions, or a hand-edited
// ticket loses the distinction the marker exists to record.
func TestSupersededRoundTrips(t *testing.T) {
	body := "## Definition of Done\n\n- [-] the old wording\n- [x] the new one\n"
	items := ParseDoDItems(body)
	if len(items) != 2 {
		t.Fatalf("parsed %d items, want 2", len(items))
	}
	if items[0].State != StateSuperseded {
		t.Errorf("state = %v, want superseded", items[0].State)
	}
	if items[0].State.String() != "superseded" || items[0].State.Marker() != "-" {
		t.Errorf("word/marker = %q/%q", items[0].State.String(), items[0].State.Marker())
	}

	out, err := SetItemState(body, SectionDoD, 1, StateSuperseded)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- [-] the new one") {
		t.Errorf("writing superseded did not produce [-]: %q", out)
	}
	if len(out) != len(body) {
		t.Errorf("length changed: %d -> %d", len(body), len(out))
	}
}

// A superseded item was never done, and nothing may report that it was — that is
// the whole reason it is not simply written as [x].
func TestSupersededIsNotChecked(t *testing.T) {
	it := DoDItem{State: StateSuperseded}
	if it.Checked() {
		t.Error("a superseded item reports as done")
	}
	if !it.Settled() {
		t.Error("a superseded item is still counted as outstanding work")
	}
}

// The point of the state: it retires a criterion instead of pretending it was
// met, so the ticket can still complete.
func TestSupersededDoesNotBlockCompletion(t *testing.T) {
	body := "## Definition of Done\n\n- [-] replaced by item 2\n- [x] the new criterion\n"
	tk := &Ticket{Body: body}
	tk.DoDItems = ParseDoDItems(body)
	complete, remaining := tk.DoDComplete()
	if !complete {
		t.Errorf("a done-plus-superseded checklist is not complete; remaining: %v", remaining)
	}
}

// An unfinished item beside a superseded one must still block, and only the
// unfinished one may be named — asking for work that was explicitly dropped is
// how the old workaround started.
func TestSupersededIsNotListedAsRemaining(t *testing.T) {
	body := "## Definition of Done\n\n- [-] dropped\n- [ ] still open\n"
	tk := &Ticket{Body: body}
	tk.DoDItems = ParseDoDItems(body)
	complete, remaining := tk.DoDComplete()
	if complete {
		t.Fatal("an open criterion did not block completion")
	}
	if len(remaining) != 1 || remaining[0] != "still open" {
		t.Errorf("remaining = %v, want only the open item", remaining)
	}
}

// The same rule for the plan, which the terminal lane checks separately.
func TestSupersededPlanStepDoesNotBlock(t *testing.T) {
	body := "## Plan\n\n- [x] did it\n- [-] turned out to be unnecessary\n"
	tk := &Ticket{Body: body}
	tk.PlanItems = ParsePlanItems(body)
	if complete, remaining := tk.PlanComplete(); !complete {
		t.Errorf("plan is not complete; remaining: %v", remaining)
	}
}

// An unrecognised marker still reads as unfinished. Dropping it would shorten
// the list and make the checklist easier to satisfy than its author wrote it,
// and adding a fourth state must not weaken that.
func TestUnknownMarkerIsStillTodo(t *testing.T) {
	items := ParseDoDItems("## Definition of Done\n\n- [?] who knows\n")
	if len(items) != 1 || items[0].State != StateTodo {
		t.Errorf("items = %+v, want one todo", items)
	}
}
