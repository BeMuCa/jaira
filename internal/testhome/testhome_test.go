package testhome

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A git hook exports GIT_DIR, or TMPDIR lies inside a checkout: either way a
// test's 'git -C <temp dir>' reached the developer's repository and read or
// overwrote its refs/jaira/people. After isolateGit a temp dir is no repository.
func TestIsolateGitKeepsTempDirsOutOfTheDevelopersRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	outer := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", outer).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	tmp := filepath.Join(outer, "tmp")
	if err := os.Mkdir(tmp, 0o755); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ name, gitDir, tmpDir string }{
		{"GIT_DIR exported", filepath.Join(outer, ".git"), t.TempDir()},
		{"TMPDIR inside a repository", "", tmp},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("GIT_DIR", c.gitDir)
			if c.gitDir == "" {
				os.Unsetenv("GIT_DIR") // set but empty, git refuses to run at all
			}
			t.Setenv("TMPDIR", c.tmpDir)
			t.Setenv("TMP", c.tmpDir) // os.TempDir on Windows
			t.Setenv("GIT_CEILING_DIRECTORIES", "")
			isolateGit()
			dir, err := os.MkdirTemp("", "probe")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			out, err := exec.Command("git", "-C", dir, "rev-parse", "--git-dir").CombinedOutput()
			if err == nil {
				t.Errorf("git found a repository from %s: %s", dir, strings.TrimSpace(string(out)))
			}
		})
	}
}

// Isolation only reaches a package whose TestMain calls Run. A package whose
// tests run git and has none sent its git calls into the developer's checkout
// under an exported GIT_DIR — branches, tags and .git/config rewritten. So
// every such package is listed here until it has one.
func TestEveryPackageThatRunsGitInTestsCallsRun(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	usesGit := map[string]bool{}
	callsRun := map[string]bool{}
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != root && (strings.HasPrefix(name, ".") || name == "node_modules" || name == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		dir := filepath.Dir(path)
		usesGit[dir] = usesGit[dir] || strings.Contains(string(src), `"git"`)
		callsRun[dir] = callsRun[dir] || strings.Contains(string(src), "testhome.Run(")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	own, _ := filepath.Abs(".")
	for dir := range usesGit {
		if usesGit[dir] && !callsRun[dir] && dir != own {
			rel, _ := filepath.Rel(root, dir)
			t.Errorf("%s runs git in its tests but has no TestMain calling testhome.Run", rel)
		}
	}
}
