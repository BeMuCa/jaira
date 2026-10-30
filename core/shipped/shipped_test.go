package shipped

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var entry = regexp.MustCompile(`^[0-9a-f]{64}  \S+  (v\d+\.\d+\.\d+|-)$`)

// entries returns the manifest's real lines.
func entries() []string {
	var out []string
	for _, l := range strings.Split(manifest, "\n") {
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		out = append(out, l)
	}
	return out
}

func TestManifestLinesAreWellFormed(t *testing.T) {
	es := entries()
	if len(es) == 0 {
		t.Fatal("the manifest has no entries")
	}
	seen := map[string]bool{}
	for _, l := range es {
		if !entry.MatchString(l) {
			t.Errorf("malformed manifest line %q: want 64 lowercase hex, two spaces, a path", l)
		}
		if seen[l] {
			t.Errorf("duplicate manifest line %q", l)
		}
		seen[l] = true
	}
}

func TestKnownMatchesPathAndBytes(t *testing.T) {
	const path = "core/lane/builtin/00-backlog.md"
	b, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	if !Known(path, b) {
		t.Errorf("the current %s is not in the manifest; run scripts/shipped-history.sh > core/shipped/shipped.sha256", path)
	}
	if Known(path, append(append([]byte{}, b...), '\n')) {
		t.Error("a file with one more byte must not be known")
	}
	if Known("core/lane/builtin/10-todo.md", b) {
		t.Error("the same bytes under another path must not be known")
	}
	if got := Line(path, b); !strings.HasSuffix(got, "  "+path) || len(got) != 64+2+len(path) {
		t.Errorf("Line = %q, want <64 hex>  %s", got, path)
	}
}

func TestReleaseReadsTheReleaseColumn(t *testing.T) {
	const path = "core/lane/builtin/00-backlog.md"
	b, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	rel, ok := Release(path, b)
	if !ok || (rel != "-" && !strings.HasPrefix(rel, "v")) {
		t.Errorf("Release(%s) = %q, %v; want a vX.Y.Z or -", path, rel, ok)
	}
	if _, ok := Release(path, []byte("not a shipped version")); ok {
		t.Error("bytes that were never shipped must have no release")
	}
}

func TestEveryManifestPathExists(t *testing.T) {
	seen := map[string]bool{}
	for _, l := range entries() {
		path := strings.Fields(l)[1]
		if seen[path] {
			continue
		}
		seen[path] = true
		if _, err := os.Stat(filepath.Join("..", "..", filepath.FromSlash(path))); err != nil {
			t.Errorf("manifest names %s, which does not exist in the tree: %v", path, err)
		}
	}
}
