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
	WorktreeDiff  string   `json:"worktree_diff"`
	Commits       []string `json:"commits"`
	CommitsSource string   `json:"commits_source"`
	Complete      bool     `json:"complete"`
	Unavailable   []string `json:"commits_unavailable"`
	WorktreeError string   `json:"worktree_error"`
	Missing       []string `json:"missing"`
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
//
// And it has to arrive in a key of its own. The two halves were one string
// once, split by the plain text "uncommitted work in the working tree", and a
// role prompt told its reader everything below that line was the work in
// progress. Nothing tells that line apart from a commit message quoted inside
// a patch: on the payload of the very ticket that fixed this it occurred nine
// times, the first some nine hundred lines above the real boundary. So the two
// halves are held apart here by content, and the old marker must be gone from
// both — which is what a single string cannot promise.
func TestForLaneCarriesTheUncommittedWorktreeInItsOwnKey(t *testing.T) {
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
		if !strings.Contains(payload.WorktreeDiff, want) {
			t.Errorf("the payload leaves uncommitted work out — %q is missing:\n%s", want, payload.WorktreeDiff)
		}
		if strings.Contains(payload.Diff, want) {
			t.Errorf("diff carries the uncommitted %q — the two halves are one string again:\n%s", want, payload.Diff)
		}
	}
	if !strings.Contains(payload.Diff, "eingecheckt\n") {
		t.Errorf("diff does not carry the committed half:\n%s", payload.Diff)
	}
	const gone = "uncommitted work in the working tree"
	if strings.Contains(payload.Diff, gone) || strings.Contains(payload.WorktreeDiff, gone) {
		t.Error("the old text marker is still in the payload — a reader can still mistake it for the boundary")
	}
	if !strings.HasSuffix(payload.CommitsSource, "+"+ticket.SourceWorktree) {
		t.Errorf("commits_source is %q — it does not say the worktree is in the diff", payload.CommitsSource)
	}
	if !payload.Complete {
		t.Error("complete is false although nothing is missing")
	}
}

// A payload built where git could not read the working tree used to be
// byte-for-byte the payload of a clean one: the error went into the bin at
// flow.go and commits_source said "git", which is what a spotless tree says
// too. That is the fraction reported as the whole all over again — the lane
// judges the commits and never learns that the uncommitted half was never
// looked at. The fixture reaches it through a repository with no commits at
// all, where "git diff HEAD" has no HEAD to diff against.
func TestForLaneSaysWhenTheWorktreeCouldNotBeRead(t *testing.T) {
	dir, s, tk, h, _, _ := forLaneGitFixture(t, time.Date(2026, 9, 18, 11, 0, 0, 0, time.UTC))

	// A sha the ticket records and this repository does not have, so the diff
	// is non-empty and the payload is complete: the point is that a complete
	// payload still admits what it could not read.
	const ghost = "0123456789012345678901234567890123456789"
	if _, err := s.Mutate(tk.ID, func(t *ticket.Ticket) error {
		return t.Doc().SetList(ticket.FieldCommits, []string{ghost})
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
	if payload.WorktreeError == "" {
		t.Errorf("the working tree could not be read and the payload does not say so:\n%s", out)
	}
	if !payload.Complete {
		t.Errorf("complete is false — an unreadable worktree informs the lane, it does not block it: missing=%v", payload.Missing)
	}
}

// repo.Diff never fails: a sha it cannot show becomes the line
// "(not available locally)" inside the patch. It is still counted in "commits",
// so anything that counts instead of reading — which is what "commits" is there
// for — reports more diffs than it was shown. The payload names them.
func TestForLaneNamesTheCommitsGitCouldNotShow(t *testing.T) {
	dir, s, tk, h, run, write := forLaneGitFixture(t, time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC))
	write("da.txt", "vorhanden\n")
	run("add", ".")
	run("commit", "-q", "-m", "feat("+h+"): the change that is here")
	here := run("rev-parse", "HEAD")

	// The rebase case the union exists to preserve: a sha only the field
	// carries, which this clone cannot resolve.
	const ghost = "0123456789012345678901234567890123456789"
	if _, err := s.Mutate(tk.ID, func(t *ticket.Ticket) error {
		return t.Doc().SetList(ticket.FieldCommits, []string{ghost})
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
	if len(payload.Commits) != 2 {
		t.Fatalf("payload names %d commit(s), want the recorded ghost beside %s: %v", len(payload.Commits), here[:7], payload.Commits)
	}
	if len(payload.Unavailable) != 1 || payload.Unavailable[0] != ghost {
		t.Errorf("commits_unavailable is %v — a sha counted in commits but never shown is not named", payload.Unavailable)
	}
}

// The plain-text lane prompt is what a worker who does not parse JSON reads,
// and there the worktree error is a sentence rather than a key. It used to be
// printed under the diff while saying "nothing uncommitted is below" — pointing
// at the empty space after itself, with the diff it qualifies above it — and,
// where there was no diff at all, a second time inside the missing line. Both
// halves of the fixture carry no commit that git can show, so the working tree
// is unreadable either way; only the recorded ghost sha differs.
func TestForLaneWorktreeErrorStandsAboveTheDiffAndOnlyOnce(t *testing.T) {
	const ghost = "0123456789012345678901234567890123456789"

	t.Run("beside a diff", func(t *testing.T) {
		dir, s, tk, h, _, _ := forLaneGitFixture(t, time.Date(2026, 9, 18, 13, 0, 0, 0, time.UTC))
		if _, err := s.Mutate(tk.ID, func(t *ticket.Ticket) error {
			return t.Doc().SetList(ticket.FieldCommits, []string{ghost})
		}); err != nil {
			t.Fatal(err)
		}
		out, err := runCLI(t, dir, "show", h, "--for-lane", "review")
		if err != nil {
			t.Fatalf("show --for-lane review: %v\n%s", err, out)
		}
		said := strings.Index(out, "The working tree could not be read")
		diffAt := strings.Index(out, "## Diff")
		if said < 0 || diffAt < 0 {
			t.Fatalf("want the worktree sentence and a diff, got:\n%s", out)
		}
		if said > diffAt {
			t.Errorf("the sentence says the uncommitted half is not below, but it stands under the diff:\n%s", out)
		}
	})

	t.Run("without a diff", func(t *testing.T) {
		dir, _, _, h, _, _ := forLaneGitFixture(t, time.Date(2026, 9, 18, 14, 0, 0, 0, time.UTC))
		out, err := runCLI(t, dir, "show", h, "--for-lane", "review")
		if err != nil {
			t.Fatalf("show --for-lane review: %v\n%s", err, out)
		}
		// Matched without its first word: the missing line spells it
		// "the working tree could not be read: …" mid-sentence, the standalone
		// line capitalises it, and a reader is told the same thing twice either
		// way.
		if n := strings.Count(out, "working tree could not be read"); n != 1 {
			t.Errorf("the same error is reported %d times, want once — the missing line carries it:\n%s", n, out)
		}
	})
}

// The other half of the same promise: a key that is only there when it carries
// something. core/release/NOTES.md and the jaira-role-lane prompt both tell a
// reader to test worktree_diff for content and never for absence — which only
// holds if a clean working tree leaves the key OUT of the payload instead of
// setting it to "". An empty key beside every payload trains the reader to
// skip it, and a reader who skips it reads a fraction as the whole. The
// struct the other tests unmarshal into cannot see the difference, so the
// payload is read as a bare map here.
func TestForLaneLeavesTheWorktreeKeyOutOfACleanTree(t *testing.T) {
	dir, s, _, h, run, write := forLaneGitFixture(t, time.Date(2026, 9, 18, 13, 0, 0, 0, time.UTC))
	write("committed.txt", "eingecheckt\n")
	run("add", ".")
	run("commit", "-q", "-m", "feat("+h+"): the committed change")

	// The CLI materialises the board's lane files on first use, so the tree is
	// only clean after it has run once: committing before that leaves them
	// behind as untracked work and the payload legitimately carries them.
	if out, err := runCLI(t, dir, "show", h, "--for-lane", "review", "--json"); err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	run("add", ".")
	run("commit", "-q", "-m", "chore("+h+"): the board the cli wrote")
	if left := run("status", "--porcelain"); left != "" {
		t.Fatalf("the fixture's tree is not clean, so the payload would carry a worktree half:\n%s", left)
	}

	out, err := runCLI(t, dir, "show", h, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("payload is not json: %v\n%s", err, out)
	}
	if v, ok := payload["worktree_diff"]; ok {
		t.Errorf("a clean tree still carries a worktree_diff key (%q) — absence is what the prompt tells its reader to rely on", v)
	}
	if v, ok := payload["worktree_error"]; ok {
		t.Errorf("a clean tree reports a worktree error: %v", v)
	}
	if got, _ := payload["commits_source"].(string); got != "git" {
		t.Errorf("commits_source is %q, want %q — nothing uncommitted was there to name", got, "git")
	}
	if got, _ := payload["complete"].(bool); !got {
		t.Error("complete is false although nothing is missing")
	}

	// The other half of the same asymmetry, and the half nothing pinned: diff
	// is ALWAYS there, empty when there is nothing, because the role prompt
	// tells its reader to test diff for content and worktree_diff for the
	// absence of the key. Read on a ticket with no commits and a clean tree,
	// where diff is the empty string — a payload that dropped the key here
	// would leave a lane unable to tell "nothing to judge" from "this build
	// does not carry a diff at all".
	empty, err := s.Create(map[string]string{
		ticket.FieldID:              ticket.NewID(time.Date(2026, 9, 18, 14, 0, 0, 0, time.UTC)),
		ticket.FieldTitle:           "ohne commits",
		ticket.FieldStatus:          "review",
		ticket.FieldGoal:            "g",
		ticket.FieldDoD:             "d",
		ticket.FieldOutcomeWhat:     "w",
		ticket.FieldOutcomeResolves: "r",
	}, nil, "## Definition of Done\n\n- [x] d\n")
	if err != nil {
		t.Fatal(err)
	}
	out, err = runCLI(t, dir, "show", ticket.Handle(empty.ID), "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	var none map[string]any
	if err := json.Unmarshal([]byte(out), &none); err != nil {
		t.Fatalf("payload is not json: %v\n%s", err, out)
	}
	v, ok := none["diff"]
	if !ok {
		t.Fatalf("the payload has no diff key although diff is promised to be always present, empty when there is nothing:\n%s", out)
	}
	if got, _ := v.(string); got != "" {
		t.Errorf("diff is %q on a ticket with no commits and a clean tree, want the empty string", got)
	}
}
