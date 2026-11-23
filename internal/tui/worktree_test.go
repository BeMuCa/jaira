package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

// mainWithWorktree is a repository with a board in its main checkout and in
// one linked worktree, the worktree's board holding a live session.
func mainWithWorktree(t *testing.T) (main, wt *ticket.Store) {
	t.Helper()
	// Resolved: session state is keyed by the store root's spelling
	// (ticket.stateKey) and project.Trees hands back resolved paths, so on
	// macOS, where t.TempDir() is under the /var -> /private/var symlink, a
	// store opened through the symlink keeps its sessions under another key.
	base := project.Canonical(t.TempDir())
	t.Setenv("JAIRA_HOME", filepath.Join(base, "home"))
	mainDir := filepath.Join(base, "repo")
	wtDir := filepath.Join(base, ".worktrees", "repo-x")
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", mainDir}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	if err := os.MkdirAll(mainDir, 0o755); err != nil {
		t.Fatal(err)
	}
	run("init", "-q")
	run("commit", "-q", "--allow-empty", "-m", "init")
	run("worktree", "add", "-q", "-b", "x", wtDir)
	stores := []*ticket.Store{}
	for _, d := range []string{mainDir, wtDir} {
		s, err := ticket.At(d)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.Init(); err != nil {
			t.Fatal(err)
		}
		stores = append(stores, s)
	}
	live := session.Session{ID: "s1", Focus: "work", UpdatedAt: ticket.FormatTime(time.Now())}
	if err := session.Save(stores[1], live); err != nil {
		t.Fatal(err)
	}
	if _, err := project.Remember(wtDir); err != nil {
		t.Fatal(err)
	}
	return stores[0], stores[1]
}

// TestWorktreeFoldedKeepsItsState: once a worktree is listed as its main
// checkout's line, that line still is the open board when the worktree is
// open, and still shows the agent working in the worktree.
func TestWorktreeFoldedKeepsItsState(t *testing.T) {
	mainStore, wtStore := mainWithWorktree(t)
	mainRoot := project.Canonical(mainStore.Root)

	for _, c := range []struct {
		name  string
		store *ticket.Store
	}{{"opened in the worktree", wtStore}, {"opened in the main checkout", mainStore}} {
		t.Run(c.name, func(t *testing.T) {
			m, err := New(c.store)
			if err != nil {
				t.Fatal(err)
			}
			if len(m.projects) != 1 || m.projects[0].Root != mainRoot {
				t.Fatalf("projects = %+v, want one line for %s", m.projects, mainRoot)
			}
			tab := m.projectTabs()[0]
			if !strings.Contains(tab, "▸") {
				t.Errorf("tab %q does not mark the open board as current", tab)
			}
			if !strings.Contains(tab, "◆") {
				t.Errorf("tab %q does not show the agent working in the worktree", tab)
			}
			if cmd := m.switchToProject(1); cmd != nil {
				t.Error("switching to the line of the open board must be a no-op")
			}
		})
	}

	h, err := NewHome(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.entries) != 1 || h.entries[0].Agents != 1 {
		t.Fatalf("launcher entries = %+v, want one line with the worktree's agent", h.entries)
	}
}
