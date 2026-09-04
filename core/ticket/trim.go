package ticket

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

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
// never touched again, so updated-at is its acceptance time. newest names the
// ticket whose arrival is being accounted for — it is pinned to the top
// regardless of its stamp, so the ticket that triggered a trim is never the
// one trimmed. Ties in updated-at are routine (the stamp has second
// resolution) and break toward the newer ULID: among tickets touched in the
// same second, the most recently created counts as newer. keep <= 0 means the
// lane has no cap, so nothing ever overflows. A board with unreadable tickets
// refuses to answer rather than guess at which ticket is oldest.
func (s *Store) Overflow(lane string, keep int, newest string) ([]*Ticket, error) {
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
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].ID == newest {
			return true
		}
		if ts[j].ID == newest {
			return false
		}
		if ts[i].UpdatedAt.Equal(ts[j].UpdatedAt) {
			return ts[i].ID > ts[j].ID
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
// beyond the newest keep leave for the given logbook folder, oldest first.
// newest is passed through to Overflow so the just-landed ticket is never
// trimmed. It returns what moved — including on error, so a caller can still
// report the tickets that had already left before the failure.
func (s *Store) TrimLane(lane string, keep int, folder, newest string) ([]Trimmed, error) {
	over, err := s.Overflow(lane, keep, newest)
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

// StampCommits writes the derived commit union onto the ticket and returns
// what was written. Derived shas come first, in git order; any sha already
// recorded that the derivation did not find is appended rather than dropped —
// a sha a person wrote down deliberately is evidence this tool has no business
// discarding. derive may be nil, the same "no derivation on offer" convention
// core/gate uses.
func (s *Store) StampCommits(t *Ticket, derive func(*Ticket) []string) ([]string, error) {
	var derived []string
	if derive != nil {
		derived = derive(t)
	}
	merged := append([]string{}, derived...)
	for _, c := range t.Commits {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		seen := false
		for _, m := range merged {
			if m == c {
				seen = true
				break
			}
		}
		if !seen {
			merged = append(merged, c)
		}
	}
	if _, err := s.Mutate(t.ID, func(t *Ticket) error {
		return t.Doc().SetList(FieldCommits, merged)
	}); err != nil {
		return nil, err
	}
	return merged, nil
}

// FileLane empties a lane into the given logbook folder, oldest first — the
// lane is a doorway (logbook-on-entry), so the just-landed ticket goes along
// with anything still sitting there from before the doorway existed. prepare
// runs on each ticket before its file moves (the callers stamp commits there);
// nil skips it. A ticket that cannot be filed — unreadable, a stamp failure, a
// name collision in the folder — is skipped and named rather than allowed to
// jam the doorway: a sweep aborts on nothing, because the same failure would
// repeat on every later landing and the lane would block for good. What moved
// and what was skipped are both always reported.
func (s *Store) FileLane(lane, folder string, prepare func(*Ticket) error) ([]Trimmed, error) {
	all, err := s.List()
	var problems []string
	if err != nil {
		var perr *PartialError
		if !errors.As(err, &perr) {
			return nil, err
		}
		// Unlike a cap, a doorway makes no ordering decision an unreadable
		// ticket could corrupt — the readable ones still leave.
		problems = append(problems, perr.Problems...)
	}
	var ts []*Ticket
	for _, t := range all {
		if t.Status == lane {
			ts = append(ts, t)
		}
	}
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].UpdatedAt.Equal(ts[j].UpdatedAt) {
			return ts[i].ID < ts[j].ID
		}
		return ts[i].UpdatedAt.Before(ts[j].UpdatedAt)
	})
	var out []Trimmed
	for _, t := range ts {
		if prepare != nil {
			if err := prepare(t); err != nil {
				problems = append(problems, fmt.Sprintf("%s: %v", Handle(t.ID), err))
				continue
			}
		}
		dst, err := s.Logbook(t.ID, folder)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", Handle(t.ID), err))
			continue
		}
		out = append(out, Trimmed{ID: t.ID, Title: t.Title, Path: dst})
	}
	if len(problems) > 0 {
		return out, &PartialError{Problems: problems}
	}
	return out, nil
}
