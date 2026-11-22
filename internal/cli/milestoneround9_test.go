package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
)

// Round 9 (ticket EBTX6X): what critique 8 found in the ways an older build's
// edit, or a refused write, is merged.

// A teammate who updated after another clone had already moved the milestone
// and removed a ticket with this build: her first write is made offline, so the
// queued create is refused at the first send. The refused write is merged from
// the base recorded when her copy moved, not from nothing, so the ticket the
// teammate removed does not come back.
func TestARefusedMovingWriteDoesNotBringBackATicketARemovalTookOut(t *testing.T) {
	ada, grace, oldGrace := oldAlpha(t)
	withT := adoptBase + "- " + adoptT1 + "\n- " + adoptT3 + "\n"
	oldGrace(withT)
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": withT})
	oldMilestoneFiles(t, grace, map[string]string{"alpha.md": withT})
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	ms, text := alphaFile(t, ada)
	if err := os.WriteFile(ms.Path(ada), []byte(strings.Replace(text, "- "+adoptT3+"\n", "", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga2"); err != nil {
		t.Fatalf("set2: %v\n%s", err, out)
	}
	offline(t, grace, func() {
		if out, err := runAndSend(t, grace, "milestone", "set", "alpha", "--goal", "gg"); err != nil {
			t.Fatalf("grace offline set: %v\n%s", err, out)
		}
	})
	for _, goal := range []string{"gg2", "gg3"} {
		if out, err := runAndSend(t, grace, "milestone", "set", "alpha", "--goal", goal); err != nil {
			t.Fatalf("grace set: %v\n%s", err, out)
		}
	}
	if ref := idRefContent(t, ada, ms.ID); strings.Contains(ref, adoptT3) {
		t.Errorf("the ticket the teammate removed came back on the ID ref:\n%s", ref)
	}
	if _, g := alphaFile(t, grace); strings.Contains(g, adoptT3) {
		t.Errorf("the ticket the teammate removed came back in grace's file:\n%s", g)
	}
}

// An older build's alpha.md written after this build renamed the milestone is
// merged into the renamed file, found by its ID, and not left as a second file
// the next fetch removes unmerged.
func TestAnOldFileWrittenAfterARenameIsMergedIntoTheRenamedMilestone(t *testing.T) {
	ada, _ := twoBoards(t)
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": adoptBase})
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "g1"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "milestone", "rename", "alpha", "alpha2"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": adoptBase + "- " + adoptT1 + "\n"})
	offline(t, ada, func() {
		out, err := runAndSend(t, ada, "milestone", "set", "alpha2", "--goal", "g2")
		t.Logf("set alpha2 offline: %v\n%s", err, out)
	})
	for _, goal := range []string{"g3", "g4"} {
		if out, err := runAndSend(t, ada, "milestone", "set", "alpha2", "--goal", goal); err != nil {
			t.Fatalf("set: %v\n%s", err, out)
		}
	}
	ms, err := milestoneLoad(ada, "alpha2")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(ms.Path(ada))
	if !strings.Contains(string(b), adoptT1) {
		t.Errorf("the older build's ticket is not in the renamed file:\n%s", b)
	}
	if ref := idRefContent(t, ada, ms.ID); !strings.Contains(ref, adoptT1) {
		t.Errorf("the older build's ticket did not reach the ID ref:\n%s", ref)
	}
}

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
	names := milestoneDirNames(t, grace)
	if !strings.Contains(names, "beta-") || strings.Contains(names, "alpha-") {
		t.Errorf("the refused rename was undone: %s", names)
	}
	ms, err := milestoneLoad(grace, "beta")
	if err != nil {
		t.Fatalf("beta not found: %v", err)
	}
	b, _ := os.ReadFile(ms.Path(grace))
	if !strings.Contains(string(b), "goal: gc") {
		t.Errorf("the later goal was not written:\n%s", b)
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

func milestoneLoad(dir, key string) (*milestone.Milestone, error) { return milestone.Load(dir, key) }
