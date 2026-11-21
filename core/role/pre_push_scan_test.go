package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// scanRepo is a scratch repository with a branch "feat" off "main", the shape
// pre-push-scan.sh sees before a push: main is the base, feat is what goes.
type scanRepo struct {
	t   *testing.T
	dir string
	env []string
}

func newScanRepo(t *testing.T) *scanRepo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("no git on this machine: %v", err)
	}
	home := t.TempDir()
	r := &scanRepo{t: t, dir: t.TempDir(), env: append(os.Environ(),
		"HOME="+home, "XDG_CONFIG_HOME="+home, "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.org",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.org")}
	r.git("init", "-q", "-b", "main")
	// A user's own diff settings must not change what the scan reads.
	r.git("config", "diff.mnemonicPrefix", "true")
	r.git("config", "core.autocrlf", "false")
	r.commit("base", map[string]string{"README": "base\n"})
	r.git("switch", "-q", "-c", "feat")
	return r
}

func (r *scanRepo) git(args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir, cmd.Env = r.dir, r.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r *scanRepo) commit(msg string, files map[string]string, env ...string) {
	r.t.Helper()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(r.dir, name), []byte(body), 0o644); err != nil {
			r.t.Fatal(err)
		}
		r.git("add", "--", name)
	}
	cmd := exec.Command("git", "commit", "-q", "--allow-empty", "-F", "-")
	cmd.Dir, cmd.Env, cmd.Stdin = r.dir, append(append([]string{}, r.env...), env...), strings.NewReader(msg)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("commit: %v\n%s", err, out)
	}
}

// The class: a mention leaves the machine through any part of any commit a
// push publishes. One case per place it can hide and per input that has
// broken a scanner of this kind before (AX8K02). The mention and the machine
// name must never appear in the output, whatever the case.
func TestPrePushScan(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "jaira-role-pr", "scripts", "pre-push-scan.sh")
	const mention, host = "acmecorp", "desktop-x1"
	deny := fixtureFile(t, "deny", "# the list\r\n\r\nACME-?CORP\r\n")
	huge := mention + " in the subject\n\n" + strings.Repeat("filler line of a long commit message\n", 4000)

	for _, c := range []struct {
		name  string
		setup func(r *scanRepo)
		base  string
		list  string // the denylist path; "" is the file above
		code  int
		want  []string
		never []string
	}{
		{"no commits", func(r *scanRepo) {}, "", "", 0, []string{"0 commit(s)"}, nil},
		{"added line", func(r *scanRepo) {
			r.commit("feat: a", map[string]string{"a.go": "package a\n// for AcmeCorp\n"})
		}, "", "", 3, []string{"a.go:2: matches the denylist", "do not push"}, nil},
		{"added then removed again", func(r *scanRepo) {
			r.commit("feat: a", map[string]string{"a.go": "acme-corp\n"})
			r.commit("fix: a", map[string]string{"a.go": "clean\n"})
		}, "", "", 3, []string{"a.go:1: matches the denylist"}, nil},
		{"empty new file named after it", func(r *scanRepo) {
			r.commit("feat: e", map[string]string{"acmecorp-notes": ""})
		}, "", "", 3, []string{"file 1 of its changed files, by its name: matches the denylist"}, nil},
		{"binary file", func(r *scanRepo) {
			r.commit("feat: b", map[string]string{"blob.bin": "\x00\x01acmecorp\x00\xff"})
		}, "", "", 3, []string{"blob.bin:1: matches"}, nil},
		{"space, non-ASCII name and CRLF", func(r *scanRepo) {
			r.commit("feat: s", map[string]string{"my notes é.txt": "first\r\nacmecorp\r\n"})
		}, "", "", 3, []string{"my notes é.txt:2: matches"}, nil},
		{"named after it, content too", func(r *scanRepo) {
			r.commit("feat: n", map[string]string{"ok.txt": "fine\n", "acmecorp.txt": "x\nacmecorp\n"})
		}, "", "", 3, []string{"file 1 of its changed files:2: matches"}, nil},
		{"merge: its own resolution, not what it brings in", func(r *scanRepo) {
			r.commit("feat: c", map[string]string{"c.txt": "feat side\n"})
			r.git("switch", "-q", "main")
			r.commit("main: c", map[string]string{"c.txt": "main side\n", "m.txt": "acmecorp on main\n"})
			r.git("switch", "-q", "feat")
			cmd := exec.Command("git", "merge", "-q", "main")
			cmd.Dir, cmd.Env = r.dir, r.env
			_ = cmd.Run() // conflicts on purpose
			r.commit("merge main", map[string]string{"c.txt": "resolved for acmecorp\n"})
		}, "main", "", 3, []string{"c.txt:1: matches"}, []string{"m.txt"}},
		{"merge: a line kept from one parent inside a conflict hunk is not its own", func(r *scanRepo) {
			r.git("switch", "-q", "main")
			r.commit("main: base", map[string]string{"f": "a\nb\nc\n"})
			r.git("switch", "-q", "feat")
			r.git("merge", "-q", "--ff-only", "main")
			r.commit("feat: f", map[string]string{"f": "a\nFEAT\nc\n"})
			r.git("switch", "-q", "main")
			r.commit("main: f", map[string]string{"f": "a\nMAIN\nacmecorp\nc\n"})
			r.git("switch", "-q", "feat")
			cmd := exec.Command("git", "merge", "-q", "main")
			cmd.Dir, cmd.Env = r.dir, r.env
			_ = cmd.Run() // conflicts on purpose
			r.commit("merge main", map[string]string{"f": "a\nresolved\nacmecorp\nc\n"})
		}, "main", "", 0, []string{"denylist: 1 pattern(s)"}, []string{"matches"}},
		{"commit message over 100 KB", func(r *scanRepo) {
			r.commit(huge, nil)
		}, "", "", 3, []string{"commit message: matches"}, nil},
		{"author email", func(r *scanRepo) {
			r.commit("feat: x", nil, "GIT_AUTHOR_EMAIL=me@acmecorp.example")
		}, "", "", 3, []string{"author email: matches"}, []string{"committer email"}},
		{"committer email", func(r *scanRepo) {
			r.commit("feat: x", nil, "GIT_COMMITTER_EMAIL=me@acme-corp.example")
		}, "", "", 3, []string{"committer email: matches"}, []string{"author email"}},
		{"author name", func(r *scanRepo) {
			r.commit("feat: x", nil, "GIT_AUTHOR_NAME=Alex AcmeCorp")
		}, "", "", 3, []string{"author name: matches"}, []string{"committer name"}},
		{"committer name", func(r *scanRepo) {
			r.commit("feat: x", nil, "GIT_COMMITTER_NAME=Acme-Corp build")
		}, "", "", 3, []string{"committer name: matches"}, []string{"author name"}},
		{"branch name", func(r *scanRepo) {
			r.git("branch", "-m", "acmecorp-fix")
			r.commit("feat: b", map[string]string{"b.txt": "clean\n"})
		}, "", "", 3, []string{"branch name: matches the denylist"}, nil},
		{"this machine's name only warns", func(r *scanRepo) {
			r.commit("feat: h", map[string]string{"t.md": "claimed-by: DESKTOP-X1-123\n"})
		}, "", "", 0, []string{"warning: ", "t.md:1: names this machine", "the push may go on"}, []string{"do not push"}},
		{"no denylist", func(r *scanRepo) {
			r.commit("feat: a", map[string]string{"a.go": "acmecorp\n"})
		}, "", "missing", 0, []string{"its strings were not checked"}, []string{"matches the denylist"}},
		{"denylist without a pattern", func(r *scanRepo) {
			r.commit("feat: a", map[string]string{"a.go": "acmecorp\n"})
		}, "", "empty", 0, []string{"holds no pattern"}, []string{"matches the denylist"}},
		{"invalid pattern never passes", func(r *scanRepo) {
			r.commit("feat: a", map[string]string{"a.go": "acmecorp\n"})
		}, "", "invalid", 2, []string{"holds a pattern grep cannot read"}, []string{"a(", "Unmatched"}},
		{"base not a commit", func(r *scanRepo) {}, "nope/main", "", 2, []string{"fetch it first"}, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := newScanRepo(t)
			c.setup(r)
			list := deny
			switch c.list {
			case "missing":
				list = filepath.Join(t.TempDir(), "none")
			case "empty":
				list = fixtureFile(t, "empty", "# only a comment\r\n\n   \n")
			case "invalid":
				list = fixtureFile(t, "invalid", "a(\n")
			}
			base := c.base
			if base == "" {
				base = "main"
			}
			cmd := exec.Command(bash, script, base)
			cmd.Dir = r.dir
			cmd.Env = append(append([]string{}, r.env...), "JAIRA_MENTION_DENYLIST="+list, "JAIRA_SCAN_HOSTNAME="+strings.ToUpper(host)+".lan")
			b, err := cmd.CombinedOutput()
			out := string(b)
			code := 0
			if ee, ok := err.(*exec.ExitError); ok {
				code = ee.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if code != c.code {
				t.Errorf("exit %d, want %d\n%s", code, c.code, out)
			}
			for _, w := range c.want {
				if !strings.Contains(out, w) {
					t.Errorf("output lacks %q\n%s", w, out)
				}
			}
			for _, n := range append(c.never, mention, "acme-corp", host) {
				if strings.Contains(strings.ToLower(out), strings.ToLower(n)) {
					t.Errorf("output carries %q\n%s", n, out)
				}
			}
		})
	}
}

// The class: the denylist is read differently than it was written, and a
// pattern lost on the way reads as "no hit" (AX8K02, bughunt 2026-11-21). A
// list is read as written — UTF-8 with or without a BOM, UTF-16 with its BOM —
// or refused with exit 2; and a non-ASCII letter folds its case whatever
// locale the caller runs in.
func TestPrePushScanReadsTheListAsWritten(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "jaira-role-pr", "scripts", "pre-push-scan.sh")
	utf16le := func(s string) string {
		b := []byte{0xff, 0xfe}
		for _, r := range s {
			b = append(b, byte(r), byte(r>>8))
		}
		return string(b)
	}
	for _, c := range []struct {
		name, list, content, locale string
		code                        int
		want                        string
	}{
		{"UTF-8", "# names\nmüller\n", "MÜLLER GmbH\n", "", 3, "matches the denylist"},
		{"UTF-8 with BOM", "\xef\xbb\xbfacmecorp\r\n", "for acmecorp\n", "", 3, "matches the denylist"},
		{"UTF-16LE with BOM", utf16le("# list\r\nacmecorp\r\n"), "for acmecorp\n", "", 3, "matches the denylist"},
		{"UTF-16 without BOM", utf16le("acmecorp\n")[2:], "for acmecorp\n", "", 2, "NUL bytes"},
		{"Latin-1", "acmecorp\nm\xfcller\n", "for acmecorp\n", "", 2, "not UTF-8"},
		{"lines ending in CR only", "foo\racmecorp\r", "for acmecorp\n", "", 3, "matches the denylist"},
		{"CRLF", "foo\r\nacmecorp\r\n", "for acmecorp.\n", "", 3, "matches the denylist"},
		{"C locale still folds non-ASCII", "müller\n", "MÜLLER GmbH\n", "C", 3, "matches the denylist"},
		{"a pattern is taken as written, its trailing blank included", "acme\\ \n", "for acme corp\n", "", 3, "matches the denylist"},
		{"a non-ASCII pattern matches UTF-8 text only, as the usage says", "müller\n", "M\xfcller GmbH\n", "", 0, "1 pattern(s)"},
		{"a commit that is not UTF-8, in a UTF-8 locale", "acmecorp\n", "caf\xe9 acmecorp\n", "C.UTF-8", 3, "matches the denylist"},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := newScanRepo(t)
			r.commit("feat: x", map[string]string{"x.txt": c.content})
			list := fixtureFile(t, "deny", c.list)
			cmd := exec.Command(bash, script, "main")
			cmd.Dir = r.dir
			cmd.Env = append(append([]string{}, r.env...), "JAIRA_MENTION_DENYLIST="+list, "JAIRA_SCAN_HOSTNAME=")
			if c.locale != "" {
				cmd.Env = append(cmd.Env, "LC_ALL="+c.locale, "LANG="+c.locale, "LC_CTYPE="+c.locale)
			}
			b, err := cmd.CombinedOutput()
			out := string(b)
			code := 0
			if ee, ok := err.(*exec.ExitError); ok {
				code = ee.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if code != c.code || !strings.Contains(out, c.want) {
				t.Errorf("exit %d, want %d with %q\n%s", code, c.code, c.want, out)
			}
			if strings.Contains(strings.ToLower(out), "acmecorp") || strings.Contains(strings.ToLower(out), "müller") {
				t.Errorf("output carries the mention\n%s", out)
			}
			if strings.Contains(out, "multibyte") {
				t.Errorf("a tool complained about the text's encoding\n%s", out)
			}
		})
	}
}
