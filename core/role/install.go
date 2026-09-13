package role

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Action is what Install did with one file.
type Action string

const (
	// Written: the file was not there.
	Written Action = "written"
	// Unchanged: the file is already byte-identical to the embedded one. A
	// second run of Install reports nothing but these.
	Unchanged Action = "unchanged"
	// Skipped: the file differs from the embedded one, so somebody edited it.
	// It is left exactly as it is and reported, because silently replacing an
	// adjusted role is the one failure this command must never have.
	Skipped Action = "skipped"
	// Overwritten: it differed and --force said to replace it anyway.
	Overwritten Action = "overwritten"
)

// Result is one file's outcome.
type Result struct {
	Role   string
	Path   string
	Action Action
}

// Install writes every embedded role into dstSkillsDir as
// <dstSkillsDir>/<role-id>/<file>, and reports what happened to each file.
//
// The comparison is against the embedded bytes, not os.Stat as lane.Export
// does. Stat answers two questions — there or not — and this needs three: a
// file identical to the embedded one is an already-finished install and must
// not be an error, while a file that differs is somebody's own version and must
// not be touched. Only --force collapses the last case back into a write.
func Install(dstSkillsDir string, force bool) ([]Result, error) {
	roles, err := Builtins()
	if err != nil {
		return nil, err
	}
	var out []Result
	for _, r := range roles {
		for _, rel := range r.Files {
			want, err := File(r.ID, rel)
			if err != nil {
				return out, err
			}
			// The destination is built from the role's own id and its embedded
			// relative path, never from anything a caller supplied, so nothing
			// can escape dstSkillsDir. The check below states that as a
			// guarantee rather than trusting the construction of it.
			dst := filepath.Join(dstSkillsDir, r.ID, filepath.FromSlash(rel))
			if !within(dstSkillsDir, dst) {
				return out, fmt.Errorf("role %s: %s would write outside %s", r.ID, rel, dstSkillsDir)
			}
			action, err := writeFile(dst, want, force)
			if err != nil {
				return out, err
			}
			out = append(out, Result{Role: r.ID, Path: dst, Action: action})
		}
	}
	return out, nil
}

// writeFile applies the three-way comparison to one file.
func writeFile(dst string, want []byte, force bool) (Action, error) {
	have, readErr := os.ReadFile(dst)
	existed := readErr == nil
	switch {
	case existed && bytes.Equal(have, want):
		return Unchanged, nil
	case existed && !force:
		return Skipped, nil
	case !existed && !errors.Is(readErr, os.ErrNotExist):
		return "", readErr
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}
	// A role may ship a script beside its prompt (teamlead/scripts/spawn.sh),
	// and a script nobody can execute is a broken role. The embedded
	// filesystem does not carry the permission bit, so the extension decides.
	mode := os.FileMode(0o644)
	if strings.HasSuffix(dst, ".sh") {
		mode = 0o755
	}
	if err := os.WriteFile(dst, want, mode); err != nil {
		return "", err
	}
	if existed {
		return Overwritten, nil
	}
	return Written, nil
}

// Twins finds roles already present in dstSkillsDir under their unprefixed
// name — a directory teamlead/ beside the installed jaira-teamlead/.
//
// The prefix is not decoration: an agent harness derives the command name from
// the directory, so an unprefixed copy is a second, older command that answers
// to a different name. Anyone who wrote these prompts by hand before this
// command existed has exactly that, and will otherwise keep invoking the stale
// one without ever being told there are two. Reported, never deleted: a
// directory this tool did not write is not a directory it removes.
func Twins(dstSkillsDir string) ([]string, error) {
	roles, err := Builtins()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, r := range roles {
		bare := strings.TrimPrefix(r.ID, Prefix)
		if bare == r.ID {
			continue
		}
		if _, err := os.Stat(filepath.Join(dstSkillsDir, bare, skillFile)); err == nil {
			out = append(out, filepath.Join(dstSkillsDir, bare))
		}
	}
	return out, nil
}

// Prefix is carried by every shipped role's directory name, and is what keeps
// jaira's roles from colliding with a skill of the same name from elsewhere.
const Prefix = "jaira-"

// Skipped reports whether any file was left alone because it had been edited —
// what the CLI turns into a non-zero exit, so a script can tell "installed" from
// "installed except the ones you changed".
func SkippedAny(results []Result) bool {
	for _, r := range results {
		if r.Action == Skipped {
			return true
		}
	}
	return false
}

// within reports whether path is dir itself or below it.
func within(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
