package role

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/mod"
	"github.com/BeMuCa/jaira/core/shipped"
)

// Action is what Install did with one file.
type Action string

const (
	// Written: the file was not there.
	Written Action = "written"
	// Unchanged: the file is already byte-identical to the embedded one. A
	// second run of Install reports nothing but these.
	Unchanged Action = "unchanged"
	// Updated: the file was a version an earlier jaira shipped, so jaira wrote it
	// and nobody edited it; it was replaced with this binary's version without
	// --force.
	Updated Action = "updated"
	// Skipped: the file differs from the embedded one and from every version
	// jaira shipped, so somebody edited it.
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

// install writes every embedded role into dstSkillsDir as
// <dstSkillsDir>/<role-id>/<file>, and every embedded mod as
// <dstSkillsDir>/<mod-id>/<file>, and reports what happened to each file.
//
// The comparison is against the embedded bytes, not os.Stat as lane.Export
// does. Stat answers two questions — there or not — and this needs three: a
// file identical to the embedded one is an already-finished install and must
// not be an error, while a file that differs is somebody's own version and must
// not be touched. Only --force collapses the last case back into a write.
//
// A file that differs but equals a version an earlier jaira shipped is not an
// edit: an old binary wrote it. It is replaced and reported as Updated, without
// --force. With dryRun nothing is written or repaired; the results say what
// would happen.
//
// The mods are written only with mods set: they are Claude Code plugin
// folders, and a harness that does not load them would hold files nothing reads.
func install(dstSkillsDir string, force, mods, dryRun bool) ([]Result, error) {
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
			action, err := writeFile(dst, want, roleKey(r.ID, rel), force, dryRun)
			if err != nil {
				return out, err
			}
			out = append(out, Result{Role: r.ID, Path: dst, Action: action})
		}
	}
	// The mods jaira ships go beside the roles: a plugin folder in a skills
	// directory is one Claude Code loads by itself, so this is their install.
	// Same three-way comparison, so an edited mod file is left alone too.
	if !mods {
		return out, nil
	}
	builtins, err := mod.Builtins()
	if err != nil {
		return out, err
	}
	for _, m := range builtins {
		for _, rel := range m.Files {
			want, err := mod.File(m.ID, rel)
			if err != nil {
				return out, err
			}
			dst := filepath.Join(dstSkillsDir, m.ID, filepath.FromSlash(rel))
			action, err := writeFile(dst, want, modKey(m.ID, rel), force, dryRun)
			if err != nil {
				return out, err
			}
			out = append(out, Result{Role: m.ID, Path: dst, Action: action})
		}
	}
	return out, nil
}

// roleKey and modKey are a file's path in the jaira repository: the key the
// list of shipped versions (core/shipped) holds it under. They sit next to the
// go:embed directories they mirror.
func roleKey(id, rel string) string { return "core/role/" + path.Join(builtinDir, id, rel) }
func modKey(id, rel string) string  { return "core/mod/builtin/" + id + "/" + rel }

// writeFile applies the comparison to one file; key is its path in the jaira
// repository.
func writeFile(dst string, want []byte, key string, force, dryRun bool) (Action, error) {
	have, readErr := os.ReadFile(dst)
	existed := readErr == nil
	if !existed && !errors.Is(readErr, os.ErrNotExist) {
		return "", readErr
	}
	mode := modeOf(dst)
	if existed && bytes.Equal(have, want) {
		// The bytes are ours, so the mode is ours: a copy that arrived without
		// its execute bit — a plain cp, an unzip, a checkout on a filesystem
		// that drops the bit — is repaired here, because nothing is written
		// for it below.
		if dryRun {
			return Unchanged, nil
		}
		return Unchanged, os.Chmod(dst, mode)
	}
	action := Written
	if existed {
		if shipped.Known(key, have) {
			action = Updated
		} else if !force {
			return Skipped, nil
		} else {
			action = Overwritten
		}
	}
	if dryRun {
		return action, nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(dst, want, mode); err != nil {
		return "", err
	}
	// os.WriteFile passes the mode to O_CREATE only, so an existing file keeps
	// the permissions it already had. Say it again for the file we just
	// replaced, or --force restores the content of a script and leaves it
	// unexecutable.
	if err := os.Chmod(dst, mode); err != nil {
		return "", err
	}
	return action, nil
}

// modeOf picks the permissions one installed file gets.
//
// A role may ship a script beside its prompt (dispatcher/scripts/spawn.sh), and
// a script nobody can execute is a broken role. The embedded filesystem does not
// carry the permission bit, so the extension decides.
func modeOf(dst string) os.FileMode {
	if strings.HasSuffix(dst, ".sh") {
		return 0o755
	}
	return 0o644
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

// Install is install for this agent: with the mods where it loads them.
// With dryRun it reports what would happen and writes nothing.
func (a Agent) Install(dstSkillsDir string, force, dryRun bool) ([]Result, error) {
	return install(dstSkillsDir, force, a.Mods, dryRun)
}
