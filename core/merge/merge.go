// Package merge resolves concurrent edits to a ticket field by field.
//
// Git merges text line by line, which is wrong for this format in a specific and
// common way: when two people move the same ticket, both sides rewrite the same
// `status:` line and git reports a conflict on the single most frequent operation
// the whole system performs. Union-merging the file instead would produce two
// `status:` lines — invalid YAML rather than a resolved value.
//
// So the merge is structural. Lists union, `status` resolves by how far along the
// lane sequence each side got, other scalars resolve by recency, and only a
// genuine overlap in prose is escalated to a human — scoped to the field that
// actually conflicts rather than the whole ticket.
//
// The trade-off worth naming: this is the one place jaira does something
// non-obvious to a repository, since git has to be told the driver exists. That
// registration is visible rather than silent, and it buys the property the whole
// storage model depends on.
package merge

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// listFields union rather than pick a side: adding a blocker or recording a
// commit is additive information, and losing either would be a silent data loss.
var listFields = map[string]bool{
	ticket.FieldBlockedBy: true,
	ticket.FieldCommits:   true,
}

// proseFields carry human writing. Two different rewrites of the same sentence
// cannot be reconciled by a rule, so these are the only fields that can produce
// a real conflict.
var proseFields = map[string]bool{
	ticket.FieldTitle:           true,
	ticket.FieldGoal:            true,
	ticket.FieldContext:         true,
	ticket.FieldDoD:             true,
	ticket.FieldQuestion:        true,
	ticket.FieldOutcomeWhat:     true,
	ticket.FieldOutcomeWhy:      true,
	ticket.FieldOutcomeResolves: true,
}

// Conflict describes one unresolvable field.
type Conflict struct {
	Field  string
	Ours   string
	Theirs string
	Base   string
}

// Result is the outcome of a merge.
type Result struct {
	// Merged is the resulting file content.
	Merged []byte
	// Conflicts is empty on a clean merge.
	Conflicts []Conflict
	// Notes explains automatic resolutions, for the caller to log.
	Notes []string
}

// Clean reports whether the merge needs no human attention.
func (r *Result) Clean() bool { return len(r.Conflicts) == 0 }

// Merge performs a three-way merge of one ticket file.
//
// ours is used as the structural base for the output so its formatting, comments
// and key order survive; only resolved values are spliced in.
func Merge(base, ours, theirs []byte, lanes *lane.Set) (*Result, error) {
	bd, err := ticket.ParseDoc(base)
	if err != nil {
		return nil, fmt.Errorf("merge: base: %w", err)
	}
	od, err := ticket.ParseDoc(ours)
	if err != nil {
		return nil, fmt.Errorf("merge: ours: %w", err)
	}
	td, err := ticket.ParseDoc(theirs)
	if err != nil {
		return nil, fmt.Errorf("merge: theirs: %w", err)
	}

	res := &Result{}

	// Recency is decided once, from each side's own updated-at, and reused for
	// every contested scalar so the result is internally consistent rather than
	// a mix of both sides.
	ourTime := stampOf(od)
	theirTime := stampOf(td)
	theirsIsNewer := theirTime.After(ourTime)

	for _, key := range unionKeys(bd, od, td) {
		switch {
		case listFields[key]:
			if err := mergeList(key, bd, od, td, res); err != nil {
				return nil, err
			}
		case key == ticket.FieldStatus:
			if err := mergeStatus(bd, od, td, lanes, res); err != nil {
				return nil, err
			}
		case key == ticket.FieldUpdatedAt:
			// The later stamp wins, so the merged file reflects the most recent
			// edit rather than whichever side git happened to call "ours".
			newest := ourTime
			if theirTime.After(newest) {
				newest = theirTime
			}
			if !newest.IsZero() {
				if err := od.SetRaw(ticket.FieldUpdatedAt, ticket.FormatTime(newest)); err != nil {
					return nil, err
				}
			}
		case key == ticket.FieldExternal:
			// Never interpreted, so never merged: whichever side changed it wins,
			// and if both did it is a conflict the tool has no basis to resolve.
			if err := mergeOpaque(key, bd, od, td, res, theirsIsNewer); err != nil {
				return nil, err
			}
		default:
			if err := mergeScalar(key, bd, od, td, res, theirsIsNewer); err != nil {
				return nil, err
			}
		}
	}

	// The markdown body is ordinary prose; the same three-way rule applies.
	bb, ob, tb := bd.Body(), od.Body(), td.Body()
	if ob != tb {
		switch {
		case ob == bb:
			od.SetBody(tb)
		case tb == bb:
			// keep ours
		default:
			res.Conflicts = append(res.Conflicts, Conflict{Field: "body", Ours: ob, Theirs: tb, Base: bb})
		}
	}

	// Record any unresolved fields inside the file itself, as valid YAML rather
	// than as conflict markers. Markers would make the frontmatter unparseable,
	// which would blank the ticket on everyone's board until someone resolved it
	// — a worse outcome than a parseable ticket carrying a visible note that two
	// versions of one field are outstanding.
	if len(res.Conflicts) > 0 {
		var names []string
		for _, c := range res.Conflicts {
			names = append(names, c.Field)
		}
		sort.Strings(names)
		if err := od.SetList(FieldConflicts, names); err != nil {
			return nil, err
		}
		for _, c := range res.Conflicts {
			if c.Field == "body" {
				continue
			}
			if err := od.SetScalar(theirsKey(c.Field), c.Theirs); err != nil {
				return nil, err
			}
		}
	}

	res.Merged = od.Bytes()
	return res, nil
}

// FieldConflicts names the frontmatter key listing fields left unresolved.
const FieldConflicts = "merge-conflicts"

// theirsKey is where the losing side of a conflict is parked so nothing is lost
// and a human can see both versions.
func theirsKey(field string) string { return "conflict-theirs-" + field }

func stampOf(d *ticket.Doc) time.Time {
	v, _, err := d.Scalar(ticket.FieldUpdatedAt)
	if err != nil || v == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t
		}
	}
	return time.Time{}
}

// unionKeys lists every top-level key appearing on any side, in a stable order
// that follows ours where possible so output ordering is predictable.
func unionKeys(bd, od, td *ticket.Doc) []string {
	seen := map[string]bool{}
	var out []string
	for _, d := range []*ticket.Doc{od, td, bd} {
		for _, k := range d.Keys() {
			if !seen[k] {
				seen[k] = true
				out = append(out, k)
			}
		}
	}
	return out
}

func scalarOf(d *ticket.Doc, key string) string {
	v, _, err := d.Scalar(key)
	if err != nil {
		return ""
	}
	return v
}

func mergeScalar(key string, bd, od, td *ticket.Doc, res *Result, theirsIsNewer bool) error {
	b, o, t := scalarOf(bd, key), scalarOf(od, key), scalarOf(td, key)
	if o == t {
		return nil // both sides agree; ours already holds it
	}
	switch {
	case o == b:
		// Only they changed it.
		return set(od, key, t)
	case t == b:
		return nil // only we changed it
	}
	// Both changed it differently.
	if proseFields[key] {
		res.Conflicts = append(res.Conflicts, Conflict{Field: key, Ours: o, Theirs: t, Base: b})
		return nil
	}
	if theirsIsNewer {
		res.Notes = append(res.Notes, fmt.Sprintf("%s: took the newer side (%q over %q)", key, t, o))
		return set(od, key, t)
	}
	res.Notes = append(res.Notes, fmt.Sprintf("%s: kept the newer side (%q over %q)", key, o, t))
	return nil
}

func mergeOpaque(key string, bd, od, td *ticket.Doc, res *Result, theirsIsNewer bool) error {
	// Opaque blocks are compared as raw text since their shape is unknown.
	b, o, t := rawBlock(bd, key), rawBlock(od, key), rawBlock(td, key)
	if o == t || t == b {
		return nil
	}
	if o == b {
		res.Conflicts = append(res.Conflicts, Conflict{
			Field: key, Ours: o, Theirs: t, Base: b,
		})
		res.Notes = append(res.Notes,
			key+": changed on the other side only, but jaira cannot rewrite an opaque block safely")
		return nil
	}
	res.Conflicts = append(res.Conflicts, Conflict{Field: key, Ours: o, Theirs: t, Base: b})
	return nil
}

// rawBlock extracts the text of a nested block so it can be compared without
// being interpreted.
func rawBlock(d *ticket.Doc, key string) string {
	fm := d.Frontmatter()
	lines := strings.Split(fm, "\n")
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, key+":") {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	end := start + 1
	for end < len(lines) {
		l := lines[end]
		if l == "" || strings.HasPrefix(l, " ") || strings.HasPrefix(l, "\t") || strings.HasPrefix(l, "-") {
			end++
			continue
		}
		break
	}
	return strings.Join(lines[start:end], "\n")
}

func mergeList(key string, bd, od, td *ticket.Doc, res *Result) error {
	o, _ := od.List(key)
	t, _ := td.List(key)
	b, _ := bd.List(key)

	// Union, minus anything deliberately removed on exactly one side. A removal
	// is only honoured when the other side did not also change the list, so a
	// concurrent addition is never silently discarded.
	removedByUs := diff(b, o)
	removedByThem := diff(b, t)

	set := map[string]bool{}
	for _, v := range append(append([]string{}, o...), t...) {
		set[v] = true
	}
	for _, v := range removedByUs {
		if !contains(t, v) || contains(removedByThem, v) {
			delete(set, v)
		}
	}
	for _, v := range removedByThem {
		if !contains(o, v) || contains(removedByUs, v) {
			delete(set, v)
		}
	}

	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)

	if equal(out, o) {
		return nil
	}
	if added := diff(o, out); len(added) > 0 {
		res.Notes = append(res.Notes, fmt.Sprintf("%s: unioned, adding %s", key, strings.Join(added, ", ")))
	}
	return od.SetList(key, out)
}

// mergeStatus resolves a lane collision by progress rather than by recency.
//
// Reverting forward progress is a worse failure than keeping a slightly
// optimistic lane: if one side moved a ticket to review and another moved it back
// to todo, silently discarding the review would lose real work.
func mergeStatus(bd, od, td *ticket.Doc, lanes *lane.Set, res *Result) error {
	b, o, t := scalarOf(bd, ticket.FieldStatus), scalarOf(od, ticket.FieldStatus), scalarOf(td, ticket.FieldStatus)
	if o == t {
		return nil
	}
	if o == b {
		return set(od, ticket.FieldStatus, t)
	}
	if t == b {
		return nil
	}
	po, pt := lanes.Precedence(o), lanes.Precedence(t)
	if pt > po {
		res.Notes = append(res.Notes, fmt.Sprintf("status: %q is further along than %q", t, o))
		return set(od, ticket.FieldStatus, t)
	}
	res.Notes = append(res.Notes, fmt.Sprintf("status: kept %q, which is at least as far along as %q", o, t))
	return nil
}

func set(d *ticket.Doc, key, val string) error {
	if val == "" {
		// Writing an empty scalar is preferable to deleting the key: a missing
		// field and an empty field mean different things to the gates.
		return d.SetScalar(key, "")
	}
	return d.SetScalar(key, val)
}

func diff(from, to []string) []string {
	in := map[string]bool{}
	for _, v := range to {
		in[v] = true
	}
	var out []string
	for _, v := range from {
		if !in[v] {
			out = append(out, v)
		}
	}
	return out
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// RenderConflicts describes what a human needs to resolve, scoped to the fields
// that actually clashed rather than dumping the whole file.
func RenderConflicts(cs []Conflict) string {
	var b strings.Builder
	for _, c := range cs {
		fmt.Fprintf(&b, "field %s changed on both sides:\n", c.Field)
		fmt.Fprintf(&b, "  ours:   %s\n", oneLine(c.Ours))
		fmt.Fprintf(&b, "  theirs: %s\n", oneLine(c.Theirs))
		if c.Base != "" {
			fmt.Fprintf(&b, "  was:    %s\n", oneLine(c.Base))
		}
	}
	return b.String()
}

func oneLine(s string) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\n", " ⏎ ")
	if len(s) > 120 {
		s = s[:119] + "…"
	}
	if s == "" {
		return "(empty)"
	}
	return s
}
