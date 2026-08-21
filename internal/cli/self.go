package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// newSelfCmd groups operations on the jaira binary itself, as distinct from
// 'jaira update', which re-applies a board's setup and downloads nothing —
// the name "update" was already taken by that command, so this is "self".
func newSelfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self",
		Short: "Manage the jaira binary itself",
		Args:  noArgs(),
	}
	cmd.AddCommand(newSelfUpgradeCmd())
	return cmd
}

// newSelfUpgradeCmd wraps core/selfupdate's resolve-download-verify-install
// path — a Go port of scripts/install.sh's own logic, so an existing install
// can move to a new release without re-running the curl pipeline.
func newSelfUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Replace the running jaira binary with the latest release",
		Long: `Downloads, verifies and installs the latest published jaira release in place
of the binary you are running.

This is not 'jaira update', which re-applies this repository's board setup
and downloads nothing — this replaces the jaira executable itself.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			exe, err := os.Executable()
			if err != nil {
				return err
			}
			// A symlinked install must be resolved before anything decides
			// what to replace, or the upgrade writes over the link and
			// orphans the real file it pointed at.
			target, err := filepath.EvalSymlinks(exe)
			if err != nil {
				return err
			}
			selfupdate.Sweep(filepath.Dir(target))

			ctx := context.Background()
			client := selfupdate.New()
			rel, err := client.Latest(ctx)
			if err != nil {
				return err
			}
			bin, err := client.Binary(ctx, rel, runtime.GOOS, runtime.GOARCH)
			if err != nil {
				return err
			}
			if err := selfupdate.Replace(target, bin); err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{
					"current": release.Current, "latest": rel.Version,
					"up_to_date": false, "upgraded": true, "target": target,
				})
			}
			fmt.Fprintf(w, "upgraded jaira %s -> %s (%s)\n", release.Current, rel.Version, target)
			fmt.Fprintf(w, "run 'jaira update' in each repository to bring its board setup in step.\n")
			return nil
		},
	}
	return cmd
}
