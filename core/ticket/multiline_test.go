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
