package project

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func mkBoard(t *testing.T, dir string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// repoWithWorktrees is a main checkout with a board in board/ (sub, "" for
// the top) and n linked worktrees, each with its own board copy, as a branch
// checked out beside the repository has.
func repoWithWorktrees(t *testing.T, base, sub string, n int) (main string, wts []string) {
	t.Helper()
	main = filepath.Join(base, "repo")
	if err := os.MkdirAll(main, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, main, "init", "-q")
	git(t, main, "commit", "-q", "--allow-empty", "-m", "init")
	mkBoard(t, filepath.Join(main, sub))
	for i := range n {
		wt := filepath.Join(base, ".worktrees", "repo-"+string(rune('a'+i)))
		git(t, main, "worktree", "add", "-q", "-b", "b"+string(rune('a'+i)), wt)
		wts = append(wts, mkBoard(t, filepath.Join(wt, sub)))
	}
	return canonical(filepath.Join(main, sub)), wts
}

// isolate points the registry and the temporary directory at two separate
// places under one test directory, so the boards the test makes are not in
// what Load counts as temporary, and returns that test directory.
func isolate(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(base, "home"))
	tmp := filepath.Join(base, "tmp")
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatal(err)
	}
	old := tempDir
	tempDir = func() string { return tmp }
	t.Cleanup(func() { tempDir = old })
	return base
}

func roots(ps []Project) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Root)
	}
	return out
}

// TestBoardIsMainCheckout covers what a root is listed under, one case per
// shape of .git a board can sit beside.
func TestBoardIsMainCheckout(t *testing.T) {
	base := isolate(t)

	main, wts := repoWithWorktrees(t, filepath.Join(base, "top"), "", 1)
	subMain, subWts := repoWithWorktrees(t, filepath.Join(base, "sub"), "tools/x", 1)

	// The main checkout has no board: the worktree's is the only one.
	bareMainDir := filepath.Join(base, "nomain", "repo")
	_, nomainWts := repoWithWorktrees(t, filepath.Join(base, "nomain"), "", 1)
	if err := os.RemoveAll(filepath.Join(bareMainDir, ".jaira")); err != nil {
		t.Fatal(err)
	}

	// A submodule's .git is a file too, but its git directory has no commondir.
	smod := mkBoard(t, filepath.Join(base, "smod", "child"))
	gd := filepath.Join(base, "smod", ".git", "modules", "child")
	if err := os.MkdirAll(gd, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(smod, ".git"), []byte("gitdir: ../.git/modules/child\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	plain := mkBoard(t, filepath.Join(base, "plain"))

	for _, c := range []struct {
		name, in, want string
	}{
		{"main checkout", main, main},
		{"linked worktree", wts[0], main},
		{"worktree, board in a subdirectory", subWts[0], subMain},
		{"worktree, main checkout without a board", nomainWts[0], canonical(nomainWts[0])},
		{"submodule", smod, canonical(smod)},
		{"no repository", plain, canonical(plain)},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := Board(c.in); got != c.want {
				t.Errorf("Board(%s) = %s, want %s", c.in, got, c.want)
			}
		})
	}
}

// TestRememberRecordsMainCheckout: a start in a worktree records the main
// checkout, and several worktrees open at once still leave one line.
func TestRememberRecordsMainCheckout(t *testing.T) {
	base := isolate(t)
	main, wts := repoWithWorktrees(t, base, "", 2)

	if got := Remember(wts[0]); got != main {
		t.Fatalf("Remember(worktree) = %s, want main checkout %s", got, main)
	}
	Remember(wts[1])
	Remember(main)
	if got := roots(Load()); len(got) != 1 || got[0] != main {
		t.Fatalf("Load() = %v, want only %s", got, main)
	}
}

// TestRememberSkipsTemp: a board in the temporary directory is not recorded,
// and one already recorded there is no longer listed.
func TestRememberSkipsTemp(t *testing.T) {
	base := isolate(t)
	scratch := mkBoard(t, filepath.Join(tempDir(), "scratchpad", "repo"))
	keep := mkBoard(t, filepath.Join(base, "keep"))

	if got := Remember(scratch); got != "" {
		t.Fatalf("Remember(temp board) = %q, want nothing recorded", got)
	}
	writeRegistry(t, scratch, keep)
	if got := roots(Load()); len(got) != 1 || got[0] != canonical(keep) {
		t.Fatalf("Load() = %v, want only %s", got, keep)
	}
}

// TestRememberTempRegistry: a registry that itself lives in the temporary
// directory - a test's, a sandbox's - keeps its temporary boards.
func TestRememberTempRegistry(t *testing.T) {
	isolate(t)
	t.Setenv("JAIRA_HOME", filepath.Join(tempDir(), "home"))
	b := mkBoard(t, filepath.Join(tempDir(), "repo"))
	if got := Remember(b); got != canonical(b) {
		t.Fatalf("Remember = %q, want %s", got, b)
	}
}

// TestLoadFoldsOldWorktreeEntries: entries written before Remember mapped
// worktrees fold into their main checkout, keep the latest open time, and a
// worktree whose folder is gone disappears.
func TestLoadFoldsOldWorktreeEntries(t *testing.T) {
	base := isolate(t)
	main, wts := repoWithWorktrees(t, base, "", 2)
	gone := filepath.Join(base, ".worktrees", "repo-gone")

	ps := []Project{
		{Root: main, Name: "repo", LastOpen: "2026-11-01T00:00:00Z"},
		{Root: wts[0], Name: "repo-a", LastOpen: "2026-11-20T00:00:00Z"},
		{Root: wts[1], Name: "repo-b", LastOpen: "2026-11-10T00:00:00Z"},
		{Root: gone, Name: "repo-gone", LastOpen: "2026-11-22T00:00:00Z"},
	}
	writeProjects(t, ps)

	got := Load()
	if len(got) != 1 || got[0].Root != main || got[0].Name != "repo" {
		t.Fatalf("Load() = %+v, want one entry for %s", got, main)
	}
	if got[0].LastOpen != "2026-11-20T00:00:00Z" {
		t.Errorf("LastOpen = %s, want the latest among the folded entries", got[0].LastOpen)
	}

	// Only a worktree entry, its main checkout never recorded: it shows as
	// the main checkout all the same.
	writeProjects(t, ps[1:2])
	if got := roots(Load()); len(got) != 1 || got[0] != main {
		t.Fatalf("Load() = %v, want %s", got, main)
	}
}

func writeRegistry(t *testing.T, rootList ...string) {
	t.Helper()
	var ps []Project
	for _, r := range rootList {
		ps = append(ps, Project{Root: r, Name: filepath.Base(r), LastOpen: "2026-11-01T00:00:00Z"})
	}
	writeProjects(t, ps)
}

func writeProjects(t *testing.T, ps []Project) {
	t.Helper()
	b, err := json.Marshal(ps)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(projectsPath()), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(projectsPath(), b, 0o644); err != nil {
		t.Fatal(err)
	}
}
