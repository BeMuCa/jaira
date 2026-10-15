package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/lane"
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
//
// Whether a newer *release* is published is a separate question this
// function deliberately does not answer: that lives in the TUI's persistent
// version indicator (internal/tui) instead, so CLI commands — the surface
// agents and scripts drive — stay silent about it. See core/selfupdate/
// cache.go for the cache both surfaces would read if a CLI-side nudge were
// ever wanted again.
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
	var agentFile string
	var setup updateSetup
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Re-apply this repository's jaira setup",
		Long: `Regenerates the jaira section in AGENTS.md/CLAUDE.md for this board's
current lanes, and prints what changed since the version that last ran it
here.

--agent-file picks where that section lives: agents (AGENTS.md), claude
(CLAUDE.md) or both. The section is taken out of the other file, and the
choice holds from then on, because jaira writes only into the files that
already carry the section. Claude Code reads AGENTS.md only where there is no
CLAUDE.md, so 'agents' deletes a CLAUDE.md left empty and puts an
'@AGENTS.md' import into one that still holds anything else.

It does not touch .gitignore. 'jaira init' makes a new board private; whether
this board is shared is a decision already taken, and re-adding the ignore
entry to a shared board would silently stop new tickets reaching teammates.

--roles brings the roles, the cockpit mod and the jaira skill up to date in every
skills directory that already holds the roles: this repository's and your
personal one. Outside a board it does the personal one only. --lanes replaces
every lane file on this board that is an older version of a lane jaira ships.
A file you edited — one that is no version jaira ever shipped — is never
replaced: it is named, and the command exits 3. Without --roles or --lanes only
the jaira section is written, as before, and it says when roles or lanes are
behind.

--dry-run says what would be written and what would be left alone, and writes
nothing — not even the version stamp. With --lanes the preview of the jaira
section is computed from the lanes as they are now.

This neither downloads anything nor upgrades the jaira binary itself; it
brings a board's own setup up to date with whichever binary you already have
installed.`,
		Args: noArgs(),
		// A named result so a skipped file can turn the exit into 3 after the notes
		// below have printed: they return nil from two places.
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			switch agentFile {
			case "", board.AgentFileAgents, board.AgentFileClaude, board.AgentFileBoth:
			default:
				return fail(ExitValidation, "bad_agent_file",
					"--agent-file %q: want agents, claude or both", agentFile)
			}
			updateInProgress = true
			defer func() { updateInProgress = false }()
			s, err := openStore()
			if errors.Is(err, ticket.ErrNoStore) && setup.roles && !setup.lanes && agentFile == "" {
				return updatePersonalRoles(cmd, setup.dry)
			}
			if err != nil {
				return err
			}
			defer func() {
				if err == nil && setup.skipped() {
					err = editedError()
				}
			}()

			previous := release.Stamped(s.StateDir())
			notes := release.Since(previous)

			if setup.lanes {
				// Before the Load below, so the section describes the refreshed lanes.
				if setup.laneResults, err = lane.RefreshShipped(s.Root, setup.dry); err != nil {
					return err
				}
			}
			// Re-generated with this board's current lanes: adopting a lane is
			// exactly when the note stops matching the board.
			boardLanes, _ := lane.Load(s.Root)
			var notesWritten []string
			var noteErr error
			if setup.dry {
				notesWritten, noteErr = board.PreviewAgentFiles(s.Root, agentFile, boardLanes.Facts())
			} else if agentFile != "" {
				notesWritten, noteErr = board.ChooseAgentFiles(s.Root, agentFile, boardLanes.Facts())
			} else {
				notesWritten, noteErr = board.AnnounceInAgentFiles(s.Root, boardLanes.Facts())
			}
			if noteErr != nil {
				return fail(ExitError, "prepare_failed", "could not update an agent instruction file: %v", noteErr)
			}
			if setup.roles {
				if setup.roleRuns, err = refreshRoles(s.Root, setup.dry); err != nil {
					return err
				}
			}
			// Stamp last: if the note had failed, the board would not actually
			// be prepared by this version, and must not claim to be.
			if !setup.dry {
				if err := release.Stamp(s.StateDir()); err != nil {
					return err
				}
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, setup.json(map[string]any{
					"root": s.Root, "version": release.Current, "previous": previous,
					"agent_notes": notesWritten,
					"notes":       entriesJSON(notes),
				}))
			}

			if setup.lanes {
				printLaneResults(w, setup.laneResults, setup.dry)
			}
			if setup.dry {
				if len(notesWritten) > 0 {
					fmt.Fprintf(w, "Would tell your agents about this board: %s\n", strings.Join(notesWritten, ", "))
				}
				if setup.roles {
					printRoleRuns(w, setup.roleRuns, s.Root, true)
				}
				fmt.Fprint(w, hintLines(s.Root, setup.roles, setup.lanes))
				fmt.Fprintln(w, "nothing was written (--dry-run)")
				return nil
			}

			if line := announceLine(notesWritten); line != "" {
				fmt.Fprint(w, line)
			}
			if setup.roles {
				printRoleRuns(w, setup.roleRuns, s.Root, false)
			}
			fmt.Fprint(w, hintLines(s.Root, setup.roles, setup.lanes))
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
	cmd.Flags().StringVar(&agentFile, "agent-file", "", "where the jaira section lives: agents (AGENTS.md), claude (CLAUDE.md) or both")
	cmd.Flags().BoolVar(&setup.roles, "roles", false, "also bring the jaira skills up to date wherever they are installed: this repository's skills directory and your personal one")
	cmd.Flags().BoolVar(&setup.lanes, "lanes", false, "also replace this board's lane files that are an older version of a lane jaira ships")
	cmd.Flags().BoolVar(&setup.dry, "dry-run", false, "say what would be written and what is left alone, and write nothing")
	return cmd
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
