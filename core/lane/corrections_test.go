package lane

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// oldBoard builds a board the way a build from before 9ad7aa9 left one: every
// shipped lane as a file, a done.md that still carries logbook-on-entry: true,
// and no corrections marker — that file did not exist yet. It returns the root
// and the path of done.md.
//
// The board is set up by Load and its done.md then overwritten, rather than
// written lane by lane, so the fixture keeps matching whatever the rest of the
// set-up path does.
func oldBoard(t *testing.T, doneFile string) (root, donePath string) {
	t.Helper()
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	root = t.TempDir()
	if err := os.MkdirAll(ProjectLanesDir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}
	donePath = filepath.Join(ProjectLanesDir(root), "done.md")
	if err := os.WriteFile(donePath, []byte(doneFile), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(correctionsPath(root)); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return root, donePath
}

// TestCorrectionRemovesTheDoorwayFromAnOldBoard: the board of a build that
// shipped the doorway loses the line on the next load, is told so, and keeps
// every other byte of the file it had.
func TestCorrectionRemovesTheDoorwayFromAnOldBoard(t *testing.T) {
	root, donePath := oldBoard(t, doneDoorway)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	done, ok := set.Get("done")
	if !ok {
		t.Fatal("the done lane is gone")
	}
	if done.LogbookOnEntry {
		t.Error("the loaded lane still files on entry")
	}

	got, err := os.ReadFile(donePath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "logbook-on-entry") {
		t.Errorf("the line is still in the file:\n%s", got)
	}
	want, _ := dropFrontmatterLine([]byte(doneDoorway), "logbook-on-entry", "true")
	if string(got) != string(want) {
		t.Errorf("the correction changed more than the one line:\ngot:\n%s\nwant:\n%s", got, want)
	}
	if !containsWarning(set.Warnings, donePath) || !containsWarning(set.Warnings, "logbook-on-entry") {
		t.Errorf("the correction must be reported, naming the file; got: %v", set.Warnings)
	}
}

// TestCorrectionRunsOncePerBoard: the marker file is written, and a line put
// back by hand afterwards is left alone — the board is never asked again, which
// is what keeps "a board is its lane directory" true.
func TestCorrectionRunsOncePerBoard(t *testing.T) {
	root, donePath := oldBoard(t, doneDoorway)
	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}

	ids, err := readIDList(correctionsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "done-logbook-on-entry" {
		t.Fatalf("marker file = %v, want the correction recorded once", ids)
	}

	if err := os.WriteFile(donePath, []byte(doneDoorway), 0o644); err != nil {
		t.Fatal(err)
	}
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(donePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != doneDoorway {
		t.Errorf("a line written back by hand was removed again:\n%s", got)
	}
	if done, _ := set.Get("done"); done == nil || !done.LogbookOnEntry {
		t.Error("the board asked for the doorway back and did not get it")
	}
	if containsWarning(set.Warnings, "logbook-on-entry") {
		t.Errorf("a correction that already ran must say nothing; got: %v", set.Warnings)
	}
}

// TestCorrectionLeavesALaneSomebodyWroteAlone: the same lane id, the same
// defective field, but not a file jaira shipped. It is reported and not
// touched — the once-only marker stops a correction repeating, not a correction
// being wrong the first time.
func TestCorrectionLeavesALaneSomebodyWroteAlone(t *testing.T) {
	mine := strings.Replace(doneDoorway,
		"description: Accepted.", "description: Ours. We want the doorway.", 1)
	if mine == doneDoorway {
		t.Fatal("fixture did not change the description")
	}
	root, donePath := oldBoard(t, mine)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(donePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != mine {
		t.Errorf("a lane written by hand was edited:\ngot:\n%s\nwant:\n%s", got, mine)
	}
	if done, _ := set.Get("done"); done == nil || !done.LogbookOnEntry {
		t.Error("the hand-written lane lost its doorway")
	}
	if !containsWarning(set.Warnings, donePath) || !containsWarning(set.Warnings, "left") {
		t.Errorf("skipping the correction must be reported, naming the file; got: %v", set.Warnings)
	}
}

// TestCorrectionLeavesTodaysLanePlusTheLineAlone: today's shipped done.md with
// the line added back by hand is not something any build ever wrote, so it is
// somebody's decision and is left standing. This is why the shipped versions a
// correction recognises are only the ones that carried the defect.
func TestCorrectionLeavesTodaysLanePlusTheLineAlone(t *testing.T) {
	current, err := builtinFS.ReadFile("builtin/50-done.md")
	if err != nil {
		t.Fatal(err)
	}
	mine := strings.Replace(string(current), "terminal: true\n",
		"terminal: true\nlogbook-on-entry: true\n", 1)
	root, donePath := oldBoard(t, mine)

	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(donePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != mine {
		t.Errorf("a line added to today's lane by hand was removed:\n%s", got)
	}
}

// TestCorrectionSaysNothingOnABoardThatNeverHadTheDefect: a board set up by
// this build is marked as corrected without a word, so nobody is told about a
// repair that had nothing to repair.
func TestCorrectionSaysNothingOnABoardThatNeverHadTheDefect(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	root := t.TempDir()
	if err := os.MkdirAll(ProjectLanesDir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if containsWarning(set.Warnings, "logbook-on-entry") {
		t.Errorf("a fresh board was told about a correction; got: %v", set.Warnings)
	}
	ids, _ := readIDList(correctionsPath(root))
	if len(ids) != 1 {
		t.Errorf("marker file = %v, want the correction recorded so it never runs here", ids)
	}
}

func TestDropFrontmatterLine(t *testing.T) {
	const body = "---\nid: done\nlogbook-on-entry: true  # doorway\nterminal: true\n---\nprompt\n"
	for _, tc := range []struct {
		name, field, value, want string
		dropped                  bool
	}{
		{"drops the field", "logbook-on-entry", "", "---\nid: done\nterminal: true\n---\nprompt\n", true},
		{"drops it on a matching value past a comment", "logbook-on-entry", "true", "---\nid: done\nterminal: true\n---\nprompt\n", true},
		{"keeps it on another value", "logbook-on-entry", "false", body, false},
		{"keeps a field that is not there", "holds", "", body, false},
		{"never reaches past the frontmatter", "prompt", "", body, false},
		{"leaves a file with no frontmatter alone", "id", "", body, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := body
			if tc.name == "leaves a file with no frontmatter alone" {
				in, tc.want = "id: done\n", "id: done\n"
			}
			got, dropped := dropFrontmatterLine([]byte(in), tc.field, tc.value)
			if dropped != tc.dropped || string(got) != tc.want {
				t.Errorf("dropped=%v content=%q, want dropped=%v content=%q", dropped, got, tc.dropped, tc.want)
			}
		})
	}
}
