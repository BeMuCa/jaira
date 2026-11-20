package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
)

// A private board is excluded in <git-common-dir>/info/exclude, never in the
// tracked .gitignore: a .gitignore jaira changed shows up in every worktree
// as a change nobody made, and a repository that refuses pushes with stray
// changes refuses every push.

// repoWithCommit is a git repository with one commit carrying a .gitignore
// of the project's own, as the repositories a board is added to have.
func repoWithCommit(t *testing.T) string {
	t.Helper()
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	dir := t.TempDir()
	gitRun(t, dir, "init", "--quiet", "-b", "master", ".")
	gitRun(t, dir, "config", "user.name", "ada")
	gitRun(t, dir, "config", "user.email", "ada@example.test")
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("/bin/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", ".gitignore")
	gitRun(t, dir, "commit", "--quiet", "-m", "project")
	return dir
}

func readExclude(t *testing.T, dir string) string {
	t.Helper()
	common := strings.TrimSpace(gitOut(t, dir, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	b, err := os.ReadFile(filepath.Join(common, "info", "exclude"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return string(b)
}

func excludesBoard(exclude string) bool {
	for _, line := range strings.Split(exclude, "\n") {
		if strings.TrimSpace(line) == board.IgnoreLine {
			return true
		}
	}
	return false
}

// statusOf is git status --porcelain for dir, minus agent instruction files
// init may announce the board in: those are a change init means to make.
func statusOf(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	for _, line := range strings.Split(strings.TrimSpace(gitOut(t, dir, "status", "--porcelain", "--untracked-files=all")), "\n") {
		if line == "" || strings.HasSuffix(line, "CLAUDE.md") || strings.HasSuffix(line, "AGENTS.md") {
			continue
		}
		out = append(out, line)
	}
	return out
}

func TestANewPrivateBoardAndItsWorktreesLeaveGitStatusClean(t *testing.T) {
	dir := repoWithCommit(t)
	before := readIgnore(t, dir)

	out, err := runCLI(t, dir, "init")
	if err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	if !strings.Contains(out, filepath.Join(".git", "info", "exclude")) && !strings.Contains(out, ".git/info/exclude") {
		t.Errorf("init does not say where the board is excluded:\n%s", out)
	}
	if got := readIgnore(t, dir); got != before {
		t.Errorf("init changed .gitignore:\n before %q\n after  %q", before, got)
	}
	if !excludesBoard(readExclude(t, dir)) {
		t.Errorf("info/exclude = %q, want it to contain %s", readExclude(t, dir), board.IgnoreLine)
	}
	if !board.Ignored(dir) {
		t.Errorf("Ignored = false on a new private board")
	}
	if st := statusOf(t, dir); len(st) != 0 {
		t.Errorf("git status after init = %q, want clean", st)
	}

	// A worktree shares the common git dir, so the one entry covers it: a
	// board there is ignored and its .gitignore is the committed one.
	wt := filepath.Join(t.TempDir(), "wt")
	gitRun(t, dir, "worktree", "add", "--quiet", "-b", "feat/x", wt)
	if err := os.MkdirAll(filepath.Join(wt, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, ".jaira", "tickets", "a.md"), []byte("---\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, wt, "init"); err != nil {
		t.Fatalf("jaira init in the worktree: %v\n%s", err, out)
	}
	if st := statusOf(t, wt); len(st) != 0 {
		t.Errorf("git status in a new worktree = %q, want clean", st)
	}
	if n := strings.Count(readExclude(t, wt), board.IgnoreLine+"\n"); n != 1 {
		t.Errorf("info/exclude holds the entry %d times after init in a worktree, want 1:\n%s", n, readExclude(t, wt))
	}
	// The worktree's own git dir (.git/worktrees/<name>) gets nothing: git
	// reads info/exclude from the common dir only.
	own := strings.TrimSpace(gitOut(t, wt, "rev-parse", "--path-format=absolute", "--git-dir"))
	if _, err := os.Stat(filepath.Join(own, "info", "exclude")); !os.IsNotExist(err) {
		t.Errorf("init in a worktree wrote %s", filepath.Join(own, "info", "exclude"))
	}

	// share and share --undo from the worktree act on the same shared file.
	if out, err := runCLI(t, wt, "share"); err != nil {
		t.Fatalf("jaira share in the worktree: %v\n%s", err, out)
	}
	if excludesBoard(readExclude(t, dir)) {
		t.Errorf("share in a worktree left the entry in the common info/exclude")
	}
	if out, err := runCLI(t, wt, "share", "--undo"); err != nil {
		t.Fatalf("jaira share --undo in the worktree: %v\n%s", err, out)
	}
	if !excludesBoard(readExclude(t, dir)) {
		t.Errorf("share --undo in a worktree did not write the common info/exclude")
	}
}

func TestShareLiftsTheExcludeAndKeepsLanesPrivate(t *testing.T) {
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "share"); err != nil {
		t.Fatalf("jaira share: %v\n%s", err, out)
	}
	if excludesBoard(readExclude(t, dir)) {
		t.Errorf("info/exclude still excludes the board after share:\n%s", readExclude(t, dir))
	}
	if board.Ignored(dir) {
		t.Errorf("Ignored = true after share")
	}
	gi := readIgnore(t, dir)
	if !strings.Contains(gi, board.LanesIgnoreLine) {
		t.Errorf(".gitignore = %q, want %s kept private", gi, board.LanesIgnoreLine)
	}
	if !strings.Contains(gi, "/bin/") {
		t.Errorf(".gitignore lost the project's own line: %q", gi)
	}
	// What git would now take: the board, not its lanes.
	st := strings.Join(statusOf(t, dir), "\n")
	if !strings.Contains(st, ".jaira/") {
		t.Errorf("git status after share shows no .jaira file:\n%s", st)
	}
	if strings.Contains(st, ".jaira/lanes/") {
		t.Errorf("git status after share shows the lanes:\n%s", st)
	}

	// And back: --undo excludes it in info/exclude again, .gitignore loses
	// only the lanes line share added.
	if out, err := runCLI(t, dir, "share", "--undo"); err != nil {
		t.Fatalf("jaira share --undo: %v\n%s", err, out)
	}
	if !excludesBoard(readExclude(t, dir)) {
		t.Errorf("share --undo did not exclude the board in info/exclude:\n%s", readExclude(t, dir))
	}
	if got := readIgnore(t, dir); strings.Contains(got, ".jaira") {
		t.Errorf(".gitignore after share --undo = %q, want no .jaira line", got)
	}
}

func TestUpdateMovesTheOldGitignoreLineToExclude(t *testing.T) {
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	// The board an older jaira left: entry in .gitignore, none in exclude.
	if _, err := board.RemoveIgnore(dir); err != nil {
		t.Fatal(err)
	}
	writeOldIgnoreLine(t, dir)
	gitRun(t, dir, "commit", "--quiet", "-am", "an older init made it private")

	dry, err := runUpdate(t, dir, "--dry-run")
	if err != nil {
		t.Fatalf("jaira update --dry-run: %v\n%s", err, dry)
	}
	if !strings.Contains(dry, "Would move") || !board.InGitignore(dir) {
		t.Errorf("--dry-run: InGitignore = %v, output:\n%s", board.InGitignore(dir), dry)
	}

	out, err := runUpdate(t, dir)
	if err != nil {
		t.Fatalf("jaira update: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Moved this private board's entry") {
		t.Errorf("update does not say it moved the entry:\n%s", out)
	}
	if got := readIgnore(t, dir); got != "/bin/\n" {
		t.Errorf(".gitignore after update = %q, want the project's own %q back", got, "/bin/\n")
	}
	if !excludesBoard(readExclude(t, dir)) || !board.Ignored(dir) {
		t.Errorf("the board is no longer excluded after update:\n%s", readExclude(t, dir))
	}

	// A teammate pulls that commit: their .gitignore loses the line, their
	// private board is excluded nowhere. update excludes it.
	if _, err := board.RemoveIgnore(dir); err != nil {
		t.Fatal(err)
	}
	if out, err := runUpdate(t, dir); err != nil || !strings.Contains(out, "Excluded this private board") {
		t.Fatalf("jaira update on an unexcluded private board: %v\n%s", err, out)
	}
	if !excludesBoard(readExclude(t, dir)) {
		t.Errorf("update left an unexcluded private board unexcluded:\n%s", readExclude(t, dir))
	}
}

// A board in git — shared, or with .jaira/ files committed — is left as it
// is by update, whether its .gitignore excludes it or nothing does: a
// .gitignore change committed from there would un-ignore it for every
// teammate who pulls it, and excluding a shared board stops new tickets
// reaching anyone.
func TestUpdateLeavesABoardInGitAsItIs(t *testing.T) {
	cases := []struct {
		name string
		make func(t *testing.T) string
	}{
		{"shared clone", sharedClone},
		{"shared clone with a committed ignore line", func(t *testing.T) string {
			dir := sharedClone(t)
			writeOldIgnoreLine(t, dir)
			gitRun(t, dir, "config", "user.name", "ada")
			gitRun(t, dir, "config", "user.email", "ada@example.test")
			gitRun(t, dir, "commit", "--quiet", "-am", "an older init made it private")
			return dir
		}},
		{"committed tickets, no attributes, ignore line", func(t *testing.T) string {
			dir := repoWithCommit(t)
			if out, err := runCLI(t, dir, "init"); err != nil {
				t.Fatalf("jaira init: %v\n%s", err, out)
			}
			if _, err := board.RemoveIgnore(dir); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, ".jaira", "tickets", "x.md"), []byte("x\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			gitRun(t, dir, "add", ".jaira/tickets/x.md")
			writeOldIgnoreLine(t, dir)
			gitRun(t, dir, "commit", "--quiet", "-am", "tickets committed, then ignored")
			return dir
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := c.make(t)
			before, excludeBefore, ignored := readIgnore(t, dir), readExclude(t, dir), board.Ignored(dir)
			for _, args := range [][]string{{"--dry-run"}, nil} {
				out, err := runUpdate(t, dir, args...)
				if err != nil {
					t.Fatalf("jaira update %v: %v\n%s", args, err, out)
				}
				if strings.Contains(out, "Moved") || strings.Contains(out, "Would move") || strings.Contains(out, "xcluded this private board") || strings.Contains(out, "Would exclude") {
					t.Errorf("update %v reports touching a board in git:\n%s", args, out)
				}
			}
			if readIgnore(t, dir) != before || readExclude(t, dir) != excludeBefore || board.Ignored(dir) != ignored {
				t.Errorf("update changed a board in git:\n gitignore %q -> %q\n exclude %q -> %q", before, readIgnore(t, dir), excludeBefore, readExclude(t, dir))
			}
		})
	}
}
