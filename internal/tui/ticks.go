package tui

import (
	"strings"
	"unicode"
)

// tickPicker is the open ticking picker — u over people, M over milestones —
// and the one mechanism both run on. Each picker builds its own rows and
// draws them its own way; what space, enter, x and esc do, and what enter
// writes, is the same for both, so it lives here once.
type tickPicker struct {
	key    string          // the filter key the ticks are written under
	idx    int             // the cursor
	picked map[string]bool // the ticked names, lower-cased
}

// openTicks opens a ticking picker that owns the filter's key condition.
// Whatever the filter already holds under key comes back ticked, so
// reopening shows the board as it is narrowed rather than a blank list that
// enter would then wipe out.
func (m *Model) openTicks(in mode, key string) {
	m.pick = tickPicker{key: key, picked: map[string]bool{}}
	_, picked := splitTerm(m.filter, key)
	for _, name := range picked {
		m.pick.picked[name] = true
	}
	m.returnTo = m.mode
	m.mode = in
}

// keyTicks drives the open ticking picker over names, its rows in the order
// it shows them; own is the key that opened it, which also closes it. esc
// leaves the filter as it was: the ticks live in the picker until enter
// writes them.
func (m *Model) keyTicks(s, own string, names []string) {
	p := &m.pick
	n := len(names)
	switch s {
	case "esc", "q", own:
		m.mode = m.returnTo
	case "j", "down":
		if n > 0 {
			p.idx = (p.idx + 1) % n
		}
	case "k", "up":
		if n > 0 {
			p.idx = (p.idx - 1 + n) % n
		}
	case "space":
		if p.idx >= 0 && p.idx < n {
			k := strings.ToLower(names[p.idx])
			p.picked[k] = !p.picked[k]
		}
	case "enter":
		m.applyTicks(names)
		m.mode = m.returnTo
	case "x":
		p.picked = map[string]bool{}
		m.applyTicks(names)
		m.mode = m.returnTo
	}
}

// applyTicks writes the ticked names into the board filter in place of any
// condition it had under the picker's key, keeping everything else typed
// there — so the people u ticked survive M, and the other way round. It
// writes the ordinary filter rather than a second kind of narrowing, so /
// shows it and esc on the board clears it.
func (m *Model) applyTicks(names []string) {
	var picked []string
	for _, name := range names {
		if !m.pick.picked[strings.ToLower(name)] {
			continue
		}
		// Quoted where the filter would otherwise cut the name apart: at a
		// space between words, or at a comma between alternatives.
		if strings.ContainsFunc(name, func(r rune) bool { return unicode.IsSpace(r) || r == ',' }) {
			name = `"` + name + `"`
		}
		picked = append(picked, name)
	}
	m.setFilterTerm(m.pick.key, picked...)
}

// tickMark is the tick in front of name's row in the open picker.
func (m *Model) tickMark(name string) string {
	if m.pick.picked[strings.ToLower(name)] {
		return styOK.Render("●")
	}
	return styMeta.Render("○")
}
