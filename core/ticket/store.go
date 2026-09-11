package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
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
	// ArchiveSubdir holds tickets taken off the board. They are moved rather
	// than deleted: the whole point of the board is that you can still answer
	// what a task was for months later, and a deleted file on a private board
	// is not in git history to recover from.
	ArchiveSubdir = "archive"
	// SharedSubdir holds lanes deliberately published to teammates. Unlike
	// LanesSubdir it is committed: publishing a lane is an opt-in act through
	// the lane settings screen, not a side effect of sharing the board.
	SharedSubdir = "shared"
	// LogbookSubdir holds tickets whose work is finished and whose commit list
	// is final, grouped into dated folders by who took them off the board. A
	// folder under here is a readable record of one person's sweep — unlike
	// ArchiveSubdir, every ticket that lands here has already had its commits
	// stamped, because leaving the board is the point at which every commit
	// is finally known.
	LogbookSubdir = "logbook"
	// legacyLogbookSubdir is what the logbook was called before it was one.
	// Nothing writes here any more; Restore still reads it so a folder written
	// by an earlier build stays restorable.
	legacyLogbookSubdir = "sync"
	// SessionsSubdir and locksSubdir live under the user's home directory rather
	// than inside the repository. tickets/, archive/ and shared/ are the parts of
	// .jaira/ meant to be committed; lanes/ (see core/lane.ProjectLanesDir) is
	// this machine's own scoping and stays gitignored even on a shared board
	// (see core/board.LanesIgnoreLine); sessions and locks stay out of the
	// repository entirely.
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

	// Actor is who this process writes as, recorded on every mutation as
	// updated-by. It is a field rather than a package-level value because a
	// process can hold two stores at once — the board switcher does — and it is
	// set by the caller that already knows the identity, so core/ticket does
	// not have to depend on core/identity. Empty records nothing.
	Actor string

	// Recorder is told about every ticket write, so a ticket can also travel
	// on a git ref of its own and reach someone who does not have the writer's
	// branch. It is an interface set by the caller for the same reason Actor
	// is a field: core/ticket must not know about git, refs or queues, and the
	// package that does (core/refsync) needs to read tickets. Nil records
	// nothing, which is the case for every board without a remote.
	Recorder WriteRecorder

	// Source supplies tickets this board can see but has no file for — a
	// ticket that is still travelling on its own git ref and has not been
	// pulled into work here.
	//
	// It is the mirror of Recorder, and it exists for the same reason: every
	// reader in this program goes through List and Load, and a reader that had
	// to remember to ask a second place would be a reader that silently shows
	// half the board. core/ticket stays unaware of git; the package that knows
	// (core/refsync) fills this in.
	Source TicketSource

	// dupIDs accumulates tickets that declare an id another file already claimed.
	// Two files with one id is an ambiguity a person has to settle, so it is
	// surfaced rather than resolved by read order.
	dupIDs []string
}

// TicketSource hands over tickets that exist for this board without a file
// here.
type TicketSource interface {
	// Extra returns those tickets, skipping any id in have. Implementations
	// mark what they return as ReadOnly: there is no file to write to, so a
	// mutation has to be refused rather than half-applied.
	Extra(have map[string]bool) ([]*Ticket, error)
}

// ErrOnRefOnly means the ticket exists for this board but not as a file here,
// so it cannot be written to until somebody pulls it.
//
// It is separate from ErrNotFound because the two ask different things of the
// user: one means the id is wrong, the other means the ticket is real and one
// command away.
var ErrOnRefOnly = errors.New("ticket: on its ref and not on your disk; pull it first")

// WriteRecorder is told what a ticket now says, after the file has been written.
type WriteRecorder interface {
	// Record is handed the ticket's id and its complete bytes. It is called
	// after the local write has succeeded, so it must never be understood as
	// a veto: by the time it runs, the ticket on disk has already changed.
	Record(id string, content []byte) error

	// RecordFiled says the ticket has been filed away here — logged or
	// archived — and hands over its final bytes.
	//
	// Filing is deliberately not a deletion. The ticket file is at that moment
	// only in the filer's own branch, and until that branch is merged nobody
	// else can see the ticket at all: taking it off the shared channel too
	// would open a window, as long as a review takes, in which somebody
	// notices the same problem and writes it down a second time. So the final
	// state goes onto the channel and stays there until the ticket has
	// arrived somewhere everybody can see.
	RecordFiled(id string, content []byte) error

	// RecordDelete says the ticket is gone for good. Only 'jaira delete' does
	// this: deleting is an intention, not a stage of finishing.
	RecordDelete(id string) error
}

// ErrNotRecorded means the local write went through but the recorder did not
// take it — the ticket is correct on this machine and has not left it.
//
// It is a separate sentinel rather than a plain error because the two halves
// need opposite handling by the caller: the mutation must be reported as done
// (it is), while the failure to hand it on is worth a line to the user. A
// recorder that returns nothing but this is still a working board.
var ErrNotRecorded = errors.New("ticket: the write was not recorded for the remote")

// onlyOnRef refuses an operation that needs a file for a ticket that has none.
//
// Every command that moves or removes a ticket file goes through Archive,
// Delete or Logbook, so the check sits in those three rather than in each
// command. Without it the empty path turns into nonsense — filepath.Base("")
// is ".", and archiving a ticket that is not here reported "already exists in
// the archive".
func onlyOnRef(t *Ticket) error {
	if t != nil && t.ReadOnly {
		return fmt.Errorf("%s: %w", Handle(t.ID), ErrOnRefOnly)
	}
	return nil
}

func (s *Store) recordDelete(id string) error {
	if s.Recorder == nil {
		return nil
	}
	if err := s.Recorder.RecordDelete(id); err != nil {
		return fmt.Errorf("%w: %w", ErrNotRecorded, err)
	}
	return nil
}

// recordFiled reports a ticket filed away, with the bytes it ended up with.
func (s *Store) recordFiled(id, path string) error {
	if s.Recorder == nil {
		return nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotRecorded, err)
	}
	if err := s.Recorder.RecordFiled(id, content); err != nil {
		return fmt.Errorf("%w: %w", ErrNotRecorded, err)
	}
	return nil
}

func (s *Store) record(id string, content []byte) error {
	if s.Recorder == nil {
		return nil
	}
	if err := s.Recorder.Record(id, content); err != nil {
		return fmt.Errorf("%w: %w", ErrNotRecorded, err)
	}
	return nil
}

// DuplicateIDs reports ids claimed by more than one file, discovered during the
// most recent lookup.
func (s *Store) DuplicateIDs() []string { return s.dupIDs }

// Discover walks up from dir looking for an existing .jaira directory.
func Discover(dir string) (*Store, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for {
		// A board is identified by its tickets directory, not merely by a
		// directory named .jaira. The user's global config lives at ~/.jaira, so
		// matching on the name alone would make the home directory look like a
		// board to anything run underneath it.
		candidate := filepath.Join(abs, DirName, TicketsSubdir)
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

// ArchiveDir is where archived tickets live.
func (s *Store) ArchiveDir() string { return filepath.Join(s.dir(), ArchiveSubdir) }

// SharedDir is where lanes published to teammates live. Nothing creates this
// directory at init: an empty shared/ folder committed to every board that
// never publishes a lane would be its own kind of confusion.
func (s *Store) SharedDir() string { return filepath.Join(s.dir(), SharedSubdir) }

// LogbookDir is where tickets taken off the board with their commits stamped
// live, grouped into dated per-person folders.
func (s *Store) LogbookDir() string { return filepath.Join(s.dir(), LogbookSubdir) }

// Archive moves a ticket out of the board, returning its new path.
//
// The file is moved, never removed. Restoring is moving it back, which is why
// this returns the destination rather than swallowing it.
func (s *Store) Archive(id string) (string, error) {
	t, err := s.Load(id)
	if err != nil {
		return "", err
	}
	if err := onlyOnRef(t); err != nil {
		return "", err
	}
	if err := os.MkdirAll(s.ArchiveDir(), 0o755); err != nil {
		return "", err
	}
	// Resolve symlinks first: archiving through a link would otherwise move the
	// link and orphan the file it pointed at.
	src, err := filepath.EvalSymlinks(t.Path)
	if err != nil {
		src = t.Path
	}
	dst := filepath.Join(s.ArchiveDir(), filepath.Base(src))
	if _, err := os.Stat(dst); err == nil {
		return "", fmt.Errorf("%s already exists in the archive", filepath.Base(dst))
	}
	if err := os.Rename(src, dst); err != nil {
		return "", err
	}
	return dst, s.recordFiled(t.ID, dst)
}

// Delete removes a ticket's file and returns the path it was at.
//
// The only irreversible operation in the store, and it exists because archiving
// is the wrong answer for a ticket that should never have been written: a
// mistyped create or a throwaway probe leaves a file that otherwise only 'rm'
// removes, and 'rm' means the caller has to know the file layout. Whether the
// caller is sure is decided above this line — the CLI asks for the handle typed
// back, and the board asks for it again.
//
// Symlinks are resolved for the same reason Archive resolves them: deleting the
// link alone would leave the ticket on disk but off the board. Both go, since a
// half-deleted ticket is worse than either outcome.
func (s *Store) Delete(id string) (string, error) {
	t, err := s.Load(id)
	if err != nil {
		return "", err
	}
	if err := onlyOnRef(t); err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(t.Path)
	if err != nil {
		real = t.Path
	}
	if err := os.Remove(real); err != nil {
		return "", err
	}
	if real != t.Path {
		os.Remove(t.Path)
	}
	return t.Path, s.recordDelete(t.ID)
}

// Logbook moves a ticket out of the board and into a dated logbook folder,
// returning its new path. Like Archive, the file is moved rather than deleted,
// and a name collision inside the destination folder is refused by name rather
// than silently overwritten — this only happens if the same ticket is logged
// into the same folder twice, which a caller should treat as a bug to look
// into, not a state to paper over.
func (s *Store) Logbook(id, folder string) (string, error) {
	t, err := s.Load(id)
	if err != nil {
		return "", err
	}
	if err := onlyOnRef(t); err != nil {
		return "", err
	}
	dir := filepath.Join(s.LogbookDir(), filepath.Base(folder))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	// Resolve symlinks first: moving through a link would otherwise move the
	// link and orphan the file it pointed at.
	src, err := filepath.EvalSymlinks(t.Path)
	if err != nil {
		src = t.Path
	}
	dst := filepath.Join(dir, filepath.Base(src))
	if _, err := os.Stat(dst); err == nil {
		return "", fmt.Errorf("%s already exists in %s", filepath.Base(dst), filepath.Join(LogbookSubdir, filepath.Base(folder)))
	}
	if err := os.Rename(src, dst); err != nil {
		return "", err
	}
	return dst, s.recordFiled(t.ID, dst)
}

// logbookFolders lists the per-person dated folders of the logbook as full
// paths, in deterministic order — those under LogbookDir and, after them,
// those under the name the logbook used to have, so a folder written by an
// earlier build is still found.
func (s *Store) logbookFolders() ([]string, error) {
	var out []string
	for _, sub := range []string{LogbookSubdir, legacyLogbookSubdir} {
		root := filepath.Join(s.dir(), sub)
		entries, err := os.ReadDir(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		var names []string
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}
		sort.Strings(names)
		for _, n := range names {
			out = append(out, filepath.Join(root, n))
		}
	}
	return out, nil
}

// LoggedPerDay counts the logbook's tickets by the day of their folder, for
// the days days ending today: out[days-1] is today, out[0] the oldest. The
// folder name carries the date — <initials>-<yyyymmdd> — so no ticket is
// read, and a folder whose name does not end in a date is skipped. A folder
// under the logbook's old name counts too, the same way Restore finds it.
func (s *Store) LoggedPerDay(now time.Time, days int) []int {
	out := make([]int, days)
	folders, err := s.logbookFolders()
	if err != nil {
		return out
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for _, folder := range folders {
		name := filepath.Base(folder)
		i := strings.LastIndex(name, "-")
		if i < 0 {
			continue
		}
		day, err := time.ParseInLocation("20060102", name[i+1:], now.Location())
		if err != nil {
			continue
		}
		// Rounded, not truncated: across a DST change two midnights are 23 or
		// 25 hours apart, and truncation would put that day off by one.
		ago := int(math.Round(today.Sub(day).Hours() / 24))
		if ago < 0 || ago >= days {
			continue
		}
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				out[days-1-ago]++
			}
		}
	}
	return out
}

// Restore moves an archived or logged ticket back onto the board, resolving
// name from the bare filename in either case — never from a caller-supplied
// path, so a name built to escape the store (e.g. "../../etc/passwd") cannot
// walk it anywhere but back into TicketsDir.
func (s *Store) Restore(name string) (string, error) {
	base := filepath.Base(name)

	var src string
	if _, err := os.Stat(filepath.Join(s.ArchiveDir(), base)); err == nil {
		src = filepath.Join(s.ArchiveDir(), base)
	}

	folders, err := s.logbookFolders()
	if err != nil {
		return "", err
	}
	var logMatches []string
	for _, folder := range folders {
		if _, err := os.Stat(filepath.Join(folder, base)); err == nil {
			logMatches = append(logMatches, folder)
		}
	}
	names := make([]string, 0, len(logMatches))
	for _, m := range logMatches {
		names = append(names, filepath.Base(m))
	}

	switch {
	case src != "" && len(logMatches) > 0:
		return "", fmt.Errorf("%s is in both the archive and %s — remove one before restoring", base, strings.Join(names, ", "))
	case len(logMatches) > 1:
		return "", fmt.Errorf("%s is in more than one logbook folder (%s) — remove one before restoring", base, strings.Join(names, ", "))
	case len(logMatches) == 1:
		src = filepath.Join(logMatches[0], base)
	case src == "":
		return "", fmt.Errorf("%s is not in the archive or in .jaira/logbook/", base)
	}

	dst := filepath.Join(s.TicketsDir(), base)
	if _, err := os.Stat(dst); err == nil {
		return "", fmt.Errorf("%s is already on the board", base)
	}
	if err := os.Rename(src, dst); err != nil {
		return "", err
	}
	return dst, nil
}

// SessionsDir is where this working tree's session state lives, outside the repo.
func (s *Store) SessionsDir() string { return filepath.Join(s.stateDir(), SessionsSubdir) }
func (s *Store) locksDir() string    { return filepath.Join(s.stateDir(), locksSubdir) }

// StateDir is the per-working-tree directory this store's session and lock
// state lives under, exported so callers outside the package (the version
// stamp, in particular) can address it too.
func (s *Store) StateDir() string { return s.stateDir() }

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
	// Duplicate ids are found by the id index, so build it here to surface them
	// alongside unreadable files.
	if _, err := s.idIndex(); err == nil {
		problems = append(problems, s.dupIDs...)
	}
	problems = append(problems, s.offBoardDuplicates(out)...)
	// And then whatever the board can see without having a file for it. Added
	// here rather than by each caller: list, next, the board and the task
	// mirror all read through this one function, and a caller that forgot would
	// show half the board.
	have := make(map[string]bool, len(out))
	for _, t := range out {
		have[t.ID] = true
	}
	out = append(out, s.extra(have)...)
	if len(problems) > 0 {
		return out, &PartialError{Problems: problems}
	}
	return out, nil
}

// offBoardDuplicates reports a ticket that is on the board here and also filed
// away — in the logbook or the archive — at the same time.
//
// It happens without anybody doing anything wrong: one clone logs a finished
// ticket while another still has it under tickets/, and the merge of those two
// branches is a rename on one side and a modification on the other, which git
// resolves by keeping both paths. The ticket then lives twice, closed and open.
//
// The existing duplicate check cannot see this, because it compares ids only
// within tickets/. Nothing here moves a file: which of the two copies is right
// is a question for the person, and a tool that guessed would sometimes reopen
// finished work.
func (s *Store) offBoardDuplicates(onBoard []*Ticket) []string {
	if len(onBoard) == 0 {
		return nil
	}
	here := make(map[string]string, len(onBoard))
	for _, t := range onBoard {
		if t.ReadOnly {
			continue // no file here; there is nothing to be a duplicate of
		}
		here[t.ID] = filepath.Base(t.Path)
	}
	var problems []string
	for _, dir := range s.filedAwayDirs() {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			id := IDFromFilename(e.Name())
			if id == "" {
				continue
			}
			if base, ok := here[id]; ok {
				problems = append(problems, fmt.Sprintf(
					"%s is on the board as %s and also filed away in %s — one of the two has to go, and only you can say which",
					Handle(id), base, filepath.Join(filepath.Base(filepath.Dir(dir)), filepath.Base(dir), e.Name())))
			}
		}
	}
	sort.Strings(problems)
	return problems
}

// FiledAwayIDs maps the ids of tickets filed away here — logged or archived —
// to the file that holds them.
//
// Exported because two questions need it and neither belongs in this package:
// whether a ticket is on the board and filed away at once, and whether a
// finished ticket has arrived anywhere everybody can see it.
func (s *Store) FiledAwayIDs() map[string]string {
	out := map[string]string{}
	for _, dir := range s.filedAwayDirs() {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if id := IDFromFilename(e.Name()); id != "" {
				out[id] = filepath.Join(dir, e.Name())
			}
		}
	}
	return out
}

// filedAwayDirs are the directories a ticket lands in when it leaves the board.
func (s *Store) filedAwayDirs() []string {
	dirs := []string{s.ArchiveDir()}
	if folders, err := s.logbookFolders(); err == nil {
		dirs = append(dirs, folders...)
	}
	return dirs
}

// PartialError reports tickets that could not be read while others succeeded.
type PartialError struct{ Problems []string }

func (e *PartialError) Error() string {
	return fmt.Sprintf("%d ticket(s) could not be read: %s", len(e.Problems), strings.Join(e.Problems, "; "))
}

// readHeader parses a ticket's frontmatter, reading only as much of the file as
// necessary.
//
// Every path that needs frontmatter goes through here. When two paths each did
// their own bounded read, they disagreed about tickets whose frontmatter exceeded
// the probe: one fell back to a full read and listed the ticket, the other gave up
// and could not resolve it — so the ticket was visible on the board yet impossible
// to open or modify.
func (s *Store) readHeader(path string) (*Doc, error) {
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
		// Frontmatter longer than the probe: read the whole file rather than
		// failing. Rare, so the cost is acceptable.
		all, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		buf = all
	}
	return ParseDoc(buf)
}

// loadHeader reads a ticket for listing.
//
// It reads the whole file rather than probing the frontmatter. Half of what a
// ticket says lives in the body — the description and both checklists — and the
// board now renders checklist progress and searches body text, so a truncated
// read produced cards that quietly disagreed with the gate: a ticket past the
// probe showed no checklist at all while the gate saw two outstanding items.
//
// The probe remains where it is still correct: idOf only needs the id.
func (s *Store) loadHeader(path string) (*Ticket, error) {
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

func hasClosingDelim(b []byte) bool {
	s := string(b)
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "\ufeff---\n") {
		return false
	}
	i := strings.Index(s, "\n---")
	return i >= 0
}

// extra asks the source for tickets with no file here, skipping the ids given.
func (s *Store) extra(have map[string]bool) []*Ticket {
	if s.Source == nil {
		return nil
	}
	out, err := s.Source.Extra(have)
	if err != nil {
		return nil
	}
	for _, t := range out {
		t.ReadOnly = true
	}
	return out
}

// Load reads one ticket in full, resolving an exact ID or unambiguous prefix.
func (s *Store) Load(idOrPrefix string) (*Ticket, error) {
	path, err := s.resolve(idOrPrefix)
	if errors.Is(err, ErrNotFound) {
		// Not here as a file — but the board may still know it, and showing a
		// ticket that is one command away is worth more than "not found".
		if t, ok := s.fromSource(idOrPrefix); ok {
			return t, nil
		}
	}
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

// fromSource looks up a ticket the board can see but has no file for, by exact
// id or unambiguous prefix, the same way resolve does for files.
func (s *Store) fromSource(idOrPrefix string) (*Ticket, bool) {
	want := NormalizeIDPrefix(idOrPrefix)
	if want == "" {
		return nil, false
	}
	var match *Ticket
	for _, t := range s.extra(nil) {
		switch {
		case t.ID == want:
			return t, true
		case strings.HasPrefix(t.ID, want), strings.HasSuffix(t.ID, want):
			if match != nil {
				return nil, false // ambiguous: let the caller report not found
			}
			match = t
		}
	}
	return match, match != nil
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
	s.dupIDs = nil
	out := make(map[string]string, len(paths))
	for _, p := range paths {
		id, err := s.idOf(p)
		if err != nil || id == "" {
			continue // unreadable or unidentified files are reported by List
		}
		// Two files declaring the same id is a genuine ambiguity, not something to
		// resolve by whichever happened to be read last: keep the first and let
		// the duplicate be reported rather than silently shadowed.
		if prev, dup := out[id]; dup {
			out[id] = prev // first wins, deterministically
			s.dupIDs = append(s.dupIDs, fmt.Sprintf("%s is declared by both %s and %s",
				Handle(id), filepath.Base(prev), filepath.Base(p)))
			continue
		}
		out[id] = p
	}
	return out, nil
}

// idOf reads just the id from a ticket file.
func (s *Store) idOf(path string) (string, error) {
	d, err := s.readHeader(path)
	if err != nil {
		return "", err
	}
	v, _, err := d.Scalar(FieldID)
	return NormalizeIDPrefix(v), err
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
	if err := WriteAtomic(path, d.Bytes()); err != nil {
		return nil, err
	}
	t, err := Decode(d, path)
	if err != nil {
		return nil, err
	}
	return t, s.record(id, d.Bytes())
}

// Mutate applies fn to a ticket under an exclusive lock, then writes it back
// atomically. Read-modify-write is done inside the lock so two concurrent
// invocations cannot each read the same original and overwrite one another —
// having a single write path prevents schema drift but does nothing about
// interleaving, which is a distinct failure mode (STORE-07).
func (s *Store) Mutate(idOrPrefix string, fn func(*Ticket) error) (*Ticket, error) {
	path, err := s.resolve(idOrPrefix)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// The board may still know this ticket: it is on its ref and has
			// not been pulled. Saying so is the difference between "your id is
			// wrong" and "you are one command away".
			if t, ok := s.fromSource(idOrPrefix); ok {
				return nil, fmt.Errorf("%s: %w", Handle(t.ID), ErrOnRefOnly)
			}
		}
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
	if err := TouchBy(t.doc, s.Actor); err != nil {
		return nil, err
	}
	if err := WriteAtomic(path, t.doc.Bytes()); err != nil {
		return nil, err
	}
	out, err := Decode(t.doc, path)
	if err != nil {
		return nil, err
	}
	// Recorded inside the lock, and only after the file is on disk: the queued
	// bytes are then exactly the bytes a reader of this repository sees, and no
	// second mutation can slip between the write and the record and queue an
	// older state on top of a newer one.
	return out, s.record(id, t.doc.Bytes())
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

// Lock takes the store's advisory lock under a caller-chosen name, for state in
// .jaira that is not one ticket.
//
// The name shares one namespace with the per-ticket locks, which are keyed by
// id, so any name that is not a 26-character ULID is safe — "tags" is not one.
// A name that IS a ULID would silently take that ticket's lock, and nothing here
// checks for it: the callers are in this repository, not a plugin surface.
//
// The caller must call the returned function, normally with defer. Failure to
// acquire within the timeout is an error, not a silent write.
func (s *Store) Lock(name string) (func(), error) { return s.lock(name) }

// WriteAtomic writes via a temporary file in the same directory and renames it
// into place, so a crash or a full disk leaves the previous file intact rather
// than a truncated one (STORE-06).
//
// Exported because the ticket files are not the only thing in .jaira that two
// sessions write at once: the tag registry (core/tag) needs the same guarantee,
// and a second copy of this dance is a second place for it to be subtly wrong.
func WriteAtomic(path string, data []byte) error {
	// Follow a symlink to its target before writing. The tmp-file-plus-rename
	// dance is what makes a write atomic, but rename replaces the link itself —
	// so a symlinked ticket would quietly fork into two diverging files, the link
	// holding new content and the target holding stale content.
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
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
func (s *Store) Save(t *Ticket) error { return WriteAtomic(t.Path, t.doc.Bytes()) }
