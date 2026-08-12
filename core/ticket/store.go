package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// DirName is the per-repo directory holding all jaira state.
	DirName = ".jaira"
	// TicketsSubdir holds one markdown file per ticket.
	TicketsSubdir = "tickets"
	// SessionsSubdir and locksSubdir live under the user's home directory rather
	// than inside the repository. Keeping them out means .jaira/ contains only
	// content that is meant to be committed, so there is no mixed directory to
	// explain and no per-repo gitignore needed to protect scratch state.
	SessionsSubdir = "sessions"
	locksSubdir    = "locks"

	// frontmatterProbe caps how much of a file is read when only the
	// frontmatter is needed. Listing the board reads every ticket, so reading
	// whole files would make startup scale with total prose rather than with
	// ticket count.
	frontmatterProbe = 16 << 10

	lockTimeout = 5 * time.Second
	lockStale   = 30 * time.Second
)

var (
	// ErrNotFound means no ticket matched.
	ErrNotFound = errors.New("ticket: not found")
	// ErrAmbiguous means an ID prefix matched more than one ticket.
	ErrAmbiguous = errors.New("ticket: ambiguous id prefix")
	// ErrNoStore means no .jaira directory was found.
	ErrNoStore = errors.New("ticket: no .jaira directory found; run 'jaira init'")
)

// Store is the ticket directory for one repository.
type Store struct {
	// Root is the directory containing .jaira.
	Root string
}

// Discover walks up from dir looking for an existing .jaira directory.
func Discover(dir string) (*Store, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for {
		candidate := filepath.Join(abs, DirName)
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			return &Store{Root: abs}, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return nil, ErrNoStore
		}
		abs = parent
	}
}

// At returns a store rooted at dir without requiring it to exist yet.
func At(dir string) (*Store, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	return &Store{Root: abs}, nil
}

func (s *Store) dir() string        { return filepath.Join(s.Root, DirName) }
func (s *Store) TicketsDir() string { return filepath.Join(s.dir(), TicketsSubdir) }

// SessionsDir is where this working tree's session state lives, outside the repo.
func (s *Store) SessionsDir() string { return filepath.Join(s.stateDir(), SessionsSubdir) }
func (s *Store) locksDir() string    { return filepath.Join(s.stateDir(), locksSubdir) }

// stateDir is a per-working-tree directory under the user's home.
//
// Keyed by working tree rather than by repository: two clones of the same project
// are two separate sets of in-flight work, and a session focused on one has
// nothing to say about the other.
func (s *Store) stateDir() string {
	home := os.Getenv("JAIRA_HOME")
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			// No home directory to write to; fall back inside the repo so the
			// tool still works, accepting the untracked directory.
			return filepath.Join(s.dir(), "state")
		}
		home = filepath.Join(h, DirName)
	}
	sum := sha256.Sum256([]byte(s.Root))
	key := filepath.Base(s.Root) + "-" + hex.EncodeToString(sum[:4])
	return filepath.Join(home, "state", key)
}

// Init creates the store layout. Safe to run repeatedly.
func (s *Store) Init() (created bool, err error) {
	if _, err := os.Stat(s.dir()); err == nil {
		created = false
	} else {
		created = true
	}
	for _, d := range []string{s.TicketsDir(), s.SessionsDir(), s.locksDir()} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return created, err
		}
	}

	return created, nil
}

// Paths lists ticket files in deterministic order. ULID filenames sort
// lexicographically by creation time, so this is chronological for free.
func (s *Store) Paths() ([]string, error) {
	entries, err := os.ReadDir(s.TicketsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoStore
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		out = append(out, filepath.Join(s.TicketsDir(), e.Name()))
	}
	sort.Strings(out)
	return out, nil
}

// List loads every ticket, reading only as much of each file as the frontmatter
// requires.
func (s *Store) List() ([]*Ticket, error) {
	paths, err := s.Paths()
	if err != nil {
		return nil, err
	}
	out := make([]*Ticket, 0, len(paths))
	var problems []string
	for _, p := range paths {
		t, err := s.loadHeader(p)
		if err != nil {
			// One malformed ticket must not blank the whole board. Collect and
			// report, but keep rendering everything else.
			problems = append(problems, fmt.Sprintf("%s: %v", filepath.Base(p), err))
			continue
		}
		out = append(out, t)
	}
	if len(problems) > 0 {
		return out, &PartialError{Problems: problems}
	}
	return out, nil
}

// PartialError reports tickets that could not be read while others succeeded.
type PartialError struct{ Problems []string }

func (e *PartialError) Error() string {
	return fmt.Sprintf("%d ticket(s) could not be read: %s", len(e.Problems), strings.Join(e.Problems, "; "))
}

// loadHeader reads a ticket's frontmatter without necessarily reading its body.
func (s *Store) loadHeader(path string) (*Ticket, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, frontmatterProbe)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return nil, err
	}
	buf = buf[:n]

	if !hasClosingDelim(buf) {
		// Rare: frontmatter longer than the probe. Fall back to a full read
		// rather than failing.
		all, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		buf = all
	}
	d, err := ParseDoc(buf)
	if err != nil {
		return nil, err
	}
	return Decode(d, path)
}

func hasClosingDelim(b []byte) bool {
	s := string(b)
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "\ufeff---\n") {
		return false
	}
	i := strings.Index(s, "\n---")
	return i >= 0
}

// Load reads one ticket in full, resolving an exact ID or unambiguous prefix.
func (s *Store) Load(idOrPrefix string) (*Ticket, error) {
	path, err := s.resolve(idOrPrefix)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d, err := ParseDoc(raw)
	if err != nil {
		return nil, err
	}
	return Decode(d, path)
}

// resolve maps a reference to exactly one ticket path.
//
// Identity comes from the frontmatter `id`, never from the filename. The filename
// is a human convenience that may be renamed or hand-written, and when the two
// disagreed an earlier version of this code made the ticket unreachable by any
// reference at all — invisible to every command while sitting in plain sight on
// disk. The file content is the only thing that can be trusted to say what a
// ticket is.
//
// An exact id wins outright. Otherwise a unique prefix or a unique suffix is
// accepted. Suffix matching is not an afterthought: a ULID's first ten characters
// encode only its millisecond timestamp, so tickets created in the same burst —
// the normal case when an agent decomposes one task — share a long common prefix
// and differ only in their random tail.
func (s *Store) resolve(ref string) (string, error) {
	want := NormalizeIDPrefix(ref)
	if want == "" {
		return "", ErrNotFound
	}
	ids, err := s.idIndex()
	if err != nil {
		return "", err
	}
	var matches []string
	for id, p := range ids {
		if id == want {
			return p, nil // exact match always wins
		}
		if strings.HasPrefix(id, want) || strings.HasSuffix(id, want) {
			matches = append(matches, p)
		}
	}
	sort.Strings(matches)
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("%w: %s", ErrNotFound, ref)
	case 1:
		return matches[0], nil
	default:
		handles := make([]string, 0, len(matches))
		for _, m := range matches {
			if id, err := s.idOf(m); err == nil {
				handles = append(handles, Handle(id))
			}
		}
		sort.Strings(handles)
		return "", fmt.Errorf("%w: %q matches %s", ErrAmbiguous, ref, strings.Join(handles, ", "))
	}
}

// idIndex maps each ticket's declared id to its path. Built by reading only the
// frontmatter of each file, so it stays cheap as ticket bodies grow.
func (s *Store) idIndex() (map[string]string, error) {
	paths, err := s.Paths()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(paths))
	for _, p := range paths {
		id, err := s.idOf(p)
		if err != nil || id == "" {
			continue // unreadable or unidentified files are reported by List
		}
		out[id] = p
	}
	return out, nil
}

// idOf reads just the id from a ticket file.
func (s *Store) idOf(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, frontmatterProbe)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	d, err := ParseDoc(buf[:n])
	if err != nil {
		return "", err
	}
	v, _, err := d.Scalar(FieldID)
	return v, err
}

// Create writes a new ticket file and returns it.
func (s *Store) Create(fields map[string]string, lists map[string][]string, body string) (*Ticket, error) {
	if err := os.MkdirAll(s.TicketsDir(), 0o755); err != nil {
		return nil, err
	}
	id := fields[FieldID]
	if id == "" {
		return nil, errors.New("ticket: cannot create without an id")
	}
	path := filepath.Join(s.TicketsDir(), Filename(id, fields[FieldTitle]))
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("ticket: %s already exists", filepath.Base(path))
	}
	d := NewDoc(fields, lists, body)
	if err := writeAtomic(path, d.Bytes()); err != nil {
		return nil, err
	}
	return Decode(d, path)
}

// Mutate applies fn to a ticket under an exclusive lock, then writes it back
// atomically. Read-modify-write is done inside the lock so two concurrent
// invocations cannot each read the same original and overwrite one another —
// having a single write path prevents schema drift but does nothing about
// interleaving, which is a distinct failure mode (STORE-07).
func (s *Store) Mutate(idOrPrefix string, fn func(*Ticket) error) (*Ticket, error) {
	path, err := s.resolve(idOrPrefix)
	if err != nil {
		return nil, err
	}
	id := IDFromFilename(filepath.Base(path))

	unlock, err := s.lock(id)
	if err != nil {
		return nil, err
	}
	defer unlock()

	// Re-read inside the lock: the file may have changed between resolve and
	// acquiring the lock.
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d, err := ParseDoc(raw)
	if err != nil {
		return nil, err
	}
	t, err := Decode(d, path)
	if err != nil {
		return nil, err
	}
	if err := fn(t); err != nil {
		return nil, err
	}
	if err := Touch(t.doc, time.Now()); err != nil {
		return nil, err
	}
	if err := writeAtomic(path, t.doc.Bytes()); err != nil {
		return nil, err
	}
	return Decode(t.doc, path)
}

// lock takes an advisory per-ticket lock. A lock file is used rather than flock
// so the same code works on Windows without cgo. Stale locks (from a process
// that died) are broken automatically so the store can never wedge permanently.
func (s *Store) lock(id string) (func(), error) {
	if err := os.MkdirAll(s.locksDir(), 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(s.locksDir(), id+".lock")
	deadline := time.Now().Add(lockTimeout)
	backoff := 2 * time.Millisecond

	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			fmt.Fprintf(f, "%d\n", os.Getpid())
			f.Close()
			return func() { os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		if fi, statErr := os.Stat(path); statErr == nil && time.Since(fi.ModTime()) > lockStale {
			os.Remove(path) // previous holder died
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("ticket: timed out waiting for lock on %s", id)
		}
		time.Sleep(backoff)
		if backoff < 100*time.Millisecond {
			backoff *= 2
		}
	}
}

// writeAtomic writes via a temporary file in the same directory and renames it
// into place, so a crash or a full disk leaves the previous file intact rather
// than a truncated one (STORE-06).
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".jaira-tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName) // no-op once renamed
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// Save writes a ticket that the caller already mutated, without taking a lock.
// Used by the merge driver, which runs single-threaded under git.
func (s *Store) Save(t *Ticket) error { return writeAtomic(t.Path, t.doc.Bytes()) }
