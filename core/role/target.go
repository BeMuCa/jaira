package role

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Agent is one harness the roles can be installed for: where it reads skills,
// in a repository and in the home directory.
//
// Only one directory per scope is written, even where the harness reads several:
// a harness that reads more than one of them would find the same role
// registered twice under the same command name, and nothing would say which
// copy answered. Installing into several places is not the same choice as
// core/board/announce.go writing both AGENTS.md and CLAUDE.md — those are two
// files one reader reads once, while these are two registrations of one
// command. A project that keeps its skills elsewhere names that directory with
// --into rather than getting a copy in every candidate.
type Agent struct {
	Name string

	// Project and Global are the directories written to: Project relative to
	// the repository root, Global relative to the home directory.
	Project, Global string

	// AlsoProject and AlsoGlobal are the other directories every surface of
	// the harness reads. Roles already in one of them are already installed for
	// it, and a write into Project or Global would be the second copy.
	AlsoProject, AlsoGlobal []string

	// Mods: whether the harness loads the jaira mods. They are Claude Code
	// plugin folders, and anywhere else they are files nothing reads.
	Mods bool
}

// agents is every harness 'jaira roles install --agent' knows, the default
// first. A further one is a further entry.
//
// Copilot, as of 2026-10: VS Code, Copilot CLI and the cloud agent all read
// .github/skills, .claude/skills and .agents/skills in a repository; in the
// home directory CLI and cloud agent read ~/.copilot/skills and
// ~/.agents/skills, and only VS Code reads ~/.claude/skills as well — which is
// why that one does not count as already installed: the CLI would be left
// without roles. Sources are on ticket E75R4Q.
var agents = []Agent{
	{Name: "claude", Project: ".claude/skills", Global: ".claude/skills", Mods: true},
	{
		Name: "copilot", Project: ".github/skills", Global: ".copilot/skills",
		AlsoProject: []string{".claude/skills", ".agents/skills"},
		AlsoGlobal:  []string{".agents/skills"},
	},
}

// DefaultAgent is the agent installed for when none is named.
const DefaultAgent = "claude"

// AgentNamed returns the agent called name, and false when there is none.
func AgentNamed(name string) (Agent, bool) {
	for _, a := range agents {
		if a.Name == name {
			return a, true
		}
	}
	return Agent{}, false
}

// AgentNames lists the known agents, for a message naming them.
func AgentNames() string {
	names := make([]string, len(agents))
	for i, a := range agents {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

// ProjectTarget returns the skills directory to install into for a project
// rooted at root, and the other directories the agent reads there.
//
// The directory is returned, not created; Install creates what it writes to.
func (a Agent) ProjectTarget(root string) (string, []string) {
	return under(root, a.Project, a.AlsoProject)
}

// GlobalTarget returns the personal skills directory — for a role installed for
// the person rather than for one repository — and the other personal
// directories the agent reads.
func (a Agent) GlobalTarget() (string, []string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil, err
	}
	dir, also := under(home, a.Global, a.AlsoGlobal)
	return dir, also, nil
}

func under(base, dir string, also []string) (string, []string) {
	out := make([]string, len(also))
	for i, d := range also {
		out[i] = filepath.Join(base, filepath.FromSlash(d))
	}
	return filepath.Join(base, filepath.FromSlash(dir)), out
}

// InstalledIn returns the first of dirs that already holds a built-in role's
// SKILL.md, or "" when none does. One role is enough: that role is the one a
// second install would register twice.
func InstalledIn(dirs []string) (string, error) {
	roles, err := Builtins()
	if err != nil {
		return "", err
	}
	for _, dir := range dirs {
		for _, r := range roles {
			_, err := os.Stat(filepath.Join(dir, r.ID, skillFile))
			if err == nil {
				return dir, nil
			}
			if !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
		}
	}
	return "", nil
}
