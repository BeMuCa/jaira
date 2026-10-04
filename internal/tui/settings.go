package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/settings"
)

// Actions settingsScreen.key can report: which entry was opened, or that the
// screen itself is finished.
const (
	settingsActionNone         = ""
	settingsActionBack         = "back"
	settingsActionLanes        = "lanes"
	settingsActionDefaultBoard = "default-board"
	// settingsActionLogbookDays opens the number in place rather than a
	// screen; settingsActionSaved reports that the screen wrote settings.json
	// and the board has to reload to show what it now says.
	settingsActionLogbookDays = "logbook-days"
	settingsActionSaved       = "saved"
)

// settingsEntry is one row in the menu: a name, what pressing enter on it
// opens, and a one-line description of what it is for.
type settingsEntry struct {
	name   string
	desc   string
	action string
}

var settingsEntries = []settingsEntry{
	{
		name:   "Lanes",
		desc:   "take a lane from the catalogue into this project, publish or adopt one",
		action: settingsActionLanes,
	},
	{
		name:   "Default board",
		desc:   "which lanes and options a new project starts with",
		action: settingsActionDefaultBoard,
	},
	{
		name:   "Logbook on the board",
		desc:   "days of the logbook the terminal lane shows below its tickets, and 'jaira logbook' lists; 0 hides them",
		action: settingsActionLogbookDays,
	},
}

// settingsScreen is the menu behind S: one door into both the lane screen and
// the default board, rather than each having its own binding. Following
// browse.go's shape: its own state, a key method, a render method, no
// bubbletea imports of its own.
type settingsScreen struct {
	idx int

	// days is logbook-days as settings.json has it now; editing is set while
	// the number is open for typing, buf holding what has been typed so far
	// and err what was wrong with the last attempt to save it.
	days    int
	editing bool
	buf     string
	err     string
}

func newSettingsScreen() *settingsScreen {
	return &settingsScreen{days: settings.Load().LogbookWindow()}
}

// editKey drives the number while it is open: digits type, backspace takes
// one back, enter saves and esc leaves it as it was. Saving reads the file
// again first, so nothing else in it — a remote, a hook — is written back from
// a copy older than the edit.
func (s *settingsScreen) editKey(k string) string {
	switch k {
	case "esc":
		s.editing, s.buf, s.err = false, "", ""
	case "backspace":
		if s.buf != "" {
			s.buf = s.buf[:len(s.buf)-1]
		}
	case "enter":
		n, err := strconv.Atoi(s.buf)
		if err != nil || n < 0 {
			s.err = "a number of days, 0 or more"
			return settingsActionNone
		}
		cur := settings.Load()
		cur.LogbookDays = &n
		if err := settings.Save(cur); err != nil {
			s.err = err.Error()
			return settingsActionNone
		}
		s.days, s.editing, s.buf, s.err = n, false, "", ""
		return settingsActionSaved
	default:
		if len(k) == 1 && k[0] >= '0' && k[0] <= '9' && len(s.buf) < 4 {
			s.buf += k
		}
	}
	return settingsActionNone
}

// key drives the screen. It reports settingsActionBack when esc/q finish the
// screen, the action of the highlighted entry when enter opens it, or
// settingsActionNone when the key only moved the cursor.
func (s *settingsScreen) key(k string) string {
	if s.editing {
		return s.editKey(k)
	}
	switch k {
	case "esc", "q":
		return settingsActionBack
	case "j", "down":
		if s.idx < len(settingsEntries)-1 {
			s.idx++
		}
	case "k", "up":
		if s.idx > 0 {
			s.idx--
		}
	case "enter":
		if a := settingsEntries[s.idx].action; a != settingsActionLogbookDays {
			return a
		}
		s.editing, s.buf, s.err = true, strconv.Itoa(s.days), ""
	}
	return settingsActionNone
}

func (s *settingsScreen) render(width, height int) string {
	w := max(20, width)
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Settings") + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	for i, e := range settingsEntries {
		lead := "  "
		name := e.name
		if i == s.idx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		if e.action == settingsActionLogbookDays {
			name += "  " + s.logbookValue()
		}
		sb.WriteString(truncate(lead+name, w) + "\n")
		sb.WriteString(styMeta.Render("      "+wrap(e.desc, max(10, w-6), 6)) + "\n")
		if e.action == settingsActionLogbookDays && s.err != "" {
			sb.WriteString(styErr.Render("      "+s.err) + "\n")
		}
	}

	hints := []string{"enter open", "esc back"}
	if s.editing {
		hints = []string{"0-9 days", "enter save", "esc cancel"}
	}
	for _, l := range wrapHints(hints, max(1, w)) {
		sb.WriteString("\n" + styMeta.Render(l))
	}
	return sb.String()
}

// logbookValue is the logbook-days row's value: the number being typed while
// it is open, otherwise what settings.json says, in words.
func (s *settingsScreen) logbookValue() string {
	switch {
	case s.editing:
		return stySelected.Render("[" + s.buf + "_]")
	case s.days == 0:
		return styMeta.Render("off")
	default:
		return styMeta.Render(fmt.Sprintf("%d days", s.days))
	}
}
