package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

func TestStaleClaimReportsTheAbandonedHolder(t *testing.T) {
	old := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRC1", ClaimedBy: "someone-else",
		ClaimedAt: time.Now().Add(-2 * ClaimTTL),
	}
	holder, since, stale := StaleClaim(old, "me")
	if !stale || holder != "someone-else" || since < ClaimTTL {
		t.Errorf("StaleClaim = %q, %v, %v; want the abandoned holder", holder, since, stale)
	}

	// A live claim is not stale — that one is refused, not stepped over — and
	// neither is your own.
	live := &ticket.Ticket{ID: old.ID, ClaimedBy: "someone-else", ClaimedAt: time.Now()}
	if _, _, stale := StaleClaim(live, "me"); stale {
		t.Error("a live claim was reported as abandoned")
	}
	mine := &ticket.Ticket{ID: old.ID, ClaimedBy: "me", ClaimedAt: time.Now().Add(-2 * ClaimTTL)}
	if _, _, stale := StaleClaim(mine, "me"); stale {
		t.Error("your own expired claim was reported as somebody's abandonment")
	}

	// ClaimActive treats a claim with no timestamp as abandoned, so this must
	// report it rather than claiming it was renewed "moments" ago.
	none := &ticket.Ticket{ID: old.ID, ClaimedBy: "someone-else"}
	if _, _, stale := StaleClaim(none, "me"); !stale {
		t.Error("a claim with no timestamp is stepped over silently")
	}
	if line := staleClaimLine(none, "me"); !strings.Contains(line, "no timestamp") {
		t.Errorf("line = %q, want it to say the claim has no timestamp", line)
	}
}

// Taking over an expired claim is allowed and must not be silent: it is the one
// moment where two sessions could be working the same ticket.
func TestClaimReportsTakingOverAnAbandonedClaim(t *testing.T) {
	dir, id := dodTestStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Mutate(id, func(tk *ticket.Ticket) error {
		if err := tk.Doc().SetScalar(ticket.FieldClaimedBy, "someone-else"); err != nil {
			return err
		}
		return tk.Doc().SetRaw(ticket.FieldClaimedAt, ticket.FormatTime(time.Now().Add(-3*time.Hour)))
	}); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "claim", id)
	if err != nil {
		t.Fatalf("claim: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Claimed") {
		t.Errorf("the claim did not go through:\n%s", out)
	}
	if !strings.Contains(out, "abandoned claim by someone-else") {
		t.Errorf("taking over an abandoned claim was silent:\n%s", out)
	}
	if !strings.Contains(out, "3h ago") {
		t.Errorf("the warning does not say how old the claim was:\n%s", out)
	}
}

// 'next' names the command that clears an abandoned claim, and running exactly
// that command clears it: a warning repeated on every 'next' with nothing to do
// about it is one people learn to ignore, and a killed session never gets a
// turn to give its claim back.
func TestNextNamesTheCommandThatClearsAnAbandonedClaim(t *testing.T) {
	dir, s, tk := claimFixture(t)
	h := ticket.Handle(tk.ID)
	if _, err := s.Mutate(tk.ID, func(x *ticket.Ticket) error {
		// Specified, so 'next' hands it out, and claimed by a session that died
		// eleven hours ago — the requirementsgenie board's case.
		for f, v := range map[string]string{
			ticket.FieldGoal: "g", ticket.FieldDoD: "d", ticket.FieldContext: "c",
			ticket.FieldAssignee: "berk", ticket.FieldClaimedBy: "dead-session",
		} {
			if err := x.Doc().SetScalar(f, v); err != nil {
				return err
			}
		}
		return x.Doc().SetRaw(ticket.FieldClaimedAt, ticket.FormatTime(time.Now().Add(-11*time.Hour)))
	}); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "next", "--all")
	if err != nil {
		t.Fatalf("next: %v\n%s", err, out)
	}
	hint := "jaira claim " + h + " --release"
	if !strings.Contains(out, "abandoned claim by dead-session") || !strings.Contains(out, hint) {
		t.Fatalf("the warning does not name %q:\n%s", hint, out)
	}

	if out, err := runCLI(t, dir, "claim", h, "--release"); err != nil {
		t.Fatalf("the named command failed: %v\n%s", err, out)
	}
	if got, _ := s.Load(tk.ID); got.ClaimedBy != "" {
		t.Errorf("the abandoned claim is still on the ticket: %q", got.ClaimedBy)
	}
	if out, _ := runCLI(t, dir, "next", "--all"); strings.Contains(out, "abandoned") {
		t.Errorf("the warning survives the release:\n%s", out)
	}
}

// 'claim' prints the same abandoned-claim line after it has taken the claim
// over; telling the new holder to release it there would be wrong.
func TestClaimTakeoverDoesNotSuggestReleasing(t *testing.T) {
	dir, s, tk := claimFixture(t)
	if _, err := s.Mutate(tk.ID, func(x *ticket.Ticket) error {
		if err := x.Doc().SetScalar(ticket.FieldClaimedBy, "dead-session"); err != nil {
			return err
		}
		return x.Doc().SetRaw(ticket.FieldClaimedAt, ticket.FormatTime(time.Now().Add(-3*time.Hour)))
	}); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "claim", ticket.Handle(tk.ID))
	if err != nil {
		t.Fatalf("claim: %v\n%s", err, out)
	}
	if strings.Contains(out, "--release") {
		t.Errorf("taking over a claim tells the new holder to release it:\n%s", out)
	}
}
