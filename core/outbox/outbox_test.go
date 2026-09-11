package outbox_test

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
)

// fakeSender answers per ticket, so a test can stage exactly the case it is
// about — a lost race on one ticket, no network on another — without a remote.
type fakeSender struct {
	err   map[string]error
	wrote []string
	seen  map[string]string // id -> lease it was pushed with
}

func newFake() *fakeSender {
	return &fakeSender{err: map[string]error{}, seen: map[string]string{}}
}

func (f *fakeSender) Write(id string, content []byte, lease string) (string, error) {
	f.seen[id] = lease
	if err := f.err[id]; err != nil {
		return "", err
	}
	f.wrote = append(f.wrote, id+":"+string(content))
	return "sha-" + id, nil
}

func (f *fakeSender) Delete(id, lease string) error {
	f.seen[id] = lease
	if err := f.err[id]; err != nil {
		return err
	}
	f.wrote = append(f.wrote, id+":deleted")
	return nil
}

func box(t *testing.T) *outbox.Box {
	t.Helper()
	return &outbox.Box{Dir: filepath.Join(t.TempDir(), "outbox")}
}

func TestAWriteWaitsUntilItIsSent(t *testing.T) {
	b := box(t)
	if err := b.Queue("01AAA", outbox.OpWrite, []byte("ticket bytes"), "lease-1", "ada"); err != nil {
		t.Fatalf("queue: %v", err)
	}
	e, ok := b.Pending("01AAA")
	if !ok {
		t.Fatal("the queued write is not pending")
	}
	if e.Lease != "lease-1" || e.Content != "ticket bytes" || e.By != "ada" {
		t.Errorf("entry came back changed: %+v", e)
	}
	if e.Ref != "refs/jaira/tickets/01AAA" {
		t.Errorf("ref name is %q", e.Ref)
	}

	f := newFake()
	res, err := b.Flush(f)
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(res) != 1 || res[0].Outcome != outbox.Sent {
		t.Fatalf("want one sent, got %+v", res)
	}
	if _, ok := b.Pending("01AAA"); ok {
		t.Error("a sent write is still queued")
	}
}

// A second local write while the first is unsent must keep the lease the remote
// confirmed. Taking the newer write's lease would compare the remote against a
// SHA it has never seen, and the push would be refused for the wrong reason.
func TestSupersedingAWriteKeepsTheConfirmedLease(t *testing.T) {
	b := box(t)
	if err := b.Queue("01AAA", outbox.OpWrite, []byte("first"), "lease-remote", "ada"); err != nil {
		t.Fatalf("queue: %v", err)
	}
	first, _ := b.Pending("01AAA")
	if err := b.Queue("01AAA", outbox.OpWrite, []byte("second"), "lease-local", "ada"); err != nil {
		t.Fatalf("re-queue: %v", err)
	}
	e, _ := b.Pending("01AAA")
	if e.Lease != "lease-remote" {
		t.Errorf("lease became %q, want the confirmed one", e.Lease)
	}
	if e.Content != "second" {
		t.Errorf("content is %q, want the newer bytes", e.Content)
	}
	if !e.QueuedAt.Equal(first.QueuedAt) {
		t.Error("superseding restarted the wait; the write has been waiting since the first one")
	}

	f := newFake()
	if _, err := b.Flush(f); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if f.seen["01AAA"] != "lease-remote" {
		t.Errorf("pushed with lease %q", f.seen["01AAA"])
	}
	if len(f.wrote) != 1 || !strings.HasSuffix(f.wrote[0], ":second") {
		t.Errorf("sent %v, want only the newer bytes once", f.wrote)
	}
}

// No network stops the flush where it stands. Walking the rest of the queue
// would make an unrelated command wait for one timeout per ticket.
func TestNoNetworkStopsTheFlushAndKeepsEverything(t *testing.T) {
	b := box(t)
	for _, id := range []string{"01AAA", "01BBB", "01CCC"} {
		if err := b.Queue(id, outbox.OpWrite, []byte(id), "lease", "ada"); err != nil {
			t.Fatalf("queue %s: %v", id, err)
		}
	}
	f := newFake()
	f.err["01BBB"] = fmt.Errorf("%w: could not resolve host", gitref.ErrOffline)

	res, err := b.Flush(f)
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("flush did not stop at the unreachable remote: %+v", res)
	}
	if res[1].Outcome != outbox.Unsent {
		t.Errorf("second entry is %q, want unsent", res[1].Outcome)
	}
	if _, ok := b.Pending("01BBB"); !ok {
		t.Error("an unsent write was dropped")
	}
	if _, ok := b.Pending("01CCC"); !ok {
		t.Error("a write the flush never reached was dropped")
	}
	if _, ok := b.Pending("01AAA"); ok {
		t.Error("the write that did go through is still queued")
	}
}

// A lost race is about one ticket and stops nothing. The entry goes, because
// replaying it would lose again; the local file stays, and the board merges the
// ref in.
func TestALostRaceDropsThatEntryAndCarriesOn(t *testing.T) {
	b := box(t)
	for _, id := range []string{"01AAA", "01BBB"} {
		if err := b.Queue(id, outbox.OpWrite, []byte(id), "lease", "ada"); err != nil {
			t.Fatalf("queue %s: %v", id, err)
		}
	}
	f := newFake()
	f.err["01AAA"] = fmt.Errorf("%w: stale info", gitref.ErrRaceLost)

	res, err := b.Flush(f)
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("want both entries attempted, got %+v", res)
	}
	if res[0].Outcome != outbox.Rejected || !errors.Is(res[0].Err, gitref.ErrRaceLost) {
		t.Errorf("first entry is %+v, want rejected as a lost race", res[0])
	}
	if _, ok := b.Pending("01AAA"); ok {
		t.Error("a rejected write is still queued and would lose again")
	}
	if res[1].Outcome != outbox.Sent {
		t.Errorf("second entry is %q; one ticket's race must not stop the queue", res[1].Outcome)
	}
}

// An unknown git failure keeps the write. Discarding on an error nobody has
// classified is the only outcome that actually loses work.
func TestAnUnknownFailureKeepsTheWrite(t *testing.T) {
	b := box(t)
	if err := b.Queue("01AAA", outbox.OpWrite, []byte("x"), "lease", "ada"); err != nil {
		t.Fatalf("queue: %v", err)
	}
	f := newFake()
	f.err["01AAA"] = errors.New("git: something nobody classified")

	res, err := b.Flush(f)
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(res) != 1 || res[0].Outcome != outbox.Failed {
		t.Fatalf("want one failed, got %+v", res)
	}
	if _, ok := b.Pending("01AAA"); !ok {
		t.Error("an unclassified failure discarded the write")
	}
}

func TestADeleteIsQueuedLikeAWrite(t *testing.T) {
	b := box(t)
	if err := b.Queue("01AAA", outbox.OpDelete, nil, "lease", "ada"); err != nil {
		t.Fatalf("queue: %v", err)
	}
	f := newFake()
	if _, err := b.Flush(f); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(f.wrote) != 1 || f.wrote[0] != "01AAA:deleted" {
		t.Errorf("sent %v", f.wrote)
	}
}

func TestListIsOldestFirstAndEmptyBoxIsNotAnError(t *testing.T) {
	b := box(t)
	got, err := b.List()
	if err != nil || got != nil {
		t.Fatalf("empty box: %v %v", got, err)
	}
	for _, id := range []string{"01CCC", "01AAA"} {
		if err := b.Queue(id, outbox.OpWrite, []byte(id), "lease", "ada"); err != nil {
			t.Fatalf("queue %s: %v", id, err)
		}
	}
	got, err = b.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 entries, got %d", len(got))
	}
	if got[0].QueuedAt.After(got[1].QueuedAt) {
		t.Errorf("newest first: %s queued %s before %s queued %s",
			got[0].ID, got[0].QueuedAt, got[1].ID, got[1].QueuedAt)
	}
	// Two writes filed in the same instant fall back to the id, so the order
	// is total rather than dependent on clock resolution.
	if got[0].QueuedAt.Equal(got[1].QueuedAt) && got[0].ID > got[1].ID {
		t.Errorf("same instant, ids out of order: %s before %s", got[0].ID, got[1].ID)
	}
}

// The case the decision was made for, end to end against real git: a write made
// while the remote is unreachable waits, and goes through unchanged once there
// is a route again.
func TestAnOfflineWriteArrivesWhenTheNetworkIsBack(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "ada")
	run(t, root, "git", "clone", "--quiet", bare, clone)

	repo := &gitref.Repo{Dir: clone, Remote: "origin", AuthorName: "ada", AuthorEmail: "ada@example.test"}
	b := box(t)

	// Offline: port 1 on the loopback refuses at once, so the test covers the
	// unreachable-remote case without waiting for a real network timeout.
	run(t, clone, "git", "remote", "set-url", "origin", "https://127.0.0.1:1/board.git")
	const content = "---\nid: 01AAA\nassignee: ada\n---\n\n# claimed offline\n"
	if _, err := repo.Write("01AAA", []byte(content), ""); !errors.Is(err, gitref.ErrOffline) {
		t.Fatalf("want ErrOffline while offline, got %v", err)
	}
	if err := b.Queue("01AAA", outbox.OpWrite, []byte(content), "", "ada"); err != nil {
		t.Fatalf("queue: %v", err)
	}
	res, err := b.Flush(repo)
	if err != nil {
		t.Fatalf("flush while offline: %v", err)
	}
	if len(res) != 1 || res[0].Outcome != outbox.Unsent {
		t.Fatalf("want unsent while offline, got %+v", res)
	}

	// The network comes back.
	run(t, clone, "git", "remote", "set-url", "origin", bare)
	res, err = b.Flush(repo)
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(res) != 1 || res[0].Outcome != outbox.Sent {
		t.Fatalf("want sent once the remote is reachable, got %+v", res)
	}
	if _, ok := b.Pending("01AAA"); ok {
		t.Error("the write is still queued after being sent")
	}
	got, _, err := repo.Read("01AAA")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(got) != content {
		t.Errorf("the ticket arrived changed:\n%q", got)
	}
}

func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.test",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.test",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}
