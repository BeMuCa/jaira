package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/role"
)

func exitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	var ce *codedError
	if errors.As(err, &ce) {
		return ce.code
	}
	return ExitError
}

func TestRolesListNamesEveryBuiltin(t *testing.T) {
	dir := t.TempDir()
	out, err := runCLI(t, dir, "roles", "list")
	if err != nil {
		t.Fatalf("roles list: %v\n%s", err, out)
	}
	roles, err := role.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) == 0 {
		t.Fatal("no built-in roles")
	}
	for _, r := range roles {
		if !strings.Contains(out, r.ID) {
			t.Errorf("roles list did not name %s:\n%s", r.ID, out)
		}
	}
}

func TestRolesListJSON(t *testing.T) {
	dir := t.TempDir()
	out, err := runCLI(t, dir, "--json", "roles", "list")
	if err != nil {
		t.Fatalf("roles list --json: %v\n%s", err, out)
	}
	var payload struct {
		Roles []struct {
			ID          string   `json:"id"`
			Description string   `json:"description"`
			Files       []string `json:"files"`
		} `json:"roles"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if len(payload.Roles) == 0 {
		t.Fatal("no roles in json output")
	}
	for _, r := range payload.Roles {
		if r.ID == "" || r.Description == "" || len(r.Files) == 0 {
			t.Errorf("incomplete role entry: %+v", r)
		}
	}
}

// --project with no agent directory in the repository still lands somewhere an
// agent looks: .claude/skills.
func TestRolesInstallProjectWritesSevenRoles(t *testing.T) {
	dir := lanesTestProject(t)
	out, err := runCLI(t, dir, "roles", "install", "--project")
	if err != nil {
		t.Fatalf("roles install --project: %v\n%s", err, out)
	}
	roles, err := role.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		if _, err := os.Stat(filepath.Join(dir, ".claude", "skills", r.ID, "SKILL.md")); err != nil {
			t.Errorf("after install: %v", err)
		}
	}
	if !strings.Contains(out, "written") {
		t.Errorf("output did not report what it wrote:\n%s", out)
	}
}

// One directory, and only the one asked for. An agent directory that happens to
// exist beside it is not a second install target: a harness reading both would
// find the same role registered twice under one command name.
func TestRolesInstallProjectIgnoresOtherAgentDirs(t *testing.T) {
	dir := lanesTestProject(t)
	if err := os.MkdirAll(filepath.Join(dir, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "roles", "install", "--project"); err != nil {
		t.Fatalf("roles install: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".codex", "skills")); !os.IsNotExist(err) {
		t.Error(".codex/skills was written to; --project installs into .claude/skills only")
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills", "jaira-teamlead", "SKILL.md")); err != nil {
		t.Errorf(".claude/skills was not written to: %v", err)
	}
}

// --into is how a project whose agent reads somewhere else gets the roles,
// rather than every candidate directory getting a copy.
func TestRolesInstallIntoNamesTheDirectory(t *testing.T) {
	dir := lanesTestProject(t)
	want := filepath.Join(dir, ".codex", "skills")
	if out, err := runCLI(t, dir, "roles", "install", "--into", want); err != nil {
		t.Fatalf("roles install --into: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(want, "jaira-teamlead", "SKILL.md")); err != nil {
		t.Errorf("--into did not write there: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills")); !os.IsNotExist(err) {
		t.Error(".claude/skills was written to as well; --into replaces the default")
	}
}

func TestRolesInstallGlobalWritesUnderHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// os.UserHomeDir reads USERPROFILE on Windows and HOME everywhere else, so
	// setting only HOME let --global install into the runner's real home and
	// the assertion looked in a temp directory nothing had been written to.
	t.Setenv("USERPROFILE", home)
	dir := t.TempDir()
	if out, err := runCLI(t, dir, "roles", "install", "--global"); err != nil {
		t.Fatalf("roles install --global: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude", "skills", "jaira-teamlead", "SKILL.md")); err != nil {
		t.Errorf("--global did not write into ~/.claude/skills: %v", err)
	}
}

// An already-installed file is not an edited one: a second run must be quiet
// and exit 0, or every 'jaira update' looks like a conflict.
func TestRolesInstallSecondRunExitsZero(t *testing.T) {
	dir := lanesTestProject(t)
	if out, err := runCLI(t, dir, "roles", "install", "--project"); err != nil {
		t.Fatalf("first install: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "roles", "install", "--project")
	if code := exitCode(err); code != ExitOK {
		t.Fatalf("second run exited %d, want %d\n%s", code, ExitOK, out)
	}
	if !strings.Contains(out, "0 written") {
		t.Errorf("second run reported writes:\n%s", out)
	}
}

func TestRolesInstallLeavesAnEditedFileAloneAndExitsThree(t *testing.T) {
	dir := lanesTestProject(t)
	if out, err := runCLI(t, dir, "roles", "install", "--project"); err != nil {
		t.Fatalf("first install: %v\n%s", err, out)
	}
	mine := filepath.Join(dir, ".claude", "skills", "jaira-role-lane", "SKILL.md")
	const edited = "---\nname: jaira-role-lane\n---\n\nmy own version\n"
	if err := os.WriteFile(mine, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "roles", "install", "--project")
	if code := exitCode(err); code != ExitValidation {
		t.Fatalf("exited %d, want %d\n%s", code, ExitValidation, out)
	}
	if !strings.Contains(out, "skipped") || !strings.Contains(out, "--force") {
		t.Errorf("output did not report the skip and the way past it:\n%s", out)
	}
	if b, _ := os.ReadFile(mine); string(b) != edited {
		t.Error("the edited file was replaced without --force")
	}

	out, err = runCLI(t, dir, "roles", "install", "--project", "--force")
	if code := exitCode(err); code != ExitOK {
		t.Fatalf("--force exited %d, want %d\n%s", code, ExitOK, out)
	}
	want, _ := role.File("jaira-role-lane", "SKILL.md")
	if b, _ := os.ReadFile(mine); string(b) != string(want) {
		t.Error("--force did not restore the built-in version")
	}
}

func TestRolesInstallJSONReportsEveryFile(t *testing.T) {
	dir := lanesTestProject(t)
	out, err := runCLI(t, dir, "--json", "roles", "install", "--project")
	if err != nil {
		t.Fatalf("roles install --json: %v\n%s", err, out)
	}
	var payload struct {
		Installed []struct {
			Role   string `json:"role"`
			Path   string `json:"path"`
			Action string `json:"action"`
		} `json:"installed"`
		Skipped bool `json:"skipped"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if len(payload.Installed) == 0 {
		t.Fatal("no files in json output")
	}
	for _, f := range payload.Installed {
		if f.Action != string(role.Written) {
			t.Errorf("%s: action %q on a first install, want %q", f.Path, f.Action, role.Written)
		}
	}
	if payload.Skipped {
		t.Error("a first install reported a skip")
	}
}

// Neither flag, or both, is a usage mistake and must exit 2 rather than doing
// something the caller did not ask for.
func TestRolesInstallNeedsExactlyOneTarget(t *testing.T) {
	dir := lanesTestProject(t)
	for _, args := range [][]string{
		{"roles", "install"},
		{"roles", "install", "--project", "--global"},
	} {
		out, err := runCLI(t, dir, args...)
		if code := exitCode(err); code != ExitUsage {
			t.Errorf("%v exited %d, want %d\n%s", args, code, ExitUsage, out)
		}
	}
}
