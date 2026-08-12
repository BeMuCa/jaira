package ticket

import "testing"

// A proof line rides directly beneath its item and must never be counted as a
// criterion of its own — that is the one hard constraint the format exists to
// satisfy.
func TestProofAttachesToItemAboveAndIsNotItsOwnItem(t *testing.T) {
	body := "## Definition of Done\n\n" +
		"- [x] 429 returned above 100/min\n" +
		"  proof: internal/x.go:12; TestRateLimit\n" +
		"- [ ] documented\n"
	items := ParseDoDItems(body)
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2 (the proof line must not count as one)", len(items))
	}
	if items[0].Proof != "internal/x.go:12; TestRateLimit" {
		t.Errorf("items[0].Proof = %q", items[0].Proof)
	}
	if items[1].Proof != "" {
		t.Errorf("items[1].Proof = %q, want empty", items[1].Proof)
	}
}

// The two checklists are independent: a proof under a Plan item must not leak
// into the Definition of Done, or vice versa.
func TestProofStaysWithinItsSection(t *testing.T) {
	body := "## Plan\n\n" +
		"- [x] write the spec\n" +
		"  proof: docs/spec.md\n\n" +
		"## Definition of Done\n\n" +
		"- [ ] implement\n"
	plan := ParsePlanItems(body)
	dod := ParseDoDItems(body)
	if len(plan) != 1 || plan[0].Proof != "docs/spec.md" {
		t.Fatalf("plan items = %#v", plan)
	}
	if len(dod) != 1 || dod[0].Proof != "" {
		t.Fatalf("dod items = %#v, want no proof leaking across sections", dod)
	}
}

// SetItemProof on an item with no proof inserts one line directly beneath it,
// leaving the item's own line byte-identical.
func TestSetItemProofInsertsForAnItemWithNone(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n- [ ] two\n"
	next, err := SetItemProof(body, SectionDoD, 0, "internal/x.go:12")
	if err != nil {
		t.Fatal(err)
	}
	want := "## Definition of Done\n\n- [ ] one\n  proof: internal/x.go:12\n- [ ] two\n"
	if next != want {
		t.Errorf("got:\n%q\nwant:\n%q", next, want)
	}
	items := ParseDoDItems(next)
	if items[0].Text != "one" || items[0].Proof != "internal/x.go:12" {
		t.Errorf("items[0] = %#v", items[0])
	}
}

// A second SetItemProof on the same item replaces the line rather than
// stacking a second one.
func TestSetItemProofReplacesRatherThanStacks(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n"
	first, err := SetItemProof(body, SectionDoD, 0, "first proof")
	if err != nil {
		t.Fatal(err)
	}
	second, err := SetItemProof(first, SectionDoD, 0, "second proof")
	if err != nil {
		t.Fatal(err)
	}
	want := "## Definition of Done\n\n- [ ] one\n  proof: second proof\n"
	if second != want {
		t.Errorf("got:\n%q\nwant:\n%q", second, want)
	}
	items := ParseDoDItems(second)
	if len(items) != 1 {
		t.Fatalf("got %d items after a replaced proof, want 1", len(items))
	}
	if items[0].Proof != "second proof" {
		t.Errorf("items[0].Proof = %q", items[0].Proof)
	}
}

// Ticking a later item must leave an earlier item's proof line attached to
// its own item, both immediately and after a re-parse.
func TestSetItemStateOnLaterItemLeavesEarlierProofIntact(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n- [ ] two\n- [ ] three\n"
	withProof, err := SetItemProof(body, SectionDoD, 0, "evidence for one")
	if err != nil {
		t.Fatal(err)
	}
	next, err := SetItemState(withProof, SectionDoD, 2, StateDone)
	if err != nil {
		t.Fatal(err)
	}
	items := ParseDoDItems(next)
	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	if items[0].Proof != "evidence for one" {
		t.Errorf("item 1's proof was lost or moved: %#v", items[0])
	}
	if items[2].State != StateDone {
		t.Errorf("item 3 was not marked done: %#v", items[2])
	}
	if items[1].Proof != "" || items[2].Proof != "" {
		t.Errorf("proof leaked onto an item that never had one: %#v %#v", items[1], items[2])
	}
}

// AddItem after a last item that carries a proof must land the new checkbox
// below the proof line, so the proof stays attached to the item it belongs to.
func TestAddItemAfterProofCarryingLastItem(t *testing.T) {
	body := "## Definition of Done\n\n- [x] one\n"
	withProof, err := SetItemProof(body, SectionDoD, 0, "evidence for one")
	if err != nil {
		t.Fatal(err)
	}
	next, err := AddItem(withProof, SectionDoD, "two")
	if err != nil {
		t.Fatal(err)
	}
	want := "## Definition of Done\n\n- [x] one\n  proof: evidence for one\n- [ ] two\n"
	if next != want {
		t.Errorf("got:\n%q\nwant:\n%q", next, want)
	}
	items := ParseDoDItems(next)
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Proof != "evidence for one" {
		t.Errorf("item 1's proof reattached to the wrong item: %#v", items[0])
	}
	if items[1].Text != "two" || items[1].Proof != "" {
		t.Errorf("item 2 = %#v", items[1])
	}
}

// Body -> parse -> write proof -> re-parse round-trips text, state and proof.
func TestProofRoundTrips(t *testing.T) {
	body := "## Definition of Done\n\n- [~] streaming write\n- [ ] documented\n"
	before := ParseDoDItems(body)
	if before[0].Text != "streaming write" || before[0].State != StateDoing {
		t.Fatalf("before = %#v", before[0])
	}

	withProof, err := SetItemProof(body, SectionDoD, 0, "internal/writer.go:88")
	if err != nil {
		t.Fatal(err)
	}
	after := ParseDoDItems(withProof)
	if len(after) != 2 {
		t.Fatalf("got %d items after writing a proof, want 2", len(after))
	}
	if after[0].Text != "streaming write" || after[0].State != StateDoing || after[0].Proof != "internal/writer.go:88" {
		t.Errorf("after = %#v", after[0])
	}
	if after[1].Text != "documented" || after[1].Proof != "" {
		t.Errorf("after[1] = %#v", after[1])
	}
}

// A newline embedded in the proof text must not let it re-parse as a
// checkbox or a heading of its own — whitespace is collapsed before writing.
func TestSetItemProofCollapsesEmbeddedNewlines(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n"
	next, err := SetItemProof(body, SectionDoD, 0, "line one\n- [x] injected\n## Injected Heading\nline two")
	if err != nil {
		t.Fatal(err)
	}
	items := ParseDoDItems(next)
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1 (a newline in the proof text must not inject new items)", len(items))
	}
	if items[0].Proof != "line one - [x] injected ## Injected Heading line two" {
		t.Errorf("items[0].Proof = %q", items[0].Proof)
	}
}

func TestSetItemProofRejectsEmptyText(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n"
	if _, err := SetItemProof(body, SectionDoD, 0, "   "); err == nil {
		t.Error("expected an error for an empty proof")
	}
}

func TestSetItemProofRejectsMissingChecklist(t *testing.T) {
	if _, err := SetItemProof("no checklist here\n", SectionDoD, 0, "x"); err == nil {
		t.Error("expected an error for a body with no definition-of-done checklist")
	}
}

func TestSetItemProofRejectsOutOfRangeItem(t *testing.T) {
	body := "## Definition of Done\n\n- [ ] one\n"
	if _, err := SetItemProof(body, SectionDoD, 5, "x"); err == nil {
		t.Error("expected an error for an item number out of range")
	}
}
