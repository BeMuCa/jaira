package tui

import "strings"

// Actions settingsScreen.key can report: which entry was opened, or that the
// screen itself is finished.
const (
	settingsActionNone         = ""
	settingsActionBack         = "back"
	settingsActionLanes        = "lanes"
	settingsActionDefaultBoard = "default-board"
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
}

// settingsScreen is the menu behind S: one door into both the lane screen and
// the default board, rather than each having its own binding. Following
// browse.go's shape: its own state, a key method, a render method, no
// bubbletea imports of its own.
type settingsScreen struct {
	idx int
}

func newSettingsScreen() *settingsScreen {
	return &settingsScreen{}
}

// key drives the screen. It reports settingsActionBack when esc/q finish the
// screen, the action of the highlighted entry when enter opens it, or
// settingsActionNone when the key only moved the cursor.
func (s *settingsScreen) key(k string) string {
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
		return settingsEntries[s.idx].action
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
		sb.WriteString(truncate(lead+name, w) + "\n")
		sb.WriteString(truncate("      "+styMeta.Render(e.desc), w) + "\n")
	}

	sb.WriteString("\n" + styMeta.Render(truncate("enter open · esc back", w)))
	return sb.String()
}
