package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// updateInProgress suppresses nudgeIfStale while 'jaira update' itself is
// running: it is the fix the nudge points at, so it must not nag about
// itself.
var updateInProgress bool

// nudgeIfStale prints one line on stderr when the running binary is newer
// than the version stamped on this board — the same hook bindDriverIfShared
// (internal/cli/share.go) uses for "this clone needs one-time setup" on the
// one path every command takes.
//
// A dev build's version says nothing about what is actually in it, so it
// never nudges: nagging a developer about their own build on every invocation
// would be pure noise. This writes to os.Stderr directly rather than through
// cmd's streams, matching bindDriverIfShared, because it must reach the
// terminal regardless of --json — stdout is reserved for the payload an agent
// parses.
func nudgeIfStale(s *ticket.Store) {
	if release.Current == "dev" || updateInProgress {
		return
	}
	dir := s.StateDir()
	if _, err := os.ReadDir(dir); err != nil && !os.IsNotExist(err) {
		// A state directory that exists but cannot be read (permissions) must
		// not turn every command into a failure over a nudge that is not
		// itself the point. A directory that simply does not exist yet is a
		// normal, unstamped board and still falls through to nudge below.
		return
	}
	if release.Stamped(dir) == release.Current {
		return
	}
	fmt.Fprintf(os.Stderr, "jaira: this board was set up by an older version of jaira — run 'jaira update' to bring it up to date\n")
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Re-apply this repository's jaira setup",
		Long: `Re-runs the same setup 'jaira init' performs — the .gitignore entry and the
jaira section in AGENTS.md/CLAUDE.md — and prints what changed since the
version that last ran it here.

This neither downloads anything nor upgrades the jaira binary itself; it
brings a board's own setup up to date with whichever binary you already have
installed.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			updateInProgress = true
			defer func() { updateInProgress = false }()
			s, err := openStore()
			if err != nil {
				return err
			}

			previous := release.Stamped(s.StateDir())
			notes := release.Since(previous)

			p := board.Prepare(s.Root)
			if p.IgnoreErr != nil {
				return fail(ExitError, "prepare_failed", "could not write .gitignore: %v", p.IgnoreErr)
			}
			if p.NoteErr != nil {
				return fail(ExitError, "prepare_failed", "could not update an agent instruction file: %v", p.NoteErr)
			}
			// Stamp last: if Prepare had failed, the board would not actually
			// be prepared by this version, and must not claim to be.
			if err := release.Stamp(s.StateDir()); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{
					"root": s.Root, "version": release.Current, "previous": previous,
					"gitignore_written": p.Ignored, "agent_notes": p.Notes,
					"notes": entriesJSON(notes),
				})
			}

			if line := announceLine(p.Notes); line != "" {
				fmt.Fprint(w, line)
			}
			if p.Ignored {
				fmt.Fprintf(w, "This board is private: .jaira/ is gitignored, so nobody else sees it.\n")
			}
			if len(notes) == 0 {
				fmt.Fprintf(w, "\nNothing has changed since the version that last set this board up.\n")
				return nil
			}
			fmt.Fprintf(w, "\nWhat's changed since then:\n")
			for _, e := range notes {
				fmt.Fprintf(w, "\n%s\n", e.Version)
				for _, c := range e.Changes {
					fmt.Fprintf(w, "  - %s\n", c)
				}
			}
			return nil
		},
	}
}

// entriesJSON renders release entries as the plain maps emit() serializes,
// rather than exporting release.Entry's field names into the JSON contract.
func entriesJSON(entries []release.Entry) []map[string]any {
	out := make([]map[string]any, 0, len(entries))
	for _, e := range entries {
		out = append(out, map[string]any{"version": e.Version, "changes": e.Changes})
	}
	return out
}
