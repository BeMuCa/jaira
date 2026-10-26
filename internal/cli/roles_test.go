package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/role"
	"github.com/BeMuCa/jaira/core/settings"
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

// The teamlead reads which role to start for a ticket off 'roles list': the
// dispatcher by default, the worker once settings.json says "dispatch":
// "worker". spawn.sh reads the model for the cockpit's chats off it too:
// empty, the account's own, until settings.json says "chat-model".
func TestRolesListReportsTheDispatchSetting(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	dir := t.TempDir()
	read := func() map[string]any {
		t.Helper()
		out, err := runCLI(t, dir, "--json", "roles", "list")
		if err != nil {
			t.Fatalf("roles list --json: %v\n%s", err, out)
		}
		var p map[string]any
		if err := json.Unmarshal([]byte(out), &p); err != nil {
			t.Fatalf("invalid json: %v\n%s", err, out)
		}
		return p
	}
	if p := read(); p["dispatch"] != "lanes" || p["dispatch_role"] != "jaira-dispatcher" || p["chat_model"] != "" {
		t.Errorf("default: dispatch=%v role=%v chat_model=%v, want lanes/jaira-dispatcher/empty", p["dispatch"], p["dispatch_role"], p["chat_model"])
	}
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(`{"dispatch": "worker", "chat-model": "sonnet"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if p := read(); p["dispatch"] != "worker" || p["dispatch_role"] != "jaira-worker" || p["chat_model"] != "sonnet" {
		t.Errorf("opted in: dispatch=%v role=%v chat_model=%v, want worker/jaira-worker/sonnet", p["dispatch"], p["dispatch_role"], p["chat_model"])
	}
	out, err := runCLI(t, dir, "roles", "list")
	if err != nil || !strings.Contains(out, "dispatch: worker — a ticket handed to an agent goes to /jaira-worker") {
		t.Errorf("roles list does not say which role a ticket goes to: %v\n%s", err, out)
	}
	if !strings.Contains(out, "chat model: sonnet") || !strings.Contains(out, "settings.json on this machine") {
		t.Errorf("roles list does not name the chat model and where it came from:\n%s", out)
	}

	// A board that names its own model wins over the machine's.
	board := lanesTestProject(t)
	if err := settings.SetBoardChatModel(board, "haiku"); err != nil {
		t.Fatal(err)
	}
	out, err = runCLI(t, board, "--json", "roles", "list")
	var p map[string]any
	if err != nil || json.Unmarshal([]byte(out), &p) != nil {
		t.Fatalf("roles list --json on a board: %v\n%s", err, out)
	}
	if p["chat_model"] != "haiku" || !strings.Contains(fmt.Sprint(p["chat_model_source"]), ".jaira/chat-model") {
		t.Errorf("on a board with its own model: chat_model=%v source=%v, want haiku from .jaira/chat-model", p["chat_model"], p["chat_model_source"])
	}
	if p["release"] != "worker" {
		t.Errorf("a board without .jaira/release: release=%v, want worker", p["release"])
	}
}

// The worker and the lane read who touches the version and the changelog off
// 'roles list': the board's .jaira/release, worker without it.
func TestRolesListReportsTheBoardsRelease(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	board := lanesTestProject(t)
	if err := os.WriteFile(settings.BoardReleasePath(board), []byte("teamlead\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, board, "--json", "roles", "list")
	var p map[string]any
	if err != nil || json.Unmarshal([]byte(out), &p) != nil {
		t.Fatalf("roles list --json on a board: %v\n%s", err, out)
	}
	if p["release"] != "teamlead" || !strings.Contains(fmt.Sprint(p["release_source"]), ".jaira/release") {
		t.Errorf("release=%v source=%v, want teamlead from .jaira/release", p["release"], p["release_source"])
	}
	out, err = runCLI(t, board, "roles", "list")
	if err != nil || !strings.Contains(out, "release: teamlead") {
		t.Errorf("roles list does not say who owns the release: %v\n%s", err, out)
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

// --agent copilot writes where every Copilot surface reads: .github/skills in
// the repository, ~/.copilot/skills at home — and not the cockpit, which only
// Claude Code loads.
func TestRolesInstallAgentCopilotWritesWhereCopilotReads(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dir := lanesTestProject(t)
	for _, c := range []struct {
		flag, want string
	}{
		{"--project", filepath.Join(dir, ".github", "skills")},
		{"--global", filepath.Join(home, ".copilot", "skills")},
	} {
		out, err := runCLI(t, dir, "roles", "install", "--agent", "copilot", c.flag)
		if err != nil {
			t.Fatalf("roles install --agent copilot %s: %v\n%s", c.flag, err, out)
		}
		if _, err := os.Stat(filepath.Join(c.want, "jaira-teamlead", "SKILL.md")); err != nil {
			t.Errorf("%s: not written into %s: %v", c.flag, c.want, err)
		}
		if _, err := os.Stat(filepath.Join(c.want, "jaira-cockpit")); !os.IsNotExist(err) {
			t.Errorf("%s: the cockpit was written for copilot, which does not load it", c.flag)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills")); !os.IsNotExist(err) {
		t.Error(".claude/skills was written to; --agent copilot writes into .github/skills")
	}
}

// Copilot reads .claude/skills too: roles already there are installed for it,
// and a copy in .github/skills would register every one of them twice.
func TestRolesInstallAgentCopilotSkipsWhereRolesAlreadyAre(t *testing.T) {
	dir := lanesTestProject(t)
	if out, err := runCLI(t, dir, "roles", "install", "--project"); err != nil {
		t.Fatalf("roles install --project: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "roles", "install", "--agent", "copilot", "--project")
	if code := exitCode(err); code != ExitOK {
		t.Fatalf("exited %d, want 0\n%s", code, out)
	}
	if lines := strings.Split(strings.TrimSpace(out), "\n"); len(lines) != 1 ||
		!strings.Contains(lines[0], filepath.Join(dir, ".claude", "skills")) ||
		!strings.Contains(lines[0], "'jaira roles install --project'") {
		t.Errorf("want one line naming .claude/skills and how to update the roles there, got:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".github", "skills")); !os.IsNotExist(err) {
		t.Error(".github/skills was written to beside the roles in .claude/skills")
	}
}

func TestRolesInstallAgentRefusesUnknownAndInto(t *testing.T) {
	dir := lanesTestProject(t)
	out, err := runCLI(t, dir, "roles", "install", "--agent", "vim", "--project")
	if code := exitCode(err); code != ExitUsage {
		t.Errorf("unknown agent exited %d, want %d\n%s", code, ExitUsage, out)
	}
	if msg := fmt.Sprint(err) + out; !strings.Contains(msg, "claude") || !strings.Contains(msg, "copilot") {
		t.Errorf("unknown agent does not name the known ones: %s", msg)
	}
	out, err = runCLI(t, dir, "roles", "install", "--agent", "copilot", "--into", t.TempDir())
	if code := exitCode(err); code != ExitUsage {
		t.Errorf("--agent with --into exited %d, want %d\n%s", code, ExitUsage, out)
	}
}

// The flag help names the directory of every agent, not only Claude's: with
// --agent copilot, .claude/skills would be the wrong answer.
func TestRolesInstallHelpNamesEveryAgentsDirectories(t *testing.T) {
	out, err := runCLI(t, t.TempDir(), "roles", "install", "--help")
	if err != nil {
		t.Fatalf("roles install --help: %v\n%s", err, out)
	}
	for _, want := range []string{"copilot: .github/skills", "copilot: ~/.copilot/skills", "claude: .claude/skills", "claude: ~/.claude/skills"} {
		if !strings.Contains(out, want) {
			t.Errorf("--help does not name %q:\n%s", want, out)
		}
	}
}

// The roles run their Python scripts with what 'roles list' names: python3
// unless settings.json says "python", and JAIRA_PYTHON over both — on Windows
// python3 is often missing and "uv run python" is what there is (issue #79).
func TestRolesListReportsThePythonCommand(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	t.Setenv(settings.PythonEnv, "")
	dir := t.TempDir()
	python := func() string {
		t.Helper()
		out, err := runCLI(t, dir, "--json", "roles", "list")
		if err != nil {
			t.Fatalf("roles list --json: %v\n%s", err, out)
		}
		var p map[string]any
		if err := json.Unmarshal([]byte(out), &p); err != nil {
			t.Fatalf("invalid json: %v\n%s", err, out)
		}
		return fmt.Sprint(p["python"])
	}
	if got := python(); got != "python3" {
		t.Errorf("default python = %q, want python3", got)
	}
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(`{"python": "uv run python"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := python(); got != "uv run python" {
		t.Errorf("from settings.json python = %q, want uv run python", got)
	}
	t.Setenv(settings.PythonEnv, "py -3")
	if got := python(); got != "py -3" {
		t.Errorf("JAIRA_PYTHON must win: python = %q, want py -3", got)
	}
}
