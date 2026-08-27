package cli

// This file implements 'jaira logbook <id>' — taking a finished ticket off
// the board with its commits stamped down, into a dated folder that records
// who finished what on which day. It was called 'jaira sync' before it was
// released: that name implied a server this tool does not have, and collided
// with 'jaira sync-tasks' (sync.go), which mirrors an agent's task list into
// the backlog and is unrelated.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newLogbookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logbook [id]",
		Short: "Take a finished ticket off the board with its commits stamped down, or list the logbook",
		Long: `Moves a terminal-lane ticket into .jaira/logbook/<initials>-<yyyymmdd>/, after
stamping it with every commit git can find for it. With no argument, lists
the logbook.

The folder is the record: who finished what, on which day. Leaving the board
is the moment every commit is finally known, so it is stamped here rather
than left to whoever remembers to run 'jaira set'. 'jaira restore <file>'
brings a logged ticket back, the same as an archived one.

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

			if len(args) == 0 {
				return listLogbook(s, w)
			}
			return logbookOut(s, args[0], w)
		},
	}
	return cmd
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
		return err
	}
	env, _, err := loadEnv(s)
	if err != nil {
		return err
	}

	term := env.Lanes.Terminal()
	if term == nil {
		return &codedError{
			code:   ExitValidation,
			reason: "not_terminal",
			message: fmt.Sprintf(
				"no terminal lane is installed, so there is nowhere for %s to be logged from", ticket.Handle(t.ID)),
		}
	}
	if t.Status != term.ID {
		return &codedError{
			code:   ExitValidation,
			reason: "not_terminal",
			message: fmt.Sprintf(
				"%s is in %q, not the terminal lane %q — move it there first with 'jaira move %s --to %s'",
				ticket.Handle(t.ID), t.Status, term.ID, ticket.Handle(t.ID), term.ID),
		}
	}

	// Stamp before moving: this is the moment every commit is finally known,
	// and the commits belong to the ticket record whether or not the move
	// that follows succeeds.
	merged, err := stampCommits(s, t, env.DeriveCommits)
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

// logbookFolder names the dated folder a ticket lands in: who took the ticket
// off and when, so the folder is a readable record of one person's sweep
// rather than a bare filename nobody can attribute.
func logbookFolder() string {
	return fmt.Sprintf("%s-%s", coreidentity.Initials(identity()), time.Now().Format("20060102"))
}

// stampCommits writes the derived commit union onto the ticket and returns
// what was written. Derived shas come first, in git order; any sha already
// recorded that the derivation did not find is appended rather than dropped —
// a sha a person wrote down deliberately is evidence this tool has no
// business discarding. derive may be nil, the same "no derivation on offer"
// convention core/gate uses.
func stampCommits(s *ticket.Store, t *ticket.Ticket, derive func(*ticket.Ticket) []string) ([]string, error) {
	var derived []string
	if derive != nil {
		derived = derive(t)
	}
	merged := append([]string{}, derived...)
	for _, c := range t.Commits {
		if c = strings.TrimSpace(c); c != "" && !contains(merged, c) {
			merged = append(merged, c)
		}
	}
	if _, err := s.Mutate(t.ID, func(t *ticket.Ticket) error {
		return t.Doc().SetList(ticket.FieldCommits, merged)
	}); err != nil {
		return nil, err
	}
	return merged, nil
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
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
				out = append(out, filepath.Join(e.Name(), f.Name()))
			}
		}
	}
	sort.Strings(out)
	return out, nil
}
