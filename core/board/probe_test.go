package board

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/BeMuCa/jaira/core/gitrepo"
)

// parseProbe reads rawPlace's lines with --is-inside-work-tree after them,
// over every shape git prints: LF and CRLF, a board at the top and in a
// subdirectory, a git directory that is no working tree, and git before 2.31,
// which prints --path-format back as an extra first line.
func TestParseProbeReadsWhatGitPrints(t *testing.T) {
	abs := "/home/a b/repo/.git"
	if runtime.GOOS == "windows" {
		abs = "C:/Users/a b/repo/.git"
	}
	cases := []struct {
		out           string
		inRepo        bool
		exclude, line string
	}{
		{abs + "\n\ntrue\n", true, abs + "/info/exclude", "/.jaira/"},
		{abs + "\r\n\r\ntrue\r\n", true, abs + "/info/exclude", "/.jaira/"},
		{abs + "\nsub/dir/\ntrue\n", true, abs + "/info/exclude", "/sub/dir/.jaira/"},
		{abs + "\r\nweb app/\r\ntrue\r\n", true, abs + "/info/exclude", "/web app/.jaira/"},
		{abs + "\n\nfalse\n", false, abs + "/info/exclude", "/.jaira/"},
		{"--path-format=absolute\n.git\nsub/\ntrue\n", true, "", ""},
		{"--path-format=absolute\n.git\n\nfalse\n", false, "", ""},
		{"", false, "", ""},
	}
	for _, c := range cases {
		inRepo, ex, prefix := parseProbe(c.out)
		line := ""
		if ex != "" {
			line = entryLine(prefix)
		}
		if inRepo != c.inRepo || filepath.ToSlash(ex) != c.exclude || line != c.line {
			t.Errorf("parseProbe(%q) = %v, %q, %q; want %v, %q, %q", c.out, inRepo, ex, line, c.inRepo, c.exclude, c.line)
		}
	}
}

// Probe gives the answers IsRepo and Ignored give, asked one after the other,
// for every place a board stands in: outside git, private at the top or in a
// subdirectory (in info/exclude or the root .gitignore), shared.
func TestProbeAgreesWithIsRepoAndIgnored(t *testing.T) {
	outside := t.TempDir()
	top, _ := newRepo(t)
	sub := filepath.Join(top, "web app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name  string
		root  string
		setup func(root string)
	}{
		{"outside git", outside, func(string) {}},
		{"outside git, private", outside, func(r string) { AddIgnore(r) }},
		{"shared at the top", top, func(string) {}},
		{"private at the top", top, func(r string) { AddIgnore(r) }},
		{"private in .gitignore", top, func(r string) { write(t, filepath.Join(r, ".gitignore"), IgnoreLine+"\n") }},
		{"shared in a subdirectory", sub, func(string) {}},
		{"private in a subdirectory", sub, func(r string) { AddIgnore(r) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			RemoveIgnore(c.root)
			c.setup(c.root)
			defer RemoveIgnore(c.root)
			inRepo, ignored := Probe(c.root)
			wantIgnored := Ignored(c.root)
			if ignored != wantIgnored {
				t.Errorf("Probe ignored = %v, Ignored = %v", ignored, wantIgnored)
			}
			if !wantIgnored {
				if want := (&gitrepo.Repo{Dir: c.root}).IsRepo(); inRepo != want {
					t.Errorf("Probe inRepo = %v, IsRepo = %v", inRepo, want)
				}
			}
		})
	}
}
