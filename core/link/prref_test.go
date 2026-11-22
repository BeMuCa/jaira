package link

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A pull request in blocked-by is no ticket: Relations lists the ticket
// blockers beside it and never a dangling pr#12.
func TestPRRefIsNoRelation(t *testing.T) {
	s := store(t)
	dep := make1(t, s, "dep", "todo", nil, nil)
	id := make1(t, s, "waits", "todo", nil, map[string][]string{ticket.FieldBlockedBy: {"pr#12", dep}})
	ix := Build(s, lanes(t), list(t, s))
	got := kinds(ix.Relations(id), KindBlockedBy)
	if len(got) != 1 || got[0] != dep {
		t.Fatalf("blocked-by relations = %v, want only %s", got, dep)
	}
}
