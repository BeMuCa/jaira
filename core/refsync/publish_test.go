package refsync_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/BeMuCa/jaira/core/milestone"
)

// A milestone made before the board had a remote is already named after its
// ID, so the move to ID names passes it by — and nothing else ever queued it.
// The first fetch that can reach a remote sends it, once, and the other clone
// has it on its next fetch, without a commit.
func TestAMilestoneOnThisDiskOnlyReachesTheOtherCloneOnTheNextFetch(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, ada, "offline-K3F9QZ.md", body)

	if !ada.syncer.FetchRefs() {
		t.Fatal("the remote was not reached")
	}
	got, err := ada.syncer.FetchMilestones(ada.store)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if !slices.Equal(got.Published, []string{"offline"}) {
		t.Fatalf("published %v, want [offline]", got.Published)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if again, err := ada.syncer.FetchMilestones(ada.store); err != nil || len(again.Published) != 0 {
		t.Errorf("a milestone with its ref was sent again: %v %v", again.Published, err)
	}

	incoming(t, ada, grace)
	if got := diskFiles(t, grace); got != "offline-K3F9QZ.md" {
		t.Errorf("grace has %q", got)
	}
}

// One ID on two files is a person's to sort out: a ref written from either
// file would pick the milestone for them.
func TestAnIDOnTwoFilesIsNotSent(t *testing.T) {
	ada, _, _ := twoSides(t)
	putDisk(t, ada, "backend-K3F9QZ.md", body)
	putDisk(t, ada, "frontend-K3F9QZ.md", body)

	ada.syncer.FetchRefs()
	got, err := ada.syncer.FetchMilestones(ada.store)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(got.Published) != 0 || ada.syncer.HasMilestoneRef("K3F9QZ") {
		t.Errorf("published %v", got.Published)
	}
}

// Without a fetch that reached the remote, a missing ref proves nothing: it
// may be a teammate's milestone this clone has not fetched yet. A write queued
// then would be merged into theirs from an empty base once it is sent.
func TestNothingIsSentWithoutAFetchThatReachedTheRemote(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putRef(t, grace, "K3F9QZ", "alpha-K3F9QZ.md", body+"\n- 01A\n")
	putDisk(t, ada, "alpha-K3F9QZ.md", body+"\n- 01A\n- 01B\n")

	got, err := ada.syncer.FetchMilestones(ada.store)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(got.Published) != 0 {
		t.Errorf("sent %v without having fetched", got.Published)
	}
}

// One file that cannot be read keeps neither the others from being sent nor
// the fetch from finishing: it comes back as a RecordError, which a fetch
// reports as a warning.
func TestAnUnreadableFileDoesNotStopTheOthers(t *testing.T) {
	ada, _, _ := twoSides(t)
	if err := os.MkdirAll(milestone.Dir(ada.store.Root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("nowhere.md", filepath.Join(milestone.Dir(ada.store.Root), "gone-K3F9QZ.md")); err != nil {
		t.Skipf("no symlinks here: %v", err)
	}
	putDisk(t, ada, "real-M7P2RT.md", body)

	ada.syncer.FetchRefs()
	got, err := ada.syncer.FetchMilestones(ada.store)
	var re *milestone.RecordError
	if !errors.As(err, &re) {
		t.Errorf("err %v, want a RecordError", err)
	}
	if !slices.Equal(got.Published, []string{"real"}) {
		t.Errorf("published %v, want [real]", got.Published)
	}
}
