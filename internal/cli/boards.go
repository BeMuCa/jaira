package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
)

func newBoardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boards",
		Short: "List this repository's boards and what each is for",
		Long: `Lists every board of this repository — today's, shown as "default", and each
one under .jaira/boards/ — with how many open tickets it holds and its
description.

Read this before creating a ticket in a repository with more than one board,
and file the ticket on the board whose description its subject matches
('jaira create --board <name>'). A description is one plain-text file,
'description', in the board's lanes directory beside its order file, written
by 'jaira boards create' or by hand — a board set up from the built-in lanes
starts with one saying it is for code; a board without one says where it
goes. Writes nothing.

'jaira boards create' adds a board, 'jaira boards types' lists what it can
be built from.`,
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
	cmd.AddCommand(newBoardsCreateCmd(), newBoardsTypesCmd())
	return cmd
}

func newBoardsCreateCmd() *cobra.Command {
	var typ, description string
	cmd := &cobra.Command{
		Use:   "create <name> --description <which tickets belong here> [--type <type>]",
		Short: "Add a board to this repository, built from a board type",
		Long: `Writes .jaira/boards/<name>/lanes/ — the lanes of the board type, their order
file and the board's description — and nothing else. Without --type the board
gets the built-in lanes for code; 'jaira boards types' lists the others.

--description is required: it says which tickets belong on this board, and it
is what 'jaira boards' shows a teamlead or an agent deciding where a new ticket
goes. A board that already has lanes is never overwritten (exit 3).

File tickets on it with 'jaira create --board <name>'.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(description) == "" {
				return fail(ExitUsage, "missing_description",
					"boards create needs --description: one line saying which tickets belong on this board — 'jaira boards' shows it to whoever files a ticket, so a teamlead can tell this board from the others")
			}
			name, err := lane.NormalizeBoard(args[0])
			if err != nil {
				return fail(ExitUsage, "usage", "%v", err)
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			dir, ids, err := lane.CreateBoard(s.Root, name, typ, description)
			switch {
			case errors.Is(err, lane.ErrBoardExists):
				return fail(ExitValidation, "board_exists", "%v", err)
			case errors.Is(err, lane.ErrNoSuchBoardType):
				return fail(ExitUsage, "unknown_type", "%v", err)
			case err != nil:
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{
					"name": name, "type": typ, "lanes": ids,
					"description": strings.TrimSpace(description), "lanes_dir": dir,
				})
			}
			rel, _ := filepath.Rel(s.Root, dir)
			fmt.Fprintf(w, "Board %s created from type %s in %s: %s\n", name, typ, filepath.ToSlash(rel), strings.Join(ids, " → "))
			fmt.Fprintf(w, "File tickets on it with 'jaira create --board %s'.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", lane.DefaultBoardType, "the board type to build from; 'jaira boards types' lists them")
	cmd.Flags().StringVar(&description, "description", "", "which tickets belong on this board (required)")
	return cmd
}

func newBoardsTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List the board types 'jaira boards create --type' can build from",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			types, err := lane.BoardTypes()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(types))
				for _, t := range types {
					arr = append(arr, map[string]any{"name": t.Name, "description": t.Description, "lanes": t.Lanes})
				}
				return emit(w, map[string]any{"types": arr, "count": len(arr)})
			}
			for _, t := range types {
				fmt.Fprintf(w, "  %-16s %s\n", t.Name, t.Description)
			}
			return nil
		},
	}
}
