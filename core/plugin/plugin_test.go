package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// install puts a shell script where an install would leave jaira-<id>, and
// the manifest beside it when events is not nil.
func install(t *testing.T, home, id, body string, events []string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fixture is a shell script")
	}
	dir := filepath.Join(home, "plugins", id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	prog := filepath.Join(dir, "jaira-"+id)
	if err := os.WriteFile(prog, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if events != nil {
		raw, _ := json.Marshal(map[string]any{"id": id, "events": events})
		if err := os.WriteFile(filepath.Join(dir, manifestFile), raw, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return prog
}

func home(t *testing.T) string {
	t.Helper()
	h := t.TempDir()
	t.Setenv("JAIRA_HOME", h)
	t.Setenv(InEventVar, "")
	return h
}

// Only what sits in the plugins directory under a well-formed id is a
// plugin; a name with a path in it never reaches the filesystem.
func TestProgramIsOnlyAnInstalledPlugin(t *testing.T) {
	h := home(t)
	prog := install(t, h, "jira", "exit 0", nil)
	if got := Program("jira"); got != prog {
		t.Errorf("Program(jira) = %q, want %q", got, prog)
	}
	for _, id := range []string{"youtrack", "../jira", "Jira", "", "-x"} {
		if got := Program(id); got != "" {
			t.Errorf("Program(%q) = %q, want nothing", id, got)
		}
	}
}

// The plugin gets every argument, the three variables and its own exit code
// back.
func TestRunPassesArgumentsEnvironmentAndExitCode(t *testing.T) {
	h := home(t)
	out := filepath.Join(t.TempDir(), "seen")
	prog := install(t, h, "jira", `printf '%s|' "$@" > `+out+`; echo "$JAIRA_ROOT $JAIRA_BOARD $JAIRA_VERSION" >> `+out+`; exit 7`, nil)
	code, err := Run(prog, []string{"pull", "--board", "ops", "a b"}, Env{Root: "/r", Board: "ops", Version: "1.2.3"})
	if err != nil {
		t.Fatal(err)
	}
	if code != 7 {
		t.Errorf("exit code %d, want 7", code)
	}
	raw, _ := os.ReadFile(out)
	if got, want := string(raw), "pull|--board|ops|a b|/r ops 1.2.3\n"; got != want {
		t.Errorf("plugin saw %q, want %q", got, want)
	}
}

// Subscribers reads the manifest: a plugin without one, or one that names
// another id, hears nothing, and the order is the ids'.
func TestSubscribersComeFromTheManifest(t *testing.T) {
	h := home(t)
	install(t, h, "zeta", "exit 0", []string{"move"})
	install(t, h, "alpha", "exit 0", []string{"move", "set"})
	install(t, h, "silent", "exit 0", nil)
	install(t, h, "liar", "exit 0", []string{"move"})
	if err := os.WriteFile(filepath.Join(h, "plugins", "liar", manifestFile), []byte(`{"id":"other","events":["move"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(Subscribers("move"), ","); got != "alpha,zeta" {
		t.Errorf("move subscribers %q, want alpha,zeta", got)
	}
	if got := strings.Join(Subscribers("set"), ","); got != "alpha" {
		t.Errorf("set subscribers %q, want alpha", got)
	}
}

// An event hands the plugin the payload on stdin, `event <name>` as its
// arguments, the board's directory and JAIRA_IN_EVENT=1.
func TestFireHandsThePluginTheEvent(t *testing.T) {
	h := home(t)
	root := t.TempDir()
	out := filepath.Join(t.TempDir(), "seen")
	install(t, h, "jira", `{ echo "$1 $2 $JAIRA_IN_EVENT $JAIRA_BOARD $(pwd)"; cat; } > `+out+`; exit 9`, []string{"move"})
	Fire(Subscribers("move"), "move", []byte(`{"id":"X"}`), Env{Root: root, Board: "default", Version: "dev"})
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	wantDir, _ := filepath.EvalSymlinks(root)
	if got, want := string(raw), "event move 1 default "+wantDir+"\n"+`{"id":"X"}`; got != want {
		t.Errorf("plugin saw %q, want %q", got, want)
	}
}

// A write made from inside an event fires nothing, or two plugins answering
// each other never stop.
func TestFireInsideAnEventCallsNoPlugin(t *testing.T) {
	h := home(t)
	out := filepath.Join(t.TempDir(), "seen")
	install(t, h, "jira", "touch "+out, []string{"set"})
	t.Setenv(InEventVar, "1")
	if ids := Subscribers("set"); len(ids) != 0 {
		t.Fatalf("inside an event the subscribers are %v, want none", ids)
	}
	Fire(Subscribers("set"), "set", nil, Env{Root: t.TempDir()})
	if _, err := os.Stat(out); err == nil {
		t.Fatal("a plugin ran from inside an event")
	}
}

// A plugin that hangs is stopped at the timeout, and the next one still runs.
func TestFireStopsAHangingPlugin(t *testing.T) {
	h := home(t)
	out := filepath.Join(t.TempDir(), "seen")
	install(t, h, "a-slow", "exec sleep 30", []string{"claim"})
	install(t, h, "b-next", "touch "+out, []string{"claim"})
	start := time.Now()
	Fire(Subscribers("claim"), "claim", nil, Env{Root: t.TempDir()})
	if took := time.Since(start); took > eventTimeout+3*time.Second {
		t.Fatalf("Fire took %v; the timeout is %v", took, eventTimeout)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal("the plugin after a hanging one did not run")
	}
}
