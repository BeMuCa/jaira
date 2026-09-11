package hook_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/hook"
)

func script(t *testing.T, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fixture is a shell script")
	}
	path := filepath.Join(t.TempDir(), "hook.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// The script is handed the ticket in the environment, which is the whole
// contract: a Slack webhook or an ntfy call needs the id, the title and where
// it went.
func TestTheScriptIsToldWhatHappened(t *testing.T) {
	out := filepath.Join(t.TempDir(), "seen")
	s := script(t, `printf '%s|%s|%s|%s|%s' "$JAIRA_EVENT" "$JAIRA_TICKET" "$JAIRA_TITLE" "$JAIRA_STATUS" "$JAIRA_ACTOR" > `+out)

	ok := hook.Run(s, hook.Event{
		Name: "claim", ID: "01AAA", Title: "session cookie dropped on 302",
		Status: "in-progress", Assignee: "grace", Actor: "ada", Root: t.TempDir(),
	})
	if !ok {
		t.Fatal("the hook did not run cleanly")
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("the script wrote nothing: %v", err)
	}
	want := "claim|01AAA|session cookie dropped on 302|in-progress|ada"
	if string(b) != want {
		t.Errorf("the script saw %q, want %q", b, want)
	}
}

// A hook nobody configured is not an event.
func TestAnEmptyHookDoesNothing(t *testing.T) {
	if hook.Run("", hook.Event{Name: "move"}) {
		t.Error("an empty hook claimed to have run")
	}
	if hook.Run("   ", hook.Event{Name: "move"}) {
		t.Error("a blank hook claimed to have run")
	}
}

// A broken or missing script is the user's business, never the command's
// failure.
func TestABrokenHookIsSilent(t *testing.T) {
	if hook.Run(filepath.Join(t.TempDir(), "nope.sh"), hook.Event{Name: "move"}) {
		t.Error("a missing script claimed to have run")
	}
	failing := script(t, "exit 3")
	if hook.Run(failing, hook.Event{Name: "move"}) {
		t.Error("a script exiting 3 was reported as clean")
	}
}

// A hook that hangs must not become the command's runtime.
func TestAHangingHookIsCutOff(t *testing.T) {
	if testing.Short() {
		t.Skip("this one waits for the timeout")
	}
	s := script(t, "sleep 30")
	done := make(chan bool, 1)
	go func() { done <- hook.Run(s, hook.Event{Name: "move"}) }()
	select {
	case ok := <-done:
		if ok {
			t.Error("a hanging hook was reported as clean")
		}
	case <-time.After(15 * time.Second):
		t.Error("the hook was not cut off")
	}
}

func TestEnvironmentIsNotReplaced(t *testing.T) {
	t.Setenv("JAIRA_TEST_MARKER", "kept")
	out := filepath.Join(t.TempDir(), "seen")
	s := script(t, `printf '%s' "$JAIRA_TEST_MARKER" > `+out)
	if !hook.Run(s, hook.Event{Name: "move", Root: t.TempDir()}) {
		t.Fatal("hook did not run")
	}
	b, _ := os.ReadFile(out)
	if strings.TrimSpace(string(b)) != "kept" {
		t.Errorf("the script lost the ambient environment: %q", b)
	}
}
