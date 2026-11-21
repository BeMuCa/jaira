package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// A ".jaira/" a person wrote into info/exclude is not jaira's to remove: it
// covers every board in the repository. share leaves it, says that git still
// ignores the board, and --undo tells the person that the entry it writes is
// this clone's alone.
func TestShareLeavesAPersonsUnanchoredExcludeLineAndSaysSo(t *testing.T) {
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	common := strings.TrimSpace(gitOut(t, dir, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	exclude := filepath.Join(common, "info", "exclude")
	if err := os.WriteFile(exclude, []byte(readExclude(t, dir)+".jaira/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "share")
	if err != nil {
		t.Fatalf("jaira share: %v\n%s", err, out)
	}
	if excludesBoard(readExclude(t, dir)) {
		t.Errorf("share left jaira's own entry:\n%s", readExclude(t, dir))
	}
	if !strings.Contains(readExclude(t, dir), ".jaira/\n") {
		t.Errorf("share took out the person's line:\n%s", readExclude(t, dir))
	}
	if !strings.Contains(out, "still ignores .jaira/") {
		t.Errorf("share did not say the board is still ignored:\n%s", out)
	}
	out, err = runCLI(t, dir, "share", "--json")
	if err != nil || !strings.Contains(out, `"still_ignored": true`) {
		t.Errorf("share --json = %v\n%s", err, out)
	}

	out, err = runCLI(t, dir, "share", "--undo")
	if err != nil {
		t.Fatalf("jaira share --undo: %v\n%s", err, out)
	}
	if !strings.Contains(out, "jaira update") {
		t.Errorf("share --undo does not tell teammates to run jaira update:\n%s", out)
	}
}

// The rule that still ignores the board need not be one jaira can parse: an
// unanchored .jaira/ also covers a board in a subdirectory. share asks git.
func TestShareOfABoardInASubdirectorySeesAnUnanchoredExcludeLine(t *testing.T) {
	dir := repoWithCommit(t)
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, sub, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	common := strings.TrimSpace(gitOut(t, dir, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	exclude := filepath.Join(common, "info", "exclude")
	if err := os.WriteFile(exclude, []byte(readExclude(t, dir)+".jaira/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, sub, "share", "--json")
	if err != nil || !strings.Contains(out, `"still_ignored": true`) {
		t.Errorf("share --json in sub = %v\n%s", err, out)
	}
	if !strings.Contains(readExclude(t, dir), ".jaira/\n") {
		t.Errorf("share took out the person's line:\n%s", readExclude(t, dir))
	}
}

// What jaira update writes and says, over the whole class "what ignores the
// board": git decides (check-ignore), every rule counts, a git that cannot
// answer leaves everything alone, and a report names only a write that
// happened. Each source × git before/after 2.31 × .jaira a directory or a
// symlink.
func TestUpdateDecidesByWhatGitIgnores(t *testing.T) {
	type outcome string
	const (
		nothing        outcome = "nothing"
		moved          outcome = "moved to info/exclude"
		addedExclude   outcome = "added to info/exclude"
		addedGitignore outcome = "added to .gitignore"
		unknown        outcome = "unknown, nothing touched"
	)
	count := func(s, line string) int {
		n := 0
		for _, l := range strings.Split(s, "\n") {
			if strings.TrimSpace(l) == line {
				n++
			}
		}
		return n
	}
	appendTo := func(t *testing.T, path, line string) {
		t.Helper()
		b, _ := os.ReadFile(path)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(b, []byte(line+"\n")...), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// still: whether share, run after update, finds a rule jaira did not
	// write still ignoring the board. A symlinked .jaira is unknown instead.
	sources := []struct {
		name     string
		prepare  func(t *testing.T, dir, exclude, xdg string)
		new, old outcome
		still    bool
	}{
		{"jaira's line in .gitignore", func(t *testing.T, dir, _, _ string) { writeOldIgnoreLine(t, dir) }, moved, nothing, false},
		{"jaira's entry in info/exclude", func(t *testing.T, dir, _, _ string) {
			if _, err := board.AddIgnore(dir); err != nil {
				t.Fatal(err)
			}
		}, nothing, nothing, false},
		{"a person's .jaira/ in info/exclude", func(t *testing.T, _, exclude, _ string) { appendTo(t, exclude, ".jaira/") }, nothing, nothing, true},
		{"core.excludesFile", func(t *testing.T, dir, _, _ string) {
			f := filepath.Join(t.TempDir(), "ignore")
			appendTo(t, f, ".jaira/")
			gitRun(t, dir, "config", "core.excludesFile", f)
		}, nothing, nothing, true},
		{"global ignore file", func(t *testing.T, _, _, xdg string) { appendTo(t, filepath.Join(xdg, "git", "ignore"), ".jaira/") }, nothing, nothing, true},
		// The person's negation wins over jaira's entry; update does not
		// write a second one on top. Old git cannot see info/exclude, so it
		// adds the line to .gitignore — below the negation, which it then
		// overrides: the one row where old git does worse.
		{"!/.jaira/ over jaira's entry", func(t *testing.T, dir, _, _ string) {
			if _, err := board.AddIgnore(dir); err != nil {
				t.Fatal(err)
			}
			appendTo(t, filepath.Join(dir, ".gitignore"), "!/.jaira/")
		}, nothing, addedGitignore, false},
		{"nothing", func(*testing.T, string, string, string) {}, addedExclude, addedGitignore, false},
	}
	for _, src := range sources {
		for _, oldGit := range []bool{false, true} {
			for _, symlink := range []bool{false, true} {
				want := src.new
				if oldGit {
					want = src.old
				}
				if symlink {
					want = unknown
				}
				name := fmt.Sprintf("%s/old git %v/symlink %v", src.name, oldGit, symlink)
				t.Run(name, func(t *testing.T) {
					if symlink && runtime.GOOS == "windows" {
						t.Skip("symlinks need privileges on Windows")
					}
					xdg := t.TempDir()
					t.Setenv("XDG_CONFIG_HOME", xdg)
					t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(xdg, "gitconfig"))
					dir := repoWithCommit(t)
					if out, err := runCLI(t, dir, "init"); err != nil {
						t.Fatalf("jaira init: %v\n%s", err, out)
					}
					if _, err := board.RemoveIgnore(dir); err != nil {
						t.Fatal(err)
					}
					exclude := filepath.Join(dir, ".git", "info", "exclude")
					src.prepare(t, dir, exclude, xdg)
					if symlink {
						store := filepath.Join(t.TempDir(), "store")
						if err := os.Rename(filepath.Join(dir, ".jaira"), store); err != nil {
							t.Fatal(err)
						}
						if err := os.Symlink(store, filepath.Join(dir, ".jaira")); err != nil {
							t.Fatal(err)
						}
					}
					if oldGit {
						countGit(t, true) // git before 2.31
					}
					read := func(p string) string { b, _ := os.ReadFile(p); return string(b) }
					gitignore := filepath.Join(dir, ".gitignore")
					exBefore, giBefore := read(exclude), read(gitignore)

					dry, err := runUpdate(t, dir, "--dry-run")
					if err != nil {
						t.Fatalf("update --dry-run: %v\n%s", err, dry)
					}
					if read(exclude) != exBefore || read(gitignore) != giBefore {
						t.Errorf("--dry-run wrote")
					}
					out, err := runUpdate(t, dir)
					if err != nil {
						t.Fatalf("update: %v\n%s", err, out)
					}
					exAfter, giAfter := read(exclude), read(gitignore)
					exGained := count(exAfter, board.IgnoreLine) > count(exBefore, board.IgnoreLine)
					giGained := count(giAfter, board.IgnoreLine) > count(giBefore, board.IgnoreLine)

					says := map[outcome][2]string{
						moved:          {"Would move this private board's entry", "Moved this private board's entry"},
						addedExclude:   {"Would exclude this private board in .git/info/exclude", "Excluded this private board in .git/info/exclude"},
						addedGitignore: {"Would exclude this private board in .gitignore", "Excluded this private board in .gitignore"},
						unknown:        {"git could not say whether .jaira/ is ignored", "git could not say whether .jaira/ is ignored"},
					}
					if want == nothing {
						for _, o := range []string{dry, out} {
							if strings.Contains(o, "private board's entry") || strings.Contains(o, "this private board in") || strings.Contains(o, "could not say") {
								t.Errorf("reports an ignore change, want none:\n%s", o)
							}
						}
					} else {
						if !strings.Contains(dry, says[want][0]) {
							t.Errorf("--dry-run does not say %q:\n%s", says[want][0], dry)
						}
						if !strings.Contains(out, says[want][1]) {
							t.Errorf("update does not say %q:\n%s", says[want][1], out)
						}
					}
					switch want {
					case nothing, unknown:
						if exAfter != exBefore || giAfter != giBefore {
							t.Errorf("wrote although it should not:\n exclude %q -> %q\n .gitignore %q -> %q", exBefore, exAfter, giBefore, giAfter)
						}
					case moved:
						if !exGained || count(giAfter, board.IgnoreLine) != 0 {
							t.Errorf("not moved: exclude %q, .gitignore %q", exAfter, giAfter)
						}
					case addedExclude:
						if !exGained || giAfter != giBefore {
							t.Errorf("not added to exclude alone: exclude %q, .gitignore %q", exAfter, giAfter)
						}
					case addedGitignore:
						if !giGained || exAfter != exBefore {
							t.Errorf("not added to .gitignore alone: exclude %q, .gitignore %q", exAfter, giAfter)
						}
					}

					// share asks the same git: a rule it did not write is
					// reported, and a git that cannot answer is unknown.
					out, err = runCLI(t, dir, "share", "--json")
					if err != nil {
						t.Fatalf("share --json: %v\n%s", err, out)
					}
					wantStill := src.still
					// Old git gives share no path to info/exclude, so
					// jaira's own entry there stays — and share says so.
					if oldGit && src.name == "jaira's entry in info/exclude" {
						wantStill = true
					}
					var got struct {
						Still   *bool  `json:"still_ignored"`
						Unknown string `json:"ignore_unknown"`
					}
					if err := json.Unmarshal([]byte(out[strings.Index(out, "{"):]), &got); err != nil {
						t.Fatalf("share --json: %v\n%s", err, out)
					}
					switch {
					case symlink:
						if got.Still != nil || got.Unknown == "" {
							t.Errorf("share on a symlinked .jaira: still_ignored %s, ignore_unknown %q; want null and a reason", nullableBool(got.Still), got.Unknown)
						}
					case got.Still == nil || *got.Still != wantStill || got.Unknown != "":
						t.Errorf("share: still_ignored %s, ignore_unknown %q; want %v", nullableBool(got.Still), got.Unknown, wantStill)
					}
				})
			}
		}
	}
}

// nullableBool prints a decoded JSON bool as JSON does, null included.
func nullableBool(b *bool) string {
	if b == nil {
		return "null"
	}
	return fmt.Sprint(*b)
}

// A .gitignore jaira cannot read stops share before it touches the exclude:
// the board stays private, the .gitignore keeps every byte, and the error
// says so. A share run again once the file is readable completes.
func TestShareWithAnUnreadableGitignoreLeavesTheBoardPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod 000 does not make a file unreadable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root reads a file with mode 000")
	}
	dir := repoWithCommit(t)
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	gi := filepath.Join(dir, ".gitignore")
	want, err := os.ReadFile(gi)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(gi, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(gi, 0o644) })

	out, err := runCLI(t, dir, "share")
	if err == nil {
		t.Fatalf("jaira share with an unreadable .gitignore: err = nil\n%s", out)
	}
	if !strings.Contains(err.Error()+out, "stays private") {
		t.Errorf("share error does not say the board stays private: %v\n%s", err, out)
	}
	if !excludesBoard(readExclude(t, dir)) {
		t.Errorf("info/exclude lost the board's entry after a failed share:\n%s", readExclude(t, dir))
	}
	if err := os.Chmod(gi, 0o644); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(gi); string(got) != string(want) {
		t.Errorf(".gitignore changed by a failed share:\n got %q\nwant %q", got, want)
	}
	if out, err := runCLI(t, dir, "share"); err != nil {
		t.Fatalf("jaira share once .gitignore is readable: %v\n%s", err, out)
	}
	if excludesBoard(readExclude(t, dir)) {
		t.Errorf("info/exclude still excludes the board after share:\n%s", readExclude(t, dir))
	}

	// Shared already: the same failure must not claim the board is private.
	if err := os.Chmod(gi, 0); err != nil {
		t.Fatal(err)
	}
	out, err = runCLI(t, dir, "share")
	if err == nil {
		t.Fatalf("second jaira share with an unreadable .gitignore: err = nil\n%s", out)
	}
	if strings.Contains(err.Error()+out, "stays private") {
		t.Errorf("share on a shared board says it stays private: %v\n%s", err, out)
	}
}

// A board in " sub" that an older jaira made private carries "/sub/.jaira/",
// a line that excludes nothing of it. update writes the right entry and takes
// that one out, and says so; --dry-run only says it would.
func TestUpdateDropsTheEntryAnOlderJairaWroteWithoutTheLeadingSpace(t *testing.T) {
	dir := repoWithCommit(t)
	root := filepath.Join(dir, " sub")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, root, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	// What the older jaira left: its own entry, without the space.
	common := strings.TrimSpace(gitOut(t, dir, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	exclude := filepath.Join(common, "info", "exclude")
	f, err := os.OpenFile(exclude, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(f, "\n# jaira board — private to this machine. Run 'jaira share' to publish it.\n/sub/.jaira/\n")
	f.Close()
	stale := func() bool { return strings.Contains(readExclude(t, dir), "\n/sub/.jaira/\n") }

	dry, err := runUpdate(t, root, "--dry-run")
	if err != nil || !strings.Contains(dry, "Would remove /sub/.jaira/") || !stale() {
		t.Fatalf("--dry-run: %v, stale line there = %v, output:\n%s", err, stale(), dry)
	}
	out, err := runUpdate(t, root)
	if err != nil || !strings.Contains(out, "Removed /sub/.jaira/") {
		t.Fatalf("jaira update: %v\n%s", err, out)
	}
	if stale() || !strings.Contains(readExclude(t, dir), "\n/ sub/.jaira/\n") {
		t.Errorf("exclude after update:\n%s", readExclude(t, dir))
	}
	if again, err := runUpdate(t, root, "--json"); err != nil || !strings.Contains(again, `"ignore_stale_dropped": ""`) {
		t.Errorf("a second update: %v\n%s", err, again)
	}
	for _, line := range statusOf(t, dir) {
		if strings.Contains(line, ".jaira") {
			t.Errorf("git sees the board: %q", line)
		}
	}
}
