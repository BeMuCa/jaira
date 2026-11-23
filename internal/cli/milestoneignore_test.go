package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
)

// Milestone files travel on their refs and never in a commit: a shared board
// keeps /.jaira/milestones/ in its .gitignore, and a repository that
// committed them is told how to stop.

// dropMilestonesLine takes the milestones entry out of dir's .gitignore, as
// a board shared by an older build has it, and commits that.
func dropMilestonesLine(t *testing.T, dir string) {
	t.Helper()
	identify(t, dir)
	if _, err := board.RemoveMilestonesIgnore(dir); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(readIgnore(t, dir), board.MilestonesIgnoreLine) {
		t.Fatalf(".gitignore still has the line:\n%s", readIgnore(t, dir))
	}
	gitRun(t, dir, "commit", "--quiet", "-am", "an older share")
}

// identify gives a clone the committer sharedClone gives its origin only.
func identify(t *testing.T, dir string) {
	t.Helper()
	gitRun(t, dir, "config", "user.name", "grace")
	gitRun(t, dir, "config", "user.email", "grace@example.test")
}

func TestShareKeepsMilestoneFilesOutOfGit(t *testing.T) {
	dir := privateBoard(t)
	if out, err := runCLI(t, dir, "share"); err != nil {
		t.Fatalf("share: %v\n%s", err, out)
	}
	if !strings.Contains(readIgnore(t, dir), board.MilestonesIgnoreLine) {
		t.Fatalf("share left the milestones line out of .gitignore:\n%s", readIgnore(t, dir))
	}
	if out, err := runCLI(t, dir, "milestone", "create", "q4"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out := gitOut(t, dir, "status", "--porcelain", "--untracked-files=all", "--", ".jaira/milestones"); out != "" {
		t.Fatalf("git sees the milestone file:\n%s", out)
	}
	out, err := runCLI(t, dir, "share", "--undo", "--json")
	if err != nil {
		t.Fatalf("share --undo: %v\n%s", err, out)
	}
	var undo map[string]any
	if err := json.Unmarshal([]byte(out), &undo); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if undo["milestones_ignore_removed"] != true {
		t.Errorf("share --undo --json does not report the removed milestones line: %s", out)
	}
	if strings.Contains(readIgnore(t, dir), board.MilestonesIgnoreLine) {
		t.Fatalf("share --undo left the milestones line behind:\n%s", readIgnore(t, dir))
	}
}

func TestUpdateAndInitAddTheMilestonesLineToASharedBoard(t *testing.T) {
	clone := sharedClone(t)
	dropMilestonesLine(t, clone)

	out, err := runCLI(t, clone, "update", "--dry-run")
	if err != nil {
		t.Fatalf("update --dry-run: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Would add "+board.MilestonesIgnoreLine) || strings.Contains(readIgnore(t, clone), board.MilestonesIgnoreLine) {
		t.Fatalf("dry run: want a 'Would add' line and no write:\n%s\n.gitignore:\n%s", out, readIgnore(t, clone))
	}
	if out, err = runCLI(t, clone, "update"); err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	if !strings.Contains(readIgnore(t, clone), board.MilestonesIgnoreLine) {
		t.Fatalf("update left the line out:\n%s\n%s", out, readIgnore(t, clone))
	}

	gitRun(t, clone, "checkout", "--quiet", "--", ".gitignore")
	if out, err = runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if !strings.Contains(readIgnore(t, clone), board.MilestonesIgnoreLine) {
		t.Fatalf("init on a shared board left the line out:\n%s\n%s", out, readIgnore(t, clone))
	}
}

// A private board needs no line: its whole .jaira/ is excluded, and init and
// update must not touch the tracked .gitignore of a repository that never
// asked for the board.
func TestAPrivateBoardGetsNoMilestonesLine(t *testing.T) {
	dir := privateBoard(t)
	if out, err := runCLI(t, dir, "update"); err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); !os.IsNotExist(err) {
		t.Fatalf("update wrote a .gitignore on a private board:\n%s", readIgnore(t, dir))
	}
}

func TestCommittedMilestoneFilesAreWarnedAbout(t *testing.T) {
	clone := sharedClone(t)
	identify(t, clone)
	if out, err := runCLI(t, clone, "milestone", "create", "q4"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	gitRun(t, clone, "add", "-f", ".jaira/milestones")
	gitRun(t, clone, "commit", "--quiet", "-m", "milestones committed by hand")

	for _, args := range [][]string{{"update"}, {"validate"}} {
		out, err := runCLI(t, clone, args...)
		if err != nil {
			t.Fatalf("%s: %v\n%s", args[0], err, out)
		}
		if !strings.Contains(out, "git rm -r --cached .jaira/milestones") {
			t.Errorf("%s does not say how to take the milestone files out of git:\n%s", args[0], out)
		}
	}
	gitRun(t, clone, "rm", "-r", "--quiet", "--cached", ".jaira/milestones")
	if out, _ := runCLI(t, clone, "validate"); strings.Contains(out, "git rm -r --cached") {
		t.Errorf("validate still warns once the files are out of git:\n%s", out)
	}
}

// KWQCYZ through the two milestone writes that take the lock themselves:
// logbook reads its gate off the file, so a copy older than the ref used to
// file a milestone whose newer ticket was still open, and write that old list
// over the ref.
func TestLogbookReadsAStaleFileOnlyAfterItCaughtUp(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "q4"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(ada, "q4")
	if err != nil {
		t.Fatal(err)
	}
	stale, err := os.ReadFile(ms.Path(ada))
	if err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, ada, "create", "open work", "--goal", "g", "--context", "c", "--dod", "d", "--json")
	if err != nil {
		t.Fatalf("create ticket: %v\n%s", err, out)
	}
	var tk struct{ ID string }
	if err := json.Unmarshal([]byte(out), &tk); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "milestone", "add", "q4", tk.ID); err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}
	if err := os.WriteFile(ms.Path(ada), stale, 0o644); err != nil {
		t.Fatal(err)
	}

	out, err = runAndSend(t, ada, "logbook", "q4")
	if err == nil || !strings.Contains(out+err.Error(), "still holds work") {
		t.Fatalf("logbook filed a milestone whose ref holds an open ticket: err %v\n%s", err, out)
	}
	file, content, _, err := (&gitref.Repo{Dir: ada, Remote: "origin"}).ReadMilestone(ms.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got := milestone.FromBytes(file, content); got.Filed() || !got.Has(tk.ID) {
		t.Fatalf("the ref after the refused logbook: filed %v, members %v", got.Filed(), got.Members())
	}
}

// blockOutboxEntry makes the outbox unable to queue milestone id: a folder
// lies where its entry file goes. Not a read-only folder: Windows ignores the
// read-only bit on folders, and root ignores it everywhere.
func blockOutboxEntry(t *testing.T, dir, id string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(stateDirOf(t, dir), "outbox", "milestone-ids", id+".json"), 0o755); err != nil {
		t.Fatal(err)
	}
}

// A filing whose mark cannot be queued for the ref is no filing: every other
// clone would keep the milestone and a fetch here would write it back. It
// fails, and the milestone stays on the board unmarked.
func TestLogbookFailsWhenTheMarkCannotBeQueued(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "q4"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(ada, "q4")
	if err != nil {
		t.Fatal(err)
	}
	blockOutboxEntry(t, ada, ms.ID)

	out, err := runAndSend(t, ada, "logbook", "q4")
	if code := exitCode(err); code != ExitError {
		t.Fatalf("logbook with an outbox it cannot write: exit %d, want %d\n%s", code, ExitError, out)
	}
	got, err := milestone.Load(ada, ms.ID)
	if err != nil || got.Filed() {
		t.Fatalf("after the failed logbook the milestone should be on the board unmarked: %v, filed %v", err, got != nil && got.Filed())
	}
}

// A restore whose unmarked file cannot be queued for the ref puts the file
// back into the logbook, marked, as logbook takes its mark back off: the
// milestone is never on no board and in no logbook.
func TestRestoreGoesBackToTheLogbookWhenTheRefCannotBeQueued(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "q4"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ms, err := milestone.Load(ada, "q4")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := runAndSend(t, ada, "logbook", "q4"); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	logged, _ := filepath.Glob(filepath.Join(ada, ".jaira", "logbook", "*", "milestones", ms.FileName()))
	if len(logged) != 1 {
		t.Fatalf("logbook copy: %v", logged)
	}
	blockOutboxEntry(t, ada, ms.ID)

	out, err := runAndSend(t, ada, "restore", ms.FileName())
	if code := exitCode(err); code != ExitError {
		t.Fatalf("restore with an outbox it cannot write: exit %d, want %d\n%s", code, ExitError, out)
	}
	b, err := os.ReadFile(logged[0])
	if err != nil {
		t.Fatalf("the file did not go back to the logbook: %v", err)
	}
	if !milestone.FromBytes(ms.FileName(), b).Filed() {
		t.Fatalf("the file back in the logbook lost its filed mark:\n%s", b)
	}
	if _, err := os.Stat(ms.Path(ada)); !os.IsNotExist(err) {
		t.Fatalf("a copy is still on the board: %v", err)
	}
}
