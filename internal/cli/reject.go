package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/refsync"
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
original instead: the kind becomes duplicate, the original goes into
related and into rejected-duplicate-of. 'jaira restore' takes both off again,
but keeps a related link that was there before the rejection.

Rejecting is a person's decision. Without --force the command refuses (exit 3)
and says what an agent does instead: move the ticket to human with
--question "reject? reason: ...", where a person decides. On the board, R
rejects with a one-line reason; in a shell, a person passes --force.

A ticket that lives only on its ref — nobody pulled it here — is rejected
on the ref: the command fetches, writes the rejection there and leaves no
file in this checkout, so no branch carries an archive file for it. Everybody
sees it after their next fetch, and it leaves every board. A ticket with a
file here is archived here, as before.

'jaira restore' brings a rejected ticket back and takes the rejection off it;
for one rejected on its ref, 'jaira pull' does, and says so.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			why = strings.Join(strings.Fields(why), " ")
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
			t, err := loadForReject(s, args[0])
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
			// Resolved before anything is stamped or written, so a wrong
			// original leaves the ticket exactly as it was.
			if duplicateOf != "" {
				orig, _, err := s.LoadAnywhere(duplicateOf)
				if errors.Is(err, ticket.ErrNotFound) {
					return &codedError{code: ExitNotFound, reason: "not_found",
						message: fmt.Sprintf("--duplicate-of %s: no such ticket on the board or filed away", duplicateOf)}
				}
				if err != nil {
					return err
				}
				if orig.ID == t.ID {
					return fail(ExitUsage, "usage", "%s cannot be a duplicate of itself", ticket.Handle(t.ID))
				}
				duplicateOf = orig.ID
			}
			if t.ReadOnly {
				return rejectOnRef(cmd, t, kind, why, duplicateOf)
			}
			// Stamped before the ticket leaves the board, as archive does.
			if env, _, envErr := loadEnv(s); envErr == nil {
				_, _ = s.StampCommits(t, env.DeriveCommits)
			}
			dst, err := s.Reject(t.ID, kind, why, duplicateOf, time.Now())
			if err != nil {
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

// loadForReject finds the ticket to reject wherever it lives: a file here, or
// only its ref. A person saying "not this" should not have to know which.
//
// Not found here, it fetches once and looks again: a ticket created in another
// clone after this one last fetched is unknown to the local refs, and "not
// found" for a ticket that exists is the wrong answer. A ticket that is already
// filed away — rejected or archived — is named as such, not as unknown.
func loadForReject(s *ticket.Store, arg string) (*ticket.Ticket, error) {
	t, err := s.Load(arg)
	var fetchErr error
	if errors.Is(err, ticket.ErrNotFound) && refs.Usable() == nil {
		if fetchErr = refs.Repo.Fetch(); fetchErr == nil {
			t, err = s.Load(arg)
		}
	}
	if !errors.Is(err, ticket.ErrNotFound) {
		if err == nil && t.ReadOnly && t.Rejected.Kind != "" {
			return nil, alreadyRejected(t)
		}
		return t, err
	}
	if filed, _, ferr := s.LoadAnywhere(arg); ferr == nil {
		if filed.Rejected.Kind != "" {
			return nil, alreadyRejected(filed)
		}
		return nil, &codedError{code: ExitValidation, reason: "filed_away",
			message: fmt.Sprintf("%s is already filed away (%s); 'jaira restore' brings it back first",
				ticket.Handle(filed.ID), filepath.Base(filed.Path))}
	}
	if fetchErr != nil {
		// Not here, and the remote that might know it could not be asked:
		// "no such ticket" would be a guess.
		return nil, &codedError{code: ExitError, reason: "no_remote",
			message: fmt.Sprintf("%v here, and the remote could not be asked: %v\n  nothing was written; try again when the remote is reachable", err, fetchErr)}
	}
	return nil, &codedError{code: ExitNotFound, reason: "not_found",
		message: fmt.Sprintf("%v — not on the board, not on a ref, not filed away", err)}
}

func alreadyRejected(t *ticket.Ticket) error {
	return &codedError{code: ExitValidation, reason: "already_rejected",
		message: fmt.Sprintf("%s is already rejected (%s by %s): %s",
			ticket.Handle(t.ID), t.Rejected.Kind, t.Rejected.By, t.Rejected.Why)}
}

// rejectOnRef rejects a ticket that has no file here on its ref, and writes no
// file anywhere: the rejection reaches everybody with their next fetch, and no
// branch carries an archive file for a ticket nobody worked.
func rejectOnRef(cmd *cobra.Command, t *ticket.Ticket, kind, why, duplicateOf string) error {
	h := ticket.Handle(t.ID)
	got, heldBy, err := refs.Reject(t.ID, kind, why, duplicateOf, time.Now())
	switch {
	case errors.Is(err, refsync.ErrAlreadyRejected):
		return alreadyRejected(got)
	case errors.Is(err, gitref.ErrRaceLost):
		return &codedError{code: ExitValidation, reason: "race",
			message: fmt.Sprintf("%s changed on the remote while it was being rejected; nothing was written — run the same command again", h)}
	case errors.Is(err, refsync.ErrUnsent):
		return &codedError{code: ExitValidation, reason: "unsent",
			message: fmt.Sprintf("%s: %v", h, err)}
	case errors.Is(err, ticket.ErrNotFound):
		return &codedError{code: ExitNotFound, reason: "not_found",
			message: fmt.Sprintf("%s: no longer on the remote; nothing was written", h)}
	case errors.Is(err, gitref.ErrOffline), errors.Is(err, gitref.ErrNoRepo):
		return &codedError{code: ExitError, reason: "no_remote",
			message: fmt.Sprintf("%s lives only on its ref, and the ref could not be reached: %v\n  nothing was written; try again when the remote is reachable", h, err)}
	case err != nil:
		return err
	}
	kind = got.Rejected.Kind
	if g.jsonOut {
		return emit(cmd.OutOrStdout(), map[string]any{
			"rejected": true, "id": got.ID, "handle": h,
			"kind": kind, "why": got.Rejected.Why, "path": "", "file": "", "on-ref-only": true,
		})
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Rejected %s (%s)  %s\n", h, kind, got.Title)
	fmt.Fprintf(cmd.OutOrStdout(), "Written to its ref, not to a file here — everybody sees it after their next 'jaira fetch'.\n")
	if heldBy != "" {
		fmt.Fprintf(cmd.ErrOrStderr(), "jaira: %s\n", refsync.HeldByWarning(h, heldBy))
	}
	return nil
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
