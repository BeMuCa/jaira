package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/settings"
)

// TestSettingsOpensFromBoardAndLNoLongerLanes asserts S is the one door to
// settings, and that L does not reach the lane screen it used to open. L now
// opens the milestone layer instead, so the guard is about where it does not go.
func TestSettingsOpensFromBoardAndLNoLongerLanes(t *testing.T) {
	m := newTestModel(t, 150, 32)

	m.key(key("L"))
	if m.mode == modeLanes {
		t.Fatal("L must no longer open the lane screen")
	}
	m.key(key("esc"))

	m.key(key("S"))
	if m.mode != modeSettings {
		t.Fatalf("S must open settings, mode = %v, want modeSettings", m.mode)
	}
	if m.settingsScreen == nil {
		t.Fatal("S must build a settingsScreen")
	}
}

// TestSettingsEnterOnLanesOpensLaneScreen asserts the Lanes entry, the first
// in the menu, reaches the existing lane screen.
func TestSettingsEnterOnLanesOpensLaneScreen(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.key(key("S"))

	m.key(key("enter"))
	if m.mode != modeLanes {
		t.Fatalf("enter on Lanes: mode = %v, want modeLanes", m.mode)
	}
	if m.laneScreen == nil {
		t.Fatal("enter on Lanes must build a laneScreen")
	}
}

// TestSettingsEnterOnDefaultBoardOpensDefaultBoardScreen asserts the second
// entry reaches the default board screen, previously reachable only from the
// launcher.
func TestSettingsEnterOnDefaultBoardOpensDefaultBoardScreen(t *testing.T) {
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))
	m := newTestModel(t, 150, 32)
	m.key(key("S"))
	m.key(key("j")) // move down to the second entry, Default board

	m.key(key("enter"))
	if m.mode != modeDefaultBoard {
		t.Fatalf("enter on Default board: mode = %v, want modeDefaultBoard", m.mode)
	}
	if m.board == nil {
		t.Fatal("enter on Default board must build a defaultBoardScreen")
	}
}

// TestSettingsEscStepsBackOneLevel asserts esc from an inner screen returns to
// the settings menu, not straight to the board, and esc again from the menu
// reaches the board.
func TestSettingsEscStepsBackOneLevel(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.key(key("S"))
	m.key(key("enter")) // into Lanes
	if m.mode != modeLanes {
		t.Fatalf("setup: mode = %v, want modeLanes", m.mode)
	}

	m.key(key("esc"))
	if m.mode != modeSettings {
		t.Fatalf("esc from lane screen: mode = %v, want modeSettings", m.mode)
	}
	if m.laneScreen != nil {
		t.Error("esc from lane screen must clear laneScreen")
	}

	m.key(key("esc"))
	if m.mode != modeBoard {
		t.Fatalf("esc from settings menu: mode = %v, want modeBoard", m.mode)
	}
}

// TestSettingsEnterSwitchesDispatch: the fourth entry is the opt-in to the
// worker mode. enter flips it and writes settings.json at once, and the row
// shows what the file now says — off ("lanes") until somebody turns it on.
func TestSettingsEnterSwitchesDispatch(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	m := newTestModel(t, 150, 32)
	m.key(key("S"))
	for range 3 {
		m.key(key("j"))
	}
	if got := m.settingsScreen.render(150, 32); m.settingsScreen.dispatch != settings.DispatchLanes || !strings.Contains(got, "lanes") {
		t.Fatalf("the dispatch row does not show the default, lanes (%q):\n%s", m.settingsScreen.dispatch, got)
	}

	m.key(key("enter"))
	if got := settings.Load().DispatchMode(); got != settings.DispatchWorker {
		t.Fatalf("after enter, settings.json dispatch = %q, want worker", got)
	}
	if m.mode != modeSettings {
		t.Errorf("switching must stay on the settings screen, mode = %v", m.mode)
	}

	m.key(key("enter"))
	if got := settings.Load().DispatchMode(); got != settings.DispatchLanes {
		t.Fatalf("after a second enter, dispatch = %q, want lanes", got)
	}
}

// TestSettingsEnterStepsTheChatModel: the fifth entry shows the board's chat
// model and where it comes from. enter steps the board's own
// (.jaira/chat-model) through sonnet, haiku, opus and none; at none the row
// says "none" and what applies instead — this machine's settings.json, or the
// account default. A value in the file that is no model name is none too.
func TestSettingsEnterStepsTheChatModel(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	m := newTestModel(t, 150, 32)
	open := func() string {
		t.Helper()
		m.key(key("esc"))
		m.key(key("S"))
		for range 4 {
			m.key(key("j"))
		}
		return stripANSI(m.settingsScreen.render(150, 32))
	}
	if got := open(); !strings.Contains(got, "Chat model  none — account default") {
		t.Fatalf("the chat model row does not show the account default:\n%s", got)
	}
	for _, want := range []string{"sonnet", "haiku", "opus", ""} {
		m.key(key("enter"))
		label := want
		if want == "" {
			label = "none — account default"
		}
		if got := stripANSI(m.settingsScreen.render(150, 32)); !strings.Contains(got, "Chat model  "+label) {
			t.Fatalf("after enter, the row does not read %q:\n%s", label, got)
		}
		if got := settings.BoardChatModel(m.store.Root); got != want {
			t.Fatalf("after enter, the board's chat model = %q, want %q", got, want)
		}
		if got := settings.Load().ChatModel; got != "" {
			t.Fatalf("enter wrote the machine's settings.json too: %q", got)
		}
	}
	if m.mode != modeSettings {
		t.Errorf("stepping must stay on the settings screen, mode = %v", m.mode)
	}

	if err := settings.Save(settings.Settings{ChatModel: "claude-sonnet-5-5"}); err != nil {
		t.Fatal(err)
	}
	if got := open(); !strings.Contains(got, "Chat model  none — claude-sonnet-5-5 from settings.json on this machine") {
		t.Fatalf("the machine's fallback is not shown with its source:\n%s", got)
	}
	m.key(key("enter"))
	if got := stripANSI(m.settingsScreen.render(150, 32)); !strings.Contains(got, "Chat model  sonnet — from .jaira/chat-model") {
		t.Fatalf("the board's model does not win over the machine's:\n%s", got)
	}

	// A hand-written value that is no model name counts as none: the machine's
	// model applies, the row reads none, and enter goes on to sonnet.
	if err := os.WriteFile(settings.BoardChatModelPath(m.store.Root), []byte("$(id)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := open(); !strings.Contains(got, "Chat model  none — claude-sonnet-5-5 from settings.json on this machine") {
		t.Fatalf("an invalid board value does not read as none:\n%s", got)
	}
	if got, _ := settings.Load().ChatModelFor(m.store.Root); got != "claude-sonnet-5-5" {
		t.Fatalf("an invalid board value must fall through to the machine's: %q", got)
	}
	m.key(key("enter"))
	if got := settings.BoardChatModel(m.store.Root); got != "sonnet" {
		t.Errorf("enter on an invalid board value: %q, want sonnet, the step after none", got)
	}
}

// TestSettingsEnterStepsTheLanguage: the sixth entry is the language the
// cockpit shows review notes in. enter steps it through Deutsch, English,
// Русский and none and writes this machine's settings.json at once - nothing
// of it reaches the board. Two people, two files: what one sets the other
// does not see. A language written by hand is shown as it stands, and enter
// steps it to none.
func TestSettingsEnterStepsTheLanguage(t *testing.T) {
	// newTestModel points JAIRA_HOME at a home of its own: that is Alex's.
	m := newTestModel(t, 150, 32)
	alex, berk := os.Getenv("JAIRA_HOME"), t.TempDir()
	open := func() string {
		t.Helper()
		m.key(key("esc"))
		m.key(key("S"))
		for range 5 {
			m.key(key("j"))
		}
		return stripANSI(m.settingsScreen.render(150, 32))
	}
	if got := open(); !strings.Contains(got, "Notes language  none — notes as the ticket has them") {
		t.Fatalf("the language row does not show none:\n%s", got)
	}
	before, err := os.ReadDir(filepath.Join(m.store.Root, ".jaira"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Deutsch", "English", "Русский", ""} {
		m.key(key("enter"))
		label := want
		if want == "" {
			label = "none — notes as the ticket has them"
		}
		if got := stripANSI(m.settingsScreen.render(150, 32)); !strings.Contains(got, "Notes language  "+label) {
			t.Fatalf("after enter, the row does not read %q:\n%s", label, got)
		}
		if got := string(settings.Load().Language); got != want {
			t.Fatalf("after enter, settings.json language = %q, want %q", got, want)
		}
	}
	after, err := os.ReadDir(filepath.Join(m.store.Root, ".jaira"))
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Errorf("stepping the language wrote into the board: %d entries before, %d after", len(before), len(after))
	}

	// Alex reads Russian; Berk, on another machine, German.
	m.key(key("enter"))
	m.key(key("enter"))
	m.key(key("enter"))
	t.Setenv("JAIRA_HOME", berk)
	m.key(key("enter"))
	if got := string(settings.Load().Language); got != "Deutsch" {
		t.Fatalf("Berk's settings.json: %q, want Deutsch", got)
	}
	t.Setenv("JAIRA_HOME", alex)
	if got := string(settings.Load().Language); got != "Русский" {
		t.Fatalf("Alex's settings.json: %q, want Русский after Berk set theirs", got)
	}

	// By hand, a language that is none of the choices: shown, and enter steps to none.
	if err := settings.Save(settings.Settings{Language: "Українська", Dispatch: "worker"}); err != nil {
		t.Fatal(err)
	}
	if got := open(); !strings.Contains(got, "Notes language  Українська") {
		t.Fatalf("a hand-written language is not shown:\n%s", got)
	}
	m.key(key("enter"))
	if got := settings.Load(); got.Language != "" || got.Dispatch != "worker" {
		t.Errorf("enter on a hand-written language: %+v, want none and the rest kept", got)
	}
}

// TestSettingsWritesNeverWipeABrokenFile: every row of S that writes
// settings.json - logbook days, agents per ticket, notes language - refuses
// a file that does not parse, says so on the row, and leaves it as it was.
func TestSettingsWritesNeverWipeABrokenFile(t *testing.T) {
	broken := `{"remote":"upstream","hook":"/bin/notify",}`
	for _, row := range []struct {
		name string
		down int
		keys []string
	}{
		{"logbook days", 2, []string{"enter", "1", "enter"}},
		{"agents per ticket", 3, []string{"enter"}},
		{"notes language", 5, []string{"enter"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			m := newTestModel(t, 150, 32)
			path := settings.Path()
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
				t.Fatal(err)
			}
			m.key(key("S"))
			for range row.down {
				m.key(key("j"))
			}
			for _, k := range row.keys {
				m.key(key(k))
			}
			if b, _ := os.ReadFile(path); string(b) != broken {
				t.Fatalf("settings.json was written over: %q", b)
			}
			if got := stripANSI(m.settingsScreen.render(150, 40)); !strings.Contains(got, "does not parse") {
				t.Fatalf("the row does not say the file does not parse:\n%s", got)
			}
		})
	}
}

// The editor and IDE rows type their value in place: enter opens it with what
// settings.json has, typed text and a paste land in it, enter saves it through
// settings.Update, esc leaves it as it was, and an empty value saves the default.
func TestSettingsTypesTheEditorAndIDE(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	m := newTestModel(t, 150, 40)
	open := func(rows int) {
		t.Helper()
		m.key(key("esc"))
		m.key(key("S"))
		for range rows {
			m.key(key("j"))
		}
	}
	row := func() string {
		t.Helper()
		return stripANSI(m.settingsScreen.render(150, 40))
	}

	open(6)
	m.key(key("enter"))
	for _, c := range "nvim" {
		m.key(key(string(c)))
	}
	m.Update(tea.PasteMsg{Content: " --clean"})
	if got := row(); !strings.Contains(got, "Editor (E)  [nvim --clean_]") {
		t.Fatalf("the editor row does not show what was typed:\n%s", got)
	}
	m.key(key("enter"))
	if got := settings.Load().Editor; got != "nvim --clean" {
		t.Fatalf("settings.json editor = %q, want %q", got, "nvim --clean")
	}
	if got := row(); !strings.Contains(got, "Editor (E)  nvim --clean") {
		t.Fatalf("after save, the editor row does not show the value:\n%s", got)
	}

	open(7)
	m.key(key("enter"))
	if got := row(); !strings.Contains(got, "IDE (o in the diff)  [_]") {
		t.Fatalf("the IDE row does not open empty when none is set:\n%s", got)
	}
	m.Update(tea.PasteMsg{Content: "x"})
	m.key(key("esc"))
	if m.settingsScreen.typing != "" || settings.Load().IDE != "" {
		t.Fatalf("esc did not leave the IDE as it was: typing=%q ide=%q", m.settingsScreen.typing, settings.Load().IDE)
	}

	m.key(key("enter"))
	m.Update(tea.PasteMsg{Content: "idea --line {line} {file}"})
	m.key(key("enter"))
	if got := settings.Load().IDETemplate(); got != "idea --line {line} {file}" {
		t.Fatalf("settings.json ide = %q", got)
	}

	m.key(key("enter"))
	for range 40 {
		m.key(key("backspace"))
	}
	m.key(key("enter"))
	if got := settings.Load().IDE; got != "" {
		t.Fatalf("an empty IDE should save the default, settings.json has %q", got)
	}
	if got := row(); !strings.Contains(got, "default — code -g {file}:{line}") {
		t.Fatalf("after an empty save, the IDE row does not show the default:\n%s", got)
	}
}
