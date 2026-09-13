package tui

import (
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// cmdKey reads a keypress as a command rather than as text.
//
// Every switch in this package matches a command against the character the key
// produced — "j", "q", "/" — and that character is whatever the active keyboard
// layout decided. Switch the layout to Cyrillic and the same physical keys send
// "о", "й", "."; no case matches, and the board stops responding to anything at
// all. cmdKey answers what the key would have been on a US PC-101 keyboard, so
// the next lane stays one "j" away whatever the layout is set to.
//
// Two sources, in this order:
//
//  1. Key.BaseCode, which the terminal reports itself when it speaks the Kitty
//     keyboard protocol (kitty, ghostty, WezTerm, foot) or comes through the
//     Windows Console API. It is correct for every layout in the world and
//     needs no table here. Key.String() ignores it whenever the key produced
//     text — which is every letter, the case that matters — so it is read
//     directly rather than through String().
//  2. usPosition, for the terminals that report nothing. Windows Terminal, where
//     most WSL2 users sit, is one of them, so this is the path that actually
//     runs for the person who reported the fault.
//
// Latin layouts — AZERTY, QWERTZ, Dvorak — need neither: they send ASCII
// letters, and the switches already match those.
//
// Text input does not come through here. A Cyrillic letter typed into a filter
// or an edit field has to stay a Cyrillic letter, so those paths keep reading
// Key.Text.
func cmdKey(k tea.KeyPressMsg) string {
	r, ok := physicalRune(k)
	if !ok {
		return k.String()
	}
	// Named keys and modifier spellings come from String() unchanged; only the
	// single character at the end of the keystroke is layout-dependent.
	var b strings.Builder
	if k.Mod.Contains(tea.ModCtrl) {
		b.WriteString("ctrl+")
	}
	if k.Mod.Contains(tea.ModAlt) {
		b.WriteString("alt+")
	}
	b.WriteRune(r)
	return b.String()
}

// physicalRune returns the US-layout character of a keypress, and whether that
// answer is worth using: false means the key already names its own position and
// cmdKey should hand back what the terminal said. Case is carried over from what was
// typed: shift+ф is Ф is "F", because the board binds "E" and "G" and "X" to
// their own commands.
func physicalRune(k tea.KeyPressMsg) (rune, bool) {
	// Only a printed character can be layout-dependent. Enter, space and the
	// arrows are already named after the physical key, and a terminal speaking
	// the Kitty protocol reports a BaseCode for those too — taking it would
	// turn "enter" into a bare rune no switch matches.
	rs := []rune(k.Text)
	if len(rs) != 1 || rs[0] == ' ' || !unicode.IsGraphic(rs[0]) {
		return 0, false
	}
	// Text is what the key produced, upper case included.
	typed := rs[0]
	base := k.BaseCode
	if base == 0 {
		var ok bool
		if base, ok = usPosition[unicode.ToLower(typed)]; !ok {
			return 0, false
		}
	}
	if unicode.IsUpper(typed) {
		base = unicode.ToUpper(base)
	}
	return base, base != typed
}

// usPosition maps a character to the one the same physical key carries on a US
// PC-101 keyboard.
//
// Only the Cyrillic ЙЦУКЕН layout is in here. The other non-Latin scripts —
// Greek, Armenian, Hebrew, Arabic — are the same mechanic and one more table
// each, to be added when somebody runs the board in one of them; every layout
// is already covered on a terminal that reports Key.BaseCode.
//
// The one Latin entry is "." on the key that is "/" on a US keyboard: on ЙЦУКЕН
// that key sends a full stop, and without this line the search the board opens
// with "/" is unreachable. Nothing in the TUI binds ".", so a US keyboard loses
// no command to it — it gains a second way into the filter.
var usPosition = map[rune]rune{
	'й': 'q', 'ц': 'w', 'у': 'e', 'к': 'r', 'е': 't', 'н': 'y',
	'г': 'u', 'ш': 'i', 'щ': 'o', 'з': 'p', 'х': '[', 'ъ': ']',
	'ф': 'a', 'ы': 's', 'в': 'd', 'а': 'f', 'п': 'g', 'р': 'h',
	'о': 'j', 'л': 'k', 'д': 'l', 'ж': ';', 'э': '\'',
	'я': 'z', 'ч': 'x', 'с': 'c', 'м': 'v', 'и': 'b', 'т': 'n',
	'ь': 'm', 'б': ',', 'ю': '.', 'ё': '`',
	'.': '/',
}
