package gitref_test

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
)

const milestoneFile = "---\nname: round-one\ncolour: 45\n---\n\n- 01TEST\n"

// A milestone has to reach everybody the way a ticket does: without a branch
// being merged, and on the same fetch. The file naming twenty tickets is the
// one file everybody is supposed to read; if it waited for a merge it would be
// the one file nobody has.
func TestMilestoneArrivesWithoutASharedBranch(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.WriteMilestone("round-one", []byte(milestoneFile), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Nothing was merged, and grace has none of ada's branches.
	if out := git(t, grace.Dir, "branch", "-a"); strings.Contains(out, "round-one") {
		t.Fatalf("the test is not testing what it says; grace has a branch:\n%s", out)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	got, sha, err := grace.ReadMilestone("round-one")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != milestoneFile {
		t.Errorf("milestone came back changed:\n%q", got)
	}
	if sha == "" {
		t.Error("read gave no lease to write back against")
	}
	if names, err := grace.ListMilestones(); err != nil || len(names) != 1 || names[0] != "round-one" {
		t.Errorf("ListMilestones() = %v, %v; want [round-one]", names, err)
	}
}

// The two namespaces must not be able to see each other. A milestone named
// after a ticket handle and that ticket are two different things, and a reader
// of ticket refs handed a milestone name would be handed something that is not
// a ulid.
func TestTheTwoNamespacesStayApart(t *testing.T) {
	ada, _ := twoClones(t)

	if _, err := ada.Write("01TEST", []byte(ticket), ""); err != nil {
		t.Fatalf("write ticket: %v", err)
	}
	// Deliberately the same key on both sides: this is the collision the path
	// segment exists to prevent.
	if _, err := ada.WriteMilestone("01TEST", []byte(milestoneFile), ""); err != nil {
		t.Fatalf("write milestone: %v", err)
	}

	ids, err := ada.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "01TEST" {
		t.Errorf("List() = %v, want only the ticket", ids)
	}
	names, err := ada.ListMilestones()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != "01TEST" {
		t.Errorf("ListMilestones() = %v, want only the milestone", names)
	}
	// And the two refs carry different bytes, which is the actual claim: one
	// flat namespace would have made the second write overwrite the first.
	gotTicket, _, err := ada.Read("01TEST")
	if err != nil {
		t.Fatal(err)
	}
	gotMilestone, _, err := ada.ReadMilestone("01TEST")
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTicket) != ticket || string(gotMilestone) != milestoneFile {
		t.Error("the two refs with the same key overwrote one another")
	}
	// The remote listing is split the same way, so a command that asks the
	// remote what exists before fetching gets the same answer.
	remote, err := ada.ListRemote()
	if err != nil {
		t.Fatal(err)
	}
	if len(remote) != 1 {
		t.Errorf("ListRemote() = %v, want only the ticket", remote)
	}
	remoteMs, err := ada.ListRemoteMilestones()
	if err != nil {
		t.Fatal(err)
	}
	if len(remoteMs) != 1 {
		t.Errorf("ListRemoteMilestones() = %v, want only the milestone", remoteMs)
	}
}

// A milestone is edited by hand as often as by jaira, so a stale write has to
// lose the same way a ticket's does rather than force-push over somebody's
// edit.
func TestAStaleMilestoneWriteLoses(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.WriteMilestone("round-one", []byte(milestoneFile), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	_, lease, err := grace.ReadMilestone("round-one")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	// ada edits first; grace still holds the lease from before.
	if _, err := ada.WriteMilestone("round-one", []byte(milestoneFile+"- 02TEST\n"), firstSHA(t, ada, "round-one")); err != nil {
		t.Fatalf("second write: %v", err)
	}
	if _, err := grace.WriteMilestone("round-one", []byte("---\nname: round-one\n---\n"), lease); err == nil {
		t.Error("a write against a stale lease was accepted")
	}
}

// The ref is what makes a milestone deletable everywhere: a milestone whose
// file is gone must stop appearing on other boards, the way a logged ticket's
// ref is reaped.
func TestDeletingAMilestoneRefRemovesItEverywhere(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.WriteMilestone("round-one", []byte(milestoneFile), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := ada.DeleteMilestone("round-one", firstSHA(t, ada, "round-one")); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if names, err := grace.ListMilestones(); err != nil || len(names) != 0 {
		t.Errorf("ListMilestones() = %v, %v; want nothing left", names, err)
	}
}

func firstSHA(t *testing.T, r *gitref.Repo, name string) string {
	t.Helper()
	sha, err := r.MilestoneSHA(name)
	if err != nil {
		t.Fatalf("MilestoneSHA(%s): %v", name, err)
	}
	return sha
}
