package pr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// counting answers every state query with the state set for its number and
// counts the graphql calls.
func counting(states map[string]string, calls *int) Run {
	return func(_ string, name string, args ...string) ([]byte, error) {
		line := name + " " + strings.Join(args, " ")
		switch {
		case strings.HasPrefix(line, "git remote get-url origin"):
			return []byte("git@github.com:me/jaira.git\n"), nil
		case strings.HasPrefix(line, "git config"):
			return nil, exitErr{1, ""}
		case strings.HasPrefix(line, "gh repo view"):
			return []byte(`{"nameWithOwner":"me/jaira"}`), nil
		case strings.HasPrefix(line, "gh api graphql"):
			*calls++
			for n, s := range states {
				if strings.Contains(line, "number="+n+" ") {
					return []byte(`{"data":{"repository":{"pullRequest":{"state":"` + s + `"}}}}`), nil
				}
			}
			return nil, exitErr{1, "no such pull request"}
		}
		return nil, exitErr{1, "unexpected: " + line}
	}
}

func TestCheckerCache(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "state", CacheFile)
	now := time.Date(2026, 11, 22, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	states := map[string]string{"12": "MERGED", "13": "OPEN"}
	calls := 0
	newChecker := func() *Checker {
		return &Checker{Dir: dir, File: file, Run: counting(states, &calls), Look: found, Now: clock}
	}

	// One process asks each number once, however often a gate asks it.
	c := newChecker()
	for range 3 {
		if a := c.Ask(12); !a.Merged() {
			t.Fatalf("12: %+v, want merged", a)
		}
		if a := c.Ask(13); a.State != "open" || a.Why() != "open" {
			t.Fatalf("13: %+v, want open", a)
		}
	}
	if calls != 2 {
		t.Fatalf("%d forge calls in one process, want 2", calls)
	}

	// Another process (another worktree, same file) asks nothing while the
	// answers are fresh.
	if newChecker().Ask(12).Merged(); calls != 2 {
		t.Fatalf("merged was asked again: %d calls", calls)
	}
	if newChecker().Ask(13); calls != 2 {
		t.Fatalf("a fresh open answer was asked again: %d calls", calls)
	}

	// After OpenFor, open is asked again — and merged never is.
	now = now.Add(OpenFor + time.Second)
	states["13"] = "MERGED"
	c = newChecker()
	if !c.Ask(12).Merged() || calls != 2 {
		t.Fatalf("merged must be kept for good: %d calls", calls)
	}
	if !c.Ask(13).Merged() || calls != 3 {
		t.Fatalf("a stale open answer must be asked again: %d calls", calls)
	}

	// A failure is not cached: the next process asks again.
	c = newChecker()
	if a := c.Ask(99); a.Err == nil || !strings.Contains(a.Why(), "state unknown") {
		t.Fatalf("99: %+v, want an error", a)
	}
	if newChecker().Ask(99); calls != 5 {
		t.Fatalf("a failed answer was cached: %d calls", calls)
	}

	// No temp file is left behind.
	left, _ := filepath.Glob(filepath.Join(filepath.Dir(file), ".jaira-tmp-*"))
	if len(left) != 0 {
		t.Fatalf("temp files left: %v", left)
	}
}

// A corrupt, unreadable or unwritable cache costs a forge call and nothing
// else; offline answers only from the cache.
func TestCheckerCacheFailures(t *testing.T) {
	dir := t.TempDir()
	calls := 0
	run := counting(map[string]string{"12": "MERGED"}, &calls)

	corrupt := filepath.Join(dir, "corrupt.json")
	if err := os.WriteFile(corrupt, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !(&Checker{Dir: dir, File: corrupt, Run: run, Look: found}).Ask(12).Merged() {
		t.Fatal("a corrupt cache must fall back to the forge")
	}
	if b, _ := os.ReadFile(corrupt); !strings.Contains(string(b), `"12"`) {
		t.Fatalf("the corrupt cache was not replaced: %s", b)
	}

	// The cache path is a directory: nothing can be written there.
	blocked := filepath.Join(dir, "isdir")
	if err := os.MkdirAll(blocked, 0o755); err != nil {
		t.Fatal(err)
	}
	if !(&Checker{Dir: dir, File: blocked, Run: run, Look: found}).Ask(12).Merged() {
		t.Fatal("an unwritable cache must not change the answer")
	}

	// Offline: no forge call, unknown until the cache knows.
	before := calls
	off := &Checker{Dir: dir, File: filepath.Join(dir, "none.json"), Run: run, Look: found, Offline: true}
	if a := off.Ask(12); a.Merged() || a.State != "" || calls != before {
		t.Fatalf("offline asked the forge or guessed: %+v, %d calls", a, calls-before)
	}
	if a := (&Checker{File: corrupt, Offline: true}).Ask(12); !a.Merged() {
		t.Fatalf("offline must read what the cache knows: %+v", a)
	}
}

func TestAnswerWhy(t *testing.T) {
	for _, c := range []struct {
		a    Answer
		want string
	}{
		{Answer{State: "open"}, "open"},
		{Answer{State: "closed"}, "closed without merge"},
		{Answer{State: "merged"}, ""},
		{Answer{}, "not asked here"},
		{Answer{Err: ErrNoTool}, "state unknown"},
	} {
		if got := c.a.Why(); !strings.Contains(got, c.want) || (c.want == "" && got != "") {
			t.Errorf("%+v.Why() = %q, want %q", c.a, got, c.want)
		}
	}
}
