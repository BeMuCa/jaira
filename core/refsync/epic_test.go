package refsync_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/epic"
	"github.com/BeMuCa/jaira/core/gitref"
)

// raceHook advances the epic's ref on the remote by one member every time a
// push to it arrives, then refuses the push: each attempt loses to a write
// that landed between its read and its push, which is the race rebaseEpic
// answers. Each member it adds is <prefix><n>, n counting the pushes.
const raceHook = `#!/bin/sh
unset GIT_QUARANTINE_PATH GIT_OBJECT_DIRECTORY GIT_ALTERNATE_OBJECT_DIRECTORIES
ref=%REF%
n=$(cat "$GIT_DIR/race-count" 2>/dev/null || echo 0)
n=$((n+1))
echo $n > "$GIT_DIR/race-count"
old=$(git rev-parse "$ref")
{ git show "$old:feature.md"; echo "- %PREFIX%$n"; } > "$GIT_DIR/race-file"
blob=$(git hash-object -w "$GIT_DIR/race-file")
tree=$(printf '100644 blob %s\tfeature.md\n' "$blob" | git mktree)
c=$(GIT_AUTHOR_NAME=hook GIT_AUTHOR_EMAIL=h@example.test GIT_COMMITTER_NAME=hook GIT_COMMITTER_EMAIL=h@example.test git commit-tree "$tree" -p "$old" -m race)
git update-ref "$ref" "$c" "$old"
exit 1
`

// An epic edit that loses every round of a flush goes back in the queue. The
// text it carries was merged onto the ref as it stood before the last lost
// push, so it must be leased on that state: leased on the ref after the
// fetch, the next flush pushes it straight over the write that won the last
// round, and that write's member is gone from the remote.
func TestAnEpicThatKeepsLosingIsRequeuedOnTheStateItWasMergedOnto(t *testing.T) {
	ada, _, remote := twoSides(t)
	ep := epic.New("feature", "Feature", "", time.Now())
	if err := ada.syncer.RecordEpic("feature", ep.Bytes()); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}

	const prefix = "01M81V2D6K5ZVMY4E9JBS1SA8"
	hook := strings.NewReplacer("%REF%", gitref.EpicRefName("feature"), "%PREFIX%", prefix).Replace(raceHook)
	hookPath := filepath.Join(remote, "hooks", "pre-receive")
	if err := os.WriteFile(hookPath, []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}

	mine := "01M81V2D6K5ZVMY4E9JBS1SA9Z"
	ep.Add(mine)
	if err := ada.syncer.RecordEpic("feature", ep.Bytes()); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(hookPath); err != nil {
		t.Fatal(err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(remote, "race-count"))
	if err != nil {
		t.Fatalf("the hook never ran: %v", err)
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(raw)))
	if n < 2 {
		t.Fatalf("the hook ran %d time(s); the test needs a write that lost more than once", n)
	}
	got := epic.FromBytes("feature", []byte(run(t, remote, "git", "show", gitref.EpicRefName("feature")+":feature.md")))
	want := []string{mine}
	for i := 1; i <= n; i++ {
		want = append(want, prefix+strconv.Itoa(i))
	}
	for _, id := range want {
		if !got.Has(id) {
			t.Errorf("the remote's epic lost %s; it holds %v", id, got.Members())
		}
	}
}
