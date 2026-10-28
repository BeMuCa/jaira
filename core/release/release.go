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
// its ticket, so two branches never edit the same lines. A file directly in
// notes/ is Unreleased; releasing moves every one of them into
// notes/<version>/ (see Cut), and the folder's name is the release.
package release

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed NOTES.md
var notesMD string

//go:embed notes
var notesFS embed.FS

// notesDir names the per-change notes inside the package directory; Cut and
// the embed share it.
const notesDir = "notes"

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
// files directly in notes/ as Unreleased, then one release per notes/<version>/
// folder, then the fixed history in NOTES.md.
func Notes() []Entry { return composeNotes(notesMD, notesFS) }

// composeNotes reads fsys as the package directory. Within a release the
// files go by name, and the releases by version number, newest first.
func composeNotes(md string, fsys fs.FS) []Entry {
	var out []Entry
	if e := entryOf(fsys, notesDir, "Unreleased"); len(e.Changes) > 0 {
		out = append(out, e)
	}
	dirs, _ := fs.ReadDir(fsys, notesDir)
	var versions []string
	for _, d := range dirs {
		if d.IsDir() {
			versions = append(versions, d.Name())
		}
	}
	sort.Slice(versions, func(i, j int) bool { return Newer(versions[i], versions[j]) })
	for _, v := range versions {
		out = append(out, entryOf(fsys, path.Join(notesDir, v), v))
	}
	return append(out, parseNotes(md)...)
}

// entryOf is one release: the changes of every *.md file directly in dir.
func entryOf(fsys fs.FS, dir, version string) Entry {
	e := Entry{Version: version}
	files, _ := fs.Glob(fsys, path.Join(dir, "*.md"))
	for _, f := range files {
		b, _ := fs.ReadFile(fsys, f)
		e.Changes = append(e.Changes, changeLines(string(b))...)
	}
	return e
}

// Newer reports whether version a is newer than b, comparing X.Y.Z numbers one by one.
func Newer(a, b string) bool {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		x, _ := strconv.Atoi(as[i])
		y, _ := strconv.Atoi(bs[i])
		if x != y {
			return x > y
		}
	}
	return len(as) > len(bs)
}

// Important splits the marker off a change: a line written "- ! text" is read
// by changeLines as the plain change "! text", and Important returns text.
// The marker is "! " with the space; "!text" is not one.
func Important(change string) (string, bool) {
	return strings.CutPrefix(change, "! ")
}

// changeLines is parseNotes for a note file, which has no headings: every
// "- " line is one change. A line "- ! text" is the change "! text": the
// marker for a change the reader must act on, read back by Important.
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
		// A release that shipped without a single note has no entry of its
		// own; it still ends the selection at the first older release.
		if e.Version == stamped || (e.Version != "Unreleased" && !Newer(e.Version, stamped)) {
			return all[:i]
		}
	}
	return all
}

// Stamped reads <dir>/version, trimmed. Errors are deliberately swallowed: a
// missing or unreadable stamp is indistinguishable from "never recorded" for
// every purpose this serves.
func Stamped(dir string) string { return readVersion(dir, "version") }

// Stamp records Current as the version that last prepared dir.
func Stamp(dir string) error { return writeCurrent(dir, "version") }

// NotesSeen reads <dir>/notes-seen: the last version whose release notes the
// user closed in the TUI, "" when none. It sits next to the stamp, per working
// tree and never in git: "I read them" is one person's fact, not the team's.
func NotesSeen(dir string) string { return readVersion(dir, "notes-seen") }

// MarkNotesSeen records Current as the version whose release notes were read.
func MarkNotesSeen(dir string) error { return writeCurrent(dir, "notes-seen") }

// readVersion reads the one-line version file <dir>/<name>, trimmed; a
// missing or unreadable file is "".
func readVersion(dir, name string) string {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// writeCurrent writes Current as the one-line version file <dir>/<name>,
// creating dir.
func writeCurrent(dir, name string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), []byte(Current+"\n"), 0o644)
}
