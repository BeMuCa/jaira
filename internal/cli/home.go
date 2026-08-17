package cli

import (
	"fmt"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/internal/tui"
)

// runHome opens the launcher, then the board of whatever was chosen.
//
// Bare `jaira` lands here rather than on a board, including inside a repository
// that has one. That costs a keypress on the common path, which is a real cost
// for something glanced at constantly — `jaira board` remains the direct route
// and is documented as such.
func runHome(cmd *cobra.Command) error {
	for {
		// Anything under the working directory counts as a board worth offering
		// even if it has never been opened, so a fresh clone shows up without
		// having to be registered first.
		found := project.Discover(g.dir, project.MaxScanDepth)

		h, err := tui.NewHome(found)
		if err != nil {
			return err
		}
		final, err := tea.NewProgram(h).Run()
		if err != nil {
			return err
		}
		fh, ok := final.(*tui.Home)
		if !ok || fh.Quit || fh.Chosen == "" {
			return nil
		}
		if err := openBoardAt(fh.Chosen); err != nil {
			return err
		}
	}
}

// openBoardAt runs the board for one project. Board switches happen inside
// the program (Model.switchBoard), so there is nothing to loop over here.
func openBoardAt(dir string) error {
	s, err := ticket.Discover(dir)
	if err != nil {
		return err
	}
	project.Remember(s.Root)

	m, err := tui.New(s)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m).Run()
	return err
}

func newProjectsAddCmd() *cobra.Command {
	var scan bool
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Register a board so it appears on the home screen",
		Long: `Adds a directory to the list of boards jaira knows about.

With --scan, searches at most two levels below the path instead, registering
every board it finds. Two levels is deliberate: repositories normally sit one or
two directories under a code folder, and an unbounded walk of a home directory
descends into caches and mounted drives to find nothing.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			w := cmd.OutOrStdout()
			var roots []string
			if scan {
				roots = project.Discover(args[0], project.MaxScanDepth)
				if len(roots) == 0 {
					return fail(ExitNotFound, "no_boards",
						"no boards found within %d levels of %s", project.MaxScanDepth, args[0])
				}
			} else {
				abs, err := filepath.Abs(args[0])
				if err != nil {
					return err
				}
				s, err := ticket.Discover(abs)
				if err != nil {
					return fail(ExitNotFound, "no_board",
						"%s is not a board and has none above it; run 'jaira init' there first", args[0])
				}
				roots = []string{s.Root}
			}
			for _, r := range roots {
				project.Remember(r)
			}
			if g.jsonOut {
				return emit(w, map[string]any{"added": roots, "count": len(roots)})
			}
			for _, r := range roots {
				fmt.Fprintf(w, "Registered %s\n", r)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&scan, "scan", false, "search two levels below the path for boards")
	return cmd
}
