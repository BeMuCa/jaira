package cli

import (
	"strings"
	"testing"
)

const notesBody = `# t

## Definition of Done

- [ ] one criterion, which must never be trimmed

## Progress

- **2026-08-01 10:00 · a** — first
- **2026-08-02 10:00 · a** — second
- **2026-08-03 10:00 · a** — third
  and a continuation line belonging to the third
- **2026-08-04 10:00 · a** — fourth

## Notes

- a line under another heading, which is not a progress note
`

// The newest notes are the ones that say where the work stands, and they sit
// furthest from the reader on a ticket that has been running for weeks.
func TestTrimNotesKeepsTheNewest(t *testing.T) {
	out, hidden := trimNotes(notesBody, 2)
	if hidden != 2 {
		t.Errorf("hidden = %d, want 2", hidden)
	}
	for _, gone := range []string{"— first", "— second"} {
		if strings.Contains(out, gone) {
			t.Errorf("%q survived the trim:\n%s", gone, out)
		}
	}
	for _, kept := range []string{"— third", "— fourth", "and a continuation line"} {
		if !strings.Contains(out, kept) {
			t.Errorf("%q was trimmed away:\n%s", kept, out)
		}
	}
}

// Only the Progress section is a log. Everything else is criteria, plan or
// outcome, and losing a line of it would be losing the ticket.
func TestTrimNotesLeavesOtherSectionsAlone(t *testing.T) {
	out, _ := trimNotes(notesBody, 1)
	if !strings.Contains(out, "one criterion, which must never be trimmed") {
		t.Errorf("the definition of done was trimmed:\n%s", out)
	}
	if !strings.Contains(out, "a line under another heading") {
		t.Errorf("a later section was trimmed:\n%s", out)
	}
	if !strings.Contains(out, "## Progress") {
		t.Errorf("the Progress heading was removed:\n%s", out)
	}
}

// 0 is "show everything", and it is the default everywhere: this board's promise
// is that nothing is lost, so a note is hidden only when the reader asks.
func TestTrimNotesShowsAllByDefault(t *testing.T) {
	out, hidden := trimNotes(notesBody, 0)
	if hidden != 0 || out != notesBody {
		t.Errorf("hidden = %d and the body changed with n = 0", hidden)
	}
	if out, hidden := trimNotes(notesBody, 9); hidden != 0 || out != notesBody {
		t.Errorf("asking for more notes than exist trimmed something: hidden = %d", hidden)
	}
}

// The count is what stops a truncated read from looking complete.
func TestShowNotesLastReportsWhatItHid(t *testing.T) {
	dir, id := dodTestStore(t)
	for _, n := range []string{"first", "second", "third"} {
		if out, err := runCLI(t, dir, "note", id, n); err != nil {
			t.Fatalf("note: %v\n%s", err, out)
		}
	}
	out, err := runCLI(t, dir, "show", id, "--notes-last", "1")
	if err != nil {
		t.Fatalf("show --notes-last: %v\n%s", err, out)
	}
	if !strings.Contains(out, "third") || strings.Contains(out, "first") {
		t.Errorf("wrong notes shown:\n%s", out)
	}
	if !strings.Contains(out, "2 earlier note(s) hidden") {
		t.Errorf("the hidden count was not reported:\n%s", out)
	}
	// Default is unchanged: everything.
	full, err := runCLI(t, dir, "show", id)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(full, "first") || strings.Contains(full, "hidden") {
		t.Errorf("the default show is no longer complete:\n%s", full)
	}
}

// A listing row said who owns a ticket and nothing about how far it had got.
func TestListShowsTheDoDCounter(t *testing.T) {
	dir, id := dodTestStore(t)
	if out, err := runCLI(t, dir, "dod", id, "1", "--done"); err != nil {
		t.Fatalf("setup: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "list")
	if err != nil {
		t.Fatalf("list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "[DoD 1/2]") {
		t.Errorf("no DoD counter in the listing:\n%s", out)
	}

	// A superseded item counts as settled, the same way it does for the gate
	// and for the board — otherwise a ticket that can be completed reads as 1/2
	// for good.
	if out, err := runCLI(t, dir, "dod", id, "2", "--superseded"); err != nil {
		t.Fatalf("supersede: %v\n%s", err, out)
	}
	if out, err = runCLI(t, dir, "list"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "[DoD 2/2]") {
		t.Errorf("a superseded item is not counted as settled:\n%s", out)
	}
}
