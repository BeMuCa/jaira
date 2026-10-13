package ticket

import (
	"os"
	"strings"
	"testing"
	"time"
)

// cacheStore is a store with two tickets whose files are old enough to sit
// outside the racy window, so the cache may trust their time and size.
func cacheStore(t *testing.T) (*Store, string) {
	t.Helper()
	s := &Store{Root: t.TempDir()}
	for _, id := range []string{"AAAAAA", "BBBBBB"} {
		if _, err := s.Create(map[string]string{FieldID: id, FieldTitle: "title " + id, FieldStatus: "todo"}, nil, "body\n"); err != nil {
			t.Fatal(err)
		}
	}
	paths, err := s.Paths()
	if err != nil || len(paths) != 2 {
		t.Fatalf("paths: %v %v", paths, err)
	}
	old := time.Now().Add(-time.Hour)
	for _, p := range paths {
		if err := os.Chtimes(p, old, old); err != nil {
			t.Fatal(err)
		}
	}
	return s, paths[0]
}

func titles(t *testing.T, s *Store) string {
	t.Helper()
	ts, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for _, x := range ts {
		b.WriteString(x.Title + ";")
	}
	return b.String()
}

// sameSizeRewrite changes a title in place, keeping the file's size and putting
// its modification time back, so only the bytes differ.
func sameSizeRewrite(t *testing.T, path, from, to string) {
	t.Helper()
	fi, _ := os.Stat(path)
	raw, _ := os.ReadFile(path)
	out := strings.Replace(string(raw), from, to, 1)
	if len(out) != len(raw) || out == string(raw) {
		t.Fatalf("bad rewrite %q -> %q", from, to)
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, fi.ModTime(), fi.ModTime()); err != nil {
		t.Fatal(err)
	}
}

// An unchanged ticket is neither read nor decoded on the next listing.
func TestCacheUnchangedReadsAndDecodesNothing(t *testing.T) {
	s, _ := cacheStore(t)
	s.EnableCache()
	first := titles(t, s)
	reads := s.cache.fileReads
	dec0, par0 := Counts()
	for i := 0; i < 3; i++ {
		if got := titles(t, s); got != first {
			t.Fatalf("listing changed: %q vs %q", got, first)
		}
	}
	dec1, par1 := Counts()
	if s.cache.fileReads != reads || dec1 != dec0 || par1 != par0 {
		t.Fatalf("unchanged tickets were reread: reads %d->%d decodes %d->%d parses %d->%d",
			reads, s.cache.fileReads, dec0, dec1, par0, par1)
	}
}

// A rewrite that keeps size and modification time is invisible to the cache
// until it is told, and then it is found.
func TestCacheSameSizeSameTimeFoundAfterForget(t *testing.T) {
	s, path := cacheStore(t)
	s.EnableCache()
	titles(t, s)
	sameSizeRewrite(t, path, "title AAAAAA", "title ZZZZZZ")
	if got := titles(t, s); strings.Contains(got, "ZZZZZZ") {
		t.Fatalf("outside the racy window the cache should trust time and size, got %q", got)
	}
	s.Forget(path)
	if got := titles(t, s); !strings.Contains(got, "title ZZZZZZ") {
		t.Fatalf("Forget did not make the change visible: %q", got)
	}
}

func TestCacheForgetAllFindsHiddenChange(t *testing.T) {
	s, path := cacheStore(t)
	s.EnableCache()
	titles(t, s)
	sameSizeRewrite(t, path, "title AAAAAA", "title ZZZZZZ")
	s.ForgetAll()
	if got := titles(t, s); !strings.Contains(got, "title ZZZZZZ") {
		t.Fatalf("ForgetAll did not make the change visible: %q", got)
	}
}

// Inside the racy window the bytes are compared, so a same-size, same-time
// rewrite is found without any help.
func TestCacheRacyWindowComparesBytes(t *testing.T) {
	s := &Store{Root: t.TempDir()}
	tk, err := s.Create(map[string]string{FieldID: "AAAAAA", FieldTitle: "title AAAAAA", FieldStatus: "todo"}, nil, "body\n")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	os.Chtimes(tk.Path, now, now)
	s.EnableCache()
	titles(t, s)
	sameSizeRewrite(t, tk.Path, "title AAAAAA", "title ZZZZZZ")
	if got := titles(t, s); !strings.Contains(got, "title ZZZZZZ") {
		t.Fatalf("racy rewrite not found: %q", got)
	}
}

// A change that alters the size is found with no help at all.
func TestCacheChangedSizeFound(t *testing.T) {
	s, path := cacheStore(t)
	s.EnableCache()
	titles(t, s)
	raw, _ := os.ReadFile(path)
	out := strings.Replace(string(raw), "title AAAAAA", "a much longer title", 1)
	os.WriteFile(path, []byte(out), 0o644)
	old := time.Now().Add(-time.Hour)
	os.Chtimes(path, old, old)
	if got := titles(t, s); !strings.Contains(got, "a much longer title") {
		t.Fatalf("size change not found: %q", got)
	}
}

// Editing a listed ticket does not edit the cached one.
func TestCacheHandsOutCopies(t *testing.T) {
	s, _ := cacheStore(t)
	s.EnableCache()
	ts, _ := s.List()
	ts[0].Title = "scribbled"
	ts[0].Tags = append(ts[0].Tags, "x")
	_ = ts[0].Doc().SetScalar("title", "scribbled")
	if got := titles(t, s); strings.Contains(got, "scribbled") {
		t.Fatalf("caller edit reached the cache: %q", got)
	}
}

func TestCacheDeletedFileDisappears(t *testing.T) {
	s, path := cacheStore(t)
	s.EnableCache()
	titles(t, s)
	os.Remove(path)
	if got := titles(t, s); strings.Count(got, ";") != 1 {
		t.Fatalf("deleted ticket still listed: %q", got)
	}
}

// Without EnableCache every listing reads and decodes everything, as before.
func TestCacheOffReadsEveryTime(t *testing.T) {
	s, _ := cacheStore(t)
	titles(t, s)
	dec0, _ := Counts()
	titles(t, s)
	if dec1, _ := Counts(); dec1-dec0 < 2 {
		t.Fatalf("cache off should decode every ticket each listing, decoded %d", dec1-dec0)
	}
}
