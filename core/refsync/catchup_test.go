package refsync_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// worktreeOf opens a second worktree of s's clone as a side of its own: the
// same refs, its own files and its own outbox.
func worktreeOf(t *testing.T, s side, name string) side {
	t.Helper()
	dir := filepath.Join(filepath.Dir(s.dir), name)
	run(t, s.dir, "git", "worktree", "add", "--quiet", "--detach", dir)
	st, err := ticket.At(dir)
	if err != nil {
		t.Fatalf("store at %s: %v", dir, err)
	}
	if _, err := st.Init(); err != nil {
		t.Fatalf("init %s: %v", dir, err)
	}
	st.Actor = s.store.Actor
	y := refsync.New(st, "origin", s.store.Actor)
	st.Recorder = y
	return side{store: st, syncer: y, dir: dir}
}

// Four ticket ids in the order they are added.
var (
	tkA = ticket.NewID(time.Unix(1, 0))
	tkB = ticket.NewID(time.Unix(2, 0))
	tkC = ticket.NewID(time.Unix(3, 0))
	tkD = ticket.NewID(time.Unix(4, 0))
)

func addTo(t *testing.T, s side, name, id string) {
	t.Helper()
	if _, err := milestone.Edit(s.store.Root, s.store, s.syncer, name, func(ms *milestone.Milestone) error {
		ms.Add(id)
		return nil
	}); err != nil {
		t.Fatalf("add %s to %s: %v", id, name, err)
	}
	if _, err := s.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
}

func refMembers(t *testing.T, s side, id string) []string {
	t.Helper()
	file, content, _, err := s.syncer.Repo.ReadMilestone(id)
	if err != nil {
		t.Fatalf("read ref %s: %v", id, err)
	}
	return milestone.FromBytes(file, content).Members()
}

// KWQCYZ, as it happened on 2026-10-22: a worktree keeps an old copy of a
// milestone file while the main checkout goes on adding to the milestone. The
// refs are one set for both worktrees, so the old worktree's lease matches,
// and an add there used to write its old file plus the new ticket over the
// ref. Every write now brings the file up to the ref first: nothing is lost.
func TestAStaleMilestoneFileNeverOverwritesItsRef(t *testing.T) {
	ada, _, _ := twoSides(t)
	run(t, ada.dir, "git", "commit", "--quiet", "--allow-empty", "-m", "start")
	ms, err := milestone.Create(ada.store.Root, ada.store, ada.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		return milestone.New("q4", 45, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "q4", tkA)

	// The old worktree takes its copy now, with AAAAAA only.
	old := worktreeOf(t, ada, "old")
	stale, err := os.ReadFile(ms.Path(ada.store.Root))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(milestone.Dir(old.store.Root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(milestone.Dir(old.store.Root), ms.FileName()), stale, 0o644); err != nil {
		t.Fatal(err)
	}

	// The main checkout goes on; the ref moves past the old copy.
	addTo(t, ada, "q4", tkB)
	addTo(t, ada, "q4", tkC)

	addTo(t, old, "q4", tkD)

	want := []string{tkA, tkB, tkC, tkD}
	if got := refMembers(t, ada, ms.ID); !slices.Equal(got, want) {
		t.Fatalf("ref holds %v, want %v: the old copy was written over the ref", got, want)
	}
	onDisk, err := milestone.Load(old.store.Root, "q4")
	if err != nil {
		t.Fatal(err)
	}
	if got := onDisk.Members(); !slices.Equal(got, want) {
		t.Fatalf("the old worktree's file holds %v, want %v", got, want)
	}
}

// A write queued here and not sent yet is newer than the ref: the catch-up
// before the next write leaves it, so the second edit builds on the first.
func TestACatchUpKeepsAWriteThatHasNotBeenSent(t *testing.T) {
	ada, _, _ := twoSides(t)
	ms, err := milestone.Create(ada.store.Root, ada.store, ada.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		return milestone.New("q4", 45, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{tkA, tkB} {
		if _, err := milestone.Edit(ada.store.Root, ada.store, ada.syncer, "q4", func(m *milestone.Milestone) error {
			m.Add(id)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := milestone.Load(ada.store.Root, ms.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{tkA, tkB}; !slices.Equal(got.Members(), want) {
		t.Fatalf("file holds %v, want %v: the catch-up threw away the unsent write", got.Members(), want)
	}
}

// A file edited by hand holds a state the ref never carried: it is what the
// next write sends, not a copy that fell behind, and the catch-up leaves it.
func TestACatchUpKeepsAHandEdit(t *testing.T) {
	ada, _, _ := twoSides(t)
	ms, err := milestone.Create(ada.store.Root, ada.store, ada.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		return milestone.New("q4", 45, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "q4", tkA)
	addTo(t, ada, "q4", tkB)
	path := ms.Path(ada.store.Root)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(b), "- "+tkA+"\n", "", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "q4", tkC)
	if got, want := refMembers(t, ada, ms.ID), []string{tkB, tkC}; !slices.Equal(got, want) {
		t.Fatalf("ref holds %v, want %v: the hand edit was thrown away", got, want)
	}
}

// A rename on the ref while this disk holds both names — the old copy fallen
// behind and a file under the new name somebody edited: the catch-up does not
// write the new name over that file. Which of two files is the milestone is
// the fetch's call.
func TestACatchUpLeavesATwinUnderTheNewName(t *testing.T) {
	ada, _, _ := twoSides(t)
	ms, err := milestone.Create(ada.store.Root, ada.store, ada.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		return milestone.New("q4", 45, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	stale, err := os.ReadFile(ms.Path(ada.store.Root))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := milestone.Rename(ada.store.Root, ada.store, ada.syncer, "q4", "q5"); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	dir := milestone.Dir(ada.store.Root)
	twin := filepath.Join(dir, "q5-"+ms.ID+".md")
	edited := "---\ncolor: 45\n---\n\n# q5\n\n- " + tkA + "\n"
	if err := os.WriteFile(twin, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ms.Path(ada.store.Root), stale, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ada.syncer.CatchUpMilestones(ada.store.Root); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(twin); string(b) != edited {
		t.Fatalf("the catch-up wrote over the file under the new name:\n%s", b)
	}
}

// The rename again, but the file under the new name is the ref itself: the
// old copy that fell behind goes, so the next write starts from the ref and
// not from whichever of two files sorts first (a1 before z9).
func TestACatchUpRemovesAnOldCopyBesideTheRenamedFile(t *testing.T) {
	ada, _, _ := twoSides(t)
	ms, err := milestone.Create(ada.store.Root, ada.store, ada.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		return milestone.New("a1", 45, time.Now()), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "a1", tkA)
	stale, err := os.ReadFile(ms.Path(ada.store.Root))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := milestone.Rename(ada.store.Root, ada.store, ada.syncer, "a1", "z9"); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "z9", tkB)
	if err := os.WriteFile(ms.Path(ada.store.Root), stale, 0o644); err != nil {
		t.Fatal(err)
	}
	addTo(t, ada, "z9", tkC)
	file, content, _, err := ada.syncer.Repo.ReadMilestone(ms.ID)
	if err != nil {
		t.Fatal(err)
	}
	got := milestone.FromBytes(file, content)
	if want := []string{tkA, tkB, tkC}; got.Name != "z9" || !slices.Equal(got.Members(), want) {
		t.Fatalf("ref is %s with %v, want z9 with %v", got.Name, got.Members(), want)
	}
	if _, err := os.Stat(ms.Path(ada.store.Root)); !os.IsNotExist(err) {
		t.Fatalf("the old copy %s is still on disk", ms.FileName())
	}
}
