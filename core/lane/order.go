package lane

import (
	"os"
	"path/filepath"
	"strings"
)

// orderFileName holds a project's column order: one lane id per line,
// position given by line number counting from 1. It is plain text, not
// markdown, so ProjectLanesActive's "*.md" glob never mistakes it for a lane
// and it lives beside the project's own lane files (ProjectLanesDir) rather
// than inside any one of them — order is a fact about the project, not about
// a lane, so a lane adopted from a teammate never arrives carrying a position
// that collides with the importing project's own layout.
const orderFileName = "order"

// orderPath is where a project's column order file lives.
func orderPath(root string) string {
	return filepath.Join(ProjectLanesDir(root), orderFileName)
}

// LoadOrder reads a project's column order file. An absent file is not an
// error: it means the project has not customised its column order, and
// callers fall back to today's after:-anchor-derived order.
func LoadOrder(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	b, err := os.ReadFile(orderPath(root))
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

// SaveOrder writes a project's column order file, one id per line.
func SaveOrder(root string, ids []string) error {
	dir := ProjectLanesDir(root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(id)
		sb.WriteByte('\n')
	}
	return os.WriteFile(orderPath(root), []byte(sb.String()), 0o644)
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
