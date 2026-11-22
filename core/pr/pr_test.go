package pr

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type exitErr struct {
	code int
	msg  string
}

func (e exitErr) Error() string { return e.msg }
func (e exitErr) ExitCode() int { return e.code }

// fake answers a command by the first rule whose key is contained in the
// command line; every call is recorded.
type fake struct {
	rules []rule
	calls []string
}

type rule struct {
	key string
	out string
	err error
}

func (f *fake) run(_ string, name string, args ...string) ([]byte, error) {
	line := name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, line)
	for _, r := range f.rules {
		if strings.Contains(line, r.key) {
			return []byte(r.out), r.err
		}
	}
	return nil, exitErr{1, "unexpected call: " + line}
}

func found(string) (string, error)   { return "/usr/bin/x", nil }
func missing(string) (string, error) { return "", errors.New("not found") }

const ghOrigin = "git remote get-url origin"

func ghRules(graphql string) []rule {
	return []rule{
		{key: ghOrigin, out: "git@github.com:me/jaira.git\n"},
		{key: "git config jaira.forge", err: exitErr{1, ""}},
		{key: "gh repo view", out: `{"nameWithOwner":"me/jaira","parent":{"name":"jaira","owner":{"login":"up"}}}`},
		{key: "gh api graphql", out: graphql},
	}
}

func ghAnswerWith(prs ...string) string {
	return `{"data":{"viewer":{"login":"me"},"repository":{"pullRequests":{"nodes":[` + strings.Join(prs, ",") + `]}}}}`
}

// ghPRJSON builds one pull request node; extra is spliced in as further fields.
func ghPRJSON(n int, branch, ci string, extra string) string {
	rollup := "null"
	if ci != "" {
		rollup = fmt.Sprintf(`{"state":%q}`, ci)
	}
	s := fmt.Sprintf(`{"number":%d,"title":"t","url":"u","state":"OPEN","isDraft":false,"headRefName":%q,
	 "subjects":{"nodes":[{"commit":{"messageHeadline":"fix(A3K9QP): x"}},{"commit":{"messageHeadline":"chore: no ticket"}}]},
	 "head":{"nodes":[{"commit":{"statusCheckRollup":%s}}]}`, n, branch, rollup)
	if extra != "" {
		s += "," + extra
	}
	return s + "}"
}

func TestReadInputClasses(t *testing.T) {
	threads := `"reviewThreads":{"nodes":[
	  {"isResolved":false,"comments":{"nodes":[{"createdAt":"2026-11-20T10:00:00Z","author":{"login":"copilot"}}]}},
	  {"isResolved":false,"comments":{"nodes":[{"createdAt":"2026-11-20T10:00:00Z","author":{"login":"copilot"}},{"createdAt":"2026-11-20T11:00:00Z","author":{"login":"me"}}]}},
	  {"isResolved":true,"comments":{"nodes":[{"createdAt":"2026-11-20T09:00:00Z","author":null}]}}]},
	 "comments":{"nodes":[{"createdAt":"2026-11-21T08:00:00Z","author":{"login":"bob"}},{"createdAt":"2026-11-21T09:00:00Z","author":{"login":"me"}}]},
	 "reviews":{"nodes":[{"state":"PENDING","submittedAt":null,"author":{"login":"carol"}},{"state":"COMMENTED","submittedAt":"2026-11-20T10:00:00Z","author":{"login":"copilot"}}]},
	 "reviewRequests":{"nodes":[{"requestedReviewer":{"login":"copilot"}},{"requestedReviewer":{"slug":"core"}},{"requestedReviewer":null}]},
	 "latestOpinionatedReviews":{"nodes":[{"state":"CHANGES_REQUESTED","author":{"login":"bob"}},{"state":"APPROVED","author":{"login":"dan"}}]}`
	cases := []struct {
		name    string
		rules   []rule
		look    func(string) (string, error)
		q       Query
		wantErr error
		errText string
		check   func(t *testing.T, st *Status, calls []string)
	}{
		{
			name:  "no pull request",
			rules: ghRules(ghAnswerWith()),
			q:     Query{Branch: "feat/x"},
			check: func(t *testing.T, st *Status, _ []string) {
				if len(st.PRs) != 0 || st.PRs == nil {
					t.Fatalf("want an empty, non-nil list, got %#v", st.PRs)
				}
			},
		},
		{
			name:  "fork reads the parent, branch passed verbatim",
			rules: ghRules(ghAnswerWith()),
			q:     Query{Branch: "feat/ü bung-$(x)"},
			check: func(t *testing.T, st *Status, calls []string) {
				if st.Repo != "up/jaira" {
					t.Fatalf("repo = %q, want the parent up/jaira", st.Repo)
				}
				if !strings.Contains(calls[len(calls)-1], "head=feat/ü bung-$(x)") || !strings.Contains(calls[len(calls)-1], "owner=up") {
					t.Fatalf("query call = %q", calls[len(calls)-1])
				}
			},
		},
		{
			name:  "several, threads, comments, reviews",
			rules: ghRules(ghAnswerWith(ghPRJSON(7, "feat/craszs", "FAILURE", threads), ghPRJSON(8, "feat/other", "", ""))),
			check: func(t *testing.T, st *Status, _ []string) {
				if len(st.PRs) != 2 {
					t.Fatalf("want 2 pull requests, got %d", len(st.PRs))
				}
				p := st.PRs[0]
				if p.Threads != (Threads{Waiting: 1, Answered: 1, Resolved: 1}) {
					t.Errorf("threads = %+v", p.Threads)
				}
				// copilot twice in threads, ghost once, bob once; me never.
				if p.Comments != 4 {
					t.Errorf("comments = %d, want 4", p.Comments)
				}
				if p.LastForeign == nil || p.LastForeign.Author != "bob" || p.LastForeign.Kind != "comment" {
					t.Errorf("last foreign = %+v", p.LastForeign)
				}
				if p.CI != "failure" {
					t.Errorf("ci = %q", p.CI)
				}
				if !reflect.DeepEqual(p.Tickets, []string{"A3K9QP", "CRASZS"}) {
					t.Errorf("tickets = %v", p.Tickets)
				}
				if !reflect.DeepEqual(p.ReviewRequested, []string{"copilot", "core"}) || !reflect.DeepEqual(p.ChangesRequested, []string{"bob"}) || !reflect.DeepEqual(p.Approved, []string{"dan"}) {
					t.Errorf("reviews: requested %v changes %v approved %v", p.ReviewRequested, p.ChangesRequested, p.Approved)
				}
				if p.Settled || len(p.Unsettled) != 4 {
					t.Errorf("settled %v, unsettled %v", p.Settled, p.Unsettled)
				}
				q := st.PRs[1]
				if q.CI != "none" || !q.Settled || q.LastForeign != nil || q.Tickets == nil {
					t.Errorf("second = %+v", q)
				}
			},
		},
		{
			name:    "no forge tool",
			rules:   ghRules(""),
			look:    missing,
			wantErr: ErrNoTool,
		},
		{
			name:    "not signed in (gh exit 4)",
			rules:   append([]rule{{key: "gh repo view", err: exitErr{4, "To get started with GitHub CLI, please run:  gh auth login"}}}, ghRules("")...),
			wantErr: ErrAuth,
		},
		{
			name:    "no origin",
			rules:   []rule{{key: ghOrigin, err: exitErr{2, "error: No such remote 'origin'"}}},
			wantErr: ErrNoForge,
		},
		{
			name:    "origin on another host",
			rules:   []rule{{key: ghOrigin, out: "https://git.example.com/a/b.git"}, {key: "git config", err: exitErr{1, ""}}},
			wantErr: ErrNoForge,
		},
		{
			name:    "empty answer",
			rules:   ghRules(""),
			errText: "unreadable answer",
		},
		{
			name:    "broken JSON",
			rules:   ghRules(`{"data":`),
			errText: "unreadable answer",
		},
		{
			name:    "graphql error",
			rules:   ghRules(`{"errors":[{"message":"Could not resolve to a Repository"}]}`),
			errText: "Could not resolve",
		},
		{
			name:    "forge fails half way",
			rules:   append([]rule{{key: "gh api graphql", err: exitErr{1, "HTTP 502: Bad Gateway"}}}, ghRules("")...),
			errText: "HTTP 502",
		},
		{
			name:    "number not found",
			rules:   ghRules(`{"data":{"viewer":{"login":"me"},"repository":{"pullRequest":null}}}`),
			q:       Query{Number: 9},
			errText: "no pull request #9",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := &fake{rules: c.rules}
			look := c.look
			if look == nil {
				look = found
			}
			st, err := Read(f.run, look, c.q)
			switch {
			case c.wantErr != nil:
				if !errors.Is(err, c.wantErr) {
					t.Fatalf("err = %v, want %v", err, c.wantErr)
				}
				return
			case c.errText != "":
				if err == nil || !strings.Contains(err.Error(), c.errText) {
					t.Fatalf("err = %v, want it to say %q", err, c.errText)
				}
				return
			case err != nil:
				t.Fatal(err)
			}
			c.check(t, st, f.calls)
		})
	}
}

func TestReadGitLab(t *testing.T) {
	disc := `[
	 {"individual_note":false,"notes":[{"author":{"username":"rev"},"created_at":"2026-11-20T10:00:00Z","resolvable":true,"resolved":false}]},
	 {"individual_note":false,"notes":[{"author":{"username":"rev"},"created_at":"2026-11-20T10:00:00Z","resolvable":true,"resolved":false},{"author":{"username":"me"},"created_at":"2026-11-20T12:00:00Z","resolvable":true,"resolved":false}]},
	 {"individual_note":false,"notes":[{"author":{"username":"rev"},"created_at":"2026-11-19T10:00:00Z","resolvable":true,"resolved":true}]},
	 {"individual_note":true,"notes":[{"author":{"username":"bob"},"created_at":"2026-11-21T10:00:00Z","resolvable":false}]},
	 {"individual_note":true,"notes":[{"author":{"username":"bot"},"created_at":"2026-11-22T10:00:00Z","system":true}]}]`
	base := func() []rule {
		return []rule{
			{key: ghOrigin, out: "git@gitlab.example.com:grp/sub/jaira.git"},
			{key: "git config jaira.forge", out: "gitlab\n"},
			{key: "glab repo view", out: `{"path_with_namespace":"grp/sub/jaira"}`},
			{key: "glab api user", out: `{"username":"me"}`},
			{key: "/discussions", out: disc},
			{key: "/reviewers", out: `[{"user":{"username":"rev"},"state":"requested_changes"},{"user":{"username":"ann"},"state":"unreviewed"}]`},
			{key: "/commits", out: `[{"title":"feat(B7XQ2M): y"}]`},
			{key: "merge_requests?state=opened", out: `[{"iid":3,"source_branch":"feat/b7xq2m"}]`},
			{key: "merge_requests/3", out: `{"iid":3,"title":"t","web_url":"u","state":"opened","source_branch":"feat/b7xq2m","head_pipeline":{"status":"running"}}`},
		}
	}
	f := &fake{rules: base()}
	st, err := Read(f.run, found, Query{Branch: "feat/b7xq2m"})
	if err != nil {
		t.Fatal(err)
	}
	if st.Repo != "grp/sub/jaira" || len(st.PRs) != 1 {
		t.Fatalf("status = %+v", st)
	}
	p := st.PRs[0]
	if p.Threads != (Threads{Waiting: 1, Answered: 1, Resolved: 1}) || p.Comments != 4 || p.CI != "pending" || p.State != "open" {
		t.Errorf("pr = %+v", p)
	}
	if !reflect.DeepEqual(p.Tickets, []string{"B7XQ2M"}) || p.LastForeign.Author != "bob" {
		t.Errorf("tickets %v, last %+v", p.Tickets, p.LastForeign)
	}
	if !reflect.DeepEqual(p.ChangesRequested, []string{"rev"}) || !reflect.DeepEqual(p.ReviewRequested, []string{"ann"}) || p.Settled {
		t.Errorf("reviews %v / %v settled %v", p.ChangesRequested, p.ReviewRequested, p.Settled)
	}
	if !strings.Contains(strings.Join(f.calls, "\n"), "projects/grp%2Fsub%2Fjaira/merge_requests?state=opened&per_page=50&order_by=updated_at&source_branch=feat%2Fb7xq2m") {
		t.Errorf("list call not escaped:\n%s", strings.Join(f.calls, "\n"))
	}

	// A call failing half way — after the list, at the discussions — is an
	// error, never a pull request read without its threads.
	f = &fake{rules: append([]rule{{key: "/discussions", err: exitErr{1, "HTTP 500"}}}, base()...)}
	if _, err := Read(f.run, found, Query{}); err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		t.Fatalf("half way: err = %v", err)
	}
	// glab says 401 when its token is gone.
	f = &fake{rules: append([]rule{{key: "glab api user", err: exitErr{1, "glab: 401 Unauthorized"}}}, base()...)}
	if _, err := Read(f.run, found, Query{}); !errors.Is(err, ErrAuth) {
		t.Fatalf("401: err = %v", err)
	}
}

func TestChanges(t *testing.T) {
	base := PR{Number: 7, Tickets: []string{"A3K9QP"}, State: "open", CI: "pending",
		LastForeign: &Activity{At: "2026-11-20T10:00:00Z", Author: "bob", Kind: "comment"}}
	with := func(f func(p *PR)) PR {
		p := base
		f(&p)
		return p
	}
	cases := []struct {
		name  string
		after PR
		want  []string
	}{
		{"nothing changed", base, nil},
		{"same activity read again", with(func(p *PR) { p.LastForeign = &Activity{At: "2026-11-20T10:00:00Z", Author: "bob", Kind: "comment"} }), nil},
		{"own reply only", with(func(p *PR) { p.Threads.Answered = 1 }), nil},
		{"new comment", with(func(p *PR) { p.LastForeign = &Activity{At: "2026-11-20T11:00:00Z", Author: "eve", Kind: "comment"} }), []string{"#7 (A3K9QP): new comment by eve"}},
		{"new thread comment", with(func(p *PR) {
			p.Threads.Waiting = 2
			p.LastForeign = &Activity{At: "2026-11-20T11:00:00Z", Author: "copilot", Kind: "thread"}
		}), []string{"#7 (A3K9QP): new review comment by copilot — threads waiting: 2"}},
		{"review", with(func(p *PR) {
			p.LastForeign = &Activity{At: "2026-11-20T11:00:00Z", Author: "dan", Kind: "review", State: "changes_requested"}
		}), []string{"#7 (A3K9QP): review by dan: changes requested"}},
		{"CI red", with(func(p *PR) { p.CI = "failure" }), []string{"#7 (A3K9QP): CI red"}},
		{"CI green from pending says nothing", with(func(p *PR) { p.CI = "success" }), nil},
		{"settled and merged", with(func(p *PR) { p.Settled = true; p.State = "merged" }), []string{"#7 (A3K9QP): settled — no thread open, no review pending", "#7 (A3K9QP): merged"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Changes(base, c.after); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
	red := with(func(p *PR) { p.CI = "failure" })
	if got := Changes(red, with(func(p *PR) { p.CI = "success" })); !reflect.DeepEqual(got, []string{"#7 (A3K9QP): CI green again"}) {
		t.Fatalf("green again: %q", got)
	}
	// A first comment on a pull request that had none is a change too.
	none := with(func(p *PR) { p.LastForeign = nil })
	if got := Changes(none, base); len(got) != 1 {
		t.Fatalf("first comment: %q", got)
	}
}
