package tui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/notify"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// refFetchEvery answers how often the board looks at the remote, from the same
// setting the CLI's background fetch uses.
//
// It used to be a constant here, and that was a quiet trap: somebody who set
// 'fetch-every' in their settings would have found the board carrying on at
// its own pace, with nothing saying why. Two seconds is right for the local
// rescan; a network round trip on that cadence is not.
func refFetchEvery() time.Duration { return settings.Load().FetchInterval() }

// refFetchedMsg carries the result of a background fetch back into the update
// loop.
type refFetchedMsg struct {
	arrivals []refsync.Arrival
	err      error
}

// newSyncer wires the board's ticket writes to the refs they travel on, with
// the same settings and the same idea of "me" the CLI uses.
func newSyncer(s *ticket.Store, me string) *refsync.Syncer {
	y := refsync.New(s, settings.Load().RemoteName(), me)
	y.IsMine = func(assignee string) bool {
		return assignee != "" && identity.IsMe(s.Root, assignee)
	}
	s.Source = y
	return y
}

// refMarks is the local, network-free half: which tickets hold a write that has
// not reached the remote.
//
// It reads the outbox, which is files already on disk, because this runs inside
// reload on the two-second timer and nothing on that path may wait for a
// remote. Tickets that live only on a ref need no counting here: the store
// hands them to List like any other ticket, so they are cards.
func (m *Model) refMarks() {
	m.unsent = map[string]bool{}
	if m.refSync == nil || m.refSync.Usable() != nil {
		return
	}
	if entries, err := m.refSync.Box.List(); err == nil {
		for _, e := range entries {
			m.unsent[e.ID] = true
		}
	}
}

// fetchRefs asks the remote what the ticket refs now say, off the update loop.
//
// It is a command rather than a call inside Update for the reason every network
// touch in a TUI has to be: View and Update must never block. Bubble Tea runs
// the returned function on its own goroutine and delivers the result as a
// message, which is the documented way to do exactly this.
func fetchRefs(y *refsync.Syncer) tea.Cmd {
	if y == nil {
		return nil
	}
	return func() tea.Msg {
		arrivals, err := y.Incoming()
		return refFetchedMsg{arrivals: arrivals, err: err}
	}
}

// refTick re-arms the background fetch.
func refTick() tea.Cmd {
	return tea.Tick(refFetchEvery(), func(time.Time) tea.Msg { return refTickMsg{} })
}

type refTickMsg struct{}

// announceArrivals says what arrived, and raises a desktop notification for a
// ticket that has just become mine.
//
// Only mine, and only changed: a ticket that has been assigned to me since last
// week is not news, and a board that pops up on every fetch teaches the user to
// ignore it.
func (m *Model) announceArrivals(arrivals []refsync.Arrival) {
	notifyOK := settings.Load().NotifyEnabled()
	var mine []refsync.Arrival
	for _, a := range arrivals {
		if a.Mine && a.Changed {
			mine = append(mine, a)
		}
	}
	if len(mine) == 0 {
		return
	}
	for _, a := range mine {
		if !notifyOK {
			break
		}
		title := a.Title
		if title == "" {
			title = ticket.Handle(a.ID)
		}
		notify.Send("jaira: a ticket is yours", fmt.Sprintf("%s — %s", ticket.Handle(a.ID), title))
	}
	// Shown as a line on the board, never as a screen of its own. An arrival is
	// news the person did not ask for, arriving while they are in the middle of
	// something: taking the screen for it would interrupt them and could
	// swallow the keystroke they were already typing.
	if len(mine) == 1 {
		m.flash(fmt.Sprintf("%s is yours: %s", ticket.Handle(mine[0].ID), mine[0].Title))
		return
	}
	m.flash(fmt.Sprintf("%d tickets were assigned to you", len(mine)))
}
