package cli

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

func stateDirOf(t *testing.T, dir string) string {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	return s.StateDir()
}

func TestUpdateDryRunWritesNothing(t *testing.T) {
	st := newStaleSetup(t)
	md := mustRead(t, st.claudeMD)
	stateDir := stateDirOf(t, st.dir)
	stampBefore := release.Stamped(stateDir)
	setCurrent(t, "2.0.0") // a newer binary than the one that stamped the board

	out, err := runUpdate(t, st.dir, "--roles", "--lanes", "--dry-run")
	if err != nil {
		t.Fatalf("update --dry-run: %v\n%s", err, out)
	}
	if string(mustRead(t, st.claudeMD)) != string(md) {
		t.Error("--dry-run wrote CLAUDE.md")
	}
	for p, old := range map[string][]byte{st.projRole: st.roleOld, st.homeRole: st.roleOld, st.lan: st.laneOld} {
		if string(mustRead(t, p)) != string(old) {
			t.Errorf("--dry-run wrote %s", p)
		}
	}
	if got := release.Stamped(stateDir); got != stampBefore {
		t.Errorf("--dry-run moved the stamp from %q to %q", stampBefore, got)
	}
	for _, want := range []string{"would update " + st.projRole, "would update " + st.homeRole, "would update " + st.lan,
		"Would tell your agents about this board: CLAUDE.md (updated)", "nothing was written (--dry-run)"} {
		if !hasLine(setupOutput(out), want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestUpdateDryRunExitCodeMatchesTheRealRun(t *testing.T) {
	st := newStaleSetup(t)
	mustWrite(t, st.projRole, []byte("my own words\n"))

	_, dryErr := runUpdate(t, st.dir, "--roles", "--dry-run")
	_, realErr := runUpdate(t, st.dir, "--roles")
	if exitCode(dryErr) != ExitValidation || exitCode(realErr) != ExitValidation {
		t.Fatalf("exit codes dry %d, real %d, want both %d", exitCode(dryErr), exitCode(realErr), ExitValidation)
	}
}

func TestUpdateDryRunJSONSaysSo(t *testing.T) {
	st := newStaleSetup(t)
	out, err := runUpdate(t, st.dir, "--json", "--dry-run")
	if err != nil {
		t.Fatalf("update --json --dry-run: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"dry_run": true`) && !strings.Contains(out, `"dry_run":true`) {
		t.Errorf("no dry_run key:\n%s", out)
	}
}

func TestUpdateWithoutSwitchesNamesWhatIsBehind(t *testing.T) {
	st := newStaleSetup(t)
	out, err := runUpdate(t, st.dir)
	if err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	if !hasLine(setupOutput(out), "lanes: ") {
		t.Errorf("output lacks the lanes hint:\n%s", out)
	}
	if string(mustRead(t, st.lan)) != string(st.laneOld) {
		t.Error("a hint wrote the lane file")
	}
}

func TestUpdateWithoutSwitchesSaysNothingWhenCurrent(t *testing.T) {
	st := newStaleSetup(t)
	if out, err := runUpdate(t, st.dir, "--roles", "--lanes"); err != nil {
		t.Fatalf("update --roles --lanes: %v\n%s", err, out)
	}
	out, err := runUpdate(t, st.dir)
	if err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	if head := setupOutput(out); hasLine(head, "roles: ") || hasLine(head, "lanes: ") ||
		strings.Contains(head, "jaira update --roles") || strings.Contains(head, "jaira update --lanes") {
		t.Errorf("a current setup got a hint:\n%s", out)
	}
}
