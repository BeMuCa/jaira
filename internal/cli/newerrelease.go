package cli

import (
	"fmt"
	"io"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// warnNewerRelease writes one line when the release cache knows a release
// newer than this binary, so an agent that never sees the TUI's pill asks the
// user whether to update, and runs the upgrade only on a yes. It is called by
// 'jaira next' and 'jaira resume', the commands a session starts with, and by
// 'jaira claim' and 'jaira move', which a session already running calls at
// every ticket and every lane — so a release that becomes known mid-session
// reaches it at its next step. Always with the command's stderr: stdout is
// the payload an agent parses.
//
// The line names both commands a yes runs: 'jaira self upgrade' brings the
// binary and, from PWAF22 on, the skills ('update --roles'); the lane files
// stay a step of their own, 'jaira update --roles --lanes' in the repository.
//
// A dev build says nothing, and so does a cache that was never written — the
// same rules as the TUI's pill. Reading the cache also refreshes it in the
// background when it is more than an hour old (selfupdate.PollCache), which is
// the only thing that checks for a release on a machine that never opens the
// TUI.
func warnNewerRelease(w io.Writer) {
	if release.Current == "dev" {
		return
	}
	latest, known := selfupdate.PollCache()
	if !known || !release.Newer(latest, release.Current) {
		return
	}
	fmt.Fprintf(w, "jaira: jaira %s is released (this is %s) - ask the user whether to update; on yes run 'jaira self upgrade', then 'jaira update --roles --lanes'\n", latest, release.Current)
}
