package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// follows: is the only record that one ticket exists because of another. It was
// written by the sign-off screen but rendered nowhere and absent from --json,
// which made the link real on disk and invisible to everyone reading the board —
// so a follow-up looked like an unrelated ticket.
func TestTicketJSONCarriesFollows(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{
		ID:      "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:   "t",
		Status:  "todo",
		Follows: "01KZZR4CBGDM5T35SZDR72PQYG",
	}

	out := ticketJSON(tk, lanes)

	if out["follows"] != "01KZZR4CBGDM5T35SZDR72PQYG" {
		t.Errorf("follows = %#v, want the full predecessor id", out["follows"])
	}
}

func TestPrintDetailShowsFollows(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{
		ID:      "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:   "t",
		Status:  "todo",
		Follows: "01KZZR4CBGDM5T35SZDR72PQYG",
	}

	var b bytes.Buffer
	printDetail(&b, tk, gate.Env{Lanes: lanes}, 0)

	// The handle, not the full id: it is what every other command takes as an
	// argument, so it is what a reader can act on.
	if !strings.Contains(b.String(), "72PQYG") {
		t.Errorf("printDetail did not show the predecessor handle:\n%s", b.String())
	}
}

// A ticket with no predecessor must not grow an empty follows: row, or every
// ticket on the board carries a line that says nothing.
func TestPrintDetailOmitsEmptyFollows(t *testing.T) {
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo"}

	var b bytes.Buffer
	printDetail(&b, tk, gate.Env{Lanes: lanes}, 0)

	if strings.Contains(b.String(), "follows") {
		t.Errorf("printDetail showed a follows row on a ticket with none:\n%s", b.String())
	}
}

// Work runs on one branch per ticket, so a ticket filed on that branch without
// --follows was found while working that one. Agents left the link out and
// the follow-up fell out of every epic; the branch now supplies it, says so,
// and --follows none turns it off. On a branch naming no ticket nothing
// happens.
func TestCreateFollowsTheTicketTheBranchIsNamedAfter(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir+"/home")
	gitRun(t, dir, "init", "--quiet", "-b", "master")
	gitRun(t, dir, "-c", "user.name=a", "-c", "user.email=a@x", "commit", "--quiet", "--allow-empty", "-m", "init")
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	origin := createTicket(t, dir, "The work")
	if got := jsonOf(t, dir, "create", "On master")["follows"]; got != "" {
		t.Errorf("on master follows = %v; want nothing", got)
	}

	gitRun(t, dir, "checkout", "--quiet", "-b", "feat/"+strings.ToLower(ticket.Handle(origin)))
	got := jsonOf(t, dir, "create", "Found on the way")
	if got["follows"] != origin || got["follows_from_branch"] != "feat/"+strings.ToLower(ticket.Handle(origin)) {
		t.Errorf("on the ticket's branch: follows = %v, from branch = %v; want %s", got["follows"], got["follows_from_branch"], origin)
	}
	if got := jsonOf(t, dir, "create", "Unrelated", "--follows", "none")["follows"]; got != "" {
		t.Errorf("--follows none still wrote %v", got)
	}
	if out, _ := runCLI(t, dir, "create", "Said aloud"); !strings.Contains(out, "Follows "+ticket.Handle(origin)+", from branch") {
		t.Errorf("create did not say where follows came from:\n%s", out)
	}
}
