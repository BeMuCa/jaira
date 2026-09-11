package refsync_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Each side is a real clone with a real store, wired the way the CLI wires it:
// the store's recorder is the syncer, so an ordinary ticket write is what
// queues the ref write. Testing through Mutate rather than through Record
// directly is the point — the claim being made is that the existing write path
// now carries tickets to the remote, not that a new function can.
type side struct {
	store  *ticket.Store
	syncer *refsync.Syncer
	dir    string
}

func twoSides(t *testing.T) (ada, grace side, remote string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)

	mk := func(name string) side {
		dir := filepath.Join(root, name)
		run(t, root, "git", "clone", "--quiet", bare, dir)
		run(t, dir, "git", "config", "user.name", name)
		run(t, dir, "git", "config", "user.email", name+"@example.test")
		s, err := ticket.At(dir)
		if err != nil {
			t.Fatalf("store at %s: %v", dir, err)
		}
		if _, err := s.Init(); err != nil {
			t.Fatalf("init %s: %v", dir, err)
		}
		s.Actor = name
		y := refsync.New(s, "origin", name)
		s.Recorder = y
		return side{store: s, syncer: y, dir: dir}
	}
	return mk("ada"), mk("grace"), bare
}

// create files a ticket the way capture works: nobody owns it. A captured
// ticket belongs to no one until somebody pulls it out of the backlog, which is
// the rule pull is built on — assignee says "this is mine to work", not "I
// wrote it".
func create(t *testing.T, s side, title string) string {
	t.Helper()
	id := ticket.NewID(time.Now())
	_, err := s.store.Create(map[string]string{
		ticket.FieldID:        id,
		ticket.FieldTitle:     title,
		ticket.FieldStatus:    "backlog",
		ticket.FieldCreator:   s.store.Actor,
		ticket.FieldUpdatedBy: s.store.Actor,
	}, nil, "# "+title+"\n")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	return id
}

// An ordinary ticket write now also travels to the remote, and arrives at
// someone who has none of the writer's branches.
func TestAWriteOnTheStoreReachesTheOtherClone(t *testing.T) {
	ada, grace, _ := twoSides(t)

	id := create(t, ada, "session cookie dropped on 302")
	if _, ok := ada.syncer.Pending(id); !ok {
		t.Fatal("creating a ticket queued nothing for the remote")
	}
	reports, err := ada.syncer.Flush()
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(reports) != 1 || reports[0].Outcome != outbox.Sent {
		t.Fatalf("want one sent, got %+v", reports)
	}
	if _, ok := ada.syncer.Pending(id); ok {
		t.Error("a sent write is still queued")
	}

	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	got, _, err := grace.syncer.Repo.Read(id)
	if err != nil {
		t.Fatalf("grace read: %v", err)
	}
	if !strings.Contains(string(got), "session cookie dropped on 302") {
		t.Errorf("the ticket did not arrive:\n%s", got)
	}
	// And it arrived without any branch of Ada's being involved.
	if out := run(t, grace.dir, "git", "branch", "--list", "--all"); strings.Contains(out, "jaira") {
		t.Errorf("the ref showed up as a branch: %q", out)
	}
}

// A move made while the other side has already moved the same ticket is
// refused, and the report says who was quicker and where the ticket now is.
// "You lost" on its own would leave the user with nothing to act on.
func TestARejectedWriteNamesWhoWasQuicker(t *testing.T) {
	ada, grace, _ := twoSides(t)

	id := create(t, ada, "shared ticket")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}

	// Grace takes the ticket over from the ref, without any branch.
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	content, lease, err := grace.syncer.Repo.Read(id)
	if err != nil {
		t.Fatalf("grace read: %v", err)
	}
	moved := strings.ReplaceAll(string(content), "status: backlog", "status: review")
	moved = withAssignee(moved, "grace")
	moved = strings.ReplaceAll(moved, "updated-by: ada", "updated-by: grace")
	if _, err := grace.syncer.Repo.Write(id, []byte(moved), lease); err != nil {
		t.Fatalf("grace write: %v", err)
	}

	// Ada, who knows nothing of that, writes the ticket again.
	if _, err := ada.store.Mutate(id, func(tk *ticket.Ticket) error {
		return tk.Doc().SetScalar(ticket.FieldStatus, "in-progress")
	}); err != nil {
		t.Fatalf("ada mutate: %v", err)
	}
	reports, err := ada.syncer.Flush()
	if err != nil {
		t.Fatalf("ada flush: %v", err)
	}
	if len(reports) != 1 || reports[0].Outcome != outbox.Rejected {
		t.Fatalf("want one rejected, got %+v", reports)
	}
	w := reports[0].Winner
	if w == nil {
		t.Fatal("a rejection with no winner tells the user nothing")
	}
	if w.UpdatedBy != "grace" && w.Assignee != "grace" {
		t.Errorf("winner is %+v, want grace", w)
	}
	if w.Status != "review" {
		t.Errorf("winner status is %q, want review", w.Status)
	}
	if line := w.Describe(); !strings.Contains(line, "grace") || !strings.Contains(line, "review") {
		t.Errorf("the line shown to the user is %q", line)
	}

	// Ada's local file is untouched by the rejection: the board reconciles it
	// with the ref, and nothing about that belongs to the queue.
	tk, err := ada.store.Load(id)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if tk.Status != "in-progress" {
		t.Errorf("the local file was changed by a rejection: %q", tk.Status)
	}
}

// Taking a ticket off the board removes its ref, and does so through the same
// queue — otherwise the last act of a finished ticket would be the only one
// needing a network.
func TestRecordDeleteTakesTheRefDown(t *testing.T) {
	ada, grace, _ := twoSides(t)

	id := create(t, ada, "finished work")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	if _, err := grace.syncer.Repo.SHA(id); err != nil {
		t.Fatalf("grace should see the ref: %v", err)
	}

	if err := ada.syncer.RecordDelete(id); err != nil {
		t.Fatalf("record delete: %v", err)
	}
	e, ok := ada.syncer.Pending(id)
	if !ok || e.Op != outbox.OpDelete {
		t.Fatalf("delete was not queued: %+v", e)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush delete: %v", err)
	}
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace re-fetch: %v", err)
	}
	if _, err := grace.syncer.Repo.SHA(id); err == nil {
		t.Error("the ref survived on the other clone")
	}
}

// A board with no remote must keep working exactly as before: writes succeed,
// nothing is queued, and no command waits for git.
func TestABoardWithoutARemoteRecordsNothing(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	dir := filepath.Join(root, "plain")
	if err := exec.Command("mkdir", "-p", dir).Run(); err != nil {
		t.Fatal(err)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.Actor = "ada"
	y := refsync.New(s, "origin", "ada")
	s.Recorder = y

	id := ticket.NewID(time.Now())
	if _, err := s.Create(map[string]string{
		ticket.FieldID:     id,
		ticket.FieldTitle:  "no remote here",
		ticket.FieldStatus: "backlog",
	}, nil, "# no remote here\n"); err != nil {
		t.Fatalf("create without a repository failed: %v", err)
	}
	if _, ok := y.Pending(id); ok {
		t.Error("a directory that is not a repository queued a ref write")
	}
	reports, err := y.Flush()
	if err != nil {
		t.Fatalf("flush on a board with no remote: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("want nothing to report, got %+v", reports)
	}
}

func TestANilSyncerIsAWorkingRecorder(t *testing.T) {
	var y *refsync.Syncer
	if err := y.Record("01AAA", []byte("x")); err != nil {
		t.Errorf("nil syncer refused a write: %v", err)
	}
	if err := y.RecordDelete("01AAA"); err != nil {
		t.Errorf("nil syncer refused a delete: %v", err)
	}
	if reports, err := y.Flush(); err != nil || reports != nil {
		t.Errorf("nil syncer flush: %v %v", reports, err)
	}
	if _, ok := y.Pending("01AAA"); ok {
		t.Error("nil syncer reported a pending write")
	}
}

func TestWinnerDescribeFallsBackToSomeoneElse(t *testing.T) {
	if got := (refsync.Winner{}).Describe(); !strings.Contains(got, "someone else") {
		t.Errorf("empty winner describes as %q", got)
	}
	w := refsync.Winner{Assignee: "grace", Status: "review"}
	if got := w.Describe(); !strings.Contains(got, "grace") {
		t.Errorf("assignee-only winner describes as %q", got)
	}
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.test",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}

var _ ticket.WriteRecorder = (*refsync.Syncer)(nil)

var _ = gitref.Prefix

// The reconciliation the board shows: one story per ticket, merged field by
// field, not "whichever side is newer". This is the case that proves it — the
// local file is newer, and it still must not drag the ticket back out of
// review.
func TestReconcileMergesFieldByFieldRatherThanTakingTheNewerSide(t *testing.T) {
	ada, grace, _ := twoSides(t)
	lanes, err := lane.Load(ada.store.Root)
	if err != nil {
		t.Fatalf("lanes: %v", err)
	}

	id := create(t, ada, "shared ticket")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}

	// Grace moves it forward to review on the ref, and adds a tag.
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	content, lease, err := grace.syncer.Repo.Read(id)
	if err != nil {
		t.Fatalf("grace read: %v", err)
	}
	// Their write is deliberately older by the clock, and further along the
	// lane chain. That is the whole point of the case.
	theirs := strings.ReplaceAll(string(content), "status: backlog", "status: review")
	theirs = withTag(t, theirs, "concurrency")
	theirs = olderStamp(theirs)
	if _, err := grace.syncer.Repo.Write(id, []byte(theirs), lease); err != nil {
		t.Fatalf("grace write: %v", err)
	}

	// Ada, later in wall-clock time, only adds a tag of her own locally.
	if _, err := ada.store.Mutate(id, func(tk *ticket.Ticket) error {
		return tk.Doc().SetList(ticket.FieldTags, []string{"cli"})
	}); err != nil {
		t.Fatalf("ada mutate: %v", err)
	}
	if err := ada.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("ada fetch: %v", err)
	}

	got, err := ada.syncer.Reconcile(id, lanes)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	text := string(got.Content)
	if !strings.Contains(text, "status: review") {
		t.Errorf("a newer local write dragged the ticket back out of review:\n%s", text)
	}
	if !strings.Contains(text, "concurrency") || !strings.Contains(text, "cli") {
		t.Errorf("tags were not unioned:\n%s", text)
	}
	if got.RefOnly {
		t.Error("a ticket with a local file was reported as ref-only")
	}
}

// A ticket that arrived on a ref alone is shown as such: it is in no branch
// this clone has, which is the situation the feature exists for.
func TestReconcileMarksATicketThatIsInNoBranchHere(t *testing.T) {
	ada, grace, _ := twoSides(t)
	lanes, err := lane.Load(grace.store.Root)
	if err != nil {
		t.Fatalf("lanes: %v", err)
	}

	id := create(t, ada, "assigned to grace")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if err := grace.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}

	got, err := grace.syncer.Reconcile(id, lanes)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if !got.RefOnly {
		t.Error("a ticket present only on the ref was not marked ref-only")
	}
	if !strings.Contains(string(got.Content), "assigned to grace") {
		t.Errorf("the ref-only ticket came back empty:\n%s", got.Content)
	}
	if got.Unsent {
		t.Error("a ticket this clone never wrote was marked unsent")
	}
}

// An unsent write is a card marker, not an error: the ticket is written here
// and has not left the machine.
func TestReconcileMarksAnUnsentWrite(t *testing.T) {
	ada, _, _ := twoSides(t)
	lanes, err := lane.Load(ada.store.Root)
	if err != nil {
		t.Fatalf("lanes: %v", err)
	}
	id := create(t, ada, "not sent yet")

	got, err := ada.syncer.Reconcile(id, lanes)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if !got.Unsent {
		t.Error("a queued write was not marked unsent")
	}
	if got.RefOnly {
		t.Error("a local ticket was marked ref-only")
	}
}

// olderStamp rewrites the updated-at line in place, so the fixture stays a
// valid ticket rather than a file with a stray line after the body.
func olderStamp(doc string) string {
	out := make([]string, 0, 32)
	for _, line := range strings.Split(doc, "\n") {
		if strings.HasPrefix(line, "updated-at:") {
			line = "updated-at: 2026-09-01T10:00:00Z"
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// withTag puts one tag on a ticket file, whatever shape the tags field has, and
// fails the test rather than silently changing nothing — a fixture that quietly
// does not set up its own case proves nothing.
func withTag(t *testing.T, doc, tag string) string {
	t.Helper()
	if strings.Contains(doc, "tags: []") {
		return strings.Replace(doc, "tags: []", "tags:\n  - "+tag, 1)
	}
	if i := strings.Index(doc, "\ntags:\n"); i >= 0 {
		return doc[:i+len("\ntags:\n")] + "  - " + tag + "\n" + doc[i+len("\ntags:\n"):]
	}
	// No tags field at all: add one inside the frontmatter, after the id.
	marker := "\nstatus:"
	i := strings.Index(doc, marker)
	if i < 0 {
		t.Fatalf("fixture has no status line to anchor tags to:\n%s", doc)
	}
	return doc[:i] + "\ntags:\n  - " + tag + doc[i:]
}

// The mechanism the whole design rests on: two people try to take the same
// ticket, exactly one gets it, and the loser is left with no file at all.
func TestOnlyOneCloneCanPullATicket(t *testing.T) {
	ada, grace, _ := twoSides(t)

	// Ada files the ticket and it reaches the remote. Nobody has claimed it.
	id := create(t, ada, "cookie dropped on 302")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}

	// A third clone would be the honest fixture, but two suffice: grace pulls
	// it, and ada — who still holds the pre-pull ref — pulls too.
	if err := ada.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("ada fetch: %v", err)
	}
	got, err := grace.syncer.Pull(id, false)
	if err != nil {
		t.Fatalf("grace pull: %v", err)
	}
	if got.AlreadyHere || got.Winner != nil {
		t.Fatalf("grace's pull did not take the ticket: %+v", got)
	}
	if _, err := os.Stat(got.Path); err != nil {
		t.Fatalf("grace has no file: %v", err)
	}
	tk, err := grace.store.Load(id)
	if err != nil {
		t.Fatalf("grace cannot load it: %v", err)
	}
	if tk.Assignee != "grace" {
		t.Errorf("the pull did not make it grace's: assignee %q", tk.Assignee)
	}

	// Ada already has the file because she created it; the case worth testing
	// is a clone that does not. Remove hers first, then pull with the stale
	// ref she still holds.
	adaTicket, err := ada.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(adaTicket.Path); err != nil {
		t.Fatal(err)
	}
	lost, err := ada.syncer.Pull(id, false)
	if !errors.Is(err, refsync.ErrTaken) {
		t.Fatalf("ada's pull should have been refused as taken, got %v", err)
	}
	if lost.Winner == nil || (lost.Winner.Assignee != "grace" && lost.Winner.UpdatedBy != "grace") {
		t.Errorf("the loser was not told who has it: %+v", lost.Winner)
	}
	// And the loser holds nothing: this is the duplicate the design exists to
	// make impossible.
	if _, err := ada.store.Load(id); err == nil {
		t.Error("the loser of the race ended up with a ticket file anyway")
	}
}

// A repeated pull is a successful no-op, like moving a ticket to the lane it is
// already in: commands here are safe to retry.
func TestPullingATicketYouAlreadyHaveIsANoOp(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := create(t, ada, "already here")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	first, err := grace.syncer.Pull(id, false)
	if err != nil {
		t.Fatalf("first pull: %v", err)
	}
	second, err := grace.syncer.Pull(id, false)
	if err != nil {
		t.Fatalf("second pull: %v", err)
	}
	if !second.AlreadyHere {
		t.Error("a repeated pull was not reported as a no-op")
	}
	if second.Path != first.Path {
		t.Errorf("the second pull moved the file: %q then %q", first.Path, second.Path)
	}
}

// withAssignee sets the assignee on a ticket file whether or not it has one
// yet, because a captured ticket has none.
func withAssignee(doc, who string) string {
	out := make([]string, 0, 32)
	done := false
	for _, line := range strings.Split(doc, "\n") {
		if strings.HasPrefix(line, "assignee:") {
			line, done = "assignee: "+who, true
		}
		if !done && strings.HasPrefix(line, "status:") {
			out = append(out, "assignee: "+who)
			done = true
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// fileOnRef captures a ticket the way a board with a remote does it: the ticket
// goes to its ref and the local file goes away again, so it exists nowhere
// until somebody pulls it. This mirrors what the create command does in
// production (see cli.fileOnRefOnly) and is what the pull rules are about.
func fileOnRef(t *testing.T, s side, title string) string {
	t.Helper()
	id := create(t, s, title)
	if _, err := s.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	tk, err := s.store.Load(id)
	if err != nil {
		t.Fatalf("load after create: %v", err)
	}
	if err := os.Remove(tk.Path); err != nil {
		t.Fatalf("removing the local file: %v", err)
	}
	return id
}

// The other half of pull: while a ticket names somebody, nobody else can take
// it, and only handing it back opens it up again.
func TestReleaseHandsTheTicketBack(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := fileOnRef(t, ada, "hand it back")
	if _, err := grace.syncer.Pull(id, false); err != nil {
		t.Fatalf("grace pull: %v", err)
	}

	// Ada cannot release grace's ticket for him.
	if _, err := ada.syncer.Release(id, false); !errors.Is(err, refsync.ErrTaken) {
		t.Fatalf("ada released somebody else's ticket: %v", err)
	}

	got, err := grace.syncer.Release(id, false)
	if err != nil {
		t.Fatalf("grace release: %v", err)
	}
	if got.Removed == "" {
		t.Error("the file was left behind, which is the duplicate this rules out")
	}
	if _, err := grace.store.Load(id); err == nil {
		t.Error("grace still has the ticket file after releasing it")
	}

	// And now it is anybody's again.
	if err := ada.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("ada fetch: %v", err)
	}
	taken, err := ada.syncer.Pull(id, false)
	if err != nil {
		t.Fatalf("ada could not pull a released ticket: %v", err)
	}
	if taken.Winner != nil || taken.AlreadyHere {
		t.Errorf("ada's pull of a released ticket: %+v", taken)
	}
	tk, err := ada.store.Load(id)
	if err != nil {
		t.Fatalf("ada has no file: %v", err)
	}
	if tk.Assignee != "ada" {
		t.Errorf("assignee is %q after ada pulled it", tk.Assignee)
	}
}

// Assigning a ticket to somebody reserves it for them: they must still pull it
// themselves, and until they do — or hand it back — nobody else can.
func TestAnAssignedTicketIsReservedForItsAssignee(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := create(t, ada, "for grace")
	if _, err := ada.store.Mutate(id, func(tk *ticket.Ticket) error {
		return tk.Doc().SetScalar(ticket.FieldAssignee, "grace")
	}); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	// Assigning it does not make it hers to hold: the file goes to the ref and
	// away from here, exactly as capture does on a board with a remote.
	assigned, err := ada.store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(assigned.Path); err != nil {
		t.Fatal(err)
	}

	// Ada, who wrote it, cannot take it back by pulling.
	if _, err := ada.syncer.Pull(id, false); !errors.Is(err, refsync.ErrTaken) {
		t.Fatalf("the assigner pulled a ticket they had given away: %v", err)
	}
	// Grace can, and only then does it exist on his disk.
	if _, err := grace.store.Load(id); err == nil {
		t.Fatal("the ticket materialised at grace without him pulling it")
	}
	got, err := grace.syncer.Pull(id, false)
	if err != nil {
		t.Fatalf("grace could not pull the ticket assigned to him: %v", err)
	}
	if got.Winner != nil {
		t.Errorf("grace was refused his own ticket: %+v", got.Winner)
	}
}

// Taking a ticket from somebody has to leave a trace. Their board still says
// the ticket is theirs until they fetch, so the note on the ticket is the only
// place they can read why it stopped being.
func TestStealingRecordsWhoItWasTakenFrom(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := fileOnRef(t, ada, "contested")

	if _, err := grace.syncer.Pull(id, false); err != nil {
		t.Fatalf("grace pull: %v", err)
	}
	got, err := ada.syncer.Pull(id, true)
	if err != nil {
		t.Fatalf("ada steal: %v", err)
	}
	if got.TakenFrom != "grace" {
		t.Errorf("the steal did not record who it was taken from: %q", got.TakenFrom)
	}
	tk, err := ada.store.Load(id)
	if err != nil {
		t.Fatalf("ada has no file after stealing: %v", err)
	}
	if tk.Assignee != "ada" {
		t.Errorf("assignee is %q after the steal", tk.Assignee)
	}
	if !strings.Contains(tk.Body, "took this ticket over from grace") {
		t.Errorf("no note on the ticket about the take-over:\n%s", tk.Body)
	}
	// And the note is on the ref too, which is the copy grace will fetch.
	onRef, _, err := grace.syncer.Repo.Read(id)
	if err == nil && !strings.Contains(string(onRef), "took this ticket over from grace") {
		if err := grace.syncer.Repo.Fetch(); err != nil {
			t.Fatalf("grace fetch: %v", err)
		}
		onRef, _, err = grace.syncer.Repo.Read(id)
		if err != nil {
			t.Fatalf("grace read: %v", err)
		}
		if !strings.Contains(string(onRef), "took this ticket over from grace") {
			t.Errorf("the ref grace fetches carries no trace of the take-over:\n%s", onRef)
		}
	}
}

// One clone takes a ticket off the board while another still has the file. The
// vanished ref is the signal, and it has to survive the very fetch that
// notices it — otherwise the report is silent exactly when it matters.
func TestATicketTakenOffTheBoardElsewhereIsReportedHere(t *testing.T) {
	ada, grace, _ := twoSides(t)
	id := fileOnRef(t, ada, "finished elsewhere")

	if _, err := grace.syncer.Pull(id, false); err != nil {
		t.Fatalf("grace pull: %v", err)
	}
	if _, err := grace.syncer.Incoming(); err != nil {
		t.Fatalf("grace fetch: %v", err)
	}
	if d := grace.syncer.Departed(); len(d) != 0 {
		t.Fatalf("nothing has left the board yet: %+v", d)
	}

	// Ada takes it over and files it away. Filing does NOT remove the ref —
	// the ticket file is only in her branch, and taking the shared copy down
	// too would hide the ticket from everybody until that branch is merged.
	if _, err := ada.syncer.Pull(id, true); err != nil {
		t.Fatalf("ada steal: %v", err)
	}
	if _, err := ada.store.Archive(id); err != nil {
		t.Fatalf("ada archive: %v", err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}
	if _, err := ada.syncer.Repo.SHA(id); err != nil {
		t.Fatalf("filing a ticket away removed its ref: %v", err)
	}
	// The ref goes when the ticket has landed somewhere everybody can see,
	// which is the snapshot run's job. Done here directly, because what is
	// under test is the report, not who deletes.
	lease, err := ada.syncer.Repo.SHA(id)
	if err != nil {
		t.Fatal(err)
	}
	if err := ada.syncer.Repo.Delete(id, lease); err != nil {
		t.Fatalf("removing the ref: %v", err)
	}

	// Grace's fetch must both notice and keep noticing.
	if _, err := grace.syncer.Incoming(); err != nil {
		t.Fatalf("grace re-fetch: %v", err)
	}
	departed := grace.syncer.Departed()
	if len(departed) != 1 || departed[0].ID != id {
		t.Fatalf("the ticket that left the board was not reported: %+v", departed)
	}
	if departed[0].Path == "" {
		t.Error("the report does not say which file is still here")
	}
	// A second fetch must not lose it: it is unfinished business until the
	// person deals with it.
	if _, err := grace.syncer.Incoming(); err != nil {
		t.Fatalf("grace third fetch: %v", err)
	}
	if d := grace.syncer.Departed(); len(d) != 1 {
		t.Errorf("the report went silent on the next fetch: %+v", d)
	}
	// And once it is dealt with, it stops being reported.
	grace.syncer.ForgetDeparted(id)
	if d := grace.syncer.Departed(); len(d) != 0 {
		t.Errorf("still reported after being dealt with: %+v", d)
	}
}

// Keeping a ref alive until the ticket lands has a cost: a branch nobody ever
// merges keeps its ref for ever. That has to be visible, or the board carries
// a finished ticket nobody can account for.
func TestAFinishedTicketThatNeverLandedIsReported(t *testing.T) {
	ada, _, _ := twoSides(t)
	id := create(t, ada, "finished but never pushed")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if _, err := ada.store.Archive(id); err != nil {
		t.Fatalf("archive: %v", err)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush after archive: %v", err)
	}
	// Filing does not remove the ref, which is the behaviour this report exists
	// to keep honest.
	if _, err := ada.syncer.Repo.SHA(id); err != nil {
		t.Fatalf("filing removed the ref: %v", err)
	}

	landing := []string{"refs/heads/master"}
	// Within the grace period there is nothing to say: a review takes time.
	if got := ada.syncer.Stranded(landing, time.Hour); len(got) != 0 {
		t.Errorf("reported inside the grace period: %+v", got)
	}
	// Past it, it is named, with how long it has been waiting.
	got := ada.syncer.Stranded(landing, time.Nanosecond)
	if len(got) != 1 || got[0].ID != id {
		t.Fatalf("a finished ticket that never landed was not reported: %+v", got)
	}
	if got[0].Path == "" {
		t.Error("the report does not say which file it is")
	}
	// And with no landing branch resolved there is nothing to be late for.
	if none := ada.syncer.Stranded(nil, time.Nanosecond); len(none) != 0 {
		t.Errorf("reported with no landing branch configured: %+v", none)
	}
}
