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
