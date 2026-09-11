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

func (b *Box) path(id string) string { return filepath.Join(b.Dir, id+".json") }

// Queue files a write, superseding any earlier unsent write for the same
// ticket. Superseding is correct rather than lossy: the entry carries the whole
// ticket, so the newer bytes already contain everything the older ones said.
func (b *Box) Queue(id string, op Op, content []byte, lease, by string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("outbox: no ticket id")
	}
	now := time.Now().UTC()
	e := Entry{
		ID: id, Op: op, Ref: gitref.RefName(id),
		Lease: lease, Content: string(content),
		QueuedAt: now, UpdatedAt: now, By: by,
	}
	if old, ok := b.Pending(id); ok {
		e.Lease = old.Lease
		e.QueuedAt = old.QueuedAt
	}
	if err := os.MkdirAll(b.Dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return ticket.WriteAtomic(b.path(id), append(data, '\n'))
}

// Pending returns the unsent write for one ticket, if there is one. The board
// asks this to mark a card as carrying something not yet sent.
func (b *Box) Pending(id string) (Entry, bool) {
	data, err := os.ReadFile(b.path(id))
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
	entries, err := os.ReadDir(b.Dir)
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
		id := strings.TrimSuffix(de.Name(), ".json")
		if e, ok := b.Pending(id); ok {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].QueuedAt.Equal(out[j].QueuedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].QueuedAt.Before(out[j].QueuedAt)
	})
	return out, nil
}

// Drop removes a queued write.
func (b *Box) Drop(id string) error {
	err := os.Remove(b.path(id))
	if err != nil && !os.IsNotExist(err) {
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
		var sendErr error
		switch e.Op {
		case OpDelete:
			sendErr = s.Delete(e.ID, e.Lease)
		case OpWrite:
			_, sendErr = s.Write(e.ID, []byte(e.Content), e.Lease)
		default:
			results = append(results, Result{ID: e.ID, Op: e.Op, Outcome: Failed,
				Err: fmt.Errorf("outbox: unknown op %q", e.Op)})
			continue
		}

		switch {
		case sendErr == nil:
			if err := b.Drop(e.ID); err != nil {
				return results, err
			}
			results = append(results, Result{ID: e.ID, Op: e.Op, Outcome: Sent})
		case errors.Is(sendErr, gitref.ErrRaceLost):
			if err := b.Drop(e.ID); err != nil {
				return results, err
			}
			results = append(results, Result{ID: e.ID, Op: e.Op, Outcome: Rejected, Err: sendErr})
		case errors.Is(sendErr, gitref.ErrOffline), errors.Is(sendErr, gitref.ErrNoGit):
			results = append(results, Result{ID: e.ID, Op: e.Op, Outcome: Unsent, Err: sendErr})
			return results, nil
		default:
			results = append(results, Result{ID: e.ID, Op: e.Op, Outcome: Failed, Err: sendErr})
		}
	}
	return results, nil
}
