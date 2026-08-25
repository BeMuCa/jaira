package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/release"
)

// A source build is what every contributor runs, and release.Current is "dev"
// there. Comparing that string to a published version can only say "different",
// so the footer used to advertise an upgrade — to code older than the code being
// run — and point at a command that refuses with dev_build.
func TestVersionLineSaysNothingOnADevBuild(t *testing.T) {
	if release.Current != "dev" {
		t.Skipf("this is a %s build, not a dev one", release.Current)
	}
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	// A cache as --check used to write it, even from a dev build. Dated far
	// ahead so staleness cannot make this test reach the network.
	if err := os.WriteFile(filepath.Join(home, "update-check.json"),
		[]byte(`{"checked_at":"2126-01-01T00:00:00Z","latest":"0.1.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got := stripANSI(versionLine())
	if strings.Contains(got, "available") || strings.Contains(got, "self upgrade") {
		t.Errorf("a dev build advertises an upgrade it will refuse: %q", got)
	}
	if strings.Contains(got, "up to date") {
		t.Errorf("a dev build claims to be up to date: %q", got)
	}
	if !strings.Contains(got, "dev") {
		t.Errorf("the footer does not say which build this is: %q", got)
	}
}
