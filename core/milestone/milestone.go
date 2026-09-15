// Package milestone holds the groups a board plans in: which tickets belong
// together for this round of work, kept in one file per group under
// .jaira/milestones/.
//
// The grouping lives in its own file rather than as a field on each ticket
// because of what happens at the end of a round. Work that did not finish has
// to appear in the next group, and a per-ticket marker means touching every
// ticket one at a time — and tickets travel on their own refs, so that is one
// pull, one edit and one release each. Twenty of those is the reason nobody
// does it. A file is copied instead: the lines that are not finished move into
// the next milestone's file in one edit.
//
// The file is read and written the way a ticket file is — frontmatter for the
// facts about the milestone, one line per member below it — and the lines are
// kept verbatim across a write. A comment, a blank line, a hand-chosen order
// and a line jaira cannot parse all survive Add and Remove, because those edit
// one line and leave the rest alone. The format is the API: it must stay
// hand-editable and readable in a diff, and one line per member is the
// smallest thing git can merge.
package milestone

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Subdir is where a board keeps its milestones, under .jaira.
const Subdir = "milestones"

// Dir is the milestone directory of a board rooted at root.
func Dir(root string) string { return filepath.Join(root, ticket.DirName, Subdir) }

// Path is the file one milestone lives in.
func Path(root, name string) string { return filepath.Join(Dir(root), name+".md") }

// Palette is the set of colours a new milestone can be given. It is the tag
// palette: the two marks sit on opposite edges of a card and are never
// confused for one another by position, so spending a second set of sixteen
// values would only push both closer to the colours the board already spends
// on status.
var Palette = tag.Palette

// NormalizeName turns a written milestone name into the form it is stored and
// filed under: lowercase kebab, exactly the tag rules. A milestone name is a
// filename, so the charset has to be one that is safe in a path on every
// platform — and reusing the tag rules means one answer to "what may I call
// this" instead of two that differ in a corner.
func NormalizeName(raw string) (string, bool, error) {
	name, changed, err := tag.Normalize(raw)
	if err != nil {
		return "", false, fmt.Errorf("milestone name %q: %w", raw, err)
	}
	return name, changed, nil
}

// Milestone is one file as it was read: every line kept verbatim beside the
// facts parsed out of them.
type Milestone struct {
	Name      string
	Colour    int
	CreatedAt time.Time

	lines   []string // the whole file, frontmatter included
	members []string // ticket ids, in file order
}

// Members lists the ticket ids this milestone holds, in the order the file
// gives them. File order is the order the card's marks are painted in, so it
// is a thing a person controls by editing.
func (m *Milestone) Members() []string { return append([]string(nil), m.members...) }

// Has reports whether a ticket id is a member.
func (m *Milestone) Has(id string) bool {
	for _, got := range m.members {
		if got == id {
			return true
		}
	}
	return false
}

// Add records a ticket as a member, appending one line. A ticket already in
// the file is left where it is: re-adding must not reorder a list somebody
// arranged by hand.
func (m *Milestone) Add(id string) bool {
	if m.Has(id) {
		return false
	}
	m.lines = append(m.lines, "- "+id)
	m.members = append(m.members, id)
	return true
}

// Remove drops a ticket's line, keeping everything around it.
func (m *Milestone) Remove(id string) bool {
	if !m.Has(id) {
		return false
	}
	kept := make([]string, 0, len(m.lines))
	for _, line := range m.lines {
		if got, ok := parseMember(line); ok && got == id {
			continue
		}
		kept = append(kept, line)
	}
	m.lines = kept
	rest := make([]string, 0, len(m.members))
	for _, got := range m.members {
		if got != id {
			rest = append(rest, got)
		}
	}
	m.members = rest
	return true
}

// New builds a milestone that has never been written, with the colour the
// caller picked. It is not on disk until Save.
func New(name string, colour int, now time.Time) *Milestone {
	m := &Milestone{Name: name, Colour: colour, CreatedAt: now}
	m.lines = []string{
		"---",
		"name: " + name,
		"color: " + strconv.Itoa(colour),
		"created-at: " + now.UTC().Format(time.RFC3339),
		"---",
		"",
		"# " + name,
		"",
		"<!-- One ticket id per line. Move the unfinished ones into the next",
		"     milestone's file when this one is done. -->",
		"",
	}
	return m
}

// Load reads one milestone by name. A missing file reports os.ErrNotExist, so
// a caller can tell "no such milestone" from "unreadable".
func Load(root, name string) (*Milestone, error) {
	path := Path(root, name)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := parse(string(b))
	if m.Name == "" {
		m.Name = name
	}
	return m, nil
}

// LoadAll reads every milestone on the board, sorted by name. An absent
// directory is not an error: it means the board has planned no milestones yet,
// which is the state every board starts in.
//
// A file that does not parse is skipped rather than failing the read: the
// board is glanced at constantly, and one malformed file must not be able to
// stop it from opening.
func LoadAll(root string) ([]*Milestone, error) {
	if root == "" {
		return nil, nil
	}
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []*Milestone
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		m, err := Load(root, name)
		if err != nil {
			continue
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// parse splits a milestone file into its lines, its frontmatter facts and its
// members. Nothing is rejected: an unparseable line is simply not a fact and
// not a member, which is how it survives a write untouched.
func parse(text string) *Milestone {
	m := &Milestone{}
	if s := strings.TrimSuffix(text, "\n"); s != "" {
		m.lines = strings.Split(s, "\n")
	}
	inFront := false
	for i, line := range m.lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if i == 0 {
				inFront = true
				continue
			}
			if inFront {
				inFront = false
				continue
			}
		}
		if inFront {
			key, value, ok := strings.Cut(trimmed, ":")
			if !ok {
				continue
			}
			value = strings.TrimSpace(value)
			switch strings.TrimSpace(key) {
			case "name":
				m.Name = value
			case "color":
				m.Colour, _ = strconv.Atoi(value)
			case "created-at":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					m.CreatedAt = t
				}
			}
			continue
		}
		if id, ok := parseMember(line); ok && !m.Has(id) {
			m.members = append(m.members, id)
		}
	}
	return m
}

// parseMember reads one member line: "- <ticket id>", with anything after the
// id — a title, a comment — ignored and kept. Only a well-formed id counts, so
// an ordinary markdown bullet in the prose is not mistaken for a member.
func parseMember(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	rest, ok := strings.CutPrefix(trimmed, "- ")
	if !ok {
		return "", false
	}
	rest = strings.TrimSpace(rest)
	if i := strings.IndexAny(rest, " \t"); i > 0 {
		rest = rest[:i]
	}
	rest = ticket.NormalizeIDPrefix(rest)
	if !ticket.ValidID(rest) {
		return "", false
	}
	return rest, true
}

// Save writes the milestone back through ticket.WriteAtomic, so a reader never
// sees a half-written file and a crash leaves the previous one intact.
// Atomicity is not exclusion: two writers still have to be serialised by the
// caller.
func (m *Milestone) Save(root string) error {
	if err := os.MkdirAll(Dir(root), 0o755); err != nil {
		return err
	}
	var b strings.Builder
	for _, l := range m.lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	if err := ticket.WriteAtomic(Path(root, m.Name), []byte(b.String())); err != nil {
		return err
	}
	return nil
}

// AssignColour picks a colour for a new milestone: a random one none of the
// existing milestones is using, so two groups on one board stay
// distinguishable. Random rather than next-in-sequence because a sequence
// gives every board the same first colours, and nobody should have to choose
// one. Once the palette is spent it repeats, derived from the name so the
// repeat is at least stable across machines.
func AssignColour(existing []*Milestone, name string) int {
	used := make(map[int]bool, len(existing))
	for _, m := range existing {
		used[m.Colour] = true
	}
	free := make([]int, 0, len(Palette))
	for _, c := range Palette {
		if !used[c] {
			free = append(free, c)
		}
	}
	if len(free) == 0 {
		return tag.Fallback(name)
	}
	return free[rand.IntN(len(free))]
}

// Index answers, for one ticket id, which milestones hold it — the direction
// every reader needs and the file does not give. Multiple membership is
// allowed: a ticket carried over from one round into the next is in both files
// until somebody removes the old line.
type Index map[string][]*Milestone

// Build indexes milestones by ticket id, keeping the board's file order within
// each ticket so the marks on a card are painted the same way everywhere.
func Build(all []*Milestone) Index {
	idx := Index{}
	for _, m := range all {
		for _, id := range m.members {
			idx[id] = append(idx[id], m)
		}
	}
	return idx
}

// For returns the milestones holding a ticket.
func (i Index) For(id string) []*Milestone { return i[id] }

// Names returns the milestone names holding a ticket.
func (i Index) Names(id string) []string {
	var out []string
	for _, m := range i[id] {
		out = append(out, m.Name)
	}
	return out
}

// Matches reports whether a ticket belongs to the named milestone. The
// comparison is equality on the normalized name, not a substring: a milestone
// is a name from a closed set, and a filter is a question whose answer to an
// impossible name is "no".
func (i Index) Matches(id, want string) bool {
	name, _, err := NormalizeName(want)
	if err != nil {
		return false
	}
	for _, m := range i[id] {
		if m.Name == name {
			return true
		}
	}
	return false
}
