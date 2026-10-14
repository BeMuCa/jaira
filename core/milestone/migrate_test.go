package milestone

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func dirNames(t *testing.T, root string) string {
	t.Helper()
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, ",")
}

// writeOld puts a file where the format before IDs kept it: <name>.md.
func writeOld(t *testing.T, root, name, body string) { writeFile(t, root, name+".md", body) }

// Two clones migrating the same old file write the same bytes under the same
// name: the ID is computed, and after: names become the IDs those names get.
func TestMigrateGivesEveryCloneTheSameFile(t *testing.T) {
	var got [2][]byte
	var rec [2]*fakeRec
	for i := range got {
		root := t.TempDir()
		writeOld(t, root, "q4", "---\ncolor: 45\n---\n\n# q4\n\n- "+idA+"\n")
		writeOld(t, root, "q5", "---\ncolor: 46\nafter: [q4, ghost]\n---\n")
		rec[i] = &fakeRec{}
		moved, err := Migrate(root, rec[i])
		if err != nil {
			t.Fatal(err)
		}
		if len(moved) != 2 || moved[0] != "q4.md -> q4-"+LegacyID("q4")+".md" {
			t.Fatalf("moved %v", moved)
		}
		if names := dirNames(t, root); names != "q4-"+LegacyID("q4")+".md,q5-"+LegacyID("q5")+".md" {
			t.Fatalf("directory %s", names)
		}
		got[i], _ = os.ReadFile(filepath.Join(Dir(root), "q5-"+LegacyID("q5")+".md"))
		if again, err := Migrate(root, rec[i]); err != nil || len(again) != 0 {
			t.Errorf("second run moved %v, %v", again, err)
		}
	}
	if string(got[0]) != string(got[1]) {
		t.Errorf("two clones wrote different bytes:\n%s\n%s", got[0], got[1])
	}
	if want := "after: [" + LegacyID("q4") + ", " + LegacyID("ghost") + "]"; !strings.Contains(string(got[0]), want) {
		t.Errorf("after: not migrated, want %s in\n%s", want, got[0])
	}
	if strings.Join(rec[0].names, ",") != LegacyID("q4")+","+LegacyID("q5") {
		t.Errorf("recorded %v", rec[0].names)
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

// The mark a rename of the name-keyed era left is removed, not migrated; a
// target already there is the newer file; an ID that already has a ref is not
// queued again.
func TestMigrateRemovesMarksKeepsNewerTargetsAndSkipsExistingRefs(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "old", "---\ncolor: 45\nstatus: filed\nrenamed-to: new\n---\n")
	writeOld(t, root, "both", "---\ncolor: 45\n---\n")
	writeFile(t, root, "both-"+LegacyID("both")+".md", "---\ncolor: 99\n---\n")
	writeOld(t, root, "ref", "---\ncolor: 47\n---\n")
	writeFile(t, root, "kept-K3F9QZ.md", "---\ncolor: 48\n---\n")
	rec := &fakeRec{have: map[string]bool{LegacyID("ref"): true}}
	if _, err := Migrate(root, rec); err != nil {
		t.Fatal(err)
	}
	if got, want := dirNames(t, root), "both-"+LegacyID("both")+".md,kept-K3F9QZ.md,ref-"+LegacyID("ref")+".md"; got != want {
		t.Errorf("directory %s, want %s", got, want)
	}
	if b, _ := os.ReadFile(filepath.Join(Dir(root), "both-"+LegacyID("both")+".md")); !strings.Contains(string(b), "color: 99") {
		t.Error("the newer target was overwritten")
	}
	if len(rec.names) != 0 {
		t.Errorf("recorded %v: the target existed or the ref did", rec.names)
	}
}

func TestMigrateReportsARecordFailureAfterTheMove(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "q4", "---\ncolor: 45\n---\n")
	_, err := Migrate(root, &fakeRec{err: errors.New("no remote")})
	var re *RecordError
	if !errors.As(err, &re) {
		t.Fatalf("err = %v", err)
	}
	if dirNames(t, root) != "q4-"+LegacyID("q4")+".md" {
		t.Error("the file did not move")
	}
}

// A ref that could not be queued for one moved file does not leave the files
// after it unmoved: every old file is moved, and the failure is reported once.
func TestMigrateMovesEveryFileDespiteARecordFailure(t *testing.T) {
	root := t.TempDir()
	writeOld(t, root, "q4", "---\ncolor: 45\n---\n")
	writeOld(t, root, "q5", "---\ncolor: 46\n---\n")
	moved, err := Migrate(root, &fakeRec{err: errors.New("no remote")})
	var re *RecordError
	if !errors.As(err, &re) {
		t.Fatalf("err = %v", err)
	}
	if len(moved) != 2 || dirNames(t, root) != "q4-"+LegacyID("q4")+".md,q5-"+LegacyID("q5")+".md" {
		t.Errorf("moved %v, directory %s", moved, dirNames(t, root))
	}
}
