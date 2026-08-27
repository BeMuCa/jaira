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

// Add brings a built-in or catalogue lane onto this board: the lane's file
// is written into the board's lane directory and its id appended to the order
// file. It refuses a lane already part of set — "already in this project" —
// rather than re-exporting it.
func Add(root string, set *Set, id string) (string, error) {
	if _, already := set.Get(id); already {
		return "", fmt.Errorf("lane %q is already part of this project", id)
	}
	installable, err := Installable(set)
	if err != nil {
		return "", err
	}
	var l *Lane
	for _, il := range installable {
		if il.ID == id {
			l = il
			break
		}
	}
	if l == nil {
		return "", fmt.Errorf("no lane %q is installed or in the catalogue", id)
	}
	dst, err := Export(l, ProjectLanesDir(root), false)
	if err != nil {
		return "", err
	}
	ids, err := effectiveOrder(root, set)
	if err != nil {
		return "", err
	}
	ids = append(ids, l.ID)
	if err := SaveOrder(root, ids); err != nil {
		return "", err
	}
	return dst, nil
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
