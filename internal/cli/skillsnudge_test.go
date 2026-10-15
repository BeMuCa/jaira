package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// skillsHome gives the test its own HOME and JAIRA_HOME (the package shares one by
// default), installs the roles into the personal claude skills directory at version
// installed and returns that directory plus the path of one skill file in it.
func skillsHome(t *testing.T, installed string) (skillsDir, skillFile string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "jaira"))
	setCurrent(t, installed)
	if out, err := runCLI(t, t.TempDir(), "roles", "install", "--global"); err != nil {
		t.Fatalf("roles install --global: %v\n%s", err, out)
	}
	skillsDir = filepath.Join(home, ".claude", "skills")
	return skillsDir, filepath.Join(skillsDir, "jaira-role-acceptance", "SKILL.md")
}

// olderShipped is the v0.4.0 file jaira once shipped: a file 'update --roles' replaces.
func olderShipped(t *testing.T) []byte {
	t.Helper()
	return mustRead(t, filepath.Join("..", "..", "core", "role", "testdata", "jaira-role-acceptance-SKILL-v0.4.0.md"))
}

func TestSkillsNudgeWarnsOnceAboutOlderPersonalSkills(t *testing.T) {
	dir, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))

	stdout, stderr := captureStdio(t, nudgeSkillsIfStale)
	if stdout != "" {
		t.Errorf("stdout = %q, want nothing", stdout)
	}
	if n := strings.Count(stderr, "\n"); n != 1 || !strings.Contains(stderr, "jaira update --roles") || !strings.Contains(stderr, dir) {
		t.Errorf("stderr = %q, want one line naming %s and 'jaira update --roles'", stderr, dir)
	}
	stdout, stderr = captureStdio(t, nudgeSkillsIfStale)
	if stdout != "" || stderr != "" {
		t.Errorf("second call: stdout=%q stderr=%q, want silence", stdout, stderr)
	}
}

func TestSkillsNudgeIgnoresAnEditedSkill(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, []byte("my own words, never shipped\n"))

	stdout, stderr := captureStdio(t, nudgeSkillsIfStale)
	if stdout != "" || stderr != "" {
		t.Errorf("an edited skill warned: stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestSkillsNudgeChecksAgainForANewVersion(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))
	if _, stderr := captureStdio(t, nudgeSkillsIfStale); stderr == "" {
		t.Fatal("1.3.0 did not warn")
	}
	setCurrent(t, "1.4.0")
	if _, stderr := captureStdio(t, nudgeSkillsIfStale); !strings.Contains(stderr, "jaira update --roles") {
		t.Errorf("1.4.0 stderr = %q, want a new warning", stderr)
	}
}

func TestSkillsNudgeSilentOnDevBuild(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))
	setCurrent(t, "dev")

	stdout, stderr := captureStdio(t, nudgeSkillsIfStale)
	if stdout != "" || stderr != "" {
		t.Errorf("dev build warned: stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestSkillsNudgeSilentDuringUpdate(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))
	dir, _ := updateTestStore(t)
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "jaira"))

	var out string
	var err error
	_, stderr := captureStdio(t, func() { out, err = runCLI(t, dir, "update", "--roles") })
	if err != nil {
		t.Fatalf("update --roles: %v\n%s", err, out)
	}
	if strings.Contains(stderr+out, "older than this jaira") {
		t.Errorf("update warned about what it is fixing: %q %q", stderr, out)
	}
}

func TestSkillsNudgeSilentDuringRolesInstall(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))
	dir, _ := updateTestStore(t)
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "jaira"))

	_, stderr := captureStdio(t, func() { runCLI(t, dir, "roles", "install", "--project") })
	if strings.Contains(stderr, "older than this jaira") {
		t.Errorf("roles install --project warned: %q", stderr)
	}
}

func TestSkillsNudgeNeverOnStdoutUnderJSON(t *testing.T) {
	_, file := skillsHome(t, "1.3.0")
	mustWrite(t, file, olderShipped(t))
	dir, _ := updateTestStore(t)
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "jaira"))

	var out string
	var err error
	stdout, stderr := captureStdio(t, func() { out, err = runCLI(t, dir, "list", "--json") })
	if err != nil {
		t.Fatalf("list --json: %v\n%s", err, out)
	}
	if !strings.Contains(stderr, "older than this jaira") {
		t.Errorf("stderr = %q, want the warning", stderr)
	}
	if strings.Contains(stdout+out, "older than this jaira") {
		t.Errorf("the warning reached stdout: %q %q", stdout, out)
	}
	var v any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("list --json output is not JSON: %v\n%s", err, out)
	}
}

func TestSkillsNudgeWithoutAHomeStampsNothing(t *testing.T) {
	setCurrent(t, "1.3.0")
	t.Setenv("JAIRA_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	cwd := t.TempDir()
	t.Chdir(cwd)

	captureStdio(t, nudgeSkillsIfStale)
	if _, err := os.Stat(filepath.Join(cwd, "skills-check")); err == nil {
		t.Error("stamped into the working directory")
	}
}
