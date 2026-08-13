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
	got := NewBody("Fix the thing", "", []BodyOption{
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
		"\n## Notes\n\n"
	if got != want {
		t.Errorf("NewBody regressed from the baseline shape:\ngot:  %q\nwant: %q", got, want)
	}
}

// TestNewBodyTicksNamedOptions asserts a ticked option renders as "- [x]"
// and everything else stays unticked.
func TestNewBodyTicksNamedOptions(t *testing.T) {
	got := NewBody("t", "", []BodyOption{
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
	got := NewBody("t", "ships behind a flag", nil)
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
	got := NewBody("t", "", nil)
	if !strings.Contains(got, "## Options\n\n\n## Plan") {
		t.Errorf("expected an empty Options section, got:\n%s", got)
	}
}
