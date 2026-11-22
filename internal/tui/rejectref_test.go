package tui

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
)

// R on a card that lives only on its ref rejected nothing: the store looked
// for a file, found none and said "not found". Now the rejection is written
// to the ref, off the update loop, and the card leaves the board.
func TestRejectACardThatLivesOnlyOnItsRef(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)

	ada, adaSync := refClone(t, root, bare, "ada")
	id := createTicket(t, ada, "not wanted", "backlog", "")
	if _, err := adaSync.Flush(); err != nil {
		t.Fatal(err)
	}
	bob, bobSync := refClone(t, root, bare, "bob")
	if _, err := bobSync.Incoming(); err != nil {
		t.Fatal(err)
	}
	m, err := New(bob)
	if err != nil {
		t.Fatal(err)
	}
	card, err := bob.Load(id)
	if err != nil || !card.ReadOnly {
		t.Fatalf("expected a card with no file here: %v %+v", err, card)
	}

	m.startReject(card)
	m.input = "nobody asked for it"
	cmd := m.confirmReject()
	if cmd == nil {
		t.Fatal("a ref-only card was not rejected through a command")
	}
	m.refRejected(cmd().(refRejectedMsg))
	if m.message == "" || m.isErr {
		t.Errorf("notice after rejecting: %q (error %v)", m.message, m.isErr)
	}

	listed, _ := bob.List()
	for _, tk := range listed {
		if tk.ID == id {
			t.Errorf("the rejected card is still on the board")
		}
	}
	// Read back as anybody would: fetched from the remote.
	if err := adaSync.Repo.Fetch(); err != nil {
		t.Fatal(err)
	}
	raw, _, err := adaSync.Repo.Read(id)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ticket.ParseDoc(raw)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ticket.Decode(d, "")
	if err != nil {
		t.Fatal(err)
	}
	if got.Rejected.Kind != ticket.RejectWontfix || got.Rejected.Why != "nobody asked for it" {
		t.Errorf("the ref does not carry the rejection: %+v", got.Rejected)
	}
}

// A milestone row for a rejected ticket opens no card, and says why: not that
// a filter hides it, which would send the reader hunting for one.
func TestOpeningARejectedMilestoneMemberSaysItIsRejected(t *testing.T) {
	m := newTestModel(t, 150, 32)
	m.openMember(memberRow{id: "01TESTREJECTEDTICKET0ABCDEF", place: link.PlaceRejected, where: "rejected"})
	if !strings.Contains(m.message, "rejected and off every board") || strings.Contains(m.message, "filter") {
		t.Errorf("message for a rejected member: %q", m.message)
	}
}
