package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// milestoneLockName keys the store lock a milestone file is written under. One
// name for the whole directory rather than one per file: creating a milestone
// reads every existing one to avoid their colours, so two creations at once
// have to be ordered even though they write different files.
const milestoneLockName = "milestones"

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
the milestone's name — rename the file to rename the milestone; the frontmatter
carries its colour and when it was created, and below it one ticket id per
line. jaira keeps every other line exactly as it found it, so comments,
blank lines and an order you chose all survive.

Each milestone is given a random free colour, which its cards then show as a
bar down their RIGHT edge — the left edge belongs to the tags. On the board, M
opens the picker and narrows everything to one milestone.

A milestone leaves the board only when somebody says so: ` + "`jaira logbook <name>`" + `
files it once all its tickets are finished, and ` + "`jaira restore <name>.md`" + ` brings
it back. Emptying one with ` + "`rm`" + ` leaves it standing.`,
	}
	cmd.AddCommand(
		newMilestoneCreateCmd(),
		newMilestoneAddCmd(),
		newMilestoneRmCmd(),
		newMilestoneLsCmd(),
	)
	return cmd
}

func newMilestoneCreateCmd() *cobra.Command {
	var colour int
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
back instead.

Names are lowercase kebab, the same rule tags follow — the name is also the
filename, so it has to be safe in a path. "Round One" is filed as "round-one"
and you are told so.`,
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

			unlock, err := s.Lock(milestoneLockName)
			if err != nil {
				return err
			}
			defer unlock()

			if have, err := milestone.Load(s.Root, name); err == nil && have.Filed() {
				// The file is here but marked, which a fetch of somebody
				// else's filing leaves behind. It is not on the board, so
				// "already exists" would send the reader looking for a
				// milestone no listing names.
				return fail(ExitValidation, "milestone_filed",
					"milestone %q has been filed: its file at %s is marked %q, which is what keeps it off the board — 'jaira restore %s.md' in the tree that filed it puts it back, and taking the line out here instead would leave that tree's copy stranded in its logbook",
					name, milestone.Path(s.Root, name), milestone.StatusFiled, name)
			} else if err == nil {
				return fail(ExitValidation, "milestone_exists",
					"milestone %q already exists at %s — 'jaira milestone add %s <id>' puts tickets in it",
					name, milestone.Path(s.Root, name), name)
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			if where, filed := milestoneFiled(s, name); filed {
				return fail(ExitValidation, "milestone_filed",
					"milestone %q has been filed into the logbook (%s) — 'jaira restore %s.md' brings it back with its ticket list and its colour; creating it again would make a second milestone with the same name",
					name, where, name)
			}
			existing, err := milestone.LoadAll(s.Root)
			if err != nil {
				return err
			}
			c := colour
			if !cmd.Flags().Changed("color") {
				c = milestone.AssignColour(existing, name)
			} else if !tag.ValidColour(c) || c == 0 {
				// 0 is not black here, it is "no colour": a milestone given it
				// would paint nothing on any card, which is not a thing to be
				// asked for by hand.
				return fail(ExitUsage, "bad_color", "--color takes an ANSI-256 value, 1-255; got %d", c)
			}
			ms := milestone.New(name, c, time.Now())
			for _, id := range ids {
				ms.Add(id)
			}
			if err := ms.Save(s.Root); err != nil {
				return err
			}
			recordMilestone(s, ms)

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, milestoneJSON(ms, milestone.Path(s.Root, name)))
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[0], name)
			}
			fmt.Fprintf(w, "Milestone %q created with %d ticket(s): %s\n",
				name, len(ms.Members()), milestone.Path(s.Root, name))
			return nil
		},
	}
	cmd.Flags().IntVar(&colour, "color", -1, "ANSI-256 colour (1-255; 0 paints no cell) instead of a random free one")
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

	unlock, err := s.Lock(milestoneLockName)
	if err != nil {
		return err
	}
	defer unlock()

	ms, err := milestone.Load(s.Root, name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Filed in this very tree: the file moved to the logbook, so Load
			// answers the same ErrNotExist a name nobody ever used answers
			// with. Saying "create it" here sends the reader into the
			// refusal create raises for a filed name, two steps for one
			// answer, and the first one points the wrong way.
			if where, filed := milestoneFiled(s, name); filed {
				return fail(ExitValidation, "milestone_filed",
					"milestone %q has been filed into the logbook (%s) — 'jaira restore %s.md' brings it back with its ticket list and its colour, and then tickets can go in and out of it again",
					name, where, name)
			}
			return fail(ExitNotFound, "no_such_milestone",
				"no milestone %q on this board — 'jaira milestone ls' lists them, 'jaira milestone create %s' starts it",
				name, name)
		}
		return err
	}
	if ms.Filed() {
		// The third door into a filed milestone's file, after create and
		// logbook. The write would not lose the mark, but it does not stay
		// local either: recordMilestone below puts it on the ref, so every
		// clone that fetches sees a milestone the filer took off the board
		// being edited. Refuse it where the other two refuse it.
		return fail(ExitValidation, "milestone_filed",
			"milestone %q has been filed: its file at %s is marked %q, which is what keeps it off the board — 'jaira restore %s.md' in the tree that filed it puts it back, and editing it here instead would push the change onto its ref for everyone who fetches",
			name, milestone.Path(s.Root, name), milestone.StatusFiled, name)
	}
	var touched []string
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
	if len(touched) > 0 {
		if err := ms.Save(s.Root); err != nil {
			return err
		}
		recordMilestone(s, ms)
	}

	w := cmd.OutOrStdout()
	if g.jsonOut {
		out := milestoneJSON(ms, milestone.Path(s.Root, name))
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
					arr = append(arr, milestoneJSON(ms, milestone.Path(s.Root, ms.Name)))
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
				fmt.Fprintf(w, "  %s %-24s %s  %3d ticket(s)\n",
					swatch(ms.Colour, known, colours), ms.Name, value, len(ms.Members()))
			}
			fmt.Fprintf(w, "\nFiles: %s\n", milestone.Dir(s.Root))
			return nil
		},
	}
}

// milestoneFiled reports whether this name belongs to a milestone that has
// been filed, and where that was found. A filed name stays taken: two
// milestones called the same thing are one identity that looks different on
// every machine, which is the confusion a name exists to prevent.
//
// Two places are looked at because neither alone covers both boards. The ref
// is the one an unshared clone does not have, and the logbook is the one a
// clone that never fetched does not see — a board that has never been shared
// has no refs at all, and would otherwise hand out the same name twice.
func milestoneFiled(s *ticket.Store, name string) (string, bool) {
	if refs != nil && refs.Usable() == nil {
		if content, _, err := refs.Repo.ReadMilestone(name); err == nil {
			if milestone.FromBytes(name, content).Filed() {
				return "on its ref", true
			}
		}
	}
	return s.FiledMilestone(name)
}

// recordMilestone puts the milestone file on its own ref, the way every ticket
// write is put on the ticket's. Queued, never pushed here: the write path must
// not wait for a network, and an unsent milestone goes out with the next
// command like everything else in the outbox.
//
// Best effort, and deliberately: a board with no repository or no remote is a
// supported way to use jaira, and refusing to group tickets because there is
// nowhere to send the grouping would break the local case to serve the shared
// one.
func recordMilestone(s *ticket.Store, ms *milestone.Milestone) {
	content, err := os.ReadFile(milestone.Path(s.Root, ms.Name))
	if err != nil {
		return
	}
	if err := refs.RecordMilestone(ms.Name, content); err != nil {
		warnRef(map[string]any{"milestone": ms.Name, "error": err.Error()},
			"jaira: warning: milestone %q was written here but could not be queued for the remote: %v", ms.Name, err)
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
	if !ms.CreatedAt.IsZero() {
		out["created_at"] = ms.CreatedAt.UTC().Format(time.RFC3339)
	}
	return out
}
