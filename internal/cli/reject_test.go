package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

func archivedTicket(t *testing.T, dir, id string) *ticket.Ticket {
	t.Helper()
	m, err := filepath.Glob(filepath.Join(dir, ".jaira", "archive", id+"*.md"))
	if err != nil || len(m) != 1 {
		t.Fatalf("archive holds %v for %s (%v)", m, id, err)
	}
	raw, err := os.ReadFile(m[0])
	if err != nil {
		t.Fatal(err)
	}
	d, err := ticket.ParseDoc(raw)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := ticket.Decode(d, m[0])
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

func TestRejectWritesTheReasonAndArchives(t *testing.T) {
	dir, id := dodTestStore(t)

	out, err := runCLI(t, dir, "reject", id, "--why", "nobody needs it", "--force")
	if err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	if got := ticketFiles(t, dir); len(got) != 0 {
		t.Fatalf("the ticket is still on the board: %v", got)
	}
	r := archivedTicket(t, dir, id).Rejected
	if r.Kind != ticket.RejectWontfix || r.Why != "nobody needs it" || r.By == "" || r.At.IsZero() {
		t.Errorf("rejection not written in full: %+v", r)
	}
}

func TestRejectWithoutWhyIsAUsageError(t *testing.T) {
	dir, id := dodTestStore(t)
	for _, args := range [][]string{
		{"reject", id, "--force"},
		{"reject", id, "--why", "  ", "--force"},
		{"reject", id, "--why", "x", "--kind", "later", "--force"},
	} {
		if _, err := runCLI(t, dir, args...); exitCode(err) != ExitUsage {
			t.Errorf("%v: exit %d, want %d (%v)", args, exitCode(err), ExitUsage, err)
		}
	}
	if got := ticketFiles(t, dir); len(got) != 1 {
		t.Errorf("a refused reject moved the ticket: %v", got)
	}
}

func TestRejectDuplicateOfLinksTheOriginal(t *testing.T) {
	dir, id := dodTestStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	orig, err := s.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now().Add(time.Second)),
		ticket.FieldTitle:  "the original",
		ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "reject", id, "--why", "same thing", "--duplicate-of", ticket.Handle(orig.ID), "--force")
	if err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	tk := archivedTicket(t, dir, id)
	if tk.Rejected.Kind != ticket.RejectDuplicate {
		t.Errorf("kind = %q, want duplicate", tk.Rejected.Kind)
	}
	if len(tk.Related) != 1 || tk.Related[0] != orig.ID {
		t.Errorf("related = %v, want the original %s", tk.Related, orig.ID)
	}
}

func TestRejectDuplicateOfAnUnknownTicketIsNotFound(t *testing.T) {
	dir, id := dodTestStore(t)
	_, err := runCLI(t, dir, "reject", id, "--why", "x", "--duplicate-of", "ZZZZZZ", "--force")
	if exitCode(err) != ExitNotFound {
		t.Fatalf("exit %d, want %d (%v)", exitCode(err), ExitNotFound, err)
	}
	if got := ticketFiles(t, dir); len(got) != 1 {
		t.Errorf("the ticket left the board: %v", got)
	}
	if tk := mustLoad(t, dir, id); tk.Rejected.Kind != "" {
		t.Errorf("a refused reject left a rejection behind: %+v", tk.Rejected)
	}
}

// The CLI is what an agent drives, so without --force it is the agent's
// path: refused, and told where a person decides instead.
func TestRejectWithoutForceSendsAnAgentToTheHumanLane(t *testing.T) {
	dir, id := dodTestStore(t)
	_, err := runCLI(t, dir, "reject", id, "--why", "x")
	if exitCode(err) != ExitValidation {
		t.Fatalf("exit %d, want %d (%v)", exitCode(err), ExitValidation, err)
	}
	if !strings.Contains(err.Error(), "--to human --question") {
		t.Errorf("the refusal does not name the human lane: %v", err)
	}
	if got := ticketFiles(t, dir); len(got) != 1 {
		t.Errorf("a refused reject moved the ticket: %v", got)
	}
}

func TestRestoreTakesTheRejectionOff(t *testing.T) {
	dir, id := dodTestStore(t)
	if out, err := runCLI(t, dir, "reject", id, "--why", "x", "--kind", "obsolete", "--force"); err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	file := filepath.Base(archivedTicket(t, dir, id).Path)
	if out, err := runCLI(t, dir, "restore", file); err != nil {
		t.Fatalf("restore: %v\n%s", err, out)
	}
	tk := mustLoad(t, dir, id)
	if tk.Rejected != (ticket.Rejection{}) {
		t.Errorf("rejection survived the restore: %+v", tk.Rejected)
	}
	raw, _ := os.ReadFile(tk.Path)
	if strings.Contains(string(raw), "rejected-") {
		t.Errorf("rejected-* keys left in the file:\n%s", raw)
	}
}

func mustLoad(t *testing.T, dir, id string) *ticket.Ticket {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

func TestArchiveListAndShowSayWhyATicketWasRejected(t *testing.T) {
	dir, id := dodTestStore(t)
	if out, err := runCLI(t, dir, "reject", id, "--why", "overtaken by TPCKZB", "--kind", "obsolete", "--force"); err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "archive")
	if err != nil {
		t.Fatalf("archive: %v\n%s", err, out)
	}
	if !strings.Contains(out, "rejected obsolete: overtaken by TPCKZB") {
		t.Errorf("the archive list does not say why:\n%s", out)
	}
	out, err = runCLI(t, dir, "archive", "--json")
	if err != nil {
		t.Fatalf("archive --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"kind": "obsolete"`) || !strings.Contains(out, `"why": "overtaken by TPCKZB"`) {
		t.Errorf("archive --json does not carry the rejection:\n%s", out)
	}

	out, err = runCLI(t, dir, "show", id)
	if err != nil {
		t.Fatalf("show: %v\n%s", err, out)
	}
	if !strings.Contains(out, "obsolete: overtaken by TPCKZB") {
		t.Errorf("show does not say why:\n%s", out)
	}
	out, err = runCLI(t, dir, "show", id, "--json")
	if err != nil {
		t.Fatalf("show --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"rejected"`) || !strings.Contains(out, `"why": "overtaken by TPCKZB"`) {
		t.Errorf("show --json does not carry the rejection:\n%s", out)
	}
}

// A ticket only archived carries no rejection anywhere it is shown.
func TestAnArchivedTicketIsNotShownAsRejected(t *testing.T) {
	dir, id := dodTestStore(t)
	if out, err := runCLI(t, dir, "archive", id); err != nil {
		t.Fatalf("archive: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "archive")
	if err != nil {
		t.Fatalf("archive: %v\n%s", err, out)
	}
	if strings.Contains(out, "rejected") {
		t.Errorf("a plain archive reads as a rejection:\n%s", out)
	}
	if out, _ := runCLI(t, dir, "show", id); strings.Contains(out, "rejected") {
		t.Errorf("show calls a plain archive rejected:\n%s", out)
	}
}
