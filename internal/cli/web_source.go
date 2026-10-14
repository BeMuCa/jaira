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

// Ticket is one ticket opened on the page, read the way 'jaira show' reads
// it: past the board, so a finished ticket a card links to still opens.
func (w *webSource) Ticket(id string) (web.TicketDetail, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	t, _, err := w.s.LoadAnywhere(id)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			return web.TicketDetail{}, fmt.Errorf("ticket %s: %w", id, web.ErrNotFound)
		}
		return web.TicketDetail{}, err
	}
	d := web.TicketDetail{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Goal: t.Goal, Context: t.Context}
	for _, it := range t.DoDItems {
		d.DoD = append(d.DoD, web.Check{Text: it.Text, Done: it.Checked()})
	}
	for _, n := range ticket.Notes(t.Body) {
		d.Notes = append(d.Notes, web.Note{At: n.At, Author: n.Author, Text: n.Text})
	}
	return d, nil
}

// CreateEpic starts an empty epic, the name normalised the way 'jaira epic
// create' normalises it.
func (w *webSource) CreateEpic(raw, title, description string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	name, _, err := epic.NormalizeName(raw)
	if err != nil {
		return fmt.Errorf("%w: %v", web.ErrInvalid, err)
	}
	_, err = epic.Create(w.s, epicRefs(), name, title, description, nil)
	if errors.Is(err, epic.ErrExists) {
		return fmt.Errorf("epic %s: %w", name, web.ErrExists)
	}
	return err
}

// Assign puts one ticket into an epic, and with from set takes it out of
// that one: a drag from one epic onto another. The add comes first, so a
// failure in between leaves the ticket in both epics rather than in none, and
// the error says which of the two writes it was. The page sends the full id
// it was given, but a handle is resolved too, the way 'jaira epic add'
// resolves one.
func (w *webSource) Assign(name, from, id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	t, err := w.load(id)
	if err != nil {
		return err
	}
	if err := w.assign(name, t.ID, true); err != nil {
		return fmt.Errorf("adding %s to %s: %w", ticket.Handle(t.ID), name, err)
	}
	if from == "" || from == name {
		return nil
	}
	if err := w.assign(from, t.ID, false); err != nil {
		return fmt.Errorf("%s is in %s now, but taking it out of %s failed: %w", ticket.Handle(t.ID), name, from, err)
	}
	return nil
}

// Unassign takes a ticket out of an epic. Only a member line can go: a
// follow-up counts through follows: or parent: and has none to remove.
func (w *webSource) Unassign(name, id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	t, err := w.load(id)
	if err != nil {
		return err
	}
	if err := w.assign(name, t.ID, false); err != nil {
		return fmt.Errorf("taking %s out of %s: %w", ticket.Handle(t.ID), name, err)
	}
	return nil
}

func (w *webSource) load(id string) (*ticket.Ticket, error) {
	t, err := w.s.Load(id)
	if errors.Is(err, ticket.ErrNotFound) {
		return nil, fmt.Errorf("ticket %s: %w", id, web.ErrNotFound)
	}
	return t, err
}

// assign is one write of one epic's ref through epic.Assign, the call 'jaira
// epic add' and 'jaira epic rm' make.
func (w *webSource) assign(name, id string, add bool) error {
	_, changed, err := epic.Assign(w.s, epicRefs(), name, []string{id}, add)
	if errors.Is(err, epic.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no epic %s: %w", name, web.ErrNotFound)
	}
	if err == nil && !add && len(changed) == 0 {
		return fmt.Errorf("it is no member of %s, it counts through follows or parent: %w", name, web.ErrConflict)
	}
	return err
}
