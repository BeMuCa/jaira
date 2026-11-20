package shipped

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fixtureRepo is a throwaway git repository holding scripts/shipped-history.sh
// and the history the cases below need, tagged like releases:
//
//	v0.1.0  core/role/builtin/r/SKILL.md = "a\n"
//	v0.2.0  the file deleted — a tag without it
//	v0.3.0  the file back, with CRLF line endings: "b\r\n"
type fixtureRepo struct {
	t   *testing.T
	dir string
}

func newFixtureRepo(t *testing.T) *fixtureRepo {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the script runs in CI on ubuntu; bash on a Windows runner may be WSL's")
	}
	for _, tool := range []string{"bash", "git", "sha256sum"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not found: %v", tool, err)
		}
	}
	r := &fixtureRepo{t: t, dir: t.TempDir()}
	script, err := os.ReadFile(filepath.Join("..", "..", "scripts", "shipped-history.sh"))
	if err != nil {
		t.Fatal(err)
	}
	r.write("scripts/shipped-history.sh", string(script))
	r.write("core/shipped/shipped.sha256", "# header\n")
	r.git("init", "-q")
	r.git("config", "core.autocrlf", "false")
	r.write("core/role/builtin/r/SKILL.md", "a\n")
	r.commit("v1")
	r.git("tag", "v0.1.0")
	r.git("rm", "-q", "core/role/builtin/r/SKILL.md")
	r.write("other.txt", "x\n")
	r.commit("deleted")
	r.git("tag", "v0.2.0")
	r.write("core/role/builtin/r/SKILL.md", "b\r\n")
	r.commit("back, with CRLF")
	r.git("tag", "v0.3.0")
	return r
}

func (r *fixtureRepo) write(rel, content string) {
	r.t.Helper()
	p := filepath.Join(r.dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o755); err != nil {
		r.t.Fatal(err)
	}
}

func (r *fixtureRepo) run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *fixtureRepo) git(args ...string) string {
	r.t.Helper()
	out, err := r.run("git", args...)
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return out
}

func (r *fixtureRepo) commit(msg string) {
	r.t.Helper()
	r.git("add", "-A")
	r.git("commit", "-q", "-m", msg)
}

// check runs the script's --check and reports its exit status and output.
func (r *fixtureRepo) check() (bool, string) {
	out, err := r.run("bash", "scripts/shipped-history.sh", "--check")
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			r.t.Fatalf("running the script: %v", err)
		}
		return false, out
	}
	return true, out
}

// generated is what the script writes as the manifest of this history.
func (r *fixtureRepo) generated() string {
	r.t.Helper()
	out, err := r.run("bash", "scripts/shipped-history.sh")
	if err != nil {
		r.t.Fatalf("generating: %v\n%s", err, out)
	}
	return out
}

const key = "core/role/builtin/r/SKILL.md"

func TestCheckHoldsTheManifestAgainstEveryTag(t *testing.T) {
	r := newFixtureRepo(t)
	full := r.generated()
	lineA := Line(key, []byte("a\n"))
	lineCRLF := Line(key, []byte("b\r\n"))
	for _, want := range []string{lineA, lineCRLF} {
		if !strings.Contains(full, want+"\n") {
			t.Fatalf("the generated manifest lacks %q:\n%s", want, full)
		}
	}

	cases := []struct {
		name     string
		manifest string
		ok       bool
		names    string // a line the output must name
	}{
		{"the generated manifest", full, true, ""},
		// v0.1.0 shipped "a\n"; the file was deleted in v0.2.0 and came back
		// different, so only the tag still knows that version.
		{"a version of a file deleted and recreated later", strings.Replace(full, lineA+"\n", "", 1), false, lineA},
		// git stores the CRLF bytes, and a build of that commit embedded them:
		// the LF version's hash is a different version.
		{"the LF hash in place of the CRLF version", strings.Replace(full, lineCRLF, Line(key, []byte("b\n")), 1), false, lineCRLF},
		{"a manifest with CRLF line endings", strings.ReplaceAll(full, "\n", "\r\n"), true, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r.write("core/shipped/shipped.sha256", c.manifest)
			ok, out := r.check()
			if ok != c.ok {
				t.Fatalf("check passed = %v, want %v:\n%s", ok, c.ok, out)
			}
			if c.names != "" && !strings.Contains(out, c.names) {
				t.Errorf("the output does not name %q:\n%s", c.names, out)
			}
		})
	}
}

// A clone without tags — actions/checkout's default — has nothing to check
// against, and a check that checked nothing must not pass.
func TestCheckFailsWithoutTags(t *testing.T) {
	r := newFixtureRepo(t)
	r.write("core/shipped/shipped.sha256", r.generated())
	for _, tag := range []string{"v0.1.0", "v0.2.0", "v0.3.0"} {
		r.git("tag", "-d", tag)
	}
	if ok, out := r.check(); ok || !strings.Contains(out, "no v* tags") {
		t.Fatalf("check without tags passed = %v:\n%s", ok, out)
	}
}

// The repository's .gitattributes merges the manifest with merge=union: two
// branches that both add a line at the same place keep both, where a plain
// merge stops on a conflict.
func TestManifestMergesWithoutConflict(t *testing.T) {
	r := newFixtureRepo(t)
	attrs, err := os.ReadFile(filepath.Join("..", "..", ".gitattributes"))
	if err != nil {
		t.Fatal(err)
	}
	r.write(".gitattributes", string(attrs))
	r.write("core/shipped/shipped.sha256", r.generated())
	r.commit("manifest")
	base := strings.TrimSpace(r.git("rev-parse", "--abbrev-ref", "HEAD"))
	one := Line("core/role/builtin/r/one.md", []byte("1\n"))
	two := Line("core/role/builtin/r/two.md", []byte("2\n"))
	manifest := r.generated()

	r.git("switch", "-q", "-c", "one")
	r.write("core/shipped/shipped.sha256", manifest+one+"\n")
	r.commit("one")
	r.git("switch", "-q", base)
	r.write("core/shipped/shipped.sha256", manifest+two+"\n")
	r.commit("two")
	if out, err := r.run("git", "merge", "-q", "--no-edit", "one"); err != nil {
		t.Fatalf("the merge stopped: %v\n%s", err, out)
	}
	got, err := os.ReadFile(filepath.Join(r.dir, "core", "shipped", "shipped.sha256"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{one, two} {
		if !strings.Contains(string(got), want+"\n") {
			t.Errorf("the merged manifest lost %q", want)
		}
	}
}

// A tag git cannot list a tree for stops the check; skipping it would compare
// fewer versions and pass.
func TestCheckFailsOnAnUnreadableTag(t *testing.T) {
	r := newFixtureRepo(t)
	r.write("core/shipped/shipped.sha256", r.generated())
	blob := strings.TrimSpace(r.git("rev-parse", "v0.1.0:"+key))
	r.git("tag", "v0.9.0", blob)
	if ok, out := r.check(); ok {
		t.Fatalf("check with a tag naming a blob passed:\n%s", out)
	}
}

// A version git cannot read back — a partial clone offline, a damaged object —
// stops the generator and the check alike. Hashed as the empty file it would
// drop the real version's line from the manifest, or have --check offer the
// empty file's hash as the line to append.
func TestUnreadableBlobStopsTheScript(t *testing.T) {
	const emptyFile = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	for _, args := range [][]string{{}, {"--check"}} {
		t.Run(strings.Join(append([]string{"shipped-history.sh"}, args...), " "), func(t *testing.T) {
			r := newFixtureRepo(t)
			r.write("core/shipped/shipped.sha256", r.generated())
			blob := strings.TrimSpace(r.git("rev-parse", "v0.1.0:"+key))
			if err := os.Remove(filepath.Join(r.dir, ".git", "objects", blob[:2], blob[2:])); err != nil {
				t.Fatal(err)
			}
			out, err := r.run("bash", append([]string{"scripts/shipped-history.sh"}, args...)...)
			if err == nil {
				t.Fatalf("the script passed over an unreadable blob:\n%s", out)
			}
			if strings.Contains(out, emptyFile) {
				t.Errorf("the script offered the empty file's hash as a version:\n%s", out)
			}
		})
	}
}
