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

// registryWith is a tag registry holding the given name/colour pairs,
// independent of any store — Colour() only ever reads the in-memory map, so
// nothing needs to be saved to disk for a test.
func registryWith(t *testing.T, pairs ...any) *tag.Registry {
	t.Helper()
	reg, err := tag.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs)%2 != 0 {
		t.Fatalf("registryWith wants name/colour pairs, got %d values", len(pairs))
	}
	for i := 0; i < len(pairs); i += 2 {
		name, ok := pairs[i].(string)
		if !ok {
			t.Fatalf("registryWith: value %d is not a tag name", i)
		}
		colour, ok := pairs[i+1].(int)
		if !ok {
			t.Fatalf("registryWith: value %d is not a colour", i+1)
		}
		reg.Set(name, colour)
	}
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

func TestCardColorsAreTheFirstTwoTagsInTicketOrder(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45)
	tk := &ticket.Ticket{ID: "x", Title: "t", Tags: []string{"ui", "backend"}}

	slots := m.cardColors(tk)
	if !slots[0].coloured || slots[0].colour != 83 {
		t.Errorf("slot 1 = %+v; want colour 83", slots[0])
	}
	if !slots[1].coloured || slots[1].colour != 45 {
		t.Errorf("slot 2 = %+v; want colour 45", slots[1])
	}
}

// Slot 3 is reserved for the sprint marker and stays empty until that decision
// is taken — a third tag must not creep into it.
func TestThirdSlotStaysUncolouredHoweverManyTags(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45, "docs", 200, "ci", 111)
	tk := &ticket.Ticket{ID: "x", Title: "t", Tags: []string{"ui", "backend", "docs", "ci"}}

	if slots := m.cardColors(tk); slots[2].coloured {
		t.Errorf("slot 3 took a colour: %+v", slots[2])
	}
}

func TestCardColorsAreAbsentWithoutTagsOrWithoutARegistryEntry(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)

	untagged := &ticket.Ticket{ID: "a", Title: "t"}
	for i, s := range m.cardColors(untagged) {
		if s.coloured {
			t.Errorf("an untagged ticket coloured slot %d", i+1)
		}
	}
	uncoloured := &ticket.Ticket{ID: "b", Title: "t", Tags: []string{"backend"}}
	if m.cardColors(uncoloured)[0].coloured {
		t.Error("a tag with no registry entry produced a card colour")
	}
	// The uncoloured first tag does not shift the second one up a slot.
	mixed := &ticket.Ticket{ID: "c", Title: "t", Tags: []string{"backend", "ui"}}
	slots := m.cardColors(mixed)
	if slots[0].coloured {
		t.Error("an uncoloured first tag took slot 1")
	}
	if !slots[1].coloured || slots[1].colour != 83 {
		t.Errorf("slot 2 = %+v; want colour 83", slots[1])
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
//
// One tag now inks the first row alone: row 2 belongs to a second tag and row 3
// is reserved, so both fall back to the lane's shade.
func TestTaggedCardCarriesItsColourInTheFirstRowOnly(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: "Coloured", Tags: []string{"ui"}}

	raw := m.renderCardBlock(tk, 40, false, false)
	stripped := stripANSI(raw)
	_, shade := m.laneShade(false)

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
	rows := strings.Split(strings.TrimSuffix(raw, "\n"), "\n")
	if !strings.HasPrefix(rows[0], "\x1b[48;5;83m") {
		t.Errorf("row 0 does not open with the tag's bar cell:\n%q", rows[0])
	}
	for _, i := range []int{1, 2} {
		if !strings.HasPrefix(rows[i], "\x1b[48;"+shade+"m") {
			t.Errorf("row %d does not fall back to the lane shade:\n%q", i, rows[i])
		}
		if strings.Contains(rows[i], "48;5;83m") {
			t.Errorf("row %d borrowed the single tag's colour:\n%q", i, rows[i])
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
			m.tags = registryWith(t, "ui", 83, "backend", 45)
			idx := todoLaneIdx(t, m)
			for _, tk := range m.cols[idx].tickets {
				tk.Tags = []string{"ui", "backend"}
			}
			win := m.boardFit(m.width)
			raw := m.renderColumn(idx, win.colW, h)
			out := stripANSI(raw)
			shown := m.cardsInBudget(m.cols[idx].tickets, 0, max(1, h-4))

			// The bar cell carries the first tag's colour on the card's first
			// row and the second tag's on its second, so counting each colour
			// counts the cards whose first and second rows reached the screen —
			// and the third row is counted by the flag line only a whole card
			// draws. A card cut short fails one of the three.
			for _, c := range []string{"48;5;83m", "48;5;45m"} {
				if got := strings.Count(raw, c); got != shown {
					t.Errorf("w=%d h=%d: colour %s drawn %d times, want %d:\n%s",
						w, h, c, got, shown, out)
				}
			}
			if got := strings.Count(out, "○ spec"); got != shown {
				t.Errorf("w=%d h=%d: %d third card rows drawn, want %d:\n%s",
					w, h, got, shown, out)
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
	if s := m.cardColors(tk)[0]; !s.coloured || s.colour != 83 {
		t.Errorf("cardColors(UI)[0] = %+v, want colour 83", s)
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

// --- three colour slots down the bar ---------------------------------------

// barRows is the bar cell of each row of a rendered card: the SGR background
// parameters the row opens with, which is the one cell of colour the card
// carries before its text.
func barRows(t *testing.T, raw string) []string {
	t.Helper()
	var params []string
	for _, row := range strings.Split(strings.TrimSuffix(raw, "\n"), "\n") {
		rest, ok := strings.CutPrefix(row, "\x1b[48;")
		if !ok {
			t.Fatalf("row does not open with a background: %q", row)
		}
		i := strings.Index(rest, "m")
		if i < 0 {
			t.Fatalf("row's opening background never ends: %q", row)
		}
		params = append(params, rest[:i])
	}
	return params
}

// Two tags, two visible colours, in the order they stand on the ticket: this is
// the whole point of the change — filtering by two axes at once without
// switching the board's filter back and forth.
func TestTwoTaggedCardShowsBothColoursInTicketOrder(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45)
	tk := &ticket.Ticket{ID: "a", Title: "Two tags", Tags: []string{"ui", "backend"}}

	rows := barRows(t, m.renderCardBlock(tk, 40, false, false))
	if len(rows) != 3 {
		t.Fatalf("card is %d rows, want 3", len(rows))
	}
	if rows[0] != "5;83" {
		t.Errorf("row 1 bar = %q, want the first tag's 5;83", rows[0])
	}
	if rows[1] != "5;45" {
		t.Errorf("row 2 bar = %q, want the second tag's 5;45", rows[1])
	}
	if rows[0] == rows[1] {
		t.Errorf("both colour slots are %q, so the two tags are indistinguishable", rows[0])
	}

	// And the other way round on the ticket, the colours swap with them.
	swapped := &ticket.Ticket{ID: "b", Title: "Two tags", Tags: []string{"backend", "ui"}}
	rows = barRows(t, m.renderCardBlock(swapped, 40, false, false))
	if rows[0] != "5;45" || rows[1] != "5;83" {
		t.Errorf("swapped tags render %q/%q, want 5;45/5;83", rows[0], rows[1])
	}
}

// Slot 3 belongs to the sprint marker, which is not decided yet, so it shows
// the lane's shade on every card — including one carrying more tags than there
// are slots.
func TestThirdRowAlwaysCarriesTheLaneShade(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45, "docs", 200, "ci", 111)
	_, shade := m.laneShade(false)

	for _, tags := range [][]string{
		nil,
		{"ui"},
		{"ui", "backend"},
		{"ui", "backend", "docs"},
		{"ui", "backend", "docs", "ci"},
	} {
		tk := &ticket.Ticket{ID: "a", Title: "Card", Tags: tags}
		rows := barRows(t, m.renderCardBlock(tk, 40, false, false))
		if len(rows) != 3 {
			t.Fatalf("tags %v: card is %d rows, want 3", tags, len(rows))
		}
		if rows[2] != shade {
			t.Errorf("tags %v: row 3 bar = %q, want the lane shade %q", tags, rows[2], shade)
		}
	}
}

// A tag with no line in .jaira/tags keeps its slot rather than letting the next
// tag slide up into it, and the card's three text lines stand exactly where an
// untagged card's do — the bar is one cell wide whatever it is painted.
func TestUncolouredSecondTagFallsBackWithoutMovingTheText(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83) // knows "ui", not "backend"
	_, shade := m.laneShade(false)

	tagged := &ticket.Ticket{ID: "a", Title: "Mixed tags", Assignee: "berk", Tags: []string{"ui", "backend"}}
	rows := barRows(t, m.renderCardBlock(tagged, 40, false, false))
	if rows[0] != "5;83" {
		t.Errorf("row 1 bar = %q, want the coloured first tag's 5;83", rows[0])
	}
	if rows[1] != shade {
		t.Errorf("row 2 bar = %q, want the lane shade %q for the uncoloured tag", rows[1], shade)
	}

	plain := &ticket.Ticket{ID: "a", Title: "Mixed tags", Assignee: "berk"}
	want := stripANSI(m.renderCardBlock(plain, 40, false, false))
	got := stripANSI(m.renderCardBlock(tagged, 40, false, false))
	if got != want {
		t.Errorf("the tagged card's text sits differently from an untagged one:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// Four tags render, and the two past the slots simply colour nothing. The limit
// is a display limit: nothing rejects the fourth tag, so no existing ticket is
// made invalid by it.
func TestFourTaggedCardRendersWithTheExtraTagsUncoloured(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45, "docs", 200, "ci", 111)
	tk := &ticket.Ticket{ID: "a", Title: "Four tags", Tags: []string{"ui", "backend", "docs", "ci"}}

	raw := m.renderCardBlock(tk, 40, false, false)
	rows := barRows(t, raw)
	if len(rows) != 3 {
		t.Fatalf("a four-tag card is %d rows, want 3", len(rows))
	}
	for _, c := range []string{"5;200", "5;111"} {
		if strings.Contains(raw, "48;"+c+"m") {
			t.Errorf("a tag past the slots coloured the bar with %q:\n%q", c, raw)
		}
	}
	if !strings.Contains(stripANSI(raw), "Four tags") {
		t.Errorf("a four-tag card lost its title:\n%s", stripANSI(raw))
	}
}

// A ticket field may carry a newline: `jaira create` takes a multi-line title
// verbatim and writes it as a YAML block scalar, and the parser reads the break
// back. The card has three rows and three colour slots, so such a value must not
// add a row — it used to, and renderCardBlock then indexed past the slots and
// took the board down with it.
func TestACardWhoseFieldsCarryNewlinesStaysThreeRows(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83, "backend", 45)

	for _, c := range []struct {
		name string
		tk   *ticket.Ticket
	}{
		{"title", &ticket.Ticket{ID: "a", Title: "Zeile eins\nZeile zwei", Tags: []string{"ui", "backend"}}},
		{"assignee", &ticket.Ticket{ID: "b", Title: "Ein Titel", Assignee: "berk\nmehr", Tags: []string{"ui"}}},
		{"updated-by", &ticket.Ticket{ID: "c", Title: "Ein Titel", UpdatedBy: "berk\nmehr"}},
		{"executed-by", &ticket.Ticket{ID: "d", Title: "Ein Titel", ExecutedBy: "agent\nmehr"}},
		{"carriage return", &ticket.Ticket{ID: "e", Title: "Zeile eins\r\nZeile zwei"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			raw := m.renderCardBlock(c.tk, 40, false, false)
			if rows := barRows(t, raw); len(rows) != cardSlots {
				t.Fatalf("card is %d rows, want %d:\n%s", len(rows), cardSlots, stripANSI(raw))
			}
			if strings.Contains(stripANSI(raw), "\r") {
				t.Errorf("a carriage return survived into the card:\n%q", raw)
			}
		})
	}
}

// renderCardBlock draws one row per line renderCard returns and reads that row's
// colour out of cardSlots slots, so the two counts have to agree. This pins the
// agreement itself rather than one field that could break it: whatever a ticket
// carries, renderCard returns exactly cardHeight lines.
func TestRenderCardAlwaysReturnsAsManyLinesAsThereAreSlots(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)

	for _, tk := range []*ticket.Ticket{
		{ID: "a", Title: "Schlicht"},
		{ID: "b", Title: "Drei\nZeilen\nTitel", Tags: []string{"ui"}},
		{ID: "c", Title: "", Assignee: "berk\n\n\nberk"},
		{ID: "d", Title: strings.Repeat("sehr langer Titel ", 20)},
	} {
		content := strings.TrimSuffix(m.renderCard(tk, 39, false), "\n")
		if got := len(strings.Split(content, "\n")); got != cardSlots {
			t.Errorf("renderCard(%s) returned %d lines, want %d (cardHeight is %d):\n%s",
				tk.ID, got, cardSlots, m.cardHeight(tk), stripANSI(content))
		}
	}
}
