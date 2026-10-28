package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
	"github.com/BeMuCa/jaira/core/ticket"
)

// selfExecutable is the binary 'self upgrade' replaces. Tests point it at a
// stand-in: replacing the running test binary fails on Windows the second
// time a process does it, because the first run's .old-<pid> is still there
// and still locked.
var selfExecutable = os.Executable

// runUpgraded runs the binary just installed at target with args and returns
// its exit code. The running process is the old binary — whatever it wrote
// from its own embedded files would be the old version — so the setup has to
// come from the new one. A variable so a test can watch the call instead of
// executing the stand-in bytes a fake release serves.
var runUpgraded = func(target string, args []string, stdout, stderr io.Writer) (int, error) {
	c := exec.Command(target, args...)
	c.Stdout, c.Stderr = stdout, stderr
	err := c.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), nil
	}
	if err != nil {
		return -1, err
	}
	return 0, nil
}

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

			exe, err := selfExecutable()
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
			selfupdate.Sweep(target)

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
				//
				// Not from a dev build, though: the cache exists to feed the
				// TUI's version line, and a dev build cannot be compared to a
				// published version at all. Writing it there would have the
				// version line report an answer to a question this build cannot ask.
				if !dev {
					_ = selfupdate.Write(selfupdate.Check{CheckedAt: now, Latest: rel.Version})
				}
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
				for _, o := range selfupdate.Overridden() {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: release host overridden by %s\n", o)
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

			// Exit 3 from the child means it skipped edited files and already
			// said which; only a child that did not run at all, or failed
			// otherwise, is worth a line here. The binary was replaced either
			// way, so the upgrade itself still succeeds.
			setupArgs := []string{"-C", g.dir, "update", "--roles"}
			// A pinned older release may not know 'update --roles' at all, and
			// its skills are older than the ones installed: leave the setup alone.
			downgrade := !release.Newer(rel.Version, release.Current)
			if g.jsonOut {
				p := payload(true)
				if downgrade {
					return emit(w, p)
				}
				var setupOut bytes.Buffer
				code, err := runUpgraded(target, append(setupArgs, "--json"), &setupOut, cmd.ErrOrStderr())
				warnSetupFailed(cmd.ErrOrStderr(), rel.Version, code, err)
				if json.Valid(setupOut.Bytes()) {
					p["setup"] = json.RawMessage(bytes.TrimSpace(setupOut.Bytes()))
				}
				p["setup_exit"] = code
				return emit(w, p)
			}
			fmt.Fprintf(w, "upgraded jaira %s -> %s (%s)\n", release.Current, rel.Version, target)
			// A binary was just replaced with bytes from somewhere other than
			// the default host. That is allowed, and it is not going to happen
			// quietly: an override lives in a shell profile or a CI job, where
			// nobody reads it, and unlike a flag it does not appear in the
			// command anyone ran.
			for _, o := range selfupdate.Overridden() {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: release host overridden by %s\n", o)
			}
			if downgrade {
				fmt.Fprintf(w, "%s is older than %s: your skills and boards are left as they are\n", rel.Version, release.Current)
				return nil
			}
			fmt.Fprintf(w, "bringing your setup in step with %s: jaira update --roles\n", rel.Version)
			code, err := runUpgraded(target, setupArgs, w, cmd.ErrOrStderr())
			warnSetupFailed(cmd.ErrOrStderr(), rel.Version, code, err)
			if _, err := ticket.Discover(g.dir); err == nil {
				fmt.Fprintf(w, "this repository is done; run 'jaira update' in each of your other repositories to bring its board setup in step.\n")
			} else {
				fmt.Fprintf(w, "run 'jaira update' in each repository to bring its board setup in step.\n")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&check, "check", false, "report what is available and install nothing")
	cmd.Flags().StringVar(&pinVersion, "version", "", "install this exact release, including an older one")
	return cmd
}

// warnSetupFailed says the new binary's setup did not complete: it did not
// start, or it exited with something other than 0. Exit 3 is not that — the
// child skipped edited files and already named them. It is a warning, not an
// error: the binary was replaced, so exiting non-zero would make an agent
// retry an upgrade that already happened.
func warnSetupFailed(w io.Writer, version string, code int, err error) {
	if err == nil && (code == 0 || code == ExitValidation) {
		return
	}
	if err == nil {
		err = fmt.Errorf("exit %d", code)
	}
	fmt.Fprintf(w, "jaira: setting up with %s failed (%v) — run 'jaira update --roles' yourself\n", version, err)
}
