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
		doing      bool
		done       bool
		todo       bool
		superseded bool
		plan       bool
		option     string
		add        []string
		proof      string
		text       string
	)
	cmd := &cobra.Command{
		Use:   "dod <id> [n]",
		Short: "Write or mark checklist items",
		Long: `Sets the state of one checklist item, numbered as 'jaira show' displays it.

Four states are recognised:

  --doing        the item being worked on right now, written as [~]
  --done         finished, written as [x]
  --todo         not started, written as [ ]
  --superseded   it did not happen and will not, written as [-]

An item marked as in progress is outstanding work: it does not satisfy the
definition of done, and a ticket carrying one cannot enter a terminal lane. Only
one item per checklist can be in progress, so marking a second moves the marker.

A superseded item is retired, not achieved: it stops blocking completion, and
nothing reports it as done. Use it when a criterion was replaced or dropped —
ticking it instead leaves a [x] beside work nobody ever did.

Reword an item with --text, which leaves its state and its proof alone:

  jaira dod <id> 3 --text "the counter shows n of m, superseded items included"

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
			if (proof != "" || text != "") && (option != "" || len(add) > 0) {
				return fail(ExitUsage, "usage", "--proof and --text address a single existing item; they cannot be combined with --add or --option")
			}
			if option != "" {
				return setOption(cmd, args[0], option, !todo)
			}
			if len(add) > 0 {
				if doing || done || todo || superseded || len(args) > 1 {
					return fail(ExitUsage, "usage", "--add writes new items; it does not take a number or a state")
				}
				return addChecklistItems(cmd, args[0], sec, add)
			}
			if len(args) != 2 {
				return fail(ExitUsage, "usage", "which item? pass its number, as 'jaira show' displays it")
			}

			// With --proof or --text present, a state flag is optional: passing none
			// records only that and leaves the marker alone. Without either, the
			// rule stays what it always was — exactly one state, so a plain
			// 'dod <id> <n>' with nothing else is a usage error rather than a
			// silent no-op.
			states := 0
			for _, on := range []bool{doing, done, todo, superseded} {
				if on {
					states++
				}
			}
			const oneState = "pass exactly one of --doing, --done, --todo or --superseded"
			var st ticket.State
			switch {
			case states > 1:
				return fail(ExitUsage, "usage", oneState)
			case doing:
				st = ticket.StateDoing
			case done:
				st = ticket.StateDone
			case todo:
				st = ticket.StateTodo
			case superseded:
				st = ticket.StateSuperseded
			case proof == "" && text == "":
				return fail(ExitUsage, "usage", oneState)
			}
			hasState := states == 1

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
			lanes, err := lane.Load(s.Root)
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
				// Before the proof, so evidence recorded in the same call attaches to
				// the wording it was written against.
				if text != "" {
					next, err := ticket.SetItemText(body, sec, n-1, text)
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
	f.BoolVar(&superseded, "superseded", false, "retire the item: it did not happen and will not, written as [-]")
	f.BoolVar(&plan, "plan", false, "address the Plan checklist instead of the Definition of Done")
	f.StringArrayVar(&add, "add", nil, "append a new item; repeat for several, in order")
	f.StringVar(&option, "option", "", "tick a step option for this ticket, e.g. --option planning (with --todo to untick)")
	f.StringVar(&proof, "proof", "", "record evidence for the item; combine with --done to tick and record in one call")
	f.StringVar(&text, "text", "", "reword the item, leaving its state and its proof alone")
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
