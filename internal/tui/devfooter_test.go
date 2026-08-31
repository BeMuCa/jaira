package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/selfupdate"
)

// A source build is what every contributor runs, and release.Current is "dev"
// there. release.Current can only ever compare as "different" from a
// published version, so the line used to advertise an upgrade — to code
// older than the code being run — and point at a command that then refuses
// with dev_build. It now says nothing at all, and neither caller draws a
// line when there is nothing to say.
func TestVersionLineSaysNothingOnADevBuild(t *testing.T) {
	setReleaseCurrent(t, "dev")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	// A cache dated far ahead, so staleness cannot make this test reach the
	// network, and so a cache existing at all does not change the outcome.
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC().AddDate(100, 0, 0), Latest: "0.1.0"}); err != nil {
		t.Fatal(err)
	}

	if got := versionLine(); got != "" {
		t.Errorf("versionLine() = %q on a dev build, want empty", got)
	}
}

// TestHomeFooterOmitsVersionLineOnDevBuild asserts the launcher draws no
// footer line at all for a dev build, where the release test in
// updatecheck_test.go (TestHomeFooterCarriesTheVersionIndicator) asserts one
// is drawn for a release build.
func TestHomeFooterOmitsVersionLineOnDevBuild(t *testing.T) {
	setReleaseCurrent(t, "dev")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")

	h, err := NewHome(nil)
	if err != nil {
		t.Fatal(err)
	}
	h.width, h.height = 100, 30
	out := h.render()
	if strings.Contains(out, "jaira dev") {
		t.Errorf("home footer = %q, dev build must not show a version line", out)
	}
	// Anchor: the hint line must still be there, so a vanished "jaira dev"
	// is the version line being omitted and not the whole footer.
	if !strings.Contains(out, "q quit") {
		t.Errorf("home footer = %q, want the hint line still present", out)
	}
}

// TestBoardStatusBarOmitsVersionLineOnDevBuild is the board-footer analogue
// of TestHomeFooterOmitsVersionLineOnDevBuild, mirroring the release-build
// coverage in TestBoardStatusBarCarriesTheVersionIndicator.
func TestBoardStatusBarOmitsVersionLineOnDevBuild(t *testing.T) {
	setReleaseCurrent(t, "dev")
	s := newTestStore(t)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	out := m.statusBar()
	if strings.Contains(out, "jaira dev") {
		t.Errorf("status bar = %q, dev build must not show a version line", out)
	}
	// Anchor: the help hint must still be there, so a vanished "jaira dev"
	// is the version line being omitted and not the whole status bar.
	if !strings.Contains(out, "? help") {
		t.Errorf("status bar = %q, want the help hint still present", out)
	}
}
