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
