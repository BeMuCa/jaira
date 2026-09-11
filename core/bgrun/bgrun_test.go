package bgrun_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/bgrun"
)

func task(t *testing.T, every time.Duration) bgrun.Task {
	t.Helper()
	return bgrun.Task{
		Stamp:      filepath.Join(t.TempDir(), "task.json"),
		DisableEnv: "JAIRA_TEST_TASK_OFF",
		Args:       []string{"fetch"},
		Every:      every,
	}
}

// The whole point of the stamp is that a task runs at most every so often, so
// the two edges are what matter: never run, and just run.
func TestDueOnlyWhenTheLastRunIsOldEnough(t *testing.T) {
	task := task(t, time.Hour)

	if !task.Due() {
		t.Error("a task that never ran is not due")
	}
	if err := task.Write(bgrun.Mark{RanAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if task.Due() {
		t.Error("a task that just ran is due again")
	}
	if err := task.Write(bgrun.Mark{RanAt: time.Now().Add(-2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if !task.Due() {
		t.Error("a task older than its interval is not due")
	}
}

// A task with no interval never runs on its own: a caller that forgot to set
// one must not turn this into a spawn after every command.
func TestATaskWithoutAnIntervalIsNeverDue(t *testing.T) {
	if task(t, 0).Due() {
		t.Error("a task with no interval reported itself due")
	}
}

// The switch that turns a task off is the same one handed to the child as a
// recursion guard, so a child can never spawn a grandchild.
func TestTheDisableSwitchStopsIt(t *testing.T) {
	task := task(t, time.Nanosecond)
	t.Setenv("JAIRA_TEST_TASK_OFF", "1")
	if !task.Disabled() {
		t.Fatal("the switch was not read")
	}
	if task.Due() {
		t.Error("a disabled task reported itself due")
	}
	// And Spawn is a no-op, which is what makes a child safe to run.
	task.Spawn()
	if _, err := os.Stat(task.Stamp); err == nil {
		t.Error("a disabled task stamped itself")
	}
}

// A damaged or missing record reads as "never ran": the cost of being wrong is
// one extra run, and refusing to parse it would cost the task entirely.
func TestADamagedMarkReadsAsNeverRan(t *testing.T) {
	task := task(t, time.Hour)
	if err := os.WriteFile(task.Stamp, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := task.Read(); !got.RanAt.IsZero() {
		t.Errorf("a damaged mark came back as %+v", got)
	}
	if !task.Due() {
		t.Error("a damaged mark did not read as never having run")
	}
}

// Spawning stamps BEFORE running, so a child that dies or hangs cannot make the
// next command start another one. From a test binary nothing is started at all,
// which is why this can be asserted here.
func TestSpawnStampsBeforeRunning(t *testing.T) {
	// An hour, not a nanosecond: with an interval that short the task is due
	// again the instant it is stamped, and the assertion below would be about
	// the clock rather than about the stamp.
	task := task(t, time.Hour)
	task.Spawn()
	if got := task.Read(); got.RanAt.IsZero() {
		t.Fatal("spawning did not stamp the task")
	}
	if task.Due() {
		t.Error("the task is still due right after spawning")
	}
}
