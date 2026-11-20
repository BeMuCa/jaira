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
	onBase(t, dir, settings.BoardPRPath(dir), value)
}

// onBase commits a board file with value into master, the base the gate
// reads it from, and returns to feat/x, where the file is then not on disk.
func onBase(t *testing.T, dir, path, value string) {
	t.Helper()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("checkout", "-q", "master")
	if err := os.WriteFile(path, []byte(value+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "-f", path)
	run("commit", "-q", "-m", "chore: board says "+value)
	run("checkout", "-q", "feat/x")
}

// The board files are read as the base holds them: a branch that writes or
// commits .jaira/pr saying agent, or pr-when saying signoff, does not open
// its own gate (JG3REW, given back from signoff, Alex 2026-11-20).
func TestPRGateReadsTheBoardFilesFromTheBase(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	done := ticketIn("done")
	commit("feat(" + done + "): accepted")
	path := settings.BoardPRPath(dir)
	if err := os.WriteFile(path, []byte("agent\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation || !strings.Contains(err.Error(), "in master does not say agent") {
		t.Fatalf("agent only on disk: exit %d, %v", exitCode(err), err)
	}
	if out, err := exec.Command("git", "-C", dir, "add", "-f", path).CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	commit("chore(" + done + "): the branch says agent itself")
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation {
		t.Fatalf("agent committed on the branch only: exit %d, %v", exitCode(err), err)
	}

	// With agent in the base, a signoff that only the branch says leaves the
	// gate at done.
	dir, ticketIn, commit = prGateFixture(t)
	setPR(t, dir, "agent")
	commit("feat(" + ticketIn("signoff") + "): waiting")
	setPRWhenOnDisk := settings.BoardPRWhenPath(dir)
	if err := os.WriteFile(setPRWhenOnDisk, []byte("signoff\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation || !strings.Contains(err.Error(), "in done") {
		t.Fatalf("signoff only on disk: exit %d, %v", exitCode(err), err)
	}
	setPRWhen(t, dir, "signoff")
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil {
		t.Fatalf("signoff in the base: %v\n%s", err, out)
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

	// A handle nothing here knows — a typo in a scope, a ticket on a ref
	// nobody fetched — was accepted by nobody: it holds the gate (review).
	commit("feat(ZZZZZZ): nobody knows it")
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation || !strings.Contains(err.Error(), "ZZZZZZ (unknown)") {
		t.Fatalf("unknown handle: exit %d, %v", exitCode(err), err)
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

func setPRWhen(t *testing.T, dir, value string) {
	t.Helper()
	onBase(t, dir, settings.BoardPRWhenPath(dir), value)
}

// Where the review happens on the pull request, it opens once every ticket
// is in signoff — or past it, in done; earlier lanes and the parking lane
// still hold it. Without .jaira/pr-when nothing short of done opens it
// (JG3REW, Alex 2026-11-20).
func TestPRGateOpensInSignoffWhereTheBoardSaysSo(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	setPR(t, dir, "agent")
	inSignoff, inDone := ticketIn("signoff"), ticketIn("done")
	commit("feat(" + inSignoff + "," + inDone + "): accepted and waiting")

	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation ||
		!strings.Contains(err.Error(), "in done") || !strings.Contains(err.Error(), inSignoff+" (signoff)") {
		t.Fatalf("pr-when done, a ticket in signoff: %v", err)
	}
	setPRWhen(t, dir, "signoff")
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil || !strings.Contains(out, "in signoff") {
		t.Fatalf("pr-when signoff, tickets in signoff and done: %v\n%s", err, out)
	}
	for _, status := range []string{"review", "human", "blocked", "in-progress"} {
		h := ticketIn(status)
		commit("fix(" + h + "): not there yet")
		_, err := runCLI(t, dir, "roles", "pr-gate", "master")
		if exitCode(err) != ExitValidation || !strings.Contains(err.Error(), h+" ("+status+")") {
			t.Errorf("pr-when signoff, a ticket in %s: exit %d, %v", status, exitCode(err), err)
		}
		if err := os.Remove(filepath.Join(dir, ".jaira", "tickets", fileOf(t, dir, h))); err != nil {
			t.Fatal(err)
		}
	}
}

// A board that has no lane signoff before its terminal lane cannot honour
// pr-when signoff: the gate says so rather than open on a guess.
func TestPRGateRefusesSignoffOnABoardWithoutIt(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	setPR(t, dir, "agent")
	setPRWhen(t, dir, "signoff")
	commit("feat(" + ticketIn("done") + "): done")
	if out, err := runCLI(t, dir, "lanes", "remove", "signoff"); err != nil {
		t.Fatalf("lanes remove signoff: %v\n%s", err, out)
	}
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation || !strings.Contains(err.Error(), "no lane signoff") {
		t.Fatalf("no signoff lane: exit %d, %v", exitCode(err), err)
	}
}

func fileOf(t *testing.T, dir, handle string) string {
	t.Helper()
	m, _ := filepath.Glob(filepath.Join(dir, ".jaira", "tickets", "*"+handle+"-*.md"))
	if len(m) != 1 {
		t.Fatalf("ticket file of %s: %v", handle, m)
	}
	return filepath.Base(m[0])
}

// A ticket whose file cannot be read — a merge conflict left in it, a broken
// frontmatter — drops out of the board's list. The gate must not read that
// as filed away and open: it refuses (JG3REW, bughunt).
func TestPRGateFailsClosedOnAnUnreadableTicket(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	setPR(t, dir, "agent")
	h := ticketIn("brainstorm")
	commit("feat(" + h + "): half done")
	path := filepath.Join(dir, ".jaira", "tickets", fileOf(t, dir, h))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append([]byte("<<<<<<< HEAD\n"), raw...), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "--json", "roles", "pr-gate", "master")
	if exitCode(err) != ExitValidation || strings.Contains(out, "may_open") {
		t.Fatalf("unreadable ticket: exit %d, %v\n%s", exitCode(err), err, out)
	}
}

// Off the board, only the logbook means accepted: a ticket filed there from
// done holds nothing up, a rejected one in the archive does. A ticket file
// that cannot be read holds the gate only when it is one of the branch's.
func TestPRGateReadsFiledTicketsByWhereTheyAre(t *testing.T) {
	dir, ticketIn, commit := prGateFixture(t)
	setPR(t, dir, "agent")
	filed, rejected := ticketIn("done"), ticketIn("backlog")
	commit("feat(" + filed + "): accepted and filed")
	if out, err := runCLI(t, dir, "logbook", filed); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil {
		t.Fatalf("a ticket in the logbook held it up: %v\n%s", err, out)
	}

	// Another ticket of the board, broken, is not this branch's business.
	other := filepath.Join(dir, ".jaira", "tickets", fileOf(t, dir, ticketIn("todo")))
	if err := os.WriteFile(other, []byte("<<<<<<< HEAD\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "roles", "pr-gate", "master"); err != nil {
		t.Fatalf("an unrelated unreadable ticket held it up: %v\n%s", err, out)
	}
	if err := os.Remove(other); err != nil {
		t.Fatal(err)
	}

	commit("fix(" + rejected + "): will not be done")
	if out, err := runCLI(t, dir, "reject", rejected, "--why", "not wanted", "--force"); err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	if _, err := runCLI(t, dir, "roles", "pr-gate", "master"); exitCode(err) != ExitValidation || !strings.Contains(err.Error(), rejected+" (archive)") {
		t.Fatalf("a rejected ticket: exit %d, %v", exitCode(err), err)
	}
}
