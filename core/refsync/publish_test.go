package refsync_test

import (
	"slices"
	"testing"
)

// A milestone made before the board had a remote is already named after its
// ID, so the move to ID names passes it by — and nothing else ever queued it.
// The first fetch that can reach a remote sends it, once, and the other clone
// has it on its next fetch, without a commit.
func TestAMilestoneOnThisDiskOnlyReachesTheOtherCloneOnTheNextFetch(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, ada, "offline-K3F9QZ.md", body)

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

	got, err := ada.syncer.FetchMilestones(ada.store)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(got.Published) != 0 || ada.syncer.HasMilestoneRef("K3F9QZ") {
		t.Errorf("published %v", got.Published)
	}
}
