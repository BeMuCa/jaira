package cli

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// shippedLane is this binary's own version of the built-in lane id.
func shippedLane(t *testing.T, id string) []byte {
	t.Helper()
	builtins, err := lane.Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range builtins {
		if b.ID == id {
			want, err := lane.Bytes(b)
			if err != nil {
				t.Fatal(err)
			}
			return want
		}
	}
	t.Fatalf("no built-in lane %s", id)
	return nil
}

// refreshLines are the lines of stderr the automatic refresh printed; the
// skills nudge, which newStaleSetup also provokes, is somebody else's line.
func refreshLines(stderr string) []string {
	var out []string
	for _, l := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(l, "jaira: brought this board up to") {
			out = append(out, l)
		}
	}
	return out
}

// upgrade is the state after 'jaira self upgrade': the board was stamped by the
// version newStaleSetup ran, and the binary is now a newer one.
func (st *staleSetup) upgrade(t *testing.T) {
	t.Helper()
	setCurrent(t, "1.0.1")
}

func TestFirstCommandAfterUpgradeReplacesAnUneditedLane(t *testing.T) {
	st := newStaleSetup(t)
	st.upgrade(t)

	_, stderr := captureStdio(t, func() {
		if out, err := runCLI(t, st.dir, "list"); err != nil {
			t.Fatalf("list: %v\n%s", err, out)
		}
	})
	if got, want := mustRead(t, st.lan), shippedLane(t, "testing"); string(got) != string(want) {
		t.Errorf("the testing lane was not replaced by this binary's version")
	}
	lines := refreshLines(stderr)
	if len(lines) != 1 || !strings.Contains(lines[0], "updated lanes testing") {
		t.Fatalf("want one line naming the updated lane, stderr:\n%s", stderr)
	}
	if !strings.Contains(lines[0], "'jaira update' shows what else changed") {
		t.Errorf("the line does not point at 'jaira update' for the release notes: %q", lines[0])
	}
	if strings.Contains(stderr, "run 'jaira update'") {
		t.Errorf("the old hint was printed besides the refresh:\n%s", stderr)
	}
	// The release notes are not swallowed: 'jaira update' still reads 1.0.0
	// as the version that last set the board up.
	s, err := ticket.At(st.dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := release.Stamped(s.StateDir()); got != "1.0.0" {
		t.Errorf("Stamped = %q after the refresh, want 1.0.0 kept for 'jaira update'", got)
	}

	// The second command refreshes nothing more and is back to the hint.
	_, stderr = captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if lines := refreshLines(stderr); len(lines) != 0 {
		t.Errorf("second command refreshed again: %q", lines)
	}
	if !strings.Contains(stderr, "run 'jaira update'") {
		t.Errorf("second command lost the hint to 'jaira update':\n%s", stderr)
	}
	// After 'jaira update' both stamps are current and nothing is said.
	if out, err := runUpdate(t, st.dir); err != nil {
		t.Fatalf("update: %v\n%s", err, out)
	}
	_, stderr = captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if strings.Contains(stderr, "set up by an older version") || len(refreshLines(stderr)) != 0 {
		t.Errorf("still nudged after 'jaira update':\n%s", stderr)
	}
}

func TestAutomaticRefreshLeavesAnEditedLaneByteForByte(t *testing.T) {
	st := newStaleSetup(t)
	edited := append(append([]byte{}, st.laneOld...), []byte("\nan extra line\n")...)
	mustWrite(t, st.lan, edited)
	st.upgrade(t)

	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.lan)) != string(edited) {
		t.Error("the edited lane file was touched")
	}
	lines := refreshLines(stderr)
	if len(lines) != 1 || !strings.Contains(lines[0], "testing") || !strings.Contains(lines[0], "'jaira update --lanes'") {
		t.Fatalf("want one line naming the edited lane and 'jaira update --lanes', stderr:\n%s", stderr)
	}
	if strings.Contains(lines[0], "updated lanes") {
		t.Errorf("the line claims an update: %q", lines[0])
	}
}

func TestAutomaticRefreshRenewsTheJairaSection(t *testing.T) {
	st := newStaleSetup(t)
	st.upgrade(t)
	if st.blockCurrent(t) {
		t.Fatal("setup: the block is already current")
	}

	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if !st.blockCurrent(t) {
		t.Error("the jaira section in CLAUDE.md was not renewed")
	}
	if strings.Contains(string(mustRead(t, st.claudeMD)), "AN OLD LINE") {
		t.Error("the stale line is still in CLAUDE.md")
	}
	if lines := refreshLines(stderr); len(lines) != 1 || !strings.Contains(lines[0], "CLAUDE.md") {
		t.Fatalf("want the line to name CLAUDE.md, stderr:\n%s", stderr)
	}
}

func TestAutomaticRefreshCreatesNoAgentFile(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, s := updateTestStore(t)
	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		os.Remove(filepath.Join(dir, name))
	}
	set, err := lane.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	written, err := board.RenewInAgentFiles(dir, set.Facts())
	if err != nil || len(written) != 0 {
		t.Fatalf("RenewInAgentFiles = %v, %v; want nothing written", written, err)
	}
	captureStdio(t, func() { nudgeIfStale(s) })
	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			t.Errorf("%s was created by an automatic refresh", name)
		}
	}
}

// Two jaira processes started at once after an upgrade: the lock lets one of
// them do the work, the other finds the stamp; and a reader loading the lane
// meanwhile sees the old file or the new one, never part of either.
func TestConcurrentRefreshesWriteNoHalfLane(t *testing.T) {
	st := newStaleSetup(t)
	st.upgrade(t)
	old, want := string(st.laneOld), string(shippedLane(t, "testing"))

	stop := make(chan struct{})
	var torn []int
	var readers sync.WaitGroup
	readers.Add(1)
	go func() {
		defer readers.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			b, err := os.ReadFile(st.lan)
			if err != nil {
				continue // between rename and nothing: a missing file is not torn
			}
			if got := string(b); got != old && got != want {
				torn = append(torn, len(got))
			}
		}
	}()

	const n = 8
	lines := make([]string, n)
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s, err := ticket.At(st.dir)
			if err != nil {
				errs[i] = err
				return
			}
			lines[i], errs[i] = refreshStaleBoard(s)
		}(i)
	}
	wg.Wait()
	close(stop)
	readers.Wait()

	spoke := 0
	for i := range lines {
		if errs[i] != nil {
			t.Fatalf("refresh %d: %v", i, errs[i])
		}
		if lines[i] != "" {
			spoke++
		}
	}
	if spoke != 1 {
		t.Errorf("%d of %d refreshes reported work, want exactly one", spoke, n)
	}
	if len(torn) > 0 {
		t.Errorf("a reader saw %d torn lane file(s), lengths %v", len(torn), torn)
	}
	if string(mustRead(t, st.lan)) != want {
		t.Error("the lane file is not this binary's version")
	}
}

// A refresh that cannot finish leaves its stamp unset and the old hint printed,
// so the next command tries again and no command fails over it.
func TestFailedRefreshFallsBackToTheHint(t *testing.T) {
	st := newStaleSetup(t)
	st.upgrade(t)
	s, err := ticket.At(st.dir)
	if err != nil {
		t.Fatal(err)
	}
	// A file where the lock directory belongs: the lock cannot be taken.
	locks := filepath.Join(s.StateDir(), "locks")
	if err := os.RemoveAll(locks); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, locks, []byte("not a directory"))

	stdout, stderr := captureStdio(t, func() { nudgeIfStale(s) })
	if stdout != "" {
		t.Fatalf("wrote to stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "run 'jaira update'") {
		t.Errorf("stderr = %q, want the old hint", stderr)
	}
	if got := release.LanesRefreshed(s.StateDir()); got == "1.0.1" {
		t.Error("a failed refresh stamped the board")
	}
	if string(mustRead(t, st.lan)) != string(st.laneOld) {
		t.Error("a refresh that could not take the lock wrote a lane")
	}
}
