package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/ticket"
)

// countGit puts a git first on PATH that writes each call's subcommand as one
// line into the file it returns, then runs the real git. With old set it also
// acts like git before 2.31, printing --path-format=absolute back in place of
// using it, as useOldGit does.
func countGit(t *testing.T, old bool) (log string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the git wrapper is a shell script")
	}
	real, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not installed")
	}
	bin := t.TempDir()
	log = filepath.Join(bin, "calls")
	wrapper := "#!/bin/sh\nif [ \"$1\" = -C ]; then echo \"$3\"; else echo \"$1\"; fi >> '" + log + "'\n"
	if old {
		wrapper += "seen=\nn=$#\ni=0\nwhile [ $i -lt $n ]; do a=$1; shift; i=$((i+1))\n" +
			"  if [ \"$a\" = --path-format=absolute ]; then seen=1; else set -- \"$@\" \"$a\"; fi\ndone\n" +
			"[ -n \"$seen\" ] && echo --path-format=absolute\n"
	}
	wrapper += "exec " + real + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

func gitCalls(t *testing.T, log string) []string {
	t.Helper()
	b, err := os.ReadFile(log)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if s := strings.TrimSpace(string(b)); s != "" {
		return strings.Split(s, "\n")
	}
	return nil
}

// The merge-driver check runs on every command. On a shared board it asks git
// once whether this is a working tree and whether jaira's entry excludes the
// board, then once whether the driver is bound — two processes, as before
// NNV39Y added the exclude lookup. Every place a board stands in: shared at
// the top or in a subdirectory, with and without .gitattributes, private in
// info/exclude or in .gitignore, on git before and after 2.31.
func TestTheMergeDriverCheckStartsNoGitForTheIgnoreCheck(t *testing.T) {
	cases := []struct {
		name    string
		sub     bool
		old     bool
		private string // "", "exclude", "gitignore"
		noAttrs bool
		want    []string // the git subcommands, in order
	}{
		{name: "shared at the top", want: []string{"rev-parse", "config"}},
		{name: "shared in a subdirectory", sub: true, want: []string{"rev-parse", "config"}},
		{name: "shared, git before 2.31", old: true, want: []string{"rev-parse", "config"}},
		{name: "shared, subdirectory, git before 2.31", sub: true, old: true, want: []string{"rev-parse", "config"}},
		{name: "without .gitattributes", noAttrs: true, want: nil},
		{name: "private in info/exclude", private: "exclude", want: []string{"rev-parse"}},
		{name: "private in info/exclude, subdirectory", sub: true, private: "exclude", want: []string{"rev-parse"}},
		{name: "private in .gitignore", private: "gitignore", want: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := repoWithCommit(t)
			root := dir
			if c.sub {
				root = filepath.Join(dir, "web app")
				if err := os.MkdirAll(root, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if out, err := runCLI(t, root, "init"); err != nil {
				t.Fatalf("jaira init: %v\n%s", err, out)
			}
			if out, err := runCLI(t, root, "share"); err != nil {
				t.Fatalf("jaira share: %v\n%s", err, out)
			}
			switch c.private {
			case "exclude":
				if _, err := board.AddIgnore(root); err != nil {
					t.Fatal(err)
				}
			case "gitignore":
				if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(board.IgnoreLine+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if c.noAttrs {
				if err := os.Remove(filepath.Join(root, ticket.DirName, ".gitattributes")); err != nil {
					t.Fatal(err)
				}
			}
			s, err := ticket.Discover(root)
			if err != nil {
				t.Fatal(err)
			}
			log := countGit(t, c.old)
			bindDriverIfShared(s)
			got := gitCalls(t, log)
			if strings.Join(got, " ") != strings.Join(c.want, " ") {
				t.Errorf("git calls = %q, want %q", got, c.want)
			}
		})
	}
}
