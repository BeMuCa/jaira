// Package epic holds the features a board is built out of: a description for a
// person and the tickets that make it up, kept on a ref of its own,
// refs/jaira/epics/<name>, and nowhere else.
//
// An epic is not a ticket. It has no definition of done, runs through no lane
// and is worked by no agent — it is a description plus a list. Nor is it a
// milestone: a milestone is a round of work, and what is unfinished moves on
// to the next one; an epic is a feature, and what is unfinished stays in it.
//
// The member list is in the epic, never on the tickets. Sorting a ticket into
// an epic therefore writes the epic's ref alone, and never the ref of a ticket
// an agent may be working at that very moment — so the two cannot conflict.
//
// There is no working file. The ref is read directly, and an edit that has not
// reached the remote yet is read back out of the outbox, which is the only copy
// of it. Both are behind the Refs interface, so nothing in here talks to git.
//
// The file the ref carries is read and written the way a milestone file is:
// frontmatter for the facts, one "- <ticket id>" line per member below it, and
// every other line kept verbatim across a write, so a person may edit it by
// hand and read it in a diff.
package epic

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// LockName keys the store lock every epic write is made under. One name for
// all of them: creating an epic has to check that the name is free, and two
// creations of one name at once must not both find it so.
const LockName = "epics"

// ErrNoRefs means the board has nowhere to keep an epic: no repository or no
// remote. Epics live on refs only, so there is no local fallback to offer.
var ErrNoRefs = errors.New("epics live on git refs, and this board has no repository with a remote")

// ErrNotFound means no epic of that name exists.
var ErrNotFound = errors.New("no such epic")

// ErrExists means an epic of that name exists already.
var ErrExists = errors.New("epic exists already")

// Refs is where epics are kept. *refsync.Syncer implements it; a test can
// hand in a map.
type Refs interface {
	// EpicNames lists every epic, unsent ones included.
	EpicNames() ([]string, error)
	// EpicContent reads one epic: an unsent write if there is one, else the
	// local ref. A missing epic wraps ErrNotFound.
	EpicContent(name string) ([]byte, error)
	// RecordEpic queues an epic's whole content for its ref.
	RecordEpic(name string, content []byte) error
}

// NormalizeName turns a written epic name into the form it is stored under:
// lowercase kebab, the tag rules. The name is a ref path segment and the name
// of the file inside the ref, so it has to be safe in both.
func NormalizeName(raw string) (string, bool, error) {
	name, changed, err := tag.Normalize(raw)
	if err != nil {
		return "", false, fmt.Errorf("epic name %q: %w", raw, err)
	}
	return name, changed, nil
}

// Epic is one epic as it was read: every line kept verbatim beside the facts
// parsed out of them.
type Epic struct {
	Name        string
	Title       string
	Description string
	CreatedAt   time.Time

	lines   []string
	members []string
}

// Members lists the ticket ids this epic holds, in file order.
func (e *Epic) Members() []string { return append([]string(nil), e.members...) }

// Has reports whether a ticket id is a member.
func (e *Epic) Has(id string) bool { return slices.Contains(e.members, id) }

// Add appends one member line. A ticket already in the epic is left where it
// is: re-adding must not reorder a list somebody arranged by hand.
func (e *Epic) Add(id string) bool {
	if e.Has(id) {
		return false
	}
	e.lines = append(e.lines, "- "+id)
	e.members = append(e.members, id)
	return true
}

// Remove drops a member's line, keeping everything around it.
func (e *Epic) Remove(id string) bool {
	if !e.Has(id) {
		return false
	}
	e.lines = slices.DeleteFunc(slices.Clone(e.lines), func(line string) bool {
		got, ok := milestone.ParseMember(line)
		return ok && got == id
	})
	e.members = slices.DeleteFunc(slices.Clone(e.members), func(got string) bool { return got == id })
	return true
}

// New builds an epic that has never been written.
func New(name, title, description string, now time.Time) *Epic {
	if title == "" {
		title = name
	}
	e := &Epic{Name: name, Title: title, Description: description, CreatedAt: now}
	e.lines = []string{
		"---",
		"title: " + oneLine(title),
		"description: " + oneLine(description),
		"created-at: " + now.UTC().Format(time.RFC3339),
		"---",
		"",
		"# " + oneLine(title),
		"",
		"<!-- One ticket id per line. Tickets that follow on from one of these",
		"     (follows: or parent:) count towards the epic without a line. -->",
		"",
	}
	return e
}

// oneLine keeps a frontmatter value on its line: a newline in it would end the
// value and start a line nothing parses.
func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// FromBytes reads an epic out of what its ref carries. The name comes from the
// ref, never from the content: a name inside the file would be a second truth
// that an edit could make disagree with the ref it lies on.
func FromBytes(name string, content []byte) *Epic {
	e := &Epic{Name: name}
	if s := strings.TrimSuffix(string(content), "\n"); s != "" {
		e.lines = strings.Split(s, "\n")
	}
	end, hasFront := e.frontEnd()
	for i, line := range e.lines {
		if hasFront && i <= end {
			if i == 0 || i == end {
				continue
			}
			key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
			if !ok {
				continue
			}
			value = strings.TrimSpace(value)
			switch strings.TrimSpace(key) {
			case "title":
				e.Title = value
			case "description":
				e.Description = value
			case "created-at":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					e.CreatedAt = t
				}
			}
			continue
		}
		if id, ok := milestone.ParseMember(line); ok && !e.Has(id) {
			e.members = append(e.members, id)
		}
	}
	if e.Title == "" {
		e.Title = name
	}
	return e
}

// frontEnd reports where the frontmatter closes; it is one only when the file
// opens with the fence.
func (e *Epic) frontEnd() (int, bool) {
	if len(e.lines) == 0 || strings.TrimSpace(e.lines[0]) != "---" {
		return 0, false
	}
	for i := 1; i < len(e.lines); i++ {
		if strings.TrimSpace(e.lines[i]) == "---" {
			return i, true
		}
	}
	return 0, false
}

// setField writes or replaces one frontmatter line and touches nothing else.
func (e *Epic) setField(key, value string) {
	value = oneLine(value)
	end, hasFront := e.frontEnd()
	if !hasFront {
		e.lines = append([]string{"---", key + ": " + value, "---", ""}, e.lines...)
		return
	}
	e.lines = slices.Clone(e.lines)
	for i := 1; i < end; i++ {
		k, _, ok := strings.Cut(strings.TrimSpace(e.lines[i]), ":")
		if ok && strings.TrimSpace(k) == key {
			e.lines[i] = key + ": " + value
			return
		}
	}
	e.lines = slices.Insert(e.lines, end, key+": "+value)
}

// Bytes renders the epic as the file its ref carries.
func (e *Epic) Bytes() []byte {
	var b strings.Builder
	for _, l := range e.lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

// Merge rebases an edit that lost the race for the ref. base is what the edit
// started from, ours the edit, theirs what the remote holds now. The result is
// theirs with our member additions and removals applied, and our title and
// description where we changed them — so two people sorting tickets into one
// epic at once both get their tickets in, and nobody's hand edit is lost.
func Merge(name string, base, ours, theirs []byte) []byte {
	b, o := FromBytes(name, base), FromBytes(name, ours)
	out := FromBytes(name, theirs)
	for _, id := range o.members {
		if !b.Has(id) {
			out.Add(id)
		}
	}
	for _, id := range b.members {
		if !o.Has(id) {
			out.Remove(id)
		}
	}
	if o.Title != b.Title {
		out.Title = o.Title
		out.setField("title", o.Title)
	}
	if o.Description != b.Description {
		out.Description = o.Description
		out.setField("description", o.Description)
	}
	return out.Bytes()
}

// Load reads one epic.
func Load(r Refs, name string) (*Epic, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	content, err := r.EpicContent(name)
	if err != nil {
		return nil, err
	}
	return FromBytes(name, content), nil
}

// LoadAll reads every epic, sorted by name. An epic whose ref cannot be read is
// skipped rather than failing the listing.
func LoadAll(r Refs) ([]*Epic, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	names, err := r.EpicNames()
	if err != nil {
		return nil, err
	}
	var out []*Epic
	for _, name := range names {
		if e, err := Load(r, name); err == nil {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Create writes a new epic with its first members. The name must be
// normalized already, and ids must be full ticket ids.
func Create(s *ticket.Store, r Refs, name, title, description string, ids []string) (*Epic, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	unlock, err := s.Lock(LockName)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if _, err := r.EpicContent(name); err == nil {
		return nil, fmt.Errorf("%w: %s", ErrExists, name)
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	e := New(name, title, description, time.Now())
	for _, id := range ids {
		e.Add(id)
	}
	if err := r.RecordEpic(name, e.Bytes()); err != nil {
		return nil, err
	}
	return e, nil
}

// Assign puts tickets into an epic (add) or takes them out, in one write of the
// epic's ref and of nothing else. It returns the epic as written and the ids
// that actually changed; nothing is written when none did. ids must be full
// ticket ids.
func Assign(s *ticket.Store, r Refs, name string, ids []string, add bool) (*Epic, []string, error) {
	if r == nil {
		return nil, nil, ErrNoRefs
	}
	unlock, err := s.Lock(LockName)
	if err != nil {
		return nil, nil, err
	}
	defer unlock()
	e, err := Load(r, name)
	if err != nil {
		return nil, nil, err
	}
	var changed []string
	for _, id := range ids {
		var ok bool
		if add {
			ok = e.Add(id)
		} else {
			ok = e.Remove(id)
		}
		if ok {
			changed = append(changed, id)
		}
	}
	if len(changed) > 0 {
		if err := r.RecordEpic(name, e.Bytes()); err != nil {
			return nil, nil, err
		}
	}
	return e, changed, nil
}

// Ticket is one ticket counted towards an epic.
type Ticket struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Lane  string `json:"lane"`
	Done  bool   `json:"done"`
	// Direct is true for a member line, false for a ticket that counts
	// because it follows on from one (follows: or parent:).
	Direct bool `json:"direct"`
}

// Progress is how far an epic is.
type Progress struct {
	Done    int      `json:"done"`
	Total   int      `json:"total"`
	Tickets []Ticket `json:"tickets"`
}

// Count says how far an epic is: its members, plus every ticket that follows
// on from one of them through follows: or parent:, however many steps down,
// wherever it is filed. Done is the index's Satisfied — a terminal lane or the
// logbook — so the epic agrees with the gate about what finished means.
//
// Follow-ups count without a line because they appear while the work goes on:
// closing one ticket often opens two, and a list that had to be kept up by hand
// would always be behind.
func Count(e *Epic, ix *link.Index) Progress {
	members := e.Members()
	ids := Covered(members, ix)
	p := Progress{Tickets: make([]Ticket, 0, len(ids))}
	for _, id := range ids {
		ref, _ := ix.Lookup(id)
		t := Ticket{ID: id, Title: ref.Title, Lane: ref.Lane, Done: ix.Satisfied(id), Direct: slices.Contains(members, id)}
		if t.Done {
			p.Done++
		}
		p.Tickets = append(p.Tickets, t)
	}
	p.Total = len(p.Tickets)
	return p
}

// Covered expands a set of member ids by every ticket that follows on from one
// of them, transitively, members first in their order, follow-ups after in id
// order.
func Covered(members []string, ix *link.Index) []string {
	kids := map[string][]string{}
	for _, t := range ix.All() {
		for _, up := range []string{t.Follows, t.Parent} {
			if up != "" && up != t.ID {
				kids[up] = append(kids[up], t.ID)
			}
		}
	}
	seen := map[string]bool{}
	var out, queue []string
	for _, id := range members {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
			queue = append(queue, id)
		}
	}
	var follow []string
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		for _, kid := range kids[id] {
			if !seen[kid] {
				seen[kid] = true
				follow = append(follow, kid)
				queue = append(queue, kid)
			}
		}
	}
	sort.Strings(follow)
	return append(out, follow...)
}

// Unassigned returns the open tickets of a board that belong to no epic,
// neither by a member line nor by following on from a member. board is the
// listing the caller already has; tickets in a terminal lane are left out,
// since finished work needs no home.
func Unassigned(epics []*Epic, ix *link.Index, board []*ticket.Ticket) []*ticket.Ticket {
	var members []string
	for _, e := range epics {
		members = append(members, e.members...)
	}
	in := map[string]bool{}
	for _, id := range Covered(members, ix) {
		in[id] = true
	}
	var out []*ticket.Ticket
	for _, t := range board {
		if !in[t.ID] && !ix.Satisfied(t.ID) {
			out = append(out, t)
		}
	}
	return out
}
