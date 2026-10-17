package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
)

func newBoardsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "boards",
		Short: "List this repository's boards and what each is for",
		Long: `Lists every board of this repository — today's, shown as "default", and each
one under .jaira/boards/ — with how many open tickets it holds and its
description.

Read this before creating a ticket in a repository with more than one board,
and file the ticket on the board whose description its subject matches
('jaira create --board <name>'). A description is one plain-text file,
'description', in the board's lanes directory beside its order file, written
by hand; a board without one says where it goes. Writes nothing.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			_, all, err := loadEnv(s)
			if err != nil {
				return err
			}
			type boardRow struct {
				name, description, file string
				open                    int
			}
			names := append([]string{lane.DefaultBoardName}, lane.Boards(s.Root)...)
			rows := make([]boardRow, 0, len(names))
			for _, name := range names {
				lanes, err := boardLanes(s, name)
				if err != nil {
					return err
				}
				// Open as in 'jaira tags': a ticket in its board's terminal
				// lane is finished work, not what the board is busy with.
				open := 0
				for _, t := range onBoard(all, name) {
					if l, ok := lanes.Get(t.Status); !ok || !l.Terminal {
						open++
					}
				}
				rows = append(rows, boardRow{
					name: name, open: open,
					description: lane.Description(s.Root, name),
					file:        lane.DescriptionPath(s.Root, name),
				})
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(rows))
				for _, r := range rows {
					arr = append(arr, map[string]any{
						"name": r.name, "open": r.open,
						"description": r.description, "description_file": r.file,
					})
				}
				return emit(w, map[string]any{"boards": arr, "count": len(arr)})
			}
			for _, r := range rows {
				desc := r.description
				if desc == "" {
					rel, _ := filepath.Rel(s.Root, r.file)
					desc = "no description yet — write one into " + filepath.ToSlash(rel)
				}
				fmt.Fprintf(w, "  %-16s %3d open   %s\n", r.name, r.open, desc)
			}
			return nil
		},
	}
}
