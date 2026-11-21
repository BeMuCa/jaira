package milestone

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// refsBeside is a fake Refs that is not the old milestone directory, so what
// Import takes off the disk and what it puts on the refs can be told apart.
func refsBeside(root string) *fakeRec { return &fakeRec{dir: filepath.Join(root, "refs")} }

func refFile(t *testing.T, rec *fakeRec, file string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(rec.dir, file))
	if err != nil {
		t.Fatalf("%s is not on its ref: %v", file, err)
	}
	return string(b)
}

// Every file an older build left is taken onto its ref and deleted, the
// directory with it; an old file is given the ID LegacyID computes and its
// after: names become IDs, so two clones importing the same file write the
// same bytes. A second run has nothing to do.
func TestImportTakesEveryFileOntoItsRefAndDeletesIt(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "q4", "---\ncolor: 45\n---\n\n# q4\n\n- "+idA+"\n")
	writeOld(t, root, "q5", "---\ncolor: 46\nafter: [q4, ghost]\n---\n")
	writeFile(t, root, "new-K3F9QZ.md", "---\ncolor: 47\n---\n")
	rec := refsBeside(root)
	res, err := Import(root, rec, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Moved) != 3 || res.Moved[0] != "new-K3F9QZ.md -> "+RefPrefix+"K3F9QZ" {
		t.Errorf("moved %v", res.Moved)
	}
	if _, err := os.Stat(Dir(root)); !os.IsNotExist(err) {
		t.Errorf("the old directory is still there: %v", err)
	}
	q5 := refFile(t, rec, "q5-"+LegacyID("q5")+".md")
	if want := "after: [" + LegacyID("q4") + ", " + LegacyID("ghost") + "]"; !strings.Contains(q5, want) {
		t.Errorf("after: not migrated, want %s in\n%s", want, q5)
	}
	if !strings.Contains(refFile(t, rec, "q4-"+LegacyID("q4")+".md"), idA) {
		t.Error("q4 lost its ticket")
	}
	if again, err := Import(root, rec, true); err != nil || len(again.Moved) != 0 {
		t.Errorf("second run: %+v, %v", again, err)
	}
}

func TestMigrateContentChangesNothingElse(t *testing.T) {
	in := "---\ncolor: 46\nafter: [K3F9QZ]\n---\n\nprose --- stays\n"
	if got := MigrateContent([]byte(in)); string(got) != in {
		t.Errorf("a file with IDs only changed:\n%s", got)
	}
	in = "---\ncolor: 46\n---\n"
	if got := MigrateContent([]byte(in)); string(got) != in {
		t.Errorf("an unplaced file changed:\n%s", got)
	}
}

// historyRec is a fake whose refs remember earlier states.
type historyRec struct {
	*fakeRec
	had map[string]string // ID -> a content its ref carried before
}

func (h historyRec) MilestoneHad(key string, content []byte) bool {
	return h.had[key] == string(content)
}

// A file holding a state its ref has carried — a copy merged in from a branch,
// or committed by mistake — is stale and goes without a write. This is the case
// that used to overwrite the ref with an old ticket list.
func TestImportDeletesAStaleCopyWithoutWriting(t *testing.T) {
	root := t.TempDir()
	old := "---\ncolor: 45\n---\n- " + idA + "\n"
	now := "---\ncolor: 45\n---\n- " + idB + "\n"
	rec := historyRec{fakeRec: refsBeside(root), had: map[string]string{"K3F9QZ": old}}
	if err := rec.RecordMilestone("K3F9QZ", "a-K3F9QZ.md", []byte(now)); err != nil {
		t.Fatal(err)
	}
	rec.names = nil
	writeFile(t, root, "a-K3F9QZ.md", old)
	writeFile(t, root, "same-AAAAAA.md", "x\n")
	if err := rec.RecordMilestone("AAAAAA", "same-AAAAAA.md", []byte("x\n")); err != nil {
		t.Fatal(err)
	}
	rec.names = nil
	res, err := Import(root, rec, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.names) != 0 || len(res.Moved) != 2 {
		t.Errorf("recorded %v, moved %v: a stale copy is deleted, not written", rec.names, res.Moved)
	}
	if got := refFile(t, rec.fakeRec, "a-K3F9QZ.md"); got != now {
		t.Errorf("the ref was overwritten by the stale copy:\n%s", got)
	}
}

// A file the ref never saw is an edit, and nothing of it is lost: its tickets
// join the ref's, and a field both say differently keeps the ref's value and
// carries the file's as a comment.
func TestImportMergesAnEditTheRefNeverSaw(t *testing.T) {
	root := t.TempDir()
	rec := refsBeside(root)
	if err := rec.RecordMilestone("K3F9QZ", "a-K3F9QZ.md", []byte("---\ncolor: 45\n---\n- "+idA+"\n")); err != nil {
		t.Fatal(err)
	}
	rec.names = nil
	writeFile(t, root, "a-K3F9QZ.md", "---\ncolor: 46\n---\n- "+idB+"\n")
	res, err := Import(root, rec, true)
	if err != nil {
		t.Fatal(err)
	}
	got := refFile(t, rec, "a-K3F9QZ.md")
	if !strings.Contains(got, idA) || !strings.Contains(got, idB) || !strings.Contains(got, "color: 45") {
		t.Errorf("merged:\n%s", got)
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Field != "color" || len(rec.names) != 1 {
		t.Errorf("conflicts %v, recorded %v", res.Conflicts, rec.names)
	}
	if _, err := os.Stat(Dir(root)); !os.IsNotExist(err) {
		t.Errorf("the file was left: %v", err)
	}
}

// Offline, a milestone with no ref here may be a teammate's this clone has not
// fetched: its file waits, and nothing is queued against "no ref".
func TestImportWaitsWhenItMayNotCreate(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-K3F9QZ.md", "---\ncolor: 45\n---\n")
	rec := refsBeside(root)
	res, err := Import(root, rec, false)
	if err != nil || len(res.Waiting) != 1 || len(res.Moved) != 0 || len(rec.names) != 0 {
		t.Fatalf("res %+v, err %v, recorded %v", res, err, rec.names)
	}
	if dirNames(t, root) != "a-K3F9QZ.md" {
		t.Error("the waiting file was removed")
	}
}

// A file that cannot be recorded stays on disk for the next run and does not
// stop the others; the mark a rename of the name-keyed era left goes.
func TestImportKeepsAFileItCouldNotRecord(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "q4", "---\ncolor: 45\n---\n")
	writeOld(t, root, "old", "---\ncolor: 45\nrenamed-to: new\n---\n")
	boom := errors.New("no repository")
	rec := refsBeside(root)
	rec.err = boom
	if _, err := Import(root, rec, true); !errors.Is(err, boom) {
		t.Fatalf("err = %v", err)
	}
	if dirNames(t, root) != "q4.md" {
		t.Errorf("directory %s, want q4.md kept and the mark gone", dirNames(t, root))
	}
}
