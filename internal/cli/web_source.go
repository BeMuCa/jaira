package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/lane"
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

// archived is the lane the page shows for a ticket filed into the archive or
// the logbook. The file keeps the lane it was filed from, and on a card
// "backlog" reads as work still to do.
const archived = "archive"

// signoff is the lane whose tickets the epic's progress counts apart: finished,
// waiting for a person to accept them. A board without it counts none.
const signoff = "signoff"

// board loads the open tickets, the index that says which are finished, and
// the epics, through the same helpers 'jaira epic ls' uses. A board with
// nowhere to keep epics has none, which the page shows as an empty column
// rather than an error.
//
// A ticket file older than its ref is replaced by the ref before the index is
// built. Lanes that change no code write only the ref, so the file in this
// checkout stays where the last code commit left it; read as it is, the page
// would show a ticket in signoff as still in progress.
func (w *webSource) board() ([]*ticket.Ticket, *link.Index, []*epic.Epic, error) {
	lanes, err := lane.Load(w.s.Root)
	if err != nil {
		return nil, nil, nil, err
	}
	all, err := w.s.List()
	if err != nil {
		var pe *ticket.PartialError
		if !errors.As(err, &pe) {
			return nil, nil, nil, err
		}
	}
	newer := refs.Newer(all)
	for i, t := range all {
		if n, ok := newer[t.ID]; ok {
			all[i] = n
		}
	}
	ix := link.Build(w.s, lanes, all)
	epics, err := epic.LoadAll(epicRefs())
	if err != nil && !errors.Is(err, epic.ErrNoRefs) {
		return nil, nil, nil, err
	}
	return all, ix, epics, nil
}

func (w *webSource) Epics() ([]web.Epic, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	all, ix, epics, err := w.board()
	if err != nil {
		return nil, err
	}
	open := make(map[string]*ticket.Ticket, len(all))
	for _, t := range all {
		open[t.ID] = t
	}
	out := make([]web.Epic, 0, len(epics))
	for _, e := range epics {
		p := epic.Count(e, ix)
		v := web.Epic{Name: e.Name, Title: e.Title, Description: e.Description, Done: p.Done, Total: p.Total}
		for _, t := range p.Tickets {
			c := web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Lane, Done: t.Done, Direct: t.Direct}
			if r, _ := ix.Lookup(t.ID); r.Place == link.PlaceArchive || r.Place == link.PlaceLogbook {
				c.Lane = archived
			}
			if c.Lane == signoff {
				v.Signoff++
			}
			// A member filed away is in no listing; its file still has the dates.
			full, ok := open[t.ID]
			if !ok {
				full, _, _ = w.s.LoadAnywhere(t.ID)
			}
			if full != nil {
				c.CreatedAt, c.UpdatedAt = full.CreatedAt, full.UpdatedAt
			}
			v.Tickets = append(v.Tickets, c)
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
		out = append(out, web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Direct: true, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt})
	}
	return out, nil
}

// Ticket is one ticket opened on the page, read the way 'jaira show' reads
// it: past the board, so a finished ticket a card links to still opens. Like
// the cards, it shows the ref when the file here is older.
func (w *webSource) Ticket(id string) (web.TicketDetail, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	t, path, err := w.s.LoadAnywhere(id)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			return web.TicketDetail{}, fmt.Errorf("ticket %s: %w", id, web.ErrNotFound)
		}
		return web.TicketDetail{}, err
	}
	if n, ok := refs.Newer([]*ticket.Ticket{t})[t.ID]; ok {
		t = n
	}
	d := web.TicketDetail{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Goal: t.Goal, Context: t.Context}
	if path != "" && filepath.Dir(path) != w.s.TicketsDir() {
		d.Lane = archived
	}
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
