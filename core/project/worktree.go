package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/gitwt"
)

// tempDir is the system's temporary directory; a variable so a test can point
// it somewhere its own boards do not live.
var tempDir = os.TempDir

// mainCheckout is the board root a board in a linked git worktree stands for:
// the same place in the repository's main checkout, when that checkout holds a
// board of its own. Anything else - a main checkout, a submodule, no repository,
// a worktree of a bare repository, a main checkout without a board - is its own
// entry and comes back unchanged.
func mainCheckout(root string) string {
	main := gitwt.MainCheckout(root)
	if main == "" || !isBoard(main) {
		return root
	}
	return canonical(main)
}

// throwaway reports whether root lies in the temporary directory: agent
// scratch pads and test sandboxes, which are gone by the next start and were
// never projects anybody wants to switch to. A registry that itself lives in
// the temporary directory - a test's, a sandbox's - is a throwaway as a whole,
// and its temporary boards are the only ones it can know, so there the rule
// stands aside.
func throwaway(root, registry string) bool {
	tmp := canonical(tempDir())
	return under(root, tmp) && !under(canonical(registry), tmp)
}

// under reports whether p is dir or lies below it.
func under(p, dir string) bool {
	rel, err := filepath.Rel(dir, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// Trees is every working tree a board listed under root stands for: root
// itself, then the same place in each linked worktree of root's repository
// that holds a board. State kept per working tree - an agent's session, its
// liveness - is read across all of them, so folding a worktree into its main
// checkout's line hides none of it.
func Trees(root string) []string {
	trees := []string{root}
	top := gitwt.Top(root)
	if top == "" {
		return trees
	}
	rel, err := filepath.Rel(top, root)
	if err != nil {
		return trees
	}
	for _, wtTop := range gitwt.Linked(top) {
		if wt := filepath.Join(wtTop, rel); isBoard(wt) {
			trees = append(trees, canonical(wt))
		}
	}
	return trees
}
