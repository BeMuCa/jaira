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
	// settingsActionLanguage steps the notes language through its choices
	// in place, the same way.
	settingsActionLanguage = "language"
	// settingsActionEditor and settingsActionIDE open a value for typing in
	// place: enter saves it, esc leaves it as it was, and an empty value saves
	// the default.
	settingsActionEditor = "editor"
	settingsActionIDE    = "ide"
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
		name:   "Notes language",
		desc:   "the language the cockpit shows a review's notes in, rewritten short by a cheap model, for you only (language in settings.json on this machine); none shows them as the ticket has them. enter steps through Deutsch, English, Русский and none; any other by hand",
		action: settingsActionLanguage,
	},
	{
		name:   "Editor (E)",
		desc:   "the editor E opens a ticket in, with its arguments ('nvim', 'code --wait'); jaira waits until it exits. enter types it, enter saves, an empty value means VISUAL, EDITOR, then git's core.editor",
		action: settingsActionEditor,
	},
	{
		name:   "IDE (o in the diff)",
		desc:   "the program o opens a file from the diff in, at its line; {file} and {line} are filled in, e.g. 'code -g {file}:{line}'. jaira does not wait for it. enter types it, enter saves, an empty value means the default",
		action: settingsActionIDE,
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

	// language is the notes language as settings.json has it now, "" for
	// none, and languageErr what went wrong saving the last step.
	language    string
	languageErr string

	// editor and ide are the two settings as settings.json has them now, "" for
	// the default; editorNow is the editor that applies when editor is "", and
	// where it comes from. typing is the action of the row whose value is open
	// for typing, "" when none; typed is what has been typed so far and textErr
	// what went wrong saving it.
	editor, ide string
	editorNow   string
	typing      string
	typed       string
	textErr     string
}

func newSettingsScreen(root string) *settingsScreen {
	cur := settings.Load()
	s := &settingsScreen{days: cur.LogbookWindow(), dispatch: cur.DispatchMode(), root: root, language: string(cur.Language)}
	s.readChatModel(cur)
	s.readText(cur)
	return s
}

// readText takes the editor and IDE as settings.json has them, and the editor
// that applies when none is set.
func (s *settingsScreen) readText(cur settings.Settings) {
	s.editor, s.ide = cur.Editor, cur.IDE
	argv, source := editorChoice(cur)
	s.editorNow = strings.Join(argv, " ") + " (" + source + ")"
}

// openText opens the editor or IDE row for typing, with the value settings.json
// has for it now.
func (s *settingsScreen) openText(action string) {
	s.typing, s.typed, s.textErr = action, s.editor, ""
	if action == settingsActionIDE {
		s.typed = s.ide
	}
}

// textKey drives the value while it is open for typing: enter saves it, esc
// leaves it as it was and backspace takes one character back. The characters
// themselves come through addText.
func (s *settingsScreen) textKey(k string) string {
	switch k {
	case "esc":
		s.typing, s.typed, s.textErr = "", "", ""
	case "backspace":
		if r := []rune(s.typed); len(r) > 0 {
			s.typed = string(r[:len(r)-1])
		}
	case "enter":
		return s.saveText()
	}
	return settingsActionNone
}

// addText types text into the open value, a key's character or a paste, on one
// line. It does nothing when no value is open.
func (s *settingsScreen) addText(text string) {
	if s.typing != "" {
		s.typed += foldToOneLine(text)
	}
}

// saveText writes the open value through settings.Update, like every write of
// this screen. The value stays open when the write fails, so nothing typed is
// lost.
func (s *settingsScreen) saveText() string {
	action, value := s.typing, strings.TrimSpace(s.typed)
	now, err := settings.Update(func(cur *settings.Settings) {
		if action == settingsActionIDE {
			cur.IDE = value
		} else {
			cur.Editor = value
		}
	})
	if err != nil {
		s.textErr = err.Error()
		return settingsActionNone
	}
	// None, not returned: nothing on the board shows these, and the diff pane and
	// the ticket editor read them when they open.
	s.readText(now)
	s.typing, s.typed, s.textErr = "", "", ""
	return settingsActionNone
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
	// Read again rather than trusting chatBoard, for the same reason editKey
	// reads settings.json again: the file may have changed since S opened.
	next := stepAfter(settings.ChatModelChoices, settings.BoardChatModel(s.root))
	if err := settings.SetBoardChatModel(s.root, next); err != nil {
		s.chatModelErr = err.Error()
		return settingsActionNone
	}
	// None, not saved: like dispatch, nothing on the board shows it.
	s.readChatModel(settings.Load())
	s.chatModelErr = ""
	return settingsActionNone
}

// stepAfter is the choice after cur, wrapping round; a cur that is none of
// them steps to the first.
func stepAfter(choices []string, cur string) string {
	for i, c := range choices {
		if c == cur {
			return choices[(i+1)%len(choices)]
		}
	}
	return choices[0]
}

// stepLanguage moves the notes language on to the next of its choices and
// saves at once, through settings.Update like every write of this screen. A
// language written by hand that is none of the choices steps to none, the
// first, as the chat model does.
func (s *settingsScreen) stepLanguage() string {
	now, err := settings.Update(func(cur *settings.Settings) {
		cur.Language = settings.Language(stepAfter(settings.LanguageChoices, string(cur.Language)))
	})
	if err != nil {
		s.languageErr = err.Error()
		return settingsActionNone
	}
	// None, not saved: nothing on the board shows it, only the cockpit.
	s.language, s.languageErr = string(now.Language), ""
	return settingsActionNone
}

// switchDispatch flips lanes and worker and saves at once, through
// settings.Update like every write of this screen.
func (s *settingsScreen) switchDispatch() string {
	now, err := settings.Update(func(cur *settings.Settings) {
		next := settings.DispatchWorker
		if cur.DispatchMode() == settings.DispatchWorker {
			next = settings.DispatchLanes
		}
		cur.Dispatch = next
	})
	if err != nil {
		s.dispatchErr = err.Error()
		return settingsActionNone
	}
	next := now.Dispatch
	// None, not saved: nothing on the board shows this setting, so there is
	// nothing to reload — the row above already says what the file now holds.
	s.dispatch, s.dispatchErr = next, ""
	return settingsActionNone
}

// editKey drives the number while it is open: digits type, backspace takes
// one back, enter saves and esc leaves it as it was. Saving goes through
// settings.Update, so nothing else in it — a remote, a hook — is written back
// from a copy older than the edit, or wiped where the file does not parse.
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
		if _, err := settings.Update(func(cur *settings.Settings) { cur.LogbookDays = &n }); err != nil {
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
		if a == settingsActionLanguage {
			return s.stepLanguage()
		}
		if a == settingsActionEditor || a == settingsActionIDE {
			s.openText(a)
			return settingsActionNone
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
		if e.action == settingsActionLanguage {
			lang := s.language
			if lang == "" {
				lang = "none — notes as the ticket has them"
			}
			name += "  " + styMeta.Render(lang)
		}
		if e.action == settingsActionEditor || e.action == settingsActionIDE {
			name += "  " + s.textValue(e.action)
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
		if e.action == settingsActionLanguage && s.languageErr != "" {
			sb.WriteString(styErr.Render("      "+s.languageErr) + "\n")
		}
		if e.action == s.typing && s.textErr != "" {
			sb.WriteString(styErr.Render("      "+s.textErr) + "\n")
		}
	}

	hints := []string{"enter open", "esc back"}
	if s.editing {
		hints = []string{"0-9 days", "enter save", "esc cancel"}
	}
	if s.typing != "" {
		hints = []string{"type", "enter save", "esc cancel"}
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

// textValue is the editor or IDE row's value: the value being typed while it is
// open, otherwise what settings.json says, and for an empty one the default in
// words.
func (s *settingsScreen) textValue(action string) string {
	if s.typing == action {
		return stySelected.Render("[" + s.typed + "_]")
	}
	set, def := s.editor, "default — "+s.editorNow
	if action == settingsActionIDE {
		set, def = s.ide, "default — "+settings.DefaultIDE
	}
	if set == "" {
		return styMeta.Render(def)
	}
	return styMeta.Render(set)
}
