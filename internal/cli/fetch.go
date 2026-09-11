package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/notify"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newFetchCmd is the read half of ticket refs: one round trip, no branch, no
// checkout, no merge.
func newFetchCmd() *cobra.Command {
	var quiet bool
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch the tickets travelling on their own git refs",
		Long: `Fetches refs/jaira/tickets/* from the board's remote and reports what they say.

This is how a ticket assigned to you arrives without anyone sharing a branch:
the ref carries the whole ticket file, so it is readable with no checkout and
nothing to merge. Refs are outside refs/heads, so the default refspec does not
carry them and a teammate who does not use jaira sees none of this.

A ticket newly assigned to you also raises a desktop notification, unless
notifications are turned off in ~/.jaira/settings.json ("notify-off": true).`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := openStore(); err != nil {
				return err
			}
			if err := refs.Usable(); err != nil {
				return fail(ExitError, "no_remote",
					"this board does not carry tickets on refs: %v", err)
			}
			arrivals, err := refs.Incoming()
			if err != nil {
				return err
			}
			// Newly assigned first, then everything else: the reason to run
			// this command is the ticket somebody just handed you.
			sort.SliceStable(arrivals, func(i, j int) bool {
				a, b := arrivals[i], arrivals[j]
				if (a.Mine && a.Changed) != (b.Mine && b.Changed) {
					return a.Mine && a.Changed
				}
				return a.ID < b.ID
			})
			if !quiet {
				announceArrivals(arrivals)
			}
			departed := refs.Departed()
			stranded := strandedHere()
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"arrivals": arrivals, "departed": departed, "stranded": stranded,
				})
			}
			printArrivals(cmd.OutOrStdout(), arrivals)
			printDeparted(cmd.OutOrStdout(), departed)
			printStranded(cmd.OutOrStdout(), stranded)
			return nil
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "do not raise a desktop notification")
	return cmd
}

// announceArrivals raises one notification per ticket newly assigned to this
// user. Only Mine and only Changed: a ticket that has been sitting there
// assigned to me since last week is not news, and announcing it on every fetch
// would teach the user to ignore the notification.
func announceArrivals(arrivals []refsync.Arrival) {
	if !settings.Load().NotifyEnabled() {
		return
	}
	for _, a := range arrivals {
		if !a.Mine || !a.Changed {
			continue
		}
		title := a.Title
		if title == "" {
			title = ticket.Handle(a.ID)
		}
		body := fmt.Sprintf("%s — %s", ticket.Handle(a.ID), title)
		if a.Status != "" {
			body += " (" + a.Status + ")"
		}
		notify.Send("jaira: a ticket is yours", body)
	}
}

// strandedHere lists tickets finished on this machine that have not arrived in
// any landing branch. Resolved the same way the snapshot run resolves it, so
// the report and the reaping can never disagree about what counts as landed.
func strandedHere() []refsync.Stranded {
	set := settings.Load()
	return refs.Stranded(set.Landing(set.RemoteName(), refs.Repo.RemoteHead), set.LandingGraceInterval())
}

// printStranded names finished tickets whose branch never arrived. This is the
// price of keeping a ref alive until the ticket lands, and it has to be
// visible: otherwise an abandoned branch keeps its ref for ever and the board
// carries a finished ticket nobody can account for.
func printStranded(w io.Writer, stranded []refsync.Stranded) {
	if len(stranded) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, s := range stranded {
		fmt.Fprintf(w, "%-8s finished %d day(s) ago and has not arrived in a landing branch: %s\n",
			ticket.Handle(s.ID), s.Days, s.Title)
		fmt.Fprintf(w, "         push the branch that holds it, or give up on it with 'jaira snapshot --drop %s'\n",
			ticket.Handle(s.ID))
	}
}

// printDeparted names the tickets somebody else has taken off the board while
// a file for them is still here. It says what to run and moves nothing: which
// copy is right is the person's call, and guessing would sometimes reopen
// finished work.
func printDeparted(w io.Writer, departed []refsync.Departure) {
	if len(departed) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, d := range departed {
		fmt.Fprintf(w, "%-8s left the board elsewhere, but is still here: %s\n", ticket.Handle(d.ID), d.Title)
		fmt.Fprintf(w, "         file it too with 'jaira logbook %s', or 'jaira archive %s' — or put it back with 'jaira pull %s'\n",
			ticket.Handle(d.ID), ticket.Handle(d.ID), ticket.Handle(d.ID))
	}
}

func printArrivals(w io.Writer, arrivals []refsync.Arrival) {
	if len(arrivals) == 0 {
		fmt.Fprintln(w, "no tickets on refs yet")
		return
	}
	for _, a := range arrivals {
		marks := ""
		if a.Mine {
			marks += " @you"
		}
		if a.Changed {
			marks += " new"
		}
		if !a.Local {
			// The property worth naming: this ticket is in no branch you have.
			marks += " ref-only"
		}
		fmt.Fprintf(w, "%-8s %-12s %s%s\n", ticket.Handle(a.ID), a.Status, a.Title, marks)
	}
}
