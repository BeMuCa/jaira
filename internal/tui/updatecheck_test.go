package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// setReleaseCurrent points release.Current at v for the duration of the
// test, mirroring internal/cli/update_test.go's setCurrent.
func setReleaseCurrent(t *testing.T, v string) {
	t.Helper()
	orig := release.Current
	release.Current = v
	t.Cleanup(func() { release.Current = orig })
}

func TestVersionLineNeverCheckedShowsVersionAlone(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "") // re-enable; TestMain disables by default

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") {
		t.Errorf("versionLine() = %q, want it to name the running version", line)
	}
	if strings.Contains(line, "up to date") || strings.Contains(line, "available") {
		t.Errorf("versionLine() = %q, must not claim a checked state when the cache has never been written", line)
	}
}

func TestVersionLineUpToDate(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.0.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") || !strings.Contains(line, "up to date") {
		t.Errorf("versionLine() = %q, want it to say up to date", line)
	}
}

func TestVersionLineUpdateAvailable(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "1.3.0") || !strings.Contains(line, "jaira self upgrade") {
		t.Errorf("versionLine() = %q, want it to name the available release and 'jaira self upgrade'", line)
	}
}

// TestVersionLineDisabledShowsVersionAlone asserts JAIRA_NO_UPDATE_CHECK=1
// suppresses the "available" half even when a stale cache would otherwise
// have something to say, and reports the version alone rather than nothing.
func TestVersionLineDisabledShowsVersionAlone(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "1")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") {
		t.Errorf("versionLine() = %q, want it to name the running version", line)
	}
	if strings.Contains(line, "1.3.0") || strings.Contains(line, "available") {
		t.Errorf("versionLine() = %q, JAIRA_NO_UPDATE_CHECK=1 must suppress the available half", line)
	}
}

// TestHomeFooterCarriesTheVersionIndicator asserts the launcher's footer
// includes the same indicator versionLine() produces.
func TestHomeFooterCarriesTheVersionIndicator(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.0.0"}); err != nil {
		t.Fatal(err)
	}

	h, err := NewHome(nil)
	if err != nil {
		t.Fatal(err)
	}
	h.width, h.height = 100, 30
	out := h.render()
	if !strings.Contains(out, "up to date") {
		t.Errorf("home footer = %q, want the version indicator", out)
	}
}

// TestBoardStatusBarCarriesTheVersionIndicator asserts the board's status
// bar includes the same indicator, on the plain (non-overlay) hints path.
//
// The cache is written only after newTestStore has run, because that helper
// sets its own isolated JAIRA_HOME — writing it first would target a
// directory the test's actual Model never reads from.
func TestBoardStatusBarCarriesTheVersionIndicator(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	s := newTestStore(t)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	out := m.statusBar(0, 0)
	if !strings.Contains(out, "1.3.0") || !strings.Contains(out, "jaira self upgrade") {
		t.Errorf("status bar = %q, want the version indicator naming the available release", out)
	}
}
