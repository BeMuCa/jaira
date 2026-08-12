package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/berk/jaira/core/session"
	"github.com/berk/jaira/core/ticket"
)

// ClaimTTL is how long a claim is honoured without being renewed.
//
// Claims expire rather than being held until released, because there is no server
// to notice that a session died. A lease that outlives its holder would wedge the
// board and require manual unlocking — which, with no coordinator, means editing a
// file by hand. Expiry makes abandonment self-healing.
const ClaimTTL = 30 * time.Minute

// ClaimActive reports whether a ticket is claimed by someone other than sessionID.
func ClaimActive(t *ticket.Ticket, sessionID string) (holder string, active bool) {
	if t.ClaimedBy == "" || t.ClaimedBy == sessionID {
		return "", false
	}
	if t.ClaimedAt.IsZero() || time.Since(t.ClaimedAt) > ClaimTTL {
		return t.ClaimedBy, false // stale: treat as abandoned
	}
	return t.ClaimedBy, true
}

func newClaimCmd() *cobra.Command {
	var release, steal bool
	var sessionID string
	cmd := &cobra.Command{
		Use:   "claim <id>",
		Short: "Take a short-lived claim on a ticket",
		Long: `Marks a ticket as being worked on by this session.

Claims exist so parallel sessions do not pick up the same ticket. They are leases,
not locks: a claim that is not renewed within 30 minutes is treated as abandoned,
so a crashed session never leaves work permanently unreachable and there is
nothing to unlock by hand.

Claiming is advisory. It keeps 'jaira next' from handing the same work to two
sessions; it does not prevent a deliberate move.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if sessionID == "" {
				sessionID = session.DefaultID()
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}

			if release {
				t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
					if err := t.Doc().SetScalar(ticket.FieldClaimedBy, ""); err != nil {
						return err
					}
					return t.Doc().SetScalar(ticket.FieldClaimedAt, "")
				})
				if err != nil {
					return err
				}
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"released": true, "id": t.ID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Released %s\n", ticket.Handle(t.ID))
				return nil
			}

			if holder, active := ClaimActive(t, sessionID); active && !steal {
				return &codedError{
					code:   ExitValidation,
					reason: "claimed",
					message: fmt.Sprintf(
						"%s is claimed by %s (%s ago); pass --steal to take it anyway",
						ticket.Handle(t.ID), holder, rough(time.Since(t.ClaimedAt))),
				}
			}

			now := time.Now()
			t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
				if err := t.Doc().SetScalar(ticket.FieldClaimedBy, sessionID); err != nil {
					return err
				}
				return t.Doc().SetRaw(ticket.FieldClaimedAt, ticket.FormatTime(now))
			})
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"claimed": true, "id": t.ID, "session": sessionID,
					"expires_at": ticket.FormatTime(now.Add(ClaimTTL)),
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Claimed %s until %s\n",
				ticket.Handle(t.ID), now.Add(ClaimTTL).Format("15:04"))
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&release, "release", false, "give up a claim")
	f.BoolVar(&steal, "steal", false, "take a claim held by another live session")
	f.StringVar(&sessionID, "session", "", "session identifier (defaults to a stable per-shell value)")
	return cmd
}

func rough(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "moments"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

// laneOutput is the structured result a lane's agent returns.
type laneOutput struct {
	Outcome struct {
		What     string `json:"what"`
		Why      string `json:"why"`
		Resolves string `json:"resolves"`
	} `json:"outcome"`
	Commits    []string `json:"commits"`
	Question   string   `json:"question"`
	ExecutedBy string   `json:"executed_by"`
	Signal     string   `json:"signal"`
}

// readLaneOutput parses an agent's structured output and reports which fields the
// lane demanded but did not receive.
//
// Validation lives here, in the tool, rather than being left to the agent's own
// judgement about whether it followed the contract. The agent supplies content;
// whether that content satisfies the contract is not its call to make.
func readLaneOutput(r io.Reader, produces []string) (*laneOutput, []string, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 {
		return nil, nil, fmt.Errorf("expected a JSON object on stdin")
	}
	var out laneOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, nil, fmt.Errorf("could not parse the lane output: %w", err)
	}
	get := func(field string) string {
		switch field {
		case ticket.FieldOutcomeWhat:
			return out.Outcome.What
		case ticket.FieldOutcomeWhy:
			return out.Outcome.Why
		case ticket.FieldOutcomeResolves:
			return out.Outcome.Resolves
		case ticket.FieldQuestion:
			return out.Question
		case ticket.FieldCommits, "diff":
			return strings.Join(out.Commits, ",")
		}
		return ""
	}
	var missing []string
	for _, p := range produces {
		if strings.TrimSpace(get(p)) == "" {
			missing = append(missing, p)
		}
	}
	return &out, missing, nil
}
