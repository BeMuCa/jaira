package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newMoveCmd() *cobra.Command {
	var (
		to, question, executedBy string
		what, why, resolves      string
		commits                  []string
		force                    bool
		fromLane                 string
	)
	cmd := &cobra.Command{
		Use:     "move <id>",
		Aliases: []string{"advance"},
		Short:   "Move a ticket to another lane",
		Long: `Moves a ticket, applying the gates for the lane being left and the lane
being entered.

The gates are checked here, at the moment of mutation, rather than only when work
is discovered. An agent that skips 'jaira next' and calls this directly is held
to exactly the same rules.

Moving to a lane the ticket is already in succeeds and changes nothing, so this
command is safe to retry.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if to == "" {
				return fail(ExitUsage, "usage", "--to is required")
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			env, _, err := loadEnv(s)
			if err != nil {
				return err
			}

			// A lane's agent returns structured output; the tool validates it
			// against what the lane declared it produces, and rejects malformed
			// output with a reason the agent can read and retry against.
			if fromLane != "" {
				l, ok := env.Lanes.Get(fromLane)
				if !ok {
					return fail(ExitUsage, "no_such_lane", "no lane %q is installed", fromLane)
				}
				out, missing, err := readLaneOutput(cmd.InOrStdin(), l.OutputProduces)
				if err != nil {
					return fail(ExitValidation, "bad_lane_output", "%v", err)
				}
				if len(missing) > 0 {
					return &codedError{
						code:   ExitValidation,
						reason: "incomplete_lane_output",
						message: fmt.Sprintf("lane %q requires %s in its output",
							l.ID, strings.Join(missing, ", ")),
					}
				}
				if what == "" {
					what = out.Outcome.What
				}
				if why == "" {
					why = out.Outcome.Why
				}
				if resolves == "" {
					resolves = out.Outcome.Resolves
				}
				if question == "" {
					question = out.Question
				}
				if executedBy == "" {
					executedBy = out.ExecutedBy
				}
				commits = append(commits, out.Commits...)
			}

			// Idempotent: re-running a satisfied move is a no-op success rather
			// than an error, because parallel agents retry.
			if t.Status == to && question == "" && what == "" && why == "" && resolves == "" && len(commits) == 0 {
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{
						"ticket": ticketJSON(t, env.Lanes), "moved": false, "reason": "already_in_lane",
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already in %s\n", ticket.Handle(t.ID), to)
				return nil
			}

			// Apply the outcome and question fields first: a gate may require
			// them, and requiring two commands to satisfy one move would be a
			// poor contract for an agent.
			staged := func(t *ticket.Ticket) error {
				set := func(field, val string) error {
					if val == "" {
						return nil
					}
					return t.Doc().SetScalar(field, val)
				}
				if err := set(ticket.FieldQuestion, question); err != nil {
					return err
				}
				if err := set(ticket.FieldOutcomeWhat, what); err != nil {
					return err
				}
				if err := set(ticket.FieldOutcomeWhy, why); err != nil {
					return err
				}
				if err := set(ticket.FieldOutcomeResolves, resolves); err != nil {
					return err
				}
				if err := set(ticket.FieldExecutedBy, executedBy); err != nil {
					return err
				}
				if len(commits) > 0 {
					merged := append([]string{}, t.Commits...)
					for _, c := range commits {
						if c = strings.TrimSpace(c); c != "" && !contains(merged, c) {
							merged = append(merged, c)
						}
					}
					if err := t.Doc().SetList(ticket.FieldCommits, merged); err != nil {
						return err
					}
				}
				return nil
			}

			if _, err := s.Mutate(t.ID, staged); err != nil {
				return err
			}
			if t, err = s.Load(t.ID); err != nil {
				return err
			}
			env, _, err = loadEnv(s)
			if err != nil {
				return err
			}

			req := gate.Request{To: to, Question: question, Actor: identity()}
			vs := gate.CheckAdvance(env, t, req)
			if len(vs) > 0 && !force {
				code := ExitValidation
				for _, v := range vs {
					if v.Code == gate.CodeBlocked || v.Code == gate.CodeSelfBlock {
						code = ExitBlocked
					}
				}
				return &codedError{
					code:       code,
					reason:     "gate_refused",
					message:    fmt.Sprintf("cannot move %s to %s:\n%s", ticket.Handle(t.ID), to, bullets(vs)),
					violations: vs,
				}
			}

			t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
				if err := t.Doc().SetScalar(ticket.FieldStatus, to); err != nil {
					return err
				}
				return ticket.SetReady(t.Doc(), gate.Ready(t))
			})
			if err != nil {
				return err
			}

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"ticket": ticketJSON(t, env.Lanes), "moved": true,
					"overridden": force && len(vs) > 0,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s → %s\n", ticket.Handle(t.ID), to)
			if force && len(vs) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Overrode %d gate refusal(s):\n%s\n", len(vs), bullets(vs))
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&to, "to", "", "target lane (required)")
	f.StringVar(&question, "question", "", "the question blocking progress, for the human lane")
	f.StringVar(&executedBy, "executed-by", "", "model that performed this run; ownership stays with the assignee")
	f.StringVar(&what, "what", "", "outcome: what was changed")
	f.StringVar(&why, "why", "", "outcome: why it was needed")
	f.StringVar(&resolves, "resolves", "", "outcome: how the change satisfies the definition of done")
	f.StringSliceVar(&commits, "commits", nil, "commit SHAs produced for this ticket")
	f.BoolVar(&force, "force", false, "override gate refusals; recorded in the output")
	f.StringVar(&fromLane, "from-lane", "",
		"read this lane's structured output as JSON on stdin and validate it against the lane's contract")
	return cmd
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func bullets(vs gate.Violations) string {
	parts := make([]string, 0, len(vs))
	for _, v := range vs {
		parts = append(parts, "  - "+v.Message)
	}
	return strings.Join(parts, "\n")
}

func newNextCmd() *cobra.Command {
	var laneFilter, assigneeFilter string
	var all bool
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Show the next actionable ticket",
		Long: `Reports work that could be started right now: past the promotion gate, not
in a terminal lane, and not waiting on an unresolved dependency.

This is a convenience for the happy path. It is not the enforcement point —
'jaira move' re-checks everything, because nothing prevents a caller from
skipping this command.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			env, tickets, err := loadEnv(s)
			if err != nil {
				return err
			}
			var ready []*ticket.Ticket
			for _, t := range tickets {
				if laneFilter != "" && t.Status != laneFilter {
					continue
				}
				if assigneeFilter != "" && !strings.EqualFold(t.Assignee, assigneeFilter) {
					continue
				}
				if !gate.Actionable(env, t) {
					continue
				}
				// Skip work another live session has claimed, so two sessions are
				// not handed the same ticket. Stale claims are ignored.
				if _, active := ClaimActive(t, session.DefaultID()); active {
					continue
				}
				ready = append(ready, t)
			}
			// Order by how far along the work already is, so in-flight tickets
			// are finished before new ones are started.
			sortByProgress(ready, env)

			if !all && len(ready) > 1 {
				ready = ready[:1]
			}
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(ready))
				for _, t := range ready {
					arr = append(arr, ticketJSON(t, env.Lanes))
				}
				return emit(cmd.OutOrStdout(), map[string]any{"tickets": arr, "count": len(arr)})
			}
			if len(ready) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing is actionable. Either everything is done, blocked, or still needs a spec.")
				return nil
			}
			for _, t := range ready {
				printDetail(cmd.OutOrStdout(), t, env)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&laneFilter, "lane", "", "only consider this lane")
	cmd.Flags().StringVar(&assigneeFilter, "assignee", "", "only consider this assignee")
	cmd.Flags().BoolVar(&all, "all", false, "list every actionable ticket instead of just the next one")
	return cmd
}

func sortByProgress(ts []*ticket.Ticket, env gate.Env) {
	prec := func(t *ticket.Ticket) int { return env.Lanes.Precedence(t.Status) }
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0; j-- {
			a, b := ts[j-1], ts[j]
			if prec(b) > prec(a) || (prec(b) == prec(a) && b.ID < a.ID) {
				ts[j-1], ts[j] = b, a
				continue
			}
			break
		}
	}
}

// showForLane assembles the bounded input a lane's agent should receive.
//
// The tool builds this, not the agent: if the agent decided what context it
// needed, the lane's declared contract would be a suggestion. Here it is the
// definition — the agent gets the fields the lane asked for and nothing else.
func showForLane(cmd *cobra.Command, s *ticket.Store, env gate.Env, t *ticket.Ticket, laneID string) error {
	l, ok := env.Lanes.Get(laneID)
	if !ok {
		return fail(ExitUsage, "no_such_lane", "no lane %q is installed", laneID)
	}

	fields := map[string]string{}
	var missing []string
	var diff string
	for _, want := range l.InputRequires {
		switch want {
		case "plan":
			if len(t.PlanItems) == 0 {
				missing = append(missing, "plan (the ticket has no Plan checklist)")
				continue
			}
			var steps []string
			for i, it := range t.PlanItems {
				steps = append(steps, fmt.Sprintf("%d. [%s] %s", i+1, it.State.Marker(), it.Text))
			}
			fields["plan"] = strings.Join(steps, "\n")
		case "diff":
			repo := &gitrepo.Repo{Dir: s.Root}
			if len(t.Commits) == 0 {
				missing = append(missing, "diff (ticket records no commits)")
				continue
			}
			d, err := repo.Diff(t.Commits)
			if err != nil {
				missing = append(missing, fmt.Sprintf("diff (%v)", err))
				continue
			}
			diff = d
		default:
			v := fieldValue(t, want)
			if strings.TrimSpace(v) == "" {
				missing = append(missing, want)
				continue
			}
			fields[want] = v
		}
	}

	if g.jsonOut {
		return emit(cmd.OutOrStdout(), map[string]any{
			"ticket_id":  t.ID,
			"lane":       l.ID,
			"model_tier": l.ModelTier,
			"prompt":     l.Prompt,
			"input":      fields,
			"diff":       diff,
			"produces":   l.OutputProduces,
			"missing":    missing,
			"complete":   len(missing) == 0,
		})
	}

	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "# Lane: %s   (tier: %s)\n\n", l.Name, dash(l.ModelTier))
	if l.Prompt != "" {
		fmt.Fprintf(w, "%s\n\n", l.Prompt)
	}
	fmt.Fprintf(w, "## Ticket %s — %s\n\n", ticket.Handle(t.ID), t.Title)
	for _, want := range l.InputRequires {
		if v, ok := fields[want]; ok {
			fmt.Fprintf(w, "**%s**\n%s\n\n", want, v)
		}
	}
	if diff != "" {
		fmt.Fprintf(w, "## Diff\n\n```diff\n%s```\n\n", diff)
	}
	if len(l.OutputProduces) > 0 {
		fmt.Fprintf(w, "## Must produce\n\n")
		for _, p := range l.OutputProduces {
			fmt.Fprintf(w, "- %s\n", p)
		}
		fmt.Fprintln(w)
	}
	if len(missing) > 0 {
		fmt.Fprintf(w, "## Missing required input\n\n")
		for _, m := range missing {
			fmt.Fprintf(w, "- %s\n", m)
		}
	}
	return nil
}

func fieldValue(t *ticket.Ticket, field string) string {
	switch field {
	case ticket.FieldGoal:
		return t.Goal
	case ticket.FieldDoD:
		return t.DoD
	case ticket.FieldContext:
		return t.Context
	case ticket.FieldTitle:
		return t.Title
	case ticket.FieldAssignee:
		return t.Assignee
	case ticket.FieldQuestion:
		return t.Question
	case ticket.FieldOutcomeWhat:
		return t.Outcome.What
	case ticket.FieldOutcomeWhy:
		return t.Outcome.Why
	case ticket.FieldOutcomeResolves:
		return t.Outcome.Resolves
	case ticket.FieldStatus:
		return t.Status
	default:
		v, _, err := t.Doc().Scalar(field)
		if err != nil {
			return ""
		}
		return v
	}
}

// laneOf is a small helper used by the board.
func laneOf(set *lane.Set, id string) *lane.Lane {
	if l, ok := set.Get(id); ok {
		return l
	}
	return lane.Passthrough(id)
}

var _ = laneOf
var _ = time.Now
