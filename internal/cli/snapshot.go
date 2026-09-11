package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/snapshot"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newSnapshotCmd runs the board's backup, and is also what a detached child
// runs when one is due.
func newSnapshotCmd() *cobra.Command {
	var drop string
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Write the board to its snapshot branch, and clear the refs of landed tickets",
		Long: `Writes every ticket that is on a ref into the snapshot branch, then removes the
refs of tickets that have arrived in a landing branch.

The snapshot is a backup, not the storage. The working state is always on the
refs, and every participant's clone holds the ones it has fetched, so the board
survives any one machine; this branch is for whoever clones for the first time
and has no refs at all. It is written with git plumbing and never checked out,
so a run cannot disturb whatever you are in the middle of.

The files go under board/, not .jaira/tickets/ — at the same path, merging this
branch by accident would collide with every working ticket at once.

Removing a ref happens here and nowhere else, immediately after the write: at
that moment the ticket is in the snapshot and in the branch it landed in, so
there is nothing left to lose. A ticket counts as landed only when it is filed
away — in the logbook or the archive — in one of the landing branches from
~/.jaira/settings.json ("landing-branches", globs allowed). With no landing
branch resolvable, nothing is removed at all: a ref left standing costs
nothing, one removed by mistake takes away the visibility this exists for.

Normally nobody runs this: it happens in the background when the last snapshot
is older than three days.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if err := refs.Usable(); err != nil {
				return fail(ExitError, "no_remote",
					"this board does not carry tickets on refs: %v", err)
			}
			runner := snapshotRunner()

			if drop != "" {
				id := resolveID(s, drop)
				if err := runner.Drop(id); err != nil {
					return err
				}
				refs.ForgetDeparted(id)
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"dropped": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: ref removed without waiting for it to land\n", ticket.Handle(id))
				fmt.Fprintf(cmd.OutOrStdout(), "  the ticket is now only wherever its file is — commit that file, or it is nowhere\n")
				return nil
			}

			// Stamped before the run, not after: a run that fails must not make
			// the next command try again at once, or an unreachable remote turns
			// into a retry on every command.
			_ = snapshot.WriteStamp(s.StateDir(), snapshot.Stamp{RanAt: time.Now().UTC()})
			res, err := runner.Run()
			if err != nil {
				return err
			}
			_ = snapshot.WriteStamp(s.StateDir(), snapshot.Stamp{RanAt: time.Now().UTC(), Note: res.Commit})

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), res)
			}
			if res.Unchanged {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already current: %d ticket(s), nothing written\n", res.Branch, res.Tickets)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %d ticket(s)", res.Branch, res.Tickets)
				if len(res.Added) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), ", added %s", strings.Join(handles(res.Added), " "))
				}
				if len(res.Removed) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), ", removed %s", strings.Join(handles(res.Removed), " "))
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
			if len(res.Reaped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "landed and cleared from the refs: %s\n", strings.Join(handles(res.Reaped), " "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&drop, "drop", "", "remove one ticket's ref without waiting for it to land")
	return cmd
}

// snapshotRunner builds the runner from the settings, resolving the landing
// branches the same way for the command and for the background run.
func snapshotRunner() *snapshot.Runner {
	set := settings.Load()
	return &snapshot.Runner{
		Repo:            refs.Repo,
		Branch:          set.SnapshotBranchName(),
		LandingBranches: set.Landing(set.RemoteName(), refs.Repo.RemoteHead),
	}
}

// maybeSnapshot spawns a detached snapshot when one is due, and never waits for
// it. Called once after a command has finished, next to the ref flush.
//
// Nothing here touches the network: deciding is a file read, and the run itself
// happens in another process. That is the same rule the release check follows —
// nothing you type ever waits on a remote for something nobody asked for.
func maybeSnapshot(s *ticket.Store) {
	if s == nil || refs.Usable() != nil {
		return
	}
	snapshot.SpawnRun(s.StateDir(), s.Root, settings.Load().SnapshotInterval())
}

func handles(ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, ticket.Handle(id))
	}
	return out
}
