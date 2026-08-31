// Package validate checks tickets at rest.
//
// The gates in core/gate fire when a ticket moves, which means a ticket damaged
// by a hand edit, a bad merge, or an agent writing something unexpected sits
// there looking fine until someone tries to use it — and then fails with an
// error about the move rather than about the damage. This package answers the
// separate question "is what is on disk actually coherent?".
package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Severities. A ticket can be incomplete without being broken: capture is meant
// to be cheap, so an unspecified backlog ticket is normal and must not fail a
// validation run.
const (
	SeverityError   = "error"
	SeverityWarning = "warning"
)

// Problem codes, stable so an agent can branch on them.
const (
	CodeBadID         = "bad_id"
	CodeNoTitle       = "no_title"
	CodeUnknownLane   = "unknown_lane"
	CodeBadTimestamp  = "bad_timestamp"
	CodeDanglingDep   = "dangling_dependency"
	CodeSelfDep       = "self_dependency"
	CodeDuplicateID   = "duplicate_id"
	CodeIncomplete    = "incomplete"
	CodeUndeclaredDep = "undeclared_dependency"
)

// handleRef matches a bare handle: the six-character tail ticket.Handle
// prints, drawn from ULID's Crockford base32 alphabet (no I, L, O, U). Almost
// no ordinary word survives that alphabet, which is what keeps this from
// firing on prose; the resolve-against-the-store check below is the real
// filter against the rest.
//
// Deliberately blind to two things: a full 26-character ULID typed out in
// prose (the \b boundaries never land six characters in from either end of
// one unbroken run of word characters, so it never matches at all), and a
// lowercase handle (the character class is uppercase only, matching what
// ticket.Handle actually prints). Both are scope, not oversight — the case
// this check exists for is the short handle a human actually types.
var handleRef = regexp.MustCompile(`\b[0-9ABCDEFGHJKMNPQRSTVWXYZ]{6}\b`)

// Problem is one finding about one ticket.
type Problem struct {
	Code     string
	Severity string
	Handle   string
	ID       string
	Path     string
	Field    string
	Message  string
}

// HasErrors reports whether any problem is severe enough to fail a run.
func HasErrors(ps []Problem) bool {
	for _, p := range ps {
		if p.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Tickets validates a whole board. Cross-ticket checks — duplicate ids and
// unresolvable dependencies — need the full set, which is why this takes a slice
// rather than validating one ticket at a time.
func Tickets(ts []*ticket.Ticket, lanes *lane.Set) []Problem {
	var ps []Problem

	byID := make(map[string]int, len(ts))
	byHandle := make(map[string]*ticket.Ticket, len(ts))
	for _, t := range ts {
		if t.ID != "" {
			byID[t.ID]++
		}
		byHandle[handleOf(t.ID)] = t
	}
	reported := map[string]bool{}

	for _, t := range ts {
		add := func(code, severity, field, format string, args ...any) {
			ps = append(ps, Problem{
				Code: code, Severity: severity, Field: field,
				ID: t.ID, Handle: handleOf(t.ID), Path: t.Path,
				Message: fmt.Sprintf(format, args...),
			})
		}

		if !ticket.ValidID(t.ID) {
			add(CodeBadID, SeverityError, ticket.FieldID,
				"id %q is not a ULID; the id is how every other ticket refers to this one", t.ID)
		} else if byID[t.ID] > 1 && !reported[t.ID] {
			reported[t.ID] = true
			add(CodeDuplicateID, SeverityError, ticket.FieldID,
				"%d files claim id %s; only the first is reachable", byID[t.ID], t.ID)
		}

		if strings.TrimSpace(t.Title) == "" {
			add(CodeNoTitle, SeverityError, ticket.FieldTitle,
				"no title, and no level-one heading in the body to fall back to")
		}

		if _, known := lanes.Get(t.Status); !known {
			add(CodeUnknownLane, SeverityError, ticket.FieldStatus,
				"lane %q is not installed, so this ticket is read-only and invisible to the pipeline", t.Status)
		}

		// The merge driver resolves competing edits by comparing updated-at. A
		// ticket without it cannot participate in that, so it is an error rather
		// than a cosmetic omission.
		if t.CreatedAt.IsZero() {
			add(CodeBadTimestamp, SeverityError, ticket.FieldCreatedAt, "created-at is missing or unparseable")
		}
		if t.UpdatedAt.IsZero() {
			add(CodeBadTimestamp, SeverityError, ticket.FieldUpdatedAt,
				"updated-at is missing or unparseable; the merge driver resolves conflicts with it")
		}

		declared := make(map[string]bool, len(t.BlockedBy))
		for _, dep := range t.BlockedBy {
			declared[handleOf(dep)] = true
			switch {
			case dep == t.ID:
				add(CodeSelfDep, SeverityError, ticket.FieldBlockedBy,
					"blocked by itself, which can never be satisfied")
			case byID[dep] == 0:
				add(CodeDanglingDep, SeverityError, ticket.FieldBlockedBy,
					"blocked by %s, which is not on this board; the dependency can never clear", handleOf(dep))
			}
		}

		// A handle typed into prose reads as a dependency to a human but is
		// invisible to the gates, which only read blocked-by. Reported only
		// when the token also resolves to a real, non-terminal ticket, so an
		// ordinary word that happens to share the six-character shape (say,
		// GOLANG) cannot fire this.
		own := handleOf(t.ID)
		follows := handleOf(t.Follows)
		// 'jaira set' replaces a list field outright rather than appending to
		// it, so each suggested command has to spell out the whole resulting
		// list — the ticket's existing blocked-by plus every handle found so
		// far, across both sources — or following an earlier one verbatim
		// would erase a dependency a later one just found.
		full := append([]string(nil), t.BlockedBy...)
		for _, src := range []struct{ name, text, field string }{
			{"context", t.Context, ticket.FieldContext},
			// No frontmatter field constant covers the body, so this
			// problem's Field is left empty rather than inventing one.
			{"note", t.Body, ""},
		} {
			seen := map[string]bool{}
			for _, m := range handleRef.FindAllString(src.text, -1) {
				// A follow-up naming its parent is declaring the relation
				// follows: already exists for, not a hidden dependency — and
				// the parent does not block the follow-up, so the fix this
				// warning suggests (adding it to blocked-by) would be wrong.
				if m == own || declared[m] || seen[m] || (t.Follows != "" && m == follows) {
					continue
				}
				ref, ok := byHandle[m]
				if !ok {
					continue
				}
				if l, known := lanes.Get(ref.Status); known && l.Terminal {
					continue
				}
				seen[m] = true
				full = append(full, ref.ID)
				add(CodeUndeclaredDep, SeverityWarning, src.field,
					"%s names %s which is not in blocked-by — declare it with 'jaira set %s blocked-by=%s' or ignore if it is not a dependency",
					src.name, m, own, strings.Join(full, ","))
			}
		}

		// Incompleteness is expected in the backlog and is reported so a migration
		// or a sweep can see it, but it never fails the run.
		if miss := missing(t); len(miss) > 0 {
			add(CodeIncomplete, SeverityWarning, "",
				"cannot leave the backlog yet; still needs %s", strings.Join(miss, ", "))
		}
	}
	return ps
}

func missing(t *ticket.Ticket) []string {
	var out []string
	if strings.TrimSpace(t.Goal) == "" {
		out = append(out, "goal")
	}
	if !t.HasDoD() {
		out = append(out, "definition-of-done")
	}
	if strings.TrimSpace(t.Context) == "" {
		out = append(out, "context")
	}
	if strings.TrimSpace(t.Assignee) == "" {
		out = append(out, "assignee")
	}
	return out
}

func handleOf(id string) string {
	if !ticket.ValidID(id) {
		return id
	}
	return ticket.Handle(id)
}
