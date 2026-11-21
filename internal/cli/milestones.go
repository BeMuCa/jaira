package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// milestoneLockName keys the store lock a milestone is written under; the name
// lives in core/milestone, where the shared write path takes it, and the
// commands that still lock by hand (restore, the import) read it from there.
const milestoneLockName = milestone.LockName

// milestoneIndex is the one place a milestone index is built for the CLI, so
// the board's filter and 'jaira list --milestone' cannot disagree about what
// belonging to a milestone means. The TUI builds the same index from the same
// call in its reload.
//
// A board with no milestones is the common case and costs one for-each-ref,
// which is why this is called on the list path without a flag guarding it.
func milestoneIndex(s *ticket.Store) (milestone.Index, error) {
	importMilestones(s)
	all, err := milestone.LoadAll(refs)
	if err != nil {
		return nil, err
	}
	return milestone.Build(all), nil
}

func newMilestoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "Group tickets into a round of work",
		Long: `A milestone is the set of tickets that belong to one round of work, kept on a
git ref of its own, ` + "`refs/jaira/milestones/<ID>`" + `, and nowhere else: there is
no file to commit, so a milestone never travels through a branch and a stale
copy can never overwrite it. ` + "`jaira fetch`" + ` brings the team's milestones, and
a board without a remote keeps them on the local ref until one is configured.
The six characters are its ID, which stays when the name changes, and every
command takes a milestone as its name, its ID or ` + "`<name>-<ID>`" + ` — a name two
milestones carry is refused (exit 3) with both forms to run it again with.

It is a document of its own rather than a field on each ticket because of what
happens at the end of a round: the work that did not finish has to appear in
the next one, and a per-ticket marker means touching every ticket one at a
time. Here the unfinished tickets move into the next milestone in one edit
instead of twenty.

The ref carries one file, ` + "`<name>-<ID>.md`" + ` — read it with
` + "`git show refs/jaira/milestones/<ID>:<name>-<ID>.md`" + `. Its frontmatter carries
the colour and when it was created, and below it one ticket id per line; jaira
keeps every other line exactly as it found it. A milestone file an older build
left under .jaira/milestones/ is taken onto its ref by the next milestone
command or fetch and deleted, with a line on stderr.

Each milestone is given a random free colour, which its cards then show as a
bar down their RIGHT edge — the left edge belongs to the tags. On the board, M
opens the picker and narrows everything to the milestones ticked in it.

A milestone leaves the board only when somebody says so: ` + "`jaira logbook <name>`" + `
marks it filed once all its tickets are finished, and ` + "`jaira restore <name>`" + `
brings it back. Emptying one with ` + "`rm`" + ` leaves it standing.

Two milestones that share a name are warned about on stderr by every milestone
command (--json: milestones_twins), with the rename that ends it. Nothing is
refused for it, but a command given the bare name exits 3.`,
		// ls says it in its own report, where --json has a field to put it in.
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			if cmd.Name() != "ls" && openedStore != nil {
				reportTwins()
			}
		},
	}
	cmd.AddCommand(
		newMilestoneCreateCmd(),
		newMilestoneAddCmd(),
		newMilestoneRmCmd(),
		newMilestoneSetCmd(),
		newMilestoneRenameCmd(),
		newMilestoneLsCmd(),
	)
	return cmd
}

func newMilestoneCreateCmd() *cobra.Command {
	var colour int
	var goal, after string
	cmd := &cobra.Command{
		Use:   "create <name> [id...]",
		Short: "Start a milestone, optionally with its first tickets",
		Long: `Creates the milestone on ` + "`refs/jaira/milestones/<ID>`" + `, with a random ID that no
milestone in this clone has, and puts any tickets named after the name into it.

The colour is picked for you, at random from the colours no other milestone on
this board is using: concurrent milestones are never many, so a clash once the
palette is spent is survivable, and nobody should have to choose one. --color
<1-255> overrides it; 0 is not black here but "no colour", and would paint no
cell on any card.

A name a milestone on the board has is refused (exit 3). A filed milestone, or
one a teammate made that no ` + "`jaira fetch`" + ` has written here yet, does not hold its
name: the ID tells them apart, and two with one name show as <name>-<ID>.

Names are lowercase kebab, the same rule tags follow — the name is part of the
file name its ref carries, so it has to be safe in a path. "Round One" is filed as "round-one"
and you are told so.

--goal <text> says in one line what the milestone is for, and --after <a,b>
puts it on the Milestone Roadmap (L on the board) under those milestones; the first one
decides its swimlane. --after "" makes it the top of a swimlane of its own.
Without --after the milestone stays off the roadmap. A predecessor has to be
a milestone on this board, and the arrows may not form a loop.`,
		Args: minArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			importMilestones(s)
			name, changed, err := milestone.NormalizeName(args[0])
			if err != nil {
				return fail(ExitUsage, "bad_milestone", "%v", err)
			}
			ids, err := resolveAll(s, args[1:])
			if err != nil {
				return err
			}
			var preds []string
			if cmd.Flags().Changed("after") {
				if preds, err = parseMilestoneNames(after); err != nil {
					return err
				}
			}

			ms, err := milestone.Create(s, refs, func(existing []*milestone.Milestone) (*milestone.Milestone, error) {
				var ee *milestone.ExistsError
				switch err := milestone.CheckFree(refs, name); {
				case errors.As(err, &ee):
					// A name two milestones carry is refused by 'add' as ambiguous, so
					// the hint spells the <name>-<ID> forms it accepts.
					labels := milestone.Labels(existing)
					var forms []string
					for _, m := range existing {
						if m.Name == name {
							forms = append(forms, labels[m.ID])
						}
					}
					target := strings.Join(forms, "|")
					if len(forms) > 1 {
						target = "<" + target + ">"
					}
					return nil, fail(ExitValidation, "milestone_exists",
						"milestone %q already exists (%s%s) — 'jaira milestone add %s <id>' puts tickets in it",
						name, milestone.RefPrefix, ee.ID, target)
				case err != nil:
					return nil, err
				}
				c := colour
				if !cmd.Flags().Changed("color") {
					c = milestone.AssignColour(existing, name)
				} else if !tag.ValidColour(c) || c == 0 {
					// 0 is not black here, it is "no colour": a milestone given it
					// would paint nothing on any card, which is not a thing to be
					// asked for by hand.
					return nil, fail(ExitUsage, "bad_color", "--color takes an ANSI-256 value, 1-255; got %d", c)
				}
				ms := milestone.New(name, c, time.Now())
				for _, id := range ids {
					ms.Add(id)
				}
				if goal != "" {
					ms.SetGoal(goal)
				}
				if cmd.Flags().Changed("after") {
					ids, err := resolvePredecessors(existing, ms.ID, preds)
					if err != nil {
						return nil, err
					}
					ms.SetAfter(ids)
				}
				return ms, nil
			})
			if errors.Is(err, refsync.ErrNoMilestoneRefs) {
				return refuseNoRefs()
			}
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, milestoneJSON(ms, nil))
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[0], name)
			}
			fmt.Fprintf(w, "Milestone %q created with %d ticket(s): %s\n",
				name, len(ms.Members()), ms.Ref())
			return nil
		},
	}
	cmd.Flags().IntVar(&colour, "color", -1, "ANSI-256 colour (1-255; 0 paints no cell) instead of a random free one")
	cmd.Flags().StringVar(&goal, "goal", "", "one line saying what the milestone is for")
	cmd.Flags().StringVar(&after, "after", "", "comma-separated milestones this one follows, placing it on the Milestone Roadmap; \"\" makes it a swimlane top")
	return cmd
}

func newMilestoneAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <id>...",
		Short: "Put tickets into a milestone",
		Long: `Adds one line per ticket to the milestone.

Every id given is added in one write, which is the whole point: grouping twenty
tickets is one write of one ref, not twenty edits of twenty ticket files each
travelling on its own ref.

A ticket already in the milestone is left where it is rather than moved to the
end — re-adding must not reorder a list somebody arranged by hand.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editMembers(cmd, args, true)
		},
	}
}

func newMilestoneRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name> <id>...",
		Short: "Take tickets out of a milestone",
		Long: `Removes one line per ticket from the milestone, in one write, leaving
every other line exactly where it was.

The milestone itself stays standing, even when its last ticket leaves. A
milestone goes away on command and never on its own — an empty one is still a
plan, and a milestone that deleted itself when emptied would take a freshly
created one with it the moment you changed your mind about its first ticket.

` + "`jaira logbook <name>`" + ` is the command that takes a milestone off the board,
and ` + "`jaira restore <name>`" + ` brings it back.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editMembers(cmd, args, false)
		},
	}
}

func newMilestoneSetCmd() *cobra.Command {
	var goal, after string
	var noRoadmap bool
	cmd := &cobra.Command{
		Use:   "set <name> [--goal <text>] [--after <a,b>] [--no-roadmap]",
		Short: "Change a milestone's goal or its place on the Milestone Roadmap",
		Long: `Rewrites the ` + "`goal:`" + ` and ` + "`after:`" + ` lines of the milestone and nothing else.

--goal <text> sets the one-line goal; --goal "" removes it.
--after <a,b> places the milestone on the Milestone Roadmap (L on the board) under those
milestones. The first one decides its swimlane, every further one is an extra
arrow. --after "" makes it the top of a swimlane of its own.
--no-roadmap takes it off the roadmap again; its tickets stay where they are,
and what came after it now comes after what came before it, as x on the
roadmap does.

At least one flag is needed, and --after and --no-roadmap exclude each other. A
predecessor has to be a milestone on this board (exit 5 otherwise), and an
arrow that would close a loop is refused (exit 3) with the loop named. A
filed milestone is refused like everywhere else.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setGoal, setAfter := cmd.Flags().Changed("goal"), cmd.Flags().Changed("after")
			if !setGoal && !setAfter && !noRoadmap {
				return fail(ExitUsage, "usage", "milestone set needs at least one of --goal, --after, --no-roadmap")
			}
			if setAfter && noRoadmap {
				return fail(ExitUsage, "usage", "--after places the milestone on the roadmap and --no-roadmap takes it off; give one")
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			importMilestones(s)
			name, err := milestoneKey(args[0])
			if err != nil {
				return err
			}
			var preds []string
			if setAfter {
				if preds, err = parseMilestoneNames(after); err != nil {
					return err
				}
			}

			var before []byte
			ms, err := milestone.Edit(s, refs, name, func(ms *milestone.Milestone) error {
				before = ms.Bytes()
				if setAfter {
					all, err := milestone.LoadAll(refs)
					if err != nil {
						return err
					}
					ids, err := resolvePredecessors(all, ms.ID, preds)
					if err != nil {
						return err
					}
					ms.SetAfter(ids)
				}
				if noRoadmap {
					ms.Unplace()
				}
				if setGoal {
					ms.SetGoal(goal)
				}
				return nil
			})
			if err := editRefusal(s, name, err, ", and then it can be changed again"); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, milestoneJSON(ms, nil))
			}
			label := milestoneLabel(ms)
			if string(before) == string(ms.Bytes()) {
				fmt.Fprintf(w, "Nothing to do: milestone %q is already as asked.\n", label)
				return nil
			}
			fmt.Fprintf(w, "Milestone %q updated: %s\n", label, ms.Ref())
			return nil
		},
	}
	cmd.Flags().StringVar(&goal, "goal", "", "one line saying what the milestone is for (\"\" removes it)")
	cmd.Flags().StringVar(&after, "after", "", "comma-separated milestones this one follows; \"\" makes it a swimlane top")
	cmd.Flags().BoolVar(&noRoadmap, "no-roadmap", false, "take the milestone off the Milestone Roadmap")
	return cmd
}

func newMilestoneRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Give a milestone a new name",
		Long: `Changes the name of a milestone and nothing else: its ref carries it as
` + "`<new>-<ID>.md`" + ` from then on, so it is one write, and the ID, the tickets, the
goal, the colour, the place on the Milestone Roadmap and the ` + "`after:`" + ` lines of
the milestones that follow it (they name the ID) all stay as they were. A clone
that has the milestone sees the new name after its next ` + "`jaira fetch`" + `.

<old> is the milestone's name, ID or <name>-<ID>. The new name follows the rules
of ` + "`create`" + `: lowercase kebab, and you are told when it was changed. A name a
milestone on the board has is refused (exit 3); an old name that is not on the
board exits 5, or 3 when it was filed.`,
		Args: exactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			importMilestones(s)
			from, err := milestoneKey(args[0])
			if err != nil {
				return err
			}
			to, changed, err := milestone.NormalizeName(args[1])
			if err != nil {
				return fail(ExitUsage, "bad_milestone", "%v", err)
			}

			fromName, fromLabel := from, from
			if old, lerr := milestone.Load(refs, from); lerr == nil {
				fromName, fromLabel = old.Name, milestoneLabel(old)
			}
			ms, err := milestone.Rename(s, refs, from, to)
			var ee *milestone.ExistsError
			if errors.As(err, &ee) {
				return fail(ExitValidation, "milestone_exists",
					"milestone %q already exists (%s%s) — a rename needs a name no milestone on the board has", to, milestone.RefPrefix, ee.ID)
			}
			if err := editRefusal(s, from, err, ", and then it can be renamed"); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				out := milestoneJSON(ms, nil)
				out["renamed_from"] = fromName
				return emit(w, out)
			}
			if to == fromName {
				fmt.Fprintf(w, "Nothing to do: milestone %q already has that name.\n", fromName)
				return nil
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[1], to)
			}
			fmt.Fprintf(w, "Milestone %q renamed to %q: %s\n", fromLabel, to, ms.Ref())
			return nil
		},
	}
}

// parseMilestoneNames is milestone.ParseNames with the CLI's exit code. The
// empty string is the empty list, which --after reads as "a swimlane of its
// own".
func parseMilestoneNames(raw string) ([]string, error) {
	names, err := milestone.ParseNames(raw)
	if err != nil {
		return nil, fail(ExitUsage, "bad_milestone", "%v", err)
	}
	return names, nil
}

// resolvePredecessors reads what --after was given as IDs and checks them
// against the board, with the CLI's exit codes: an unknown predecessor is not
// found (5), a name two milestones carry or an arrow that closes a loop is a
// validation error (3).
func resolvePredecessors(all []*milestone.Milestone, id string, preds []string) ([]string, error) {
	ids, err := milestone.ResolveRefs(all, id, preds)
	var amb *milestone.AmbiguousError
	if errors.As(err, &amb) {
		return nil, fail(ExitValidation, "milestone_ambiguous", "%v", amb)
	}
	var pe *milestone.PredecessorError
	if !errors.As(err, &pe) {
		return ids, err
	}
	if pe.Cycle != nil {
		return nil, fail(ExitValidation, "milestone_cycle", "%v", pe)
	}
	return nil, fail(ExitNotFound, "no_such_milestone",
		"%v — 'jaira milestone ls' lists them, 'jaira milestone create %s' starts it", pe, pe.Pred)
}

// editMembers is add and rm, which differ only in which way the lines move.
// One implementation because the locking, the resolving and the reporting are
// the whole of both.
func editMembers(cmd *cobra.Command, args []string, add bool) error {
	s, err := openStore()
	if err != nil {
		return err
	}
	importMilestones(s)
	name, err := milestoneKey(args[0])
	if err != nil {
		return err
	}
	ids, err := resolveAll(s, args[1:])
	if err != nil {
		return err
	}

	var touched []string
	ms, err := milestone.Edit(s, refs, name, func(ms *milestone.Milestone) error {
		for _, id := range ids {
			changed := false
			if add {
				changed = ms.Add(id)
			} else {
				changed = ms.Remove(id)
			}
			if changed {
				touched = append(touched, id)
			}
		}
		return nil
	})
	if err := editRefusal(s, name, err, ", and then tickets can go in and out of it again"); err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	if g.jsonOut {
		out := milestoneJSON(ms, nil)
		out["changed"] = strOrEmpty(touched)
		return emit(w, out)
	}
	label := milestoneLabel(ms)
	verb := "added to"
	if !add {
		verb = "removed from"
	}
	if len(touched) == 0 {
		fmt.Fprintf(w, "Nothing to do: milestone %q is already as asked.\n", label)
		return nil
	}
	for _, id := range touched {
		fmt.Fprintf(w, "%s %s %q\n", ticket.Handle(id), verb, label)
	}
	fmt.Fprintf(w, "%s now holds %d ticket(s).\n", label, len(ms.Members()))
	return nil
}

func newMilestoneLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List this board's milestones",
		Long: `Lists every milestone on the board with its colour and how many tickets it
holds.

Read it before creating one, for the same reason 'jaira tags' is read before
tagging: "q4" and "quarter-four" on one board are two names for one round of
work and each filters to half of it.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			importMilestones(s)
			all, err := milestone.LoadAll(refs)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(all))
				for _, ms := range all {
					arr = append(arr, milestoneJSON(ms, all))
				}
				return emit(w, map[string]any{
					"milestones": arr, "count": len(arr), "refs": milestone.RefPrefix,
					"milestones_twins": twinsJSON(milestone.Twins(all)),
				})
			}
			if len(all) == 0 {
				fmt.Fprintf(w, "No milestones on this board yet. 'jaira milestone create <name>' starts one.\n")
				return nil
			}
			colours := colourable(w)
			labels := milestone.Labels(all)
			fmt.Fprintln(w)
			for _, ms := range all {
				known := ms.HasColour()
				value := "  -"
				if known {
					value = fmt.Sprintf("%3d", ms.Colour)
				}
				goal := ""
				if ms.Goal != "" {
					goal = "  " + truncate(ms.Goal, 50)
				}
				fmt.Fprintf(w, "  %s %-24s %s  %3d ticket(s)%s\n",
					swatch(ms.Colour, known, colours), labels[ms.ID], value, len(ms.Members()), goal)
			}
			fmt.Fprintf(w, "\nRefs: %s<ID> — 'git show <ref>:<name>-<ID>.md' reads one\n", milestone.RefPrefix)
			printTwins(milestone.Twins(all))
			return nil
		},
	}
}

// These are the refusals a filed milestone raises, each written in one place.
// create, add/rm/set/rename and logbook all stand in front of the same
// milestone and owe the reader the same facts; the last clause is what
// differs, so the caller passes it: it says what THIS command would have done
// instead. The refusal is a validation error, so a script sees exit 3 and a
// reason it can branch on.

// refuseFiled is a milestone whose ref says it is filed. The ref is the same
// in every clone, so the way back runs anywhere: 'jaira restore <name>-<ID>'.
func refuseFiled(name, id, instead string) error {
	return fail(ExitValidation, "milestone_filed",
		"milestone %q has been filed: %s%s is marked %q, which is what keeps it off the board — 'jaira restore %s-%s' puts it back, and %s",
		name, milestone.RefPrefix, id, milestone.StatusFiled, name, id, instead)
}

// refuseFiledInLogbook is a milestone an older build filed into this tree's
// logbook as a file, whose ref this clone does not have. `where` is the
// logbook folder the copy lies in, and `then` is the caller's closing clause,
// punctuation included — it follows the colour with no separator of its own.
func refuseFiledInLogbook(name, file, where, then string) error {
	return fail(ExitValidation, "milestone_filed",
		"milestone %q has been filed into the logbook (%s) — 'jaira restore %s' brings it back with its ticket list and its colour%s",
		name, where, file, then)
}

// refuseIfFiled asks whether a name nothing on the refs carries was filed into
// this tree's logbook as a file by an older build, and hands back the refusal
// that owes the reader, or nil when it was not and the caller may carry on.
// `inLogbook` is the caller's closing clause, punctuation included.
func refuseIfFiled(s *ticket.Store, name, inLogbook string) error {
	if where, file, ok := s.FiledMilestone(milestone.FileNamed(name)); ok {
		return refuseFiledInLogbook(name, file, where, inLogbook)
	}
	return nil
}

// reportConflicts tells, on stderr in the shape the caller asked for, the
// fields two sides of a milestone both changed, where the ref's value was kept.
func reportConflicts(conflicts []milestone.Conflict) {
	for _, c := range conflicts {
		warnRef(map[string]any{"milestones_conflicts": []milestone.Conflict{c}}, "jaira: %s", c)
	}
}

// importMilestones takes the milestone files an older build left under
// .jaira/milestones/ onto their refs and deletes them, once, for every
// milestone command and the list path: from then on nothing reads that
// directory. One readdir when there is nothing to do. Where the board has a
// remote it fetches first, so "no ref here" can be told to mean "none on the
// remote" (refsync.Syncer.ImportMilestones); a file that still cannot be told
// apart waits, and says so.
//
// Best effort: what could not be taken stays on disk for the next command and
// is reported, and the command itself goes on.
func importMilestones(s *ticket.Store) {
	if !milestone.Pending(s.Root) {
		return
	}
	refs.FetchRefs()
	unlock, err := s.Lock(milestoneLockName)
	if err != nil {
		warnRef(map[string]any{"error": err.Error()}, "jaira: warning: milestone files were not taken onto their refs: %v", err)
		return
	}
	defer unlock()
	res, err := refs.ImportMilestones(s.Root)
	reportImported(res.Moved, res.Waiting)
	reportConflicts(res.Conflicts)
	if err != nil {
		warnRef(map[string]any{"error": err.Error()}, "jaira: warning: a milestone file could not be taken onto its ref and stays under .jaira/milestones/: %v", err)
	}
}

// reportImported tells what an import took off the disk and what it left. On
// stderr, in the shape the caller asked for, like every other warning a
// command carries: the deletion is a write nobody asked for by name.
func reportImported(moved, waiting []string) {
	if len(moved) > 0 {
		warnRef(map[string]any{"milestones_migrated": moved},
			"jaira: milestones live on refs now; took %d file(s) off the disk: %s", len(moved), strings.Join(moved, ", "))
	}
	if len(waiting) > 0 {
		warnRef(map[string]any{"milestones_waiting": waiting},
			"jaira: %d milestone file(s) stay under .jaira/milestones/ until a 'jaira fetch' reaches the remote: %s", len(waiting), strings.Join(waiting, ", "))
	}
}

// twinsJSON is the twins as --json gives them: the name, the labels that tell
// its milestones apart, and the command that renames one.
func twinsJSON(twins []milestone.Twin) []map[string]any {
	out := make([]map[string]any, 0, len(twins))
	for _, t := range twins {
		out = append(out, map[string]any{"name": t.Name, "labels": t.Labels, "rename": t.RenameCommand()})
	}
	return out
}

// printTwins is the warning for each name two milestones share, on stderr.
func printTwins(twins []milestone.Twin) {
	for _, t := range twins {
		fmt.Fprintf(os.Stderr, "jaira: warning: %d milestones are named %q: %s - rename one: %s\n",
			len(t.Labels), t.Name, strings.Join(t.Labels, ", "), t.RenameCommand())
	}
}

// reportTwins is the warning a command that has no JSON field of its own
// gives, on stderr in the shape the caller asked for.
func reportTwins() {
	all, err := milestone.LoadAll(refs)
	if err != nil {
		return
	}
	twins := milestone.Twins(all)
	if g.jsonOut {
		for _, j := range twinsJSON(twins) {
			warnRef(map[string]any{"milestones_twins": []map[string]any{j}}, "")
		}
		return
	}
	printTwins(twins)
}

// editRefusal words what milestone.Edit refused for a command that changes a
// milestone. inLogbook is the command's closing clause for a name no ref
// carries, as refuseIfFiled takes it.
func editRefusal(s *ticket.Store, name string, err error, inLogbook string) error {
	var fe *milestone.FiledError
	var amb *milestone.AmbiguousError
	switch {
	case err == nil:
		return nil
	case errors.Is(err, refsync.ErrNoMilestoneRefs):
		return refuseNoRefs()
	case errors.As(err, &amb):
		return refuseAmbiguous(amb)
	case errors.As(err, &fe):
		// The write would put the change on the ref, so every clone that
		// fetches sees a milestone the filer took off the board being
		// edited. Refuse it the way logbook does.
		return refuseFiled(fe.Name, fe.ID,
			"editing it instead would push the change onto its ref for everyone who fetches")
	case errors.Is(err, os.ErrNotExist):
		// Filed into this tree's logbook by an older build, whose ref this
		// clone never had: Load answers the same ErrNotExist a name nobody
		// ever used answers with, so point at restore first.
		if err := refuseIfFiled(s, name, inLogbook); err != nil {
			return err
		}
		return fail(ExitNotFound, "no_such_milestone",
			"no milestone %q on this board — 'jaira milestone ls' lists them, 'jaira milestone create %s' starts it",
			name, name)
	default:
		return err
	}
}

// resolveAll turns handles, prefixes and full ids into full ids, refusing the
// whole call if one of them names nothing. A milestone line pointing at no
// ticket is a dead line that looks like a plan.
func resolveAll(s *ticket.Store, args []string) ([]string, error) {
	var out []string
	for _, a := range args {
		id, err := resolveRef(s, a)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

// milestoneJSON is one milestone as --json gives it. all is the board's
// milestones when the caller has read them already; nil reads them only when
// a predecessor has to be spelled.
func milestoneJSON(ms *milestone.Milestone, all []*milestone.Milestone) map[string]any {
	members := memberRefs(ms.Members())
	out := map[string]any{
		"id":      ms.ID,
		"name":    ms.Name,
		"color":   ms.Colour,
		"tickets": members,
		"count":   len(members),
		"ref":     ms.Ref(),
		"file":    ms.FileName(),
	}
	if ms.Goal != "" {
		out["goal"] = ms.Goal
	}
	if ms.Placed {
		// Present only when placed, and then even when empty: [] is "top of a
		// swimlane of its own", which a reader must tell from "not in the layer".
		out["after"] = afterIDs(ms, all)
	}
	if !ms.CreatedAt.IsZero() {
		out["created_at"] = ms.CreatedAt.UTC().Format(time.RFC3339)
	}
	return out
}

// memberRefs is a member list as --json prints it, for milestones and epics
// alike: the full id and the handle a person types back.
func memberRefs(ids []string) []map[string]string {
	out := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, map[string]string{"id": id, "handle": ticket.Handle(id)})
	}
	return out
}

// afterIDs is a milestone's predecessors as IDs. A token a hand wrote as a
// name is spelled as the ID it names, and one that names nothing as it stands.
func afterIDs(ms *milestone.Milestone, all []*milestone.Milestone) []string {
	if len(ms.After) == 0 {
		return []string{}
	}
	if all == nil {
		all, _ = milestone.LoadAll(refs)
	}
	return milestone.SpellAsIDs(all, ms.After)
}

// milestoneLabel is what a message calls ms: its name, or <name>-<ID> where two
// share the name, as everywhere else.
func milestoneLabel(ms *milestone.Milestone) string {
	if all, err := milestone.LoadAll(refs); err == nil {
		if l, ok := milestone.Labels(all)[ms.ID]; ok {
			return l
		}
	}
	return ms.Name
}

// milestoneKey is a milestone argument as typed: an ID, a <name>-<ID> or a name,
// which milestone.Load reads in that order. It is refused when it could not be
// any of them; it is not lower-cased, because an ID is upper case.
func milestoneKey(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if _, _, err := milestone.NormalizeName(raw); err != nil {
		return "", fail(ExitUsage, "bad_milestone", "%v", err)
	}
	return raw, nil
}

// refuseNoRefs is a milestone write on a board outside git: a milestone lives
// on a ref and nowhere else, so there is nowhere to put it. The same refusal
// an epic gives (epicFail).
func refuseNoRefs() error {
	return fail(ExitValidation, "no_refs",
		"milestones live on git refs, and this board is not in a git repository — 'git init' (and 'jaira share' to send them to a remote) sets one up")
}

// refuseAmbiguous is a name two milestones carry: nothing was written, and the
// refusal names both <name>-<ID> forms to say again.
func refuseAmbiguous(amb *milestone.AmbiguousError) error {
	return fail(ExitValidation, "milestone_ambiguous",
		"%v — nothing was written; run it again with one of them", amb)
}
