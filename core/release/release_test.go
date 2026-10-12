package release

import (
	"path/filepath"
	"reflect"
	"testing"
)

const fixture = `<!-- format notes, ignored -->

## 0.1.0

- second release change one
- second release change two

Some prose that is not a change line, also ignored.

## 0.0.1

- first release change
`

func TestParseNotesKeepsOnlyDashLines(t *testing.T) {
	entries := parseNotes(fixture)
	want := []Entry{
		{Version: "0.1.0", Changes: []string{"second release change one", "second release change two"}},
		{Version: "0.0.1", Changes: []string{"first release change"}},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("parseNotes = %#v, want %#v", entries, want)
	}
}

func TestSinceEmptyStampReturnsEverything(t *testing.T) {
	all := parseNotes(fixture)
	got := sinceEntries(all, "")
	if !reflect.DeepEqual(got, all) {
		t.Fatalf("sinceEntries(all, \"\") = %#v, want everything", got)
	}
}

func TestSinceTopEntryReturnsNothing(t *testing.T) {
	all := parseNotes(fixture)
	got := sinceEntries(all, "0.1.0")
	if len(got) != 0 {
		t.Fatalf("sinceEntries(all, top) = %#v, want none", got)
	}
}

func TestSinceSecondEntryReturnsOnlyEntriesAboveIt(t *testing.T) {
	all := parseNotes(fixture)
	got := sinceEntries(all, "0.0.1")
	want := all[:1]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sinceEntries(all, second) = %#v, want %#v", got, want)
	}
}

func TestSinceUnknownStampReturnsEverything(t *testing.T) {
	all := parseNotes(fixture)
	for _, stamp := range []string{"dev", "nonexistent"} {
		got := sinceEntries(all, stamp)
		if !reflect.DeepEqual(got, all) {
			t.Errorf("sinceEntries(all, %q) = %#v, want everything", stamp, got)
		}
	}
}

func TestStampedOnMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	if got := Stamped(dir); got != "" {
		t.Fatalf("Stamped(no file) = %q, want empty", got)
	}
}

func TestStampRoundTripsThroughStamped(t *testing.T) {
	orig := Current
	Current = "1.2.3"
	t.Cleanup(func() { Current = orig })

	dir := filepath.Join(t.TempDir(), "nested", "state")
	if err := Stamp(dir); err != nil {
		t.Fatalf("Stamp: %v", err)
	}
	if got := Stamped(dir); got != Current {
		t.Fatalf("Stamped(after Stamp) = %q, want %q", got, Current)
	}
}

func TestEmbeddedNotesParseToAtLeastOneRealEntry(t *testing.T) {
	notes := Notes()
	if len(notes) == 0 {
		t.Fatal("Notes() is empty — NOTES.md failed to parse into any entry")
	}
	for _, e := range notes {
		if e.Version != "" && len(e.Changes) > 0 {
			return
		}
	}
	// Entry zero is deliberately not the one checked: it is the note files'
	// Unreleased, which a release build may not have at all. What must hold is
	// that some released section carries changes.
	t.Fatalf("no entry in NOTES.md has both a version and changes: %#v", notes)
}

// A heading with nothing under it yet must stay an entry of its own rather
// than be dropped or let the next section's changes fall into it.
func TestParseNotesKeepsAnEmptyLeadingSection(t *testing.T) {
	entries := parseNotes("## Unreleased\n\n## 0.1.0\n- a released change\n")
	want := []Entry{
		{Version: "Unreleased"},
		{Version: "0.1.0", Changes: []string{"a released change"}},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("parseNotes = %#v, want %#v", entries, want)
	}
}

// A heading with no points under it is an entry the parser keeps, but nobody
// can read it. It must not reach Since()'s result — and on a board stamped
// with the version that was just cut, the result must be empty, so "jaira
// update" can say nothing has changed instead of printing a bare heading.
func TestSinceDropsEntriesWithoutChanges(t *testing.T) {
	all := parseNotes("## Unreleased\n\n## 0.1.0\n- a released change\n\n## 0.0.1\n- an older change\n")

	if got := sinceEntries(all, "0.1.0"); len(got) != 0 {
		t.Fatalf("sinceEntries(all, cut version) = %#v, want none", got)
	}
	want := []Entry{
		{Version: "0.1.0", Changes: []string{"a released change"}},
		{Version: "0.0.1", Changes: []string{"an older change"}},
	}
	if got := sinceEntries(all, ""); !reflect.DeepEqual(got, want) {
		t.Fatalf("sinceEntries(all, \"\") = %#v, want %#v", got, want)
	}
}
