package gate

import (
	"slices"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// huntCodes keeps only the refusals this rule writes, so a ticket that owes
// its review lane other fields still reads clean here.
func huntCodes(vs Violations) Violations {
	var out Violations
	for _, v := range vs {
		if v.Code == CodeNeedsBughunt || v.Code == CodeHuntUncommitted {
			out = append(out, v)
		}
	}
	return out
}

// A ticket that changed code reaches signoff only with a bug hunt that
// covers HEAD (G8HNHT): 1MYFKE got there on 21.11. with a commit after its
// last hunt and bughunt.state fail, and nothing stopped it.
func TestSignoffNeedsHuntOnHEAD(t *testing.T) {
	cases := []struct {
		name string
		rep  HuntReport
		code string // "" means the move passes
		says []string
	}{
		{"stale hunt", HuntReport{CodeChanged: true, State: "fail", At: "10e5c341c", Head: "c3728ba70"}, CodeNeedsBughunt,
			[]string{"10e5c341c", "c3728ba70", "/jaira-role-bughunt G8HNHT", "--force"}},
		{"no hunt", HuntReport{CodeChanged: true, Head: "c3728ba70"}, CodeNeedsBughunt,
			[]string{"no bug hunt finished", "c3728ba70", "/jaira-role-bughunt G8HNHT"}},
		{"uncommitted", HuntReport{CodeChanged: true, State: "unknown", At: "10e5c341c", Head: "c3728ba70"}, CodeHuntUncommitted,
			[]string{"commit them first", "10e5c341c", "c3728ba70"}},
		{"covers HEAD", HuntReport{CodeChanged: true, State: "ok", At: "c3728ba70", Head: "c3728ba70"}, "", nil},
		{"stop-after-fix", HuntReport{CodeChanged: true, State: "skipped", Head: "c3728ba70"}, "", nil},
		{"no code changed", HuntReport{}, "", nil},
		{"no code changed, stale hunt", HuntReport{State: "fail", At: "10e5c341c"}, "", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			env := testEnv(t)
			env.Hunt = func(*ticket.Ticket) HuntReport { return c.rep }
			tk := ticketWith("")
			tk.ID = "01M7ZYEXTPHSC0PNW5STG8HNHT"
			vs := huntCodes(CheckAdvance(env, tk, Request{To: "signoff"}))
			if c.code == "" {
				if len(vs) != 0 {
					t.Fatalf("refused: %v", vs)
				}
				return
			}
			if len(vs) != 1 || vs[0].Code != c.code {
				t.Fatalf("got %v, want one %s", vs, c.code)
			}
			for _, s := range c.says {
				if !strings.Contains(vs[0].Message, s) {
					t.Errorf("message %q does not say %q", vs[0].Message, s)
				}
			}
		})
	}
}

// The human lane takes questions — testing's own max-rounds escape among
// them — and a person moving in the board decides for themselves: neither
// is refused for a stale hunt. Nor is a board that never runs hunts: no lane
// there owes hunt-verdict, and it behaves exactly as before.
func TestHuntRuleScope(t *testing.T) {
	stale := func(*ticket.Ticket) HuntReport { return HuntReport{CodeChanged: true, State: "fail", At: "a1b2c3d"} }
	called := false

	env := testEnv(t)
	env.Hunt = stale
	if vs := huntCodes(CheckAdvance(env, ticketWith(""), Request{To: "human", Question: "q"})); len(vs) != 0 {
		t.Errorf("move into human refused: %v", vs)
	}
	if vs := huntCodes(CheckAdvance(env, ticketWith(""), Request{To: "signoff", Interactive: true})); len(vs) != 0 {
		t.Errorf("move in the board refused: %v", vs)
	}

	env = testEnv(t)
	for _, l := range env.Lanes.Lanes {
		l.OutputProduces = slices.DeleteFunc(l.OutputProduces, func(f string) bool { return f == "hunt-verdict" })
	}
	env.Hunt = func(tk *ticket.Ticket) HuntReport { called = true; return stale(tk) }
	if vs := huntCodes(CheckAdvance(env, ticketWith(""), Request{To: "signoff"})); len(vs) != 0 {
		t.Errorf("board without hunts refused: %v", vs)
	}
	if called {
		t.Error("a board without hunts asked git about one")
	}

	env.Hunt = nil
	if vs := huntCodes(CheckAdvance(testEnv(t), ticketWith(""), Request{To: "signoff"})); len(vs) != 0 {
		t.Errorf("no Hunt on offer refused: %v", vs)
	}
}
