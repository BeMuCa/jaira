package role

import (
	"strings"
	"testing"
)

// The roles pass jaira's suggestion lines on and never open a pull request
// themselves: the line is the person's to act on, and a role that ran it would
// give the decision away.
func TestRolesPassThePRSuggestionOnAndNeverRunIt(t *testing.T) {
	for _, id := range []string{"jaira-dispatcher", "jaira-worker", "jaira-teamlead"} {
		body, err := File(id, "SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		text := string(body)
		for _, want := range []string{"suggestions", "word for word", "You never run `/jaira-role-pr` yourself"} {
			if !strings.Contains(text, want) {
				t.Errorf("%s must pass jaira's suggestion lines on: missing %q", id, want)
			}
		}
	}
}

// A milestone is one pull request over its tickets, so its branch is built from
// every member's branch, and a ticket with no commit on it stops the push.
func TestPRTakesSeveralTicketsOrAMilestone(t *testing.T) {
	body, err := File("jaira-role-pr", "SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, want := range []string{
		"`milestone <name>`",
		"jaira list --milestone <name>\n--json",
		"## One branch for several tickets",
		"git merge --no-ff",
		"A ticket with no commit on the branch is a finding",
		"gh pr edit",
		"glab mr update",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("jaira-role-pr must take several tickets or a milestone: missing %q", want)
		}
	}
	one := strings.Index(text, "## One branch for several tickets")
	before := strings.Index(text, "## Before you push anything")
	if one < 0 || before < 0 || one > before {
		t.Errorf("the one-branch section must come before ## Before you push anything (at %d, %d)", one, before)
	}
}
