package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/spf13/cobra"
)

// noteHeading is where progress notes are written.
//
// They go in the ticket body rather than a sidecar file on purpose. A separate
// memory file would be one more thing to lose, would not travel with the ticket
// through git, and could drift out of sync with it. The ticket already is the
// memory; this just gives an agent a way to add to it mid-run.
const noteHeading = "## Progress"

func newNoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "note <id> <text>",
		Short: "Record what you have worked out so far, so a later session can pick it up",
		Long: `Appends a timestamped note to the ticket's Progress section.

Write one when you learn something a later session would otherwise have to
rediscover: what you tried, what did not work, what you were about to do next.
A session that stops abruptly — a limit, a crash, a closed laptop — leaves
nothing behind except what was written down, and the point of this board is that
the work does not have to be reconstructed from memory.

Three parts, in order: what you were doing, what you found — especially what
did not work, the part everyone skips — and what the next step is. Skip the
third and the reader has to re-derive a decision you already made.

Write it as if the reader has mild ADHD and knows none of what you know — what
you were doing first, short concrete lines with one point each, names and paths
rather than adjectives, no jargon and no preamble.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			text := strings.TrimSpace(strings.Join(args[1:], " "))
			if text == "" {
				return fail(ExitUsage, "usage", "a note with no text records nothing")
			}
			stamp := time.Now().UTC().Format("2006-01-02 15:04")
			who := identity()

			t, err := s.Mutate(args[0], func(t *ticket.Ticket) error {
				body := t.Doc().Body()
				entry := fmt.Sprintf("- **%s · %s** — %s", stamp, who, text)
				if idx := strings.Index(body, noteHeading); idx >= 0 {
					// Append under the existing heading, after anything already there.
					rest := body[idx:]
					nextHeading := strings.Index(rest[len(noteHeading):], "\n## ")
					if nextHeading < 0 {
						body = strings.TrimRight(body, "\n") + "\n" + entry + "\n"
					} else {
						at := idx + len(noteHeading) + nextHeading
						body = strings.TrimRight(body[:at], "\n") + "\n" + entry + "\n" + body[at:]
					}
				} else {
					body = strings.TrimRight(body, "\n") + "\n\n" + noteHeading + "\n\n" + entry + "\n"
				}
				t.Doc().SetBody(body)
				return nil
			})
			if err != nil {
				return err
			}
			if g.jsonOut {
				env, _, err := loadEnv(s)
				if err != nil {
					return err
				}
				return emit(cmd.OutOrStdout(), ticketJSON(t, env.Lanes))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Noted on %s.\n", ticket.Handle(t.ID))
			return nil
		},
	}
}

func newResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Show work that was left in progress, with everything recorded about it",
		Long: `Lists tickets that were being worked on and stopped.

A session that ends abruptly leaves tickets mid-flight: a plan step marked as in
progress, a claim that has expired, a ticket parked in a working lane. This
gathers them with their progress notes and the step each was on, so a new session
can carry on rather than starting by working out where things were left.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			lanes, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			all, err := s.List()
			if err != nil {
				if pe, ok := err.(*ticket.PartialError); ok {
					fmt.Fprintf(cmd.ErrOrStderr(), "jaira: warning: %v\n", pe)
				} else {
					return err
				}
			}

			type inflight struct {
				t     *ticket.Ticket
				step  string
				why   string
				notes []string
			}
			var out []inflight
			for _, t := range all {
				l, known := lanes.Get(t.Status)
				var why, step string
				for _, it := range t.PlanItems {
					if it.State == ticket.StateDoing {
						step = it.Text
					}
				}
				for _, it := range t.DoDItems {
					if it.State == ticket.StateDoing && step == "" {
						step = it.Text
					}
				}
				switch {
				case step != "":
					why = "a step is still marked in progress"
				case t.ClaimedBy != "" && !t.ClaimedAt.IsZero() && time.Since(t.ClaimedAt) > ClaimTTL:
					why = fmt.Sprintf("claim by %s expired %s ago", t.ClaimedBy, rough(time.Since(t.ClaimedAt)))
				case known && l.Agentic && !l.Terminal:
					why = fmt.Sprintf("parked in %s, which is a working lane", l.Name)
				default:
					continue
				}
				// Reload in full: listing gives the body, but notes are worth
				// reading from the file that is actually on disk right now.
				full, err := s.Load(t.ID)
				if err != nil {
					full = t
				}
				out = append(out, inflight{t: full, step: step, why: why, notes: progressNotes(full.Body)})
			}

			if g.jsonOut {
				items := make([]map[string]any, 0, len(out))
				for _, i := range out {
					items = append(items, map[string]any{
						"id": i.t.ID, "handle": ticket.Handle(i.t.ID), "title": i.t.Title,
						"status": i.t.Status, "reason": i.why, "current_step": i.step,
						"goal": i.t.Goal, "notes": i.notes,
					})
				}
				return emit(cmd.OutOrStdout(), map[string]any{"in_flight": items, "count": len(items)})
			}

			w := cmd.OutOrStdout()
			if len(out) == 0 {
				fmt.Fprintf(w, "Nothing was left in progress.\n")
				return nil
			}
			for _, i := range out {
				fmt.Fprintf(w, "%s  %s\n", ticket.Handle(i.t.ID), i.t.Title)
				fmt.Fprintf(w, "    %s · %s\n", i.t.Status, i.why)
				if i.step != "" {
					fmt.Fprintf(w, "    was on: %s\n", i.step)
				}
				for _, n := range i.notes {
					fmt.Fprintf(w, "    %s\n", n)
				}
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, "%d ticket(s) left in progress.\n", len(out))
			return nil
		},
	}
}

// progressNotes reads the entries under the Progress heading.
func progressNotes(body string) []string {
	var out []string
	in := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			in = strings.EqualFold(trimmed, noteHeading)
			continue
		}
		if in && strings.HasPrefix(trimmed, "- ") {
			out = append(out, strings.TrimPrefix(trimmed, "- "))
		}
	}
	return out
}
