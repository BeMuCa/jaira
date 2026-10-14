// Package outbox holds the ticket writes that have not reached the remote yet.
//
// It exists because of a decision about how the tool should feel: claiming a
// ticket works offline. Being offline is rare enough that it must not dictate
// the interaction — so a claim is written locally straight away and the push is
// caught up when there is a network again. Something has to hold the write in
// between, and that is this.
//
// It also settles a problem that would otherwise sit in the write path. Every
// ticket write funnels through core/ticket's Mutate and Create, and those are
// called from the TUI and from the merge driver as well as from the CLI. A
// synchronous push there would hang the board on the network and would have the
// merge driver pushing in the middle of a git merge. With an outbox, the write
// path only files the write; sending happens outside it.
//
// The outbox lives under the per-working-tree state directory, not in the
// repository: an unsent write is this machine's business, and committing it
// would publish it in the very moment it is meant to be waiting.
package outbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Kind is which namespace a queued write belongs to: a ticket, keyed by its
// ulid, or a milestone, keyed by its name.
//
// It is part of the path an entry is filed under, not only a field inside it.
// One flat directory keyed by the bare name would collide the moment somebody
// named a milestone after a ticket handle — and the collision would be silent,
// one unsent write overwriting the other.
type Kind string

const (
	// KindTicket is a ticket's own ref. It is the zero value's meaning too, so
	// an entry queued by an older build — which had no kinds and no
	// subdirectory — still reads back as the ticket write it is.
	KindTicket Kind = "tickets"
	// KindMilestone is a milestone file's ref.
	KindMilestone Kind = "milestones"
)

func (k Kind) or(def Kind) Kind {
	if k == "" {
		return def
	}
	return k
}

// Op is what the queued write does to the ticket's ref.
type Op string

const (
	// OpWrite puts the ticket's current bytes on its ref.
	OpWrite Op = "write"
	// OpDelete removes the ref, for a ticket that has been logged or archived.
	OpDelete Op = "delete"
)

// Entry is one ticket's pending write.
type Entry struct {
	ID  string `json:"id"`
	Op  Op     `json:"op"`
	Ref string `json:"ref"`

	// Kind says which namespace ID belongs to. Absent on an entry written
	// before there were two, which means a ticket.
	Kind Kind `json:"kind,omitempty"`

	// Lease is the ref SHA the remote last confirmed to us, and it is the
	// compare-and-swap value. It survives being superseded: a second local
	// write while the first is still unsent replaces the content but keeps
	// this, because the question the remote answers is "has anyone written
	// since I last saw it", and our own unsent write is not an answer to that.
	Lease string `json:"lease"`

	// Content is the whole ticket file, not a patch. The file is the API, and
	// a queued write that replayed a diff could apply cleanly onto a state
	// nobody ever reviewed.
	Content string `json:"content,omitempty"`

	// File is the file name a milestone is committed under on its ref. Empty
	// on a ticket, and on a milestone entry an older build queued, whose ref
	// carried <key>.md.
	File string `json:"file,omitempty"`

	// Migration is set on a milestone's first write when it moves an older
	// milestone file to its ID: what that file migrated to, and where it came
	// from. It is sent as a commit of its own before Content, whatever edits
	// replaced Content in the meantime.
	Migration *Migration `json:"migration,omitempty"`

	QueuedAt  time.Time `json:"queued-at"`
	UpdatedAt time.Time `json:"updated-at"`
	By        string    `json:"by,omitempty"`
}

// Migration is the move of an older milestone file to its ID, as the commit that
// records it holds it.
type Migration struct {
	File    string `json:"file"`
	Content string `json:"content"`
	// From is the tip of the old name ref the file came from, "none" without one.
	From string `json:"from"`
}

// Age is how long this write has been waiting to be sent.
func (e Entry) Age() time.Duration { return time.Since(e.QueuedAt) }

// Box is the queue for one working tree.
type Box struct{ Dir string }

// At returns the box for a store's working tree.
func At(s *ticket.Store) *Box { return &Box{Dir: filepath.Join(s.StateDir(), "outbox")} }

// path is where one entry is filed: a subdirectory per kind, so a milestone
// named after a ticket handle is a different file rather than the same one.
// Every caller passes a kind it named itself or one that came back normalized
// off disk, so there is nothing to default here.
func (b *Box) path(kind Kind, key string) string {
	if kind == KindMilestone {
		return filepath.Join(b.Dir, milestoneDir, key+".json")
	}
	return filepath.Join(b.Dir, string(kind), key+".json")
}

// milestoneDir is where a milestone write keyed by its ID is filed. It is NOT
// "milestones": an older build lists that folder and sends what it finds under
// <key>.md, which for an ID key is a ghost milestone named after the ID. A
// folder it has never heard of is one it cannot flush.
const milestoneDir = "milestone-ids"

// oldMilestonePath is the folder an older build files name-keyed milestone
// writes in. Still read; nothing new is written there.
func (b *Box) oldMilestonePath(key string) string {
	return filepath.Join(b.Dir, string(KindMilestone), key+".json")
}

// legacyPath is where a ticket entry was filed before there were kinds. Still
// read, so an unsent write queued by an older build is sent rather than
// quietly forgotten the first time the new build runs; never written.
func (b *Box) legacyPath(id string) string { return filepath.Join(b.Dir, id+".json") }

// Queue files a write, superseding any earlier unsent write for the same
// ticket. Superseding is correct rather than lossy: the entry carries the whole
// ticket, so the newer bytes already contain everything the older ones said.
func (b *Box) Queue(id string, op Op, content []byte, lease, by string) error {
	return b.queueKind(KindTicket, id, "", op, content, lease, by)
}

// QueueMilestone files a pending write of a milestone file, under the same
// rules: whole content, the lease the remote last confirmed, superseding any
// earlier unsent write of the same milestone.
func (b *Box) QueueMilestone(key, file string, op Op, content []byte, lease, by string) error {
	return b.queueKind(KindMilestone, key, file, op, content, lease, by)
}

// QueueMigration files the first write of a milestone's ref when it moves an
// older file to its ID: the migrated content, and the old ref's tip (or
// "none"), which are sent as a commit of their own. Later edits replace the
// content and leave the migration as it is.
func (b *Box) QueueMigration(key, file string, migrated []byte, from, by string) error {
	return b.queue(KindMilestone, key, file, OpWrite, migrated, "", by, &Migration{File: file, Content: string(migrated), From: from})
}

// queueKind is what both of those are.
func (b *Box) queueKind(kind Kind, key, file string, op Op, content []byte, lease, by string) error {
	return b.queue(kind, key, file, op, content, lease, by, nil)
}

func (b *Box) queue(kind Kind, key, file string, op Op, content []byte, lease, by string, mig *Migration) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("outbox: no id")
	}
	ref := gitref.RefName(key)
	if kind == KindMilestone {
		ref = gitref.MilestoneRefName(key)
	}
	now := time.Now().UTC()
	e := Entry{
		ID: key, Op: op, Ref: ref, Kind: kind,
		Lease: lease, Content: string(content), File: file,
		QueuedAt: now, UpdatedAt: now, By: by, Migration: mig,
	}
	if old, ok := b.pendingKind(kind, key); ok {
		e.Lease = old.Lease
		e.QueuedAt = old.QueuedAt
		if e.Migration == nil && e.Lease == "" {
			e.Migration = old.Migration
		}
	}
	path := b.path(kind, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	if err := ticket.WriteAtomic(path, append(data, '\n')); err != nil {
		return err
	}
	// An entry an older build filed at the flat path has just been superseded
	// by this one, so it goes. Left behind it would be a second entry for the
	// same ticket: List would return both, the flush would send the stale
	// content first, and the newer write would then carry a lease the remote
	// no longer has.
	if kind == KindTicket {
		if err := os.Remove(b.legacyPath(key)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if kind == KindMilestone {
		if err := os.Remove(b.oldMilestonePath(key)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Pending returns the unsent write for one ticket, if there is one. The board
// asks this to mark a card as carrying something not yet sent.
func (b *Box) Pending(id string) (Entry, bool) { return b.pendingKind(KindTicket, id) }

// PendingMilestone returns the unsent write for one milestone key, if there is
// one, in either folder a milestone write is filed in.
func (b *Box) PendingMilestone(key string) (Entry, bool) { return b.pendingKind(KindMilestone, key) }

// pendingKind reads one entry. A ticket is also looked for at the flat path an
// older build used, so upgrading does not strand a queued write.
func (b *Box) pendingKind(kind Kind, key string) (Entry, bool) {
	e, ok := readEntry(b.path(kind, key))
	if !ok && kind == KindTicket {
		e, ok = readEntry(b.legacyPath(key))
	}
	if !ok && kind == KindMilestone {
		e, ok = readEntry(b.oldMilestonePath(key))
	}
	if !ok {
		return Entry{}, false
	}
	e.Kind = e.Kind.or(KindTicket)
	return e, true
}

func readEntry(path string) (Entry, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, false
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, false
	}
	return e, true
}

// List returns every unsent write, oldest first. Ticket ids are ULIDs, so
// sorting by id is chronological by creation; QueuedAt orders by when the write
// was filed, which is what a flush should follow.
func (b *Box) List() ([]Entry, error) {
	var out []Entry
	// The flat level first: entries an older build left behind, which are
	// tickets by definition.
	flat, err := b.readDir(b.Dir, KindTicket)
	if err != nil {
		return nil, err
	}
	out = append(out, flat...)
	for _, d := range []struct {
		dir  string
		kind Kind
	}{
		{string(KindTicket), KindTicket},
		{milestoneDir, KindMilestone},
		{string(KindMilestone), KindMilestone},
	} {
		got, err := b.readDir(filepath.Join(b.Dir, d.dir), d.kind)
		if err != nil {
			return nil, err
		}
		out = append(out, got...)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].QueuedAt.Equal(out[j].QueuedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].QueuedAt.Before(out[j].QueuedAt)
	})
	return out, nil
}

// readDir reads the entries filed in one directory, as one kind. A missing
// directory is not an error: it means nothing of that kind is waiting.
func (b *Box) readDir(dir string, kind Kind) ([]Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Entry
	for _, de := range entries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".json") {
			continue
		}
		e, ok := readEntry(filepath.Join(dir, de.Name()))
		if !ok {
			continue
		}
		e.Kind = e.Kind.or(kind)
		out = append(out, e)
	}
	return out, nil
}

// Drop removes a queued write.
func (b *Box) Drop(id string) error { return b.dropKind(KindTicket, id) }

// DropMilestone removes a queued write of a milestone.
func (b *Box) DropMilestone(key string) error { return b.dropKind(KindMilestone, key) }

// RekeyMilestone replaces the write an older build queued under oldKey — a
// milestone's name — with one under newKey, its ID, naming the file to commit.
// The new entry keeps when and by whom the old one was queued. If newKey
// already has a write queued, that one stands: it is newer than anything an
// old build could have said, and the old entry is dropped.
func (b *Box) RekeyMilestone(oldKey, newKey, file string, content []byte, lease string) error {
	old, ok := b.pendingKind(KindMilestone, oldKey)
	if !ok {
		return nil
	}
	if _, exists := b.pendingKind(KindMilestone, newKey); exists {
		return b.DropMilestone(oldKey)
	}
	e := Entry{
		ID: newKey, Op: OpWrite, Ref: gitref.MilestoneRefName(newKey), Kind: KindMilestone,
		Lease: lease, Content: string(content), File: file,
		QueuedAt: old.QueuedAt, UpdatedAt: time.Now().UTC(), By: old.By,
	}
	path := b.path(KindMilestone, newKey)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	if err := ticket.WriteAtomic(path, append(data, '\n')); err != nil {
		return err
	}
	return b.DropMilestone(oldKey)
}

// dropKind removes a queued write of either kind, including one an older build
// filed at the flat path.
func (b *Box) dropKind(kind Kind, key string) error {
	if err := os.Remove(b.path(kind, key)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if kind == KindMilestone {
		if err := os.Remove(b.oldMilestonePath(key)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if kind != KindTicket {
		return nil
	}
	if err := os.Remove(b.legacyPath(key)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Sender is the transport a flush pushes through. It is an interface so the
// queue's behaviour — what it keeps, what it drops, when it stops — can be
// tested without a remote, and so nothing in here depends on a network being
// reachable.
type Sender interface {
	Write(id string, content []byte, lease string) (string, error)
	Delete(id, lease string) error
	WriteMilestone(key, file string, content []byte, lease string) (string, error)
}

// migrationSender is a Sender that can write a milestone's migration as a
// commit of its own (gitref.Repo.WriteMigration).
type migrationSender interface {
	WriteMigration(key, file string, migrated []byte, from, editFile string, edit []byte) (string, error)
}

// Outcome is what became of one queued write during a flush.
type Outcome string

const (
	// Sent means the remote accepted it and the entry is gone.
	Sent Outcome = "sent"
	// Rejected means someone else wrote the ticket first. The entry is
	// dropped: replaying it would only lose again, and the local file is still
	// there for the board to merge the ref into.
	Rejected Outcome = "rejected"
	// Unsent means the remote could not be reached. The entry stays.
	Unsent Outcome = "unsent"
	// Failed means git refused for some other reason. The entry stays, because
	// the cause is unknown and discarding a write on an unknown error is the
	// one outcome that loses work.
	Failed Outcome = "failed"
)

// Result reports one entry's fate.
type Result struct {
	ID      string
	Kind    Kind
	Op      Op
	Outcome Outcome
	Err     error
	// File is the file a milestone entry was to be committed under, which is
	// what carries the name a person knows it by; empty for a ticket.
	File string
}

// Flush sends what it can and reports each entry.
//
// It stops at the first unreachable remote instead of walking the rest of the
// queue. The network is a property of the machine, not of the ticket: once one
// push has proved there is no route, every following attempt would wait for the
// same timeout, and a command the user ran for another reason entirely would
// sit there for as long as the queue is.
//
// A lost race is the opposite case and does not stop anything: it is about one
// ticket, and the next one may well go through.
func (b *Box) Flush(s Sender) ([]Result, error) {
	entries, err := b.List()
	if err != nil {
		return nil, err
	}
	var results []Result
	for _, e := range entries {
		kind := e.Kind.or(KindTicket)
		sendErr := send(s, kind, e)

		switch {
		case sendErr == nil:
			if err := b.dropKind(kind, e.ID); err != nil {
				return results, err
			}
			results = append(results, Result{ID: e.ID, File: e.File, Kind: kind, Op: e.Op, Outcome: Sent})
		case errors.Is(sendErr, gitref.ErrRaceLost):
			// A refused milestone write stays queued: its caller merges it into
			// the winner and replaces it, and a failure on the way must not lose
			// the edit.
			if kind != KindMilestone {
				if err := b.dropKind(kind, e.ID); err != nil {
					return results, err
				}
			}
			results = append(results, Result{ID: e.ID, File: e.File, Kind: kind, Op: e.Op, Outcome: Rejected, Err: sendErr})
		case errors.Is(sendErr, gitref.ErrOffline), errors.Is(sendErr, gitref.ErrNoGit):
			results = append(results, Result{ID: e.ID, File: e.File, Kind: kind, Op: e.Op, Outcome: Unsent, Err: sendErr})
			return results, nil
		default:
			results = append(results, Result{ID: e.ID, File: e.File, Kind: kind, Op: e.Op, Outcome: Failed, Err: sendErr})
		}
	}
	return results, nil
}

// send picks the transport call for one entry. An op nothing recognises comes
// back as an error, which Flush reports as Failed and leaves queued: the entry
// is unreadable, not unwanted.
//
// A milestone is only ever written, never deleted: no command deletes one —
// removing a milestone is rm on its file — so there is nothing to queue.
func send(s Sender, kind Kind, e Entry) error {
	if kind == KindMilestone {
		if e.Op != OpWrite {
			return fmt.Errorf("outbox: unknown milestone op %q", e.Op)
		}
		file := e.File
		if file == "" {
			file = e.ID + ".md"
		}
		if m, ok := s.(migrationSender); ok && e.Migration != nil && e.Lease == "" {
			_, err := m.WriteMigration(e.ID, e.Migration.File, []byte(e.Migration.Content), e.Migration.From, file, []byte(e.Content))
			return err
		}
		_, err := s.WriteMilestone(e.ID, file, []byte(e.Content), e.Lease)
		return err
	}
	switch e.Op {
	case OpDelete:
		return s.Delete(e.ID, e.Lease)
	case OpWrite:
		_, err := s.Write(e.ID, []byte(e.Content), e.Lease)
		return err
	}
	return fmt.Errorf("outbox: unknown op %q", e.Op)
}
