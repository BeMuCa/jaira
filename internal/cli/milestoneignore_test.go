package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
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
	if out, err := runCLI(t, dir, "share", "--undo"); err != nil {
		t.Fatalf("share --undo: %v\n%s", err, out)
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
