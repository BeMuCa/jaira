package lane

import "github.com/BeMuCa/jaira/core/ticket"

// SendsBack reports a move to an earlier lane of the board, the one that
// starts a new round. Parking does not: a ticket going to blocked and coming
// back from it is the same round, whichever side of the work lane blocked is
// drawn on. A lane this installation does not know counts nothing.
//
// One definition for the move that counts the round and the gate that limits
// it, so the two can never disagree about what a send-back is.
func (s *Set) SendsBack(from, to string) bool {
	f, t := s.Index(from), s.Index(to)
	if f < 0 || t < 0 || t >= f {
		return false
	}
	for _, id := range []string{from, to} {
		if l, ok := s.Get(id); ok && l.RequiresBlockedReason {
			return false
		}
	}
	return true
}

// RoundsLeft is how many more times l may send t back before a person
// decides, and whether l declares a limit at all. It reads back-count, the
// returns in a row the move recorded: only a run of this lane's own
// send-backs uses rounds up, so a person's return or one from another lane
// gives the lane its full count again.
func (l *Lane) RoundsLeft(t *ticket.Ticket) (int, bool) {
	if l.MaxRounds <= 0 {
		return 0, false
	}
	used := 0
	if t.BackFrom == l.ID && t.BackBy == ticket.BackByLane {
		used = t.BackCount
	}
	return max(l.MaxRounds-used, 0), true
}

// AsksPerson reports a lane a person answers in: it takes a question and no
// agent works it. The one reading for both ends of max-rounds — where the
// refusal sends the ticket, and which move the refusal lets through — so an
// agentic lane that merely requires a question is never a way around it.
func (l *Lane) AsksPerson() bool {
	return l.RequiresQuestion && !l.Agentic
}

// PersonLane names the lane a person answers in when l has used up its
// rounds: the first of l's back edges a person answers in, else the board's
// first such lane. Empty when the board has none: a sign-off lane is not
// offered instead, since it is for reviewed work, not an undecided loop.
func (s *Set) PersonLane(l *Lane) string {
	for _, id := range l.RejectsTo {
		if r, ok := s.Get(id); ok && r.AsksPerson() {
			return r.ID
		}
	}
	for _, r := range s.Lanes {
		if r.AsksPerson() {
			return r.ID
		}
	}
	return ""
}
