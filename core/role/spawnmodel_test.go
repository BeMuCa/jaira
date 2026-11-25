package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Each agent spawn.sh starts is started on the model and effort it is given:
// the dispatcher and a worker on the worker model 'jaira roles list' names, any
// other lane on the model its tier has, which 'jaira show --for-lane' names.
// Herdr and jaira are stubs that answer from files; the command each tab is
// started with is what is checked.
func TestSpawnStartsEachAgentOnItsModel(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	real := strings.SplitN(path, string(os.PathListSeparator), 2)[0]
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	ran := filepath.Join(t.TempDir(), "ran")
	answers := t.TempDir()
	herdr := `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") echo '{"result":{"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}' ;;
  "pane run")   printf '%s\n' "$4" > "` + ran + `" ;;
esac
`
	jaira := `#!/usr/bin/env bash
case "$*" in
  "json "*) exec "` + filepath.Join(real, "jaira") + `" "$@" ;;
  "-C "*" roles list --json") cat "` + answers + `/roles.json" ;;
  *" show AB12CD --for-lane "*"--json")
    args="$*"; lane="${args##*--for-lane }"; lane="${lane%% *}"
    cat "` + answers + `/$lane.json" || exit 1 ;;
  *) exit 1 ;;
esac
`
	for name, script := range map[string]string{"herdr": herdr, "jaira": jaira} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	run := func(lane, roles, laneJSON string) string {
		t.Helper()
		if err := os.WriteFile(filepath.Join(answers, "roles.json"), []byte(roles), 0o644); err != nil {
			t.Fatal(err)
		}
		laneFile := filepath.Join(answers, lane+".json")
		_ = os.Remove(laneFile)
		if laneJSON != "" {
			if err := os.WriteFile(laneFile, []byte(laneJSON), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		cmd := spawnOn(bash, "Linux", spawn, "--no-worktree", "ab12cd", "AB12CD", lane)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
			"HERDR_WORKSPACE_ID=", "PATH="+bin+string(os.PathListSeparator)+path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("spawn.sh %s: %v\n%s", lane, err, out)
		}
		r, _ := os.ReadFile(ran)
		return strings.TrimSpace(string(r))
	}

	for _, c := range []struct{ name, lane, roles, laneJSON, want string }{
		{"the dispatcher starts on the worker model", "dispatch",
			`{"worker_model":"sonnet","worker_effort":"xhigh"}`, "", "&& claude --model 'sonnet' --effort xhigh"},
		{"a worker starts on the worker model", "worker",
			`{"worker_model":"haiku","worker_effort":"high"}`, "", "&& claude --model 'haiku' --effort high"},
		{"a lane starts on its tier's model", "review",
			`{}`, `{"model":"opus","effort":""}`, "&& claude --model 'opus'"},
		{"a lane's effort follows its model", "in-progress",
			`{}`, `{"model":"haiku","effort":"high"}`, "&& claude --model 'haiku' --effort high"},
		{"a lane with no model named keeps the account default", "testing",
			`{}`, "", "&& claude"},
		{"a model or effort that is no name or level is dropped", "review",
			`{}`, `{"model":"x;rm -rf ~","effort":"huge"}`, "&& claude"},
		{"pr starts on none", "pr",
			`{"worker_model":"sonnet"}`, "", "&& claude"},
	} {
		if start := run(c.lane, c.roles, c.laneJSON); !strings.HasSuffix(start, c.want) {
			t.Errorf("%s: started %q, want it to end %q", c.name, start, c.want)
		}
	}
}
