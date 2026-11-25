package settings_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
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

// Who opens a pull request is the board's: .jaira/pr. Without it, or with
// anything but agent in it, an agent pushes and hands the create line back —
// a typo must never be what lets an agent open one (JG3REW).
func TestPRDefaultsToThePerson(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, src := settings.PRFor(root); got != settings.PRPerson || src != "the default" {
		t.Errorf("no file: %q (%s), want person, the default", got, src)
	}
	if got, _ := settings.PRFor(""); got != settings.PRPerson {
		t.Errorf("outside a board: %q, want person", got)
	}
	for in, want := range map[string]string{
		"":                            settings.PRPerson,
		"agent\n":                     settings.PRAgent,
		"agent\r\n":                   settings.PRAgent,
		"# agents open\n\n  agent \n": settings.PRAgent,
		"person\n":                    settings.PRPerson,
		"Agent\n":                     settings.PRPerson,
		"agents\n":                    settings.PRPerson,
		"person\nagent\n":             settings.PRPerson,
		"agent; rm -rf /\n":           settings.PRPerson,
		"none\n":                      settings.PRNone,
		"none\r\n":                    settings.PRNone,
		"# no PRs\n\n\t none  \n":     settings.PRNone,
		"None\n":                      settings.PRPerson,
		"nnone\n":                     settings.PRPerson,
		"none agent\n":                settings.PRPerson,
		"person\nnone\n":              settings.PRPerson,
	} {
		if err := os.WriteFile(settings.BoardPRPath(root), []byte(in), 0o644); err != nil {
			t.Fatal(err)
		}
		if got, _ := settings.PRFor(root); got != want {
			t.Errorf(".jaira/pr %q: %q, want %q", in, got, want)
		}
		// The gate reads the same content out of the base commit.
		if got := settings.PRIn(in); got != want {
			t.Errorf("PRIn(%q): %q, want %q", in, got, want)
		}
	}
	if err := os.WriteFile(settings.BoardPRPath(root), []byte("agent\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, src := settings.PRFor(root); src != "from .jaira/pr, shared with the board" {
		t.Errorf("a file naming agent: source %q", src)
	}
}

// pr-when opens a pull request before acceptance only when the board says
// exactly signoff; anything else — no file, a typo, a second line — keeps the
// order this repository runs: done, then the pull request (JG3REW).
func TestPRWhenDefaultsToDone(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, src := settings.PRWhenFor(root); got != settings.PRWhenDone || src != "the default" {
		t.Errorf("no file: %q (%s), want done, the default", got, src)
	}
	if got, _ := settings.PRWhenFor(""); got != settings.PRWhenDone {
		t.Errorf("outside a board: %q, want done", got)
	}
	for in, want := range map[string]string{
		"":                                  settings.PRWhenDone,
		"signoff\n":                         settings.PRWhenSignoff,
		"signoff\r\n":                       settings.PRWhenSignoff,
		"# review on the PR\n\n signoff \n": settings.PRWhenSignoff,
		"done\n":                            settings.PRWhenDone,
		"Signoff\n":                         settings.PRWhenDone,
		"sign-off\n":                        settings.PRWhenDone,
		"done\nsignoff\n":                   settings.PRWhenDone,
		"signoff; rm -rf /\n":               settings.PRWhenDone,
		"\xef\xbb\xbfsignoff\r\n":           settings.PRWhenSignoff,
	} {
		if err := os.WriteFile(settings.BoardPRWhenPath(root), []byte(in), 0o644); err != nil {
			t.Fatal(err)
		}
		if got, _ := settings.PRWhenFor(root); got != want {
			t.Errorf(".jaira/pr-when %q: %q, want %q", in, got, want)
		}
		if got := settings.PRWhenIn(in); got != want {
			t.Errorf("PRWhenIn(%q): %q, want %q", in, got, want)
		}
	}
}

// The reviewer goes onto a gh or glab command line, so only a login survives.
func TestPRReviewerIsALoginOrNothing(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := settings.PRReviewerFor(root); got != "" {
		t.Errorf("no file: %q, want empty", got)
	}
	if got := settings.PRReviewerFor(""); got != "" {
		t.Errorf("outside a board: %q, want empty", got)
	}
	for in, want := range map[string]string{
		"":                       "",
		"copilot\n":              "copilot",
		"copilot\r\n":            "copilot",
		"# who reviews\ncopilot": "copilot",
		"monalisa\n":             "monalisa",
		"myorg/team-name\n":      "myorg/team-name",
		"@copilot\n":             "@copilot",
		"copilot $(id)\n":        "",
		"copilot;rm\n":           "",
		"a b\n":                  "",
		"\"copilot\"\n":          "",
		"кто-то\n":               "",
	} {
		if err := os.WriteFile(settings.BoardPRReviewerPath(root), []byte(in), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := settings.PRReviewerFor(root); got != want {
			t.Errorf(".jaira/pr-reviewer %q: %q, want %q", in, got, want)
		}
	}
}

// A .jaira/pr that is there but reads as person says so in its source, so a
// typo or a stray word is visible in 'jaira roles list' instead of looking
// like no file at all.
func TestPRSourceSaysWhenTheFileDidNotTake(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	for in, want := range map[string]string{
		"nnone\n":              `the default — .jaira/pr says "nnone", which is neither agent nor none`,
		"# only a comment\n":   "the default — .jaira/pr is there but says nothing",
		"person\n":             "from .jaira/pr, shared with the board",
		"\xef\xbb\xbfnone\r\n": "from .jaira/pr, shared with the board",
	} {
		if err := os.WriteFile(settings.BoardPRPath(root), []byte(in), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, src := settings.PRFor(root); src != want {
			t.Errorf(".jaira/pr %q: source %q, want %q", in, src, want)
		}
	}
}

// TestUpdateNeverWipesTheFile: a screen's write goes through Update. Over a
// file that does not parse it writes nothing and says where the file is; keys
// this build does not know stay; a field emptied by the change leaves the file;
// a "language" that is no string reads as none and costs nothing beside it.
func TestUpdateNeverWipesTheFile(t *testing.T) {
	cases := []struct {
		name, file string
		wantErr    bool
		want       map[string]any
	}{
		{"trailing comma", `{"remote":"upstream","hook":"/bin/notify",}`, true, nil},
		{"not json", "remote = upstream", true, nil},
		{"a non-string language", `{"remote":"upstream","hook":"/bin/notify","language":5}`, false,
			map[string]any{"remote": "upstream", "hook": "/bin/notify", "language": "Deutsch"}},
		{"an unknown key", `{"remote":"upstream","future-key":{"a":[1,2]}}`, false,
			map[string]any{"remote": "upstream", "future-key": map[string]any{"a": []any{1.0, 2.0}}, "language": "Deutsch"}},
		{"empty file", "", false, map[string]any{"language": "Deutsch"}},
		{"no file", "-", false, map[string]any{"language": "Deutsch"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("JAIRA_HOME", home)
			path := filepath.Join(home, "settings.json")
			if c.file != "-" {
				if err := os.WriteFile(path, []byte(c.file), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			_, err := settings.Update(func(s *settings.Settings) { s.Language = "Deutsch" })
			b, _ := os.ReadFile(path)
			if c.wantErr {
				if err == nil || !strings.Contains(err.Error(), path) {
					t.Fatalf("want an error naming %s, got %v", path, err)
				}
				if string(b) != c.file {
					t.Fatalf("the file was written over: %q", b)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("file %v, want %v", got, c.want)
			}
		})
	}

	// A field the change empties leaves the file; the rest stays.
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(`{"language":"Русский","dispatch":"worker","x":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := settings.Update(func(s *settings.Settings) { s.Language = "" }); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(home, "settings.json"))
	var got map[string]any
	_ = json.Unmarshal(b, &got)
	if want := map[string]any{"dispatch": "worker", "x": 1.0}; !reflect.DeepEqual(got, want) {
		t.Fatalf("after emptying the language: %v, want %v", got, want)
	}
	// And Load reads a non-string language as none without losing the rest.
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(`{"remote":"upstream","language":["x"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if s := settings.Load(); s.Remote != "upstream" || s.Language != "" {
		t.Fatalf("Load: %+v, want remote upstream and no language", s)
	}
}

// TestUpdateTwoWritersAtOnce: two settings screens writing at the same
// moment take turns. Every write lands, and nothing beside them is lost.
func TestUpdateTwoWritersAtOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	path := filepath.Join(home, "settings.json")
	if err := os.WriteFile(path, []byte(`{"remote":"upstream","hook":"/h","zz-unknown":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range 50 {
				_, err := settings.Update(func(s *settings.Settings) {
					if i == 0 {
						s.Language = settings.Language(fmt.Sprint("L", n))
					} else {
						s.Python = fmt.Sprint("py", n)
					}
				})
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
	b, _ := os.ReadFile(path)
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("settings.json does not parse after two writers: %v\n%s", err, b)
	}
	want := map[string]any{"remote": "upstream", "hook": "/h", "zz-unknown": 1.0, "language": "L49", "python": "py49"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("after two writers: %v, want %v", got, want)
	}
}

// TestUpdateDropsAKnownKeyInAnyCase: encoding/json reads "Language" as the
// language, so emptying the language takes that key out too.
func TestUpdateDropsAKnownKeyInAnyCase(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	path := filepath.Join(home, "settings.json")
	if err := os.WriteFile(path, []byte(`{"Language":"Deutsch","REMOTE":"upstream"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := settings.Update(func(s *settings.Settings) { s.Language = "" }); err != nil {
		t.Fatal(err)
	}
	if s := settings.Load(); s.Language != "" || s.Remote != "upstream" {
		t.Fatalf("after emptying: %+v, want no language and remote upstream", s)
	}
}

// Without a settings file the teamlead runs on opus, workers and cheap lanes on
// haiku at high effort, and strong lanes on opus.
func TestModelsDefault(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())

	s := settings.Load()
	want := map[string]string{
		"orchestrator": "opus",
		"worker":       "haiku high",
		"strong":       "opus",
		"cheap":        "haiku high",
	}
	for key, w := range want {
		m, src := s.ModelFor(key)
		if m.String() != w || src != "the default" {
			t.Errorf("%s: %q (%s), want %q from the default", key, m, src, w)
		}
	}
	if m, _ := s.ModelFor("no-such-tier"); m.String() != "haiku high" {
		t.Errorf("unknown key: %q, want the worker's haiku high", m)
	}
}

// An entry in settings.json wins whole: naming a worker model without an effort
// does not keep the default's high.
func TestModelsFromFileWinWhole(t *testing.T) {
	s := settings.Settings{Models: map[string]settings.AgentModel{
		"worker": {Model: "sonnet"},
		"strong": {Model: "sonnet", Effort: "xhigh"},
	}}
	if m, src := s.ModelFor("worker"); m.String() != "sonnet" || !strings.Contains(src, "settings.json") {
		t.Errorf("worker: %q (%s), want sonnet with no effort from settings.json", m, src)
	}
	if m, _ := s.ModelFor("strong"); m.String() != "sonnet xhigh" {
		t.Errorf("strong: %q, want sonnet xhigh", m)
	}
	if m, _ := s.ModelFor("cheap"); m.String() != "haiku high" {
		t.Errorf("cheap is not named, got %q, want the default", m)
	}
}

// A model that is not a name falls back to the default; an effort that is not a
// level is dropped and the model kept.
func TestModelsRejectBadValues(t *testing.T) {
	s := settings.Settings{Models: map[string]settings.AgentModel{
		"worker": {Model: "x;rm"},
		"cheap":  {Model: "sonnet", Effort: "huge"},
	}}
	if m, src := s.ModelFor("worker"); m.String() != "haiku high" || src != "the default" {
		t.Errorf("bad model: %q (%s), want the default haiku high", m, src)
	}
	if m, _ := s.ModelFor("cheap"); m.String() != "sonnet" {
		t.Errorf("bad effort: %q, want sonnet with no effort", m)
	}
}
