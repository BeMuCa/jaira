package tui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The board shows what recently went into the logbook at the foot of the
// terminal lane, below the tickets still waiting to be filed. Filing takes a
// ticket off the board on purpose, and that is right for the work in hand — but
// it also made the board the one place that could not say what had just been
// finished. These cards answer that, for logbook-days days (settings.json,
// edited on the settings screen), and nothing more: they are read, never moved
// or edited. Bringing one back stays 'jaira restore', which is a decision about
// the record rather than a gesture on a card.

// Logbook cards are filled in their own pair of shades instead of the lane's
// greys: a done card still on the board and a filed one must not be told apart
// by their position alone. Teal against the neutral greys, darker on a dark
// terminal and paler on a light one, and the third row says "filed <day>" so
// the difference never rests on colour alone.
const (
	logbookShadeADark  = "23"
	logbookShadeBDark  = "24"
	logbookShadeALight = "195"
	logbookShadeBLight = "153"
)

// isLogged reports whether a card on the board is one of the logbook's rather
// than a ticket on the board. By pointer, not by id: a ticket that is on the
// board and in the logbook at once is a known merge outcome, and the board copy
// must stay workable while its filed twin stays read-only.
func (m *Model) isLogged(t *ticket.Ticket) bool {
	if t == nil {
		return false
	}
	_, ok := m.loggedDay[t]
	return ok
}

// logbookShade is the background of an unselected logbook card, alternating
// the way laneShade does so two stacked filed cards stay apart.
func (m *Model) logbookShade(alt bool) string {
	a, b := logbookShadeADark, logbookShadeBDark
	if !m.darkBG {
		a, b = logbookShadeALight, logbookShadeBLight
	}
	if alt {
		return "5;" + b
	}
	return "5;" + a
}

// logbookFlag is the marker a logbook card carries first among its flags, so
// a narrow column truncates everything else before it.
func (m *Model) logbookFlag(t *ticket.Ticket) string {
	return styOK.Render("⎙ filed " + m.loggedDay[t].Format("2 Jan"))
}

// refuseLogged answers a gesture that would change a logbook card, and
// reports whether it did. The answer names the way back, because the reader
// who pressed m on a filed card wanted to work it again.
func (m *Model) refuseLogged(t *ticket.Ticket) bool {
	if !m.isLogged(t) {
		return false
	}
	m.notify(fmt.Sprintf("%s is in the logbook, filed %s — the board only shows it.\n\nTo work on it again, bring it back with:\n\n  jaira restore %s",
		ticket.Handle(t.ID), m.loggedDay[t].Format("2006-01-02"), filepath.Base(t.Path)), false)
	return true
}

// loadLogbook reads the window of the logbook the board shows. It is called on
// every reload, which is cheap enough: only folders inside the window are
// opened, and the window is days, not the history.
func (m *Model) loadLogbook() {
	m.logged = m.store.LoggedSince(time.Now(), m.logbookDays)
	m.loggedDay = make(map[*ticket.Ticket]time.Time, len(m.logged))
	for _, l := range m.logged {
		m.loggedDay[l.Ticket] = l.Day
	}
}
