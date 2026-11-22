package pr

import (
	"errors"
	"strings"
	"testing"
)

// TestState covers the answers a blocked-by entry pr#<n> can get: the three
// states on both forges, and every way of not getting one — no forge, no
// tool, no sign-in, no such pull request, an unreadable or unknown answer.
func TestState(t *testing.T) {
	gl := func(mr string) []rule {
		return []rule{
			{key: ghOrigin, out: "git@gitlab.com:me/jaira.git\n"},
			{key: "git config jaira.forge", err: exitErr{1, ""}},
			{key: "glab repo view", out: `{"id":7,"path_with_namespace":"me/jaira","forked_from_project":{"path_with_namespace":"up/jaira"}}`},
			{key: "glab api projects/up%2Fjaira/merge_requests/12", out: mr},
		}
	}
	ghState := func(s string) string {
		return `{"data":{"repository":{"pullRequest":{"state":"` + s + `"}}}}`
	}
	cases := []struct {
		name  string
		rules []rule
		look  func(string) (string, error)
		want  string
		err   string
		is    error
	}{
		{name: "github open", rules: ghRules(ghState("OPEN")), want: "open"},
		{name: "github merged", rules: ghRules(ghState("MERGED")), want: "merged"},
		{name: "github closed", rules: ghRules(ghState("CLOSED")), want: "closed"},
		{name: "github no such pr", rules: ghRules(`{"data":{"repository":{"pullRequest":null}}}`), err: "has no pull request #12"},
		{name: "github graphql error", rules: ghRules(`{"errors":[{"message":"Could not resolve"}]}`), err: "Could not resolve"},
		{name: "github unreadable", rules: ghRules(`<html>`), err: "unreadable answer"},
		{name: "github unknown state", rules: ghRules(ghState("DRAFTY")), err: `"drafty" is not open`},
		{name: "github signed out", rules: append([]rule{{key: "gh repo view", err: exitErr{4, "HTTP 401"}}}, ghRules("")...), is: ErrAuth},
		{name: "gh missing", rules: ghRules(""), look: missing, is: ErrNoTool},
		{name: "no origin", rules: []rule{{key: ghOrigin, err: exitErr{2, "no such remote"}}}, is: ErrNoForge},
		{name: "foreign host", rules: []rule{{key: ghOrigin, out: "git@example.org:x/y.git"}, {key: "git config jaira.forge", err: exitErr{1, ""}}}, is: ErrNoForge},
		{name: "gitlab opened", rules: gl(`{"state":"opened"}`), want: "open"},
		{name: "gitlab locked", rules: gl(`{"state":"locked"}`), want: "open"},
		{name: "gitlab merged", rules: gl(`{"state":"merged"}`), want: "merged"},
		{name: "gitlab closed", rules: gl(`{"state":"closed"}`), want: "closed"},
		{name: "gitlab unreadable", rules: gl(`nope`), err: "unreadable answer"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := &fake{rules: c.rules}
			look := c.look
			if look == nil {
				look = found
			}
			got, err := State(f.run, look, "/repo", 12)
			switch {
			case c.is != nil:
				if !errors.Is(err, c.is) {
					t.Fatalf("err = %v, want %v", err, c.is)
				}
			case c.err != "":
				if err == nil || !strings.Contains(err.Error(), c.err) {
					t.Fatalf("err = %v, want it to say %q", err, c.err)
				}
			default:
				if err != nil || got != c.want {
					t.Fatalf("State = %q, %v; want %q", got, err, c.want)
				}
			}
			if c.is == nil && c.err == "" && strings.HasPrefix(c.name, "github") {
				// The parent of the fork is asked, with the number as an Int.
				last := f.calls[len(f.calls)-1]
				if !strings.Contains(last, "owner=up") || !strings.Contains(last, "-F number=12") {
					t.Fatalf("asked %q, want the parent up/jaira and number 12", last)
				}
			}
		})
	}
}
