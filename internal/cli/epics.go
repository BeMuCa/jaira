package cli

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newEpicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic",
		Short: "Group tickets into features and see how far each one is",
		Long: `An epic describes one feature for a person and collects the tickets that make
it up. It is not a ticket: it has no definition of done, runs through no lane
and no agent works it. Nor is it a milestone: a milestone is a round of work,
and what is unfinished moves on to the next; an epic is a feature, and what is
unfinished stays in it.

An epic lives on its own git ref, refs/jaira/epics/<name>, and nowhere else —
there is no file in the working tree. It reaches everybody on their next
fetch, without a branch being merged. Read one by hand with

  git show refs/jaira/epics/<name>:<name>.md

The member list is in the epic, never on the tickets: sorting a ticket into an
epic writes the epic's ref only, never the ref of a ticket somebody is working.

A ticket that follows on from a member (follows: or parent:, any number of
steps down) counts towards the epic without a line of its own, so follow-up
work that appears while a feature is built is never lost from it.

Epics need a repository with a remote; a board without one has none.`,
	}
	cmd.AddCommand(newEpicCreateCmd(), newEpicAddCmd(), newEpicRmCmd(), newEpicLsCmd(), newEpicShowCmd())
	return cmd
}

// epicFail turns the epic package's errors into the exit codes and reasons a
// script branches on.
func epicFail(name string, err error) error {
	switch {
	case errors.Is(err, epic.ErrNoRefs):
		return fail(ExitValidation, "no_refs", "%v — 'jaira share' sets one up", err)
	case errors.Is(err, epic.ErrNotFound):
		return fail(ExitNotFound, "no_such_epic",
			"no epic %q — 'jaira epic ls' lists them, 'jaira epic create %s' starts it", name, name)
	case errors.Is(err, epic.ErrExists):
		return fail(ExitValidation, "epic_exists",
			"epic %q exists already — 'jaira epic add %s <id>' puts tickets in it", name, name)
	}
	return err
}

func epicName(raw string) (string, bool, error) {
	name, changed, err := epic.NormalizeName(raw)
	if err != nil {
		return "", false, fail(ExitUsage, "bad_epic", "%v", err)
	}
	return name, changed, nil
}

func newEpicCreateCmd() *cobra.Command {
	var title, description string
	cmd := &cobra.Command{
		Use:   "create <name> [id...]",
		Short: "Start an epic, optionally with its first tickets",
		Long: `Writes refs/jaira/epics/<name> and puts any tickets named after the name in it.

The name is lowercase kebab, the tag rule; "Magic Login" is stored as
"magic-login" and you are told so. --title is what a person reads (the name
when left out), --description says in one line what the feature is for.`,
		Args: minArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, changed, err := epicName(args[0])
			if err != nil {
				return err
			}
			ids, err := resolveAll(s, args[1:])
			if err != nil {
				return err
			}
			e, err := epic.Create(s, epicRefs(), name, title, description, ids)
			if err != nil {
				return epicFail(name, err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, epicJSON(e, nil))
			}
			if changed {
				fmt.Fprintf(w, "Filed %q as %q.\n", args[0], name)
			}
			fmt.Fprintf(w, "Epic %q created with %d ticket(s) on %s\n", name, len(e.Members()), gitref.EpicRefName(name))
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "what a person reads (default: the name)")
	cmd.Flags().StringVar(&description, "description", "", "one line on what the feature is for")
	return cmd
}

func newEpicAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <id>...",
		Short: "Put tickets into an epic",
		Long: `Adds one line per ticket to the epic, in one write of the epic's ref. The
tickets themselves are not touched. A ticket already in the epic stays where
it is.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editEpic(cmd, args, true) },
	}
}

func newEpicRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name> <id>...",
		Short: "Take tickets out of an epic",
		Long: `Removes one line per ticket from the epic, in one write of the epic's ref,
leaving every other line where it was. The epic stays standing when it is
emptied.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editEpic(cmd, args, false) },
	}
}

func editEpic(cmd *cobra.Command, args []string, add bool) error {
	s, err := openStore()
	if err != nil {
		return err
	}
	name, _, err := epicName(args[0])
	if err != nil {
		return err
	}
	ids, err := resolveAll(s, args[1:])
	if err != nil {
		return err
	}
	e, changed, err := epic.Assign(s, epicRefs(), name, ids, add)
	if err != nil {
		return epicFail(name, err)
	}
	w := cmd.OutOrStdout()
	if g.jsonOut {
		out := epicJSON(e, nil)
		out["changed"] = strOrEmpty(changed)
		return emit(w, out)
	}
	if len(changed) == 0 {
		fmt.Fprintf(w, "Nothing to do: epic %q is already as asked.\n", name)
		return nil
	}
	verb := "added to"
	if !add {
		verb = "removed from"
	}
	for _, id := range changed {
		fmt.Fprintf(w, "%s %s %q\n", ticket.Handle(id), verb, name)
	}
	fmt.Fprintf(w, "%s now lists %d ticket(s).\n", name, len(e.Members()))
	return nil
}

// epicRefs is the store epics are read from and written to. A nil syncer is
// passed on as a nil interface, never as a typed nil, so the epic package
// sees "no refs" instead of calling a method on nothing.
func epicRefs() epic.Refs {
	if refs == nil {
		return nil
	}
	return refs
}

func newEpicLsCmd() *cobra.Command {
	var unassigned bool
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List epics with how far each one is",
		Long: `Lists every epic with done X of Y. Y counts the member tickets and every
ticket that follows on from one of them; X counts those in a terminal lane or
filed in the logbook.

--unassigned lists instead the open tickets that belong to no epic — neither by
a line nor by following on from a member — which is the list to work through
when sorting a board into features.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			all, err := epic.LoadAll(epicRefs())
			if err != nil {
				return epicFail("", err)
			}
			ix, board, _, err := boardIndex(s)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if unassigned {
				return printUnassigned(w, epic.Unassigned(all, ix, board))
			}
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(all))
				for _, e := range all {
					p := epic.Count(e, ix)
					arr = append(arr, epicJSON(e, &p))
				}
				return emit(w, map[string]any{"epics": arr, "count": len(arr)})
			}
			if len(all) == 0 {
				fmt.Fprintln(w, "No epics on this board yet. 'jaira epic create <name>' starts one.")
				return nil
			}
			fmt.Fprintln(w)
			for _, e := range all {
				p := epic.Count(e, ix)
				fmt.Fprintf(w, "  %-24s %3d of %-3d done  %s\n", e.Name, p.Done, p.Total, e.Title)
			}
			fmt.Fprintln(w)
			return nil
		},
	}
	cmd.Flags().BoolVar(&unassigned, "unassigned", false, "list the open tickets that are in no epic instead")
	return cmd
}

func printUnassigned(w io.Writer, open []*ticket.Ticket) error {
	if g.jsonOut {
		arr := make([]map[string]string, 0, len(open))
		for _, t := range open {
			arr = append(arr, map[string]string{"id": t.ID, "handle": ticket.Handle(t.ID), "title": t.Title, "lane": t.Status})
		}
		return emit(w, map[string]any{"unassigned": arr, "count": len(arr)})
	}
	if len(open) == 0 {
		fmt.Fprintln(w, "Every open ticket is in an epic.")
		return nil
	}
	fmt.Fprintf(w, "%d open ticket(s) in no epic:\n\n", len(open))
	for _, t := range open {
		fmt.Fprintf(w, "  %s  %-12s %s\n", ticket.Handle(t.ID), t.Status, t.Title)
	}
	return nil
}

func newEpicShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show one epic: what it is for and every ticket counted towards it",
		Args:  exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, _, err := epicName(args[0])
			if err != nil {
				return err
			}
			e, err := epic.Load(epicRefs(), name)
			if err != nil {
				return epicFail(name, err)
			}
			ix, _, _, err := boardIndex(s)
			if err != nil {
				return err
			}
			p := epic.Count(e, ix)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, epicJSON(e, &p))
			}
			fmt.Fprintf(w, "%s — %s\n", e.Name, e.Title)
			if e.Description != "" {
				fmt.Fprintf(w, "%s\n", e.Description)
			}
			fmt.Fprintf(w, "\n%d of %d done\n\n", p.Done, p.Total)
			for _, t := range p.Tickets {
				mark := " "
				if t.Done {
					mark = "x"
				}
				via := ""
				if !t.Direct {
					via = "  (follow-up)"
				}
				title := t.Title
				if title == "" {
					title = "(not found on this board or in the logbook)"
				}
				fmt.Fprintf(w, "  [%s] %s  %-12s %s%s\n", mark, ticket.Handle(t.ID), t.Lane, title, via)
			}
			return nil
		},
	}
}

func epicJSON(e *epic.Epic, p *epic.Progress) map[string]any {
	out := map[string]any{
		"name":        e.Name,
		"title":       e.Title,
		"description": e.Description,
		"members":     memberRefs(e.Members()),
		"ref":         gitref.EpicRefName(e.Name),
	}
	if !e.CreatedAt.IsZero() {
		out["created_at"] = e.CreatedAt.UTC().Format(time.RFC3339)
	}
	if p != nil {
		out["done"], out["total"], out["tickets"] = p.Done, p.Total, p.Tickets
	}
	return out
}
