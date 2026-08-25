package tui

import (
	"strings"
	"testing"
)

// The board's whole vocabulary for state is a glyph plus a word on a card, and
// until this existed the only place that explained any of it was the source.
func TestHelpExplainsEveryCardMark(t *testing.T) {
	m := newTestModel(t, 150, 60)
	m.mode = modeHelp
	out := stripANSI(m.render())

	for _, want := range []string{
		"○ spec", "■ blocked", "▲ asks", "◆ sign off", "◇ unworked", "✎ name",
		"Plan 2/5", "DoD 2/5", "✓ 3", "@name",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the help does not explain %q:\n%s", want, out)
		}
	}
	// Each one is explained, not merely listed.
	for _, want := range []string{
		"not specified enough", "waiting on a ticket", "question is waiting",
		"waiting for a person to accept", "has not produced what it declares",
		"somebody else wrote this ticket last",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the help lists a mark without saying what it means: %q missing", want)
		}
	}
}
