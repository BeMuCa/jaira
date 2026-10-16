package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/selfupdate"
	"github.com/BeMuCa/jaira/core/ticket"
)

// releaseLineWants is what the stderr line must carry: the release, the two
// previews the agent shows the user before it asks, that it asks before
// anything is written, and the two commands a yes runs.
var releaseLineWants = []string{
	"1.3.0", "jaira self upgrade --check", "jaira update --roles --lanes --dry-run", "show the user",
	"ask the user", "only on yes run 'jaira self upgrade', then 'jaira update --roles --lanes'",
}

// releaseLineOf is the one stderr line that names the release: what the agent
// must do stands on that line whole, so it reads as one instruction.
func releaseLineOf(stderr string) string {
	for _, l := range strings.Split(stderr, "\n") {
		if strings.Contains(l, "is released") {
			return l
		}
	}
	return ""
}

// backlogTicket files one ticket in backlog and returns its handle, so claim
// and move have something to act on.
func backlogTicket(t *testing.T, s *ticket.Store) string {
	t.Helper()
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return ticket.Handle(tk.ID)
}

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
		line := releaseLineOf(stderr)
		for _, want := range releaseLineWants {
			if !strings.Contains(line, want) {
				t.Errorf("jaira %v: stderr = %q, want its release line to name %q", args, stderr, want)
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

// TestClaimAndMoveAskAboutANewerReleaseOnStderr is the session that is already
// running when a release becomes known: it does not start again, but it claims
// tickets and moves them between lanes, so those commands carry the line too.
func TestClaimAndMoveAskAboutANewerReleaseOnStderr(t *testing.T) {
	setCurrent(t, "1.0.0")
	t.Setenv("JAIRA_USER", "berk")
	dir, s := updateTestStore(t)
	h := backlogTicket(t, s)
	cacheLatest(t, "1.3.0")

	for _, args := range [][]string{{"claim", h, "--session", "me"}, {"move", h, "--to", "backlog"}, {"move", h, "--to", "backlog", "--json"}} {
		stdout, stderr := runSplit(t, dir, args...)
		line := releaseLineOf(stderr)
		for _, want := range releaseLineWants {
			if !strings.Contains(line, want) {
				t.Errorf("jaira %v: stderr = %q, want its release line to name %q", args, stderr, want)
			}
		}
		if bytes.Contains([]byte(stdout), []byte("1.3.0")) {
			t.Errorf("jaira %v: stdout = %q, the release line belongs on stderr", args, stdout)
		}
	}

	stdout, _ := runSplit(t, dir, "move", h, "--to", "backlog", "--json")
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("move --json stdout is not JSON: %v\n%s", err, stdout)
	}
}

// TestReleaseLineNeverStandsBeforeAnError: a command that fails says only
// why. Under --json its stderr is the structured error a caller parses — the
// cockpit shows a failed move's stderr as its message — so the release line
// waits for a command that worked.
func TestReleaseLineNeverStandsBeforeAnError(t *testing.T) {
	setCurrent(t, "1.0.0")
	t.Setenv("JAIRA_USER", "berk")
	dir, s := updateTestStore(t)
	h := backlogTicket(t, s)
	cacheLatest(t, "1.3.0")

	root := newRoot("test")
	var out, errOut bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs([]string{"-C", dir, "move", h, "--to", "no-such-lane", "--json"})
	if err := root.Execute(); err == nil {
		t.Fatalf("move to an unknown lane succeeded:\n%s", out.String())
	}
	if bytes.Contains(errOut.Bytes(), []byte("1.3.0")) {
		t.Errorf("stderr of a failed move = %q, want no release line", errOut.String())
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
			t.Setenv("JAIRA_USER", "berk")
			dir, s := updateTestStore(t)
			h := backlogTicket(t, s)
			cacheLatest(t, c.latest)
			for _, args := range [][]string{{"next"}, {"resume"}, {"claim", h, "--session", "me"}, {"move", h, "--to", "backlog"}} {
				if _, stderr := runSplit(t, dir, args...); bytes.Contains([]byte(stderr), []byte("jaira self upgrade")) {
					t.Errorf("jaira %v: stderr = %q, want no release line", args, stderr)
				}
			}
		})
	}
}
