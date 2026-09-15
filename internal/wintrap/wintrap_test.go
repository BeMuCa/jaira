package wintrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEachPatternFires proves the checker is not vacuous: every one of the five
// patterns has a fixture that trips exactly it. Without this the green run over
// the repository below would prove nothing — a checker that finds nothing
// anywhere is also green.
func TestEachPatternFires(t *testing.T) {
	cases := []struct {
		dir  string
		rule int
		want string // a phrase the remedy must carry, because a finding that
		// only names the problem leaves the reader where CI left them
	}{
		{"testdata/rule1", 1, "USERPROFILE"},
		{"testdata/rule2", 2, ".gitattributes"},
		{"testdata/rule3", 3, "runtime.GOOS"},
		{"testdata/rule4", 4, ".exe"},
		{"testdata/rule5", 5, "filepath.Join"},
	}
	for _, tc := range cases {
		found, err := Scan(tc.dir)
		if err != nil {
			t.Fatalf("Scan(%s): %v", tc.dir, err)
		}
		if len(found) == 0 {
			t.Errorf("Scan(%s) found nothing, want pattern %d to fire", tc.dir, tc.rule)
			continue
		}
		for _, f := range found {
			if f.Rule != tc.rule {
				t.Errorf("Scan(%s) reported rule %d, want only rule %d: %s", tc.dir, f.Rule, tc.rule, f)
			}
			if f.Problem == "" {
				t.Errorf("Scan(%s) finding has no problem text: %+v", tc.dir, f)
			}
			if !strings.Contains(f.Fix, tc.want) {
				t.Errorf("Scan(%s) remedy %q does not name %q", tc.dir, f.Fix, tc.want)
			}
		}
	}
}

// TestCoveredEmbedIsSilent checks the other half of pattern 2: a .gitattributes
// line that does cover the embedded files silences it. A matcher that never
// matches would make every embed a finding and the rule useless.
func TestCoveredEmbedIsSilent(t *testing.T) {
	found, err := Scan("testdata/clean")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Errorf("Scan(testdata/clean) = %v, want nothing: the embed is pinned to LF", found)
	}
}

// TestExemptionIsHonoured checks the escape hatch: a site that only looks like
// a trap carries //wintrap:ok and its reason instead of loosening the rule.
func TestExemptionIsHonoured(t *testing.T) {
	found, err := Scan("testdata/exempt")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Errorf("Scan(testdata/exempt) = %v, want nothing: the site is marked //wintrap:ok", found)
	}
}

// TestRepositoryIsClean is the one that matters day to day: it is what turns a
// red windows-latest job eight minutes in into a red Linux test in a second.
func TestRepositoryIsClean(t *testing.T) {
	found, err := Scan(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range found {
		t.Errorf("\n%s", f)
	}
}

// repoRoot climbs to the directory holding go.mod, so the scan covers the whole
// module however the test is invoked.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test's working directory")
		}
		dir = parent
	}
}

func TestMatchPattern(t *testing.T) {
	cases := []struct {
		pattern, path string
		want          bool
	}{
		{"core/lane/builtin/*.md", "core/lane/builtin/todo.md", true},
		{"core/lane/builtin/*.md", "core/lane/builtin/sub/todo.md", false},
		{"core/role/builtin/**/*.md", "core/role/builtin/jaira-teamlead/SKILL.md", true},
		{"core/role/builtin/**/*.md", "core/role/builtin/SKILL.md", true},
		{"core/role/builtin/**/*.md", "core/lane/builtin/todo.md", false},
		{"*.md", "core/release/NOTES.md", true},
		{"/NOTES.md", "NOTES.md", true},
	}
	for _, tc := range cases {
		if got := matchPattern(tc.pattern, tc.path); got != tc.want {
			t.Errorf("matchPattern(%q, %q) = %v, want %v", tc.pattern, tc.path, got, tc.want)
		}
	}
}

// TestNonTextEmbedIsNotExempt is the other half of what .gitattributes can say.
// Rule 2 used to skip embedded files by extension, so a .png was never a
// finding and a -text or binary line was never read; now the file is a finding
// like any other, and only a line in .gitattributes silences it. Without these
// two fixtures the "-text"/"binary" arm of loadAttributes could be deleted and
// every test would stay green.
func TestNonTextEmbedIsNotExempt(t *testing.T) {
	found, err := Scan("testdata/rule2bin")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || found[0].Rule != 2 {
		t.Fatalf("Scan(testdata/rule2bin) = %v, want one rule 2 finding: the embedded .png is pinned by nothing", found)
	}
	// A non-text file told to carry "text eol=lf" would be the wrong remedy, so
	// the finding has to offer the other one too.
	if !strings.Contains(found[0].Fix, `"assets/*.png binary"`) {
		t.Errorf("remedy %q does not offer the binary attribute for a non-text embed", found[0].Fix)
	}

	found, err = Scan("testdata/cleanbinary")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Errorf("Scan(testdata/cleanbinary) = %v, want nothing: binary and -text pin the bytes as firmly as eol=lf", found)
	}
}
