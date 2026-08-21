package selfupdate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPathHonorsJairaHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir)
	got := Path()
	want := filepath.Join(dir, "update-check.json")
	if got != want {
		t.Errorf("Path() = %q, want %q", got, want)
	}
}

// TestPathNeverLandsInARepository asserts the cache path is under the
// user's home, not the current working directory — a repository must never
// end up with this file inside it. JAIRA_HOME is deliberately unset here so
// the fallback (~/.jaira) is what gets exercised against a changed cwd.
func TestPathNeverLandsInARepository(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("JAIRA_HOME", "")
	repo := t.TempDir()
	restore := chdir(t, repo)
	defer restore()

	got := Path()
	want := filepath.Join(home, ".jaira", "update-check.json")
	if got != want {
		t.Errorf("Path() = %q, want %q — independent of the working directory %q", got, want, repo)
	}
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() { os.Chdir(orig) }
}

func TestWriteThenReadRoundTrips(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	want := Check{CheckedAt: time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC), Latest: "1.3.0"}
	if err := Write(want); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, ok := Read()
	if !ok {
		t.Fatal("Read() ok = false after Write")
	}
	if !got.CheckedAt.Equal(want.CheckedAt) {
		t.Errorf("CheckedAt = %v, want %v", got.CheckedAt, want.CheckedAt)
	}
	if got.Latest != want.Latest {
		t.Errorf("Latest = %q, want %q", got.Latest, want.Latest)
	}
}

func TestWriteLeavesNoTempFileBeside(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir)
	if err := Write(Check{CheckedAt: time.Now(), Latest: "1.0.0"}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "update-check.json" {
			t.Errorf("unexpected leftover entry %q in %s", e.Name(), dir)
		}
	}
}

func TestStaleAtThresholds(t *testing.T) {
	cases := []struct {
		name string
		age  time.Duration
		ok   bool
		want bool
	}{
		{"25 hours old is stale", 25 * time.Hour, true, true},
		{"23 hours old is not stale", 23 * time.Hour, true, false},
		{"absent is stale", 0, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := Check{CheckedAt: time.Now().Add(-tc.age)}
			if got := Stale(c, tc.ok); got != tc.want {
				t.Errorf("Stale(age=%v, ok=%v) = %v, want %v", tc.age, tc.ok, got, tc.want)
			}
		})
	}
}

func TestReadAbsentOrUnparseableIsStale(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir)

	c, ok := Read()
	if ok {
		t.Fatalf("Read() on an absent cache: ok = true, c = %+v", c)
	}
	if !Stale(c, ok) {
		t.Error("an absent cache must read as stale")
	}

	if err := os.WriteFile(Path(), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	c2, ok2 := Read()
	if ok2 {
		t.Fatalf("Read() on a garbage cache: ok = true, c = %+v", c2)
	}
	if !Stale(c2, ok2) {
		t.Error("an unparseable cache must read as stale")
	}
}

func TestPollCacheFreshReturnsKnownLatest(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := Write(Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}
	latest, known := PollCache()
	if !known {
		t.Fatal("known = false with a fresh cache present")
	}
	if latest != "1.3.0" {
		t.Errorf("latest = %q, want %q", latest, "1.3.0")
	}
}

func TestPollCacheNeverCheckedReportsUnknown(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	latest, known := PollCache()
	if known {
		t.Fatalf("known = true with no cache ever written; latest = %q", latest)
	}
	if latest != "" {
		t.Errorf("latest = %q, want empty when unknown", latest)
	}
}

func TestPollCacheDisabledReportsUnknownAndDoesNotTouchTheCache(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "1")

	latest, known := PollCache()
	if known {
		t.Fatalf("known = true while disabled; latest = %q", latest)
	}
	if latest != "" {
		t.Errorf("latest = %q, want empty while disabled", latest)
	}
	if _, err := os.Stat(Path()); !os.IsNotExist(err) {
		t.Errorf("PollCache wrote a cache file while disabled: stat err = %v", err)
	}
}

// TestPollCacheOnStaleCacheReturnsWithoutBlocking asserts a stale (here:
// absent) cache does not make PollCache wait on the network — it only ever
// spawns a detached child and returns.
func TestPollCacheOnStaleCacheReturnsWithoutBlocking(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	done := make(chan struct{})
	go func() {
		PollCache()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("PollCache did not return promptly on a stale cache — it must never block on the network")
	}
}

func TestDisabledOnlyOnExactEnvValue(t *testing.T) {
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "1")
	if !Disabled() {
		t.Error("Disabled() = false with JAIRA_NO_UPDATE_CHECK=1")
	}
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if Disabled() {
		t.Error("Disabled() = true with JAIRA_NO_UPDATE_CHECK unset/empty")
	}
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "true")
	if Disabled() {
		t.Error(`Disabled() = true with JAIRA_NO_UPDATE_CHECK="true", want only "1" to count`)
	}
}
