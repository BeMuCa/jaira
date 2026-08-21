package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

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
	var check bool
	var pinVersion string
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Replace the running jaira binary with the latest release",
		Long: `Downloads, verifies and installs the latest published jaira release in place
of the binary you are running.

This is not 'jaira update', which re-applies this repository's board setup
and downloads nothing — this replaces the jaira executable itself.

Refuses to touch a Homebrew or 'go install' build, or an install directory
it cannot write to, naming the right way to upgrade that install instead.
--check reports what is available without installing anything; --version
pins (or downgrades to) an exact release.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if check && pinVersion != "" {
				// Pinning a version and refusing to install one is a
				// contradiction; silently honoring one over the other would
				// be worse than saying so.
				return fail(ExitUsage, "usage", "--check and --version cannot be used together")
			}

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

			// Every guard below runs before any network call, because a
			// refusal must be reachable on a machine with no network at
			// all.
			dev := release.Current == "dev"
			if dev && !check {
				return fail(ExitValidation, "dev_build",
					"a dev build has no published release to upgrade to — build a real release, or install one with `go install github.com/BeMuCa/jaira/cmd/jaira@latest`")
			}

			kind, instr := selfupdate.Detect(target)
			if kind != selfupdate.SelfManaged && !check {
				return fail(ExitValidation, string(kind), "%s", instr)
			}

			ctx := context.Background()
			client := selfupdate.New()

			var rel selfupdate.Release
			if pinVersion != "" {
				rel, err = client.At(ctx, pinVersion)
				if err != nil {
					return fail(ExitValidation, "no_release", "%v", err)
				}
			} else {
				rel, err = client.Latest(ctx)
				if err != nil {
					return err
				}
			}

			// A version already installed is a success that downloaded
			// nothing, matching this project's idempotent-retry convention:
			// re-invoking a command that has nothing left to do is not an
			// error, because agents retry.
			upToDate := rel.Version == release.Current
			now := time.Now().UTC()

			w := cmd.OutOrStdout()
			payload := func(upgraded bool) map[string]any {
				p := map[string]any{
					"current": release.Current, "latest": rel.Version,
					"up_to_date": upToDate, "upgraded": upgraded, "target": target,
					"install_method": string(kind), "checked_at": now.Format(time.RFC3339),
				}
				if dev {
					p["dev"] = true
				}
				return p
			}

			if check {
				// --check is the only path (besides a successful upgrade,
				// below) that ever populates the cache, which keeps the
				// "what did we last learn" logic in exactly one place. The
				// write error is ignored: a best-effort cache must never
				// turn a successful --check into a failure.
				_ = selfupdate.Write(selfupdate.Check{CheckedAt: now, Latest: rel.Version})
				if g.jsonOut {
					return emit(w, payload(false))
				}
				switch {
				case dev:
					fmt.Fprintf(w, "this is a dev build; the latest published release is %s\n", rel.Version)
				case upToDate:
					fmt.Fprintf(w, "jaira %s is already the latest release\n", release.Current)
				default:
					fmt.Fprintf(w, "jaira %s is available (you have %s)\n", rel.Version, release.Current)
				}
				return nil
			}

			if upToDate {
				if g.jsonOut {
					return emit(w, payload(false))
				}
				fmt.Fprintf(w, "jaira %s is already the latest release; nothing to do\n", release.Current)
				return nil
			}

			bin, err := client.Binary(ctx, rel, runtime.GOOS, runtime.GOARCH)
			if err != nil {
				return err
			}
			if err := selfupdate.Replace(target, bin); err != nil {
				return err
			}
			_ = selfupdate.Write(selfupdate.Check{CheckedAt: now, Latest: rel.Version})

			if g.jsonOut {
				return emit(w, payload(true))
			}
			fmt.Fprintf(w, "upgraded jaira %s -> %s (%s)\n", release.Current, rel.Version, target)
			fmt.Fprintf(w, "run 'jaira update' in each repository to bring its board setup in step.\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&check, "check", false, "report what is available and install nothing")
	cmd.Flags().StringVar(&pinVersion, "version", "", "install this exact release, including an older one")
	return cmd
}
