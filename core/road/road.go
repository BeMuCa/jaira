// Package road holds the swimlanes of the milestone roadmap as objects of their
// own: a title, a description and the milestones listed in them, kept on a ref
// of its own, refs/jaira/roads/<name>, and nowhere else.
//
// A road is a whole swimlane. The milestones listed in it are drawn in its
// lane, and what is built inside the lane still comes from after: — a road
// groups its milestones, it does not order them. A milestone never knows which
// road it is in: the list lives on the road alone, so placing a milestone after
// a listed one draws it in that lane without writing the road at all.
//
// A milestone is listed in at most one road. Create and Assign refuse a second
// listing with a TakenError and write nothing; a git merge can still produce
// one, and then the older road (created-at, then name) draws it.
//
// The file the ref carries is read the way an epic file is: frontmatter for
// the facts, one "- <milestone ID>" line per member below it, and every other
// line kept verbatim across a write. Both are behind the Refs interface, so
// nothing in here talks to git.
package road

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// lockName keys the store lock every road write is made under: creating a road
// has to check that its name is free, and two creations at once must not both
// find it so.
const lockName = "roads"

// ErrNoRefs means the board has nowhere to keep a road: no repository or no
// remote. Roads live on refs only, so there is no local fallback to offer.
var ErrNoRefs = errors.New("roads live on git refs, and this board has no repository with a remote")

// ErrNotFound means no road of that name exists.
var ErrNotFound = errors.New("no such road")

// ErrExists means a road of that name exists already.
var ErrExists = errors.New("road exists already")

// TakenError means a milestone was to be listed in a road while another road
// already lists it. Nothing was written.
type TakenError struct {
	Milestone string
	Road      string
}

func (e *TakenError) Error() string {
	return fmt.Sprintf("milestone %s is already in road %q", e.Milestone, e.Road)
}

// Refs is where roads are kept. *refsync.Syncer implements it; a test can hand
// in a map.
type Refs interface {
	// RoadNames lists every road, unsent ones included.
	RoadNames() ([]string, error)
	// RoadContent reads one road: an unsent write if there is one, else the
	// local ref. A missing road wraps ErrNotFound.
	RoadContent(name string) ([]byte, error)
	// RecordRoad queues a road's whole content for its ref.
	RecordRoad(name string, content []byte) error
}

// NormalizeName turns a written road name into the form it is stored under:
// lowercase kebab, the tag rules, since it is a ref path segment.
func NormalizeName(raw string) (string, bool, error) {
	name, changed, err := tag.Normalize(raw)
	if err != nil {
		return "", false, fmt.Errorf("road name %q: %w", raw, err)
	}
	return name, changed, nil
}

// Road is one road as it was read: every line kept verbatim beside the facts
// parsed out of them.
type Road struct {
	Name        string
	Title       string
	Description string
	CreatedAt   time.Time

	lines   []string
	members []string
}

// Label is how a message names the road: its title, with the name beside it
// when they differ. A road whose title is its name is named once.
func (r *Road) Label() string {
	if r.Title == "" || r.Title == r.Name {
		return fmt.Sprintf("%q", r.Name)
	}
	return fmt.Sprintf("%q (%s)", r.Title, r.Name)
}

// Members lists the milestone IDs this road lists, in file order.
func (r *Road) Members() []string { return append([]string(nil), r.members...) }

// Has reports whether a milestone ID is listed.
func (r *Road) Has(id string) bool { return slices.Contains(r.members, id) }

// Add appends one member line. A milestone already listed is left where it is.
func (r *Road) Add(id string) bool {
	if r.Has(id) {
		return false
	}
	r.lines = append(r.lines, "- "+id)
	r.members = append(r.members, id)
	return true
}

// Remove drops a member's line, keeping everything around it.
func (r *Road) Remove(id string) bool {
	if !r.Has(id) {
		return false
	}
	r.lines = slices.DeleteFunc(slices.Clone(r.lines), func(line string) bool {
		got, ok := parseMember(line)
		return ok && got == id
	})
	r.members = slices.DeleteFunc(slices.Clone(r.members), func(got string) bool { return got == id })
	return true
}

// New builds a road that has never been written. An empty title is the name.
func New(name, title, description string, now time.Time) *Road {
	if title == "" {
		title = name
	}
	r := &Road{Name: name, Title: title, Description: description, CreatedAt: now}
	r.lines = []string{
		"---",
		"title: " + oneLine(title),
		"description: " + oneLine(description),
		"created-at: " + now.UTC().Format(time.RFC3339),
		"---",
		"",
		"<!-- One milestone ID per line. A milestone placed after one of these is drawn in this road without a line. -->",
		"",
	}
	return r
}

// oneLine keeps a frontmatter value on its line: a newline in it would end the
// value and start a line nothing parses.
func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// FromBytes reads a road out of what its ref carries. The name comes from the
// ref, never from the content.
func FromBytes(name string, content []byte) *Road {
	r := &Road{Name: name}
	if s := strings.TrimSuffix(string(content), "\n"); s != "" {
		r.lines = strings.Split(s, "\n")
	}
	end, hasFront := r.frontEnd()
	for i, line := range r.lines {
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
				r.Title = value
			case "description":
				r.Description = value
			case "created-at":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					r.CreatedAt = t
				}
			}
			continue
		}
		if id, ok := parseMember(line); ok && !r.Has(id) {
			r.members = append(r.members, id)
		}
	}
	if r.Title == "" {
		r.Title = name
	}
	return r
}

// parseMember reads one member line, "- <milestone ID>", with anything after
// the ID kept and ignored. It is not milestone.ParseMember, which reads ticket
// IDs: a road lists milestones, whose IDs are the six-character kind.
func parseMember(line string) (string, bool) {
	rest, ok := strings.CutPrefix(strings.TrimSpace(line), "- ")
	if !ok {
		return "", false
	}
	rest = strings.TrimSpace(rest)
	if i := strings.IndexAny(rest, " \t"); i > 0 {
		rest = rest[:i]
	}
	id := strings.ToUpper(rest)
	if !milestone.IsID(id) {
		return "", false
	}
	return id, true
}

// frontEnd reports where the frontmatter closes; it is one only when the file
// opens with the fence.
func (r *Road) frontEnd() (int, bool) {
	if len(r.lines) == 0 || strings.TrimSpace(r.lines[0]) != "---" {
		return 0, false
	}
	for i := 1; i < len(r.lines); i++ {
		if strings.TrimSpace(r.lines[i]) == "---" {
			return i, true
		}
	}
	return 0, false
}

// setField writes or replaces one frontmatter line and touches nothing else.
func (r *Road) setField(key, value string) {
	value = oneLine(value)
	end, hasFront := r.frontEnd()
	if !hasFront {
		r.lines = append([]string{"---", key + ": " + value, "---", ""}, r.lines...)
		return
	}
	r.lines = slices.Clone(r.lines)
	for i := 1; i < end; i++ {
		k, _, ok := strings.Cut(strings.TrimSpace(r.lines[i]), ":")
		if ok && strings.TrimSpace(k) == key {
			r.lines[i] = key + ": " + value
			return
		}
	}
	r.lines = slices.Insert(r.lines, end, key+": "+value)
}

// Bytes is the file the ref carries.
func (r *Road) Bytes() []byte {
	var b strings.Builder
	for _, l := range r.lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

// Merge rebases an edit that lost the race for the ref. base is what the edit
// started from, ours the edit, theirs what the remote holds now. The result is
// theirs with our member additions and removals applied, and our title and
// description where we changed them, so two people listing milestones in one
// road at once both get theirs in. Unlike an epic's text there is nothing a
// machine cannot pick: a road has no prose, and a title one side changed wins.
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

// Load reads one road.
func Load(r Refs, name string) (*Road, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	content, err := r.RoadContent(name)
	if err != nil {
		return nil, err
	}
	return FromBytes(name, content), nil
}

// LoadAll reads every road, oldest first (created-at, then name): that order is
// the one a milestone listed in two roads is drawn in, and the one Groups keeps.
// A road whose ref cannot be read is skipped rather than failing the listing.
func LoadAll(r Refs) ([]*Road, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	names, err := r.RoadNames()
	if err != nil {
		return nil, err
	}
	var out []*Road
	for _, name := range names {
		if x, err := Load(r, name); err == nil {
			out = append(out, x)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}

// checkFree refuses a milestone that a road other than self already lists.
func checkFree(all []*Road, self string, ids []string) error {
	for _, id := range ids {
		for _, x := range all {
			if x.Name != self && x.Has(id) {
				return &TakenError{Milestone: id, Road: x.Name}
			}
		}
	}
	return nil
}

// Create writes a new road with its first milestones. The name must be
// normalized already, and ids must be milestone IDs. Nothing is written when a
// milestone is taken by another road.
func Create(s *ticket.Store, r Refs, name, title, description string, msIDs []string) (*Road, error) {
	if r == nil {
		return nil, ErrNoRefs
	}
	unlock, err := s.Lock(lockName)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if _, err := r.RoadContent(name); err == nil {
		return nil, fmt.Errorf("%w: %s", ErrExists, name)
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	all, err := LoadAll(r)
	if err != nil {
		return nil, err
	}
	if err := checkFree(all, name, msIDs); err != nil {
		return nil, err
	}
	x := New(name, title, description, time.Now())
	for _, id := range msIDs {
		x.Add(id)
	}
	if err := r.RecordRoad(name, x.Bytes()); err != nil {
		return nil, err
	}
	return x, nil
}

// Assign lists milestones in a road (add) or takes them out, in one write of
// the road's ref. It returns the road as written and the IDs that actually
// changed; nothing is written when none did, or when one is taken by another
// road.
func Assign(s *ticket.Store, r Refs, name string, msIDs []string, add bool) (*Road, []string, error) {
	if r == nil {
		return nil, nil, ErrNoRefs
	}
	unlock, err := s.Lock(lockName)
	if err != nil {
		return nil, nil, err
	}
	defer unlock()
	x, err := Load(r, name)
	if err != nil {
		return nil, nil, err
	}
	if add {
		all, err := LoadAll(r)
		if err != nil {
			return nil, nil, err
		}
		if err := checkFree(all, name, msIDs); err != nil {
			return nil, nil, err
		}
	}
	var changed []string
	for _, id := range msIDs {
		var ok bool
		if add {
			ok = x.Add(id)
		} else {
			ok = x.Remove(id)
		}
		if ok {
			changed = append(changed, id)
		}
	}
	if len(changed) > 0 {
		if err := r.RecordRoad(name, x.Bytes()); err != nil {
			return nil, nil, err
		}
	}
	return x, changed, nil
}

// Set rewrites a road's title and description in one write of its ref; a nil
// field is left as it is, and a blank title is the name. It returns the road as
// it stands afterwards and whether anything was written.
func Set(s *ticket.Store, r Refs, name string, title, description *string) (*Road, bool, error) {
	if r == nil {
		return nil, false, ErrNoRefs
	}
	unlock, err := s.Lock(lockName)
	if err != nil {
		return nil, false, err
	}
	defer unlock()
	x, err := Load(r, name)
	if err != nil {
		return nil, false, err
	}
	written := false
	if title != nil {
		t := oneLine(*title)
		if t == "" {
			t = name
		}
		if t != x.Title {
			x.Title = t
			x.setField("title", t)
			written = true
		}
	}
	if description != nil {
		if d := oneLine(*description); d != x.Description {
			x.Description = d
			x.setField("description", d)
			written = true
		}
	}
	if written {
		if err := r.RecordRoad(name, x.Bytes()); err != nil {
			return nil, false, err
		}
	}
	return x, written, nil
}

// Duplicate is a milestone listed in more than one road. Roads names them, the
// oldest first: the first one draws the milestone.
type Duplicate struct {
	Milestone string
	Roads     []string
}

// Duplicates lists the milestones that more than one road lists, by milestone
// ID. A git merge is the only way to get one; validate reports it.
func Duplicates(roads []*Road) []Duplicate {
	holders := map[string][]string{}
	for _, x := range roads {
		for _, id := range x.members {
			holders[id] = append(holders[id], x.Name)
		}
	}
	var ids []string
	for id, names := range holders {
		if len(names) > 1 {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	var out []Duplicate
	for _, id := range ids {
		out = append(out, Duplicate{Milestone: id, Roads: holders[id]})
	}
	return out
}

// Groups is the member list of each road, in LoadAll order: what the milestone
// layout takes, with a road's index as its position here.
func Groups(roads []*Road) [][]string {
	out := make([][]string, len(roads))
	for i, x := range roads {
		out[i] = x.Members()
	}
	return out
}

// NewName is the name a new road gets from its title: the title normalized, and
// made unique among roads with -2, -3 and so on.
func NewName(roads []*Road, title string) string {
	base, _, err := NormalizeName(title)
	if err != nil {
		base = "road"
	}
	taken := map[string]bool{}
	for _, x := range roads {
		taken[x.Name] = true
	}
	name := base
	for n := 2; taken[name]; n++ {
		name = fmt.Sprintf("%s-%d", base, n)
	}
	return name
}
