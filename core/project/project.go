// Package project keeps a list of boards the user has opened, so the board can
// offer a switcher. It is a convenience list rather than a source of truth: a
// project is just a directory containing .jaira, and deleting this registry loses
// nothing but the shortcut.
package project

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Project is one repository the user has opened a board in.
//
// The registry exists only so the board can offer a switcher. It is a convenience
// list, not a source of truth: a project is still just a directory containing
// .jaira, and deleting this file loses nothing but the shortcut.
type Project struct {
	Root     string `json:"root"`
	Name     string `json:"name"`
	LastOpen string `json:"last_open"`
}

func projectsPath() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "projects.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "projects.json")
}

// Load returns known projects, most recently opened first, dropping any that
// no longer exist on disk.
//
// Each repository is one entry, under its main checkout: an entry recorded in a
// linked worktree before Remember mapped them, or two worktrees of one
// repository, come back as that main checkout once, with the latest open time
// among them. Boards in the temporary directory are left out (see throwaway).
func Load() []Project {
	path := projectsPath()
	if path == "" {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var ps []Project
	if err := json.Unmarshal(b, &ps); err != nil {
		return nil
	}
	var out []Project
	at := map[string]int{}
	for _, p := range ps {
		if fi, err := os.Stat(filepath.Join(p.Root, ticket.DirName)); err != nil || !fi.IsDir() {
			continue
		}
		c := canonical(p.Root)
		root := mainCheckout(c)
		if throwaway(root, path) {
			continue
		}
		if root != c {
			p.Root, p.Name = root, filepath.Base(root)
		}
		if i, ok := at[root]; ok {
			if p.LastOpen > out[i].LastOpen {
				out[i].LastOpen = p.LastOpen
			}
			continue
		}
		at[root] = len(out)
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastOpen > out[j].LastOpen })
	return out
}

// Remember records that a board was opened here and returns the root it was
// recorded under: a board in a linked worktree is recorded as its main
// checkout (see Board), and one in the temporary directory not at all, which
// returns "" and no error. An error means the registry could not be written.
//
// Every start writes this file, so two jaira processes starting at once write
// it together: the read-modify-write runs under the lock file the ticket
// store's own locks are made of, and the file is replaced whole
// (ticket.WriteAtomic) - written in place, two writers interleaved their bytes
// and left a file that no longer parsed, which Load reads as an empty list.
func Remember(root string) (string, error) {
	path := projectsPath()
	if path == "" {
		return "", errors.New("no home directory to keep the project list in")
	}
	root = Board(root)
	if throwaway(root, path) {
		return "", nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	unlock, err := ticket.LockFile(path + ".lock")
	if err != nil {
		return "", err
	}
	defer unlock()
	ps := Load()
	now := ticket.FormatTime(time.Now())
	found := false
	for i := range ps {
		if canonical(ps[i].Root) == root {
			ps[i].LastOpen = now
			found = true
		}
	}
	if !found {
		ps = append(ps, Project{Root: root, Name: filepath.Base(root), LastOpen: now})
	}
	if len(ps) > 20 {
		sort.Slice(ps, func(i, j int) bool { return ps[i].LastOpen > ps[j].LastOpen })
		ps = ps[:20]
	}
	b, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return "", err
	}
	if err := ticket.WriteAtomic(path, append(b, '\n')); err != nil {
		return "", err
	}
	return root, nil
}

// Board is the one root a board is listed under: Canonical, and for a board in
// a linked git worktree the same place in the repository's main checkout, when
// that holds a board. Every list of boards reduces its roots with it, so one
// repository is one line wherever boards are listed.
func Board(p string) string { return mainCheckout(canonical(p)) }

// Canonical is exported because every caller that compares board paths has to
// use the same reduction. The launcher deduping with filepath.Abs while the
// registry stored a symlink-resolved path listed the same board twice on macOS,
// where /var is a symlink to /private/var.
func Canonical(p string) string { return canonical(p) }

// canonical reduces a path to one spelling, so the same board added twice - once
// with a trailing slash, once through a symlink, once relative - is recognised
// as the same board rather than listed twice.
func canonical(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		abs = p
	}
	abs = filepath.Clean(abs)
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	return abs
}
