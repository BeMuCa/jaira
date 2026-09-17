package tui

import "strings"

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
//
// Carriage returns are already gone: insertText normalises them to "\n" before
// any buffer sees the text.
func foldToOneLine(text string) string {
	lines := strings.Split(text, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if line != "" {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, " ")
}
