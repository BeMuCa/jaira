package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/berk/jaira/core/ticket"
)

// ignoreLine is the entry that keeps a board private.
const ignoreLine = "/.jaira/"

// A board starts private. Committing tickets is what makes a board visible to
// everyone who can read the repository, and that should be a decision rather than
// a default — the tool has no way to know whether these notes are ready to be
// seen, or whether this repository is even yours to publish into.
//
// Private and shared differ only in whether `.jaira/` is gitignored. Nothing about
// the tickets themselves changes, so going public later is not a migration.

// isShared reports whether the board is committed rather than local-only.
func isShared(s *ticket.Store) bool {
	return !ignored(s.Root) && hasAttributes(s)
}

// ignored reports whether the root .gitignore excludes the board.
func ignored(root string) bool {
	b, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		switch strings.TrimSpace(line) {
		case ignoreLine, ".jaira/", ".jaira":
			return true
		}
	}
	return false
}

func hasAttributes(s *ticket.Store) bool {
	b, err := os.ReadFile(filepath.Join(s.Root, ticket.DirName, ".gitattributes"))
	return err == nil && strings.Contains(string(b), "merge="+mergeDriverName)
}

// addIgnore keeps the board out of git.
func addIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, _ := os.ReadFile(path)
	if ignored(root) {
		return false, nil
	}
	body := string(existing)
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	body += "\n# jaira board — private to this machine. Run 'jaira share' to publish it.\n" + ignoreLine + "\n"
	return true, os.WriteFile(path, []byte(body), 0o644)
}

// removeIgnore stops excluding the board.
func removeIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	var out []string
	var skipNextBlank bool
	for _, line := range strings.Split(string(existing), "\n") {
		t := strings.TrimSpace(line)
		if t == ignoreLine || t == ".jaira/" || t == ".jaira" {
			changed = true
			continue
		}
		if strings.Contains(t, "jaira board — private") {
			changed = true
			skipNextBlank = true
			continue
		}
		if skipNextBlank && t == "" {
			skipNextBlank = false
			continue
		}
		out = append(out, line)
	}
	if !changed {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}

func newShareCmd() *cobra.Command {
	var undo bool
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Publish the board so everyone with the repository can see it",
		Long: `Makes the board part of the repository.

A board starts private: 'jaira init' gitignores it, so your tickets are yours
alone until you decide otherwise. Publishing is a decision, not a default — the
tool cannot know whether these notes are ready to be read by whoever can clone
this repository.

Publishing stops ignoring .jaira/, writes the .gitattributes that tells git to
merge tickets field by field, and registers the merge driver for this clone. The
tickets themselves are unchanged, so this is a decision you can reverse with
--undo rather than a migration.

You still have to commit. jaira does not commit on your behalf.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()

			if undo {
				changed, err := addIgnore(s.Root)
				if err != nil {
					return err
				}
				if g.jsonOut {
					return emit(w, map[string]any{"shared": false, "changed": changed})
				}
				fmt.Fprintf(w, "The board is private again.\n")
				if changed {
					fmt.Fprintf(w, "Added %s to .gitignore.\n", ignoreLine)
				}
				fmt.Fprintf(w, "\nTickets already committed stay in history. To remove them:\n"+
					"  git rm -r --cached .jaira && git commit -m \"make jaira board private\"\n")
				return nil
			}

			unignored, err := removeIgnore(s.Root)
			if err != nil {
				return err
			}
			attrs, err := writeGitAttributes(s)
			if err != nil {
				return err
			}
			installed, driverErr := ensureMergeDriver(s.Root)

			if g.jsonOut {
				return emit(w, map[string]any{
					"shared": true, "unignored": unignored,
					"gitattributes_written": attrs, "merge_driver_installed": installed,
				})
			}
			fmt.Fprintf(w, "The board is now part of the repository.\n")
			if unignored {
				fmt.Fprintf(w, "Removed %s from .gitignore.\n", ignoreLine)
			}
			if attrs {
				fmt.Fprintf(w, "Wrote .jaira/.gitattributes so git merges tickets field by field.\n")
			}
			if driverErr != nil {
				fmt.Fprintf(os.Stderr, "jaira: warning: %v\n", driverErr)
			} else if installed {
				fmt.Fprintf(w, "Registered the merge driver in .git/config for this clone.\n")
			}
			count := 0
			if paths, err := s.Paths(); err == nil {
				count = len(paths)
			}
			fmt.Fprintf(w, "\nCommit to publish %d ticket(s):\n  git add .jaira .gitignore && git commit -m \"share jaira board\"\n", count)
			fmt.Fprintf(w, "\nTeammates then clone, and jaira binds the merge driver on their first command.\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&undo, "undo", false, "make the board private again")
	return cmd
}

// bindDriverIfShared registers the merge driver whenever the board is a committed,
// shared one and this clone has not bound it yet.
//
// This runs on every command rather than only at init, because the machine that
// needs the driver most is a teammate's: they clone a repository that already
// contains .gitattributes, never run init, and would otherwise get git's
// line-based merge with no indication anything was missing. The committed
// attributes file is the signal that the board is shared, so its presence is what
// triggers binding.
func bindDriverIfShared(s *ticket.Store) {
	if !hasAttributes(s) || ignored(s.Root) {
		return
	}
	installed, err := ensureMergeDriver(s.Root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "jaira: warning: %v\n", err)
		fmt.Fprintf(os.Stderr, "jaira: ticket merges will fall back to git's line-based merge\n")
		return
	}
	if installed {
		fmt.Fprintf(os.Stderr, "jaira: registered the field-aware merge driver for this clone\n")
	}
}
