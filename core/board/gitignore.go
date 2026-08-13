package board

import (
	"os"
	"path/filepath"
	"strings"
)

// IgnoreLine is the entry that keeps a board private.
const IgnoreLine = "/.jaira/"

// LanesIgnoreLine keeps a project's own lane directory out of git even once
// the rest of the board is shared. .jaira/lanes/ is scoped to this machine
// (see D-03) and was never meant to reach a teammate. This is an ordinary
// ignore rule on a path outside any ignored tree, not a negation pattern: by
// the time this line matters, RemoveIgnore has already stopped ignoring
// .jaira/ itself, so there is no ignored parent to negate out of.
const LanesIgnoreLine = "/.jaira/lanes/"

// A board starts private. Committing tickets is what makes a board visible to
// everyone who can read the repository, and that should be a decision rather than
// a default — the tool has no way to know whether these notes are ready to be
// seen, or whether this repository is even yours to publish into.
//
// Private and shared differ only in whether `.jaira/` is gitignored. Nothing about
// the tickets themselves changes, so going public later is not a migration.

// Ignored reports whether the root .gitignore excludes the board.
func Ignored(root string) bool {
	b, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		switch strings.TrimSpace(line) {
		case IgnoreLine, ".jaira/", ".jaira":
			return true
		}
	}
	return false
}

// AddIgnore keeps the board out of git.
func AddIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, _ := os.ReadFile(path)
	if Ignored(root) {
		return false, nil
	}
	body := string(existing)
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	body += "\n# jaira board — private to this machine. Run 'jaira share' to publish it.\n" + IgnoreLine + "\n"
	return true, os.WriteFile(path, []byte(body), 0o644)
}

// RemoveIgnore stops excluding the board.
func RemoveIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	var out []string
	var skipNextBlank bool
	for _, line := range strings.Split(string(existing), "\n") {
		t := strings.TrimSpace(line)
		if t == IgnoreLine || t == ".jaira/" || t == ".jaira" {
			changed = true
			continue
		}
		if strings.Contains(t, "jaira board — private") {
			changed = true
			skipNextBlank = true
			continue
		}
		if skipNextBlank && t == "" {
			skipNextBlank = false
			continue
		}
		out = append(out, line)
	}
	if !changed {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}

// AddLanesIgnore keeps this project's lane directory private, independent of
// whether the rest of the board is shared. Idempotent: calling it again once
// the line is already present reports changed=false and writes nothing.
func AddLanesIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, _ := os.ReadFile(path)
	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(line) == LanesIgnoreLine {
			return false, nil
		}
	}
	body := string(existing)
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	body += "\n# this project's lane files are this machine's, not the team's\n" + LanesIgnoreLine + "\n"
	return true, os.WriteFile(path, []byte(body), 0o644)
}

// RemoveLanesIgnore removes the lanes-only ignore line. Used by 'share
// --undo': once the whole board is private again, /.jaira/ already covers
// .jaira/lanes/, and a leftover line would be a puzzle for the next reader.
func RemoveLanesIgnore(root string) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	var out []string
	var skipNextBlank bool
	for _, line := range strings.Split(string(existing), "\n") {
		t := strings.TrimSpace(line)
		if t == LanesIgnoreLine {
			changed = true
			continue
		}
		if strings.Contains(t, "lane files are this machine's") {
			changed = true
			skipNextBlank = true
			continue
		}
		if skipNextBlank && t == "" {
			skipNextBlank = false
			continue
		}
		out = append(out, line)
	}
	if !changed {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}
