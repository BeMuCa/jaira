package cli

import (
	"github.com/spf13/cobra"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/internal/tui"
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
			s, err := ticket.Discover(g.dir)
			if err != nil {
				return err
			}
			project.Remember(s.Root)

			m, err := tui.New(s)
			if err != nil {
				return err
			}
			// Board switches happen inside the program (Model.switchBoard):
			// restarting the process per switch dropped the alternate screen
			// for a frame, flashing the terminal through on every switch.
			_, err = tea.NewProgram(m).Run()
			return err
		},
	}
}
