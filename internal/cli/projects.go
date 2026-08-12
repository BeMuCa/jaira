package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/berk/jaira/core/project"
)

func newProjectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: "List boards you have opened",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ps := project.Load()
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"projects": ps, "count": len(ps)})
			}
			if len(ps) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No boards recorded yet. Open one with 'jaira' inside a repository.")
				return nil
			}
			for _, p := range ps {
				fmt.Fprintf(cmd.OutOrStdout(), "%-24s %s\n", p.Name, p.Root)
			}
			return nil
		},
	}
}
