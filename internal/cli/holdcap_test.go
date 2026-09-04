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
	// These tests pin the holds mechanics, but the builtin done is a doorway
	// (logbook-on-entry) — materialise the board's lane files and give done a
	// cap instead.
	if out, err := runCLI(t, dir, "list"); err != nil {
		t.Fatalf("materialise lanes: %v\n%s", err, out)
	}
	laneFile := filepath.Join(dir, ".jaira", "lanes", "done.md")
	if _, err := os.Stat(laneFile); err != nil {
		t.Fatalf("lane files not materialised: %v", err)
	}
	capped := "---\nid: done\nname: Done\nafter: signoff\nprecedence: 60\nagentic: false\nterminal: true\nrequires-outcome: true\nrequires-nonmodel-signal: true\nrequires-commits: true\nholds: 10\ndescription: capped for these tests\n---\n"
	if err := os.WriteFile(laneFile, []byte(capped), 0o644); err != nil {
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

// The builtin done is a doorway: the move that lands a ticket there stamps its
// commits and files it — and anything still sitting in the lane — straight
// into the logbook. The lane holds nothing; the logbook is the record.
func TestMoveIntoDoneFilesEverythingToTheLogbook(t *testing.T) {
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
	for i := 0; i < 3; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		if _, err := s.Create(map[string]string{
			ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: fmt.Sprintf("t%02d", i),
			ticket.FieldStatus: "done", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, ""); err != nil {
			t.Fatal(err)
		}
	}
	stamp := base.Add(time.Hour)
	mover, err := s.Create(map[string]string{
		ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: "mover",
		ticket.FieldStatus: "todo", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
		ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
	}, map[string][]string{ticket.FieldCommits: {"deadbeef"}}, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "move", ticket.Handle(mover.ID), "--to", "done", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if got := strings.Count(out, "filed to the logbook"); got != 4 {
		t.Errorf("%d filing lines, want 4 (three residents + the mover):\n%s", got, out)
	}
	if n := countInLane(t, s, "done"); n != 0 {
		t.Errorf("done still holds %d tickets, want 0 — it is a doorway", n)
	}
	// The mover's commits survived the filing, stamped onto the filed copy.
	matches, err := filepath.Glob(filepath.Join(s.LogbookDir(), "*", filepath.Base(mover.Path)))
	if err != nil || len(matches) != 1 {
		t.Fatalf("mover not found in the logbook: %v (%d matches)", err, len(matches))
	}
	filed, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(filed), "deadbeef") {
		t.Errorf("the filed ticket lost its commits:\n%s", filed)
	}
}
