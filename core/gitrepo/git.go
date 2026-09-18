// Package gitrepo is a thin wrapper over the system git binary.
//
// Shelling out to real git is deliberate. A reimplementation (go-git and
// friends) would produce diffs that differ in formatting and rename detection
// from what the reviewer sees when they run git themselves, and linking a C
// library would reintroduce the runtime dependency this tool exists without.
// Requiring git on PATH costs nothing: without a repository there is nothing for
// jaira to operate on.
package gitrepo

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Repo is a git working tree.
type Repo struct{ Dir string }

// ErrNoGit means the git binary is unavailable.
var ErrNoGit = errors.New("gitrepo: git is not available on PATH")

// Available reports whether git can be used at all.
func Available() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func (r *Repo) run(args ...string) (string, error) {
	if !Available() {
		return "", ErrNoGit
	}
	cmd := exec.Command("git", append([]string{"-C", r.Dir}, args...)...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return out.String(), nil
}

// IsRepo reports whether Dir is inside a git working tree.
func (r *Repo) IsRepo() bool {
	out, err := r.run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// Root returns the working tree root.
func (r *Repo) Root() (string, error) {
	out, err := r.run("rev-parse", "--show-toplevel")
	return strings.TrimSpace(out), err
}

// Commit describes one commit for display on a ticket.
type Commit struct {
	SHA     string `json:"sha"`
	Short   string `json:"short"`
	Subject string `json:"subject"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// Commits resolves metadata for the SHAs recorded on a ticket. Unknown SHAs are
// returned with empty metadata rather than dropped, so a ticket referencing a
// commit from an unfetched branch still shows what it claims.
func (r *Repo) Commits(shas []string) ([]Commit, error) {
	out := make([]Commit, 0, len(shas))
	for _, sha := range shas {
		if strings.TrimSpace(sha) == "" {
			continue
		}
		line, err := r.run("show", "-s", "--format=%H%x1f%h%x1f%s%x1f%an%x1f%aI", sha)
		if err != nil {
			out = append(out, Commit{SHA: sha, Short: shortSHA(sha)})
			continue
		}
		parts := strings.Split(strings.TrimSpace(line), "\x1f")
		if len(parts) < 5 {
			out = append(out, Commit{SHA: sha, Short: shortSHA(sha)})
			continue
		}
		out = append(out, Commit{SHA: parts[0], Short: parts[1], Subject: parts[2], Author: parts[3], Date: parts[4]})
	}
	return out, nil
}

// Diff returns the combined patch for a set of commits, scoped to those commits
// rather than to the working tree — a reviewer is judging what the ticket
// shipped, not whatever happens to be uncommitted right now.
func (r *Repo) Diff(shas []string) (string, error) {
	var b strings.Builder
	for _, sha := range shas {
		if strings.TrimSpace(sha) == "" {
			continue
		}
		out, err := r.run("show", "--patch", "--stat", "--format=commit %H%n%s%n", sha)
		if err != nil {
			b.WriteString(fmt.Sprintf("commit %s\n  (not available locally)\n\n", sha))
			continue
		}
		b.WriteString(out)
		b.WriteString("\n")
	}
	return b.String(), nil
}

// Stat returns the per-file summary for a set of commits.
func (r *Repo) Stat(shas []string) (string, error) {
	if len(shas) == 0 {
		return "", nil
	}
	var b strings.Builder
	for _, sha := range shas {
		out, err := r.run("show", "--stat=200", "--format=", sha)
		if err != nil {
			continue
		}
		b.WriteString(strings.TrimRight(out, "\n"))
		b.WriteString("\n")
	}
	return b.String(), nil
}

// HeadSHA returns the current commit.
func (r *Repo) HeadSHA() (string, error) {
	out, err := r.run("rev-parse", "HEAD")
	return strings.TrimSpace(out), err
}

func shortSHA(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

// runTolerating runs git and accepts one non-zero exit status as a normal
// result. `git diff --no-index` reports "the files differ" as exit 1, which is
// the whole point of calling it, so run() — which treats any non-zero status as
// a failure — would turn every untracked file into an error.
func (r *Repo) runTolerating(code int, args ...string) (string, error) {
	if !Available() {
		return "", ErrNoGit
	}
	cmd := exec.Command("git", append([]string{"-C", r.Dir}, args...)...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	if err == nil {
		return out.String(), nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.ExitCode() == code {
		return out.String(), nil
	}
	msg := strings.TrimSpace(errb.String())
	if msg == "" {
		msg = err.Error()
	}
	return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
}

// WorktreeDiff returns the patch for work that is not committed yet: every
// tracked change against HEAD, staged or not, followed by every untracked file
// in full.
//
// A lane judges what is in front of it, and on a conversational ticket — where
// the worker hands back a commit line instead of committing — that is the
// worktree and not the commits. The rule "a lane that changed no code commits
// nothing" leaves work lying there on an autonomous board too, so this is not
// a conversational-mode special case. Telling the reviewer to go and look for
// it by hand is the same hand instruction this file's Diff comment exists to
// have removed.
//
// .jaira/tickets is excluded on purpose: the worker rewrites the ticket file
// with every 'jaira dod' and every 'jaira note', so including it would bury
// the few lines of code under the ticket's own prose — prose the payload
// already carries as goal, definition-of-done and notes. ",top" anchors the
// exclusion at the repository root so it holds whatever directory Dir is.
func (r *Repo) WorktreeDiff() (string, error) {
	const notTickets = ":(exclude,top).jaira/tickets"
	var b strings.Builder
	tracked, err := r.run("diff", "HEAD", "--patch", "--stat", "--", ":/", notTickets)
	if err != nil {
		return "", err
	}
	b.WriteString(tracked)
	others, err := r.run("ls-files", "--others", "--exclude-standard", "--", ":/", notTickets)
	if err != nil {
		return "", err
	}
	for _, path := range strings.Split(strings.TrimSpace(others), "\n") {
		if strings.TrimSpace(path) == "" {
			continue
		}
		// git diff alone is blind to a file the index has never seen, and a
		// new test or a new package is the common shape of a change, not a
		// corner case. --no-index is what sees it without 'git add -N', which
		// would write to an index another session may be holding.
		out, err := r.runTolerating(1, "diff", "--no-index", "--", "/dev/null", path)
		if err != nil {
			continue
		}
		b.WriteString(out)
	}
	return b.String(), nil
}
