package ticket

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// deletePending creates the file at path and leaves it "delete pending": it is
// marked for deletion on a handle that stays open until the returned function
// runs. That is the state an unlocked lock file is in while some other handle
// on it — the antivirus, the indexer — has not been closed yet, and the state
// in which creating it again fails with Access is denied instead of "exists".
//
// os.Remove does not reproduce it: current Windows deletes with POSIX
// semantics where it can, which frees the name at once. The legacy
// FileDispositionInfo class keeps the name until the last handle closes.
func deletePending(t *testing.T, path string) (release func()) {
	t.Helper()
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := windows.CreateFile(name,
		windows.GENERIC_WRITE|windows.DELETE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.CREATE_NEW, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatal(err)
	}
	var del uint32 = 1 // FILE_DISPOSITION_INFO.DeleteFile
	if err := windows.SetFileInformationByHandle(h, windows.FileDispositionInfo,
		(*byte)(unsafe.Pointer(&del)), uint32(unsafe.Sizeof(del))); err != nil {
		windows.CloseHandle(h)
		t.Fatal(err)
	}
	// The premise of both tests: creating the file now is denied, not "exists".
	if _, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644); !os.IsPermission(err) {
		windows.CloseHandle(h)
		t.Fatalf("a delete-pending file was not refused with Access is denied: %v", err)
	}
	return sync.OnceFunc(func() { windows.CloseHandle(h) })
}

// A lock file that is still being deleted is a busy lock: lockFile waits for
// it, and takes the lock as soon as the last handle on the old file closes.
// Before 32TTSS it returned Access is denied at once (8X9SHE).
func TestLockFileWaitsForADeletePendingLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "milestones.lock")
	release := deletePending(t, path)
	defer release()

	got := make(chan error, 1)
	go func() {
		unlock, err := lockFile(path)
		if err == nil {
			unlock()
		}
		got <- err
	}()

	select {
	case err := <-got:
		t.Fatalf("lockFile did not wait for the delete-pending lock (returned %v)", err)
	case <-time.After(150 * time.Millisecond):
	}
	release()

	select {
	case err := <-got:
		if err != nil {
			t.Fatalf("lockFile after the old lock file went away: %v", err)
		}
	case <-time.After(lockTimeout):
		t.Fatal("lockFile did not take the lock after the old lock file went away")
	}
}

// A denial that does not go away is not waited for forever: at lockTimeout
// lockFile gives up with the denial itself, not with a timeout that would
// blame a holder there is none of.
func TestLockFileReportsAPersistentDenialAtTheDeadline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "milestones.lock")
	release := deletePending(t, path)
	defer release()

	start := time.Now()
	unlock, err := lockFile(path)
	if err == nil {
		unlock()
		t.Fatal("lockFile took a lock whose file is still being deleted")
	}
	if waited := time.Since(start); waited < lockTimeout {
		t.Errorf("lockFile gave up after %v, before lockTimeout (%v)", waited, lockTimeout)
	}
	if errors.Is(err, errLockTimeout) || !os.IsPermission(err) {
		t.Errorf("want the Access is denied error, got %v", err)
	}
}
