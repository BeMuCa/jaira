package lane

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The repository ships lane files under lanes/ that are NOT built-ins: they are
// catalogue lanes, adopted deliberately rather than injected into every board.
// Nothing loads them at build time, so without this test a typo in one of them
// is found by whoever adopts it. A single unquoted colon in a description is
// enough — YAML reads it as a nested mapping and the lane loses its id.
func TestShippedCatalogueLanesParseAndValidate(t *testing.T) {
	dir := filepath.Join("..", "..", "lanes")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("no lanes/ directory: %v", err)
	}

	cat := t.TempDir()
	var names []string
	for _, e := range entries {
		// README.md documents the directory. Everything else in it must be a lane,
		// skipped by name rather than by "does it look like one" so a lane file
		// with broken frontmatter cannot slip through as prose.
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "README.md" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cat, e.Name()), raw, 0o644); err != nil {
			t.Fatal(err)
		}
		names = append(names, e.Name())
	}
	if len(names) == 0 {
		t.Skip("lanes/ holds no lane files")
	}

	t.Setenv("JAIRA_LANES_DIR", cat)
	s, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("shipped catalogue lanes do not load: %v", err)
	}
	if len(s.Warnings) > 0 {
		t.Errorf("shipped catalogue lanes load with warnings:\n%s", strings.Join(s.Warnings, "\n"))
	}
	for _, n := range names {
		id := strings.TrimSuffix(n, ".md")
		l, ok := s.Get(id)
		if !ok {
			t.Errorf("lanes/%s did not resolve to lane %q; check its id: field", n, id)
			continue
		}
		if l.Builtin {
			t.Errorf("lane %q is a built-in; a file under lanes/ must not shadow one silently", id)
		}
		// These are the fields a stranger reads to decide whether to adopt it.
		if strings.TrimSpace(l.Description) == "" {
			t.Errorf("lane %q has no description", id)
		}
		if l.Agentic && strings.TrimSpace(l.Prompt) == "" {
			t.Errorf("lane %q is agentic with no prompt", id)
		}
	}
}
