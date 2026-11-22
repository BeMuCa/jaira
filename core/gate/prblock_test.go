package gate

import (
	"strings"
	"testing"
)

// TestPRBlocker: a pull request in blocked-by clears only when the injected
// PR answer says merged. Every other answer — open, closed, unknown, no
// answer at all, a malformed entry — keeps the ticket waiting and says on
// what, under its own code.
func TestPRBlocker(t *testing.T) {
	cases := []struct {
		name string
		dep  string
		pr   func(int) PRWait
		want string // "" = clears
	}{
		{"merged", "pr#12", func(int) PRWait { return PRWait{Merged: true} }, ""},
		{"open", "pr#12", func(int) PRWait { return PRWait{Why: "open"} }, "waits for PR #12: open"},
		{"closed", "pr#12", func(int) PRWait { return PRWait{Why: "closed without merge"} }, "waits for PR #12: closed without merge"},
		{"unknown", "pr#12", func(int) PRWait { return PRWait{Why: "state unknown (gh missing)"} }, "state unknown"},
		{"no hook", "pr#12", nil, "waits for PR #12: not asked here"},
		{"malformed", "pr#12a", func(int) PRWait { return PRWait{Merged: true} }, `waits for pr#12a: "pr#12a" is not pr#<number>`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			env := testEnv(t)
			asked := 0
			if c.pr != nil {
				env.PR = func(n int) PRWait {
					asked++
					if n != 12 {
						t.Fatalf("asked PR %d, want 12", n)
					}
					return c.pr(n)
				}
			}
			tk := blocking(c.dep)
			vs := blockedBy(env, tk)
			if c.want == "" {
				if len(vs) != 0 || !Actionable(env, tk) {
					t.Fatalf("a merged pull request must clear the ticket, got %v", vs)
				}
				return
			}
			if len(vs) != 1 || vs[0].Code != CodeWaitsOnPR || !strings.Contains(vs[0].Message, c.want) {
				t.Fatalf("got %v, want one %s saying %q", vs, CodeWaitsOnPR, c.want)
			}
			if Actionable(env, tk) {
				t.Fatal("a ticket waiting on a pull request is not actionable")
			}
			if w := WaitingOnPR(env, tk); len(w) != 1 || w[0] != vs[0].Message {
				t.Fatalf("WaitingOnPR = %v", w)
			}
			if c.name == "malformed" && asked != 0 {
				t.Fatal("a malformed entry must not be asked about")
			}
		})
	}
}

// A ticket waiting on a ticket and on a pull request clears only when both
// do, and WaitingOnPR names only the pull request.
func TestPRAndTicketBlocker(t *testing.T) {
	env := testEnv(t)
	env.Satisfied = func(string) bool { return false }
	env.PR = func(int) PRWait { return PRWait{Merged: true} }
	tk := blocking(absentDep)
	tk.BlockedBy = append(tk.BlockedBy, "pr#3")
	if Actionable(env, tk) {
		t.Fatal("the ticket blocker still holds")
	}
	if w := WaitingOnPR(env, tk); len(w) != 0 {
		t.Fatalf("merged pr listed as waiting: %v", w)
	}
	env.PR = func(int) PRWait { return PRWait{Why: "open"} }
	if w := Waiting(env, tk); len(w) != 2 {
		t.Fatalf("Waiting = %v, want the ticket and the pull request", w)
	}
}
