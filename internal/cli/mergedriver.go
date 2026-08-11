package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/merge"
	"github.com/berk/jaira/core/ticket"
)

// mergeDriverName is the git merge driver identifier. It appears in both
// .gitattributes (committed, travels with a clone) and .git/config (local, does
// not), which is why registration has to happen on the machine.
const mergeDriverName = "jaira"

func newMergeDriverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "merge-driver <base> <ours> <theirs>",
		Short:  "Internal: three-way merge for ticket files, invoked by git",
		Hidden: true,
		Args:   cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath, oursPath, theirsPath := args[0], args[1], args[2]

			read := func(p string) ([]byte, error) {
				b, err := os.ReadFile(p)
				if err != nil {
					return nil, fmt.Errorf("merge-driver: %w", err)
				}
				return b, nil
			}
			base, err := read(basePath)
			if err != nil {
				return err
			}
			ours, err := read(oursPath)
			if err != nil {
				return err
			}
			theirs, err := read(theirsPath)
			if err != nil {
				return err
			}

			lanes, err := lane.Load()
			if err != nil {
				return err
			}
			res, err := merge.Merge(base, ours, theirs, lanes)
			if err != nil {
				// A merge driver that errors makes git fall back to reporting a
				// conflict, which is the safe outcome — but say why on stderr.
				fmt.Fprintf(os.Stderr, "jaira merge-driver: %v\n", err)
				return &codedError{code: 1, reason: "merge_failed", message: err.Error()}
			}

			// Git takes the merged result from the "ours" path.
			if err := os.WriteFile(oursPath, res.Merged, 0o644); err != nil {
				return err
			}
			for _, n := range res.Notes {
				fmt.Fprintf(os.Stderr, "jaira: merged %s\n", n)
			}
			if !res.Clean() {
				fmt.Fprintf(os.Stderr, "jaira: %s\n", merge.RenderConflicts(res.Conflicts))
				// A non-zero exit tells git the merge is unresolved.
				return &codedError{
					code:    1,
					reason:  "merge_conflict",
					message: fmt.Sprintf("%d field(s) need manual resolution", len(res.Conflicts)),
				}
			}
			return nil
		},
	}
	return cmd
}

// ensureMergeDriver registers the driver for this clone if it is not already
// registered, and reports whether it did.
//
// Registration cannot travel with the repository: git deliberately does not let a
// clone configure an executable that the cloning user has not opted into. The
// pattern in .gitattributes is committed, but pointing that pattern at a program
// is a local decision — so the tool makes it on first use and says so out loud
// rather than modifying git configuration silently.
func ensureMergeDriver(root string) (installed bool, err error) {
	key := "merge." + mergeDriverName + ".driver"
	out, err := exec.Command("git", "-C", root, "config", "--local", "--get", key).Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return false, nil
	}

	self, err := os.Executable()
	if err != nil || self == "" {
		self = "jaira"
	}
	driver := fmt.Sprintf("%q merge-driver %%O %%A %%B", self)
	if err := exec.Command("git", "-C", root, "config", "--local", key, driver).Run(); err != nil {
		return false, fmt.Errorf("could not register the git merge driver: %w", err)
	}
	if err := exec.Command("git", "-C", root, "config", "--local",
		"merge."+mergeDriverName+".name", "jaira field-aware ticket merge").Run(); err != nil {
		return true, nil // cosmetic only
	}
	return true, nil
}

// writeGitAttributes points the committed attributes file at the driver, so a
// teammate's clone knows which files want structural merging.
func writeGitAttributes(s *ticket.Store) (changed bool, err error) {
	path := filepath.Join(s.Root, ticket.DirName, ".gitattributes")
	want := fmt.Sprintf("%s/*.md merge=%s\n", ticket.TicketsSubdir, mergeDriverName)
	existing, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(existing), "merge="+mergeDriverName) {
		return false, nil
	}
	header := "# Ticket files are merged field by field by the jaira merge driver.\n" +
		"# Without it, two people moving the same ticket collide on the status line.\n"
	body := header + want
	if err == nil {
		body = string(existing)
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body += want
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func newResolveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve <id>",
		Short: "Show the conflicting fields of a ticket left in conflict",
		Long: `Reports which fields git could not merge, side by side.

Most concurrent edits never reach this command: lanes, dependencies and commit
lists are merged structurally. What lands here is the case no rule can settle —
two people rewriting the same prose.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			raw, err := os.ReadFile(t.Path)
			if err != nil {
				return err
			}
			if !strings.Contains(string(raw), "<<<<<<<") {
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"conflicted": false, "id": t.ID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s has no conflict markers.\n", ticket.Handle(t.ID))
				return nil
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"conflicted": true, "id": t.ID, "path": t.Path,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(),
				"%s still contains conflict markers.\nEdit %s, then commit.\n",
				ticket.Handle(t.ID), t.Path)
			return nil
		},
	}
}
