package lane

import (
	"strings"
	"testing"
)

func TestHoldsParsesAndDefaultsToUnlimited(t *testing.T) {
	l, err := parse([]byte("---\nid: closed\nholds: 7\n---\n"), "closed.md", false)
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
		_, err := parse([]byte("---\nid: closed\nholds: "+v+"\n---\n"), "closed.md", false)
		if err == nil || !strings.Contains(err.Error(), "holds") {
			t.Errorf("holds: %s was accepted, want a refusal naming holds (got %v)", v, err)
		}
	}
}

// TestBuiltinDoneDeclaresTheCap pins the shipped default: done keeps ten. The
// number lives in the lane file, not the code, so this is the test that notices
// if the declaration is ever lost.
func TestBuiltinDoneDeclaresTheCap(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("done")
	if !ok {
		t.Fatal("built-in done lane missing")
	}
	if l.Holds != 10 {
		t.Errorf("built-in done Holds = %d, want 10", l.Holds)
	}
}
