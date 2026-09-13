package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/role"
	"github.com/BeMuCa/jaira/core/ticket"
)

// boardAt gives a test a real store, which 'roles install --project' needs to
// find the repository root.
func boardAt(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir
}

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
			Name        string   `json:"name"`
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
		if r.ID == "" || r.Name != r.ID || r.Description == "" || len(r.Files) == 0 {
			t.Errorf("incomplete role entry: %+v", r)
		}
	}
}

// --project with no agent directory in the repository still lands somewhere an
// agent looks: .claude/skills.
func TestRolesInstallProjectWritesSevenRoles(t *testing.T) {
	dir := boardAt(t)
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

// Every agent directory that exists gets the roles, because a project may be
// worked with more than one tool.
func TestRolesInstallProjectFollowsExistingAgentDirs(t *testing.T) {
	dir := boardAt(t)
	if err := os.MkdirAll(filepath.Join(dir, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "roles", "install", "--project"); err != nil {
		t.Fatalf("roles install: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".codex", "skills", "jaira-teamlead", "SKILL.md")); err != nil {
		t.Errorf(".codex was not written to: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills", "jaira-teamlead")); !os.IsNotExist(err) {
		t.Error(".claude/skills was created although .codex already existed")
	}
}

func TestRolesInstallGlobalWritesUnderHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
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
	dir := boardAt(t)
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
	dir := boardAt(t)
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
	dir := boardAt(t)
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
		Skipped    bool     `json:"skipped"`
		Unprefixed []string `json:"unprefixed"`
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
	dir := boardAt(t)
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

// An older, unprefixed copy of a role answers to a different command name. It
// is reported so nobody keeps invoking the stale one, and never deleted.
func TestRolesInstallReportsAnUnprefixedCopy(t *testing.T) {
	dir := boardAt(t)
	bare := filepath.Join(dir, ".claude", "skills", "teamlead")
	if err := os.MkdirAll(bare, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bare, "SKILL.md"), []byte("---\nname: teamlead\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "roles", "install", "--project")
	if err != nil {
		t.Fatalf("roles install: %v\n%s", err, out)
	}
	if !strings.Contains(out, bare) {
		t.Errorf("the unprefixed copy was not reported:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(bare, "SKILL.md")); err != nil {
		t.Errorf("the unprefixed copy was removed: %v", err)
	}
}
