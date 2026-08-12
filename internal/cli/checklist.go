package cli

import (
	"fmt"
	"strconv"

	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/ticket"
	"github.com/spf13/cobra"
)

func newDoDCmd() *cobra.Command {
	var (
		doing bool
		done  bool
		todo  bool
		plan  bool
	)
	cmd := &cobra.Command{
		Use:   "dod <id> <n>",
		Short: "Mark a checklist item as in progress, done, or not started",
		Long: `Sets the state of one checklist item, numbered as 'jaira show' displays it.

Three states are recognised:

  --doing   the item being worked on right now, written as [~]
  --done    finished, written as [x]
  --todo    not started, written as [ ]

An item marked as in progress is outstanding work: it does not satisfy the
definition of done, and a ticket carrying one cannot enter a terminal lane. Only
one item per checklist can be in progress, so marking a second moves the marker.

By default this addresses the Definition of Done. Use --plan for the Plan
checklist, which records the method being followed rather than the criteria for
acceptance.`,
		Args: exactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var st ticket.State
			switch {
			case doing && !done && !todo:
				st = ticket.StateDoing
			case done && !doing && !todo:
				st = ticket.StateDone
			case todo && !doing && !done:
				st = ticket.StateTodo
			default:
				return fail(ExitUsage, "usage", "pass exactly one of --doing, --done or --todo")
			}

			n, err := strconv.Atoi(args[1])
			if err != nil || n < 1 {
				return fail(ExitUsage, "usage", "item number must be a positive integer, got %q", args[1])
			}

			sec := ticket.SectionDoD
			if plan {
				sec = ticket.SectionPlan
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
				body, err := ticket.SetItemState(t.Doc().Body(), sec, n-1, st)
				if err != nil {
					return &codedError{code: ExitValidation, reason: "no_such_item", message: err.Error()}
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
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&doing, "doing", false, "mark the item as being worked on right now")
	f.BoolVar(&done, "done", false, "mark the item as finished")
	f.BoolVar(&todo, "todo", false, "mark the item as not started")
	f.BoolVar(&plan, "plan", false, "address the Plan checklist instead of the Definition of Done")
	return cmd
}
