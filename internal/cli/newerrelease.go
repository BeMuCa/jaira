package cli

import (
	"fmt"
	"io"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// warnNewerRelease writes one line when the release cache knows a release
// newer than this binary, so an agent that never sees the TUI's pill still
// tells the user. It is called by 'jaira next' and 'jaira resume' only, the
// commands a session starts with, and with the command's stderr: stdout is the
// payload an agent parses.
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
	fmt.Fprintf(w, "jaira: jaira %s is released (this is %s) - tell the user; 'jaira self upgrade' installs it\n", latest, release.Current)
}
