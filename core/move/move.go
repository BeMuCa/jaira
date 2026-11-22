// Package move is the one place a ticket changes lanes. Four write-sites grew
// up around the board — the CLI's move, the TUI's gated and forced moves, and
// the sign-off accept — and every rule that had to run "after a move" needed
// all four found by hand; a review caught the fourth only by searching. The
// gate check, the claim-on-pull, the status write and the lane's settle now
// run here, and the callers keep what is genuinely theirs: flags, staging,
// wording, and the pending/confirm dance.
package move

import (
	"slices"
	"strconv"
	"strings"
	"time"

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
	// ByPerson marks a move a person made in the board, whatever lane it
	// leaves: a send-back is recorded as theirs (back-by: person) and never
	// counts toward a loop lane's third return, and the gate never refuses
	// it for a lane's max-rounds. Unlike Interactive it opens no other gate.
	ByPerson bool
	// Force writes despite refusals; what was overridden is returned, the
	// caller says it out loud. Nothing marks the ticket itself (deliberate).
	Force bool
	// Question and Reason ride into the gate check: entering a question lane
	// or a parking lane is judged on what the move brings along.
	Question string
	Reason   string
	// Note is appended to the ticket's progress notes in the same write as
	// the move, so a person's send-back and the reason for it land together
	// or not at all. Empty writes none.
	Note string
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
	// DoDUnticked is the size of a definition of done not one item of which
	// is ticked, when the move left the implementing lane forward; 0 when
	// there is nothing to warn about. The move lands regardless — the caller
	// says it out loud. See untickedDoD.
	DoDUnticked int
}

// Moved, when set, is called after every move that wrote: the CLI's, the
// TUI's and any later caller's alike, so a plugin subscribed to moves hears
// about each of them however it was made. from is the lane the ticket left.
// It is set once at startup by the binary (internal/cli), which owns what a
// plugin is told; a nil Moved calls nothing.
var Moved func(s *ticket.Store, t *ticket.Ticket, from string, lanes *lane.Set)

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
		Interactive: req.Interactive, ByPerson: req.ByPerson,
	})
	if len(vs) > 0 && !req.Force {
		return &Result{Ticket: t, Refused: vs}, nil
	}

	from := t.Status
	t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
		if claiming {
			if err := t.Doc().SetScalar(ticket.FieldAssignee, req.Actor); err != nil {
				return err
			}
			t.Assignee = req.Actor
		}
		if env.Lanes.SendsBack(t.Status, req.To) {
			if err := t.Doc().SetRaw(ticket.FieldRound, strconv.Itoa(ticket.CurrentRound(t)+1)); err != nil {
				return err
			}
			if err := recordBack(t, env.Lanes, t.Status, req.Interactive || req.ByPerson); err != nil {
				return err
			}
		}
		if err := t.Doc().SetScalar(ticket.FieldStatus, req.To); err != nil {
			return err
		}
		if note := strings.TrimSpace(req.Note); note != "" {
			if err := ticket.AppendNote(t.Doc(), note, time.Now(), req.Actor); err != nil {
				return err
			}
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	})
	if err != nil {
		return nil, err
	}

	res := &Result{Ticket: t, Claimed: claiming, DoDUnticked: untickedDoD(env.Lanes, t, from, req.To)}
	if req.Force && len(vs) > 0 {
		res.Overrode = vs
	}
	if l, ok := env.Lanes.Get(req.To); ok {
		res.Holds = l.Holds
		res.Trimmed, res.Filed, res.SettleErr = lane.Settle(s, l, req.Folder, t.ID, req.Prepare)
	}
	if Moved != nil {
		Moved(s, t, from, env.Lanes)
	}
	return res, nil
}

// untickedDoD returns the number of definition-of-done items when a move out
// of the implementing lane goes on with none of them ticked, and 0 otherwise.
//
// The terminal lane's gate reads the definition of done, not the plan, and an
// agent once carried a ticket through eight lanes ticking only the plan: the
// nine refusals arrived at once, at signoff, when finding out cost the most.
// The gate stays as strict as it is; this only moves the news to the step
// where the work was supposed to tick them.
//
// The implementing lane is the one whose output names outcome-what. Only a
// move forward counts — a send-back or a parked ticket has not left the work
// behind it. Ticked means settled: a superseded item satisfies the gate too.
func untickedDoD(lanes *lane.Set, t *ticket.Ticket, from, to string) int {
	l, ok := lanes.Get(from)
	if !ok || !slices.Contains(l.OutputProduces, ticket.FieldOutcomeWhat) || len(t.DoDItems) == 0 {
		return 0
	}
	dst, ok := lanes.Get(to)
	if !ok || dst.RequiresBlockedReason || lanes.Index(to) <= lanes.Index(from) {
		return 0
	}
	for _, it := range t.DoDItems {
		if it.Settled() {
			return 0
		}
	}
	return len(t.DoDItems)
}

// recordBack writes who sent the ticket back: the lane it leaves, whether a
// person or a lane sent it — a person's lane, or a person at the board
// (byPerson) — and how many returns in a row came from that lane. A person's return, or one from another lane, starts the count again:
// the third-return stop is about a loop lane that keeps finding more, and a
// person's wish is not a loop. Written by the move itself, so every path that
// sends work back — the CLI, the cockpit, --force — records it the same way.
func recordBack(t *ticket.Ticket, lanes *lane.Set, from string, byPerson bool) error {
	by := ticket.BackByLane
	if l, ok := lanes.Get(from); byPerson || !ok || isPersons(l) {
		by = ticket.BackByPerson
	}
	n := 1
	if by == ticket.BackByLane && t.BackFrom == from && t.BackBy == ticket.BackByLane {
		n = t.BackCount + 1
	}
	d := t.Doc()
	if err := d.SetScalar(ticket.FieldBackFrom, from); err != nil {
		return err
	}
	if err := d.SetScalar(ticket.FieldBackBy, by); err != nil {
		return err
	}
	return d.SetRaw(ticket.FieldBackCount, strconv.Itoa(n))
}

// isPersons reports a lane whose work a person does: one no agent steps
// through, or one only a person may leave or enter with a question.
func isPersons(l *lane.Lane) bool {
	return !l.Agentic || l.RequiresHumanExit || l.RequiresQuestion
}
