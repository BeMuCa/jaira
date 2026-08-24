package cli

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

func createdTicket(t *testing.T, dir string) *ticket.Ticket {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatalf("expected one ticket, got %d", len(all))
	}
	return all[0]
}

func emptyStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", dir+"/home")
	t.Setenv("JAIRA_LANES_DIR", dir+"/no-lanes")
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir
}

// Capture belongs to nobody, and adding a way to opt in must not quietly make
// opting in the default. TestCreateLeavesTheTicketUnassigned in claiming_test.go
// is the invariant's own test; this one asserts it survives beside --mine.
func TestPlainCreateStillLeavesTheAssigneeEmpty(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something someone noticed"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if got := createdTicket(t, dir).Assignee; got != "" {
		t.Errorf("assignee = %q, want empty: capture does not claim", got)
	}
}

func TestCreateMineAssignsYou(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "I am taking this one", "--mine"); err != nil {
		t.Fatalf("create --mine: %v\n%s", err, out)
	}
	got := createdTicket(t, dir).Assignee
	if got == "" || got != identity() {
		t.Errorf("assignee = %q, want %q", got, identity())
	}
}

// Naming someone is more specific than naming yourself, so --assignee wins.
func TestCreateAssigneeWinsOverMine(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "for someone else", "--mine", "--assignee", "teammate"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if got := createdTicket(t, dir).Assignee; got != "teammate" {
		t.Errorf("assignee = %q, want teammate", got)
	}
	if strings.EqualFold(createdTicket(t, dir).Assignee, identity()) {
		t.Error("--mine overrode an explicit --assignee")
	}
}
