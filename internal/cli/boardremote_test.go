package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

// cloneWithRemotes builds a real clone whose first remote is origin and which
// carries whatever else is asked for, plus a settings.json of its own.
func cloneWithRemotes(t *testing.T, machineRemote string, extra ...string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_HOME", home)
	if machineRemote != "" {
		settings := `{"remote": "` + machineRemote + `"}` + "\n"
		if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(settings), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	bare := filepath.Join(root, "board.git")
	gitRun(t, root, "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "clone")
	gitRun(t, root, "clone", "--quiet", bare, clone)
	gitRun(t, clone, "config", "user.name", "ada")
	gitRun(t, clone, "config", "user.email", "ada@example.test")
	for _, name := range extra {
		gitRun(t, clone, "remote", "add", name, bare)
	}
	return clone
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// ticketIn initialises a board in the clone and files one ticket in it,
// answering with its handle.
//
// Both steps go through the CLI, because a created ticket lives on its git ref
// and not on disk until somebody pulls it — which is exactly the mechanism
// these tests are about. The handle is then read back off that ref rather than
// picked out of what create printed: a six-character uppercase word turns up in
// titles and lane names too, so scanning stdout for one finds the wrong word
// sooner or later.
func ticketIn(t *testing.T, clone, title string) string {
	t.Helper()
	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "create", title); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ids, err := (&gitref.Repo{Dir: clone}).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		// No ref means the remote refused it, and create kept the ticket on
		// disk instead — which is the very case the third test sets up.
		st, err := ticket.At(clone)
		if err != nil {
			t.Fatal(err)
		}
		all, err := st.List()
		if err != nil {
			t.Fatal(err)
		}
		for _, tk := range all {
			ids = append(ids, tk.ID)
		}
	}
	if len(ids) != 1 {
		t.Fatalf("expected one ticket on the board, found %d: %v", len(ids), ids)
	}
	return ticket.Handle(ids[0])
}

// The whole bug in one test: a repository whose only remote is origin, on a
// machine whose settings.json says "upstream" because one other checkout needed
// it. Every ref command used to stop here with no remote "upstream", which took
// away the half of the mechanism a person needs most — handing a ticket back.
func TestRefCommandsWorkWhenTheMachineSettingNamesAnAbsentRemote(t *testing.T) {
	clone := cloneWithRemotes(t, "upstream")
	t.Setenv("JAIRA_USER", "ada")

	h := ticketIn(t, clone, "hand me back")
	if out, err := runCLI(t, clone, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "release", h); err != nil {
		t.Fatalf("release refused on a single-remote repository: %v\n%s", err, out)
	}
}

// And the counter-check that keeps the fix from being "always fall back to
// origin": where the configured remote exists, that is the one resolved, not the
// first remote in the list. This checks the resolution and nothing more — both
// remotes here point at the same bare repository, so which one the push went
// through is not observable from the outside. That the ref actually travels is
// what TestRefCommandsWorkWhenTheMachineSettingNamesAnAbsentRemote pulls and
// releases through, and TestABoardRemoteThatIsGoneStopsLoudly is what keeps a
// wrong name from being quietly swapped for a working one.
func TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt(t *testing.T) {
	clone := cloneWithRemotes(t, "upstream", "upstream")
	t.Setenv("JAIRA_USER", "ada")

	h := ticketIn(t, clone, "goes upstream")
	if out, err := runCLI(t, clone, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if refs == nil || refs.Repo.Remote != "upstream" {
		t.Fatalf("the ref was resolved to %q, not upstream", refs.Repo.Remote)
	}
}

// A name set for this clone is never second-guessed. With jaira.remote pointing
// at a remote that is not there, the command stops and the message says what to
// do about it — rather than quietly writing the ticket to origin.
func TestABoardRemoteThatIsGoneStopsLoudly(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	gitRun(t, clone, "config", "--local", "jaira.remote", "upstream")
	t.Setenv("JAIRA_USER", "ada")

	h := ticketIn(t, clone, "nowhere to go")
	out, err := runCLI(t, clone, "release", h)
	if err == nil {
		t.Fatalf("release went through with a remote that does not exist:\n%s", out)
	}
	msg := err.Error() + out
	for _, want := range []string{`"upstream"`, "origin", "config jaira.remote"} {
		if !strings.Contains(msg, want) {
			t.Errorf("the refusal does not mention %s:\n%s", want, msg)
		}
	}
}
