package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/people"
)

// newPeopleCmd exists because one human used to appear on a board under every
// name their git configurations carried, and nothing could say that two of
// those names were one person. The list it writes is shared on a ref, so the
// answer is the same on every clone that fetches.
func newPeopleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "people",
		Short: "List who is who on this board: each person's jaira name and the names that mean them",
		Long: `Prints every person the board knows: the jaira name it shows and writes for
them, and every other name that means them — a git user.name from another
machine, a git user.email, an old name.

jaira acts as the person your git user.name or user.email is a name of, and
writes their jaira name into creator, assignee, claims and notes. Tickets
written under an older name are shown under the jaira name without being
rewritten. A name the list does not know is used as it is, as before.

The list lives on ` + gitref.PeopleRef + `, one file, one person a line:
"<jaira name>: <alias>, <alias>". Every write fetches it first and pushes at
once, so it needs the board's remote; two people writing at once both arrive.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			r, err := peopleRepo()
			if err != nil {
				return err
			}
			l, err := people.Load(r)
			if err != nil {
				return fail(ExitError, "error", "%v", err)
			}
			return printPeople(cmd, l)
		},
	}
	cmd.AddCommand(newPeopleAliasCmd(), newPeopleRenameCmd())
	return cmd
}

func newPeopleAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "alias <name> <other-name>...",
		Short: "Record that other names — a git user.name, an email — mean the person <name>",
		Long: `Adds each <other-name> to the person <name> means. <name> may be their jaira
name or any name they already have; a name nobody has yet starts a new person
under it.

A name that already means somebody else is refused with exit 3: two people
are made one by 'jaira people rename', never by taking a name off one of them.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, changed, err := updatePeople(func(l *people.List) (bool, error) {
				did := false
				for _, a := range args[1:] {
					c, err := l.Alias(args[0], a)
					if err != nil {
						return false, err
					}
					did = did || c
				}
				return did, nil
			})
			if err != nil {
				return err
			}
			// nil when nothing changed for a name nobody has: "alias zed zed"
			// names nobody new, so there is no person to show.
			p := l.Find(args[0])
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"person": p, "changed": changed})
			}
			if !changed {
				fmt.Fprintf(cmd.OutOrStdout(), "Nothing to change: %s already goes by %s.\n", strings.TrimSpace(args[0]), strings.Join(args[1:], ", "))
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s now also goes by %s, on %s.\n", p.Name, strings.Join(args[1:], ", "), gitref.PeopleRef)
			return nil
		},
	}
}

func newPeopleRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Give a person a new jaira name; the old one stays theirs as an alias",
		Long: `Makes <new> the jaira name of the person <old> means. The old name becomes one
of their aliases, so every ticket written under it is shown under <new>, is
found by 'list --assignee <new>' and stays theirs — no ticket file is rewritten.

<new> may be one of their own aliases. A name of anybody else is refused with
exit 3, and an <old> nobody goes by with exit 5.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, changed, err := updatePeople(func(l *people.List) (bool, error) {
				before := l.Canonical(args[0])
				if err := l.Rename(args[0], args[1]); err != nil {
					return false, err
				}
				return before != l.Canonical(args[1]), nil
			})
			if err != nil {
				return err
			}
			p := l.Find(args[1])
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"person": p, "changed": changed})
			}
			if !changed {
				fmt.Fprintf(cmd.OutOrStdout(), "Nothing to change: %s is already called that.\n", p.Name)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s is now %s; %s stays theirs as an alias, on %s.\n", args[0], p.Name, args[0], gitref.PeopleRef)
			return nil
		},
	}
}

// peopleRepo is the repository the list travels in: the board's, with the
// remote it uses.
func peopleRepo() (*gitref.Repo, error) {
	if _, err := openStore(); err != nil {
		return nil, err
	}
	if refs == nil || refs.Repo == nil {
		return nil, fail(ExitValidation, "no_refs", "the people list lives on a git ref, and this board has no repository")
	}
	return refs.Repo, nil
}

// updatePeople runs one change against the newest list and turns the
// package's errors into exit codes a script branches on.
func updatePeople(change func(*people.List) (bool, error)) (*people.List, bool, error) {
	r, err := peopleRepo()
	if err != nil {
		return nil, false, err
	}
	l, changed, err := people.Update(r, change)
	switch {
	case err == nil:
		return l, changed, nil
	case errors.Is(err, people.ErrBadName):
		return nil, false, fail(ExitUsage, "bad_name", "%v", err)
	case errors.Is(err, people.ErrTaken):
		return nil, false, fail(ExitValidation, "name_taken", "%v — 'jaira people' lists who has which name", err)
	case errors.Is(err, people.ErrUnknown):
		return nil, false, fail(ExitNotFound, "no_such_person", "%v — 'jaira people' lists them", err)
	case errors.Is(err, gitref.ErrNoRepo), errors.Is(err, gitref.ErrNoGitRepo):
		return nil, false, fail(ExitValidation, "no_refs", "%v — the people list lives on a git ref and needs the board's remote; 'jaira share' sets one up", err)
	case errors.Is(err, gitref.ErrOffline):
		return nil, false, fail(ExitError, "offline", "%v — nothing was written; run it again once the remote is reachable", err)
	}
	return nil, false, fail(ExitError, "error", "%v", err)
}

func printPeople(cmd *cobra.Command, l *people.List) error {
	sorted := l.Sorted()
	if g.jsonOut {
		if sorted == nil {
			sorted = []people.Person{}
		}
		for i := range sorted {
			if sorted[i].Aliases == nil {
				sorted[i].Aliases = []string{}
			}
		}
		return emit(cmd.OutOrStdout(), map[string]any{"people": sorted, "ref": gitref.PeopleRef})
	}
	w := cmd.OutOrStdout()
	if len(sorted) == 0 {
		fmt.Fprintln(w, "Nobody is named yet: everybody acts under their git user.name.")
		fmt.Fprintln(w, "'jaira people alias <name> <git-name|email>' starts the list.")
		return nil
	}
	for _, p := range sorted {
		if len(p.Aliases) == 0 {
			fmt.Fprintln(w, p.Name)
			continue
		}
		fmt.Fprintf(w, "%s — also %s\n", p.Name, strings.Join(p.Aliases, ", "))
	}
	return nil
}
