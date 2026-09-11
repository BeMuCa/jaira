package gitref_test

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
)

// The fixtures are real repositories built at test time: a bare repo standing in
// for the remote and two clones standing in for two people. Nothing here is
// faked, because the whole mechanism under test is git's own ref lock, and a
// fake would only prove that the fake behaves as assumed.
func twoClones(t *testing.T) (a, b *gitref.Repo) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	bare := filepath.Join(root, "board.git")
	git(t, root, "init", "--bare", "--quiet", bare)

	mk := func(name string) *gitref.Repo {
		dir := filepath.Join(root, name)
		git(t, root, "clone", "--quiet", bare, dir)
		git(t, dir, "config", "user.name", name)
		git(t, dir, "config", "user.email", name+"@example.test")
		return &gitref.Repo{Dir: dir, Remote: "origin", AuthorName: name, AuthorEmail: name + "@example.test"}
	}
	return mk("ada"), mk("grace")
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.test",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

const ticket = "---\nid: 01TEST\nstatus: in-progress\nassignee: ada\n---\n\n# A ticket\n"

// A ticket written by one person is readable by another who has none of their
// branches — the property the whole feature exists for.
func TestTicketArrivesWithoutASharedBranch(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.Write("01TEST", []byte(ticket), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	got, sha, err := grace.Read("01TEST")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != ticket {
		t.Errorf("ticket came back changed:\n%q", got)
	}
	if sha == "" {
		t.Error("read returned no lease sha")
	}
	if out := git(t, grace.Dir, "branch", "--list", "--all"); strings.Contains(out, "jaira") {
		t.Errorf("the ref showed up as a branch: %q", out)
	}
}

// Two people writing the same ticket: the second is told it lost, and is told it
// as a race rather than as a transport failure, because the two need opposite
// handling.
func TestSecondWriterLosesTheRace(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.Write("01TEST", []byte(ticket), ""); err != nil {
		t.Fatalf("ada write: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	_, lease, err := grace.Read("01TEST")
	if err != nil {
		t.Fatalf("grace read: %v", err)
	}

	// Ada moves the ticket again while Grace holds the older sha.
	if _, err := ada.Write("01TEST", []byte(strings.Replace(ticket, "in-progress", "review", 1)), mustSHA(t, ada, "01TEST")); err != nil {
		t.Fatalf("ada second write: %v", err)
	}
	_, err = grace.Write("01TEST", []byte(strings.Replace(ticket, "ada", "grace", 1)), lease)
	if !errors.Is(err, gitref.ErrRaceLost) {
		t.Fatalf("want ErrRaceLost, got %v", err)
	}

	// And the remote still holds Ada's write, not a half-applied mixture.
	if err := grace.Fetch(); err != nil {
		t.Fatalf("re-fetch: %v", err)
	}
	got, lease2, err := grace.Read("01TEST")
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	if !strings.Contains(string(got), "status: review") {
		t.Errorf("remote state is not the winner's:\n%s", got)
	}

	// Re-reading is all it takes to take the ticket over.
	if _, err := grace.Write("01TEST", []byte(strings.Replace(string(got), "assignee: ada", "assignee: grace", 1)), lease2); err != nil {
		t.Fatalf("takeover after re-read: %v", err)
	}
	if err := ada.Fetch(); err != nil {
		t.Fatalf("ada fetch: %v", err)
	}
	back, _, err := ada.Read("01TEST")
	if err != nil {
		t.Fatalf("ada read: %v", err)
	}
	if !strings.Contains(string(back), "assignee: grace") {
		t.Errorf("takeover did not arrive:\n%s", back)
	}
}

// The first write is a compare-and-swap too: two people creating the same
// ticket id at once must not silently overwrite one another.
func TestFirstWriteRefusesAnExistingRef(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.Write("01TEST", []byte(ticket), ""); err != nil {
		t.Fatalf("ada write: %v", err)
	}
	_, err := grace.Write("01TEST", []byte(ticket), "")
	if !errors.Is(err, gitref.ErrRaceLost) {
		t.Fatalf("want ErrRaceLost for a ref that already exists, got %v", err)
	}
}

// A logged ticket deletes its ref, and the deletion has to reach the other
// clone — otherwise the ticket is gone for its owner and immortal for everyone
// else.
func TestDeleteTakesTheTicketOffEveryBoard(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.Write("01TEST", []byte(ticket), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if err := ada.Delete("01TEST", mustSHA(t, ada, "01TEST")); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := ada.SHA("01TEST"); !errors.Is(err, gitref.ErrNoRef) {
		t.Errorf("local ref survived the delete: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("re-fetch: %v", err)
	}
	if _, err := grace.SHA("01TEST"); !errors.Is(err, gitref.ErrNoRef) {
		t.Errorf("the other clone still has the ref: %v", err)
	}
}

// ls-remote lists every ticket in one roundtrip, which is what a board needs
// before it decides what to fetch.
func TestListRemoteSeesEveryTicket(t *testing.T) {
	ada, _ := twoClones(t)

	for _, id := range []string{"01AAA", "01BBB"} {
		if _, err := ada.Write(id, []byte(ticket), ""); err != nil {
			t.Fatalf("write %s: %v", id, err)
		}
	}
	ids, err := ada.ListRemote()
	if err != nil {
		t.Fatalf("ls-remote: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("want 2 ids, got %v", ids)
	}
	local, err := ada.List()
	if err != nil {
		t.Fatalf("for-each-ref: %v", err)
	}
	if len(local) != 2 {
		t.Fatalf("want 2 local refs, got %v", local)
	}
}

// An unreachable remote is not a lost race: the write is unsent, and the caller
// has to be able to tell the difference to know whether to re-read or to keep
// the write for later.
//
// Port 1 on the loopback refuses immediately, which is the offline case without
// a test that waits for a real timeout.
func TestUnreachableRemoteIsOfflineNotARace(t *testing.T) {
	ada, _ := twoClones(t)
	git(t, ada.Dir, "remote", "set-url", "origin", "https://127.0.0.1:1/board.git")

	_, err := ada.Write("01TEST", []byte(ticket), "")
	if !errors.Is(err, gitref.ErrOffline) {
		t.Fatalf("want ErrOffline, got %v", err)
	}
	if errors.Is(err, gitref.ErrRaceLost) {
		t.Error("an unreachable remote was reported as a lost race")
	}
	// Nothing was written locally either: the next attempt must still be a
	// first write, not one holding a lease the remote has never seen.
	if _, err := ada.SHA("01TEST"); !errors.Is(err, gitref.ErrNoRef) {
		t.Errorf("a refused push moved the local ref anyway: %v", err)
	}
}

// A remote that answers but is not a repository is neither offline nor a lost
// race: waiting does not fix it and re-reading does not either, so it must not
// be reported as either.
func TestARemoteThatIsNotARepositorySaysSo(t *testing.T) {
	ada, _ := twoClones(t)
	git(t, ada.Dir, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "not-a-repo"))

	_, err := ada.Write("01TEST", []byte(ticket), "")
	if !errors.Is(err, gitref.ErrNoRepo) {
		t.Fatalf("want ErrNoRepo, got %v", err)
	}
}

// Usable answers the permanent question — can this checkout carry refs at all —
// separately from what one push did.
func TestUsableRejectsAMissingRemote(t *testing.T) {
	ada, _ := twoClones(t)
	if err := ada.Usable(); err != nil {
		t.Fatalf("a clone with an origin should be usable: %v", err)
	}
	git(t, ada.Dir, "remote", "remove", "origin")
	if err := ada.Usable(); !errors.Is(err, gitref.ErrNoRepo) {
		t.Errorf("want ErrNoRepo without a remote, got %v", err)
	}

	plain := &gitref.Repo{Dir: t.TempDir(), Remote: "origin"}
	if err := plain.Usable(); !errors.Is(err, gitref.ErrNoRepo) {
		t.Errorf("a directory that is not a repository: want ErrNoRepo, got %v", err)
	}
}

func TestReadWithoutARefSaysSo(t *testing.T) {
	ada, _ := twoClones(t)
	if _, _, err := ada.Read("01NOPE"); !errors.Is(err, gitref.ErrNoRef) {
		t.Fatalf("want ErrNoRef, got %v", err)
	}
}

func mustSHA(t *testing.T, r *gitref.Repo, id string) string {
	t.Helper()
	sha, err := r.SHA(id)
	if err != nil {
		t.Fatalf("sha %s: %v", id, err)
	}
	return sha
}
