package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

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

// newLanesUseCmd wraps lane.Export, the same call the lane settings screen's
// 'u' key makes — one implementation of "take a lane from the catalogue into
// this project", not two.
func newLanesUseCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "use <id>",
		Short: "Copy a catalogue lane into this project's own lane directory",
		Long: `Copies the named lane's file, verbatim, into <root>/.jaira/lanes — the same
move the lane settings screen's 'u' key performs. Once a project has its own
lane file, that directory becomes authoritative for the board (see 'jaira
lanes path').

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
			dst, err := lane.Export(l, lane.ProjectLanesDir(s.Root), force)
			if err != nil {
				return writeConflictError(err)
			}
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
