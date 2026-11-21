package cli

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/role"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
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
	cmd.AddCommand(newRolesListCmd(), newRolesInstallCmd(), newRolesPRGateCmd())
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
release_source.

And who opens a pull request: "agent" when the board's .jaira/pr says so —
an agent then opens it itself once every ticket on the branch is in done
('jaira roles pr-gate' checks that),
asks for the reviewer .jaira/pr-reviewer names and answers one round of
comments — and "person" otherwise, an agent pushing the branch and handing
the create line back. "none" when .jaira/pr says so: the project takes no
pull requests, and a branch whose tickets are all in done lands on the
default branch with a fast-forward push ('jaira roles pr-gate' checks that
too). Carried as pr, pr_source and pr_reviewer. And when:
"signoff" when the board's .jaira/pr-when says so — the pull request opens
once every ticket is in signoff, and its review round has to be finished
before done — and "done" otherwise. Carried as pr_when and pr_when_source.
Merging and approving stay a person's either way.`,
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
			pr, prSource := settings.PRFor(root)
			reviewer := settings.PRReviewerFor(root)
			prWhen, prWhenSource := settings.PRWhenFor(root)
			if pr == settings.PRNone {
				prWhen, prWhenSource = settings.PRWhenDone, "no pull request under .jaira/pr none"
			}
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
					"pr": pr, "pr_source": prSource, "pr_reviewer": reviewer,
					"pr_when": prWhen, "pr_when_source": prWhenSource,
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
			fmt.Fprintf(w, "pr: %s — who opens a pull request (%s); reviewer: %s; opens in: %s (%s)\n",
				pr, prSource, cmp.Or(reviewer, "none"), prWhen, prWhenSource)
			return nil
		},
	}
}

func newRolesPRGateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pr-gate <base>",
		Short: "Say whether an agent may open this branch's pull request",
		Long: `Says whether an agent may open the pull request of the checked-out branch
against <base> — the default branch of the repository it goes to, as
<remote>/<branch>. Exit 0 only when .jaira/pr says agent or none in <base> itself, at
least one commit in <base>..HEAD names a ticket in its subject
(type(A3K9QP): ...), and every ticket named that is still on the board sits
in its board's terminal lane — or, where .jaira/pr-when says signoff, in the
lane signoff or any lane after it up to the terminal one. Both files are
read as <base> holds them, not off the branch's disk: a branch cannot open
its own gate by committing them, and a board whose .jaira/ is not shared yet
keeps a person opening it. A ticket filed into the logbook is past done and
holds nothing up; one in the archive (rejected) or nowhere here (a typo, an
unfetched ref) holds it. Anything else exits 3 and says which ticket holds
it.

jaira-role-pr runs this before it opens a pull request on its own, so the
rule does not rest on a role's text alone. A person who invokes the role
needs no gate: their invocation is the decision.

Where .jaira/pr says none the same gate decides whether the branch may land
on <base> without a pull request: every ticket in the terminal lane, and
.jaira/pr-when is not read — with no pull request there is no review on one
to wait for. And every commit in <base>..HEAD must name a ticket — a merge
only when it carries content of its own: nobody reads the branch after the
gate, so a commit no ticket owns would land unseen. --json carries may_open,
pr, pr_when and tickets; a refusal is the usual error, its reason
pr_person, no_ticket, no_signoff_lane, unreadable, pending or, under none,
unnamed.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			repo := &gitrepo.Repo{Dir: s.Root}
			named, err := repo.TicketsNamed(args[0])
			if errors.Is(err, gitrepo.ErrNoSuchRef) {
				return fail(ExitUsage, "no_such_ref", "pr-gate: %q is no commit here — fetch the remote first", args[0])
			}
			if err != nil {
				return err
			}
			// Both board files are read as the base holds them, never off
			// the disk of the branch being judged: a branch that commits
			// .jaira/pr saying agent, or a pr-when of signoff, would
			// otherwise open its own gate. A change to either takes effect
			// once it has been merged into the base.
			pr := settings.PRIn(repo.FileAt(args[0], settings.BoardPRPath(".")))
			if pr == settings.PRPerson {
				return fail(ExitValidation, "pr_person",
					"pr-gate: .jaira/pr in %s says neither agent nor none — a person opens the pull request; hand the create line back", args[0])
			}
			// What the agent does on a no: hand the create line back where a
			// person opens pull requests; under none there is no such line.
			then := "hand the create line back"
			if pr == settings.PRNone {
				then = "push the branch and stop"
			}
			if len(named) == 0 {
				return fail(ExitValidation, "no_ticket",
					"pr-gate: no commit in %s..HEAD names a ticket — nobody accepted anything on this branch; %s", args[0], then)
			}
			// Under none nobody reads the branch after the gate: a pull
			// request would show a commit that names no ticket to a person
			// before it merges, a fast-forward push shows it to nobody. So
			// every commit must belong to a ticket a person accepted.
			if pr == settings.PRNone {
				unnamed, err := repo.UnnamedCommits(args[0])
				if err != nil {
					return err
				}
				if len(unnamed) > 0 {
					return fail(ExitValidation, "unnamed",
						"pr-gate: %d commit(s) in %s..HEAD name no ticket, and without a pull request nobody would see them before they land — %s; name a ticket in each, type(<ID>): ... — a revert too, as fix(<ID>): revert <what> (git rebase -i, reword) — or %s",
						len(unnamed), args[0], strings.Join(unnamed, "; "), then)
				}
			}
			// A ticket file that cannot be read is not a ticket filed away:
			// it would drop out of the list and look past done. A problem
			// naming a handle of this branch fails the gate closed; one
			// about another ticket of the board is not this branch's.
			all, err := s.List()
			var pe *ticket.PartialError
			if errors.As(err, &pe) {
				for _, p := range pe.Problems {
					for _, h := range named {
						if strings.Contains(p, h) {
							return fail(ExitValidation, "unreadable",
								"pr-gate: %s cannot be read, so whether it is accepted cannot be told — %s; fix it, or %s", h, p, then)
						}
					}
				}
			} else if err != nil {
				return err
			}
			when := settings.PRWhenIn(repo.FileAt(args[0], settings.BoardPRWhenPath(".")))
			// Without a pull request there is no review on one to wait for:
			// under none the branch lands on the default branch, and only a
			// person's acceptance in done is the go-ahead for that.
			if pr == settings.PRNone {
				when = settings.PRWhenDone
			}
			// reached is, per board, every lane a ticket may sit in for the
			// pull request to open: the terminal lane, and under signoff
			// also signoff and whatever lies between it and the terminal.
			// A parking lane ordered after the terminal one is not among
			// them: blocked work has not reached anything.
			reached := map[string]map[string]bool{}
			var pending []string
			for _, h := range named {
				onBoard := false
				for _, t := range all {
					if ticket.Handle(t.ID) != h {
						continue
					}
					onBoard = true
					b := lane.BoardOf(t)
					ok, known := reached[b]
					if !known {
						if ok, err = lanesReached(s, b, when); err != nil {
							return err
						}
						reached[b] = ok
					}
					if !ok[t.Status] {
						pending = append(pending, h+" ("+t.Status+")")
					}
				}
				// Off the board, a handle counts only where the logbook holds
				// it: filed there from done. Archived means rejected, and
				// unknown may be a typo in a scope or a ticket on a ref
				// nobody fetched — neither was accepted by anyone.
				if !onBoard {
					if where := filedWhere(s, h); where != "logbook" {
						pending = append(pending, h+" ("+where+")")
					}
				}
			}
			if len(pending) > 0 {
				if pr == settings.PRNone {
					then = "push the branch and stop; it lands once they are"
				}
				return fail(ExitValidation, "pending",
					"pr-gate: not every ticket on the branch is in %s — %s; %s", when, strings.Join(pending, ", "), then)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"may_open": true, "pr": pr, "pr_when": when, "tickets": named})
			}
			if pr == settings.PRNone {
				fmt.Fprintf(w, "pr-gate: no pull request on this board — the branch may land on %s with a fast-forward push; every ticket on it is in done: %s\n", args[0], strings.Join(named, ", "))
				return nil
			}
			fmt.Fprintf(w, "pr-gate: an agent may open it — every ticket on the branch is in %s: %s\n", when, strings.Join(named, ", "))
			return nil
		},
	}
}

// filedWhere says where a handle that is not on the board lives: "logbook",
// "archive" (rejected), or "unknown" — not here at all.
func filedWhere(s *ticket.Store, handle string) string {
	_, path, err := s.LoadAnywhere(handle)
	switch {
	case err != nil:
		return "unknown"
	case strings.HasPrefix(path, s.LogbookDir()+string(filepath.Separator)):
		return "logbook"
	case strings.HasPrefix(path, s.ArchiveDir()+string(filepath.Separator)):
		return "archive"
	}
	return "unknown"
}

// lanesReached is the set of lanes on board b in which a ticket counts as
// far enough along for the pull request to open under when.
func lanesReached(s *ticket.Store, b, when string) (map[string]bool, error) {
	lanes, err := boardLanes(s, b)
	if err != nil {
		return nil, err
	}
	term := lanes.Terminal()
	if term == nil {
		return nil, fail(ExitValidation, "no_terminal_lane", "pr-gate: board %q has no terminal lane, so nothing on it is ever done", cmp.Or(b, "default"))
	}
	from := term.ID
	if when == settings.PRWhenSignoff {
		from = settings.PRWhenSignoff
	}
	out := map[string]bool{}
	in := false
	for _, l := range lanes.Lanes {
		in = in || l.ID == from
		if in {
			out[l.ID] = true
		}
		if l.ID == term.ID {
			break
		}
	}
	if !out[term.ID] {
		return nil, fail(ExitValidation, "no_signoff_lane",
			"pr-gate: .jaira/pr-when says signoff, but board %q has no lane signoff before its terminal lane — fix the file or the board", cmp.Or(b, "default"))
	}
	return out, nil
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
