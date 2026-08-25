package tui

import (
	"fmt"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// versionLine renders the persistent "which version am I, and is a newer one
// published" indicator shown in the launcher's and the board's footer.
//
// It is meant to be computed once, at construction (see Home.versionLine and
// Model.versionLine), rather than on every render: selfupdate.PollCache also
// decides whether to touch the cache and spawn a detached refresh when it is
// stale (see core/selfupdate/cache.go's own doc comments for why), and a
// bubbletea program can render dozens of times per keypress — nothing here
// needs to run more than once per session to stay correct within the ~24h
// staleness window PollCache itself enforces.
//
// This is the TUI's own read of the cache, independent of any CLI command:
// CLI commands stay silent about a published release entirely (see
// internal/cli/update.go's nudgeIfStale), because a line on every 'jaira
// list' or 'jaira next' is noise for the scripts and agents that drive the
// CLI. The person actually watching a screen is the one this is for.
//
// A cache that has never been written to at all — never checked, or the
// check is disabled via JAIRA_NO_UPDATE_CHECK — reports the running version
// alone rather than claiming "up to date": that would assert a fact nobody
// has actually verified.
func versionLine() string {
	latest, known := selfupdate.PollCache()
	switch {
	case !known:
		return styMeta.Render(fmt.Sprintf("jaira %s", release.Current))
	case latest == release.Current:
		return styMeta.Render(fmt.Sprintf("jaira %s · up to date", release.Current))
	default:
		return styMeta.Render(fmt.Sprintf("jaira %s · %s available — run: jaira self upgrade", release.Current, latest))
	}
}
