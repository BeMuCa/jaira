package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ci-code-changed.sh decides whether a pull request runs the test jobs in
// .github/workflows/ci.yaml. Every case is a branch cut from a base commit;
// the script compares the two the way the 'changes' job does.
func TestCICodeChanged(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the script runs in CI on ubuntu; bash on a Windows runner may be WSL's")
	}
	for _, tool := range []string{"bash", "git"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not found: %v", tool, err)
		}
	}
	script, err := filepath.Abs("ci-code-changed.sh")
	if err != nil {
		t.Fatal(err)
	}

	ticket := ".jaira/tickets/01ABC-a-ticket.md"
	cases := []struct {
		name  string
		files []string
		want  string
	}{
		{"a ticket moved lane", []string{ticket}, "false"},
		{"a logbook file with an umlaut", []string{".jaira/logbook/as-1/01X-schließt-die-tabs.md"}, "false"},
		{"tickets and planning notes", []string{ticket, ".planning/research/x.md", "docs/LANE-EVAL.md"}, "false"},
		{"the agent files", []string{"CLAUDE.md", "AGENTS.md", "lanes/README.md"}, "false"},
		{"nothing changed", nil, "false"},
		{"a ticket with its code", []string{ticket, "core/ticket/x.go"}, "true"},
		{"a ticket with its release note", []string{ticket, "core/release/notes/01ABC.md"}, "true"},
		{"embedded markdown", []string{"core/lane/builtin/todo.md"}, "true"},
		{"a catalogue lane", []string{"lanes/secrets-scan.md"}, "true"},
		{"the repository's jaira skill", []string{".claude/skills/jaira/SKILL.md"}, "true"},
		{"the README a test reads", []string{"README.md"}, "true"},
		{"the command reference a test reads", []string{"docs/COMMANDS.md"}, "true"},
		{"a markdown file outside the list", []string{"NOTES.md"}, "true"},
		{"the workflow itself", []string{".github/workflows/ci.yaml"}, "true"},
		{"a name with a space", []string{"core/a b.go"}, "true"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			git := func(args ...string) string {
				t.Helper()
				cmd := exec.Command("git", args...)
				cmd.Dir = dir
				cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v\n%s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			write := func(rel string) {
				t.Helper()
				p := filepath.Join(dir, filepath.FromSlash(rel))
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(rel+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			git("init", "-q", "-b", "master")
			write("go.mod")
			git("add", "-A")
			git("commit", "-q", "-m", "base")
			base := git("rev-parse", "HEAD")
			git("checkout", "-q", "-b", "feature")
			for _, f := range c.files {
				write(f)
			}
			git("add", "-A")
			git("commit", "-q", "--allow-empty", "-m", "change")
			head := git("rev-parse", "HEAD")
			// master moves on with code of its own after the branch was cut:
			// that is master's change, not the pull request's.
			git("checkout", "-q", "master")
			write("core/moved-on.go")
			git("add", "-A")
			git("commit", "-q", "-m", "master moved on")
			base2 := git("rev-parse", "HEAD")

			for _, b := range []string{base, base2} {
				cmd := exec.Command("bash", script, b, head)
				cmd.Dir = dir
				cmd.Env = append(os.Environ(), "GITHUB_OUTPUT=")
				out, err := cmd.Output()
				if err != nil {
					t.Fatalf("script: %v", err)
				}
				if got := strings.TrimSpace(string(out)); got != "code="+c.want {
					t.Errorf("base %s: got %q, want code=%s", b[:7], got, c.want)
				}
			}
		})
	}
}

// A base the repository does not know is an error, not an empty diff that
// skips every test.
func TestCICodeChangedFailsOnUnknownBase(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the script runs in CI on ubuntu")
	}
	dir := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--allow-empty", "-m", "x"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	script, _ := filepath.Abs("ci-code-changed.sh")
	cmd := exec.Command("bash", script, strings.Repeat("0", 40), "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err == nil {
		t.Fatalf("want a failure, got %q", out)
	}
	if strings.Contains(string(out), "code=") {
		t.Errorf("printed a verdict on failure: %q", out)
	}
}
