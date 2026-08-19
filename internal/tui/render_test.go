package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// newTestStore builds a store with a handful of tickets spread across lanes.
func newTestStore(t *testing.T) *ticket.Store {
	t.Helper()
	dir := t.TempDir()
	// Redirect BOTH the config and state locations before anything touches the
	// filesystem, so a test run never writes into the real user's home.
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))

	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	mk := func(title, status string, ready bool, blockedBy []string, commits []string) *ticket.Ticket {
		f := map[string]string{
			ticket.FieldID:        ticket.NewID(now),
			ticket.FieldTitle:     title,
			ticket.FieldStatus:    status,
			ticket.FieldReady:     boolStr(ready),
			ticket.FieldCreator:   "berk",
			ticket.FieldAssignee:  "berk",
			ticket.FieldCreatedAt: ticket.FormatTime(now),
			ticket.FieldUpdatedAt: ticket.FormatTime(now),
		}
		if ready {
			f[ticket.FieldGoal] = "make the thing work"
			f[ticket.FieldDoD] = "a test covers it"
			f[ticket.FieldContext] = "came up while debugging"
		}
		l := map[string][]string{ticket.FieldBlockedBy: blockedBy, ticket.FieldCommits: commits}
		tk, err := s.Create(f, l, "")
		if err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Millisecond)
		return tk
	}

	mk("Add session test harness", "done", true, nil, []string{"abc1234"})
	open := mk("Fix session cookie dropped on 302", "in-progress", true, nil, []string{"def5678"})
	// Blocked by work that is still in flight, so the dependency is unsatisfied.
	mk("Refactor auth middleware", "todo", true, []string{open.ID}, nil)
	mk("Investigate flaky logout test", "backlog", false, nil, nil)
	mk("Rate limit the login endpoint", "backlog", false, nil, nil)
	mk("Decide on cookie SameSite policy", "human", true, nil, nil)
	return s
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func newTestModel(t *testing.T, w, h int) *Model {
	t.Helper()
	s := newTestStore(t)
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = w, h
	return m
}

func TestBoardRenders(t *testing.T) {
	m := newTestModel(t, 150, 32)
	out := m.render()
	if out == "" {
		t.Fatal("empty render")
	}
	// Every built-in lane that holds a ticket must be visible.
	for _, want := range []string{"Backlog", "Todo", "Implementing", "HITL"} {
		if !strings.Contains(out, want) {
			t.Errorf("lane %q missing from board", want)
		}
	}
	if !strings.Contains(out, "Refactor auth") {
		t.Error("ticket title missing")
	}
	t.Logf("\n%s", out)
}

func TestBoardShowsStateWithGlyphsNotColourAlone(t *testing.T) {
	m := newTestModel(t, 150, 32)
	out := stripANSI(m.render())
	// An under-specified ticket and a blocked ticket must be distinguishable
	// with colour removed, or the board is unreadable to anyone who cannot see
	// the colours.
	if !strings.Contains(out, "spec") {
		t.Error("no textual marker for an under-specified ticket")
	}
	if !strings.Contains(out, "blocked") {
		t.Error("no textual marker for a blocked ticket")
	}
}

func TestFilterNarrowsTheBoard(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.filter = "auth"
	m.rebuild()
	out := stripANSI(m.render())
	if !strings.Contains(out, "Refactor auth") {
		t.Error("matching ticket was filtered out")
	}
	if strings.Contains(out, "Rate limit") {
		t.Error("non-matching ticket survived the filter")
	}
}

func TestNavigationStaysInBounds(t *testing.T) {
	m := newTestModel(t, 150, 32)
	for i := 0; i < 40; i++ {
		m.moveLane(1)
		m.moveCard(1)
	}
	for i := 0; i < 40; i++ {
		m.moveLane(-1)
		m.moveCard(-1)
	}
	if m.laneIdx < 0 || m.laneIdx >= len(m.cols) {
		t.Errorf("laneIdx %d out of range (%d lanes)", m.laneIdx, len(m.cols))
	}
	if m.render() == "" {
		t.Error("render broke after navigation")
	}
}

func TestRefreshPreservesSelection(t *testing.T) {
	m := newTestModel(t, 150, 32)
	// Land on a specific ticket, then reload as the background timer would.
	m.laneIdx, m.cardIdx = 0, 0
	before := m.selected()
	if before == nil {
		t.Fatal("nothing selected")
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	after := m.selected()
	if after == nil || after.ID != before.ID {
		t.Errorf("selection moved across a refresh: %v -> %v", idOf(before), idOf(after))
	}
}

func TestUnknownLaneBecomesReadOnlyColumn(t *testing.T) {
	m := newTestModel(t, 170, 32)
	// Simulate a teammate's custom lane this installation does not have.
	tk := m.tickets[0]
	if _, err := m.store.Mutate(tk.ID, func(x *ticket.Ticket) error {
		return x.Doc().SetScalar(ticket.FieldStatus, "critique")
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, c := range m.cols {
		if c.lane.ID == "critique" {
			found = true
			if !c.lane.Unknown {
				t.Error("column should be marked unknown")
			}
			if len(c.tickets) != 1 {
				t.Errorf("ticket was hidden: %d in column", len(c.tickets))
			}
		}
	}
	if !found {
		t.Fatal("no passthrough column was created; the ticket would be invisible")
	}
	// Navigate to it so the column actually renders.
	for i, c := range m.cols {
		if c.lane.ID == "critique" {
			m.laneIdx = i
		}
	}
	out := stripANSI(m.render())
	if !strings.Contains(out, "critique") {
		t.Errorf("passthrough column not rendered:\n%s", out)
	}
	if !strings.Contains(out, "read-only") {
		t.Error("passthrough column is not labelled read-only")
	}
	t.Logf("\n%s", out)
}

func TestGateRefusalIsSurfacedInTheUI(t *testing.T) {
	m := newTestModel(t, 150, 32)
	// Select an under-specified backlog ticket and try to move it forward.
	var target *ticket.Ticket
	for _, tk := range m.tickets {
		if tk.Status == "backlog" {
			target = tk
			break
		}
	}
	if target == nil {
		t.Fatal("no backlog ticket in fixture")
	}
	m.selectByID(target.ID)
	m.openMove()
	for i, l := range m.lanes.Lanes {
		if l.ID == "todo" {
			m.moveTarget = i
		}
	}
	m.applyMove()

	if m.mode != modeMessage || !m.isErr {
		t.Fatalf("expected a refusal message, got mode=%v isErr=%v", m.mode, m.isErr)
	}
	out := stripANSI(m.render())
	for _, want := range []string{"goal", "definition-of-done", "context"} {
		if !strings.Contains(out, want) {
			t.Errorf("refusal did not name the missing field %q:\n%s", want, out)
		}
	}
	// And the ticket must not have moved.
	reloaded, err := m.store.Load(target.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Status != "backlog" {
		t.Errorf("ticket moved despite the refusal: now %q", reloaded.Status)
	}
	t.Logf("\n%s", out)
}

func TestHelpRenders(t *testing.T) {
	m := newTestModel(t, 100, 40)
	m.mode = modeHelp
	out := stripANSI(m.render())
	for _, want := range []string{"Move around", "filter", "create a ticket", "Gates are enforced"} {
		if !strings.Contains(out, want) {
			t.Errorf("help missing %q", want)
		}
	}
}

func TestNarrowTerminalStillRenders(t *testing.T) {
	// A board with seven lanes cannot fit in 40 columns; it must scroll rather
	// than corrupt the layout.
	m := newTestModel(t, 40, 20)
	out := stripANSI(m.render())
	if out == "" {
		t.Fatal("empty render at 40 columns")
	}
	if strings.Contains(out, "off-screen") {
		t.Error("off-screen notice should be gone")
	}
	for _, line := range strings.Split(out, "\n") {
		if len([]rune(line)) > 60 {
			t.Errorf("line overflows a narrow terminal by a lot (%d cols): %q", len([]rune(line)), line)
			break
		}
	}
	t.Logf("\n%s", out)
}

func TestKeyRoutingDoesNotPanic(t *testing.T) {
	m := newTestModel(t, 120, 30)
	keys := []string{"j", "k", "h", "l", "g", "G", "enter", "esc", "/", "a", "backspace", "enter",
		"?", "esc", "n", "x", "esc", "m", "esc", "d", "esc", "r", "esc"}
	var model tea.Model = m
	for _, k := range keys {
		model, _ = model.Update(tea.KeyPressMsg{Code: rune(k[0]), Text: k})
	}
	if model == nil {
		t.Fatal("model became nil")
	}
}

// The detail pane is the only place a human (or another agent) can find the
// full id to hand to the CLI, so it has to print the whole thing, and y has
// to put it on the clipboard without closing the pane the id came from. The
// clipboard write itself (OSC52) goes straight to the terminal, so there is
// nothing here to read back — the test stops at "a command was returned".
func TestDetailShowsAndCopiesFullID(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := withChecklists(twoChecklists)
	m.detail = tk
	m.mode = modeDetail

	out := stripANSI(m.renderDetail())
	if !strings.Contains(out, tk.ID) {
		t.Fatalf("detail pane missing full id %q:\n%s", tk.ID, out)
	}

	_, cmd := m.key(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if cmd == nil {
		t.Fatal("y did not return a clipboard command")
	}
	if m.mode != modeDetail {
		t.Fatal("y closed the detail pane")
	}

	out = stripANSI(m.renderDetail())
	if !strings.Contains(out, "copied") {
		t.Fatalf("no copy confirmation after y:\n%s", out)
	}

	// "x" is unbound in the detail pane, so this is purely "any other key",
	// not a key that happens to also dismiss the pane.
	m.key(tea.KeyPressMsg{Code: 'x', Text: "x"})
	out = stripANSI(m.renderDetail())
	if strings.Contains(out, "copied") {
		t.Fatalf("copy confirmation still present after next keypress:\n%s", out)
	}
}

func idOf(t *ticket.Ticket) string {
	if t == nil {
		return "<nil>"
	}
	return ticket.Handle(t.ID)
}

// stripANSI removes escape sequences so assertions test content, not styling.
func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

var _ = os.Getenv

// A ticket captured in one session and picked up in another is only legible if
// the board says when each of those happened.
func TestDetailShowsWhenTheTicketAppearedAndWasTouched(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := withChecklists(twoChecklists)
	tk.CreatedAt = time.Now().Add(-72 * time.Hour)
	tk.UpdatedAt = time.Now().Add(-2 * time.Hour)
	m.detail = tk
	m.mode = modeDetail

	out := stripANSI(m.renderDetail())
	if !strings.Contains(out, "when") {
		t.Fatalf("detail pane has no when row:\n%s", out)
	}
	if !strings.Contains(out, "3d old") {
		t.Fatalf("when row does not say how old the ticket is:\n%s", out)
	}
	if !strings.Contains(out, "touched 2h ago") {
		t.Fatalf("when row does not say when it was last touched:\n%s", out)
	}
}

// A ticket written the moment it was created must not claim it was touched
// separately, or every card ever made carries a meaningless second timestamp.
func TestWhenRowOmitsTouchedForAFreshTicket(t *testing.T) {
	now := time.Now()
	if got := timespan(now, now); strings.Contains(got, "touched") {
		t.Fatalf("fresh ticket reported a separate touch: %q", got)
	}
	if got := timespan(time.Time{}, now); got != "" {
		t.Fatalf("a ticket with no creation time should render nothing, got %q", got)
	}
}
