package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// linksModel builds a board whose selected card has one of every link, with
// one of them already filed into the logbook — the case the window exists
// for.
func linksModel(t *testing.T) (*Model, string) {
	t.Helper()
	m := newTestModel(t, 150, 40)

	now := time.Now()
	mk := func(title, status string, fields map[string]string, lists map[string][]string) string {
		f := map[string]string{
			ticket.FieldID:     ticket.NewID(now),
			ticket.FieldTitle:  title,
			ticket.FieldStatus: status,
		}
		for k, v := range fields {
			f[k] = v
		}
		tk, err := m.store.Create(f, lists, "")
		if err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Millisecond)
		return tk.ID
	}

	epic := mk("Epic that holds things", "todo", nil, nil)
	child := mk("A child still in play", "todo", map[string]string{ticket.FieldParent: epic}, nil)
	shipped := mk("A child already shipped", "done", map[string]string{ticket.FieldParent: epic}, nil)
	mk("A grandchild", "todo", map[string]string{ticket.FieldParent: child}, nil)
	mk("Something it relates to", "todo", nil, map[string][]string{ticket.FieldRelated: {epic}})
	if _, err := m.store.Logbook(shipped, "as-20260911"); err != nil {
		t.Fatal(err)
	}

	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	m.selectByID(epic)
	return m, epic
}

// L opens the window, esc closes it and puts the board back.
func TestLOpensAndClosesTheLinkWindow(t *testing.T) {
	m, _ := linksModel(t)

	m.key(key("L"))
	if m.mode != modeLinks {
		t.Fatalf("L must open the link window, mode = %v", m.mode)
	}
	m.key(key("esc"))
	if m.mode != modeBoard || m.links != nil {
		t.Fatalf("esc must close it, mode = %v links = %v", m.mode, m.links)
	}
}

// The window shows the whole tree and says where each end lives — including
// the child that is finished and no longer on the board, which is the thing
// that used to be invisible.
func TestLinkWindowShowsChildrenWhereverTheyLive(t *testing.T) {
	m, _ := linksModel(t)
	m.key(key("L"))
	out := m.render()

	for _, want := range []string{
		"contains",
		"A child still in play",
		"A child already shipped",
		"A grandchild",
		"logbook · done",
		"related to",
		"Something it relates to",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("link window is missing %q\n%s", want, out)
		}
	}
}

// j and k walk the tickets and never stop on a group heading.
func TestLinkWindowCursorSkipsHeadings(t *testing.T) {
	m, _ := linksModel(t)
	m.key(key("L"))

	for i := 0; i < 8; i++ {
		if _, ok := m.links.current(); !ok {
			t.Fatalf("cursor landed on a heading after %d moves down", i)
		}
		m.key(key("j"))
	}
	for i := 0; i < 12; i++ {
		if _, ok := m.links.current(); !ok {
			t.Fatalf("cursor landed on a heading after %d moves up", i)
		}
		m.key(key("k"))
	}
}

// enter on a ticket that is on the board jumps to its card.
func TestEnterJumpsToALinkedCard(t *testing.T) {
	m, _ := linksModel(t)
	m.key(key("L"))
	// The first selectable row is the first child, which is on the board.
	e, ok := m.links.current()
	if !ok {
		t.Fatal("nothing selected in the link window")
	}
	m.key(key("enter"))
	if m.mode != modeBoard {
		t.Fatalf("enter must return to the board, mode = %v", m.mode)
	}
	if sel := m.selected(); sel == nil || sel.ID != e.Ref.ID {
		t.Fatalf("enter must select %s, got %v", ticket.Handle(e.Ref.ID), sel)
	}
}

// enter on a ticket that is not on the board says so rather than appearing to
// do nothing.
func TestEnterOnAFiledTicketExplainsItself(t *testing.T) {
	m, _ := linksModel(t)
	m.key(key("L"))
	for {
		e, ok := m.links.current()
		if !ok {
			t.Fatal("no filed ticket in the window")
		}
		if strings.Contains(e.Ref.Title, "already shipped") {
			break
		}
		before := m.links.cursor
		m.key(key("j"))
		if m.links.cursor == before {
			t.Fatal("no filed ticket in the window")
		}
	}
	m.key(key("enter"))
	if m.mode != modeMessage {
		t.Fatalf("enter on a filed ticket must explain itself, mode = %v", m.mode)
	}
}

// The window is a modal: the board stays visible behind it. A dialog that
// replaces the screen costs the reader the thing they were reasoning about —
// which card they were on, and which lane it sat in.
func TestLinkWindowIsDrawnOverTheBoard(t *testing.T) {
	m, _ := linksModel(t)
	board := m.render()
	m.key(key("L"))
	out := m.render()

	// A lane header from the board behind, and the box's own border, on one
	// screen.
	for _, want := range []string{"Backlog", "Implementing", "\u256d", "Links · "} {
		if !strings.Contains(out, want) {
			t.Errorf("the modal must leave %q visible\n%s", want, out)
		}
	}
	if out == board {
		t.Error("L must change what is on screen")
	}
	// And closing it puts the board back exactly as it was.
	m.key(key("esc"))
	if got := m.render(); got != board {
		t.Errorf("esc must restore the board unchanged\n%s", got)
	}
}
