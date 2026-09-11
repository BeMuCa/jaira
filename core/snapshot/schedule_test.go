package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/snapshot"
)

// The cadence decides how often a network round trip happens in the
// background, so the two edges — never run, and just run — are what matter.
func TestDueOnlyWhenTheLastSnapshotIsOldEnough(t *testing.T) {
	dir := t.TempDir()

	if !snapshot.Due(dir, time.Hour) {
		t.Error("a board that never ran a snapshot is not due")
	}
	if err := snapshot.WriteStamp(dir, snapshot.Stamp{RanAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if snapshot.Due(dir, time.Hour) {
		t.Error("a snapshot taken just now is due again")
	}
	if err := snapshot.WriteStamp(dir, snapshot.Stamp{RanAt: time.Now().Add(-2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if !snapshot.Due(dir, time.Hour) {
		t.Error("a snapshot older than the interval is not due")
	}
	// Zero means the default, three days, so a caller that forgot to configure
	// one does not turn this into a per-command network call.
	if snapshot.Due(dir, 0) {
		t.Error("a two-hour-old snapshot is due under the default interval")
	}
}

// The switch that turns snapshots off is the same one handed to a spawned
// child as a recursion guard, so there is one mechanism rather than two.
func TestTheDisableSwitchStopsItBeingDue(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_NO_SNAPSHOT", "1")
	if !snapshot.Disabled() {
		t.Fatal("the switch was not read")
	}
	if snapshot.Due(dir, time.Nanosecond) {
		t.Error("a disabled board reported a snapshot as due")
	}
}

// A damaged or missing record must read as "never ran" rather than fail: the
// cost of getting it wrong is one extra snapshot.
func TestADamagedStampReadsAsNeverRan(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(filepath.Join(dir, "snapshot.json"), "{not json"); err != nil {
		t.Fatal(err)
	}
	if got := snapshot.ReadStamp(dir); !got.RanAt.IsZero() {
		t.Errorf("a damaged stamp came back as %+v", got)
	}
	if !snapshot.Due(dir, time.Hour) {
		t.Error("a damaged stamp did not read as never having run")
	}
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
