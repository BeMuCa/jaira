package board

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// newRepo is a git repository at a path with a space in it, as on a Windows
// home directory, and its info/exclude.
func newRepo(t *testing.T) (root, exclude string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root = filepath.Join(t.TempDir(), "my repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", root, "init", "--quiet", ".").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return root, filepath.Join(root, ".git", "info", "exclude")
}

func count(s, line string) int {
	n := 0
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) == line {
			n++
		}
	}
	return n
}

// Every shape info/exclude is found in: AddIgnore leaves exactly one entry and
// the person's own lines as they were; RemoveIgnore gives the file back.
func TestAddAndRemoveIgnoreOverEveryShapeOfExclude(t *testing.T) {
	const own = "*.log"
	cases := []struct {
		name    string
		prepare func(t *testing.T, exclude string) // nil: as git init left it
		want    string                             // the file after Add then Remove; "-" for: keep what was there
	}{
		{"as git init left it", nil, "-"},
		{"missing", func(t *testing.T, ex string) { os.Remove(ex) }, ""},
		{"no info directory", func(t *testing.T, ex string) { os.RemoveAll(filepath.Dir(ex)) }, ""},
		{"empty", func(t *testing.T, ex string) { write(t, ex, "") }, ""},
		{"no final newline", func(t *testing.T, ex string) { write(t, ex, own) }, own + "\n"},
		{"CRLF", func(t *testing.T, ex string) { write(t, ex, "# mine\r\n"+own+"\r\n") }, "# mine\r\n" + own + "\r\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root, exclude := newRepo(t)
			if c.prepare != nil {
				c.prepare(t, exclude)
			}
			before, _ := os.ReadFile(exclude)

			changed, err := AddIgnore(root)
			if err != nil || !changed {
				t.Fatalf("AddIgnore = %v, %v", changed, err)
			}
			got, _ := os.ReadFile(exclude)
			if n := count(string(got), IgnoreLine); n != 1 {
				t.Fatalf("exclude holds the entry %d times: %q", n, got)
			}
			if !strings.HasPrefix(string(got), string(before)) {
				t.Errorf("AddIgnore changed the lines that were there:\n before %q\n after  %q", before, got)
			}
			if strings.Contains(string(before), "\r\n") && strings.Count(string(got), "\n") != strings.Count(string(got), "\r\n") {
				t.Errorf("a CRLF file got LF lines: %q", got)
			}
			if !Ignored(root) {
				t.Errorf("Ignored = false after AddIgnore")
			}
			if _, err := os.Stat(filepath.Join(root, ".gitignore")); !os.IsNotExist(err) {
				t.Errorf("AddIgnore wrote a .gitignore in a git repository")
			}
			if changed, _ := AddIgnore(root); changed {
				t.Errorf("a second AddIgnore changed the file")
			}

			if changed, err := RemoveIgnore(root); err != nil || !changed {
				t.Fatalf("RemoveIgnore = %v, %v", changed, err)
			}
			after, _ := os.ReadFile(exclude)
			want := c.want
			if want == "-" {
				want = string(before)
			}
			if string(after) != want {
				t.Errorf("after Add and Remove:\n got  %q\n want %q", after, want)
			}
			if Ignored(root) {
				t.Errorf("Ignored = true after RemoveIgnore")
			}
		})
	}
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A board below the top of the working tree is named from that top in
// info/exclude, whose patterns are relative to it.
func TestABoardInASubdirectoryIsExcludedByItsOwnPath(t *testing.T) {
	root, exclude := newRepo(t)
	sub := filepath.Join(root, "tools", "web app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := AddIgnore(sub); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(exclude)
	if count(string(got), "/tools/web app/.jaira/") != 1 {
		t.Fatalf("exclude = %q, want /tools/web app/.jaira/", got)
	}
	if !Ignored(sub) {
		t.Errorf("Ignored(sub) = false")
	}
	if Ignored(root) {
		t.Errorf("Ignored(root) = true: the subdirectory's entry must not cover the top")
	}
	if err := os.MkdirAll(filepath.Join(sub, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(sub, ".jaira", "x.md"), "x")
	out, err := exec.Command("git", "-C", root, "status", "--porcelain", "--untracked-files=all").Output()
	if err != nil || strings.Contains(string(out), ".jaira") {
		t.Errorf("git sees the subdirectory's board: %v\n%s", err, out)
	}
}

// Two processes making the same board private at once.
func TestParallelAddIgnoreWritesTheEntryOnce(t *testing.T) {
	root, exclude := newRepo(t)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := AddIgnore(root); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	got, _ := os.ReadFile(exclude)
	if n := count(string(got), IgnoreLine); n != 1 {
		t.Errorf("exclude holds the entry %d times after parallel AddIgnore:\n%s", n, got)
	}
}

func TestParsePlaceReadsWhatGitPrints(t *testing.T) {
	// git prints an absolute path the way this system writes one: from / on
	// Unix, from a drive with forward slashes on Windows ("C:/Users/a b").
	abs := "/home/a b/repo/.git"
	if runtime.GOOS == "windows" {
		abs = "C:/Users/a b/repo/.git"
	}
	cases := []struct{ out, exclude, line string }{
		{abs + "\n\n", abs + "/info/exclude", "/.jaira/"},
		{abs + "\r\nweb app/\r\n", abs + "/info/exclude", "/web app/.jaira/"},
		{abs + "\nsub/dir/\n", abs + "/info/exclude", "/sub/dir/.jaira/"},
		{"", "", ""},
		// git before 2.31 prints the option it does not know back.
		{"--path-format=absolute\n.git\nsub/\n", "", ""},
		{".git\n\n", "", ""},
		// Glob characters in a directory name stay literal.
		{abs + "\na[1]/b*?/\n", abs + "/info/exclude", `/a\[1]/b\*\?/.jaira/`},
	}
	for _, c := range cases {
		ex, line := parsePlace(c.out)
		if filepath.ToSlash(ex) != c.exclude || line != c.line {
			t.Errorf("parsePlace(%q) = %q, %q; want %q, %q", c.out, ex, line, c.exclude, c.line)
		}
	}
}

// Outside a git repository there is no info/exclude: the entry goes into
// .gitignore, which is not tracked there.
func TestOutsideGitTheEntryGoesIntoGitignore(t *testing.T) {
	root := t.TempDir()
	if _, err := AddIgnore(root); err != nil {
		t.Fatal(err)
	}
	if !InGitignore(root) || !Ignored(root) {
		t.Errorf("no .gitignore entry outside git")
	}
	if moved, _ := MoveIgnore(root); moved {
		t.Errorf("MoveIgnore moved the entry outside git")
	}
}

// Two boards in one repository share info/exclude: lifting one leaves the
// other's entry with its comment.
func TestRemoveIgnoreLeavesTheOtherBoardsEntry(t *testing.T) {
	root, exclude := newRepo(t)
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{root, sub} {
		if _, err := AddIgnore(dir); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := RemoveIgnore(sub); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(exclude)
	if !strings.Contains(string(got), privateComment+"\n"+IgnoreLine+"\n") || strings.Contains(string(got), "/sub/.jaira/") {
		t.Errorf("exclude after lifting the subdirectory's board:\n%s", got)
	}
	if count(string(got), privateComment) != 1 {
		t.Errorf("want one comment left, for the top board:\n%s", got)
	}
}

// A directory named with glob characters is excluded as itself.
func TestAGlobCharacterDirectoryIsExcludedLiterally(t *testing.T) {
	root, _ := newRepo(t)
	sub := filepath.Join(root, "a[1]")
	if err := os.MkdirAll(filepath.Join(sub, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(sub, ".jaira", "x.md"), "x")
	if _, err := AddIgnore(sub); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", root, "status", "--porcelain", "--untracked-files=all").Output()
	if err != nil || strings.Contains(string(out), ".jaira") {
		t.Errorf("git sees the board in a[1]: %v\n%s", err, out)
	}
}

// dropEntry takes out the blank line directly after the entry it removes and
// no other: a blank line further down belongs to what stands above it.
func TestDropEntryTakesOnlyTheBlankLineThatFollowsTheEntry(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"entry at the end", "a\n\n# jaira board — private\n/.jaira/\n", "a\n"},
		{"blank right after", "# jaira board — private\n/.jaira/\n\nfoo\n", "foo\n"},
		{"no blank right after", "# jaira board — private\n/.jaira/\nfoo\n\nbar\n", "foo\n\nbar\n"},
		{"entry without its comment", "/.jaira/\n\nfoo\n", "\nfoo\n"},
		{"no blank line before, at the end", "a\n# jaira board — private\n/.jaira/\n", "a\n"},
		{"CRLF, no blank line before, at the end", "a\r\n# jaira board — private\r\n/.jaira/\r\n", "a\r\n"},
		{"CRLF, no blank right after", "# jaira board — private\r\n/.jaira/\r\nfoo\r\n\r\nbar\r\n", "foo\r\n\r\nbar\r\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "exclude")
			write(t, path, c.in)
			if changed, err := dropEntry(path, privateMarker, func(s string) bool { return s == IgnoreLine }); err != nil || !changed {
				t.Fatalf("dropEntry = %v, %v", changed, err)
			}
			if got, _ := os.ReadFile(path); string(got) != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

// An unanchored ".jaira/" a person wrote into info/exclude excludes every
// board in the repository. Sharing the top board takes out only jaira's own
// anchored entry: the person's line stays, the board in the subdirectory
// stays private, and Ignored still reports the top board as excluded.
func TestRemoveIgnoreLeavesAnUnanchoredLineAPersonWrote(t *testing.T) {
	root, exclude := newRepo(t)
	write(t, exclude, ".jaira/\n")
	if !Ignored(root) {
		t.Fatalf("Ignored = false with .jaira/ in info/exclude")
	}
	if changed, _ := AddIgnore(root); changed {
		t.Errorf("AddIgnore wrote an entry beside a line that already excludes the board")
	}
	if changed, err := RemoveIgnore(root); err != nil || changed {
		t.Fatalf("RemoveIgnore = %v, %v; want nothing changed", changed, err)
	}
	if got, _ := os.ReadFile(exclude); string(got) != ".jaira/\n" {
		t.Errorf("exclude after RemoveIgnore: %q, want the person's line alone", got)
	}
	if !Ignored(root) {
		t.Errorf("Ignored = false although the person's line still excludes the board")
	}
}

// A file that exists but cannot be read is never replaced by "the entry
// alone": every writer that appends (AddIgnore into info/exclude or the
// .gitignore fallback, MoveIgnore, AddLanesIgnore) fails instead, and the
// file keeps every byte.
func TestAppendEntryNeverOverwritesAFileItCannotRead(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod 000 does not make a file unreadable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root reads a file with mode 000")
	}
	for _, name := range []string{"info/exclude", ".gitignore"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			want := []byte("# a person's rules\n*.log\n")
			if err := os.WriteFile(path, want, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, 0); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { os.Chmod(path, 0o644) })

			if err := appendEntry(path, privateComment, IgnoreLine); err == nil {
				t.Fatal("appendEntry on an unreadable file: err = nil, want the read error")
			}
			if err := os.Chmod(path, 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want) {
				t.Fatalf("file changed:\n got %q\nwant %q", got, want)
			}
		})
	}
	t.Run("missing file is created", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "exclude")
		if err := appendEntry(path, privateComment, IgnoreLine); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if !strings.Contains(string(got), IgnoreLine) {
			t.Fatalf("got %q, want the entry", got)
		}
	})
}

// TestIgnoreFileNameIsRelativeThroughASymlinkedRoot: git prints the common dir
// resolved, as /private/var for macOS's /var temp dirs, so a root reached
// through a symlink must still name the file .git/info/exclude, not an
// absolute path.
func TestIgnoreFileNameIsRelativeThroughASymlinkedRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need privileges on Windows")
	}
	root, _ := newRepo(t)
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	if got := IgnoreFileName(link); got != ".git/info/exclude" {
		t.Errorf("IgnoreFileName through a symlink = %q, want .git/info/exclude", got)
	}
	if got := IgnoreFileName(root); got != ".git/info/exclude" {
		t.Errorf("IgnoreFileName = %q, want .git/info/exclude", got)
	}
}
