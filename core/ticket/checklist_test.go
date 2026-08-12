package ticket

import (
	"strings"
	"testing"
)

const checklistBody = `# Add rate limiting

## Plan

- [x] write the spec
- [ ] design the interface
- [ ] implement

## Definition of Done

- [ ] 429 returned above 100 req/min
- [ ] documented in README

## Notes

- [ ] this is not a criterion and must never be touched
`

// A write must be surgical. The frontmatter path already guarantees that unrelated
// fields survive a single-field edit byte for byte; the body deserves the same,
// because a checklist tick that reflows the prose around it makes every diff
// unreadable and every merge worse.
func TestSetItemStateChangesExactlyOneCharacter(t *testing.T) {
	out, err := SetItemState(checklistBody, SectionDoD, 1, StateDone)
	if err != nil {
		t.Fatalf("SetItemState: %v", err)
	}
	if len(out) != len(checklistBody) {
		t.Fatalf("length changed: %d -> %d", len(checklistBody), len(out))
	}
	var diffs []int
	for i := range out {
		if out[i] != checklistBody[i] {
			diffs = append(diffs, i)
		}
	}
	if len(diffs) != 1 {
		t.Fatalf("%d bytes differ, want exactly 1: %v", len(diffs), diffs)
	}
	if out[diffs[0]] != 'x' {
		t.Errorf("changed byte is %q, want 'x'", out[diffs[0]])
	}
	if !strings.Contains(out, "- [x] documented in README") {
		t.Error("the intended item was not ticked")
	}
}

// The Notes section holds checkboxes that are not acceptance criteria. Indexes
// are counted within the addressed section only, so a DoD index must never reach
// into a later section.
func TestSetItemStateIgnoresOtherSections(t *testing.T) {
	if _, err := SetItemState(checklistBody, SectionDoD, 2, StateDone); err == nil {
		t.Fatal("expected an error: the DoD has 2 items, index 2 is out of range")
	}
	out, err := SetItemState(checklistBody, SectionPlan, 1, StateDoing)
	if err != nil {
		t.Fatalf("SetItemState: %v", err)
	}
	if !strings.Contains(out, "- [~] design the interface") {
		t.Error("the plan item was not marked as in progress")
	}
	if !strings.Contains(out, "- [ ] this is not a criterion and must never be touched") {
		t.Error("a checkbox outside the addressed section was modified")
	}
}

// "Which step is the agent on" must have exactly one answer, so marking a second
// item as in progress moves the marker rather than adding another.
func TestOnlyOneItemIsDoingPerSection(t *testing.T) {
	out, err := SetItemState(checklistBody, SectionPlan, 1, StateDoing)
	if err != nil {
		t.Fatal(err)
	}
	out, err = SetItemState(out, SectionPlan, 2, StateDoing)
	if err != nil {
		t.Fatal(err)
	}
	items := ParsePlanItems(out)
	var doing []string
	for _, it := range items {
		if it.State == StateDoing {
			doing = append(doing, it.Text)
		}
	}
	if len(doing) != 1 || doing[0] != "implement" {
		t.Errorf("doing = %v, want exactly [implement]", doing)
	}
}

func TestSetItemStatePreservesCRLF(t *testing.T) {
	body := "## Definition of Done\r\n\r\n- [ ] first\r\n- [ ] second\r\n"
	out, err := SetItemState(body, SectionDoD, 0, StateDone)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- [x] first\r\n") {
		t.Errorf("CRLF line ending was lost: %q", out)
	}
	if strings.Count(out, "\r\n") != strings.Count(body, "\r\n") {
		t.Error("line ending count changed")
	}
}

func TestSetItemStateOutOfRangeLeavesBodyAlone(t *testing.T) {
	out, err := SetItemState(checklistBody, SectionPlan, 9, StateDone)
	if err == nil {
		t.Fatal("expected an out-of-range error")
	}
	if out != "" && out != checklistBody {
		t.Error("body was modified despite the error")
	}
}

func TestSetItemStateMissingSection(t *testing.T) {
	if _, err := SetItemState("# no checklists here\n", SectionDoD, 0, StateDone); err == nil {
		t.Fatal("expected an error when the section does not exist")
	}
}
