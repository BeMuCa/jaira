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
