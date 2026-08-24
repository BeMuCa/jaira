package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/ticket"
)

func newDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Remove a ticket's file for good",
		Long: `Deletes a ticket. The file is removed, not moved.

Archive is almost always what you want: it takes a ticket off the board and
'jaira restore' brings it back. This is for a ticket that should never have
existed — a mistyped create, a throwaway probe — where there is nothing worth
keeping and archiving would only file it somewhere.

It asks for the ticket's handle typed back, not a yes. Deletion is the one
irreversible thing this tool does, and a board whose whole promise is that
nothing is lost should make it cost a moment's attention. On a private board the
file is not in git history to recover from either.

--force skips the question, for scripts.

A ticket other tickets still point at is refused: a dangling 'blocked-by' is an
error 'jaira validate' reports and a dependency that can never clear. Delete or
repair those first, or pass --force and accept the board it leaves.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			handle := ticket.Handle(t.ID)

			if refs, err := referencesTo(s, t.ID); err != nil {
				return err
			} else if len(refs) > 0 && !force {
				return &codedError{
					code:   ExitValidation,
					reason: "referenced",
					message: fmt.Sprintf("%s is still referenced by %s; delete or repair those first, or pass --force",
						handle, strings.Join(refs, ", ")),
				}
			}

			if !force {
				// The prompt goes to stderr so --json on stdout stays parseable,
				// and the answer is read from stdin: a caller with nothing to say
				// gets EOF, which is not the handle, which deletes nothing.
				fmt.Fprintf(cmd.ErrOrStderr(), "Delete %s — %s\nThis cannot be undone. Type the handle to confirm: ", handle, t.Title)
				line, _ := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
				if !strings.EqualFold(strings.TrimSpace(line), handle) {
					return &codedError{
						code:    ExitValidation,
						reason:  "not_confirmed",
						message: fmt.Sprintf("that is not %s; nothing was deleted", handle),
					}
				}
			}

			path, err := s.Delete(t.ID)
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"deleted": true, "id": t.ID, "handle": handle, "path": path, "title": t.Title,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s  %s\n%s is gone.\n", handle, t.Title, path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "skip the confirmation, for scripts")
	return cmd
}

// referencesTo lists the tickets that point at id, as "HANDLE (field)".
//
// Deleting something another ticket is blocked by leaves a dependency that can
// never clear, which 'jaira validate' reports as an error — so the delete is
// refused rather than quietly breaking the board it just cleaned up. 'follows'
// is reported the same way: it is not a validation error, but the chain it
// records is the answer to "why does this ticket exist", and that is the whole
// point of the board.
func referencesTo(s *ticket.Store, id string) ([]string, error) {
	all, err := s.List()
	if err != nil {
		// A board with one unreadable ticket can still answer this question for
		// the rest, and refusing to delete because something unrelated is
		// damaged would be its own trap.
		if _, ok := err.(*ticket.PartialError); !ok {
			return nil, err
		}
	}
	var out []string
	for _, o := range all {
		if o.ID == id {
			continue
		}
		if contains(o.BlockedBy, id) {
			out = append(out, ticket.Handle(o.ID)+" (blocked-by)")
		}
		if o.Follows == id {
			out = append(out, ticket.Handle(o.ID)+" (follows)")
		}
	}
	return out, nil
}
