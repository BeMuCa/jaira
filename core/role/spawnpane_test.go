package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// Two dispatchers open a tab at the same moment, and the Herdr stub answers
// the critique's 'tab create' with the pane of the optimize tab - the race
// 3YRPXJ found. The spawn that was handed the other's pane types nothing, says
// which pane and which tab, and closes its own empty tab; the other one starts.
func TestSpawnTypesOnlyIntoItsOwnTab(t *testing.T) {
	state := t.TempDir()
	// Which spawn loses is fixed by its label, not by who reaches the stub
	// first: the critique gets tab t1, the optimize t2, and both are handed
	// p2, which pane get places in t2. The critique waits until t2 exists so
	// the two overlap, but the outcome does not hang on that wait - a counter
	// and a ten-second bound here made both spawns get t1 once on a loaded CI
	// runner (JCBHBE).
	root, spawn := spawnWith(t, `#!/usr/bin/env bash
s="`+state+`"
case "$1 $2" in
  "tab create")
    case "$*" in *AB12CD/critique*) n=1 ;; *) n=2 ;; esac
    mkdir "$s/t$n"
    if [ "$n" = 1 ]; then for _ in $(seq 2400); do [ -d "$s/t2" ] && break; sleep 0.05; done; fi
    echo "{\"result\":{\"tab\":{\"tab_id\":\"t$n\"},\"root_pane\":{\"pane_id\":\"p2\"}}}" ;;
  "pane get") echo "{\"result\":{\"pane\":{\"agent\":\"claude\",\"agent_status\":\"idle\",\"tab_id\":\"t${3#p}\"}}}" ;;
  "pane run") echo "$3" >> "$s/run" ;;
  "pane send-text") echo "$3" >> "$s/sent" ;;
  "tab close") echo "$3" >> "$s/closed" ;;
esac
`)

	var wg sync.WaitGroup
	outs := make([]string, 2)
	errs := make([]error, 2)
	for i, lane := range []string{"critique", "optimize"} {
		wg.Add(1)
		go func(i int, lane string) {
			defer wg.Done()
			b, err := spawn(lane).CombinedOutput()
			outs[i], errs[i] = string(b), err
		}(i, lane)
	}
	wg.Wait()

	if errs[0] == nil || errs[1] != nil {
		t.Fatalf("want the critique, handed optimize's pane, refused and the optimize started; critique: %v, optimize: %v\n%s\n---\n%s", errs[0], errs[1], outs[0], outs[1])
	}
	if !strings.Contains(outs[0], "pane p2 is not in tab t1 (AB12CD/critique)") || !strings.Contains(outs[0], "nothing was typed") {
		t.Errorf("the refused spawn names neither the pane nor its tab:\n%s", outs[0])
	}
	won := "optimize"
	read := func(name string) string {
		b, _ := os.ReadFile(filepath.Join(state, name))
		return strings.TrimSpace(string(b))
	}
	if got := read("run"); got != "p2" {
		t.Errorf("claude was started in %q, want only p2", got)
	}
	if got := read("sent"); got != "p2" {
		t.Errorf("a role was typed into %q, want only p2", got)
	}
	if got := read("closed"); got != "t1" {
		t.Errorf("tabs closed: %q, want the refused spawn's own t1", got)
	}
	// reap.sh's record of p2 belongs to the spawn that started in it, and
	// the refused spawn left none of its own.
	recs, _ := filepath.Glob(filepath.Join(root, ".git", "jaira", "panes", "*"))
	if len(recs) != 1 || filepath.Base(recs[0]) != "p2" {
		t.Fatalf("pane records %v, want only p2", recs)
	}
	if rec, _ := os.ReadFile(recs[0]); !strings.Contains(string(rec), "lane="+won+"\n") {
		t.Errorf("the record of p2 names another lane than %s:\n%s", won, rec)
	}
}

// claude that never comes up is reported with its pane, its tab and the
// ticket/lane the tab carries, not with a bare pane id.
func TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp(t *testing.T) {
	_, spawn := spawnWith(t, `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") echo '{"result":{"tab":{"tab_id":"w1:t7"},"root_pane":{"pane_id":"w1:p9"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"","agent_status":"","tab_id":"w1:t7"}}}' ;;
esac
`)
	out, err := spawn("critique").CombinedOutput()
	if err == nil {
		t.Fatalf("spawn.sh succeeded without claude:\n%s", out)
	}
	for _, want := range []string{"pane w1:p9", "tab w1:t7", "AB12CD/critique"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("the failure does not name %q:\n%s", want, out)
		}
	}
}

// spawnWith installs the roles, makes a repository and returns spawn.sh with
// the Herdr stub herdr, for one lane of AB12CD in that repository.
func spawnWith(t *testing.T, herdr string) (root string, spawn func(lane string) *exec.Cmd) {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	root = filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if b, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, b)
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "herdr"), []byte(herdr), 0o755); err != nil {
		t.Fatal(err)
	}
	return root, func(lane string) *exec.Cmd {
		cmd := exec.Command(bash, script, "--no-worktree", "ab12cd", "AB12CD", lane, root)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
			"HERDR_WORKSPACE_ID=", "JAIRA_SPAWN_PAUSE=0", "PATH="+bin+string(os.PathListSeparator)+path)
		return cmd
	}
}

// A 'pane get' that fails is no race: the spawn stops without typing, and the
// tab - where claude may already run - stays standing.
func TestSpawnKeepsItsTabWhenHerdrDoesNotAnswer(t *testing.T) {
	state := t.TempDir()
	_, spawn := spawnWith(t, `#!/usr/bin/env bash
s="`+state+`"
case "$1 $2" in
  "tab create") echo '{"result":{"tab":{"tab_id":"t1"},"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")
    echo x >> "$s/gets"
    # The first two answer, the one before send-text does not.
    [ "$(wc -l < "$s/gets" | tr -d ' ')" -le 2 ] || exit 1
    echo '{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}' ;;
  "pane send-text") echo "$3" >> "$s/sent" ;;
  "tab close") echo "$3" >> "$s/closed" ;;
esac
`)
	out, err := spawn("critique").CombinedOutput()
	if err == nil {
		t.Fatalf("spawn.sh succeeded although herdr did not answer:\n%s", out)
	}
	if !strings.Contains(string(out), "could not check pane p1 of tab t1") {
		t.Errorf("the failure is not reported as an unanswered check:\n%s", out)
	}
	for _, f := range []string{"sent", "closed"} {
		if _, err := os.Stat(filepath.Join(state, f)); err == nil {
			t.Errorf("herdr saw a %s call after its 'pane get' failed", f)
		}
	}
}

// Each 'pane get' after the first ones fails or names another tab, in the
// three places spawn.sh asks: the check before the record, the come-up loop
// and the check before send-text. Every one stops with a line that says
// where and why, types nothing, and leaves neither an empty tab nobody reaps
// nor a record that names a pane which is not ours.
func TestSpawnStopsCleanlyWhenAPaneGetGoesWrong(t *testing.T) {
	for _, c := range []struct {
		name    string
		answer  string // what 'pane get' number $n prints; exit 1 when empty
		want    string
		closed  string
		records int
	}{
		{"first check unanswered", ``, "the empty tab is closed", "t1", 0},
		{"come-up loop unanswered", `[ "$n" -le 1 ] && echo "$ok"`, "claude did not come up in pane p1 (tab t1, AB12CD/critique) after the wait: state 'none'", "", 1},
		{"late check finds another tab", `if [ "$n" -le 2 ]; then echo "$ok"; else echo '{"result":{"pane":{"tab_id":"t9"}}}'; fi`, "pane p1 is not in tab t1 (AB12CD/critique) any more but in 't9'", "t1", 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			state := t.TempDir()
			answer := c.answer
			if answer == "" {
				answer = "false"
			}
			root, spawn := spawnWith(t, `#!/usr/bin/env bash
s="`+state+`"
ok='{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}'
case "$1 $2" in
  "tab create") echo '{"result":{"tab":{"tab_id":"t1"},"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")
    echo x >> "$s/gets"; n=$(wc -l < "$s/gets" | tr -d ' ')
    `+answer+` || exit 1 ;;
  "pane send-text") echo "$3" >> "$s/sent" ;;
  "tab close") echo "$3" >> "$s/closed" ;;
esac
`)
			out, err := spawn("critique").CombinedOutput()
			if err == nil {
				t.Fatalf("spawn.sh succeeded:\n%s", out)
			}
			if !strings.Contains(string(out), c.want) {
				t.Errorf("the failure does not say %q:\n%s", c.want, out)
			}
			if _, err := os.Stat(filepath.Join(state, "sent")); err == nil {
				t.Errorf("a role was typed although the spawn failed")
			}
			closed, _ := os.ReadFile(filepath.Join(state, "closed"))
			if got := strings.TrimSpace(string(closed)); got != c.closed {
				t.Errorf("tabs closed: %q, want %q", got, c.closed)
			}
			recs, _ := filepath.Glob(filepath.Join(root, ".git", "jaira", "panes", "*"))
			if len(recs) != c.records {
				t.Errorf("pane records %v, want %d", recs, c.records)
			}
		})
	}
}
