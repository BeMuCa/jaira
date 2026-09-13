//go:build unix

package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// exampleScriptFile writes what 'hook example' prints to a file and returns its
// path. What is under test is the text a user redirects into ~/.jaira, not a Go
// paraphrase of it.
func exampleScriptFile(t *testing.T) string {
	t.Helper()
	out, err := runCLI(t, t.TempDir(), "hook", "example")
	if err != nil {
		t.Fatalf("hook example: %v\n%s", err, out)
	}
	path := filepath.Join(t.TempDir(), "notify.sh")
	if err := os.WriteFile(path, []byte(out), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// runExampleHook runs the printed script the way jaira's own hook runner does —
// with the JAIRA_* variables in its environment and nothing else — and returns
// its exit code and stdout.
//
// It runs in a session of its own, which takes the controlling terminal away.
// That is what makes the bell observable at all: on a terminal the script
// writes to /dev/tty, precisely so a caller that discards stdout (jaira does)
// cannot swallow it, and a test capturing stdout would see nothing.
func runExampleHook(t *testing.T, script string, env []string) (int, string) {
	t.Helper()
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no POSIX shell on PATH")
	}
	c := exec.Command(sh, script)
	c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	c.Env = env
	c.Dir = t.TempDir()
	var stdout bytes.Buffer
	// No stdin: the script reads none, and a test that hands it one would hide
	// a script that started to.
	c.Stdin = strings.NewReader("")
	c.Stdout, c.Stderr = &stdout, io.Discard
	err = c.Run()
	var ee *exec.ExitError
	switch {
	case err == nil:
		return 0, stdout.String()
	case errors.As(err, &ee):
		return ee.ExitCode(), stdout.String()
	default:
		t.Fatalf("could not run the example: %v", err)
		return -1, ""
	}
}

// hookEnv is what jaira hands a hook, with the fields a test cares about
// overridden.
func hookEnv(event, status string) []string {
	return []string{
		// Deliberately no PATH: the example must work on a machine with no
		// delivery tool installed, which is the state everyone reading it for
		// the first time is in.
		"JAIRA_EVENT=" + event,
		"JAIRA_TICKET=01M2E86MS66ZZ58RRV3EY0A7DT",
		"JAIRA_TITLE=An example ticket",
		"JAIRA_STATUS=" + status,
		"JAIRA_ASSIGNEE=someone",
		"JAIRA_ACTOR=someone",
		"JAIRA_ROOT=/tmp",
	}
}

// The example exists to make one thing audible: a lane where nothing moves
// without a person. Every lane an agent works must stay silent, or the sound
// stops meaning anything and gets switched off.
func TestHookExampleSoundsOnlyForThePersonsLanes(t *testing.T) {
	script := exampleScriptFile(t)

	cases := []struct {
		status string
		// bells is how many times the script rings; 0 is silence.
		bells int
	}{
		{"human", 2},
		{"signoff", 2},
		{"done", 1},
		{"in-progress", 0},
		{"todo", 0},
		{"blocked", 0},
	}

	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			code, out := runExampleHook(t, script, hookEnv("move", tc.status))
			if code != 0 {
				t.Errorf("exit = %d, want 0 — a hook's failure is nobody's problem but it must not be one", code)
			}
			if got := strings.Count(out, "\a"); got != tc.bells {
				t.Errorf("%s rang %d times, want %d (output: %q)", tc.status, got, tc.bells, out)
			}
			if tc.bells == 0 && out != "" {
				t.Errorf("%s is an agent's lane and must be silent, got: %q", tc.status, out)
			}
		})
	}

	// And the discriminating half: asking for a person has to be audibly
	// different from a ticket simply finishing, or installing this buys the
	// reader nothing over a bell on every move.
	_, human := runExampleHook(t, script, hookEnv("move", "human"))
	_, done := runExampleHook(t, script, hookEnv("move", "done"))
	if strings.Count(human, "\a") <= strings.Count(done, "\a") {
		t.Errorf("a human lane is not audibly louder than a finish: %q vs %q", human, done)
	}
}

// A claim moves nothing and waits for nobody, and jaira calls the hook on it
// just the same.
func TestHookExampleIgnoresClaims(t *testing.T) {
	script := exampleScriptFile(t)
	code, out := runExampleHook(t, script, hookEnv("claim", "human"))
	if code != 0 || out != "" {
		t.Errorf("a claim produced exit %d and %q, want a silent 0", code, out)
	}
}

// The empty machine: no PATH at all, so no curl, no notify-send, no herdr. The
// script has to run through cleanly and do nothing rather than fail, because
// that is the state of every machine before the reader edits the delivery line.
func TestHookExampleRunsOnAMachineWithNothingInstalled(t *testing.T) {
	script := exampleScriptFile(t)
	for _, status := range []string{"human", "done", "in-progress"} {
		env := append(hookEnv("move", status), "PATH=")
		if code, _ := runExampleHook(t, script, env); code != 0 {
			t.Errorf("%s with an empty PATH exited %d, want 0", status, code)
		}
	}
}
