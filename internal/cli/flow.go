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
		reason                   string
		what, why, resolves      string
		commits                  []string
		force                    bool
		dryRun                   bool
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
command is safe to retry.

--dry-run answers "would this be refused, and why" without writing anything: the
same fields are staged in memory, the same gates run, and the exit code is the
one the real move would have returned.`,
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
			if t.Status == to && question == "" && reason == "" && what == "" && why == "" && resolves == "" && len(commits) == 0 {
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
				if err := set(ticket.FieldBlockedReason, reason); err != nil {
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
				// Moving an unassigned ticket claims it: capture leaves the
				// assignee empty, and pulling the ticket into work is the act
				// that makes it yours. Staged before the gate so the promotion
				// gate's assignee requirement is satisfied by the pull itself.
				if strings.TrimSpace(t.Assignee) == "" {
					if err := t.Doc().SetScalar(ticket.FieldAssignee, identity()); err != nil {
						return err
					}
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

			// A dry run stages the same fields the real move would, in memory,
			// and re-decodes so the gate sees exactly the ticket it would have
			// seen. Nothing reaches the disk. This exists because the only way
			// to ask "would this be refused, and why" used to be to try it: an
			// agent created and deleted a throwaway ticket to find out, and one
			// of the user's tickets was moved by accident probing a gate that
			// turned out not to refuse.
			if dryRun {
				if err := staged(t); err != nil {
					return err
				}
				if t, err = ticket.Decode(t.Doc(), t.Path); err != nil {
					return err
				}
			} else {
				if _, err := s.Mutate(t.ID, staged); err != nil {
					return err
				}
				if t, err = s.Load(t.ID); err != nil {
					return err
				}
			}
			env, _, err = loadEnv(s)
			if err != nil {
				return err
			}

			req := gate.Request{
				To: to, Question: question, Reason: reason,
				Actor: identity(), ActorAliases: identityAliases(),
			}
			vs := gate.CheckAdvance(env, t, req)
			if dryRun {
				return reportDryRun(cmd, env, t, to, vs, force)
			}
			if len(vs) > 0 && !force {
				return &codedError{
					code:       refusalCode(vs),
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
	f.StringVar(&reason, "reason", "", "what the ticket is waiting on, for the blocked lane")
	f.StringVar(&executedBy, "executed-by", "", "model that performed this run; ownership stays with the assignee")
	f.StringVar(&what, "what", "", "outcome: what was changed")
	f.StringVar(&why, "why", "", "outcome: why it was needed")
	f.StringVar(&resolves, "resolves", "", "outcome: how the change satisfies the definition of done")
	f.StringSliceVar(&commits, "commits", nil, "commit SHAs produced for this ticket")
	f.BoolVar(&force, "force", false, "override gate refusals; recorded in the output")
	f.BoolVar(&dryRun, "dry-run", false, "run the gates and report, writing nothing")
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

// reportDryRun says what the move would have done. The refusal is byte-identical
// to the real one, down to the exit code, so a caller branches on a dry run
// exactly as it branches on the move itself — the only difference is the line
// saying nothing was written.
func reportDryRun(cmd *cobra.Command, env gate.Env, t *ticket.Ticket, to string, vs gate.Violations, force bool) error {
	if len(vs) > 0 && !force {
		return &codedError{
			code:       refusalCode(vs),
			reason:     "gate_refused",
			message:    fmt.Sprintf("cannot move %s to %s:\n%s\n\nnothing was written (--dry-run)", ticket.Handle(t.ID), to, bullets(vs)),
			violations: vs,
		}
	}
	if g.jsonOut {
		return emit(cmd.OutOrStdout(), map[string]any{
			"ticket": ticketJSON(t, env.Lanes), "moved": false, "reason": "dry_run",
			"would_move": true, "overridden": force && len(vs) > 0,
		})
	}
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "%s → %s would be allowed; nothing was written (--dry-run)\n", ticket.Handle(t.ID), to)
	if force && len(vs) > 0 {
		fmt.Fprintf(w, "It would override %d gate refusal(s):\n%s\n", len(vs), bullets(vs))
	}
	return nil
}

// refusalCode maps a refusal to its exit code, so a dry run and the move it
// stands in for cannot report different ones.
func refusalCode(vs gate.Violations) int {
	for _, v := range vs {
		if v.Code == gate.CodeBlocked || v.Code == gate.CodeSelfBlock {
			return ExitBlocked
		}
	}
	return ExitValidation
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
	var all, perLane bool
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

			if perLane {
				return emitPerLane(cmd, env, ready)
			}

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
	cmd.Flags().BoolVar(&perLane, "per-lane", false,
		"one ticket per lane that has work waiting, in pipeline order")
	return cmd
}

// emitPerLane answers "where is work waiting", one line per lane, instead of
// "what is the single furthest-along ticket".
//
// Both questions are real and they have different answers. A board with a deep
// queue in a late lane starves every earlier lane under the default ordering: an
// agent driving 'jaira next' is handed the late ticket every time, and a step
// inserted in the middle of the pipeline never sees traffic until the queue
// ahead of it drains. Measured on a real board: 28 tickets waiting in one lane
// hid two waiting in another entirely.
//
// Lanes are reported in pipeline order, not by how far along they are, because
// this is a map of the front line rather than a recommendation. Each entry
// carries the lane's own agentic flag, so the caller can tell which of them it
// is allowed to work at all.
func emitPerLane(cmd *cobra.Command, env gate.Env, ready []*ticket.Ticket) error {
	byLane := map[string][]*ticket.Ticket{}
	for _, t := range ready {
		byLane[t.Status] = append(byLane[t.Status], t)
	}

	type entry struct {
		lane    *lane.Lane
		waiting int
		first   *ticket.Ticket
	}
	var entries []entry
	for _, l := range env.Lanes.Lanes {
		ts := byLane[l.ID]
		if len(ts) == 0 {
			continue
		}
		entries = append(entries, entry{lane: l, waiting: len(ts), first: ts[0]})
	}

	if g.jsonOut {
		out := make([]map[string]any, 0, len(entries))
		for _, e := range entries {
			out = append(out, map[string]any{
				"lane": e.lane.ID, "name": e.lane.Name, "agentic": e.lane.Agentic,
				"model_tier": e.lane.ModelTier, "waiting": e.waiting,
				"ticket": ticketJSON(e.first, env.Lanes),
			})
		}
		return emit(cmd.OutOrStdout(), map[string]any{"lanes": out, "count": len(out)})
	}

	w := cmd.OutOrStdout()
	if len(entries) == 0 {
		fmt.Fprintln(w, "Nothing is actionable. Either everything is done, blocked, or still needs a spec.")
		return nil
	}
	fmt.Fprintf(w, "%-14s %-7s %8s  %s\n", "LANE", "AGENTIC", "WAITING", "NEXT")
	for _, e := range entries {
		fmt.Fprintf(w, "%-14s %-7v %8d  %s  %s\n", e.lane.ID, e.lane.Agentic, e.waiting,
			ticket.Handle(e.first.ID), e.first.Title)
	}
	return nil
}

func sortByProgress(ts []*ticket.Ticket, env gate.Env) {
	prec := func(t *ticket.Ticket) int { return env.Lanes.Precedence(t.Status) }
	// Within one lane, a ticket that still owes the lane's declared output is
	// the one waiting to be worked; one that has produced it is waiting to be
	// *moved*, and handing it back means doing the lane's work twice. Across
	// lanes this changes nothing — precedence still finishes in-flight work
	// first — so it only decides who goes first inside a single lane, which is
	// what lets 'next --lane <id>' drain one.
	owes := func(t *ticket.Ticket) bool {
		l, ok := env.Lanes.Get(t.Status)
		return ok && l.Agentic && len(gate.OutputOwed(l, t)) > 0
	}
	// Later keys only break ties in the earlier ones: lane precedence, then
	// unworked before worked, then oldest first.
	ahead := func(b, a *ticket.Ticket) bool {
		switch {
		case prec(b) != prec(a):
			return prec(b) > prec(a)
		case owes(b) != owes(a):
			return owes(b)
		default:
			return b.ID < a.ID
		}
	}
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0; j-- {
			a, b := ts[j-1], ts[j]
			if ahead(b, a) {
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
