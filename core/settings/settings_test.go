package settings_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/settings"
)

// Settings nobody has ever opened must not change how anything behaves: the
// remote is origin and the notification the ticket asked for is on.
func TestDefaultsWithNoFile(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())

	s := settings.Load()
	if got := s.RemoteFor(t.TempDir()); got != "origin" {
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
	if name := got.RemoteFor(t.TempDir()); name != "board" {
		t.Errorf("remote is %q", name)
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
	if s.RemoteFor(t.TempDir()) != "origin" || !s.NotifyEnabled() {
		t.Errorf("a malformed file changed behaviour: %+v", s)
	}
}

func TestBlankRemoteIsStillOrigin(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{Remote: "   "}); err != nil {
		t.Fatal(err)
	}
	if got := settings.Load().RemoteFor(t.TempDir()); got != "origin" {
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

// logbook-days has three answers that must stay apart: missing means four
// weeks, 0 means none, and a number means that many days — and 0 survives a
// save, which an int with omitempty would have dropped back to the default.
func TestLogbookWindow(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)

	if got := settings.Load().LogbookWindow(); got != settings.DefaultLogbookDays {
		t.Errorf("no file: window = %d, want %d", got, settings.DefaultLogbookDays)
	}
	for _, tc := range []struct {
		raw  string
		want int
	}{
		{`{}`, settings.DefaultLogbookDays},
		{`{"logbook-days": 0}`, 0},
		{`{"logbook-days": 7}`, 7},
		{`{"logbook-days": -3}`, settings.DefaultLogbookDays},
	} {
		if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(tc.raw), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := settings.Load().LogbookWindow(); got != tc.want {
			t.Errorf("%s: window = %d, want %d", tc.raw, got, tc.want)
		}
	}

	zero := 0
	if err := settings.Save(settings.Settings{LogbookDays: &zero}); err != nil {
		t.Fatal(err)
	}
	if got := settings.Load().LogbookWindow(); got != 0 {
		t.Errorf("a saved 0 came back as %d", got)
	}
}

// The worker mode is the default: nothing written, or anything but "lanes",
// starts one agent through every lane. Only an explicit "lanes" keeps the
// one-agent-per-lane dispatcher.
func TestDispatchDefaultsToWorker(t *testing.T) {
	for in, want := range map[string]string{
		"": settings.DispatchWorker, "workr": settings.DispatchWorker,
		"worker": settings.DispatchWorker, " worker ": settings.DispatchWorker,
		"lanes": settings.DispatchLanes, " lanes ": settings.DispatchLanes,
	} {
		s := settings.Settings{Dispatch: in}
		if got := s.DispatchMode(); got != want {
			t.Errorf("dispatch %q: mode %q, want %q", in, got, want)
		}
	}
	if r := (settings.Settings{}).DispatchRole(); r != "jaira-worker" {
		t.Errorf("default role %q, want jaira-worker", r)
	}
	if r := (settings.Settings{Dispatch: "lanes"}).DispatchRole(); r != "jaira-dispatcher" {
		t.Errorf("lanes role %q, want jaira-dispatcher", r)
	}
	if r := (settings.Settings{Dispatch: "worker"}).DispatchRole(); r != "jaira-worker" {
		t.Errorf("worker role %q, want jaira-worker", r)
	}
}

// The chat model goes into a command line, so only a model name passes; no
// setting means no --model, the account default.
func TestChatModelName(t *testing.T) {
	for in, want := range map[string]string{
		"": "", "  ": "", "sonnet": "sonnet", " haiku ": "haiku",
		"claude-sonnet-5-5": "claude-sonnet-5-5", "opus[1m]": "opus[1m]",
		"sonnet; rm -rf ~": "", "it's": "", "$(id)": "", "a b": "",
	} {
		if got, _ := (settings.Settings{ChatModel: in}).ChatModelFor(""); got != want {
			t.Errorf("chat-model %q: %q, want %q", in, got, want)
		}
	}
}

// The board's own chat model wins over the machine's, and removing it lets the
// machine's apply again; a board file holding no model name counts as none.
func TestChatModelForPrefersTheBoard(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	s := settings.Settings{ChatModel: "haiku"}
	if got, src := s.ChatModelFor(root); got != "haiku" || !strings.Contains(src, "settings.json") {
		t.Errorf("no board file: %q (%s), want haiku from settings.json", got, src)
	}
	if err := settings.SetBoardChatModel(root, "sonnet"); err != nil {
		t.Fatal(err)
	}
	if got, src := s.ChatModelFor(root); got != "sonnet" || !strings.Contains(src, ".jaira/chat-model") {
		t.Errorf("board file: %q (%s), want sonnet from .jaira/chat-model", got, src)
	}
	if got, _ := (settings.Settings{}).ChatModelFor(root); got != "sonnet" {
		t.Errorf("board file, no machine setting: %q, want sonnet", got)
	}
	if err := os.WriteFile(settings.BoardChatModelPath(root), []byte("# comment\n\n  opus  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := settings.BoardChatModel(root); got != "opus" {
		t.Errorf("hand-written file: %q, want opus", got)
	}
	if err := os.WriteFile(settings.BoardChatModelPath(root), []byte("$(id)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.ChatModelFor(root); got != "haiku" {
		t.Errorf("a board file with no model name: %q, want the machine's haiku", got)
	}
	if err := settings.SetBoardChatModel(root, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(settings.BoardChatModelPath(root)); !os.IsNotExist(err) {
		t.Errorf("setting the board's model to empty must remove the file: %v", err)
	}
	if got, src := (settings.Settings{}).ChatModelFor(root); got != "" || src != "nothing configured" {
		t.Errorf("nothing set: %q (%s), want the account default", got, src)
	}
}

// Who touches the version and the changelog is the board's: .jaira/release.
// Without it, or with anything but teamlead in it, a worker does what it did
// before there was a choice.
func TestReleaseDefaultsToTheWorker(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, src := settings.ReleaseFor(root); got != settings.ReleaseWorker || src != "the default" {
		t.Errorf("no file: %q (%s), want worker, the default", got, src)
	}
	if got, _ := settings.ReleaseFor(""); got != settings.ReleaseWorker {
		t.Errorf("outside a board: %q, want worker", got)
	}
	for in, want := range map[string]string{
		"teamlead\n":                   settings.ReleaseTeamlead,
		"# who bumps\n\n  teamlead \n": settings.ReleaseTeamlead,
		"worker\n":                     settings.ReleaseWorker,
		"team-lead\n":                  settings.ReleaseWorker,
		"worker\nteamlead\n":           settings.ReleaseWorker,
	} {
		if err := os.WriteFile(settings.BoardReleasePath(root), []byte(in), 0o644); err != nil {
			t.Fatal(err)
		}
		if got, _ := settings.ReleaseFor(root); got != want {
			t.Errorf(".jaira/release %q: %q, want %q", in, got, want)
		}
	}
	if err := os.WriteFile(settings.BoardReleasePath(root), []byte("worker\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, src := settings.ReleaseFor(root); src != "the default" {
		t.Errorf("a file naming worker: source %q, want the default", src)
	}
}
