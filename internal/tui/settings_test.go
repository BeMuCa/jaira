package tui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/settings"
)

// TestSettingsOpensFromBoardAndLNoLongerLanes asserts S is the one door to
// settings, and that L does not reach the lane screen it used to open. L now
// opens the link window instead, so the guard is about where it does not go.
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

// TestSettingsEnterStepsTheChatModel: the fifth entry shows the chat model and
// where it comes from. enter steps the board's own (.jaira/chat-model) through
// sonnet, haiku, opus and none; with none, this machine's settings.json
// applies, and the row says so.
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
	if got := open(); !strings.Contains(got, "Chat model  account default — account default") {
		t.Fatalf("the chat model row does not show the account default:\n%s", got)
	}
	for _, want := range []string{"sonnet", "haiku", "opus", ""} {
		m.key(key("enter"))
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
	if got := open(); !strings.Contains(got, "Chat model  claude-sonnet-5-5 — from settings.json on this machine") {
		t.Fatalf("the machine's fallback is not shown with its source:\n%s", got)
	}
	m.key(key("enter"))
	if got := stripANSI(m.settingsScreen.render(150, 32)); !strings.Contains(got, "Chat model  sonnet — from .jaira/chat-model") {
		t.Fatalf("the board's model does not win over the machine's:\n%s", got)
	}
}
