package tui

import (
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
