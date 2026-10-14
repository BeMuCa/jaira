package refsync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/milestone"
)

// The ID refs of these tests are written by hand through gitref, the way
// another clone's build would put them on the remote.

func putRef(t *testing.T, s side, key, file, body string) {
	t.Helper()
	if _, err := s.syncer.Repo.WriteMilestone(key, file, []byte(body), ""); err != nil {
		t.Fatalf("write ref %s: %v", key, err)
	}
}

func diskFiles(t *testing.T, s side) string {
	t.Helper()
	entries, _ := os.ReadDir(milestone.Dir(s.store.Root))
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, ",")
}

func putDisk(t *testing.T, s side, file, body string) {
	t.Helper()
	if err := os.MkdirAll(milestone.Dir(s.store.Root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(milestone.Dir(s.store.Root), file), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func incoming(t *testing.T, ada, grace side) (wrote, filed []string, renamed map[string]string) {
	t.Helper()
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	wrote, filed, renamed, _, err := grace.syncer.IncomingMilestones(grace.store.Root)
	if err != nil {
		t.Fatal(err)
	}
	return wrote, filed, renamed
}

const body = "---\ncolor: 45\n---\n"

// A rename is the same ref with another file name: it arrives under the new
// name and the old file, with the same ID, is gone. A file of the same name
// with another ID is another milestone and stays.
func TestARenameArrivesAndTheOldFileWithItsIDIsGone(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, grace, "backend-K3F9QZ.md", body)
	putDisk(t, grace, "frontend-M7P2RT.md", body) // another milestone, same new name
	putRef(t, ada, "K3F9QZ", "frontend-K3F9QZ.md", body)

	wrote, _, renamed := incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "frontend-K3F9QZ.md,frontend-M7P2RT.md" {
		t.Errorf("disk has %s", got)
	}
	if len(wrote) != 1 || wrote[0] != "frontend" || renamed["backend"] != "frontend" {
		t.Errorf("wrote %v renamed %v", wrote, renamed)
	}
}

// The ref decides even when its file is already there byte for byte: a copy
// under another name (a rename merged by git) is removed all the same.
func TestADoubleAfterAGitMergeIsResolvedByTheRef(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, grace, "backend-K3F9QZ.md", body)
	putDisk(t, grace, "frontend-K3F9QZ.md", body)
	putRef(t, ada, "K3F9QZ", "frontend-K3F9QZ.md", body)
	incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "frontend-K3F9QZ.md" {
		t.Errorf("disk has %s", got)
	}
}

// An old name ref is ignored once the ID ref exists, on this disk or in this
// tree; otherwise it is taken as an old file, under its old name.
func TestAnOldNameRefIsIgnoredOnceTheIDRefIsThere(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putRef(t, ada, "q4", "q4.md", body)
	putRef(t, ada, "q5", "q5.md", body)
	putRef(t, ada, "q6", "q6.md", "---\nstatus: filed\nrenamed-to: q7\n---\n")
	putRef(t, ada, milestone.LegacyID("q5"), "q5-"+milestone.LegacyID("q5")+".md", body)

	wrote, _, _ := incoming(t, ada, grace)
	want := "q4.md,q5-" + milestone.LegacyID("q5") + ".md"
	if got := diskFiles(t, grace); got != want {
		t.Errorf("disk has %s, want %s (q5 moved to its ID, q6 is a mark)", got, want)
	}
	if len(wrote) != 2 {
		t.Errorf("wrote %v", wrote)
	}
	// the ID ref of q4 arrives later: it takes q4.md's place, the ref decides
	putRef(t, ada, milestone.LegacyID("q4"), "q4-"+milestone.LegacyID("q4")+".md", body)
	incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "q4-"+milestone.LegacyID("q4")+".md,"+want[len("q4.md,"):] {
		t.Errorf("disk has %s", got)
	}
}

// A filed ID ref marks the file this tree has, and is not written where
// nothing of that milestone is.
func TestAFiledIDRefMarksTheFileHereAndNothingElse(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, grace, "backend-K3F9QZ.md", body)
	filedBody := "---\ncolor: 45\nstatus: filed\n---\n"
	putRef(t, ada, "K3F9QZ", "backend-K3F9QZ.md", filedBody)
	putRef(t, ada, "M7P2RT", "other-M7P2RT.md", filedBody)

	_, filed, _ := incoming(t, ada, grace)
	if len(filed) != 1 || filed[0] != "backend" {
		t.Errorf("filed %v", filed)
	}
	if got := diskFiles(t, grace); got != "backend-K3F9QZ.md" {
		t.Errorf("disk has %s", got)
	}
}

// A ref whose file name does not carry its ID is not a milestone.
func TestAnIDRefWhoseFileNamesAnotherIDIsSkipped(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putRef(t, ada, "K3F9QZ", "backend-M7P2RT.md", body)
	incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "" {
		t.Errorf("disk has %s", got)
	}
}

// The mark a rename of the name-keyed era left removes the old file this disk
// still has, says where it went, and is never written; with no such file there
// is nothing to do.
func TestAnOldRenameMarkRemovesTheOldFileAndIsReported(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, grace, "q6.md", body)
	putRef(t, ada, "q6", "q6.md", "---\ncolor: 45\nstatus: filed\nrenamed-to: q7\n---\n")
	putRef(t, ada, "q8", "q8.md", "---\ncolor: 45\nstatus: filed\nrenamed-to: q9\n---\n")

	wrote, filed, renamed := incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "" {
		t.Errorf("disk has %s", got)
	}
	if len(wrote)+len(filed) != 0 || len(renamed) != 1 || renamed["q6"] != "q7" {
		t.Errorf("wrote %v filed %v renamed %v", wrote, filed, renamed)
	}
}

func oldEntry(t *testing.T, s side, name, content string) {
	t.Helper()
	dir := filepath.Join(s.syncer.Box.Dir, "milestones")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `{"id":"` + name + `","op":"write","ref":"refs/jaira/milestones/` + name + `","kind":"milestones","lease":"","content":` + quote(content) + `,"queued-at":"2026-10-01T10:00:00Z","updated-at":"2026-10-01T10:00:00Z","by":"ada"}`
	if err := os.WriteFile(filepath.Join(dir, name+".json"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
}

func quote(s string) string { return `"` + strings.NewReplacer("\n", `\n`, `"`, `\"`).Replace(s) + `"` }

// A write an older build queued under a name goes out on the ID ref its file
// moves to; the mark of a rename is dropped and sent nowhere.
func TestAnOldMilestoneEntryIsRekeyedToItsID(t *testing.T) {
	ada, _, _ := twoSides(t)
	oldEntry(t, ada, "backend", "---\ncolor: 45\nafter: [ui]\n---\n")
	oldEntry(t, ada, "gone", "---\ncolor: 45\nstatus: filed\nrenamed-to: elsewhere\n---\n")

	reports, err := ada.syncer.Flush()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range reports {
		if r.Outcome != "sent" {
			t.Errorf("report %+v", r)
		}
	}
	id := milestone.LegacyID("backend")
	file, content, _, err := ada.syncer.Repo.ReadMilestone(id)
	if err != nil || file != "backend-"+id+".md" || !strings.Contains(string(content), "after: ["+milestone.LegacyID("ui")+"]") {
		t.Errorf("the ID ref carries %q %q (%v)", file, content, err)
	}
	if _, err := ada.syncer.Repo.MilestoneSHA("backend"); err == nil {
		t.Error("the old name got a ref")
	}
	if _, err := ada.syncer.Repo.MilestoneSHA("gone"); err == nil {
		t.Error("the mark of a rename was sent")
	}
	if left, _ := ada.syncer.Box.List(); len(left) != 0 {
		t.Errorf("entries left: %+v", left)
	}
}

// An older build on the same clone queues a write keyed by the file's base
// name, <name>-<ID>: it goes out on that ID's ref in that file, not as a new
// ref for a ghost milestone whose ID is the hash of the whole base name.
func TestAnOldEntryKeyedByNameAndIDKeepsItsID(t *testing.T) {
	ada, _, _ := twoSides(t)
	oldEntry(t, ada, "alpha-MB9HR0", "---\ncolor: 45\n---\n")

	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	file, _, _, err := ada.syncer.Repo.ReadMilestone("MB9HR0")
	if err != nil || file != "alpha-MB9HR0.md" {
		t.Errorf("the ID ref carries %q (%v)", file, err)
	}
	if _, err := ada.syncer.Repo.MilestoneSHA(milestone.LegacyID("alpha-MB9HR0")); err == nil {
		t.Error("the base name got a ghost ref")
	}
}

// Two clones that moved the same old file write the same first entry; the one
// that pushes second finds the ref already says it and drops its entry, with no
// report of a lost race.
func TestAnIdenticalFirstMilestoneEntryIsDroppedSilently(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := milestone.LegacyID("q4")
	file := "q4-" + id + ".md"
	if err := ada.syncer.RecordMilestone(id, file, []byte(body)); err != nil {
		t.Fatal(err)
	}
	if err := grace.syncer.RecordMilestone(id, file, []byte(body)); err != nil {
		t.Fatal(err)
	}
	if reports, err := ada.syncer.Flush(); err != nil || len(reports) != 1 || reports[0].Outcome != "sent" {
		t.Fatalf("ada: %+v %v", reports, err)
	}
	reports, err := grace.syncer.Flush()
	if err != nil || len(reports) != 0 {
		t.Errorf("grace's identical entry was reported: %+v %v", reports, err)
	}
	if left, _ := grace.syncer.Box.List(); len(left) != 0 {
		t.Errorf("entries left: %+v", left)
	}
}
