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
	"slices"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/tag"
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
	CodeBadTag        = "bad_tag"
	CodeBadMode       = "bad_mode"
	// CodeDanglingParent, CodeSelfParent and CodeParentCycle guard the parent
	// chain. It is followed recursively wherever children are rendered, so a
	// ring in it is a hung view rather than a wrong one — which is why a
	// cycle is an error and not a warning.
	CodeDanglingParent = "dangling_parent"
	CodeSelfParent     = "self_parent"
	CodeParentCycle    = "parent_cycle"
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
// Tickets checks every ticket it is given.
//
// known answers whether an id exists anywhere this repository can see it —
// on the board, on a ref, or filed away in the logbook or the archive. It is
// injected because finding that out means reading directories and this
// package is pure. nil means only the given tickets are known, which is the
// right answer for a caller that has no store and the reason a filed-away
// blocker used to be reported as unreachable.
func Tickets(ts []*ticket.Ticket, lanes *lane.Set, known func(id string) bool) []Problem {
	if known == nil {
		known = func(string) bool { return false }
	}
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
	byParent := make(map[string]*ticket.Ticket, len(ts))
	for _, t := range ts {
		if t.ID != "" {
			byParent[t.ID] = t
		}
	}

	for _, t := range ts {
		add := func(code, severity, field, format string, args ...any) {
			where := t.Path
			if where == "" {
				// A ticket the board can see but has no file for. Naming that
				// instead of leaving the path blank is the difference between
				// "go and fix this file" and "this is not yours to fix here".
				where = "(on its ref; pull it to work on it)"
			}
			ps = append(ps, Problem{
				Code: code, Severity: severity, Field: field,
				ID: t.ID, Handle: handleOf(t.ID), Path: where,
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

		// A tag that cannot be normalised is on the ticket and reachable by
		// nothing: 'jaira tags' skips it when counting, and --tag and the board
		// filter both compare normalised names, so it is a label only a reader
		// of the raw file ever sees. The same shape as the dangling dependency
		// above — the value is written, the thing it points at is not there.
		//
		// A case difference is not this problem: "UI" normalises to "ui", so the
		// listing counts it and the filters match it. Only a name Normalize
		// refuses outright is inert.
		//
		// A warning, not an error: the ticket itself is intact, and capture must
		// stay cheap.
		for _, raw := range t.Tags {
			_, _, badTag := tag.Normalize(raw)
			if badTag == nil {
				continue
			}
			msg := "tag %q cannot be stored, so nothing can filter on it: %v"
			args := []any{raw, badTag}
			if fixed, ok := tag.Suggest(raw); ok {
				msg += "; rename it with 'jaira set %s tags=%s'"
				args = append(args, handleOf(t.ID),
					strings.Join(normalizedTags(t.Tags, raw, fixed), ","))
			}
			add(CodeBadTag, SeverityWarning, ticket.FieldTags, msg, args...)
		}

		// mode takes a closed set of values, and the two CLI write paths
		// already refuse anything else. This check is for the three ways that
		// go past them — a hand-edited file, a bad merge, an agent writing
		// something unexpected — because a value outside the set is not read
		// as "no mode" by a reader: it prints in the 'mode' row and in the
		// JSON key, so a person can believe they are in the conversational
		// mode while the worker, comparing against exactly one word, runs on
		// autonomously and commits.
		//
		// A warning, not an error, for the same reason as the bad tag above:
		// the ticket itself is intact.
		// The canonical form is compared, not only the verdict: CanonicalMode
		// accepts " conversational " and trims it, so a value that only
		// differs in its whitespace passes the verdict while t.Mode reaches
		// the JSON key and the header untrimmed — the same silent failure as
		// "chat", because the worker compares against exactly one word.
		if canon, ok := ticket.CanonicalMode(t.Mode); !ok || canon != t.Mode {
			add(CodeBadMode, SeverityWarning, ticket.FieldMode,
				"mode %q is not a mode: an agent reads it as no mode at all and runs autonomously; set it with 'jaira set %s mode=%s' or clear it with 'mode='",
				t.Mode, handleOf(t.ID), ticket.ModeConversational)
		}

		declared := make(map[string]bool, len(t.BlockedBy))
		for _, dep := range t.BlockedBy {
			declared[handleOf(dep)] = true
			switch {
			case dep == t.ID:
				add(CodeSelfDep, SeverityError, ticket.FieldBlockedBy,
					"blocked by itself, which can never be satisfied")
			case byID[dep] == 0 && !known(dep):
				add(CodeDanglingDep, SeverityError, ticket.FieldBlockedBy,
					"blocked by %s, which exists nowhere here; the dependency can never clear", handleOf(dep))
			}
		}

		// A handle typed into prose reads as a dependency to a human but is
		// invisible to the gates, which only read blocked-by. Reported only
		// when the token also resolves to a real, non-terminal ticket, so an
		// ordinary word that happens to share the six-character shape (say,
		// GOLANG) cannot fire this.
		own := handleOf(t.ID)
		// Relations this ticket already declares. A handle in the prose that
		// names one of them is the relation being explained, not a
		// dependency somebody forgot: telling a child to put its epic in
		// blocked-by would contradict the rule that a parent is not a gate.
		stated := map[string]bool{}
		if t.Follows != "" {
			stated[handleOf(t.Follows)] = true
		}
		if t.Parent != "" {
			stated[handleOf(t.Parent)] = true
		}
		for _, r := range t.Related {
			stated[handleOf(r)] = true
		}
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
				if m == own || declared[m] || seen[m] || stated[m] {
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
				// The same handle can turn up in both context and note; only
				// list it once in the suggested command, or following the
				// note's warning verbatim after the context's would write a
				// duplicate blocked-by entry.
				if !slices.Contains(full, ref.ID) {
					full = append(full, ref.ID)
				}
				add(CodeUndeclaredDep, SeverityWarning, src.field,
					"%s names %s which is not in blocked-by — declare it with 'jaira set %s blocked-by=%s' or ignore if it is not a dependency",
					src.name, m, own, strings.Join(full, ","))
			}
		}

		switch {
		case t.Parent == "":
		case t.Parent == t.ID:
			add(CodeSelfParent, SeverityError, ticket.FieldParent,
				"is its own parent")
		case byID[t.Parent] == 0 && !known(t.Parent):
			add(CodeDanglingParent, SeverityError, ticket.FieldParent,
				"part of %s, which exists nowhere here", handleOf(t.Parent))
		default:
			if ring := parentCycle(t, byParent, byHandle); len(ring) > 0 {
				add(CodeParentCycle, SeverityError, ticket.FieldParent,
					"parent chain runs in a circle: %s", strings.Join(ring, " → "))
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

// normalizedTags is the ticket's whole tag list with one member replaced by its
// normalised form, so the suggested command can be followed verbatim.
//
// 'jaira set' replaces a list field outright rather than appending to it, so a
// suggestion naming only the offending tag would erase every other tag on the
// ticket — the same trap the blocked-by suggestion above spells the full list to
// avoid. A member that cannot normalise at all is dropped from the suggestion:
// there is nothing to put in its place.
func normalizedTags(all []string, offender, replacement string) []string {
	out := make([]string, 0, len(all))
	for _, raw := range all {
		if raw == offender {
			out = append(out, replacement)
			continue
		}
		if n, _, err := tag.Normalize(raw); err == nil {
			out = append(out, n)
		}
	}
	return out
}

// parentCycle walks a ticket's parent chain and returns the handles on it
// once it comes back to something already seen. A parent outside the given
// set ends the walk: this package only sees the tickets it was handed, and
// inventing a verdict about a file it cannot read would be worse than
// staying quiet.
func parentCycle(t *ticket.Ticket, byID, byHandle map[string]*ticket.Ticket) []string {
	seen := map[string]bool{}
	var chain []string
	for cur := t; cur != nil; {
		if seen[cur.ID] {
			return append(chain, handleOf(cur.ID))
		}
		seen[cur.ID] = true
		chain = append(chain, handleOf(cur.ID))
		if cur.Parent == "" {
			return nil
		}
		next, ok := byID[cur.Parent]
		if !ok && !ticket.ValidID(cur.Parent) {
			// The file is hand-editable and a handle is what a person reads
			// off the board, so a chain written in handles is a chain that
			// exists. Only a reference that is not already a full id is
			// widened this way: a full id that names nothing here is a
			// dangling parent, reported as such, and resolving it to
			// whichever ticket happens to share its last six characters
			// would invent a ring that is not there.
			next = byHandle[handleOf(cur.Parent)]
		}
		cur = next
	}
	return nil
}
