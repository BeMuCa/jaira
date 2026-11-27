// Package gitwt reads git's linked-worktree bookkeeping straight off the disk.
//
// It asks git nothing: board discovery runs on every command and the project
// registry once per entry, and a git process each time is what makes a board
// slow to open on Windows. What is read is narrow - a linked worktree's .git is
// a file naming its private git directory, and only a linked worktree's private
// git directory holds a commondir file. A submodule's .git is such a file too,
// but its git directory has no commondir.
//
// A leaf package so that both core/ticket and core/project can use it: project
// imports ticket, so neither can borrow it from the other.
package gitwt

import (
	"os"
	"path/filepath"
	"strings"
)

// Top is the top of the working tree dir sits in - the nearest directory at or
// above it holding a .git entry - or "" when there is none.
func Top(dir string) string {
	for {
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

// MainCheckout is the place in the repository's main checkout that dir, inside
// a linked worktree, stands for: the same path relative to the working tree's
// top. It is "" for anything else - a main checkout, a submodule, no
// repository, a worktree whose git directory is gone, a worktree of a bare
// repository, which has no main checkout. Whether anything exists there is the
// caller's question.
func MainCheckout(dir string) string {
	top := Top(dir)
	if top == "" {
		return ""
	}
	b, err := os.ReadFile(filepath.Join(top, ".git"))
	if err != nil {
		return "" // a directory: this is a main checkout
	}
	gitdir, ok := strings.CutPrefix(strings.TrimSpace(firstLine(b)), "gitdir:")
	if !ok {
		return ""
	}
	gitdir = absUnder(top, strings.TrimSpace(gitdir))
	c, err := os.ReadFile(filepath.Join(gitdir, "commondir"))
	if err != nil {
		return "" // a submodule, or a worktree whose git directory is gone
	}
	common := absUnder(gitdir, strings.TrimSpace(firstLine(c)))
	if filepath.Base(common) != ".git" {
		return "" // a bare repository has no main checkout
	}
	rel, err := filepath.Rel(top, dir)
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(common), rel)
}

// Linked is the top of every linked worktree registered in the main checkout
// at top, read off its git directory: each one has a directory under
// .git/worktrees whose gitdir file names that worktree's .git file. Nothing
// for a top that is not a main checkout or has no linked worktrees.
func Linked(top string) []string {
	gitDir := filepath.Join(top, ".git")
	entries, err := os.ReadDir(filepath.Join(gitDir, "worktrees"))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(gitDir, "worktrees", e.Name(), "gitdir"))
		if err != nil {
			continue
		}
		wtGit := absUnder(filepath.Join(gitDir, "worktrees", e.Name()), strings.TrimSpace(firstLine(b)))
		out = append(out, filepath.Dir(wtGit))
	}
	return out
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
