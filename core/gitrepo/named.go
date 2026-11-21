package gitrepo

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrNoSuchRef is returned when the base a branch is counted against does not
// name a commit here.
var ErrNoSuchRef = errors.New("no such commit")

// subjectScope is the scope of a commit subject the way this board's commits
// are written: type(A3K9QP): ..., or type(A3K9QP,B7XQ2M)!: ... for several.
// The subject git itself writes for 'git revert' — Revert "type(A3K9QP): ..."
// — and for reverting that again — Reapply "...", or Revert "Revert "..."" in
// a git before 2.43 — names the same tickets: the revert is work on the ticket
// it undoes, and under .jaira/pr none an unnamed commit holds the gate.
var subjectScope = regexp.MustCompile(`^(?:(?:Revert|Reapply) ")*[a-z]+\(([^)]*)\)!?:`)

var handleShape = regexp.MustCompile(`^[0-9A-Z]{6}$`)

// TicketsNamed returns every ticket handle the subjects of base..HEAD name in
// their scope, once each, in the order git lists them. Only the subject is
// read: a merge commit quotes the subjects it brought in in its body, and
// would name every ticket it merged. It is the same reading as ticketsNamed
// in the cockpit's board.ts, which decides when an Accept opens a pull
// request.
func (r *Repo) TicketsNamed(base string) ([]string, error) {
	if !Available() {
		return nil, ErrNoGit
	}
	// A base starting with a dash would be read as an option of git log.
	if base == "" || strings.HasPrefix(base, "-") {
		return nil, ErrNoSuchRef
	}
	if _, err := r.run("rev-parse", "--verify", "--quiet", base+"^{commit}"); err != nil {
		return nil, ErrNoSuchRef
	}
	out, err := r.run("log", "--format=%s", base+"..HEAD", "--")
	if err != nil {
		return nil, err
	}
	return HandlesInSubjects(out), nil
}

// UnnamedCommits returns the commits of base..HEAD whose subject names no
// ticket, as "<short sha> <subject>" in the order git lists them. A revert
// with git's own subject names a ticket only where it takes back a commit of
// the branch (revertsBranchWork); a revert of the base's work counts here. A merge
// that names none passes only when its tree is exactly what git's own merge
// of its two parents gives ('git merge-tree --write-tree'): such a merge
// brings in nothing not listed on its own — a branch merged in shows its
// commits here, the base merged into the branch nothing, since the base is
// excluded. Anything else counts: a conflict resolution, a file added during
// 'merge --no-commit', and a merge that drops one side ('-s ours', a
// conflict settled with 'checkout --ours'), which 'git diff-tree --cc' hides
// because the result matches a parent. An octopus merge, or a git too old to
// answer (before 2.38), counts too: the gate fails closed.
// base must already be known to name a commit (TicketsNamed checks it).
func (r *Repo) UnnamedCommits(base string) ([]string, error) {
	if base == "" || strings.HasPrefix(base, "-") {
		return nil, ErrNoSuchRef
	}
	out, err := r.run("log", "--format=%H%x00%h%x00%P%x00%s", base+"..HEAD", "--")
	if err != nil {
		return nil, err
	}
	var commits [][]string
	onBranch := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		if f := strings.SplitN(strings.TrimSpace(line), "\x00", 4); len(f) == 4 {
			commits = append(commits, f)
			onBranch[f[0]] = true
		}
	}
	var unnamed []string
	for _, f := range commits {
		if len(HandlesInSubjects(f[3])) > 0 && (!revertSubject.MatchString(f[3]) || r.revertsBranchWork(f[0], onBranch)) {
			continue
		}
		if parents := strings.Fields(f[2]); len(parents) > 1 && r.isPlainMerge(f[0], parents) {
			continue
		}
		unnamed = append(unnamed, f[1]+" "+f[3])
	}
	return unnamed, nil
}

// revertSubject is the subject git writes for 'git revert' and for reverting
// a revert again, as subjectScope reads it.
var revertSubject = regexp.MustCompile(`^(?:Revert|Reapply) "`)

// revertedCommit is the line 'git revert' writes into the body.
var revertedCommit = regexp.MustCompile(`(?m)^This reverts commit ([0-9a-f]{40})\b`)

// revertsBranchWork reports whether the revert sha takes back a commit of
// the branch itself — one in base..HEAD, the set onBranch — and so names that
// commit's ticket as its own: the two cancel out on the branch. A revert of
// work the base already has is new work nobody accepted, whatever ticket the
// reverted commit named; it needs a ticket of its own. No line saying which
// commit it reverts, or one that cannot be read, is the same: the gate fails
// closed.
func (r *Repo) revertsBranchWork(sha string, onBranch map[string]bool) bool {
	body, err := r.run("log", "-1", "--format=%b", sha)
	if err != nil {
		return false
	}
	m := revertedCommit.FindStringSubmatch(body)
	return m != nil && onBranch[m[1]]
}

// isPlainMerge reports whether the merge commit sha has exactly two parents
// and the tree git's own merge of them produces, without a conflict.
func (r *Repo) isPlainMerge(sha string, parents []string) bool {
	if len(parents) != 2 {
		return false
	}
	auto, err := r.run("merge-tree", "--write-tree", parents[0], parents[1])
	if err != nil {
		return false // a conflict (exit 1) or a git without --write-tree
	}
	got, err := r.run("rev-parse", sha+"^{tree}")
	if err != nil {
		return false
	}
	first, _, _ := strings.Cut(strings.TrimSpace(auto), "\n")
	return first != "" && strings.TrimSpace(first) == strings.TrimSpace(got)
}

// FileAt returns the content of path, relative to Dir, as the commit rev
// holds it, and "" when rev holds no such file. rev must already be known to
// name a commit: any failure to read the file reads as empty, so a caller
// whose default for an empty file is the safe one fails closed.
func (r *Repo) FileAt(rev, path string) string {
	if rev == "" || strings.HasPrefix(rev, "-") {
		return ""
	}
	out, err := r.run("cat-file", "blob", rev+":./"+filepath.ToSlash(path))
	if err != nil {
		return ""
	}
	return out
}

// HandlesInSubjects reads the handles out of a list of subjects, one per line.
func HandlesInSubjects(subjects string) []string {
	var out []string
	seen := map[string]bool{}
	for _, line := range strings.Split(subjects, "\n") {
		m := subjectScope.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		for _, h := range strings.FieldsFunc(m[1], func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }) {
			if handleShape.MatchString(h) && !seen[h] {
				seen[h] = true
				out = append(out, h)
			}
		}
	}
	return out
}
