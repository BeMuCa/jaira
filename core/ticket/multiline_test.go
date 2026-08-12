package ticket

import "testing"

// Fields edited in the board can hold more than one line, so the frontmatter
// writer has to survive newlines without corrupting the document. It stores them
// as an escaped double-quoted scalar, which is lossless but noisy in a diff —
// which is why long prose belongs in the body rather than in a field.
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
