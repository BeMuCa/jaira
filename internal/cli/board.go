package cli

import (
	"github.com/spf13/cobra"

	tea "charm.land/bubbletea/v2"

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
			s, err := openStore()
			if err != nil {
				return err
			}
			m, err := tui.New(s)
			if err != nil {
				return err
			}
			_, err = tea.NewProgram(m).Run()
			return err
		},
	}
}
