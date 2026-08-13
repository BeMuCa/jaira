package tui

import (
	"path/filepath"
	"testing"
)

// TestSettingsOpensFromBoardAndLNoLonger asserts S is the one door: it opens
// the settings menu, and L — the old, removed binding — does nothing.
func TestSettingsOpensFromBoardAndLNoLonger(t *testing.T) {
	m := newTestModel(t, 150, 32)

	m.key(key("L"))
	if m.mode != modeBoard {
		t.Fatalf("L must no longer open anything, mode = %v, want modeBoard", m.mode)
	}

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
