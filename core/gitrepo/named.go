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

// FileAt returns the content of path, relative to Dir, as the commit rev
// holds it, and false when rev holds no such file. rev must already be known
// to name a commit: any failure to read the file is reported as absent, so a
// caller that treats absent as its safe default fails closed.
func (r *Repo) FileAt(rev, path string) (string, bool) {
	if rev == "" || strings.HasPrefix(rev, "-") {
		return "", false
	}
	out, err := r.run("cat-file", "blob", rev+":./"+filepath.ToSlash(path))
	if err != nil {
		return "", false
	}
	return out, true
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
