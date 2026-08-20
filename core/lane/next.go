package lane

import "github.com/BeMuCa/jaira/core/ticket"

// Next reports the lane this ticket goes to when the step it is in is finished.
//
// It exists because the route was only ever implicit. The column order carries
// it, and a ticket's Options decide which optional steps are part of that
// ticket's path at all — so an agent wanting to know "where does this go next"
// had to load the lane set, work out the order, and apply the ticket's Options
// itself. Every agent doing that separately is the same derivation written four
// times, each free to get it wrong.
//
// Three kinds of lane are never the answer:
//
//   - a lane this ticket opted out of (RequiresOption not ticked): it is not on
//     this ticket's path, and moving into it is refused by the gate.
//   - a parking lane (RequiresBlockedReason): being blocked is a thing that
//     happens to a ticket, not the next step in its work.
//   - a human checkpoint that only asks a question (RequiresQuestion): a ticket
//     goes there because something stopped, not because it advanced.
//
// A ticket that is itself parked or waiting on an answer also has no next lane.
// It resumes where it stopped, and nothing on the board records where that was —
// so walking forward from there would send work past the very steps that raised
// the question. Measured against a real board: a ticket in HITL reported "review"
// and would have skipped critique and optimize.
//
// A ticket in a terminal lane, in an unknown lane, or in the last lane on its
// path has no next lane either, and nil says so. Nothing enforces the answer: 'jaira
// move' re-checks the gates whatever route the caller took.
func (s *Set) Next(t *ticket.Ticket) *Lane {
	from := s.Index(t.Status)
	if from < 0 {
		return nil
	}
	if cur, ok := s.Get(t.Status); ok &&
		(cur.Terminal || cur.RequiresBlockedReason || cur.RequiresQuestion) {
		return nil
	}
	for _, l := range s.Lanes[from+1:] {
		if l.RequiresBlockedReason || l.RequiresQuestion {
			continue
		}
		if l.RequiresOption != "" && !t.OptionSet(l.RequiresOption) {
			continue
		}
		return l
	}
	return nil
}
