package snapshot

import (
	"path/filepath"
	"time"

	"github.com/BeMuCa/jaira/core/bgrun"
)

// stampName is the per-working-tree record of when this board was last
// snapshotted. Per tree rather than per machine, unlike the release check: how
// old a board's backup is says nothing about any other board's.
const stampName = "snapshot.json"

// DisableEnv switches snapshots off, and is also the recursion guard handed to
// the spawned child, so there is one mechanism here rather than two.
const DisableEnv = "JAIRA_NO_SNAPSHOT"

// StampPath is where the record lives, under the state directory the caller
// passes — never inside the repository, because when the backup last ran is
// this machine's business and committing it would be one more thing for two
// clones to disagree about.
func StampPath(stateDir string) string { return filepath.Join(stateDir, stampName) }

// Task is the background job: take a snapshot when the last one is older than
// every. Everything careful about running it — stamping before the run, the
// recursion guard, not waiting, not spawning from a test binary — lives in
// core/bgrun, which the ref fetch uses too.
func Task(stateDir, boardRoot string, every time.Duration) bgrun.Task {
	if every <= 0 {
		every = DefaultEvery
	}
	return bgrun.Task{
		Stamp:      StampPath(stateDir),
		DisableEnv: DisableEnv,
		Args:       []string{"snapshot", "--json"},
		Dir:        boardRoot,
		Every:      every,
	}
}

// Stamp is what gets recorded, kept as this package's own name for it.
type Stamp = bgrun.Mark

// ReadStamp returns the record, or a zero Stamp.
func ReadStamp(stateDir string) Stamp { return Task(stateDir, "", DefaultEvery).Read() }

// WriteStamp records that a run happened.
func WriteStamp(stateDir string, s Stamp) error { return Task(stateDir, "", DefaultEvery).Write(s) }

// Due reports whether the last snapshot is older than every.
func Due(stateDir string, every time.Duration) bool { return Task(stateDir, "", every).Due() }

// Disabled reports whether snapshots are switched off for this process.
func Disabled() bool { return bgrun.Task{DisableEnv: DisableEnv}.Disabled() }

// SpawnRun starts a detached child that takes the snapshot, if one is due.
func SpawnRun(stateDir, boardRoot string, every time.Duration) {
	Task(stateDir, boardRoot, every).Spawn()
}
