package refsync_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

const body = "---\ncolor: 45\n---\n"

const (
	tk1 = "01M2GGKMFB1XXK7V8AFW0YGWX1"
	tk2 = "01M2GGKMFB1XXK7V8AFW0YGWX2"
)

// putRef writes a milestone ref by hand through gitref, the way another
// clone's build would put it on the remote.
func putRef(t *testing.T, s side, key, file, content string) {
	t.Helper()
	if _, err := s.syncer.Repo.WriteMilestone(key, file, []byte(content), ""); err != nil {
		t.Fatalf("write ref %s: %v", key, err)
	}
}

// putDisk leaves a milestone file where an older build kept it.
func putDisk(t *testing.T, s side, file, content string) {
	t.Helper()
	if err := os.MkdirAll(milestone.Dir(s.store.Root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(milestone.Dir(s.store.Root), file), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func diskFiles(s side) string {
	entries, _ := os.ReadDir(milestone.Dir(s.store.Root))
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, ",")
}

// fetch is what 'jaira fetch' does: the fetch, then the milestones.
func fetch(t *testing.T, s side) refsync.MilestonesFetched {
	t.Helper()
	if _, err := s.syncer.Incoming(); err != nil {
		t.Fatal(err)
	}
	got, err := s.syncer.FetchMilestones(s.store)
	if err != nil {
		t.Fatalf("fetch milestones: %v", err)
	}
	return got
}

func flush(t *testing.T, s side) []refsync.Report {
	t.Helper()
	reports, err := s.syncer.Flush()
	if err != nil {
		t.Fatal(err)
	}
	return reports
}

func createMilestone(t *testing.T, s side, name string, ids ...string) *milestone.Milestone {
	t.Helper()
	ms, err := milestone.Create(s.store, s.syncer, func([]*milestone.Milestone) (*milestone.Milestone, error) {
		m := milestone.New(name, 45, time.Now())
		for _, id := range ids {
			m.Add(id)
		}
		return m, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return ms
}

func load(t *testing.T, s side, key string) *milestone.Milestone {
	t.Helper()
	ms, err := milestone.Load(s.syncer, key)
	if err != nil {
		t.Fatalf("load %s: %v", key, err)
	}
	return ms
}

func onBoard(t *testing.T, s side) string {
	t.Helper()
	all, err := milestone.LoadAll(s.syncer)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, m := range all {
		names = append(names, m.Name)
	}
	return strings.Join(names, ",")
}

// A milestone is its ref and nothing else: creating one writes no file, and
// the other clone has it after a fetch, without a commit or a branch.
func TestAMilestoneLivesOnItsRefAndReachesTheOtherClone(t *testing.T) {
	ada, grace, _ := twoSides(t)
	ms := createMilestone(t, ada, "q1", tk1)
	if _, err := os.Stat(milestone.Dir(ada.store.Root)); !os.IsNotExist(err) {
		t.Errorf("a file was written: %v", err)
	}
	if got := load(t, ada, "q1"); !got.Has(tk1) {
		t.Error("the unsent write is not read back")
	}
	flush(t, ada)
	got := fetch(t, grace)
	if !slices.Equal(got.Wrote, []string{"q1"}) {
		t.Errorf("wrote %v", got.Wrote)
	}
	if m := load(t, grace, ms.ID); !m.Has(tk1) || m.Name != "q1" {
		t.Errorf("grace reads %s", m.Bytes())
	}
	if diskFiles(grace) != "" {
		t.Errorf("grace has files %s", diskFiles(grace))
	}
	if again := fetch(t, grace); len(again.Wrote) != 0 {
		t.Errorf("an unchanged ref is news again: %v", again.Wrote)
	}
}

// A rename and a filing arrive as what they are; a filed milestone leaves the
// board and stays readable.
func TestARenameAndAFilingArrive(t *testing.T) {
	ada, grace, _ := twoSides(t)
	a := createMilestone(t, ada, "alpha")
	c := createMilestone(t, ada, "gamma")
	flush(t, ada)
	fetch(t, grace)

	if _, err := milestone.Rename(ada.store, ada.syncer, "alpha", "beta"); err != nil {
		t.Fatal(err)
	}
	if _, err := milestone.EditFiled(ada.store, ada.syncer, c.ID, func(m *milestone.Milestone) error {
		m.SetStatus(milestone.StatusFiled)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	flush(t, ada)
	got := fetch(t, grace)
	if got.Renamed["alpha"] != "beta" || !slices.Equal(got.Filed, []string{"gamma"}) {
		t.Errorf("renamed %v filed %v", got.Renamed, got.Filed)
	}
	if onBoard(t, grace) != "beta" || load(t, grace, a.ID).Name != "beta" || !load(t, grace, c.ID).Filed() {
		t.Errorf("board %s", onBoard(t, grace))
	}
}

// A fetch never hides a milestone edit that has not been sent: the queued
// write is newer than the ref, and goes out on the next flush under its lease.
func TestAFetchKeepsALocalMilestoneEditThatIsStillQueued(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const id, file = "K3F9QZ", "backend-K3F9QZ.md"
	putRef(t, ada, id, file, body)
	fetch(t, grace)
	edited := "---\ncolor: 99\n---\n"
	if err := grace.syncer.RecordMilestone(id, file, []byte(edited)); err != nil {
		t.Fatal(err)
	}
	got := fetch(t, grace)
	if m := load(t, grace, id); m.Colour != 99 || len(got.Wrote) != 0 {
		t.Errorf("colour %d wrote %v, want the queued edit", m.Colour, got.Wrote)
	}
	flush(t, grace)
	if _, onRef, _, err := grace.syncer.Repo.ReadMilestone(id); err != nil || string(onRef) != edited {
		t.Errorf("ref carries %q (%v) after the flush", onRef, err)
	}
}

// A milestone file an older build left — made before the board had a remote,
// or never sent — is taken onto its ref by the next fetch, deleted, and
// reaches the other clone.
func TestAFileOnDiskIsTakenOntoItsRefAndDeleted(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putDisk(t, ada, "alpha-K3F9QZ.md", body+"- "+tk1+"\n")
	putDisk(t, ada, "q4.md", body)
	got := fetch(t, ada)
	if len(got.Moved) != 2 || diskFiles(ada) != "" {
		t.Fatalf("moved %v, left %q", got.Moved, diskFiles(ada))
	}
	flush(t, ada)
	fetch(t, grace)
	if onBoard(t, grace) != "alpha,q4" || !load(t, grace, "K3F9QZ").Has(tk1) || load(t, grace, "q4").ID != milestone.LegacyID("q4") {
		t.Errorf("grace sees %s", onBoard(t, grace))
	}
}

// Without a fetch that reached the remote, a missing ref proves nothing: it
// may be a teammate's milestone this clone has not fetched yet. The file
// waits, and nothing is queued against "no ref" — a write that would merge
// into theirs from an empty base and undo their removals.
func TestAFileWaitsForAFetchThatReachedTheRemote(t *testing.T) {
	ada, grace, remote := twoSides(t)
	putRef(t, grace, "K3F9QZ", "alpha-K3F9QZ.md", body+"- "+tk1+"\n")
	putDisk(t, ada, "alpha-K3F9QZ.md", body+"- "+tk1+"\n- "+tk2+"\n")

	got, err := ada.syncer.FetchMilestones(ada.store)
	if err != nil || len(got.Waiting) != 1 || diskFiles(ada) != "alpha-K3F9QZ.md" {
		t.Fatalf("no fetch: %+v %v, left %q", got, err, diskFiles(ada))
	}
	if left, _ := ada.syncer.Box.List(); len(left) != 0 {
		t.Errorf("queued %+v", left)
	}
	// One good fetch an hour ago says nothing about the milestones made since.
	if !ada.syncer.FetchRefs() {
		t.Fatal("the remote was not reached")
	}
	if err := os.Rename(remote, remote+".away"); err != nil {
		t.Fatal(err)
	}
	if ada.syncer.FetchRefs() {
		t.Fatal("the remote was reached")
	}
	putDisk(t, ada, "late-M7P2RT.md", body)
	if got, _ := ada.syncer.FetchMilestones(ada.store); len(got.Waiting) != 1 {
		t.Errorf("after a failed fetch: %+v", got)
	}
}

// A stale copy of a milestone — merged in from a branch, or committed by
// mistake — never overwrites its ref: it holds a state the ref has carried,
// and is deleted without a write. A ticket a teammate took out stays out.
func TestAStaleFileDoesNotOverwriteTheRef(t *testing.T) {
	ada, grace, _ := twoSides(t)
	ms := createMilestone(t, ada, "alpha", tk1, tk2)
	flush(t, ada)
	fetch(t, grace)
	stale := string(load(t, grace, ms.ID).Bytes())
	if _, err := milestone.Edit(ada.store, ada.syncer, ms.ID, func(m *milestone.Milestone) error {
		m.Remove(tk2)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	flush(t, ada)
	putDisk(t, grace, ms.FileName(), stale)
	got := fetch(t, grace)
	if len(got.Moved) != 1 || diskFiles(grace) != "" {
		t.Fatalf("moved %v left %q", got.Moved, diskFiles(grace))
	}
	if load(t, grace, ms.ID).Has(tk2) {
		t.Error("the stale copy brought the removed ticket back")
	}
	if left, _ := grace.syncer.Box.List(); len(left) != 0 {
		t.Errorf("queued %+v", left)
	}
}

// A board without a remote keeps its milestones on local refs; once a remote
// is configured, the next flush sends them and another clone has them.
func TestABoardWithoutARemoteKeepsMilestonesOnALocalRef(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	dir := filepath.Join(root, "plain")
	run(t, root, "git", "init", "--quiet", dir)
	run(t, dir, "git", "config", "user.name", "ada")
	run(t, dir, "git", "config", "user.email", "ada@example.test")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ada := side{store: s, syncer: refsync.New(s, "origin", "ada"), dir: dir}
	if ada.syncer.Usable() == nil {
		t.Fatal("the test needs a board without a remote")
	}
	ms := createMilestone(t, ada, "q1", tk1)
	if _, content, _, err := ada.syncer.Repo.ReadMilestone(ms.ID); err != nil || !strings.Contains(string(content), tk1) {
		t.Fatalf("not on the local ref: %q %v", content, err)
	}
	if _, err := milestone.Edit(ada.store, ada.syncer, "q1", func(m *milestone.Milestone) error { m.Add(tk2); return nil }); err != nil {
		t.Fatal(err)
	}
	if reports := flush(t, ada); len(reports) != 0 {
		t.Errorf("a board without a remote sent %+v", reports)
	}

	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)
	run(t, dir, "git", "remote", "add", "origin", bare)
	ada.syncer = refsync.New(s, "origin", "ada") // a later command sees the remote
	if reports := flush(t, ada); len(reports) != 1 || reports[0].Outcome != outbox.Sent {
		t.Fatalf("reports %+v", reports)
	}
	other := filepath.Join(root, "grace")
	run(t, root, "git", "clone", "--quiet", bare, other)
	gs, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := gs.Init(); err != nil {
		t.Fatal(err)
	}
	grace := side{store: gs, syncer: refsync.New(gs, "origin", "grace"), dir: other}
	fetch(t, grace)
	if m := load(t, grace, ms.ID); !m.Has(tk1) || !m.Has(tk2) {
		t.Errorf("grace reads %s", m.Bytes())
	}
}

// A board that is not in a git repository has nowhere to keep a milestone:
// it lists none, and a write says why.
func TestABoardOutsideGitHasNoMilestones(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	t.Setenv("GIT_CEILING_DIRECTORIES", root)
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	y := refsync.New(s, "origin", "ada")
	if files, err := y.MilestoneFiles(); err != nil || len(files) != 0 {
		t.Errorf("files %v, %v", files, err)
	}
	if err := y.RecordMilestone("K3F9QZ", "a-K3F9QZ.md", []byte(body)); !errors.Is(err, refsync.ErrNoMilestoneRefs) {
		t.Errorf("err %v", err)
	}
}

// A ref named after an old name — written before 0.5.0 gave milestones IDs —
// is read under the ID it moves to, and the next fetch moves it there; once
// the ID ref is there, the old one is not read any more.
func TestAnOldNameRefMovesToItsID(t *testing.T) {
	ada, grace, _ := twoSides(t)
	putRef(t, ada, "q4", "q4.md", body+"after: [q3]\n")
	putRef(t, ada, "q6", "q6.md", "---\nstatus: filed\nrenamed-to: q7\n---\n")
	id := milestone.LegacyID("q4")
	got := fetch(t, grace)
	if len(got.Moved) != 1 || onBoard(t, grace) != "q4" || load(t, grace, "q4").ID != id {
		t.Fatalf("moved %v board %s", got.Moved, onBoard(t, grace))
	}
	flush(t, grace)
	if file, _, _, err := grace.syncer.Repo.ReadMilestone(id); err != nil || file != "q4-"+id+".md" {
		t.Errorf("the ID ref carries %q (%v)", file, err)
	}
	putRef(t, ada, "q4b", "q4b.md", body) // unrelated, keeps the fetch busy
	if again := fetch(t, grace); len(again.Moved) != 1 {
		t.Errorf("q4 moved twice: %v", again.Moved)
	}
}

func oldEntry(t *testing.T, s side, name, content string) {
	t.Helper()
	dir := filepath.Join(s.syncer.Box.Dir, "milestones")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `{"id":"` + name + `","op":"write","ref":"refs/jaira/milestones/` + name + `","kind":"milestones","lease":"","content":` + quote(content) + `,"queued-at":"2026-10-01T10:00:00Z","updated-at":"2026-10-01T10:00:00Z","by":"ada"}`
	if err := os.WriteFile(filepath.Join(dir, name+".json"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
}

func quote(s string) string { return `"` + strings.NewReplacer("\n", `\n`, `"`, `\"`).Replace(s) + `"` }

// A write an older build queued under a name goes out on the ID ref its file
// moves to; the mark of a rename is dropped and sent nowhere.
func TestAnOldMilestoneEntryIsRekeyedToItsID(t *testing.T) {
	ada, _, _ := twoSides(t)
	oldEntry(t, ada, "backend", "---\ncolor: 45\nafter: [ui]\n---\n")
	oldEntry(t, ada, "gone", "---\ncolor: 45\nstatus: filed\nrenamed-to: elsewhere\n---\n")

	for _, r := range flush(t, ada) {
		if r.Outcome != "sent" {
			t.Errorf("report %+v", r)
		}
	}
	id := milestone.LegacyID("backend")
	file, content, _, err := ada.syncer.Repo.ReadMilestone(id)
	if err != nil || file != "backend-"+id+".md" || !strings.Contains(string(content), "after: ["+milestone.LegacyID("ui")+"]") {
		t.Errorf("the ID ref carries %q %q (%v)", file, content, err)
	}
	if _, err := ada.syncer.Repo.MilestoneSHA("backend"); err == nil {
		t.Error("the old name got a ref")
	}
	if _, err := ada.syncer.Repo.MilestoneSHA("gone"); err == nil {
		t.Error("the mark of a rename was sent")
	}
	if left, _ := ada.syncer.Box.List(); len(left) != 0 {
		t.Errorf("entries left: %+v", left)
	}
}

// An older build on the same clone queues a write keyed by the file's base
// name, <name>-<ID>: it goes out on that ID's ref in that file, not as a new
// ref for a ghost milestone whose ID is the hash of the whole base name.
func TestAnOldEntryKeyedByNameAndIDKeepsItsID(t *testing.T) {
	ada, _, _ := twoSides(t)
	oldEntry(t, ada, "alpha-MB9HR0", "---\ncolor: 45\n---\n")
	flush(t, ada)
	file, _, _, err := ada.syncer.Repo.ReadMilestone("MB9HR0")
	if err != nil || file != "alpha-MB9HR0.md" {
		t.Errorf("the ID ref carries %q (%v)", file, err)
	}
	if _, err := ada.syncer.Repo.MilestoneSHA(milestone.LegacyID("alpha-MB9HR0")); err == nil {
		t.Error("the base name got a ghost ref")
	}
}

// Two clones that moved the same old file write the same first entry; the one
// that pushes second finds the ref already says it and drops its entry, with no
// report of a lost race.
func TestAnIdenticalFirstMilestoneEntryIsDroppedSilently(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := milestone.LegacyID("q4")
	file := "q4-" + id + ".md"
	if err := ada.syncer.RecordMilestone(id, file, []byte(body)); err != nil {
		t.Fatal(err)
	}
	if err := grace.syncer.RecordMilestone(id, file, []byte(body)); err != nil {
		t.Fatal(err)
	}
	if reports := flush(t, ada); len(reports) != 1 || reports[0].Outcome != "sent" {
		t.Fatalf("ada: %+v", reports)
	}
	if reports := flush(t, grace); len(reports) != 0 {
		t.Errorf("grace's identical entry was reported: %+v", reports)
	}
	if left, _ := grace.syncer.Box.List(); len(left) != 0 {
		t.Errorf("entries left: %+v", left)
	}
}

// A refused milestone write is merged under the milestone lock, and a write
// another session queues while the merge waits for it is the one that is
// merged, not overwritten.
func TestARefusedMilestoneMergeWaitsForTheLockAndTakesAWriteQueuedMeanwhile(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const id, file = "K3F9QZ", "backend-K3F9QZ.md"
	putRef(t, ada, id, file, body)
	fetch(t, grace)
	mine := "---\ncolor: 45\ngoal: mine\n---\n"
	if err := grace.syncer.RecordMilestone(id, file, []byte(mine)); err != nil {
		t.Fatal(err)
	}
	lease, _ := ada.syncer.Repo.MilestoneSHA(id)
	if _, err := ada.syncer.Repo.WriteMilestone(id, file, []byte("---\ncolor: 46\n---\n"), lease); err != nil {
		t.Fatal(err)
	}

	unlock, err := grace.store.Lock(milestone.LockName)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan []refsync.Report, 1)
	go func() {
		reports, _ := grace.syncer.Flush()
		done <- reports
	}()
	select {
	case <-done:
		t.Fatal("the merge did not wait for the milestone lock")
	case <-time.After(300 * time.Millisecond):
	}
	if _, ok := grace.syncer.Box.PendingMilestone(id); !ok {
		t.Error("the refused entry was dropped before the merge replaced it")
	}
	newer := "---\ncolor: 45\ngoal: mine2\n---\n"
	if err := grace.syncer.Box.QueueMilestone(id, file, outbox.OpWrite, []byte(newer), "", "grace"); err != nil {
		t.Fatal(err)
	}
	unlock()
	reports := <-done
	if len(reports) != 1 || !reports[0].Merged {
		t.Fatalf("reports %+v", reports)
	}
	e, ok := grace.syncer.Box.PendingMilestone(id)
	if !ok || !strings.Contains(e.Content, "goal: mine2") || !strings.Contains(e.Content, "color: 46") {
		t.Errorf("queued %+v (%v), want the newer write merged into the winner", e, ok)
	}
	if diskFiles(grace) != "" {
		t.Errorf("the merge wrote a file: %s", diskFiles(grace))
	}
}

// A moving write (no lease, it carries the migration) that is refused keeps a
// rename made after the move: the migration's file is the name the write was
// based on, so the rename is one-sided and is taken, not dropped for the
// winner's name.
func TestARefusedMovingWriteKeepsTheRenameMadeAfterTheMove(t *testing.T) {
	ada, grace, _ := twoSides(t)
	const id, old, renamed = "K3F9QZ", "alpha-K3F9QZ.md", "beta-K3F9QZ.md"
	if err := grace.syncer.RecordMigration(id, old, []byte(body), ""); err != nil {
		t.Fatal(err)
	}
	edited := "---\ncolor: 45\ngoal: mine\n---\n"
	if err := grace.syncer.RecordMilestone(id, renamed, []byte(edited)); err != nil {
		t.Fatal(err)
	}
	putRef(t, ada, id, old, "---\ncolor: 46\n---\n")

	reports := flush(t, grace)
	if len(reports) != 1 || !reports[0].Merged {
		t.Fatalf("reports %+v", reports)
	}
	e, ok := grace.syncer.Box.PendingMilestone(id)
	if !ok || e.File != renamed || !strings.Contains(e.Content, "color: 46") || !strings.Contains(e.Content, "goal: mine") {
		t.Errorf("queued %+v (%v), want the rename and both edits", e, ok)
	}
	for _, c := range reports[0].Conflicts {
		if c.Field == "name" {
			t.Errorf("a one-sided rename was reported as a clash: %+v", c)
		}
	}
}
