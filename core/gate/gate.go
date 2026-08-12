// Package gate holds the rules about when a ticket may move.
//
// These are pure functions over already-loaded state, deliberately separated
// from the store. Two reasons: they can be unit-tested without touching a
// filesystem, and there is exactly one file to change when a rule changes — the
// CLI and the TUI both call in here rather than each carrying their own copy,
// which is the only way "enforced identically by both interfaces" is true rather
// than merely intended.
package gate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/ticket"
)

// Violation codes. Agents branch on these, so they are stable strings rather
// than prose that might be reworded.
const (
	CodeMissingField    = "missing_field"
	CodeBlocked         = "blocked"
	CodeUnknownLane     = "unknown_lane"
	CodeNoSuchLane      = "no_such_lane"
	CodeMissingOutcome  = "missing_outcome"
	CodeNeedsSignal     = "needs_nonmodel_signal"
	CodeNeedsQuestion   = "needs_question"
	CodeSelfBlock       = "self_block"
	CodeMissingProduces = "missing_lane_output"
	CodePlanIncomplete  = "plan_incomplete"
)

// Violation is one reason a move was refused. Field is set when the fix is to
// supply a value, so an agent can act on it without parsing the message.
type Violation struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (v Violation) String() string { return v.Message }

// Violations is an ordered set of refusals.
type Violations []Violation

func (vs Violations) Err() error {
	if len(vs) == 0 {
		return nil
	}
	parts := make([]string, 0, len(vs))
	for _, v := range vs {
		parts = append(parts, "  - "+v.Message)
	}
	return fmt.Errorf("refused:\n%s", strings.Join(parts, "\n"))
}

// Fields lists the field names that need filling in.
func (vs Violations) Fields() []string {
	var out []string
	for _, v := range vs {
		if v.Field != "" {
			out = append(out, v.Field)
		}
	}
	sort.Strings(out)
	return out
}

// PromotionFields are required before a ticket may leave the backlog.
//
// The point of the gate is that capture should be frictionless while execution
// should not be: anything can be thrown into the backlog, but nothing starts
// without a goal, a definition of done, the context it came from, and a human
// who owns the outcome. Acquiring those is the price of admission to a run.
var PromotionFields = []string{
	ticket.FieldGoal,
	ticket.FieldDoD,
	ticket.FieldContext,
	ticket.FieldAssignee,
}

// Request describes a proposed move.
type Request struct {
	To string

	// Question accompanies a move into a lane that requires one.
	Question string

	// NonModelSignal is evidence that something other than a language model
	// judged the work complete — a passing command, or an explicit human
	// sign-off. Required to enter a lane that demands it.
	NonModelSignal string

	// Actor is who is performing the move, used only in messages.
	Actor string
}

// Env is the surrounding state a decision needs.
type Env struct {
	Lanes *lane.Set
	// All is every ticket, used to evaluate blockers. Only id and status are read.
	All []*ticket.Ticket
}

// CheckAdvance decides whether t may move to req.To.
func CheckAdvance(env Env, t *ticket.Ticket, req Request) Violations {
	var vs Violations

	target, ok := env.Lanes.Get(req.To)
	if !ok {
		return append(vs, Violation{
			Code:    CodeNoSuchLane,
			Message: fmt.Sprintf("no lane %q is installed; available lanes: %s", req.To, strings.Join(env.Lanes.IDs(), ", ")),
		})
	}

	// A ticket sitting in a lane this installation does not know about has no
	// contract to enforce, so moving it either way is refused rather than
	// guessed at. Obtaining the lane file is the fix.
	if _, known := env.Lanes.Get(t.Status); !known && t.Status != "" {
		vs = append(vs, Violation{
			Code:    CodeUnknownLane,
			Message: fmt.Sprintf("ticket is in unrecognized lane %q; install that lane file before moving it", t.Status),
		})
		return vs
	}

	from := env.Lanes.Index(t.Status)
	to := env.Lanes.Index(req.To)

	// Leaving the backlog is the moment work becomes real, so the promotion gate
	// applies to any forward move out of it.
	leavingBacklog := t.Status == "backlog" && req.To != "backlog"
	if leavingBacklog {
		vs = append(vs, missingPromotionFields(t)...)
		vs = append(vs, blockedBy(env, t)...)
	}

	// Any move that starts work also respects dependencies, not just the first.
	if to > from && !leavingBacklog {
		vs = append(vs, blockedBy(env, t)...)
	}

	if target.RequiresQuestion && strings.TrimSpace(req.Question) == "" && strings.TrimSpace(t.Question) == "" {
		vs = append(vs, Violation{
			Code:    CodeNeedsQuestion,
			Field:   ticket.FieldQuestion,
			Message: fmt.Sprintf("lane %q needs the question that is blocking progress", target.ID),
		})
	}

	if target.RequiresOutcome && !t.Outcome.Filled() {
		for _, f := range []struct {
			name, val string
		}{
			{ticket.FieldOutcomeWhat, t.Outcome.What},
			{ticket.FieldOutcomeWhy, t.Outcome.Why},
			{ticket.FieldOutcomeResolves, t.Outcome.Resolves},
		} {
			if strings.TrimSpace(f.val) == "" {
				vs = append(vs, Violation{
					Code:    CodeMissingOutcome,
					Field:   f.name,
					Message: fmt.Sprintf("lane %q requires %s to be filled in", target.ID, f.name),
				})
			}
		}
	}

	// A model must not be able to certify its own work as finished. Review agents
	// are measurably poor at catching real defects, so the terminal transition
	// needs evidence that did not come from a model.
	//
	// A fully ticked definition-of-done checklist is exactly that evidence:
	// ticking a box is a human editing a file, which a model asserting "it works"
	// cannot manufacture. So a complete checklist satisfies the requirement, and
	// --signal remains available for tickets whose criteria live elsewhere.
	if target.RequiresNonModelSignal {
		complete, remaining := t.DoDComplete()
		if !complete && strings.TrimSpace(req.NonModelSignal) == "" {
			if len(remaining) > 0 {
				for _, r := range remaining {
					vs = append(vs, Violation{
						Code:    CodeNeedsSignal,
						Field:   ticket.FieldDoD,
						Message: fmt.Sprintf("definition of done is not met: %s", r),
					})
				}
			} else {
				vs = append(vs, Violation{
					Code: CodeNeedsSignal,
					Message: fmt.Sprintf(
						"lane %q cannot be entered on a model's own assessment; tick the definition-of-done checklist, or pass --signal with a passing command or a human sign-off",
						target.ID),
				})
			}
		}
	}

	// The plan and the definition of done are not independent. The plan is the
	// method — write the spec, design it, implement it — and the criteria cannot
	// have been met while the work that meets them is still under way. A ticket
	// accepted with "implement" unfinished is the "I thought it was built, only
	// the design was done" failure the board exists to prevent, so a terminal
	// lane refuses it rather than merely noting it.
	//
	// This applies only at terminal lanes: a plan is expected to be in progress
	// while the ticket moves through the pipeline.
	if target.Terminal {
		if complete, remaining := t.PlanComplete(); !complete {
			for _, r := range remaining {
				vs = append(vs, Violation{
					Code:    CodePlanIncomplete,
					Message: fmt.Sprintf("the plan still has unfinished steps: %s", r),
				})
			}
		}
	}

	// The lane being left must have produced what it promised, so a pipeline
	// step cannot be skipped by simply advancing past it.
	if cur, ok := env.Lanes.Get(t.Status); ok && to > from {
		for _, f := range cur.OutputProduces {
			if !fieldFilled(t, f) {
				vs = append(vs, Violation{
					Code:    CodeMissingProduces,
					Field:   f,
					Message: fmt.Sprintf("lane %q declares it produces %s, which is still empty", cur.ID, f),
				})
			}
		}
	}

	return vs
}

// missingPromotionFields reports which gate fields are absent.
func missingPromotionFields(t *ticket.Ticket) Violations {
	var vs Violations
	for _, f := range PromotionFields {
		if !fieldFilled(t, f) {
			vs = append(vs, Violation{
				Code:    CodeMissingField,
				Field:   f,
				Message: fmt.Sprintf("%s is required before this ticket can leave the backlog", f),
			})
		}
	}
	return vs
}

func fieldFilled(t *ticket.Ticket, field string) bool {
	var v string
	switch field {
	case ticket.FieldGoal:
		v = t.Goal
	case ticket.FieldDoD:
		// Either form counts: a checklist in the body, or a one-line scalar. The
		// checklist is the form humans actually write, so it is checked first.
		if len(t.DoDItems) > 0 {
			return true
		}
		v = t.DoD
	case ticket.FieldContext:
		v = t.Context
	case ticket.FieldAssignee:
		v = t.Assignee
	case ticket.FieldTitle:
		v = t.Title
	case ticket.FieldQuestion:
		v = t.Question
	case ticket.FieldOutcomeWhat:
		v = t.Outcome.What
	case ticket.FieldOutcomeWhy:
		v = t.Outcome.Why
	case ticket.FieldOutcomeResolves:
		v = t.Outcome.Resolves
	case "diff", ticket.FieldCommits:
		return len(t.Commits) > 0
	default:
		s, _, err := t.Doc().Scalar(field)
		if err != nil {
			return false
		}
		v = s
	}
	return strings.TrimSpace(v) != ""
}

// blockedBy reports unresolved dependencies.
func blockedBy(env Env, t *ticket.Ticket) Violations {
	if len(t.BlockedBy) == 0 {
		return nil
	}
	byID := map[string]*ticket.Ticket{}
	for _, o := range env.All {
		byID[o.ID] = o
	}
	var vs Violations
	for _, dep := range t.BlockedBy {
		if dep == t.ID {
			vs = append(vs, Violation{
				Code:    CodeSelfBlock,
				Message: "ticket lists itself in blocked-by",
			})
			continue
		}
		other, ok := byID[dep]
		if !ok {
			vs = append(vs, Violation{
				Code:    CodeBlocked,
				Message: fmt.Sprintf("blocked by %s, which does not exist in this store", ticket.Handle(dep)),
			})
			continue
		}
		if l, ok := env.Lanes.Get(other.Status); !ok || !l.Terminal {
			vs = append(vs, Violation{
				Code:    CodeBlocked,
				Message: fmt.Sprintf("blocked by %s (%s) which is in %q, not a terminal lane", ticket.Handle(dep), other.Title, other.Status),
			})
		}
	}
	return vs
}

// Ready reports whether a ticket satisfies the promotion gate, without
// reference to any particular target lane.
func Ready(t *ticket.Ticket) bool {
	return len(missingPromotionFields(t)) == 0
}

// Actionable reports whether a ticket could be picked up right now.
func Actionable(env Env, t *ticket.Ticket) bool {
	if l, ok := env.Lanes.Get(t.Status); ok && l.Terminal {
		return false
	}
	if !Ready(t) {
		return false
	}
	return len(blockedBy(env, t)) == 0
}
