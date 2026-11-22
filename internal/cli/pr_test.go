package cli

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

type prExit struct {
	code int
	msg  string
}

func (e prExit) Error() string { return e.msg }
func (e prExit) ExitCode() int { return e.code }

// fakeForge stands in for git and gh: the branch checked out, and one
// GraphQL answer per read in turn — the last repeats. An answer starting
// with "!" fails with that text.
func fakeForge(t *testing.T, branch string, answers ...string) *[]string {
	t.Helper()
	var calls []string
	reads := 0
	oldRun, oldLook, oldSleep := prRun, prLook, prSleep
	t.Cleanup(func() { prRun, prLook, prSleep = oldRun, oldLook, oldSleep })
	prLook = func(string) (string, error) { return "/usr/bin/gh", nil }
	prSleep = func(time.Duration) {}
	prRun = func(_ string, name string, args ...string) ([]byte, error) {
		line := name + " " + strings.Join(args, " ")
		calls = append(calls, line)
		switch {
		case strings.HasPrefix(line, "git remote get-url origin"):
			return []byte("https://github.com/me/jaira.git\n"), nil
		case strings.HasPrefix(line, "git config"):
			return nil, prExit{1, ""}
		case strings.HasPrefix(line, "git symbolic-ref"):
			if branch == "" {
				return nil, prExit{1, "fatal: ref HEAD is not a symbolic ref"}
			}
			return []byte(branch + "\n"), nil
		case strings.HasPrefix(line, "gh repo view"):
			return []byte(`{"nameWithOwner":"me/jaira"}`), nil
		case strings.HasPrefix(line, "gh api graphql"):
			a := answers[min(reads, len(answers)-1)]
			reads++
			if strings.HasPrefix(a, "!") {
				return nil, prExit{1, a[1:]}
			}
			return []byte(a), nil
		}
		return nil, errors.New("unexpected: " + line)
	}
	return &calls
}

// prNode is one pull request in a GraphQL answer.
func prNode(state, ci, lastComment string) string {
	comments := ""
	if lastComment != "" {
		comments = fmt.Sprintf(`{"createdAt":%q,"author":{"login":"bob"}}`, lastComment)
	}
	return fmt.Sprintf(`{"number":7,"title":"t","url":"https://x/7","state":%q,"headRefName":"feat/craszs",
	 "comments":{"nodes":[%s]},"head":{"nodes":[{"commit":{"statusCheckRollup":{"state":%q}}}]}}`, state, comments, ci)
}

func listAnswer(nodes ...string) string {
	return `{"data":{"viewer":{"login":"me"},"repository":{"pullRequests":{"nodes":[` + strings.Join(nodes, ",") + `]}}}}`
}

func oneAnswer(node string) string {
	return `{"data":{"viewer":{"login":"me"},"repository":{"pullRequest":` + node + `}}}`
}

func TestPRShowsTheCheckedOutBranch(t *testing.T) {
	calls := fakeForge(t, "feat/craszs", listAnswer(prNode("OPEN", "FAILURE", "2026-11-20T10:00:00Z")))
	out, err := runCLI(t, t.TempDir(), "pr")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	for _, want := range []string{"#7 (CRASZS) feat/craszs", "CI: failure", "comments from others: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(strings.Join(*calls, "\n"), "head=feat/craszs") {
		t.Errorf("the read was not for the checked-out branch:\n%s", strings.Join(*calls, "\n"))
	}
}

func TestPRWatchReportsEachChangeOnceAndEnds(t *testing.T) {
	open := prNode("OPEN", "PENDING", "2026-11-20T10:00:00Z")
	fakeForge(t, "feat/craszs",
		listAnswer(open),
		oneAnswer(open), // nothing new
		"!HTTP 502: Bad Gateway",
		"!HTTP 502: Bad Gateway", // said once
		oneAnswer(prNode("OPEN", "FAILURE", "2026-11-20T11:00:00Z")),
		oneAnswer(prNode("OPEN", "FAILURE", "2026-11-20T11:00:00Z")), // the same again: nothing
		oneAnswer(prNode("MERGED", "FAILURE", "2026-11-20T11:00:00Z")),
	)
	out, err := runCLI(t, t.TempDir(), "pr", "--watch")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	want := []string{
		"watching #7 (CRASZS) https://x/7: threads: 0 waiting, 0 answered, 0 resolved · comments from others: 1 · CI: pending · settled",
		"#7 (CRASZS): read failed, trying again in 1m0s: gh: HTTP 502: Bad Gateway",
		"#7 (CRASZS): read again",
		"#7 (CRASZS): new comment by bob",
		"#7 (CRASZS): CI red",
		"#7 (CRASZS): merged",
	}
	if got := strings.Split(strings.TrimSpace(out), "\n"); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("lines:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestPRExitCodes(t *testing.T) {
	two := listAnswer(prNode("OPEN", "SUCCESS", ""), strings.Replace(prNode("OPEN", "SUCCESS", ""), `"number":7`, `"number":8`, 1))
	cases := []struct {
		name   string
		branch string
		answer string
		args   []string
		code   int
		says   string
	}{
		{"no pull request shows and exits 0", "feat/x", listAnswer(), []string{"pr"}, 0, "No open pull request for feat/x"},
		{"no pull request to watch", "feat/x", listAnswer(), []string{"pr", "--watch"}, ExitNotFound, "nothing to watch"},
		{"several to watch", "feat/x", two, []string{"pr", "--watch"}, ExitUsage, "#7, #8"},
		{"detached head", "", listAnswer(), []string{"pr"}, ExitUsage, "no branch is checked out"},
		{"signed out", "feat/x", "!To get started with GitHub CLI, please run:  gh auth login", []string{"pr"}, ExitValidation, "gh auth login"},
		{"forge error", "feat/x", "!HTTP 502", []string{"pr"}, ExitError, "HTTP 502"},
		{"watch with json", "feat/x", listAnswer(), []string{"pr", "--watch", "--json"}, ExitUsage, "not JSON"},
		{"every too short", "feat/x", listAnswer(), []string{"pr", "--watch", "--every", "1s"}, ExitUsage, "too short"},
		{"all with a branch", "feat/x", listAnswer(), []string{"pr", "--all", "feat/x"}, ExitUsage, "--all"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeForge(t, c.branch, c.answer)
			out, err := runCLI(t, t.TempDir(), c.args...)
			if got := exitCode(err); got != c.code {
				t.Fatalf("exit %d, want %d: %v\n%s", got, c.code, err, out)
			}
			text := out
			if err != nil {
				text += err.Error()
			}
			if !strings.Contains(text, c.says) {
				t.Fatalf("says %q, want %q", text, c.says)
			}
		})
	}
}
