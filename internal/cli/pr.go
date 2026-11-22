package cli

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/pr"
)

// prRun and prLook are how 'jaira pr' reaches git, gh and glab; tests swap
// them for fakes. prSleep is the watch's pause between two reads.
var (
	prRun pr.Run = func(dir, name string, args ...string) ([]byte, error) {
		c := exec.Command(name, args...)
		c.Dir = dir
		return c.Output()
	}
	prLook  = exec.LookPath
	prSleep = time.Sleep
)

// minEvery keeps a watch from asking the forge more often than its rate
// limit can carry: a GraphQL read of one pull request costs a point, and
// gh's budget is 5000 an hour.
const minEvery = 10 * time.Second

func newPRCmd() *cobra.Command {
	var all, watch bool
	var number int
	var every time.Duration
	cmd := &cobra.Command{
		Use:   "pr [<branch>]",
		Short: "Show where this branch's pull request stands, or watch it",
		Long: `Shows the open pull request of <branch> — the checked-out branch when none
is named — as the forge has it now: review threads waiting for an answer,
answered but not resolved, and resolved; how many comments others wrote; the
CI state of its head commit; reviews still requested, changes requested and
approvals; the tickets its commits name; and whether its review round is
settled (no thread open, no review pending, no changes requested — what
jaira-role-pr's 'pr-round.sh settled' asks). --all shows every open pull
request of the repository, --number one pull request in any state. Every
pull request and every review thread is read, page after page; comments
from others are counted over the newest 50 on the conversation and the
newest 20 of each thread. A branch is the one in origin: a fork's branch of
the same name is somebody else's pull request.

Read through gh on GitHub and glab on GitLab — the host of origin decides, or
git config jaira.forge. In a fork the pull requests are the parent's. jaira
stores nothing: every call asks the forge.

--watch follows one pull request and prints one line per change, made to
run under a Monitor: a new comment, review comment or review from someone
else, CI turning red or green again, the round becoming settled, and merged
or closed — after which it ends with exit 0. Its first line is the stand it
starts from; nothing in it is reported as a change, so a watch restarted
does not report again what an earlier one did. A read that fails is said
once and tried again at the next interval; a lost sign-in ends the watch.

Exit codes: 0 read (also when there is no pull request to show); 2 usage,
or --watch on a branch with more than one open pull request (name it with
--number); 3 no forge (origin is neither GitHub nor GitLab), gh or glab
missing, or not signed in; 5 --watch found no open pull request; 1 the forge
answered with an error.`,
		Args: maxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			w := cmd.OutOrStdout()
			q := pr.Query{Dir: g.dir, Number: number}
			if len(args) == 1 {
				q.Branch = args[0]
			}
			switch {
			case all && (q.Branch != "" || number > 0):
				return fail(ExitUsage, "usage", "pr: --all reads every open pull request — leave the branch and --number out")
			case q.Branch != "" && number > 0:
				return fail(ExitUsage, "usage", "pr: name a branch or --number, not both")
			case watch && all:
				return fail(ExitUsage, "usage", "pr: --watch follows one pull request — leave --all out")
			case watch && g.jsonOut:
				return fail(ExitUsage, "usage", "pr: --watch prints lines, not JSON")
			case every < minEvery:
				return fail(ExitUsage, "usage", "pr: --every %s is too short — at least %s", every, minEvery)
			}
			if !all && q.Branch == "" && number == 0 {
				b, err := currentBranch(g.dir)
				if err != nil {
					return err
				}
				q.Branch = b
			}
			st, err := pr.Read(prRun, prLook, q)
			if err != nil {
				return prFailure(err)
			}
			if watch {
				return watchPR(w, q, st, every)
			}
			if g.jsonOut {
				return emit(w, st)
			}
			printPRs(w, st)
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "every open pull request of the repository")
	cmd.Flags().IntVar(&number, "number", 0, "one pull request by its number, open or not")
	cmd.Flags().BoolVar(&watch, "watch", false, "follow the pull request and print one line per change until it is merged or closed")
	cmd.Flags().DurationVar(&every, "every", time.Minute, "with --watch: how long to wait between two reads")
	return cmd
}

func currentBranch(dir string) (string, error) {
	out, err := prRun(dir, "git", "symbolic-ref", "--quiet", "--short", "HEAD")
	b := strings.TrimSpace(string(out))
	if err != nil || b == "" {
		return "", fail(ExitUsage, "no_branch", "pr: no branch is checked out here — name one, or use --all or --number")
	}
	return b, nil
}

func prFailure(err error) error {
	switch {
	case errors.Is(err, pr.ErrNoForge):
		return fail(ExitValidation, "no_forge", "pr: %v", err)
	case errors.Is(err, pr.ErrNoTool):
		return fail(ExitValidation, "no_tool", "pr: %v", err)
	case errors.Is(err, pr.ErrAuth):
		return fail(ExitValidation, "not_signed_in", "pr: %v", err)
	}
	return fail(ExitError, "forge_failed", "pr: %v", err)
}

func printPRs(w io.Writer, st *pr.Status) {
	if len(st.PRs) == 0 {
		if st.Branch != "" {
			fmt.Fprintf(w, "No open pull request for %s in %s.\n", st.Branch, st.Repo)
		} else {
			fmt.Fprintf(w, "No open pull request in %s.\n", st.Repo)
		}
		return
	}
	for i, p := range st.PRs {
		if i > 0 {
			fmt.Fprintln(w)
		}
		draft := ""
		if p.Draft {
			draft = " [draft]"
		}
		state := ""
		if p.State != "open" {
			state = " [" + p.State + "]"
		}
		fmt.Fprintf(w, "%s %s%s%s — %s\n  %s\n  %s\n", pr.Label(p), p.Branch, draft, state, p.Title, p.URL, pr.Summary(p))
		if len(p.Unsettled) > 0 {
			fmt.Fprintf(w, "  in the way: %s\n", strings.Join(p.Unsettled, "; "))
		}
	}
}

// watchPR follows the one pull request st found until it is merged or
// closed. It reads it by number from then on, so a pull request that stops
// being open is still found.
func watchPR(w io.Writer, q pr.Query, st *pr.Status, every time.Duration) error {
	switch len(st.PRs) {
	case 0:
		return fail(ExitNotFound, "no_pull_request", "pr: no open pull request for %s in %s — nothing to watch", q.Branch, st.Repo)
	case 1:
	default:
		var nums []string
		for _, p := range st.PRs {
			nums = append(nums, fmt.Sprintf("#%d", p.Number))
		}
		return fail(ExitUsage, "several_pull_requests", "pr: %s has %d open pull requests (%s) — watch one with --number", q.Branch, len(st.PRs), strings.Join(nums, ", "))
	}
	last := st.PRs[0]
	fmt.Fprintf(w, "watching %s %s: %s\n", pr.Label(last), last.URL, pr.Summary(last))
	if pr.Over(last) {
		fmt.Fprintf(w, "%s: %s\n", pr.Label(last), last.State)
		return nil
	}
	q = pr.Query{Dir: q.Dir, Number: last.Number}
	failed := ""
	for {
		prSleep(every)
		next, err := pr.Read(prRun, prLook, q)
		if err != nil {
			if errors.Is(err, pr.ErrAuth) || errors.Is(err, pr.ErrNoTool) || errors.Is(err, pr.ErrNoForge) {
				return prFailure(err)
			}
			// Said once, not every interval: a forge that is down for an
			// hour would otherwise fill the watcher's screen with one line.
			if msg := err.Error(); msg != failed {
				fmt.Fprintf(w, "%s: read failed, trying again in %s: %s\n", pr.Label(last), every, msg)
				failed = msg
			}
			continue
		}
		if failed != "" {
			fmt.Fprintf(w, "%s: read again\n", pr.Label(last))
			failed = ""
		}
		if len(next.PRs) != 1 {
			continue
		}
		cur := next.PRs[0]
		for _, line := range pr.Changes(last, cur) {
			fmt.Fprintln(w, line)
		}
		last = pr.Carry(last, cur)
		if pr.Over(cur) {
			return nil
		}
	}
}
