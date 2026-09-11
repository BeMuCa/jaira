package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// registryWith is a tag registry holding one coloured entry, independent of
// any store — Colour() only ever reads the in-memory map, so nothing needs
// to be saved to disk for a test.
func registryWith(t *testing.T, name string, colour int) *tag.Registry {
	t.Helper()
	reg, err := tag.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg.Set(name, colour)
	return reg
}

// manyTodoTicketsModel builds a board with n plain tickets, all in the "todo"
// lane, so a column-level test can put more cards in one lane than a small
// terminal has room for.
func manyTodoTicketsModel(t *testing.T, n, w, h int) *Model {
	t.Helper()
	dir := t.TempDir()
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
	for i := 0; i < n; i++ {
		f := map[string]string{
			ticket.FieldID:        ticket.NewID(now),
			ticket.FieldTitle:     fmt.Sprintf("Ticket %d", i),
			ticket.FieldStatus:    "todo",
			ticket.FieldReady:     "false",
			ticket.FieldCreator:   "berk",
			ticket.FieldCreatedAt: ticket.FormatTime(now),
			ticket.FieldUpdatedAt: ticket.FormatTime(now),
		}
		if _, err := s.Create(f, nil, ""); err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Millisecond)
	}
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = w, h
	return m
}

func todoLaneIdx(t *testing.T, m *Model) int {
	t.Helper()
	for i, c := range m.cols {
		if c.lane.ID == "todo" {
			return i
		}
	}
	t.Fatal("no todo column on the board")
	return -1
}

// --- cardColor / cardHeight -------------------------------------------------

func TestCardColorIsTheFirstTagsRegistryColour(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "x", Title: "t", Tags: []string{"ui", "backend"}}

	c, ok := m.cardColor(tk)
	if !ok || c != 83 {
		t.Errorf("cardColor = %d, %v; want 83, true", c, ok)
	}
}

func TestCardColorIsAbsentWithoutTagsOrWithoutARegistryEntry(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)

	untagged := &ticket.Ticket{ID: "a", Title: "t"}
	if _, ok := m.cardColor(untagged); ok {
		t.Error("an untagged ticket has a card colour")
	}
	uncoloured := &ticket.Ticket{ID: "b", Title: "t", Tags: []string{"backend"}}
	if _, ok := m.cardColor(uncoloured); ok {
		t.Error("a tag with no registry entry produced a card colour")
	}
}

func TestCardHeightIsTheThreeContentRows(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)

	plain := &ticket.Ticket{ID: "a", Title: "t"}
	coloured := &ticket.Ticket{ID: "b", Title: "t", Tags: []string{"ui"}}
	uncoloured := &ticket.Ticket{ID: "c", Title: "t", Tags: []string{"backend"}}

	for _, tk := range []*ticket.Ticket{plain, coloured, uncoloured} {
		if h := m.cardHeight(tk); h != 3 {
			t.Errorf("cardHeight(%s) = %d, want 3 — a band has no border rows", tk.ID, h)
		}
	}
}

// --- cardsInBudget -----------------------------------------------------------

func TestCardsInBudgetCountsEachCardsOwnHeight(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tickets := []*ticket.Ticket{
		{ID: "a", Title: "t", Tags: []string{"ui"}}, // 3
		{ID: "b", Title: "t"},                       // 3 — a band too, just unmarked
		{ID: "c", Title: "t"},                       // 3
		{ID: "d", Title: "t", Tags: []string{"ui"}}, // 3
	}
	cases := []struct {
		budget int
		want   int
	}{
		{budget: 0, want: 1},   // a lone card always counts, however tight
		{budget: 3, want: 1},   // exactly the first card, no room for the next
		{budget: 5, want: 1},   // one row short of two
		{budget: 6, want: 2},   // 3 + 3
		{budget: 9, want: 3},   // 3 + 3 + 3
		{budget: 12, want: 4},  // every card fits
		{budget: 100, want: 4}, // more than enough
	}
	for _, c := range cases {
		if got := m.cardsInBudget(tickets, 0, c.budget); got != c.want {
			t.Errorf("cardsInBudget(budget=%d) = %d, want %d", c.budget, got, c.want)
		}
	}
}

func TestCardsInBudgetFromANonZeroStart(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tickets := []*ticket.Ticket{
		{ID: "a", Title: "t"},
		{ID: "b", Title: "t", Tags: []string{"ui"}}, // 5
		{ID: "c", Title: "t", Tags: []string{"ui"}}, // 5
	}
	if got := m.cardsInBudget(tickets, 1, 6); got != 2 {
		t.Errorf("cardsInBudget(start=1, budget=6) = %d, want 2 (3 + 3)", got)
	}
	if got := m.cardsInBudget(tickets, 3, 10); got != 0 {
		t.Errorf("cardsInBudget(start past the end) = %d, want 0", got)
	}
}

// --- renderCardBlock ---------------------------------------------------------

// No card is framed any more. A card with no colour to show simply has nothing
// in the bar cell — its own shade goes there, so its text still lines up with
// every other card in the lane.
func TestUntaggedCardCarriesNoColour(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: "Untagged ticket", Assignee: "berk"}

	raw := m.renderCardBlock(tk, 40, false, false)
	stripped := stripANSI(raw)
	if got := strings.Count(stripped, "\n"); got != 3 {
		t.Errorf("untagged card is %d lines, want 3:\n%s", got, stripped)
	}
	if strings.ContainsAny(stripped, "┌└│─") {
		t.Errorf("untagged card still draws a frame:\n%s", stripped)
	}
	if strings.Contains(raw, ";83m") {
		t.Errorf("untagged card borrowed a tag colour:\n%q", raw)
	}
}

// A tag with no registry colour is the same case: the tag exists (the legend
// tests cover that), but there is no colour to put in the bar.
func TestTaggedCardWithoutColourCarriesNoColour(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83) // registry knows "ui", not "backend"
	tk := &ticket.Ticket{ID: "a", Title: "Uncoloured tag", Tags: []string{"backend"}}

	raw := m.renderCardBlock(tk, 40, false, false)
	if strings.ContainsAny(stripANSI(raw), "┌└│─") {
		t.Error("a card whose tag has no registry colour is framed")
	}
	if strings.Contains(raw, ";83m") {
		t.Errorf("a colourless tag borrowed colour 83:\n%q", raw)
	}
}

// The tag's colour fills a whole cell as a background. As a border glyph it
// inked about half a cell and read as a differently-coloured frame rather than
// as a marker, which is why it moved.
func TestTaggedCardCarriesItsColourAsAFilledCell(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: "Coloured", Tags: []string{"ui"}}

	raw := m.renderCardBlock(tk, 40, false, false)
	stripped := stripANSI(raw)

	if got := strings.Count(stripped, "\n"); got != 3 {
		t.Errorf("card is %d lines, want 3:\n%s", got, stripped)
	}
	if !strings.Contains(stripped, "Coloured") {
		t.Errorf("card lost the title:\n%s", stripped)
	}
	// A background, not a foreground: 48, not 38.
	if !strings.Contains(raw, "48;5;83m") {
		t.Errorf("the tag's colour is not filling a cell:\n%q", raw)
	}
	// Every row of the card carries the bar, not only the first.
	for i, l := range strings.Split(strings.TrimSuffix(raw, "\n"), "\n") {
		if !strings.Contains(l, "48;5;83m") {
			t.Errorf("row %d has no bar cell:\n%q", i, l)
		}
	}
}

func TestCardBandFitsTheWidthItIsGiven(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: strings.Repeat("a very long title ", 5), Tags: []string{"ui"}}

	for _, w := range []int{10, 18, 24, 40, 80} {
		out := stripANSI(m.renderCardBlock(tk, w, false, false))
		for _, l := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
			checkLineWidths(t, w, fmt.Sprintf("card band w=%d", w), l)
		}
	}
}

// --- column budgeting: no card is ever cut short ---------------------------

// A band has no border rows to count, so what pins the budget now is that every
// card the column claims to show is drawn whole: exactly cardHeight rows of it,
// no card ending halfway down the lane. A card cut short is worse here than it
// was with frames — with nothing outlining it, half a card reads as a whole
// one, and the lane silently lies about what it holds.
func TestColumnDrawsEveryCardItCountsInFull(t *testing.T) {
	for _, w := range []int{40, 80} {
		for _, h := range []int{8, 9, 10, 11, 12, 14, 18, 24, 32} {
			m := manyTodoTicketsModel(t, 8, w, h)
			m.tags = registryWith(t, "ui", 83)
			idx := todoLaneIdx(t, m)
			for _, tk := range m.cols[idx].tickets {
				tk.Tags = []string{"ui"}
			}
			win := m.boardFit(m.width)
			raw := m.renderColumn(idx, win.colW, h)
			out := stripANSI(raw)
			shown := m.cardsInBudget(m.cols[idx].tickets, 0, max(1, h-4))

			// The bar cell is drawn once per card row, so counting it counts
			// the rows that actually reached the screen.
			if got := strings.Count(raw, "48;5;83m"); got != shown*3 {
				t.Errorf("w=%d h=%d: %d card rows drawn, want %d (%d cards × 3):\n%s",
					w, h, got, shown*3, shown, out)
			}
			for _, l := range strings.Split(out, "\n") {
				if got := len([]rune(l)); got > win.colW+2 {
					t.Errorf("w=%d h=%d: column line %d wide, want at most %d: %q", w, h, got, win.colW+2, l)
				}
			}
		}
	}
}

// --- t / the legend -----------------------------------------------------------

func TestTOpensAndClosesTheLegend(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.mode = modeBoard

	press(m, "t")
	if m.mode != modeLegend {
		t.Fatalf("t did not open the legend, mode = %v", m.mode)
	}
	press(m, "t")
	if m.mode != modeBoard {
		t.Errorf("t did not close the legend, mode = %v", m.mode)
	}

	press(m, "t")
	press(m, "esc")
	if m.mode != modeBoard {
		t.Errorf("esc did not close the legend, mode = %v", m.mode)
	}
}

func TestLegendListsSwatchAndNameOfActiveTags(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	if len(m.tickets) < 2 {
		t.Fatal("test store needs at least two tickets")
	}
	m.tickets[0].Tags = []string{"ui"}
	m.tickets[1].Tags = []string{"backend"} // no registry colour

	m.mode = modeLegend
	out := stripANSI(m.render())
	if !strings.Contains(out, "ui") {
		t.Errorf("legend missing the coloured tag's name:\n%s", out)
	}
	if !strings.Contains(out, "backend") {
		t.Errorf("legend missing the colourless tag's name:\n%s", out)
	}
	raw := m.render()
	if !strings.Contains(raw, "38;5;83") {
		t.Errorf("legend does not show ui's swatch in its registry colour:\n%q", raw)
	}
}

// --- key line -----------------------------------------------------------------

func TestBoardKeyLineNamesT(t *testing.T) {
	m := newTestModel(t, 150, 32)
	out := stripANSI(m.render())
	if !strings.Contains(out, "t tags") {
		t.Errorf("board key line does not name t:\n%s", out)
	}
}

// A hand-written "UI" works as "ui" in the filters and the count, so it must
// wear ui's colour on the board and be one legend line, not a colourless twin.
func TestHandWrittenCaseWearsTheRegistryColour(t *testing.T) {
	m := newTestModel(t, 120, 40)
	reg, err := tag.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg.Set("ui", 83)
	m.tags = reg

	tk := m.tickets[0]
	tk.Tags = []string{"UI"}
	if c, ok := m.cardColor(tk); !ok || c != 83 {
		t.Errorf("cardColor(UI) = %d,%v, want 83,true", c, ok)
	}
	m.tickets[1].Tags = []string{"ui"}
	tags := m.activeTags()
	n := 0
	for _, name := range tags {
		if name == "ui" {
			n++
		}
	}
	if n != 1 || len(tags) != 1 {
		t.Errorf("activeTags = %v, want exactly [ui]", tags)
	}
}

// The wrap class the review caught: a card heavy with flags must still render
// exactly its three rows — a line wider than the band wraps, the card outgrows
// cardHeight, and the column hides cards behind an honest-looking count. Pins
// the reviewer's probe as a permanent test.
func TestACardHeavyWithFlagsStaysThreeRows(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{
		ID: "a", Title: "A realistically long ticket title that will truncate",
		Tags: []string{"ui"}, Assignee: "someone-with-a-name", UpdatedBy: "someone-else",
		ExecutedBy: "opus", Commits: []string{"a1", "b2"},
		PlanItems: []ticket.DoDItem{{Text: "a", State: ticket.StateDone}, {Text: "b"}},
		DoDItems:  []ticket.DoDItem{{Text: "c"}},
	}
	for _, w := range []int{12, 18, 24, 40} {
		out := stripANSI(m.renderCardBlock(tk, w, false, false))
		lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
		if len(lines) != 3 {
			t.Errorf("w=%d: card renders %d rows, want 3:\n%s", w, len(lines), out)
		}
		for _, l := range lines {
			checkLineWidths(t, w, fmt.Sprintf("flag-heavy card w=%d", w), l)
		}
	}
}

// Two stacked cards are told apart by their shades alternating and by nothing
// else, so two neighbours must never land on the same one — there is no frame
// left to fall back on if they do.
func TestStackedCardsAlternateTheirShade(t *testing.T) {
	m := newTestModel(t, 150, 40) // the fixture's backlog holds two cards

	_, a := m.laneShade(false)
	_, b := m.laneShade(true)
	if a == b {
		t.Fatalf("both lane shades are %q, so stacked cards merge", a)
	}

	// And the selection sits above both, or the cursor is lost on whichever
	// shade happens to match it.
	_, sel := m.selectionFill(0, false)
	if sel == a || sel == b {
		t.Errorf("the neutral selection fill %q is also a lane shade", sel)
	}
}
