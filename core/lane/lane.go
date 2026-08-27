// Package lane resolves the board's columns.
//
// A board is its lane directory: the lane files under <root>/.jaira/lanes/
// are the board, in the order of the order file beside them, and nothing is
// added underneath. The shipped lanes are compiled into the binary as the
// catalogue's standing offer, together with the files under ~/.jaira/lanes;
// a board opened for the first time gets its files written from the default
// board's selection or from the shipped lanes, and from then on 'jaira lanes
// add' and 'remove' are what change it. Both kinds of file go through this
// same parser: a shipped lane is not a special case, it is a lane that
// happens to ship inside the executable.
//
// A board's file that shares a shipped lane's id is that lane — the user
// owns the pipeline and nothing about it is off limits. The loader's one
// reservation is a protection the shipped lane carried and the file drops;
// that is reported, because it is exactly what the protection guards against.
// See Load.
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

// Load resolves this board's lanes.
//
// A board is its lane directory. When <root>/.jaira/lanes/ holds lane files,
// those files are the board — all of it, and nothing is added to them. The
// built-ins are compiled into the binary as the catalogue's standing offer
// (Installable lists them beside ~/.jaira/lanes); 'jaira lanes add' copies one
// onto a board, and nothing else does. The directory is per person and
// gitignored, so two people on one shared board may run two different boards.
//
// A board whose lane directory is empty is set up on first load: the default
// board's selection, or the built-ins when there is none, are written into the
// directory as files with an order file — once, reported in Warnings. This is
// the one place a read writes. It was asked for: a board must never sit
// half-configured, and the next command reads files rather than a fallback.
// Boards from before this rule are brought over by migrateLegacy.
//
// An empty root, or a root with no .jaira/ at all, means no board is in hand —
// the launcher spans many, a test has none, 'jaira lanes template' needs the
// shape only — and Load returns the offer: built-ins plus the catalogue.
//
// A file whose id matches a shipped lane is that lane, not an override of
// anything: nothing is on the board to override. It is still compared with the
// shipped lane of that name, for one reason — if it drops a protection the
// shipped lane carried (requires-human-exit, terminal, requires-outcome,
// requires-nonmodel-signal) a warning names which, because those exist to stop
// an agent accepting its own work, and losing one silently is exactly what
// they guard against.
func Load(root string) (*Set, error) {
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	builtinByID := make(map[string]*Lane, len(builtins))
	for _, l := range builtins {
		builtinByID[l.ID] = l
	}

	var warnings []string
	var lanes []*Lane
	if root == "" || !boardExists(root) {
		// The offer: built-ins, with the catalogue laid over them. A catalogue
		// lane sharing a built-in's id stands in for it here, so what 'jaira
		// lanes add' offers under that id is the file the user wrote.
		lanes = builtins
		var w []string
		lanes, w = readLaneDir(UserLanesDir(), lanes, builtinByID)
		warnings = append(warnings, w...)
	} else {
		if !ProjectLanesActive(root) {
			w, err := setUp(root)
			if err != nil {
				return nil, err
			}
			warnings = append(warnings, w...)
		} else {
			warnings = append(warnings, migrateLegacy(root)...)
		}
		var w []string
		lanes, w = readLaneDir(ProjectLanesDir(root), nil, builtinByID)
		warnings = append(warnings, w...)
	}

	ordered, orderWarn := order(lanes)

	// A board's order file decides column order — after: is a hint for where
	// a lane first lands, consulted only for a lane the file does not name.
	// So with an order file, a dangling anchor is not worth a warning: the
	// lane it pointed at was removed, and the order file already says where
	// everything goes. Without one, the anchors are all there is.
	if orderIDs, err := LoadOrder(root); err == nil && len(orderIDs) > 0 {
		var applyWarn []string
		ordered, applyWarn = applyOrder(ordered, orderIDs)
		warnings = append(warnings, applyWarn...)
	} else {
		warnings = append(warnings, orderWarn...)
	}

	warnings = append(warnings, checkContracts(ordered)...)

	if !anyRequiresSpecified(ordered) {
		// A board could take requires-specified out of circulation entirely
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

// boardExists reports whether root carries a .jaira/ directory at all — the
// difference between "this board has no lanes yet" and "this is not a board".
func boardExists(root string) bool {
	fi, err := os.Stat(filepath.Join(root, ticket.DirName))
	return err == nil && fi.IsDir()
}

// readLaneDir parses every lane file in dir onto lanes, in filename order. A
// file sharing an id with something already in lanes replaces it (that is how
// the catalogue lays over the built-ins in the offer); on a board lanes starts
// empty, so every file simply is a lane. Either way a file named like a shipped
// lane is checked for protections it drops, per Load's doc comment.
func readLaneDir(dir string, lanes []*Lane, builtinByID map[string]*Lane) ([]*Lane, []string) {
	var warnings []string
	if dir == "" {
		return lanes, nil
	}
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

		if base, shipped := builtinByID[l.ID]; shipped && !lanesEquivalent(base, l) {
			// Overrides labels the file in the settings screen as differing
			// from the shipped lane of that name. Only a dropped protection
			// is worth a warning; a lane the user changed on purpose is not.
			l.Overrides = base.ID
			if lost := droppedProtections(base, l); len(lost) > 0 {
				warnings = append(warnings, fmt.Sprintf(
					"lane %s: %q drops %s that the shipped lane of that name carries — an agent could now accept its own work here undetected",
					m, l.ID, strings.Join(lost, ", ")))
			}
		}
		lanes = replaceLane(lanes, l)
	}
	return lanes, warnings
}

// setUp writes a board's first lane directory: the default board's selection
// or, without one, the built-ins — each as a file, plus the order file. The
// returned warning says what was written, because a command that quietly
// creates a dozen files is a command whose output lies by omission.
func setUp(root string) ([]string, error) {
	offer, err := Load("")
	if err != nil {
		return nil, err
	}
	db, err := LoadDefaultBoard()
	if err != nil {
		return nil, err
	}
	ids, source := db.Lanes, "your default board"
	if len(ids) == 0 {
		source = "the built-in lanes"
		for _, l := range offer.Lanes {
			if l.Builtin {
				ids = append(ids, l.ID)
			}
		}
	}
	written, warnings, err := Materialise(root, offer, ids)
	if err != nil {
		return nil, err
	}
	warnings = append(warnings, fmt.Sprintf(
		"set up %s with %d lane(s) from %s — this board's lanes are its own now; 'jaira lanes' lists them, 'jaira lanes add|remove' changes them",
		ProjectLanesDir(root), len(written), source))
	return warnings, nil
}

// migrateLegacy brings a lane directory from before "a board is its lane
// directory" over to it. Back then Load injected the built-ins under whatever
// files the directory held, so a directory with three files was a
// thirteen-lane board, a "removed" file listed the built-ins to leave out, and
// the order file named lanes that had no file. Two signs mark such a
// directory: a "removed" file, or no order file at all.
//
// For each shipped lane the board used — every id the order file names, or
// every built-in when there is no order file — that has no file and is not in
// "removed", the shipped file is written beside the others, so the board keeps
// looking exactly as it did. Then the order file is written and "removed" is
// deleted, so this runs once. Everything done is reported.
func migrateLegacy(root string) []string {
	dir := ProjectLanesDir(root)
	_, orderErr := os.Stat(orderPath(root))
	_, removedErr := os.Stat(removedPath(root))
	if orderErr == nil && removedErr != nil {
		return nil // set up under the current rule; nothing to do
	}
	builtins, err := Builtins()
	if err != nil {
		return []string{fmt.Sprintf("could not read the built-in lanes to migrate %s: %v", dir, err)}
	}
	builtinByID := make(map[string]*Lane, len(builtins))
	for _, l := range builtins {
		builtinByID[l.ID] = l
	}
	removed := map[string]bool{}
	if ids, err := readIDList(removedPath(root)); err == nil {
		for _, id := range ids {
			removed[id] = true
		}
	}
	used := map[string]bool{}
	if ids, err := LoadOrder(root); err == nil && len(ids) > 0 {
		for _, id := range ids {
			used[id] = true
		}
	} else {
		for _, l := range builtins {
			used[l.ID] = true
		}
	}

	var wrote []string
	var warnings []string
	for _, l := range builtins {
		if !used[l.ID] || removed[l.ID] {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, l.ID+".md")); err == nil {
			continue
		}
		if _, err := Export(l, dir, false); err != nil {
			warnings = append(warnings, fmt.Sprintf("migrating %s: could not write %s: %v", dir, l.ID, err))
			continue
		}
		wrote = append(wrote, l.ID)
	}
	if orderErr != nil {
		// Order it as the old board showed it: shipped lanes in shipped order,
		// custom lanes by their anchors. readLaneDir returns filename order,
		// so the shipped ones are re-sequenced by Builtins() first. These Lane
		// values are for ordering only and are discarded; Load reads the
		// directory afresh.
		parsed, _ := readLaneDir(dir, nil, nil)
		byID := make(map[string]*Lane, len(parsed))
		for _, l := range parsed {
			byID[l.ID] = l
		}
		var lanes []*Lane
		for _, b := range builtins {
			if l, ok := byID[b.ID]; ok {
				l.Builtin = true
				lanes = append(lanes, l)
			}
		}
		for _, l := range parsed {
			if _, shipped := builtinByID[l.ID]; !shipped {
				lanes = append(lanes, l)
			}
		}
		ordered, _ := order(lanes)
		ids := make([]string, 0, len(ordered))
		for _, l := range ordered {
			ids = append(ids, l.ID)
		}
		if err := SaveOrder(root, ids); err != nil {
			warnings = append(warnings, fmt.Sprintf("migrating %s: could not write the order file: %v", dir, err))
		}
	}
	if removedErr == nil {
		if err := os.Remove(removedPath(root)); err != nil {
			warnings = append(warnings, fmt.Sprintf("migrating %s: could not delete the obsolete 'removed' file: %v", dir, err))
		}
	}
	msg := fmt.Sprintf("migrated %s: a board is its lane directory now, nothing is implied", dir)
	if len(wrote) > 0 {
		msg += fmt.Sprintf("; wrote the %d shipped lane(s) this board used without a file: %s", len(wrote), strings.Join(wrote, ", "))
	}
	if orderErr != nil {
		msg += "; wrote the order file"
	}
	if removedErr == nil {
		msg += "; deleted the obsolete 'removed' file"
	}
	return append([]string{msg}, warnings...)
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
			// No anchor is a statement, not a broken reference: park before
			// the terminal lane, or at the front when nothing is placed yet.
			// That last case is what a materialised lane directory looks like
			// — every lane is a file, so nothing is Builtin and out starts
			// empty — and it used to fall through into the unresolved-anchor
			// branch below, because terminalIndex of an empty list is 0 and
			// present[""] is false.
			if l.After == "" {
				idx = terminalIndex(out) - 1
			} else if idx < 0 {
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
