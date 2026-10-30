package suggest

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// testLanes is a board in miniature: work, an acceptance lane a person owns,
// and the terminal lane. The names are deliberately not the shipped ones, so
// a rule that reads a name instead of a property fails here.
func testLanes() *lane.Set {
	return &lane.Set{Lanes: []*lane.Lane{
		{ID: "work", Agentic: true},
		{ID: "accept", RequiresHumanExit: true},
		{ID: "finish", Terminal: true},
	}}
}

const (
	idA = "01M4ZQ715HGZN2GHX5NWAAAAAA"
	idB = "01M4ZQ715HGZN2GHX5NWBBBBBB"
	idC = "01M4ZQ715HGZN2GHX5NWCCCCCC"
)

func kinds(ls []Line) []string {
	var out []string
	for _, l := range ls {
		out = append(out, l.Kind)
	}
	return out
}

func TestTicketReachingAcceptanceSuggestsItsPullRequest(t *testing.T) {
	got := ForTicket(testLanes(), idA, "accept", nil)
	if len(got) != 1 || got[0].Kind != KindPullRequest || got[0].By != "person" {
		t.Fatalf("got %+v, want one pull_request by person", got)
	}
	if !strings.Contains(got[0].Text, "/jaira-role-pr "+ticket.Handle(idA)) {
		t.Errorf("line %q does not name the ticket", got[0].Text)
	}
}

func TestTicketReachingTheTerminalLaneSuggestsItsPullRequest(t *testing.T) {
	got := ForTicket(testLanes(), idA, "finish", nil)
	if len(got) != 1 || got[0].Kind != KindPullRequest {
		t.Fatalf("got %+v, want one pull_request", got)
	}
}

func TestTicketFiledOffTheBoardSuggestsItsPullRequest(t *testing.T) {
	// A logbook entry or a file-on-entry move has no status left: it is done.
	got := ForTicket(testLanes(), idA, "", nil)
	if len(got) != 1 || got[0].Kind != KindPullRequest {
		t.Fatalf("got %+v, want one pull_request", got)
	}
}

func TestNoSuggestionOnAnAgenticLane(t *testing.T) {
	if got := ForTicket(testLanes(), idA, "work", nil); len(got) != 0 {
		t.Fatalf("got %+v, want nothing", got)
	}
}

func TestTicketInAnOpenMilestoneSuggestsNoPullRequest(t *testing.T) {
	groups := []Group{{Name: "m", Members: []Member{{ID: idA, Status: "work"}, {ID: idB, Status: "work"}}}}
	got := ForTicket(testLanes(), idA, "accept", groups)
	for _, l := range got {
		if l.Kind == KindPullRequest {
			t.Fatalf("a ticket in an open milestone got a pull request line: %+v", got)
		}
	}
	if len(got) != 1 || got[0].Kind != KindProgress || !strings.Contains(got[0].Text, "Milestone m: 0 of 2 done") {
		t.Fatalf("got %+v, want one progress line 0 of 2", got)
	}
}

func TestMilestoneCompleteSuggestsOnePullRequestThenFiling(t *testing.T) {
	groups := []Group{{Name: "m", Members: []Member{{ID: idA, Status: "work"}, {ID: idB, Status: "finish"}}}}
	// idA is the ticket that just finished: its status comes from the call.
	got := ForTicket(testLanes(), idA, "finish", groups)
	if want := []string{KindPullRequest, KindFileMilestone}; strings.Join(kinds(got), ",") != strings.Join(want, ",") {
		t.Fatalf("kinds = %v, want %v (%+v)", kinds(got), want, got)
	}
	if !strings.Contains(got[0].Text, ticket.Handle(idA)+" "+ticket.Handle(idB)) {
		t.Errorf("pull request line %q does not name every ticket", got[0].Text)
	}
	if got[1].After != "merge" || got[1].Command != "jaira logbook m" {
		t.Errorf("filing line = %+v, want after merge, jaira logbook m", got[1])
	}
}

func TestTicketInTwoMilestonesGetsTheLinesOfEachComplete(t *testing.T) {
	groups := []Group{
		{Name: "one", Members: []Member{{ID: idA, Status: "finish"}}},
		{Name: "two", Members: []Member{{ID: idA, Status: "finish"}, {ID: idC, Status: "work"}}},
	}
	got := ForTicket(testLanes(), idA, "finish", groups)
	if want := []string{KindPullRequest, KindFileMilestone, KindProgress}; strings.Join(kinds(got), ",") != strings.Join(want, ",") {
		t.Fatalf("kinds = %v, want %v (%+v)", kinds(got), want, got)
	}
	if !strings.Contains(got[0].Text, "Milestone one is complete") {
		t.Errorf("first line %q is not about milestone one", got[0].Text)
	}
}

func TestEmptyMilestoneIsNeverComplete(t *testing.T) {
	g := Group{Name: "empty"}
	if got := ForMilestone(testLanes(), g); len(got) != 0 {
		t.Fatalf("empty milestone suggested %+v", got)
	}
}

func TestFiledMilestoneSuggestsNothing(t *testing.T) {
	g := Group{Name: "m", Filed: true, Members: []Member{{ID: idA, Status: "finish"}}}
	if got := ForMilestone(testLanes(), g); len(got) != 0 {
		t.Fatalf("filed milestone suggested %+v", got)
	}
	if got := ForTicket(testLanes(), idA, "finish", []Group{g}); len(got) != 1 || got[0].Kind != KindPullRequest {
		t.Fatalf("a ticket whose only milestone is filed is a group of one: got %+v", got)
	}
}

func TestLanesAreFoundByWhatTheyDoNotByName(t *testing.T) {
	lanes := &lane.Set{Lanes: []*lane.Lane{
		{ID: "signoff"}, // the name says accept, the lane does not
		{ID: "review-x", RequiresHumanExit: true}, // the lane is acceptance under any name
	}}
	if got := ForTicket(lanes, idA, "signoff", nil); len(got) != 0 {
		t.Errorf("a lane named signoff without the property suggested %+v", got)
	}
	if got := ForTicket(lanes, idA, "review-x", nil); len(got) != 1 {
		t.Errorf("a lane with requires-human-exit suggested %+v, want one line", got)
	}
}

func TestTheStatusPassedInOverridesTheStoredOne(t *testing.T) {
	// The group still shows idA in work, as the file on disk did before the move.
	groups := []Group{{Name: "m", Members: []Member{{ID: idA, Status: "work"}}}}
	got := ForTicket(testLanes(), idA, "finish", groups)
	if len(got) != 2 || got[0].Kind != KindPullRequest {
		t.Fatalf("got %+v, want the milestone complete", got)
	}
}

func TestLinesSayOnlyAPersonOpensThePullRequest(t *testing.T) {
	groups := []Group{{Name: "m", Members: []Member{{ID: idA, Status: "finish"}}}}
	for _, l := range append(ForTicket(testLanes(), idA, "finish", groups), ForTicket(testLanes(), idB, "accept", nil)...) {
		if l.Kind == KindPullRequest && !strings.Contains(l.Text, "Only a person opens it") {
			t.Errorf("pull request line %q does not say a person opens it", l.Text)
		}
	}
}
