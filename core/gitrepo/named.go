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
// that names none counts only when it carries content of its own — a
// conflict resolution, a file added during 'merge --no-commit' — which is
// what 'git diff-tree --cc' shows; a clean merge brings in nothing that is
// not listed on its own: a branch merged in shows its commits here, the base
// merged into the branch shows nothing, since the base is excluded.
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
		if len(strings.Fields(f[2])) > 1 {
			own, err := r.run("diff-tree", "--cc", "-p", "--no-commit-id", f[0])
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(own) == "" {
				continue
			}
		}
		unnamed = append(unnamed, f[1]+" "+f[3])
	}
	return unnamed, nil
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
