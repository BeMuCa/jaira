package milestone

import (
	"bytes"
	"errors"
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

// Pending reports whether a file written before IDs existed lies in the
// milestone directory, which is when Migrate has something to do.
func Pending(root string) bool {
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		return false
	}
	for _, e := range entries {
		base, ok := strings.CutSuffix(e.Name(), ".md")
		if _, _, idNamed := FileIdentity(base); ok && !e.IsDir() && !idNamed {
			return true
		}
	}
	return false
}

// Result is what Migrate did: the files it moved ("old.md -> new.md"), the
// milestones that took an edit an older build had made to their old file, and
// the fields where that edit clashed with this side's.
type Result struct {
	Moved     []string
	Adopted   []string
	Conflicts []Conflict
}

// oldRefTip is what Migrate asks a Recorder that knows refs, to record the
// base: the tip SHA of the old name's ref, "" when there is none.
type oldRefTip interface{ OldRefSHA(name string) string }

// oldRefHistory is what Migrate asks a Recorder that knows refs: whether the
// old name ref has carried this very content.
type oldRefHistory interface {
	OldRefHad(name string, content []byte) bool
}

// StaleOld reports whether an old <name>.md holds a state the old name ref
// already carried. That is a copy a clone never updated, not an edit made
// here: it adopts nothing, and what the ref carries since the move is adopted
// from the ref.
func StaleOld(rec Recorder, name string, content []byte) bool {
	r, ok := rec.(oldRefHistory)
	return ok && r.OldRefHad(name, content)
}

// migrationRecorder is a Recorder that queues the move as a commit of its own
// (outbox.Box.QueueMigration), with the old ref's tip it came from.
type migrationRecorder interface {
	RecordMigration(key, file string, migrated []byte, from string) error
	MigrationBase(id string) (Base, bool)
}

func recordMigration(rec Recorder, id, file string, migrated []byte, from string) error {
	if m, ok := rec.(migrationRecorder); ok {
		return m.RecordMigration(id, file, migrated, from)
	}
	return rec.RecordMilestone(id, file, migrated)
}

// BaseFor is the base for the milestone with this ID: the one recorded here
// when this clone moved it or last merged an older build's edit, else the one
// the migration commit on the ID ref gives, so a clone that never moved it
// merges from the same state.
func BaseFor(root string, rec Recorder, id string) (Base, bool) {
	if b, ok := LoadBase(root, id); ok {
		return b, true
	}
	if m, ok := rec.(migrationRecorder); ok {
		return m.MigrationBase(id)
	}
	return Base{}, false
}

// Migrate moves every milestone file written before IDs existed to its
// <name>-<ID>.md name, once, with the ID LegacyID computes from the old name —
// the same on every clone, so clones that migrate without ever speaking to
// one another agree. The caller holds the milestone lock. It returns one
// "old.md -> new.md" per file it moved.
//
// Each move records a Base, the old content and the tip of the old name ref,
// here and as the commit of its own the ref's first write is (RecordMigration).
// An old <name>.md that turns up beside its ID file later — an older build
// wrote it again — is not simply removed: what it changed since the base is
// adopted into the ID file (Adopt), and the ref write is queued. An old file
// that is the very file the ID file is (a case-insensitive file system) is
// left alone.
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
func Migrate(root string, rec Recorder) (Result, error) {
	var res Result
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		if os.IsNotExist(err) {
			return res, nil
		}
		return res, err
	}
	var recErr *RecordError
	tip := func(name string) string {
		if r, ok := rec.(oldRefTip); ok {
			return r.OldRefSHA(name)
		}
		return ""
	}
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
			return res, err
		}
		if RenamedTo(content) != "" {
			if err := os.Remove(old); err != nil {
				return res, err
			}
			continue
		}
		target := base + "-" + id + ".md"
		targetPath := filepath.Join(Dir(root), target)
		// The ID file is found by its ID: it may have been renamed since, and
		// then the old file is still its older side, never a second milestone.
		if tInfo, err := os.Stat(idFilePath(root, id, targetPath)); err == nil {
			if oInfo, err := os.Stat(old); err == nil && os.SameFile(oInfo, tInfo) {
				// One file under two spellings: removing the old name would
				// remove the milestone.
				continue
			}
			if !StaleOld(rec, base, content) {
				adopted, conflicts, err := AdoptOld(root, rec, id, base, content, tip(base))
				var re *RecordError
				if errors.As(err, &re) {
					if recErr == nil {
						recErr = re
					}
				} else if err != nil {
					return res, err
				}
				if adopted {
					res.Adopted = append(res.Adopted, base)
				}
				res.Conflicts = append(res.Conflicts, conflicts...)
			}
			if err := os.Remove(old); err != nil {
				return res, err
			}
			res.Moved = append(res.Moved, e.Name()+" -> "+target)
			continue
		}
		migrated := MigrateContent(content)
		if err := ticket.WriteAtomic(targetPath, migrated); err != nil {
			return res, err
		}
		_ = SaveBase(root, id, Base{Content: string(migrated), OldSHA: tip(base)})
		if err := os.Remove(old); err != nil {
			return res, err
		}
		res.Moved = append(res.Moved, e.Name()+" -> "+target)
		if !rec.HasMilestoneRef(id) {
			if err := recordMigration(rec, id, target, migrated, tip(base)); err != nil {
				if recErr == nil {
					recErr = &RecordError{Name: base, Err: err}
				}
			}
		}
	}
	sort.Strings(res.Moved)
	if recErr != nil {
		return res, recErr
	}
	return res, nil
}

// idFilePath is the path of the file in the milestone directory that carries
// this ID in its name, whatever the name before it is, and fallback when there
// is none.
func idFilePath(root, id, fallback string) string {
	entries, _ := os.ReadDir(Dir(root))
	for _, e := range entries {
		if b, ok := strings.CutSuffix(e.Name(), ".md"); ok && !e.IsDir() {
			if _, got, idNamed := FileIdentity(b); idNamed && got == id {
				return filepath.Join(Dir(root), e.Name())
			}
		}
	}
	return fallback
}

// AdoptOld merges what an older build wrote to its old file or old name ref
// (old, the whole file; tip, the ref's SHA when it came from one) into the ID
// file of the milestone with this ID, and queues the ref write. Without a
// base (BaseFor) the old side is merged against an empty one. It does nothing
// when the old side is what the base says, or when no ID file lies here. The
// base then becomes the old side as it is now, so the same edit is merged, and
// reported, once.
func AdoptOld(root string, rec Recorder, id, name string, old []byte, tip string) (adopted bool, conflicts []Conflict, err error) {
	// Neither recorded here nor a migration commit on the ID ref: an empty
	// base, so everything the old side has is merged in (the union). An old
	// edit is adopted, never dropped.
	base, _ := BaseFor(root, rec, id)
	migrated := MigrateContent(old)
	if base.Content == string(migrated) {
		if tip != base.OldSHA {
			_ = SaveBase(root, id, Base{Content: base.Content, OldSHA: tip})
		}
		return false, nil, nil
	}
	path := idFilePath(root, id, "")
	if path == "" {
		return false, nil, nil
	}
	cur, err := os.ReadFile(path)
	if err != nil {
		return false, nil, err
	}
	all, err := LoadAll(root)
	if err != nil {
		return false, nil, err
	}
	merged, conflicts := Adopt(name, []byte(base.Content), old, cur, all)
	if bytes.Equal(merged, cur) {
		_ = SaveBase(root, id, Base{Content: string(migrated), OldSHA: tip})
		return false, conflicts, nil
	}
	if err := ticket.WriteAtomic(path, merged); err != nil {
		// Not written, so not merged: the base stays and the next run tries
		// again.
		return false, conflicts, err
	}
	_ = SaveBase(root, id, Base{Content: string(migrated), OldSHA: tip})
	if err := rec.RecordMilestone(id, filepath.Base(path), merged); err != nil {
		return true, conflicts, &RecordError{Name: name, Err: err}
	}
	return true, conflicts, nil
}
