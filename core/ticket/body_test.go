package ticket

import (
	"strings"
	"testing"
)

// TestNewBodyMatchesRegressionBaseline asserts NewBody with the two options
// jaira has always shipped, unticked, is byte-identical to the body a ticket
// got before the default board existed — the existing shape is the
// regression baseline this function must not silently drift from.
func TestNewBodyMatchesRegressionBaseline(t *testing.T) {
	got := NewBody("Fix the thing", nil, []BodyOption{
		{Name: "brainstorm", Ticked: false},
		{Name: "planning", Ticked: false},
	})
	want := "# Fix the thing\n\n" +
		"## Definition of Done\n\n" +
		"- [ ] <A checkable statement, readable by someone who was not here>\n" +
		"\n## Options\n\n" +
		"- [ ] brainstorm\n" +
		"- [ ] planning\n" +
		"\n## Plan\n\n" +
		"<Steps, in order — filled in by the pre-process step, or by you.>\n" +
		"\n## Progress\n\n"
	if got != want {
		t.Errorf("NewBody regressed from the baseline shape:\ngot:  %q\nwant: %q", got, want)
	}
}

// TestNewBodyTicksNamedOptions asserts a ticked option renders as "- [x]"
// and everything else stays unticked.
func TestNewBodyTicksNamedOptions(t *testing.T) {
	got := NewBody("t", nil, []BodyOption{
		{Name: "brainstorm", Ticked: true},
		{Name: "planning", Ticked: false},
	})
	if !strings.Contains(got, "- [x] brainstorm\n") {
		t.Errorf("expected brainstorm ticked, got:\n%s", got)
	}
	if !strings.Contains(got, "- [ ] planning\n") {
		t.Errorf("expected planning unticked, got:\n%s", got)
	}
}

// TestNewBodyWithDoDUsesIt asserts a supplied definition of done replaces the
// placeholder.
func TestNewBodyWithDoDUsesIt(t *testing.T) {
	got := NewBody("t", []string{"ships behind a flag"}, nil)
	if !strings.Contains(got, "- [ ] ships behind a flag\n") {
		t.Errorf("expected the supplied dod, got:\n%s", got)
	}
	if strings.Contains(got, "checkable statement") {
		t.Errorf("placeholder dod must not appear when one was supplied, got:\n%s", got)
	}
}

// TestNewBodyWithNoOptionsOmitsChecklistItems asserts a nil or empty options
// slice still produces a valid (empty) Options section rather than panicking
// or emitting stale entries.
func TestNewBodyWithNoOptionsOmitsChecklistItems(t *testing.T) {
	got := NewBody("t", nil, nil)
	if !strings.Contains(got, "## Options\n\n\n## Plan") {
		t.Errorf("expected an empty Options section, got:\n%s", got)
	}
}

// TestNewBodyWritesOneBoxPerCriterion asserts several criteria become several
// boxes. The terminal lane's gate counts boxes, so criteria folded into one
// line are criteria it cannot check separately.
func TestNewBodyWritesOneBoxPerCriterion(t *testing.T) {
	got := NewBody("t", []string{"first", "  ", "second", "third"}, nil)
	for _, want := range []string{"- [ ] first\n", "- [ ] second\n", "- [ ] third\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if n := strings.Count(got, "- [ ] "); n != 3 {
		t.Errorf("got %d boxes, want 3 — a blank criterion must not become one:\n%s", n, got)
	}
	if strings.Contains(got, "checkable statement") {
		t.Errorf("the placeholder appeared although criteria were supplied:\n%s", got)
	}
}

// Notes reads back what AppendNote wrote, in that order, with a note that
// runs over several lines kept whole and a section after Progress left out.
func TestNotesReadsBackWhatAppendNoteWrote(t *testing.T) {
	body := "# T\n\n## Progress\n\n" +
		"- **2026-10-14 08:46 · Ada Lovelace** — first\n" +
		"- **2026-10-14 09:01 · Grace** — second line one\n- a list item\n```\ncode\n```\n" +
		"\n## Later\n\n- **2026-10-14 10:00 · X** — not a note\n"
	got := Notes(body)
	if len(got) != 2 {
		t.Fatalf("Notes = %+v; want 2", got)
	}
	if got[0].Author != "Ada Lovelace" || got[0].Text != "first" || got[0].At.Format("15:04") != "08:46" {
		t.Errorf("first = %+v", got[0])
	}
	if got[1].Text != "second line one\n- a list item\n```\ncode\n```" {
		t.Errorf("second kept %q", got[1].Text)
	}
	if got[1].Raw != "- **2026-10-14 09:01 · Grace** — "+got[1].Text {
		t.Errorf("raw = %q", got[1].Raw)
	}
	if n := Notes("## Progress\n- **2026-10-14 08:00 by hand** x\n"); len(n) != 1 || n[0].Author != "" || n[0].Text != n[0].Raw {
		t.Errorf("a dated note in another shape = %+v; want it kept whole, without an author", n)
	}
	if Notes("# no progress\n") != nil {
		t.Error("a body without Progress has notes")
	}
}
