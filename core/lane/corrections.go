package lane

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// correctionsFileName lists the corrections that have already run against this
// board, one id per line, beside the order file (ProjectLanesDir). It is plain
// text, not markdown, so ProjectLanesActive's "*.md" glob never mistakes it for
// a lane, and it is a fact about the board rather than about any one lane.
//
// A correction writes its id here whether it edited the file, skipped it or
// found no such lane, which is the whole point: the entry is what makes a
// correction visit a board once. Somebody who writes the corrected line back
// afterwards keeps it — the board is never asked again.
const correctionsFileName = "corrections"

func correctionsPath(root string) string {
	return filepath.Join(ProjectLanesDir(root), correctionsFileName)
}

// correction is one named, once-only, field-precise repair of a lane file a
// board was given by an older build.
//
// It exists because of the tension between two rules that are both right. "A
// board is its lane directory" (743737f) says a lane file you hold IS the lane:
// nothing may quietly rewrite it from the built-ins, or a prompt you changed
// would be reset on the next load. But a lane that shipped with a defect is
// still on every board made before the fix, and telling each person to delete a
// line by hand in every checkout — .jaira/lanes/ is gitignored, so the edit does
// not travel with the repository — is not a fix.
//
// The way between them is to keep the repair as small as it can possibly be and
// to only ever apply it to a file jaira itself wrote:
//
//   - one named defect, one lane, one frontmatter field, removed once, said out
//     loud. Nothing else in the file is read or written, so a prompt, a
//     description, an added field and a lane you wrote yourself all come through
//     untouched;
//   - and the file has to be recognisable as the shipped lane it claims to
//     correct — equal to a version jaira shipped, give or take the very field
//     being corrected. Anything else is somebody's own lane and is reported,
//     never edited, however familiar its id looks. The once-only marker does not
//     cover this: it stops a correction repeating, not a correction being wrong
//     the first time.
//
// Recognition deliberately does not use creator:. stampCreatorLine writes that
// line, but only onto a lane published into a catalogue and adopted from one —
// the shipped lanes carry no creator at all — so its absence says nothing about
// who wrote the file.
type correction struct {
	// ID goes into the marker file. It never changes once shipped: changing it
	// would make the correction run a second time on boards that already had it.
	ID string
	// Lane is the lane file this touches, by lane id.
	Lane string
	// Field is the one frontmatter key it removes — the whole repair.
	Field string
	// Value limits the removal to a field carrying exactly this value.
	Value string
	// Shipped holds, verbatim, every version of this lane file jaira shipped
	// that carried the defect. A file on a board is corrected only when it
	// equals one of them apart from the Field line. Versions that never carried
	// the defect are deliberately absent: a file equal to one of those plus the
	// Field line was written by a person, since no build ever produced it.
	Shipped []string
	// Says is what the person is told when the line was removed: %s is the file.
	// It names the defect and how to get the old behaviour back, because a tool
	// that edits a file of yours and stays quiet is worse than the defect.
	Says string
	// Skipped is what the person is told when the lane is theirs, not jaira's:
	// %s is the file. Nothing was touched, so it has to say what to do by hand.
	Skipped string
}

// doneDoorway is the one version of the done lane jaira ever shipped carrying
// logbook-on-entry: true (2ecc670, superseded by 9ad7aa9). Kept here verbatim
// rather than as a hash, so that what a correction will accept as "the lane we
// shipped" can be read and reviewed in the diff like any other lane file.
const doneDoorway = `---
id: done
name: Done
after: signoff
precedence: 60
agentic: false
terminal: true
requires-outcome: true
requires-nonmodel-signal: true
requires-commits: true
logbook-on-entry: true
description: Accepted. Every definition-of-done item must be marked done, the plan finished if there is one, and the commits that carry the change recorded. The move that lands here stamps the commits and files the ticket straight into the logbook — 'jaira restore' brings it back.
---
`

// corrections is the shipped list. Append-only: a correction that has run on a
// board can never be taken back, so an entry here is history.
var corrections = []correction{{
	ID:      "done-logbook-on-entry",
	Lane:    "done",
	Field:   "logbook-on-entry",
	Value:   "true",
	Shipped: []string{doneDoorway},
	Says: "corrected %s: removed 'logbook-on-entry: true'. With it, finishing one ticket filed the " +
		"whole terminal lane into the logbook, everybody else's finished tickets with it. Finished " +
		"tickets now stay in the lane until you cut them with 'jaira logbook --all'. Write the line " +
		"back yourself if you want the old behaviour; this runs once and will not remove it again.",
	Skipped: "%s carries 'logbook-on-entry: true' and is not the lane jaira shipped, so it was left " +
		"exactly as it is. With that line, finishing one ticket files the whole terminal lane into the " +
		"logbook, everybody else's finished tickets with it; remove the line yourself to stop it. This " +
		"is said once.",
}}

// applyCorrections runs every shipped correction this board has not seen and
// returns what it did, one warning per lane it changed or refused to change. It
// is called from Load beside migrateLegacy: both repair what a board carries
// from an older build, and both must happen before the directory is read.
//
// Errors are warnings, never failures. A board that cannot be corrected — a
// read-only checkout, a lane file nobody can write — must still open.
func applyCorrections(root string) []string {
	applied, err := readIDList(correctionsPath(root))
	if err != nil {
		return []string{fmt.Sprintf("could not read %s: %v", correctionsPath(root), err)}
	}
	done := make(map[string]bool, len(applied))
	for _, id := range applied {
		done[id] = true
	}

	var warnings []string
	var ran []string
	for _, c := range corrections {
		if done[c.ID] {
			continue
		}
		path := filepath.Join(ProjectLanesDir(root), c.Lane+".md")
		b, err := os.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				// Unreadable, not absent: leave it unmarked and try again next
				// time rather than record a correction that never happened.
				warnings = append(warnings, fmt.Sprintf("correcting %s: %v", path, err))
				continue
			}
			ran = append(ran, c.ID) // no such lane on this board; nothing to correct, ever
			continue
		}
		out, dropped := dropFrontmatterLine(b, c.Field, c.Value)
		switch {
		case !dropped:
			// The defect is not in this file: already corrected by hand, never
			// present, or the field carries a different value. Nothing to say.
		case !c.recognises(b):
			warnings = append(warnings, fmt.Sprintf(c.Skipped, path))
		default:
			if err := os.WriteFile(path, out, 0o644); err != nil {
				warnings = append(warnings, fmt.Sprintf("correcting %s: %v", path, err))
				continue
			}
			warnings = append(warnings, fmt.Sprintf(c.Says, path))
		}
		ran = append(ran, c.ID)
	}

	if len(ran) > 0 {
		if err := writeIDList(root, correctionsPath(root), append(applied, ran...)); err != nil {
			warnings = append(warnings, fmt.Sprintf("could not record the corrections applied to %s: %v", ProjectLanesDir(root), err))
		}
	}
	return warnings
}

// recognises reports whether b is a lane file jaira itself wrote: equal to one
// of the shipped versions that carried the defect, apart from the corrected
// field, which is the only line either side is allowed to differ by.
//
// Comparison is on the file's bytes rather than on the parsed lane, and that is
// the point — a parsed comparison sees only the fields this binary knows, so a
// prompt somebody rewrote, a comment they left, or a field a later build will
// add would all compare equal and the file would be edited anyway.
func (c correction) recognises(b []byte) bool {
	mine, _ := dropFrontmatterLine(b, c.Field, "")
	got := normaliseLaneFile(mine)
	for _, s := range c.Shipped {
		theirs, _ := dropFrontmatterLine([]byte(s), c.Field, "")
		if got == normaliseLaneFile(theirs) {
			return true
		}
	}
	return false
}

// normaliseLaneFile makes two lane files comparable across the things a
// checkout can change without anybody meaning to: line endings, trailing
// whitespace on a line, and whether the file ends in a newline.
func normaliseLaneFile(b []byte) string {
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

// dropFrontmatterLine removes the first frontmatter line whose key is field —
// and, when value is non-empty, only when the line carries exactly that value.
// It reports whether it removed anything.
//
// A line edit rather than a YAML rewrite, for the same reason as
// stampCreatorLine: parsing and re-serialising the frontmatter would reorder
// fields, drop comments and lose anything this tool does not recognise. The
// file format is the API, and a correction must leave every part of the file it
// is not about exactly as it found it.
func dropFrontmatterLine(b []byte, field, value string) ([]byte, bool) {
	s := string(b)
	if !strings.HasPrefix(s, "---") {
		return b, false // not well-formed frontmatter; leave it alone rather than guess
	}
	lines := strings.SplitAfter(s, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(trimmed) == "---" {
			return b, false // end of the frontmatter: the field is not in it
		}
		key, rest, ok := strings.Cut(trimmed, ":")
		if !ok || strings.TrimSpace(key) != field {
			continue
		}
		if value != "" {
			got, _, _ := strings.Cut(rest, "#") // a trailing comment is not part of the value
			if strings.TrimSpace(got) != value {
				return b, false
			}
		}
		return []byte(strings.Join(append(append([]string{}, lines[:i]...), lines[i+1:]...), "")), true
	}
	return b, false
}
