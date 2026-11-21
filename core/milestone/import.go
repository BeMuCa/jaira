package milestone

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// MigrateContent is a milestone document as the ID format wants it: every
// after: token that is not an ID becomes the ID LegacyID computes for that
// name, which is the ID the milestone it names was given when it moved to an
// ID. Nothing else is touched, and a document with nothing to change comes back
// byte for byte — two clones migrating the same file write the same bytes.
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

// Pending reports whether a milestone file of an older build lies in the old
// milestone directory, which is when Import has something to do. One readdir:
// it is asked on the read path of every milestone command.
func Pending(root string) bool {
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			return true
		}
	}
	return false
}

// Result is what Import did: one "<file> -> <ref>" per file it took off the
// disk, the files it left because their milestone may be on the remote
// already, and the fields where a file and its ref both said something
// different.
type Result struct {
	Moved     []string
	Waiting   []string
	Conflicts []Conflict
}

// Base is the old side of a migrated milestone as it stood when it moved to
// its ID: the migrated content, and the tip of the old name ref it came from.
// The migration commit on the ID ref carries it for every clone.
type Base struct {
	Content string `json:"content"`
	OldSHA  string `json:"old-sha,omitempty"`
}

// What Import asks of a Refs that knows more than the interface promises.
// refsync's Syncer answers all of them; a Refs that does not simply gets the
// plainer behaviour.
type (
	// refHistory: whether the ref of key has carried content at some point.
	refHistory interface {
		MilestoneHad(key string, content []byte) bool
	}
	// migrationRecorder queues a move to an ID as a commit of its own, with
	// the old name ref's tip it came from, and reads that commit back.
	migrationRecorder interface {
		RecordMigration(key, file string, migrated []byte, from string) error
		MigrationBase(id string) (Base, bool)
	}
	// oldRefs: the tip of an old name ref, and whether it has carried content.
	oldRefs interface {
		OldRefSHA(name string) string
		OldRefHad(name string, content []byte) bool
	}
)

// Import takes every milestone file an older build left under
// .jaira/milestones/ onto its ref, and deletes it. The caller holds the
// milestone lock. A file is deleted only once its content is on the ref or
// queued for it, so an import cut short leaves the file for the next one.
//
// Per file, by the ID in its name (or LegacyID for one written before IDs):
//
//   - no milestone with that ID here: the file is the milestone. It is
//     recorded — but only when mayCreate says "no ref here" really means none
//     on the remote; offline it may be a teammate's milestone this clone has
//     not fetched yet, and the file waits for a fetch that reaches the remote;
//   - the ref carries the same, or has carried it before: the file is a stale
//     copy — a branch merged in, or a commit made by mistake — and goes. This
//     is the case that used to overwrite the ref;
//   - anything else is an edit the ref never saw. It is merged into the ref's
//     version (Adopt): tickets either side added stay, a field both changed
//     keeps the ref's value and the file's is written in as a comment and
//     reported. Nothing in the file is lost.
//
// The mark a rename of the name-keyed era left is not a milestone and goes.
// A file that cannot be read or recorded is left where it is and does not stop
// the others; the first such error comes back after all of them.
func Import(root string, refs Refs, mayCreate bool) (Result, error) {
	var res Result
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		if os.IsNotExist(err) {
			return res, nil
		}
		return res, err
	}
	var files []diskFile
	for _, e := range entries {
		base, ok := strings.CutSuffix(e.Name(), ".md")
		if !ok || e.IsDir() {
			continue
		}
		f := diskFile{base: base}
		f.name, f.id, f.idNamed = FileIdentity(base)
		files = append(files, f)
	}
	// A file named with its ID first: of two with one ID, it is the newer.
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].idNamed != files[j].idNamed {
			return files[i].idNamed
		}
		return files[i].base < files[j].base
	})
	if len(files) == 0 {
		return res, nil
	}

	current := map[string]File{}
	have, err := refs.MilestoneFiles()
	if err != nil {
		return res, err
	}
	for _, f := range have {
		_, id, _ := FileIdentity(strings.TrimSuffix(f.Name, ".md"))
		if _, ok := current[id]; !ok {
			current[id] = f
		}
	}
	all, _ := scan(refs)
	var first error
	keep := func(err error) {
		if first == nil {
			first = err
		}
	}

	for _, f := range files {
		path := filepath.Join(Dir(root), f.base+".md")
		content, err := os.ReadFile(path)
		if err != nil {
			keep(err)
			continue
		}
		if RenamedTo(content) != "" {
			if err := os.Remove(path); err != nil {
				keep(err)
			}
			continue
		}
		migrated := content
		if !f.idNamed {
			migrated = MigrateContent(content)
		}
		target := f.name + "-" + f.id + ".md"
		cur, known := current[f.id]
		switch {
		case known && (bytes.Equal(cur.Content, migrated) || hadBefore(refs, f, content, migrated)):
			// A copy of a state the ref has had: nothing to take.
		case known:
			var base []byte
			if m, ok := refs.(migrationRecorder); ok && !f.idNamed {
				if b, ok := m.MigrationBase(f.id); ok {
					base = []byte(b.Content)
				}
			}
			merged, cs := Adopt(f.name, base, content, cur.Content, all)
			if !bytes.Equal(merged, cur.Content) {
				if err := refs.RecordMilestone(f.id, cur.Name, merged); err != nil {
					keep(err)
					continue
				}
				current[f.id] = File{Name: cur.Name, Content: merged}
			}
			res.Conflicts = append(res.Conflicts, cs...)
		case !mayCreate:
			res.Waiting = append(res.Waiting, f.base+".md")
			continue
		default:
			var err error
			if m, ok := refs.(migrationRecorder); ok && !f.idNamed {
				tip := ""
				if o, ok := refs.(oldRefs); ok {
					tip = o.OldRefSHA(f.name)
				}
				err = m.RecordMigration(f.id, target, migrated, tip)
			} else {
				err = refs.RecordMilestone(f.id, target, migrated)
			}
			if err != nil {
				keep(err)
				continue
			}
			current[f.id] = File{Name: target, Content: migrated}
		}
		if err := os.Remove(path); err != nil {
			keep(err)
			continue
		}
		res.Moved = append(res.Moved, f.base+".md -> "+RefPrefix+f.id)
	}
	// Gone when it is empty, so nothing is left to commit or to look at.
	_ = os.Remove(Dir(root))
	return res, first
}

// diskFile is one file in the old milestone directory, read for what its name
// says: its base name, and the milestone name and ID it stands for — the ID in
// the name, or the one LegacyID computes for a file written before IDs.
type diskFile struct {
	base, name, id string
	idNamed        bool
}

// hadBefore reports whether a file holds a state the milestone's ref, or the
// old name ref of a file written before IDs, has carried.
func hadBefore(refs Refs, f diskFile, content, migrated []byte) bool {
	if h, ok := refs.(refHistory); ok && h.MilestoneHad(f.id, migrated) {
		return true
	}
	if o, ok := refs.(oldRefs); ok && !f.idNamed && o.OldRefHad(f.name, content) {
		return true
	}
	return false
}
