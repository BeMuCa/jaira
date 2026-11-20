package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// prGateFixture is a repository with a board, a base branch and a feature
// branch checked out ahead of it. ticketIn files a ticket in a lane and
// returns its handle; commit adds an empty commit with the given subject.
func prGateFixture(t *testing.T) (dir string, ticketIn func(status string) string, commit func(subject string)) {
	t.Helper()
	if !gitrepo.Available() {
		t.Skip("git is not on PATH")
	}
	dir = t.TempDir()
	t.Setenv("JAIRA_USER", "Alexander Sacharov")
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q", "-b", "master")
	run("config", "user.email", "fixture@example.com")
	run("config", "user.name", "Fixture")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	run("commit", "-q", "--allow-empty", "-m", "chore: base")
	run("checkout", "-q", "-b", "feat/x")
	at := time.Now()
	ticketIn = func(status string) string {
		t.Helper()
		at = at.Add(time.Millisecond)
		tk, err := s.Create(map[string]string{
			ticket.FieldID: ticket.NewID(at), ticket.FieldTitle: "t", ticket.FieldStatus: status,
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		return ticket.Handle(tk.ID)
	}
	commit = func(subject string) { run("commit", "-q", "--allow-empty", "-m", subject) }
	return dir, ticketIn, commit
}

func setPR(t *testing.T, dir, value string) {
	t.Helper()
	if err := os.WriteFile(settings.BoardPRPath(dir), []byte(value+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The gate the role used to keep in its text alone: an agent opens the pull
// request only on an agent board, and only once every ticket the branch's
// commits name is in done or already filed away.
func TestPRGateOpensOnlyWhenEveryTicketOfTheBranchIsDone(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	done, working := ticketIn("done"), ticketIn("in-progress")
	commit("feat(" + done + "): the accepted one")

	// The default board: a person opens it, whatever the tickets say.
	out, err := runCLI(t, dir, "roles", "pr-gate", "master")
	if exitCode(err) != ExitValidation || !strings.Contains(err.Error(), ".jaira/pr") {
		t.Fatalf("person board: exit %d, %v\n%s", exitCode(err), err, out)
	}

	setPR(t, dir, "agent")
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil || !strings.Contains(out, done) {
		t.Fatalf("every ticket done: %v\n%s", err, out)
	}

	// A second ticket on the branch, still being worked, holds it — and is
	// named, so the agent can say which one.
	commit("fix(" + working + "," + done + "): half of the work")
	_, err = runCLI(t, dir, "roles", "pr-gate", "master")
	if exitCode(err) != ExitValidation || !strings.Contains(err.Error(), working+" (in-progress)") ||
		strings.Contains(err.Error(), done+" (") {
		t.Fatalf("pending ticket: exit %d, %v", exitCode(err), err)
	}
}

func TestPRGateRefusesWhatItCannotJudge(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	setPR(t, dir, "agent")

	// Nothing on the branch names a ticket: nobody accepted anything.
	commit("chore: no ticket")
	_, err := runCLI(t, dir, "roles", "pr-gate", "master")
	if exitCode(err) != ExitValidation || !strings.Contains(err.Error(), "names a ticket") {
		t.Fatalf("no ticket: exit %d, %v", exitCode(err), err)
	}

	// A handle only in the body — a merge quoting what it brought in — is not
	// a ticket of this branch.
	commit("chore: merge\n\nfeat(" + ticketIn("review") + "): quoted")
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation {
		t.Fatalf("handle only in a body counted: %v", err)
	}

	// A handle the board does not hold is filed away, past done.
	commit("feat(ZZZZZZ): filed long ago")
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil {
		t.Fatalf("filed-away ticket held it up: %v\n%s", err, out)
	}

	for _, base := range []string{"nosuchbranch", "--all", ""} {
		if _, err := runCLI(t, dir, "roles", "pr-gate", base); exitCode(err) != ExitUsage {
			t.Errorf("base %q: exit %d, want usage", base, exitCode(err))
		}
	}
}

func TestHandlesInSubjects(t *testing.T) {
	got := gitrepo.HandlesInSubjects("feat(A3K9QP): x\nfix(B7XQ2M, A3K9QP)!: y\nchore: z\nfeat(lowerc): w\nfeat(TOOLONG1): v\n  docs(C4D5E6): indented")
	want := []string{"A3K9QP", "B7XQ2M", "C4D5E6"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("handles = %v, want %v", got, want)
	}
}
