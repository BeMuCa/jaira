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

func TestCardHeightGrowsByTwoRowsOnlyWhenColoured(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)

	plain := &ticket.Ticket{ID: "a", Title: "t"}
	coloured := &ticket.Ticket{ID: "b", Title: "t", Tags: []string{"ui"}}
	uncoloured := &ticket.Ticket{ID: "c", Title: "t", Tags: []string{"backend"}}

	if h := m.cardHeight(plain); h != 3 {
		t.Errorf("untagged cardHeight = %d, want 3", h)
	}
	if h := m.cardHeight(uncoloured); h != 3 {
		t.Errorf("colourless-tag cardHeight = %d, want 3", h)
	}
	if h := m.cardHeight(coloured); h != 5 {
		t.Errorf("coloured-tag cardHeight = %d, want 5 (3 content + 2 border)", h)
	}
}

// --- cardsInBudget -----------------------------------------------------------

func TestCardsInBudgetCountsEachCardsOwnHeight(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tickets := []*ticket.Ticket{
		{ID: "a", Title: "t", Tags: []string{"ui"}}, // 5
		{ID: "b", Title: "t"},                       // 3
		{ID: "c", Title: "t"},                       // 3
		{ID: "d", Title: "t", Tags: []string{"ui"}}, // 5
	}
	cases := []struct {
		budget int
		want   int
	}{
		{budget: 0, want: 1},   // a lone card always counts, however tight
		{budget: 5, want: 1},   // exactly the first (5-row) card, no room for the next
		{budget: 8, want: 2},   // 5 + 3
		{budget: 11, want: 3},  // 5 + 3 + 3
		{budget: 16, want: 4},  // every card fits
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
	if got := m.cardsInBudget(tickets, 1, 10); got != 2 {
		t.Errorf("cardsInBudget(start=1, budget=10) = %d, want 2", got)
	}
	if got := m.cardsInBudget(tickets, 3, 10); got != 0 {
		t.Errorf("cardsInBudget(start past the end) = %d, want 0", got)
	}
}

// --- renderCardBlock ---------------------------------------------------------

// Untagged rendering must not move: renderCardBlock has to be byte-identical
// to what renderCard alone produced before this feature existed.
func TestUntaggedCardRenderingIsUnchanged(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: "Untagged ticket", Assignee: "berk"}

	plain := m.renderCard(tk, 40, false)
	boxed := m.renderCardBlock(tk, 40, false)
	if boxed != plain {
		t.Errorf("renderCardBlock changed an untagged card:\nplain: %q\nboxed: %q", plain, boxed)
	}
	if strings.Count(plain, "\n") != 3 {
		t.Errorf("an untagged card is not three lines: %q", plain)
	}
}

// A tag with no registry colour costs nothing on the card: same rendering as
// no tag at all, and the tag still exists (checked separately in the legend
// tests) — it just draws no box.
func TestTaggedCardWithoutColourIsUnboxed(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83) // registry knows "ui", not "backend"
	tk := &ticket.Ticket{ID: "a", Title: "Uncoloured tag", Tags: []string{"backend"}}

	plain := m.renderCard(tk, 40, false)
	boxed := m.renderCardBlock(tk, 40, false)
	if boxed != plain {
		t.Error("a card whose tag has no registry colour was still boxed")
	}
}

func TestTaggedCardWithColourIsBoxedInThatColour(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: "Coloured", Tags: []string{"ui"}}

	raw := m.renderCardBlock(tk, 40, false)
	stripped := stripANSI(raw)

	if got := strings.Count(stripped, "\n"); got != 5 {
		t.Errorf("boxed card is %d lines (incl. trailing newline), want 5:\n%s", got, stripped)
	}
	if !strings.Contains(stripped, "┌") || !strings.Contains(stripped, "└") {
		t.Errorf("boxed card has no border:\n%s", stripped)
	}
	if !strings.Contains(stripped, "Coloured") {
		t.Errorf("boxed card lost the title:\n%s", stripped)
	}
	// The border is drawn in the registry's colour, not left to the default.
	if !strings.Contains(raw, "38;5;83") {
		t.Errorf("boxed card border is not in colour 83:\n%q", raw)
	}
}

func TestTaggedCardBoxFitsTheWidthItIsGiven(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.tags = registryWith(t, "ui", 83)
	tk := &ticket.Ticket{ID: "a", Title: strings.Repeat("a very long title ", 5), Tags: []string{"ui"}}

	for _, w := range []int{10, 18, 24, 40, 80} {
		out := stripANSI(m.renderCardBlock(tk, w, false))
		for _, l := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
			checkLineWidths(t, w, fmt.Sprintf("tagged card w=%d", w), l)
		}
	}
}

// --- column budgeting: no card is ever cut mid-border -----------------------

// borderGlyphsBalance reports whether every opened border in s is closed —
// the invariant that must hold however renderColumn's height budget is
// computed. clampBlock cuts a column's body at a hard line count; if the
// budget under-counts a tall card's height, clampBlock lands inside that
// card's box instead of before it, and an opening corner survives with no
// matching closing one. This is the assertion that actually fails when the
// budgeting regresses to a fixed rows-per-card constant — unlike a bare
// width/height bound, which the outer clamp satisfies by construction no
// matter what renderColumn does internally.
func borderGlyphsBalance(t *testing.T, s string) {
	t.Helper()
	if strings.Count(s, "┌") != strings.Count(s, "└") {
		t.Errorf("a box was opened without a matching close:\n%s", s)
	}
	if strings.Count(s, "┐") != strings.Count(s, "┘") {
		t.Errorf("a box's right side was opened without a matching close:\n%s", s)
	}
}

func TestColumnNeverCutsATaggedCardInHalf(t *testing.T) {
	for _, w := range []int{40, 80} {
		for _, h := range []int{8, 9, 10, 11, 12, 14, 18, 24, 32} {
			m := manyTodoTicketsModel(t, 8, w, h)
			m.tags = registryWith(t, "ui", 83)
			idx := todoLaneIdx(t, m)
			for _, tk := range m.cols[idx].tickets {
				tk.Tags = []string{"ui"}
			}
			win := m.boardFit(m.width)
			out := stripANSI(m.renderColumn(idx, win.colW, h))
			borderGlyphsBalance(t, out)
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
