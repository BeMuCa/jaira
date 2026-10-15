package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/lane"
)

// staleSetup is a board and a home whose block, project roles, personal roles and
// testing lane are all older than this binary: the state after an upgrade.
type staleSetup struct {
	dir, home                         string
	claudeMD, projRole, homeRole, lan string
	roleOld, laneOld                  []byte
}

func newStaleSetup(t *testing.T) *staleSetup {
	t.Helper()
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	roleOld, err := os.ReadFile(filepath.Join("..", "..", "core", "role", "testdata", "jaira-role-acceptance-SKILL-v0.4.0.md"))
	if err != nil {
		t.Fatal(err)
	}
	laneOld, err := os.ReadFile(filepath.Join("..", "..", "core", "lane", "testdata", "testing-catalogue-2fb8f177.md"))
	if err != nil {
		t.Fatal(err)
	}
	st := &staleSetup{dir: dir, home: home, roleOld: roleOld, laneOld: laneOld,
		claudeMD: filepath.Join(dir, "CLAUDE.md"),
		projRole: filepath.Join(dir, ".claude", "skills", "jaira-role-acceptance", "SKILL.md"),
		homeRole: filepath.Join(home, ".claude", "skills", "jaira-role-acceptance", "SKILL.md"),
	}
	for _, args := range [][]string{{"update"}, {"lanes", "add", "testing"}, {"roles", "install", "--project"}, {"roles", "install", "--global"}} {
		if out, err := runCLI(t, dir, args...); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	set, err := lane.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("testing")
	if !ok {
		t.Fatal("testing lane not on the board")
	}
	st.lan = l.Source

	md, err := os.ReadFile(st.claudeMD)
	if err != nil {
		t.Fatal(err)
	}
	stale := strings.Replace(string(md), "<!-- jaira:start -->\n", "<!-- jaira:start -->\nAN OLD LINE\n", 1)
	if stale == string(md) {
		t.Fatal("no jaira block in CLAUDE.md to make stale")
	}
	mustWrite(t, st.claudeMD, []byte(stale))
	mustWrite(t, st.projRole, roleOld)
	mustWrite(t, st.homeRole, roleOld)
	mustWrite(t, st.lan, laneOld)
	return st
}

func mustWrite(t *testing.T, path string, b []byte) {
	t.Helper()
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func (st *staleSetup) blockCurrent(t *testing.T) bool {
	t.Helper()
	set, err := lane.Load(st.dir)
	if err != nil {
		t.Fatal(err)
	}
	cur, _ := board.NoteIsCurrent(st.dir, set.Facts())
	return cur
}

func TestUpdateRolesLanesBringsAStaleBoardUpToDate(t *testing.T) {
	st := newStaleSetup(t)

	out, err := runUpdate(t, st.dir, "--roles", "--lanes")
	if err != nil {
		t.Fatalf("update --roles --lanes: %v\n%s", err, out)
	}
	for _, p := range []string{st.projRole, st.homeRole, st.lan} {
		if !strings.Contains(out, "updated "+p) {
			t.Errorf("output does not name %s as updated:\n%s", p, out)
		}
	}
	if string(mustRead(t, st.projRole)) == string(st.roleOld) || string(mustRead(t, st.homeRole)) == string(st.roleOld) {
		t.Error("a role file is still the old version")
	}
	if string(mustRead(t, st.lan)) == string(st.laneOld) {
		t.Error("the lane file is still the old version")
	}
	if !st.blockCurrent(t) {
		t.Error("the jaira block is still stale")
	}

	out, err = runUpdate(t, st.dir, "--roles", "--lanes")
	if err != nil {
		t.Fatalf("second run: %v\n%s", err, out)
	}
	if strings.Contains(out, "updated /") || strings.Contains(out, "would ") {
		t.Errorf("second run still updated something:\n%s", out)
	}
	if !strings.Contains(out, "lanes: 0 updated") {
		t.Errorf("second run lane summary:\n%s", out)
	}
}

func TestUpdateWithoutSwitchesWritesOnlyTheBlock(t *testing.T) {
	st := newStaleSetup(t)
	if out, err := runUpdate(t, st.dir); err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	for p, old := range map[string][]byte{st.projRole: st.roleOld, st.homeRole: st.roleOld, st.lan: st.laneOld} {
		if string(mustRead(t, p)) != string(old) {
			t.Errorf("a bare update wrote %s", p)
		}
	}
	if !st.blockCurrent(t) {
		t.Error("the jaira block is still stale")
	}
}

func TestUpdateRolesSkipsAnEditedFileAndExitsThree(t *testing.T) {
	st := newStaleSetup(t)
	edited := []byte("my own words\n")
	mustWrite(t, st.projRole, edited)

	out, err := runUpdate(t, st.dir, "--json", "--roles")
	if exitCode(err) != ExitValidation {
		t.Fatalf("exit = %d (%v), want %d", exitCode(err), err, ExitValidation)
	}
	if string(mustRead(t, st.projRole)) != string(edited) {
		t.Error("the edited role file was touched")
	}
	if string(mustRead(t, st.homeRole)) == string(st.roleOld) {
		t.Error("the stale personal file was not updated beside the edited one")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out)
	}
	if payload["skipped"] != true {
		t.Errorf("skipped = %#v, want true", payload["skipped"])
	}
}

func TestUpdateLanesSkipsAnEditedLaneAndExitsThree(t *testing.T) {
	st := newStaleSetup(t)
	edited := append(append([]byte{}, st.laneOld...), []byte("\nan extra line\n")...)
	mustWrite(t, st.lan, edited)

	out, err := runUpdate(t, st.dir, "--lanes")
	if exitCode(err) != ExitValidation {
		t.Fatalf("exit = %d (%v), want %d\n%s", exitCode(err), err, ExitValidation, out)
	}
	if string(mustRead(t, st.lan)) != string(edited) {
		t.Error("the edited lane file was touched")
	}
	if !strings.Contains(out, "skipped "+st.lan) {
		t.Errorf("output does not name the skipped lane:\n%s", out)
	}
}

func TestUpdateRolesOutsideABoardUpdatesPersonalSkillsOnly(t *testing.T) {
	st := newStaleSetup(t)
	elsewhere := t.TempDir()

	out, err := runUpdate(t, elsewhere, "--roles")
	if err != nil {
		t.Fatalf("update --roles outside a board: %v\n%s", err, out)
	}
	if !strings.Contains(out, "no board here") {
		t.Errorf("output lacks the no-board line:\n%s", out)
	}
	if string(mustRead(t, st.homeRole)) == string(st.roleOld) {
		t.Error("the personal role file was not updated")
	}
	if string(mustRead(t, st.projRole)) != string(st.roleOld) {
		t.Error("the project role file of another board was touched")
	}
}

func TestUpdateLanesOutsideABoardStillFails(t *testing.T) {
	newStaleSetup(t)
	out, err := runUpdate(t, t.TempDir(), "--lanes")
	if err == nil {
		t.Fatalf("update --lanes outside a board succeeded:\n%s", out)
	}
	if exitCode(err) == ExitOK {
		t.Errorf("exit code = %d", exitCode(err))
	}
}

func TestUpdateJSONCarriesRolesAndLanesOnlyWhenAsked(t *testing.T) {
	st := newStaleSetup(t)

	out, err := runUpdate(t, st.dir, "--json")
	if err != nil {
		t.Fatalf("update --json: %v\n%s", err, out)
	}
	var bare map[string]any
	if err := json.Unmarshal([]byte(out), &bare); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	for _, k := range []string{"roles", "lanes", "skipped", "dry_run"} {
		if _, ok := bare[k]; ok {
			t.Errorf("a bare update carries %q", k)
		}
	}

	out, err = runUpdate(t, st.dir, "--json", "--roles", "--lanes")
	if err != nil {
		t.Fatalf("update --json --roles --lanes: %v\n%s", err, out)
	}
	var full map[string]any
	if err := json.Unmarshal([]byte(out), &full); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	for _, k := range []string{"roles", "lanes", "skipped"} {
		if _, ok := full[k]; !ok {
			t.Errorf("update --roles --lanes lacks %q: %v", k, full)
		}
	}
}
