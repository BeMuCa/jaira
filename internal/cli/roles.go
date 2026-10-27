package cli

import (
	"cmp"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/role"
	"github.com/BeMuCa/jaira/core/settings"
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
dispatcher who carries one ticket lane by lane — or, with "dispatch": "worker"
in ~/.jaira/settings.json, a worker who carries it through every lane itself —
and the agents they hand lanes to. They are compiled into this binary, and
'jaira roles install' writes them into your agent's skills directory.`,
	}
	cmd.AddCommand(newRolesListCmd(), newRolesInstallCmd())
	return cmd
}

func newRolesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the role prompts built into this binary",
		Long: `Lists the role prompts built into this binary, and which one a teamlead
starts for a whole ticket: the "dispatch" setting in ~/.jaira/settings.json.
"lanes" (the default) starts jaira-dispatcher, one new agent per lane;
"worker" starts jaira-worker, one agent through every lane, with only a lane
marked context: fresh handed to a new one. --json carries it as dispatch and
dispatch_role.

It also says which model the cockpit's explain and discuss chats start on:
the board's own .jaira/chat-model first, then "chat-model" in settings.json —
carried as chat_model, empty for the account's own default, and where it came
from as chat_model_source.

And which command the roles run their Python scripts with: JAIRA_PYTHON first,
then "python" in settings.json, else python3 — carried as python, and where it
came from as python_source. A command, so "uv run python" works too.

And who touches the repository's version and changelog: "teamlead" when the
board's .jaira/release says so — one line, committed, shared with everybody
on the board — and "worker" otherwise, every ticket bringing its own as the
repository's instructions say. Carried as release, and where it came from as
release_source.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			roles, err := role.Builtins()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			// Here rather than in a command of its own: the teamlead reads
			// this list to know what it can start, and which of two it starts
			// for a ticket is the same question.
			cur := settings.Load()
			root := bestEffortRoot()
			chat, chatSource := cur.ChatModelFor(root)
			py, pySource := cur.PythonCommand()
			release, releaseSource := settings.ReleaseFor(root)
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(roles))
				for _, r := range roles {
					arr = append(arr, map[string]any{
						"id": r.ID, "description": r.Description,
						"files": r.Files,
					})
				}
				return emit(w, map[string]any{
					"roles": arr, "dispatch": cur.DispatchMode(), "dispatch_role": cur.DispatchRole(),
					"chat_model": chat, "chat_model_source": chatSource,
					"python": py, "python_source": pySource,
					"release": release, "release_source": releaseSource,
				})
			}
			fmt.Fprintf(w, "%-24s %s\n", "ID", "DESCRIPTION")
			for _, r := range roles {
				fmt.Fprintf(w, "%-24s %s\n", r.ID, board.FirstSentence(r.Description))
			}
			fmt.Fprintf(w, "\ndispatch: %s — a ticket handed to an agent goes to /%s (\"dispatch\" in %s)\n",
				cur.DispatchMode(), cur.DispatchRole(), settings.Path())
			fmt.Fprintf(w, "chat model: %s — what the cockpit's explain and discuss chats start on (%s)\n",
				cmp.Or(chat, "account default"), chatSource)
			fmt.Fprintf(w, "python: %s — what the roles run their Python scripts with (%s)\n", py, pySource)
			fmt.Fprintf(w, "release: %s — who touches the version and the changelog (%s)\n", release, releaseSource)
			return nil
		},
	}
}

func newRolesInstallCmd() *cobra.Command {
	var project, global, force bool
	var into, agentName string
	cmd := &cobra.Command{
		Use:   "install --project|--global|--into <dir> [--agent <name>]",
		Short: "Write the role prompts into a skills directory",
		Long: `Writes every built-in role as <skills>/<role-id>/SKILL.md, together with any
files the role ships beside its prompt, the jaira skill as <skills>/jaira/SKILL.md
(not a role: it teaches the agent the CLI), and the jaira-cockpit mod as
<skills>/jaira-cockpit/ — a Claude Code plugin folder, which Claude Code loads
by itself from the next session on (from .claude/skills once you trust the
workspace). It keeps a teamlead chat readable and shows what waits on the
board; /cockpit inside Claude Code lists what it does.

--project writes into the agent's skills directory in this repository, so the
roles arrive with a clone. --global writes into the agent's personal skills
directory, for every repository you work in. --into names a skills directory
instead, for an agent jaira does not know.

--agent names the agent to install for, and so the directory: claude by
default, into .claude/skills or ~/.claude/skills. --agent copilot writes into
.github/skills with --project and ~/.copilot/skills with --global, and leaves
the cockpit out, which only Claude Code loads. Copilot also reads
.claude/skills and .agents/skills in a repository and ~/.agents/skills at home:
when the roles already sit in one of those, nothing is written, because a
second copy would register every role twice. --agent and --into exclude each
other: --into already names the directory.

A file you have edited is never replaced. It is left exactly as it is and
reported, and the command exits 3 so a script can tell "installed" from
"installed except the ones you changed"; --force replaces it. A file already
identical to the built-in one is not an edit — a second run reports it as
unchanged and exits 0. Neither is a file identical to a version an earlier
jaira shipped: it is replaced and reported as updated, without --force.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// --project opens the store below; it must not warn about the skills it is installing.
			updateInProgress = true
			defer func() { updateInProgress = false }()
			agent, ok := role.AgentNamed(agentName)
			if !ok {
				return fail(ExitUsage, "usage", "unknown agent %q; known agents: %s", agentName, role.AgentNames())
			}
			var target, root string
			var also []string
			switch {
			case into != "":
				if project || global {
					return fail(ExitUsage, "usage", "--into replaces --project and --global; pass one of the three")
				}
				if cmd.Flags().Changed("agent") {
					return fail(ExitUsage, "usage", "--into names the directory itself; pass --agent with --project or --global")
				}
				target = into
			case project == global:
				return fail(ExitUsage, "usage", "choose exactly one of --project, --global or --into")
			case global:
				dir, rest, err := agent.GlobalTarget()
				if err != nil {
					return err
				}
				target, also = dir, rest
			default:
				s, err := openStore()
				if err != nil {
					return err
				}
				root = s.Root
				target, also = agent.ProjectTarget(root)
			}

			// The roles already sitting where this agent reads them is an
			// install that is finished, not one to repeat beside it.
			present, err := role.InstalledIn(also)
			if err != nil {
				return err
			}
			if present != "" {
				// Without the command a later upgrade runs --agent copilot
				// again, is told the same thing, and the copy it reads stays old.
				update, err := role.UpdateCommand(present, root, global)
				if err != nil {
					return err
				}
				w := cmd.OutOrStdout()
				if g.jsonOut {
					return emit(w, map[string]any{
						"agent": agent.Name, "directory": target, "already_in": present,
						"update": update, "installed": []map[string]any{}, "skipped": false,
					})
				}
				fmt.Fprintf(w, "%s already reads the roles in %s; nothing written to %s — update them there with '%s'\n",
					agent.Name, present, target, update)
				return nil
			}

			results, err := agent.Install(target, force, false)
			if err != nil {
				return err
			}
			return reportRoleInstall(cmd, agent.Name, target, results)
		},
	}
	cmd.Flags().BoolVar(&project, "project", false, "write into the agent's skills directory in this repository ("+role.Dirs(false)+")")
	cmd.Flags().BoolVar(&global, "global", false, "write into the agent's personal skills directory ("+role.Dirs(true)+")")
	cmd.Flags().StringVar(&into, "into", "", "write into this skills directory instead")
	cmd.Flags().StringVar(&agentName, "agent", role.DefaultAgent, "the agent to install for: "+role.AgentNames())
	cmd.Flags().BoolVar(&force, "force", false, "replace a file you have edited")
	return cmd
}

func reportRoleInstall(cmd *cobra.Command, agent, target string, results []role.Result) error {
	w := cmd.OutOrStdout()
	skipped := role.SkippedAny(results)
	if g.jsonOut {
		if err := emit(w, map[string]any{"agent": agent, "directory": target, "installed": roleResultsJSON(results), "skipped": skipped}); err != nil {
			return err
		}
	} else {
		printRoleResults(w, target, results, false)
	}
	if skipped {
		// A refusal, not a failure: everything else was installed, and the
		// files left alone were left alone on purpose.
		return fail(ExitValidation, "edited", "some roles were edited here and left alone; --force to replace them")
	}
	return nil
}

// roleResultsJSON is one entry per file, as 'roles install --json' lists them.
func roleResultsJSON(results []role.Result) []map[string]any {
	arr := make([]map[string]any, 0, len(results))
	for _, r := range results {
		arr = append(arr, map[string]any{
			"role": r.Role, "path": r.Path, "action": string(r.Action),
		})
	}
	return arr
}

// printRoleResults writes the text report of one install. With dry the lines
// say what would happen instead of what did.
func printRoleResults(w io.Writer, target string, results []role.Result, dry bool) {
	counts := map[role.Action]int{}
	for _, r := range results {
		counts[r.Action]++
		// Only the actions that changed or withheld something are worth a line
		// each; a run that writes nothing new should be quiet.
		switch r.Action {
		case role.Written, role.Updated, role.Overwritten:
			if dry {
				verb := map[role.Action]string{role.Written: "write", role.Updated: "update", role.Overwritten: "overwrite"}[r.Action]
				fmt.Fprintf(w, "would %s %s\n", verb, r.Path)
			} else {
				fmt.Fprintf(w, "%s %s\n", r.Action, r.Path)
			}
		case role.Skipped:
			fmt.Fprintf(w, "%s %s (edited here; --force to replace)\n", r.Action, r.Path)
		}
	}
	fmt.Fprintf(w, "%d written, %d updated, %d unchanged, %d skipped, %d overwritten\n",
		counts[role.Written], counts[role.Updated], counts[role.Unchanged], counts[role.Skipped], counts[role.Overwritten])
	fmt.Fprintf(w, "roles in %s\n", target)
}
