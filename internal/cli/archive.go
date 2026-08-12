package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/spf13/cobra"
)

func newArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive [id]",
		Short: "Take a ticket off the board, or list what has been taken off",
		Long: `Moves a ticket into .jaira/archive/. With no argument, lists the archive.

The file is moved, not deleted. The board exists so you can still answer what a
task was for months later, and on a private board a deleted file is not in git
history to recover from — so nothing here removes anything. 'jaira restore' moves
a ticket back.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fail(ExitUsage, "usage", "archive takes at most one ticket id, received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()

			if len(args) == 0 {
				names, err := archivedNames(s)
				if err != nil {
					return err
				}
				if g.jsonOut {
					return emit(w, map[string]any{"archived": names, "count": len(names)})
				}
				if len(names) == 0 {
					fmt.Fprintf(w, "The archive is empty.\n")
					return nil
				}
				for _, n := range names {
					fmt.Fprintf(w, "%s\n", n)
				}
				fmt.Fprintf(w, "\n%d archived. Bring one back with 'jaira restore <file>'.\n", len(names))
				return nil
			}

			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			dst, err := s.Archive(args[0])
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(w, map[string]any{
					"archived": true, "id": t.ID, "handle": ticket.Handle(t.ID),
					"path": dst, "file": filepath.Base(dst),
				})
			}
			fmt.Fprintf(w, "Archived %s  %s\n", ticket.Handle(t.ID), t.Title)
			fmt.Fprintf(w, "Moved to %s — restore it with 'jaira restore %s'.\n",
				filepath.Join(ticket.DirName, ticket.ArchiveSubdir), filepath.Base(dst))
			return nil
		},
	}
	return cmd
}

func newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <file>",
		Short: "Put an archived ticket back on the board",
		Args:  exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			dst, err := s.Restore(args[0])
			if err != nil {
				return &codedError{code: ExitNotFound, reason: "not_archived", message: err.Error()}
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"restored": true, "path": dst})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restored %s to the board.\n", filepath.Base(dst))
			return nil
		},
	}
}

func archivedNames(s *ticket.Store) ([]string, error) {
	entries, err := os.ReadDir(s.ArchiveDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
