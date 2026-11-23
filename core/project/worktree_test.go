package project

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

	if got, err := Remember(wts[0]); err != nil || got != main {
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

	if got, err := Remember(scratch); err != nil || got != "" {
		t.Fatalf("Remember(temp board) = %q, want nothing recorded", got)
	}
	writeProjects(t, []Project{{Root: scratch, Name: "repo"}, {Root: keep, Name: "keep"}})
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
	if got, err := Remember(b); err != nil || got != canonical(b) {
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

// TestRememberParallel: starts that record at the same moment take turns -
// every board arrives and the file still parses. Written in place without a
// lock, two writers interleaved their bytes and Load read the result as an
// empty list.
func TestRememberParallel(t *testing.T) {
	base := isolate(t)
	const n = 12
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := range n {
		b := mkBoard(t, filepath.Join(base, "b"+strconv.Itoa(i)))
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := Remember(b); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(projectsPath())
	if err != nil {
		t.Fatal(err)
	}
	var ps []Project
	if err := json.Unmarshal(raw, &ps); err != nil {
		t.Fatalf("projects.json does not parse after parallel writes: %v", err)
	}
	if got := len(Load()); got != n {
		t.Fatalf("Load() has %d boards after %d parallel Remember, want all", got, n)
	}
}

// TestTrees: a main checkout's line stands for every linked worktree that
// holds a board, which is where per-tree state such as sessions is read.
func TestTrees(t *testing.T) {
	base := isolate(t)
	main, wts := repoWithWorktrees(t, base, "sub", 2)
	if err := os.RemoveAll(filepath.Join(wts[1], ".jaira")); err != nil {
		t.Fatal(err)
	}
	got := Trees(main)
	if len(got) != 2 || got[0] != main || got[1] != canonical(wts[0]) {
		t.Fatalf("Trees(%s) = %v, want [main, %s]", main, got, wts[0])
	}
	if got := Trees(wts[0]); len(got) != 1 {
		t.Errorf("Trees(worktree) = %v, want only itself", got)
	}
}

// TestRememberKeepsBrokenList: a projects.json that does not parse is copied
// aside before Remember replaces it, once, with one line on Stderr; a second
// broken file does not overwrite the first copy.
func TestRememberKeepsBrokenList(t *testing.T) {
	base := isolate(t)
	var out bytes.Buffer
	old := Stderr
	Stderr = &out
	t.Cleanup(func() { Stderr = old })

	path := projectsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	broken := []byte(`[{"root": "/a", "name": "a",}]`)
	if err := os.WriteFile(path, broken, 0o644); err != nil {
		t.Fatal(err)
	}
	b := mkBoard(t, filepath.Join(base, "b"))
	if _, err := Remember(b); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path + ".broken")
	if err != nil || !bytes.Equal(got, broken) {
		t.Fatalf("projects.json.broken = %q, %v; want the broken file", got, err)
	}
	if !strings.Contains(out.String(), path+".broken") || strings.Count(out.String(), "\n") != 1 {
		t.Fatalf("stderr = %q; want one line naming the copy", out.String())
	}
	if ps := Load(); len(ps) != 1 || ps[0].Root != Board(b) {
		t.Fatalf("Load() = %+v; want a new list with the board just opened", ps)
	}

	out.Reset()
	if _, err := Remember(b); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("a readable list printed %q", out.String())
	}

	second := []byte("{not json")
	if err := os.WriteFile(path, second, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Remember(b); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(path + ".broken"); !bytes.Equal(got, broken) {
		t.Fatalf("the first copy was overwritten: %q", got)
	}
	matches, _ := filepath.Glob(path + ".broken-*")
	if len(matches) != 1 {
		t.Fatalf("second copy: %v; want one timestamped file", matches)
	}
	if got, _ := os.ReadFile(matches[0]); !bytes.Equal(got, second) {
		t.Fatalf("second copy = %q", got)
	}
}

// TestRememberMissingOrEmptyListIsQuiet: no file and an empty file are the
// first start, not a broken list - no copy, no message.
func TestRememberMissingOrEmptyListIsQuiet(t *testing.T) {
	for _, content := range []*string{nil, ptr(""), ptr(" \n")} {
		base := isolate(t)
		var out bytes.Buffer
		old := Stderr
		Stderr = &out
		path := projectsPath()
		if content != nil {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(*content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		b := mkBoard(t, filepath.Join(base, "b"))
		_, err := Remember(b)
		Stderr = old
		if err != nil {
			t.Fatal(err)
		}
		if out.Len() != 0 {
			t.Errorf("content %v: stderr = %q", content, out.String())
		}
		if m, _ := filepath.Glob(path + ".broken*"); len(m) != 0 {
			t.Errorf("content %v: copy made: %v", content, m)
		}
		if len(Load()) != 1 {
			t.Errorf("content %v: board not recorded", content)
		}
	}
}

// TestRememberLeavesUnreadableList: a projects.json that is there but cannot
// be read is not replaced - Remember returns an error and writes nothing. A
// directory stands in for any read error (permissions do not stop root).
func TestRememberLeavesUnreadableList(t *testing.T) {
	base := isolate(t)
	path := projectsPath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	b := mkBoard(t, filepath.Join(base, "b"))
	if _, err := Remember(b); err == nil {
		t.Fatal("Remember over an unreadable projects.json returned no error")
	}
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		t.Fatalf("projects.json was replaced: %v", err)
	}
}

// TestRememberLeavesPermissionDeniedList: the case the directory above cannot
// show, because a rename over a directory fails anyway - a file with no read
// permission, which rename would replace.
func TestRememberLeavesPermissionDeniedList(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permission bits do not stop this reader")
	}
	base := isolate(t)
	path := projectsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`[]`), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(path, 0o644) })
	b := mkBoard(t, filepath.Join(base, "b"))
	if _, err := Remember(b); err == nil {
		t.Fatal("Remember over an unreadable projects.json returned no error")
	}
	os.Chmod(path, 0o644)
	if got, _ := os.ReadFile(path); string(got) != "[]" {
		t.Fatalf("projects.json was replaced: %q", got)
	}
}

func ptr(s string) *string { return &s }
