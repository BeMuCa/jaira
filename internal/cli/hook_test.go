package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// stopHookCommand decodes 'hook print' the way Claude Code would read it, and
// returns the one shell command it registers.
func stopHookCommand(t *testing.T) string {
	t.Helper()
	out, err := runCLI(t, t.TempDir(), "hook", "print")
	if err != nil {
		t.Fatalf("hook print: %v\n%s", err, out)
	}
	var cfg struct {
		Hooks struct {
			Stop []struct {
				Hooks []struct {
					Type    string `json:"type"`
					Command string `json:"command"`
				} `json:"hooks"`
			} `json:"Stop"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(out), &cfg); err != nil {
		t.Fatalf("output is not JSON a settings file could hold: %v\n%s", err, out)
	}
	if len(cfg.Hooks.Stop) != 1 || len(cfg.Hooks.Stop[0].Hooks) != 1 {
		t.Fatalf("expected exactly one Stop hook, got %#v", cfg.Hooks.Stop)
	}
	h := cfg.Hooks.Stop[0].Hooks[0]
	if h.Type != "command" {
		t.Errorf("hook type = %q, want \"command\"", h.Type)
	}
	return h.Command
}

// The snippet is pasted into someone's settings.json and then never looked at
// again, so its text is a contract. Asserted verbatim: an edit to the constant
// that changes what the hook does has to be made deliberately, in two places.
func TestHookPrintIsTheStopHookVerbatim(t *testing.T) {
	want := `grep -q '"stop_hook_active": *true' && exit 0; jaira next --per-lane --json 2>/dev/null | grep -q '"agentic": *true' || exit 0; echo 'jaira: an agentic lane still has work waiting. Run jaira next --per-lane, work the lane, and move the ticket on before stopping.' >&2; exit 2`
	if got := stopHookCommand(t); got != want {
		t.Errorf("hook command changed\n got: %s\nwant: %s", got, want)
	}
}

// Exit 2 is the only code Claude Code treats as blocking a stop — exit 1 is a
// non-blocking error and lets the session end, which would make the hook look
// installed while enforcing nothing.
func TestStopHookBlocksWithExitTwo(t *testing.T) {
	cmd := stopHookCommand(t)
	if !strings.HasSuffix(cmd, "exit 2") {
		t.Errorf("hook does not end by exiting 2:\n%s", cmd)
	}
	if !strings.Contains(cmd, ">&2") {
		t.Errorf("hook does not write its reason to stderr, so Claude is told nothing:\n%s", cmd)
	}
}

// Three ways this hook must get out of the way, each of which has cost somebody
// a wedged session somewhere: a stop it has already blocked once, a machine
// with no jq, and a directory with no board.
func TestStopHookFailsOpen(t *testing.T) {
	cmd := stopHookCommand(t)
	if !strings.Contains(cmd, "stop_hook_active") {
		t.Errorf("hook does not check stop_hook_active and can block a condition it cannot resolve:\n%s", cmd)
	}
	if strings.Contains(cmd, "jq") {
		t.Errorf("hook depends on jq, which is not installed everywhere:\n%s", cmd)
	}
	if !strings.Contains(cmd, "2>/dev/null") || !strings.Contains(cmd, "|| exit 0") {
		t.Errorf("hook does not let a board it cannot read through:\n%s", cmd)
	}
}

// The read the hook branches on is 'next --per-lane --json' and the field is
// "agentic": if either is renamed, the hook silently stops enforcing anything,
// so the shape it depends on is asserted here as well as in perlane_test.go.
func TestStopHookReadsThePerLaneAgenticFlag(t *testing.T) {
	cmd := stopHookCommand(t)
	if !strings.Contains(cmd, "jaira next --per-lane --json") {
		t.Errorf("hook does not read 'jaira next --per-lane --json':\n%s", cmd)
	}
	if !strings.Contains(cmd, `'"agentic": *true'`) {
		t.Errorf("hook does not match the per-lane agentic flag:\n%s", cmd)
	}

	// And the pattern against the real bytes: the hook greps rather than parses,
	// so a change to how the payload is encoded breaks it just as a renamed
	// field would.
	dir, id := movableTicket(t)
	if out, err := runCLI(t, dir, "move", id, "--to", "todo"); err != nil {
		t.Fatalf("move to todo: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "move", id, "--to", "in-progress"); err != nil {
		t.Fatalf("move to in-progress: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "next", "--per-lane", "--json")
	if err != nil {
		t.Fatalf("next --per-lane --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"agentic": true`) {
		t.Errorf("per-lane output no longer carries the bytes the hook greps for:\n%s", out)
	}
}
