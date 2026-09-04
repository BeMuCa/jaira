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

// TestBuiltinDoneIsADoorway pins the shipped default: a ticket landing in done
// goes straight to the logbook. The rule lives in the lane file, not the code,
// so this is the test that notices if the declaration is ever lost.
func TestBuiltinDoneIsADoorway(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("done")
	if !ok {
		t.Fatal("built-in done lane missing")
	}
	if !l.LogbookOnEntry {
		t.Error("built-in done does not file on entry")
	}
	if l.Holds != 0 {
		t.Errorf("built-in done still carries a cap (%d) beside the doorway", l.Holds)
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
