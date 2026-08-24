package cli

import (
	"testing"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// A ticket whose lane has already produced its declared output is waiting to be
// moved on, not to be worked. Handed back by 'next', it gets the lane's work
// done twice while the tickets the lane has not touched keep waiting — which is
// how nine tickets sat in a critique lane uncritiqued.
func TestNextPrefersATicketItsLaneHasNotWorked(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	env := gate.Env{Lanes: lanes}

	// The older ID sorts first on the existing tiebreak, so it is the worked one:
	// if the preference did not exist, this test would see it first.
	worked := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRC1", Title: "worked", Status: "in-progress",
		Outcome: ticket.Outcome{What: "w", Why: "y", Resolves: "r"},
	}
	unworked := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRC2", Title: "unworked", Status: "in-progress"}

	ts := []*ticket.Ticket{worked, unworked}
	sortByProgress(ts, env)
	if ts[0] != unworked {
		t.Errorf("first = %q, want the ticket its lane has not worked", ts[0].Title)
	}

	// And once its output is produced it loses the preference, so a lane cannot
	// keep handing back the same ticket.
	unworked.Outcome = ticket.Outcome{What: "w", Why: "y", Resolves: "r"}
	sortByProgress(ts, env)
	if ts[0] != worked {
		t.Errorf("with both worked, first = %q, want the older ID", ts[0].Title)
	}
}

// Lane precedence still decides across lanes: "finish in-flight work first" is
// the stronger rule, and preferring unworked tickets must not promote a fresh
// todo above work already under way.
func TestUnworkedPreferenceDoesNotOutrankLanePrecedence(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	env := gate.Env{Lanes: lanes}

	inFlight := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRC1", Title: "in flight", Status: "in-progress",
		Outcome: ticket.Outcome{What: "w", Why: "y", Resolves: "r"},
	}
	fresh := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRC2", Title: "fresh", Status: "todo"}

	ts := []*ticket.Ticket{fresh, inFlight}
	sortByProgress(ts, env)
	if ts[0] != inFlight {
		t.Errorf("first = %q, want the further-along lane", ts[0].Title)
	}
}
