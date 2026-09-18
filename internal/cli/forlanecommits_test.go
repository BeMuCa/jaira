package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/ticket"
)

// forLanePayload is the slice of the --for-lane JSON both tests below read: the
// diff they were handed and the three fields that say what it was built from.
type forLanePayload struct {
	Diff          string   `json:"diff"`
	Commits       []string `json:"commits"`
	CommitsSource string   `json:"commits_source"`
	Complete      bool     `json:"complete"`
}

// forLaneGitFixture stands up what both tests below need: a git repo with a
// store in it, and a ticket in the review lane whose every declared input but
// the diff is already filled — so the only thing that can be missing from the
// payload is the diff itself. It hands back a git runner that returns stdout,
// which gitRun (boardremote_test.go) does not, because rev-parse is read here.
func forLaneGitFixture(t *testing.T, at time.Time) (string, *ticket.Store, *ticket.Ticket, string, func(...string) string, func(string, string)) {
	t.Helper()
	if !gitrepo.Available() {
		t.Skip("git is not on PATH")
	}
	dir := t.TempDir()
	t.Setenv("JAIRA_USER", "berk")
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))

	run := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", "-q")
	run("config", "user.email", "fixture@example.com")
	run("config", "user.name", "Fixture")

	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:              ticket.NewID(at),
		ticket.FieldTitle:           "t",
		ticket.FieldStatus:          "review",
		ticket.FieldGoal:            "g",
		ticket.FieldDoD:             "d",
		ticket.FieldOutcomeWhat:     "w",
		ticket.FieldOutcomeResolves: "r",
	}, nil, "## Definition of Done\n\n- [x] d\n")
	if err != nil {
		t.Fatal(err)
	}

	write := func(name, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir, s, tk, ticket.Handle(tk.ID), run, write
}

// A ticket's commits: field is a snapshot, never the truth: it is written once,
// by 'jaira move --out --commits', and every commit made after that never joins
// it. A lane whose whole job is judging a diff must therefore not be handed the
// field as if it were the answer — it saw three of twenty-one commits and was
// told complete:true, which is the failure this test pins.
func TestForLaneDiffIsNotLimitedToTheRecordedCommits(t *testing.T) {
	dir, s, tk, h, run, write := forLaneGitFixture(t, time.Date(2026, 9, 18, 8, 0, 0, 0, time.UTC))
	write("first.txt", "erster\n")
	run("add", ".")
	run("commit", "-q", "-m", "feat("+h+"): the first change")
	first := run("rev-parse", "HEAD")

	write("second.txt", "zweiter\n")
	run("add", ".")
	run("commit", "-q", "-m", "feat("+h+"): the second change")

	// The snapshot: the ticket records only the first commit, exactly as a
	// 'move --out --commits' made before the second one would have left it.
	if _, err := s.Mutate(tk.ID, func(t *ticket.Ticket) error {
		return t.Doc().SetList(ticket.FieldCommits, []string{first})
	}); err != nil {
		t.Fatal(err)
	}

	var payload forLanePayload
	out, err := runCLI(t, dir, "show", h, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("payload is not json: %v\n%s", err, out)
	}
	if !strings.Contains(payload.Diff, "zweiter") {
		t.Errorf("the reviewed diff stops at the recorded commit — the second change is missing:\n%s", payload.Diff)
	}
	if len(payload.Commits) != 2 {
		t.Errorf("payload names %d commit(s), want the 2 the diff was built from: %v", len(payload.Commits), payload.Commits)
	}
	if payload.CommitsSource == "" {
		t.Error("payload does not say where the commit list came from")
	}
	if !payload.Complete {
		t.Error("complete is false although nothing is missing")
	}
}

// The same failure one step on: work that is not committed yet is work the lane
// is there to judge. A lane that changed no code commits nothing, and a
// conversational ticket commits nothing at all until a person pastes the line —
// so the payload that leaves the worktree out reports a fraction as complete.
func TestForLaneDiffCarriesTheUncommittedWorktree(t *testing.T) {
	dir, _, _, h, run, write := forLaneGitFixture(t, time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC))
	write("committed.txt", "eingecheckt\n")
	run("add", ".")
	run("commit", "-q", "-m", "feat("+h+"): the committed change")

	// Three shapes the payload has to carry, and each one used to vanish:
	// a tracked file edited, a tracked file edited and staged, and a brand-new
	// file no index has ever seen — the shape of a new test or a new package.
	write("committed.txt", "eingecheckt und geaendert\n")
	write("staged.txt", "vorgemerkt\n")
	run("add", "staged.txt")
	write("brandneu.txt", "unverfolgt\n")

	var payload forLanePayload
	out, err := runCLI(t, dir, "show", h, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("payload is not json: %v\n%s", err, out)
	}
	for _, want := range []string{"eingecheckt und geaendert", "vorgemerkt", "unverfolgt"} {
		if !strings.Contains(payload.Diff, want) {
			t.Errorf("the payload leaves uncommitted work out — %q is missing:\n%s", want, payload.Diff)
		}
	}
	if !strings.HasSuffix(payload.CommitsSource, "+"+ticket.SourceWorktree) {
		t.Errorf("commits_source is %q — it does not say the worktree is in the diff", payload.CommitsSource)
	}
	if !payload.Complete {
		t.Error("complete is false although nothing is missing")
	}
}
