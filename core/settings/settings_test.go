package settings_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/settings"
)

// Settings nobody has ever opened must not change how anything behaves: the
// remote is origin and the notification the ticket asked for is on.
func TestDefaultsWithNoFile(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())

	s := settings.Load()
	if got := s.RemoteName(); got != "origin" {
		t.Errorf("remote defaults to %q", got)
	}
	if !s.NotifyEnabled() {
		t.Error("notifications are off by default; the feature was asked for")
	}
	if s.Hook != "" {
		t.Errorf("a hook appeared out of nowhere: %q", s.Hook)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)

	want := settings.Settings{Remote: "board", NotifyOff: true, Hook: "/usr/local/bin/tell-me"}
	if err := settings.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, "settings.json")); err != nil {
		t.Fatalf("settings file: %v", err)
	}
	got := settings.Load()
	if got.Remote != want.Remote || got.NotifyOff != want.NotifyOff || got.Hook != want.Hook {
		t.Errorf("round trip changed the settings: %+v", got)
	}
	if got.RemoteName() != "board" {
		t.Errorf("remote is %q", got.RemoteName())
	}
	if got.NotifyEnabled() {
		t.Error("notify-off was not honoured")
	}
}

// A broken settings file is one person's preferences, not a reason to refuse to
// open a board.
func TestAMalformedFileFallsBackToDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	s := settings.Load()
	if s.RemoteName() != "origin" || !s.NotifyEnabled() {
		t.Errorf("a malformed file changed behaviour: %+v", s)
	}
}

func TestBlankRemoteIsStillOrigin(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{Remote: "   "}); err != nil {
		t.Fatal(err)
	}
	if got := settings.Load().RemoteName(); got != "origin" {
		t.Errorf("whitespace remote resolved to %q", got)
	}
}

// The landing branches decide when a ref may be removed, so what happens when
// nobody configured them is the important case: nothing is removable.
func TestLandingFallsBackToTheRemoteHeadAndThenToNothing(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())

	s := settings.Load()
	if got := s.Landing("origin", func() string { return "origin/master" }); len(got) != 1 || got[0] != "origin/master" {
		t.Errorf("with no config, want the remote head, got %v", got)
	}
	if got := s.Landing("origin", func() string { return "" }); len(got) != 0 {
		t.Errorf("with nothing resolvable, want no branches at all, got %v", got)
	}
	if got := s.Landing("origin", nil); len(got) != 0 {
		t.Errorf("want no branches when there is nobody to ask, got %v", got)
	}
}

// A configured list wins over the remote's head: the host's default branch is
// not necessarily the one a team cares about.
func TestConfiguredLandingBranchesWinAndAreRemotePrefixed(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{LandingBranches: []string{"main", " develop ", "", "origin/release/1.x"}}); err != nil {
		t.Fatal(err)
	}
	got := settings.Load().Landing("origin", func() string { return "origin/master" })
	want := []string{"origin/main", "origin/develop", "origin/release/1.x"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
			break
		}
	}
}

func TestIntervalDefaults(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	s := settings.Load()
	if s.SnapshotBranchName() != "jaira/board" {
		t.Errorf("snapshot branch defaults to %q", s.SnapshotBranchName())
	}
	if s.SnapshotInterval() != 72*time.Hour {
		t.Errorf("snapshot interval defaults to %v", s.SnapshotInterval())
	}
	if s.FetchInterval() != 10*time.Minute {
		t.Errorf("fetch interval defaults to %v", s.FetchInterval())
	}
	if s.LandingGraceInterval() != 3*24*time.Hour {
		t.Errorf("landing grace defaults to %v", s.LandingGraceInterval())
	}
}

// The intervals are durations, so the same field says "three seconds" for a
// screen recording and "three days" for a backup without a second field per
// unit.
func TestIntervalsAreDurations(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{
		FetchEvery: "3s", SnapshotEvery: "90m", LandingGrace: "48h",
	}); err != nil {
		t.Fatal(err)
	}
	s := settings.Load()
	if s.FetchInterval() != 3*time.Second {
		t.Errorf("fetch interval is %v", s.FetchInterval())
	}
	if s.SnapshotInterval() != 90*time.Minute {
		t.Errorf("snapshot interval is %v", s.SnapshotInterval())
	}
	if s.LandingGraceInterval() != 48*time.Hour {
		t.Errorf("landing grace is %v", s.LandingGraceInterval())
	}
}

// A value nobody can parse is one person's typo, not a reason to refuse to
// open a board: the background job keeps its usual pace and says nothing.
func TestAnUnreadableIntervalFallsBackToTheDefault(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{
		FetchEvery: "ten minutes", SnapshotEvery: "0s", LandingGrace: "-5h",
	}); err != nil {
		t.Fatal(err)
	}
	s := settings.Load()
	if s.FetchInterval() != 10*time.Minute {
		t.Errorf("nonsense fetch interval became %v", s.FetchInterval())
	}
	if s.SnapshotInterval() != 72*time.Hour {
		t.Errorf("a zero snapshot interval became %v", s.SnapshotInterval())
	}
	if s.LandingGraceInterval() != 3*24*time.Hour {
		t.Errorf("a negative grace became %v", s.LandingGraceInterval())
	}
}
