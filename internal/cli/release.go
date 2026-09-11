package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newReleaseCmd hands a ticket back to the board.
func newReleaseCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "release <id>",
		Short: "Hand a ticket back, so somebody else can take it",
		Long: `Clears you as the ticket's assignee on its ref and removes the file from here.

It is the other half of 'jaira pull'. Without it an assignment is a
reservation nobody can give back: while a ticket names you, nobody else can
pull it.

The ref is cleared first and the file removed second, so a failure leaves the
ticket still yours rather than unowned and still on your disk. The file is
removed rather than kept, because a second copy of a ticket somebody else may
now pull is exactly the duplicate this design rules out — nothing is lost, the
ref carries the ticket and committed work stays in the history.

A ticket that is not yours is refused, naming who has it; --force releases it
anyway, for the case where that person is not coming back.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if err := refs.Usable(); err != nil {
				return fail(ExitError, "no_remote",
					"this board does not carry tickets on refs: %v", err)
			}
			id := resolveID(s, args[0])
			got, relErr := refs.Release(id, force)
			err = relErr
			if errors.Is(err, refsync.ErrTaken) {
				who := "somebody else has it"
				if got != nil && got.NotYours != nil {
					who = got.NotYours.Holds()
				}
				return &codedError{
					code:   ExitValidation,
					reason: "not_yours",
					message: fmt.Sprintf("%s is not yours to release: %s\n  pass --force if they are not coming back",
						ticket.Handle(id), who),
				}
			}
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), got)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s is back on the board: %s\n", ticket.Handle(got.ID), got.Title)
			if got.Removed != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  removed %s\n", got.Removed)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "release a ticket somebody else holds")
	return cmd
}
