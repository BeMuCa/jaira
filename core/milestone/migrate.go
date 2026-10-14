package milestone

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// MigrateContent is a milestone file as the ID format wants it: every
// after: token that is not an ID becomes the ID LegacyID computes for that
// name, which is the ID the milestone it names gets when its own file moves.
// Nothing else is touched, and a file with nothing to change comes back byte
// for byte — two clones migrating the same file write the same bytes.
func MigrateContent(content []byte) []byte {
	m := parse(string(content))
	if !m.Placed || !slices.ContainsFunc(m.After, func(tok string) bool { return !IsID(tok) }) {
		return content
	}
	ids := make([]string, len(m.After))
	for i, tok := range m.After {
		ids[i] = tok
		if !IsID(tok) {
			ids[i] = LegacyID(tok)
		}
	}
	m.SetAfter(ids)
	return m.Bytes()
}

// Migrate moves every milestone file written before IDs existed to its
// <name>-<ID>.md name, once, with the ID LegacyID computes from the old name —
// the same on every clone, so clones that migrate without ever speaking to
// one another agree. The caller holds the milestone lock. It returns one
// "old.md -> new.md" per file it moved.
//
// The new file is written before the old one is removed, so a crash leaves
// the milestone twice, never nowhere, and a second run finishes the move. A
// target that is already there is the newer file: only the old one goes. The
// mark a rename of the name-keyed era left under its old name is not a
// milestone and is removed. A migrated file whose ID has no ref in this tree
// is queued for it, so the move reaches the remote; the old name's ref stays
// where it is. A ref that cannot be queued does not stop the walk: the move
// is done, and only that milestone's own next write queues it, so the
// remaining files are moved too and the first such failure is returned, as a
// RecordError, after all of them.
func Migrate(root string, rec Recorder) ([]string, error) {
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var moved []string
	var recErr *RecordError
	for _, e := range entries {
		base, ok := strings.CutSuffix(e.Name(), ".md")
		if !ok || e.IsDir() {
			continue
		}
		_, id, idNamed := FileIdentity(base)
		if idNamed {
			continue
		}
		old := filepath.Join(Dir(root), e.Name())
		content, err := os.ReadFile(old)
		if err != nil {
			return moved, err
		}
		if RenamedTo(content) != "" {
			if err := os.Remove(old); err != nil {
				return moved, err
			}
			continue
		}
		target := base + "-" + id + ".md"
		if _, err := os.Stat(filepath.Join(Dir(root), target)); err == nil {
			if err := os.Remove(old); err != nil {
				return moved, err
			}
			moved = append(moved, e.Name()+" -> "+target)
			continue
		}
		migrated := MigrateContent(content)
		if err := ticket.WriteAtomic(filepath.Join(Dir(root), target), migrated); err != nil {
			return moved, err
		}
		if err := os.Remove(old); err != nil {
			return moved, err
		}
		moved = append(moved, e.Name()+" -> "+target)
		if !rec.HasMilestoneRef(id) {
			if err := rec.RecordMilestone(id, target, migrated); err != nil {
				if recErr == nil {
					recErr = &RecordError{Name: base, Err: err}
				}
			}
		}
	}
	sort.Strings(moved)
	if recErr != nil {
		return moved, recErr
	}
	return moved, nil
}
