package cli

import (
	"errors"
	"sync"

	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/core/web"
)

// webSource is the page's view of this board. Requests arrive concurrently and
// the store was written for one command at a time, so every call is serialised.
type webSource struct {
	mu sync.Mutex
	s  *ticket.Store
}

func newWebSource(s *ticket.Store) *webSource { return &webSource{s: s} }

// board loads the open tickets and the index that says which are finished,
// the same two things 'jaira list' and 'jaira epic ls' start from.
func (w *webSource) board() ([]*ticket.Ticket, *link.Index, error) {
	env, all, err := loadEnv(w.s)
	if err != nil {
		return nil, nil, err
	}
	return all, link.Build(w.s, env.Lanes, all), nil
}

func webTicket(t *ticket.Ticket, done, direct bool) web.Ticket {
	return web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Done: done, Direct: direct}
}

func (w *webSource) Epics() ([]web.Epic, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return nil, nil
}

func (w *webSource) Unassigned() ([]web.Ticket, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	all, ix, err := w.board()
	if err != nil {
		return nil, err
	}
	out := make([]web.Ticket, 0, len(all))
	for _, t := range all {
		if ix.Satisfied(t.ID) {
			continue
		}
		out = append(out, webTicket(t, false, true))
	}
	return out, nil
}

func (w *webSource) Assign(epic, id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return errors.New("epics are not on this build yet")
}
