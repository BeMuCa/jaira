package refsync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
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

// A fetch never resets a milestone edit that has not been sent: the queued
// write is newer than the ref, and goes out on the next flush under its lease.
func TestAFetchKeepsALocalMilestoneEditThatIsStillQueued(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const file = "backend-K3F9QZ.md"
	putRef(t, ada, "K3F9QZ", file, body)
	incoming(t, ada, grace) // grace now has the old value on disk
	edited := "---\ncolor: 99\n---\n"
	putDisk(t, grace, file, edited)
	if err := grace.syncer.RecordMilestone("K3F9QZ", file, []byte(edited)); err != nil {
		t.Fatal(err)
	}

	wrote, _, _ := incoming(t, ada, grace)
	got, _ := os.ReadFile(filepath.Join(milestone.Dir(grace.store.Root), file))
	if string(got) != edited || len(wrote) != 0 {
		t.Errorf("file is %q wrote %v, want the queued edit kept", got, wrote)
	}
	if _, err := grace.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	_, onRef, _, err := grace.syncer.Repo.ReadMilestone("K3F9QZ")
	if err != nil || string(onRef) != edited {
		t.Errorf("ref carries %q (%v) after the flush", onRef, err)
	}
}

// An old-format file beside an ID ref is merged before it is removed: a fetch
// that fails on the way leaves it for the next one.
func TestAnOldFileIsRemovedOnlyAfterTheRefIsOnDisk(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := milestone.LegacyID("alpha")
	file := "alpha-" + id + ".md"
	putRef(t, ada, id, file, body)
	putDisk(t, grace, "alpha.md", body+"\n- 01M2GGKMFB1XXK7V8AFW0YGWX1\n")
	if err := os.MkdirAll(filepath.Join(milestone.Dir(grace.store.Root), file, "x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	if _, _, _, _, err := grace.syncer.IncomingMilestones(grace.store.Root); err == nil {
		t.Fatal("writing the ID file should have failed")
	}
	if _, err := os.Stat(filepath.Join(milestone.Dir(grace.store.Root), "alpha.md")); err != nil {
		t.Errorf("the old file was removed before its content was merged: %v", err)
	}
}

// A refused milestone write is merged under the milestone lock, and a write
// another session queues while the merge waits for it is the one that is
// merged, not overwritten.
func TestARefusedMilestoneMergeWaitsForTheLockAndTakesAWriteQueuedMeanwhile(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const id, file = "K3F9QZ", "backend-K3F9QZ.md"
	putRef(t, ada, id, file, body)
	incoming(t, ada, grace)
	mine := "---\ncolor: 45\ngoal: mine\n---\n"
	putDisk(t, grace, file, mine)
	if err := grace.syncer.RecordMilestone(id, file, []byte(mine)); err != nil {
		t.Fatal(err)
	}
	lease, _ := ada.syncer.Repo.MilestoneSHA(id)
	if _, err := ada.syncer.Repo.WriteMilestone(id, file, []byte("---\ncolor: 46\n---\n"), lease); err != nil {
		t.Fatal(err)
	}

	unlock, err := grace.store.Lock(milestone.LockName)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan []refsync.Report, 1)
	go func() {
		reports, _ := grace.syncer.Flush()
		done <- reports
	}()
	select {
	case <-done:
		t.Fatal("the merge did not wait for the milestone lock")
	case <-time.After(300 * time.Millisecond):
	}
	if _, ok := grace.syncer.Box.PendingMilestone(id); !ok {
		t.Error("the refused entry was dropped before the merge replaced it")
	}
	newer := "---\ncolor: 45\ngoal: mine2\n---\n"
	if err := grace.syncer.Box.QueueMilestone(id, file, outbox.OpWrite, []byte(newer), "", "grace"); err != nil {
		t.Fatal(err)
	}
	putDisk(t, grace, file, newer)
	unlock()
	reports := <-done
	if len(reports) != 1 || !reports[0].Merged {
		t.Fatalf("reports %+v", reports)
	}
	e, ok := grace.syncer.Box.PendingMilestone(id)
	if !ok || !strings.Contains(e.Content, "goal: mine2") || !strings.Contains(e.Content, "color: 46") {
		t.Errorf("queued %+v (%v), want the newer write merged into the winner", e, ok)
	}
}

// A moving write (no lease, it carries the migration) that is refused keeps a
// rename made after the move: the migration's file is the name the write was
// based on, so the rename is one-sided and is taken, not dropped for the
// winner's name.
func TestARefusedMovingWriteKeepsTheRenameMadeAfterTheMove(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const id, old, renamed = "K3F9QZ", "alpha-K3F9QZ.md", "beta-K3F9QZ.md"
	if err := grace.syncer.RecordMigration(id, old, []byte(body), ""); err != nil {
		t.Fatal(err)
	}
	edited := "---\ncolor: 45\ngoal: mine\n---\n"
	if err := grace.syncer.RecordMilestone(id, renamed, []byte(edited)); err != nil {
		t.Fatal(err)
	}
	putDisk(t, grace, renamed, edited)
	putRef(t, ada, id, old, "---\ncolor: 46\n---\n")

	reports, err := grace.syncer.Flush()
	if err != nil || len(reports) != 1 || !reports[0].Merged {
		t.Fatalf("reports %+v, %v", reports, err)
	}
	if got := diskFiles(t, grace); got != renamed {
		t.Errorf("files %s, want only %s", got, renamed)
	}
	e, ok := grace.syncer.Box.PendingMilestone(id)
	if !ok || e.File != renamed || !strings.Contains(e.Content, "color: 46") || !strings.Contains(e.Content, "goal: mine") {
		t.Errorf("queued %+v (%v), want the rename and both edits", e, ok)
	}
	for _, c := range reports[0].Conflicts {
		if c.Field == "name" {
			t.Errorf("a one-sided rename was reported as a clash: %+v", c)
		}
	}
}

// A clone with no milestones directory yet gets an old name ref that sorts
// before its ID ref ("0-alpha" < "<ID>"), carrying an edit an older build made
// after the move: one fetch brings the ID file and merges the edit, and the
// missing directory is not an error.
func TestAnOldNameRefThatSortsFirstIsAdoptedIntoAnIDFileThatIsNotThereYet(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := milestone.LegacyID("0-alpha")
	const tk = "01M2GGKMFB1XXK7V8AFW0YGWX1"
	putRef(t, ada, id, "0-alpha-"+id+".md", body)
	putRef(t, ada, "0-alpha", "0-alpha.md", body+"\n- "+tk+"\n")
	if err := os.RemoveAll(milestone.Dir(grace.store.Root)); err != nil {
		t.Fatal(err)
	}
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	if _, _, _, _, err := grace.syncer.IncomingMilestones(grace.store.Root); err != nil {
		t.Fatalf("incoming: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(milestone.Dir(grace.store.Root), "0-alpha-"+id+".md"))
	if err != nil {
		t.Fatalf("the ID file did not arrive: %v", err)
	}
	if !strings.Contains(string(got), tk) {
		t.Errorf("the old edit is not in the ID file:\n%s", got)
	}
}

// An old file that is a state the old name ref already carried is a copy this
// clone never updated, not an edit: the merge is not run the wrong way round.
func staleCopy(t *testing.T, s0, s1, newSide string) string {
	t.Helper()
	ada, grace, _ := twoSides(t)
	id := milestone.LegacyID("alpha")
	putRef(t, ada, "alpha", "alpha.md", s0)
	putRef(t, ada, "alpha", "alpha.md", s1)
	from, err := ada.syncer.Repo.MilestoneSHA("alpha")
	if err != nil {
		t.Fatal(err)
	}
	file := "alpha-" + id + ".md"
	if _, err := ada.syncer.Repo.WriteMigration(id, file, []byte(s1), from, file, []byte(newSide)); err != nil {
		t.Fatal(err)
	}
	putDisk(t, grace, "alpha.md", s0)
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	if _, _, _, _, err := grace.syncer.IncomingMilestones(grace.store.Root); err != nil {
		t.Fatalf("incoming: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(milestone.Dir(grace.store.Root), file))
	if err != nil {
		t.Fatalf("the ID file did not arrive: %v", err)
	}
	if _, err := os.Stat(filepath.Join(milestone.Dir(grace.store.Root), "alpha.md")); err == nil {
		t.Errorf("the old file is still there")
	}
	return string(got)
}

func TestAStaleOldFileDoesNotBringARemovedTicketBack(t *testing.T) {
	const t1, t2 = "01M2GGKMFB1XXK7V8AFW0YGWX1", "01M2GGKMFB1XXK7V8AFW0YGWX2"
	s0 := body + "\n- " + t1 + "\n"
	s1 := s0 + "- " + t2 + "\n"
	got := staleCopy(t, s0, s1, s0)
	if strings.Contains(got, t2) || !strings.Contains(got, t1) {
		t.Errorf("a ticket removed with the new build came back:\n%s", got)
	}
}

func TestAStaleOldFileDoesNotDropATicketAddedBack(t *testing.T) {
	const t1, t3 = "01M2GGKMFB1XXK7V8AFW0YGWX1", "01M2GGKMFB1XXK7V8AFW0YGWX3"
	s0 := body + "\n- " + t1 + "\n- " + t3 + "\n"
	s1 := body + "\n- " + t1 + "\n"
	got := staleCopy(t, s0, s1, s0)
	if !strings.Contains(got, t3) || !strings.Contains(got, t1) {
		t.Errorf("a ticket added back with the new build was dropped:\n%s", got)
	}
}
