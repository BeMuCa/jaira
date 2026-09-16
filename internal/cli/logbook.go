package cli

// This file implements 'jaira logbook <id>' and 'jaira logbook --all' —
// taking a finished ticket, or the whole terminal lane, off the board with
// its commits stamped down, into a dated folder that records who finished
// what on which day. It was called 'jaira sync' before it was
// released: that name implied a server this tool does not have, and collided
// with 'jaira sync-tasks' (sync.go), which mirrors an agent's task list into
// the backlog and is unrelated.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newLogbookCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "logbook [id]",
		Short: "Take finished tickets off the board with their commits stamped down, or list the logbook",
		Long: `Moves a terminal-lane ticket into .jaira/logbook/<initials>-<yyyymmdd>/, after
stamping it with every commit git can find for it. With no argument, lists
the logbook.

--all files everything that has reached the terminal lane into today's folder,
which is the usual way: finished tickets pile up there, and whoever enters
their hours cuts them in one go.

Filing is a decision, and nothing does it for you. Reaching a terminal lane
says the work is accepted; it says nothing about whether anybody is ready to
account for it, which happens days later and covers a set somebody assembles.
A board that filed on its own once swept forty-nine other people's tickets
into a commit named after a single handle.

The folder is the record: who finished what, on which day. Leaving the board
is the moment every commit is finally known, so it is stamped here rather
than left to whoever remembers to run 'jaira set'. 'jaira restore <file>'
brings a logged ticket back, the same as an archived one.

Naming a milestone instead of a ticket files the milestone: its file moves
into .jaira/logbook/<initials>-<yyyymmdd>/milestones/, 'jaira milestone ls'
stops naming it, no card carries its colour and the board's M filter forgets
it. Its ref stays up carrying the status "filed": that line is what keeps a
milestone off a board, so a clone that already has the file has it marked on
the next fetch and stops showing the group, and the name stays taken until
somebody restores it. Only a milestone you name by hand — --all sweeps
the terminal lane and never takes a group with it. A milestone is filed only
once every ticket in it has reached the terminal lane or left the board.

'jaira archive' is for a ticket that is not being worked — abandoned,
duplicate, obsolete — and works from any lane. This command is for finished
work and refuses a ticket that has not reached the terminal lane.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fail(ExitUsage, "usage", "logbook takes at most one ticket id, received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()

			switch {
			case all && len(args) > 0:
				return fail(ExitUsage, "usage", "--all files the whole terminal lane; naming a ticket as well says two different things")
			case all:
				return logbookAll(s, w, cmd.ErrOrStderr())
			case len(args) == 0:
				return listLogbook(s, w)
			}
			return logbookOut(s, args[0], w)
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "file everything in the terminal lane into today's folder")
	return cmd
}

// logbookAll files the whole terminal lane, which is the cut somebody makes
// when they enter their hours.
//
// It is the same sweep a lane can be told to do on entry, moved to where a
// person asks for it: the set is the same, the moment is not, and the moment
// was the problem.
func logbookAll(s *ticket.Store, w, errw io.Writer) error {
	// loadEnv already loads the lanes, and it is the only loader that prints
	// the warnings they carry. Loading them a second time here dropped those
	// warnings on the floor and named the same condition differently from
	// logbookOut below.
	env, _, err := loadEnv(s)
	if err != nil {
		return err
	}
	terminal := env.Lanes.Terminal()
	if terminal == nil {
		return fail(ExitValidation, "not_terminal", "this board has no terminal lane, so nothing can be finished into the logbook")
	}
	var skipped string
	filed, err := s.FileLane(terminal.ID, logbookFolder(), func(t *ticket.Ticket) error {
		_, err := s.StampCommits(t, env.DeriveCommits)
		return err
	})
	if err != nil {
		// One unreadable file must not hold the rest of the cut hostage: the
		// readable tickets still leave, and the problem is reported rather
		// than swallowed. Refusing the whole cut would leave somebody entering
		// hours with no way through but 'git mv'.
		var pe *ticket.PartialError
		if !errors.As(err, &pe) {
			return err
		}
		skipped = pe.Error()
		fmt.Fprintf(errw, "jaira: warning: %v\n", pe)
	}
	if g.jsonOut {
		// One rendering of a filed ticket, shared with the sweep 'move' reports
		// (flow.go, trimmedJSON): both answer "which tickets left the board",
		// and two spellings of that answer drift the moment one gains a field.
		// The handle it carries is why the reader wants it — the next thing an
		// agent does after a cut is name what it filed in a commit message,
		// which is written by handle.
		res := map[string]any{"filed": trimmedJSON(filed), "count": len(filed), "lane": terminal.ID}
		// A cut that skipped something is still a successful cut, so the
		// skipping rides along with the result rather than replacing it.
		// 'move' carries its sweep failure the same way (flow.go, trim_error),
		// and prose on stderr is the one channel a --json reader does not read.
		if skipped != "" {
			res["trim_error"] = skipped
		}
		return emit(w, res)
	}
	if len(filed) == 0 {
		fmt.Fprintf(w, "nothing in %s to file\n", terminal.ID)
		return nil
	}
	fmt.Fprintf(w, "filed %d ticket(s) from %s:\n", len(filed), terminal.ID)
	for _, f := range filed {
		fmt.Fprintf(w, "  %-8s %s\n", ticket.Handle(f.ID), filepath.Base(f.Path))
	}
	fmt.Fprintf(w, "restore one with 'jaira restore <file>'\n")
	return nil
}

func listLogbook(s *ticket.Store, w io.Writer) error {
	names, err := logbookNames(s)
	if err != nil {
		return err
	}
	if g.jsonOut {
		return emit(w, map[string]any{"logbook": names, "count": len(names)})
	}
	if len(names) == 0 {
		fmt.Fprintf(w, "The logbook is empty.\n")
		return nil
	}
	for _, n := range names {
		fmt.Fprintf(w, "%s\n", n)
	}
	fmt.Fprintf(w, "\n%d in the logbook. Bring one back with 'jaira restore <file>'.\n", len(names))
	return nil
}

func logbookOut(s *ticket.Store, idArg string, w io.Writer) error {
	t, err := s.Load(idArg)
	if err != nil {
		// A ticket first, a milestone second. A milestone name is chosen
		// freely and could be spelled like a handle, and the ticket is both
		// the older and the far commoner meaning of an argument here.
		name, nerr := milestoneNamed(s, idArg)
		switch {
		case nerr == nil:
			return logbookMilestone(s, name, w)
		case errors.Is(nerr, os.ErrNotExist):
			// The file has left the board, which is exactly what filing does
			// to it: this tree filed it, or this clone only ever saw the
			// marked ref. milestone.Load answers with the same ErrNotExist a
			// name nobody ever used answers with, so without this the reader
			// is told the name is not a ticket — in the two states where
			// create and add/rm refuse it by name. This is the door's only
			// route past the ms.Filed() check below, which needs an ms.
			if err := refuseIfFiled(s, name,
				"filing it again here would only stamp today's folder on somebody else's record of it",
				", and only a milestone that is on the board can be filed"); err != nil {
				return err
			}
		}
		return err
	}
	env, _, err := loadEnv(s)
	if err != nil {
		return err
	}

	term := env.Lanes.Terminal()
	if term == nil {
		return fail(ExitValidation, "not_terminal",
			"no terminal lane is installed, so there is nowhere for %s to be logged from", ticket.Handle(t.ID))
	}
	if t.Status != term.ID {
		return fail(ExitValidation, "not_terminal",
			"%s is in %q, not the terminal lane %q — move it there first with 'jaira move %s --to %s'",
			ticket.Handle(t.ID), t.Status, term.ID, ticket.Handle(t.ID), term.ID)
	}

	// Stamp before moving: this is the moment every commit is finally known,
	// and the commits belong to the ticket record whether or not the move
	// that follows succeeds.
	merged, err := s.StampCommits(t, env.DeriveCommits)
	if err != nil {
		return err
	}

	folder := logbookFolder()
	dst, err := s.Logbook(t.ID, folder)
	if err != nil {
		return err
	}

	if g.jsonOut {
		return emit(w, map[string]any{
			"logged": true, "id": t.ID, "handle": ticket.Handle(t.ID),
			"path": dst, "file": filepath.Base(dst), "commits": merged,
		})
	}
	fmt.Fprintf(w, "Logged %s  %s\n", ticket.Handle(t.ID), t.Title)
	fmt.Fprintf(w, "Stamped %d commit(s). Moved to %s — restore it with 'jaira restore %s'.\n",
		len(merged), filepath.Join(ticket.DirName, ticket.LogbookSubdir, folder), filepath.Base(dst))
	if len(merged) == 0 {
		fmt.Fprintf(w, "No commits were found for this ticket. Record them by hand with 'jaira set %s commits=<sha>'.\n",
			ticket.Handle(t.ID))
	}
	return nil
}

// milestoneNamed resolves an argument to a milestone on this board, or reports
// why it is not one. Only an exact name: a milestone is a name from a closed
// set, the same way the --milestone filter treats it.
func milestoneNamed(s *ticket.Store, arg string) (string, error) {
	name, _, err := milestone.NormalizeName(arg)
	if err != nil {
		return "", err
	}
	if _, err := milestone.Load(s.Root, name); err != nil {
		// The name travels with the error: a caller that has to ask whether
		// this name was filed needs the normalized spelling, and the file
		// being gone is one of the answers it asks about.
		return name, err
	}
	return name, nil
}

// logbookMilestone files a milestone: off the board, into the logbook, with
// its ref left standing and marked.
//
// The ref is the part worth reading twice. Taking it down would be the obvious
// move and is the wrong one: refsync.IncomingMilestones writes to disk every
// milestone the refs carry, so a removed ref frees the name for a second
// milestone with the same identity, while a ref that says "filed" is read by
// every clone — it stops the file coming back and holds the name until
// 'jaira restore' gives it up.
func logbookMilestone(s *ticket.Store, name string, w io.Writer) error {
	env, _, err := loadEnv(s)
	if err != nil {
		return err
	}
	term := env.Lanes.Terminal()
	if term == nil {
		return fail(ExitValidation, "not_terminal",
			"no terminal lane is installed, so there is nowhere for milestone %q to be logged from", name)
	}

	unlock, err := s.Lock(milestoneLockName)
	if err != nil {
		return err
	}
	defer unlock()

	ms, err := milestone.Load(s.Root, name)
	if err != nil {
		return err
	}
	// A file that already says filed is one a fetch wrote back here, or one a
	// filing marked and got no further with. Either way it is off the board
	// already, and filing it a second time would only stamp today's folder on
	// somebody else's record of it.
	//
	// This check is the route where the marked file is lying here, and only
	// that one: it took a fetch of somebody else's filing to put it here, so
	// this tree's logbook holds nothing and 'jaira restore' here would only
	// answer that the file is not in the archive. The copy that can come back
	// is in the tree that filed it — which is what refuseFiledOnDisk says, so
	// this door tells the reader what the other two tell them.
	//
	// The routes where the file is NOT here are answered above, before
	// milestone.Load is even asked for an ms: this check cannot see them.
	if ms.Filed() {
		return refuseFiledOnDisk(s.Root, name,
			"there is no copy of it here to bring back, so filing it again would only stamp today's folder on somebody else's record of it")
	}
	// The same gate a ticket passes, asked of a group: a milestone with
	// unfinished work in it is a plan somebody is still working, and filing it
	// takes the plan off the board while the work stays on it.
	var open []string
	for _, id := range ms.Members() {
		t, err := s.Load(id)
		if err != nil {
			// Already off the board — filed or archived — which is as finished
			// as this can ask for.
			continue
		}
		if t.Status != term.ID {
			open = append(open, fmt.Sprintf("%s (%s)", ticket.Handle(id), t.Status))
		}
	}
	if len(open) > 0 {
		return fail(ExitValidation, "milestone_unfinished",
			"milestone %q still holds work that has not reached %q: %s — finish or take those out with 'jaira milestone rm %s <id>' first",
			name, term.ID, strings.Join(open, ", "), name)
	}

	// Mark, then put the marked file on the ref, then move it: recordMilestone
	// reads the file from the board, so writing has to come before moving.
	// Nothing rides on it beyond that — the marked line is what takes the
	// milestone off the board, so a move that fails here leaves a file that is
	// already invisible to the board and already filed on its ref.
	ms.SetStatus(milestone.StatusFiled)
	if err := ms.Save(s.Root); err != nil {
		return err
	}
	recordMilestone(ms)

	folder := logbookFolder()
	dst, err := s.LogbookMilestone(name, folder)
	if err != nil {
		return err
	}

	if g.jsonOut {
		return emit(w, map[string]any{
			"logged": true, "milestone": name, "path": dst, "file": filepath.Base(dst),
			"count": len(ms.Members()),
		})
	}
	fmt.Fprintf(w, "Logged milestone %s (%d ticket(s)).\n", name, len(ms.Members()))
	fmt.Fprintf(w, "Moved to %s — restore it with 'jaira restore %s'.\n",
		filepath.Join(ticket.DirName, ticket.LogbookSubdir, folder, ticket.MilestonesSubdir), filepath.Base(dst))
	fmt.Fprintf(w, "Its ref stays up marked %q, so the name stays taken until you do.\n", milestone.StatusFiled)
	return nil
}

// logbookFolder names the dated folder a ticket lands in: who took the ticket
// off and when, so the folder is a readable record of one person's sweep
// rather than a bare filename nobody can attribute.
func logbookFolder() string {
	return fmt.Sprintf("%s-%s", coreidentity.Initials(identity()), time.Now().Format("20060102"))
}

func logbookNames(s *ticket.Store) ([]string, error) {
	entries, err := os.ReadDir(s.LogbookDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sub, err := os.ReadDir(filepath.Join(s.LogbookDir(), e.Name()))
		if err != nil {
			continue
		}
		for _, f := range sub {
			// A filed milestone sits one level deeper, in milestones/. It is
			// listed with the rest or the file is there and the list denies it.
			if f.IsDir() && f.Name() == ticket.MilestonesSubdir {
				inner, err := os.ReadDir(filepath.Join(s.LogbookDir(), e.Name(), f.Name()))
				if err != nil {
					continue
				}
				for _, mf := range inner {
					if !mf.IsDir() && strings.HasSuffix(mf.Name(), ".md") {
						out = append(out, filepath.Join(e.Name(), f.Name(), mf.Name()))
					}
				}
				continue
			}
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
				out = append(out, filepath.Join(e.Name(), f.Name()))
			}
		}
	}
	sort.Strings(out)
	return out, nil
}
