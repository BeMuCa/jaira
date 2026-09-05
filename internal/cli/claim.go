package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
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

// StaleClaim reports a claim that has expired and is therefore about to be
// stepped over: who held it, and how long ago it was last renewed.
//
// Expiry is deliberate — there is no server to notice a session died, so a lease
// that outlived its holder would wedge the board — but until now it happened in
// silence. Two sessions could work the same ticket with nothing anywhere saying
// the second one walked over the first.
func StaleClaim(t *ticket.Ticket, sessionID string) (holder string, since time.Duration, stale bool) {
	if t.ClaimedBy == "" || t.ClaimedBy == sessionID {
		return "", 0, false
	}
	if t.ClaimedAt.IsZero() {
		// A claim with no timestamp cannot be aged, so it is treated as
		// abandoned by ClaimActive. Report it as such rather than saying "0m".
		return t.ClaimedBy, 0, true
	}
	if d := time.Since(t.ClaimedAt); d > ClaimTTL {
		return t.ClaimedBy, d, true
	}
	return "", 0, false
}

// staleClaimLine is the one-line warning, or "" when there is nothing to say.
func staleClaimLine(t *ticket.Ticket, sessionID string) string {
	holder, since, stale := StaleClaim(t, sessionID)
	if !stale {
		return ""
	}
	if since == 0 {
		return fmt.Sprintf("%s carries an abandoned claim by %s (no timestamp)", ticket.Handle(t.ID), holder)
	}
	return fmt.Sprintf("%s carries an abandoned claim by %s, last renewed %s ago", ticket.Handle(t.ID), holder, rough(since))
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
sessions; it does not prevent a deliberate move.

Taking over a claim that has expired is allowed and needs no flag, but it is
reported: whose claim it was and how long ago it was last renewed. Silence there
would let two sessions work the same ticket with nothing anywhere saying so.

--release gives a claim up: your own at any time, somebody else's only once it
has expired — that clears an abandoned claim's warning without waiting. A claim
another session is actively renewing is refused.`,
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
				// Your own claim goes at will; somebody else's only once it has
				// expired. A live claim is another session mid-work — releasing
				// it from outside would invite exactly the double-work claims
				// exist to prevent.
				if holder, active := ClaimActive(t, sessionID); active {
					return &codedError{
						code:   ExitValidation,
						reason: "claimed",
						message: fmt.Sprintf(
							"%s is claimed by %s and the claim is still live (renewed %s ago); only an expired claim may be released by someone else",
							ticket.Handle(t.ID), holder, rough(time.Since(t.ClaimedAt))),
					}
				}
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

			// Read before the write, since claiming overwrites the evidence.
			tookOver := staleClaimLine(t, sessionID)

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
				out := map[string]any{
					"claimed": true, "id": t.ID, "session": sessionID,
					"expires_at": ticket.FormatTime(now.Add(ClaimTTL)),
				}
				if tookOver != "" {
					out["took_over"] = tookOver
				}
				return emit(cmd.OutOrStdout(), out)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Claimed %s until %s\n",
				ticket.Handle(t.ID), now.Add(ClaimTTL).Format("15:04"))
			if tookOver != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", tookOver)
			}
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
