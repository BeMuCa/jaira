package mod

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// The cockpit is embedded whole — manifest, hooks module, state contract —
// and its tests are not: they run here, not where the mod is installed.
func TestBuiltinsCarryTheCockpitWithoutItsTests(t *testing.T) {
	mods, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(mods, func(m Mod) bool { return m.ID == "jaira-cockpit" })
	if i < 0 {
		t.Fatalf("jaira-cockpit not embedded: %+v", mods)
	}
	files := mods[i].Files
	for _, want := range []string{".claude-plugin/plugin.json", "hooks/hooks.json", "hooks/register.tsx", "types/index.d.ts"} {
		if !slices.Contains(files, want) {
			t.Errorf("%s missing from %v", want, files)
		}
	}
	for _, f := range files {
		if strings.HasPrefix(f, "tests/") {
			t.Errorf("test file %s would be installed", f)
		}
	}
}

// Claude Code loads the folder under the name in its manifest; the folder is
// named the same, so the two never drift.
func TestManifestNameIsTheFolderName(t *testing.T) {
	mods, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mods {
		raw, err := File(m.ID, ".claude-plugin/plugin.json")
		if err != nil {
			t.Fatal(err)
		}
		var manifest struct{ Name string }
		if err := json.Unmarshal(raw, &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.Name != m.ID {
			t.Errorf("folder %s, manifest name %q", m.ID, manifest.Name)
		}
	}
}
