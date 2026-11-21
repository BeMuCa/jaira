package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// A rename whose write the remote refused, because a teammate edited the same
// milestone first, is merged like a field: the file name stays the new one and
// the CLI's "merged" is true.
func TestARefusedRenameIsMergedAndKept(t *testing.T) {
	ada, grace := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "alpha"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	offline(t, grace, func() {
		if out, err := runAndSend(t, grace, "milestone", "rename", "alpha", "beta"); err != nil {
			t.Fatalf("grace rename offline: %v\n%s", err, out)
		}
	})
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	out, err := runAndSend(t, grace, "milestone", "set", "beta", "--goal", "gb")
	t.Logf("grace set: %v\n%s", err, out)
	if out, err := runAndSend(t, grace, "milestone", "set", "beta", "--goal", "gc"); err != nil {
		t.Fatalf("grace set 2: %v\n%s", err, out)
	}
	ms, err := milestoneLoad(t, grace, "beta")
	if err != nil {
		t.Fatalf("beta not found, the refused rename was undone: %v", err)
	}
	if !strings.Contains(string(ms.Bytes()), "goal: gc") {
		t.Errorf("the later goal was not written:\n%s", ms.Bytes())
	}
	repo := &gitref.Repo{Dir: ada, Remote: "origin"}
	if err := repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	file, _, _, err := repo.ReadMilestone(ms.ID)
	if err != nil || !strings.HasPrefix(file, "beta-") {
		t.Errorf("the ID ref's file is %q (%v), want the renamed beta-<ID>.md", file, err)
	}
}

// refsOf is the milestones of the board at dir, as a command there reads
// them: its refs, and what its outbox still holds.
func refsOf(t *testing.T, dir string) *refsync.Syncer {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	return refsync.New(s, "origin", "test")
}

func milestoneLoad(t *testing.T, dir, key string) (*milestone.Milestone, error) {
	t.Helper()
	return milestone.Load(refsOf(t, dir), key)
}

// idRefContent is what the remote's ID ref carries for the milestone, read
// through a fresh fetch.
func idRefContent(t *testing.T, dir, id string) string {
	t.Helper()
	repo := &gitref.Repo{Dir: dir, Remote: "origin"}
	if err := repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	_, content, _, err := repo.ReadMilestone(id)
	if err != nil {
		t.Fatalf("no ID ref %s: %v", id, err)
	}
	return string(content)
}

// offline runs fn with the board's remote pointing nowhere.
func offline(t *testing.T, dir string, fn func()) {
	t.Helper()
	url := strings.TrimSpace(gitOut(t, dir, "remote", "get-url", "origin"))
	gitRun(t, dir, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "nowhere.git"))
	fn()
	gitRun(t, dir, "remote", "set-url", "origin", url)
}
