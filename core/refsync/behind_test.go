package refsync_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// sentTicket files a ticket on ada's side and sends it, so the file and the
// ref hold the same bytes. It returns the id and the file read back as a
// listing reads it.
func sentTicket(t *testing.T, ada side, title string) (string, *ticket.Ticket) {
	t.Helper()
	id := create(t, ada, title)
	// A write stamps updated-at, which create leaves out and the verdict needs.
	if _, err := ada.store.Mutate(id, func(*ticket.Ticket) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	f, err := ada.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	return id, f
}

// moveRef puts a new lane on the ref the way another worktree's move in the
// same second does: the stamp stays.
func moveRef(t *testing.T, ada side, id string) {
	t.Helper()
	content, sha, err := ada.syncer.Repo.Read(id)
	if err != nil {
		t.Fatal(err)
	}
	moved := strings.Replace(string(content), "\nstatus: backlog\n", "\nstatus: todo\n", 1)
	if moved == string(content) {
		t.Fatalf("no status line to move:\n%s", content)
	}
	if _, err := ada.syncer.Repo.Write(id, []byte(moved), sha); err != nil {
		t.Fatal(err)
	}
}

// gitCalls counts the git commands named in a GIT_TRACE log since it was last
// truncated.
func gitCalls(t *testing.T, trace string, names ...string) int {
	t.Helper()
	b, _ := os.ReadFile(trace)
	n := 0
	for _, line := range strings.Split(string(b), "\n") {
		for _, name := range names {
			if strings.Contains(line, "git "+name) || strings.Contains(line, "'"+name+"'") {
				n++
			}
		}
	}
	return n
}

// K6BGB0 review: a tie in updated-at with different bytes asks git whether the
// ref carried the file before, and the board asks on every two-second tick. The
// answer is cached, so a board that is not changing runs no git for it.
func TestATiedFileIsJudgedOnceNotOnEveryListing(t *testing.T) {
	ada, _, _ := twoSides(t)
	id, file := sentTicket(t, ada, "tied")
	moveRef(t, ada, id)

	trace := filepath.Join(t.TempDir(), "trace")
	t.Setenv("GIT_TRACE", trace)
	if got := ada.syncer.Ahead([]*ticket.Ticket{file}); got[id] == nil || got[id].Status != "todo" {
		t.Fatalf("a file behind its ref in the same second is not ahead: %v", got)
	}
	if n := gitCalls(t, trace, "rev-list"); n == 0 {
		t.Fatal("the first judgement did not ask git; the trace counts nothing and this test proves nothing")
	}
	if err := os.WriteFile(trace, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ada.syncer.Ahead([]*ticket.Ticket{file}); got[id] == nil {
		t.Fatal("the second listing lost the verdict")
	}
	if n := gitCalls(t, trace, "rev-list", "hash-object"); n != 0 {
		t.Errorf("the second listing asked git again: %d calls", n)
	}
}

// A checkout with core.autocrlf holds every ticket with CRLF while the ref
// holds LF. That is the same ticket, not one behind its ref, and it costs no
// git to say so.
func TestALineEndingDifferenceIsNotBehind(t *testing.T) {
	ada, _, _ := twoSides(t)
	id, _ := sentTicket(t, ada, "crlf")
	content, _, err := ada.syncer.Repo.Read(id)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ticket.ParseDoc([]byte(strings.ReplaceAll(string(content), "\n", "\r\n")))
	if err != nil {
		t.Fatal(err)
	}
	file, err := ticket.Decode(d, "")
	if err != nil {
		t.Fatal(err)
	}
	trace := filepath.Join(t.TempDir(), "trace")
	t.Setenv("GIT_TRACE", trace)
	if got := ada.syncer.Ahead([]*ticket.Ticket{file}); got[id] != nil {
		t.Error("a CRLF copy of the ref is shown as behind it")
	}
	if n := gitCalls(t, trace, "rev-list", "hash-object"); n != 0 {
		t.Errorf("a line-ending difference asked git: %d calls", n)
	}
}

// A write still queued here while the ref moved on is not "behind: pull it":
// pull refuses it (ErrUnsent). The card keeps the file, and the board's unsent
// mark says what is going on.
func TestAnUnsentWriteIsNotShownAsBehind(t *testing.T) {
	ada, _, _ := twoSides(t)
	id, _ := sentTicket(t, ada, "unsent")
	tk, err := ada.store.Mutate(id, func(tk *ticket.Ticket) error {
		return ticket.AppendNote(tk.Doc(), "not sent yet", tk.UpdatedAt, "ada")
	})
	if err != nil {
		t.Fatalf("mutate: %v", err)
	}
	// The ref moves on later than the local write, as a teammate's would.
	content, sha, err := ada.syncer.Repo.Read(id)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "updated-at:") {
			line = "updated-at: 2099-01-01T00:00:00Z"
		}
		lines = append(lines, line)
	}
	if _, err := ada.syncer.Repo.Write(id, []byte(strings.Join(lines, "\n")), sha); err != nil {
		t.Fatal(err)
	}
	if _, pending := ada.syncer.Pending(id); !pending {
		t.Fatal("the local write is not queued; this test proves nothing")
	}
	if got := ada.syncer.Ahead([]*ticket.Ticket{tk}); got[id] != nil {
		t.Error("an unsent write is shown as behind its ref, though pull will refuse it")
	}
	if _, err := ada.syncer.Pull(id, false); !errors.Is(err, refsync.ErrUnsent) {
		t.Errorf("pull over the unsent write: %v, want ErrUnsent", err)
	}
}
