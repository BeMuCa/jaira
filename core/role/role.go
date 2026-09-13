// Package role ships the agent role prompts that drive a jaira board.
//
// A lane carries the instructions for one step of the work; a role carries the
// instructions for one worker who performs steps. The lanes travel with the
// repository, so a teammate who clones gets the board — but until now the roles
// lived in one person's ~/.claude/skills, so the teammate got a board and
// nobody to drive it. The roles are compiled into the binary for the same
// reason the lanes are, and 'jaira roles install' writes them out.
//
// A role is a directory, not a file: its SKILL.md is the prompt, and anything
// beside it (a reference document, a script) belongs to the same role and is
// installed with it. The directory name is the API — an agent harness derives
// the command name from it, not from the frontmatter — so the embedded
// directory names carry the jaira- prefix verbatim and nothing here rewrites
// them.
package role

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// all: rather than builtin/*, because a role's supporting files sit in
// subdirectories (teamlead/scripts/spawn.sh) and a plain glob would ship the
// prompt that references them without the files themselves.
//
//go:embed all:builtin
var builtinFS embed.FS

const builtinDir = "builtin"

// Role is one agent prompt, with every file that belongs to it.
type Role struct {
	// ID is the directory name, and therefore the name the harness invokes:
	// jaira-teamlead is /jaira-teamlead.
	ID          string
	Description string

	// Files are paths relative to the role's own directory, SKILL.md first
	// and the rest sorted. Always at least one entry.
	Files []string
}

// Builtins returns every embedded role, sorted by id.
func Builtins() ([]Role, error) {
	entries, err := fs.ReadDir(builtinFS, builtinDir)
	if err != nil {
		return nil, err
	}
	roles := make([]Role, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		r, err := load(e.Name())
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].ID < roles[j].ID })
	return roles, nil
}

// File returns the embedded bytes of one file of a role. rel is a path
// relative to the role's directory, as it appears in Role.Files.
func File(id, rel string) ([]byte, error) {
	// path.Join, not filepath.Join: builtinFS is an embedded filesystem and
	// always uses forward slashes, on every platform.
	return builtinFS.ReadFile(path.Join(builtinDir, id, rel))
}

// load walks one embedded role directory and reads its SKILL.md frontmatter.
func load(id string) (Role, error) {
	r := Role{ID: id}
	root := path.Join(builtinDir, id)
	err := fs.WalkDir(builtinFS, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		r.Files = append(r.Files, strings.TrimPrefix(p, root+"/"))
		return nil
	})
	if err != nil {
		return Role{}, err
	}
	sort.Slice(r.Files, func(i, j int) bool {
		// SKILL.md first: it is the role, the rest supports it.
		if (r.Files[i] == skillFile) != (r.Files[j] == skillFile) {
			return r.Files[i] == skillFile
		}
		return r.Files[i] < r.Files[j]
	})
	if len(r.Files) == 0 || r.Files[0] != skillFile {
		return Role{}, fmt.Errorf("role %s has no %s", id, skillFile)
	}
	b, err := File(id, skillFile)
	if err != nil {
		return Role{}, err
	}
	r.Description = frontmatterDescription(b)
	return r, nil
}

const skillFile = "SKILL.md"

// frontmatterDescription reads description: out of a SKILL.md header.
//
// Only description: — the frontmatter also carries name:, but a role's name is
// its directory, so reading the field back would only ever restate the id.
//
// The same parser the lane definitions use (core/lane/lane.go), because a
// SKILL.md header is the same frontmatter: a hand-rolled line scan would be a
// second reader of one file format, and would hand back the escapes of a
// quoted scalar verbatim.
func frontmatterDescription(b []byte) string {
	d, err := ticket.ParseDoc(b)
	if err != nil {
		return ""
	}
	v, _, err := d.Scalar("description")
	if err != nil {
		return ""
	}
	return v
}
