package refsync

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

// A detached fetch writes the seen record after every command while the
// foreground one may be reading it. A reader must get the old record or the
// new one, never a torn file that parses as empty: an empty record forgets
// every ticket that has left the board, and Departed never reports it again.
func TestTheSeenRecordIsNeverReadHalfWritten(t *testing.T) {
	y := &Syncer{SeenPath: filepath.Join(t.TempDir(), "state", "refs-seen.json")}
	rec := seen{}
	for i := 0; i < 200; i++ {
		rec[fmt.Sprintf("01TICKET%04d", i)] = fmt.Sprintf("%040d", i)
	}
	y.saveSeen(rec)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				y.saveSeen(rec)
			}
		}
	}()
	torn := 0
	for i := 0; i < 2000; i++ {
		if got := y.loadSeen(); len(got) != len(rec) {
			torn++
		}
	}
	close(stop)
	wg.Wait()
	if torn > 0 {
		t.Errorf("%d of 2000 reads saw a half-written seen record", torn)
	}
}
