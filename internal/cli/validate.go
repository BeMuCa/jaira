package cli

import (
	"fmt"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/core/validate"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var strict bool
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Check every ticket on the board for damage",
		Long: `Reads every ticket and reports what is wrong with it.

The lane gates fire when a ticket moves. This checks the board at rest, which is
where damage from a hand edit, a bad merge, or an agent writing something
unexpected actually shows up — a ticket with an unparseable id, a lane that is
not installed, a dependency on something that is not there.

Errors exit 3. Warnings — a backlog ticket that is not specified yet — exit 0,
because capture is meant to be cheap. Use --strict to fail on those too.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			lanes, err := lane.Load()
			if err != nil {
				return err
			}
			tickets, err := s.List()
			if err != nil {
				// A ticket that cannot be read at all is the most severe finding
				// there is, so it is reported rather than aborting the run.
				pe, ok := err.(*ticket.PartialError)
				if !ok {
					return err
				}
				for _, p := range pe.Problems {
					fmt.Fprintf(cmd.ErrOrStderr(), "jaira: unreadable: %s\n", p)
				}
			}

			problems := validate.Tickets(tickets, lanes)

			if g.jsonOut {
				out := make([]map[string]any, 0, len(problems))
				for _, p := range problems {
					out = append(out, map[string]any{
						"code": p.Code, "severity": p.Severity, "handle": p.Handle,
						"id": p.ID, "path": p.Path, "field": p.Field, "message": p.Message,
					})
				}
				if err := emit(cmd.OutOrStdout(), map[string]any{
					"checked": len(tickets), "problems": out,
					"errors": validate.HasErrors(problems),
				}); err != nil {
					return err
				}
			} else {
				w := cmd.OutOrStdout()
				for _, p := range problems {
					label := "warning"
					if p.Severity == validate.SeverityError {
						label = "error"
					}
					where := p.Handle
					if where == "" {
						where = p.Path
					}
					if p.Field != "" {
						fmt.Fprintf(w, "%-7s %s  %s: %s\n", label, where, p.Field, p.Message)
					} else {
						fmt.Fprintf(w, "%-7s %s  %s\n", label, where, p.Message)
					}
				}
				if len(problems) == 0 {
					fmt.Fprintf(w, "%d ticket(s) checked, nothing wrong.\n", len(tickets))
				} else {
					fmt.Fprintf(w, "\n%d ticket(s) checked, %d problem(s).\n", len(tickets), len(problems))
				}
			}

			if validate.HasErrors(problems) || (strict && len(problems) > 0) {
				return &codedError{code: ExitValidation, reason: "invalid_tickets",
					message: fmt.Sprintf("%d problem(s) found", len(problems))}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&strict, "strict", false, "fail on warnings as well as errors")
	return cmd
}
