package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// mkMilestone writes a milestone file holding the given tickets and returns it.
func mkMilestone(t *testing.T, m *Model, name string, colour int, ids ...string) {
	t.Helper()
	ms := milestone.New(name, colour, nowish())
	for _, id := range ids {
		ms.Add(id)
	}
	if err := ms.Save(m.store.Root); err != nil {
		t.Fatal(err)
	}
}

func nowish() time.Time { return time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC) }

// A card must not change width, and its text must not change width, depending
// on whether the ticket is in a milestone — otherwise every title in a lane
// shifts by a column as tickets join and leave a group, which is exactly the
// flutter the reserved right cell exists to prevent.
func TestRightBarWidthIsTheSameWithAndWithoutMilestones(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := m.tickets[0]

	bare := m.renderCardBlock(tk, 30, false, false)
	bareRows := rowWidths(bare)

	mkMilestone(t, m, "round-one", 45, tk.ID)
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	tk = byID(t, m, tk.ID)
	one := m.renderCardBlock(tk, 30, false, false)
	if got := rowWidths(one); !sameInts(got, bareRows) {
		t.Errorf("one milestone changed the card: rows %v, want %v", got, bareRows)
	}

	// Four milestones: three slots are coloured, the fourth colours nothing,
	// and the card is still the same size.
	for i, name := range []string{"round-two", "round-three", "round-four"} {
		mkMilestone(t, m, name, milestone.Palette[i+1], tk.ID)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	tk = byID(t, m, tk.ID)
	four := m.renderCardBlock(tk, 30, false, false)
	if got := rowWidths(four); !sameInts(got, bareRows) {
		t.Errorf("four milestones changed the card: rows %v, want %v", got, bareRows)
	}

	slots := m.milestoneColors(tk)
	coloured := 0
	for _, s := range slots {
		if s.coloured {
			coloured++
		}
	}
	if coloured != cardSlots {
		t.Errorf("four milestones coloured %d slots, want all %d and no more", coloured, cardSlots)
	}
}

// The milestone colour belongs on the right edge, separated from the tag
// colours on the left: that separation is the whole reason two edges are used
// rather than one longer bar.
func TestMilestoneColourPaintsTheRightEdge(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := m.tickets[0]
	mkMilestone(t, m, "round-one", 45, tk.ID)
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	tk = byID(t, m, tk.ID)

	row := strings.Split(strings.TrimSuffix(m.renderCardBlock(tk, 30, false, false), "\n"), "\n")[0]
	// Three painted spans per row: left bar, text, right bar. The last one
	// carries the milestone's colour.
	spans := strings.Split(row, "\x1b[48;")
	if len(spans) != 4 {
		t.Fatalf("card row has %d painted spans, want 3 (left bar, text, right bar)", len(spans)-1)
	}
	if !strings.HasPrefix(spans[3], "5;45m") {
		t.Errorf("right-edge span = %q, want it opened in the milestone's colour 45", spans[3][:min(12, len(spans[3]))])
	}
	// And a ticket in no milestone leaves that cell in the card's own shade
	// rather than inking it.
	other := byID(t, m, m.tickets[1].ID)
	if other.ID == tk.ID {
		t.Skip("need a second ticket")
	}
	otherRow := strings.Split(strings.TrimSuffix(m.renderCardBlock(other, 30, false, false), "\n"), "\n")[0]
	otherSpans := strings.Split(otherRow, "\x1b[48;")
	if otherSpans[3] == spans[3] {
		t.Error("a ticket in no milestone painted the same right edge as one in a milestone")
	}
}

// The board's gesture writes milestone:<name> into the ordinary filter, so the
// filter has to answer that key exactly — and answer nothing to a name no
// milestone carries.
func TestFilterNarrowsToAMilestone(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := m.tickets[0]
	mkMilestone(t, m, "round-one", 45, tk.ID)
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.filter = "milestone:round-one"
	m.rebuild()
	var seen []string
	for _, c := range m.cols {
		for _, got := range c.tickets {
			seen = append(seen, got.ID)
		}
	}
	if len(seen) != 1 || seen[0] != tk.ID {
		t.Errorf("filter milestone:round-one showed %v, want only %s", seen, tk.ID)
	}
	m.filter = "milestone:round-nine"
	m.rebuild()
	for _, c := range m.cols {
		if len(c.tickets) != 0 {
			t.Errorf("filter on an unknown milestone still showed %d cards", len(c.tickets))
		}
	}
	if len(m.milestones) != 1 || m.milestones[0].Name != "round-one" {
		t.Errorf("the picker offers %d milestones, want only round-one", len(m.milestones))
	}
}

// The whole point of a milestone is pulling the board down to it in one
// gesture, so the picker has to do that from the keyboard and let go of it
// again without the user having to remember that esc on the board clears a
// filter.
func TestPickerNarrowsTheBoardAndReleasesIt(t *testing.T) {
	m := newTestModel(t, 150, 32)
	tk := m.tickets[0]
	mkMilestone(t, m, "round-one", 45, tk.ID)
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}

	m.Update(key("M"))
	if m.mode != modeMilestones {
		t.Fatalf("M did not open the picker; mode = %v", m.mode)
	}
	if !strings.Contains(stripANSI(m.render()), "round-one") {
		t.Error("the picker does not name the board's milestone")
	}
	m.Update(key("enter"))
	if m.mode != modeBoard {
		t.Errorf("enter left the picker open; mode = %v", m.mode)
	}
	if m.filter != "milestone:round-one" {
		t.Errorf("filter = %q, want milestone:round-one", m.filter)
	}
	shown := 0
	for _, c := range m.cols {
		shown += len(c.tickets)
	}
	if shown != 1 {
		t.Errorf("board shows %d cards after the gesture, want only the one in the milestone", shown)
	}

	m.Update(key("M"))
	m.Update(key("x"))
	if m.filter != "" {
		t.Errorf("x in the picker left the filter at %q", m.filter)
	}
	shown = 0
	for _, c := range m.cols {
		shown += len(c.tickets)
	}
	if shown != len(m.tickets) {
		t.Errorf("board shows %d of %d cards after releasing the filter", shown, len(m.tickets))
	}
}

// A board with no milestones must say so rather than opening an empty picker
// that looks broken.
func TestPickerOnABoardWithNoMilestonesExplainsItself(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.Update(key("M"))
	if m.mode == modeMilestones {
		t.Fatal("M opened an empty picker instead of explaining")
	}
	if !strings.Contains(stripANSI(m.render()), "jaira milestone create") {
		t.Error("the message does not say how to make one")
	}
}

func byID(t *testing.T, m *Model, id string) *ticket.Ticket {
	t.Helper()
	for _, tk := range m.tickets {
		if tk.ID == id {
			return tk
		}
	}
	t.Fatalf("ticket %s gone after reload", id)
	return nil
}

// rowWidths is the visible width of each row of a rendered card, colour
// stripped — the thing that must not move.
func rowWidths(block string) []int {
	var out []int
	for _, line := range strings.Split(strings.TrimSuffix(block, "\n"), "\n") {
		out = append(out, len([]rune(stripANSI(line))))
	}
	return out
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
