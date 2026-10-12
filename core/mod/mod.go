// Package mod ships the Claude Code mods jaira brings along, compiled into the
// binary the way the roles are.
//
// A mod is a plugin folder — .claude-plugin/plugin.json, hooks/, types/ — and
// Claude Code loads a plugin folder it finds in a skills directory by itself:
// ~/.claude/skills/<name> in every session, a project's .claude/skills/<name>
// once the workspace is trusted. So 'jaira roles install' writes each mod
// beside the roles, and nothing else is needed for it to load: no marketplace,
// no flag, no new command. Without that install nothing changes for anybody.
package mod

import (
	"embed"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// all: because the manifest sits in .claude-plugin/, which a plain pattern
// skips as a dot directory.
//
//go:embed all:builtin
var builtinFS embed.FS

const builtinDir = "builtin"

// What Claude Code writes into a mod folder it loads from disk (a
// --plugin-dir while developing): this build's API types and a tsconfig.json
// that points at them. go:embed reads the disk, not git, so a build made
// after such a run would carry them; it writes them again wherever the mod
// is installed.
const (
	generatedTypes  = ".claude-plugin/types"
	generatedConfig = "tsconfig.json"
)

// Mod is one plugin folder and the files that are installed of it.
type Mod struct {
	// ID is the folder name, which is also the plugin's name in its manifest.
	ID string
	// Files are paths relative to the mod's folder, sorted. Its tests are
	// left out: they run here, with 'claude plugin test', not where it is
	// installed. So is what Claude Code generated into the folder.
	Files []string
}

// Builtins returns every embedded mod, sorted by id.
func Builtins() ([]Mod, error) {
	entries, err := fs.ReadDir(builtinFS, builtinDir)
	if err != nil {
		return nil, err
	}
	var mods []Mod
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := Mod{ID: e.Name()}
		root := path.Join(builtinDir, e.Name())
		err := fs.WalkDir(builtinFS, root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// An io/fs path: always forward slashes, on every platform.
			//wintrap:ok
			rel := strings.TrimPrefix(p, root+"/")
			if d.IsDir() {
				if rel == "tests" || rel == generatedTypes {
					return fs.SkipDir
				}
				return nil
			}
			if rel == generatedConfig {
				return nil
			}
			m.Files = append(m.Files, rel)
			return nil
		})
		if err != nil {
			return nil, err
		}
		sort.Strings(m.Files)
		mods = append(mods, m)
	}
	return mods, nil
}

// File returns the embedded bytes of one file of a mod, rel as in Mod.Files.
func File(id, rel string) ([]byte, error) {
	return builtinFS.ReadFile(path.Join(builtinDir, id, rel))
}
