package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// refClone is a clone of a shared remote with a board in it, wired to its refs
// the way the CLI wires every store it opens.
func refClone(t *testing.T, root, bare, who string) (*ticket.Store, *refsync.Syncer) {
	t.Helper()
	clone := filepath.Join(root, who)
	run(t, root, "git", "clone", "--quiet", bare, clone)
	run(t, clone, "git", "config", "user.name", who)
	run(t, clone, "git", "config", "user.email", who+"@example.test")
	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.Actor = who
	y := refsync.New(s, "origin", who)
	s.Recorder = y
	s.Source = y
	return s, y
}

func createTicket(t *testing.T, s *ticket.Store, title, status, board string) string {
	t.Helper()
	now := time.Now()
	id := ticket.NewID(now)
	fields := map[string]string{
		ticket.FieldID: id, ticket.FieldTitle: title, ticket.FieldStatus: status,
		ticket.FieldCreatedAt: ticket.FormatTime(now), ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	if board != "" {
		fields[ticket.FieldBoard] = board
	}
	if _, err := s.Create(fields, nil, "# "+title+"\n"); err != nil {
		t.Fatalf("create %q: %v", title, err)
	}
	return id
}

// A ticket stays on its ref until it is filed away, so most of a board's
// tickets have no file in a given clone. The launcher read only the files and
// showed such a board as "0 open / 0" while 'jaira list' listed every ticket.
func TestHomeCountsTicketsThatLiveOnlyOnTheirRef(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)

	ada, adaSync := refClone(t, root, bare, "ada")
	createTicket(t, ada, "open on the ref", "backlog", "")
	createTicket(t, ada, "finished on the ref", "done", "")
	if _, err := adaSync.Flush(); err != nil {
		t.Fatal(err)
	}

	bob, bobSync := refClone(t, root, bare, "bob")
	if _, err := bobSync.Incoming(); err != nil {
		t.Fatal(err)
	}
	// The case from the video project: refs, and not a single ticket file.
	if e := describe(bob.Root, catalogue(t), true); e.Err != nil || e.Total != 2 || e.Open != 1 {
		t.Fatalf("a board with refs and no files: %d open / %d (err %v), want 1 open / 2", e.Open, e.Total, e.Err)
	}

	// A ticket that is a file here and on its ref as well is one ticket.
	createTicket(t, bob, "a file and a ref", "backlog", "")
	if _, err := bobSync.Flush(); err != nil {
		t.Fatal(err)
	}
	e := describe(bob.Root, catalogue(t), true)
	if e.Err != nil || e.Total != 3 || e.Open != 2 {
		t.Errorf("got %d open / %d (err %v), want 2 open / 3", e.Open, e.Total, e.Err)
	}
	// And the same number 'jaira list' gives, which reads through this store.
	listed, err := bob.List()
	if err != nil {
		t.Fatal(err)
	}
	if e.Total != len(listed) {
		t.Errorf("launcher counts %d, the store lists %d", e.Total, len(listed))
	}
}

// A ticket runs through the lanes of its own board. Counted against today's
// lanes, a ticket in a named board's terminal lane was open forever, and one
// in its human lane was never waiting.
func TestHomeCountsATicketAgainstItsOwnBoard(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	dir := makeBoard(t, root, "repo")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	lanes := filepath.Join(dir, ticket.DirName, "boards", "video", "lanes")
	if err := os.MkdirAll(lanes, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"script.md":    "---\nid: script\nname: Skript\nprecedence: 10\n---\n",
		"approve.md":   "---\nid: approve\nname: Abnahme\nafter: script\nprecedence: 20\nrequires-human-exit: true\n---\n",
		"published.md": "---\nid: published\nname: Veroeffentlicht\nafter: approve\nprecedence: 30\nterminal: true\n---\n",
	} {
		if err := os.WriteFile(filepath.Join(lanes, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	createTicket(t, s, "being written", "script", "video")
	createTicket(t, s, "waiting for a person", "approve", "video")
	createTicket(t, s, "out", "published", "video")
	createTicket(t, s, "today's board", "backlog", "")

	e := describe(dir, catalogue(t), true)
	if e.Err != nil || e.Total != 4 || e.Open != 3 || e.Waiting != 1 {
		t.Errorf("got %d open / %d, %d waiting (err %v), want 3 open / 4, 1 waiting", e.Open, e.Total, e.Waiting, e.Err)
	}
}

// catalogue is the lane set the launcher counts today's board against.
func catalogue(t *testing.T) *lane.Set {
	t.Helper()
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return lanes
}

// Counting reads every board's refs, so it runs after the first frame. Until
// it is in, a row says so instead of showing the files on disk: that count is
// a wrong number that then jumps, on a board living on refs often "0 / 0".
func TestHomeShowsAPlaceholderUntilItHasCounted(t *testing.T) {
	root := t.TempDir()
	dir := makeBoard(t, root, "repo")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	createTicket(t, s, "on disk", "backlog", "")

	h := newHome(t, []string{dir}, 100, 30)
	out := stripANSI(h.render())
	if !strings.Contains(out, "… open") || strings.Contains(out, "open /") {
		t.Fatalf("before counting, the row shows a count instead of the placeholder:\n%s", out)
	}
	if !strings.Contains(out, "repo") {
		t.Fatalf("the board's name is missing before counting:\n%s", out)
	}

	h.Update(h.Init()())
	out = stripANSI(h.render())
	if !strings.Contains(out, "1 open / 1") || strings.Contains(out, "… open") {
		t.Errorf("after counting, the row does not show the count:\n%s", out)
	}
}
