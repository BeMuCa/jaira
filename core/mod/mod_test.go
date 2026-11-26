package mod

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
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
		if strings.HasPrefix(f, "tests/") || strings.HasPrefix(f, ".claude-plugin/types/") || f == "tsconfig.json" {
			t.Errorf("%s would be installed: tests and what Claude Code generates stay out", f)
		}
	}
}

// A test run in the mod folder leaves vitest's cache in node_modules/; a
// build made after it embeds the folder, and the cache must not be installed.
func TestBuiltinsLeaveATestRunsCacheOut(t *testing.T) {
	mods, err := builtinsIn(fstest.MapFS{
		"builtin/m/.claude-plugin/plugin.json":                 {Data: []byte(`{"name":"m"}`)},
		"builtin/m/hooks/register.tsx":                         {},
		"builtin/m/node_modules/.vite/vitest/results.json":     {},
		"builtin/m/hooks/node_modules/kept-only-at-the-top.ts": {},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".claude-plugin/plugin.json", "hooks/node_modules/kept-only-at-the-top.ts", "hooks/register.tsx"}
	if len(mods) != 1 || !slices.Equal(mods[0].Files, want) {
		t.Fatalf("got %+v, want files %v", mods, want)
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
