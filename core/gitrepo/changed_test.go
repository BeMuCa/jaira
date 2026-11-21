package gitrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// A hunt still covers the code when nothing the branch itself did since
// touched what it read: other files, a merge from master and its conflict
// resolution do not count; a commit or an uncommitted edit on a hunted file
// does, and so does a commit git does not know.
func TestChangedSince(t *testing.T) {
	at := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name    string
		then    func(t *testing.T, dir string)
		paths   []string
		changed bool
	}{
		{"nothing since", func(*testing.T, string) {}, []string{"a.go"}, false},
		{"another file", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "b", at.Add(time.Hour))
		}, []string{"a.go"}, false},
		{"a hunted file committed", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "a", at.Add(time.Hour))
		}, []string{"a.go"}, true},
		{"a hunted file not committed", func(t *testing.T, dir string) {
			write(t, dir, "a.go", "a3")
		}, []string{"a.go"}, true},
		{"a path with a space", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"my dir/c.go": "c2"}, "c", at.Add(time.Hour))
		}, []string{"my dir/c.go"}, true},
		{"no paths reads all code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "b", at.Add(time.Hour))
		}, nil, true},
		{"a ticket file only, no paths", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/tickets/x.md": "t2"}, "t", at.Add(time.Hour))
		}, nil, false},
		{"master merged in, touching the hunted file", func(t *testing.T, dir string) {
			git(t, dir, "checkout", "-q", "-b", "side", "HEAD~0")
			fixtureCommit(t, dir, map[string]string{"a.go": "from master"}, "m", at.Add(time.Hour))
			git(t, dir, "checkout", "-q", "-")
			git(t, dir, "merge", "-q", "--no-ff", "-m", "merge", "side")
		}, []string{"a.go"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, repo := fixtureRepo(t)
			hunt := fixtureCommit(t, dir, map[string]string{"a.go": "a", "b.go": "b", "my dir/c.go": "c", ".jaira/tickets/x.md": "t"}, "start", at)
			tc.then(t, dir)
			got, err := repo.ChangedSince(hunt[:8], tc.paths)
			if err != nil || got != tc.changed {
				t.Errorf("ChangedSince = %v, %v; want %v", got, err, tc.changed)
			}
		})
	}
	t.Run("an unknown commit", func(t *testing.T) {
		dir, repo := fixtureRepo(t)
		fixtureCommit(t, dir, map[string]string{"a.go": "a"}, "start", at)
		if _, err := repo.ChangedSince("deadbeef", []string{"a.go"}); err == nil {
			t.Error("an unknown commit read as unchanged")
		}
	})
}

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
