package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// warnNewerRelease writes one line when the release cache knows a release
// newer than this binary, so an agent that never sees the TUI's pill asks the
// user whether to update, and runs the upgrade only on a yes. It is the
// PostRun of 'jaira next' and 'jaira resume', the commands a session starts
// with, and of 'jaira claim' and 'jaira move', which a session already running
// calls at every ticket and every lane — so a release that becomes known
// mid-session reaches it at its next step. Always on the command's stderr:
// stdout is the payload an agent parses.
//
// PostRun, not a call at the top of RunE: cobra runs PostRun only when RunE
// returned nil, so the line follows a command that worked and never stands in
// front of an error — under --json stderr is the structured error a caller
// parses, and the cockpit shows a failed move's stderr as its message.
//
// The line names both commands a yes runs: 'jaira self upgrade' brings the
// binary and, from PWAF22 on, the skills ('update --roles'); the lane files
// stay a step of their own, 'jaira update --roles --lanes' in the repository.
// Before it asks, the agent runs their previews, 'jaira self upgrade --check'
// and 'jaira update --roles --lanes --dry-run', and shows the user what they
// report (EDSQZP): a yes is given to a change the user has seen, as in the
// cockpit's Update pane.
//
// A dev build says nothing, and so does a cache that was never written — the
// same rules as the TUI's pill. Reading the cache also refreshes it in the
// background when it is more than an hour old (selfupdate.PollCache), which is
// the only thing that checks for a release on a machine that never opens the
// TUI.
func warnNewerRelease(cmd *cobra.Command, _ []string) {
	if release.Current == "dev" {
		return
	}
	latest, known := selfupdate.PollCache()
	if !known || !release.Newer(latest, release.Current) {
		return
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "jaira: jaira %s is released (this is %s) - run 'jaira self upgrade --check' and 'jaira update --roles --lanes --dry-run', show the user what they report and ask the user whether to update; only on yes run 'jaira self upgrade', then 'jaira update --roles --lanes'\n", latest, release.Current)
}
