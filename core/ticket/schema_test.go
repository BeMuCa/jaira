package ticket

import "testing"

// An in-progress item must never make a checklist look finished. The parser used
// to recognise only "[ ]" and "[x]", dropping anything else outright, so a "[~]"
// item vanished and DoDComplete counted the survivors as the whole list. That
// let a ticket with outstanding work enter the terminal lane.
func TestInProgressItemIsNotComplete(t *testing.T) {
	body := `## Definition of Done

- [x] first item done
- [~] second item in progress
- [x] third item done
`
	tk := &Ticket{Body: body, DoDItems: ParseDoDItems(body)}
	if got := len(tk.DoDItems); got != 3 {
		t.Fatalf("expected 3 items, got %d — an item was dropped by the parser", got)
	}
	complete, remaining := tk.DoDComplete()
	if complete {
		t.Error("checklist reported complete while an item is still in progress")
	}
	if len(remaining) != 1 || remaining[0] != "second item in progress" {
		t.Errorf("remaining = %q, want the in-progress item", remaining)
	}
}

// An unrecognised marker is unfinished work, not absent work. Dropping it is the
// dangerous reading: it makes the list shorter and therefore easier to complete.
func TestUnknownMarkerParsesAsTodo(t *testing.T) {
	body := `## Definition of Done

- [x] done
- [/] cancelled?
- [?] unsure
`
	items := ParseDoDItems(body)
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d — unknown markers were dropped", len(items))
	}
	for _, it := range items[1:] {
		if it.State != StateTodo {
			t.Errorf("item %q: state = %v, want todo", it.Text, it.State)
		}
	}
}

func TestCheckboxStates(t *testing.T) {
	body := `## Definition of Done

- [ ] not started
- [~] being worked on
- [x] finished
- [X] also finished
`
	items := ParseDoDItems(body)
	want := []State{StateTodo, StateDoing, StateDone, StateDone}
	if len(items) != len(want) {
		t.Fatalf("got %d items, want %d", len(items), len(want))
	}
	for i, w := range want {
		if items[i].State != w {
			t.Errorf("item %d (%q): state = %v, want %v", i, items[i].Text, items[i].State, w)
		}
	}
	// Checked stays meaningful for existing callers: only done counts.
	if items[1].Checked() {
		t.Error("an in-progress item must not report as checked")
	}
	if !items[2].Checked() {
		t.Error("a done item must report as checked")
	}
}

// The Plan records method — spec, design, implement — while the Definition of
// Done records acceptance criteria. Only the latter gates the terminal lane, so
// the two must not be read out of the same section.
func TestPlanAndDoDAreParsedSeparately(t *testing.T) {
	body := `# Add rate limiting

## Plan

- [x] write the spec
- [~] design the interface
- [ ] implement

## Definition of Done

- [ ] 429 returned above 100 req/min
`
	plan := ParsePlanItems(body)
	if len(plan) != 3 {
		t.Fatalf("plan items = %d, want 3", len(plan))
	}
	if plan[1].State != StateDoing || plan[1].Text != "design the interface" {
		t.Errorf("plan[1] = %+v, want the in-progress design step", plan[1])
	}

	dod := ParseDoDItems(body)
	if len(dod) != 1 {
		t.Fatalf("dod items = %d, want 1 — the Plan section leaked into the DoD", len(dod))
	}
	if dod[0].Text != "429 returned above 100 req/min" {
		t.Errorf("dod[0].Text = %q", dod[0].Text)
	}
}

// A ticket with no Plan section is the normal case and must not be an error.
func TestPlanAbsentIsEmpty(t *testing.T) {
	if got := ParsePlanItems("## Definition of Done\n\n- [ ] a\n"); len(got) != 0 {
		t.Errorf("expected no plan items, got %d", len(got))
	}
}
