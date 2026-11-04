package gitrepo

import "strings"

// BranchSegments splits a branch name into its words: runs of letters and
// digits, so "feat/eh6z0n-diff-view" gives "feat", "eh6z0n", "diff", "view".
func BranchSegments(branch string) []string {
	return strings.FieldsFunc(branch, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
}

// BranchFor returns the branch whose name carries handle as one of its
// segments ("feat/eh6z0n-diff-view" carries "EH6Z0N"), using the same rule as
// branchTicket in internal/cli. When several do, the one with the newest commit
// wins: an old experiment under the same handle must not shadow the work. ""
// when none does.
func (r *Repo) BranchFor(handle string) string {
	out, err := r.run("for-each-ref", "--sort=-committerdate", "--format=%(refname:short)", "refs/heads", "refs/remotes")
	if err != nil {
		return ""
	}
	for _, name := range strings.Split(out, "\n") {
		name = strings.TrimSpace(name)
		if name == "" || strings.HasSuffix(name, "/HEAD") {
			continue
		}
		for _, seg := range BranchSegments(name) {
			if strings.EqualFold(seg, handle) {
				return name
			}
		}
	}
	return ""
}

// Base returns the branch a ticket branch is reviewed against: master, else
// main, else the same two on origin. "" when none exists.
func (r *Repo) Base() string {
	for _, ref := range []string{"refs/heads/master", "refs/heads/main", "refs/remotes/origin/master", "refs/remotes/origin/main"} {
		if _, err := r.run("rev-parse", "--verify", "--quiet", ref); err == nil {
			return strings.TrimPrefix(strings.TrimPrefix(ref, "refs/heads/"), "refs/remotes/")
		}
	}
	return ""
}

// BranchDiff returns the patch of branch against the merge base with base,
// which is what GitHub shows for a pull request: changes made on the branch,
// not changes made on base since it forked. The .jaira directory is left out,
// as the ticket files there are bookkeeping, not code under review. The
// --src-prefix and --dst-prefix flags pin the a/ and b/ prefixes the TUI
// parser reads, whatever diff.noprefix says.
func (r *Repo) BranchDiff(base, branch string) (string, error) {
	return r.run("-c", "core.quotePath=false", "diff", "--patch",
		"--src-prefix=a/", "--dst-prefix=b/",
		base+"..."+branch, "--", ":/", ":(exclude,top).jaira")
}

// WorktreeOf returns the path of the worktree that has branch checked out
// (short name, e.g. feat/x), or "" when no worktree has it.
func (r *Repo) WorktreeOf(branch string) string {
	out, err := r.run("worktree", "list", "--porcelain")
	if err != nil {
		return ""
	}
	want := "branch refs/heads/" + branch
	path := ""
	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			path = strings.TrimPrefix(line, "worktree ")
		case line == want:
			return path
		}
	}
	return ""
}
