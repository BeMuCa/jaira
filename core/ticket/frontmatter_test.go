package ticket

import (
	"strings"
	"testing"
)

// A deliberately awkward ticket: aligned trailing comments, mixed quoting, a
// block scalar, a flow-style empty list, a block list, and an opaque nested
// block the tool must never interpret. Byte fidelity is asserted against this.
const awkward = `---
# Created by jaira from a Claude session. Do not reorder these fields by hand.
id: 01J8QK3ZXG4M7YP2VN8QRTKD9F
title: "Fix session cookie dropped on 302"
status: todo          # lane id, managed by jaira
ready: true

# --- ownership ---------------------------------------------------
creator: berk
assignee: berk        # always a human, even when an agent does the work
executed-by: null

goal: |
  Session must survive the OAuth round-trip so users are not
  silently logged out mid-flow.
definition-of-done: 'session survives OAuth round-trip, covered by a test'

blocked-by:
  - 01J8QK1AAA0000000000000000
  - 01J8QK2BBB0000000000000000

commits: []
model-tier: cheap

# Reserved for a future Jira/YouTrack adapter. jaira must round-trip this
# block byte-for-byte without understanding any of it.
external:
  provider: jira
  key: PROJ-1234
  fields:
    customfield_10021: "3"
    labels: [auth, regression]

updated-at: 2026-08-11T20:14:03Z
---

# Notes

Reproduced on Safari only. See thread with @teammate.
`

func mustParse(t *testing.T, src string) *Doc {
	t.Helper()
	d, err := ParseDoc([]byte(src))
	if err != nil {
		t.Fatalf("ParseDoc: %v", err)
	}
	return d
}

// changedLines reports which 1-based line numbers differ between two files.
func changedLines(a, b string) []int {
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	n := len(al)
	if len(bl) > n {
		n = len(bl)
	}
	var out []int
	for i := 0; i < n; i++ {
		var x, y string
		if i < len(al) {
			x = al[i]
		}
		if i < len(bl) {
			y = bl[i]
		}
		if x != y {
			out = append(out, i+1)
		}
	}
	return out
}

func TestParseRoundTripIsExact(t *testing.T) {
	d := mustParse(t, awkward)
	if got := d.String(); got != awkward {
		t.Errorf("round trip changed the file\n--- got ---\n%s", got)
	}
}

func TestSetScalarTouchesExactlyOneLine(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("status", "in-progress"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	got := changedLines(awkward, d.String())
	if len(got) != 1 || got[0] != 5 {
		t.Errorf("expected only line 5 to change, got lines %v", got)
	}
	if !strings.Contains(d.String(), "# lane id, managed by jaira") {
		t.Error("trailing comment was lost")
	}
}

func TestSetScalarPreservesUnrelatedCommentAlignment(t *testing.T) {
	// This is the exact failure that ruled out AST re-printing: patching one
	// field must not re-align a comment on a different field.
	d := mustParse(t, awkward)
	if err := d.SetScalar("status", "review"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	const untouched = "assignee: berk        # always a human, even when an agent does the work"
	if !strings.Contains(d.String(), untouched) {
		t.Errorf("unrelated field was reformatted; expected to still contain:\n  %q", untouched)
	}
}

func TestSetScalarPreservesQuotingStyle(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("title", "Fix cookie on 302 (Safari)"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !strings.Contains(d.String(), `title: "Fix cookie on 302 (Safari)"`) {
		t.Error("double-quoted field did not stay double-quoted")
	}

	d2 := mustParse(t, awkward)
	if err := d2.SetScalar("definition-of-done", "it works"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !strings.Contains(d2.String(), `definition-of-done: 'it works'`) {
		t.Errorf("single-quoted field did not stay single-quoted:\n%s", d2.Frontmatter())
	}
}

func TestOpaqueBlockSurvivesWrites(t *testing.T) {
	d := mustParse(t, awkward)
	for _, kv := range [][2]string{
		{"status", "done"}, {"assignee", "teammate"}, {"model-tier", "strong"},
		{"executed-by", "opus"}, {"ready", "false"},
	} {
		if err := d.SetScalar(kv[0], kv[1]); err != nil {
			t.Fatalf("SetScalar(%s): %v", kv[0], err)
		}
	}
	const block = "external:\n  provider: jira\n  key: PROJ-1234\n  fields:\n    customfield_10021: \"3\"\n    labels: [auth, regression]\n"
	if !strings.Contains(d.String(), block) {
		t.Errorf("reserved external: block was modified\n%s", d.Frontmatter())
	}
	if !strings.HasSuffix(d.String(), "Reproduced on Safari only. See thread with @teammate.\n") {
		t.Error("markdown body was damaged")
	}
	if !strings.Contains(d.String(), "updated-at: 2026-08-11T20:14:03Z\n---\n") {
		t.Error("closing delimiter was damaged")
	}
}

func TestSetScalarAppendsMissingKey(t *testing.T) {
	d := mustParse(t, awkward)
	if d.Has("claimed-by") {
		t.Fatal("precondition: key should be absent")
	}
	if err := d.SetScalar("claimed-by", "session-7"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !d.Has("claimed-by") {
		t.Fatal("key was not added")
	}
	got, _, err := d.Scalar("claimed-by")
	if err != nil || got != "session-7" {
		t.Errorf("read back %q, err %v", got, err)
	}
	// Appending must be a pure insertion: deleting the one new line from the
	// result must reproduce the original file byte-for-byte. This is stronger
	// than comparing line numbers, which shift under insertion.
	out := d.String()
	lines := strings.Split(out, "\n")
	var idx = -1
	for i, l := range lines {
		if strings.HasPrefix(l, "claimed-by:") {
			idx = i
			break
		}
	}
	if idx < 0 {
		t.Fatalf("appended line not found in output:\n%s", out)
	}
	withoutNew := strings.Join(append(append([]string{}, lines[:idx]...), lines[idx+1:]...), "\n")
	if withoutNew != awkward {
		t.Errorf("append was not a pure insertion; diff at lines %v", changedLines(awkward, withoutNew))
	}
}

func TestSetScalarRefusesCollections(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("blocked-by", "nope"); err == nil {
		t.Error("expected an error setting a list field via SetScalar")
	}
	if err := d.SetScalar("external", "nope"); err == nil {
		t.Error("expected an error setting a mapping field via SetScalar")
	}
}

func TestScalarQuotingIsSafe(t *testing.T) {
	// Values that would change meaning if written bare must be quoted.
	cases := map[string]string{
		"true": `"true"`, "42": `"42"`, "null": `"null"`,
		"has: colon": `"has: colon"`, "": `""`, "1.5": `"1.5"`,
		"plain": "plain", "with-dash": "with-dash",
	}
	for in, want := range cases {
		if got := encodeScalar(in); got != want {
			t.Errorf("encodeScalar(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSetListReplacesBlock(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetList("blocked-by", []string{"01AAA", "01BBB", "01CCC"}); err != nil {
		t.Fatalf("SetList: %v", err)
	}
	got, err := d.List("blocked-by")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 3 || got[0] != "01AAA" || got[2] != "01CCC" {
		t.Errorf("read back %v", got)
	}
	// Neighbouring keys must survive.
	for _, k := range []string{"definition-of-done", "commits", "model-tier"} {
		if !d.Has(k) {
			t.Errorf("neighbouring key %q disappeared\n%s", k, d.Frontmatter())
		}
	}
}

func TestSetListOnEmptyFlowList(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetList("commits", []string{"abc123", "def456"}); err != nil {
		t.Fatalf("SetList: %v", err)
	}
	got, err := d.List("commits")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("read back %v\n%s", got, d.Frontmatter())
	}
	if !d.Has("model-tier") {
		t.Errorf("following key was consumed\n%s", d.Frontmatter())
	}
}

func TestSetListToEmpty(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetList("blocked-by", nil); err != nil {
		t.Fatalf("SetList: %v", err)
	}
	if !strings.Contains(d.Frontmatter(), "blocked-by: []") {
		t.Errorf("expected flow-style empty list\n%s", d.Frontmatter())
	}
	if !d.Has("definition-of-done") || !d.Has("commits") {
		t.Errorf("neighbours damaged\n%s", d.Frontmatter())
	}
}

func TestValidateRejectsAnchors(t *testing.T) {
	src := "---\na: &anchor value\nb: *anchor\n---\nbody\n"
	d := mustParse(t, src)
	if err := d.Validate(); err == nil {
		t.Error("expected anchors to be rejected")
	}
}

func TestValidateAcceptsAwkwardButSafe(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.Validate(); err != nil {
		t.Errorf("Validate rejected a safe ticket: %v", err)
	}
}

func TestParseErrors(t *testing.T) {
	if _, err := ParseDoc([]byte("no frontmatter here\n")); err != ErrNoFrontmatter {
		t.Errorf("got %v, want ErrNoFrontmatter", err)
	}
	if _, err := ParseDoc([]byte("---\nid: x\nnever closed\n")); err != ErrUnterminated {
		t.Errorf("got %v, want ErrUnterminated", err)
	}
}

func TestBodyAccessors(t *testing.T) {
	d := mustParse(t, awkward)
	if !strings.Contains(d.Body(), "Reproduced on Safari") {
		t.Errorf("Body() = %q", d.Body())
	}
	d.SetBody("# New\n\nreplaced\n")
	if strings.Contains(d.Body(), "Safari") {
		t.Error("SetBody did not replace the body")
	}
	if !d.Has("status") {
		t.Error("SetBody damaged the frontmatter")
	}
}

func TestCRLFAndBOMTolerated(t *testing.T) {
	src := "\ufeff---\nid: x\nstatus: todo\n---\nbody\n"
	d, err := ParseDoc([]byte(src))
	if err != nil {
		t.Fatalf("ParseDoc with BOM: %v", err)
	}
	if err := d.SetScalar("status", "done"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !strings.HasPrefix(d.String(), "\ufeff---") {
		t.Error("BOM was lost")
	}
	if !strings.Contains(d.String(), "status: done") {
		t.Error("write failed")
	}
}
