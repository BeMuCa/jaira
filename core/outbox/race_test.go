package outbox_test

import (
	"fmt"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
)

// interleaving is a Sender that, while the push of one key is in flight,
// queues a second write for that key — what another session does when it runs
// 'jaira epic add' while this one's flush pushes the same epic.
type interleaving struct {
	*fakeSender
	during func()
}

func (s *interleaving) mid() {
	if s.during != nil {
		s.during()
		s.during = nil
	}
}

func (s *interleaving) Write(id string, c []byte, lease string) (string, error) {
	s.mid()
	return s.fakeSender.Write(id, c, lease)
}

func (s *interleaving) Delete(id, lease string) error {
	s.mid()
	return s.fakeSender.Delete(id, lease)
}

func (s *interleaving) WriteMilestone(k, f string, c []byte, lease string) (string, error) {
	s.mid()
	return s.fakeSender.WriteMilestone(k, f, c, lease)
}

func (s *interleaving) WriteEpic(n string, c []byte, lease string) (string, error) {
	s.mid()
	return s.fakeSender.WriteEpic(n, c, lease)
}

type raceKind struct {
	name    string
	queue   func(b *outbox.Box, key, content, lease string) error
	pending func(b *outbox.Box, key string) (outbox.Entry, bool)
}

var raceKinds = []raceKind{
	{"ticket",
		func(b *outbox.Box, k, c, l string) error { return b.Queue(k, outbox.OpWrite, []byte(c), l, "ada") },
		(*outbox.Box).Pending},
	{"milestone",
		func(b *outbox.Box, k, c, l string) error {
			return b.QueueMilestone(k, "m-"+k+".md", outbox.OpWrite, []byte(c), l, "ada")
		},
		(*outbox.Box).PendingMilestone},
	{"epic",
		func(b *outbox.Box, k, c, l string) error { return b.QueueEpic(k, []byte(c), l, "ada") },
		(*outbox.Box).PendingEpic},
}

// A write queued for a key while that key's push is in flight survives the
// flush, for every kind: it is newer than what was pushed, and dropping the
// entry by key would lose it — for an epic it is the only copy of the edit.
// Once the push went through, it is leased on the SHA that push made, so the
// next flush sends it rather than losing a race against our own write.
func TestAWriteQueuedDuringThePushOfItsKeySurvivesTheFlush(t *testing.T) {
	for _, k := range raceKinds {
		t.Run(k.name, func(t *testing.T) {
			b := box(t)
			const key = "01AAA"
			if err := k.queue(b, key, "first", "lease-remote"); err != nil {
				t.Fatal(err)
			}
			s := &interleaving{fakeSender: newFake(), during: func() {
				if err := k.queue(b, key, "second", "lease-remote"); err != nil {
					t.Errorf("queue during the push: %v", err)
				}
			}}
			res, err := b.Flush(s)
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 1 || res[0].Outcome != outbox.Sent {
				t.Fatalf("want the first write sent, got %+v", res)
			}
			e, ok := k.pending(b, key)
			if !ok {
				t.Fatal("the write queued during the push is gone from the outbox")
			}
			if e.Content != "second" || e.Lease != "sha-"+key {
				t.Fatalf("queued entry is %q leased on %q, want %q leased on the pushed SHA", e.Content, e.Lease, "second")
			}

			if _, err := b.Flush(s); err != nil {
				t.Fatal(err)
			}
			if s.seen[key] != "sha-"+key {
				t.Errorf("the second write was pushed with lease %q", s.seen[key])
			}
			if _, ok := k.pending(b, key); ok {
				t.Error("the second write is still queued after it was sent")
			}
		})
	}
}

// A write queued for a ticket while its delete is in flight survives the
// flush leased on no ref: the remote has none after the delete, and the
// deleted ref's SHA would make the next flush lose its push to nothing.
func TestAWriteQueuedDuringItsDeleteIsLeasedOnNoRef(t *testing.T) {
	b := box(t)
	const key = "01AAA"
	if err := b.Queue(key, outbox.OpDelete, nil, "lease-remote", "ada"); err != nil {
		t.Fatal(err)
	}
	s := &interleaving{fakeSender: newFake(), during: func() {
		if err := b.Queue(key, outbox.OpWrite, []byte("reborn"), "lease-remote", "ada"); err != nil {
			t.Errorf("queue during the delete: %v", err)
		}
	}}
	if res, err := b.Flush(s); err != nil || len(res) != 1 || res[0].Outcome != outbox.Sent {
		t.Fatalf("want the delete sent, got %+v, %v", res, err)
	}
	e, ok := b.Pending(key)
	if !ok || e.Op != outbox.OpWrite || e.Content != "reborn" || e.Lease != "" {
		t.Fatalf("want the write queued on no ref, got %+v, %v", e, ok)
	}
}

// A ticket write the remote refused is dropped, but only that one: a write
// queued during the refused push is a newer edit and stays queued.
func TestALostRaceDropsOnlyTheRefusedWrite(t *testing.T) {
	b := box(t)
	if err := b.Queue("01AAA", outbox.OpWrite, []byte("first"), "lease", "ada"); err != nil {
		t.Fatal(err)
	}
	s := &interleaving{fakeSender: newFake(), during: func() {
		if err := b.Queue("01AAA", outbox.OpWrite, []byte("second"), "lease", "ada"); err != nil {
			t.Errorf("queue during the push: %v", err)
		}
	}}
	s.err["01AAA"] = fmt.Errorf("%w: stale info", gitref.ErrRaceLost)
	if _, err := b.Flush(s); err != nil {
		t.Fatal(err)
	}
	e, ok := b.Pending("01AAA")
	if !ok || e.Content != "second" {
		t.Fatalf("the write queued during the refused push is gone: %+v, %v", e, ok)
	}
}

// Settle and Replace act on the entry the caller read, never on a newer one
// filed under the same key in the meantime.
func TestSettleAndReplaceLeaveANewerEntryAlone(t *testing.T) {
	for _, k := range raceKinds {
		t.Run(k.name, func(t *testing.T) {
			b := box(t)
			if err := k.queue(b, "01AAA", "first", "lease-1"); err != nil {
				t.Fatal(err)
			}
			seen, _ := k.pending(b, "01AAA")
			if err := k.queue(b, "01AAA", "second", "lease-1"); err != nil {
				t.Fatal(err)
			}
			if err := b.Settle(seen, ""); err != nil {
				t.Fatal(err)
			}
			if err := b.Replace(seen, []byte("merged"), "lease-2", "ada"); err != nil {
				t.Fatal(err)
			}
			e, ok := k.pending(b, "01AAA")
			if !ok || e.Content != "second" || e.Lease != "lease-1" {
				t.Fatalf("the newer entry was touched: %+v, %v", e, ok)
			}

			// The entry read is still there: Replace swaps it, lease included.
			if err := b.Replace(e, []byte("merged"), "lease-2", "ada"); err != nil {
				t.Fatal(err)
			}
			got, _ := k.pending(b, "01AAA")
			if got.Content != "merged" || got.Lease != "lease-2" || !got.QueuedAt.Equal(e.QueuedAt) {
				t.Fatalf("replace gave %+v", got)
			}
			if err := b.Settle(got, ""); err != nil {
				t.Fatal(err)
			}
			if _, ok := k.pending(b, "01AAA"); ok {
				t.Error("settling the entry that was read left it queued")
			}
		})
	}
}
