package gate

import "testing"

func hasCode(vs Violations, code string) bool {
	for _, v := range vs {
		if v.Code == code {
			return true
		}
	}
	return false
}

// Sending work back from the sign-off lane decides nothing: it is the opposite
// of accepting it. Refusing it made every send-back a --force, and --force
// silences every other refusal on the move along with this one.
func TestHumanExitLaneSendsBackWithoutForce(t *testing.T) {
	tk := ticketWith("")
	tk.Status = "signoff"

	vs := CheckAdvance(testEnv(t), tk, Request{To: "in-progress"})
	if hasCode(vs, CodeNeedsHuman) {
		t.Errorf("an agent was refused sending work back from signoff: %v", vs)
	}
}

// Forward out of the lane is the acceptance itself, and stays a person's.
func TestHumanExitLaneStillRefusesForward(t *testing.T) {
	for _, to := range []string{"done", "blocked"} {
		tk := ticketWith("")
		tk.Status = "signoff"

		vs := CheckAdvance(testEnv(t), tk, Request{To: to, Reason: "r"})
		if !hasCode(vs, CodeNeedsHuman) {
			t.Errorf("an agent was allowed from signoff to %s: %v", to, vs)
		}
	}
}

// The rule follows the lane contract, not the name signoff: a board may put a
// person's checkpoint anywhere in its order.
func TestHumanExitRuleFollowsTheContract(t *testing.T) {
	env := testEnv(t)
	l, ok := env.Lanes.Get("testing")
	if !ok {
		t.Fatal("the shipped board has no testing lane")
	}
	l.RequiresHumanExit = true

	tk := ticketWith("")
	tk.Status = "testing"

	if vs := CheckAdvance(env, tk, Request{To: "in-progress"}); hasCode(vs, CodeNeedsHuman) {
		t.Errorf("an agent was refused sending work back from a human-exit lane: %v", vs)
	}
	if vs := CheckAdvance(env, tk, Request{To: "review"}); !hasCode(vs, CodeNeedsHuman) {
		t.Errorf("an agent was allowed forward out of a human-exit lane: %v", vs)
	}
}
