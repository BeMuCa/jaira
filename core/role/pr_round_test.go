package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fakeGH is a gh that records every call, one line per call with its
// arguments joined by a tab, and answers from the environment: GH_AUTH=fail
// makes 'auth status' fail, GH_POST=fail makes every --method POST fail, and
// GH_FIXTURE is the JSON a read answers with, run through the real jq with
// the filter the script hands --jq — so the filter itself is what is tested;
// GH_ISSUES answers the conversation's comments instead, and GH_STATE the
// query reply asks before it writes.
const fakeGH = `#!/usr/bin/env bash
printf '%s\t' "$@" >> "$GH_LOG"; printf '\n' >> "$GH_LOG"
[ "$1" = auth ] && { [ "${GH_AUTH:-}" = fail ] && exit 1; exit 0; }
for a in "$@"; do [ "$a" = POST ] && [ "${GH_POST:-}" = fail ] && { echo "HTTP 422: Reviews may only be requested from collaborators" >&2; exit 1; }; done
filter="" fixture=${GH_FIXTURE:-}
case "$*" in *"/issues/"*) fixture=${GH_ISSUES:-} ;; *"viewer{login}"*) fixture=${GH_STATE:-} ;; esac
while [ $# -gt 0 ]; do [ "$1" = --jq ] && filter=$2; shift; done
if [ -n "$filter" ] && [ -n "$fixture" ]; then jq -r "$filter" "$fixture"; else echo '{}'; fi
`

type prRound struct {
	t      *testing.T
	bash   string
	script string
	log    string
	env    []string
}

func newPRRound(t *testing.T) *prRound {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skipf("no jq on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(fakeGH), 0o755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "gh.log")
	return &prRound{
		t:      t,
		bash:   bash,
		script: filepath.Join(dir, "jaira-role-pr", "scripts", "pr-round.sh"),
		log:    log,
		env:    append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "GH_LOG="+log),
	}
}

func (p *prRound) run(env []string, args ...string) (string, int) {
	p.t.Helper()
	_ = os.Remove(p.log)
	cmd := exec.Command(p.bash, append([]string{p.script}, args...)...)
	cmd.Env = append(append([]string{}, p.env...), env...)
	out, err := cmd.CombinedOutput()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		p.t.Fatal(err)
	}
	return string(out), code
}

func (p *prRound) calls() []string {
	b, err := os.ReadFile(p.log)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(b)), "\n")
}

func fixtureFile(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// The input classes the script must refuse before it writes anything: a
// malformed target, an empty answer, and gh that is not signed in (JG3REW).
func TestPRRoundRefusesBeforeItWrites(t *testing.T) {
	p := newPRRound(t)
	answer := fixtureFile(t, "answer with space.md", "Fixed in abc123: renamed it.\n")
	empty := fixtureFile(t, "empty.md", "")
	for _, c := range []struct {
		name string
		env  []string
		args []string
		code int
	}{
		{"no owner/repo", nil, []string{"threads", "repo", "7"}, 2},
		{"number not a number", nil, []string{"threads", "o/r", "7;id"}, 2},
		{"empty number", nil, []string{"threads", "o/r", ""}, 2},
		{"unknown command", nil, []string{"merge", "o/r", "7"}, 2},
		{"too few arguments", nil, []string{"reply", "o/r", "7", "PRRT_1"}, 2},
		{"empty answer", nil, []string{"reply", "o/r", "7", "PRRT_1", empty}, 2},
		{"missing answer", nil, []string{"comment", "o/r", "7", filepath.Join(t.TempDir(), "nope")}, 2},
		{"not signed in", []string{"GH_AUTH=fail"}, []string{"reply", "o/r", "7", "PRRT_1", answer}, 3},
	} {
		out, code := p.run(c.env, c.args...)
		if code != c.code {
			t.Errorf("%s: exit %d, want %d\n%s", c.name, code, c.code, out)
		}
		for _, call := range p.calls() {
			if strings.Contains(call, "mutation") || strings.Contains(call, "POST") || strings.HasPrefix(call, "pr\t") {
				t.Errorf("%s: wrote before refusing: %q", c.name, call)
			}
		}
	}
	if out, code := p.run([]string{"GH_AUTH=fail"}, "threads", "o/r", "7"); code != 3 || !strings.Contains(out, "gh auth login") {
		t.Errorf("not signed in: exit %d, want 3 naming gh auth login\n%s", code, out)
	}
}

// A second round lists only what is still open: a thread resolved in the
// first round is not answered twice. Comments without a thread are listed
// apart, because they cannot be resolved; an empty review body is not one.
func TestPRRoundListsOnlyWhatIsOpen(t *testing.T) {
	p := newPRRound(t)
	threads := fixtureFile(t, "threads.json", `{"data":{"repository":{"pullRequest":{"reviewThreads":{
		"pageInfo":{"hasNextPage":false,"endCursor":null},
		"nodes":[
		 {"id":"PRRT_done","isResolved":true,"isOutdated":false,"path":"a.go","line":3,"comments":{"nodes":[{"author":{"login":"copilot"},"body":"x","url":"u1"}]}},
		 {"id":"PRRT_open","isResolved":false,"isOutdated":true,"path":"dir with space/ü.go","line":null,"comments":{"nodes":[{"author":null,"body":"line\r\nwith \"quotes\" and $(id)","url":"u2"}]}}
		]}}}}}`)
	out, code := p.run([]string{"GH_FIXTURE=" + threads}, "threads", "o/r", "7")
	if !strings.Contains(strings.Join(p.calls(), "\n"), "-f\towner=o\t-f\tname=r\t") {
		t.Errorf("threads must hand owner and name as strings (-f):\n%s", strings.Join(p.calls(), "\n"))
	}
	if code != 0 || strings.Contains(out, "PRRT_done") || !strings.Contains(out, "PRRT_open") || !strings.Contains(out, `"author":"ghost"`) {
		t.Errorf("threads: exit %d, want only PRRT_open with a ghost author\n%s", code, out)
	}
	if lines := strings.Count(strings.TrimSpace(out), "\n") + 1; lines != 1 {
		t.Errorf("threads: %d lines, want one JSON object per thread\n%s", lines, out)
	}

	reviews := fixtureFile(t, "reviews.json", `[{"user":{"login":"copilot-pull-request-reviewer[bot]"},"body":"Overview","html_url":"r1"},{"user":{"login":"me"},"body":"","html_url":"r2"}]`)
	issues := fixtureFile(t, "issues.json", `[{"user":{"login":"alex"},"body":"Why this name?","html_url":"c1"}]`)
	out, code = p.run([]string{"GH_FIXTURE=" + reviews, "GH_ISSUES=" + issues}, "comments", "o/r", "7")
	if code != 0 || !strings.Contains(out, `"url":"r1"`) || strings.Contains(out, `"url":"r2"`) || !strings.Contains(out, `"kind":"review"`) ||
		!strings.Contains(out, `{"kind":"conversation","author":"alex","body":"Why this name?","url":"c1"}`) {
		t.Errorf("comments: exit %d, want the review with a body and not the empty one\n%s", code, out)
	}
}

// Answering a thread replies and then resolves it, with the answer read from
// a file so no quoting of the reviewer's words or ours reaches a shell.
func TestPRRoundRepliesThenResolves(t *testing.T) {
	p := newPRRound(t)
	answer := fixtureFile(t, "answer.md", "Fixed in abc123: `x` -> \"y\"; $(id)\r\n")
	open := fixtureFile(t, "open.json", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"copilot-pull-request-reviewer"}}]}}}}`)
	out, code := p.run([]string{"GH_STATE=" + open}, "reply", "o/r", "7", "PRRT_open", answer)
	if code != 0 {
		t.Fatalf("reply: exit %d\n%s", code, out)
	}
	calls := strings.Join(p.calls(), "\n")
	reply := strings.Index(calls, "addPullRequestReviewThreadReply")
	resolve := strings.Index(calls, "resolveReviewThread")
	if reply < 0 || resolve < 0 || reply > resolve {
		t.Errorf("reply must answer the thread and then resolve it:\n%s", calls)
	}
	if !strings.Contains(calls, "body=@"+answer) {
		t.Errorf("reply must hand gh the answer as a file:\n%s", calls)
	}
	if strings.Contains(calls, "merge") || strings.Contains(calls, "APPROVE") || strings.Contains(calls, "--approve") {
		t.Errorf("pr-round.sh must never merge or approve:\n%s", calls)
	}
}

// Copilot is asked for by its bot login; a reviewer that cannot be requested
// is reported with its own exit code and leaves the pull request as it is.
func TestPRRoundRequestsTheReviewer(t *testing.T) {
	p := newPRRound(t)
	for reviewer, want := range map[string]string{
		"copilot":         "reviewers[]=copilot-pull-request-reviewer[bot]",
		"@copilot":        "reviewers[]=copilot-pull-request-reviewer[bot]",
		"monalisa":        "reviewers[]=monalisa",
		"@monalisa":       "reviewers[]=monalisa",
		"myorg/team-name": "team_reviewers[]=team-name",
	} {
		out, code := p.run(nil, "request", "o/r", "7", reviewer)
		calls := strings.Join(p.calls(), "\n")
		if code != 0 || !strings.Contains(calls, "repos/o/r/pulls/7/requested_reviewers") || !strings.Contains(calls, want) {
			t.Errorf("request %s: exit %d, want %s\n%s\n%s", reviewer, code, want, out, calls)
		}
	}
	out, code := p.run([]string{"GH_POST=fail"}, "request", "o/r", "7", "copilot")
	if code != 4 || !strings.Contains(out, "422") {
		t.Errorf("unavailable reviewer: exit %d, want 4 with gh's reason\n%s", code, out)
	}
}

// The script is the role's, never a way to accept work: no merge, no approve.
func TestPRRoundNeverMergesOrApproves(t *testing.T) {
	body, err := File("jaira-role-pr", "scripts/pr-round.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"pr merge", "mr merge", "--approve", "APPROVE", "mr approve", "mergePullRequest"} {
		if strings.Contains(string(body), bad) {
			t.Errorf("pr-round.sh contains %q", bad)
		}
	}
}

// Without .jaira/pr an agent still pushes and hands the create line back;
// with "agent" jaira-role-pr opens the pull request only once every ticket on
// the branch is in done, asks for the board's reviewer, and runs one round
// through pr-round.sh. Merge and approve stay forbidden under any setting
// (JG3REW).
func TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts(t *testing.T) {
	body, err := File("jaira-role-pr", "SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	setting := strings.Index(text, "## Whether this board lets an agent open it")
	open := strings.Index(text, "## Open it")
	if setting < 0 || open < 0 || setting > open {
		t.Fatalf("jaira-role-pr must settle the board's pr setting before ## Open it (at %d, %d)", setting, open)
	}
	for _, want := range []string{
		"jaira roles list --json | jaira json get pr",
		"jaira roles list --json | jaira json get pr_reviewer",
		"**`person`** — the default",
		"every ticket the branch's commits name is in the board's",
		"you are back to `person` for\n  this run",
		"or the board says\n`agent` and every ticket on the branch is in done",
		"scripts/pr-round.sh request <owner/repo> <number> <pr_reviewer>",
		"scripts/pr-round.sh threads  <owner/repo> <number>",
		"scripts/pr-round.sh reply <owner/repo> <number> <thread id> <file>",
		"`Fixed in <sha>: <what changed>`",
		"`.jaira/pr` never reaches merge or approve.",
		"never `gh pr merge`. Never `gh pr review --approve`.",
		"never `glab mr merge`. Never `glab mr approve`.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("jaira-role-pr does not say %q", want)
		}
	}
	for _, id := range []string{"jaira-teamlead", "jaira-dispatcher"} {
		b, err := File(id, "SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "`agent`") || !strings.Contains(string(b), "in done") {
			t.Errorf("%s still says only a person opens a pull request, without the board's exception", id)
		}
	}
}

// Answer and resolve are two calls. A round run again after the resolve
// failed finds the thread open with our answer last, and only resolves it; a
// thread already resolved is left alone; a thread that does not exist is
// refused before anything is written. The class: no thread ever gets the same
// answer twice, and no question goes unanswered because it was written under
// the same login the agent works with — the maintainer's (JG3REW, bughunt).
func TestPRRoundReplyIsSafeToRunAgain(t *testing.T) {
	p := newPRRound(t)
	answer := fixtureFile(t, "answer.md", "Fixed in abc123.\n")
	for _, c := range []struct {
		name, state     string
		code            int
		posts, resolves bool
	}{
		{"answered, resolve failed", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"me"},"body":"Fixed in abc123."}]}}}}`, 0, false, true},
		{"answered with CRLF", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"me"},"body":"Fixed in abc123.\r\n"}]}}}}`, 0, false, true},
		{"the maintainer's own remark", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"me"},"body":"please rename this"}]}}}}`, 0, true, true},
		{"someone else wrote the same words", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"copilot"},"body":"Fixed in abc123."}]}}}}`, 0, true, true},
		{"already resolved", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":true,"comments":{"nodes":[{"author":{"login":"me"}}]}}}}`, 0, false, false},
		{"no such thread", `{"data":{"viewer":{"login":"me"},"node":null}}`, 1, false, false},
		{"ghost author", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":null}]}}}}`, 0, true, true},
	} {
		state := fixtureFile(t, "state.json", c.state)
		out, code := p.run([]string{"GH_STATE=" + state}, "reply", "o/r", "7", "PRRT_x", answer)
		calls := strings.Join(p.calls(), "\n")
		if code != c.code {
			t.Errorf("%s: exit %d, want %d\n%s", c.name, code, c.code, out)
		}
		if got := strings.Contains(calls, "addPullRequestReviewThreadReply"); got != c.posts {
			t.Errorf("%s: posted an answer %v, want %v\n%s", c.name, got, c.posts, calls)
		}
		if got := strings.Contains(calls, "resolveReviewThread"); got != c.resolves {
			t.Errorf("%s: resolved %v, want %v\n%s", c.name, got, c.resolves, calls)
		}
	}
}
