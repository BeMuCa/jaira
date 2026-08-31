package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	// And the discriminating half, which is where the hook's correctness
	// actually lives: a board whose only waiting work is in a lane no agent
	// works must not produce those bytes, or the hook blocks every stop
	// forever and looks like it is working.
	quiet, id2 := movableTicket(t)
	if out, err := runCLI(t, quiet, "move", id2, "--to", "todo"); err != nil {
		t.Fatalf("move to todo: %v\n%s", err, out)
	}
	out, err = runCLI(t, quiet, "next", "--per-lane", "--json")
	if err != nil {
		t.Fatalf("next --per-lane --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"agentic": false`) {
		t.Fatalf("expected the non-agentic lane to be reported at all:\n%s", out)
	}
	if strings.Contains(out, `"agentic": true`) {
		t.Errorf("per-lane says a lane is agentic on a board whose only work is not:\n%s", out)
	}
}

// perLaneAgentic and perLaneQuiet are the two answers the hook branches on,
// encoded the way 'next --per-lane --json' encodes them. They are canned rather
// than produced by the real command because this test is about the shell
// command's behaviour; that the real command still emits these bytes is what
// TestStopHookReadsThePerLaneAgenticFlag asserts.
const perLaneAgentic = `{
  "lanes": [
    {
      "lane": "critique",
      "name": "Critique",
      "agentic": true,
      "model_tier": "strong",
      "waiting": 2,
      "ticket": {}
    }
  ],
  "count": 1
}`

const perLaneQuiet = `{
  "lanes": [
    {
      "lane": "todo",
      "name": "Todo",
      "agentic": false,
      "waiting": 1,
      "ticket": {}
    }
  ],
  "count": 1
}`

// runStopHook executes the snippet's own shell command, so what is under test is
// the text a user pastes into settings.json rather than a Go paraphrase of it.
//
// jaira is stubbed by a script on PATH that answers with perLane, which keeps
// the test off a real board and lets it produce the one answer it wants. An
// empty perLane installs no stub at all, which is the "no jaira here" case.
func runStopHook(t *testing.T, perLane, stdin string) (int, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the snippet is a POSIX shell command")
	}
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no POSIX shell on PATH")
	}
	bin := t.TempDir()
	if perLane != "" {
		stub := "#!/bin/sh\ncat <<'JSON'\n" + perLane + "\nJSON\n"
		if err := os.WriteFile(filepath.Join(bin, "jaira"), []byte(stub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	c := exec.Command(sh, "-c", stopHookCommand(t))
	// bin first so the stub wins over any jaira actually installed here; the
	// rest of PATH is still needed for grep.
	c.Env = []string{"PATH=" + bin + string(os.PathListSeparator) + os.Getenv("PATH")}
	c.Dir = t.TempDir()
	c.Stdin = strings.NewReader(stdin)
	var stderr bytes.Buffer
	c.Stdout, c.Stderr = io.Discard, &stderr
	err = c.Run()
	var ee *exec.ExitError
	switch {
	case err == nil:
		return 0, stderr.String()
	case errors.As(err, &ee):
		return ee.ExitCode(), stderr.String()
	default:
		t.Fatalf("could not run the snippet: %v", err)
		return -1, ""
	}
}

// The snippet is the deliverable, and quoting it back is not evidence that it
// runs. Every branch is executed here against a stubbed board.
func TestStopHookSnippetRuns(t *testing.T) {
	cases := []struct {
		name     string
		perLane  string
		stdin    string
		wantCode int
		wantErr  string
	}{{
		name: "a stop it has already blocked once goes through",
		// Agentic work is waiting and it still must not block: without this the
		// session is stuck on a condition the hook cannot resolve itself.
		perLane: perLaneAgentic, stdin: `{"session_id":"s","stop_hook_active": true}`,
		wantCode: 0,
	}, {
		name:    "no agentic lane with work goes through",
		perLane: perLaneQuiet, stdin: `{"session_id":"s","stop_hook_active": false}`,
		wantCode: 0,
	}, {
		name:    "an agentic lane with work blocks, and says why",
		perLane: perLaneAgentic, stdin: `{"session_id":"s","stop_hook_active": false}`,
		wantCode: 2, wantErr: "an agentic lane still has work waiting",
	}, {
		name:    "no jaira on PATH goes through",
		perLane: "", stdin: `{"session_id":"s","stop_hook_active": false}`,
		wantCode: 0,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, stderr := runStopHook(t, tc.perLane, tc.stdin)
			if code != tc.wantCode {
				t.Errorf("exit = %d, want %d (stderr: %s)", code, tc.wantCode, stderr)
			}
			if tc.wantErr == "" && stderr != "" {
				t.Errorf("expected nothing on stderr, got: %s", stderr)
			}
			if tc.wantErr != "" && !strings.Contains(stderr, tc.wantErr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tc.wantErr)
			}
		})
	}
}
