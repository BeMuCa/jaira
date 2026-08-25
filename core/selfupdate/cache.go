package selfupdate

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Check is what gets cached: the answer to "what is the latest published
// release", and when that answer was last obtained.
type Check struct {
	CheckedAt time.Time `json:"checked_at"`
	Latest    string    `json:"latest"`
}

// MaxAge is how long a cached Check is trusted before it counts as stale.
const MaxAge = 24 * time.Hour

// Path is where the cache file lives: JAIRA_HOME if set, else ~/.jaira,
// never inside a repository. Which release is published is a fact about the
// machine, not about a working tree — putting it per-tree would make
// someone with ten clones run ten checks a day for the same answer.
func Path() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "update-check.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "update-check.json")
}

// Read returns the cached Check, or (Check{}, false) on any error. A
// missing, truncated or garbage cache is indistinguishable from "never
// checked" for every purpose this serves, so every error is swallowed here
// the same way release.Stamped swallows its own — a best-effort
// informational cache must never be able to fail a command.
func Read() (Check, bool) {
	path := Path()
	if path == "" {
		return Check{}, false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Check{}, false
	}
	var c Check
	if err := json.Unmarshal(b, &c); err != nil {
		return Check{}, false
	}
	return c, true
}

// Write records c, through a temp file in the same directory plus a
// rename — the same atomic-write idiom the rest of this project uses for
// anything it cares about.
func Write(c Check) error {
	path := Path()
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".jaira-update-check-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if _, err := f.Write(b); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// Stale reports whether c should be refreshed: an absent or unparseable
// cache (ok == false) always is, and so is one whose CheckedAt is more than
// MaxAge old.
func Stale(c Check, ok bool) bool {
	if !ok {
		return true
	}
	return time.Since(c.CheckedAt) > MaxAge
}

// Disabled reports whether the background release check has been switched
// off.
func Disabled() bool {
	return os.Getenv("JAIRA_NO_UPDATE_CHECK") == "1"
}

// SpawnRefresh starts a detached child that refreshes the cache and does
// not wait for it.
//
// The child gets JAIRA_NO_UPDATE_CHECK=1 as a recursion guard: it must never
// itself decide the cache is stale and spawn a grandchild. That happens to
// be the same switch a user flips to turn the whole thing off, so there is
// one mechanism here, not two.
//
// There is no Wait() — that is the whole point. The parent is about to
// exit; the child is reparented to init and reaped there. There is
// deliberately no attempt to setsid: that needs per-OS SysProcAttr and buys
// only immunity from a Ctrl-C sent to the process group, and losing a
// best-effort release check to a Ctrl-C costs nothing — the cache stays
// stale and the next command retries.
//
// Stdio goes to os.DevNull rather than being inherited, because a child
// writing to the parent's terminal after the parent has already exited is
// the exact failure this whole design exists to avoid.
func SpawnRefresh() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	// Never from inside a test binary. The child is invoked as
	// "<exe> self upgrade --check --json", which only means anything if exe is
	// jaira; a Go test binary has no such command and simply runs its whole
	// suite again, detached, reaching the network on the way. On Windows it
	// then holds its own image open and 'go test' fails to remove it after
	// every test has passed — which is exactly how this was found.
	if isTestBinary(exe) {
		return nil
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer devNull.Close()

	cmd := exec.Command(exe, "self", "upgrade", "--check", "--json")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = devNull, devNull, devNull
	cmd.Env = append(os.Environ(), "JAIRA_NO_UPDATE_CHECK=1")
	return cmd.Start()
}

// isTestBinary reports whether exe is a Go test binary, which 'go test' names
// <package>.test (plus .exe on Windows).
func isTestBinary(exe string) bool {
	base := filepath.Base(exe)
	base = strings.TrimSuffix(base, ".exe")
	return strings.HasSuffix(base, ".test")
}

// PollCache returns the last known "latest published release" answer,
// refreshing it in the background first if it is stale.
//
// This is the one non-network entry point a consumer should read through
// instead of composing Read/Stale/Write/SpawnRefresh itself — currently
// that consumer is the TUI's persistent version indicator, not any CLI
// command. It never blocks on the network: a stale cache gets touched and
// handed to a detached child (see SpawnRefresh's own doc comment for why),
// and this function returns immediately with whatever was already known,
// which may be nothing at all.
//
// Disabled() short-circuits before touching anything, so an opted-out
// machine gets neither a spawned refresh nor a reported answer.
//
// known is false exactly when nothing has actually been learned yet — the
// check is disabled, or the cache is absent/unparseable and no prior write
// ever recorded a real Latest. A caller must not render "up to date" on a
// false known: that would assert a fact nobody has checked.
func PollCache() (latest string, known bool) {
	if Disabled() {
		return "", false
	}
	c, ok := Read()
	if Stale(c, ok) {
		// Touching the cache before spawning is what stops N parallel
		// sessions — several TUI instances, or a mix of TUI and CLI
		// invocations — from each spawning their own refresher for the same
		// answer: the first to notice staleness writes a freshly-timestamped
		// entry (carrying whatever Latest it already had, possibly none), so
		// every other caller in the same window sees a fresh cache and skips
		// its own spawn. The child overwrites this with the real answer
		// moments later; if it dies first, the worst case is one missed day
		// of the indicator updating, which is the correct trade against ever
		// blocking on the network.
		_ = Write(Check{CheckedAt: time.Now().UTC(), Latest: c.Latest})
		_ = SpawnRefresh()
	}
	if !ok || c.Latest == "" {
		return "", false
	}
	return c.Latest, true
}
