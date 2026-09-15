package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	coreidentity "github.com/BeMuCa/jaira/core/identity"
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
anyway, for the case where that person is not coming back.

A ticket that has no ref at all is put on one first. That is how a ticket
created while the board had no usable remote gets back where it belongs: give
the repository the remote, then release the ticket.`,
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
			if errors.Is(relErr, gitref.ErrNoRef) {
				// The ticket never reached a ref: it was created while this
				// board had no usable remote, so it only exists as a file. That
				// is the same handover one step earlier — put it on its ref
				// first, then hand it back — and it is the only way back for a
				// board that spent a while in file mode.
				got, relErr = releaseFromFile(s, id, force)
			}
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

// releaseFromFile hands back a ticket that has no ref at all, by putting it on
// one.
//
// It is deliberately the same command rather than a second one: 'release' has
// always meant "the file goes, the ticket is on its ref, nobody holds it", and a
// ticket written while the board had no remote is that sentence with the middle
// part missing. Clearing the assignee is not a side effect here but the reason —
// the tickets this exists for are unworked ones that were never meant to live in
// one person's checkout.
//
// The write itself goes through the store, so the queue, the lock and the
// timestamps are the ones every other write uses; putOnRef then sends it and
// takes the file away.
func releaseFromFile(s *ticket.Store, id string, force bool) (*refsync.Released, error) {
	t, err := s.Load(id)
	if err != nil {
		// No ref and no file: the original ErrNoRef was the honest answer.
		return nil, gitref.ErrNoRef
	}
	if holder := strings.TrimSpace(t.Assignee); holder != "" && !force && !coreidentity.IsMe(s.Root, holder) {
		return &refsync.Released{ID: t.ID, Title: t.Title},
			fmt.Errorf("%w: %s has it", refsync.ErrTaken, holder)
	}
	t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
		return t.Doc().SetScalar(ticket.FieldAssignee, "")
	})
	if err != nil {
		return nil, err
	}
	out := &refsync.Released{ID: t.ID, Title: t.Title}
	onRef, why := putOnRef(t)
	if !onRef {
		if why == nil {
			why = errors.New("the remote did not accept it")
		}
		return out, fmt.Errorf("%s could not be put on its ref: %w", ticket.Handle(t.ID), why)
	}
	out.Removed = t.Path
	return out, nil
}
