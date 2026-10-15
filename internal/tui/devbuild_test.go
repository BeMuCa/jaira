package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/selfupdate"
)

// A source build is what every contributor runs, and release.Current is "dev"
// there. release.Current can only ever compare as "different" from a
// published version, so the line must never advertise an upgrade — to code
// *older* than the code being run — pointing at a command that then refuses
// with dev_build.
//
// Since P1AE82 it does name the build. DNAEPN had it say nothing at all,
// because "jaira dev" in a footer was noise; in the top left corner it is the
// whole point, since "which binary is this" is exactly what a contributor
// switching between 'self upgrade', 'go build' and ~/.local/bin is asking.
func TestVersionLineNamesADevBuildWithoutAnUpgradeClaim(t *testing.T) {
	setReleaseCurrent(t, "dev")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	// A cache dated far ahead, so staleness cannot make this test reach the
	// network, and so a cache that does hold a published release cannot leak
	// into a dev build's line.
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC().AddDate(100, 0, 0), Latest: "0.1.0"}); err != nil {
		t.Fatal(err)
	}

	line, pill := versionLine()
	if pill != "" {
		t.Errorf("pill = %q on a dev build, want none", pill)
	}
	if !strings.Contains(line, "jaira dev") {
		t.Errorf("versionLine() = %q on a dev build, want it to name the build", line)
	}
	for _, claim := range []string{"up to date", "available", "0.1.0"} {
		if strings.Contains(line, claim) {
			t.Errorf("versionLine() = %q on a dev build, must not claim %q", line, claim)
		}
	}
}

// TestHomeHeadNamesADevBuild asserts the launcher names a dev build in its top
// left corner and nowhere else — in particular not back in the footer it was
// moved out of.
func TestHomeHeadNamesADevBuild(t *testing.T) {
	setReleaseCurrent(t, "dev")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")

	h, err := NewHome(nil)
	if err != nil {
		t.Fatal(err)
	}
	h.width, h.height = 100, 30
	out := h.render()
	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[0], "jaira dev") {
		t.Errorf("home first line = %q, want it to name the dev build", lines[0])
	}
	if n := strings.Count(out, "jaira dev"); n != 1 {
		t.Errorf("home says %q %d times, want exactly once — the top left corner", "jaira dev", n)
	}
	// Anchor: the hint line must still be there, so a version found only at the
	// top is the line having moved and not the footer having vanished.
	if !strings.Contains(out, "q quit") {
		t.Errorf("home = %q, want the hint line still present", out)
	}
}

// TestBoardHeadNamesADevBuild is the board analogue of
// TestHomeHeadNamesADevBuild.
func TestBoardHeadNamesADevBuild(t *testing.T) {
	setReleaseCurrent(t, "dev")
	s := newTestStore(t)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	first := strings.Split(m.renderBoard(), "\n")[0]
	if !strings.Contains(first, "jaira dev") {
		t.Errorf("board first line = %q, want it to name the dev build", first)
	}
	sb := m.statusBar()
	if strings.Contains(sb, "jaira dev") {
		t.Errorf("status bar = %q, the version moved to the top left and must not be here", sb)
	}
	// Anchor: the help hint must still be there, so a bar without a version is
	// the version having moved and not the whole bar having vanished.
	if !strings.Contains(sb, "? help") {
		t.Errorf("status bar = %q, want the help hint still present", sb)
	}
}
