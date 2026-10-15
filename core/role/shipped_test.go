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
		t.Fatal("the v0.4.0 fixture is not in core/shipped/shipped.sha256; regenerate it: scripts/shipped-history.sh > core/shipped/shipped.sha256")
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
		t.Errorf("core/shipped/shipped.sha256 lacks the version of %d file(s) this build installs; run scripts/shipped-history.sh > core/shipped/shipped.sha256, or append:\n%s",
			len(missing), strings.Join(missing, "\n"))
	}
	// the hard-coded prefixes must name real files
	for _, key := range []string{roleKey(acceptanceID, acceptanceRel), modKey(mods[0].ID, mods[0].Files[0])} {
		if _, err := os.Stat(filepath.Join("..", "..", filepath.FromSlash(key))); err != nil {
			t.Errorf("key %s names no file in the repository: %v", key, err)
		}
	}
}

func TestInstallReplacesAnOlderShippedVersion(t *testing.T) {
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	old := olderAcceptance(t)
	dst := filepath.Join(dir, acceptanceID, acceptanceRel)
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
	if SkippedAny(results) {
		t.Error("an older shipped version was reported as edited")
	}
	got, _ := os.ReadFile(dst)
	want, _ := File(acceptanceID, acceptanceRel)
	if string(got) != string(want) {
		t.Error("the older shipped version was not replaced")
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
