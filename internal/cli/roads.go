package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/road"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newRoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "road",
		Short: "Name the swimlanes of the milestone roadmap and list their milestones",
		Long: `A road is one swimlane of the milestone roadmap: its title, a description and
the milestones listed in it. The milestones in a road are drawn in its lane;
what is built inside the lane still comes from 'after:' on each milestone.

A milestone does not know its road. The list lives on the road, so a milestone
placed after a listed one is drawn in that lane without the road being written.
A milestone can be listed in one road only: adding it to a second is refused
(exit 3) and names the road it is in. A git merge can still list one in two;
'jaira validate' reports it, and the older road draws it.

A road lives on its own git ref, refs/jaira/roads/<name>, and nowhere else,
like an epic. A lane nobody has named is called Road-1, Road-2 and so on, in the
order the roadmap draws them, until somebody names it with 'jaira road set' or
on the roadmap itself.

Roads need a repository with a remote; a board without one has none.`,
	}
	cmd.AddCommand(newRoadCreateCmd(), newRoadAddCmd(), newRoadRmCmd(), newRoadSetCmd(), newRoadLsCmd())
	return cmd
}

// roadRefs is where roads are kept, or nil when the board has no refs.
func roadRefs() road.Refs {
	if refs == nil {
		return nil
	}
	return refs
}

// roadFail turns the road package's errors into the exit codes and reasons a
// script branches on.
func roadFail(name string, err error) error {
	var taken *road.TakenError
	switch {
	case errors.Is(err, road.ErrNoRefs):
		return fail(ExitValidation, "no_refs", "%v — 'jaira share' sets one up", err)
	case errors.Is(err, road.ErrNotFound):
		return fail(ExitNotFound, "no_such_road",
			"no road %q — 'jaira road ls' lists them, 'jaira road create %s' starts it", name, name)
	case errors.Is(err, road.ErrExists):
		return fail(ExitValidation, "road_exists",
			"road %q exists already — 'jaira road add %s <milestone>...' lists milestones in it", name, name)
	case errors.As(err, &taken):
		return fail(ExitValidation, "milestone_in_road",
			"%v — 'jaira road rm %s %s' takes it out there first; nothing was written", taken, taken.Road, taken.Milestone)
	}
	return err
}

func roadName(raw string) (string, bool, error) {
	name, changed, err := road.NormalizeName(raw)
	if err != nil {
		return "", false, fail(ExitUsage, "bad_road", "%v", err)
	}
	return name, changed, nil
}

// roadMilestones resolves milestone arguments, each an ID, a name or a
// <name>-<ID>, to the IDs a road lists. An unknown one is not found (5) and a
// name two milestones carry is refused (3), as the milestone commands do.
func roadMilestones(s *ticket.Store, args []string) ([]string, error) {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		m, err := milestoneNamed(s, arg)
		var amb *milestone.AmbiguousError
		switch {
		case errors.As(err, &amb):
			return nil, refuseAmbiguous(amb)
		case errors.Is(err, fs.ErrNotExist):
			return nil, fail(ExitNotFound, "no_such_milestone",
				"no milestone %q — 'jaira milestone ls' lists them, 'jaira milestone create %s' starts it", arg, arg)
		case err != nil:
			return nil, fail(ExitUsage, "bad_milestone", "%v", err)
		}
		if !slices.Contains(out, m.ID) {
			out = append(out, m.ID)
		}
	}
	return out, nil
}

// roadMilestoneNames maps each milestone on the board to its name. A listed ID
// that is not there (a road can outlive a milestone file) has no entry.
func roadMilestoneNames(s *ticket.Store) map[string]string {
	names := map[string]string{}
	if all, err := milestone.LoadAll(s.Root); err == nil {
		for _, m := range all {
			names[m.ID] = m.Name
		}
	}
	return names
}

// roadMemberLabel is a listed milestone as a person types it: name-ID, or the
// bare ID when the milestone is not on the board.
func roadMemberLabel(names map[string]string, id string) string {
	if n := names[id]; n != "" {
		return n + "-" + id
	}
	return id
}

func newRoadCreateCmd() *cobra.Command {
	var title, description string
	cmd := &cobra.Command{
		Use:   "create <name> [milestone...]",
		Short: "Start a road, optionally with its first milestones",
		Long: `Writes refs/jaira/roads/<name> and lists any milestones named after the name in it.

The name is lowercase kebab, the tag rule; "Auth rollout" is stored as
"auth-rollout" and you are told so. --title is what the roadmap shows (the name
when left out), --description says in one line what the road is for.`,
		Args: minArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, changed, err := roadName(args[0])
			if err != nil {
				return err
			}
			ids, err := roadMilestones(s, args[1:])
			if err != nil {
				return err
			}
			x, err := road.Create(s, roadRefs(), name, title, description, ids)
			if err != nil {
				return roadFail(name, err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, roadJSON(s, x))
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[0], name)
			}
			fmt.Fprintf(w, "Road %s created with %d milestone(s) on %s\n", x.Label(), len(x.Members()), gitref.RoadRefName(name))
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "what the roadmap shows (default: the name)")
	cmd.Flags().StringVar(&description, "description", "", "one line on what the road is for")
	return cmd
}

func newRoadAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <milestone>...",
		Short: "List milestones in a road",
		Long: `Adds one line per milestone to the road, in one write of the road's ref. The
milestones themselves are not touched. A milestone the road lists already stays
where it is. A milestone another road lists is refused with exit 3.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editRoad(cmd, args, true) },
	}
}

func newRoadRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name> <milestone>...",
		Short: "Take milestones out of a road",
		Long: `Removes one line per milestone from the road, in one write of the road's ref,
leaving every other line where it was. The road stays standing when it is
emptied; the milestones are not touched.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editRoad(cmd, args, false) },
	}
}

func editRoad(cmd *cobra.Command, args []string, add bool) error {
	s, err := openStore()
	if err != nil {
		return err
	}
	name, _, err := roadName(args[0])
	if err != nil {
		return err
	}
	ids, err := roadMilestones(s, args[1:])
	if err != nil {
		return err
	}
	x, changed, err := road.Assign(s, roadRefs(), name, ids, add)
	if err != nil {
		return roadFail(name, err)
	}
	w := cmd.OutOrStdout()
	if g.jsonOut {
		out := roadJSON(s, x)
		out["changed"] = strOrEmpty(changed)
		return emit(w, out)
	}
	if len(changed) == 0 {
		fmt.Fprintf(w, "Nothing to do: road %s is already as asked.\n", x.Label())
		return nil
	}
	names := roadMilestoneNames(s)
	verb := "added to"
	if !add {
		verb = "removed from"
	}
	for _, id := range changed {
		fmt.Fprintf(w, "%s %s %s\n", roadMemberLabel(names, id), verb, x.Label())
	}
	fmt.Fprintf(w, "%s now lists %d milestone(s).\n", x.Label(), len(x.Members()))
	return nil
}

func newRoadSetCmd() *cobra.Command {
	var title, description string
	cmd := &cobra.Command{
		Use:   "set <name> [--title <t>] [--description <d>]",
		Short: "Rename a road or change what it says",
		Long: `Rewrites a road's title and description, in one write of its ref. The milestones
it lists are left as they are. A blank title is the name again.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, _, err := roadName(args[0])
			if err != nil {
				return err
			}
			var t, d *string
			if cmd.Flags().Changed("title") {
				t = &title
			}
			if cmd.Flags().Changed("description") {
				d = &description
			}
			if t == nil && d == nil {
				return fail(ExitUsage, "usage", "nothing to set: give --title or --description")
			}
			x, written, err := road.Set(s, roadRefs(), name, t, d)
			if err != nil {
				return roadFail(name, err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				out := roadJSON(s, x)
				out["changed"] = written
				return emit(w, out)
			}
			if !written {
				fmt.Fprintf(w, "Nothing to do: road %s is already as asked.\n", x.Label())
				return nil
			}
			fmt.Fprintf(w, "Road %s saved on %s\n", x.Label(), gitref.RoadRefName(name))
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "what the roadmap shows")
	cmd.Flags().StringVar(&description, "description", "", "one line on what the road is for")
	return cmd
}

func newRoadLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List roads and the milestones each one lists",
		Long: `Lists every road, oldest first, with the milestones it lists. --json gives
name, title, description, created_at and milestones (id and name) per road.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			all, err := road.LoadAll(roadRefs())
			if err != nil {
				return roadFail("", err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(all))
				for _, x := range all {
					arr = append(arr, roadJSON(s, x))
				}
				return emit(w, map[string]any{"roads": arr, "count": len(arr)})
			}
			if len(all) == 0 {
				fmt.Fprintln(w, "No roads on this board. 'jaira road create <name>' starts one.")
				return nil
			}
			names := roadMilestoneNames(s)
			width := 0
			for _, x := range all {
				width = max(width, utf8.RuneCountInString(x.Label()))
			}
			fmt.Fprintln(w)
			for _, x := range all {
				var listed []string
				for _, id := range x.Members() {
					listed = append(listed, roadMemberLabel(names, id))
				}
				label := x.Label()
				fmt.Fprintf(w, "  %s%s  %d milestone(s)%s\n", label, strings.Repeat(" ", width-utf8.RuneCountInString(label)), len(listed), listSuffix(listed))
			}
			fmt.Fprintln(w)
			return nil
		},
	}
}

// listSuffix is ": a, b" for a non-empty list, and nothing for an empty one.
func listSuffix(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return ": " + strings.Join(items, ", ")
}

// roadJSON is one road as the commands print it.
func roadJSON(s *ticket.Store, x *road.Road) map[string]any {
	names := roadMilestoneNames(s)
	ms := make([]map[string]string, 0, len(x.Members()))
	for _, id := range x.Members() {
		ms = append(ms, map[string]string{"id": id, "name": names[id]})
	}
	out := map[string]any{
		"name":        x.Name,
		"title":       x.Title,
		"description": x.Description,
		"milestones":  ms,
		"ref":         gitref.RoadRefName(x.Name),
	}
	if !x.CreatedAt.IsZero() {
		out["created_at"] = x.CreatedAt.UTC().Format(time.RFC3339)
	}
	return out
}
