package release

import (
	"os"
	"os/exec"
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

// Files directly in notes/ are Unreleased, every notes/<version>/ folder is
// one release named after it, newest version first — 0.10.0 above 0.9.0, which
// a sort by name would get wrong — and NOTES.md comes last. Inside a release
// the files go by name, and only "- " lines are changes.
func TestComposeNotesReadsVersionFolders(t *testing.T) {
	fsys := fstest.MapFS{
		"notes/BBBBBB.md":        {Data: []byte("- second\n")},
		"notes/AAAAAA.md":        {Data: []byte("- first\nprose is ignored\n- first, again\n")},
		"notes/0.9.0/CCCCCC.md":  {Data: []byte("- older release\n")},
		"notes/0.10.0/EEEEEE.md": {Data: []byte("- newer release, second file\n")},
		"notes/0.10.0/DDDDDD.md": {Data: []byte("- newer release, first file\n")},
	}
	got := composeNotes("## 0.0.1\n- a fixed change\n", fsys)
	want := []Entry{
		{Version: "Unreleased", Changes: []string{"first", "first, again", "second"}},
		{Version: "0.10.0", Changes: []string{"newer release, first file", "newer release, second file"}},
		{Version: "0.9.0", Changes: []string{"older release"}},
		{Version: "0.0.1", Changes: []string{"a fixed change"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("composeNotes = %#v, want %#v", got, want)
	}
	if got := sinceEntries(got, "0.9.0"); !reflect.DeepEqual(got, want[:2]) {
		t.Fatalf("sinceEntries(0.9.0) = %#v, want %#v", got, want[:2])
	}
}

// With no note files at all there is no Unreleased entry, only the history.
func TestComposeNotesWithoutNotesShowsOnlyHistory(t *testing.T) {
	got := composeNotes("## 0.0.1\n- a fixed change\n", fstest.MapFS{})
	if want := []Entry{{Version: "0.0.1", Changes: []string{"a fixed change"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("composeNotes without notes = %#v, want %#v", got, want)
	}
}

// Cut is what the release workflow runs: it moves every note directly in
// notes/ into notes/<version>/, and afterwards those changes read as that
// release while a note added later is Unreleased again.
func TestCutMovesEveryNoteIntoTheVersionFolder(t *testing.T) {
	pkg := t.TempDir()
	if err := os.MkdirAll(filepath.Join(pkg, notesDir, "0.1.0"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, notesDir, "0.1.0", "OLDOLD.md"), "- released before\n")
	writeFile(t, filepath.Join(pkg, notesDir, "AAAAAA.md"), "- change a\n")
	writeFile(t, filepath.Join(pkg, notesDir, "BBBBBB.md"), "- change b\n")

	version, err := Cut(pkg, "v0.2.0")
	if err != nil {
		t.Fatalf("Cut: %v", err)
	}
	if version != "0.2.0" {
		t.Fatalf("Cut returned %q, want the leading v dropped", version)
	}
	writeFile(t, filepath.Join(pkg, notesDir, "CCCCCC.md"), "- after the cut\n")

	got := composeNotes("", os.DirFS(pkg))
	want := []Entry{
		{Version: "Unreleased", Changes: []string{"after the cut"}},
		{Version: "0.2.0", Changes: []string{"change a", "change b"}},
		{Version: "0.1.0", Changes: []string{"released before"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("composeNotes after Cut = %#v, want %#v", got, want)
	}
}

// Cut refuses what would make a broken release, and moves nothing when it does.
func TestCutRefuses(t *testing.T) {
	for _, tc := range []struct{ name, version string }{
		{"not X.Y.Z", "0.2"},
		{"not a number", "latest"},
		{"released before", "0.1.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pkg := t.TempDir()
			if err := os.MkdirAll(filepath.Join(pkg, notesDir, "0.1.0"), 0o755); err != nil {
				t.Fatal(err)
			}
			writeFile(t, filepath.Join(pkg, notesDir, "AAAAAA.md"), "- change a\n")
			if _, err := Cut(pkg, tc.version); err == nil {
				t.Fatalf("Cut(%q) succeeded", tc.version)
			}
			if _, err := os.Stat(filepath.Join(pkg, notesDir, "AAAAAA.md")); err != nil {
				t.Fatalf("the note moved although Cut refused: %v", err)
			}
		})
	}
	t.Run("nothing to release", func(t *testing.T) {
		pkg := t.TempDir()
		if err := os.MkdirAll(filepath.Join(pkg, notesDir), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := Cut(pkg, "0.2.0"); err == nil {
			t.Fatal("Cut with no notes succeeded")
		}
		if _, err := os.Stat(filepath.Join(pkg, notesDir, "0.2.0")); err == nil {
			t.Fatal("Cut with no notes left an empty release folder")
		}
	})
}

// Moving the Unreleased lines out of NOTES.md into note files must not change
// what jaira update shows: every released section exactly, and the same
// Unreleased changes. Their order inside Unreleased now follows the file
// names, so that one is compared as a set. Once a release folder exists this
// test has done its job: those lines then rightly read as that release.
func TestNotesUnchangedByMovingUnreleasedIntoFiles(t *testing.T) {
	before, err := os.ReadFile(filepath.Join("testdata", "NOTES-before-AR97K9.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := parseNotes(string(before))

	got := Notes()

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
