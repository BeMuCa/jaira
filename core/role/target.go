package role

import (
	"os"
	"path/filepath"
)

// agentDirs are the per-tool configuration directories a skills/ folder lives
// under. There is no single convention — .claude/ is Claude Code's, .codex/ is
// Codex's, .agents/ is the cross-tool attempt — and a project may have any
// combination of them. Roles are written into the ones that exist, on the same
// reasoning as the agent instruction files in core/board/announce.go: a role
// nobody's agent can find is a role that does not get used, and an extra
// markdown directory costs nothing.
var agentDirs = []string{".claude", ".codex", ".agents"}

const skillsSubdir = "skills"

// ProjectTargets returns the skills directories to install into for a project
// rooted at root: one per agent directory that already exists there, and
// .claude/skills when none does — a project that has never configured an agent
// still gets a working default rather than nothing.
//
// The directories are returned, not created; Install creates what it writes to.
func ProjectTargets(root string) []string {
	var out []string
	for _, d := range agentDirs {
		if fi, err := os.Stat(filepath.Join(root, d)); err == nil && fi.IsDir() {
			out = append(out, filepath.Join(root, d, skillsSubdir))
		}
	}
	if len(out) == 0 {
		out = append(out, filepath.Join(root, agentDirs[0], skillsSubdir))
	}
	return out
}

// GlobalTarget returns ~/.claude/skills — the personal skills directory, which
// is Claude Code's and has no per-tool spread to consider: a role installed
// globally is installed for the person, and this is where that lives.
func GlobalTarget() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, agentDirs[0], skillsSubdir), nil
}
