package lane

import "github.com/BeMuCa/jaira/core/ticket"

// Settle applies a lane's file rules after a move landed a ticket there — the
// one place that knows which rule wins. A doorway (logbook-on-entry) files the
// whole lane into the logbook folder, prepare running on each ticket first
// (the callers stamp commits there); a capped lane (holds) trims the oldest
// beyond the newest Holds, never the just-moved ticket. filed reports which
// rule ran, so a caller can word its message. Both UIs call this instead of
// carrying the switch themselves — four status write-sites were one too many
// places for the same decision.
func Settle(s *ticket.Store, l *Lane, folder, moved string, prepare func(*ticket.Ticket) error) (trimmed []ticket.Trimmed, filed bool, err error) {
	switch {
	case l == nil:
		return nil, false, nil
	case l.LogbookOnEntry:
		out, err := s.FileLane(l.ID, folder, prepare)
		return out, true, err
	case l.Holds > 0:
		out, err := s.TrimLane(l.ID, l.Holds, folder, moved)
		return out, false, err
	}
	return nil, false, nil
}
