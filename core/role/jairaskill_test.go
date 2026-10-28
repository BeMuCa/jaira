package role

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestJairaSkillMatchesTheRepositoryCopy(t *testing.T) {
	repo, err := os.ReadFile(filepath.Join("..", "..", ".claude", "skills", "jaira", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(repo, jairaSkill) {
		t.Fatal("core/role/jaira/SKILL.md and .claude/skills/jaira/SKILL.md differ: edit one and copy it over the other — the binary installs the first, this repository's own sessions read the second")
	}
}

func TestInstallWritesTheJairaSkill(t *testing.T) {
	dir := t.TempDir()
	results, err := install(dir, false, false, false)
	if err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "jaira", "SKILL.md")
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, jairaSkill) {
		t.Error("the installed jaira skill is not the embedded one")
	}
	found := false
	for _, r := range results {
		if r.Role == "jaira" && r.Path == dst && r.Action == Written {
			found = true
		}
	}
	if !found {
		t.Errorf("no Written result for the jaira skill in %+v", results)
	}
	// an edited copy is left alone, like any role file
	if err := os.WriteFile(dst, []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	results, err = install(dir, false, false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Role == "jaira" && r.Action != Skipped {
			t.Errorf("edited jaira skill: got %s, want %s", r.Action, Skipped)
		}
	}
}

func TestInstalledDoesNotCountTheJairaSkillAlone(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "jaira", "SKILL.md"), jairaSkill, 0o644); err != nil {
		t.Fatal(err)
	}
	found, err := InstalledIn([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if found != "" {
		t.Errorf("InstalledIn = %q, want none: the jaira skill is not a role", found)
	}
}
