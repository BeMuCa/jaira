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
// the first 'tab create' with the pane of the second tab - the race 3YRPXJ
// found. The spawn that was handed the other's pane types nothing, says which
// pane and which tab, and closes its own empty tab; the other one starts.
func TestSpawnTypesOnlyIntoItsOwnTab(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	root := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if b, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, b)
	}
	bin := t.TempDir()
	state := t.TempDir()
	// tab create n makes tab tn with pane pn, waits until both tabs exist and
	// then hands back the newest pane, whoever asked. pane get pk says tk.
	herdr := `#!/usr/bin/env bash
s="` + state + `"
case "$1 $2" in
  "tab create")
    n=1; until mkdir "$s/t$n" 2>/dev/null; do n=$((n+1)); done
    for _ in $(seq 200); do [ -d "$s/t2" ] && break; sleep 0.05; done
    last=2
    echo "{\"result\":{\"tab\":{\"tab_id\":\"t$n\"},\"root_pane\":{\"pane_id\":\"p$last\"}}}" ;;
  "pane get") echo "{\"result\":{\"pane\":{\"agent\":\"claude\",\"agent_status\":\"idle\",\"tab_id\":\"t${3#p}\"}}}" ;;
  "pane run") echo "$3" >> "$s/run" ;;
  "pane send-text") echo "$3" >> "$s/sent" ;;
  "tab close") echo "$3" >> "$s/closed" ;;
esac
`
	if err := os.WriteFile(filepath.Join(bin, "herdr"), []byte(herdr), 0o755); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	outs := make([]string, 2)
	errs := make([]error, 2)
	for i, lane := range []string{"critique", "optimize"} {
		wg.Add(1)
		go func(i int, lane string) {
			defer wg.Done()
			cmd := exec.Command(bash, spawn, "--no-worktree", "ab12cd", "AB12CD", lane, root)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
				"HERDR_WORKSPACE_ID=", "JAIRA_SPAWN_PAUSE=0", "PATH="+bin+string(os.PathListSeparator)+path)
			b, err := cmd.CombinedOutput()
			outs[i], errs[i] = string(b), err
		}(i, lane)
	}
	wg.Wait()

	failed := 0
	for i := range outs {
		if errs[i] == nil {
			continue
		}
		failed++
		if !strings.Contains(outs[i], "pane p2 is not in tab t1") || !strings.Contains(outs[i], "nothing was typed") {
			t.Errorf("the refused spawn names neither the pane nor its tab:\n%s", outs[i])
		}
	}
	if failed != 1 {
		t.Fatalf("%d spawns failed, want exactly the one handed the other's pane:\n%s\n---\n%s", failed, outs[0], outs[1])
	}
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
}

// claude that never comes up is reported with its pane, its tab and the
// ticket/lane the tab carries, not with a bare pane id.
func TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	root := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if b, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, b)
	}
	bin := t.TempDir()
	herdr := `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") echo '{"result":{"tab":{"tab_id":"w1:t7"},"root_pane":{"pane_id":"w1:p9"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"","agent_status":"","tab_id":"w1:t7"}}}' ;;
esac
`
	if err := os.WriteFile(filepath.Join(bin, "herdr"), []byte(herdr), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bash, spawn, "--no-worktree", "ab12cd", "AB12CD", "critique", root)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
		"HERDR_WORKSPACE_ID=", "JAIRA_SPAWN_PAUSE=0", "PATH="+bin+string(os.PathListSeparator)+path)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("spawn.sh succeeded without claude:\n%s", out)
	}
	for _, want := range []string{"pane w1:p9", "tab w1:t7", "AB12CD/critique"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("the failure does not name %q:\n%s", want, out)
		}
	}
}
