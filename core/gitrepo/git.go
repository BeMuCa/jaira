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
