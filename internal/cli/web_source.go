package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

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
	all, ix, epics, err := w.board()
	if err != nil {
		return nil, err
	}
	cards := w.cards(all, ix)
	out := make([]web.Epic, 0, len(epics))
	for _, e := range epics {
		p := epic.Count(e, ix)
		v := web.Epic{Name: e.Name, Title: e.Title, Description: e.Description, Done: p.Done, Total: p.Total}
		for _, t := range p.Tickets {
			c, ok := cards[t.ID]
			if !ok {
				c = web.Ticket{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Lane}
			}
			c.Done, c.Direct = t.Done, t.Direct
			v.Tickets = append(v.Tickets, c)
		}
		out = append(out, v)
	}
	return out, nil
}

func (w *webSource) Unassigned(everything bool) ([]web.Ticket, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	all, ix, epics, err := w.board()
	if err != nil {
		return nil, err
	}
	cards := w.cards(all, ix)
	if !everything {
		loose := epic.Unassigned(epics, ix, all)
		out := make([]web.Ticket, 0, len(loose))
		for _, t := range loose {
			out = append(out, cards[t.ID])
		}
		return out, nil
	}
	// Every ticket no epic counts, the board's in its own order and the
	// filed ones after them, most recently filed first.
	var members []string
	for _, e := range epics {
		members = append(members, e.Members()...)
	}
	in := map[string]bool{}
	for _, id := range epic.Covered(members, ix) {
		in[id] = true
	}
	var out, filed []web.Ticket
	for _, t := range all {
		if !in[t.ID] {
			out = append(out, cards[t.ID])
		}
	}
	for _, c := range cards {
		if !in[c.ID] && (c.Place == string(link.PlaceLogbook) || c.Place == string(link.PlaceArchive)) {
			filed = append(filed, c)
		}
	}
	sort.Slice(filed, func(i, j int) bool {
		if !filed[i].FiledAt.Equal(filed[j].FiledAt) {
			return filed[i].FiledAt.After(filed[j].FiledAt)
		}
		return filed[i].ID > filed[j].ID
	})
	return append(out, filed...), nil
}

// cards is every ticket this clone knows as the page lists it, by id: the
// board and the refs from all, the logbook and the archive read here. A filed
// ticket is read for its frontmatter only, as the link index reads it.
func (w *webSource) cards(all []*ticket.Ticket, ix *link.Index) map[string]web.Ticket {
	out := make(map[string]web.Ticket, len(all))
	for _, t := range all {
		place := link.PlaceBoard
		if t.ReadOnly {
			place = link.PlaceRef
		}
		c := card(t, place)
		c.Done = ix.Satisfied(t.ID)
		out[t.ID] = c
	}
	for id, path := range w.s.FiledAwayIDs() {
		if _, ok := out[id]; ok {
			continue // on the board and filed at once; the board's copy wins, as 'jaira show' has it
		}
		t, err := ticket.LoadHeader(path)
		if err != nil {
			continue
		}
		place := link.PlaceOf(path)
		c := card(t, place)
		c.Done = ix.Satisfied(id)
		if place == link.PlaceLogbook {
			c.FiledAt, _ = ticket.LogbookFolderDay(filepath.Base(filepath.Dir(path)), time.Local)
		} else if r := rejection(path); r != nil {
			c.Rejected, c.FiledAt = r.Kind, r.At
		}
		out[id] = c
	}
	return out
}

func card(t *ticket.Ticket, place link.Place) web.Ticket {
	return web.Ticket{
		ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Direct: true,
		Place: string(place), Goal: t.Goal, Tags: t.Tags,
		Creator: t.Creator, Assignee: t.Assignee, UpdatedBy: t.UpdatedBy,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

// rejection reads what 'jaira reject' wrote onto a ticket, nil when it was
// never rejected. The keys are read straight off the frontmatter rather than
// through ticket.Ticket, so a ticket rejected by a newer build shows its
// reason here too.
func rejection(path string) *web.Rejection {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	d, err := ticket.ParseDoc(raw)
	if err != nil {
		return nil
	}
	f, err := d.Fields()
	if err != nil {
		return nil
	}
	str := func(key string) string {
		v, _, _ := f.Scalar(key)
		return v
	}
	r := &web.Rejection{Kind: str("rejected-kind"), Why: str("rejected-why"), By: str("rejected-by")}
	if r.Kind == "" {
		return nil
	}
	r.At, _ = time.Parse(time.RFC3339, str("rejected-at"))
	return r
}

// Ticket is one ticket opened on the page, read the way 'jaira show' reads
// it: past the board, so a finished ticket a card links to still opens.
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
	place := link.PlaceBoard
	switch {
	case t.ReadOnly:
		place = link.PlaceRef
	case w.s.FiledAwayIDs()[t.ID] == path:
		place = link.PlaceOf(path)
	}
	d := web.TicketDetail{ID: t.ID, Handle: ticket.Handle(t.ID), Title: t.Title, Lane: t.Status, Place: string(place), Goal: t.Goal, Context: t.Context}
	if place == link.PlaceArchive {
		d.Rejected = rejection(path)
	}
	for _, it := range t.DoDItems {
		d.DoD = append(d.DoD, web.Check{Text: it.Text, Done: it.Checked()})
	}
	for _, n := range ticket.Notes(t.Body) {
		d.Notes = append(d.Notes, web.Note{At: n.At, Author: n.Author, Text: n.Text})
	}
	return d, nil
}

// Lanes is the board's lanes in its order, each with the kind the page
// colours it by.
func (w *webSource) Lanes() ([]web.Lane, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	set, err := lane.Load(w.s.Root)
	if err != nil {
		return nil, err
	}
	var out []web.Lane
	for _, l := range set.Lanes {
		kind := web.LaneWork
		switch {
		case l.Terminal:
			kind = web.LaneTerminal
		case l.RequiresBlockedReason:
			kind = web.LaneParking
		case l.RequiresHumanExit || l.RequiresQuestion:
			// a person's lane, as the agent block in CLAUDE.md marks it
			kind = web.LaneHuman
		}
		out = append(out, web.Lane{ID: l.ID, Name: l.Name, Kind: kind})
	}
	return out, nil
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
