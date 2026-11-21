package gitrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// A hunt still covers the code when nothing the ticket itself did since
// touched what it read: other files, other people's commits and a merge of
// master do not count; a ticket commit or an uncommitted edit on a hunted
// file does, and so does a hunt this HEAD does not descend from, or a file
// list that names no file.
func TestChangedSince(t *testing.T) {
	const id = "01M7WWZWSDHK37DGGHJFDF0YD6" // handle DF0YD6
	at := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	later := at.Add(time.Hour)
	hunted := time.Now().Add(-time.Minute) // the untracked files below are written after it
	for _, tc := range []struct {
		name    string
		then    func(t *testing.T, dir string)
		paths   []string
		changed bool
	}{
		{"nothing since", func(*testing.T, string) {}, []string{"a.go"}, false},
		{"another file", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "fix(DF0YD6): b", later)
		}, []string{"a.go"}, false},
		{"a hunted file, by this ticket", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "fix(DF0YD6): a", later)
		}, []string{"a.go"}, true},
		{"a hunted file, by another ticket", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "fix(ZZZZZZ): a", later)
		}, []string{"a.go"}, false},
		{"a hunted file not committed", func(t *testing.T, dir string) {
			write(t, dir, "a.go", "a3")
		}, []string{"a.go"}, true},
		{"a path with a space", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"my dir/c.go": "c2"}, "fix(DF0YD6): c", later)
		}, []string{"my dir/c.go"}, true},
		{"a list split at a space names no file", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"my dir/c.go": "c2"}, "fix(DF0YD6): c", later)
		}, []string{"my", "dir/c.go"}, true},
		{"no paths reads all code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "fix(DF0YD6): b", later)
		}, nil, true},
		{"an uncommitted edit from before the hunt (conversational mode)", func(t *testing.T, dir string) {
			write(t, dir, "a.go", "a3")
			old := time.Now().Add(-time.Hour)
			if err := os.Chtimes(filepath.Join(dir, "a.go"), old, old); err != nil {
				t.Fatal(err)
			}
		}, []string{"a.go"}, false},
		{"a hunted file deleted, not committed", func(t *testing.T, dir string) {
			if err := os.Remove(filepath.Join(dir, "a.go")); err != nil {
				t.Fatal(err)
			}
		}, []string{"a.go"}, true},
		{"untracked board data under .jaira, no paths", func(t *testing.T, dir string) {
			if err := os.MkdirAll(filepath.Join(dir, ".jaira", "milestones"), 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, dir, ".jaira/milestones/q4-ABC123.md", "m")
		}, nil, false},
		{"untracked board data under .jaira, listed", func(t *testing.T, dir string) {
			if err := os.MkdirAll(filepath.Join(dir, ".jaira", "milestones"), 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, dir, ".jaira/milestones/q4-ABC123.md", "m")
		}, []string{"a.go", ".jaira/milestones/q4-ABC123.md"}, false},
		{"tracked board data changed after the hunt, no paths", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/tags": "ui red"}, "chore(DF0YD6): tag", later)
			write(t, dir, ".jaira/tags", "ui red\ngates blue")
		}, nil, false},
		{"tracked board data changed after the hunt, listed", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/milestones/q4-ABC123.md": "m"}, "chore(DF0YD6): milestone", later)
		}, []string{"a.go", ".jaira/milestones/q4-ABC123.md"}, false},
		{"a lane prompt is code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/lanes/testing.md": "p2"}, "fix(DF0YD6): prompt", later)
		}, nil, true},
		{"a named board's lane prompt is code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/boards/video/lanes/testing.md": "p2"}, "fix(DF0YD6): prompt", later)
		}, nil, true},
		{"a named board's other data is not", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/boards/video/notes.md": "n"}, "chore(DF0YD6): board", later)
		}, nil, false},
		{"an untracked file written after the hunt", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, nil, true},
		{"an untracked hunted file written after the hunt", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, []string{"a.go", "new.go"}, true},
		{"an untracked file written before the hunt", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
			old := time.Now().Add(-time.Hour)
			if err := os.Chtimes(filepath.Join(dir, "new.go"), old, old); err != nil {
				t.Fatal(err)
			}
		}, []string{"a.go", "new.go"}, false},
		{"an untracked file outside the hunted paths", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, []string{"a.go"}, false},
		{"a ticket file only, no paths", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/tickets/x.md": "t2"}, "DF0YD6 note", later)
		}, nil, false},
		{"master merged in, touching the hunted file", func(t *testing.T, dir string) {
			git(t, dir, "checkout", "-q", "-b", "side")
			fixtureCommit(t, dir, map[string]string{"a.go": "from master"}, "fix(ZZZZZZ): a", later)
			git(t, dir, "checkout", "-q", "-")
			git(t, dir, "merge", "-q", "--no-ff", "-m", "chore(DF0YD6): merge master", "side")
		}, []string{"a.go"}, false},
		{"a checkout that does not descend from the hunt", func(t *testing.T, dir string) {
			git(t, dir, "checkout", "-q", "--orphan", "other")
			fixtureCommit(t, dir, map[string]string{"z.go": "z"}, "other", later)
		}, []string{"a.go"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, repo := fixtureRepo(t)
			hunt := fixtureCommit(t, dir, map[string]string{"a.go": "a", "b.go": "b", "my dir/c.go": "c", ".jaira/tickets/x.md": "t"}, "feat(DF0YD6): start", at)
			tc.then(t, dir)
			got, err := repo.ChangedSince(hunt[:8], id, hunted, tc.paths)
			if err != nil || got != tc.changed {
				t.Errorf("ChangedSince = %v, %v; want %v", got, err, tc.changed)
			}
		})
	}
	t.Run("an unknown commit", func(t *testing.T) {
		dir, repo := fixtureRepo(t)
		fixtureCommit(t, dir, map[string]string{"a.go": "a"}, "start", at)
		if _, err := repo.ChangedSince("deadbeef", id, time.Now(), []string{"a.go"}); err == nil {
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
