package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/spf13/cobra"
)

func newDoDCmd() *cobra.Command {
	var (
		doing  bool
		done   bool
		todo   bool
		plan   bool
		option string
		add    []string
		proof  string
	)
	cmd := &cobra.Command{
		Use:   "dod <id> [n]",
		Short: "Write or mark checklist items",
		Long: `Sets the state of one checklist item, numbered as 'jaira show' displays it.

Three states are recognised:

  --doing   the item being worked on right now, written as [~]
  --done    finished, written as [x]
  --todo    not started, written as [ ]

An item marked as in progress is outstanding work: it does not satisfy the
definition of done, and a ticket carrying one cannot enter a terminal lane. Only
one item per checklist can be in progress, so marking a second moves the marker.

Add items with --add, which may be repeated. An added item lands inside the
right section — a checkbox only counts under its own heading, so appending to
the end of the file would put it somewhere it does not count:

  jaira dod <id> --plan --add "read the exporter" --add "design the chunking"

By default this addresses the Definition of Done. Use --plan for the Plan
checklist, which records the method being followed rather than the criteria for
acceptance.

Record the evidence for an item in the same call that ticks it, with --proof:

  jaira dod <id> 2 --done --proof "internal/tui/browse.go:144 ... covered by TestX"

--proof can also be passed on its own, with no state flag, to record evidence
without changing the marker. Setting it again on the same item replaces the
line rather than stacking a second one.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fail(ExitUsage, "usage",
					"dod takes a ticket id and, unless adding, an item number\n\nUsage: %s", cmd.UseLine())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sec := ticket.SectionDoD
			if plan {
				sec = ticket.SectionPlan
			}
			if proof != "" && (option != "" || len(add) > 0) {
				return fail(ExitUsage, "usage", "--proof addresses a single existing item; it cannot be combined with --add or --option")
			}
			if option != "" {
				return setOption(cmd, args[0], option, !todo)
			}
			if len(add) > 0 {
				if doing || done || todo || len(args) > 1 {
					return fail(ExitUsage, "usage", "--add writes new items; it does not take a number or a state")
				}
				return addChecklistItems(cmd, args[0], sec, add)
			}
			if len(args) != 2 {
				return fail(ExitUsage, "usage", "which item? pass its number, as 'jaira show' displays it")
			}

			// With --proof present, a state flag is optional: passing none records
			// only the proof and leaves the marker alone. Without --proof, the rule
			// stays what it always was — exactly one state, so a plain 'dod <id> <n>'
			// with nothing else is a usage error rather than a silent no-op.
			hasState := doing || done || todo
			var st ticket.State
			switch {
			case hasState && doing && !done && !todo:
				st = ticket.StateDoing
			case hasState && done && !doing && !todo:
				st = ticket.StateDone
			case hasState && todo && !doing && !done:
				st = ticket.StateTodo
			case hasState:
				return fail(ExitUsage, "usage", "pass exactly one of --doing, --done or --todo")
			case proof == "":
				return fail(ExitUsage, "usage", "pass exactly one of --doing, --done or --todo")
			}

			n, err := strconv.Atoi(args[1])
			if err != nil || n < 1 {
				return fail(ExitUsage, "usage", "item number must be a positive integer, got %q", args[1])
			}

			s, err := openStore()
			if err != nil {
				return err
			}
			cur, err := s.Load(args[0])
			if err != nil {
				return err
			}
			lanes, err := lane.Load()
			if err != nil {
				return err
			}
			// Matches 'set' and 'move': a ticket in a lane this installation does
			// not have is read-only on every mutation path, not just some of them.
			if _, known := lanes.Get(cur.Status); !known && cur.Status != "" {
				return &codedError{
					code:   ExitValidation,
					reason: "unknown_lane",
					message: fmt.Sprintf(
						"%s sits in unrecognized lane %q and is read-only; install that lane file first",
						ticket.Handle(cur.ID), cur.Status),
				}
			}

			t, err := s.Mutate(args[0], func(t *ticket.Ticket) error {
				body := t.Doc().Body()
				if hasState {
					next, err := ticket.SetItemState(body, sec, n-1, st)
					if err != nil {
						return &codedError{code: ExitValidation, reason: "no_such_item", message: err.Error()}
					}
					body = next
				}
				if proof != "" {
					next, err := ticket.SetItemProof(body, sec, n-1, proof)
					if err != nil {
						return &codedError{code: ExitValidation, reason: "no_such_item", message: err.Error()}
					}
					body = next
				}
				t.Doc().SetBody(body)
				return nil
			})
			if err != nil {
				return err
			}
			// The body changed, so readiness may have too. Re-read rather than
			// reasoning about it: the reload repopulates the parsed checklists.
			t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
				return ticket.SetReady(t.Doc(), gate.Ready(t))
			})
			if err != nil {
				return err
			}

			items := t.DoDItems
			if plan {
				items = t.PlanItems
			}
			if g.jsonOut {
				env, _, err := loadEnv(s)
				if err != nil {
					return err
				}
				return emit(cmd.OutOrStdout(), ticketJSON(t, env.Lanes))
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "%s %s\n", ticket.Handle(t.ID), sec)
			for i, it := range items {
				marker := " "
				if i == n-1 {
					marker = "→"
				}
				fmt.Fprintf(w, "%s %d. [%s] %s\n", marker, i+1, it.State.Marker(), it.Text)
				if it.Proof != "" {
					fmt.Fprintf(w, "       proof: %s\n", it.Proof)
				}
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&doing, "doing", false, "mark the item as being worked on right now")
	f.BoolVar(&done, "done", false, "mark the item as finished")
	f.BoolVar(&todo, "todo", false, "mark the item as not started")
	f.BoolVar(&plan, "plan", false, "address the Plan checklist instead of the Definition of Done")
	f.StringArrayVar(&add, "add", nil, "append a new item; repeat for several, in order")
	f.StringVar(&option, "option", "", "tick a step option for this ticket, e.g. --option planning (with --todo to untick)")
	f.StringVar(&proof, "proof", "", "record evidence for the item; combine with --done to tick and record in one call")
	return cmd
}

// addChecklistItems writes new items into a section.
func addChecklistItems(cmd *cobra.Command, id string, sec ticket.Section, items []string) error {
	s, err := openStore()
	if err != nil {
		return err
	}
	t, err := s.Mutate(id, func(t *ticket.Ticket) error {
		body := t.Doc().Body()
		for _, text := range items {
			next, err := ticket.AddItem(body, sec, text)
			if err != nil {
				return fail(ExitUsage, "empty_item", "%v", err)
			}
			body = next
		}
		t.Doc().SetBody(body)
		return nil
	})
	if err != nil {
		return err
	}
	// The definition of done feeds readiness, so it has to be recomputed.
	t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
		return ticket.SetReady(t.Doc(), gate.Ready(t))
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
	list := t.DoDItems
	if sec == ticket.SectionPlan {
		list = t.PlanItems
	}
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "%s %s\n", ticket.Handle(t.ID), sec)
	for i, it := range list {
		fmt.Fprintf(w, "  %d. [%s] %s\n", i+1, it.State.Marker(), it.Text)
	}
	return nil
}

// setOption ticks or unticks one entry in a ticket's Options checklist, which is
// what decides whether an optional step is part of that ticket's path.
func setOption(cmd *cobra.Command, id, name string, on bool) error {
	s, err := openStore()
	if err != nil {
		return err
	}
	t, err := s.Mutate(id, func(t *ticket.Ticket) error {
		body := t.Doc().Body()
		// Add the option if the ticket has never heard of it, so turning a step on
		// does not require hand-editing the body first.
		found := -1
		for i, o := range ticket.ParseOptions(body) {
			if strings.EqualFold(strings.TrimSpace(o.Text), strings.TrimSpace(name)) {
				found = i
				break
			}
		}
		if found < 0 {
			next, err := ticket.AddItem(body, ticket.SectionOptions, name)
			if err != nil {
				return err
			}
			body = next
			found = len(ticket.ParseOptions(body)) - 1
		}
		st := ticket.StateDone
		if !on {
			st = ticket.StateTodo
		}
		next, err := ticket.SetItemState(body, ticket.SectionOptions, found, st)
		if err != nil {
			return err
		}
		t.Doc().SetBody(next)
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
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "%s options\n", ticket.Handle(t.ID))
	for _, o := range t.Options {
		fmt.Fprintf(w, "  [%s] %s\n", o.State.Marker(), o.Text)
	}
	return nil
}
