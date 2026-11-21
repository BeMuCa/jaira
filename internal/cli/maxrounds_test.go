package cli

import (
	"strings"
	"testing"
)

// newLoopTicket is a ticket on a board that has adopted the shipped testing
// lane, whose max-rounds is 3, sitting in that lane.
func newLoopTicket(t *testing.T, terminal bool) (dir, h string) {
	t.Helper()
	was := atATerminal
	atATerminal = func() bool { return terminal }
	t.Cleanup(func() { atATerminal = was })

	dir, h = newModeTicket(t)
	if out, err := runCLI(t, dir, "lanes", "add", "testing"); err != nil {
		t.Fatalf("lanes add testing: %v\n%s", err, out)
	}
	return dir, h
}

func roundsLeft(t *testing.T, dir, h string) any {
	t.Helper()
	return jsonCLI(t, dir, "show", h, "--for-lane", "testing")["rounds_left"]
}

// Three send-backs in a row from testing go through, the fourth is refused
// with exit 3 naming the human lane, and --for-lane counts down beside them.
func TestMaxRoundsRefusesTheFourthSendBack(t *testing.T) {
	dir, h := newLoopTicket(t, false)
	for i := 0; i < 3; i++ {
		if out, err := runCLI(t, dir, "move", h, "--to", "testing", "--force"); err != nil {
			t.Fatalf("to testing: %v\n%s", err, out)
		}
		if got := roundsLeft(t, dir, h); got != float64(3-i) {
			t.Errorf("before send-back %d: rounds_left %v, want %d", i+1, got, 3-i)
		}
		if out, err := runCLI(t, dir, "move", h, "--to", "in-progress"); err != nil {
			t.Fatalf("send-back %d refused: %v\n%s", i+1, err, out)
		}
	}
	if out, err := runCLI(t, dir, "move", h, "--to", "testing", "--force"); err != nil {
		t.Fatalf("to testing: %v\n%s", err, out)
	}
	if got := roundsLeft(t, dir, h); got != float64(0) {
		t.Errorf("rounds_left %v after three, want 0", got)
	}
	plain, _ := runCLI(t, dir, "show", h, "--for-lane", "testing")
	if !strings.Contains(plain, "max-rounds 3: 0 send-back(s) left") {
		t.Errorf("the plain brief does not count down:\n%s", plain)
	}

	out, err := runCLI(t, dir, "move", h, "--to", "in-progress")
	if exitCode(err) != ExitValidation {
		t.Fatalf("fourth send-back: exit %d (%v), want %d\n%s", exitCode(err), err, ExitValidation, out)
	}
	if msg := out + err.Error(); !strings.Contains(msg, "human") || !strings.Contains(msg, "max-rounds") {
		t.Errorf("refusal does not name max-rounds and the human lane: %s", msg)
	}
	if out, err := runCLI(t, dir, "move", h, "--to", "in-progress", "--dry-run"); exitCode(err) != ExitValidation {
		t.Errorf("the dry run predicts exit %d, want %d\n%s", exitCode(err), ExitValidation, out)
	}
	if out, err := runCLI(t, dir, "move", h, "--to", "in-progress", "--force"); err != nil {
		t.Errorf("--force did not override: %v\n%s", err, out)
	}
}

// A person at a terminal sending it back is the decision the limit asks for.
func TestMaxRoundsLetsAPersonSendItBack(t *testing.T) {
	dir, h := newLoopTicket(t, true)
	for i := 0; i < 5; i++ {
		if out, err := runCLI(t, dir, "move", h, "--to", "testing", "--force"); err != nil {
			t.Fatalf("to testing: %v\n%s", err, out)
		}
		if out, err := runCLI(t, dir, "move", h, "--to", "in-progress"); err != nil {
			t.Fatalf("a person's send-back %d refused: %v\n%s", i+1, err, out)
		}
	}
}

// A lane without max-rounds carries no rounds keys at all: a rounds_left
// beside it would read as a limit of zero.
func TestForLaneWithoutMaxRoundsCarriesNoRounds(t *testing.T) {
	dir, h := newModeTicket(t)
	got := jsonCLI(t, dir, "show", h, "--for-lane", "in-progress")
	if _, ok := got["rounds_left"]; ok {
		t.Errorf("in-progress declares no max-rounds but hands on rounds_left %v", got["rounds_left"])
	}
	if _, ok := got["max_rounds"]; ok {
		t.Errorf("in-progress hands on max_rounds %v", got["max_rounds"])
	}
}
