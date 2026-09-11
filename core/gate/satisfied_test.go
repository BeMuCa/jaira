package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// blocking builds a ticket that waits on one other id.
func blocking(dep string) *ticket.Ticket {
	tk := ticketWith("")
	tk.Status = "todo"
	tk.DoD = "it works"
	tk.BlockedBy = []string{dep}
	return tk
}

const absentDep = "01KZTT3XZ2YQBX93TTSR7BVRC9"

// The blocker finished, so it left the board for the logbook — and that used
// to be the moment the ticket waiting on it became permanently blocked and
// dropped out of 'jaira list --actionable'. Finishing work must not look like
// breaking it.
func TestFiledAwayBlockerDoesNotBlock(t *testing.T) {
	env := testEnv(t)
	env.Satisfied = func(id string) bool { return id == absentDep }

	if vs := blockedBy(env, blocking(absentDep)); len(vs) != 0 {
		t.Errorf("a finished, filed-away blocker must clear the dependency, got %v", vs)
	}
	if !Actionable(env, blocking(absentDep)) {
		t.Error("with its blocker done, the ticket is actionable")
	}
}

// The other half of the same rule: an id nothing can vouch for still blocks.
// Treating every unknown id as satisfied would turn a typo into a silently
// ignored dependency.
func TestUnknownBlockerStillBlocks(t *testing.T) {
	env := testEnv(t)
	env.Satisfied = func(string) bool { return false }

	vs := blockedBy(env, blocking(absentDep))
	if len(vs) != 1 || vs[0].Code != CodeBlocked {
		t.Errorf("an unknown blocker must still block, got %v", vs)
	}
}

// A caller with no store injects nothing, and the old behaviour has to stand
// rather than silently becoming permissive.
func TestNoSatisfiedFuncKeepsTheOldRefusal(t *testing.T) {
	env := testEnv(t)

	vs := blockedBy(env, blocking(absentDep))
	if len(vs) != 1 || vs[0].Code != CodeBlocked {
		t.Errorf("without a Satisfied func the refusal stands, got %v", vs)
	}
}

// A blocker that is on the board and unfinished is unaffected by any of this.
func TestOnBoardUnfinishedBlockerStillBlocks(t *testing.T) {
	env := testEnv(t)
	dep := ticketWith("")
	dep.ID = absentDep
	dep.Status = "in-progress"
	env.All = []*ticket.Ticket{dep}
	env.Satisfied = func(string) bool { return true }

	vs := blockedBy(env, blocking(absentDep))
	if len(vs) != 1 || vs[0].Code != CodeBlocked {
		t.Errorf("a blocker sitting in a working lane blocks whatever the index says, got %v", vs)
	}
}
