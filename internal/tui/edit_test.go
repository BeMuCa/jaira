package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/berk/jaira/core/ticket"
)

func key(s string) tea.KeyPressMsg {
	if len(s) == 1 {
		return tea.KeyPressMsg{Code: rune(s[0]), Text: s}
	}
	return tea.KeyPressMsg{Code: keyCodeFor(s), Text: ""}
}

func keyCodeFor(s string) rune {
	switch s {
	case "enter":
		return tea.KeyEnter
	case "esc":
		return tea.KeyEscape
	case "tab":
		return tea.KeyTab
	case "backspace":
		return tea.KeyBackspace
	}
	return 0
}

// The complaint that started this: a ticket could be created from the board but
// not filled in, so the only way to give it a goal was to quit and use the CLI.
func TestEditingAFieldFromTheBoard(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.mode = modeBoard
	// Select a real ticket and open it.
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	if m.mode != modeDetail || m.detail == nil {
		t.Fatalf("enter did not open the detail pane (mode=%v)", m.mode)
	}
	before := m.detail.ID

	m.key(key("e"))
	if m.mode != modeEdit {
		t.Fatalf("e did not start editing (mode=%v)", m.mode)
	}
	for _, r := range "hello" {
		m.key(key(string(r)))
	}
	m.key(key("enter"))

	reloaded, err := m.store.Load(before)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reloaded.Goal, "hello") {
		t.Errorf("goal was not written: %q", reloaded.Goal)
	}
}

// Escape must leave the ticket exactly as it was.
func TestEscapeAbandonsAnEdit(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	original := m.detail.Goal
	id := m.detail.ID

	m.key(key("e"))
	for _, r := range "zzz" {
		m.key(key(string(r)))
	}
	m.key(key("esc"))

	reloaded, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Goal != original {
		t.Errorf("goal changed after escape: %q -> %q", original, reloaded.Goal)
	}
	if m.mode != modeDetail {
		t.Errorf("escape did not return to the detail pane (mode=%v)", m.mode)
	}
}

// Tab moves between fields so a whole ticket can be filled in without leaving
// the board — the round trip that made the CLI mandatory.
func TestTabCyclesFields(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	m.key(key("e"))
	first := editableFields[m.editIdx].name
	m.key(key("tab"))
	second := editableFields[m.editIdx].name
	if first == second {
		t.Error("tab did not move to another field")
	}
}

// Non-ASCII input must survive. The previous line editor gated on len(s) == 1,
// which counts bytes, so every umlaut was silently dropped — and the tickets
// being migrated into this board are German.
func TestMultiByteInputIsAccepted(t *testing.T) {
	m := newTestModel(t, 120, 32)
	m.laneIdx, m.cardIdx = 0, 0
	m.key(key("enter"))
	id := m.detail.ID
	m.key(key("e"))
	for _, r := range "Knöpfe klären" {
		m.key(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	m.key(key("enter"))

	reloaded, err := m.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reloaded.Goal, "Knöpfe klären") {
		t.Errorf("multi-byte characters were dropped: %q", reloaded.Goal)
	}
}

var _ = ticket.FieldGoal
