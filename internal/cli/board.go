package cli

import (
	"github.com/spf13/cobra"

	tea "charm.land/bubbletea/v2"

	"github.com/berk/jaira/core/project"
	"github.com/berk/jaira/core/ticket"
	"github.com/berk/jaira/internal/tui"
)

// newBoardCmd launches the interactive board. It is also what running jaira with
// no subcommand does, since looking at the board is the common case.
func newBoardCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "board",
		Aliases: []string{"ui"},
		Short:   "Open the interactive board",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			dir := g.dir
			for {
				s, err := ticket.Discover(dir)
				if err != nil {
					return err
				}
				project.Remember(s.Root)

				m, err := tui.New(s)
				if err != nil {
					return err
				}
				final, err := tea.NewProgram(m).Run()
				if err != nil {
					return err
				}
				// Switching boards reopens the program rather than swapping the
				// store underneath it, so no per-board state survives the move.
				if fm, ok := final.(*tui.Model); ok && fm.SwitchTo != "" {
					dir = fm.SwitchTo
					continue
				}
				return nil
			}
		},
	}
}
