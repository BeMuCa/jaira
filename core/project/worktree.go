package project

import (
	"os"
	"path/filepath"
	"strings"
)

// tempDir is the system's temporary directory; a variable so a test can point
// it somewhere its own boards do not live.
var tempDir = os.TempDir

// mainCheckout is the board root a board in a linked git worktree stands for:
// the same place in the repository's main checkout, when that checkout holds a
// board of its own. Anything else - a main checkout, a submodule, no repository,
// a worktree of a bare repository, a main checkout without a board - is its own
// entry and comes back unchanged.
//
// Read off the disk rather than asked of git: the registry is read on every
// start, once per entry, and a git process per entry is what makes a board slow
// to open on Windows. What is read is narrow - a linked worktree's .git is a file
// naming its private git directory, and only a linked worktree's private git
// directory holds a commondir file. A submodule's .git is such a file too, but
// its git directory has no commondir, so it stays its own entry.
func mainCheckout(root string) string {
	top := enclosingGit(root)
	if top == "" {
		return root
	}
	b, err := os.ReadFile(filepath.Join(top, ".git"))
	if err != nil {
		return root // a directory: this is a main checkout
	}
	gitdir, ok := strings.CutPrefix(strings.TrimSpace(firstLine(b)), "gitdir:")
	if !ok {
		return root
	}
	gitdir = absUnder(top, strings.TrimSpace(gitdir))
	c, err := os.ReadFile(filepath.Join(gitdir, "commondir"))
	if err != nil {
		return root // a submodule, or a worktree whose git directory is gone
	}
	common := absUnder(gitdir, strings.TrimSpace(firstLine(c)))
	if filepath.Base(common) != ".git" {
		return root // a bare repository has no main checkout
	}
	rel, err := filepath.Rel(top, root)
	if err != nil {
		return root
	}
	main := filepath.Join(filepath.Dir(common), rel)
	if !isBoard(main) {
		return root
	}
	return canonical(main)
}

// enclosingGit is the top of the working tree root sits in - the nearest
// directory at or above it holding a .git entry - or "" when there is none.
func enclosingGit(root string) string {
	for dir := root; ; {
		if _, err := os.Lstat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func firstLine(b []byte) string {
	s, _, _ := strings.Cut(string(b), "\n")
	return strings.TrimSuffix(s, "\r")
}

// absUnder resolves p against base when git wrote it relative.
func absUnder(base, p string) string {
	p = filepath.FromSlash(p)
	if !filepath.IsAbs(p) {
		p = filepath.Join(base, p)
	}
	return filepath.Clean(p)
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
//
// Read off the main checkout's git directory, as mainCheckout reads the other
// direction: each linked worktree has a directory under .git/worktrees whose
// gitdir file names that worktree's .git file.
func Trees(root string) []string {
	trees := []string{root}
	top := enclosingGit(root)
	if top == "" {
		return trees
	}
	rel, err := filepath.Rel(top, root)
	if err != nil {
		return trees
	}
	gitDir := filepath.Join(top, ".git")
	entries, err := os.ReadDir(filepath.Join(gitDir, "worktrees"))
	if err != nil {
		return trees // not a main checkout, or no linked worktrees
	}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(gitDir, "worktrees", e.Name(), "gitdir"))
		if err != nil {
			continue
		}
		wtGit := absUnder(filepath.Join(gitDir, "worktrees", e.Name()), strings.TrimSpace(firstLine(b)))
		if wt := filepath.Join(filepath.Dir(wtGit), rel); isBoard(wt) {
			trees = append(trees, canonical(wt))
		}
	}
	return trees
}
