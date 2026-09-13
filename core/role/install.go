package role

import (
	"bytes"
	"errors"
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
			// Both path elements come out of the embedded filesystem — the
			// role's own directory name and its relative path — and never
			// from a caller, so the destination is fixed at compile time and
			// stays under dstSkillsDir by construction.
			dst := filepath.Join(dstSkillsDir, r.ID, filepath.FromSlash(rel))
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

// SkippedAny reports whether any file was left alone because it had been edited —
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
