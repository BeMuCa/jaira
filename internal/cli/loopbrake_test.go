package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The loop brake reaches every lane through the payload: the blocker rule —
// the ticket's own, else the default — and stop-after-fix, beside mode.
func TestForLaneCarriesTheLoopBrake(t *testing.T) {
	dir, h := newModeTicket(t)

	for _, lane := range []string{"in-progress", "review"} {
		got := jsonCLI(t, dir, "show", h, "--for-lane", lane)
		if got["blocker_rule"] != ticket.DefaultBlockerRule {
			t.Errorf("%s: a ticket without a rule hands on %v, want the default", lane, got["blocker_rule"])
		}
		if got["stop_after_fix"] != false {
			t.Errorf("%s: stop_after_fix is %v on a ticket that never set it, want false", lane, got["stop_after_fix"])
		}
	}

	if out, err := runCLI(t, dir, "set", h, "blocker-rule=only data loss", "stop-after-fix=true"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	got := jsonCLI(t, dir, "show", h, "--for-lane", "review")
	if got["blocker_rule"] != "only data loss" || got["stop_after_fix"] != true {
		t.Errorf("the payload carries %v / %v, want the ticket's own rule and true", got["blocker_rule"], got["stop_after_fix"])
	}
	show := jsonCLI(t, dir, "show", h)
	if show["blocker_rule"] != "only data loss" || show["stop_after_fix"] != true {
		t.Errorf("show --json carries %v / %v", show["blocker_rule"], show["stop_after_fix"])
	}
	plain, err := runCLI(t, dir, "show", h, "--for-lane", "review")
	if err != nil || !strings.Contains(plain, "blockers: only data loss") || !strings.Contains(plain, "stop-after-fix") {
		t.Errorf("the plain-text lane brief does not name the brake: %v\n%s", err, plain)
	}
}

// A blocker rule is a sentence a person writes: quotes, umlauts and more
// than one line come back as written, a CRLF as the line break it means.
func TestBlockerRuleRoundTrip(t *testing.T) {
	dir, h := newModeTicket(t)
	for _, rule := range []string{
		`a "wrong" result: data loss, crash`,
		"Datenverlust, Absturz, falsches Ergebnis (äöü ß)",
		"wrong result\nor data loss",
		"wrong result\r\nor data loss",
		"wrong result\ror data loss",
		"wrong result: 'quoted' # not a comment",
	} {
		if out, err := runCLI(t, dir, "set", h, "blocker-rule="+rule); err != nil {
			t.Fatalf("set blocker-rule=%q: %v\n%s", rule, err, out)
		}
		want := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(rule)
		if got := jsonCLI(t, dir, "show", h)["blocker_rule"]; got != want {
			t.Errorf("blocker-rule %q came back as %q", rule, got)
		}
	}
}

// stop-after-fix is read as exactly "true": anything else a person might type
// is refused rather than stored and silently read as off.
func TestSetRefusesUnknownStopAfterFix(t *testing.T) {
	dir, h := newModeTicket(t)
	for _, ok := range []string{"true", "FALSE", " true ", ""} {
		if out, err := runCLI(t, dir, "set", h, "stop-after-fix="+ok); err != nil {
			t.Errorf("set stop-after-fix=%q refused: %v\n%s", ok, err, out)
		}
	}
	for _, bad := range []string{"yes", "on", "1", "tru"} {
		out, err := runCLI(t, dir, "set", h, "stop-after-fix="+bad)
		if err == nil {
			t.Fatalf("set stop-after-fix=%s was accepted:\n%s", bad, out)
		}
		if !strings.Contains(out+err.Error(), "stop-after-fix") {
			t.Errorf("refusal does not name the field: %v\n%s", err, out)
		}
	}
	if got := jsonCLI(t, dir, "show", h)["stop_after_fix"]; got != false {
		t.Errorf("after the last good write (empty) and refused ones, stop_after_fix is %v", got)
	}
}

// A send-back typed at a terminal is a person's; the same move through a
// pipe — how every agent runs jaira — is the loop lane's and counts.
func TestCLIMoveAtATerminalIsAPersons(t *testing.T) {
	for _, tc := range []struct {
		name     string
		terminal bool
		by       string
		count    float64
	}{
		{"at a terminal", true, ticket.BackByPerson, 1},
		{"through a pipe", false, ticket.BackByLane, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			was := atATerminal
			atATerminal = func() bool { return tc.terminal }
			t.Cleanup(func() { atATerminal = was })

			dir, h := newModeTicket(t)
			for _, to := range []string{"review", "in-progress"} {
				if out, err := runCLI(t, dir, "move", h, "--to", to, "--force"); err != nil {
					t.Fatalf("move --to %s: %v\n%s", to, err, out)
				}
			}
			back, _ := jsonCLI(t, dir, "show", h)["back"].(map[string]any)
			if back["from"] != "review" || back["by"] != tc.by || back["count"] != tc.count {
				t.Errorf("back is %v, want from review, by %s, count %v", back, tc.by, tc.count)
			}
		})
	}
}

// The review's gaps stay beside its verdict, whatever order they were
// written in, while a later round's loop lane rewrites review-gaps.
func TestVerdictKeepsItsGaps(t *testing.T) {
	for _, tc := range []struct {
		name   string
		writes [][]string
	}{
		{"one call", [][]string{{"review-gaps=note: a\nblocker: b", "review-verdict=back"}}},
		{"gaps first", [][]string{{"review-gaps=note: a\nblocker: b"}, {"review-verdict=back"}}},
		{"verdict first", [][]string{{"review-verdict=back"}, {"review-gaps=note: a\nblocker: b"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, h := newModeTicket(t)
			for _, w := range tc.writes {
				if out, err := runCLI(t, dir, append([]string{"set", h}, w...)...); err != nil {
					t.Fatalf("set %v: %v\n%s", w, err, out)
				}
			}
			for _, to := range []string{"review", "in-progress"} {
				if out, err := runCLI(t, dir, "move", h, "--to", to, "--force"); err != nil {
					t.Fatalf("move --to %s: %v\n%s", to, err, out)
				}
			}
			if out, err := runCLI(t, dir, "set", h, "review-gaps=removed helper fooBar"); err != nil {
				t.Fatalf("set: %v\n%s", err, out)
			}
			review, _ := jsonCLI(t, dir, "show", h)["review"].(map[string]any)
			if review["verdict_gaps"] != "note: a\nblocker: b" || review["gaps"] != "removed helper fooBar" {
				t.Errorf("review is %v, want the review's gaps kept beside the verdict and the later round's in gaps", review)
			}
		})
	}
}

// The brief says whether the newest finished hunt still covers the code, so
// an Explain chat handed four notes does not have to find the hunt's note.
func TestBriefCarriesTheHunt(t *testing.T) {
	dir, h := newModeTicket(t)
	git := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t"}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	edit := func(content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	hunt := func() any {
		t.Helper()
		out, err := runCLI(t, dir, "show", h, "--brief")
		if err != nil {
			t.Fatalf("show --brief: %v\n%s", err, out)
		}
		var b map[string]any
		if err := json.Unmarshal([]byte(out[strings.Index(out, "{"):]), &b); err != nil {
			t.Fatalf("brief is no JSON: %v\n%s", err, out)
		}
		return b["bughunt"]
	}
	git("init", "-q")
	edit("a")
	git("add", "a.go")
	git("commit", "-q", "-m", "a")
	if got := hunt(); got != nil {
		t.Errorf("a ticket without a hunt reports %v, want null", got)
	}
	at := git("rev-parse", "--short", "HEAD")
	if out, err := runCLI(t, dir, "note", h, "bughunt: done at "+at+" - 0 findings; files: a.go"); err != nil {
		t.Fatalf("note: %v\n%s", err, out)
	}
	for i := 0; i < 5; i++ { // the hunt's note scrolls out of notes_last
		if out, err := runCLI(t, dir, "note", h, "later note"); err != nil {
			t.Fatalf("note: %v\n%s", err, out)
		}
	}
	state := func(want, when string) {
		t.Helper()
		if got, _ := hunt().(map[string]any); got["at"] != at || got["state"] != want {
			t.Errorf("%s: %v, want at %s, state %s", when, got, at, want)
		}
	}
	state("ok", "right after the hunt")
	edit("a2")
	state("unknown", "a hunted file changed, not committed")
	git("commit", "-q", "-am", "fix("+h+"): a2")
	state("fail", "a ticket commit on a hunted file")
	// stop-after-fix: a fix round runs no hunt by choice - skipped, not a fault
	if out, err := runCLI(t, dir, "set", h, "stop-after-fix=true"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	state("skipped", "stop-after-fix after a ticket commit")
	if out, err := runCLI(t, dir, "set", h, "stop-after-fix="); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	// a hunt at a commit this clone does not have (rebased away) is fail
	if out, err := runCLI(t, dir, "note", h, "bughunt: done at deadbeef1 - 0 findings; files: a.go"); err != nil {
		t.Fatalf("note: %v\n%s", err, out)
	}
	if got, _ := hunt().(map[string]any); got["at"] != "deadbeef1" || got["state"] != "fail" {
		t.Errorf("a hunt at a missing commit: %v, want state fail", got)
	}
}

// A ticket that runs no hunt after a fix and has none yet reads skipped.
func TestBriefSkipsTheHuntOnStopAfterFix(t *testing.T) {
	dir, h := newModeTicket(t)
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "set", h, "stop-after-fix=true"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "show", h, "--brief")
	if err != nil {
		t.Fatalf("show --brief: %v\n%s", err, out)
	}
	var b map[string]any
	if err := json.Unmarshal([]byte(out[strings.Index(out, "{"):]), &b); err != nil {
		t.Fatalf("brief is no JSON: %v\n%s", err, out)
	}
	if got, _ := b["bughunt"].(map[string]any); got["state"] != "skipped" || got["at"] != "" {
		t.Errorf("bughunt is %v, want skipped", b["bughunt"])
	}
}
