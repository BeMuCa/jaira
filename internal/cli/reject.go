package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/ticket"
)

func newRejectCmd() *cobra.Command {
	var why, kind, duplicateOf string
	var force bool
	cmd := &cobra.Command{
		Use:   "reject <id> --why <text>",
		Short: "Say why a ticket will not be done, and take it off the board",
		Long: `Writes rejected-kind, rejected-why, rejected-by and rejected-at onto a ticket
and archives it, exactly as 'jaira archive' does. Archived alone, a ticket
looks the same as one merely tidied away; rejected, the archive says why.

--kind is wontfix (the default) or obsolete. --duplicate-of <id> names the
original instead: the kind becomes duplicate and the original goes into
related.

Rejecting is a person's decision. Without --force the command refuses (exit 3)
and says what an agent does instead: move the ticket to human with
--question "reject? reason: ...", where a person decides. On the board, R
rejects with a one-line reason; in a shell, a person passes --force.

'jaira restore' brings a rejected ticket back and takes the rejection off it.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			why = strings.TrimSpace(why)
			if why == "" {
				return fail(ExitUsage, "usage", "reject needs --why <text>: the reason is the whole point of rejecting rather than archiving")
			}
			if duplicateOf != "" && cmd.Flags().Changed("kind") {
				return fail(ExitUsage, "usage", "--duplicate-of sets the kind to duplicate; drop --kind")
			}
			if kind != ticket.RejectWontfix && kind != ticket.RejectObsolete {
				return fail(ExitUsage, "usage", "--kind is %s or %s, not %q", ticket.RejectWontfix, ticket.RejectObsolete, kind)
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			if !force {
				handle, lane := ticket.Handle(t.ID), questionLane(s, t)
				return &codedError{
					code:   ExitValidation,
					reason: "needs_human",
					message: fmt.Sprintf("rejecting %s is a person's decision. An agent moves it to the %s lane instead: "+
						"jaira move %s --to %s --question \"reject? reason: ...\". "+
						"A person rejects it on the board (R), or here with --force",
						handle, lane, handle, lane),
				}
			}
			// Stamped before the ticket leaves the board, as archive does.
			if env, _, envErr := loadEnv(s); envErr == nil {
				_, _ = s.StampCommits(t, env.DeriveCommits)
			}
			dst, err := s.Reject(t.ID, kind, why, duplicateOf, time.Now())
			if err != nil {
				if errors.Is(err, ticket.ErrNotFound) && duplicateOf != "" {
					return &codedError{code: ExitNotFound, reason: "not_found",
						message: fmt.Sprintf("--duplicate-of %s: no such ticket on the board or filed away", duplicateOf)}
				}
				return err
			}
			if duplicateOf != "" {
				kind = ticket.RejectDuplicate
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"rejected": true, "id": t.ID, "handle": ticket.Handle(t.ID),
					"kind": kind, "why": why, "path": dst, "file": filepath.Base(dst),
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Rejected %s (%s)  %s\n", ticket.Handle(t.ID), kind, t.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "Moved to %s — restore it with 'jaira restore %s'.\n",
				filepath.Join(ticket.DirName, ticket.ArchiveSubdir), filepath.Base(dst))
			return nil
		},
	}
	cmd.Flags().StringVar(&why, "why", "", "why the ticket will not be done (required)")
	cmd.Flags().StringVar(&kind, "kind", ticket.RejectWontfix, "wontfix or obsolete")
	cmd.Flags().StringVar(&duplicateOf, "duplicate-of", "", "the ticket this one duplicates; sets the kind to duplicate")
	cmd.Flags().BoolVar(&force, "force", false, "you are a person deciding this; an agent asks in the human lane instead")
	return cmd
}

// questionLane is the lane on t's board that stops for a person's answer —
// "human" on the shipped board — so the refusal names a lane that exists.
func questionLane(s *ticket.Store, t *ticket.Ticket) string {
	if env, err := envFor(s, t); err == nil && env.Lanes != nil {
		for _, l := range env.Lanes.Lanes {
			if l.RequiresQuestion {
				return l.ID
			}
		}
	}
	return "human"
}
