package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/milestone"
)

// An older build keeps writing after this one has moved a milestone to its ID:
// its edits land on the old name ref and in an old <name>.md. These tests stand
// in for it with the same plumbing it used — the old file and the old name ref
// — and check that no such edit is ever dropped (ticket EBTX6X, round 8).

const (
	adoptT1 = "01M2GGKMFB1XXK7V8AFW0YGWX1"
	adoptT2 = "01M2GGKMFB1XXK7V8AFW0YGWX2"
	adoptT3 = "01M2GGKMFB1XXK7V8AFW0YGWX3"
)

const adoptBase = "---\ncolor: 45\ncreated-at: 2026-01-01T00:00:00Z\nafter: []\n---\n\n# alpha\n\n"

// oldAlpha gives both clones the file an older build wrote and puts it on its
// name ref, the way that build did. oldGrace is the older build running on
// grace's clone; it returns the new tip of the name ref.
func oldAlpha(t *testing.T) (ada, grace string, oldGrace func(content string) string) {
	t.Helper()
	ada, grace = twoBoards(t)
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": adoptBase})
	oldMilestoneFiles(t, grace, map[string]string{"alpha.md": adoptBase})
	adaRepo := &gitref.Repo{Dir: ada, Remote: "origin"}
	if _, err := adaRepo.WriteMilestone("alpha", "alpha.md", []byte(adoptBase), ""); err != nil {
		t.Fatal(err)
	}
	graceRepo := &gitref.Repo{Dir: grace, Remote: "origin"}
	if err := graceRepo.Fetch(); err != nil {
		t.Fatal(err)
	}
	return ada, grace, func(content string) string {
		t.Helper()
		lease, err := graceRepo.MilestoneSHA("alpha")
		if err != nil {
			t.Fatal(err)
		}
		sha, err := graceRepo.WriteMilestone("alpha", "alpha.md", []byte(content), lease)
		if err != nil {
			t.Fatal(err)
		}
		return sha
	}
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

func alphaFile(t *testing.T, dir string) (*milestone.Milestone, string) {
	t.Helper()
	ms, err := milestone.Load(dir, "alpha")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(ms.Path(dir))
	if err != nil {
		t.Fatal(err)
	}
	return ms, string(b)
}

func offline(t *testing.T, dir string, fn func()) {
	t.Helper()
	url := strings.TrimSpace(gitOut(t, dir, "remote", "get-url", "origin"))
	gitRun(t, dir, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "nowhere.git"))
	fn()
	gitRun(t, dir, "remote", "set-url", "origin", url)
}

// x3 "nofetch": a teammate's older build added a ticket to the old name ref
// which this clone had not fetched when its first write moved the milestone —
// it was offline. The ID ref was written without the ticket and every clone of
// this build then ignored the old ref: the ticket was gone. The next fetch has
// to bring the edit into the ID milestone.
func TestAnOldEditNotFetchedBeforeTheMovingWriteIsAdoptedByTheNextFetch(t *testing.T) {
	ada, _, oldGrace := oldAlpha(t)
	oldGrace(adoptBase + "- " + adoptT1 + "\n")

	offline(t, ada, func() {
		if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
			t.Fatalf("set: %v\n%s", err, out)
		}
	})
	if ms, _ := alphaFile(t, ada); ms.Has(adoptT1) || ms.Goal != "ga" {
		t.Fatalf("the offline write should not know the edit yet: goal %q members %v", ms.Goal, ms.Members())
	}

	if out, err := runAndSend(t, ada, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	ms, _ := alphaFile(t, ada)
	if !ms.Has(adoptT1) || ms.Goal != "ga" {
		t.Errorf("after the fetch: goal %q members %v, want the old ticket and the new goal", ms.Goal, ms.Members())
	}
	if ref := idRefContent(t, ada, ms.ID); !strings.Contains(ref, adoptT1) || !strings.Contains(ref, "goal: ga") {
		t.Errorf("the ID ref lacks the adopted edit:\n%s", ref)
	}
}

// x3 "fetch": with the remote reachable the write that moves the milestone
// brings the old edit in first, so the ID ref is born with it.
func TestTheMovingWriteFetchesFirstWhenTheRemoteIsReachable(t *testing.T) {
	ada, _, oldGrace := oldAlpha(t)
	oldGrace(adoptBase + "- " + adoptT1 + "\n")

	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	ms, _ := alphaFile(t, ada)
	if !ms.Has(adoptT1) || ms.Goal != "ga" {
		t.Fatalf("goal %q members %v, want the fetched ticket and the new goal", ms.Goal, ms.Members())
	}
	if ref := idRefContent(t, ada, ms.ID); !strings.Contains(ref, adoptT1) {
		t.Errorf("the ID ref was born without the old edit:\n%s", ref)
	}
}

// x3 tail: after the move the older build keeps editing the old name ref. This
// build adopts: tickets added and removed, a field only the old side changed,
// the body; a field both sides changed keeps this side and is told, on stderr,
// in the fetch report and in --json.
func TestLaterOldBuildEditsAreAdoptedAndAClashIsToldNotDropped(t *testing.T) {
	ada, _, oldGrace := oldAlpha(t)
	withT1 := adoptBase + "- " + adoptT1 + "\n- " + adoptT3 + "\n"
	oldGrace(withT1)
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ours"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}

	// the old build: removes T3, adds T2, writes its own goal and a note
	edit := strings.Replace(withT1, "after: []\n", "after: []\ngoal: theirs\n", 1)
	edit = strings.Replace(edit, "- "+adoptT3+"\n", "- "+adoptT2+"\n", 1)
	edit = strings.Replace(edit, "# alpha\n", "# alpha\n\nA note from the older build.\n", 1)
	oldGrace(edit)

	var out string
	var err error
	_, stderr := captureStdio(t, func() { out, err = runCLI(t, ada, "fetch", "--json") })
	if err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	flushRefs()
	ms, text := alphaFile(t, ada)
	if !ms.Has(adoptT1) || !ms.Has(adoptT2) || ms.Has(adoptT3) {
		t.Errorf("members %v, want T1 and T2 and not T3", ms.Members())
	}
	if ms.Goal != "ours" {
		t.Errorf("goal %q: a field both sides changed keeps this side", ms.Goal)
	}
	if !strings.Contains(text, "A note from the older build.") {
		t.Errorf("a body change only the old side made was not taken:\n%s", text)
	}
	for where, got := range map[string]string{"stderr": stderr, "--json": out} {
		if !strings.Contains(got, "theirs") {
			t.Errorf("the clashing old goal is not in the %s report:\n%s", where, got)
		}
	}
	if !strings.Contains(out, "milestones_conflicts") {
		t.Errorf("--json has no milestones_conflicts:\n%s", out)
	}
	if !strings.Contains(text, "theirs") {
		t.Errorf("the clashing old value is not kept in the file, so a background fetch would drop it silently:\n%s", text)
	}
	ref := idRefContent(t, ada, ms.ID)
	if !strings.Contains(ref, adoptT2) || strings.Contains(ref, "- "+adoptT3) {
		t.Errorf("the ID ref does not carry the adoption:\n%s", ref)
	}

	// told once: the same old state is not merged or reported again
	again, err := runCLI(t, ada, "fetch", "--json")
	if err != nil || strings.Contains(again, "theirs") {
		t.Errorf("the clash came back on a second fetch (%v):\n%s", err, again)
	}
}

// An older after: holds names. They are taken as IDs.
func TestAnAdoptedOldAfterIsSpelledAsIDs(t *testing.T) {
	ada, _, oldGrace := oldAlpha(t)
	oldMilestoneFiles(t, ada, map[string]string{"beta.md": "---\ncolor: 46\nafter: []\n---\n\n# beta\n"})
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "ga"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	oldGrace(strings.Replace(adoptBase, "after: []", "after: [beta]", 1))
	if out, err := runAndSend(t, ada, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	ms, _ := alphaFile(t, ada)
	if want := milestone.LegacyID("beta"); strings.Join(ms.After, ",") != want {
		t.Errorf("after %v, want [%s]", ms.After, want)
	}
}

// x1: an older build on the SAME clone writes <name>.md again beside the ID
// file. The next command moves nothing away: the edit goes into the ID file.
func TestAnOldFileBesideTheIDFileIsAdoptedNotRemoved(t *testing.T) {
	ada, _ := twoBoards(t)
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": adoptBase})
	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "g1"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	oldMilestoneFiles(t, ada, map[string]string{"alpha.md": adoptBase + "- " + adoptT1 + "\n"})

	if out, err := runAndSend(t, ada, "milestone", "set", "alpha", "--goal", "g2"); err != nil {
		t.Fatalf("second set: %v\n%s", err, out)
	}
	ms, _ := alphaFile(t, ada)
	if !ms.Has(adoptT1) || ms.Goal != "g2" {
		t.Errorf("goal %q members %v, want the old file's ticket and the new goal", ms.Goal, ms.Members())
	}
	if got := milestoneDirNames(t, ada); got != "alpha-"+ms.ID+".md" {
		t.Errorf("directory %s, want the ID file only", got)
	}
	if ref := idRefContent(t, ada, ms.ID); !strings.Contains(ref, adoptT1) {
		t.Errorf("the ID ref lacks the adopted ticket:\n%s", ref)
	}
}
