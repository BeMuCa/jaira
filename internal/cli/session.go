package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newCheckpointCmd() *cobra.Command {
	var focus, reasoning, ticketID, model, sessionID string
	var clear bool
	cmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Record what this session is currently working on",
		Long: `Records the current focus so the board can show it.

This is the memory made visible: a session that stops and restarts loses its
conversation, but the board still shows what was being worked on and why. Call it
when the topic changes rather than on every turn — the value is a readable trail,
not a transcript.

Session state is per working tree and is never committed.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if sessionID == "" {
				sessionID = session.DefaultID()
			}
			if clear {
				if err := session.Remove(s, sessionID); err != nil {
					return err
				}
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"cleared": true, "session": sessionID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Cleared session %s\n", sessionID)
				return nil
			}
			if focus == "" && reasoning == "" && ticketID == "" {
				return fail(ExitUsage, "usage", "give at least one of --focus, --why or --ticket")
			}

			// Merge with what is already recorded so a partial update does not
			// erase the rest of the context.
			sess := session.Read(s, sessionID)
			if focus != "" {
				sess.Focus = focus
			}
			if reasoning != "" {
				sess.Reasoning = reasoning
			}
			if model != "" {
				sess.Model = model
			}
			if ticketID != "" {
				t, err := s.Load(ticketID)
				if err != nil {
					return err
				}
				sess.TicketID = t.ID
			}
			sess.UpdatedAt = ticket.FormatTime(time.Now())
			if err := session.Save(s, sess); err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), sess)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint recorded for session %s\n", sessionID)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&focus, "focus", "", "what this session is working on right now")
	f.StringVar(&reasoning, "why", "", "why it is being worked on")
	f.StringVar(&ticketID, "ticket", "", "the ticket this session is on")
	f.StringVar(&model, "model", "", "model doing the work")
	f.StringVar(&sessionID, "session", "", "session identifier (defaults to a stable per-shell value)")
	f.BoolVar(&clear, "clear", false, "remove this session's state")
	return cmd
}

func newSessionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sessions",
		Short: "List sessions working in this tree",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			sess, err := session.Load(s)
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"sessions": sess, "count": len(sess)})
			}
			if len(sess) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No sessions have checkpointed in this working tree.")
				return nil
			}
			for _, x := range sess {
				state := "live"
				if x.Stale() {
					state = "stale"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-6s %s\n", x.ID, state, x.Focus)
				if x.Reasoning != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "%-27s %s\n", "", x.Reasoning)
				}
				if x.TicketID != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "%-27s on %s\n", "", ticket.Handle(x.TicketID))
				}
			}
			return nil
		},
	}
}
