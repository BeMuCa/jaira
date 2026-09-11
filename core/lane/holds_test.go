package lane

import (
	"strings"
	"testing"
)

func TestHoldsParsesAndDefaultsToUnlimited(t *testing.T) {
	l, err := parse([]byte("---\nid: closed\nterminal: true\nholds: 7\n---\n"), "closed.md", false)
	if err != nil {
		t.Fatal(err)
	}
	if l.Holds != 7 {
		t.Errorf("Holds = %d, want 7", l.Holds)
	}
	l, err = parse([]byte("---\nid: open\n---\n"), "open.md", false)
	if err != nil {
		t.Fatal(err)
	}
	if l.Holds != 0 {
		t.Errorf("absent holds = %d, want 0 (unlimited)", l.Holds)
	}
}

func TestHoldsRefusesWhatIsNotACount(t *testing.T) {
	for _, v := range []string{"ten", "-1"} {
		_, err := parse([]byte("---\nid: closed\nterminal: true\nholds: "+v+"\n---\n"), "closed.md", false)
		if err == nil || !strings.Contains(err.Error(), "holds") {
			t.Errorf("holds: %s was accepted, want a refusal naming holds (got %v)", v, err)
		}
	}
}

// A cap on a working lane would file unfinished tickets into the logbook —
// past the terminal-lane guard 'jaira logbook' enforces, and without their
// commits. Refused at parse, where the lane author can read why.
func TestHoldsRefusesANonTerminalLane(t *testing.T) {
	_, err := parse([]byte("---\nid: busy\nholds: 3\n---\n"), "busy.md", false)
	if err == nil || !strings.Contains(err.Error(), "terminal") {
		t.Errorf("holds on a non-terminal lane was accepted, want a refusal naming terminal (got %v)", err)
	}
	// holds: 0 stays legal anywhere — it declares nothing.
	if _, err := parse([]byte("---\nid: busy\nholds: 0\n---\n"), "busy.md", false); err != nil {
		t.Errorf("holds: 0 on a non-terminal lane refused: %v", err)
	}
}

// TestBuiltinDoneIsNotADoorway pins the shipped default: a ticket landing in
// done stays there until somebody files it. The rule lives in the lane file,
// not the code, so this is the test that notices if a doorway is ever declared
// there again.
//
// It was the other way round, and that is what this pins against: filing on
// entry meant one person finishing one ticket swept forty-nine of other
// people's into the logbook, in a commit named after a single handle. Filing is
// a statement about bookkeeping, made days later about a set somebody
// assembles; reaching a terminal lane is a statement about the work.
func TestBuiltinDoneIsNotADoorway(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("done")
	if !ok {
		t.Fatal("built-in done lane missing")
	}
	if l.LogbookOnEntry {
		t.Error("built-in done files on entry: finishing one ticket would file everybody's")
	}
	if l.Holds != 0 {
		t.Errorf("built-in done carries a cap (%d), which would file the oldest without being asked", l.Holds)
	}
}

// A doorway on a working lane would file unfinished tickets as done.
func TestLogbookOnEntryParsesAndRefusesANonTerminalLane(t *testing.T) {
	l, err := parse([]byte("---\nid: closed\nterminal: true\nlogbook-on-entry: true\n---\n"), "closed.md", false)
	if err != nil {
		t.Fatal(err)
	}
	if !l.LogbookOnEntry {
		t.Error("logbook-on-entry: true did not parse")
	}
	_, err = parse([]byte("---\nid: busy\nlogbook-on-entry: true\n---\n"), "busy.md", false)
	if err == nil || !strings.Contains(err.Error(), "logbook-on-entry") {
		t.Errorf("logbook-on-entry on a non-terminal lane was accepted, want a refusal (got %v)", err)
	}
}
