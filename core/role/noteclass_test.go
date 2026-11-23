package role

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
)

// GY8T52: the review lane writes each note's class as a fixed word, and the
// signoff role and the teamlead name the same four - the cockpit colours a
// note by that word, so a table in other words shows the person another
// colour than the pane.
func TestReviewSignoffAndTeamleadShareTheFourNoteClasses(t *testing.T) {
	lanes, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	prompts := map[string]string{}
	for _, l := range lanes {
		prompts[l.ID] = l.Prompt
	}
	review, ok := prompts["review"]
	if !ok {
		t.Fatal("no built-in review lane")
	}
	for _, prefix := range []string{"`blocker:`", "`note: edge:`", "`note: limit:`", "`note: elsewhere:`", "`note: edge: fix-now:`"} {
		if !strings.Contains(review, prefix) {
			t.Errorf("the review lane does not ask for %s", prefix)
		}
	}
	// Its examples write the classes too: an example without one teaches a
	// note the pane shows as an edge.
	for _, line := range strings.Split(review, "\n") {
		if i := strings.Index(line, "review-gaps="); i >= 0 && strings.Contains(line[i:], "note: ") {
			for _, n := range strings.Split(line[i:], "note: ")[1:] {
				if !strings.HasPrefix(n, "edge:") && !strings.HasPrefix(n, "limit:") && !strings.HasPrefix(n, "elsewhere:") {
					t.Errorf("an example note in the review lane carries no class: %q", line)
				}
			}
		}
	}
	if !strings.Contains(prompts["testing"], "testing: note: edge:") {
		t.Error("the testing lane's note for the person carries no class")
	}
	for _, id := range []string{"jaira-role-signoff", "jaira-teamlead"} {
		b, err := File(id, "SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		for _, word := range []string{"🔴 **blocker**", "🟡 **edge**", "⚪ **limit**", "🔵 **elsewhere**", "`note: edge:`"} {
			plain := strings.ReplaceAll(word, "**", "")
			if !strings.Contains(string(b), word) && !strings.Contains(string(b), plain) {
				t.Errorf("%s does not name %s", id, plain)
			}
		}
	}
}
