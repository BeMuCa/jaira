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
	e := notes[0]
	if e.Version == "" {
		t.Error("first entry has an empty Version")
	}
	if len(e.Changes) == 0 {
		t.Error("first entry has no Changes")
	}
}
