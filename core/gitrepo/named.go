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
// The subject git writes for 'git revert', Revert "type(A3K9QP): ...", names
// no ticket: whose work a revert takes back depends on where the reverted
// commit sits, which no subject says. A revert is written as
// fix(A3K9QP): revert ... (jaira-role-pr says so).
var subjectScope = regexp.MustCompile(`^[a-z]+\(([^)]*)\)!?:`)

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
// ticket, as "<short sha> <subject>" in the order git lists them. A merge
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
	var unnamed []string
	for _, line := range strings.Split(out, "\n") {
		f := strings.SplitN(strings.TrimSpace(line), "\x00", 4)
		if len(f) != 4 || len(HandlesInSubjects(f[3])) > 0 {
			continue
		}
		if parents := strings.Fields(f[2]); len(parents) > 1 && r.isPlainMerge(f[0], parents) {
			continue
		}
		unnamed = append(unnamed, f[1]+" "+f[3])
	}
	return unnamed, nil
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

// IsCommit reports whether rev names a commit here. A rev starting with a
// dash is never one: git would read it as an option.
func (r *Repo) IsCommit(rev string) bool {
	if rev == "" || strings.HasPrefix(rev, "-") {
		return false
	}
	_, err := r.run("rev-parse", "--verify", "--quiet", rev+"^{commit}")
	return err == nil
}

// MergeBase is the best common ancestor of a and b, or "" when there is
// none or either is not a commit here.
func (r *Repo) MergeBase(a, b string) string {
	if strings.HasPrefix(a, "-") || strings.HasPrefix(b, "-") {
		return ""
	}
	out, err := r.run("merge-base", a, b)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// SymbolicRef is the full ref name the symbolic ref name points at, or ""
// when it is not one here.
func (r *Repo) SymbolicRef(name string) string {
	out, err := r.run("symbolic-ref", "-q", name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
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
