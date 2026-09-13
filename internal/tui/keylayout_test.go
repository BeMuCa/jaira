package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// A board that only answers a US layout is a board a Russian-speaking user
// cannot steer at all, so every command key has to survive the trip through
// ЙЦУКЕН.
func TestCmdKeyMapsCyrillicToItsPhysicalKey(t *testing.T) {
	for _, tc := range []struct {
		typed rune
		want  string
	}{
		{'о', "j"}, {'л', "k"}, {'р', "h"}, {'д', "l"},
		{'й', "q"}, {'у', "e"}, {'г', "u"}, {'ф', "a"},
		{'ч', "x"}, {'м', "v"}, {'т', "n"}, {'.', "/"},
	} {
		got := cmdKey(tea.KeyPressMsg{Code: tc.typed, Text: string(tc.typed)})
		if got != tc.want {
			t.Errorf("cmdKey(%q) = %q, want %q", tc.typed, got, tc.want)
		}
	}
}

// The board binds "G", "E" and "X" separately from their lower case, so the
// case that was typed has to survive the mapping.
func TestCmdKeyKeepsCase(t *testing.T) {
	got := cmdKey(tea.KeyPressMsg{Code: 'П', Text: "П"})
	if got != "G" {
		t.Errorf("cmdKey(П) = %q, want %q", got, "G")
	}
}

// Mapping the character must not eat the modifier that was held with it.
func TestCmdKeyKeepsModifiers(t *testing.T) {
	got := cmdKey(tea.KeyPressMsg{Code: 'в', Text: "в", Mod: tea.ModCtrl})
	if got != "ctrl+d" {
		t.Errorf("cmdKey(ctrl+в) = %q, want %q", got, "ctrl+d")
	}
}

// A US layout must come out exactly as it went in — named keys included, which
// are spelled out rather than carried as a character.
func TestCmdKeyLeavesLatinAndNamedKeysAlone(t *testing.T) {
	for _, k := range []tea.KeyPressMsg{
		{Code: 'j', Text: "j"},
		{Code: 'G', Text: "G"},
		{Code: '/', Text: "/"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEscape},
		{Code: tea.KeyUp},
		{Code: 'd', Mod: tea.ModCtrl},
	} {
		if got, want := cmdKey(k), k.String(); got != want {
			t.Errorf("cmdKey(%q) = %q, want it unchanged", want, got)
		}
	}
}

// A terminal that speaks the Kitty keyboard protocol reports the physical key
// itself, and that answer outranks the table: it is right for layouts this
// package has never heard of.
func TestCmdKeyPrefersTheTerminalsOwnBaseCode(t *testing.T) {
	// Greek: no table here, and none needed when the terminal reports it.
	got := cmdKey(tea.KeyPressMsg{Code: 'ξ', Text: "ξ", BaseCode: 'j'})
	if got != "j" {
		t.Errorf("cmdKey(ξ with BaseCode j) = %q, want %q", got, "j")
	}
}

// The fault this fixes end to end: with a Cyrillic layout selected, the board
// used to answer nothing at all, because every key arrived as a letter no case
// matched.
func TestBoardAnswersACyrillicLayout(t *testing.T) {
	m := newTestModel(t, 150, 32)
	// Put the cursor somewhere it can move down from.
	m.laneIdx, m.cardIdx = 0, 0
	for len(m.cols[m.laneIdx].tickets) < 2 {
		m.laneIdx++
		if m.laneIdx >= len(m.cols) {
			t.Skip("no lane in the fixture holds two tickets")
		}
	}

	press(m, "о") // the physical "j"
	if m.cardIdx != 1 {
		t.Errorf("о did not move the cursor down: cardIdx = %d, want 1", m.cardIdx)
	}
	press(m, "л") // the physical "k"
	if m.cardIdx != 0 {
		t.Errorf("л did not move the cursor back up: cardIdx = %d, want 0", m.cardIdx)
	}
	press(m, ".") // the physical "/"
	if m.mode != modeFilter {
		t.Errorf("the full stop key did not open the filter: mode = %v", m.mode)
	}
}

// The other half of the same fault: the mapping must not follow the user into a
// text field, or the board becomes unable to write down anything in Russian.
func TestTypingStaysCyrillicInTheFilter(t *testing.T) {
	m := newTestModel(t, 150, 32)
	press(m, ".") // opens the filter
	if m.mode != modeFilter {
		t.Fatalf("filter did not open: mode = %v", m.mode)
	}
	for _, r := range "отчёт" {
		press(m, string(r))
	}
	if m.input != "отчёт" {
		t.Errorf("filter holds %q, want %q", m.input, "отчёт")
	}
}

// A terminal that reports a BaseCode reports one for enter and space as well.
// Reading it there would hand the switches a bare rune instead of the name they
// match on, and the fault would only show up on the terminals the BaseCode
// branch exists for.
func TestCmdKeyIgnoresBaseCodeOnNamedKeys(t *testing.T) {
	for _, k := range []tea.KeyPressMsg{
		{Code: tea.KeyEnter, BaseCode: '\r'},
		{Code: tea.KeySpace, Text: " ", BaseCode: ' '},
		{Code: tea.KeyUp, BaseCode: tea.KeyUp},
		{Code: tea.KeyTab, BaseCode: '\t'},
	} {
		if got, want := cmdKey(k), k.String(); got != want {
			t.Errorf("cmdKey(%q) = %q, want it unchanged", want, got)
		}
	}
}
