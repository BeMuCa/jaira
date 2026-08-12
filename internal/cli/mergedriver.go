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
		Args:   exactArgs(3),
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
	var takeOurs, takeTheirs bool
	cmd := &cobra.Command{
		Use:   "resolve <id>",
		Short: "Settle the fields a merge could not resolve",
		Long: `Shows the fields git could not merge, and can settle them.

Most concurrent edits never reach this command: lanes, dependencies and commit
lists are merged structurally. What lands here is the case no rule can settle —
two people rewriting the same prose.

The ticket stays valid YAML while a conflict is outstanding, so the board keeps
working. The losing value is parked in a conflict-theirs-<field> key rather than
being discarded.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			fields, _ := t.Doc().List(merge.FieldConflicts)
			if len(fields) == 0 {
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"conflicted": false, "id": t.ID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s has no outstanding conflicts.\n", ticket.Handle(t.ID))
				return nil
			}

			if takeOurs || takeTheirs {
				_, err := s.Mutate(t.ID, func(t *ticket.Ticket) error {
					for _, f := range fields {
						key := "conflict-theirs-" + f
						theirs, ok, _ := t.Doc().Scalar(key)
						if takeTheirs && ok {
							if err := t.Doc().SetScalar(f, theirs); err != nil {
								return err
							}
						}
						// Either way the parked value is no longer needed.
						if ok {
							if err := t.Doc().SetScalar(key, ""); err != nil {
								return err
							}
						}
					}
					return t.Doc().SetList(merge.FieldConflicts, nil)
				})
				if err != nil {
					return err
				}
				side := "ours"
				if takeTheirs {
					side = "theirs"
				}
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{"resolved": true, "took": side, "fields": fields})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Resolved %d field(s) in favour of %s.\n", len(fields), side)
				return nil
			}

			type row struct {
				Field  string `json:"field"`
				Ours   string `json:"ours"`
				Theirs string `json:"theirs"`
			}
			var rows []row
			for _, f := range fields {
				ours := fieldValue(t, f)
				theirs, _, _ := t.Doc().Scalar("conflict-theirs-" + f)
				rows = append(rows, row{Field: f, Ours: ours, Theirs: theirs})
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"conflicted": true, "id": t.ID, "fields": rows,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s has %d unresolved field(s):\n\n", ticket.Handle(t.ID), len(rows))
			for _, r := range rows {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n    ours:   %s\n    theirs: %s\n\n", r.Field, r.Ours, r.Theirs)
			}
			fmt.Fprintf(cmd.OutOrStdout(),
				"Edit %s directly, or run this command with --take-ours or --take-theirs.\n", t.Path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&takeOurs, "take-ours", false, "keep this side's values")
	cmd.Flags().BoolVar(&takeTheirs, "take-theirs", false, "adopt the incoming values")
	return cmd
}
