package role

import (
	"os"
	"path/filepath"
)

// skillsDir is where an agent harness looks for skills. Only .claude/ is
// written to by default, even though .codex/ and .agents/ hold the same kind of
// directory: a harness that reads more than one of them would find the same
// role registered twice under the same command name, and nothing would say
// which copy answered. Installing into several places is not the same choice as
// core/board/announce.go writing both AGENTS.md and CLAUDE.md — those are two
// files one reader reads once, while these are two registrations of one
// command. A project that keeps its skills elsewhere names that directory with
// --into rather than getting a copy in every candidate.
const skillsDir = ".claude/skills"

// ProjectTarget returns the skills directory to install into for a project
// rooted at root.
//
// The directory is returned, not created; Install creates what it writes to.
func ProjectTarget(root string) string {
	return filepath.Join(root, filepath.FromSlash(skillsDir))
}

// GlobalTarget returns ~/.claude/skills — the personal skills directory, for a
// role installed for the person rather than for one repository.
func GlobalTarget() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, filepath.FromSlash(skillsDir)), nil
}
