// Package move is the one place a ticket changes lanes. Four write-sites grew
// up around the board — the CLI's move, the TUI's gated and forced moves, and
// the sign-off accept — and every rule that had to run "after a move" needed
// all four found by hand; a review caught the fourth only by searching. The
// gate check, the claim-on-pull, the status write and the lane's settle now
// run here, and the callers keep what is genuinely theirs: flags, staging,
// wording, and the pending/confirm dance.
package move

import (
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Request carries one lane change. The zero value of every optional field
// means "not this caller's behaviour": Stage is the CLI's persisted pre-gate
// write, ClaimOnPull is the TUI's in-memory claim, Reload refreshes the env
// after staging changed the board.
type Request struct {
	To           string
	Actor        string
	ActorAliases []string
	// Interactive marks a human pressing a key — the one thing an agent
	// cannot claim, and the signal requires-human-exit listens for.
	Interactive bool
	// Force writes despite refusals; what was overridden is returned, the
	// caller says it out loud. Nothing marks the ticket itself (deliberate).
	Force bool
	// Question and Reason ride into the gate check: entering a question lane
	// or a parking lane is judged on what the move brings along.
	Question string
	Reason   string
	// Stage is applied and PERSISTED before the gate runs — the CLI stages
	// outcome, question and commits this way, and a refused move keeps them:
	// requiring two commands to satisfy one move would be a poor contract for
	// an agent. nil skips it.
	Stage func(*ticket.Ticket) error
	// Reload re-reads the gate env after Stage wrote to the board. nil keeps
	// the env the caller passed.
	Reload func() (gate.Env, error)
	// ClaimOnPull stages assignee=Actor in memory for the gate and writes it
	// in the move when the ticket is unassigned — capture belongs to nobody,
	// the pull claims (invariant 12). The CLI claims inside Stage instead.
	ClaimOnPull bool
	// Folder is the logbook folder a settle files into (initials-date).
	Folder string
	// Prepare runs on each ticket a doorway files — the callers stamp commits
	// there. nil skips it.
	Prepare func(*ticket.Ticket) error
}

// Result reports what happened. Refused non-empty means nothing was written
// (and Force was off); everything else describes the write that landed.
type Result struct {
	Ticket   *ticket.Ticket
	Refused  gate.Violations
	Overrode gate.Violations
	Claimed  bool
	// The lane's settle, straight from lane.Settle: what left for the
	// logbook, whether the doorway (true) or the cap ran, the cap's size for
	// the caller's wording, and the never-silent partial error.
	Trimmed   []ticket.Trimmed
	Filed     bool
	Holds     int
	SettleErr error
}

// Move performs one lane change: optional persisted staging, the gate, the
// write, the settle — in that order, the same order every caller ran alone.
func Move(s *ticket.Store, env gate.Env, id string, req Request) (*Result, error) {
	t, err := s.Load(id)
	if err != nil {
		return nil, err
	}
	if req.Stage != nil {
		if _, err := s.Mutate(t.ID, req.Stage); err != nil {
			return nil, err
		}
		if t, err = s.Load(t.ID); err != nil {
			return nil, err
		}
		if req.Reload != nil {
			if env, err = req.Reload(); err != nil {
				return nil, err
			}
		}
	}

	// Before the gate, not only in the write: the gate would otherwise judge
	// the stale file and list its gaps as the ticket's — gaps --force then
	// overrides, writing the stale file onto the ref after all.
	if err := s.CheckFresh(t); err != nil {
		return nil, err
	}

	claiming := req.ClaimOnPull && strings.TrimSpace(t.Assignee) == ""
	gated := t
	if claiming {
		// The gate judges the pull as if the claim had happened — staged in
		// memory, written only if the move lands.
		copy := *t
		copy.Assignee = req.Actor
		gated = &copy
	}
	vs := gate.CheckAdvance(env, gated, gate.Request{
		To: req.To, Question: req.Question, Reason: req.Reason,
		Actor: req.Actor, ActorAliases: req.ActorAliases,
		Interactive: req.Interactive,
	})
	if len(vs) > 0 && !req.Force {
		return &Result{Ticket: t, Refused: vs}, nil
	}

	t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
		if claiming {
			if err := t.Doc().SetScalar(ticket.FieldAssignee, req.Actor); err != nil {
				return err
			}
			t.Assignee = req.Actor
		}
		if sendsBack(env.Lanes, t.Status, req.To) {
			if err := t.Doc().SetScalar(ticket.FieldRound, strconv.Itoa(ticket.CurrentRound(t)+1)); err != nil {
				return err
			}
		}
		if err := t.Doc().SetScalar(ticket.FieldStatus, req.To); err != nil {
			return err
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	})
	if err != nil {
		return nil, err
	}

	res := &Result{Ticket: t, Claimed: claiming}
	if req.Force && len(vs) > 0 {
		res.Overrode = vs
	}
	if l, ok := env.Lanes.Get(req.To); ok {
		res.Holds = l.Holds
		res.Trimmed, res.Filed, res.SettleErr = lane.Settle(s, l, req.Folder, t.ID, req.Prepare)
	}
	return res, nil
}

// sendsBack reports a move to an earlier lane of the board, the one that
// starts a new round. Parking does not: a ticket going to blocked and coming
// back from it is the same round, whichever side of the work lane blocked is
// drawn on. A lane this installation does not know counts nothing.
func sendsBack(lanes *lane.Set, from, to string) bool {
	f, t := lanes.Index(from), lanes.Index(to)
	if f < 0 || t < 0 || t >= f {
		return false
	}
	for _, id := range []string{from, to} {
		if l, ok := lanes.Get(id); ok && l.RequiresBlockedReason {
			return false
		}
	}
	return true
}
