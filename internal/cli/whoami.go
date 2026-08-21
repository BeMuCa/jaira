package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	coreidentity "github.com/BeMuCa/jaira/core/identity"
)

// newWhoamiCmd exists because the ownership rail compares strings and a person
// is several. When a move is refused with "belongs to <someone>" and that
// someone is you under another name, nothing else on the board would tell you
// which name jaira thinks you have.
func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the identity jaira acts as, and the other names that mean you",
		Long: `Prints the name jaira attributes work to, followed by every other string it
treats as the same person: git's user.email as well as user.name, plus anything
listed in the alias file.

A ticket assigned under a name that is not in this list is somebody else's as
far as the gates are concerned, and moving it will be refused. Add the name to
the alias file rather than passing --force to every move.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			me := identity()
			aliases := coreidentity.Aliases(g.dir)
			path := coreidentity.AliasesPath()

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"identity": me, "aliases": aliases, "alias_file": path,
				})
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "Acting as:   %s\n", me)
			for i, a := range aliases {
				if i == 0 {
					continue // the canonical name, already printed above
				}
				fmt.Fprintf(w, "Also me:     %s\n", a)
			}
			fmt.Fprintf(w, "Alias file:  %s\n", dash(path))
			return nil
		},
	}
}
