package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/ticket"
)

// ticketCount is the store's own count, so a test can prove nothing was written.
func ticketCount(t *testing.T, m *Model) int {
	t.Helper()
	ts, err := m.store.List()
	if err != nil {
		t.Fatal(err)
	}
	return len(ts)
}

// openATicket puts a ticket on screen from the board, the ordinary way in.
func openATicket(t *testing.T, m *Model) *ticket.Ticket {
	t.Helper()
	m.mode = modeBoard
	src := focusTicket(t, m, "in-progress", "Fix session cookie dropped on 302")
	m.openDetail()
	if m.detail == nil {
		t.Fatal("detail did not open")
	}
	return src
}

func TestFollowUpOpensBesideItsTicketWithoutWriting(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	before := ticketCount(t, m)

	press(m, "n")

	if m.follow == nil {
		t.Fatal("n did not open the split")
	}
	if m.follow.src == nil || m.follow.src.ID != src.ID {
		t.Errorf("the wrong ticket is on the left: %v", m.follow.src)
	}
	if m.follow.draft == nil {
		t.Fatal("no draft to write into")
	}
	if m.mode != modeEdit {
		t.Errorf("mode = %v, want modeEdit — the follow-up opens ready to write", m.mode)
	}
	if got := m.follow.draft.fields[ticket.FieldTitle]; got != "Follow-up: "+src.Title {
		t.Errorf("draft title = %q", got)
	}
	if got := m.follow.draft.fields[ticket.FieldFollows]; got != src.ID {
		t.Errorf("draft follows = %q, want %s", got, src.ID)
	}
	if got := m.follow.draft.fields[ticket.FieldContext]; !strings.Contains(got, ticket.Handle(src.ID)) {
		t.Errorf("draft context does not name the ticket it follows:\n%s", got)
	}
	if after := ticketCount(t, m); after != before {
		t.Errorf("a ticket was written before it was saved: %d -> %d", before, after)
	}
}

// Both tickets have to be on screen — seeing the one the follow-up is for while
// writing it is the entire reason for the split.
func TestSplitShowsBothTickets(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	press(m, "n")

	out := stripANSI(m.render())

	if !strings.Contains(out, ticket.Handle(src.ID)) {
		t.Errorf("the ticket being followed is not on screen:\n%s", out)
	}
	if !strings.Contains(out, "goal") || !strings.Contains(out, "context") {
		t.Errorf("the editor is not on screen:\n%s", out)
	}
	if !strings.Contains(out, "ctrl+s save") {
		t.Errorf("the footer does not say how to save:\n%s", out)
	}
}

// Typing goes into the draft, not into a file.
func TestTypingADraftTouchesNothingOnDisk(t *testing.T) {
	m := newTestModel(t, 140, 40)
	openATicket(t, m)
	before := ticketCount(t, m)

	press(m, "n")
	for _, r := range "make it fast" {
		m.key(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	press(m, "tab")

	if got := m.follow.draft.fields[ticket.FieldGoal]; got != "make it fast" {
		t.Errorf("draft goal = %q", got)
	}
	if after := ticketCount(t, m); after != before {
		t.Errorf("typing wrote to disk: %d -> %d", before, after)
	}
}

func TestDiscardedFollowUpLeavesNoTrace(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	before := ticketCount(t, m)

	press(m, "n")
	for _, r := range "abandoned" {
		m.key(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	press(m, "esc")

	if after := ticketCount(t, m); after != before {
		t.Errorf("a discarded follow-up was written anyway: %d -> %d", before, after)
	}
	if m.follow != nil {
		t.Error("the split is still open after discarding")
	}
	if m.mode != modeDetail || m.detail == nil || m.detail.ID != src.ID {
		t.Errorf("discarding did not come back to the ticket: mode = %v, detail = %v", m.mode, m.detail)
	}
}

func TestSavedFollowUpLandsInTheDefaultLaneLinkedBack(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	before := ticketCount(t, m)

	press(m, "n")
	for _, r := range "stop dropping it" {
		m.key(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})

	if after := ticketCount(t, m); after != before+1 {
		t.Fatalf("expected exactly one new ticket: %d -> %d", before, after)
	}
	if m.mode != modeDetail {
		t.Errorf("mode = %v, want modeDetail — a saved follow-up reads as a ticket", m.mode)
	}
	if m.follow == nil || m.follow.draft != nil {
		t.Fatalf("the draft is still open: %v", m.follow)
	}
	if m.follow.src == nil || m.follow.src.ID != src.ID {
		t.Error("the ticket it follows left the left pane")
	}
	if m.detail == nil {
		t.Fatal("the saved follow-up is not the open ticket")
	}
	if m.detail.Follows != src.ID {
		t.Errorf("follows = %q, want %s", m.detail.Follows, src.ID)
	}
	if m.detail.Status != m.lanes.Default().ID {
		t.Errorf("landed in %q, want the default lane %q", m.detail.Status, m.lanes.Default().ID)
	}
	if m.detail.Goal != "stop dropping it" {
		t.Errorf("goal = %q, want what was typed", m.detail.Goal)
	}
	// ready is derived, never typed: whatever the gate says about the saved file is
	// what the file must claim.
	full, err := m.store.Load(m.detail.ID)
	if err != nil {
		t.Fatal(err)
	}
	if full.Ready != gate.Ready(full) {
		t.Errorf("ready = %v, but the gate derives %v", full.Ready, gate.Ready(full))
	}
}

// A saved follow-up reads as a ticket, not as a form.
func TestSavedFollowUpRendersAsATicket(t *testing.T) {
	m := newTestModel(t, 140, 40)
	openATicket(t, m)
	press(m, "n")
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})

	out := stripANSI(m.render())

	if strings.Contains(out, "ctrl+s save") {
		t.Errorf("the editor footer is still up after saving:\n%s", out)
	}
	if !strings.Contains(out, "tab other pane") || !strings.Contains(out, "esc close") {
		t.Errorf("the footer does not name the split's keys:\n%s", out)
	}
	if !strings.Contains(out, "n follow-up") {
		t.Errorf("a saved follow-up cannot be followed up:\n%s", out)
	}
}

// Chaining: the follow-up you just wrote becomes the ticket the next one follows.
func TestFollowUpChainsAndSlidesLeft(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	press(m, "n")
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	first := m.detail.ID

	press(m, "n")

	if m.follow == nil || m.follow.src == nil {
		t.Fatal("the chain did not open a second split")
	}
	if m.follow.src.ID != first {
		t.Errorf("left pane is %s, want the follow-up just written (%s)", m.follow.src.ID, first)
	}
	if m.follow.src.ID == src.ID {
		t.Error("the original ticket did not slide off")
	}
	if got := m.follow.draft.fields[ticket.FieldFollows]; got != first {
		t.Errorf("the chained draft follows %q, want %s", got, first)
	}
}

func TestClosingASavedFollowUpComesBackToItsTicket(t *testing.T) {
	m := newTestModel(t, 140, 40)
	src := openATicket(t, m)
	press(m, "n")
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})

	press(m, "esc")

	if m.follow != nil {
		t.Error("the split is still open")
	}
	if m.mode != modeDetail || m.detail == nil || m.detail.ID != src.ID {
		t.Errorf("esc did not come back to the ticket: mode = %v, detail = %v", m.mode, m.detail)
	}
	// And one more esc leaves the ticket, as it always did.
	press(m, "esc")
	if m.mode != modeBoard {
		t.Errorf("a second esc landed on %v, want the board", m.mode)
	}
}

func TestTabMovesFocusBetweenPanesAndArrowsFollowIt(t *testing.T) {
	m := newTestModel(t, 140, 40)
	openATicket(t, m)
	press(m, "n")
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})

	m.key(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.detailScroll == 0 {
		t.Error("the arrows did not scroll the follow-up")
	}
	right := m.detailScroll

	press(m, "tab")
	if !m.follow.focusLeft {
		t.Fatal("tab did not move focus to the left pane")
	}
	m.key(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.follow.srcScroll == 0 {
		t.Error("with the left pane focused the arrows did not scroll it")
	}
	if m.detailScroll != right {
		t.Error("the arrows scrolled both panes at once")
	}
}

// While the editor holds tab, the ticket beside it still has to be readable.
func TestShiftArrowsScrollTheLeftPaneWhileWriting(t *testing.T) {
	m := newTestModel(t, 140, 40)
	openATicket(t, m)
	press(m, "n")

	m.key(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModShift})

	if m.follow.srcScroll != 1 {
		t.Errorf("srcScroll = %d, want 1", m.follow.srcScroll)
	}
	if m.mode != modeEdit {
		t.Errorf("shift+down left the editor: mode = %v", m.mode)
	}
}

// Two unreadable columns are worse than one readable one.
func TestNarrowTerminalDropsTheSplit(t *testing.T) {
	m := newTestModel(t, 60, 30)
	openATicket(t, m)
	press(m, "n")

	out := stripANSI(m.render())

	// The narrow path is the full-screen editor, exactly as it renders unsplit.
	if want := stripANSI(m.renderEdit()); out != want {
		t.Errorf("a 60-column terminal did not fall back to the plain editor:\n%s", out)
	}
	if strings.Contains(out, "╮ ╭") {
		t.Errorf("two panes were drawn into 60 columns:\n%s", out)
	}
	if !strings.Contains(out, "ctrl+s save") {
		t.Errorf("the follow-up editor is not on screen:\n%s", out)
	}
	for _, l := range strings.Split(out, "\n") {
		if len([]rune(l)) > 60 {
			t.Errorf("line wider than the terminal (%d):\n%s", len([]rune(l)), l)
		}
	}
}

// A wide but short terminal has no rows to spare for a second pane either.
func TestShortTerminalDropsTheSplit(t *testing.T) {
	m := newTestModel(t, 140, 12)
	openATicket(t, m)
	press(m, "n")

	if out := stripANSI(m.render()); strings.Contains(out, "╮ ╭") {
		t.Errorf("two panes were drawn into 12 rows:\n%s", out)
	}
}

// Walking to the neighbouring ticket leaves the whole follow-up screen behind.
func TestJumpingToTheNextTicketClosesTheSplit(t *testing.T) {
	m := newTestModel(t, 140, 40)
	openATicket(t, m)
	press(m, "n")
	m.key(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})

	press(m, "j")

	if m.follow != nil {
		t.Error("the split survived j")
	}
	if m.detail != nil {
		t.Error("the ticket is still open after j")
	}
}
