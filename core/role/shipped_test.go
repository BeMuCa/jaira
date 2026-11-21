package role

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/mod"
	"github.com/BeMuCa/jaira/core/shipped"
)

const acceptanceRel = "SKILL.md"
const acceptanceID = "jaira-role-acceptance"

// An older jaira's role file, as it lay on disk before the current one.
func olderAcceptance(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "jaira-role-acceptance-SKILL-v0.4.0.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !shipped.Known(roleKey(acceptanceID, acceptanceRel), b) {
		t.Fatal("the v0.4.0 fixture is not in core/shipped/shipped.sha256; regenerate it: scripts/shipped-history.sh --write")
	}
	want, err := File(acceptanceID, acceptanceRel)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) == string(want) {
		t.Fatal("the fixture equals the current file; it must be an older version")
	}
	return b
}

// Every file this binary installs is a version jaira ships, or the next release
// reports its own files as edited. The failure prints the lines to append.
func TestEveryInstalledFileIsAShippedVersion(t *testing.T) {
	var missing []string
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		for _, rel := range r.Files {
			b, err := File(r.ID, rel)
			if err != nil {
				t.Fatal(err)
			}
			if !shipped.Known(roleKey(r.ID, rel), b) {
				missing = append(missing, shipped.Line(roleKey(r.ID, rel), b))
			}
		}
	}
	if !shipped.Known(jairaSkillKey, jairaSkill) {
		missing = append(missing, shipped.Line(jairaSkillKey, jairaSkill))
	}
	mods, err := mod.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mods {
		for _, rel := range m.Files {
			b, err := mod.File(m.ID, rel)
			if err != nil {
				t.Fatal(err)
			}
			if !shipped.Known(modKey(m.ID, rel), b) {
				missing = append(missing, shipped.Line(modKey(m.ID, rel), b))
			}
		}
	}
	if len(missing) > 0 {
		t.Errorf("core/shipped/shipped.sha256 lacks the version of %d file(s) this build installs; run scripts/shipped-history.sh --write, or append:\n%s",
			len(missing), strings.Join(missing, "\n"))
	}
	// the hard-coded prefixes must name real files
	for _, key := range []string{jairaSkillKey, roleKey(acceptanceID, acceptanceRel), modKey(mods[0].ID, mods[0].Files[0])} {
		if _, err := os.Stat(filepath.Join("..", "..", filepath.FromSlash(key))); err != nil {
			t.Errorf("key %s names no file in the repository: %v", key, err)
		}
	}
}

// An older jaira's file, untouched on disk, is replaced, not skipped as edited.
// The cockpit's board.ts from v0.5.2 is the version the manifest once lost in a
// hand-merged conflict (0K8Y6S); its hash pins the fixture to what the tag
// carried: git show v0.5.2:<key> | sha256sum.
func TestInstallReplacesAnOlderShippedVersion(t *testing.T) {
	cases := []struct {
		name    string
		dst     string // under the skills directory
		key     string
		fixture string
		pin     string // sha256 of the fixture, when it stands for one tag
		want    func() ([]byte, error)
	}{
		{"acceptance role v0.4.0", filepath.Join(acceptanceID, acceptanceRel), roleKey(acceptanceID, acceptanceRel),
			"jaira-role-acceptance-SKILL-v0.4.0.md", "", func() ([]byte, error) { return File(acceptanceID, acceptanceRel) }},
		{"cockpit board.ts v0.5.2", filepath.Join("jaira-cockpit", "hooks", "board.ts"), modKey("jaira-cockpit", "hooks/board.ts"),
			"jaira-cockpit-hooks-board-v0.5.2.ts", "a06a3c1cb1b417918cf78e6995583cd472df472db843ae09cc5a67ca597cbf8f",
			func() ([]byte, error) { return mod.File("jaira-cockpit", "hooks/board.ts") }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			old, err := os.ReadFile(filepath.Join("testdata", c.fixture))
			if err != nil {
				t.Fatal(err)
			}
			if c.pin != "" && shipped.Line(c.key, old) != c.pin+"  "+c.key {
				t.Fatalf("the fixture is not the version it is named after: %s", shipped.Line(c.key, old))
			}
			want, err := c.want()
			if err != nil {
				t.Fatal(err)
			}
			if string(old) == string(want) {
				t.Fatal("the fixture equals the current file; it must be an older version")
			}

			dir := t.TempDir()
			if _, err := install(dir, false, true, false); err != nil {
				t.Fatal(err)
			}
			dst := filepath.Join(dir, c.dst)
			if err := os.WriteFile(dst, old, 0o644); err != nil {
				t.Fatal(err)
			}
			results, err := install(dir, false, true, false)
			if err != nil {
				t.Fatal(err)
			}
			for _, res := range results {
				want := Unchanged
				if res.Path == dst {
					want = Updated
				}
				if res.Action != want {
					t.Errorf("%s: got %s, want %s", res.Path, res.Action, want)
				}
			}
			if got, _ := os.ReadFile(dst); string(got) != string(want) {
				t.Error("the older shipped version was not replaced")
			}
		})
	}
}

func TestInstallDryRunWritesNothing(t *testing.T) {
	dir := t.TempDir()
	results, err := install(dir, false, true, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, res := range results {
		if res.Action != Written {
			t.Errorf("%s: got %s, want %s on an empty directory", res.Path, res.Action, Written)
		}
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("a dry run wrote %d entries into the directory", len(entries))
	}

	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	old := olderAcceptance(t)
	dst := filepath.Join(dir, acceptanceID, acceptanceRel)
	if err := os.WriteFile(dst, old, 0o644); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(dir, "jaira-role-lane", "SKILL.md")
	const mine = "my own version\n"
	if err := os.WriteFile(edited, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err = install(dir, false, true, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, res := range results {
		switch res.Path {
		case dst:
			if res.Action != Updated {
				t.Errorf("older version: got %s, want %s", res.Action, Updated)
			}
		case edited:
			if res.Action != Skipped {
				t.Errorf("edited file: got %s, want %s", res.Action, Skipped)
			}
		}
	}
	if b, _ := os.ReadFile(dst); string(b) != string(old) {
		t.Error("a dry run replaced the older version")
	}
	if b, _ := os.ReadFile(edited); string(b) != mine {
		t.Error("a dry run touched the edited file")
	}
}
