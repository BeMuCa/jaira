// Package release tracks which version of jaira last prepared a board, and
// carries the change notes a user reads when a newer binary catches up an
// older one.
//
// The notes are never copied to disk anywhere: a copy can drift from the
// binary that reads it, but the binary always knows what is in itself. So the
// notes are embedded with go:embed and read straight out of the executable.
package release

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
)

//go:embed NOTES.md
var notesMD string

// Current is the version of the running binary. It is a package-level var,
// set once by cli.Execute, because the TUI also creates boards and plumbing a
// version string through the whole TUI to stamp one file is not worth it.
var Current = "dev"

// Entry is one release's worth of change notes.
type Entry struct {
	Version string
	Changes []string
}

// Notes parses the embedded NOTES.md, newest release first, in file order.
func Notes() []Entry { return parseNotes(notesMD) }

// parseNotes is a line scan: "## " starts a release, "- " inside one is a
// change. Everything else — prose, blank lines, HTML comments — is ignored,
// which keeps the file free to explain the format to a future contributor.
func parseNotes(s string) []Entry {
	var out []Entry
	for _, line := range strings.Split(s, "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			out = append(out, Entry{Version: strings.TrimSpace(strings.TrimPrefix(line, "## "))})
		case strings.HasPrefix(line, "- ") && len(out) > 0:
			last := &out[len(out)-1]
			last.Changes = append(last.Changes, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		}
	}
	return out
}

// Since returns the entries newer than stamped, selecting by position in the
// notes file rather than by parsing version numbers. A stamp that appears
// nowhere in the file — empty, "dev", a snapshot build, or a version older
// than the oldest entry — cannot be reasoned about any other way, so it
// selects everything. Showing too much on an unknown stamp is the safe
// direction; showing too little is not.
func Since(stamped string) []Entry { return sinceEntries(Notes(), stamped) }

// sinceEntries holds the selection logic apart from Notes() so tests can
// drive it against a fixture instead of the real embedded NOTES.md.
func sinceEntries(all []Entry, stamped string) []Entry {
	if stamped == "" {
		return all
	}
	for i, e := range all {
		if e.Version == stamped {
			return all[:i]
		}
	}
	return all
}

// Stamped reads <dir>/version, trimmed. Errors are deliberately swallowed: a
// missing or unreadable stamp is indistinguishable from "never recorded" for
// every purpose this serves.
func Stamped(dir string) string {
	b, err := os.ReadFile(filepath.Join(dir, "version"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// Stamp records Current as the version that last prepared dir.
func Stamp(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "version"), []byte(Current+"\n"), 0o644)
}
