package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

// pasteInto hands the model the event a terminal sends for a bracketed paste —
// which is not a key, which is the whole reason this path exists.
func pasteInto(m *Model, s string) {
	m.Update(tea.PasteMsg{Content: s})
}

// The reported fault: nothing could be put into the search filter from the
// clipboard. The pasted text has to land in the field and narrow the board at
// once, exactly as typing it does.
func TestPasteIntoFilterLandsAndFilters(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.key(key("/"))
	if m.mode != modeFilter {
		t.Fatalf("/ opened mode %v, want modeFilter", m.mode)
	}

	pasteInto(m, "review")

	if m.input != "review" {
		t.Errorf("the filter buffer holds %q after a paste, want %q", m.input, "review")
	}
	if m.filter != "review" {
		t.Errorf("the board filter is %q after a paste, want %q — it did not narrow", m.filter, "review")
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "review") {
		t.Errorf("the pasted text is not on screen:\n%s", out)
	}
}

// Typing and pasting have to reach the same buffer, or half a query typed and
// half of it pasted comes out scrambled.
func TestPasteAppendsToTypedText(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.key(key("/"))
	typeInto(m, "lane:")
	pasteInto(m, "review")

	if m.input != "lane:review" {
		t.Errorf("typed text plus a paste gave %q, want %q", m.input, "lane:review")
	}
}

// Every mode that collects text read k.Text and nothing else, so all four
// swallowed a paste. One case per mode, because they are four separate
// branches.
func TestPasteReachesEveryInputMode(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		m := newTestModel(t, 150, 32)
		m.key(key("n"))
		if m.mode != modeCreate {
			t.Fatalf("n opened mode %v, want modeCreate", m.mode)
		}
		pasteInto(m, "a pasted title")
		if m.input != "a pasted title" {
			t.Errorf("create buffer holds %q after a paste", m.input)
		}
	})

	t.Run("delete", func(t *testing.T) {
		m := newTestModel(t, 150, 32)
		tk := openFirstTicket(t, m)
		m.key(key("X"))
		if m.mode != modeDelete {
			t.Fatalf("X opened mode %v, want modeDelete", m.mode)
		}
		// The confirmation is the handle typed back; pasting it counts too.
		pasteInto(m, ticket.Handle(tk.ID))
		if m.input != ticket.Handle(tk.ID) {
			t.Errorf("delete buffer holds %q, want the handle %q", m.input, ticket.Handle(tk.ID))
		}
	})

	t.Run("edit", func(t *testing.T) {
		m := newTestModel(t, 150, 32)
		openFirstTicket(t, m)
		m.key(key("e"))
		if m.mode != modeEdit {
			t.Fatalf("e opened mode %v, want modeEdit", m.mode)
		}
		m.editBuf = ""
		pasteInto(m, "pasted context")
		if m.editBuf != "pasted context" {
			t.Errorf("edit buffer holds %q after a paste", m.editBuf)
		}
	})
}

// The buffer is rune-based everywhere else, and a paste must not be the one
// place that cuts by byte. The backspace afterwards is the check that matters:
// it removes one character, not one byte of one.
func TestPasteKeepsMultiByteTextWhole(t *testing.T) {
	const pasted = "Grüße Привет 🙂"

	m := newTestModel(t, 150, 32)
	m.key(key("/"))
	pasteInto(m, pasted)

	if m.input != pasted {
		t.Fatalf("multi-byte paste arrived as %q, want %q", m.input, pasted)
	}
	m.key(key("backspace"))
	want := string([]rune(pasted)[:len([]rune(pasted))-1])
	if m.input != want {
		t.Errorf("backspace after a paste left %q, want %q — it cut a byte, not a character", m.input, want)
	}
}

// The filter is one line, so a multi-line paste is folded to spaces rather than
// discarded: a line copied out of a terminal carries its trailing newline along,
// and discarding on a newline would swallow the whole paste — which is the very
// fault this handling fixes.
func TestMultiLinePasteIsFoldedToSpaces(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"trailing newline is dropped", "review\n", "review"},
		{"inner newline becomes a space", "lane\nreview", "lane review"},
		{"a run of newlines becomes one space", "lane\n\n\nreview", "lane review"},
		{"carriage returns count as newlines", "lane\r\nreview\r\n", "lane review"},
		{"a lone carriage return counts too", "lane\rreview\r", "lane review"},
		{"leading newlines fall away", "\n\nreview", "review"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestModel(t, 150, 32)
			m.key(key("/"))
			pasteInto(m, c.in)
			if m.input != c.want {
				t.Errorf("pasting %q gave %q, want %q", c.in, m.input, c.want)
			}
			if strings.Contains(m.input, "\n") {
				t.Errorf("a newline survived into the one-line filter: %q", m.input)
			}
		})
	}
}

// The field editor is the one multi-line buffer — enter inserts a line there —
// so a pasted paragraph keeps its shape instead of being folded flat.
func TestPasteKeepsLinesInTheFieldEditor(t *testing.T) {
	m := newTestModel(t, 150, 32)
	openFirstTicket(t, m)
	m.key(key("e"))
	m.editBuf = ""
	pasteInto(m, "first line\r\nsecond line")

	if m.editBuf != "first line\nsecond line" {
		t.Errorf("the editor folded a pasted paragraph: %q", m.editBuf)
	}
}

// A terminal may send a bare CR inside a bracketed paste. Raw it would reach the
// editor buffer and from there the ticket file, so it is normalised like CRLF.
func TestPasteNormalisesALoneCarriageReturnInTheFieldEditor(t *testing.T) {
	m := newTestModel(t, 150, 32)
	openFirstTicket(t, m)
	m.key(key("e"))
	m.editBuf = ""
	pasteInto(m, "first line\rsecond line")

	if m.editBuf != "first line\nsecond line" {
		t.Errorf("a bare carriage return survived into the editor: %q", m.editBuf)
	}
}
