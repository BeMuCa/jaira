package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newPullCmd takes a ticket off its ref and onto this machine.
func newPullCmd() *cobra.Command {
	var steal bool
	cmd := &cobra.Command{
		Use:   "pull <id>",
		Short: "Take a ticket that lives on its ref and start working it here",
		Long: `Takes a ticket over on its ref and writes it here as a file.

A captured ticket belongs to nobody. Pulling it is what makes it yours: this
records you as its assignee on the ref, and only then puts the file in
.jaira/tickets/ — in that order, so that losing the race leaves you with no
file rather than with a ticket somebody else is working.

Two things can refuse it, and they mean different things:

  taken     somebody else already pulled it; the message names them, and
            --steal takes it anyway
  race      two clones pulled in the same instant and the other one landed
            first; re-read and decide

Pulling a ticket you already have is a successful no-op, so this is safe to
repeat.

This is the one command that needs the network. That is not an exception to
working offline but a consequence of it: with no route to the remote there is
no ref to read, so there is nothing to take over.`,
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
			got, pullErr := refs.Pull(id, steal)
			err = pullErr

			switch {
			case errors.Is(err, refsync.ErrTaken):
				return refusePull(got, "taken",
					"%s is not yours to pull: %s\n  pass --steal to take it anyway",
					ticket.Handle(id), holder(got))
			case errors.Is(err, gitref.ErrRaceLost):
				return refusePull(got, "race",
					"%s was taken while you were pulling it: %s\n  nothing was written here; run 'jaira fetch' and decide",
					ticket.Handle(id), describe(got))
			case err != nil:
				return err
			}

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), got)
			}
			if got.AlreadyHere {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already here: %s\n", ticket.Handle(got.ID), got.Path)
				return nil
			}
			if got.TakenFrom != "" {
				// Loud on purpose: the other person's board still says the
				// ticket is theirs until they fetch, and the only record of
				// why it stopped being is the note this wrote on the ticket.
				fmt.Fprintf(cmd.OutOrStdout(), "%s taken over from %s — a note on the ticket says so\n",
					ticket.Handle(got.ID), got.TakenFrom)
				fmt.Fprintf(cmd.OutOrStdout(), "  tell them: their board still shows it as theirs until they run 'jaira fetch'\n")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s is yours: %s\n", ticket.Handle(got.ID), got.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", got.Path)
			fmt.Fprintf(cmd.OutOrStdout(), "  commit it with the work, so the change and what it was for arrive together\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&steal, "steal", false, "take a ticket somebody else has pulled")
	return cmd
}

// holder names who the ticket belongs to, for the refusal that is about
// ownership rather than about a lost push.
func holder(got *refsync.Pulled) string {
	if got != nil && got.Winner != nil {
		return got.Winner.Holds()
	}
	return "somebody else has it"
}

// describe renders the other side of a refusal, or says plainly that we could
// not find out — never a blank where a name should be.
func describe(got *refsync.Pulled) string {
	if got != nil && got.Winner != nil {
		return got.Winner.Describe()
	}
	return "somebody else has it"
}

// refusePull returns the refusal as a validation error, so a script sees exit 3
// and, with --json, a reason it can branch on rather than a sentence.
func refusePull(got *refsync.Pulled, reason, format string, args ...any) error {
	err := &codedError{
		code:    ExitValidation,
		reason:  reason,
		message: fmt.Sprintf(format, args...),
	}
	if got != nil && got.Winner != nil {
		err.violations = gate.Violations{{Message: got.Winner.Describe()}}
	}
	return err
}
