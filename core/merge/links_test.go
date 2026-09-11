package merge

import (
	"strings"
	"testing"
)

// linkDoc is doc() with the two link fields this file is about.
func linkDoc(parent, related, updated string) string {
	var b strings.Builder
	b.WriteString("---\nid: 01TEST\ntitle: A ticket\nstatus: todo\nassignee: berk\ngoal: g\n")
	b.WriteString("parent: " + parent + "\n")
	b.WriteString("related: [" + related + "]\n")
	b.WriteString("updated-at: " + updated + "\n---\n\nbody\n")
	return b.String()
}

// related is a list, and two sessions relating the same ticket to two
// different neighbours are both right. Picking a side would drop one of them
// silently, which is the failure a union exists to rule out.
func TestRelatedUnionsBothSides(t *testing.T) {
	base := linkDoc("01PARENT", "", "2026-08-11T10:00:00Z")
	ours := linkDoc("01PARENT", "01AAA", "2026-08-11T10:01:00Z")
	theirs := linkDoc("01PARENT", "01BBB", "2026-08-11T10:02:00Z")

	r := mergeStr(t, base, ours, theirs)

	if !r.Clean() {
		t.Fatalf("a union must not conflict: %v", r.Conflicts)
	}
	got := string(r.Merged)
	for _, want := range []string{"01AAA", "01BBB"} {
		if !strings.Contains(got, want) {
			t.Errorf("merged related lost %q:\n%s", want, got)
		}
	}
}

// parent is a scalar: a ticket is part of one thing. Two answers is a
// disagreement resolved by recency, not two facts to be kept side by side —
// a unioned parent would make the tree ambiguous everywhere it is drawn.
func TestParentIsAScalarNotAUnion(t *testing.T) {
	base := linkDoc("01AAA", "", "2026-08-11T10:00:00Z")
	ours := linkDoc("01AAA", "", "2026-08-11T10:01:00Z")
	theirs := linkDoc("01BBB", "", "2026-08-11T10:02:00Z")

	r := mergeStr(t, base, ours, theirs)

	got := string(r.Merged)
	if strings.Contains(got, "01AAA") && strings.Contains(got, "01BBB") {
		t.Errorf("parent must not end up carrying both answers:\n%s", got)
	}
	if !strings.Contains(got, "01BBB") {
		t.Errorf("the newer side changed it; want 01BBB:\n%s", got)
	}
}
