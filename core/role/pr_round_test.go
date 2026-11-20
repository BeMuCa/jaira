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
// query reply asks before it writes, GH_AFTER the thread it reads back after
// it posted; GH_POSTED answers a conversation comment's POST, and
// GH_ISSUES_AFTER the conversation read back once something was posted;
// GH_REQUESTED the reviewers a pull request still waits for, GH_OPINIONS each
// reviewer's latest verdict. GH_LOST=1 makes every call after 'auth status'
// end with 4, the way gh does once its sign-in has run out; GH_LOST=after
// does so only once something was posted.
const fakeGH = `#!/usr/bin/env bash
printf '%s\t' "$@" >> "$GH_LOG"; printf '\n' >> "$GH_LOG"
[ "$1" = auth ] && { [ "${GH_AUTH:-}" = fail ] && exit 1; exit 0; }
[ -n "${GH_LOST:-}" ] && { [ "$GH_LOST" != after ] || { grep -q POST "$GH_LOG" && ! printf '%s ' "$@" | grep -q POST; }; } && { echo "To get started with GitHub CLI, please run:  gh auth login" >&2; exit 4; }
for a in "$@"; do [ "$a" = POST ] && [ "${GH_POST:-}" = fail ] && { echo "HTTP 422: Reviews may only be requested from collaborators" >&2; exit 1; }; done
filter="" fixture=${GH_FIXTURE:-}
case "$*" in *"latestOpinionatedReviews"*) fixture=${GH_OPINIONS:-} ;; *"requested_reviewers"*) fixture=${GH_REQUESTED:-} ;; *"POST repos/"*"/issues/"*) fixture=${GH_POSTED:-} ;; *"/issues/"*) fixture=${GH_ISSUES:-}; [ -n "${GH_ISSUES_AFTER:-}" ] && grep -q POST "$GH_LOG" && fixture=$GH_ISSUES_AFTER ;; *"comments(last:100)"*) fixture=${GH_AFTER:-} ;; *"viewer{login}"*) fixture=${GH_STATE:-} ;; esac
while [ $# -gt 0 ]; do [ "$1" = --jq ] && filter=$2; shift; done
if [ -n "$filter" ] && [ -n "$fixture" ]; then jq -r "$filter" "$fixture"; else echo '{}'; fi
`

// fakeJaira stands in for the binary's gate: JAIRA_GATE=fail refuses the way
// 'jaira roles pr-gate' does, with exit 3 and the ticket that holds it.
const fakeJaira = `#!/usr/bin/env bash
printf '%s\t' jaira "$@" >> "$GH_LOG"; printf '\n' >> "$GH_LOG"
[ "${JAIRA_GATE:-}" = fail ] && { echo "jaira: pr-gate: not every ticket on the branch is in done — B7XQ2M (review)" >&2; exit 3; }
echo "pr-gate: an agent may open it — every ticket on the branch is done: A3K9QP"
`

type prRound struct {
	t      *testing.T
	bash   string
	script string
	log    string
	env    []string
	dir    string
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
	if err := os.WriteFile(filepath.Join(bin, "jaira"), []byte(fakeJaira), 0o755); err != nil {
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
	cmd.Dir = p.dir
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
	blank := fixtureFile(t, "blank.md", "  \r\n\t\n")
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
		{"blanks only", nil, []string{"comment", "o/r", "7", blank}, 2},
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
		"jaira roles pr-gate <remote>/<branch>` at **Open",
		"scripts/pr-round.sh open <owner/repo> <remote> <branch>",
		"jaira roles pr-gate <remote>/<branch> && glab mr create",
		"scripts/pr-round.sh request <owner/repo> <number> <pr_reviewer>",
		"jaira roles list --json | jaira json get pr_when",
		"**done waits for it**",
		"scripts/pr-round.sh settled <owner/repo> <number>",
		"blocking_discussions_resolved",
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
		if !strings.Contains(string(b), "`agent`") || !strings.Contains(string(b), "in done") || !strings.Contains(string(b), "signoff") {
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

// An agent opens the pull request through the script, and the script asks the
// binary first: the gate 'every ticket on the branch is in done' no longer
// rests on the role's text alone. Refused, nothing is opened (JG3REW, given
// back from signoff).
func TestPROpenRunsTheGateFirst(t *testing.T) {
	p := newPRRound(t)
	p.dir = pushedCheckout(t)
	body := fixtureFile(t, "body.md", "## Why\n")
	out, code := p.run([]string{"JAIRA_GATE=fail"}, "open", "up/r", "upstream", "master", "feat(A3K9QP): x", body)
	calls := strings.Join(p.calls(), "\n")
	if code != 3 || !strings.Contains(out, "B7XQ2M (review)") || strings.Contains(calls, "pr\tcreate") {
		t.Errorf("gate refused: exit %d, want 3 naming the ticket and no pr create\n%s\n%s", code, out, calls)
	}
	if !strings.Contains(calls, "jaira\troles\tpr-gate\tupstream/master") {
		t.Errorf("open must ask the gate about <remote>/<branch>:\n%s", calls)
	}
	out, code = p.run(nil, "open", "up/r", "upstream", "master", "feat(A3K9QP): x", body)
	calls = strings.Join(p.calls(), "\n")
	gate, create := strings.Index(calls, "pr-gate"), strings.Index(calls, "pr\tcreate\t--repo\tup/r\t--base\tmaster\t")
	if code != 0 || gate < 0 || create < gate {
		t.Errorf("gate passed: exit %d, want the gate and then pr create into up/r master\n%s\n%s", code, out, calls)
	}
	for _, args := range [][]string{
		{"open", "up/r", "upstream", "master", "title"},
		{"open", "up/r", "upstream", "master", "title", filepath.Join(t.TempDir(), "nope")},
	} {
		if out, code := p.run(nil, args...); code != 2 || strings.Contains(strings.Join(p.calls(), "\n"), "create") {
			t.Errorf("%v: exit %d, want 2 and nothing opened\n%s", args, code, out)
		}
	}
}

// pushedCheckout is a repository whose branch feat/x is pushed to origin, a
// bare repository beside it, so the checkout is the pushed tip.
func pushedCheckout(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	work := filepath.Join(root, "work")
	gitIn(t, root, "init", "-q", "--bare", "origin.git")
	gitIn(t, root, "init", "-q", "-b", "master", "work")
	gitIn(t, work, "commit", "-q", "--allow-empty", "-m", "chore: base")
	gitIn(t, work, "remote", "add", "origin", filepath.Join(root, "origin.git"))
	gitIn(t, work, "checkout", "-q", "-b", "feat/x")
	gitIn(t, work, "commit", "-q", "--allow-empty", "-m", "feat(A3K9QP): x")
	gitIn(t, work, "push", "-q", "-u", "origin", "HEAD")
	return work
}

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.email=a@b", "-c", "user.name=x", "-c", "commit.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// The gate reads the checkout's commits while gh opens the branch as pushed:
// open refuses unless they are the same commit, so commits another session
// pushed since are never opened past a gate that did not see them, and a
// local commit not pushed yet never passes it on the pull request's behalf
// (JG3REW review-gaps 1).
func TestPROpenOnlyFromThePushedTip(t *testing.T) {
	p := newPRRound(t)
	body := fixtureFile(t, "body.md", "## Why\n")
	work := pushedCheckout(t)
	p.dir = work
	gitIn(t, work, "commit", "-q", "--allow-empty", "-m", "feat(B7XQ2M): not pushed")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 ||
		!strings.Contains(out, "pull or push until they agree") || strings.Contains(strings.Join(p.calls(), "\n"), "create") {
		t.Errorf("checkout ahead of the push: exit %d, want 3 and nothing opened\n%s", code, out)
	}
	gitIn(t, work, "push", "-q")
	gitIn(t, work, "reset", "-q", "--hard", "HEAD~1")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 ||
		!strings.Contains(out, "pull or push until they agree") || strings.Contains(strings.Join(p.calls(), "\n"), "pr-gate") {
		t.Errorf("push ahead of the checkout: exit %d, want 3 before the gate runs\n%s", code, out)
	}
	// A single-branch clone has no remote-tracking ref for feat/x, and the
	// branch is pushed all the same: the push target comes from the
	// configuration.
	clone := filepath.Join(t.TempDir(), "clone")
	gitIn(t, work, "push", "-q", "origin", "master")
	gitIn(t, filepath.Dir(clone), "clone", "-q", "--single-branch", "-b", "master", filepath.Join(filepath.Dir(work), "origin.git"), clone)
	gitIn(t, clone, "checkout", "-q", "-b", "feat/y")
	gitIn(t, clone, "commit", "-q", "--allow-empty", "-m", "feat(A3K9QP): y")
	gitIn(t, clone, "push", "-q", "-u", "origin", "HEAD")
	p.dir = clone
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 0 || !strings.Contains(strings.Join(p.calls(), "\n"), "pr\tcreate") {
		t.Errorf("single-branch clone, pushed: exit %d, want the pull request opened\n%s", code, out)
	}
	// push.default upstream, or tracking, its older name, pushes a branch
	// under the name it tracks, not its own.
	gitIn(t, clone, "checkout", "-q", "-b", "feat/z")
	gitIn(t, clone, "push", "-q", "origin", "HEAD:refs/heads/remote-z")
	gitIn(t, clone, "config", "branch.feat/z.remote", "origin")
	gitIn(t, clone, "config", "branch.feat/z.merge", "refs/heads/remote-z")
	for _, mode := range []string{"upstream", "tracking"} {
		gitIn(t, clone, "config", "push.default", mode)
		if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 0 {
			t.Errorf("push.default=%s, branch pushed as remote-z: exit %d, want the pull request opened\n%s", mode, code, out)
		}
	}
	p.dir = work
	gitIn(t, work, "checkout", "-q", "-b", "feat/local")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 || !strings.Contains(out, "push it first") {
		t.Errorf("branch never pushed: exit %d, want 3 saying to push first\n%s", code, out)
	}
	gitIn(t, work, "checkout", "-q", "--detach")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 || !strings.Contains(out, "detached") {
		t.Errorf("detached HEAD: exit %d, want 3\n%s", code, out)
	}
}

// open names the head itself — the branch the tip check read, on the account
// it was pushed to — instead of leaving gh to guess it from the checkout: in
// a fork the guess need not be the commits the gate judged (JG3REW, given
// back from signoff, Alex 2026-11-20).
func TestPROpenNamesTheHeadItChecked(t *testing.T) {
	p := newPRRound(t)
	body := fixtureFile(t, "body.md", "## Why\n")
	work := pushedCheckout(t)
	p.dir = work
	bare := filepath.Join(filepath.Dir(work), "origin.git")
	for _, c := range []struct{ url, repo, head string }{
		{"https://github.com/me/r.git", "up/r", "me:feat/x"},
		{"git@github.com:me/r.git", "up/r", "me:feat/x"},
		{"ssh://git@github.com/me/r/", "up/r", "me:feat/x"},
		{"https://github.com/me/r", "me/r", "feat/x"},
		{"git@github.com:Me/r.git", "me/r", "feat/x"},
	} {
		// The configured URL is GitHub's; insteadOf sends the push and the
		// tip check to the bare repository beside the checkout.
		gitIn(t, work, "config", "remote.origin.url", c.url)
		gitIn(t, work, "config", "url."+bare+".insteadOf", c.url)
		out, code := p.run(nil, "open", c.repo, "origin", "master", "t", body)
		calls := strings.Join(p.calls(), "\n")
		if code != 0 || !strings.Contains(calls, "pr\tcreate\t--repo\t"+c.repo+"\t--base\tmaster\t--head\t"+c.head+"\t") {
			t.Errorf("%s into %s: exit %d, want --head %s\n%s\n%s", c.url, c.repo, code, c.head, out, calls)
		}
		gitIn(t, work, "config", "--unset", "url."+bare+".insteadOf")
	}
	// A URL with no account in it opens nothing rather than a guess.
	gitIn(t, work, "config", "remote.origin.url", "nowhere")
	gitIn(t, work, "config", "url."+bare+".insteadOf", "nowhere")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 ||
		!strings.Contains(out, "which account") || strings.Contains(strings.Join(p.calls(), "\n"), "pr\tcreate") {
		t.Errorf("URL without an account: exit %d, want 3 and nothing opened\n%s", code, out)
	}
	// A pushurl apart from the URL fetched from: the tip is read, and the
	// account taken, where git push sends the branch.
	gitIn(t, work, "config", "remote.origin.url", "https://github.com/up/r.git")
	gitIn(t, work, "config", "remote.origin.pushurl", "git@github.com:me/r.git")
	gitIn(t, work, "config", "url."+bare+".insteadOf", "git@github.com:me/r.git")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 0 ||
		!strings.Contains(strings.Join(p.calls(), "\n"), "--head\tme:feat/x\t") {
		t.Errorf("pushurl apart from the url: exit %d, want --head me:feat/x\n%s", code, out)
	}
	gitIn(t, work, "config", "remote.origin.pushurl", "nowhere")
	gitIn(t, work, "config", "--add", "url."+bare+".insteadOf", "nowhere")
	if out, code := p.run(nil, "open", "up/r", "origin", "master", "t", body); code != 3 ||
		!strings.Contains(out, "which account") || strings.Contains(strings.Join(p.calls(), "\n"), "pr\tcreate") {
		t.Errorf("URL without an account: exit %d, want 3 and nothing opened\n%s", code, out)
	}
}

// gh ends with 4 when its sign-in has run out, and 4 is the script's
// 'reviewer not requested'. A sign-in lost after the check before writing is
// reported as what it is, exit 3, in every command (JG3REW review-gaps 4).
func TestPRRoundLostSignInIsNeverExit4(t *testing.T) {
	p := newPRRound(t)
	for _, args := range [][]string{
		{"threads", "o/r", "7"},
		{"settled", "o/r", "7"},
		{"request", "o/r", "7", "copilot"},
	} {
		if out, code := p.run([]string{"GH_LOST=1"}, args...); code != 3 || !strings.Contains(out, "gh auth login") {
			t.Errorf("%v with the sign-in gone: exit %d, want 3 naming gh auth login\n%s", args, code, out)
		}
	}
	// Lost between the post and the read-back: the answer went out, the run
	// still ends with 3, never as a finished 0.
	user := fixtureFile(t, "user.json", `{"data":{"viewer":{"login":"me"}}}`)
	issues := fixtureFile(t, "issues.json", `[{"id":1,"user":{"login":"copilot"},"body":"why?"}]`)
	posted := fixtureFile(t, "posted.json", `{"id":9,"html_url":"c9"}`)
	answer := fixtureFile(t, "answer.md", "Renamed it.\n")
	out, code := p.run([]string{"GH_LOST=after", "GH_STATE=" + user, "GH_ISSUES=" + issues, "GH_POSTED=" + posted}, "comment", "o/r", "7", answer)
	if code != 3 || !strings.Contains(out, "c9") || !strings.Contains(out, "gh auth login") || strings.Contains(out, "could not read") {
		t.Errorf("comment, sign-in lost after posting: exit %d, want 3 after printing c9\n%s", code, out)
	}
}

// Two sessions answering the same thread at once both pass the check before
// posting. Read back afterwards, the earliest answer of this login with this
// text stays; the later copy removes itself, whichever session it was
// (JG3REW, given back from signoff).
func TestPRRoundReplyLeavesOneAnswerWhenTwoSessionsRace(t *testing.T) {
	p := newPRRound(t)
	answer := fixtureFile(t, "answer.md", "Fixed in abc123.\r\n")
	open := fixtureFile(t, "open.json", `{"data":{"viewer":{"login":"me"},"node":{"isResolved":false,"comments":{"nodes":[{"author":{"login":"copilot"},"body":"rename"}]}}}}`)
	posted := fixtureFile(t, "posted.json", `{"data":{"addPullRequestReviewThreadReply":{"comment":{"id":"C2","url":"u2"}}}}`)
	for _, c := range []struct {
		name, after string
		deletes     bool
	}{
		{"the other session posted first", `[{"id":"C0","author":{"login":"copilot"},"body":"rename"},{"id":"C1","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, true},
		{"this session posted first", `[{"id":"C0","author":{"login":"copilot"},"body":"rename"},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123.\r\n"},{"id":"C3","author":{"login":"me"},"body":"Fixed in abc123."}]`, false},
		{"an earlier answer of other words", `[{"id":"C1","author":{"login":"me"},"body":"Will look at it."},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, false},
		{"the same words from someone else", `[{"id":"C1","author":{"login":"copilot"},"body":"Fixed in abc123."},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, false},
		// The agent writes under the person's login: a comment of theirs with
		// the same words is indistinguishable from another session's copy,
		// and is never what gets removed — only this run's own post is.
		{"the person wrote the same words before", `[{"id":"P1","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, true},
		{"the person wrote the same words after", `[{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"P3","author":{"login":"me"},"body":"Fixed in abc123."}]`, false},
		// The reviewer objected to an earlier answer of the same words: the
		// repeat is a new answer to the objection, not a copy (bughunt).
		{"the same answer again after an objection", `[{"id":"C1","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"R1","author":{"login":"copilot"},"body":"still wrong"},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, false},
		{"two copies after an objection", `[{"id":"C1","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"R1","author":{"login":"copilot"},"body":"still wrong"},{"id":"C0","author":{"login":"me"},"body":"Fixed in abc123."},{"id":"C2","author":{"login":"me"},"body":"Fixed in abc123."}]`, true},
	} {
		after := fixtureFile(t, "after.json", `{"data":{"viewer":{"login":"me"},"node":{"comments":{"nodes":`+c.after+`}}}}`)
		out, code := p.run([]string{"GH_STATE=" + open, "GH_FIXTURE=" + posted, "GH_AFTER=" + after}, "reply", "o/r", "7", "PRRT_x", answer)
		calls := strings.Join(p.calls(), "\n")
		if code != 0 || !strings.Contains(calls, "resolveReviewThread") {
			t.Errorf("%s: exit %d, want 0 and the thread resolved\n%s\n%s", c.name, code, out, calls)
		}
		deleted := strings.Contains(calls, "deletePullRequestReviewComment") && strings.Contains(calls, "id=C2")
		if deleted != c.deletes || (!c.deletes && strings.Contains(calls, "deletePullRequestReviewComment")) {
			t.Errorf("%s: removed this copy %v, want %v\n%s", c.name, deleted, c.deletes, calls)
		}
		for _, other := range []string{"id=C0", "id=C1", "id=C3", "id=P1", "id=P3"} {
			if strings.Contains(calls, other) {
				t.Errorf("%s: touched %s, which this run did not post\n%s", c.name, other, calls)
			}
		}
	}
}

// The conversation has nothing to resolve: the text is what keeps a second
// round from posting the same answer twice, and what keeps two sessions at
// once to one copy. Other words under the same login — the maintainer's own
// remark — never stop an answer. A forge that trims blanks at the ends of
// lines still matches (JG3REW, review-gaps closed in round 5).
func TestPRRoundCommentSaysItOnce(t *testing.T) {
	p := newPRRound(t)
	answer := fixtureFile(t, "answer.md", "Renamed it.  \r\nSee abc123.\r\n")
	user := fixtureFile(t, "user.json", `{"data":{"viewer":{"login":"me"}}}`)
	for _, c := range []struct {
		name, issues   string
		posts, deletes bool
	}{
		{"nothing said yet", `[{"id":1,"user":{"login":"copilot"},"body":"why?"}]`, true, false},
		{"said in an earlier round, trimmed by the forge", `[{"id":1,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."}]`, false, false},
		{"the maintainer's own remark", `[{"id":1,"user":{"login":"me"},"body":"please rename"}]`, true, false},
		{"the same words from someone else", `[{"id":1,"user":{"login":"copilot"},"body":"Renamed it.\nSee abc123."}]`, true, false},
		{"said before, then objected to", `[{"id":1,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."},{"id":2,"user":{"login":"copilot"},"body":"no"}]`, true, false},
	} {
		issues := fixtureFile(t, "issues.json", c.issues)
		out, code := p.run([]string{"GH_STATE=" + user, "GH_ISSUES=" + issues}, "comment", "o/r", "7", answer)
		calls := strings.Join(p.calls(), "\n")
		if code != 0 {
			t.Errorf("%s: exit %d\n%s", c.name, code, out)
		}
		if got := strings.Contains(calls, "--method\tPOST\trepos/o/r/issues/7/comments\t-F\tbody=@"+answer); got != c.posts {
			t.Errorf("%s: posted %v, want %v\n%s", c.name, got, c.posts, calls)
		}
		if strings.Contains(calls, "DELETE") {
			t.Errorf("%s: removed a comment\n%s", c.name, calls)
		}
	}

	// Two sessions at once: both found nothing before posting. Read back, the
	// earlier copy stays and the later one removes itself.
	posted := fixtureFile(t, "posted.json", `{"id":9,"html_url":"c9"}`)
	before := fixtureFile(t, "before.json", `[]`)
	for name, after := range map[string]struct {
		issues  string
		deletes bool
	}{
		"the other session posted first": {`[{"id":8,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."},{"id":9,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."}]`, true},
		"this session posted first":      {`[{"id":9,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."},{"id":10,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."}]`, false},
		"the person said it after":       {`[{"id":9,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."},{"id":11,"user":{"login":"me"},"body":"Renamed it.\nSee abc123."}]`, false},
	} {
		afterFile := fixtureFile(t, "after.json", after.issues)
		out, code := p.run([]string{"GH_STATE=" + user, "GH_ISSUES=" + before, "GH_ISSUES_AFTER=" + afterFile, "GH_POSTED=" + posted}, "comment", "o/r", "7", answer)
		calls := strings.Join(p.calls(), "\n")
		deleted := strings.Contains(calls, "DELETE\trepos/o/r/issues/comments/9")
		for _, call := range p.calls() {
			if strings.Contains(call, "DELETE") && strings.TrimRight(call, "\t") != "api\t--method\tDELETE\trepos/o/r/issues/comments/9" {
				t.Errorf("%s: removed a comment this run did not post: %q", name, call)
			}
		}
		if code != 0 || deleted != after.deletes || (!after.deletes && strings.Contains(calls, "DELETE")) {
			t.Errorf("%s: exit %d, removed this copy %v, want %v\n%s\n%s", name, code, deleted, after.deletes, out, calls)
		}
	}
}

// On a board whose review happens on the pull request, done waits until the
// round is over: no thread open, no reviewer still asked. Each thing in the
// way is named; a resolved thread and an empty request list are not
// (JG3REW, pr-when signoff).
func TestPRRoundSettledOnlyWhenNothingIsOpen(t *testing.T) {
	p := newPRRound(t)
	threads := func(resolved bool) string {
		return fixtureFile(t, "threads.json", `{"data":{"repository":{"pullRequest":{"reviewThreads":{
			"pageInfo":{"hasNextPage":false,"endCursor":null},
			"nodes":[{"id":"PRRT_1","isResolved":`+map[bool]string{true: "true", false: "false"}[resolved]+`,"isOutdated":false,"path":"dir with space/a.go","line":7,"comments":{"nodes":[{"author":{"login":"copilot"},"body":"x","url":"u1"}]}}]}}}}}`)
	}
	none := fixtureFile(t, "none.json", `{"users":[],"teams":[]}`)
	opinions := func(states ...string) string {
		var nodes []string
		for i, st := range states {
			nodes = append(nodes, `{"state":"`+st+`","author":{"login":"r`+string(rune('0'+i))+`"}}`)
		}
		return fixtureFile(t, "opinions.json", `{"data":{"repository":{"pullRequest":{"latestOpinionatedReviews":{"nodes":[`+strings.Join(nodes, ",")+`]}}}}}`)
	}
	approved := opinions("APPROVED", "COMMENTED")
	asked := fixtureFile(t, "asked.json", `{"users":[{"login":"Copilot"}],"teams":[{"slug":"core"}]}`)
	for _, c := range []struct {
		name, threads, requested, opinions string
		code                               int
		says                               []string
	}{
		{"everything answered", threads(true), none, approved, 0, []string{"settled"}},
		{"a thread open", threads(false), none, approved, 3, []string{"open thread: dir with space/a.go:7 u1"}},
		{"a review still asked", threads(true), asked, approved, 3, []string{"review pending: Copilot", "review pending: core"}},
		{"changes requested in a review body alone", threads(true), none, opinions("APPROVED", "CHANGES_REQUESTED"), 3, []string{"changes requested: r1"}},
		{"all of it", threads(false), asked, opinions("CHANGES_REQUESTED"), 3, []string{"open thread:", "review pending: Copilot", "changes requested: r0"}},
	} {
		out, code := p.run([]string{"GH_FIXTURE=" + c.threads, "GH_REQUESTED=" + c.requested, "GH_OPINIONS=" + c.opinions}, "settled", "o/r", "7")
		if code != c.code {
			t.Errorf("%s: exit %d, want %d\n%s", c.name, code, c.code, out)
		}
		for _, want := range c.says {
			if !strings.Contains(out, want) {
				t.Errorf("%s: does not say %q\n%s", c.name, want, out)
			}
		}
	}
	if out, code := p.run(nil, "settled", "o/r", "7", "extra"); code != 2 {
		t.Errorf("settled with an extra argument: exit %d, want 2\n%s", code, out)
	}
}
