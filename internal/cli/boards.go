package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	coreidentity "github.com/BeMuCa/jaira/core/identity"
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
be built from. 'jaira boards publish' shares a board with the team,
'jaira boards shared' lists what teammates shared, and 'jaira boards adopt'
takes one into your catalogue.`,
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
	cmd.AddCommand(newBoardsCreateCmd(), newBoardsTypesCmd(),
		newBoardsPublishCmd(), newBoardsSharedCmd(), newBoardsAdoptCmd())
	return cmd
}

func newBoardsCreateCmd() *cobra.Command {
	var typ, description, lanesFlag string
	cmd := &cobra.Command{
		Use:   "create <name> --description <which tickets belong here> [--type <type> | --lanes <a,b,...>]",
		Short: "Add a board to this repository, built from a board type or from lanes you pick",
		Long: `Writes .jaira/boards/<name>/lanes/ — the lanes of the board type, their order
file and the board's description. Without --type the board gets the built-in
lanes for code; 'jaira boards types' lists the others, the ones jaira ships
and the boards in your catalogue.

--lanes picks the lanes instead of a type, in column order: any lane of
today's board, of your catalogue or built into jaira ('jaira lanes add'
lists the last two).

--description is required: it says which tickets belong on this board, and it
is what 'jaira boards' shows a teamlead or an agent deciding where a new ticket
goes. A board that already has lanes is never overwritten (exit 3).

The new board is also saved into your catalogue, ~/.jaira/boards/<name>/, so
every repository on this machine can build it with --type <name>. A
catalogue board of that name already there is left as it is, and a board
named like a type jaira ships is not saved.

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
			if lanesFlag != "" && cmd.Flags().Changed("type") {
				return fail(ExitUsage, "usage", "--type and --lanes each say which lanes the board gets; give one")
			}
			s, err := openStore()
			if err != nil {
				return err
			}
			var dir string
			var ids []string
			if lanesFlag != "" {
				typ = ""
				set, lerr := lane.Load(s.Root)
				if lerr != nil {
					return lerr
				}
				ids = splitCSV(lanesFlag)
				dir, err = lane.CreateBoardWithLanes(s.Root, name, set, ids, description)
				if err != nil && !errors.Is(err, lane.ErrBoardExists) {
					return fail(ExitUsage, "unknown_lane", "%v", err)
				}
			} else {
				dir, ids, err = lane.CreateBoard(s.Root, name, typ, description)
			}
			switch {
			case errors.Is(err, lane.ErrBoardExists):
				return fail(ExitValidation, "board_exists", "%v", err)
			case errors.Is(err, lane.ErrNoSuchBoardType):
				return fail(ExitUsage, "unknown_type", "%v", err)
			case err != nil:
				return err
			}
			catalogue, saved, err := lane.SaveToCatalogue(s.Root, name)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{
					"name": name, "type": typ, "lanes": ids,
					"description": strings.TrimSpace(description), "lanes_dir": dir,
					"catalogue": catalogue, "catalogue_saved": saved,
				})
			}
			rel, _ := filepath.Rel(s.Root, dir)
			from := "from type " + typ
			if typ == "" {
				from = "from the lanes you picked"
			}
			fmt.Fprintf(w, "Board %s created %s in %s: %s\n", name, from, filepath.ToSlash(rel), strings.Join(ids, " → "))
			switch {
			case saved:
				fmt.Fprintf(w, "Saved to your catalogue in %s: every repository here can build it with 'jaira boards create <name> --type %s'.\n", catalogue, name)
			case catalogue != "":
				fmt.Fprintf(w, "Your catalogue already has a board %s, in %s; it was left as it is.\n", name, catalogue)
			}
			fmt.Fprintf(w, "File tickets on it with 'jaira create --board %s'.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", lane.DefaultBoardType, "the board type to build from; 'jaira boards types' lists them")
	cmd.Flags().StringVar(&lanesFlag, "lanes", "", "comma-separated lane ids, in column order, instead of a type")
	cmd.Flags().StringVar(&description, "description", "", "which tickets belong on this board (required)")
	return cmd
}

func newBoardsPublishCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "publish <name>",
		Short: "Copy a board into .jaira/shared/<you>/boards/ for teammates to adopt",
		Long: `Copies the named board — its lane files, order and description — into
.jaira/shared/<your identity>/boards/<name>/, the board counterpart of
'jaira lanes publish'. "default" publishes today's board. Commit the folder
for teammates to see it; nothing reaches anyone's board until they run
'jaira boards adopt' and create a board from it.

Refuses to overwrite an earlier publish unless --force is given.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			who := coreidentity.Slug(identity())
			dst, err := lane.PublishBoard(s.Root, args[0], s.SharedDir(), who, force)
			switch {
			case errors.Is(err, lane.ErrBoardExists):
				return fail(ExitValidation, "exists", "%v — use --force to overwrite", err)
			case errors.Is(err, lane.ErrNoSuchBoard):
				return fail(ExitUsage, "no_such_board", "%v", err)
			case err != nil:
				return fail(ExitUsage, "usage", "%v", err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"name": filepath.Base(dst), "path": dst, "who": who})
			}
			fmt.Fprintf(w, "published %s\n", dst)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an earlier publish")
	return cmd
}

func newBoardsSharedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shared",
		Short: "List the boards teammates published under .jaira/shared/",
		Long: `Lists every board under .jaira/shared/<who>/boards/, with who published it,
its lanes and its PATH — the argument 'jaira boards adopt' takes. A folder
that does not read as a board is named and skipped. Writes nothing.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			boards, warnings, err := lane.SharedBoards(s.Root)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(boards))
				for _, b := range boards {
					arr = append(arr, map[string]any{
						"name": b.Name, "folder": b.Folder, "description": b.Description,
						"lanes": b.Lanes, "path": b.Path,
					})
				}
				return emit(w, map[string]any{"boards": arr, "count": len(arr), "warnings": warnings})
			}
			if len(boards) == 0 {
				fmt.Fprintln(w, "No shared boards yet — 'jaira boards publish <name>' shares one.")
			}
			for _, b := range boards {
				fmt.Fprintf(w, "  %-16s %-12s %s\n      lanes: %s\n      %s\n", b.Name, b.Folder, b.Description, strings.Join(b.Lanes, " → "), b.Path)
			}
			for _, warn := range warnings {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", warn)
			}
			return nil
		},
	}
}

func newBoardsAdoptCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "adopt <path>",
		Short: "Copy a teammate's shared board into your catalogue",
		Long: `Copies a published board — the PATH 'jaira boards shared' prints — into your
catalogue, ~/.jaira/boards/<name>/. From there every repository on this
machine builds it with 'jaira boards create <new> --type <name>'.

Adopting means agreeing to run whatever its lanes' prompts say at whatever
model tier they declare — read the lane files in the folder first.

Refuses a board named like a type jaira ships, and an existing catalogue
board unless --force is given.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bt, dst, err := lane.AdoptBoard(args[0], force)
			switch {
			case errors.Is(err, lane.ErrBoardExists):
				return fail(ExitValidation, "exists", "%v — use --force to overwrite", err)
			case err != nil:
				return fail(ExitUsage, "usage", "%v", err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"name": bt.Name, "lanes": bt.Lanes, "path": dst})
			}
			fmt.Fprintf(w, "adopted %s; build it with 'jaira boards create <name> --type %s'\n", dst, bt.Name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing catalogue board")
	return cmd
}

func newBoardsTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List the board types 'jaira boards create --type' can build from",
		Long: `Lists the board types jaira ships and the boards in your catalogue,
~/.jaira/boards/ — the boards 'jaira boards create' saved there and the
ones 'jaira boards adopt' took in. The SOURCE column says which. A
catalogue board named like a shipped type is not listed: the shipped one is
what --type builds.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			types, err := lane.BoardTypes()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(types))
				for _, t := range types {
					arr = append(arr, map[string]any{"name": t.Name, "description": t.Description, "lanes": t.Lanes, "source": typeSource(t)})
				}
				return emit(w, map[string]any{"types": arr, "count": len(arr)})
			}
			for _, t := range types {
				fmt.Fprintf(w, "  %-16s %-10s %s\n", t.Name, typeSource(t), t.Description)
			}
			return nil
		},
	}
}

// typeSource says where a board type comes from.
func typeSource(t lane.BoardType) string {
	if t.Catalogue {
		return "catalogue"
	}
	return "shipped"
}
