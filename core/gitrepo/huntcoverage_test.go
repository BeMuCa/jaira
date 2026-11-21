package gitrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// What a hunt says about the code now, from commits only: a commit of the
// ticket on a hunted file makes it stale; other people's commits, a merge of
// master and board data do not; a hunted file changed and not committed is
// unknown, whenever it was written; a checkout that does not descend from
// the hunt is stale, and a file list that names no file reads all code.
func TestHuntCoverage(t *testing.T) {
	const id = "01M7WWZWSDHK37DGGHJFDF0YD6" // handle DF0YD6
	at := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	later := at.Add(time.Hour)
	mkdir := func(t *testing.T, dir, sub string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		name  string
		then  func(t *testing.T, dir string)
		paths []string
		want  HuntState
	}{
		{"nothing since", func(*testing.T, string) {}, []string{"a.go"}, HuntCovers},
		{"another file", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "fix(DF0YD6): b", later)
		}, []string{"a.go"}, HuntCovers},
		{"a hunted file, by this ticket", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "fix(DF0YD6): a", later)
		}, []string{"a.go"}, HuntStale},
		{"a hunted file, by another ticket", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "fix(ZZZZZZ): a", later)
		}, []string{"a.go"}, HuntCovers},
		{"a hunted file changed, not committed", func(t *testing.T, dir string) {
			write(t, dir, "a.go", "a3")
		}, []string{"a.go"}, HuntUnknown},
		{"an uncommitted edit from long before the hunt (conversational mode)", func(t *testing.T, dir string) {
			write(t, dir, "a.go", "a3")
			old := time.Now().Add(-24 * time.Hour)
			if err := os.Chtimes(filepath.Join(dir, "a.go"), old, old); err != nil {
				t.Fatal(err)
			}
		}, []string{"a.go"}, HuntUnknown},
		{"a hunted file deleted, not committed", func(t *testing.T, dir string) {
			if err := os.Remove(filepath.Join(dir, "a.go")); err != nil {
				t.Fatal(err)
			}
		}, []string{"a.go"}, HuntUnknown},
		{"an untracked hunted file", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, []string{"a.go", "new.go"}, HuntUnknown},
		{"an untracked file outside the hunted paths", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, []string{"a.go"}, HuntCovers},
		{"an untracked file, no paths", func(t *testing.T, dir string) {
			write(t, dir, "new.go", "n")
		}, nil, HuntUnknown},
		{"a path with a space", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"my dir/c.go": "c2"}, "fix(DF0YD6): c", later)
		}, []string{"my dir/c.go"}, HuntStale},
		{"a list split at a space names no file", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"my dir/c.go": "c2"}, "fix(DF0YD6): c", later)
		}, []string{"my", "dir/c.go"}, HuntStale},
		{"no paths reads all code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{"b.go": "b2"}, "fix(DF0YD6): b", later)
		}, nil, HuntStale},
		{"a ticket file only, no paths", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/tickets/x.md": "t2"}, "DF0YD6 note", later)
		}, nil, HuntCovers},
		{"tracked board data changed and committed, no paths", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/tags": "ui red"}, "chore(DF0YD6): tag", later)
			write(t, dir, ".jaira/tags", "ui red\ngates blue")
		}, nil, HuntCovers},
		{"untracked board data, listed", func(t *testing.T, dir string) {
			mkdir(t, dir, ".jaira/milestones")
			write(t, dir, ".jaira/milestones/q4-ABC123.md", "m")
		}, []string{"a.go", ".jaira/milestones/q4-ABC123.md"}, HuntCovers},
		{"a lane prompt is code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/lanes/testing.md": "p2"}, "fix(DF0YD6): prompt", later)
		}, nil, HuntStale},
		{"a named board's lane prompt is code", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/boards/video/lanes/testing.md": "p2"}, "fix(DF0YD6): prompt", later)
		}, nil, HuntStale},
		{"a named board's other data is not", func(t *testing.T, dir string) {
			fixtureCommit(t, dir, map[string]string{".jaira/boards/video/notes.md": "n"}, "chore(DF0YD6): board", later)
		}, nil, HuntCovers},
		{"master merged in, touching the hunted file", func(t *testing.T, dir string) {
			git(t, dir, "checkout", "-q", "-b", "side")
			fixtureCommit(t, dir, map[string]string{"a.go": "from master"}, "fix(ZZZZZZ): a", later)
			git(t, dir, "checkout", "-q", "-")
			git(t, dir, "merge", "-q", "--no-ff", "-m", "chore(DF0YD6): merge master", "side")
		}, []string{"a.go"}, HuntCovers},
		{"a checkout that does not descend from the hunt", func(t *testing.T, dir string) {
			git(t, dir, "checkout", "-q", "--orphan", "other")
			fixtureCommit(t, dir, map[string]string{"z.go": "z"}, "other", later)
		}, []string{"a.go"}, HuntStale},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, repo := fixtureRepo(t)
			hunt := fixtureCommit(t, dir, map[string]string{"a.go": "a", "b.go": "b", "my dir/c.go": "c", ".jaira/tickets/x.md": "t"}, "feat(DF0YD6): start", at)
			tc.then(t, dir)
			got, err := repo.HuntCoverage(hunt[:8], id, tc.paths)
			if err != nil || got != tc.want {
				t.Errorf("HuntCoverage = %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	t.Run("an unknown commit", func(t *testing.T) {
		dir, repo := fixtureRepo(t)
		fixtureCommit(t, dir, map[string]string{"a.go": "a"}, "start", at)
		if _, err := repo.HuntCoverage("deadbeef", id, []string{"a.go"}); err == nil {
			t.Error("an unknown commit read without an error")
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
