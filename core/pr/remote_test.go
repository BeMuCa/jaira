package pr

import (
	"strings"
	"testing"
)

func TestParseRemote(t *testing.T) {
	for _, c := range []struct {
		in, host, repo string
	}{
		{"https://github.com/me/jaira.git", "github.com", "me/jaira"},
		{"https://github.com/me/jaira", "github.com", "me/jaira"},
		{"https://token@github.com/me/jaira.git/", "github.com", "me/jaira"},
		{"https://ghe.example.com:8443/team/app.git", "ghe.example.com", "team/app"},
		{"ssh://git@github.com/me/jaira.git", "github.com", "me/jaira"},
		{"ssh://alice@ghe.example.com:2222/team/app.git", "ghe.example.com", "team/app"},
		{"git+ssh://git@GHE.Example.com/team/app", "ghe.example.com", "team/app"},
		{"git@github.com:me/jaira.git", "github.com", "me/jaira"},
		{"git@github.com:me/jaira", "github.com", "me/jaira"},
		{"alice@ghe.example.com:team/app.git", "ghe.example.com", "team/app"},
		{"svc-build@ghe.example.com:team/app", "ghe.example.com", "team/app"},
		{"ghe.example.com:team/app.git", "ghe.example.com", "team/app"},
		{"git@gitlab.example.com:group/sub/project.git", "gitlab.example.com", "group/sub/project"},
		{"  git@github.com:me/jaira.git\n", "github.com", "me/jaira"},
	} {
		r, ok := parseRemote(c.in)
		if !ok || r.Host != c.host || r.Repo != c.repo {
			t.Errorf("parseRemote(%q) = %+v, %v; want %s %s", c.in, r, ok, c.host, c.repo)
		}
	}
	for _, in := range []string{"", "/srv/git/app.git", "../app", `C:\repos\app`, "C:/repos/app", "https://ghe.example.com/app", "git@host:app.git", "file:///srv/git/app.git"} {
		if r, ok := parseRemote(in); ok {
			t.Errorf("parseRemote(%q) = %+v, want no remote", in, r)
		}
	}
}

// On GitHub Enterprise gh is told the host: the repository as
// host/owner/repo, --hostname on every api call. github.com stays as it was.
func TestReadGitHubEnterprise(t *testing.T) {
	for _, c := range []struct {
		origin, view, api string
	}{
		{"alice@ghe.example.com:team/app.git", "gh repo view ghe.example.com/team/app ", "gh api --hostname ghe.example.com graphql"},
		{"git@github.com:me/jaira.git", "gh repo view git@github.com:me/jaira.git ", "gh api graphql"},
	} {
		f := &fake{rules: []rule{
			{key: ghOrigin, out: c.origin + "\n"},
			{key: "git config jaira.forge", out: "github\n"},
			{key: "gh repo view", out: `{"nameWithOwner":"team/app"}`},
			{key: "gh api", out: ghAnswerWith()},
		}}
		if _, err := Read(f.run, found, Query{}); err != nil {
			t.Fatalf("%s: %v", c.origin, err)
		}
		calls := strings.Join(f.calls, "\n")
		if !strings.Contains(calls, c.view) || !strings.Contains(calls, c.api) {
			t.Fatalf("%s: want %q and %q in\n%s", c.origin, c.view, c.api, calls)
		}
		f.calls = nil
		f.rules[3] = rule{key: "gh api", out: `{"data":{"repository":{"pullRequest":{"state":"MERGED"}}}}`}
		if st, err := State(f.run, found, "", 7); err != nil || st != "merged" {
			t.Fatalf("%s: State = %q, %v", c.origin, st, err)
		}
		if calls := strings.Join(f.calls, "\n"); !strings.Contains(calls, c.api) {
			t.Fatalf("%s: State asked without the host:\n%s", c.origin, calls)
		}
	}
}
