package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// shift+down shows the next of the repository's boards: its lanes, and only
// its tickets. Without it the board is today's, unchanged, and one more
// shift+down comes back round to it.
func TestShiftDownPicksABoard(t *testing.T) {
	m := newTestModel(t, 150, 32)
	root := m.store.Root
	lanes := filepath.Join(root, ticket.DirName, "boards", "video", "lanes")
	if err := os.MkdirAll(lanes, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"script.md": "---\nid: script\nname: Skript\nprecedence: 10\nrequires-specified: true\n---\n",
		"voice.md":  "---\nid: voice\nname: Stimme\nafter: script\nprecedence: 20\n---\n",
	} {
		if err := os.WriteFile(filepath.Join(lanes, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	if _, err := m.store.Create(map[string]string{
		ticket.FieldID: ticket.NewID(now), ticket.FieldTitle: "Intro video",
		ticket.FieldBoard: "video", ticket.FieldStatus: "voice",
		ticket.FieldCreatedAt: ticket.FormatTime(now), ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}, nil, ""); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}

	out := stripANSI(m.render())
	if strings.Contains(out, "Intro video") || strings.Contains(out, "Stimme") {
		t.Errorf("today's board shows the video board's ticket or lane:\n%s", out)
	}
	if !strings.Contains(out, "Implementing") || !strings.Contains(out, "Refactor") {
		t.Errorf("today's board lost its own lanes or tickets:\n%s", out)
	}

	pressShift(m, "down")
	out = stripANSI(m.render())
	for _, want := range []string{"board: video", "Skript", "Stimme", "Intro video"} {
		if !strings.Contains(out, want) {
			t.Errorf("after shift+down, %q is missing:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Implementing") || strings.Contains(out, "Refactor") {
		t.Errorf("the video board shows today's lanes or tickets:\n%s", out)
	}

	// A ticket captured on the video board belongs to it and stays in view.
	m.createTicket("Outro video")
	found := false
	for _, tk := range m.tickets {
		if tk.Title == "Outro video" {
			found = tk.Board == "video"
		}
	}
	if !found {
		t.Errorf("a ticket created on the video board is not on it")
	}
	m.mode = modeBoard

	pressShift(m, "down")
	out = stripANSI(m.render())
	if strings.Contains(out, "board: video") || !strings.Contains(out, "Refactor") {
		t.Errorf("a second shift+down did not come back to today's board:\n%s", out)
	}
}

// With one board, shift+down has nothing to switch to and says how another
// comes about.
func TestShiftDownOnASingleBoardExplains(t *testing.T) {
	m := newTestModel(t, 150, 32)
	pressShift(m, "down")
	if m.boardName != "" || !strings.Contains(m.message, ".jaira/boards/") {
		t.Errorf("shift+down on a single board: board %q, message %q", m.boardName, m.message)
	}
}
