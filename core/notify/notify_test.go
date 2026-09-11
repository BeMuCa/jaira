package notify

import (
	"strings"
	"testing"
)

// The quoting is the only part with teeth: a ticket title is arbitrary text
// from a person, and it is handed to osascript and to PowerShell as source.
func TestQuotingClosesTheHoleATitleCouldOpen(t *testing.T) {
	for _, in := range []string{
		`plain`,
		`with "quotes"`,
		`back\slash`,
		`" & do something else; echo`,
	} {
		got := quote(in)
		if !strings.HasPrefix(got, `"`) || !strings.HasSuffix(got, `"`) {
			t.Errorf("quote(%q) = %q, not a literal", in, got)
		}
		// No bare quote may survive inside the literal.
		if inner := got[1 : len(got)-1]; strings.Contains(strings.ReplaceAll(inner, `\"`, ""), `"`) {
			t.Errorf("quote(%q) = %q leaves a bare quote", in, got)
		}
	}
}

func TestPowerShellQuotingDoublesTheApostrophe(t *testing.T) {
	if got := psQuote("it's"); got != "'it''s'" {
		t.Errorf("psQuote = %q", got)
	}
}

// Send must never fail loudly, on any machine, including one with no notifier
// at all — a headless CI runner is exactly the case.
func TestSendNeverPanicsWithoutANotifier(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if Send("jaira", "a ticket was assigned to you") {
		t.Error("something claimed to notify with an empty PATH")
	}
	if Available() {
		t.Error("Available said yes with an empty PATH")
	}
}
