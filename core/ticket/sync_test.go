package ticket

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// syncTestStore builds a fresh store with one ticket, and returns the store
// and the ticket's id.
func syncTestStore(t *testing.T) (s *Store, id string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		FieldID:     NewID(time.Now()),
		FieldTitle:  "t",
		FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return s, tk.ID
}

func TestSyncMovesTicketIntoDatedFolder(t *testing.T) {
	s, id := syncTestStore(t)
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(tk.Path)

	dst, err := s.Sync(id, "as-20260823")
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if dst != filepath.Join(s.SyncDir(), "as-20260823", base) {
		t.Errorf("Sync() = %q, want the file under %s", dst, filepath.Join(s.SyncDir(), "as-20260823"))
	}
	if _, err := os.Stat(tk.Path); !os.IsNotExist(err) {
		t.Errorf("original ticket path %q still exists after Sync", tk.Path)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("synced file not found at %q: %v", dst, err)
	}
}

func TestSyncRefusesNameCollision(t *testing.T) {
	s, id := syncTestStore(t)
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(tk.Path)

	dst, err := s.Sync(id, "as-20260823")
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	// Recreate a ticket with the same filename so a second sync into the same
	// folder collides.
	if err := os.MkdirAll(s.TicketsDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.TicketsDir(), base), []byte("---\nid: "+id+"\ntitle: t\nstatus: backlog\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Sync(id, "as-20260823"); err == nil {
		t.Fatal("expected Sync to refuse a name collision, got nil error")
	} else if !strings.Contains(err.Error(), base) {
		t.Errorf("collision error %q does not name the file %q", err, base)
	}
	// The first sync's file is untouched by the refused second attempt.
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("first synced file disappeared after a refused second sync: %v", err)
	}
}

func TestRestoreFindsArchivedTicket(t *testing.T) {
	s, id := syncTestStore(t)
	dst, err := s.Archive(id)
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	base := filepath.Base(dst)

	restored, err := s.Restore(base)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if restored != filepath.Join(s.TicketsDir(), base) {
		t.Errorf("Restore() = %q, want it back under %s", restored, s.TicketsDir())
	}
}

func TestRestoreFindsSyncedTicket(t *testing.T) {
	s, id := syncTestStore(t)
	dst, err := s.Sync(id, "as-20260823")
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	base := filepath.Base(dst)

	restored, err := s.Restore(base)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if restored != filepath.Join(s.TicketsDir(), base) {
		t.Errorf("Restore() = %q, want it back under %s", restored, s.TicketsDir())
	}
}

func TestRestoreOfUnknownNameNamesBothPlaces(t *testing.T) {
	s, _ := syncTestStore(t)
	_, err := s.Restore("nope.md")
	if err == nil {
		t.Fatal("expected Restore of an unknown name to fail")
	}
	if !strings.Contains(err.Error(), "archive") || !strings.Contains(err.Error(), "sync") {
		t.Errorf("Restore error %q does not name both the archive and sync as possible places", err)
	}
}

func TestRestoreAmbiguousAcrossSyncFoldersIsRefused(t *testing.T) {
	s, id := syncTestStore(t)
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(tk.Path)

	dst, err := s.Sync(id, "as-20260823")
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	// Place a second file with the same base name in a different sync folder,
	// so Restore has two candidates and must refuse rather than guess.
	otherDir := filepath.Join(s.SyncDir(), "amr-20260824")
	if err := os.MkdirAll(otherDir, 0o755); err != nil {
		t.Fatal(err)
	}
	otherPath := filepath.Join(otherDir, base)
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(otherPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Restore(base); err == nil {
		t.Fatal("expected Restore to refuse when the name matches two sync folders")
	} else if !strings.Contains(err.Error(), "as-20260823") || !strings.Contains(err.Error(), "amr-20260824") {
		t.Errorf("ambiguous-restore error %q does not name both folders", err)
	}
}

// TestRestoreCannotEscapeStore asserts a traversal attempt in the name
// argument cannot walk Restore anywhere outside the store: only the base name
// is ever used to build a path, in the archive lookup as well as the sync
// lookup.
func TestRestoreCannotEscapeStore(t *testing.T) {
	s, _ := syncTestStore(t)
	if _, err := s.Restore("../../etc/passwd"); err == nil {
		t.Fatal("expected Restore(\"../../etc/passwd\") to fail rather than escape the store")
	}
}

// TestPathsAndListIgnoreSyncDir asserts a populated .jaira/sync/ is invisible
// to Paths (and therefore List and core/validate): a synced ticket is not a
// board ticket.
func TestPathsAndListIgnoreSyncDir(t *testing.T) {
	s, id := syncTestStore(t)
	if _, err := s.Sync(id, "as-20260823"); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	paths, err := s.Paths()
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	if len(paths) != 0 {
		t.Errorf("Paths() = %v after the only ticket was synced off the board, want empty", paths)
	}

	all, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("List() = %v after the only ticket was synced off the board, want empty", all)
	}
}
