package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// holdsFixture builds a board with inDone tickets already accepted (updated a
// minute apart, oldest first) and one mover in todo, and returns everything a
// test needs to push the mover over the done lane's cap.
func holdsFixture(t *testing.T, inDone int) (string, *ticket.Store, []*ticket.Ticket, *ticket.Ticket) {
	t.Helper()
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(dir, "no-lanes"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	var done []*ticket.Ticket
	for i := 0; i < inDone; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			ticket.FieldID:        ticket.NewID(stamp),
			ticket.FieldTitle:     fmt.Sprintf("t%02d", i),
			ticket.FieldStatus:    "done",
			ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		done = append(done, tk)
	}
	stamp := base.Add(time.Duration(inDone) * time.Minute)
	mover, err := s.Create(map[string]string{
		ticket.FieldID:        ticket.NewID(stamp),
		ticket.FieldTitle:     "mover",
		ticket.FieldStatus:    "todo",
		ticket.FieldCreatedAt: ticket.FormatTime(stamp),
		ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return dir, s, done, mover
}

func countInLane(t *testing.T, s *ticket.Store, lane string) int {
	t.Helper()
	ts, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	n := 0
	for _, tk := range ts {
		if tk.Status == lane {
			n++
		}
	}
	return n
}

// The eleventh arrival pushes the oldest out: done holds ten, and the move that
// lands a ticket there says which ticket left and how to get it back.
func TestMoveIntoACappedLaneTrimsTheOldest(t *testing.T) {
	dir, s, done, mover := holdsFixture(t, 10)
	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if !strings.Contains(out, "left for the logbook") || !strings.Contains(out, "jaira restore") {
		t.Errorf("move output does not report the trim:\n%s", out)
	}
	if !strings.Contains(out, ticket.Handle(done[0].ID)) {
		t.Errorf("move output does not name the oldest ticket %s:\n%s", ticket.Handle(done[0].ID), out)
	}
	if _, err := os.Stat(done[0].Path); !os.IsNotExist(err) {
		t.Errorf("oldest ticket %s is still on the board", ticket.Handle(done[0].ID))
	}
	if n := countInLane(t, s, "done"); n != 10 {
		t.Errorf("done holds %d tickets after the move, want 10", n)
	}
}

// Exactly at the cap nothing moves: ten in done is the promise, not nine.
func TestMoveUpToTheCapTrimsNothing(t *testing.T) {
	dir, s, _, mover := holdsFixture(t, 9)
	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if strings.Contains(out, "logbook") {
		t.Errorf("a move up to the cap reported a trim:\n%s", out)
	}
	if n := countInLane(t, s, "done"); n != 10 {
		t.Errorf("done holds %d tickets after the move, want 10", n)
	}
}

func TestMoveJSONCarriesTheTrim(t *testing.T) {
	dir, _, done, mover := holdsFixture(t, 10)
	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force", "--json")
	if err != nil {
		t.Fatalf("move --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"trimmed"`) || !strings.Contains(out, ticket.Handle(done[0].ID)) {
		t.Errorf("--json output lacks the trimmed ticket %s:\n%s", ticket.Handle(done[0].ID), out)
	}
}

// A failed trim must not fail a move that landed: exit stays 0, the success
// line and next step survive, and the failure is named on stderr. A non-zero
// exit would send an agent into a retry that short-circuits already-in-lane
// and never trims — the lane would stay over cap for good.
func TestMoveSurvivesATrimFailureAndSaysSo(t *testing.T) {
	dir, s, _, mover := holdsFixture(t, 10)
	// One unreadable ticket makes List refuse, so the trim aborts by design.
	if err := os.WriteFile(filepath.Join(s.TicketsDir(), "broken.md"), []byte("not frontmatter at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force")
	if err != nil {
		t.Fatalf("a trim failure failed the move: %v\n%s", err, out)
	}
	if !strings.Contains(out, "→ done") {
		t.Errorf("the move's success line is missing:\n%s", out)
	}
	if !strings.Contains(out, "trimming the done lane failed") {
		t.Errorf("the trim failure is not reported:\n%s", out)
	}
	after, err := s.Load(mover.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != "done" {
		t.Errorf("mover is in %q, want done", after.Status)
	}
}

func TestMoveJSONCarriesTheTrimError(t *testing.T) {
	dir, s, _, mover := holdsFixture(t, 10)
	if err := os.WriteFile(filepath.Join(s.TicketsDir(), "broken.md"), []byte("not frontmatter at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force", "--json")
	if err != nil {
		t.Fatalf("move --json: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"trim_error"`) || !strings.Contains(out, `"moved"`) {
		t.Errorf("--json output lacks trim_error next to the successful move:\n%s", out)
	}
}
