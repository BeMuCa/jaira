// Package lane resolves the board's columns.
//
// Seven base lanes are compiled into the binary so a fresh clone renders a
// working board with no repository configuration. Custom lanes are single
// markdown files under ~/.jaira/lanes (the catalogue) or, when a project has
// scoped itself down, under <root>/.jaira/lanes (see ProjectLanesDir) — which
// makes sharing one "send someone the file" rather than "reproduce my config".
// Both kinds go through this same parser: a built-in lane is not a special
// case, it is just a lane that happens to ship inside the executable.
//
// A custom lane whose id matches a built-in overrides it, prompt included,
// rather than being refused. That is deliberate: the user owns the pipeline
// and nothing about it is off limits, including the protections a built-in
// carries. The loader's job is not to stop an override, only to make sure it
// is never quiet — see Load.
package lane

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
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

	// RequiresBlockedReason means a ticket may not be parked here without saying
	// what it is waiting on. Parking is cheap and forgetting is silent, so this
	// is the one place the board charges for it.
	RequiresBlockedReason bool

	// RequiresCommits means a ticket cannot enter this lane without naming the
	// commits that carry its change. It sits on the done lane rather than on
	// the implementing lane's output contract, so work can move through review
	// uncommitted — but nothing is accepted that cannot be checked.
	RequiresCommits bool

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

	// RejectsTo names the lane a ticket goes back to when this lane finds the
	// work wanting: the loop's back edge, written down. Moving backwards has
	// always been allowed — the gate only checks a move that advances (see
	// core/gate) — but nothing declared where back was, so the route lived in
	// the prompt's prose and an agent reading the board could not see it.
	//
	// It is documentation the tool validates, not a rail: nothing forces a
	// rejection to go here, exactly as nothing forces a lane's prompt to be
	// followed.
	RejectsTo string

	// RequiresSpecified marks the first lane, in precedence order, that sets
	// this: the boundary between thinking about a ticket and working on it.
	// Everything before it is a place a half-formed ticket may sit; nothing at
	// or after it may hold a ticket that is missing its promotion fields.
	RequiresSpecified bool

	// Prompt is the markdown body: the instruction given to the subagent.
	Prompt string

	// Creator names who wrote this lane, for provenance once a lane has been
	// adopted or copied between machines. A built-in with no creator: field
	// defaults to "jaira" rather than having the field written into all nine
	// shipped files; a custom lane with no creator: field is left empty, since
	// "shipped by the tool" and "author unknown" are different facts.
	Creator string

	// Builtin distinguishes shipped lanes from user-installed ones.
	Builtin bool
	// Source is the path a custom lane came from, for error messages.
	Source string
	// Overrides names the built-in this lane displaced, empty if it does not
	// override anything. Set by Load, not by parse: overriding is a fact about
	// resolution, not about the file itself.
	Overrides string
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

// Options lists the distinct RequiresOption values declared by the installed
// lanes, in display order. A ticket's Options checklist is derived from what
// is actually installed rather than a fixed list, so a user-written optional
// lane shows up in it too, without the loader needing to know its name.
func (s *Set) Options() []string {
	seen := map[string]bool{}
	var out []string
	for _, l := range s.Lanes {
		if l.RequiresOption == "" || seen[l.RequiresOption] {
			continue
		}
		seen[l.RequiresOption] = true
		out = append(out, l.RequiresOption)
	}
	return out
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
		Creator:           str("creator"),
		Builtin:           builtin,
		Source:            source,
	}
	if l.Creator == "" && builtin {
		// The nine shipped lane files carry no creator: line — defaulting it
		// here is the same observable result for none of the file churn, and
		// no collision with lane files quick tasks are editing in parallel.
		l.Creator = "jaira"
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
	l.RequiresBlockedReason = boolOf("requires-blocked-reason")
	l.RequiresCommits = boolOf("requires-commits")
	l.RequiresOption = strings.TrimSpace(str("requires-option"))
	l.RejectsTo = strings.TrimSpace(str("rejects-to"))

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

// UserLanesDir is where custom lane files live. This is the catalogue: every
// lane a user has built or adopted, independent of any one project.
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

// ProjectLanesDir is where a project scopes itself down to the lanes it
// actually uses. Per D-03 it is authoritative when present and non-empty: it
// replaces the catalogue as this project's second lane source rather than
// adding to it, so a project directory holding only three lane files means a
// three-lane board, not the catalogue plus three overrides.
func ProjectLanesDir(root string) string {
	return filepath.Join(root, ticket.DirName, "lanes")
}

// ProjectLanesActive reports whether root's project lane directory is
// currently authoritative per D-03: present and holding at least one lane
// file. Load and 'jaira lanes path' both need this same test, so it lives
// here once rather than as two copies that could drift apart.
func ProjectLanesActive(root string) bool {
	if root == "" {
		return false
	}
	projDir := ProjectLanesDir(root)
	if fi, err := os.Stat(projDir); err != nil || !fi.IsDir() {
		return false
	}
	matches, _ := filepath.Glob(filepath.Join(projDir, "*.md"))
	return len(matches) > 0
}

// Load resolves the full lane set: built-ins first, then either the project's
// own lane directory or the catalogue, per D-03.
//
// root is the board's root. An empty root means no project is in hand — the
// launcher spans many boards and cannot scope to one — and Load falls back to
// built-ins plus the catalogue, exactly as it would for any project with no
// <root>/.jaira/lanes directory of its own.
//
// A custom lane whose id collides with a built-in overrides it, prompt
// included, rather than being refused. The override is always reported as a
// warning naming the file, the id and the built-in it displaced. If the
// override also drops a protection the built-in carried — requires-human-exit,
// terminal, requires-outcome or requires-nonmodel-signal — a second, separate
// warning names which protection is gone: those exist to stop an agent
// accepting its own work, and that is exactly what losing one silently would
// allow.
func Load(root string) (*Set, error) {
	lanes, err := Builtins()
	if err != nil {
		return nil, err
	}
	builtinByID := make(map[string]*Lane, len(lanes))
	for _, l := range lanes {
		builtinByID[l.ID] = l
	}

	var warnings []string

	// D-03: a project's own lane directory, when present and non-empty, is
	// authoritative over the catalogue for that project. Present-but-empty is
	// treated the same as absent (falls back to the catalogue) but is worth a
	// warning of its own: from the user's side it looks identical to "my lanes
	// vanished", and nothing else would say why.
	dir := UserLanesDir()
	if root != "" {
		projDir := ProjectLanesDir(root)
		if ProjectLanesActive(root) {
			dir = projDir
		} else if fi, statErr := os.Stat(projDir); statErr == nil && fi.IsDir() {
			warnings = append(warnings, fmt.Sprintf(
				"%s exists but holds no lane files; falling back to the catalogue", projDir))
		}
	}

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
			if prev, dup := seen[l.ID]; dup {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: id %q already defined by %s and was ignored", m, l.ID, prev))
				continue
			}
			seen[l.ID] = m

			if base, overrides := builtinByID[l.ID]; overrides {
				// A file that behaves exactly like the built-in it shadows is not an
				// override in any sense the reader cares about — it is a copy, which
				// is exactly what 'jaira lanes use' produces. Only warn, and only mark
				// it as an override for the settings screen's label, when something
				// that actually changes behaviour differs.
				if !lanesEquivalent(base, l) {
					l.Overrides = base.ID
					warnings = append(warnings, fmt.Sprintf(
						"lane %s: id %q overrides the built-in lane of the same name", m, l.ID))
					if lost := droppedProtections(base, l); len(lost) > 0 {
						warnings = append(warnings, fmt.Sprintf(
							"lane %s: overriding %q drops %s — an agent could now accept its own work here undetected",
							m, l.ID, strings.Join(lost, ", ")))
					}
				}
				lanes = replaceLane(lanes, l)
			} else {
				lanes = append(lanes, l)
			}
		}
	}

	// A project's removed-lanes tombstone excludes a lane even when it is a
	// built-in — built-ins are always injected above, so a project has no
	// file whose mere absence could mean "removed"; see removedFileName.
	if removedIDs, err := LoadRemoved(root); err == nil && len(removedIDs) > 0 {
		removedSet := make(map[string]bool, len(removedIDs))
		for _, id := range removedIDs {
			removedSet[id] = true
		}
		kept := lanes[:0:0]
		for _, l := range lanes {
			if !removedSet[l.ID] {
				kept = append(kept, l)
			}
		}
		lanes = kept
	}

	ordered, orderWarn := order(lanes)
	warnings = append(warnings, orderWarn...)

	// A project's own order file, when present, decides column order for the
	// lanes it names — after: is no longer consulted for those. See order.go.
	if orderIDs, err := LoadOrder(root); err == nil && len(orderIDs) > 0 {
		var applyWarn []string
		ordered, applyWarn = applyOrder(ordered, orderIDs)
		warnings = append(warnings, applyWarn...)
	}

	warnings = append(warnings, checkContracts(ordered)...)

	if !anyRequiresSpecified(ordered) {
		// An override could take requires-specified out of circulation entirely
		// (there is no built-in fallback for it the way there is for terminal or
		// requires-outcome): without it nothing marks the boundary between a
		// half-formed ticket and one ready to work, and a ticket can never be
		// promoted.
		warnings = append(warnings,
			"no installed lane declares requires-specified; a ticket can never be promoted into work")
	}

	set := &Set{Lanes: ordered, byID: map[string]*Lane{}, Warnings: warnings}
	for _, l := range ordered {
		set.byID[l.ID] = l
	}
	return set, nil
}

// lanesEquivalent reports whether replacement behaves exactly like base, so a
// file that merely re-ships a built-in unchanged is not treated as an
// override.
//
// Every field that can change what the lane actually does is compared: what
// gates it (RequiresQuestion, RequiresOutcome, RequiresNonModelSignal,
// RequiresHumanExit, RequiresOption, RequiresSpecified, Terminal), what it
// takes and gives an agent (InputRequires, OutputProduces, ModelTier,
// Agentic), what it says (Prompt, Name, Description — a person or an agent
// reads these to decide what the lane is for, so a changed description is a
// real change even though it does not touch a gate), and where it sits
// (After, Precedence).
//
// Creator is deliberately excluded: it records who made the file, not what it
// does, and a plain copy — 'jaira lanes use' — never carries the built-in's
// "jaira" default forward, so comparing it would flag every unmodified copy as
// changed. ID, Builtin, Source and Overrides are excluded for a different
// reason: ID is what matched base and replacement to each other in the
// caller, Builtin/Source describe provenance rather than behaviour, and
// Overrides is set by the caller after this comparison returns, so it carries
// no information yet on either side.
func lanesEquivalent(base, replacement *Lane) bool {
	return base.Name == replacement.Name &&
		base.Description == replacement.Description &&
		base.After == replacement.After &&
		base.Precedence == replacement.Precedence &&
		base.Agentic == replacement.Agentic &&
		base.Terminal == replacement.Terminal &&
		base.ModelTier == replacement.ModelTier &&
		base.RequiresQuestion == replacement.RequiresQuestion &&
		base.RequiresBlockedReason == replacement.RequiresBlockedReason &&
		base.RequiresCommits == replacement.RequiresCommits &&
		base.RequiresOutcome == replacement.RequiresOutcome &&
		base.RequiresNonModelSignal == replacement.RequiresNonModelSignal &&
		base.RequiresHumanExit == replacement.RequiresHumanExit &&
		base.RequiresOption == replacement.RequiresOption &&
		base.RejectsTo == replacement.RejectsTo &&
		base.RequiresSpecified == replacement.RequiresSpecified &&
		base.Prompt == replacement.Prompt &&
		slices.Equal(base.InputRequires, replacement.InputRequires) &&
		slices.Equal(base.OutputProduces, replacement.OutputProduces)
}

// droppedProtections reports which of a built-in's protections a replacement
// no longer carries, in a fixed order so the warning is deterministic.
func droppedProtections(base, replacement *Lane) []string {
	var lost []string
	if base.RequiresHumanExit && !replacement.RequiresHumanExit {
		lost = append(lost, "requires-human-exit")
	}
	if base.Terminal && !replacement.Terminal {
		lost = append(lost, "terminal")
	}
	if base.RequiresOutcome && !replacement.RequiresOutcome {
		lost = append(lost, "requires-outcome")
	}
	if base.RequiresNonModelSignal && !replacement.RequiresNonModelSignal {
		lost = append(lost, "requires-nonmodel-signal")
	}
	return lost
}

// replaceLane drops the lane sharing replacement's ID and appends replacement
// in its place. The replacement is not itself marked Builtin, so order() will
// position it by its own "after" like any other custom lane — an override
// does not inherit the shipped slot of the lane it displaced.
func replaceLane(lanes []*Lane, replacement *Lane) []*Lane {
	out := make([]*Lane, 0, len(lanes))
	for _, l := range lanes {
		if l.ID != replacement.ID {
			out = append(out, l)
		}
	}
	return append(out, replacement)
}

func anyRequiresSpecified(lanes []*Lane) bool {
	for _, l := range lanes {
		if l.RequiresSpecified {
			return true
		}
	}
	return false
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

// checkContracts walks lanes in display order and warns when a lane's
// input-requires names a field that nothing before it — the ticket itself,
// per ticket.SuppliedFields, or an earlier lane's output-produces — supplies.
// A lane ordered before its own producer otherwise fails later and opaquely,
// with "missing: plan" at the moment an agent asks for its input; this is
// the same load-time treatment the cycle check above already gives to a bad
// anchor.
//
// column order still follows after:, unchanged by this check — precedence is
// the merge rank (see Lane.Precedence), not a position, and this function
// does not touch either.
func checkContracts(lanes []*Lane) []string {
	available := make(map[string]bool, len(ticket.SuppliedFields))
	for _, f := range ticket.SuppliedFields {
		available[f] = true
	}
	// A global producer lookup, independent of position, so a violation can
	// name which installed lane produces the field — even when that lane is
	// simply ordered too late to help the one asking for it yet.
	producer := map[string]string{}
	for _, l := range lanes {
		for _, p := range l.OutputProduces {
			if _, exists := producer[p]; !exists {
				producer[p] = l.ID
			}
		}
	}

	installed := make(map[string]bool, len(lanes))
	for _, l := range lanes {
		installed[l.ID] = true
	}

	var warnings []string
	for _, l := range lanes {
		// A back edge is only useful if it points somewhere a ticket can actually
		// go, and a lane that sends work back to itself is not a loop, it is a
		// stall.
		if l.RejectsTo != "" {
			switch {
			case l.RejectsTo == l.ID:
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: rejects-to names itself, which would stall rather than loop", l.ID))
			case !installed[l.RejectsTo]:
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: rejects-to %q is not installed, so rejected work has nowhere to go",
					l.ID, l.RejectsTo))
			}
		}
		for _, want := range l.InputRequires {
			if available[want] {
				continue
			}
			if prod, ok := producer[want]; ok {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: requires %q, but %s (which produces it) is ordered after it",
					l.ID, want, prod))
			} else {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: requires %q, which no installed lane produces", l.ID, want))
			}
		}
		for _, p := range l.OutputProduces {
			available[p] = true
		}
	}
	return warnings
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
