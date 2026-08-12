package project

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/BeMuCa/jaira/core/ticket"
)

// MaxScanDepth is how many directory levels below a root are searched for
// boards.
//
// Two is deliberate and not a placeholder: repositories are normally one or two
// levels under a code directory, and an unbounded walk of a home directory is
// slow enough to feel broken — it would descend into node_modules, caches and
// mounted drives to find nothing.
const MaxScanDepth = 2

// Discover finds boards at most MaxScanDepth levels below root.
//
// A board is a directory containing .jaira/tickets, matching how Discover works
// in the ticket store: the name .jaira alone is not enough, because the user's
// own configuration lives at ~/.jaira and would otherwise look like a board.
func Discover(root string, maxDepth int) []string {
	if maxDepth <= 0 {
		maxDepth = MaxScanDepth
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil
	}
	var found []string
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if isBoard(dir) {
			found = append(found, dir)
			// A board is not searched for nested boards. Tickets live inside it,
			// and a repository inside a repository is not this tool's business.
			return
		}
		if depth >= maxDepth {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() || skipDir(e.Name()) {
				continue
			}
			walk(filepath.Join(dir, e.Name()), depth+1)
		}
	}
	walk(abs, 0)
	sort.Strings(found)
	return found
}

func isBoard(dir string) bool {
	fi, err := os.Stat(filepath.Join(dir, ticket.DirName, ticket.TicketsSubdir))
	return err == nil && fi.IsDir()
}

// skipDir keeps the scan out of directories that never hold a board but do hold
// enormous numbers of files.
func skipDir(name string) bool {
	switch name {
	case "node_modules", "vendor", "target", "dist", "build", ".cache", ".venv", "venv", "__pycache__":
		return true
	}
	// Dot directories are configuration and state, not projects. .jaira itself is
	// covered by this, which is what keeps ~/.jaira from being scanned.
	return len(name) > 1 && name[0] == '.'
}
