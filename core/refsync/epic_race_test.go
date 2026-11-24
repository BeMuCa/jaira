package refsync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
)

// queueDuringPushHook files a prepared outbox entry into the pushing clone's
// outbox while a push to the epic's ref is in flight — another session's
// 'jaira epic add' landing during this one's flush. On the push numbered
// %LOSE% (0: none) it first advances the ref by one member and refuses, so the
// flush goes through rebaseEpic; the entry is filed on the push after that.
const queueDuringPushHook = `#!/bin/sh
unset GIT_QUARANTINE_PATH GIT_OBJECT_DIRECTORY GIT_ALTERNATE_OBJECT_DIRECTORIES
ref=%REF%
n=$(cat "$GIT_DIR/race-count" 2>/dev/null || echo 0)
n=$((n+1))
echo $n > "$GIT_DIR/race-count"
if [ "$n" = "%LOSE%" ]; then
	old=$(git rev-parse "$ref")
	{ git show "$old:feature.md"; echo "- %THEIRS%"; } > "$GIT_DIR/race-file"
	blob=$(git hash-object -w "$GIT_DIR/race-file")
	tree=$(printf '100644 blob %s\tfeature.md\n' "$blob" | git mktree)
	c=$(GIT_AUTHOR_NAME=hook GIT_AUTHOR_EMAIL=h@example.test GIT_COMMITTER_NAME=hook GIT_COMMITTER_EMAIL=h@example.test git commit-tree "$tree" -p "$old" -m race)
	git update-ref "$ref" "$c" "$old"
	exit 1
fi
if [ ! -e "$GIT_DIR/filed" ]; then
	mkdir -p "%DSTDIR%" && cp "%SRC%" "%DST%" && touch "$GIT_DIR/filed"
fi
exit 0
`

// An epic edit queued by a second session while the first one's flush pushes
// the same epic is not lost — neither when that push goes straight through,
// nor when it first loses a race and goes through as a merge. It stays in the
// outbox, and the next flush puts it on the remote beside everything else.
func TestAnEpicEditQueuedDuringThePushReachesTheRemote(t *testing.T) {
	const (
		mine   = "01M81V2D6K5ZVMY4E9JBS1SA9Z"
		second = "01M81V2D6K5ZVMY4E9JBS1SA8X"
		theirs = "01M81V2D6K5ZVMY4E9JBS1SA7W"
	)
	for _, tc := range []struct {
		name string
		lose string
	}{{"push goes through", "0"}, {"push is merged first", "1"}} {
		t.Run(tc.name, func(t *testing.T) {
			ada, _, remote := twoSides(t)
			ep := epic.New("feature", "Feature", "", time.Now())
			if err := ada.syncer.RecordEpic("feature", ep.Bytes()); err != nil {
				t.Fatal(err)
			}
			if _, err := ada.syncer.Flush(); err != nil {
				t.Fatal(err)
			}
			lease, err := ada.syncer.Repo.EpicSHA("feature")
			if err != nil {
				t.Fatal(err)
			}

			// What the second session queues: my edit and its own, on the
			// lease the queue held — as RecordEpic files it from EpicContent.
			ep.Add(mine)
			mid := epic.FromBytes("feature", ep.Bytes())
			mid.Add(second)
			scratch := &outbox.Box{Dir: filepath.Join(t.TempDir(), "outbox")}
			if err := scratch.QueueEpic("feature", mid.Bytes(), lease, "grace"); err != nil {
				t.Fatal(err)
			}
			box := outbox.At(ada.store)
			hook := strings.NewReplacer(
				"%REF%", gitref.EpicRefName("feature"),
				"%LOSE%", tc.lose,
				"%THEIRS%", theirs,
				"%SRC%", filepath.Join(scratch.Dir, string(outbox.KindEpic), "feature.json"),
				"%DSTDIR%", filepath.Join(box.Dir, string(outbox.KindEpic)),
				"%DST%", filepath.Join(box.Dir, string(outbox.KindEpic), "feature.json"),
			).Replace(queueDuringPushHook)
			hookPath := filepath.Join(remote, "hooks", "pre-receive")
			if err := os.WriteFile(hookPath, []byte(hook), 0o755); err != nil {
				t.Fatal(err)
			}

			if err := ada.syncer.RecordEpic("feature", ep.Bytes()); err != nil {
				t.Fatal(err)
			}
			if _, err := ada.syncer.Flush(); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(remote, "filed")); err != nil {
				t.Fatal("the hook never filed the second write; the test proves nothing")
			}
			if e, ok := box.PendingEpic("feature"); !ok || !epic.FromBytes("feature", []byte(e.Content)).Has(second) {
				t.Fatalf("the edit queued during the push is gone from the outbox: %+v, %v", e, ok)
			}

			if err := os.Remove(hookPath); err != nil {
				t.Fatal(err)
			}
			if _, err := ada.syncer.Flush(); err != nil {
				t.Fatal(err)
			}
			got := epic.FromBytes("feature", []byte(run(t, remote, "git", "show", gitref.EpicRefName("feature")+":feature.md")))
			want := []string{mine, second}
			if tc.lose != "0" {
				want = append(want, theirs)
			}
			for _, id := range want {
				if !got.Has(id) {
					t.Errorf("the remote's epic lost %s; it holds %v", id, got.Members())
				}
			}
			if _, ok := box.PendingEpic("feature"); ok {
				t.Error("the epic is still queued after the last flush")
			}
		})
	}
}
