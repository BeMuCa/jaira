package cli

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/selfupdate"
)

// runSplit runs a command with stdout and stderr kept apart, which is the
// whole point of the tests below: stdout is what an agent parses.
func runSplit(t *testing.T, dir string, args ...string) (stdout, stderr string) {
	t.Helper()
	root := newRoot("test")
	var out, errOut bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs(append([]string{"-C", dir}, args...))
	if err := root.Execute(); err != nil {
		t.Fatalf("jaira %v: %v\n%s%s", args, err, out.String(), errOut.String())
	}
	return out.String(), errOut.String()
}

// cacheLatest writes a fresh update-check.json naming latest (empty: none).
func cacheLatest(t *testing.T, latest string) {
	t.Helper()
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if latest == "" {
		return
	}
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: latest}); err != nil {
		t.Fatal(err)
	}
}

func TestNextAndResumeNameANewerReleaseOnStderr(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)
	cacheLatest(t, "1.3.0")

	for _, args := range [][]string{{"next"}, {"next", "--json"}, {"next", "--per-lane", "--json"}, {"resume"}} {
		stdout, stderr := runSplit(t, dir, args...)
		for _, want := range []string{"1.3.0", "jaira self upgrade"} {
			if !bytes.Contains([]byte(stderr), []byte(want)) {
				t.Errorf("jaira %v: stderr = %q, want it to name %q", args, stderr, want)
			}
		}
		if bytes.Contains([]byte(stdout), []byte("1.3.0")) {
			t.Errorf("jaira %v: stdout = %q, the release line belongs on stderr", args, stdout)
		}
	}

	stdout, _ := runSplit(t, dir, "next", "--json")
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("next --json stdout is not JSON: %v\n%s", err, stdout)
	}
}

func TestNoReleaseLineWhenUpToDateUnknownOrDev(t *testing.T) {
	cases := []struct {
		name, current, latest string
	}{
		{"same version", "1.0.0", "1.0.0"},
		{"older cache", "1.0.0", "0.9.0"},
		{"never checked", "1.0.0", ""},
		{"dev build", "dev", "1.3.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			setCurrent(t, c.current)
			dir, _ := updateTestStore(t)
			cacheLatest(t, c.latest)
			for _, args := range [][]string{{"next"}, {"resume"}} {
				if _, stderr := runSplit(t, dir, args...); bytes.Contains([]byte(stderr), []byte("jaira self upgrade")) {
					t.Errorf("jaira %v: stderr = %q, want no release line", args, stderr)
				}
			}
		})
	}
}
