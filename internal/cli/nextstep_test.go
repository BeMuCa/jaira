package cli

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// movableTicket returns a specified, claimed ticket — one the promotion gate
// lets out of the backlog — so a test can move it and read what the move says.
func movableTicket(t *testing.T) (dir, id string) {
	t.Helper()
	t.Setenv("JAIRA_USER", "berk")
	dir = t.TempDir()
	t.Setenv("JAIRA_HOME", dir+"/home")
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "create", "drive the lane",
		"--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	ts, err := s.List()
	if err != nil || len(ts) != 1 {
		t.Fatalf("expected one ticket, got %d (%v)", len(ts), err)
	}
	id = ts[0].ID
	if out, err := runCLI(t, dir, "claim", id); err != nil {
		t.Fatalf("claim: %v\n%s", err, out)
	}
	return dir, id
}

// The move is the moment a session should carry on working: the lane's prompt
// and the agent block both already said so and were measured not to be enough,
// because the block is read once at the start of a session. So the move names
// the command itself.
func TestMoveIntoAnAgenticLaneNamesTheNextCommand(t *testing.T) {
	dir, id := movableTicket(t)
	if out, err := runCLI(t, dir, "move", id, "--to", "todo"); err != nil {
		t.Fatalf("move to todo: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "move", id, "--to", "in-progress")
	if err != nil {
		t.Fatalf("move to in-progress: %v\n%s", err, out)
	}
	if !strings.Contains(out, "in-progress is an agentic lane") {
		t.Errorf("move output does not say the lane is agentic:\n%s", out)
	}
	if !strings.Contains(out, "--for-lane in-progress --json") {
		t.Errorf("move output does not name 'jaira show --for-lane':\n%s", out)
	}
}

// A lane with no prompt has no step to run, and the human lanes are the ones a
// session must not work at all — naming a command for either would invite
// exactly the wrong thing.
func TestMoveIntoANonAgenticLaneNamesNoNextCommand(t *testing.T) {
	dir, id := movableTicket(t)

	out, err := runCLI(t, dir, "move", id, "--to", "todo")
	if err != nil {
		t.Fatalf("move to todo: %v\n%s", err, out)
	}
	if strings.Contains(out, "agentic lane") || strings.Contains(out, "--for-lane") {
		t.Errorf("non-agentic lane 'todo' still got a next-step line:\n%s", out)
	}
}

// --json is what an agent parses; a sentence for a human on the same stream
// would break it. next_lane is already in the payload.
func TestMoveJSONCarriesNoNextStepSentence(t *testing.T) {
	dir, id := movableTicket(t)
	if out, err := runCLI(t, dir, "move", id, "--to", "todo"); err != nil {
		t.Fatalf("move to todo: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "move", id, "--to", "in-progress", "--json")
	if err != nil {
		t.Fatalf("move --json: %v\n%s", err, out)
	}
	if strings.Contains(out, "agentic lane") {
		t.Errorf("--json output carries the human sentence:\n%s", out)
	}
}
