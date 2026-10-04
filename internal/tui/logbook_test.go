package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// fileIntoLogbook puts a finished ticket into the logbook folder of the day
// ago days back, the way 'jaira logbook' leaves it.
func fileIntoLogbook(t *testing.T, s *ticket.Store, title string, ago int) string {
	t.Helper()
	tk, err := s.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  title,
		ticket.FieldStatus: "done",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	folder := "as-" + time.Now().AddDate(0, 0, -ago).Format("20060102")
	if _, err := s.Logbook(tk.ID, folder); err != nil {
		t.Fatal(err)
	}
	return tk.ID
}

func setLogbookDays(t *testing.T, n int) {
	t.Helper()
	if err := settings.Save(settings.Settings{LogbookDays: &n}); err != nil {
		t.Fatal(err)
	}
}

// doneColumn is the terminal lane's column of the board.
func doneColumn(t *testing.T, m *Model) (int, column) {
	t.Helper()
	for i, c := range m.cols {
		if c.lane.Terminal {
			return i, c
		}
	}
	t.Fatal("no terminal lane on the board")
	return 0, column{}
}

// The terminal lane shows the board's own tickets first and then the logbook
// of the window, newest first, nothing older — and its count stays the count
// of what is on the board.
func TestLogbookCardsSitBelowTheTerminalLane(t *testing.T) {
	s := newTestStore(t)
	fileIntoLogbook(t, s, "Filed five days ago", 5)
	fileIntoLogbook(t, s, "Filed yesterday", 1)
	fileIntoLogbook(t, s, "Filed long ago", 60)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 200, 40

	idx, col := doneColumn(t, m)
	var titles []string
	for _, tk := range col.tickets {
		titles = append(titles, tk.Title)
	}
	want := []string{"Add session test harness", "Filed yesterday", "Filed five days ago"}
	if strings.Join(titles, "|") != strings.Join(want, "|") {
		t.Fatalf("terminal lane = %q, want %q", titles, want)
	}
	if col.filed != 2 {
		t.Errorf("filed = %d, want 2", col.filed)
	}

	// The lane window follows the cursor; put it on the terminal lane so the
	// render has to draw it.
	m.laneIdx, m.cardIdx = idx, 0
	out := stripANSI(m.render())
	if !strings.Contains(out, "filed") {
		t.Errorf("a logbook card carries no textual marker:\n%s", out)
	}
	if strings.Contains(out, "Filed long ago") {
		t.Error("a ticket filed sixty days ago is on the board; the default window is 28 days")
	}

	// The card is filled in the logbook's own shade, not the lane's grey.
	logged := col.tickets[1]
	if block := m.renderCardBlock(logged, 30, false, false); !strings.Contains(block, "48;"+m.logbookShade(false)+"m") {
		t.Errorf("logbook card is not filled in its own shade: %q", block)
	}
	if block := m.renderCardBlock(col.tickets[0], 30, false, false); strings.Contains(block, "48;"+m.logbookShade(false)+"m") {
		t.Error("a card still on the board is filled like a logbook card")
	}
}

// logbook-days 0 takes the logbook off the board altogether.
func TestLogbookDaysZeroHidesTheCards(t *testing.T) {
	s := newTestStore(t)
	fileIntoLogbook(t, s, "Filed yesterday", 1)
	setLogbookDays(t, 0)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	if _, col := doneColumn(t, m); col.filed != 0 || len(col.tickets) != 1 {
		t.Errorf("logbook-days 0: terminal lane holds %d cards, %d filed; want only the board's own", len(col.tickets), col.filed)
	}
}

// A logbook card opens with enter and refuses everything that would change
// it, naming the way back; the file stays where it is.
func TestALogbookCardIsReadOnly(t *testing.T) {
	s := newTestStore(t)
	id := fileIntoLogbook(t, s, "Filed yesterday", 1)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 200, 40
	if !m.selectByID(id) {
		t.Fatal("the logbook card is not on the board")
	}

	m.key(key("m"))
	if m.mode == modeMove {
		t.Fatal("m on a logbook card opened the move picker")
	}
	if !strings.Contains(m.message, "jaira restore") {
		t.Errorf("refusal = %q, want it to name 'jaira restore'", m.message)
	}
	m.key(key("esc"))

	m.key(key("x"))
	if m.mode != modeMessage {
		t.Error("x on a logbook card did not answer with the refusal")
	}
	m.key(key("esc"))

	m.key(key("enter"))
	if m.mode != modeDetail || m.detail == nil || m.detail.ID != id {
		t.Fatalf("enter on a logbook card: mode %v, detail %v; want it open", m.mode, m.detail)
	}
	for _, k := range []string{"e", "E", "m", "X"} {
		m.key(key(k))
		if m.mode != modeMessage {
			t.Errorf("%s in an open logbook ticket: mode %v, want the refusal", k, m.mode)
		}
		m.key(key("esc"))
		if m.mode != modeDetail {
			t.Fatalf("esc after the refusal: mode %v, want the ticket still open", m.mode)
		}
	}

	if _, err := s.Load(id); err == nil {
		t.Error("the logbook ticket is back on the board")
	}
	if files, _ := filepath.Glob(filepath.Join(s.LogbookDir(), "*", "*.md")); len(files) != 1 {
		t.Errorf("logbook holds %d files, want the one it had", len(files))
	}
}

// The settings screen shows logbook-days, writes what is typed into
// settings.json, and the board shows the new window at once.
func TestSettingsEditsLogbookDays(t *testing.T) {
	s := newTestStore(t)
	fileIntoLogbook(t, s, "Filed five days ago", 5)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	if _, col := doneColumn(t, m); col.filed != 1 {
		t.Fatalf("setup: %d logbook cards, want 1", col.filed)
	}

	m.key(key("S"))
	if out := stripANSI(m.render()); !strings.Contains(out, "28 days") {
		t.Errorf("settings screen does not show the default window:\n%s", out)
	}
	m.key(key("j"))
	m.key(key("j"))
	m.key(key("enter"))
	if !m.settingsScreen.editing {
		t.Fatal("enter on the logbook row did not open the number")
	}
	m.key(key("backspace"))
	m.key(key("backspace"))
	m.key(key("3"))
	m.key(key("enter"))
	if m.mode != modeSettings || m.settingsScreen.editing {
		t.Fatalf("after saving: mode %v, editing %v", m.mode, m.settingsScreen.editing)
	}
	if got := settings.Load().LogbookWindow(); got != 3 {
		t.Errorf("settings.json logbook-days = %d, want 3", got)
	}
	if _, col := doneColumn(t, m); col.filed != 0 {
		t.Errorf("after logbook-days 3 the board still shows %d card(s) filed five days ago", col.filed)
	}
	b, err := os.ReadFile(settings.Path())
	if err != nil || !strings.Contains(string(b), `"logbook-days": 3`) {
		t.Errorf("settings.json = %q (%v), want logbook-days 3", b, err)
	}
}
