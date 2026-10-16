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

// Target is one skills directory jaira's roles are already installed in.
type Target struct {
	Agent  Agent
	Dir    string
	Global bool
}

// Installed returns every directory the built-in roles already sit in: each agent's own
// directory in the repository at root (none when root is ""), then each agent's personal
// one. A directory holding none of them is not listed — writing into it would be a first
// install, which is 'jaira roles install' and a choice, not an update. A directory reached
// twice (a board rooted at the home directory) is listed once.
func Installed(root string) ([]Target, error) {
	var out []Target
	seen := map[string]bool{}
	add := func(a Agent, dir string, global bool) error {
		dir = filepath.Clean(dir)
		if seen[dir] {
			return nil
		}
		found, err := InstalledIn([]string{dir})
		if err != nil {
			return err
		}
		if found == "" {
			return nil
		}
		seen[dir] = true
		out = append(out, Target{Agent: a, Dir: dir, Global: global})
		return nil
	}
	if root != "" {
		for _, a := range agents {
			dir, _ := a.ProjectTarget(root)
			if err := add(a, dir, false); err != nil {
				return nil, err
			}
		}
	}
	for _, a := range agents {
		dir, _, err := a.GlobalTarget()
		if err != nil {
			return nil, err
		}
		if err := add(a, dir, true); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Behind counts the files an install would write or replace: missing ones and older
// shipped versions. Edited files are not behind; they are the person's.
func Behind(results []Result) int {
	n := 0
	for _, r := range results {
		if r.Action == Written || r.Action == Updated {
			n++
		}
	}
	return n
}

// Dirs lists, for a flag's help, the directory each agent writes in one scope:
// "claude: .claude/skills, copilot: .github/skills", or with ~/ for global.
func Dirs(global bool) string {
	parts := make([]string, len(agents))
	for i, a := range agents {
		dir := a.Project
		if global {
			dir = "~/" + a.Global
		}
		parts[i] = a.Name + ": " + dir
	}
	return strings.Join(parts, ", ")
}

// UpdateCommand returns the jaira command that updates the roles already in
// dir: the --project or --global run of the agent that writes there, and
// --into for a directory none of them writes. root is the repository root
// for a project scope, ignored for global.
func UpdateCommand(dir, root string, global bool) (string, error) {
	scope := "--project"
	if global {
		scope = "--global"
	}
	for _, a := range agents {
		var target string
		if global {
			var err error
			if target, _, err = a.GlobalTarget(); err != nil {
				return "", err
			}
		} else {
			target, _ = a.ProjectTarget(root)
		}
		if target != dir {
			continue
		}
		if a.Name == DefaultAgent {
			return "jaira roles install " + scope, nil
		}
		return "jaira roles install --agent " + a.Name + " " + scope, nil
	}
	return "jaira roles install --into " + dir, nil
}
