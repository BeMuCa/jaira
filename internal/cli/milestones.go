package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// milestoneLockName keys the store lock a milestone file is written under; the
// name lives in core/milestone, where the shared write path takes it, and the
// commands that still lock by hand (logbook, restore, fetch) read it from there.
const milestoneLockName = milestone.LockName

// milestoneIndex is the one place a milestone index is built for the CLI, so
// the board's filter and 'jaira list --milestone' cannot disagree about what
// belonging to a milestone means. The TUI builds the same index from the same
// call in its reload.
//
// A board with no milestones is the common case and costs one failed readdir,
// which is why this is called on the list path without a flag guarding it.
func milestoneIndex(root string) (milestone.Index, error) {
	all, err := milestone.LoadAll(root)
	if err != nil {
		return nil, err
	}
	return milestone.Build(all), nil
}

func newMilestoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "Group tickets into a round of work",
		Long: `A milestone is the set of tickets that belong to one round of work, kept in its
own file under ` + "`.jaira/milestones/<name>.md`" + `.

It is a file rather than a field on each ticket because of what happens at the
end of a round: the work that did not finish has to appear in the next one, and
a per-ticket marker means touching every ticket one at a time. Here you open
the next milestone's file and move the unfinished lines into it — one edit
instead of twenty.

The file is hand-editable and reads in a diff, like a ticket. The file name IS
the milestone's name — ` + "`jaira milestone rename <old> <new>`" + ` renames it, since
renaming the file by hand leaves its ref behind and the next fetch writes the
old file back; the frontmatter carries its colour and when it was created, and
below it one ticket id per line. jaira keeps every other line exactly as it
found it, so comments, blank lines and an order you chose all survive.

Each milestone is given a random free colour, which its cards then show as a
bar down their RIGHT edge — the left edge belongs to the tags. On the board, M
opens the picker and narrows everything to the milestones ticked in it.

A milestone leaves the board only when somebody says so: ` + "`jaira logbook <name>`" + `
files it once all its tickets are finished, and ` + "`jaira restore <name>.md`" + ` brings
it back. Emptying one with ` + "`rm`" + ` leaves it standing.`,
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
		Long: `Creates ` + "`.jaira/milestones/<name>.md`" + ` and puts any tickets named after the
name into it.

The colour is picked for you, at random from the colours no other milestone on
this board is using: concurrent milestones are never many, so a clash once the
palette is spent is survivable, and nobody should have to choose one. --color
<1-255> overrides it; 0 is not black here but "no colour", and would paint no
cell on any card.

A name that belongs to a filed milestone is refused rather than reused: its
ref still holds it, and a second milestone under one name is one identity that
looks different on every machine. ` + "`jaira restore <name>.md`" + ` brings the filed one
back instead. A name a teammate's milestone has on its ref is refused too, while
no ` + "`jaira fetch`" + ` has written it here yet: creating it would overwrite theirs.

Names are lowercase kebab, the same rule tags follow — the name is also the
filename, so it has to be safe in a path. "Round One" is filed as "round-one"
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

			ms, err := milestone.Create(s.Root, s, refs, func(existing []*milestone.Milestone) (*milestone.Milestone, error) {
				var ee *milestone.ExistsError
				var fe *milestone.FiledError
				switch err := milestone.CheckFree(s.Root, name); {
				case errors.As(err, &fe):
					// The file is here but marked, which a fetch of somebody
					// else's filing leaves behind. It is not on the board, so
					// "already exists" would send the reader looking for a
					// milestone no listing names.
					return nil, refuseFiledOnDisk(s.Root, name,
						"taking the line out here instead would leave that tree's copy stranded in its logbook")
				case errors.As(err, &ee):
					return nil, fail(ExitValidation, "milestone_exists",
						"milestone %q already exists at %s — 'jaira milestone add %s <id>' puts tickets in it",
						name, milestonePath(s.Root, name), name)
				case err != nil:
					return nil, err
				}
				// Filed under this name, here or in the tree that wrote the
				// ref: the name stays taken until 'jaira restore' gives it up.
				// Or a teammate's milestone no fetch has written here yet.
				if err := refuseIfTaken(s, name,
					"creating it again here would make a second milestone with the same name",
					"; creating it again would make a second milestone with the same name"); err != nil {
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
					if err := checkPredecessors(existing, ms.ID, preds); err != nil {
						return nil, err
					}
					ms.SetAfter(preds)
				}
				return ms, nil
			})
			if err := warnIfRecordOnly(err); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, milestoneJSON(ms, milestonePath(s.Root, name)))
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[0], name)
			}
			fmt.Fprintf(w, "Milestone %q created with %d ticket(s): %s\n",
				name, len(ms.Members()), milestonePath(s.Root, name))
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
		Long: `Adds one line per ticket to the milestone's file.

Every id given is added in one write, which is the whole point: grouping twenty
tickets is one edit of one file, not twenty edits of twenty ticket files each
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
		Long: `Removes one line per ticket from the milestone's file, in one write, leaving
every other line exactly where it was.

The milestone itself stays standing, even when its last ticket leaves. A
milestone goes away on command and never on its own — an empty one is still a
plan, and a milestone that deleted itself when emptied would take a freshly
created one with it the moment you changed your mind about its first ticket.

` + "`jaira logbook <name>`" + ` is the command that takes a milestone off the board,
and ` + "`jaira restore <name>.md`" + ` brings it back.`,
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
		Long: `Rewrites the ` + "`goal:`" + ` and ` + "`after:`" + ` lines of the milestone's file and nothing else.

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
			name, _, err := milestone.NormalizeName(args[0])
			if err != nil {
				return fail(ExitUsage, "bad_milestone", "%v", err)
			}
			var preds []string
			if setAfter {
				if preds, err = parseMilestoneNames(after); err != nil {
					return err
				}
			}

			var before []byte
			ms, err := milestone.Edit(s.Root, s, refs, name, func(ms *milestone.Milestone) error {
				before = ms.Bytes()
				if setAfter {
					all, err := milestone.LoadAll(s.Root)
					if err != nil {
						return err
					}
					if err := checkPredecessors(all, ms.ID, preds); err != nil {
						return err
					}
					ms.SetAfter(preds)
				}
				if noRoadmap {
					ms.Unplace()
				}
				if setGoal {
					ms.SetGoal(goal)
				}
				return nil
			})
			if err := editRefusal(s, name, err,
				"there is no copy of it here to change",
				", and then it can be changed again"); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, milestoneJSON(ms, milestonePath(s.Root, name)))
			}
			if string(before) == string(ms.Bytes()) {
				fmt.Fprintf(w, "Nothing to do: milestone %q is already as asked.\n", name)
				return nil
			}
			fmt.Fprintf(w, "Milestone %q updated: %s\n", name, milestonePath(s.Root, name))
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
		Long: `Renames the milestone everywhere its name is: the file becomes
` + "`.jaira/milestones/<new>.md`" + `, its ref ` + "`refs/jaira/milestones/<new>`" + `, and every
milestone that comes after it names <new> in its ` + "`after:`" + ` line. Tickets, goal,
colour and place on the Milestone Roadmap stay as they were.

The old name's ref stays up, marked ` + "`status: filed`" + ` and ` + "`renamed-to: <new>`" + `: a
clone that has the old file gets it marked on its next ` + "`jaira fetch`" + ` and stops
showing it, which a ref taken down would never tell it. The old name is free
again afterwards — there is nothing to restore — so a rename can be taken back
by renaming once more.

The new name follows the rules of ` + "`create`" + `: lowercase kebab, and you are told
when it was changed. A name a milestone on the board, a filed one, or a
teammate's milestone not fetched yet has is refused (exit 3); an old name that
is not on the board exits 5, or 3 when it was filed.`,
		Args: exactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			from, _, err := milestone.NormalizeName(args[0])
			if err != nil {
				return fail(ExitUsage, "bad_milestone", "%v", err)
			}
			to, changed, err := milestone.NormalizeName(args[1])
			if err != nil {
				return fail(ExitUsage, "bad_milestone", "%v", err)
			}

			const twice = "renaming another milestone to it would make a second milestone with the same name"
			ms, err := milestone.Rename(s.Root, s, refs, from, to, func(name string) error {
				return refuseIfTaken(s, name, twice, "; "+twice)
			})
			var ee *milestone.ExistsError
			var fe *milestone.FiledError
			switch {
			case errors.As(err, &ee):
				return fail(ExitValidation, "milestone_exists",
					"milestone %q already exists at %s — a rename needs a name no milestone has", to, milestonePath(s.Root, to))
			case errors.As(err, &fe) && fe.Name == to && to != from:
				return refuseFiledOnDisk(s.Root, to, twice)
			}
			if err := editRefusal(s, from, err,
				"there is no copy of it here to rename",
				", and then it can be renamed"); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				out := milestoneJSON(ms, milestonePath(s.Root, ms.Name))
				out["renamed_from"] = from
				return emit(w, out)
			}
			if to == from {
				fmt.Fprintf(w, "Nothing to do: milestone %q already has that name.\n", from)
				return nil
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[1], to)
			}
			fmt.Fprintf(w, "Milestone %q renamed to %q: %s\n", from, to, milestonePath(s.Root, to))
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

// checkPredecessors is milestone.CheckAfter with the CLI's exit codes: an
// unknown predecessor is not found (5), an arrow that closes a loop is a
// validation error (3).
func checkPredecessors(all []*milestone.Milestone, id string, preds []string) error {
	var pe *milestone.PredecessorError
	if err := milestone.CheckAfter(all, id, preds); !errors.As(err, &pe) {
		return err
	}
	if pe.Cycle != nil {
		return fail(ExitValidation, "milestone_cycle", "%v", pe)
	}
	return fail(ExitNotFound, "no_such_milestone",
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
	name, _, err := milestone.NormalizeName(args[0])
	if err != nil {
		return fail(ExitUsage, "bad_milestone", "%v", err)
	}
	ids, err := resolveAll(s, args[1:])
	if err != nil {
		return err
	}

	var touched []string
	ms, err := milestone.Edit(s.Root, s, refs, name, func(ms *milestone.Milestone) error {
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
	if err := editRefusal(s, name, err,
		"there is no copy of it here for tickets to go in and out of",
		", and then tickets can go in and out of it again"); err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	if g.jsonOut {
		out := milestoneJSON(ms, milestonePath(s.Root, name))
		out["changed"] = strOrEmpty(touched)
		return emit(w, out)
	}
	verb := "added to"
	if !add {
		verb = "removed from"
	}
	if len(touched) == 0 {
		fmt.Fprintf(w, "Nothing to do: milestone %q is already as asked.\n", name)
		return nil
	}
	for _, id := range touched {
		fmt.Fprintf(w, "%s %s %q\n", ticket.Handle(id), verb, name)
	}
	fmt.Fprintf(w, "%s now holds %d ticket(s).\n", name, len(ms.Members()))
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
			all, err := milestone.LoadAll(s.Root)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(all))
				for _, ms := range all {
					arr = append(arr, milestoneJSON(ms, milestonePath(s.Root, ms.Name)))
				}
				return emit(w, map[string]any{
					"milestones": arr, "count": len(arr), "dir": milestone.Dir(s.Root),
				})
			}
			if len(all) == 0 {
				fmt.Fprintf(w, "No milestones on this board yet. 'jaira milestone create <name>' starts one.\n")
				return nil
			}
			colours := colourable(w)
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
					swatch(ms.Colour, known, colours), ms.Name, value, len(ms.Members()), goal)
			}
			fmt.Fprintf(w, "\nFiles: %s\n", milestone.Dir(s.Root))
			return nil
		},
	}
}

// These are the refusals a filed milestone raises, each written in one place.
// Which one a reader gets depends on the state they are in and not on the
// command they typed: create, add/rm and logbook all stand in front of the
// same milestone and owe the reader the same facts. The last clause is what
// differs, so the caller passes it: it says what THIS command would have done
// instead.
//
// Written once because they were written four times, and a refusal that drifts
// between doors teaches the reader that the doors are different problems.
// Modelled on refusePull (pull.go) — the refusal is a validation error, so a
// script sees exit 3 and a reason it can branch on.
//
// The split that matters is WHERE THE WAY BACK RUNS. refuseFiledElsewhere is
// for a filing that happened in another tree — the copy 'jaira restore' needs
// is over there, whether the mark reached us as a marked file or as a marked
// ref — and refuseFiledInLogbook is the one case where the copy is here.
// Naming the wrong one sends the reader to a restore that answers "not in the
// archive".
//
// `carries` names the thing that holds the mark, in the reader's terms.
func refuseFiledElsewhere(name, carries, instead string) error {
	return fail(ExitValidation, "milestone_filed",
		"milestone %q has been filed: %s is marked %q, which is what keeps it off the board — 'jaira restore %s.md' in the tree that filed it puts it back, and %s",
		name, carries, milestone.StatusFiled, name, instead)
}

// refuseFiledOnDisk: the file is lying here and carries the mark, which is
// what a fetch of somebody else's filing leaves behind.
func refuseFiledOnDisk(root, name, instead string) error {
	return refuseFiledElsewhere(name, "its file at "+milestonePath(root, name), instead)
}

// refuseRenamed is the door of a name a rename left behind: its file lies
// here marked, put there by a fetch of somebody's rename. Nothing was filed
// and 'jaira restore' brings nothing back — the milestone is where the mark
// points, so that is the whole answer.
func refuseRenamed(name, to string) error {
	return fail(ExitNotFound, "milestone_renamed",
		"milestone %q was renamed to %q — use that name; 'jaira milestone ls' lists them", name, to)
}

// refuseFiledInLogbook is the other state: the file is not on disk any more
// because this very tree filed it, so the way back runs here. `where` is the
// logbook folder the copy lies in, and `then` is the caller's closing clause,
// punctuation included — it follows the colour with no separator of its own.
func refuseFiledInLogbook(name, where, then string) error {
	return fail(ExitValidation, "milestone_filed",
		"milestone %q has been filed into the logbook (%s) — 'jaira restore %s.md' brings it back with its ticket list and its colour%s",
		name, where, name, then)
}

// refuseIfFiled is the door itself: it asks where the name is filed and hands
// back the refusal that state owes the reader, or nil when the name is not
// filed at all and the caller may carry on. Every door into a filed milestone
// — create, add/rm and 'jaira logbook <name>' — stands in front of the same
// two states and differs only in its closing clause, so each one passes the
// two clauses and is a single `if err := ...; err != nil`.
//
// Written here because the state→refusal mapping was written out at each door
// instead: the refusal TEXT was already shared, but a fourth door had to
// rediscover which state maps to which helper, and a state added later would
// have to be added at every door or be silently dropped at one.
//
// `onRef` and `inLogbook` are the caller's closing clauses for the two states,
// punctuation included, exactly as refuseFiledElsewhere and
// refuseFiledInLogbook take them — the wording differs because the way back
// differs.
func refuseIfFiled(s *ticket.Store, name, onRef, inLogbook string) error {
	switch where, at := refs.MilestoneFiled(s, name); at {
	case refsync.MilestoneFiledOnRef:
		// Nothing is filed under this name in this tree's logbook and the mark
		// is on the ref, so the copy 'jaira restore' needs lies in whichever
		// tree filed it. The ref is what `carries` names because, with no file
		// and no logbook entry here, it is the only thing the reader can go
		// and look at.
		return refuseFiledElsewhere(name, "its ref "+where, onRef)
	case refsync.MilestoneFiledHere:
		return refuseFiledInLogbook(name, where, inLogbook)
	}
	return nil
}

// refuseIfTaken is refuseIfFiled for a door that hands out a name, plus the
// one state only such a door has to fear: a teammate's milestone on its ref
// that no fetch has written to this disk yet (refsync MilestoneOnRef).
// Writing the name would overwrite theirs on the remote.
func refuseIfTaken(s *ticket.Store, name, onRef, inLogbook string) error {
	if err := refuseIfFiled(s, name, onRef, inLogbook); err != nil {
		return err
	}
	if refs.MilestoneOnRef(name) {
		return fail(ExitValidation, "milestone_exists",
			"milestone %q exists on its ref %s but not on this disk yet — 'jaira fetch' brings it here; writing the name now would overwrite it",
			name, gitref.MilestoneRefName(name))
	}
	return nil
}

// recordMilestone queues the file for its ref for the commands that write a
// milestone themselves (logbook, restore, fetch); create, add and rm go through
// milestone.Create and milestone.Edit, which record it. Best effort, see
// warnIfRecordOnly.
func recordMilestone(ms *milestone.Milestone) {
	if err := refs.RecordMilestone(ms.Name, ms.Name+".md", ms.Bytes()); err != nil {
		warnIfRecordOnly(&milestone.RecordError{Name: ms.Name, Err: err})
	}
}

// warnIfRecordOnly turns the one failure the write path reports beside a
// successful save — the ref could not be queued — into the warning it has
// always been, and passes every other error on. Best effort, and deliberately:
// a board with no repository or no remote is a supported way to use jaira, and
// refusing to group tickets because there is nowhere to send the grouping would
// break the local case to serve the shared one. The queue is never pushed
// here: the write path must not wait for a network.
func warnIfRecordOnly(err error) error {
	var re *milestone.RecordError
	if !errors.As(err, &re) {
		return err
	}
	warnRef(map[string]any{"milestone": re.Name, "error": re.Err.Error()},
		"jaira: warning: milestone %q was written here but could not be queued for the remote: %v", re.Name, re.Err)
	return nil
}

// editRefusal words what milestone.Edit refused for a command that changes a
// milestone, and downgrades a ref that could not be queued to its warning.
// onRef and inLogbook are the command's closing clauses for a name no file
// carries, as refuseIfFiled takes them.
func editRefusal(s *ticket.Store, name string, err error, onRef, inLogbook string) error {
	var fe *milestone.FiledError
	switch err := warnIfRecordOnly(err); {
	case err == nil:
		return nil
	case errors.As(err, &fe) && fe.RenamedTo != "":
		return refuseRenamed(name, fe.RenamedTo)
	case errors.As(err, &fe):
		// The third door into a filed milestone's file, after create and
		// logbook. The write would not lose the mark, but it does not stay
		// local either: the ref write would put it on the ref, so every
		// clone that fetches sees a milestone the filer took off the board
		// being edited. Refuse it where the other two refuse it.
		return refuseFiledOnDisk(s.Root, name,
			"editing it here instead would push the change onto its ref for everyone who fetches")
	case errors.Is(err, os.ErrNotExist):
		// Filed in this very tree: the file moved to the logbook, so Load
		// answers the same ErrNotExist a name nobody ever used answers
		// with. Saying "create it" here sends the reader into the
		// refusal create raises for a filed name, two steps for one
		// answer, and the first one points the wrong way.
		if err := refuseIfFiled(s, name, onRef, inLogbook); err != nil {
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

func milestoneJSON(ms *milestone.Milestone, path string) map[string]any {
	ids := ms.Members()
	members := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		members = append(members, map[string]string{"id": id, "handle": ticket.Handle(id)})
	}
	out := map[string]any{
		"name":    ms.Name,
		"color":   ms.Colour,
		"tickets": members,
		"count":   len(members),
		"file":    path,
	}
	if ms.Goal != "" {
		out["goal"] = ms.Goal
	}
	if ms.Placed {
		// Present only when placed, and then even when empty: [] is "top of a
		// swimlane of its own", which a reader must tell from "not in the layer".
		out["after"] = append([]string{}, ms.After...)
	}
	if !ms.CreatedAt.IsZero() {
		out["created_at"] = ms.CreatedAt.UTC().Format(time.RFC3339)
	}
	return out
}

// milestonePath is the file a milestone lives in, found by whatever names it:
// the file's name carries an ID that the argument may not. A milestone nobody
// has written yet gets the name-only path, the place the old format kept it.
func milestonePath(root, key string) string {
	if ms, err := milestone.Load(root, key); err == nil {
		return ms.Path(root)
	}
	return filepath.Join(milestone.Dir(root), key+".md")
}
