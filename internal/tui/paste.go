package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// paste takes text the terminal handed over as a bracketed paste and puts it
// where a typed character would have gone.
//
// A paste is not a keypress. The terminal wraps it in its own escape sequence
// and the decoder turns that into a tea.PasteMsg, so none of it ever reaches
// Model.key — which is why every input field swallowed pasted text until this
// existed.
//
// It cannot be done through the key layout either. cmdKey reads a key by its
// physical position, but the decoder clears Key.Text as soon as a modifier
// beyond shift is down (see the note at keylayout.go's cmdKey), so a ctrl+v
// carries no character to place and there is nothing for a layout to map.
// Handling the event instead makes the layout irrelevant: a Cyrillic or German
// keyboard pastes through the same branch as a US one, because the paste never
// was a key combination to begin with.
func (m *Model) paste(text string) (tea.Model, tea.Cmd) {
	m.insertText(text)
	return m, nil
}

// foldToOneLine folds a block into the one line a single-line field can hold.
//
// Line breaks become a single space rather than being dropped: copying one line
// out of a terminal or a file usually takes the trailing newline with it, and
// discarding on a newline would then swallow the whole paste — which looks
// exactly like the fault this handling fixes. Leading and trailing breaks fall
// away, and a run of them collapses to one space.
//
// Nothing here cuts the text by byte offset, so Cyrillic, umlauts and emoji
// arrive character for character, the same way k.Text does.
func foldToOneLine(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if line != "" {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, " ")
}
