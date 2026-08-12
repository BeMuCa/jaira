// Package project keeps a list of boards the user has opened, so the board can
// offer a switcher. It is a convenience list rather than a source of truth: a
// project is just a directory containing .jaira, and deleting this registry loses
// nothing but the shortcut.
package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/berk/jaira/core/ticket"
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

// LoadProjects returns known projects, most recently opened first, dropping any
// that no longer exist on disk.
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
	for _, p := range ps {
		if fi, err := os.Stat(filepath.Join(p.Root, ticket.DirName)); err == nil && fi.IsDir() {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastOpen > out[j].LastOpen })
	return out
}

// RememberProject records that a board was opened here.
func Remember(root string) {
	path := projectsPath()
	if path == "" {
		return
	}
	root = canonical(root)
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	if b, err := json.MarshalIndent(ps, "", "  "); err == nil {
		_ = os.WriteFile(path, append(b, '\n'), 0o644)
	}
}

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
