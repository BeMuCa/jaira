package ticket

import (
	"strings"
	"testing"
)

// Fields edited in the board can hold more than one line, so the frontmatter
// writer has to survive newlines without corrupting the document. A value
// eligible for the block form (see blockable) is stored as an indented block
// literal, which is readable and diffs as a one-line change; anything not
// provably safe as a block still falls back to an escaped double-quoted scalar.
func TestMultiLineScalarRoundTrips(t *testing.T) {
	src := "---\nid: 01AAA\ntitle: t\ngoal: short\n---\n\nbody\n"
	d, err := ParseDoc([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := "first line\nsecond line\n\nfourth after a blank"
	if err := d.SetScalar("goal", want); err != nil {
		t.Fatalf("SetScalar with newlines: %v", err)
	}
	t.Logf("document after write:\n%s", d.String())
	d2, err := ParseDoc(d.Bytes())
	if err != nil {
		t.Fatalf("reparse: %v", err)
	}
	got, ok, err := d2.Scalar("goal")
	if err != nil || !ok {
		t.Fatalf("read back: ok=%v err=%v", ok, err)
	}
	if got != want {
		t.Errorf("round trip changed the value:\n got %q\nwant %q", got, want)
	}
}

// A block scalar's own token is its header ("|" or ">-"), not its content. A
// naive read via GetToken().Value returns the header, so a ticket hand-written
// (or merge-produced) with a block scalar silently reads back as that literal
// string instead of the text it holds.
func TestScalarReadsBlockContentNotHeader(t *testing.T) {
	src := "---\nid: 01AAA\nstatus: todo\n" +
		"context: |\n  first line\n  second line\n\n  fourth\n" +
		"outcome-resolves: >-\n  folded one\n  folded two\n" +
		"dashed: |-\n  no trailing newline\n" +
		"goal: after\n" +
		"blocked-by:\n  - 01BBB\n---\n\nbody\n"
	d, err := ParseDoc([]byte(src))
	if err != nil {
		t.Fatal(err)
	}

	if got, ok, err := d.Scalar("context"); err != nil || !ok {
		t.Fatalf("context: ok=%v err=%v", ok, err)
	} else if want := "first line\nsecond line\n\nfourth\n"; got != want {
		t.Errorf("context = %q, want %q", got, want)
	}

	if got, ok, err := d.Scalar("outcome-resolves"); err != nil || !ok {
		t.Fatalf("outcome-resolves: ok=%v err=%v", ok, err)
	} else if want := "folded one folded two"; got != want {
		t.Errorf("outcome-resolves = %q, want %q", got, want)
	}

	if got, ok, err := d.Scalar("dashed"); err != nil || !ok {
		t.Fatalf("dashed: ok=%v err=%v", ok, err)
	} else if want := "no trailing newline"; got != want {
		t.Errorf("dashed = %q, want %q", got, want)
	}

	if got, ok, err := d.Scalar("goal"); err != nil || !ok {
		t.Fatalf("goal: ok=%v err=%v", ok, err)
	} else if got != "after" {
		t.Errorf("goal = %q, want %q", got, "after")
	}

	if _, ok, err := d.Scalar("missing"); err != nil || ok {
		t.Errorf("missing key: ok=%v err=%v, want ok=false err=nil", ok, err)
	}
	if _, ok, err := d.Scalar("blocked-by"); err == nil {
		t.Errorf("scalar on a list-typed key should error, not silently succeed: ok=%v", ok)
	}
}

// Decode uses Scalar under the hood; a ticket whose context is a block literal
// must project to its content in Ticket.Context, not "|".
func TestDecodeBlockContext(t *testing.T) {
	src := "---\nid: 01AAA\nstatus: todo\ncontext: |\n  why this exists\n  more detail\n---\n\nbody\n"
	d, err := ParseDoc([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	tk, err := Decode(d, "ticket.md")
	if err != nil {
		t.Fatal(err)
	}
	if want := "why this exists\nmore detail\n"; tk.Context != want {
		t.Errorf("Ticket.Context = %q, want %q", tk.Context, want)
	}
}

// NewDoc emits a multi-line field as an indented block literal, not a quoted
// one-liner with \n escapes.
func TestNewDocWritesBlockLiteral(t *testing.T) {
	d := NewDoc(map[string]string{
		"id":      "01AAA",
		"status":  "todo",
		"context": "first line\nsecond line\n\nfourth",
	}, nil, "")
	const want = "context: |-\n  first line\n  second line\n\n  fourth\n"
	if !strings.Contains(d.Frontmatter(), want) {
		t.Errorf("frontmatter does not contain the expected block:\n%s\n--- got ---\n%s", want, d.Frontmatter())
	}
	if strings.Contains(d.Frontmatter(), `\n`) {
		t.Errorf("frontmatter still contains an escaped newline:\n%s", d.Frontmatter())
	}
}

// TestNewDocSingleLineUnchanged states, byte for byte, that this change does
// not alter how a fully single-line ticket is written.
func TestNewDocSingleLineUnchanged(t *testing.T) {
	d := NewDoc(map[string]string{
		"id":      "01AAA",
		"title":   "Fix the thing",
		"status":  "todo",
		"creator": "berk",
		"goal":    "it works",
		"context": "one line",
	}, map[string][]string{
		"blocked-by": {"01BBB"},
	}, "# Fix the thing\n")
	const want = "---\n" +
		"id: 01AAA\n" +
		"title: Fix the thing\n" +
		"status: todo\n" +
		"creator: berk\n" +
		"goal: it works\n" +
		"context: one line\n" +
		"blocked-by:\n  - 01BBB\n" +
		"---\n\n# Fix the thing\n"
	if got := d.String(); got != want {
		t.Errorf("single-line NewDoc output changed:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// Round trip through NewDoc -> ParseDoc -> Scalar must reproduce the original
// value byte for byte, with an internal blank line and both chomping styles.
func TestNewDocBlockRoundTrips(t *testing.T) {
	cases := map[string]string{
		"with_trailing_newline":    "first\n\nthird\n",
		"without_trailing_newline": "first\n\nthird",
	}
	for name, val := range cases {
		t.Run(name, func(t *testing.T) {
			d := NewDoc(map[string]string{"id": "01AAA", "status": "todo", "context": val}, nil, "")
			d2, err := ParseDoc(d.Bytes())
			if err != nil {
				t.Fatalf("reparse: %v", err)
			}
			got, ok, err := d2.Scalar("context")
			if err != nil || !ok {
				t.Fatalf("read back: ok=%v err=%v", ok, err)
			}
			if got != val {
				t.Errorf("round trip changed the value:\n got %q\nwant %q", got, val)
			}
		})
	}
}

// SetScalar with a multi-line value over an existing single-line field
// converts that entry to a block and leaves every other line untouched.
func TestSetScalarSingleToBlock(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("definition-of-done", "line one\n\nline three"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !strings.Contains(d.String(), "definition-of-done: |-\n  line one\n\n  line three\n") {
		t.Errorf("expected a block literal:\n%s", d.Frontmatter())
	}
	// The blank line and blocked-by list right after definition-of-done must
	// survive untouched.
	if !strings.Contains(d.String(), "blocked-by:\n  - 01J8QK1AAA0000000000000000\n  - 01J8QK2BBB0000000000000000\n") {
		t.Errorf("blocked-by list was disturbed:\n%s", d.Frontmatter())
	}
	if !strings.Contains(d.String(), `assignee: berk        # always a human, even when an agent does the work`) {
		t.Error("unrelated field above the edit was disturbed")
	}
}

// SetScalar with a single-line value over an existing block collapses the
// whole block to one line, leaving no orphaned content lines.
func TestSetScalarBlockToSingle(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("goal", "one line now"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if strings.Contains(d.String(), "Session must survive") {
		t.Errorf("old block content was not removed:\n%s", d.Frontmatter())
	}
	if !strings.Contains(d.String(), "goal: one line now\n") {
		t.Errorf("expected a single-line goal:\n%s", d.Frontmatter())
	}
	// definition-of-done, on the very next non-blank line, must be untouched.
	if !strings.Contains(d.String(), `definition-of-done: 'session survives OAuth round-trip, covered by a test'`) {
		t.Error("field after the collapsed block was disturbed")
	}
	if _, err := ParseDoc(d.Bytes()); err != nil {
		t.Fatalf("document is no longer valid YAML after collapsing: %v", err)
	}
}

// SetScalar with a multi-line value over an existing block replaces the whole
// block, leaving no orphaned content lines.
func TestSetScalarBlockToBlock(t *testing.T) {
	d := mustParse(t, awkward)
	if err := d.SetScalar("goal", "new first\n\nnew third"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if strings.Contains(d.String(), "Session must survive") {
		t.Errorf("old block content was not fully replaced:\n%s", d.Frontmatter())
	}
	if !strings.Contains(d.String(), "goal: |-\n  new first\n\n  new third\n") {
		t.Errorf("expected the replacement block:\n%s", d.Frontmatter())
	}
	if !strings.Contains(d.String(), `definition-of-done: 'session survives OAuth round-trip, covered by a test'`) {
		t.Error("field after the replaced block was disturbed")
	}
}

// A trailing comment on a key line that is about to become a block must
// survive on the block header, and must still parse back out.
func TestSetScalarBlockCarriesTrailingComment(t *testing.T) {
	src := "---\nid: 01AAA\nstatus: todo\ncontext: value # from chat\ngoal: g\n---\n\nbody\n"
	d := mustParse(t, src)
	if err := d.SetScalar("context", "first\n\nthird"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	if !strings.Contains(d.String(), "context: |- # from chat\n") {
		t.Errorf("expected the comment to survive on the block header:\n%s", d.Frontmatter())
	}
	d2, err := ParseDoc(d.Bytes())
	if err != nil {
		t.Fatalf("reparse: %v", err)
	}
	got, ok, err := d2.Scalar("context")
	if err != nil || !ok {
		t.Fatalf("read back: ok=%v err=%v", ok, err)
	}
	if want := "first\n\nthird"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// A value whose content includes a line reading "---" must still round-trip:
// ParseDoc scans for the closing delimiter without YAML awareness, so this
// only works because the line is indented as block content.
func TestBlockContentWithDashLine(t *testing.T) {
	d := NewDoc(map[string]string{"id": "01AAA", "status": "todo", "context": "before\n---\nafter"}, nil, "")
	d2, err := ParseDoc(d.Bytes())
	if err != nil {
		t.Fatalf("reparse: %v", err)
	}
	got, ok, err := d2.Scalar("context")
	if err != nil || !ok {
		t.Fatalf("read back: ok=%v err=%v", ok, err)
	}
	if want := "before\n---\nafter"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBlockable(t *testing.T) {
	cases := []struct {
		name string
		val  string
		want bool
	}{
		{"no newline", "single line", false},
		{"empty", "", false},
		{"simple multiline", "a\nb", true},
		{"trailing single newline", "a\nb\n", true},
		{"trailing double newline", "a\nb\n\n", false},
		{"trailing space on a line", "a \nb", false},
		{"trailing tab on a line", "a\nb\t", false},
		{"leading space on a line", "a\n b", false},
		{"carriage return", "a\r\nb", false},
		{"control rune", "a\nb\x01c", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := blockable(c.val); got != c.want {
				t.Errorf("blockable(%q) = %v, want %v", c.val, got, c.want)
			}
		})
	}
}
