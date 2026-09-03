package ticket

import "sort"

// Trimmed records one ticket a lane cap moved off the board, so the caller can
// say what happened — the cap is allowed to move files only because it never
// does so silently.
type Trimmed struct {
	ID    string
	Title string
	// Path is where the file now lives, under the logbook.
	Path string
}

// Overflow returns the tickets of a lane beyond the newest keep, oldest first.
// "Newest" is measured by updated-at: a ticket in a terminal lane is normally
// never touched again, so updated-at is its acceptance time. keep <= 0 means
// the lane has no cap, so nothing ever overflows. A board with unreadable
// tickets refuses to answer rather than guess at which ticket is oldest.
func (s *Store) Overflow(lane string, keep int) ([]*Ticket, error) {
	if keep <= 0 {
		return nil, nil
	}
	all, err := s.List()
	if err != nil {
		return nil, err
	}
	var ts []*Ticket
	for _, t := range all {
		if t.Status == lane {
			ts = append(ts, t)
		}
	}
	if len(ts) <= keep {
		return nil, nil
	}
	// Newest first, ties broken by id — the same order the board displays.
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].UpdatedAt.Equal(ts[j].UpdatedAt) {
			return ts[i].ID < ts[j].ID
		}
		return ts[i].UpdatedAt.After(ts[j].UpdatedAt)
	})
	over := ts[keep:]
	// Oldest first: the order they leave in.
	for i, j := 0, len(over)-1; i < j; i, j = i+1, j-1 {
		over[i], over[j] = over[j], over[i]
	}
	return over, nil
}

// TrimLane enforces a lane's cap after something new landed in it: the tickets
// beyond the newest keep leave for the given logbook folder, oldest first. It
// returns what moved — including on error, so a caller can still report the
// tickets that had already left before the failure.
func (s *Store) TrimLane(lane string, keep int, folder string) ([]Trimmed, error) {
	over, err := s.Overflow(lane, keep)
	if err != nil {
		return nil, err
	}
	var out []Trimmed
	for _, t := range over {
		dst, err := s.Logbook(t.ID, folder)
		if err != nil {
			return out, err
		}
		out = append(out, Trimmed{ID: t.ID, Title: t.Title, Path: dst})
	}
	return out, nil
}
