package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/release"
)

// styImportant carries a release-note line marked "- ! text": the warning
// colour, bold, so it stands out from the plain changes around it.
var styImportant = lipgloss.NewStyle().Foreground(colWarn).Bold(true)

// notesView shows release notes one version per page, the newest first:
// left/h pages to the older version, right/l back to the newer one, and j/k
// scroll a version longer than the screen. The launcher's "w" opens it with
// every release; the popup after a version change opens it with the releases
// since the last one the user read.
type notesView struct {
	title   string
	entries []release.Entry
	// at is the version on screen, an index into entries: 0 is the newest.
	at     int
	scroll int
	// page is how far pgup/pgdown move: the body height of the last render.
	page int
}

func newNotesView(title string, entries []release.Entry) *notesView {
	return &notesView{title: title, entries: entries, page: 10}
}

// key handles one key and reports whether the view asked to be closed.
// Scrolling uses the same keys the help screen does; paging past either end
// does nothing.
func (v *notesView) key(k string) (done bool) {
	switch k {
	case "esc", "q", "enter":
		return true
	case "left", "h":
		if v.at < len(v.entries)-1 {
			v.at, v.scroll = v.at+1, 0
		}
	case "right", "l":
		if v.at > 0 {
			v.at, v.scroll = v.at-1, 0
		}
	case "down", "j":
		v.scroll++
	case "up", "k":
		v.scroll = max(0, v.scroll-1)
	case "pgdown", "ctrl+d":
		v.scroll += v.page
	case "pgup", "ctrl+u":
		v.scroll = max(0, v.scroll-v.page)
	case "g":
		v.scroll = 0
	}
	return false
}

// lines lays out the changes of the version on screen for width. Each line is
// styled on its own, so a wrapped change never leaves a padded block behind a
// prefix.
func (v *notesView) lines(width int) []string {
	if v.at >= len(v.entries) {
		return nil
	}
	var out []string
	for _, c := range v.entries[v.at].Changes {
		text, important := release.Important(c)
		bullet, sty := "- ", lipgloss.NewStyle()
		if important {
			bullet, sty = "! ", styImportant
		}
		for j, l := range strings.Split(wrap(text, max(1, width-2), 0), "\n") {
			lead := "  "
			if j == 0 {
				lead = bullet
			}
			out = append(out, sty.Render(lead+l))
		}
	}
	return out
}

// header names the version on screen and where it sits: "0.5.4 · 1/12".
func (v *notesView) header(width int) string {
	if v.at >= len(v.entries) {
		return ""
	}
	pos := fmt.Sprintf(" · %d/%d", v.at+1, len(v.entries))
	return styLaneTitle.Render(truncate(v.entries[v.at].Version, max(1, width-lipgloss.Width(pos)))) + styMeta.Render(pos)
}

// render draws the title, the version header, the visible slice of its
// changes and the key hint, at most height lines. The scroll offset is clamped
// here, where the length is known.
func (v *notesView) render(width, height int) string {
	body := max(1, height-5)
	v.page = max(1, body-1)
	all := v.lines(width)
	v.scroll = min(v.scroll, max(0, len(all)-body))
	end := min(len(all), v.scroll+body)

	var b strings.Builder
	b.WriteString(styLaneTitle.Render(truncate(v.title, width)) + "\n\n")
	b.WriteString(v.header(width) + "\n")
	b.WriteString(strings.Join(all[v.scroll:end], "\n"))
	b.WriteString("\n\n" + styMeta.Render(truncate("←/→ older/newer · j/k scroll · esc close", width)))
	return b.String()
}

// markSeen records the running version as read (release.MarkNotesSeen). A dev
// build has no version to record, so it writes nothing. Best effort: failing
// to remember means the popup shows once more, which is not worth an error on
// a board.
func markSeen(dir string) {
	if release.Current == "dev" || release.Current == "" {
		return
	}
	_ = release.MarkNotesSeen(dir)
}

// notesBase is the version whose notes the popup should show everything after,
// or "" for no popup. The seen record wins; without one the board's stamp
// stands in. A stamp of "" (a board never prepared) or "dev" says nothing about
// what the user last saw, and a base that is not older than current is a
// downgrade or no change.
func notesBase(current, seen, stamped string) string {
	if current == "dev" || seen == current {
		return ""
	}
	base := seen
	if base == "" {
		base = stamped
	}
	if base == "" || base == "dev" || !release.Newer(current, base) {
		return ""
	}
	return base
}

// openReleaseNotes shows the notes since the last version the user read, once,
// when the board opens. Called from New only: a board reached with 1-9 shows it
// the next time it is opened.
func (m *Model) openReleaseNotes() {
	dir := m.store.StateDir()
	seen := release.NotesSeen(dir)
	base := notesBase(release.Current, seen, release.Stamped(dir))
	if base == "" {
		if seen != release.Current {
			markSeen(dir)
		}
		return
	}
	entries := release.Since(base)
	if len(entries) == 0 {
		markSeen(dir)
		return
	}
	m.notes = newNotesView(fmt.Sprintf("What is new since %s · %d unread", base, len(entries)), entries)
	m.returnTo = modeBoard
	m.mode = modeReleaseNotes
}
