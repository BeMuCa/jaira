package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// updateInProgress suppresses nudgeIfStale while 'jaira update' itself is
// running: it is the fix the nudge points at, so it must not nag about
// itself.
var updateInProgress bool

// nudgeIfStale keeps a board in step with a newer binary on the one path
// every command takes — the same hook bindDriverIfShared
// (internal/cli/share.go) uses for "this clone needs one-time setup".
//
// Two things, with a stamp each. The first command after an upgrade brings
// the board's unedited lanes and its jaira section up to date
// (refreshStaleBoard) and says so in one line. And as long as 'jaira update'
// has not run under this version, a line points at it: it alone shows the
// release notes since, including the lines the reader must act on, so the
// refresh does not take that stamp. The command that refreshed folds the
// pointer into its own line rather than printing two.
//
// A dev build's version says nothing about what is actually in it, so it
// never acts: nagging a developer about their own build on every invocation
// would be pure noise. This writes to os.Stderr directly rather than through
// cmd's streams, matching bindDriverIfShared, because it must reach the
// terminal regardless of --json — stdout is reserved for the payload an agent
// parses.
//
// Whether a newer *release* is published is a separate question this
// function deliberately does not answer: the TUI's pill (internal/tui) shows
// it, and 'jaira next', 'jaira resume', 'jaira claim' and 'jaira move' print
// it (warnNewerRelease). Not every command, because a line on each of them is
// noise for the scripts that drive the CLI.
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
	setUp := release.Stamped(dir) == release.Current
	if !g.dryRun && release.LanesRefreshed(dir) != release.Current {
		line, err := refreshStaleBoard(s)
		if err != nil {
			// What was written before the failure is said all the same: the
			// next run finds those lanes current and would never name them.
			if line != "" {
				fmt.Fprintln(os.Stderr, line)
			}
			fmt.Fprintf(os.Stderr, "jaira: could not bring this board's lanes up to %s (%v) — 'jaira update --lanes' does it by hand\n", release.Current, err)
			return
		}
		if line != "" {
			if !setUp {
				line += " — 'jaira update' shows what else changed"
			}
			fmt.Fprintln(os.Stderr, line)
			return
		}
	}
	if !setUp {
		fmt.Fprintf(os.Stderr, "jaira: this board was set up by an older version of jaira — run 'jaira update' to bring it up to date\n")
	}
}

// refreshStaleBoard does to a board an older jaira left behind what
// 'jaira update --lanes' does to its files: every lane file that is an older
// version jaira shipped is replaced, and the jaira section is rewritten where
// a file already carries it. Lanes stay in the repository, so the whole team
// runs one flow; what this writes are ordinary changes in the working tree
// for a person to commit.
//
// It returns the one line to print, "" when nothing on disk changed. A lane
// somebody edited is never replaced, only named. An error leaves the
// lanes-refreshed stamp unset, so the next command tries again, and still
// comes with the line for what was written before it.
//
// The lock makes two processes started at once do the work once: the second
// waits, finds the stamp the first wrote, and returns "". The stamp is what
// keeps every later command at one small file read.
func refreshStaleBoard(s *ticket.Store) (string, error) {
	unlock, err := s.Lock("board-refresh")
	if err != nil {
		return "", err
	}
	defer unlock()
	if release.LanesRefreshed(s.StateDir()) == release.Current {
		return "", nil
	}
	results, err := lane.RefreshShipped(s.Root, false)
	if err != nil {
		return refreshedLine(results, nil, nil), err
	}
	set, err := lane.Load(s.Root)
	if err != nil {
		return refreshedLine(results, nil, nil), err
	}
	written, newer, err := board.RenewInAgentFiles(s.Root, set.Facts())
	if err == nil {
		err = release.MarkLanesRefreshed(s.StateDir())
	}
	return refreshedLine(results, written, newer), err
}

// refreshedLine says in one line what refreshStaleBoard changed, "" for nothing.
func refreshedLine(results []lane.Refreshed, written, newer []string) string {
	var updated, skipped []string
	for _, r := range results {
		name := r.ID
		if r.Board != "" {
			name = r.Board + "/" + r.ID
		}
		switch r.Action {
		case lane.RefreshUpdated:
			updated = append(updated, name)
		case lane.RefreshSkipped:
			skipped = append(skipped, name)
		}
	}
	var parts []string
	if len(updated) > 0 {
		parts = append(parts, "updated lanes "+strings.Join(updated, ", "))
	}
	if len(written) > 0 {
		parts = append(parts, "rewrote the jaira section in "+strings.Join(written, ", "))
	}
	if len(newer) > 0 {
		parts = append(parts, "left the jaira section alone, a newer jaira wrote it: "+strings.Join(newer, ", "))
	}
	if len(skipped) > 0 {
		// An older jaira cannot tell an edited lane from one a newer jaira
		// shipped: both are no version it knows.
		parts = append(parts, "left alone, not a version this jaira shipped (edited, or from a newer jaira): "+
			strings.Join(skipped, ", ")+" — 'jaira update --lanes' names them")
	}
	if len(parts) == 0 {
		return ""
	}
	return "jaira: brought this board up to " + release.Current + ": " + strings.Join(parts, "; ")
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

It never changes whether this board is shared: 'jaira init' makes a new board
private, 'jaira share' publishes it, and re-adding the ignore entry to a shared
board would silently stop new tickets reaching teammates. A private board
whose entry an older jaira wrote into the tracked .gitignore has it moved to
.git/info/exclude, which no worktree shows as a change; commit the .gitignore
that leaves behind once. A private board excluded nowhere — a teammate whose
.gitignore lost the line in that commit — is excluded there too. A private
board in a directory whose name starts with a space loses the entry an older
jaira wrote for it without that space, which excluded nothing; an entry a
person wrote, or one that names a directory holding a board, stays. A board
with .jaira/.gitattributes or committed files under .jaira/ is never touched.

--roles brings the roles, the cockpit mod and the jaira skill up to date in every
skills directory that already holds the roles: this repository's and your
personal one. Outside a board it does the personal one only. --lanes replaces
every lane file of this repository's boards that is an older version of a lane
jaira ships. The first command after an upgrade does the same by itself, and
rewrites the jaira section where it stands, unless it was given --dry-run.
A file you edited — one that is no version jaira ever shipped — is never
replaced: it is named, and the command exits 3. Without --roles or --lanes only
the jaira section is written, as before, and it says when roles or lanes are
behind.

--dry-run says what would be written and what would be left alone, and writes
nothing — not even the version stamp. With --lanes the preview of the jaira
section is computed from the lanes as --lanes would leave them, so it names
the same files the real run writes.

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
			// A dry run wrote no lane, so it previews from the lanes as the
			// refresh would leave them, not as they still are on disk.
			var boardLanes *lane.Set
			if setup.dry && setup.lanes {
				boardLanes, _ = lane.PreviewRefresh(s.Root, setup.laneResults)
			} else {
				boardLanes, _ = lane.Load(s.Root)
			}
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
			ignore, err := moveIgnore(s, setup.dry)
			if err != nil {
				return fail(ExitError, "prepare_failed", "could not move the board's ignore entry: %v", err)
			}
			// A shared board's milestone files stay out of git; a private
			// board's are excluded with the rest of .jaira/.
			var msIgnored bool
			if isShared(s) {
				if msIgnored, err = board.AddMilestonesIgnore(s.Root, setup.dry); err != nil {
					return fail(ExitError, "prepare_failed", "could not write .gitignore: %v", err)
				}
			}
			msTracked := milestonesTracked(s.Root)
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
					"agent_notes":  notesWritten,
					"ignore_moved": ignore.moved, "ignore_added": ignore.added, "ignore_unknown": ignore.unknown,
					"ignore_stale_dropped": ignore.stale,
					"milestones_ignored":   msIgnored, "milestones_tracked": msTracked,
					"notes": entriesJSON(notes),
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
				fmt.Fprint(w, ignore.line(s.Root, true))
				fmt.Fprint(w, milestonesLine(msIgnored, msTracked, true))
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
			fmt.Fprint(w, ignore.line(s.Root, false))
			fmt.Fprint(w, milestonesLine(msIgnored, msTracked, false))
			fmt.Fprint(w, hintLines(s.Root, setup.roles, setup.lanes))
			if len(notes) == 0 {
				fmt.Fprintf(w, "\nNothing has changed since the version that last set this board up.\n")
				return nil
			}
			fmt.Fprintf(w, "\nWhat's changed since then:\n")
			writeNotes(w, notes)
			return nil
		},
	}
	cmd.Flags().StringVar(&agentFile, "agent-file", "", "where the jaira section lives: agents (AGENTS.md), claude (CLAUDE.md) or both")
	cmd.Flags().BoolVar(&setup.roles, "roles", false, "also bring the jaira skills up to date wherever they are installed: this repository's skills directory and your personal one")
	cmd.Flags().BoolVar(&setup.lanes, "lanes", false, "also replace this board's lane files that are an older version of a lane jaira ships")
	cmd.Flags().BoolVar(&setup.dry, "dry-run", false, "say what would be written and what is left alone, and write nothing")
	return cmd
}

// ignoreMove is what moveIgnore did, or would do on a dry run.
type ignoreMove struct {
	moved, added bool
	unknown      string // why git could not say whether the board is ignored; nothing was touched
	stale        string // the entry without the leading space an older jaira wrote, taken out
}

// moveIgnore keeps a private board private without the tracked .gitignore:
// an entry an older jaira wrote there moves to .git/info/exclude, and a
// private board excluded nowhere — a teammate who pulled the commit that took
// the line out of .gitignore — is excluded there. Private means no
// .jaira/.gitattributes and nothing under .jaira/ committed; a shared board is
// left exactly as it is, as is every board outside a git repository.
func moveIgnore(s *ticket.Store, dry bool) (ignoreMove, error) {
	var m ignoreMove
	repo := &gitrepo.Repo{Dir: s.Root}
	// One test for both branches: a board in git keeps its .gitignore, since
	// the commit that took the line out would un-ignore it for every
	// teammate who pulls it.
	if !repo.IsRepo() || hasAttributes(s) || repo.Tracks(ticket.DirName) {
		return m, nil
	}
	// git decides whether the board is ignored, and nothing else does. A git
	// that cannot answer — a symlinked .jaira — leaves everything as it is.
	by, err := board.GitIgnores(s.Root)
	if err != nil {
		m.unknown = err.Error()
		return m, nil
	}
	if err := m.exclude(s.Root, by, dry); err != nil {
		return m, err
	}
	// The entry an older jaira wrote for a board in " sub" lacks the space
	// and excludes nothing; it goes after the right one is in place.
	if dry {
		_, m.stale = board.StaleEntry(s.Root)
		return m, nil
	}
	m.stale, err = board.DropStaleEntry(s.Root)
	return m, err
}

// exclude moves or writes the board's entry; by is the rule git says
// ignores the board, "" for none.
func (m *ignoreMove) exclude(root, by string, dry bool) (err error) {
	if board.InGitignore(root) {
		// Where git names no info/exclude — before 2.31 it does not know
		// --path-format — .gitignore is the place, and there is nothing to move.
		if board.IgnoreFileName(root) == ".gitignore" {
			return nil
		}
		if dry {
			m.moved = true
			return nil
		}
		m.moved, err = board.MoveIgnore(root)
		return err
	}
	if by != "" {
		return nil
	}
	// Not ignored. AddIgnore writes only where jaira's entry is missing: a
	// "!/.jaira/" a person wrote above it is theirs, and a second entry
	// would not change it.
	if dry {
		m.added = !board.Ignored(root)
		return nil
	}
	m.added, err = board.AddIgnore(root)
	return err
}

// line says what moveIgnore did, "" for nothing.
func (m ignoreMove) line(root string, dry bool) string {
	out := m.entryLine(root, dry)
	switch {
	case m.stale != "" && dry:
		out += fmt.Sprintf("Would remove %s from %s: an older jaira wrote it for this board without the leading space, and it excludes nothing.\n", m.stale, board.IgnoreFileName(root))
	case m.stale != "":
		out += fmt.Sprintf("Removed %s from %s: an older jaira wrote it for this board without the leading space, and it excluded nothing.\n", m.stale, board.IgnoreFileName(root))
	}
	return out
}

func (m ignoreMove) entryLine(root string, dry bool) string {
	file := board.IgnoreFileName(root)
	switch {
	case m.unknown != "":
		return fmt.Sprintf("Left this private board's ignore entry alone: git could not say whether .jaira/ is ignored (%s).\n", m.unknown)
	case m.moved && dry:
		return fmt.Sprintf("Would move this private board's entry from .gitignore to %s.\n", file)
	case m.moved:
		return fmt.Sprintf("Moved this private board's entry from .gitignore to %s — commit the changed .gitignore once.\n", file)
	case m.added && dry:
		return fmt.Sprintf("Would exclude this private board in %s: nothing excludes it yet.\n", file)
	case m.added:
		return fmt.Sprintf("Excluded this private board in %s: nothing excluded it.\n", file)
	}
	return ""
}

// writeNotes prints the release notes. A change marked important ("- ! text"
// in its note file) reads "  ! text" instead of "  - text", and a feature
// ("- + text") reads "  + text": plain text, because agents and pipes read this
// output too.
func writeNotes(w io.Writer, notes []release.Entry) {
	for _, e := range notes {
		fmt.Fprintf(w, "\n%s\n", e.Version)
		for _, c := range e.Changes {
			rest, feature := release.Feature(c)
			if text, ok := release.Important(rest); ok {
				fmt.Fprintf(w, "  ! %s\n", text)
				continue
			}
			if feature {
				fmt.Fprintf(w, "  + %s\n", rest)
				continue
			}
			fmt.Fprintf(w, "  - %s\n", c)
		}
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

// milestonesLine is what update says about the milestone files: the
// .gitignore line it added, and the warning about committed ones.
func milestonesLine(added, tracked, dry bool) string {
	out := ""
	switch {
	case added && dry:
		out += fmt.Sprintf("Would add %s to .gitignore: milestone files travel on their refs and stay out of git.\n", board.MilestonesIgnoreLine)
	case added:
		out += fmt.Sprintf("Added %s to .gitignore: milestone files travel on their refs and stay out of git.\n", board.MilestonesIgnoreLine)
	}
	if tracked {
		out += "warning: " + trackedMilestonesMessage + "\n"
	}
	return out
}
