// Package lane resolves the board's columns.
//
// Seven base lanes are compiled into the binary so a fresh clone renders a
// working board with no repository configuration. Custom lanes are single
// markdown files under ~/.jaira/lanes, which makes sharing one "send someone the
// file" rather than "reproduce my config". Both kinds go through this same
// parser: a built-in lane is not a special case, it is just a lane that happens
// to ship inside the executable.
package lane

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

//go:embed builtin/*.md
var builtinFS embed.FS

// Lane is one column of the board.
type Lane struct {
	ID          string
	Name        string
	Description string

	// After anchors ordering to another lane's ID. Anchoring is used instead of
	// a numeric position because two independently written lane files would
	// otherwise both claim the same slot.
	After string

	// Precedence orders lanes by how far along the work is, and is what a merge
	// consults so concurrent edits never revert forward progress. It is
	// deliberately separate from display order: Blocked appears last on the
	// board but must not outrank active work in a merge.
	Precedence int

	// Agentic marks a lane whose work is performed by a subagent.
	Agentic   bool
	Terminal  bool
	ModelTier string

	// InputRequires names the ticket fields assembled into a subagent's bounded
	// input. OutputProduces names the fields it must return.
	InputRequires  []string
	OutputProduces []string

	RequiresQuestion       bool
	RequiresOutcome        bool
	RequiresNonModelSignal bool

	// RequiresHumanExit means no agent may move a ticket out of this lane. It is
	// the difference between a review step that is conventionally respected and
	// one that actually stops: an agent that decides its own work passed review
	// would make the step decorative.
	RequiresHumanExit bool

	// RequiresOption names an entry in a ticket's Options checklist. When set,
	// this lane is only part of a ticket's path if that option is ticked — which
	// is how one ticket is planned and reviewed while another goes straight to
	// implementation, without either lane being uninstalled.
	RequiresOption string

	// RequiresSpecified marks the first lane, in precedence order, that sets
	// this: the boundary between thinking about a ticket and working on it.
	// Everything before it is a place a half-formed ticket may sit; nothing at
	// or after it may hold a ticket that is missing its promotion fields.
	RequiresSpecified bool

	// Prompt is the markdown body: the instruction given to the subagent.
	Prompt string

	// Builtin distinguishes shipped lanes from user-installed ones.
	Builtin bool
	// Source is the path a custom lane came from, for error messages.
	Source string
	// Unknown marks a synthetic passthrough lane created for a ticket whose
	// status names a lane this installation does not have.
	Unknown bool
}

// Set is the resolved, ordered collection of lanes.
type Set struct {
	Lanes []*Lane
	byID  map[string]*Lane
	// Warnings records non-fatal problems (a custom lane anchored to a lane that
	// is not installed, for instance) so the caller can surface them without
	// failing to draw a board.
	Warnings []string
}

// Get returns a lane by ID.
func (s *Set) Get(id string) (*Lane, bool) {
	l, ok := s.byID[id]
	return l, ok
}

// IDs lists lane IDs in display order.
func (s *Set) IDs() []string {
	out := make([]string, 0, len(s.Lanes))
	for _, l := range s.Lanes {
		out = append(out, l.ID)
	}
	return out
}

// Index reports a lane's display position, or -1.
func (s *Set) Index(id string) int {
	for i, l := range s.Lanes {
		if l.ID == id {
			return i
		}
	}
	return -1
}

// Precedence returns a lane's merge precedence. An unknown lane sorts below
// everything known, so a merge never promotes a ticket into a lane this
// installation cannot reason about.
func (s *Set) Precedence(id string) int {
	if l, ok := s.byID[id]; ok {
		return l.Precedence
	}
	return -1
}

// Terminal is the lane a signed-off ticket lands in — the first lane declaring
// itself terminal, which is where accepting work moves it to.
func (s *Set) Terminal() *Lane {
	for _, l := range s.Lanes {
		if l.Terminal {
			return l
		}
	}
	return nil
}

// Default is the lane new tickets enter.
func (s *Set) Default() *Lane {
	if l, ok := s.byID["backlog"]; ok {
		return l
	}
	if len(s.Lanes) > 0 {
		return s.Lanes[0]
	}
	return nil
}

// parse reads one lane definition. The frontmatter carries the contract; the
// markdown body is the prompt.
func parse(src []byte, source string, builtin bool) (*Lane, error) {
	d, err := ticket.ParseDoc(src)
	if err != nil {
		return nil, fmt.Errorf("lane %s: %w", source, err)
	}
	str := func(k string) string {
		v, _, err := d.Scalar(k)
		if err != nil {
			return ""
		}
		return v
	}
	list := func(k string) []string {
		v, err := d.List(k)
		if err != nil {
			return nil
		}
		return v
	}
	boolOf := func(k string) bool { return strings.EqualFold(str(k), "true") }
	// Distinguishing "absent" from "explicitly false" matters for the defaults
	// below: a lane that says nothing about evidence must not thereby escape it.
	boolOr := func(k string, def bool) bool {
		if !d.Has(k) {
			return def
		}
		return boolOf(k)
	}

	l := &Lane{
		ID:                str("id"),
		Name:              str("name"),
		Description:       str("description"),
		After:             str("after"),
		Agentic:           boolOf("agentic"),
		ModelTier:         str("model-tier"),
		InputRequires:     list("input-requires"),
		OutputProduces:    list("output-produces"),
		RequiresQuestion:  boolOf("requires-question"),
		RequiresSpecified: boolOf("requires-specified"),
		Prompt:            strings.TrimSpace(d.Body()),
		Builtin:           builtin,
		Source:            source,
	}
	// A terminal lane is where work is declared finished, so by default it
	// demands the same evidence the built-in Done lane does. Without this, anyone
	// — including an agent writing its own lane file — could define a terminal
	// lane with no requirements and use it to mark work complete on nothing but
	// its own assessment. Opting out is possible but must be deliberate.
	l.Terminal = boolOf("terminal")
	l.RequiresOutcome = boolOr("requires-outcome", l.Terminal)
	l.RequiresNonModelSignal = boolOr("requires-nonmodel-signal", l.Terminal)
	l.RequiresHumanExit = boolOf("requires-human-exit")
	l.RequiresOption = strings.TrimSpace(str("requires-option"))

	if l.ID == "" {
		return nil, fmt.Errorf("lane %s: missing required field \"id\"", source)
	}
	if !validID(l.ID) {
		return nil, fmt.Errorf("lane %s: id %q must be lowercase letters, digits and dashes", source, l.ID)
	}
	if l.Name == "" {
		l.Name = l.ID
	}
	if p := str("precedence"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("lane %s: precedence %q is not a number", source, p)
		}
		l.Precedence = n
	}
	if l.Agentic && l.Prompt == "" {
		return nil, fmt.Errorf("lane %s: lane is marked agentic but has no prompt body", source)
	}
	if l.Agentic && l.ModelTier == "" {
		return nil, fmt.Errorf("lane %s: agentic lane must declare a model-tier", source)
	}
	return l, nil
}

func validID(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return s != ""
}

// Builtins returns the seven lanes shipped in the binary.
func Builtins() ([]*Lane, error) {
	entries, err := builtinFS.ReadDir("builtin")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names) // numeric filename prefixes give the canonical order

	var out []*Lane
	for _, n := range names {
		// path.Join, not filepath.Join: an embedded filesystem always uses forward
		// slashes, so joining with the OS separator looks for "builtin\\00-backlog.md"
		// on Windows and finds nothing.
		b, err := builtinFS.ReadFile(path.Join("builtin", n))
		if err != nil {
			return nil, err
		}
		l, err := parse(b, "builtin:"+n, true)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

// UserLanesDir is where custom lane files live.
func UserLanesDir() string {
	if v := os.Getenv("JAIRA_LANES_DIR"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "lanes")
}

// Load resolves the full lane set: built-ins first, then any custom lanes.
func Load() (*Set, error) {
	lanes, err := Builtins()
	if err != nil {
		return nil, err
	}
	builtinIDs := map[string]bool{}
	for _, l := range lanes {
		builtinIDs[l.ID] = true
	}

	var warnings []string
	dir := UserLanesDir()
	if dir != "" {
		matches, _ := filepath.Glob(filepath.Join(dir, "*.md"))
		sort.Strings(matches)
		seen := map[string]string{}
		for _, m := range matches {
			b, err := os.ReadFile(m)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("could not read lane %s: %v", m, err))
				continue
			}
			l, err := parse(b, m, false)
			if err != nil {
				warnings = append(warnings, err.Error())
				continue
			}
			// A custom lane must never shadow a built-in: the promise that a
			// teammate's board behaves predictably depends on the base lanes
			// meaning the same thing everywhere.
			if builtinIDs[l.ID] {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: id %q collides with a built-in lane and was ignored", m, l.ID))
				continue
			}
			if prev, dup := seen[l.ID]; dup {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: id %q already defined by %s and was ignored", m, l.ID, prev))
				continue
			}
			seen[l.ID] = m
			lanes = append(lanes, l)
		}
	}

	ordered, orderWarn := order(lanes)
	warnings = append(warnings, orderWarn...)

	set := &Set{Lanes: ordered, byID: map[string]*Lane{}, Warnings: warnings}
	for _, l := range ordered {
		set.byID[l.ID] = l
	}
	return set, nil
}

// order places each lane after its anchor, falling back to just before the
// terminal lane when the anchor is not installed — a shared lane file that
// references someone else's custom lane still loads rather than breaking.
func order(lanes []*Lane) ([]*Lane, []string) {
	var warnings []string
	present := map[string]bool{}
	for _, l := range lanes {
		present[l.ID] = true
	}

	// Built-ins keep their shipped order; custom lanes are inserted relative to
	// their anchors.
	var out []*Lane
	for _, l := range lanes {
		if l.Builtin {
			out = append(out, l)
		}
	}

	var pending []*Lane
	for _, l := range lanes {
		if !l.Builtin {
			pending = append(pending, l)
		}
	}

	// Repeatedly place any lane whose anchor is already positioned. This
	// resolves chains of custom lanes anchored to each other.
	for progress := true; progress && len(pending) > 0; {
		progress = false
		var still []*Lane
		for _, l := range pending {
			idx := indexOf(out, l.After)
			if l.After == "" {
				// No anchor: park before the terminal lane.
				idx = terminalIndex(out) - 1
			}
			if idx < 0 {
				if !present[l.After] {
					warnings = append(warnings, fmt.Sprintf(
						"lane %s: anchor %q is not installed; placed before the terminal lane",
						l.Source, l.After))
					idx = terminalIndex(out) - 1
				} else {
					still = append(still, l) // anchor exists but is not placed yet
					continue
				}
			}
			at := idx + 1
			if at > len(out) {
				at = len(out)
			}
			out = append(out[:at], append([]*Lane{l}, out[at:]...)...)
			progress = true
		}
		pending = still
	}
	for _, l := range pending { // cyclic anchors
		warnings = append(warnings, fmt.Sprintf(
			"lane %s: anchor %q forms a cycle; appended at the end", l.Source, l.After))
		out = append(out, l)
	}
	return out, warnings
}

func indexOf(ls []*Lane, id string) int {
	if id == "" {
		return -1
	}
	for i, l := range ls {
		if l.ID == id {
			return i
		}
	}
	return -1
}

// terminalIndex finds the first terminal lane, or the end of the list.
func terminalIndex(ls []*Lane) int {
	for i, l := range ls {
		if l.Terminal {
			return i
		}
	}
	return len(ls)
}

// Passthrough synthesises a read-only column for a status no installed lane
// claims. Hiding such tickets would be a worse failure than showing an inert
// column: the ticket still exists and the user needs to know it does.
func Passthrough(id string) *Lane {
	return &Lane{
		ID:          id,
		Name:        id,
		Description: "Unrecognized lane. Install the lane file to work these tickets.",
		Precedence:  -1,
		Unknown:     true,
	}
}

// Columns returns the lanes to draw for a given set of ticket statuses,
// appending passthrough columns for any status with no matching lane. Unknown
// lanes are sorted alphabetically among themselves so the board is deterministic.
func (s *Set) Columns(statuses []string) []*Lane {
	out := append([]*Lane{}, s.Lanes...)
	var unknown []string
	seen := map[string]bool{}
	for _, st := range statuses {
		if st == "" || seen[st] {
			continue
		}
		seen[st] = true
		if _, ok := s.byID[st]; !ok {
			unknown = append(unknown, st)
		}
	}
	sort.Strings(unknown)
	for _, u := range unknown {
		out = append(out, Passthrough(u))
	}
	return out
}
