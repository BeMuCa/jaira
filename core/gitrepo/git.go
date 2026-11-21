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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
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

// run treats any non-zero exit status as a failure. No exit status is
// negative, so noTolerance tolerates none of them.
func (r *Repo) run(args ...string) (string, error) {
	return r.runTolerating(noTolerance, args...)
}

const noTolerance = -1

// IsRepo reports whether Dir is inside a git working tree.
func (r *Repo) IsRepo() bool {
	out, err := r.run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// Tracks reports whether git tracks at least one file under path, relative to
// Dir. Outside a repository, or without git, nothing is tracked.
func (r *Repo) Tracks(path string) bool {
	out, err := r.run("ls-files", "--", path)
	return err == nil && strings.TrimSpace(out) != ""
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
		if !IsSHA(sha) {
			out = append(out, Commit{SHA: sha, Short: shortSHA(sha)})
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

// codeOnly is the pathspec that keeps ticket files out of Diff and WorktreeDiff:
// the whole repository, minus every place a ticket file lives — on the board,
// filed in the logbook, rejected into the archive. A filed ticket carries the
// same outcome and notes as one on the board, and a commit that filed it can
// be in a ticket's commit list by the file's history alone. ",top" anchors
// each exclusion at the repository root so it holds whatever directory Dir is.
var codeOnly = []string{
	":/",
	":(exclude,top).jaira/tickets",
	":(exclude,top).jaira/logbook",
	":(exclude,top).jaira/archive",
}

// pathspec returns args, then "--", then codeOnly.
func pathspec(args ...string) []string {
	return append(append(args, "--"), codeOnly...)
}

// Diff returns the combined patch for a set of commits, scoped to those commits
// rather than to the working tree — a reviewer is judging what the ticket
// shipped, not whatever happens to be uncommitted right now.
//
// The second return names the shas git could not show: rebased away,
// cherry-picked elsewhere, or on a branch this clone never fetched. They stay
// in the patch as a placeholder line, but a caller counting commits would
// otherwise report them as diffs it had shown — the same fraction-reported-whole
// failure this package was fixed for. There is no error return because there is
// no failure: an unshowable sha is an outcome, not a broken repository.
//
// Ticket files (codeOnly) are left out for the reason WorktreeDiff leaves them
// out, and one
// more: the ticket file rides in the same commits as the code, so its patch
// carries the implementer's outcome and notes — and a lane that declares it
// does not take them (review reads goal, definition-of-done and diff only, so
// it judges the code without the author's account of it) would read them in
// the diff all the same. A commit that touched nothing but its ticket file
// therefore leaves no text here at all — it carried no code to judge — while
// it is still counted in the commits the caller names.
func (r *Repo) Diff(shas []string) (string, []string) {
	var b strings.Builder
	var unavailable []string
	for _, sha := range shas {
		if strings.TrimSpace(sha) == "" {
			continue
		}
		if !IsSHA(sha) {
			b.WriteString(fmt.Sprintf("commit %s\n  (not a commit id)\n\n", sha))
			unavailable = append(unavailable, sha)
			continue
		}
		out, err := r.run(pathspec("show", "--patch", "--stat", "--format=commit %H%n%s%n", sha)...)
		if err != nil {
			b.WriteString(fmt.Sprintf("commit %s\n  (not available locally)\n\n", sha))
			unavailable = append(unavailable, sha)
			continue
		}
		if out == "" {
			// Only the ticket file changed; see above.
			continue
		}
		b.WriteString(out)
		b.WriteString("\n")
	}
	return b.String(), unavailable
}

// Stat returns the per-file summary for a set of commits.
func (r *Repo) Stat(shas []string) (string, error) {
	if len(shas) == 0 {
		return "", nil
	}
	var b strings.Builder
	for _, sha := range shas {
		if !IsSHA(sha) {
			continue
		}
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

// Branch returns the checked-out branch's short name, or "" on a detached
// HEAD or outside a repository.
func (r *Repo) Branch() string {
	out, err := r.run("symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// IsSHA reports whether s is a commit id: 4 to 64 hex digits. The shas these
// functions take come from a ticket's hand-editable commits: field, and a git
// argument is an option when it starts with a dash — "--output=<path>" in that
// field would make 'git show' write over a file. Only an id is handed to git.
func IsSHA(s string) bool {
	if len(s) < 4 || len(s) > 64 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

func shortSHA(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

// runTolerating runs git and accepts one non-zero exit status as a normal
// result. `git diff --no-index` reports "the files differ" as exit 1, which is
// the whole point of calling it, so run — which tolerates nothing — would turn
// every untracked file into an error. Every call in this package goes through
// here; run is this with nothing tolerated.
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
	// Exited() is not belt and braces: a git killed by a signal has no exit
	// status and ExitCode reports that as -1, which would otherwise match a
	// caller that tolerates nothing and hand back a truncated result as a
	// success.
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.Exited() && ee.ExitCode() == code {
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
	var b strings.Builder
	tracked, err := r.run(pathspec("-c", "core.quotePath=false", "diff", "HEAD", "--patch", "--stat")...)
	if err != nil {
		return "", err
	}
	b.WriteString(tracked)
	// -z is not a convenience, and neither is the core.quotePath=false above
	// it and on the call below. Without it git quotes any path that is not
	// plain ASCII — "\303\204nderung.txt" for a file with an umlaut — and the
	// --no-index call below then looks for a file of that literal name, does
	// not find it, and says so with exit 1, the very code tolerated here as
	// "the files differ". The file would drop out of the patch in silence
	// while the payload still claims to carry the worktree: the fragment
	// reported as the whole that this package exists to have ended. -z turns
	// the quoting off on the way in and settles paths with spaces in them at
	// the same time; core.quotePath=false turns it off on the way out, so the
	// patch a reviewer reads names the file rather than its escape sequence.
	others, err := r.run(pathspec("ls-files", "--others", "--exclude-standard", "-z")...)
	if err != nil {
		return "", err
	}
	for _, path := range strings.Split(others, "\x00") {
		if path == "" {
			continue
		}
		// git diff alone is blind to a file the index has never seen, and a
		// new test or a new package is the common shape of a change, not a
		// corner case. --no-index is what sees it without 'git add -N', which
		// would write to an index another session may be holding.
		out, err := r.runTolerating(1, "-c", "core.quotePath=false", "diff", "--no-index", "--", "/dev/null", path)
		if err != nil {
			return "", err
		}
		// Exit 1 carries two meanings here — "the files differ", which is the
		// expected one, and "the path could not be read". Only the second
		// produces no patch, and a file that exists with bytes in it always
		// produces one, so an empty patch beside a non-empty file is that
		// second meaning wearing the first one's exit code. Reporting it is
		// what keeps the caller from labelling a worktree it never saw.
		if out == "" {
			if info, statErr := os.Stat(filepath.Join(r.Dir, path)); statErr != nil || info.Size() > 0 {
				return "", fmt.Errorf("git diff --no-index -- /dev/null %s: no patch for an unreadable or non-empty file", path)
			}
		}
		b.WriteString(out)
	}
	return b.String(), nil
}

// HuntState is what a bug hunt at one commit says about the code now.
type HuntState int

const (
	// HuntCovers: nothing the ticket committed since touched what it read,
	// and nothing it read is changed and not committed.
	HuntCovers HuntState = iota
	// HuntStale: a commit of the ticket since then touched what it read, or
	// this checkout does not descend from the hunt's commit.
	HuntStale
	// HuntUnknown: no commit since, but a file it read is changed and not
	// committed. Whether the hunt saw that change git cannot tell — in
	// conversational mode the worker commits after the lanes, and the hunt
	// read the change as it was — so this is neither covered nor stale.
	HuntUnknown
)

// HuntCoverage compares commits only, never clocks: a hunt is stale when a
// commit since then that names the ticket and is no merge touched the paths
// it read. Every commit of a ticket names it (the house rule its commit list
// is built from), so master merged into the branch — work other people
// reviewed — does not count, on the branch or after it was merged; what a
// merge's conflict resolution wrote does not either. The limit of that rule:
// a commit of this ticket that does not name it is not seen. A commit that is
// no ancestor of HEAD — the hunt ran on a branch this checkout does not have
// — reads as stale. Board data — everything under .jaira/ but the lane
// prompts — never counts: a tag, a milestone or a pulled ticket changes no
// code. Paths are relative to the repository root; none, or one that is
// neither in commit nor on disk (a file list split at a space inside a
// name), means all code (codeOnly).
func (r *Repo) HuntCoverage(commit, ticketID string, paths []string) (HuntState, error) {
	if !ticket.ValidID(ticketID) {
		return HuntStale, fmt.Errorf("gitrepo: %q is no ticket id", ticketID)
	}
	// A commit this repository does not have — rebased away and
	// force-pushed, seen from a fresh clone — is as far from HEAD as one
	// that is no ancestor: the hunt read code this checkout does not hold.
	if _, err := r.run("rev-parse", "--verify", "--quiet", commit+"^{commit}"); err != nil {
		return HuntStale, nil
	}
	if _, err := r.run("merge-base", "--is-ancestor", commit, "HEAD"); err != nil {
		return HuntStale, nil
	}
	specs, err := r.huntedSpecs(commit, paths)
	if err != nil {
		return HuntStale, err
	}
	own, err := r.run(append([]string{"-c", "core.quotePath=false", "log", "--no-merges", "--format=", "--name-only", "-z", "--extended-regexp", "--regexp-ignore-case", "--grep=" + namesTicket(ticketID), commit + "..HEAD", "--"}, specs...)...)
	if err != nil {
		return HuntStale, err
	}
	if anyCode(own) {
		return HuntStale, nil
	}
	changed, err := r.run(append([]string{"-c", "core.quotePath=false", "diff", "--name-only", "-z", "HEAD", "--"}, specs...)...)
	if err != nil {
		return HuntStale, err
	}
	others, err := r.run(append([]string{"-c", "core.quotePath=false", "ls-files", "--others", "--exclude-standard", "--full-name", "-z", "--"}, specs...)...)
	if err != nil {
		return HuntStale, err
	}
	if anyCode(changed) || anyCode(others) {
		return HuntUnknown, nil
	}
	return HuntCovers, nil
}

// anyCode reports a path in git's NUL-separated list that is no board data.
func anyCode(list string) bool {
	for _, p := range strings.Split(list, "\x00") {
		if p = strings.TrimSpace(p); p != "" && !isBoardData(p) {
			return true
		}
	}
	return false
}

// isBoardData reports a path under .jaira/ that is no lane prompt: tickets,
// tags, milestones, settings — what the board writes, not code. Lane prompts
// are today's board's in .jaira/lanes/ and a named board's in
// .jaira/boards/<name>/lanes/ (core/lane boardLanesDir).
func isBoardData(p string) bool {
	rest, ok := strings.CutPrefix(p, ".jaira/")
	if !ok {
		return false
	}
	if strings.HasPrefix(rest, "lanes/") {
		return false
	}
	if board, ok := strings.CutPrefix(rest, "boards/"); ok {
		if _, file, ok := strings.Cut(board, "/"); ok && strings.HasPrefix(file, "lanes/") {
			return false
		}
	}
	return true
}

// huntedSpecs turns a hunt's file list into pathspecs: each one literal, or
// all code when the list is empty or names a path that is neither in commit
// nor on disk — a list a reader cannot trust is not narrowed by.
func (r *Repo) huntedSpecs(commit string, paths []string) ([]string, error) {
	if len(paths) == 0 {
		return codeOnly, nil
	}
	root, err := r.Root()
	if err != nil {
		return nil, err
	}
	listed, err := r.run(append([]string{"-c", "core.quotePath=false", "ls-tree", "-r", "--name-only", "--full-tree", "-z", commit, "--"}, paths...)...)
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, p := range strings.Split(listed, "\x00") {
		known[p] = true
	}
	specs := make([]string, 0, len(paths))
	for _, p := range paths {
		if !known[p] {
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(p))); err != nil {
				return codeOnly, nil
			}
		}
		specs = append(specs, ":(top,literal)"+p)
	}
	return specs, nil
}
