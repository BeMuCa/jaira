package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
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

// The board's own fetch brings in a teammate's milestone, as 'jaira fetch'
// does, and says so on its status line: nobody should have to run a command
// beside the board to see the team's plan.
func TestTheBoardFetchBringsInMilestones(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	open := func(who string) (*ticket.Store, *refsync.Syncer, *gitref.Repo) {
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
		return s, y, &gitref.Repo{Dir: clone, Remote: "origin", AuthorName: who}
	}
	_, _, adaRepo := open("ada")
	bobStore, bobSync, _ := open("bob")

	ms := milestone.New("backend", 45, nowish())
	ms.ID = "K3F9QZ"
	if _, err := adaRepo.WriteMilestone(ms.ID, ms.FileName(), ms.Bytes(), ""); err != nil {
		t.Fatal(err)
	}

	m, err := New(bobStore)
	if err != nil {
		t.Fatal(err)
	}
	msg := fetchRefs(bobSync)()
	if got, ok := msg.(refFetchedMsg); !ok || got.err != nil {
		t.Fatalf("background run: %#v", msg)
	}
	m.Update(msg)

	if _, err := os.Stat(filepath.Join(milestone.Dir(bobStore.Root), ms.FileName())); err != nil {
		t.Errorf("the milestone did not arrive on disk: %v", err)
	}
	if line := m.flashLine(); !strings.Contains(line, "milestones updated: backend") {
		t.Errorf("status line is %q", line)
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

// A milestone edit made in the board whose send is refused, because a teammate
// wrote the ref first, must survive: the board's next fetch used to replace the
// file with the teammate's version and the edit was gone. It is merged into the
// teammate's version instead (base: the ref as the edit was leased on) and goes
// out with the next send.
func TestARejectedMilestoneFlushKeepsTheLocalEdit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	open := func(who string) (*ticket.Store, *refsync.Syncer, *gitref.Repo) {
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
		return s, y, &gitref.Repo{Dir: clone, Remote: "origin", AuthorName: who}
	}
	_, _, adaRepo := open("ada")
	bobStore, bobSync, _ := open("bob")

	ms := milestone.New("backend", 45, nowish())
	ms.ID = "K3F9QZ"
	if _, err := adaRepo.WriteMilestone(ms.ID, ms.FileName(), ms.Bytes(), ""); err != nil {
		t.Fatal(err)
	}
	m, err := New(bobStore)
	if err != nil {
		t.Fatal(err)
	}
	m.Update(fetchRefs(bobSync)())

	// bob edits the goal in the board: written to the file, queued.
	if _, err := milestone.Edit(bobStore.Root, bobStore, bobSync, "backend", func(ms *milestone.Milestone) error {
		ms.SetGoal("bob's goal")
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	// ada writes the ref first: another field.
	lease, err := adaRepo.MilestoneSHA(ms.ID)
	if err != nil {
		t.Fatal(err)
	}
	theirs := milestone.FromBytes(ms.FileName(), ms.Bytes())
	theirs.SetColour(46)
	if _, err := adaRepo.WriteMilestone(ms.ID, ms.FileName(), theirs.Bytes(), lease); err != nil {
		t.Fatal(err)
	}

	msg := fetchRefs(bobSync)()
	if got, ok := msg.(refFetchedMsg); !ok || got.err != nil {
		t.Fatalf("background run: %#v", msg)
	}
	m.Update(msg)

	got, err := milestone.Load(bobStore.Root, "backend")
	if err != nil {
		t.Fatal(err)
	}
	if got.Goal != "bob's goal" || got.Colour != 46 {
		t.Errorf("goal %q colour %d: bob's edit and ada's must both be there", got.Goal, got.Colour)
	}
	if line := m.flashLine(); !strings.Contains(line, "backend") || !strings.Contains(line, "merged") {
		t.Errorf("status line is %q, want it to say bob's edit was merged into the teammate's version", line)
	}
	// and it goes out with the next send
	m.Update(fetchRefs(bobSync)())
	adaRepo.Fetch()
	_, content, _, err := adaRepo.ReadMilestone(ms.ID)
	if err != nil || !strings.Contains(string(content), "goal: bob's goal") {
		t.Errorf("the merged edit did not reach the remote (%v):\n%s", err, content)
	}
}
