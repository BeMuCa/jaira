package validate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A pull request in blocked-by is no ticket and never dangles; a malformed
// one would wait forever and is an error of its own.
func TestPRRefInBlockedBy(t *testing.T) {
	for _, c := range []struct {
		dep string
		bad bool
	}{
		{"pr#12", false},
		{"PR#12", false},
		{"pr#", true},
		{"pr#0", true},
		{"pr#12a", true},
		{"pr# 12", true},
	} {
		a := ok1(idA)
		a.BlockedBy = []string{c.dep}
		ps := Tickets([]*ticket.Ticket{a}, lanes(t), nil)
		if has(ps, CodeDanglingDep) {
			t.Errorf("%q: a pull request is not a dangling ticket, got %v", c.dep, ps)
		}
		if has(ps, CodeBadPRRef) != c.bad {
			t.Errorf("%q: bad_pr_ref reported = %v, want %v (%v)", c.dep, !c.bad, c.bad, ps)
		}
	}
}
