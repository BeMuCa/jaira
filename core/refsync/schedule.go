package refsync

import (
	"path/filepath"
	"time"

	"github.com/BeMuCa/jaira/core/bgrun"
)

// DefaultFetchEvery is how often the refs are brought up to date in the
// background.
//
// Ten minutes, and the case it exists for decides the number: somebody who
// only uses the CLI and only reads — jaira list, next, show — never fetches at
// all, because reading must never wait for a remote. Without this, a ticket
// assigned to them does not exist until they happen to type 'jaira fetch'.
const DefaultFetchEvery = 10 * time.Minute

// FetchDisableEnv switches the background fetch off, and is also the recursion
// guard handed to the spawned child.
const FetchDisableEnv = "JAIRA_NO_FETCH"

// FetchTask is the background job: fetch the ticket refs when the last fetch is
// older than every.
//
// Everything careful about running it in the background lives in core/bgrun,
// shared with the snapshot, so the stamp-before-run rule and the recursion
// guard exist once rather than once per task.
func FetchTask(stateDir, boardRoot string, every time.Duration) bgrun.Task {
	if every <= 0 {
		every = DefaultFetchEvery
	}
	return bgrun.Task{
		Stamp:      filepath.Join(stateDir, "fetch.json"),
		DisableEnv: FetchDisableEnv,
		Args:       []string{"fetch", "--json"},
		Dir:        boardRoot,
		Every:      every,
	}
}

// SpawnFetch starts a detached fetch if one is due, and never waits for it.
//
// boardRoot is not optional: the child inherits this process's working
// directory and knows nothing of a -C flag, so without it a fetch spawned from
// 'jaira -C somewhere list' would refresh whatever board the shell is standing
// in — or none at all.
func (y *Syncer) SpawnFetch(stateDir, boardRoot string, every time.Duration) {
	if y == nil || y.Usable() != nil {
		return
	}
	FetchTask(stateDir, boardRoot, every).Spawn()
}
