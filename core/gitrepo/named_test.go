package gitrepo

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestHandlesInSubjectsReadsTheScopeAndGitsRevertSubjects(t *testing.T) {
	for subjects, want := range map[string][]string{
		"fix(A3K9QP): x":                                  {"A3K9QP"},
		"chore(A3K9QP,B7XQ2M)!: x":                        {"A3K9QP", "B7XQ2M"},
		"chore(A3K9QP, B7XQ2M): x\nfix(A3K9QP): y":        {"A3K9QP", "B7XQ2M"},
		`Revert "fix(A3K9QP): x"`:                         {"A3K9QP"},
		`Revert "Revert "fix(A3K9QP): x""`:                {"A3K9QP"},
		`Reapply "fix(A3K9QP): x"`:                        {"A3K9QP"},
		"fix: no scope":                                   nil,
		"Fix(A3K9QP): capital type":                       nil,
		"fix(a3k9qp): lower-case handle":                  nil,
		"fix(A3K9QPX): seven":                             nil,
		`revert "fix(A3K9QP): x"`:                         nil,
		`Revert fix(A3K9QP): x`:                           nil,
		"Merge branch 'feat/a3k9qp'\n\nfix(B7XQ2M): body": {"B7XQ2M"},
	} {
		if got := HandlesInSubjects(subjects); !reflect.DeepEqual(got, want) {
			t.Errorf("HandlesInSubjects(%q) = %v, want %v", subjects, got, want)
		}
	}
}

// namedRig: a repository with base on main and a branch feat off it, where
// each test adds what it needs on top.
func namedRig(t *testing.T) (dir string, repo *Repo, git func(args ...string) string, at func() time.Time) {
	t.Helper()
	dir, repo = fixtureRepo(t)
	clock := time.Date(2026, 11, 21, 9, 0, 0, 0, time.UTC)
	at = func() time.Time { clock = clock.Add(time.Minute); return clock }
	git = func(args ...string) string {
		t.Helper()
		stamp := at().Format(time.RFC3339)
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(cmd.Environ(), "GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("checkout", "-q", "-b", "main")
	fixtureCommit(t, dir, map[string]string{"a": "1\n", "b": "1\n"}, "chore(BASE00): base", at())
	git("checkout", "-q", "-b", "feat")
	return dir, repo, git, at
}

func TestUnnamedCommitsCountsOnlyWhatNoTicketOwns(t *testing.T) {
	dir, repo, git, at := namedRig(t)
	fixtureCommit(t, dir, map[string]string{"a": "2\n"}, "fix(A3K9QP): the work", at())
	fixtureCommit(t, dir, map[string]string{"b": "2\n"}, "tidy up", at())
	git("revert", "--no-edit", "HEAD~1") // Revert "fix(A3K9QP): the work"
	got, err := repo.UnnamedCommits("main")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !strings.HasSuffix(got[0], " tidy up") {
		t.Errorf("UnnamedCommits = %q, want only the 'tidy up' commit (a revert names what it reverts)", got)
	}
	if named, _ := repo.TicketsNamed("main"); !reflect.DeepEqual(named, []string{"A3K9QP"}) {
		t.Errorf("TicketsNamed = %v, want [A3K9QP]", named)
	}
	for _, base := range []string{"", "-p"} {
		if _, err := repo.UnnamedCommits(base); err != ErrNoSuchRef {
			t.Errorf("UnnamedCommits(%q) err = %v, want ErrNoSuchRef", base, err)
		}
	}
}

func TestIsPlainMergeTakesOnlyGitsOwnMerge(t *testing.T) {
	if out, err := exec.Command("git", "merge-tree", "-h").CombinedOutput(); err == nil || !strings.Contains(string(out), "--write-tree") {
		t.Skip("git before 2.38: no merge-tree --write-tree, the gate fails closed there")
	}
	dir, repo, git, at := namedRig(t)
	fixtureCommit(t, dir, map[string]string{"a": "feat\n"}, "fix(A3K9QP): feat side", at())
	git("checkout", "-q", "main")
	fixtureCommit(t, dir, map[string]string{"b": "main\n"}, "fix(B7XQ2M): main side", at())
	mainTip := git("rev-parse", "HEAD")
	git("checkout", "-q", "feat")

	merge := func(args ...string) (sha string, parents []string) {
		t.Helper()
		git(args...)
		sha = git("rev-parse", "HEAD")
		return sha, strings.Fields(git("log", "-1", "--format=%P", sha))
	}
	// git's own merge of the two: plain.
	plain, ps := merge("merge", "-q", "--no-edit", "-m", "Merge main", "main")
	if !repo.isPlainMerge(plain, ps) {
		t.Errorf("a clean 'git merge' did not read as plain")
	}
	if got, _ := repo.UnnamedCommits(mainTip); len(got) != 0 {
		t.Errorf("a clean merge held the gate: %q", got)
	}

	// -s ours: the result is one parent's tree, which diff-tree --cc hides.
	git("reset", "-q", "--hard", "HEAD~1")
	ours, ps := merge("merge", "-q", "-s", "ours", "-m", "Merge main, ours", "main")
	if repo.isPlainMerge(ours, ps) {
		t.Errorf("a '-s ours' merge read as plain")
	}
	if got, _ := repo.UnnamedCommits(mainTip); len(got) != 1 || !strings.HasSuffix(got[0], "Merge main, ours") {
		t.Errorf("UnnamedCommits = %q, want the -s ours merge", got)
	}

	// Content smuggled in during merge --no-commit.
	git("reset", "-q", "--hard", "HEAD~1")
	git("merge", "-q", "--no-commit", "main")
	fixtureCommit(t, dir, map[string]string{"c": "extra\n"}, "Merge main with extra", at())
	sha := git("rev-parse", "HEAD")
	if repo.isPlainMerge(sha, strings.Fields(git("log", "-1", "--format=%P", sha))) {
		t.Errorf("a merge with a file of its own read as plain")
	}

	// A conflict, settled by hand: git's own merge has no clean tree to
	// compare with, so it is never plain.
	git("reset", "-q", "--hard", "HEAD~1")
	fixtureCommit(t, dir, map[string]string{"b": "feat\n"}, "fix(A3K9QP): touch b", at())
	if cmd := exec.Command("git", "-C", dir, "merge", "-q", "main"); cmd.Run() == nil {
		t.Fatal("expected a conflict on b")
	}
	fixtureCommit(t, dir, map[string]string{"b": "settled\n"}, "Merge main, settled", at())
	sha = git("rev-parse", "HEAD")
	if repo.isPlainMerge(sha, strings.Fields(git("log", "-1", "--format=%P", sha))) {
		t.Errorf("a merge with a conflict settled by hand read as plain")
	}

	// Not two parents, or not a commit at all: never plain.
	if repo.isPlainMerge(plain, ps[:1]) || repo.isPlainMerge(plain, append(ps, ps[0])) {
		t.Errorf("a parent count other than two read as plain")
	}
	if repo.isPlainMerge("0000000000000000000000000000000000000000", ps) {
		t.Errorf("an unknown merge sha read as plain")
	}
}
