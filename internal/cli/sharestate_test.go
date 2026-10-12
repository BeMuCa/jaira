package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
)

// Whether a board is shared is decided once, by 'jaira init' (private) and
// 'jaira share' (shared). Nothing a user runs to bring jaira up to date may
// undo that decision. These tests take each of those paths on a board of
// either kind and require .gitignore to come out byte for byte as it went in.

// sharedClone is what a teammate has after cloning a repository whose board
// was shared and committed, and running 'jaira init' in it — the first step
// the README gives them. A fresh clone holds .jaira/.gitattributes and no
// tickets directory (git keeps no empty directories, and tickets travel on
// their refs), so every other command answers "run 'jaira init'" until they
// do. That init is the share-state hazard this fixture checks first.
func sharedClone(t *testing.T) string {
	t.Helper()
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	origin := filepath.Join(t.TempDir(), "origin")
	if err := os.MkdirAll(origin, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, origin, "init", "--quiet", "-b", "master", ".")
	gitRun(t, origin, "config", "user.name", "ada")
	gitRun(t, origin, "config", "user.email", "ada@example.test")
	if out, err := runCLI(t, origin, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, origin, "share"); err != nil {
		t.Fatalf("jaira share: %v\n%s", err, out)
	}
	gitRun(t, origin, "add", ".jaira", ".gitignore")
	gitRun(t, origin, "commit", "--quiet", "-m", "share jaira board")

	clone := filepath.Join(t.TempDir(), "clone")
	gitRun(t, filepath.Dir(clone), "clone", "--quiet", origin, clone)
	cloned := readIgnore(t, clone)
	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("jaira init in the clone: %v\n%s", err, out)
	}
	if got := readIgnore(t, clone); got != cloned {
		t.Fatalf("init in a fresh clone changed .gitignore:\n before %q\n after  %q", cloned, got)
	}
	return clone
}

// privateBoard is a board 'jaira init' made and nobody shared.
func privateBoard(t *testing.T) string {
	t.Helper()
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	dir := t.TempDir()
	gitRun(t, dir, "init", "--quiet", "-b", "master", ".")
	gitRun(t, dir, "config", "user.name", "ada")
	gitRun(t, dir, "config", "user.email", "ada@example.test")
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	gitRun(t, dir, "add", ".gitignore")
	gitRun(t, dir, "commit", "--quiet", "-m", "jaira board stays private")
	return dir
}

func readIgnore(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestUpdatePathsKeepTheShareState(t *testing.T) {
	paths := []struct {
		name string
		run  func(t *testing.T, dir string)
	}{
		{"jaira init run again", func(t *testing.T, dir string) {
			if out, err := runCLI(t, dir, "init"); err != nil {
				t.Fatalf("jaira init: %v\n%s", err, out)
			}
		}},
		{"jaira update", func(t *testing.T, dir string) {
			if out, err := runUpdate(t, dir); err != nil {
				t.Fatalf("jaira update: %v\n%s", err, out)
			}
		}},
		{"jaira self upgrade", func(t *testing.T, dir string) {
			selfTestRelease(t, "v1.1.0", "new-bytes", "v1.1.0", "new-bytes")
			selfUpgradeTarget(t)
			if out, err := runSelf(t, dir, "upgrade"); err != nil {
				t.Fatalf("jaira self upgrade: %v\n%s", err, out)
			}
		}},
		// A newer binary meets a board stamped by an older one: openStore runs
		// the merge-driver binding, the stale-version nudge and lane loading.
		{"first command after a binary switch", func(t *testing.T, dir string) {
			setCurrent(t, "1.1.0")
			if out, err := runCLI(t, dir, "list"); err != nil {
				t.Fatalf("jaira list: %v\n%s", err, out)
			}
		}},
	}
	boards := []struct {
		name    string
		make    func(t *testing.T) string
		ignored bool
	}{
		{"shared clone", sharedClone, false},
		{"private board", privateBoard, true},
	}
	for _, b := range boards {
		for _, p := range paths {
			t.Run(b.name+"/"+p.name, func(t *testing.T) {
				setCurrent(t, "1.0.0")
				dir := b.make(t)
				before := readIgnore(t, dir)
				if got := board.Ignored(dir); got != b.ignored {
					t.Fatalf("fixture: Ignored = %v, want %v\n%s", got, b.ignored, before)
				}

				p.run(t, dir)

				if after := readIgnore(t, dir); after != before {
					t.Errorf(".gitignore changed:\n before %q\n after  %q", before, after)
				}
			})
		}
	}
}

// The share test above must not turn into "init never gitignores an existing
// .jaira/": 'jaira fetch' and lane loading create the directory on boards that
// were never shared, and init still has to make those private.
func TestInitStillMakesAnUntrackedExistingBoardPrivate(t *testing.T) {
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	dir := t.TempDir()
	gitRun(t, dir, "init", "--quiet", "-b", "master", ".")
	if err := os.MkdirAll(filepath.Join(dir, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}

	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("jaira init: %v\n%s", err, out)
	}
	if !board.Ignored(dir) {
		t.Errorf("init left an untracked, never-shared board committable:\n%s", readIgnore(t, dir))
	}
}

// init's closing line and --json 'private' report what .gitignore says, not
// why init left it alone: a board can be in git and still ignored.
func TestInitReportsTheShareStateItLeaves(t *testing.T) {
	const sharedLine = "This board is shared"
	cases := []struct {
		name    string
		make    func(t *testing.T) string
		private bool
	}{
		{"shared clone", sharedClone, false},
		// 'share --undo' re-adds the ignore entry but leaves .gitattributes and
		// the tracked files where they are.
		{"after share --undo", func(t *testing.T) string {
			dir := sharedClone(t)
			if out, err := runCLI(t, dir, "share", "--undo"); err != nil {
				t.Fatalf("jaira share --undo: %v\n%s", err, out)
			}
			return dir
		}, true},
		// A teammate's init in an older build committed the ignore line.
		{"clone with a committed ignore line", func(t *testing.T) string {
			dir := sharedClone(t)
			if _, err := board.AddIgnore(dir); err != nil {
				t.Fatal(err)
			}
			gitRun(t, dir, "config", "user.name", "ada")
			gitRun(t, dir, "config", "user.email", "ada@example.test")
			gitRun(t, dir, "commit", "--quiet", "-am", "an older init made it private")
			return dir
		}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := c.make(t)
			before := readIgnore(t, dir)

			out, err := runCLI(t, dir, "init")
			if err != nil {
				t.Fatalf("jaira init: %v\n%s", err, out)
			}
			if got := strings.Contains(out, sharedLine); got == c.private {
				t.Errorf("output says %q = %v on a board with private = %v:\n%s", sharedLine, got, c.private, out)
			}
			js, err := runCLI(t, dir, "init", "--json")
			if err != nil {
				t.Fatalf("jaira init --json: %v\n%s", err, js)
			}
			var payload map[string]any
			if err := json.Unmarshal([]byte(js), &payload); err != nil {
				t.Fatalf("init --json: %v\n%s", err, js)
			}
			if payload["private"] != c.private {
				t.Errorf("--json private = %v, want %v", payload["private"], c.private)
			}
			if after := readIgnore(t, dir); after != before {
				t.Errorf(".gitignore changed:\n before %q\n after  %q", before, after)
			}
		})
	}
}
