package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// updateTestStore builds a fresh, initialized board under its own JAIRA_HOME,
// the same isolation core/merge_test.go and browse_test.go use.
func updateTestStore(t *testing.T) (dir string, s *ticket.Store) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir, s
}

// setCurrent points release.Current at v for the duration of the test.
func setCurrent(t *testing.T, v string) {
	t.Helper()
	orig := release.Current
	release.Current = v
	t.Cleanup(func() { release.Current = orig })
}

// captureStdio swaps the real os.Stdout/os.Stderr for the duration of fn.
// nudgeIfStale, like bindDriverIfShared, writes straight to os.Stderr rather
// than through a cobra command's streams, so proving it never touches stdout
// needs the real file descriptors, not a buffer set via cmd.SetOut.
func captureStdio(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origOut, origErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW
	fn()
	outW.Close()
	errW.Close()
	os.Stdout, os.Stderr = origOut, origErr

	var ob, eb bytes.Buffer
	ob.ReadFrom(outR)
	eb.ReadFrom(errR)
	return ob.String(), eb.String()
}

func TestNudgeIfStaleNeverFiresOnDevBuild(t *testing.T) {
	setCurrent(t, "dev")
	_, s := updateTestStore(t)
	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" || stderr != "" {
		t.Fatalf("dev build nudged: stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestNudgeIfStaleSilentOnMatchingStamp(t *testing.T) {
	setCurrent(t, "1.0.0")
	_, s := updateTestStore(t)
	if err := release.Stamp(s.StateDir()); err != nil {
		t.Fatal(err)
	}
	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" || stderr != "" {
		t.Fatalf("matching stamp nudged: stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestNudgeIfStaleOnDifferingStampPrintsOneStderrLineOnly(t *testing.T) {
	setCurrent(t, "1.0.0")
	_, s := updateTestStore(t) // never stamped: an absent stamp counts as differing

	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" {
		t.Fatalf("nudge wrote to stdout: %q", stdout)
	}
	lines := strings.Split(strings.TrimRight(stderr, "\n"), "\n")
	if len(lines) != 1 || lines[0] == "" {
		t.Fatalf("stderr = %q, want exactly one line", stderr)
	}
	if !strings.Contains(stderr, "jaira update") {
		t.Errorf("nudge line = %q, want it to point at 'jaira update'", stderr)
	}
}

func TestNudgeIfStaleSkipsSilentlyWhenStateDirUnreadable(t *testing.T) {
	// Windows has no equivalent of a directory whose permission bits deny a
	// read: os.Chmod maps onto the read-only attribute and leaves the directory
	// listable, so the condition under test cannot be created there. The nudge
	// then fires correctly on an unstamped board, and the assertion below would
	// be measuring the setup rather than the behaviour.
	if runtime.GOOS == "windows" {
		t.Skip("chmod cannot make a directory unreadable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits do not block reads")
	}
	setCurrent(t, "1.0.0")
	_, s := updateTestStore(t)
	dir := s.StateDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" || stderr != "" {
		t.Fatalf("unreadable state dir nudged: stdout=%q stderr=%q", stdout, stderr)
	}
}

// runUpdate runs the real 'update' command against dir, the same path a user
// or agent invokes it through.
func runUpdate(t *testing.T, dir string, args ...string) (stdout string, err error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "update"}, args...))
	err = root.Execute()
	return out.String(), err
}

func TestUpdateStampsSoASecondNudgeCheckIsSilent(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, s := updateTestStore(t)

	if _, err := runUpdate(t, dir); err != nil {
		t.Fatalf("jaira update: %v", err)
	}
	if got := release.Stamped(s.StateDir()); got != "1.0.0" {
		t.Fatalf("Stamped after update = %q, want %q", got, "1.0.0")
	}

	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" || stderr != "" {
		t.Fatalf("nudge fired after update: stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestUpdateJSONCarriesTheDocumentedFieldsOnStdoutOnly(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)

	out, err := runUpdate(t, dir, "--json")
	if err != nil {
		t.Fatalf("jaira update --json: %v\n%s", err, out)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout was not a single JSON payload: %v\n%s", err, out)
	}
	for _, key := range []string{"version", "previous", "agent_notes", "notes"} {
		if _, ok := payload[key]; !ok {
			t.Errorf("payload missing %q: %#v", key, payload)
		}
	}
	if payload["version"] != "1.0.0" {
		t.Errorf("version = %#v, want %q", payload["version"], "1.0.0")
	}
	if payload["previous"] != "" {
		t.Errorf("previous = %#v, want empty on a never-stamped board", payload["previous"])
	}
}

// TestUpdateLeavesASharedBoardShared covers the measured defect: a board whose
// owner shared it by hand — the ignore line commented out, 'jaira share' never
// run — got "/.jaira/" appended again by every 'jaira update', so new tickets
// silently stopped reaching teammates. update writes the agent note only.
func TestUpdateLeavesASharedBoardShared(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)
	ignore := filepath.Join(dir, ".gitignore")
	before := "# jaira board — private to this machine. Run 'jaira share' to publish it.\n#/.jaira/\n"
	if err := os.WriteFile(ignore, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}

	if out, err := runUpdate(t, dir); err != nil {
		t.Fatalf("jaira update: %v\n%s", err, out)
	}

	after, err := os.ReadFile(ignore)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != before {
		t.Errorf("update changed .gitignore:\n before %q\n after  %q", before, string(after))
	}
}
