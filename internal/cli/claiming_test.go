package cli

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

func runCLI(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir}, args...))
	err := root.Execute()
	return out.String(), err
}

// Capturing and claiming are two acts: a created ticket belongs to nobody, and
// pulling it out of the backlog is what makes it yours. That is the team flow —
// pull, drag into todo, push — and it also keeps a shared backlog from filling
// up with tickets pre-assigned to whoever happened to write them down.
func TestCreateLeavesTheTicketUnassigned(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}

	if out, err := runCLI(t, dir, "create", "capture", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ts, err := s.List()
	if err != nil || len(ts) != 1 {
		t.Fatalf("list: %v (%d tickets)", err, len(ts))
	}
	if ts[0].Assignee != "" {
		t.Errorf("a captured ticket is assigned to %q, want nobody", ts[0].Assignee)
	}
	if ts[0].Creator != "berk" {
		t.Errorf("creator = %q, want berk", ts[0].Creator)
	}
}

func TestMovingAnUnassignedTicketClaimsIt(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "create", "pull me", "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ts, _ := s.List()
	id := ts[0].ID

	// The promotion gate demands an assignee to leave the backlog; the pull
	// itself must satisfy it rather than refusing the very move that claims.
	if out, err := runCLI(t, dir, "move", id, "--to", "todo"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	got, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Assignee != "berk" {
		t.Errorf("assignee after the pull = %q, want berk", got.Assignee)
	}
	if got.Status != "todo" {
		t.Errorf("status = %q, want todo", got.Status)
	}
}

// An explicit assignee is never overwritten by a move: the claim rule applies
// only to tickets that belong to nobody.
func TestMovingAnAssignedTicketDoesNotReassign(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "create", "sams ticket", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "sam"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ts, _ := s.List()
	id := ts[0].ID

	// berk is not sam, so this is a foreign ticket; the owner guard applies and
	// --force records the deliberate override. The point here is only that the
	// claim rule does not touch an existing assignee.
	if out, err := runCLI(t, dir, "move", id, "--to", "todo", "--force"); err != nil {
		t.Fatalf("move --force: %v\n%s", err, out)
	}
	got, err := s.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Assignee != "sam" {
		t.Errorf("assignee = %q, the move must not steal sam's ticket", got.Assignee)
	}
}
