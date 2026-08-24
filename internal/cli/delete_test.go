package cli

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// runCLIStdin is runCLI with something to answer a prompt with.
func runCLIStdin(t *testing.T, dir, stdin string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetIn(strings.NewReader(stdin))
	root.SetArgs(append([]string{"-C", dir}, args...))
	err := root.Execute()
	return out.String(), err
}

func ticketFiles(t *testing.T, dir string) []string {
	t.Helper()
	m, err := filepath.Glob(filepath.Join(dir, ".jaira", "tickets", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// Deletion is the one irreversible thing this tool does, so the confirmation is
// the handle typed back rather than a yes: a stray return key must not be able
// to remove a ticket.
func TestDeleteRefusesAMistypedConfirmation(t *testing.T) {
	dir, id := dodTestStore(t)
	handle := ticket.Handle(id)
	before := ticketFiles(t, dir)

	for _, answer := range []string{"\n", "y\n", "yes\n", strings.ToLower(handle[:3]) + "\n", ""} {
		out, err := runCLIStdin(t, dir, answer, "delete", id)
		if err == nil {
			t.Fatalf("answer %q deleted the ticket:\n%s", answer, out)
		}
		if !strings.Contains(err.Error(), "nothing was deleted") {
			t.Errorf("answer %q: unhelpful refusal: %v", answer, err)
		}
	}
	if got := ticketFiles(t, dir); len(got) != len(before) {
		t.Errorf("files changed after refusals: %v", got)
	}
}

func TestDeleteRemovesTheFileWhenConfirmed(t *testing.T) {
	dir, id := dodTestStore(t)
	handle := ticket.Handle(id)

	out, err := runCLIStdin(t, dir, handle+"\n", "delete", id)
	if err != nil {
		t.Fatalf("delete: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Deleted "+handle) {
		t.Errorf("the deletion was not reported:\n%s", out)
	}
	if got := ticketFiles(t, dir); len(got) != 0 {
		t.Errorf("the file survived: %v", got)
	}

	// The board is still coherent afterwards — that is the whole reason the
	// referenced-by check exists.
	if out, err := runCLI(t, dir, "validate"); err != nil {
		t.Errorf("validate is not clean after a delete: %v\n%s", err, out)
	}
}

// --force is for scripts, and it still says what it removed.
func TestDeleteForceSkipsThePromptAndSaysWhatItDid(t *testing.T) {
	dir, id := dodTestStore(t)
	out, err := runCLIStdin(t, dir, "", "delete", id, "--force")
	if err != nil {
		t.Fatalf("delete --force: %v\n%s", err, out)
	}
	if !strings.Contains(out, ticket.Handle(id)) || !strings.Contains(out, "Deleted") {
		t.Errorf("--force did not report the deletion:\n%s", out)
	}
	if got := ticketFiles(t, dir); len(got) != 0 {
		t.Errorf("the file survived --force: %v", got)
	}
}

// A dangling blocked-by is an error 'validate' reports and a dependency that can
// never clear. Deleting into that state is refused, so the board is never made
// invalid by the command that was tidying it.
func TestDeleteRefusesAReferencedTicket(t *testing.T) {
	dir, id := dodTestStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	dependent, err := s.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  "waits on the other",
		ticket.FieldStatus: "backlog",
	}, map[string][]string{ticket.FieldBlockedBy: {id}}, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCLIStdin(t, dir, ticket.Handle(id)+"\n", "delete", id)
	if err == nil {
		t.Fatalf("a ticket another one depends on was deleted:\n%s", out)
	}
	if !strings.Contains(err.Error(), ticket.Handle(dependent.ID)) {
		t.Errorf("the refusal does not name what still points at it: %v", err)
	}
	if len(ticketFiles(t, dir)) != 2 {
		t.Error("something was deleted despite the refusal")
	}

	// --force is the user's call, and it goes through.
	if out, err := runCLIStdin(t, dir, "", "delete", id, "--force"); err != nil {
		t.Fatalf("--force did not override the reference check: %v\n%s", err, out)
	}
}

// Deleting a ticket must not turn every other command into an error. The board
// still lists, and a ticket that merely followed the deleted one still reads.
func TestOtherCommandsSurviveADelete(t *testing.T) {
	dir, id := dodTestStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	follower, err := s.Create(map[string]string{
		ticket.FieldID:      ticket.NewID(time.Now()),
		ticket.FieldTitle:   "came out of the other",
		ticket.FieldStatus:  "backlog",
		ticket.FieldFollows: id,
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	// follows is reported by the reference check, so --force is the way past it.
	if out, err := runCLIStdin(t, dir, "", "delete", id, "--force"); err != nil {
		t.Fatalf("delete --force: %v\n%s", err, out)
	}
	for _, args := range [][]string{{"list"}, {"next"}, {"show", ticket.Handle(follower.ID)}} {
		if out, err := runCLI(t, dir, args...); err != nil {
			t.Errorf("%v failed after a delete: %v\n%s", args, err, out)
		}
	}
}
