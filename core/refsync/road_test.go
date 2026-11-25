package refsync_test

import (
	"slices"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/road"
)

// Two clones list milestones in one road at once: the one that loses the race
// merges its listing onto the winner's, as an epic's does, and both listings
// reach the ref. Nothing is dropped.
func TestAConcurrentRoadListingIsMergedNotLost(t *testing.T) {
	ada, grace, _ := twoSides(t)
	now := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	const m1, m2, m3 = "ABC123", "DEF456", "GHK789"

	r := road.New("auth", "Auth", "", now)
	r.Add(m1)
	if err := ada.syncer.RecordRoad("auth", r.Bytes()); err != nil {
		t.Fatalf("record: %v", err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	// grace reads the road as it is, so her lease is the commit ada wrote.
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	gc, err := grace.syncer.RoadContent("auth")
	if err != nil {
		t.Fatalf("grace read: %v", err)
	}
	g := road.FromBytes("auth", gc)
	g.Add(m2)
	if err := grace.syncer.RecordRoad("auth", g.Bytes()); err != nil {
		t.Fatalf("grace record: %v", err)
	}

	// ada gets there first with a listing of her own.
	a := road.FromBytes("auth", r.Bytes())
	a.Add(m3)
	if err := ada.syncer.RecordRoad("auth", a.Bytes()); err != nil {
		t.Fatalf("ada record: %v", err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}

	reports, err := grace.syncer.Flush()
	if err != nil {
		t.Fatalf("grace flush: %v", err)
	}
	if len(reports) != 1 || reports[0].Outcome != outbox.Sent {
		t.Fatalf("grace's flush = %+v; want the road merged and sent", reports)
	}
	if _, ok := grace.syncer.Box.PendingRoad("auth"); ok {
		t.Error("a merged road is still queued")
	}

	if err := ada.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("ada fetch: %v", err)
	}
	final, err := ada.syncer.RoadContent("auth")
	if err != nil {
		t.Fatalf("ada read: %v", err)
	}
	got := road.FromBytes("auth", final).Members()
	slices.Sort(got)
	if !slices.Equal(got, []string{m1, m2, m3}) {
		t.Errorf("road on the ref lists %v; want all three", got)
	}
}
