package ticket

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func staleFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
}

func exists(path string) bool { _, err := os.Stat(path); return err == nil }

// A dead holder's lock is removed, and the guard taken for it is not left behind.
func TestBreakStaleRemovesADeadLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tags.lock")
	staleFile(t, path)
	if !breakStale(path) {
		t.Fatal("stale lock not broken")
	}
	if exists(path) || exists(path+".break") {
		t.Fatal("lock or guard left behind")
	}
}

// While another waiter holds the guard - it is breaking the same dead lock and
// may already have taken a fresh one - this waiter removes nothing. Without the
// guard both waiters removed, the second one the first one's fresh lock.
func TestBreakStaleLeavesTheLockToTheWaiterHoldingTheGuard(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tags.lock")
	staleFile(t, path)
	if err := os.WriteFile(path+".break", nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if breakStale(path) {
		t.Fatal("broke a lock while another waiter held the guard")
	}
	if !exists(path) || !exists(path+".break") {
		t.Fatal("lock or the other waiter's guard removed")
	}
}

// A live lock is never broken.
func TestBreakStaleLeavesAFreshLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tags.lock")
	if err := os.WriteFile(path, []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if breakStale(path) || !exists(path) {
		t.Fatal("fresh lock broken")
	}
}

// A guard left by a waiter that died while holding it does not wedge the lock:
// it is cleared once stale, and the next look breaks the dead lock.
func TestBreakStaleClearsADeadGuard(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tags.lock")
	staleFile(t, path)
	staleFile(t, path+".break")
	if breakStale(path) {
		t.Fatal("broke the lock in the same look that cleared the dead guard")
	}
	if exists(path + ".break") {
		t.Fatal("dead guard not cleared")
	}
	if !breakStale(path) || exists(path) {
		t.Fatal("dead lock not broken after the guard was cleared")
	}
}
