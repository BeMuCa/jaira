// Package tag holds a board's shared tag vocabulary: the rules a tag name must
// obey, and the colour registry in .jaira/tags.
//
// Colours live in one file for the whole board rather than on each ticket. A
// colour is a property of the tag, not of the ticket wearing it, and storing it
// per ticket would mean the same tag could render two ways on one board and
// that recolouring "ui" meant rewriting every ticket carrying it.
//
// The file is read and written the way core/lane's order file is — plain text,
// one line per entry, hand-editable, no YAML — because that format is already
// the one this project uses for "a small fact about the board, not about a
// ticket", and a second mechanism for the same job would be a second thing to
// learn.
package tag

import (
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// FileName is the registry's name under .jaira. Not "tags.yaml" and not a
// directory: there is one fact per tag, and a line is the smallest thing git
// can merge.
const FileName = "tags"

// Path is where a board keeps its tag colours.
func Path(root string) string { return filepath.Join(root, ticket.DirName, FileName) }

// header is written once, when the file is created. It says what the number is
// and what jaira will and will not do to the file, because the format is the
// API and the person editing it by hand is the audience.
const header = `# jaira tag colours — one "name: <ansi256>" line per tag.
# The number is an ANSI-256 colour, 0-255. Hand-edit it freely: jaira only ever
# adds a line for a name it has not seen yet, rewrites the one line whose colour
# you change, and never reorders or reformats the rest of this file.
# A tag with no line here is real; it just renders without a colour.`

// Palette is the set of colours a new tag can be given: sixteen ANSI-256
// values chosen to stay apart from one another, and away from the colours the
// board already spends on meaning — accent 39, warning 214, error 203, ok 78,
// agentic 141 — so a tag swatch never reads as a status.
//
// Sixteen is a budget, not a limit on how many tags a board may have: past
// sixteen, colours repeat (see Registry.Assign), which is a legibility cost and
// not a failure.
var Palette = []int{
	33,  // blue
	45,  // cyan
	40,  // green
	71,  // moss
	100, // olive
	136, // amber
	178, // gold
	208, // orange
	209, // salmon
	197, // pink
	168, // rose
	170, // orchid
	135, // violet
	63,  // indigo
	111, // periwinkle
	73,  // teal
}

// Normalize turns a written tag name into the single form the board stores:
// lowercase kebab. It reports whether it had to change anything, so a caller can
// tell the person that "My UI" was filed as "my-ui" rather than leaving them to
// discover it on the ticket.
//
// Whitespace and underscores become dashes. Anything else outside [a-z0-9-] is
// refused rather than dropped: silently deleting a character invents a second
// name for one topic, which is the exact thing a shared vocabulary exists to
// prevent. The surviving charset is also a plain YAML scalar by construction, so
// a normalized name always survives a round trip through frontmatter.
func Normalize(raw string) (name string, changed bool, err error) {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(raw)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		case r == ' ' || r == '\t' || r == '_':
			b.WriteRune('-')
		default:
			return "", false, fmt.Errorf("tag name %q contains %q: a tag is lowercase letters, digits and dashes, and spaces become dashes", raw, string(r))
		}
	}
	name = strings.Trim(collapseDashes(b.String()), "-")
	if name == "" {
		return "", false, fmt.Errorf("tag name %q has nothing left once normalized", raw)
	}
	return name, name != raw, nil
}

// collapseDashes squeezes runs of dashes, so "ui  ux" does not become "ui--ux".
func collapseDashes(s string) string {
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

// Fallback is the colour a name gets once every palette colour is taken: a
// stable choice derived from the name itself, so an exhausted palette still
// gives one tag one colour on every machine, instead of a fresh random one per
// run and a file that never stops changing.
func Fallback(name string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return Palette[int(h.Sum32()%uint32(len(Palette)))]
}

// Registry is .jaira/tags as it was read: every line kept verbatim beside the
// colours parsed out of them.
//
// Keeping the raw lines is what makes a write safe on a file a person owns. A
// comment, a blank line, a line jaira cannot parse and a colour somebody
// hand-picked all survive a Set, because Set edits one line and leaves the slice
// alone.
type Registry struct {
	lines   []string
	colours map[string]int
	order   []string
	exists  bool
}

// Load reads a board's registry. An absent file is not an error: it means the
// board has no tag colours yet, which is the state every board starts in.
func Load(root string) (*Registry, error) {
	r := &Registry{colours: map[string]int{}}
	if root == "" {
		return r, nil
	}
	b, err := os.ReadFile(Path(root))
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}
		return nil, err
	}
	r.exists = true
	if s := strings.TrimSuffix(string(b), "\n"); s != "" {
		r.lines = strings.Split(s, "\n")
	}
	for _, line := range r.lines {
		name, colour, ok := parseLine(line)
		if !ok {
			continue
		}
		if _, dup := r.colours[name]; !dup {
			r.order = append(r.order, name)
		}
		r.colours[name] = colour
	}
	return r, nil
}

// parseLine reads one "name: <ansi256>" line. Anything else — a comment, a
// blank line, a typo — is not an entry and is reported as such, which is how it
// ends up preserved rather than rewritten.
//
// The name is normalized on the way in, so a file that says "UI: 39" colours
// the tag "ui" instead of quietly colouring nothing.
func parseLine(line string) (name string, colour int, ok bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return "", 0, false
	}
	rawName, rawColour, cut := strings.Cut(trimmed, ":")
	if !cut {
		return "", 0, false
	}
	name, _, err := Normalize(rawName)
	if err != nil {
		return "", 0, false
	}
	value := strings.TrimSpace(rawColour)
	if i := strings.Index(value, "#"); i >= 0 {
		value = strings.TrimSpace(value[:i])
	}
	colour, err = strconv.Atoi(value)
	if err != nil || !ValidColour(colour) {
		return "", 0, false
	}
	return name, colour, true
}

// ValidColour reports whether n is an ANSI-256 colour.
func ValidColour(n int) bool { return n >= 0 && n <= 255 }

// Colour returns the colour recorded for a name, and whether there is one.
func (r *Registry) Colour(name string) (int, bool) {
	c, ok := r.colours[name]
	return c, ok
}

// Names lists the names the registry has a colour for, in the order the file
// gives them.
func (r *Registry) Names() []string { return append([]string(nil), r.order...) }

// Assign picks a colour for a name the registry does not know: a random one no
// entry is using yet, so two tags on one board stay distinguishable. Random
// rather than next-in-sequence because sequence would give every board the same
// first three colours, and two teammates adding two tags at once would both
// pick the same "next" one.
//
// Once every palette colour is taken it falls back to Fallback — a repeat, not
// a refusal.
func (r *Registry) Assign(name string) int {
	used := make(map[int]bool, len(r.colours))
	for _, c := range r.colours {
		used[c] = true
	}
	free := make([]int, 0, len(Palette))
	for _, c := range Palette {
		if !used[c] {
			free = append(free, c)
		}
	}
	if len(free) == 0 {
		return Fallback(name)
	}
	return free[rand.IntN(len(free))]
}

// Set records a colour for a normalized name.
//
// A name already in the file has its own line rewritten in place, keeping any
// trailing comment. A new name is inserted at the position that keeps a sorted
// file sorted — appending instead would put every concurrent addition on the
// same last line, which is the one place two teammates' commits are guaranteed
// to conflict.
func (r *Registry) Set(name string, colour int) {
	if _, known := r.colours[name]; known {
		for i, line := range r.lines {
			if n, _, ok := parseLine(line); ok && n == name {
				r.lines[i] = entryLine(name, colour) + trailingComment(line)
				break
			}
		}
		r.colours[name] = colour
		return
	}
	r.lines = insertEntry(r.lines, name, entryLine(name, colour))
	r.colours[name] = colour
	r.order = append(r.order, name)
}

func entryLine(name string, colour int) string {
	return name + ": " + strconv.Itoa(colour)
}

// trailingComment returns a line's "#" comment together with the whitespace
// that separated it from the value, so recolouring a hand-annotated entry
// changes the number and nothing else — including its alignment.
func trailingComment(line string) string {
	i := strings.Index(line, "#")
	if i < 0 {
		return ""
	}
	j := i
	for j > 0 && (line[j-1] == ' ' || line[j-1] == '\t') {
		j--
	}
	return line[j:]
}

// insertEntry places a new entry line among the existing ones alphabetically,
// before the first entry that sorts after it. A file whose entries are not
// sorted is not reordered: the new line lands after the last entry instead, so
// hand-grouped tags stay grouped.
func insertEntry(lines []string, name, line string) []string {
	last := -1
	for i, l := range lines {
		n, _, ok := parseLine(l)
		if !ok {
			continue
		}
		last = i
		if n > name {
			out := append([]string{}, lines[:i]...)
			out = append(out, line)
			return append(out, lines[i:]...)
		}
	}
	if last < 0 {
		return append(append([]string{}, lines...), line)
	}
	out := append([]string{}, lines[:last+1]...)
	out = append(out, line)
	return append(out, lines[last+1:]...)
}

// Save writes the registry back, creating .jaira and the header if this is the
// board's first tag colour.
func (r *Registry) Save(root string) error {
	if err := os.MkdirAll(filepath.Join(root, ticket.DirName), 0o755); err != nil {
		return err
	}
	lines := r.lines
	if !r.exists {
		lines = append(strings.Split(header, "\n"), lines...)
	}
	var b strings.Builder
	for _, l := range lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(Path(root), []byte(b.String()), 0o644); err != nil {
		return err
	}
	r.exists = true
	r.lines = lines
	return nil
}
