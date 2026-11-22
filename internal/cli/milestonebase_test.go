package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
)

// The base an older build's edits are merged from is carried by the ID ref
// itself: its first commit is the migration, nothing else (ticket EBTX6X,
// round 9, option C). Every clone derives it from there, so a clone that never
// moved the milestone itself merges as exactly as the one that did.

// A teammate who updates late: ada moved the milestone, removed T3 with this
// build and renamed it; no older build edited anything. grace's first command
// on the new build must not bring T3 back or leave a clash comment.
func TestALateTeammateGetsNoComebackAndNoFalseClashComment(t *testing.T) {
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
	if out, err := runAndSend(t, ada, "milestone", "rename", "alpha", "alpha2"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	flushRefs()
	ref := idRefContent(t, ada, ms.ID)
	if strings.Contains(ref, adoptT3) {
		t.Errorf("T3, removed with this build, came back on the ID ref:\n%s", ref)
	}
	if strings.Contains(ref, "<!--") {
		t.Errorf("a clash comment was pushed though no older build edited anything:\n%s", ref)
	}
	gms, _ := milestoneLoad(grace, "alpha2")
	if b, _ := os.ReadFile(gms.Path(grace)); strings.Contains(string(b), adoptT3) || strings.Contains(string(b), "<!--") {
		t.Errorf("grace's file has a comeback or a comment:\n%s", b)
	}
}

// What option B would lose: an older build edits the old name ref after the
// milestone moved. A clone that has already synced, and never recorded
// anything of its own, still adopts it, and does not bring back what a removal
// on the ID side took out.
func TestAnOldEditAfterTheMigrationReachesACloneThatAlreadySynced(t *testing.T) {
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
	if out, err := runCLI(t, grace, "fetch"); err != nil { // grace syncs first
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	flushRefs()
	oldGrace(withT + "- " + adoptT2 + "\n") // an older build adds T2 afterwards
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("second fetch: %v\n%s", err, out)
	}
	flushRefs()
	gms, _ := milestoneLoad(grace, "alpha")
	b, _ := os.ReadFile(gms.Path(grace))
	if !strings.Contains(string(b), adoptT2) {
		t.Errorf("the older build's edit after the migration was lost:\n%s", b)
	}
	if strings.Contains(string(b), adoptT3) {
		t.Errorf("T3, removed with this build, came back:\n%s", b)
	}
	if ref := idRefContent(t, ada, ms.ID); !strings.Contains(ref, adoptT2) || strings.Contains(ref, adoptT3) {
		t.Errorf("the ID ref after the adoption:\n%s", ref)
	}
}

// Two clones migrate the same old file offline, each with an edit of its own.
// Both push a root commit; the second is refused and merged. They end with the
// same file, nothing is lost and nothing comes back, and the ref has a
// migration commit a third clone can derive its base from.
func TestTwoOfflineMigrationsConverge(t *testing.T) {
	ada, grace := twoBoards(t)
	withT := adoptBase + "- " + adoptT1 + "\n- " + adoptT3 + "\n"
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": withT})
	oldMilestoneFiles(t, grace, map[string]string{"alpha.md": withT})
	offline(t, ada, func() {
		if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
			t.Fatalf("ada: %v\n%s", err, out)
		}
	})
	offline(t, grace, func() {
		if out, err := runAndSend(t, grace, "milestone", "set", "alpha", "--goal", "gg"); err != nil {
			t.Fatalf("grace: %v\n%s", err, out)
		}
	})
	// Back online, each clone's next write sends its queued move together with
	// the edit that followed: one push, two commits.
	for _, w := range []struct{ dir, goal string }{{ada, "ga2"}, {grace, "gg2"}} {
		if out, err := runAndSend(t, w.dir, "milestone", "set", "alpha", "--goal", w.goal); err != nil {
			t.Fatalf("set: %v\n%s", err, out)
		}
	}
	for _, d := range []string{ada, grace, ada, grace} {
		if out, err := runAndSend(t, d, "fetch"); err != nil {
			t.Fatalf("fetch: %v\n%s", err, out)
		}
	}
	_, a := alphaFile(t, ada)
	ms, g := alphaFile(t, grace)
	ref := idRefContent(t, ada, ms.ID)
	// grace's merged edit waits for her next write; what it adds is a comment.
	for name, text := range map[string]string{"ada": a, "grace": g, "ref": ref} {
		for _, id := range []string{adoptT1, adoptT3} {
			if !strings.Contains(text, id) {
				t.Errorf("%s: %s was lost:\n%s", name, id, text)
			}
		}
		if !strings.Contains(text, "goal: ga2") {
			t.Errorf("%s: the first pusher's goal is not kept:\n%s", name, text)
		}
	}
	repo := &gitref.Repo{Dir: ada, Remote: "origin"}
	content, _, err := repo.ReadMigration(ms.ID)
	if err != nil {
		t.Fatalf("the ID ref has no migration commit: %v", err)
	}
	if !strings.Contains(string(content), adoptT3) || strings.Contains(string(content), "goal:") {
		t.Errorf("the migration commit holds more than the migrated content:\n%s", content)
	}
}
