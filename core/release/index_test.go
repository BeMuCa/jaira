package release

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"
)

// fixtureRepo makes a git repository whose core/release holds a NOTES.md with
// one fixed release and an empty notes/ directory, and returns the repository
// and that package directory.
func fixtureRepo(t *testing.T) (repo, pkg string) {
	t.Helper()
	repo = t.TempDir()
	pkg = filepath.Join(repo, "core", "release")
	if err := os.MkdirAll(filepath.Join(pkg, notesDir), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, "NOTES.md"), "## 0.0.1\n- a fixed change\n")
	run(t, repo, "init", "-q", "-b", "main")
	run(t, repo, "add", ".")
	run(t, repo, "commit", "-q", "-m", "start")
	return repo, pkg
}

func run(t *testing.T, dir string, args ...string) error {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@example.test", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil && args[0] != "merge" {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return err
}

func writeFile(t *testing.T, p, s string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// addNote commits one note file on the current branch.
func addNote(t *testing.T, repo, pkg, name, line string) {
	t.Helper()
	// git keeps no empty directory, so a branch from before the first note
	// checks out without notes/.
	if err := os.MkdirAll(filepath.Join(pkg, notesDir), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, notesDir, name), "- "+line+"\n")
	run(t, repo, "add", ".")
	run(t, repo, "commit", "-q", "-m", "note "+name)
}

// Two branches that each add a note merge without a conflict — the reason
// notes are one file per change. The same two branches appending a line to one
// shared file, as NOTES.md used to be written, do conflict; that half proves
// the fixture can see a conflict at all.
func TestTwoBranchesAddingNotesMergeWithoutConflict(t *testing.T) {
	repo, pkg := fixtureRepo(t)
	run(t, repo, "branch", "a")
	run(t, repo, "branch", "b")

	run(t, repo, "checkout", "-q", "a")
	addNote(t, repo, pkg, "AAAAAA.md", "change on a")
	run(t, repo, "checkout", "-q", "b")
	addNote(t, repo, pkg, "BBBBBB.md", "change on b")

	run(t, repo, "checkout", "-q", "main")
	if err := run(t, repo, "merge", "-q", "--no-edit", "a"); err != nil {
		t.Fatalf("merge a: %v", err)
	}
	if err := run(t, repo, "merge", "-q", "--no-edit", "b"); err != nil {
		t.Fatalf("merge b after a: %v — two note files conflicted", err)
	}

	run(t, repo, "checkout", "-q", "-b", "c", "main~2")
	writeFile(t, filepath.Join(pkg, "NOTES.md"), "## 0.0.1\n- a fixed change\n- appended on c\n")
	run(t, repo, "commit", "-q", "-am", "append on c")
	run(t, repo, "checkout", "-q", "-b", "d", "main~2")
	writeFile(t, filepath.Join(pkg, "NOTES.md"), "## 0.0.1\n- a fixed change\n- appended on d\n")
	run(t, repo, "commit", "-q", "-am", "append on d")
	if err := run(t, repo, "merge", "-q", "--no-edit", "c"); err == nil {
		t.Fatal("two appends to one shared file merged cleanly; the fixture cannot tell a conflict apart")
	}
}

// A tagged build files each note under the oldest tag that holds it, without
// anybody having renamed anything before the tag; a note no tag holds stays
// Unreleased. Rebuilding the older tag does not see the newer tag's version.
func TestWriteIndexFilesNotesUnderTheTagThatFirstHeldThem(t *testing.T) {
	repo, pkg := fixtureRepo(t)
	addNote(t, repo, pkg, "AAAAAA.md", "first release change")
	run(t, repo, "tag", "v0.1.0")
	addNote(t, repo, pkg, "BBBBBB.md", "second release change")
	run(t, repo, "tag", "-a", "-m", "annotated", "v0.2.0")
	addNote(t, repo, pkg, "CCCCCC.md", "not released yet")

	if err := WriteIndex(pkg); err != nil {
		t.Fatalf("WriteIndex: %v", err)
	}
	all := composeNotes("## 0.0.1\n- a fixed change\n", os.DirFS(pkg))
	want := []Entry{
		{Version: "Unreleased", Changes: []string{"not released yet"}},
		{Version: "0.2.0", Changes: []string{"second release change"}},
		{Version: "0.1.0", Changes: []string{"first release change"}},
		{Version: "0.0.1", Changes: []string{"a fixed change"}},
	}
	if !reflect.DeepEqual(all, want) {
		t.Fatalf("composeNotes after WriteIndex = %#v, want %#v", all, want)
	}
	if got := sinceEntries(all, "0.1.0"); !reflect.DeepEqual(got, want[:2]) {
		t.Fatalf("sinceEntries(0.1.0) = %#v, want %#v", got, want[:2])
	}

	run(t, repo, "checkout", "-q", "v0.1.0")
	if err := WriteIndex(pkg); err != nil {
		t.Fatalf("WriteIndex at v0.1.0: %v", err)
	}
	got := composeNotes("", os.DirFS(pkg))
	if want := []Entry{{Version: "0.1.0", Changes: []string{"first release change"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("composeNotes at v0.1.0 = %#v, want %#v", got, want)
	}
}

// A build that never ran go generate — go build, go install — has no index
// and shows every note file as Unreleased, in file-name order, above the fixed
// history.
func TestComposeNotesWithoutIndexShowsEveryNoteUnreleased(t *testing.T) {
	fsys := fstest.MapFS{
		"notes/BBBBBB.md": {Data: []byte("- second\n")},
		"notes/AAAAAA.md": {Data: []byte("- first\nprose is ignored\n- first, again\n")},
	}
	got := composeNotes("## 0.0.1\n- a fixed change\n", fsys)
	want := []Entry{
		{Version: "Unreleased", Changes: []string{"first", "first, again", "second"}},
		{Version: "0.0.1", Changes: []string{"a fixed change"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("composeNotes without index = %#v, want %#v", got, want)
	}
}

// Moving the Unreleased lines out of NOTES.md into note files must not change
// what jaira update shows: every released section exactly, and the same
// Unreleased changes. Their order inside Unreleased now follows the file
// names, so that one is compared as a set. The index is left out because a
// checkout that ran go generate after a later tag rightly files these lines
// under that release.
func TestNotesUnchangedByMovingUnreleasedIntoFiles(t *testing.T) {
	before, err := os.ReadFile(filepath.Join("testdata", "NOTES-before-AR97K9.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := parseNotes(string(before))

	noIndex := fstest.MapFS{}
	files, _ := fs.Glob(notesFS, path.Join(notesDir, "*.md"))
	for _, f := range files {
		b, _ := fs.ReadFile(notesFS, f)
		noIndex[f] = &fstest.MapFile{Data: b}
	}
	got := composeNotes(notesMD, noIndex)

	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}
	if got[0].Version != "Unreleased" || !sameSet(got[0].Changes, want[0].Changes) {
		t.Fatalf("Unreleased = %#v, want the same changes as %#v", got[0], want[0])
	}
	if !reflect.DeepEqual(got[1:], want[1:]) {
		t.Fatal("a released section changed")
	}
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[string]int{}
	for _, s := range a {
		seen[s]++
	}
	for _, s := range b {
		if seen[s]--; seen[s] < 0 {
			return false
		}
	}
	return true
}
