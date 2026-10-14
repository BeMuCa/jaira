package cli

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/BeMuCa/jaira/core/web"
)

// webSource is the page's view of this board, built from the same calls
// 'jaira epic ls' and 'jaira epic add' make, so the page and the CLI cannot
// disagree about what an epic holds or how far it is. Requests arrive
// concurrently and the store was written for one command at a time, so every
// call is serialised.
type webSource struct {
	mu sync.Mutex
	s  *ticket.Store
}

func newWebSource(s *ticket.Store) *webSource { return &webSource{s: s} }

// board loads the open tickets, the index that says which are finished, and
// the epics, through the same helpers 'jaira epic ls' uses. A board with
// nowhere to keep epics has none, which the page shows as an empty column
// rather than an error.
func (w *webSource) board() ([]*ticket.Ticket, *link.Index, []*epic.Epic, error) {
	ix, all, err := boardIndex(w.s)
	if err != nil {
		return nil, nil, nil, err
	}
	epics, err := epic.LoadAll(epicRefs())
	if err != nil && !errors.Is(err, epic.ErrNoRefs) {
		return nil, nil, nil, err
	}
	return all, ix, epics, nil
}

func (w *webSource) Epics() ([]web.Epic, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, ix, epics, err := w.board()
	if err != nil {
		return nil, err
	}
	out := make([]web.Epic, 0, len(epics))
	for _, e := range epics {
		p := epic.Count(e, ix)
		v := web.Epic{Name: e.Name, Title: e.Title, Description: e.Description, Done: p.Done, Total: p.Total}
		for _, t := range p.Tickets {
			v.Tickets = append(v.Tickets, web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Lane, Done: t.Done, Direct: t.Direct})
		}
		out = append(out, v)
	}
	return out, nil
}

func (w *webSource) Unassigned() ([]web.Ticket, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	all, ix, epics, err := w.board()
	if err != nil {
		return nil, err
	}
	loose := epic.Unassigned(epics, ix, all)
	out := make([]web.Ticket, 0, len(loose))
	for _, t := range loose {
		out = append(out, web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Direct: true})
	}
	return out, nil
}

// Assign puts one ticket into an epic. The page sends the full id it was
// given, but a handle is resolved too, the way 'jaira epic add' resolves one.
func (w *webSource) Assign(name, id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	t, err := w.s.Load(id)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			return fmt.Errorf("ticket %s: %w", id, web.ErrNotFound)
		}
		return err
	}
	_, _, err = epic.Assign(w.s, epicRefs(), name, []string{t.ID}, true)
	if errors.Is(err, epic.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("epic %s: %w", name, web.ErrNotFound)
	}
	return err
}
