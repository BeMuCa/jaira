package milestone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// An older build keeps writing a milestone under its old name — the old
// <name>.md and the old name ref — after this build moved it to its ID. An
// update never makes an existing edit disappear, so what that build wrote is
// merged into the ID milestone instead of being ignored. The merge needs the
// state the two sides last agreed on: the Base, recorded when the old file
// moved, local only and never committed.

// Base is the old side of a migrated milestone as it stood when it was last
// merged: the migrated content, and the tip of the old name ref when there was
// one. A later old-side state is compared with Content; OldSHA only spares the
// comparison when the ref has not moved.
type Base struct {
	Content string `json:"content"`
	OldSHA  string `json:"old-sha,omitempty"`
}

// basesDir is where this board's bases are kept: per working tree, beside the
// outbox, because the milestone files they describe are per working tree.
func basesDir(root string) string {
	return filepath.Join((&ticket.Store{Root: root}).StateDir(), "milestone-bases")
}

// LoadBase reads the base recorded for the milestone with this ID.
func LoadBase(root, id string) (Base, bool) {
	b, err := os.ReadFile(filepath.Join(basesDir(root), id+".json"))
	if err != nil {
		return Base{}, false
	}
	var base Base
	if json.Unmarshal(b, &base) != nil {
		return Base{}, false
	}
	return base, true
}

// SaveBase records the base for the milestone with this ID. Best effort for
// its callers: without a base the old side is simply not adopted.
func SaveBase(root, id string, base Base) error {
	if err := os.MkdirAll(basesDir(root), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(base)
	if err != nil {
		return err
	}
	return ticket.WriteAtomic(filepath.Join(basesDir(root), id+".json"), append(b, '\n'))
}

// Conflict is a field the older build and this one both changed since the
// base. This side is kept; Old is what the older build wrote, so it can be
// told rather than dropped.
type Conflict struct {
	Milestone string `json:"milestone"`
	Field     string `json:"field"`
	Old       string `json:"old"`
}

func (c Conflict) String() string {
	return fmt.Sprintf("%s: %s was changed by an older build too; kept as it is here, the older build wrote %s", c.Milestone, c.Field, c.Old)
}

// Adopt merges the old side into the ID milestone. base is what the old side
// looked like when last merged, old what it is now, cur the ID milestone's
// file; all are the milestones on the board, for spelling names as IDs.
//
//   - tickets added or removed on the old side since base are added to or
//     removed from cur;
//   - a frontmatter field (goal, after, any other) or the body that only the
//     old side changed is taken; after: names are spelled as IDs;
//   - one that both sides changed keeps cur, is returned as a Conflict, and is
//     written into the file as a comment, because a fetch that runs in the
//     background has nobody to tell.
//
// The result equals cur byte for byte when there is nothing to adopt.
func Adopt(name string, base, old, cur []byte, all []*Milestone) ([]byte, []Conflict) {
	b, o, c := parse(string(base)), parse(string(old)), parse(string(cur))
	res := parse(string(cur))
	bf, of, cf := fields(b), fields(o), fields(c)
	var conflicts []Conflict

	keys := slices.Clone(cf.keys)
	for _, k := range of.keys {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		bv, bok := bf.val(k)
		ov, ook := of.val(k)
		cv, cok := cf.val(k)
		if ook == bok && ov == bv {
			continue // the old side left it alone
		}
		if cok == ook && cv == ov {
			continue // both sides made the same change
		}
		if cok != bok || cv != bv {
			conflicts = append(conflicts, Conflict{Milestone: name, Field: k, Old: quoteOld(ov, ook)})
			continue
		}
		switch {
		case !ook:
			res.setFront(k, "", true)
		case k == "after":
			res.SetAfter(SpellAsIDs(all, parseAfter(of.raw[k])))
		default:
			res.setFront(k, of.raw[k], false)
		}
	}

	notes := conflictNotes(res)
	bp, op, cp := prose(b), prose(o), prose(c)
	if bp != op && cp != op {
		if cp != bp {
			conflicts = append(conflicts, Conflict{Milestone: name, Field: "body", Old: quoteOld(op, true)})
		} else {
			// The body comes with the comments the old side carries; this
			// side's own comments stay, one line per field.
			res = withBody(res, o)
			theirs := conflictNotes(res)
			res.lines = slices.DeleteFunc(res.lines, isConflictNote)
			for _, l := range append(slices.Clone(notes), theirs...) {
				res.addConflictNote(l)
			}
		}
	}

	target := slices.Clone(c.members)
	for _, id := range b.members {
		if !o.Has(id) {
			target = slices.DeleteFunc(target, func(m string) bool { return m == id })
		}
	}
	for _, id := range o.members {
		if !b.Has(id) && !slices.Contains(target, id) {
			target = append(target, id)
		}
	}
	for _, id := range res.Members() {
		if !slices.Contains(target, id) {
			res.Remove(id)
		}
	}
	for _, id := range target {
		res.Add(id)
	}

	for _, k := range conflicts {
		res.setConflictNote(k)
	}
	if bytes.Equal(res.Bytes(), c.Bytes()) {
		return cur, nil
	}
	return res.Bytes(), conflicts
}

// frontFields is a file's frontmatter as key and value lines, in file order.
type frontFields struct {
	keys []string
	raw  map[string]string
}

func fields(m *Milestone) frontFields {
	f := frontFields{raw: map[string]string{}}
	end, ok := m.frontEnd()
	if !ok {
		return f
	}
	for _, line := range m.lines[1:end] {
		k, v, ok := strings.Cut(strings.TrimSpace(line), ":")
		k = strings.TrimSpace(k)
		if !ok || k == "" {
			continue
		}
		if _, dup := f.raw[k]; !dup {
			f.keys = append(f.keys, k)
		}
		f.raw[k] = strings.TrimSpace(v)
	}
	return f
}

// val is a field's value as it is compared. after: is compared as IDs, the way
// a migration spells it, so a base that was migrated equals an old file that
// still holds names.
func (f frontFields) val(k string) (string, bool) {
	v, ok := f.raw[k]
	if ok && k == "after" {
		v = strings.Join(afterIDs(parseAfter(v)), ", ")
	}
	return v, ok
}

func afterIDs(tokens []string) []string {
	out := make([]string, len(tokens))
	for i, tok := range tokens {
		out[i] = tok
		if !IsID(tok) {
			out[i] = LegacyID(tok)
		}
	}
	return out
}

// prose is everything below the frontmatter that is not a member line.
func prose(m *Milestone) string {
	start := 0
	if end, ok := m.frontEnd(); ok {
		start = end + 1
	}
	var kept []string
	for _, l := range m.lines[start:] {
		if _, member := parseMember(l); !member && !isConflictNote(l) {
			kept = append(kept, l)
		}
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// The comment a clash leaves in the file, one line per field. It is neither a
// field nor body: prose skips it, so it never shows up in a later comparison.
const conflictNotePrefix = "<!-- another edit changed "

func isConflictNote(line string) bool { return strings.HasPrefix(line, conflictNotePrefix) }

func conflictNoteLine(c Conflict) string {
	return conflictNotePrefix + c.Field + " too and its value was not taken: " +
		strings.ReplaceAll(c.Old, "--", "-") + " -->"
}

func conflictNotes(m *Milestone) []string {
	var out []string
	for _, l := range m.lines {
		if isConflictNote(l) {
			out = append(out, l)
		}
	}
	return out
}

// setConflictNote writes the line for c's field, replacing the one an earlier
// clash on the same field left.
func (m *Milestone) setConflictNote(c Conflict) {
	line := conflictNoteLine(c)
	if slices.Contains(m.lines, line) {
		return
	}
	own := conflictNotePrefix + c.Field + " too "
	m.lines = slices.DeleteFunc(slices.Clone(m.lines), func(l string) bool { return strings.HasPrefix(l, own) })
	m.lines = append(m.lines, line)
}

// NoteConflict writes the comment for c into a milestone file, replacing the
// one an earlier clash on the same field left.
func NoteConflict(content []byte, c Conflict) []byte {
	m := parse(string(content))
	m.setConflictNote(c)
	return m.Bytes()
}

// addConflictNote appends a comment line unless the file already has one for
// its field.
func (m *Milestone) addConflictNote(line string) {
	field, _, _ := strings.Cut(strings.TrimPrefix(line, conflictNotePrefix), " too ")
	own := conflictNotePrefix + field + " too "
	if !slices.ContainsFunc(m.lines, func(l string) bool { return strings.HasPrefix(l, own) }) {
		m.lines = append(m.lines, line)
	}
}

// withBody is res with the lines below its frontmatter replaced by o's.
func withBody(res, o *Milestone) *Milestone {
	head := func(m *Milestone) []string {
		if end, ok := m.frontEnd(); ok {
			return m.lines[:end+1]
		}
		return nil
	}
	start := len(head(o))
	lines := append(slices.Clone(head(res)), o.lines[start:]...)
	return parse(strings.Join(lines, "\n") + "\n")
}

func quoteOld(v string, present bool) string {
	if !present {
		return "(removed)"
	}
	v = strings.Join(strings.Fields(v), " ")
	if len(v) > 120 {
		v = v[:120] + "..."
	}
	return fmt.Sprintf("%q", v)
}
