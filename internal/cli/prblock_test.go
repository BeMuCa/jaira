package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// prStates stands in for git and gh: each pull request number answers with
// its own state, "!" fails with that text; look reports gh installed or not.
func prStates(t *testing.T, ghInstalled bool, states map[string]string) *int {
	t.Helper()
	calls := 0
	oldRun, oldLook := prRun, prLook
	t.Cleanup(func() { prRun, prLook = oldRun, oldLook })
	prLook = func(string) (string, error) {
		if !ghInstalled {
			return "", errors.New("executable file not found in $PATH")
		}
		return "/usr/bin/gh", nil
	}
	prRun = func(_ string, name string, args ...string) ([]byte, error) {
		line := name + " " + strings.Join(args, " ")
		switch {
		case strings.HasPrefix(line, "git remote get-url origin"):
			return []byte("https://github.com/me/jaira.git\n"), nil
		case strings.HasPrefix(line, "git config"):
			return nil, prExit{1, ""}
		case strings.HasPrefix(line, "gh repo view"):
			return []byte(`{"nameWithOwner":"me/jaira"}`), nil
		case strings.HasPrefix(line, "gh api graphql"):
			calls++
			for n, s := range states {
				if strings.Contains(line, "number="+n+" ") {
					if strings.HasPrefix(s, "!") {
						return nil, prExit{1, s[1:]}
					}
					return []byte(`{"data":{"repository":{"pullRequest":{"state":"` + s + `"}}}}`), nil
				}
			}
		}
		return nil, errors.New("unexpected: " + line)
	}
	return &calls
}

// waitingTicket is a specified, assigned ticket whose only blocker is ref.
func waitingTicket(t *testing.T, dir, title, ref string) string {
	t.Helper()
	h := ticket.Handle(createTicket(t, dir, title, "--blocked-by", ref))
	if out, err := runCLI(t, dir, "set", h, "assignee=ada"); err != nil {
		t.Fatalf("set assignee: %v\n%s", err, out)
	}
	return h
}

// set and create take pr#<n> in any case and write pr#<n>; every malformed
// one is refused with exit 3, the code an unknown ticket gets, and nothing is
// written.
func TestBlockedByTakesAPullRequest(t *testing.T) {
	prStates(t, true, map[string]string{"12": "OPEN"})
	dir, h, _ := prBoard(t)
	for _, in := range []string{"pr#12", "PR#12", " Pr#12 "} {
		if out, err := runCLI(t, dir, "set", h, "blocked-by="+in); err != nil {
			t.Fatalf("set blocked-by=%q: %v\n%s", in, err, out)
		}
		got := jsonOf(t, dir, "show", h)["blocked_by"].([]any)
		if len(got) != 1 || got[0] != "pr#12" {
			t.Fatalf("blocked-by=%q stored as %v, want [pr#12]", in, got)
		}
	}
	for _, bad := range []string{"pr#", "pr#0", "pr#-1", "pr#abc", "pr# 12", "pr#99999999999999999999"} {
		out, err := runCLI(t, dir, "set", h, "blocked-by="+bad)
		if exitCode(err) != ExitValidation || !strings.Contains(out+fmt.Sprint(err), "pr#<number>") {
			t.Errorf("set blocked-by=%q: exit %d, %v\n%s", bad, exitCode(err), err, out)
		}
		if got := jsonOf(t, dir, "show", h)["blocked_by"].([]any); len(got) != 1 || got[0] != "pr#12" {
			t.Errorf("a refused set wrote %v", got)
		}
		out, err = runCLI(t, dir, "create", "bad "+bad, "--goal", "g", "--context", "c", "--dod", "d", "--blocked-by", bad)
		if exitCode(err) != ExitValidation {
			t.Errorf("create --blocked-by %q: exit %d, %v\n%s", bad, exitCode(err), err, out)
		}
	}
	// A handle given to create is stored as the full id, as set does.
	other := createTicket(t, dir, "other", "--blocked-by", h+",PR#7")
	got := jsonOf(t, dir, "show", other)["blocked_by"].([]any)
	if len(got) != 2 || !strings.HasSuffix(got[0].(string), h) || len(got[0].(string)) != 26 || got[1] != "pr#7" {
		t.Fatalf("create stored blocked-by %v, want [<full id of %s> pr#7]", got, h)
	}
}

// list and next show a ticket as waiting on its pull request while it is
// open, closed or unreadable, and next hands it out once the pull request
// is merged.
func TestTicketWaitsForItsPullRequest(t *testing.T) {
	calls := prStates(t, true, map[string]string{"12": "OPEN", "13": "MERGED", "14": "CLOSED", "15": "!HTTP 401: Bad credentials"})
	dir, _, _ := prBoard(t)
	open := waitingTicket(t, dir, "waits on open", "pr#12")
	merged := waitingTicket(t, dir, "waits on merged", "pr#13")
	closed := waitingTicket(t, dir, "waits on closed", "pr#14")
	auth := waitingTicket(t, dir, "waits on signed out", "pr#15")

	out, err := runCLI(t, dir, "list")
	if err != nil {
		t.Fatalf("list: %v\n%s", err, out)
	}
	for _, want := range []string{"[waits for PR #12: open]", "[waits for PR #14: closed without merge", "[waits for PR #15: state unknown (not signed in"} {
		if !strings.Contains(out, want) {
			t.Errorf("list lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "PR #13") {
		t.Errorf("a merged pull request is still listed as waiting:\n%s", out)
	}

	var next struct {
		Tickets []struct {
			Handle string `json:"handle"`
		} `json:"tickets"`
		Count       int `json:"count"`
		WaitingOnPR []struct {
			Handle    string   `json:"handle"`
			WaitingOn []string `json:"waiting_on"`
		} `json:"waiting_on_pr"`
	}
	raw, err := runCLI(t, dir, "next", "--all", "--json")
	if err != nil {
		t.Fatalf("next: %v\n%s", err, raw)
	}
	if err := json.Unmarshal([]byte(raw[strings.Index(raw, "{"):]), &next); err != nil {
		t.Fatalf("%v\n%s", err, raw)
	}
	var handed []string
	for _, tk := range next.Tickets {
		handed = append(handed, tk.Handle)
	}
	if !strings.Contains(strings.Join(handed, " "), merged) {
		t.Errorf("next does not hand out %s, whose pull request merged: %v", merged, handed)
	}
	waiting := map[string]string{}
	for _, w := range next.WaitingOnPR {
		waiting[w.Handle] = strings.Join(w.WaitingOn, "; ")
	}
	for _, h := range []string{open, closed, auth} {
		if strings.Contains(strings.Join(handed, " "), h) {
			t.Errorf("next hands out %s, whose pull request has not merged", h)
		}
		if waiting[h] == "" {
			t.Errorf("next --json waiting_on_pr lacks %s: %v", h, waiting)
		}
	}

	text, _ := runCLI(t, dir, "next", "--all")
	if !strings.Contains(text, "Waiting on a pull request:") || !strings.Contains(text, open+"  waits on open — waits for PR #12: open") {
		t.Errorf("next text does not name the waiting ticket:\n%s", text)
	}

	// Answers are cached: a second command asks only the failed one again.
	before := *calls
	if out, err := runCLI(t, dir, "list", "--json"); err != nil {
		t.Fatalf("list --json: %v\n%s", err, out)
	} else if !strings.Contains(out, `"waits for PR #12: open"`) {
		t.Errorf("list --json lacks waiting_on:\n%s", out)
	}
	if *calls-before != 1 {
		t.Errorf("list --json asked the forge %d times, want 1 (only the failed answer)", *calls-before)
	}
}

// Without gh nothing is guessed: the ticket waits and says why.
func TestNoForgeToolKeepsTheTicketWaiting(t *testing.T) {
	prStates(t, false, nil)
	dir, _, _ := prBoard(t)
	h := waitingTicket(t, dir, "waits", "pr#12")
	out, err := runCLI(t, dir, "list", "--actionable")
	if err != nil {
		t.Fatalf("list: %v\n%s", err, out)
	}
	if strings.Contains(out, h) {
		t.Errorf("without gh the ticket must not be actionable:\n%s", out)
	}
	out, _ = runCLI(t, dir, "list")
	if !strings.Contains(out, "[waits for PR #12: state unknown (forge tool missing: gh is not installed)]") {
		t.Errorf("list does not say why:\n%s", out)
	}
}
