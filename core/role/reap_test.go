package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// reapRig is a repository with one ticket's worktree, opened the way a
// teamlead opens it: spawn.sh starts the worker, which starts its review in
// the same worktree. Herdr, gh and 'jaira show' are stubs; git is real.
type reapRig struct {
	bash, scripts, repo, wt string
	closed                  string // every tab the herdr stub closed, a line each
	merged                  string // what the gh stub answers for feat/ab12cd
	shown                   string // what 'jaira show AB12CD --json' answers; absent: not found
	labels                  string // the label of each pane's tab, a file per pane
	env                     []string
}

func newReapRig(t *testing.T) *reapRig {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	// Git Bash spells the worktree /c/..., Go and git.exe C:\...: the records
	// and the worktree list never compare equal there, and spawn.sh's Windows
	// path is covered by TestSpawnStartsClaudeWhereItsShellRuns.
	if runtime.GOOS == "windows" {
		t.Skip("reap.sh compares POSIX paths; Git Bash and git.exe spell them differently")
	}
	path := jairaPath(t)
	tmp, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	r := &reapRig{
		bash:    bash,
		scripts: filepath.Join(tmp, "skills"),
		repo:    filepath.Join(tmp, "src", "jaira"),
		closed:  filepath.Join(tmp, "closed"),
		merged:  filepath.Join(tmp, "merged.json"),
		shown:   filepath.Join(tmp, "shown.json"),
		labels:  filepath.Join(tmp, "labels"),
	}
	r.wt = filepath.Join(tmp, "src", ".worktrees", "jaira-ab12cd")
	if _, err := install(r.scripts, false, true); err != nil {
		t.Fatal(err)
	}
	r.scripts = filepath.Join(r.scripts, "jaira-dispatcher", "scripts")

	bin := filepath.Join(tmp, "bin")
	stubs := map[string]string{
		// A pane per tab create, numbered; its tab is t<pane>, and the tab's
		// label is kept in labels/<pane>.
		"herdr": `#!/usr/bin/env bash
n="` + tmp + `/panes-made"
l="` + r.labels + `"
case "$1 $2" in
  "tab create") echo x >> "$n"; p="p$(wc -l < "$n" | tr -d ' ')"
                while [ $# -gt 0 ]; do [ "$1" = --label ] && { mkdir -p "$l"; echo "$2" > "$l/$p"; }; shift; done
                echo "{\"result\":{\"root_pane\":{\"pane_id\":\"$p\"}}}" ;;
  "pane get")   echo "{\"result\":{\"pane\":{\"agent\":\"claude\",\"agent_status\":\"idle\",\"tab_id\":\"t$3\"}}}" ;;
  "tab get")    echo "{\"result\":{\"tab\":{\"tab_id\":\"$3\",\"label\":\"$(cat "$l/${3#t}")\"}}}" ;;
  "tab close")  echo "$3" >> "` + r.closed + `" ;;
esac
`,
		// origin is a fork: the pull requests are asked for on its parent.
		"gh": `#!/usr/bin/env bash
case "$*" in
  "repo view git@github.com:someone/jaira.git --json"*)
    echo '{"isFork":true,"nameWithOwner":"someone/jaira","parent":{"name":"jaira","owner":{"login":"upstream"}}}' ;;
  "pr list --head feat/ab12cd --state merged --limit 1 --json headRefOid --repo upstream/jaira")
    cat "` + r.merged + `" 2>/dev/null || echo '[]' ;;
  *) echo "unexpected gh $*" >&2; exit 1 ;;
esac
`,
		// Everything but 'show' goes to this tree's jaira: the scripts read
		// their JSON with it.
		"jaira": `#!/usr/bin/env bash
if [ "$1" = -C ] && [ "$3" = show ]; then
  cat "` + r.shown + `" 2>/dev/null || exit 5
  exit 0
fi
exec "` + filepath.Join(strings.SplitN(path, string(os.PathListSeparator), 2)[0], "jaira") + `" "$@"
`,
	}
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range stubs {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	r.env = append(os.Environ(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
		"HERDR_WORKSPACE_ID=", "HERDR_PANE_ID=", "JAIRA_BRANCH_PREFIX=", "TMPDIR="+tmp,
		"PATH="+bin+string(os.PathListSeparator)+path)

	if err := os.MkdirAll(r.repo, 0o755); err != nil {
		t.Fatal(err)
	}
	r.git(t, r.repo, "init", "-q", "-b", "master")
	r.git(t, r.repo, "config", "user.email", "t@example.com")
	r.git(t, r.repo, "config", "user.name", "t")
	r.git(t, r.repo, "remote", "add", "origin", "git@github.com:someone/jaira.git")
	r.write(t, filepath.Join(r.repo, "README"), "one\n")
	r.git(t, r.repo, "add", "README")
	r.git(t, r.repo, "commit", "-q", "-m", "one")

	r.run(t, r.repo, "spawn.sh", "ab12cd", "AB12CD", "worker", r.repo)
	r.run(t, r.wt, "spawn.sh", "--no-worktree", "ab12cd", "ab12cd", "review", r.wt)
	return r
}

func (r *reapRig) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r *reapRig) write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (r *reapRig) run(t *testing.T, dir, script string, args ...string) string {
	t.Helper()
	cmd := exec.Command(r.bash, append([]string{filepath.Join(r.scripts, script)}, args...)...)
	cmd.Dir = dir
	cmd.Env = r.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", script, args, err, out)
	}
	return string(out)
}

// mergedAt makes the gh stub report a merged pull request whose head is sha.
func (r *reapRig) mergedAt(t *testing.T, sha string) {
	r.write(t, r.merged, `[{"headRefOid":"`+sha+`"}]`)
}

func (r *reapRig) records(t *testing.T) []string {
	t.Helper()
	got, _ := filepath.Glob(filepath.Join(r.repo, ".git", "jaira", "panes", "*"))
	return got
}

func (r *reapRig) closedTabs() string {
	b, _ := os.ReadFile(r.closed)
	return strings.TrimSpace(string(b))
}

// spawn.sh leaves a record of every tab it opens — the pane, the ticket, the
// lane, the worktree — in the repository's git directory, where every
// worktree of it finds the same one.
func TestSpawnRecordsItsTabForReap(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	recs := r.records(t)
	if len(recs) != 2 {
		t.Fatalf("want a record per tab, got %v", recs)
	}
	for i, want := range []string{
		"pane=p1\nticket=AB12CD\nlane=worker\nworktree=" + r.wt + "\n",
		"pane=p2\nticket=ab12cd\nlane=review\nworktree=" + r.wt + "\n",
	} {
		b, err := os.ReadFile(recs[i])
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != want {
			t.Errorf("record %s:\n%s\nwant:\n%s", recs[i], b, want)
		}
	}
}

// Merged and done: every tab of the ticket closes, its worktree and branch go,
// and the records with them. --dry-run says the same and changes nothing.
func TestReapClosesAMergedTicket(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	r.write(t, r.shown, `{"status":"done"}`)
	// Untracked board files are no work: AYBDJ2's rule.
	r.write(t, filepath.Join(r.wt, ".jaira", "milestones", "m.md"), "x\n")

	out := r.run(t, r.repo, "reap.sh", "--dry-run")
	for _, want := range []string{"would close AB12CD/worker p1", "would close AB12CD/review p2",
		"would remove " + r.wt, "would delete branch feat/ab12cd"} {
		if !strings.Contains(out, want) {
			t.Errorf("--dry-run lacks %q:\n%s", want, out)
		}
	}
	if r.closedTabs() != "" || len(r.records(t)) != 2 {
		t.Fatalf("--dry-run changed something: closed %q, records %v", r.closedTabs(), r.records(t))
	}
	if _, err := os.Stat(r.wt); err != nil {
		t.Fatalf("--dry-run removed the worktree: %v", err)
	}

	out = r.run(t, r.repo, "reap.sh")
	if got := r.closedTabs(); got != "tp1\ntp2" {
		t.Errorf("closed tabs %q, want tp1 and tp2\n%s", got, out)
	}
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("worktree still there (%v):\n%s", err, out)
	}
	if b := r.git(t, r.repo, "branch", "--list", "feat/ab12cd"); b != "" {
		t.Errorf("branch still there: %q\n%s", b, out)
	}
	if recs := r.records(t); len(recs) != 0 {
		t.Errorf("records left: %v", recs)
	}
	if _, err := os.Stat(r.repo); err != nil {
		t.Errorf("the main worktree went: %v", err)
	}
}

// Anything that may still be work keeps the whole ticket — tabs, worktree,
// branch — and says why.
func TestReapKeepsWhatIsNotFinished(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		setup func(t *testing.T, r *reapRig)
		why   string
	}{
		"not merged": {
			func(t *testing.T, r *reapRig) {}, "waiting AB12CD: feat/ab12cd has no merged pull request",
		},
		"uncommitted": {
			func(t *testing.T, r *reapRig) {
				r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
				r.write(t, filepath.Join(r.wt, "README"), "changed\n")
			}, "kept AB12CD: " + r0 + " has uncommitted changes: README",
		},
		"committed after the merge": {
			func(t *testing.T, r *reapRig) {
				r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
				r.write(t, filepath.Join(r.wt, "README"), "later\n")
				r.git(t, r.wt, "commit", "-q", "-am", "later")
			}, "kept AB12CD: feat/ab12cd is not at the merged head",
		},
		"given back": {
			func(t *testing.T, r *reapRig) {
				r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
				r.write(t, r.shown, `{"status":"in-progress"}`)
			}, "kept AB12CD: the ticket is in in-progress, not done",
		},
		// Outside Herdr, or with a stale HERDR_BIN_PATH, an open tab looks
		// closed: everything stays until herdr answers.
		"herdr not reachable": {
			func(t *testing.T, r *reapRig) {
				r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
				r.write(t, r.shown, `{"status":"done"}`)
				r.env = append(r.env, "HERDR_BIN_PATH="+filepath.Join(r.repo, "no-herdr"))
			}, "kept AB12CD: herdr not reachable",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := newReapRig(t)
			c.setup(t, r)
			out := r.run(t, r.repo, "reap.sh")
			want := strings.ReplaceAll(c.why, r0, r.wt)
			if !strings.Contains(out, want) {
				t.Errorf("want %q, got:\n%s", want, out)
			}
			if r.closedTabs() != "" || len(r.records(t)) != 2 {
				t.Errorf("a tab was closed: %q, records %v", r.closedTabs(), r.records(t))
			}
			if _, err := os.Stat(r.wt); err != nil {
				t.Errorf("the worktree went: %v", err)
			}
		})
	}
}

// herdr gives a closed tab's ids out again after a restart: a record whose
// pane now sits in a tab labelled for something else closes nothing, and
// goes, since the tab it named is gone.
func TestReapLeavesATabItDidNotOpen(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	r.write(t, r.shown, `{"status":"done"}`)
	r.write(t, filepath.Join(r.labels, "p2"), "ZZ99ZZ/worker\n")

	if out := r.run(t, r.repo, "reap.sh", "--dry-run"); strings.Contains(out, "p2") {
		t.Errorf("--dry-run promises to close a tab the run leaves open:\n%s", out)
	}
	out := r.run(t, r.repo, "reap.sh")
	if got := r.closedTabs(); got != "tp1" {
		t.Errorf("closed tabs %q, want only tp1\n%s", got, out)
	}
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("worktree still there (%v):\n%s", err, out)
	}
	if recs := r.records(t); len(recs) != 0 {
		t.Errorf("records left: %v", recs)
	}
}

// A worktree git will not remove keeps its records, its branch with them, and
// the next round finishes what this one could not.
func TestReapComesBackToAWorktreeItCouldNotRemove(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	r.write(t, r.shown, `{"status":"done"}`)
	// A locked worktree resists a single --force.
	r.git(t, r.repo, "worktree", "lock", r.wt)

	out := r.run(t, r.repo, "reap.sh")
	if !strings.Contains(out, "kept AB12CD: removing "+r.wt+" failed") {
		t.Errorf("no reason given:\n%s", out)
	}
	if len(r.records(t)) != 2 {
		t.Errorf("records went with the worktree still there: %v", r.records(t))
	}
	if b := r.git(t, r.repo, "branch", "--list", "feat/ab12cd"); b == "" {
		t.Errorf("branch went with its worktree still there:\n%s", out)
	}

	r.git(t, r.repo, "worktree", "unlock", r.wt)
	out = r.run(t, r.repo, "reap.sh")
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("the next round left the worktree (%v):\n%s", err, out)
	}
	if b := r.git(t, r.repo, "branch", "--list", "feat/ab12cd"); b != "" {
		t.Errorf("the next round left the branch: %q\n%s", b, out)
	}
	if recs := r.records(t); len(recs) != 0 {
		t.Errorf("records left: %v", recs)
	}
}

// r0 stands for the rig's worktree in a reason written before the rig exists.
const r0 = "<worktree>"

// The teamlead runs reap.sh every round; a rename of the script has to fail
// here, not in a round where nothing gets closed and nobody notices.
func TestTeamleadReapsEveryRound(t *testing.T) {
	b, err := File("jaira-teamlead", skillFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "jaira-dispatcher/scripts/reap.sh") {
		t.Error("jaira-teamlead does not name jaira-dispatcher/scripts/reap.sh")
	}
	if _, err := File("jaira-dispatcher", "scripts/reap.sh"); err != nil {
		t.Errorf("jaira-dispatcher/scripts/reap.sh does not ship: %v", err)
	}
}

// Run from inside the ticket's own worktree — a subdirectory of it, even —
// reap.sh closes the tabs but leaves the worktree it runs in standing.
func TestReapNeverRemovesTheWorktreeItRunsIn(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	if out := r.run(t, r.wt, "reap.sh", "--dry-run"); strings.Contains(out, "would remove") || strings.Contains(out, "would delete") {
		t.Errorf("--dry-run promises what the run will not do:\n%s", out)
	}
	sub := filepath.Join(r.wt, "sub")
	r.write(t, filepath.Join(sub, ".keep"), "")
	r.git(t, r.wt, "add", "sub")
	r.git(t, r.wt, "commit", "-q", "-m", "sub")
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	out := r.run(t, sub, "reap.sh")
	if r.closedTabs() != "tp1\ntp2" {
		t.Errorf("closed tabs %q, want tp1 and tp2\n%s", r.closedTabs(), out)
	}
	if _, err := os.Stat(r.wt); err != nil || strings.Contains(out, "removed") {
		t.Errorf("removed the worktree it runs in (%v):\n%s", err, out)
	}
}
