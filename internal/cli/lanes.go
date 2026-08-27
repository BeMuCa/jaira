package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/lane"
)

// writeConflictError classes an error from Export/Publish/Adopt as a gate
// refusal (exit 3) when it is the "already exists" refusal those functions
// already return, naming --force as the way past it. Any other error (a
// permission problem, a bad source file) is passed through unchanged, so it
// still exits 1 as an unexpected failure rather than being mislabelled as a
// refusal.
func writeConflictError(err error) error {
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return fail(ExitValidation, "exists", "%v — use --force to overwrite", err)
	}
	return err
}

// newLanesUseCmd copies a lane's catalogue or shipped version onto this board,
// verbatim. Without --force it is 'lanes add' for a lane the board does not
// have; with --force it overwrites the board's copy — the way a board's lane
// is reset to the shipped one, or brought up to date after an upgrade.
func newLanesUseCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "use <id>",
		Short: "Copy a lane's catalogue or shipped version onto this board, --force to overwrite the board's copy",
		Long: `Copies the named lane's file, verbatim, from the catalogue — or the binary,
for a shipped lane — into <root>/.jaira/lanes, and puts it on the board if it
was not. A board is its lane directory, so the copy is the lane from then on.

Refuses to overwrite a file already on the board unless --force is given. With
--force, this is how a board's copy of a lane is reset to the shipped one, or
brought up to date after a jaira upgrade changed it.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			// Load the board first: a board opened for the first time gets its
			// lane files here, so 'use' of a lane it already has refuses like
			// it should, instead of writing into an empty directory.
			if _, err := lane.Load(s.Root); err != nil {
				return err
			}
			// Resolve against the offer, not the board: the board's copy is
			// what --force is meant to replace, so resolving through it would
			// copy the stale file onto itself and call it a refresh.
			lanes, err := lane.Load("")
			if err != nil {
				return err
			}
			l, ok := lanes.Get(args[0])
			if !ok {
				return fail(ExitUsage, "no_such_lane", "no lane %q is in the catalogue; available: %s",
					args[0], strings.Join(lanes.IDs(), ", "))
			}
			dst, err := lane.Export(l, lane.ProjectLanesDir(s.Root), force)
			if err != nil {
				return writeConflictError(err)
			}
			// A lane new to the board takes the last column, like 'lanes add'.
			if ids, err := lane.LoadOrder(s.Root); err == nil && len(ids) > 0 && !slices.Contains(ids, l.ID) {
				if err := lane.SaveOrder(s.Root, append(ids, l.ID)); err != nil {
					return err
				}
			}
			refreshAgentNote(cmd, s.Root)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": l.ID, "path": dst})
			}
			fmt.Fprintf(w, "wrote %s\n", dst)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}

// newLanesAddCmd wraps lane.Add, the same call the settings screen's '+'
// column makes once a lane is chosen from its catalogue — one implementation
// of "add a lane to this project", not two.
func newLanesAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add a built-in or catalogue lane to this project's board",
		Long: `Adds the named lane to this board, appending it to the column order — the
same move the lane settings screen's '+' column makes once a lane is chosen.
A board is its lane directory: the lane's file is written there, and that is
what puts it on the board.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			set, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			dst, err := lane.Add(s.Root, set, args[0])
			if err != nil {
				if strings.Contains(err.Error(), "already part of this project") ||
					strings.Contains(err.Error(), "no lane") {
					return fail(ExitUsage, "no_such_lane", "%v", err)
				}
				return err
			}
			refreshAgentNote(cmd, s.Root)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": args[0], "path": dst})
			}
			fmt.Fprintf(w, "added %s to this project (%s)\n", args[0], dst)
			return nil
		},
	}
	return cmd
}

// newLanesRemoveCmd wraps lane.Remove, the same call the settings screen's
// 'x' key makes.
func newLanesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a lane from this project (it stays in the catalogue)",
		Long: `Removes the named lane from this project's board and its column order —
the same move the lane settings screen's 'x' key makes. The lane itself is
untouched in the catalogue; only this project stops using it.

Refused, naming the tickets, when any ticket currently sits in the lane: a
lane that vanishes under a ticket would leave it in a lane nothing knows.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			set, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			path, err := lane.Remove(s.Root, set, s, args[0])
			if err != nil {
				if strings.Contains(err.Error(), "cannot be removed") {
					return fail(ExitValidation, "lane_occupied", "%v", err)
				}
				if strings.Contains(err.Error(), "is part of this project") {
					return fail(ExitUsage, "no_such_lane", "%v", err)
				}
				return err
			}
			refreshAgentNote(cmd, s.Root)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": args[0], "path": path})
			}
			fmt.Fprintf(w, "removed %s from this project\n", args[0])
			return nil
		},
	}
	return cmd
}

// newLanesMoveCmd wraps lane.MoveLane, the same call the settings screen's
// H/L keys make.
func newLanesMoveCmd() *cobra.Command {
	var left, right bool
	cmd := &cobra.Command{
		Use:   "move <id> --left|--right",
		Short: "Move a lane one column left or right in this project",
		Long: `Shifts the named lane one step in this project's column order, swapping it
with its neighbour — the same move the lane settings screen's H/L keys make.
Moving past either end is a no-op, not an error.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if left == right {
				return fail(ExitUsage, "usage", "exactly one of --left or --right is required")
			}
			delta := 1
			if left {
				delta = -1
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			set, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			if err := lane.MoveLane(s.Root, set, args[0], delta); err != nil {
				if strings.Contains(err.Error(), "is part of this project") {
					return fail(ExitUsage, "no_such_lane", "%v", err)
				}
				return err
			}
			refreshAgentNote(cmd, s.Root)
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": args[0], "delta": delta})
			}
			fmt.Fprintf(w, "moved %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVar(&left, "left", false, "move one column left")
	cmd.Flags().BoolVar(&right, "right", false, "move one column right")
	return cmd
}

// newLanesPublishCmd wraps lane.Publish, the same call the lane settings
// screen's 'p' key makes.
func newLanesPublishCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "publish <id>",
		Short: "Copy a lane into .jaira/shared/<you>/ for teammates to adopt",
		Long: `Copies the named lane's file into .jaira/shared/<your identity>/, stamping
creator: with your name when the lane does not already declare one — the same
move the lane settings screen's 'p' key performs. Nothing is loaded onto
anyone's board until a teammate runs 'jaira lanes adopt'.

Refuses to overwrite an existing file unless --force is given.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			lanes, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			l, ok := lanes.Get(args[0])
			if !ok {
				return fail(ExitUsage, "no_such_lane", "no lane %q is installed; available: %s",
					args[0], strings.Join(lanes.IDs(), ", "))
			}
			// The slug, not the raw identity, is both the directory name and the
			// stamped creator: value — matching lane.Publish's other caller, the
			// lane settings screen's 'p' key.
			who := coreidentity.Slug(identity())
			dstDir := filepath.Join(s.SharedDir(), who)
			dst, err := lane.Publish(l, dstDir, who, force)
			if err != nil {
				return writeConflictError(err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": l.ID, "path": dst, "who": who})
			}
			fmt.Fprintf(w, "published %s\n", dst)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}

// newLanesAdoptCmd wraps lane.Adopt, the same call the lane settings screen's
// 'a' key makes.
func newLanesAdoptCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "adopt <path>",
		Short: "Copy a teammate's shared lane into your catalogue",
		Long: `Copies a lane file into your catalogue (see 'jaira lanes path'), under the
lane's own id — the same move the lane settings screen's 'a' key performs.

Takes the PATH column 'jaira lanes shared' prints, not a lane id: two
teammates can publish a lane under the same id, and only the path names one
file unambiguously.

Adopting means agreeing to run whatever the file's prompt says at whatever
model tier it declares — read it first, with 'jaira lanes show' once it is
installed or by opening the file directly before that.

Refuses to overwrite an existing catalogue entry unless --force is given.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, dst, err := lane.Adopt(args[0], lane.UserLanesDir(), force)
			if err != nil {
				return writeConflictError(err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": l.ID, "path": dst})
			}
			fmt.Fprintf(w, "adopted %s\n", dst)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing catalogue entry")
	return cmd
}

// splitCSV splits a comma-separated flag value, trimming whitespace and
// dropping empty entries — the same treatment tickets.go's list-field
// assignments already give a comma-separated value.
func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// newLanesDefaultCmd wraps LoadDefaultBoard/SaveDefaultBoard, the same calls
// the launcher's 'd' screen makes, so a default board can be read and
// written without opening the TUI.
func newLanesDefaultCmd() *cobra.Command {
	var lanesFlag, optionsFlag string
	var clear bool
	cmd := &cobra.Command{
		Use:   "default",
		Short: "Show or set which lanes and options a new board starts with",
		Long: `With no flags, prints which lanes and which ticket Options a freshly
initialised board gets, and the file's path.

--lanes and --options set the selection (a comma-separated list of ids or
option names); --clear removes the file, returning to the built-ins. An
unknown lane id or option name is refused rather than written, since a
default board naming something nobody has installed is how a board ends up
looking empty.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			w := cmd.OutOrStdout()

			if clear {
				path := lane.DefaultBoardPath()
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return err
				}
				if g.jsonOut {
					return emit(w, map[string]any{"cleared": true, "path": path})
				}
				fmt.Fprintf(w, "cleared the default board (%s)\n", path)
				return nil
			}

			board, err := lane.LoadDefaultBoard()
			if err != nil {
				return err
			}

			lanesSet, optionsSet := cmd.Flags().Changed("lanes"), cmd.Flags().Changed("options")
			if lanesSet || optionsSet {
				set, err := lane.Load(bestEffortRoot())
				if err != nil {
					return err
				}
				if lanesSet {
					ids := splitCSV(lanesFlag)
					for _, id := range ids {
						if _, ok := set.Get(id); !ok {
							return fail(ExitUsage, "no_such_lane", "no lane %q is installed; available: %s",
								id, strings.Join(set.IDs(), ", "))
						}
					}
					board.Lanes = ids
				}
				if optionsSet {
					opts := splitCSV(optionsFlag)
					known := make(map[string]bool, len(set.Options()))
					for _, o := range set.Options() {
						known[o] = true
					}
					for _, o := range opts {
						if !known[o] {
							return fail(ExitUsage, "no_such_option", "no installed lane requires option %q; available: %s",
								o, strings.Join(set.Options(), ", "))
						}
					}
					board.Options = opts
				}
				if err := lane.SaveDefaultBoard(board); err != nil {
					return err
				}
			}

			usingBuiltins := len(board.Lanes) == 0
			effectiveLanes := board.Lanes
			if usingBuiltins {
				builtins, err := lane.Builtins()
				if err != nil {
					return err
				}
				effectiveLanes = make([]string, 0, len(builtins))
				for _, l := range builtins {
					effectiveLanes = append(effectiveLanes, l.ID)
				}
			}

			if g.jsonOut {
				return emit(w, map[string]any{
					"path": board.Path, "lanes": effectiveLanes, "options": board.Options,
					"using_builtins": usingBuiltins,
				})
			}
			fmt.Fprintf(w, "Default board: %s\n", board.Path)
			if usingBuiltins {
				fmt.Fprintf(w, "Lanes (built-in): %s\n", strings.Join(effectiveLanes, ", "))
			} else {
				fmt.Fprintf(w, "Lanes: %s\n", strings.Join(effectiveLanes, ", "))
			}
			if len(board.Options) == 0 {
				fmt.Fprintf(w, "Options: none pre-ticked\n")
			} else {
				fmt.Fprintf(w, "Options: %s\n", strings.Join(board.Options, ", "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&lanesFlag, "lanes", "", "comma-separated lane ids to set as the default board's lanes")
	cmd.Flags().StringVar(&optionsFlag, "options", "", "comma-separated ticket options to pre-tick on a new board")
	cmd.Flags().BoolVar(&clear, "clear", false, "remove the default board, returning to the built-ins")
	return cmd
}

// refreshAgentNote regenerates the jaira block in CLAUDE.md and AGENTS.md after
// the board's lanes changed.
//
// The block names this board's lanes, so adopting, adding, removing or
// reordering one is exactly the moment it stops being true. Leaving that to
// whoever remembers to run 'jaira update' means the note is wrong precisely
// when someone changed the pipeline — the case where an agent most needs it to
// be right.
//
// Best-effort and quiet: a lane command's own success does not depend on an
// agent file being writable, and a failure here is reported rather than turned
// into a failed lane operation.
func refreshAgentNote(cmd *cobra.Command, root string) {
	lanes, err := lane.Load(root)
	if err != nil {
		return
	}
	if _, err := board.AnnounceInAgentFiles(root, laneFacts(lanes)); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "jaira: warning: the lanes changed but the agent note could not be updated: %v\n", err)
	}
}
