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

// setupOutput is what a plain-text update printed before its release-notes
// block. The notes are not the setup's: under a test version every shipped note
// is printed, and they quote the very words (would, jaira update --roles) the
// setup lines are checked for.
func setupOutput(out string) string {
	cut := len(out)
	for _, m := range []string{"\nWhat's changed since then:", "\nNothing has changed since"} {
		if i := strings.Index(out, m); i >= 0 && i < cut {
			cut = i
		}
	}
	return out[:cut]
}

// hasLine reports whether a line of out starts with prefix.
func hasLine(out, prefix string) bool {
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, prefix) {
			return true
		}
	}
	return false
}

func TestUpdateRolesLanesBringsAStaleBoardUpToDate(t *testing.T) {
	st := newStaleSetup(t)

	out, err := runUpdate(t, st.dir, "--roles", "--lanes")
	if err != nil {
		t.Fatalf("update --roles --lanes: %v\n%s", err, out)
	}
	for _, p := range []string{st.projRole, st.homeRole, st.lan} {
		if !hasLine(setupOutput(out), "updated "+p) {
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
	head := setupOutput(out)
	if hasLine(head, "updated ") || hasLine(head, "would ") {
		t.Errorf("second run still updated something:\n%s", out)
	}
	if !hasLine(head, "lanes: 0 updated") {
		t.Errorf("second run lane summary:\n%s", out)
	}
}

func TestUpdateWithoutSwitchesUpdatesRolesButNotLanes(t *testing.T) {
	st := newStaleSetup(t)
	if out, err := runUpdate(t, st.dir); err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	for _, p := range []string{st.projRole, st.homeRole} {
		if string(mustRead(t, p)) == string(st.roleOld) {
			t.Errorf("a bare update left the stale role %s", p)
		}
	}
	if string(mustRead(t, st.lan)) != string(st.laneOld) {
		t.Error("a bare update wrote the lane file; lanes need --lanes")
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
	if !hasLine(setupOutput(out), "skipped "+st.lan) {
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
	if !hasLine(setupOutput(out), "no board here") {
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
	if _, ok := bare["roles"]; !ok {
		t.Error("a bare update lacks \"roles\": roles are updated by default")
	}
	for _, k := range []string{"lanes", "dry_run"} {
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
