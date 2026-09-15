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

	QueuedAt  time.Time `json:"queued-at"`
	UpdatedAt time.Time `json:"updated-at"`
	By        string    `json:"by,omitempty"`
}

// Age is how long this write has been waiting to be sent.
func (e Entry) Age() time.Duration { return time.Since(e.QueuedAt) }

// Box is the queue for one working tree.
type Box struct{ Dir string }

// At returns the box for a store's working tree.
func At(s *ticket.Store) *Box { return &Box{Dir: filepath.Join(s.StateDir(), "outbox")} }

// path is where one entry is filed: a subdirectory per kind, so a milestone
// named after a ticket handle is a different file rather than the same one.
func (b *Box) path(kind Kind, key string) string {
	return filepath.Join(b.Dir, string(kind.or(KindTicket)), key+".json")
}

// legacyPath is where a ticket entry was filed before there were kinds. Still
// read, so an unsent write queued by an older build is sent rather than
// quietly forgotten the first time the new build runs; never written.
func (b *Box) legacyPath(id string) string { return filepath.Join(b.Dir, id+".json") }

// Queue files a write, superseding any earlier unsent write for the same
// ticket. Superseding is correct rather than lossy: the entry carries the whole
// ticket, so the newer bytes already contain everything the older ones said.
func (b *Box) Queue(id string, op Op, content []byte, lease, by string) error {
	return b.QueueKind(KindTicket, id, op, content, lease, by)
}

// QueueMilestone files a pending write of a milestone file, under the same
// rules: whole content, the lease the remote last confirmed, superseding any
// earlier unsent write of the same milestone.
func (b *Box) QueueMilestone(name string, op Op, content []byte, lease, by string) error {
	return b.QueueKind(KindMilestone, name, op, content, lease, by)
}

// QueueKind is what both of those are.
func (b *Box) QueueKind(kind Kind, key string, op Op, content []byte, lease, by string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("outbox: no id")
	}
	kind = kind.or(KindTicket)
	ref := gitref.RefName(key)
	if kind == KindMilestone {
		ref = gitref.MilestoneRefName(key)
	}
	now := time.Now().UTC()
	e := Entry{
		ID: key, Op: op, Ref: ref, Kind: kind,
		Lease: lease, Content: string(content),
		QueuedAt: now, UpdatedAt: now, By: by,
	}
	if old, ok := b.PendingKind(kind, key); ok {
		e.Lease = old.Lease
		e.QueuedAt = old.QueuedAt
	}
	path := b.path(kind, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return ticket.WriteAtomic(path, append(data, '\n'))
}

// Pending returns the unsent write for one ticket, if there is one. The board
// asks this to mark a card as carrying something not yet sent.
func (b *Box) Pending(id string) (Entry, bool) { return b.PendingKind(KindTicket, id) }

// PendingKind reads one entry. A ticket is also looked for at the flat path an
// older build used, so upgrading does not strand a queued write.
func (b *Box) PendingKind(kind Kind, key string) (Entry, bool) {
	kind = kind.or(KindTicket)
	e, ok := readEntry(b.path(kind, key))
	if !ok && kind == KindTicket {
		e, ok = readEntry(b.legacyPath(key))
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
	for _, kind := range []Kind{KindTicket, KindMilestone} {
		got, err := b.readDir(filepath.Join(b.Dir, string(kind)), kind)
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
func (b *Box) Drop(id string) error { return b.DropKind(KindTicket, id) }

// DropKind removes a queued write of either kind, including one an older build
// filed at the flat path.
func (b *Box) DropKind(kind Kind, key string) error {
	for _, path := range []string{b.path(kind, key), b.legacyPath(key)} {
		if kind != KindTicket && path == b.legacyPath(key) {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
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
	WriteMilestone(name string, content []byte, lease string) (string, error)
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
			if err := b.DropKind(kind, e.ID); err != nil {
				return results, err
			}
			results = append(results, Result{ID: e.ID, Kind: kind, Op: e.Op, Outcome: Sent})
		case errors.Is(sendErr, gitref.ErrRaceLost):
			if err := b.DropKind(kind, e.ID); err != nil {
				return results, err
			}
			results = append(results, Result{ID: e.ID, Kind: kind, Op: e.Op, Outcome: Rejected, Err: sendErr})
		case errors.Is(sendErr, gitref.ErrOffline), errors.Is(sendErr, gitref.ErrNoGit):
			results = append(results, Result{ID: e.ID, Kind: kind, Op: e.Op, Outcome: Unsent, Err: sendErr})
			return results, nil
		default:
			results = append(results, Result{ID: e.ID, Kind: kind, Op: e.Op, Outcome: Failed, Err: sendErr})
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
		_, err := s.WriteMilestone(e.ID, []byte(e.Content), e.Lease)
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
