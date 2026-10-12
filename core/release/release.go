// Package release tracks which version of jaira last prepared a board, and
// carries the change notes a user reads when a newer binary catches up an
// older one.
//
// The notes are never copied to disk anywhere: a copy can drift from the
// binary that reads it, but the binary always knows what is in itself. So the
// notes are embedded with go:embed and read straight out of the executable.
//
// They come from two places. NOTES.md holds every release up to 0.3.6 as fixed
// history. Every change since is a file of its own under notes/, named after
// its ticket, so two branches never edit the same lines. Which release a file
// belongs to is not written in the repository at all: the release build runs
// go generate, which works it out from the tags and writes notes/versions,
// and that file is embedded with the rest. A build without it — go build,
// go install — shows every note file as Unreleased.
package release

//go:generate go run ./gen

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed NOTES.md
var notesMD string

//go:embed notes
var notesFS embed.FS

// notesDir and indexFile name the per-change notes and the release index
// inside the package directory; WriteIndex and the embed share them.
const (
	notesDir  = "notes"
	indexFile = "versions"
)

// Current is the version of the running binary. It is a package-level var,
// set once by cli.Execute, because the TUI also creates boards and plumbing a
// version string through the whole TUI to stamp one file is not worth it.
var Current = "dev"

// Entry is one release's worth of change notes.
type Entry struct {
	Version string
	Changes []string
}

// Notes returns the embedded release notes, newest release first: the note
// files nothing has assigned to a release yet as Unreleased, then the releases
// notes/versions names, then the fixed history in NOTES.md.
func Notes() []Entry { return composeNotes(notesMD, notesFS) }

// composeNotes reads fsys as the package directory: notes/*.md for the
// changes, notes/versions for which release each belongs to. Within a release
// the files keep the index's order and the unassigned ones go by name, which
// is the only order a build without git can know.
func composeNotes(md string, fsys fs.FS) []Entry {
	files, _ := fs.Glob(fsys, path.Join(notesDir, "*.md"))
	changes := map[string][]string{}
	for _, f := range files {
		b, _ := fs.ReadFile(fsys, f)
		changes[path.Base(f)] = changeLines(string(b))
	}

	var released []Entry
	assigned := map[string]bool{}
	index, _ := fs.ReadFile(fsys, path.Join(notesDir, indexFile))
	for _, line := range strings.Split(string(index), "\n") {
		version, name, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok || strings.HasPrefix(version, "#") {
			continue
		}
		c, exists := changes[name]
		if !exists {
			continue
		}
		assigned[name] = true
		if n := len(released); n == 0 || released[n-1].Version != version {
			released = append(released, Entry{Version: version})
		}
		last := &released[len(released)-1]
		last.Changes = append(last.Changes, c...)
	}

	var out []Entry
	unreleased := Entry{Version: "Unreleased"}
	for _, f := range files {
		if name := path.Base(f); !assigned[name] {
			unreleased.Changes = append(unreleased.Changes, changes[name]...)
		}
	}
	if len(unreleased.Changes) > 0 {
		out = append(out, unreleased)
	}
	out = append(out, released...)
	return append(out, parseNotes(md)...)
}

// changeLines is parseNotes for a note file, which has no headings: every
// "- " line is one change.
func changeLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "- ") {
			out = append(out, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		}
	}
	return out
}

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
//
// Entries without changes are dropped from the result. A heading with no
// points under it — once the empty "## Unreleased" every release cut left in
// NOTES.md, now a note file without a "- " line — gives a reader of "jaira
// update" nothing to read, and an entry like that also hid the "Nothing has
// changed since ..." branch from every board stamped with the current
// version. Notes() keeps it; it just does not travel.
func sinceEntries(all []Entry, stamped string) []Entry {
	var out []Entry
	for _, e := range selectSince(all, stamped) {
		if len(e.Changes) > 0 {
			out = append(out, e)
		}
	}
	return out
}

func selectSince(all []Entry, stamped string) []Entry {
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
