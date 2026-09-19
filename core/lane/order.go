package lane

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// orderFileName holds a project's column order: one lane id per line,
// position given by line number counting from 1. It is plain text, not
// markdown, so ProjectLanesActive's "*.md" glob never mistakes it for a lane
// and it lives beside the project's own lane files (ProjectLanesDir) rather
// than inside any one of them — order is a fact about the project, not about
// a lane, so a lane adopted from a teammate never arrives carrying a position
// that collides with the importing project's own layout.
const orderFileName = "order"

// removedFileName is the file a board from before "a board is its lane
// directory" used to list the built-ins it left out — needed then because Load
// injected every built-in under the directory's files. Nothing writes it any
// more; migrateLegacy reads it once and deletes it.
const removedFileName = "removed"

func orderPath(root string) string   { return filepath.Join(ProjectLanesDir(root), orderFileName) }
func removedPath(root string) string { return filepath.Join(ProjectLanesDir(root), removedFileName) }

// readIDList reads a plain-text, one-id-per-line file. An absent file is not
// an error: it returns a nil slice.
func readIDList(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			ids = append(ids, line)
		}
	}
	return ids, nil
}

// writeIDList writes ids, one per line, creating the project's lane
// directory if needed.
func writeIDList(root, path string, ids []string) error {
	if err := os.MkdirAll(ProjectLanesDir(root), 0o755); err != nil {
		return err
	}
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(id)
		sb.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// LoadOrder reads a project's column order file. An absent file is not an
// error: it means the project has not customised its column order, and
// callers fall back to today's after:-anchor-derived order.
func LoadOrder(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	return readIDList(orderPath(root))
}

// SaveOrder writes a project's column order file, one id per line.
func SaveOrder(root string, ids []string) error {
	return writeIDList(root, orderPath(root), ids)
}

// Move returns a copy of ids with id shifted delta positions — -1 is one
// step toward the front, +1 one step toward the back. It is the user's model
// exactly: moving a lane one step swaps it with its neighbour. Moving past
// either end is a no-op, not an error and not a wrap-around: asking to move
// further than the board allows is not the same as asking for a mistake.
func Move(ids []string, id string, delta int) []string {
	out := append([]string{}, ids...)
	i := indexOfID(out, id)
	if i < 0 {
		return out
	}
	j := i + delta
	if j < 0 || j >= len(out) {
		return out
	}
	out[i], out[j] = out[j], out[i]
	return out
}

func indexOfID(ids []string, id string) int {
	for i, v := range ids {
		if v == id {
			return i
		}
	}
	return -1
}

// applyOrder reorders lanes (already after:-resolved by order()) per ids: a
// lane named in ids takes that position, in ids' own sequence; a lane present
// in lanes but missing from ids is appended after all the named ones, in the
// order lanes already has it — so hand-editing the order file can never make
// a lane vanish from the board. An id in ids with no lane behind it produces
// a warning and is otherwise skipped.
func applyOrder(lanes []*Lane, ids []string) ([]*Lane, []string) {
	byID := make(map[string]*Lane, len(lanes))
	for _, l := range lanes {
		byID[l.ID] = l
	}
	var warnings []string
	seen := make(map[string]bool, len(ids))
	out := make([]*Lane, 0, len(lanes))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		l, ok := byID[id]
		if !ok {
			warnings = append(warnings, "order file names lane \""+id+"\", which is not installed")
			continue
		}
		out = append(out, l)
	}
	for _, l := range lanes {
		if !seen[l.ID] {
			out = append(out, l)
		}
	}
	return out, warnings
}

// effectiveOrder is the column order a mutation should build on: the order
// file if one exists, otherwise the currently loaded set's own order.
func effectiveOrder(root string, set *Set) ([]string, error) {
	ids, err := LoadOrder(root)
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		return ids, nil
	}
	out := make([]string, 0, len(set.Lanes))
	for _, l := range set.Lanes {
		out = append(out, l.ID)
	}
	return out, nil
}

// withoutID returns ids with id removed, preserving order.
func withoutID(ids []string, id string) []string {
	out := make([]string, 0, len(ids))
	for _, v := range ids {
		if v != id {
			out = append(out, v)
		}
	}
	return out
}

// Installable lists every built-in and catalogue lane not already part of
// set — the "add a lane to this project" catalogue, for both 'jaira lanes
// add' and the settings screen's '+' column. A file that fails to read or
// parse is skipped rather than failing the whole listing, matching Shared's
// treatment of the same class of problem.
func Installable(set *Set) ([]*Lane, error) {
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(set.Lanes))
	for _, l := range set.Lanes {
		seen[l.ID] = true
	}
	var out []*Lane
	add := func(l *Lane) {
		if !seen[l.ID] {
			seen[l.ID] = true
			out = append(out, l)
		}
	}
	for _, l := range builtins {
		add(l)
	}
	matches, _ := filepath.Glob(filepath.Join(UserLanesDir(), "*.md"))
	sort.Strings(matches)
	for _, m := range matches {
		b, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		l, err := parse(b, m, false)
		if err != nil {
			continue
		}
		add(l)
	}
	return out, nil
}

// Add brings a built-in or catalogue lane onto this board: the lane's file is
// written into the board's lane directory and its id placed where its after:
// field says it belongs. It refuses a lane already part of set — "already in
// this project" — rather than re-exporting it. The returned warnings say what
// the placement had to assume, the same way Load reports an anchor it could
// not resolve.
//
// The second return is the id the lane now follows, empty when it went to the
// front. Placement is the whole point of the call and the one thing a caller
// cannot see: a chain resolved through uninstalled lanes lands 'jaira lanes
// add testing' between in-progress and human, and silently, because warning on
// the path the board itself advertises would be noise. So the neighbour is
// handed back to be said out loud instead.
func Add(root string, set *Set, id string) (string, string, []string, error) {
	if _, already := set.Get(id); already {
		return "", "", nil, fmt.Errorf("lane %q is already part of this project", id)
	}
	installable, err := Installable(set)
	if err != nil {
		return "", "", nil, err
	}
	var l *Lane
	for _, il := range installable {
		if il.ID == id {
			l = il
			break
		}
	}
	if l == nil {
		return "", "", nil, fmt.Errorf("no lane %q is installed or in the catalogue", id)
	}
	dst, err := Export(l, ProjectLanesDir(root), false)
	if err != nil {
		return "", "", nil, err
	}
	ids, err := effectiveOrder(root, set)
	if err != nil {
		return "", "", nil, err
	}
	ids, warnings := insertAfterAnchor(ids, l, set, installable)
	if err := SaveOrder(root, ids); err != nil {
		return "", "", nil, err
	}
	var after string
	for i, got := range ids {
		if got == id && i > 0 {
			after = ids[i-1]
		}
	}
	return dst, after, warnings, nil
}

// insertAfterAnchor places a newly added lane where its after: field says it
// belongs, rather than at the end of the board. Appending was harmless while
// every shipped lane was already installed and only custom lanes arrived this
// way; with critique, optimize and testing shipping uninstalled, appending
// puts a review loop behind done and blocked, where no ticket ever reaches it
// — a lane installed but out of the flow is not an installed lane.
//
// Only the new id moves: the rest of the order is left exactly as it is,
// because a board's column order is the user's arrangement and adding one
// lane is no reason to re-derive it.
//
// Once nothing can resolve the anchor, the fallback is order()'s, decided the
// same way it decides it, because a lane must not land in one place when Load
// derives the order and another when Add writes it:
//
//   - no anchor at all is a statement, not an omission: park the lane before
//     the terminal lane, where work still flows through it.
//   - an anchor this board does not have is the same placement plus a
//     warning. Appending here was the bug this function was written against,
//     one board removed or renamed away from biting: 'jaira lanes remove
//     in-progress' followed by 'jaira lanes add critique' put critique behind
//     done and blocked, installed and unreachable.
//
// The anchor may also be a lane that ships uninstalled, which is the common
// case now rather than a corner: 'jaira lanes add testing' on a fresh board
// names optimize, which names critique, which names in-progress. Only the last
// of those is on the board, and stopping at the first missing name would park
// testing before the terminal lane — behind signoff, a test lane after the
// human acceptance. So the chain is followed through the offer until it
// reaches a lane the board has; a name that is nowhere is the unresolvable
// case above, warning and all. That step is Add's alone and no break with
// order(): order() sees only the lanes on the board and has no chain to
// follow, while Add is holding the offer the lane came out of.
func insertAfterAnchor(ids []string, l *Lane, set *Set, installable []*Lane) ([]string, []string) {
	var warnings []string
	at := anchorIndex(ids, l, installable)
	if at < 0 {
		if l.After != "" {
			// l.After, not the name the chain ended on: the user wrote this one,
			// and a warning naming a link they never typed — or the empty string
			// a chain ending in an anchor-less lane leaves behind — reads as a bug
			// in jaira rather than a lane this board does not have.
			warnings = append(warnings, fmt.Sprintf(
				"lane %s: anchor %q is not on this board; placed before the terminal lane",
				l.ID, l.After))
		}
		at = terminalIDIndex(ids, set)
	}
	out := make([]string, 0, len(ids)+1)
	out = append(out, ids[:at]...)
	out = append(out, l.ID)
	return append(out, ids[at:]...), warnings
}

// anchorIndex resolves where l's after: chain lands in ids: the position just
// past the first anchor the board actually has, or -1 when nothing in the
// chain is on this board. Following the chain through the offer is what a board that has
// not installed the whole review loop needs — its links are lanes it does not
// have yet, and each one knows where the next belongs.
func anchorIndex(ids []string, l *Lane, installable []*Lane) int {
	byID := make(map[string]*Lane, len(installable))
	for _, il := range installable {
		byID[il.ID] = il
	}
	// A chain that loops — two uninstalled lanes naming each other — would spin
	// here, so every link is visited once.
	seen := map[string]bool{l.ID: true}
	anchor := l.After
	for anchor != "" && !seen[anchor] {
		seen[anchor] = true
		for i, id := range ids {
			if id == anchor {
				return i + 1
			}
		}
		next, ok := byID[anchor]
		if !ok {
			break
		}
		anchor = next.After
	}
	return -1
}

// terminalIDIndex is terminalIndex over a list of ids: the position of the
// first terminal lane, or the end when this board has none. Placing before it
// is what keeps an anchor-less lane in the flow instead of behind done.
func terminalIDIndex(ids []string, set *Set) int {
	for i, id := range ids {
		if l, ok := set.Get(id); ok && l.Terminal {
			return i
		}
	}
	return len(ids)
}

// ticketsIn lists the handles of tickets currently sitting in lane id, so
// Remove's refusal can name them.
func ticketsIn(store *ticket.Store, id string) ([]string, error) {
	tickets, err := store.List()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, t := range tickets {
		if t.Status == id {
			out = append(out, ticket.Handle(t.ID))
		}
	}
	return out, nil
}

// Remove takes a lane off this board — never out of the catalogue — refusing
// when any ticket currently sits in it, naming them: a lane that vanishes
// under a ticket leaves it in a lane nothing knows. The board is its lane
// directory, so removing is deleting the file and its line in the order file;
// a shipped lane stays on offer in the catalogue and 'jaira lanes add' brings
// it back.
func Remove(root string, set *Set, store *ticket.Store, id string) (string, error) {
	held, err := ticketsIn(store, id)
	if err != nil {
		return "", err
	}
	if len(held) > 0 {
		return "", fmt.Errorf("lane %q holds %d ticket(s) and cannot be removed: %s",
			id, len(held), strings.Join(held, ", "))
	}
	l, ok := set.Get(id)
	if !ok {
		return "", fmt.Errorf("no lane %q is part of this project", id)
	}
	path := filepath.Join(ProjectLanesDir(root), l.ID+".md")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	ids, err := effectiveOrder(root, set)
	if err != nil {
		return "", err
	}
	ids = withoutID(ids, l.ID)
	if err := SaveOrder(root, ids); err != nil {
		return "", err
	}
	return path, nil
}

// MoveLane shifts id one step in this project's column order and writes the
// order file — the one implementation the CLI's 'lanes move' and the
// settings screen's H/L keys both call, so "move a lane" never has two
// bodies to drift apart.
func MoveLane(root string, set *Set, id string, delta int) error {
	if _, ok := set.Get(id); !ok {
		return fmt.Errorf("no lane %q is part of this project", id)
	}
	ids, err := effectiveOrder(root, set)
	if err != nil {
		return err
	}
	ids = Move(ids, id, delta)
	return SaveOrder(root, ids)
}
