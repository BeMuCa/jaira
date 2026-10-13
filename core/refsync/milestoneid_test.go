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
	wrote, filed, renamed, err := grace.syncer.IncomingMilestones(grace.store.Root)
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
