package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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
	if !strings.Contains(string(mustRead(t, st.claudeMD)), "<!-- jaira:start -->\n<!-- jaira:version 1.0.1 -->\n") {
		t.Error("the renewed block does not open with the version that wrote it")
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
	written, _, err := board.RenewInAgentFiles(dir, set.Facts())
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

// A refresh that cannot start leaves its stamp unset and says so, so the next
// command tries again and no command fails over it.
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
	if !strings.Contains(stderr, "could not bring this board's lanes up to 1.0.1") || !strings.Contains(stderr, "'jaira update --lanes'") {
		t.Errorf("stderr = %q, want the failure named with 'jaira update --lanes'", stderr)
	}
	if got := release.LanesRefreshed(s.StateDir()); got == "1.0.1" {
		t.Error("a failed refresh stamped the board")
	}
	if string(mustRead(t, st.lan)) != string(st.laneOld) {
		t.Error("a refresh that could not take the lock wrote a lane")
	}
}

// A refresh that fails half way — lanes written, the agent file refused —
// still names the lanes it wrote: the next run finds them current and would
// never mention them again.
func TestRefreshFailingHalfWayNamesWhatItWrote(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a read-only file mode does not refuse the write on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits do not block writes")
	}
	st := newStaleSetup(t)
	st.upgrade(t)
	if err := os.Chmod(st.claudeMD, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(st.claudeMD, 0o644) })

	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.lan)) != string(shippedLane(t, "testing")) {
		t.Fatal("setup: the lane was not written before the agent file failed")
	}
	if lines := refreshLines(stderr); len(lines) != 1 || !strings.Contains(lines[0], "updated lanes testing") {
		t.Errorf("the written lane is not named, stderr:\n%s", stderr)
	}
	if !strings.Contains(stderr, "could not bring this board's lanes up") {
		t.Errorf("the failure is not said, stderr:\n%s", stderr)
	}
	s, err := ticket.At(st.dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := release.LanesRefreshed(s.StateDir()); got == "1.0.1" {
		t.Error("a refresh that failed half way stamped the board")
	}
}

// A command given --dry-run promised to write nothing; the refresh waits for
// the next command without one.
func TestDryRunCommandDoesNotRefresh(t *testing.T) {
	st := newStaleSetup(t)
	id, _ := jsonCLI(t, st.dir, "create", "a ticket")["id"].(string)
	st.upgrade(t)
	claude := mustRead(t, st.claudeMD)

	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "move", id, "--to", "todo", "--dry-run") })
	if string(mustRead(t, st.lan)) != string(st.laneOld) {
		t.Error("a --dry-run command replaced a lane")
	}
	if string(mustRead(t, st.claudeMD)) != string(claude) {
		t.Error("a --dry-run command rewrote CLAUDE.md")
	}
	if lines := refreshLines(stderr); len(lines) != 0 {
		t.Errorf("a --dry-run command reported a refresh: %q", lines)
	}

	captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.lan)) != string(shippedLane(t, "testing")) {
		t.Error("the next command without --dry-run did not refresh")
	}
}

// A named board's lanes are the repository's as much as today's are.
func TestRefreshReachesANamedBoard(t *testing.T) {
	st := newStaleSetup(t)
	lanes := filepath.Join(st.dir, ticket.DirName, "boards", "video", "lanes")
	if err := os.MkdirAll(lanes, 0o755); err != nil {
		t.Fatal(err)
	}
	video := filepath.Join(lanes, "testing.md")
	mustWrite(t, video, st.laneOld)
	mustWrite(t, filepath.Join(lanes, "order"), []byte("testing\n"))
	st.upgrade(t)

	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, video)) != string(shippedLane(t, "testing")) {
		t.Error("the named board's lane was not replaced")
	}
	lines := refreshLines(stderr)
	if len(lines) != 1 || !strings.Contains(lines[0], "video/testing") {
		t.Fatalf("want the line to name video/testing, stderr:\n%s", stderr)
	}
}

// 'jaira lanes use' writes today's board only, so an edited lane of a named
// board must not be offered it: following that hint would overwrite today's
// lane of the same id and leave the named one as it was.
func TestSkippedNamedBoardLaneIsNotOfferedLanesUse(t *testing.T) {
	st := newStaleSetup(t)
	lanes := filepath.Join(st.dir, ticket.DirName, "boards", "video", "lanes")
	if err := os.MkdirAll(lanes, 0o755); err != nil {
		t.Fatal(err)
	}
	video := filepath.Join(lanes, "testing.md")
	mustWrite(t, video, append(append([]byte{}, st.laneOld...), []byte("\nmine\n")...))
	mustWrite(t, filepath.Join(lanes, "order"), []byte("testing\n"))

	out, _ := runUpdate(t, st.dir, "--lanes", "--dry-run")
	setup := setupOutput(out)
	if !hasLine(setup, "skipped "+video+" (edited here; a lane of board video is replaced by hand)") {
		t.Errorf("the named board's lane is not named with the by-hand hint:\n%s", out)
	}
	for _, l := range strings.Split(setup, "\n") {
		if strings.HasPrefix(l, "skipped "+video) && strings.Contains(l, "lanes use") {
			t.Errorf("a named board's lane is offered 'jaira lanes use': %q", l)
		}
	}

	// Exit 3 is expected: an edited lane was left alone.
	raw, _ := runCLI(t, st.dir, "update", "--lanes", "--dry-run", "--json")
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("update --json: %v\n%s", err, raw)
	}
	found := false
	for _, r := range m["lanes"].([]any) {
		r := r.(map[string]any)
		if r["path"] == video {
			found = true
			if r["board"] != "video" {
				t.Errorf("JSON result for %s has board %v, want video", video, r["board"])
			}
		} else if _, ok := r["board"]; ok {
			t.Errorf("today's lane %v carries a board key", r["path"])
		}
	}
	if !found {
		t.Errorf("JSON lanes do not list %s: %v", video, m["lanes"])
	}
}

// forgetRefresh is a new working tree on the same repository: the stamps are
// per working tree, the files are the repository's.
func forgetRefresh(t *testing.T, dir string) {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(s.StateDir(), "lanes-refreshed")); err != nil {
		t.Fatal(err)
	}
}

// Two people on two versions in one repository: the older jaira must not
// write back what the newer one wrote, or the block flips with every commit.
func TestOlderBinaryLeavesANewerBlockByteForByte(t *testing.T) {
	st := newStaleSetup(t)
	setCurrent(t, "1.0.2")
	captureStdio(t, func() { runCLI(t, st.dir, "list") })
	newer := strings.Replace(string(mustRead(t, st.claudeMD)),
		"<!-- jaira:version 1.0.2 -->\n", "<!-- jaira:version 1.0.2 -->\nA LINE ONLY 1.0.2 WRITES\n", 1)
	if !strings.Contains(newer, "A LINE ONLY 1.0.2 WRITES") {
		t.Fatalf("setup: no version line in the block:\n%s", newer)
	}
	mustWrite(t, st.claudeMD, []byte(newer))

	setCurrent(t, "1.0.1")
	forgetRefresh(t, st.dir)
	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.claudeMD)) != newer {
		t.Error("an older jaira rewrote the block a newer one wrote")
	}
	lines := refreshLines(stderr)
	if len(lines) != 1 || !strings.Contains(lines[0], "a newer jaira wrote it: CLAUDE.md (1.0.2)") {
		t.Fatalf("want the line to name the newer block, stderr:\n%s", stderr)
	}
}

// The version line follows the text, not the release: a newer jaira whose
// block says the same leaves the file alone, so a release with no change to
// the block leaves no diff in anybody's repository.
func TestNewerBinaryWithTheSameTextLeavesTheVersionLine(t *testing.T) {
	st := newStaleSetup(t)
	st.upgrade(t)
	captureStdio(t, func() { runCLI(t, st.dir, "list") })
	before := mustRead(t, st.claudeMD)

	setCurrent(t, "1.0.2")
	forgetRefresh(t, st.dir)
	captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.claudeMD)) != string(before) {
		t.Error("a newer jaira with the same text rewrote the block for its version line")
	}
	if !st.blockCurrent(t) {
		t.Error("a block that differs only in its version line reads as stale")
	}
}

// A Windows checkout with core.autocrlf=true holds the section with CRLF line
// endings; the version line must still be read, or the guard is gone there.
func TestOlderBinaryLeavesANewerCRLFBlockByteForByte(t *testing.T) {
	st := newStaleSetup(t)
	setCurrent(t, "1.0.2")
	captureStdio(t, func() { runCLI(t, st.dir, "list") })
	lf := strings.Replace(string(mustRead(t, st.claudeMD)),
		"<!-- jaira:version 1.0.2 -->\n", "<!-- jaira:version 1.0.2 -->\nA LINE ONLY 1.0.2 WRITES\n", 1)
	crlf := strings.ReplaceAll(lf, "\n", "\r\n")
	mustWrite(t, st.claudeMD, []byte(crlf))

	setCurrent(t, "1.0.1")
	forgetRefresh(t, st.dir)
	_, stderr := captureStdio(t, func() { runCLI(t, st.dir, "list") })
	if string(mustRead(t, st.claudeMD)) != crlf {
		t.Error("an older jaira rewrote a newer CRLF block")
	}
	if lines := refreshLines(stderr); len(lines) != 1 || !strings.Contains(lines[0], "CLAUDE.md (1.0.2)") {
		t.Fatalf("want the line to name the newer CRLF block, stderr:\n%s", stderr)
	}
}
