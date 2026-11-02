package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

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

A finished epic is closed with 'jaira epic close', which writes "status:
closed" onto its ref and deletes nothing. A closed epic takes no new tickets
without --force and drops out of 'jaira epic ls'; 'jaira epic reopen' takes the
closing back.

An epic may carry the keys of the tasks it stands for in another system —
JIRA-1234, a YouTrack issue — each with a link: 'jaira epic set <name>
--external "JIRA-1234 <url>"'. They are on the epic's ref and never on a
ticket, so they never ride in a code commit. 'jaira epic ls --external <key>'
and 'jaira list --external <key>' find the work done for one of them. The ref
is pushed: everybody who can read the repository can read the keys.

Epics need a repository with a remote; a board without one has none.`,
	}
	cmd.AddCommand(newEpicCreateCmd(), newEpicAddCmd(), newEpicRmCmd(), newEpicSetCmd(), newEpicLsCmd(), newEpicShowCmd(), newEpicCloseCmd(), newEpicReopenCmd())
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
	case errors.Is(err, epic.ErrNoTitle):
		return fail(ExitUsage, "usage", "%v — leave --title out to keep the one it has", err)
	case errors.Is(err, epic.ErrClosed):
		return fail(ExitValidation, "epic_closed",
			"epic %s is closed — 'jaira epic add %s <id>... --force' puts the ticket in anyway and leaves it closed, 'jaira epic reopen %s' opens it again", epic.LabelOf(epicRefs(), name), name, name)
	case errors.Is(err, epic.ErrUnfinished):
		return fail(ExitValidation, "epic_unfinished",
			"%v — 'jaira epic show %s' lists them, 'jaira epic close %s --force' closes it anyway", err, name, name)
	case errors.Is(err, epic.ErrBadExternal):
		return fail(ExitUsage, "bad_external", "%v", err)
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

// externalFlags reads --external values, each "<KEY> [<URL>]".
func externalFlags(raw []string) ([]epic.External, error) {
	out := make([]epic.External, 0, len(raw))
	for _, r := range raw {
		x, err := epic.ParseExternal(r)
		if err != nil {
			return nil, epicFail("", err)
		}
		out = append(out, x)
	}
	return out, nil
}

// externalUsage is the one wording of --external on create and set.
const externalUsage = `the key of the task in another system with its link, "<KEY> [<URL>]"; repeat for several`

func newEpicCreateCmd() *cobra.Command {
	var title, description string
	var external []string
	cmd := &cobra.Command{
		Use:   "create <name> [id...]",
		Short: "Start an epic, optionally with its first tickets",
		Long: `Writes refs/jaira/epics/<name> and puts any tickets named after the name in it.

The name is lowercase kebab, the tag rule; "Magic Login" is stored as
"magic-login" and you are told so. --title is what a person reads (the name
when left out), --description says in one line what the feature is for.
--external "<KEY> <URL>" records the task in another system the epic stands
for; repeat it for several.`,
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
			keys, err := externalFlags(external)
			if err != nil {
				return err
			}
			e, err := epic.Create(s, epicRefs(), name, title, description, ids, keys...)
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
			fmt.Fprintf(w, "Epic %s created with %d ticket(s) on %s\n", e.Label(), len(e.Members()), gitref.EpicRefName(name))
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "what a person reads (default: the name)")
	cmd.Flags().StringVar(&description, "description", "", "one line on what the feature is for")
	cmd.Flags().StringArrayVar(&external, "external", nil, externalUsage)
	return cmd
}

func newEpicAddCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "add <name> <id>...",
		Short: "Put tickets into an epic",
		Long: `Adds one line per ticket to the epic, in one write of the epic's ref. The
tickets themselves are not touched. A ticket already in the epic stays where
it is.

A closed epic refuses a new ticket with exit 3. --force puts it in anyway and
leaves the epic closed; 'jaira epic reopen' is the way to open it again.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editEpic(cmd, args, true, force) },
	}
	cmd.Flags().BoolVar(&force, "force", false, "put tickets into a closed epic, which stays closed")
	return cmd
}

func newEpicRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name> <id>...",
		Short: "Take tickets out of an epic",
		Long: `Removes one line per ticket from the epic, in one write of the epic's ref,
leaving every other line where it was. The epic stays standing when it is
emptied.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { return editEpic(cmd, args, false, false) },
	}
}

func editEpic(cmd *cobra.Command, args []string, add, force bool) error {
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
	e, changed, err := epic.Assign(s, epicRefs(), name, ids, add, force)
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
		fmt.Fprintf(w, "Nothing to do: epic %s is already as asked.\n", e.Label())
		return nil
	}
	verb := "added to"
	if !add {
		verb = "removed from"
	}
	for _, id := range changed {
		fmt.Fprintf(w, "%s %s %s\n", ticket.Handle(id), verb, e.Label())
	}
	fmt.Fprintf(w, "%s now lists %d ticket(s).\n", name, len(e.Members()))
	return nil
}

func newEpicSetCmd() *cobra.Command {
	var title, description, body, bodyFile string
	var external, externalRm []string
	cmd := &cobra.Command{
		Use:   "set <name> [--title <t>] [--description <d>] [--body <text> | --body-file <path|->] [--external \"<KEY> <URL>\"]... [--external-rm <KEY>]...",
		Short: "Rewrite an epic's title, description, text or external keys",
		Long: `Rewrites what an epic says about itself, in one write of its ref. The name
stays — it is the ref — and the member lines stay as they are, in their order.

--title is what a person reads; a heading at the top of the text that carried
the old title is renamed with it. --description is the one line under the
title. --body replaces the markdown text below them, --body-file reads it from
a file, or from stdin with '-'. 'jaira epic show <name>' prints the text as it
stands.

--external "<KEY> <URL>" puts the key of a task in another system on the epic,
or gives a key it carries already the new link; --external-rm <KEY> takes one
off. Keys compare in any case. Both repeat.

Only what you name is touched. Asking for what is there already writes
nothing.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			if !f.Changed("title") && !f.Changed("description") && !f.Changed("body") && !f.Changed("body-file") &&
				!f.Changed("external") && !f.Changed("external-rm") {
				return fail(ExitUsage, "usage", "nothing to set — name at least one of --title, --description, --body, --body-file, --external, --external-rm")
			}
			if f.Changed("body") && f.Changed("body-file") {
				return fail(ExitUsage, "usage", "--body and --body-file both give the text; name one")
			}
			keys, err := externalFlags(external)
			if err != nil {
				return err
			}
			p := epic.Patch{AddExternal: keys, DropExternal: externalRm}
			if f.Changed("title") {
				p.Title = &title
			}
			if f.Changed("description") {
				p.Description = &description
			}
			if f.Changed("body") {
				p.Body = &body
			}
			if f.Changed("body-file") {
				var raw []byte
				var err error
				if bodyFile == "-" {
					raw, err = io.ReadAll(cmd.InOrStdin())
				} else {
					raw, err = os.ReadFile(bodyFile)
				}
				if err != nil {
					return fail(ExitUsage, "usage", "--body-file: %v", err)
				}
				text := string(raw)
				p.Body = &text
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			name, _, err := epicName(args[0])
			if err != nil {
				return err
			}
			e, changed, err := epic.Set(s, epicRefs(), name, p)
			if err != nil {
				return epicFail(name, err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				out := epicJSON(e, nil)
				out["changed"] = changed
				return emit(w, out)
			}
			if !changed {
				fmt.Fprintf(w, "Nothing to do: epic %s is already as asked.\n", e.Label())
				return nil
			}
			fmt.Fprintf(w, "Epic %s rewritten on %s\n", e.Label(), gitref.EpicRefName(name))
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "what a person reads")
	cmd.Flags().StringVar(&description, "description", "", "one line on what the feature is for")
	cmd.Flags().StringVar(&body, "body", "", "the markdown text below title and description")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "read the text from a file, '-' for stdin")
	cmd.Flags().StringArrayVar(&external, "external", nil, externalUsage)
	cmd.Flags().StringArrayVar(&externalRm, "external-rm", nil, "take the key of a task in another system off the epic; repeat for several")
	return cmd
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

func newEpicCloseCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "close <name>",
		Short: "Close a finished epic: it takes no new tickets and leaves 'epic ls'",
		Long: `Writes "status: closed" onto the epic's ref, in one write of it. Nothing is
deleted: the members, the text and the ref stay, every clone reads the closing
on its next fetch, and 'jaira epic reopen <name>' takes it back.

An epic that still counts an unfinished ticket — a member, or a ticket that
follows on from one — is refused with exit 3. --force closes it anyway.

A closed epic refuses 'jaira epic add' without --force and is left out of
'jaira epic ls'; 'jaira epic ls --all' shows it, marked.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, _, err := epicName(args[0])
			if err != nil {
				return err
			}
			ix, _, err := boardIndex(s)
			if err != nil {
				return err
			}
			e, changed, err := epic.Close(s, epicRefs(), name, ix, force)
			if err != nil {
				return epicFail(name, err)
			}
			return reportClosing(cmd.OutOrStdout(), e, changed)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "close it although tickets counted towards it are unfinished")
	return cmd
}

func newEpicReopenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reopen <name>",
		Short: "Take back the closing of an epic",
		Long: `Takes "status: closed" off the epic's ref again, in one write of it. The epic
takes new tickets again and is back in 'jaira epic ls'.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			name, _, err := epicName(args[0])
			if err != nil {
				return err
			}
			e, changed, err := epic.Reopen(s, epicRefs(), name)
			if err != nil {
				return epicFail(name, err)
			}
			return reportClosing(cmd.OutOrStdout(), e, changed)
		},
	}
}

func reportClosing(w io.Writer, e *epic.Epic, changed bool) error {
	if g.jsonOut {
		out := epicJSON(e, nil)
		out["changed"] = changed
		return emit(w, out)
	}
	if !changed {
		state := "open"
		if e.Closed {
			state = "closed"
		}
		fmt.Fprintf(w, "Nothing to do: epic %s is %s already.\n", e.Label(), state)
		return nil
	}
	if e.Closed {
		fmt.Fprintf(w, "Epic %s closed on %s — 'jaira epic reopen %s' takes it back.\n", e.Label(), gitref.EpicRefName(e.Name), e.Name)
		return nil
	}
	fmt.Fprintf(w, "Epic %s is open again on %s.\n", e.Label(), gitref.EpicRefName(e.Name))
	return nil
}

func newEpicLsCmd() *cobra.Command {
	var unassigned, everything bool
	var externalKey string
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List epics with how far each one is",
		Long: `Lists every epic with done X of Y. Y counts the member tickets and every
ticket that follows on from one of them; X counts those in a terminal lane or
filed in the logbook. A closed epic is left out; --all lists it too, marked
"closed".

--unassigned lists instead the open tickets that belong to no epic — neither by
a line nor by following on from a member — which is the list to work through
when sorting a board into features.

--external <key> lists only the epics that carry that key of a task in another
system, in any case; with --all, closed ones too.`,
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
			ix, board, err := boardIndex(s)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if unassigned {
				return printUnassigned(w, epic.Unassigned(all, ix, board))
			}
			if !everything {
				all = slices.DeleteFunc(all, func(e *epic.Epic) bool { return e.Closed })
			}
			if externalKey != "" {
				all = slices.DeleteFunc(all, func(e *epic.Epic) bool { return !e.HasExternal(externalKey) })
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
				if externalKey != "" {
					fmt.Fprintf(w, "No epic carries %q. 'jaira epic set <name> --external \"%s <url>\"' puts it on one.\n", externalKey, externalKey)
					return nil
				}
				if !everything {
					fmt.Fprintln(w, "No open epics on this board. 'jaira epic create <name>' starts one, 'jaira epic ls --all' shows closed ones.")
					return nil
				}
				fmt.Fprintln(w, "No epics on this board yet. 'jaira epic create <name>' starts one.")
				return nil
			}
			fmt.Fprintln(w)
			width := 0
			for _, e := range all {
				width = max(width, utf8.RuneCountInString(titleAndName(e)))
			}
			for _, e := range all {
				p := epic.Count(e, ix)
				label := titleAndName(e)
				fmt.Fprintf(w, "  %s%s  %3d of %-3d done%s%s\n", label, strings.Repeat(" ", width-utf8.RuneCountInString(label)), p.Done, p.Total, closedMark(e), keysMark(e))
			}
			fmt.Fprintln(w)
			return nil
		},
	}
	cmd.Flags().BoolVar(&unassigned, "unassigned", false, "list the open tickets that are in no epic instead")
	cmd.Flags().BoolVar(&everything, "all", false, "list closed epics too, marked")
	cmd.Flags().StringVar(&externalKey, "external", "", "only epics carrying this key of a task in another system")
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
			ix, _, err := boardIndex(s)
			if err != nil {
				return err
			}
			p := epic.Count(e, ix)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, epicJSON(e, &p))
			}
			fmt.Fprintf(w, "%s%s\n", titleAndName(e), closedMark(e))
			if e.Description != "" {
				fmt.Fprintf(w, "%s\n", e.Description)
			}
			for _, x := range e.Externals() {
				fmt.Fprintf(w, "external: %s\n", x)
			}
			if body := e.Text(); body != "" {
				fmt.Fprintf(w, "\n%s\n", body)
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

// titleAndName is an epic as a listing names it: the title a person reads,
// the name to type beside it — "Web 0.2  (web-0-2)". A title that is the name
// is printed once.
func titleAndName(e *epic.Epic) string {
	if e.Title == "" || e.Title == e.Name {
		return e.Name
	}
	return e.Title + "  (" + e.Name + ")"
}

// closedMark is what 'epic ls --all' and 'epic show' print after a closed
// epic's title.
func closedMark(e *epic.Epic) string {
	if e.Closed {
		return "  [closed]"
	}
	return ""
}

// keysMark is what 'epic ls' prints after a title: the keys the epic carries.
func keysMark(e *epic.Epic) string {
	var keys []string
	for _, x := range e.Externals() {
		keys = append(keys, x.Key)
	}
	if len(keys) == 0 {
		return ""
	}
	return "  [" + strings.Join(keys, ", ") + "]"
}

func epicJSON(e *epic.Epic, p *epic.Progress) map[string]any {
	out := map[string]any{
		"name":        e.Name,
		"title":       e.Title,
		"description": e.Description,
		"body":        e.Body(),
		"members":     memberRefs(e.Members()),
		"ref":         gitref.EpicRefName(e.Name),
		"closed":      e.Closed,
		"external":    append([]epic.External{}, e.Externals()...),
	}
	if !e.CreatedAt.IsZero() {
		out["created_at"] = e.CreatedAt.UTC().Format(time.RFC3339)
	}
	if p != nil {
		out["done"], out["total"], out["tickets"] = p.Done, p.Total, p.Tickets
	}
	return out
}

// ticketsForKey is every ticket counted towards an epic carrying the key:
// members and their follow-ups, as 'jaira epic show' counts them. Closed
// epics count too — the key asks what was done for a task, not what is open.
func ticketsForKey(s *ticket.Store, key string) (map[string]bool, error) {
	epics, err := epic.LoadAll(epicRefs())
	if err != nil {
		return nil, epicFail("", err)
	}
	ix, _, err := boardIndex(s)
	if err != nil {
		return nil, err
	}
	var members []string
	for _, e := range epics {
		if e.HasExternal(key) {
			members = append(members, e.Members()...)
		}
	}
	out := map[string]bool{}
	for _, id := range epic.Covered(members, ix) {
		out[id] = true
	}
	return out, nil
}
