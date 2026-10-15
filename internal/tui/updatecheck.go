package tui

import (
	"fmt"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// versionLine renders the persistent "which version am I, and is a newer one
// published" indicator drawn at the top left of the board's project head and
// of the launcher: the version on a line of its own, and, when the cache
// knows a newer release, a pill "^ <version>" for the row under it. The pill
// is empty in every other case — same version, an older cache, a cache never
// written, a dev build — because only a known newer release is worth a mark.
//
// It used to sit in both screens' footers. It moved up because the question
// it answers is "which binary am I looking at" — asked while switching
// between 'jaira self upgrade', a 'go build' and ~/.local/bin — and a footer
// under a wrapping hint bar is not where anyone looks for an identity. Why
// it gets a row of its own up there is at the two call sites.
//
// It is meant to be computed once, at construction (see Home.versionLine and
// Model.versionLine), rather than on every render: selfupdate.PollCache also
// decides whether to touch the cache and spawn a detached refresh when it is
// stale (see core/selfupdate/cache.go's own doc comments for why), and a
// bubbletea program can render dozens of times per keypress — nothing here
// needs to run more than once per session to stay correct within the ~24h
// staleness window PollCache itself enforces.
//
// This is the TUI's own read of the cache. The CLI reads it too, but only in
// 'jaira next' and 'jaira resume' (see internal/cli/newerrelease.go), the
// commands an agent session starts with: a line on every 'jaira list' would
// be noise for the scripts that drive the CLI. The person actually watching
// a screen is the one this is for.
//
// A cache that has never been written to at all — never checked, or the
// check is disabled via JAIRA_NO_UPDATE_CHECK — reports the running version
// alone rather than claiming "up to date": that would assert a fact nobody
// has actually verified.
//
// A dev build names itself and stops there — "jaira dev", never an upgrade
// claim. This half-reverses DNAEPN, which had this function return nothing
// at all on a dev build: in a footer "jaira dev" was noise, but in the top
// left corner it is the whole point, because "which binary is this" is
// exactly the question a contributor switching between builds is asking.
// The other half of DNAEPN stands, below.
func versionLine() (line, pill string) {
	// A build that is not a published release has nothing to compare itself to.
	// release.Current is "dev" in every source build — which is what every
	// contributor runs — and comparing that string to a published version can
	// only ever say "different", so the line would advertise an upgrade to
	// code *older* than the code being run, pointing at a command that then
	// refuses with dev_build. So a dev build reports its identity and nothing
	// more, following the same rule as the !known case below: never assert a
	// fact nobody has checked. It returns before PollCache deliberately —
	// there is no answer worth having here, so a contributor's every run
	// should not stamp the cache or spawn a detached release check.
	if release.Current == "dev" {
		return styMeta.Render("jaira dev"), ""
	}
	latest, known := selfupdate.PollCache()
	switch {
	case !known:
		return styMeta.Render(fmt.Sprintf("jaira %s", release.Current)), ""
	case release.Newer(latest, release.Current):
		return styMeta.Render(fmt.Sprintf("jaira %s", release.Current)), styPill.Render("^ " + latest)
	default:
		return styMeta.Render(fmt.Sprintf("jaira %s · up to date", release.Current)), ""
	}
}
