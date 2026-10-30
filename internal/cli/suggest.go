package cli

import (
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/suggest"
	"github.com/BeMuCa/jaira/core/ticket"
)

// suggestionsFor is what jaira says after ticket id stands in status, where an
// empty status means it is off the board. The milestones consulted are the ones
// that hold the ticket, and the members of each are read as the board has them.
func suggestionsFor(s *ticket.Store, lanes *lane.Set, id, status string) ([]suggest.Line, error) {
	all, err := milestone.LoadAll(s.Root)
	if err != nil {
		return nil, err
	}
	var groups []suggest.Group
	for _, m := range all {
		if m.Has(id) {
			groups = append(groups, groupOf(s, m))
		}
	}
	return suggest.ForTicket(lanes, id, status, groups), nil
}

// groupOf reads a milestone's members the way 'jaira logbook <milestone>' does:
// a member that does not load is already off the board, which counts as done.
func groupOf(s *ticket.Store, m *milestone.Milestone) suggest.Group {
	g := suggest.Group{Name: m.Name, Filed: m.Filed()}
	for _, id := range m.Members() {
		status := ""
		if t, err := s.Load(id); err == nil {
			status = t.Status
		}
		g.Members = append(g.Members, suggest.Member{ID: id, Status: status})
	}
	return g
}
