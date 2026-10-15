package tui

import (
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/release"
)

// styImportant carries a release-note line marked "- ! text": the warning
// colour, bold, so it stands out from the plain changes around it.
var styImportant = lipgloss.NewStyle().Foreground(colWarn).Bold(true)

// notesView is a scrollable list of release notes, one block per version. The
// launcher's "w" opens it with every release; the popup after a version change
// opens it with the releases since the last one the user read.
type notesView struct {
	title   string
	entries []release.Entry
	scroll  int
	// page is how far pgup/pgdown move: the body height of the last render.
	page int
}

func newNotesView(title string, entries []release.Entry) *notesView {
	return &notesView{title: title, entries: entries, page: 10}
}

// key handles one key and reports whether the view asked to be closed.
// Scrolling uses the same keys the help screen does.
func (v *notesView) key(k string) (done bool) {
	switch k {
	case "esc", "q", "enter":
		return true
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

// lines lays the notes out for width. Each line is styled on its own, so a
// wrapped change never leaves a padded block behind a prefix.
func (v *notesView) lines(width int) []string {
	var out []string
	for i, e := range v.entries {
		if i > 0 {
			out = append(out, "")
		}
		out = append(out, styLaneTitle.Render(e.Version))
		for _, c := range e.Changes {
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
	}
	return out
}

// render draws the title, the visible slice of the notes and the key hint,
// at most height lines. The scroll offset is clamped here, where the length is
// known.
func (v *notesView) render(width, height int) string {
	body := max(1, height-4)
	v.page = max(1, body-1)
	all := v.lines(width)
	v.scroll = min(v.scroll, max(0, len(all)-body))
	end := min(len(all), v.scroll+body)

	var b strings.Builder
	b.WriteString(styLaneTitle.Render(truncate(v.title, width)) + "\n\n")
	b.WriteString(strings.Join(all[v.scroll:end], "\n"))
	b.WriteString("\n\n" + styMeta.Render("j/k scroll · esc close"))
	return b.String()
}

// seenFile records, per working tree and next to the "version" stamp, the last
// version whose release notes the user closed. It lives in the state directory,
// never in git: "I read them" is one person's fact, not the team's.
const seenFile = "notes-seen"

// readSeen reads <dir>/notes-seen, trimmed; a missing file is "".
func readSeen(dir string) string {
	b, err := os.ReadFile(filepath.Join(dir, seenFile))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// markSeen records the running version as read. A dev build has no version to
// record, so it writes nothing. Best effort: failing to remember means the
// popup shows once more, which is not worth an error on a board.
func markSeen(dir string) {
	if release.Current == "dev" || release.Current == "" {
		return
	}
	if os.MkdirAll(dir, 0o755) != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, seenFile), []byte(release.Current+"\n"), 0o644)
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
	base := notesBase(release.Current, readSeen(dir), release.Stamped(dir))
	if base == "" {
		markSeen(dir)
		return
	}
	entries := release.Since(base)
	if len(entries) == 0 {
		markSeen(dir)
		return
	}
	m.notes = newNotesView("What is new since "+base, entries)
	m.returnTo = modeBoard
	m.mode = modeReleaseNotes
}
