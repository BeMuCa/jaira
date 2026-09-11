package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Two people work the same ticket on two branches and the branches are merged.
// This is the case line-based merging gets wrong on the single most common
// operation the whole tool performs — one line, status: — so it is tested
// against real git with the real driver rather than by calling the merge
// function directly.
func TestTwoBranchesOnOneTicketMergeFieldAware(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the binary")
	}
	for _, tool := range []string{"git", "go"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not on PATH", tool)
		}
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	// The driver is an executable git invokes, so the test needs a real one.
	// With the .exe on Windows: without it the file is built but exec cannot
	// find it, and the failure reads as "executable file not found in %PATH%"
	// rather than as the naming rule it is.
	bin := filepath.Join(root, "jaira"+exeSuffix())
	build := exec.Command("go", "build", "-o", bin, "github.com/BeMuCa/jaira/cmd/jaira")
	build.Dir = repoRoot(t)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building the binary: %v\n%s", err, out)
	}

	repo := filepath.Join(root, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	run(t, repo, "git", "init", "--quiet", "-b", "master", ".")
	run(t, repo, "git", "config", "user.name", "ada")
	run(t, repo, "git", "config", "user.email", "ada@example.test")
	run(t, repo, bin, "init")
	run(t, repo, bin, "share")

	id := strings.TrimSpace(runOut(t, repo, bin, "create", "session cookie dropped on 302",
		"--goal", "the cookie survives the OAuth round-trip",
		"--context", "reported while debugging Safari logouts; dropped cross-site on the redirect",
		"--dod", "session survives the round-trip",
		"--json"))
	ticketID := jsonField(t, id, "id")
	path := ticketPath(t, filepath.Join(repo, ".jaira", "tickets"), ticketID)

	run(t, repo, "git", "add", ".jaira", ".gitignore")
	run(t, repo, "git", "commit", "--quiet", "-m", "share the board")

	// Ada tags it on her branch and leaves it where it is. Her write is the
	// later one in wall-clock time, which is exactly what must not decide the
	// lane.
	run(t, repo, "git", "checkout", "--quiet", "-b", "ada")
	run(t, repo, bin, "tag", ticketID, "concurrency")
	run(t, repo, "git", "commit", "--quiet", "-am", "ada tags it")

	// Grace, on his own branch off the same commit, moves it forward and tags
	// it with something else.
	run(t, repo, "git", "checkout", "--quiet", "master")
	run(t, repo, "git", "checkout", "--quiet", "-b", "grace")
	run(t, repo, bin, "move", ticketID, "--to", "todo")
	run(t, repo, bin, "tag", ticketID, "cli")
	run(t, repo, "git", "commit", "--quiet", "-am", "grace moves it to todo")

	// The merge git would call a conflict on one line.
	run(t, repo, "git", "checkout", "--quiet", "ada")
	merge := exec.Command("git", "merge", "--no-edit", "grace")
	merge.Dir = repo
	if out, err := merge.CombinedOutput(); err != nil {
		t.Fatalf("the merge failed instead of resolving field by field: %v\n%s", err, out)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the merged ticket: %v", err)
	}
	got := string(b)
	if strings.Contains(got, "<<<<<<<") {
		t.Fatalf("the file has conflict markers, which would blank the card on every board:\n%s", got)
	}
	// Progress is never reverted: the further lane wins, whichever side's
	// timestamp is newer.
	if !strings.Contains(got, "status: todo") {
		t.Errorf("the merge did not keep the further lane:\n%s", got)
	}
	// Neither person's tag is silently dropped.
	for _, tag := range []string{"concurrency", "cli"} {
		if !strings.Contains(got, tag) {
			t.Errorf("tag %q was lost in the merge:\n%s", tag, got)
		}
	}
}

// repoRoot finds the module root from the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test's working directory")
		}
		dir = parent
	}
}

func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	runOut(t, dir, name, args...)
}

func runOut(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=ada", "GIT_AUTHOR_EMAIL=ada@example.test",
		"GIT_COMMITTER_NAME=ada", "GIT_COMMITTER_EMAIL=ada@example.test",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}

// jsonField reads one top-level string field, without a struct for a shape only
// this test cares about.
func jsonField(t *testing.T, blob, field string) string {
	t.Helper()
	needle := `"` + field + `": "`
	i := strings.Index(blob, needle)
	if i < 0 {
		t.Fatalf("no %q in %s", field, blob)
	}
	rest := blob[i+len(needle):]
	j := strings.Index(rest, `"`)
	if j < 0 {
		t.Fatalf("unterminated %q in %s", field, blob)
	}
	return rest[:j]
}

func ticketPath(t *testing.T, dir, id string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), id) {
			return filepath.Join(dir, e.Name())
		}
	}
	t.Fatalf("no ticket file for %s in %s", id, dir)
	return ""
}

// Two branches that each created the same ticket file — which is what git calls
// an add/add conflict, and what it hands the merge driver an empty base for.
// This is the case that used to fail silently: git left "ours" behind with no
// conflict markers, so it read as a clean merge while the other side's lane and
// updated-by were gone.
func TestTwoBranchesThatBothCreateTheTicketStillMergeFieldAware(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the binary")
	}
	for _, tool := range []string{"git", "go"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not on PATH", tool)
		}
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	bin := filepath.Join(root, "jaira"+exeSuffix())
	build := exec.Command("go", "build", "-o", bin, "github.com/BeMuCa/jaira/cmd/jaira")
	build.Dir = repoRoot(t)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building the binary: %v\n%s", err, out)
	}

	repo := filepath.Join(root, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	run(t, repo, "git", "init", "--quiet", "-b", "master", ".")
	run(t, repo, "git", "config", "user.name", "ada")
	run(t, repo, "git", "config", "user.email", "ada@example.test")
	run(t, repo, bin, "init")
	run(t, repo, bin, "share")
	run(t, repo, "git", "add", ".jaira", ".gitignore")
	run(t, repo, "git", "commit", "--quiet", "-m", "share the board")

	// Ada creates the ticket on her branch. Her write is the later one by the
	// clock, which is exactly what must not decide the lane.
	run(t, repo, "git", "checkout", "--quiet", "-b", "ada")
	out := runOut(t, repo, bin, "create", "cookie dropped on 302",
		"--goal", "the cookie survives the OAuth round-trip",
		"--context", "reported while debugging Safari logouts",
		"--dod", "survives the round-trip", "--json")
	id := jsonField(t, out, "id")
	run(t, repo, bin, "tag", id, "concurrency")
	run(t, repo, "git", "add", ".jaira")
	run(t, repo, "git", "commit", "--quiet", "-m", "ada creates the ticket")

	path := ticketPath(t, filepath.Join(repo, ".jaira", "tickets"), id)
	name := filepath.Base(path)

	// Grace lands the same ticket on his own branch, one lane further on. The
	// file exists on both branches with no common ancestor: add/add.
	run(t, repo, "git", "checkout", "--quiet", "master")
	run(t, repo, "git", "checkout", "--quiet", "-b", "grace")
	run(t, repo, "git", "checkout", "--quiet", "ada", "--", ".jaira/tickets/"+name)
	landed, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	theirs := strings.ReplaceAll(string(landed), "status: backlog", "status: todo")
	theirs = strings.ReplaceAll(theirs, "updated-by: ada", "updated-by: grace")
	theirs = strings.ReplaceAll(theirs, "  - concurrency", "  - cli")
	if err := os.WriteFile(path, []byte(theirs), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, repo, "git", "add", ".jaira")
	run(t, repo, "git", "commit", "--quiet", "-m", "grace lands the same ticket")

	run(t, repo, "git", "checkout", "--quiet", "ada")
	merge := exec.Command("git", "merge", "--no-edit", "grace")
	merge.Dir = repo
	mergeOut, mergeErr := merge.CombinedOutput()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the merged ticket: %v", err)
	}
	got := string(b)
	if strings.Contains(got, "<<<<<<<") {
		t.Fatalf("conflict markers would blank the card on every board:\n%s", got)
	}
	if !strings.Contains(got, "status: todo") {
		t.Errorf("the other side's lane was lost — the silent failure this test exists for.\ngit said: %s\nfile:\n%s", mergeOut, got)
	}
	for _, tag := range []string{"concurrency", "cli"} {
		if !strings.Contains(got, tag) {
			t.Errorf("tag %q was lost:\n%s", tag, got)
		}
	}
	if mergeErr != nil {
		// git still reports add/add as needing a commit; what matters is that
		// the file it left behind is the field-aware merge, not one side.
		t.Logf("git reported the add/add path as conflicted (expected): %s", mergeOut)
	}
}

// exeSuffix is what an executable has to be called on this platform.
func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
