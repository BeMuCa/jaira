package cli

import (
	"strings"
	"testing"
)

// reviewedTicket is a ticket in review with every review output written and
// the given verdict, moved there the way an agent moves it.
func reviewedTicket(t *testing.T, verdict string) (dir, h string) {
	t.Helper()
	dir, h = newModeTicket(t)
	if out, err := runCLI(t, dir, "move", h, "--to", "review", "--force"); err != nil {
		t.Fatalf("move --to review: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "set", h, "review-summary=s", "review-gaps=blocker: DoD 2", "review-check=c", "review-verdict="+verdict); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	return dir, h
}

// The reviewer said send back: the ticket does not go on to sign-off, the
// refusal is exit 3 and names the verdict and where the ticket goes instead.
func TestCLISendBackVerdictRefusesSignoff(t *testing.T) {
	dir, h := reviewedTicket(t, "send back: DoD 2 is unmet")
	out, err := runCLI(t, dir, "move", h, "--to", "signoff")
	if exitCode(err) != ExitValidation {
		t.Fatalf("move to signoff with a send-back verdict: exit %d, want %d\n%s", exitCode(err), ExitValidation, out)
	}
	msg := out + err.Error()
	for _, want := range []string{`"send-back"`, `"in-progress"`, "review-verdict"} {
		if !strings.Contains(msg, want) {
			t.Errorf("refusal does not name %s:\n%s", want, msg)
		}
	}
	if got := jsonCLI(t, dir, "show", h)["status"]; got != "review" {
		t.Errorf("the refused move left the ticket in %v", got)
	}
}

// ... and goes back to in-progress without --force.
func TestCLISendBackVerdictGoesBack(t *testing.T) {
	dir, h := reviewedTicket(t, "send-back: DoD 2 is unmet")
	if out, err := runCLI(t, dir, "move", h, "--to", "in-progress"); err != nil {
		t.Fatalf("move back to in-progress refused: %v\n%s", err, out)
	}
}

// An accepting verdict passes, so the rail costs a correct review nothing.
func TestCLIAcceptVerdictGoesOn(t *testing.T) {
	dir, h := reviewedTicket(t, "accept: the diff matches the criteria")
	if out, err := runCLI(t, dir, "move", h, "--to", "signoff"); err != nil {
		t.Fatalf("move to signoff with an accepting verdict refused: %v\n%s", err, out)
	}
}

// The words are in the lane's brief, so the reviewer reads what it is held to.
func TestForLaneReviewNamesTheVerdicts(t *testing.T) {
	dir, h := newModeTicket(t)
	got, _ := jsonCLI(t, dir, "show", h, "--for-lane", "review")["verdicts"].(map[string]any)
	acc, _ := got["accept"].([]any)
	back, _ := got["send_back"].([]any)
	if len(acc) != 1 || acc[0] != "accept" || len(back) != 1 || back[0] != "send-back" {
		t.Errorf("show --for-lane review carries verdicts %v, want accept / send-back", got)
	}
	if _, ok := jsonCLI(t, dir, "show", h, "--for-lane", "in-progress")["verdicts"]; ok {
		t.Error("a lane without verdict words carries a verdicts key")
	}
}
