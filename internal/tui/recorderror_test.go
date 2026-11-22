package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// A milestone the board's background fetch wrote here but could not queue for
// its ref reaches nobody, and the fetch retries every few minutes. 'jaira
// fetch' warns about it; the board set the error to nil and said nothing, so
// the person believed the milestone was on its way to their teammates.
func TestTheBoardSaysWhenAMilestoneCouldNotBeQueued(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	if os.Geteuid() == 0 {
		t.Skip("root writes into a read-only directory")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "bob")
	run(t, root, "git", "clone", "--quiet", bare, clone)
	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.Actor = "bob"
	y := refsync.New(s, "origin", "bob")
	s.Recorder = y
	s.Source = y

	// A milestone file an older build left under its name: the fetch moves
	// it to its ID name and queues its ref, and the queueing is what fails.
	ms := milestone.New("backend", 45, nowish())
	if err := os.MkdirAll(milestone.Dir(clone), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(milestone.Dir(clone), "backend.md"), ms.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	box := filepath.Join(s.StateDir(), "outbox")
	if err := os.MkdirAll(box, 0o555); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(box, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(box, 0o755) })

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	msg := fetchRefs(y)()
	got, ok := msg.(refFetchedMsg)
	if !ok || got.err != nil {
		t.Fatalf("background run: %#v", msg)
	}
	// Typing a ticket's title when the background fetch comes back: the
	// warning must not take the keys, or Enter closes it instead of creating
	// the ticket.
	m.mode = modeCreate
	m.Update(msg)
	if line := m.flashLine(); !strings.Contains(line, `milestone "backend"`) || !strings.Contains(line, "could not be queued") {
		t.Errorf("the status line is %q; want the milestone it could not queue", line)
	}
	if m.mode != modeCreate {
		t.Errorf("the warning switched the board to mode %v; the person was typing a ticket", m.mode)
	}
}
