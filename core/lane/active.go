package lane

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// activeBoardsFileName lists the boards a teamlead may put new tickets on and
// shift+up/down on the board steps through: one board name per line, today's
// board written as "default". It lives in the repository rather than in a
// person's settings because which board a ticket belongs on is the team's
// rule — two teamleads reading two private lists would file the same ticket
// onto two different boards. Plain text beside the tickets, so it is read and
// fixed by hand and merges line by line.
const activeBoardsFileName = "active-boards"

// ErrLastActiveBoard refuses to switch off the only board still active: a
// teamlead would have nowhere left to put a ticket.
var ErrLastActiveBoard = errors.New("the last active board stays active")

// ActiveBoardsPath is where root's list of active boards lives, whether or not
// it exists yet.
func ActiveBoardsPath(root string) string {
	return filepath.Join(root, ticket.DirName, activeBoardsFileName)
}

// AllBoards is every board of root in the one fixed order the board window
// and shift+up/down share: today's board ("") first, then the others by
// name.
func AllBoards(root string) []string {
	return append([]string{""}, Boards(root)...)
}

// ActiveBoards is root's active boards in AllBoards order. Without the file
// every board is active, so a repository that never chose keeps working as
// before; so is a file naming no board that still exists, for the same
// reason. A name in the file with no board behind it is left out.
func ActiveBoards(root string) []string {
	all := AllBoards(root)
	names, _ := readIDList(ActiveBoardsPath(root))
	if len(names) == 0 {
		return all
	}
	listed := map[string]bool{}
	for _, n := range names {
		if b, err := NormalizeBoard(n); err == nil {
			listed[b] = true
		}
	}
	var out []string
	for _, b := range all {
		if listed[b] {
			out = append(out, b)
		}
	}
	if len(out) == 0 {
		return all
	}
	return out
}

// SetBoardActive switches board name on or off in root's list of active
// boards and writes the whole list back in AllBoards order. Switching off the
// last active board is ErrLastActiveBoard and writes nothing.
func SetBoardActive(root, name string, on bool) error {
	name, err := NormalizeBoard(name)
	if err != nil {
		return err
	}
	active := map[string]bool{}
	for _, b := range ActiveBoards(root) {
		active[b] = true
	}
	active[name] = on
	var lines []string
	for _, b := range AllBoards(root) {
		if !active[b] {
			continue
		}
		if b == "" {
			b = DefaultBoardName
		}
		lines = append(lines, b)
	}
	if len(lines) == 0 {
		return ErrLastActiveBoard
	}
	return os.WriteFile(ActiveBoardsPath(root), []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
