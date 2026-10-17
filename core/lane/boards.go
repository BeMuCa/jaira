package lane

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// DefaultBoardName stands for the board a repository has always had, the one
// in <root>/.jaira/lanes/. That board carries no name on disk and a ticket on
// it carries no board: field, so every board from before there were several
// keeps working untouched. The name exists only for the command line and the
// TUI, where "no board given" and "today's board" have to be said out loud.
const DefaultBoardName = "default"

// boardNameRe is what a further board may be called: it is a directory name
// and a frontmatter value, so nothing that needs quoting in either.
var boardNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// ErrNoSuchBoard is returned for a board name with no lane files behind it.
var ErrNoSuchBoard = errors.New("no such board")

// NormalizeBoard maps the ways of naming today's board — nothing, or
// "default" — to "", which is how a ticket on it stores its board, and checks
// any other name. It does not look at the disk; LoadBoard does.
func NormalizeBoard(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || name == DefaultBoardName {
		return "", nil
	}
	if !boardNameRe.MatchString(name) {
		return "", fmt.Errorf("board name %q: use lower-case letters, digits and dashes", name)
	}
	return name, nil
}

// boardsDir holds every board a repository has beyond today's, one directory
// per board, each with a lanes/ directory of its own.
func boardsDir(root string) string {
	return filepath.Join(root, ticket.DirName, "boards")
}

// boardLanesDir is where the lanes of the named board name live.
func boardLanesDir(root, name string) string {
	return filepath.Join(boardsDir(root), name, "lanes")
}

// Boards lists the named boards of root, sorted. A board is there when its
// lanes/ directory holds a lane file; an empty directory is not a board yet.
func Boards(root string) []string {
	if root == "" {
		return nil
	}
	entries, err := os.ReadDir(boardsDir(root))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() || !boardNameRe.MatchString(e.Name()) || e.Name() == DefaultBoardName {
			continue
		}
		if hasLaneFiles(boardLanesDir(root, e.Name())) {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

// LoadBoard resolves the lanes of board name in root; "" or "default" is
// today's board and goes through Load unchanged.
//
// A named board is read as it stands. Unlike today's board it is never set up
// on first load: a mistyped --board would otherwise leave a new board behind.
// Nobody creates one but the person who fills its lanes/ directory, and a
// name with no lane files behind it is ErrNoSuchBoard. The migrations and
// corrections Load runs belong to today's board's history; a named board has
// none.
func LoadBoard(root, name string) (*Set, error) {
	name, err := NormalizeBoard(name)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return Load(root)
	}
	dir := boardLanesDir(root, name)
	if root == "" || !hasLaneFiles(dir) {
		return nil, fmt.Errorf("%w %q; boards here: %s", ErrNoSuchBoard, name, strings.Join(append([]string{DefaultBoardName}, Boards(root)...), ", "))
	}
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	lanes, warnings := readLaneDir(dir, nil, indexByID(builtins))
	orderIDs, _ := readIDList(filepath.Join(dir, orderFileName))
	return assemble(lanes, orderIDs, warnings), nil
}

// BoardOf is the board t runs on, "" for today's — however its board: line
// spells that.
func BoardOf(t *ticket.Ticket) string {
	b, _ := NormalizeBoard(t.Board)
	return b
}

// hasLaneFiles reports whether dir holds a lane file: what makes a directory
// under .jaira/boards/ a board.
func hasLaneFiles(dir string) bool {
	m, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	return len(m) > 0
}

// descriptionFileName holds what a board is for, in plain text beside its
// order file: .jaira/lanes/description for today's board,
// .jaira/boards/<name>/lanes/description for a named one. It is what a
// teamlead reads to put a ticket on the board its subject belongs to. Plain
// text for the same reason order is: a "*.md" file there would be a lane.
const descriptionFileName = "description"

// builtinDescription is what a board set up from the built-in lanes is for.
const builtinDescription = "Changes to this repository's code: features, fixes, refactors, tests and docs."

// writeDefaultDescription gives today's board builtinDescription, unless it
// already has a description of its own.
func writeDefaultDescription(root string) error {
	path := DescriptionPath(root, "")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, []byte(builtinDescription+"\n"), 0o644)
}

// DescriptionPath is where board name's description lives, whether or not it
// has one yet — so a listing can say where to write it.
func DescriptionPath(root, name string) string {
	name, _ = NormalizeBoard(name)
	if name == "" {
		return filepath.Join(ProjectLanesDir(root), descriptionFileName)
	}
	return filepath.Join(boardLanesDir(root, name), descriptionFileName)
}

// Description is board name's description on one line: its lines trimmed and
// joined by a space, so a description wrapped by hand still reads as the one
// sentence it is. "" when the board has none.
func Description(root, name string) string {
	lines, _ := readIDList(DescriptionPath(root, name))
	return strings.Join(lines, " ")
}

// LanesOf returns a lookup of each ticket's own board's lanes in root, each
// board loaded once, for a command that walks tickets of every board. A
// board that cannot be loaded gives nil.
func LanesOf(root string) func(*ticket.Ticket) *Set {
	boards := map[string]*Set{}
	return func(t *ticket.Ticket) *Set {
		b := BoardOf(t)
		if set, ok := boards[b]; ok {
			return set
		}
		set, _ := LoadBoard(root, b)
		boards[b] = set
		return set
	}
}
