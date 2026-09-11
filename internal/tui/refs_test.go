package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The board queues every write it makes, and used to have no moment to send
// them: Flush ran only after a CLI command. Somebody working a whole day in the
// board left the team looking at yesterday's tickets, and their own accepted
// work came back as a card because the ref still said what it said that
// morning.
func TestTheBoardSendsWhatItQueuedWithoutACommand(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	clone := filepath.Join(root, "ada")
	run(t, root, "git", "clone", "--quiet", bare, clone)
	run(t, clone, "git", "config", "user.name", "ada")
	run(t, clone, "git", "config", "user.email", "ada@example.test")

	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.Actor = "ada"
	y := refsync.New(s, "origin", "ada")
	s.Recorder = y
	s.Source = y

	id := ticket.NewID(time.Now())
	if _, err := s.Create(map[string]string{
		ticket.FieldID: id, ticket.FieldTitle: "written in the board",
		ticket.FieldStatus: "backlog", ticket.FieldCreator: "ada",
	}, nil, "# written in the board\n"); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Nothing has been sent yet: this is the state the board used to sit in
	// for as long as it stayed open.
	repo := &gitref.Repo{Dir: clone, Remote: "origin", AuthorName: "ada"}
	if _, err := repo.SHA(id); err == nil {
		t.Fatal("the write reached the remote before anything sent it")
	}

	// What the board's own background run does, with no command in sight.
	msg := fetchRefs(y)()
	if got, ok := msg.(refFetchedMsg); ok && got.err != nil {
		t.Fatalf("background run: %v", got.err)
	}

	if _, err := repo.SHA(id); err != nil {
		t.Errorf("the board's background run did not send the queued write: %v", err)
	}
	if _, ok := y.Pending(id); ok {
		t.Error("the write is still queued after the board ran")
	}
}

func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=ada", "GIT_AUTHOR_EMAIL=ada@example.test",
		"GIT_COMMITTER_NAME=ada", "GIT_COMMITTER_EMAIL=ada@example.test",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}
