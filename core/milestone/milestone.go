// Package milestone holds the groups a board plans in: which tickets belong
// together for this round of work, kept on a ref of their own,
// refs/jaira/milestones/<ID>, and nowhere else.
//
// The grouping lives in its own document rather than as a field on each ticket
// because of what happens at the end of a round. Work that did not finish has
// to appear in the next group, and a per-ticket marker means touching every
// ticket one at a time — and tickets travel on their own refs, so that is one
// pull, one edit and one release each. Twenty of those is the reason nobody
// does it. A milestone is edited instead: the lines that are not finished move
// into the next milestone in one edit.
//
// There is no working file. A milestone used to be a file under
// .jaira/milestones/ as well as a ref, and the file was a second truth: a stale
// copy merged in from a branch or committed by mistake overwrote the ref, and
// files went back and forth through commits. Now the ref is read directly, and
// an edit that has not reached the remote yet is read back out of the outbox,
// which is the only copy of it — both behind the Refs interface, so nothing in
// here talks to git. A file still lying under .jaira/milestones/ from an older
// build is taken onto its ref once and deleted (Import).
//
// A milestone that has been filed into the logbook carries `status: filed` in
// its frontmatter, and that line is what keeps it off the board. The ref stays,
// so every clone reads the filing on its next fetch and a restore has
// something to take back. No line means the milestone is on the board. A filed
// milestone does not hold its name: the ID is what tells two milestones apart,
// so a new one may take the name and a restore beside it is a visible double,
// not a loss.
//
// The document the ref carries is read and written the way a ticket file is —
// frontmatter for the facts about the milestone, one line per member below it —
// and the lines are kept verbatim across a write. A comment, a blank line, a
// hand-chosen order and a line jaira cannot parse all survive Add and Remove,
// because those edit one line and leave the rest alone. Its name is
// <name>-<ID>.md inside the ref's tree: the name lives in the file name, never
// in the frontmatter, so there is one truth for it.
package milestone

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Subdir is where a board kept its milestone files, under .jaira, before they
// lived on refs only. It is the name core/ticket files a milestone under
// inside an old logbook folder too, and is taken from there so the two cannot
// drift apart.
const Subdir = ticket.MilestonesSubdir

// Dir is the directory an older build kept a board's milestone files in. Only
// Import reads it, to take what lies there onto the refs.
func Dir(root string) string { return filepath.Join(root, ticket.DirName, Subdir) }

// RefPrefix is the namespace a milestone's ref lives in, as gitref names it.
const RefPrefix = "refs/jaira/milestones/"

// File is one milestone as the refs hold it: the file name it is carried
// under, <name>-<ID>.md, and the document.
type File struct {
	Name    string
	Content []byte
}

// Refs is where milestones live; refsync's Syncer is one. MilestoneFiles lists
// every milestone — an unsent write where there is one, else the local ref —
// and RecordMilestone writes one: queued for the remote, and on a board
// without one put on the local ref straight away. A board that is not in a git
// repository at all has no milestones: MilestoneFiles lists none and
// RecordMilestone refuses.
type Refs interface {
	MilestoneFiles() ([]File, error)
	RecordMilestone(key, file string, content []byte) error
	// HasMilestoneRef reports whether the milestone with this ID already has
	// its ref in this tree.
	HasMilestoneRef(id string) bool
}

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
	ID        string // the six characters that stay when the name changes
	Colour    int
	CreatedAt time.Time
	Status    string
	Goal      string   // free text, one line; empty when the file has none
	After     []string // predecessor IDs (a name only for a hand-written line), first one decides the swimlane
	Placed    bool     // an `after:` line exists; `after: []` is a swimlane top

	file    string   // the file name its ref carries it under, .md included; empty for a milestone never written
	lines   []string // the whole file, frontmatter included
	members []string // ticket ids, in file order
}

// FileName is the name its ref carries this milestone under: the name it was
// read under, or <name>-<ID>.md for one never written.
func (m *Milestone) FileName() string {
	if m.file != "" {
		return m.file
	}
	return m.Name + "-" + m.ID + ".md"
}

// Ref is the ref this milestone lives on, which is what a message points a
// person at: 'git show <ref>:<file>' reads it.
func (m *Milestone) Ref() string { return RefPrefix + m.ID }

// Members lists the ticket ids this milestone holds, in the order the file
// gives them. File order is the order the card's marks are painted in, so it
// is a thing a person controls by editing.
func (m *Milestone) Members() []string { return append([]string(nil), m.members...) }

// HasColour reports whether this milestone paints a cell on a card. Zero is
// "no colour" rather than ANSI black: the board reserves the right-hand column
// whether a ticket is in a milestone or not, so something has to mean "leave
// the cell blank", and 0 is that value everywhere it is asked.
func (m *Milestone) HasColour() bool { return tag.ValidColour(m.Colour) && m.Colour > 0 }

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
		if got, ok := ParseMember(line); ok && got == id {
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
// caller picked. It is on no ref until Create records it.
func New(name string, colour int, now time.Time) *Milestone {
	m := &Milestone{Name: name, ID: NewID(), Colour: colour, CreatedAt: now}
	m.lines = []string{
		"---",
		"color: " + strconv.Itoa(colour),
		"created-at: " + now.UTC().Format(time.RFC3339),
		"---",
		"",
		"# " + name,
		"",
		"<!-- One ticket id per line. Move the unfinished ones into the next",
		"     milestone when this one is done. -->",
		"",
	}
	return m
}

// AmbiguousError is Load's answer when a name is carried by more than one
// milestone: nothing is chosen for the caller, who is handed the <name>-<ID>
// forms that tell the milestones apart.
type AmbiguousError struct {
	Name    string
	Choices []string
}

func (e *AmbiguousError) Error() string {
	return fmt.Sprintf("milestone %q is the name of %d milestones; say which one: %s",
		e.Name, len(e.Choices), strings.Join(e.Choices, " or "))
}

// Load reads one milestone, filed ones included, by whatever the caller has:
// its ID (any case), its <name>-<ID> file name, or a name. A name that two
// milestones carry is an *AmbiguousError, and a name on the board is preferred
// over the same name on a milestone that is filed. Nothing found reports
// os.ErrNotExist, so a caller can tell "no such milestone" from "unreadable".
//
// The name and the ID come from the file name the ref carries — the
// frontmatter does not carry either, so there is one truth for each.
func Load(refs Refs, key string) (*Milestone, error) {
	all, err := scan(refs)
	if err != nil {
		return nil, err
	}
	return lookup(all, key)
}

// lookup is Load's resolution over a list already read.
func lookup(all []*Milestone, key string) (*Milestone, error) {
	if id := strings.ToUpper(key); IsID(id) {
		for _, m := range all {
			if m.ID == id {
				return m, nil
			}
		}
	}
	if i := strings.LastIndexByte(key, '-'); i > 0 {
		if name, id, ok := SplitFile(key[:i+1] + strings.ToUpper(key[i+1:])); ok {
			for _, m := range all {
				if m.ID == id && m.Name == name {
					return m, nil
				}
			}
		}
	}
	if name, _, err := NormalizeName(key); err == nil {
		var live, filed []*Milestone
		for _, m := range all {
			switch {
			// A name as typed also finds a hand-named file as it stands, Q4.md.
			case m.Name != name && m.Name != key:
			case m.Filed():
				filed = append(filed, m)
			default:
				live = append(live, m)
			}
		}
		found := live
		if len(found) == 0 {
			found = filed
		}
		switch len(found) {
		case 0:
		case 1:
			return found[0], nil
		default:
			e := &AmbiguousError{Name: name}
			for _, m := range found {
				e.Choices = append(e.Choices, m.Name+"-"+m.ID)
			}
			return nil, e
		}
	}
	return nil, &fs.PathError{Op: "load", Path: key, Err: fs.ErrNotExist}
}

// ResolveRefs turns the predecessors a person typed for the milestone with
// this ID — each an ID, a <name>-<ID> or a name, as Load reads them — into the
// IDs after: stores, and puts them through CheckAfter. all is the milestones
// on the board. A token nobody carries is a *PredecessorError, a name two
// milestones carry an *AmbiguousError and an arrow that closes a loop
// CheckAfter's refusal, so a write that refuses has written nothing; repeats
// are dropped.
func ResolveRefs(all []*Milestone, id string, tokens []string) ([]string, error) {
	out := []string{}
	for _, tok := range tokens {
		m, err := lookup(all, tok)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return nil, &PredecessorError{Pred: tok}
		case err != nil:
			return nil, err
		}
		if !slices.Contains(out, m.ID) {
			out = append(out, m.ID)
		}
	}
	if err := CheckAfter(all, id, out); err != nil {
		return nil, err
	}
	return out, nil
}

// FromBytes reads a milestone out of the bytes a ref carries, under the file
// name the ref carries them as.
func FromBytes(file string, content []byte) *Milestone {
	m := parse(string(content))
	base := strings.TrimSuffix(file, ".md")
	m.file = base + ".md"
	m.Name, m.ID, _ = FileIdentity(base)
	return m
}

// scan reads every milestone the refs hold, filed ones included, sorted by
// name and then ID. Two that carry one ID are one milestone, and the first is
// the one that counts.
func scan(refs Refs) ([]*Milestone, error) {
	if refs == nil {
		return nil, nil
	}
	files, err := refs.MilestoneFiles()
	if err != nil {
		return nil, err
	}
	var out []*Milestone
	seen := map[string]bool{}
	for _, f := range files {
		m := FromBytes(f.Name, f.Content)
		if seen[m.ID] {
			continue
		}
		seen[m.ID] = true
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// LoadAll reads every milestone that is ON the board, sorted by name. No
// milestones is not an error: it is the state every board starts in.
//
// A filed milestone is not on the board and is left out here, which is what
// makes `status: filed` the line that decides it. Load reads a filed
// milestone, which is how a restore gets at it.
func LoadAll(refs Refs) ([]*Milestone, error) {
	all, err := scan(refs)
	if err != nil {
		return nil, err
	}
	out := all[:0]
	for _, m := range all {
		if !m.Filed() {
			out = append(out, m)
		}
	}
	return out, nil
}

// LoadFiled reads every milestone that has been filed, sorted by name: what
// the logbook listing names, since a filed milestone lies nowhere on disk.
func LoadFiled(refs Refs) ([]*Milestone, error) {
	all, err := scan(refs)
	if err != nil {
		return nil, err
	}
	var out []*Milestone
	for _, m := range all {
		if m.Filed() {
			out = append(out, m)
		}
	}
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
			case "color":
				m.Colour, _ = strconv.Atoi(value)
			case "status":
				m.Status = value
			case "goal":
				m.Goal = value
			case "after":
				m.After, m.Placed = parseAfter(value), true
			case "created-at":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					m.CreatedAt = t
				}
			}
			continue
		}
		if id, ok := ParseMember(line); ok && !m.Has(id) {
			m.members = append(m.members, id)
		}
	}
	return m
}

// ParseMember reads one member line: "- <ticket id>", with anything after the
// id — a title, a comment — ignored and kept. Only a well-formed id counts, so
// an ordinary markdown bullet in the prose is not mistaken for a member.
// core/epic reads its member lines with this too, so the two lists cannot
// disagree about what a member line is.
func ParseMember(line string) (string, bool) {
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

// StatusFiled is the one status a milestone can carry. A milestone is either
// on the board or filed into the logbook, so the field is a flag with a name
// rather than a state machine — and a name, because a reader of the file has
// to be able to tell what it means without a table.
const StatusFiled = "filed"

// Filed reports whether this milestone has been put into the logbook. A filed
// milestone is off the board: 'jaira milestone ls' does not name it, no card
// carries its colour.
func (m *Milestone) Filed() bool { return m.Status == StatusFiled }

// RenamedTo is the value of the `renamed-to:` line a rename of the name-keyed
// era left in the file it marked under the old name, or "" when the file has
// none. Such a file is not a milestone any more. It reads the raw lines — the
// field it was parsed into is gone.
func RenamedTo(content []byte) string {
	m := &Milestone{}
	m.lines = strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	end, ok := m.frontEnd()
	if !ok {
		return ""
	}
	for _, line := range m.lines[1:end] {
		if k, v, _ := strings.Cut(strings.TrimSpace(line), ":"); strings.TrimSpace(k) == "renamed-to" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// SetStatus writes, replaces or removes the status line in the frontmatter and
// touches nothing else — the same rule Add and Remove follow. Rebuilding the
// frontmatter from the parsed facts would be shorter and would throw away the
// comments and the order somebody chose, which is exactly what the format
// promises to keep.
//
// The empty status removes the line, which is how a restore puts a milestone
// back on the board. A file with no frontmatter at all is given one, because
// there is nowhere else for the status to go and a hand-written file is
// allowed to have started without it.
func (m *Milestone) SetStatus(status string) {
	m.Status = status
	m.setFront("status", status, status == "")
}

// setFront is the one-line rewrite SetStatus pioneered, shared by every
// frontmatter key jaira writes: replace the key's line where it stands, insert
// it before the closing fence, or — remove — drop it. Nothing else in the file
// moves.
func (m *Milestone) setFront(key, value string, remove bool) {
	end, hasFront := m.frontEnd()
	idx := -1
	if hasFront {
		for i := 1; i < end; i++ {
			k, _, ok := strings.Cut(strings.TrimSpace(m.lines[i]), ":")
			if ok && strings.TrimSpace(k) == key {
				idx = i
				break
			}
		}
	}
	line := key + ": " + value
	switch {
	case remove:
		if idx >= 0 {
			m.lines = slices.Delete(slices.Clone(m.lines), idx, idx+1)
		}
	case idx >= 0:
		m.lines = slices.Clone(m.lines)
		m.lines[idx] = line
	case hasFront:
		m.lines = slices.Insert(slices.Clone(m.lines), end, line)
	default:
		m.lines = append([]string{"---", line, "---", ""}, m.lines...)
	}
}

// SetColour writes the colour line where it stands — the layer's c, which
// steps a milestone through the palette.
func (m *Milestone) SetColour(colour int) {
	m.Colour = colour
	m.setFront("color", strconv.Itoa(colour), false)
}

// SetGoal writes, replaces or (empty text) removes the goal line. One line of
// text: a newline in it would end the value and leave the rest as a stray
// frontmatter line, so line breaks become spaces.
func (m *Milestone) SetGoal(text string) {
	text = strings.Join(strings.Fields(text), " ")
	m.Goal = text
	m.setFront("goal", text, text == "")
}

// SetAfter places the milestone in the layer under the given predecessors; no
// names makes it the top of a swimlane of its own (`after: []`). Names are
// normalized and deduplicated, so what is written is what a reader gets back.
func (m *Milestone) SetAfter(names []string) {
	m.After, m.Placed = cleanNames(names), true
	m.setFront("after", "["+strings.Join(m.After, ", ")+"]", false)
}

// Unplace takes the milestone out of the layer by removing its after line.
func (m *Milestone) Unplace() {
	m.After, m.Placed = nil, false
	m.setFront("after", "", true)
}

// ParseNames reads a comma-separated list of milestone names the way --after
// and the layer's form take it. Unlike a file being read, typed input is
// refused when a name cannot be one, rather than having it quietly left out.
// Empty is the empty list: a swimlane top.
func ParseNames(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			continue
		}
		if _, _, err := NormalizeName(p); err != nil {
			return nil, err
		}
	}
	return cleanNames(parts), nil
}

// parseAfter reads the value of an after line: `[a, b]`, or the same without
// brackets, which is what a hand-edit tends to produce.
func parseAfter(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	return cleanNames(strings.Split(value, ","))
}

// cleanNames cleans a list of predecessor tokens: an ID stays exactly as it
// is — it is upper case, and lowering it would make it a name — and anything
// else is normalized like a milestone name, dropping what does not normalize,
// empties and repeats. The order is kept — the first one is the one that
// decides the swimlane.
func cleanNames(raw []string) []string {
	var out []string
	for _, r := range raw {
		if strings.TrimSpace(r) == "" {
			continue
		}
		if id := strings.TrimSpace(r); IsID(id) {
			if !slices.Contains(out, id) {
				out = append(out, id)
			}
			continue
		}
		n, _, err := NormalizeName(r)
		if err != nil || slices.Contains(out, n) {
			continue
		}
		out = append(out, n)
	}
	return out
}

// frontEnd reports where the frontmatter closes. It is a frontmatter only when
// the file opens with the fence, which is the rule parse reads by: a "---"
// further down is a horizontal rule in somebody's prose.
func (m *Milestone) frontEnd() (int, bool) {
	if len(m.lines) == 0 || strings.TrimSpace(m.lines[0]) != "---" {
		return 0, false
	}
	for i := 1; i < len(m.lines); i++ {
		if strings.TrimSpace(m.lines[i]) == "---" {
			return i, true
		}
	}
	return 0, false
}

// Bytes renders the milestone as the document it is: every line as it was
// read or edited, newline-terminated. It is what the ref carries.
func (m *Milestone) Bytes() []byte {
	var b strings.Builder
	for _, l := range m.lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return []byte(b.String())
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

// Labels is how each of these milestones is called on screen, by ID: its
// name, and <name>-<ID> only for the ones that share a name with another —
// the one case in which a name does not tell them apart.
func Labels(all []*Milestone) map[string]string {
	count := make(map[string]int, len(all))
	for _, m := range all {
		count[m.Name]++
	}
	out := make(map[string]string, len(all))
	for _, m := range all {
		out[m.ID] = m.Name
		if count[m.Name] > 1 {
			out[m.ID] = m.Name + "-" + m.ID
		}
	}
	return out
}

// Twin is one name that two or more milestones on the board share: the
// <name>-<ID> labels that tell them apart, in file order.
type Twin struct {
	Name   string
	Labels []string
}

// RenameCommand is the command that ends the twinning, renaming the first.
func (t Twin) RenameCommand() string {
	return "jaira milestone rename " + t.Labels[0] + " <new-name>"
}

// Twins is every name these milestones share, in the order the names first
// appear. nil when no two milestones are called the same — the one rule every
// warning about twins reads, so the CLI and the board cannot disagree.
func Twins(all []*Milestone) []Twin {
	labels := Labels(all)
	var out []Twin
	at := map[string]int{}
	for _, m := range all {
		if labels[m.ID] == m.Name {
			continue
		}
		i, ok := at[m.Name]
		if !ok {
			i = len(out)
			at[m.Name] = i
			out = append(out, Twin{Name: m.Name})
		}
		out[i].Labels = append(out[i].Labels, labels[m.ID])
	}
	return out
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

// Matches reports whether a ticket belongs to the milestone want names: by its
// ID, by its <name>-<ID>, or by its name — and a name hits every milestone that
// carries it, since filtering writes nothing and two milestones of one name are
// the same question. Letter case is not told apart. The comparison is equality,
// not a substring: a milestone is one of a closed set, and a filter is a
// question whose answer to an impossible name is "no".
func (i Index) Matches(id, want string) bool {
	want = strings.TrimSpace(want)
	name, _, nameErr := NormalizeName(want)
	for _, m := range i[id] {
		if strings.EqualFold(m.ID, want) || strings.EqualFold(m.Name+"-"+m.ID, want) {
			return true
		}
		if nameErr == nil && (m.Name == name || m.Name == want) {
			return true
		}
	}
	return false
}
