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
	// settingsActionDispatch flips the dispatch setting in place: there are
	// two values, so enter is the whole editor.
	settingsActionDispatch = "dispatch"
	// settingsActionChatModel steps the chat model through its choices in
	// place, the same way.
	settingsActionChatModel = "chat-model"
	// settingsActionModel steps one agent's model and effort through
	// settings.ModelChoices in place; the entry's key says which agent.
	settingsActionModel = "model"
)

// settingsEntry is one row in the menu: a name, what pressing enter on it
// opens, and a one-line description of what it is for. key is the agent a
// settingsActionModel row is for, as settings.ModelFor takes it.
type settingsEntry struct {
	name   string
	desc   string
	action string
	key    string
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
	{
		name:   "Agents per ticket",
		desc:   "lanes: one new agent per lane (jaira-dispatcher); worker: one agent carries the ticket, only a fresh-context lane such as review gets a new one (jaira-worker). enter switches",
		action: settingsActionDispatch,
	},
	{
		name:   "Chat model",
		desc:   "the model the cockpit's explain and discuss chats start on, for everybody on this board (.jaira/chat-model); without it this machine's chat-model in settings.json, else the account default. enter steps the board's through sonnet, haiku, opus and none; any other model by hand",
		action: settingsActionChatModel,
	},
	{
		name:   "Orchestrator model",
		desc:   "model + effort for the teamlead session; enter steps through opus, opus high, sonnet, sonnet high, haiku, haiku high, then the default; any other by hand in settings.json",
		action: settingsActionModel,
		key:    "orchestrator",
	},
	{
		name:   "Worker model",
		desc:   "model + effort for the worker, the session that carries a ticket through its lanes; enter steps as above; any other by hand in settings.json",
		action: settingsActionModel,
		key:    "worker",
	},
	{
		name:   "Strong lanes model",
		desc:   "model + effort for lanes of tier strong (brainstorm, pre-process, critique, optimize, review); enter steps as above; any other by hand in settings.json",
		action: settingsActionModel,
		key:    "strong",
	},
	{
		name:   "Cheap lanes model",
		desc:   "model + effort for lanes of tier cheap (in-progress, testing); enter steps as above; any other by hand in settings.json",
		action: settingsActionModel,
		key:    "cheap",
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

	// dispatch is the dispatch setting as settings.json has it now, and
	// dispatchErr what went wrong saving the last switch.
	dispatch    string
	dispatchErr string

	// root is the board the screen was opened on. chatBoard is the model its
	// .jaira/chat-model names, "" for none; chatModel is the model its chats
	// actually start on, chatModelSource where that came from, and
	// chatModelErr what went wrong saving the last step.
	root            string
	chatBoard       string
	chatModel       string
	chatModelSource string
	chatModelErr    string

	// models is settings.json as it stands, for the four model rows. modelErr
	// is what went wrong saving the last step, and modelErrKey the agent whose
	// row it belongs under.
	models      settings.Settings
	modelErr    string
	modelErrKey string
}

func newSettingsScreen(root string) *settingsScreen {
	cur := settings.Load()
	s := &settingsScreen{days: cur.LogbookWindow(), dispatch: cur.DispatchMode(), root: root, models: cur}
	s.readChatModel(cur)
	return s
}

// readChatModel takes the board's chat model and the one that applies from
// the files as they are now.
func (s *settingsScreen) readChatModel(cur settings.Settings) {
	if s.root != "" {
		s.chatBoard = settings.BoardChatModel(s.root)
	}
	s.chatModel, s.chatModelSource = cur.ChatModelFor(s.root)
}

// chatModelValue is the chat model row's value: the board's model and where
// it comes from, or "none" and what applies instead — the step enter is on
// is always the word the row starts with.
func (s *settingsScreen) chatModelValue() string {
	switch {
	case s.chatBoard != "":
		return s.chatBoard + " — " + s.chatModelSource
	case s.chatModel != "":
		return "none — " + s.chatModel + " " + s.chatModelSource
	default:
		return "none — account default"
	}
}

// stepChatModel moves the board's chat model on to the next of its choices
// and saves at once — the board's, because it is the one that wins; the
// machine's fallback is written by hand. A model written by hand that is none
// of the choices (claude-sonnet-5-5) steps to none, the first, which lets the
// machine's apply again. A value that is no model name at all already counts
// as none, so it steps on to sonnet like none does.
func (s *settingsScreen) stepChatModel() string {
	if s.root == "" {
		return settingsActionNone
	}
	choices := settings.ChatModelChoices
	// Read again rather than trusting chatBoard, for the same reason editKey
	// reads settings.json again: the file may have changed since S opened.
	cur, next := settings.BoardChatModel(s.root), choices[0]
	for i, c := range choices {
		if c == cur {
			next = choices[(i+1)%len(choices)]
		}
	}
	if err := settings.SetBoardChatModel(s.root, next); err != nil {
		s.chatModelErr = err.Error()
		return settingsActionNone
	}
	// None, not saved: like dispatch, nothing on the board shows it.
	s.readChatModel(settings.Load())
	s.chatModelErr = ""
	return settingsActionNone
}

// switchDispatch flips lanes and worker and saves at once, reading the file
// again first for the same reason editKey does.
func (s *settingsScreen) switchDispatch() string {
	cur := settings.Load()
	next := settings.DispatchWorker
	if cur.DispatchMode() == settings.DispatchWorker {
		next = settings.DispatchLanes
	}
	cur.Dispatch = next
	if err := settings.Save(cur); err != nil {
		s.dispatchErr = err.Error()
		return settingsActionNone
	}
	// None, not saved: nothing on the board shows this setting, so there is
	// nothing to reload — the row above already says what the file now holds.
	s.dispatch, s.dispatchErr = next, ""
	return settingsActionNone
}

// nextModelChoice is the entry an agent's model steps on to: the first choice
// when the file names none, the one after the current choice, and no entry
// (keep=false) after the last choice. A value written by hand that is none of
// the choices steps to no entry too, the default.
func nextModelChoice(now settings.AgentModel, set bool) (next settings.AgentModel, keep bool) {
	if !set {
		return settings.ModelChoices[0], true
	}
	for i, c := range settings.ModelChoices {
		if c == now && i+1 < len(settings.ModelChoices) {
			return settings.ModelChoices[i+1], true
		}
	}
	return settings.AgentModel{}, false
}

// stepModel moves one agent's model on to the next of its choices and saves at
// once, reading the file again first for the same reason editKey does. The
// last step removes the agent's entry, so the default applies again.
func (s *settingsScreen) stepModel(key string) string {
	cur := settings.Load()
	now, set := cur.Models[key]
	next, keep := nextModelChoice(now, set)
	if keep {
		if cur.Models == nil {
			cur.Models = map[string]settings.AgentModel{}
		}
		cur.Models[key] = next
	} else {
		delete(cur.Models, key)
	}
	if err := settings.Save(cur); err != nil {
		s.modelErr, s.modelErrKey = err.Error(), key
		return settingsActionNone
	}
	s.models, s.modelErr, s.modelErrKey = cur, "", ""
	return settingsActionNone
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
		a := settingsEntries[s.idx].action
		if a == settingsActionDispatch {
			return s.switchDispatch()
		}
		if a == settingsActionChatModel {
			return s.stepChatModel()
		}
		if a == settingsActionModel {
			return s.stepModel(settingsEntries[s.idx].key)
		}
		if a != settingsActionLogbookDays {
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
		if e.action == settingsActionDispatch {
			name += "  " + styMeta.Render(s.dispatch)
		}
		if e.action == settingsActionChatModel {
			name += "  " + styMeta.Render(s.chatModelValue())
		}
		if e.action == settingsActionModel {
			name += "  " + styMeta.Render(s.modelValue(e.key))
		}
		sb.WriteString(truncate(lead+name, w) + "\n")
		sb.WriteString(styMeta.Render("      "+wrap(e.desc, max(10, w-6), 6)) + "\n")
		if e.action == settingsActionLogbookDays && s.err != "" {
			sb.WriteString(styErr.Render("      "+s.err) + "\n")
		}
		if e.action == settingsActionDispatch && s.dispatchErr != "" {
			sb.WriteString(styErr.Render("      "+s.dispatchErr) + "\n")
		}
		if e.action == settingsActionChatModel && s.chatModelErr != "" {
			sb.WriteString(styErr.Render("      "+s.chatModelErr) + "\n")
		}
		if e.action == settingsActionModel && s.modelErr != "" && e.key == s.modelErrKey {
			sb.WriteString(styErr.Render("      "+s.modelErr) + "\n")
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

// modelValue is one model row's value: the model and effort the agent starts
// on, and "(default)" when settings.json does not name it.
func (s *settingsScreen) modelValue(key string) string {
	m, src := s.models.ModelFor(key)
	if src == "the default" {
		return m.String() + " (default)"
	}
	return m.String()
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
