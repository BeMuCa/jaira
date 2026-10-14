package ticket

import (
	"bytes"
	"os"
	"sync"
	"time"
)

// racyWindow is how recent a file's modification time may be, relative to when
// it was read, before a matching time and size no longer proves the content is
// unchanged. FAT and exFAT keep modification times in 2-second steps and Linux
// stamps with a coarse clock, so a same-sized rewrite inside that window can
// keep the old time. Git calls this the "racy" case and handles it the same way:
// inside the window, compare the bytes.
const racyWindow = 2 * time.Second

// headerCache remembers decoded ticket headers by path.
type headerCache struct {
	mu        sync.Mutex
	m         map[string]*cacheEntry
	fileReads int64 // files actually read, for tests
}

type cacheEntry struct {
	mod   time.Time
	size  int64
	seen  time.Time // when the bytes were read
	raw   []byte
	t     *Ticket // never handed out; callers get a Clone
	err   error
	id    string
	idErr error
}

// EnableCache makes List, LoggedSince and the id lookup skip a ticket file whose
// modification time and size are unchanged since it was last decoded. It is off
// until called: a one-shot CLI process has nothing to remember, and only the
// board, which lists every 2 seconds, gains from it. Load is never cached, since
// it is the path a write starts from.
func (s *Store) EnableCache() {
	if s.cache == nil {
		s.cache = &headerCache{m: map[string]*cacheEntry{}}
	}
}

// Forget drops the cached headers of these paths, so the next read decodes them
// again. A path that was never cached, or any call with the cache off, is a no-op.
func (s *Store) Forget(paths ...string) {
	if s.cache == nil {
		return
	}
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	for _, p := range paths {
		delete(s.cache.m, p)
	}
}

// ForgetAll drops every cached header.
func (s *Store) ForgetAll() {
	if s.cache == nil {
		return
	}
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	clear(s.cache.m)
}

// cached returns the decoded header of path, reading and decoding the file only
// when its time or size changed, or when it is too recent for them to be trusted.
func (s *Store) cached(path string) (*cacheEntry, error) {
	c := s.cache
	c.mu.Lock()
	defer c.mu.Unlock()
	fi, err := os.Stat(path)
	if err != nil {
		delete(c.m, path)
		return nil, err
	}
	e := c.m[path]
	same := e != nil && e.mod.Equal(fi.ModTime()) && e.size == fi.Size()
	if same && e.seen.Sub(e.mod) >= racyWindow {
		return e, nil
	}
	seen := time.Now()
	raw, err := os.ReadFile(path)
	c.fileReads++
	if err != nil {
		delete(c.m, path)
		return nil, err
	}
	if same && bytes.Equal(raw, e.raw) {
		e.seen = seen
		return e, nil
	}
	e = &cacheEntry{mod: fi.ModTime(), size: fi.Size(), seen: seen, raw: raw}
	d, err := ParseDoc(raw)
	if err != nil {
		e.err, e.idErr = err, err
	} else if f, err := d.ticketFields(); err != nil {
		e.err, e.idErr = err, err
	} else {
		if e.idErr = f.mappingErr(); e.idErr == nil {
			var v string
			v, _, e.idErr = f.Scalar(FieldID)
			e.id = NormalizeIDPrefix(v)
		}
		e.t, e.err = decodeFields(d, f, path)
	}
	c.m[path] = e
	return e, nil
}
