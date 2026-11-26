package role

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

// reapRig is a repository with one ticket's worktree, opened the way a
// teamlead opens it: spawn.sh starts the worker, which starts its review in
// the same worktree. Herdr, gh and 'jaira show' are stubs; git is real.
type reapRig struct {
	bash, scripts, repo, wt string
	closed                  string // every tab the herdr stub closed, a line each
	merged                  string // what the gh stub answers for feat/ab12cd
	shown                   string // what 'jaira show AB12CD --json' answers; absent: not found
	lanes                   string // what 'jaira -C <dir> lanes --json' answers
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
		lanes:   filepath.Join(tmp, "lanes.json"),
		labels:  filepath.Join(tmp, "labels"),
	}
	r.wt = filepath.Join(tmp, "src", ".worktrees", "jaira-ab12cd")
	if _, err := install(r.scripts, false, true, false); err != nil {
		t.Fatal(err)
	}
	r.scripts = filepath.Join(r.scripts, "jaira-dispatcher", "scripts")

	bin := filepath.Join(tmp, "bin")
	stubs := map[string]string{
		// A pane per tab create, numbered; its tab is t<pane>, and the tab's
		// label is kept in labels/<pane>, its agent's state, idle unless
		// said otherwise, in labels/<pane>.status.
		"herdr": `#!/usr/bin/env bash
n="` + tmp + `/panes-made"
l="` + r.labels + `"
# labels/<id>.fails, where it exists, answers 'pane get' or 'tab get' for that
# id instead: a script fragment that prints herdr's error and exits.
case "$1 $2" in "pane get"|"tab get") [ -f "$l/$3.fails" ] && . "$l/$3.fails" ;; esac
case "$1 $2" in
  "tab create") echo x >> "$n"; p="p$(wc -l < "$n" | tr -d ' ')"
                while [ $# -gt 0 ]; do [ "$1" = --label ] && { mkdir -p "$l"; echo "$2" > "$l/$p"; }; shift; done
                echo "{\"result\":{\"root_pane\":{\"pane_id\":\"$p\"}}}" ;;
  "pane get")   st=idle; [ -f "$l/$3.status" ] && st="$(cat "$l/$3.status")"
                echo "{\"result\":{\"pane\":{\"agent\":\"claude\",\"agent_status\":\"$st\",\"tab_id\":\"t$3\"}}}" ;;
  "tab get")    echo "{\"result\":{\"tab\":{\"tab_id\":\"$3\",\"label\":\"$(cat "$l/${3#t}")\"}}}" ;;
  "tab close")  echo "$3" >> "` + r.closed + `" ;;
esac
`,
		// origin is a fork: the pull requests are asked for on its parent.
		"gh": `#!/usr/bin/env bash
if [ -e "` + r.merged + `.gh-down" ]; then echo "error connecting to api.github.com" >&2; exit 1; fi
case "$*" in
  "repo view git@github.com:someone/jaira.git --json"*)
    if [ -e "` + r.merged + `.notfork" ]; then
      echo '{"isFork":false,"nameWithOwner":"renamed/jaira"}'
    else
      echo '{"isFork":true,"nameWithOwner":"someone/jaira","parent":{"name":"jaira","owner":{"login":"upstream"}}}'
    fi ;;
  "repo view git@github.com:formerly/jaira.git --json nameWithOwner")
    echo '{"nameWithOwner":"upstream/jaira"}' ;;
  "pr list --head feat/ab12cd --state merged --limit 1 --json headRefOid --repo upstream/jaira")
    cat "` + r.merged + `" 2>/dev/null || echo '[]' ;;
  *) echo "unexpected gh $*" >&2; exit 1 ;;
esac
`,
		// Everything but 'show' goes to this tree's jaira: the scripts read
		// their JSON with it - unless shown.json.json-fails exists, which
		// stands for a jaira whose 'json get' refuses (an older build).
		"jaira": `#!/usr/bin/env bash
if [ "$1" = json ] && [ -e "` + r.shown + `.json-fails" ]; then
  exit 7
fi
if [ "$1" = roles ] && [ -e "` + r.shown + `.roles-silent" ]; then
  echo "jaira: this board was set up by an older version of jaira" >&2
  exit 4
fi
if [ "$1" = roles ] && [ -e "` + r.shown + `.roles-fails" ]; then
  echo "jaira: this board was set up by an older version of jaira" >&2
  echo "jaira: roles list: cannot read .jaira/lanes" >&2
  exit 4
fi
if [ "$1" = -C ] && [ "$3" = show ]; then
  cat "` + r.shown + `" 2>/dev/null || exit 5
  exit 0
fi
if [ "$1" = -C ] && [ "$3" = lanes ]; then
  cat "` + r.lanes + `" 2>/dev/null || exit 5
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

// noneOnBase commits .jaira/pr none into master, where 'jaira roles list'
// reads it — the base decides, never the checkout's disk — and brings the
// ticket's branch up to it, so every origin ref a test points at master or
// past it holds the file.
func (r *reapRig) noneOnBase(t *testing.T) {
	t.Helper()
	r.write(t, filepath.Join(r.repo, ".jaira", "pr"), "none\n")
	r.git(t, r.repo, "add", "-f", ".jaira/pr")
	r.git(t, r.repo, "commit", "-q", "-m", "chore: no pull requests")
	r.git(t, r.wt, "merge", "-q", "--ff-only", "master")
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

// A real run keeps the pauses a terminal needs: 2 s between the looks at the
// tab, 1 s before the Enter. Only JAIRA_SPAWN_PAUSE, which TestMain sets to 0,
// shortens them; a sleep on PATH records what the script asked for.
func TestSpawnPausesAsLongAsBeforeWithoutTheVariable(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	bin := filepath.Join(filepath.Dir(r.closed), "bin")
	slept := filepath.Join(bin, "slept")
	if err := os.WriteFile(filepath.Join(bin, "sleep"), []byte("#!/usr/bin/env bash\necho \"$*\" >> "+slept+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	r.env = append(r.env, "JAIRA_SPAWN_PAUSE=")
	r.run(t, r.wt, "spawn.sh", "--no-worktree", "ab12cd", "ab12cd", "critique", r.wt)
	b, err := os.ReadFile(slept)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(b); got != "2\n1\n" {
		t.Errorf("spawn.sh slept %q, want 2 s between looks and 1 s before Enter", got)
	}
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
	// A German title, as ticket files have: git quotes a non-ASCII path
	// unless told not to, and the filter would miss it.
	ticket := filepath.Join(r.wt, ".jaira", "tickets", "ab12cd-größe.md")
	r.write(t, ticket, "in review\n")
	r.git(t, r.wt, "add", ".jaira/tickets")
	r.git(t, r.wt, "commit", "-q", "-m", "feat(AB12CD): ticket")
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	r.write(t, r.shown, `{"status":"done"}`)
	// Untracked milestone files are no work: AYBDJ2's rule; a ref holds them.
	r.write(t, filepath.Join(r.wt, ".jaira", "milestones", "m.md"), "x\n")
	// Neither is a ticket file a lane without code left modified: after the
	// merge no commit comes to carry it, and its state is on the ticket's ref.
	r.write(t, ticket, "done\n")

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

// Every person's cockpit and every teamlead run reap.sh: while one runs, a
// second skips the round instead of asking the forge again and racing it over
// 'git worktree remove'. The lock dies with the run that held it.
func TestReapRunsOneAtATime(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("perl"); err != nil {
		t.Skipf("no perl to hold the lock with: %v", err)
	}
	r := newReapRig(t)
	r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	r.write(t, r.shown, `{"status":"done"}`)
	lock := filepath.Join(r.repo, ".git", "jaira", "reap.lock")

	// Another run, as far as the lock can tell: a process holding it.
	holder := exec.Command("perl", "-MFcntl=:flock", "-e",
		`open(my $f, ">>", $ARGV[0]) or die; flock($f, LOCK_EX) or die; $| = 1; print "held\n"; sleep 60`, lock)
	stdout, err := holder.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = holder.Process.Kill(); _ = holder.Wait() })
	if line, err := bufio.NewReader(stdout).ReadString('\n'); err != nil || line != "held\n" {
		t.Fatalf("the holder did not take the lock: %q %v", line, err)
	}

	out := r.run(t, r.repo, "reap.sh")
	if !strings.Contains(out, "skipped: reap.sh runs already") {
		t.Errorf("a second run did not skip:\n%s", out)
	}
	if r.closedTabs() != "" || len(r.records(t)) != 2 {
		t.Fatalf("the skipped run changed something: closed %q, records %v", r.closedTabs(), r.records(t))
	}
	// --dry-run changes nothing and waits for nobody.
	if out := r.run(t, r.repo, "reap.sh", "--dry-run"); !strings.Contains(out, "would remove "+r.wt) {
		t.Errorf("--dry-run stopped at the lock:\n%s", out)
	}

	// The holder dies, as a run killed halfway would: the next run goes on.
	_ = holder.Process.Kill()
	_ = holder.Wait()
	out = r.run(t, r.repo, "reap.sh")
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("the lock outlived its holder (%v):\n%s", err, out)
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
		"detached HEAD": {
			func(t *testing.T, r *reapRig) {
				r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
				r.write(t, r.shown, `{"status":"done"}`)
				r.git(t, r.wt, "checkout", "-q", "--detach")
			}, "kept AB12CD: " + r0 + " has a detached HEAD",
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

// A merged ticket's tab counts as closed only when herdr itself says it does
// not know the pane or tab: hc's kind not_found on stderr, old herdr's code
// pane_not_found. Any other failure - a timeout, a herdr gone mid-round, no
// message at all - keeps the worktree, the branch and the records for the next
// round: claude may still sit in that tab.
func TestReapReadsOnlyNotFoundAsAClosedTab(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		id, fails string
		removes   bool
		says      string
	}{
		"pane get times out": {"p2", "echo 'request timed out' >&2; exit 1", false,
			"kept AB12CD: herdr gave no answer for pane p2 (request timed out)"},
		"pane get fails silently": {"p2", "exit 1", false,
			"kept AB12CD: herdr gave no answer for pane p2 (exit 1, no message)"},
		// Old herdr has more *_not_found codes; only the pane's says it is gone.
		"herdr without pane get": {"p2", `echo '{"error":{"code":"command_not_found"}}'; exit 1`, false,
			`kept AB12CD: herdr gave no answer for pane p2 ({"error":{"code":"command_not_found"}})`},
		"tab get times out": {"tp2", "echo 'request timed out' >&2; exit 1", false,
			"kept AB12CD: herdr gave no answer for tab tp2 (request timed out)"},
		"pane unknown to hc": {"p2", `echo '{"error":{"kind":"not_found"}}' >&2; exit 3`, true,
			"removed "},
		"pane unknown to the old herdr": {"p2", `echo '{"error":{"code":"pane_not_found"}}'; exit 1`, true,
			"removed "},
		"tab unknown to hc": {"tp2", `echo '{"error":{"kind":"not_found"}}' >&2; exit 3`, true,
			"removed "},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := newReapRig(t)
			r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
			r.write(t, r.shown, `{"status":"done"}`)
			r.write(t, filepath.Join(r.labels, c.id+".fails"), c.fails+"\n")

			out := r.run(t, r.repo, "reap.sh")
			if !strings.Contains(out, c.says) {
				t.Errorf("output lacks %q:\n%s", c.says, out)
			}
			_, err := os.Stat(r.wt)
			if c.removes {
				if !os.IsNotExist(err) {
					t.Errorf("worktree still there (%v):\n%s", err, out)
				}
				if got := r.closedTabs(); got != "tp1" {
					t.Errorf("closed tabs %q, want only tp1\n%s", got, out)
				}
				return
			}
			if err != nil {
				t.Errorf("worktree removed beside a tab herdr gave no answer about: %v\n%s", err, out)
			}
			if b := r.git(t, r.repo, "branch", "--list", "feat/ab12cd"); b == "" {
				t.Errorf("branch deleted:\n%s", out)
			}
			if recs := r.records(t); len(recs) != 2 {
				t.Errorf("records %v, want both kept\n%s", recs, out)
			}
		})
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

// A board whose .jaira/pr says none has no pull request to ask about: the
// ticket's worktree goes once origin's default branch contains its branch,
// and waits until then — it is never kept forever for a merge that cannot
// come (XB97NJ).
func TestReapOnANoneBoardWaitsForTheDefaultBranch(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.noneOnBase(t)
	r.write(t, filepath.Join(r.repo, ".jaira", "tickets", ".keep"), "")
	r.write(t, r.shown, `{"status":"done"}`)
	// origin is a fork (the rig's gh says so): work lands on its parent, and
	// with no remote here pointing at it, it says how to add one rather
	// than measuring the fork, which never receives the landing.
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "no remote here points at upstream/jaira") {
		t.Fatalf("fork without its parent as a remote, want the add hint:\n%s", out)
	}
	// A remote added before the parent was renamed keeps its old URL; gh
	// follows GitHub's redirect and says it is the parent all the same.
	r.git(t, r.repo, "remote", "add", "old", "git@github.com:formerly/jaira.git")
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "cannot tell old's default branch") {
		t.Fatalf("parent known only under its old name, want it found:\n%s", out)
	}
	r.git(t, r.repo, "remote", "remove", "old")
	r.git(t, r.repo, "remote", "add", "up", "https://github.com/Upstream/jaira.git")
	// No ref of it at all: nothing to measure against, and it says how to
	// fix that rather than claiming the branch has not landed.
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "cannot tell up's default branch") {
		t.Fatalf("no up refs, want the set-head hint:\n%s", out)
	}
	// up/HEAD left dangling (the remote renamed master to main and a
	// prune took up/master away) counts as unset, not as a base.
	r.git(t, r.repo, "symbolic-ref", "refs/remotes/up/HEAD", "refs/remotes/up/gone")
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "cannot tell up's default branch") {
		t.Fatalf("dangling up/HEAD, want the set-head hint:\n%s", out)
	}
	// up/master beside it (or without up/HEAD at all) is enough.
	base := r.git(t, r.repo, "rev-parse", "master")
	r.git(t, r.repo, "update-ref", "refs/remotes/up/master", base)
	// A valid up/HEAD beside a local branch that is also called
	// up/master: the remote one is still the base (bughunt round 6).
	r.git(t, r.repo, "symbolic-ref", "refs/remotes/up/HEAD", "refs/remotes/up/master")
	r.git(t, r.repo, "branch", "up/master", base)

	// A branch with no work of its own yet sits on the default branch from
	// the start: that is not landed, whatever the ticket's lane says.
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "waiting AB12CD") {
		t.Fatalf("no work on the branch yet, want waiting:\n%s", out)
	}
	if _, err := os.Stat(r.wt); err != nil {
		t.Fatalf("worktree removed before any work landed: %v", err)
	}
	// A body line naming the ticket is not a subject that does: --grep finds
	// it, the subject check does not take it.
	r.git(t, r.wt, "commit", "-q", "--allow-empty", "-m", "chore: tidy", "-m", "feat(AB12CD): only mentioned")
	r.git(t, r.repo, "update-ref", "refs/remotes/up/master", r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "waiting AB12CD: feat/ab12cd has not landed") {
		t.Fatalf("only a body line names the ticket, want waiting:\n%s", out)
	}
	r.git(t, r.wt, "commit", "-q", "--allow-empty", "-m", "feat(AB12CD): the work")

	out := r.run(t, r.repo, "reap.sh")
	if !strings.Contains(out, "waiting AB12CD: feat/ab12cd has not landed") {
		t.Fatalf("not landed yet, want waiting:\n%s", out)
	}
	if _, err := os.Stat(r.wt); err != nil {
		t.Fatalf("worktree removed before it landed: %v", err)
	}

	// The fork's own default branch carrying it is not landing: work lands
	// on the parent.
	tip := r.git(t, r.repo, "rev-parse", "feat/ab12cd")
	r.git(t, r.repo, "update-ref", "refs/remotes/origin/master", tip)
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "waiting AB12CD: feat/ab12cd has not landed on up's") {
		t.Fatalf("only the fork has it, want waiting:\n%s", out)
	}
	r.git(t, r.repo, "update-ref", "refs/remotes/up/master", tip)
	out = r.run(t, r.repo, "reap.sh")
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("landed and done, worktree still there (%v):\n%s", err, out)
	}
	if strings.Contains(out, "pull request") {
		t.Errorf("a none board asked about a pull request:\n%s", out)
	}
}

// Not a fork, but renamed since it was cloned: gh answers with the new name,
// origin's URL still has the old one, and origin is still where work lands —
// never 'no remote here points at it' (bughunt round 8).
func TestReapOnANoneBoardMeasuresARenamedOrigin(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.write(t, r.merged+".notfork", "")
	r.noneOnBase(t)
	r.write(t, filepath.Join(r.repo, ".jaira", "tickets", ".keep"), "")
	r.write(t, r.shown, `{"status":"done"}`)
	r.git(t, r.wt, "commit", "-q", "--allow-empty", "-m", "feat(AB12CD): the work")
	r.git(t, r.repo, "update-ref", "refs/remotes/origin/master", r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
	out := r.run(t, r.repo, "reap.sh")
	if strings.Contains(out, "no remote here points at") {
		t.Fatalf("renamed origin read as missing:\n%s", out)
	}
	if _, err := os.Stat(r.wt); !os.IsNotExist(err) {
		t.Errorf("landed on origin and done, worktree still there (%v):\n%s", err, out)
	}
}

// When 'jaira roles list' fails, how work arrives cannot be told: reap keeps
// every ticket and says so, rather than reading the board as person and
// waiting for a pull request a none board never opens.
func TestReapNamesAFailingRolesListAsTheCause(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.noneOnBase(t)
	r.write(t, filepath.Join(r.repo, ".jaira", "tickets", ".keep"), "")
	r.write(t, r.shown, `{"status":"done"}`)
	r.write(t, r.shown+".roles-fails", "")
	out := r.run(t, r.repo, "reap.sh")
	if !strings.Contains(out, "kept AB12CD: cannot tell how work arrives on this board — 'jaira roles list' failed: jaira: roles list: cannot read .jaira/lanes") {
		t.Fatalf("roles list failed, want it named:\n%s", out)
	}
	if strings.Contains(out, "pull request") {
		t.Errorf("fell back to the pull request path:\n%s", out)
	}
	if _, err := os.Stat(r.wt); err != nil {
		t.Fatalf("worktree removed although how work arrives is unknown: %v", err)
	}
	// A failure with nothing on stderr but the older-version notice: still
	// named, and the run goes on to the end under pipefail (bughunt round 10).
	r.write(t, r.shown+".roles-silent", "")
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "'jaira roles list' failed: no message (exit status only)") {
		t.Fatalf("roles list failed silently, want it named:\n%s", out)
	}
}

// When gh does not answer, whether origin is a fork cannot be told: reap
// measures at origin, and a ticket that waits there names gh as the cause
// rather than only saying it has not landed (review note, 21.11.).
func TestReapOnANoneBoardNamesAnUnansweringGh(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	r.write(t, r.merged+".gh-down", "")
	r.noneOnBase(t)
	r.write(t, filepath.Join(r.repo, ".jaira", "tickets", ".keep"), "")
	r.write(t, r.shown, `{"status":"done"}`)
	r.git(t, r.wt, "commit", "-q", "--allow-empty", "-m", "feat(AB12CD): the work")
	r.git(t, r.repo, "update-ref", "refs/remotes/origin/master", r.git(t, r.repo, "rev-parse", "master"))
	out := r.run(t, r.repo, "reap.sh")
	if !strings.Contains(out, "waiting AB12CD: feat/ab12cd has not landed on origin's default branch; gh repo view did not answer (error connecting to api.github.com)") {
		t.Fatalf("gh down, want it named as the cause:\n%s", out)
	}
}

// What .jaira/worktree-setup built, .jaira/worktree-teardown takes down: reap.sh
// calls it with the arguments spawn.sh recorded for the setup, once the tabs
// are closed and before the worktree goes. The record decides, never the
// worktree's name: a worktree spawned from another one carries a slug its name
// does not spell. A teardown that fails or cannot run keeps everything; a
// worktree with no record gets no teardown call.
func TestReapTearsDownWhatTheSetupBuilt(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		script      string // the teardown; "" for none
		mode        os.FileMode
		record      bool   // whether spawn.sh left a jaira-setup record
		want        string // a line reap.sh prints
		called      bool   // whether the teardown ran
		gone        bool   // whether the worktree went
		timeout     string // JAIRA_TEARDOWN_TIMEOUT; "" for the default
		noTimeout   bool   // PATH without timeout and gtimeout
		fakeTimeout bool   // a timeout first on PATH that refuses -k, as Windows' does
	}{
		{name: "none", record: true, gone: true},
		{name: "succeeds", script: "exit 0", mode: 0o755, record: true, want: "tore down <wt>", called: true, gone: true},
		{name: "fails", script: "echo stack busy >&2; exit 3", mode: 0o755, record: true,
			want: "kept AB12CD: teardown failed: <main>/.jaira/worktree-teardown <wt> bbb 9 <root> (exit 3)", called: true},
		// A child that would outlive the teardown holds reap.sh's stderr
		// open: the run returns early only when timeout killed it too.
		{name: "hangs", script: "sleep 30; exit 0", mode: 0o755, record: true, timeout: "1",
			want: "kept AB12CD: teardown failed: <main>/.jaira/worktree-teardown <wt> bbb 9 <root> (timed out after 1s)", called: true},
		// A machine without timeout(1): a teardown cannot be bounded, so a
		// worktree that needs one stays; one without a record goes as before.
		{name: "no timeout", script: "exit 0", mode: 0o755, record: true, noTimeout: true,
			want: "kept AB12CD: no timeout(1) that takes -k to bound <main>/.jaira/worktree-teardown"},
		// Git Bash finding Windows' timeout.exe first: it refuses -k.
		{name: "windows timeout", script: "exit 0", mode: 0o755, record: true, noTimeout: true, fakeTimeout: true,
			want: "kept AB12CD: no timeout(1) that takes -k to bound <main>/.jaira/worktree-teardown"},
		// ... and coreutils' own behind it on PATH, which is the one used.
		{name: "windows timeout first", script: "exit 0", mode: 0o755, record: true, fakeTimeout: true,
			want: "tore down <wt>", called: true, gone: true},
		{name: "no timeout no record", script: "exit 0", mode: 0o755, noTimeout: true, gone: true},
		{name: "not executable", script: "exit 0", mode: 0o644, record: true,
			want: "kept AB12CD: <main>/.jaira/worktree-teardown is not executable"},
		{name: "no record", script: "exit 0", mode: 0o755, gone: true},
		{name: "not executable no record", script: "exit 0", mode: 0o644, gone: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			r := newReapRig(t)
			// As spawn.sh writes it when started from another worktree: the
			// slug bbb, not the ab12cd the worktree's name ends in.
			root := filepath.Join(filepath.Dir(r.wt), "jaira-aaa")
			if c.record {
				gd := r.git(t, r.wt, "rev-parse", "--absolute-git-dir")
				r.write(t, filepath.Join(gd, "jaira-setup"), "worktree="+r.wt+"\nslug=bbb\noffset=9\nroot="+root+"\n")
			}
			calls := filepath.Join(filepath.Dir(r.closed), "teardown-calls")
			if c.script != "" {
				// It records its arguments, whether the worktree still stood
				// and which tabs were closed by then.
				body := "#!/usr/bin/env bash\n" +
					"{ printf '%s|' \"$@\"; [ -d \"$1\" ] && printf 'present|'; tr '\\n' ' ' < " + r.closed + "; echo; } >> " + calls + "\n" +
					c.script + "\n"
				p := filepath.Join(r.repo, ".jaira", "worktree-teardown")
				r.write(t, p, body)
				if err := os.Chmod(p, c.mode); err != nil {
					t.Fatal(err)
				}
			}
			if c.timeout != "" {
				r.env = append(r.env, "JAIRA_TEARDOWN_TIMEOUT="+c.timeout)
			}
			{
				dir := pathOf(r.env)
				if c.noTimeout {
					dir = pathWithout(t, r.env, "timeout", "gtimeout")
				} else {
					// macOS runners have neither timeout nor gtimeout: one
					// last on PATH, so the host's own still goes first.
					dir += string(os.PathListSeparator) + stubTimeout(t, r)
				}
				if c.fakeTimeout {
					fake := t.TempDir()
					r.write(t, filepath.Join(fake, "timeout"), "#!/usr/bin/env bash\necho 'ERROR: Invalid syntax. Default option is not allowed more than '1' time(s).' >&2\nexit 1\n")
					if err := os.Chmod(filepath.Join(fake, "timeout"), 0o755); err != nil {
						t.Fatal(err)
					}
					dir = fake + string(os.PathListSeparator) + dir
				}
				r.env = append(r.env, "PATH="+dir)
			}
			r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))
			r.write(t, r.shown, `{"status":"done"}`)

			dry := r.run(t, r.repo, "reap.sh", "--dry-run")
			if _, err := os.Stat(calls); err == nil {
				t.Fatalf("--dry-run called the teardown:\n%s", dry)
			}
			if c.called && !strings.Contains(dry, "would tear down "+r.wt) {
				t.Errorf("--dry-run does not name the teardown:\n%s", dry)
			}

			start := time.Now()
			got := r.run(t, r.repo, "reap.sh")
			if c.timeout != "" && time.Since(start) > 20*time.Second {
				t.Errorf("reap.sh waited %v for a teardown bounded to %ss: its child outlived the timeout", time.Since(start), c.timeout)
			}
			b, _ := os.ReadFile(calls)
			want := strings.NewReplacer("<wt>", r.wt, "<main>", r.repo, "<root>", root).Replace(c.want)
			if want != "" && !strings.Contains(got, want) {
				t.Errorf("output lacks %q:\n%s", want, got)
			}
			if c.called {
				line := r.wt + "|bbb|9|" + root + "|present|tp1 tp2 \n"
				if string(b) != line {
					t.Errorf("teardown called as %q, want %q (recorded args, worktree still there, tabs closed first)", b, line)
				}
			} else if len(b) != 0 {
				t.Errorf("teardown called: %q\n%s", b, got)
			}
			if c.mode == 0o644 && c.record && r.closedTabs() != "" {
				t.Errorf("a teardown that cannot run closed tabs anyway: %q", r.closedTabs())
			}
			_, err := os.Stat(r.wt)
			if gone := os.IsNotExist(err); gone != c.gone {
				t.Errorf("worktree gone %v, want %v:\n%s", gone, c.gone, got)
			}
			if !c.gone {
				if len(r.records(t)) != 2 || r.git(t, r.repo, "branch", "--list", "feat/ab12cd") == "" {
					t.Errorf("kept ticket lost its records or branch:\n%s", got)
				}
			}
			// A failed teardown is called again in the next round.
			if c.name == "fails" {
				again := r.run(t, r.repo, "reap.sh")
				if b, _ := os.ReadFile(calls); strings.Count(string(b), "\n") != 2 {
					t.Errorf("the next round did not call the teardown again: %q\n%s", b, again)
				}
			}
		})
	}
}

// pathOf is env's PATH: the last one is what exec sees.
func pathOf(env []string) string {
	var path string
	for _, e := range env {
		if v, ok := strings.CutPrefix(e, "PATH="); ok {
			path = v
		}
	}
	return path
}

// stubTimeout is a directory holding a timeout(1) that takes -k as
// coreutils' does: it signals the command's whole process group and exits
// 124 when the time ran out.
func stubTimeout(t *testing.T, r *reapRig) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "timeout")
	r.write(t, p, `#!/usr/bin/env bash
[ "$1" = -k ] || exit 125
k=$2 s=$3; shift 3
flag=$(mktemp) && rm -f "$flag" || exit 125
set -m
"$@" & pid=$!
( sleep "$s"; : > "$flag"; kill -TERM -- -"$pid"; sleep "$k"; kill -KILL -- -"$pid" ) </dev/null >/dev/null 2>&1 & w=$!
wait "$pid"; rc=$?
kill -- -"$w" 2>/dev/null
if [ -e "$flag" ]; then rm -f "$flag"; exit 124; fi
exit "$rc"
`)
	if err := os.Chmod(p, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// pathWithout is env's PATH as one directory of links to every program on it
// but the named ones: the first of each name wins, as in a lookup.
func pathWithout(t *testing.T, env []string, drop ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, d := range filepath.SplitList(pathOf(env)) {
		entries, _ := os.ReadDir(d)
		for _, e := range entries {
			if slices.Contains(drop, e.Name()) {
				continue
			}
			link := filepath.Join(dir, e.Name())
			if _, err := os.Lstat(link); err == nil {
				continue
			}
			if err := os.Symlink(filepath.Join(d, e.Name()), link); err != nil {
				t.Fatal(err)
			}
		}
	}
	return dir
}
