package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/role"
)

// newRolesCmd groups the role prompts the binary ships.
//
// Lanes carry the instructions for a step; roles carry the instructions for the
// worker that performs steps. The lanes already travel with the repository, so
// a teammate who clones gets the board — this is what gets them the drivers too.
func newRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "The agent role prompts this binary ships",
		Long: `Roles are the agent prompts that drive a board: a teamlead who talks to you, a
dispatcher who carries one ticket lane by lane, and the workers they hand lanes
to. They are compiled into this binary, and 'jaira roles install' writes them
into your agent's skills directory.`,
	}
	cmd.AddCommand(newRolesListCmd(), newRolesInstallCmd())
	return cmd
}

func newRolesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the role prompts built into this binary",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			roles, err := role.Builtins()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(roles))
				for _, r := range roles {
					arr = append(arr, map[string]any{
						"id": r.ID, "name": r.Name,
						"description": r.Description, "files": r.Files,
					})
				}
				return emit(w, map[string]any{"roles": arr})
			}
			fmt.Fprintf(w, "%-24s %s\n", "ID", "DESCRIPTION")
			for _, r := range roles {
				fmt.Fprintf(w, "%-24s %s\n", r.ID, firstSentence(r.Description))
			}
			return nil
		},
	}
}

// firstSentence keeps the listing to one line per role. The full description is
// in --json and in the file itself.
func firstSentence(s string) string {
	for i, c := range s {
		if c == '.' {
			return s[:i+1]
		}
	}
	return s
}

func newRolesInstallCmd() *cobra.Command {
	var project, global, force bool
	cmd := &cobra.Command{
		Use:   "install --project|--global",
		Short: "Write the role prompts into a skills directory",
		Long: `Writes every built-in role as <skills>/<role-id>/SKILL.md, together with any
files the role ships beside its prompt.

--project writes into this repository, so the roles arrive with a clone: into
each of .claude/, .codex/ and .agents/ that already exists, and into .claude/
when none does. --global writes into ~/.claude/skills, for every repository you
work in.

A file you have edited is never replaced. It is left exactly as it is and
reported, and the command exits 3 so a script can tell "installed" from
"installed except the ones you changed"; --force replaces it. A file already
identical to the built-in one is not an edit — a second run reports it as
unchanged and exits 0.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if project == global {
				return fail(ExitUsage, "usage", "choose exactly one of --project or --global")
			}
			var targets []string
			if global {
				dir, err := role.GlobalTarget()
				if err != nil {
					return err
				}
				targets = []string{dir}
			} else {
				s, err := openStore()
				if err != nil {
					return err
				}
				targets = role.ProjectTargets(s.Root)
			}

			var all []role.Result
			var twins []string
			for _, dir := range targets {
				results, err := role.Install(dir, force)
				all = append(all, results...)
				if err != nil {
					return err
				}
				found, err := role.Twins(dir)
				if err != nil {
					return err
				}
				twins = append(twins, found...)
			}
			return reportRoleInstall(cmd, all, twins)
		},
	}
	cmd.Flags().BoolVar(&project, "project", false, "write into this repository's agent skills directories")
	cmd.Flags().BoolVar(&global, "global", false, "write into ~/.claude/skills")
	cmd.Flags().BoolVar(&force, "force", false, "replace a file you have edited")
	return cmd
}

func reportRoleInstall(cmd *cobra.Command, results []role.Result, twins []string) error {
	w := cmd.OutOrStdout()
	skipped := role.SkippedAny(results)
	if g.jsonOut {
		arr := make([]map[string]any, 0, len(results))
		for _, r := range results {
			arr = append(arr, map[string]any{
				"role": r.Role, "path": r.Path, "action": string(r.Action),
			})
		}
		if twins == nil {
			twins = []string{}
		}
		if err := emit(w, map[string]any{"installed": arr, "skipped": skipped, "unprefixed": twins}); err != nil {
			return err
		}
	} else {
		counts := map[role.Action]int{}
		for _, r := range results {
			counts[r.Action]++
			// Only the two actions that changed or withheld something are
			// worth a line each; a run that writes nothing new should be quiet.
			switch r.Action {
			case role.Written, role.Overwritten:
				fmt.Fprintf(w, "%s %s\n", r.Action, r.Path)
			case role.Skipped:
				fmt.Fprintf(w, "%s %s (edited here; --force to replace)\n", r.Action, r.Path)
			}
		}
		fmt.Fprintf(w, "%d written, %d unchanged, %d skipped, %d overwritten\n",
			counts[role.Written], counts[role.Unchanged], counts[role.Skipped], counts[role.Overwritten])
		for _, dir := range installDirs(results) {
			fmt.Fprintf(w, "roles in %s\n", dir)
		}
		for _, t := range twins {
			fmt.Fprintf(w, "note: %s is an older copy under its unprefixed name; it answers to a different command and was left alone\n", t)
		}
	}
	if skipped {
		// A refusal, not a failure: everything else was installed, and the
		// files left alone were left alone on purpose.
		return fail(ExitValidation, "edited", "some roles were edited here and left alone; --force to replace them")
	}
	return nil
}

// installDirs names each skills directory written to, once, in order.
func installDirs(results []role.Result) []string {
	var out []string
	seen := map[string]bool{}
	for _, r := range results {
		// A role's supporting file sits one level deeper than its prompt, so
		// the skills directory is found by climbing to the role id, not by
		// counting parents.
		dir := skillsDirOf(r.Path, r.Role)
		if dir == "" || seen[dir] {
			continue
		}
		seen[dir] = true
		out = append(out, dir)
	}
	return out
}

// skillsDirOf walks up from an installed file to the skills directory that
// holds the role, by finding the role-id element in the path.
func skillsDirOf(path, roleID string) string {
	dir := filepath.Dir(path)
	for dir != "" && dir != string(filepath.Separator) && dir != "." {
		if filepath.Base(dir) == roleID {
			return filepath.Dir(dir)
		}
		dir = filepath.Dir(dir)
	}
	return ""
}
