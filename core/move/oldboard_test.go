package move

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// doneDoorway is the done.md a build from before 9ad7aa9 wrote onto a board,
// verbatim. Kept here so this test runs against what those boards actually
// hold, not against a paraphrase of it.
const doneDoorway = `---
id: done
name: Done
after: signoff
precedence: 60
agentic: false
terminal: true
requires-outcome: true
requires-nonmodel-signal: true
requires-commits: true
logbook-on-entry: true
description: Accepted. Every definition-of-done item must be marked done, the plan finished if there is one, and the commits that carry the change recorded. The move that lands here stamps the commits and files the ticket straight into the logbook — 'jaira restore' brings it back.
---
`

// TestMoveIntoDoneOnAnOldBoardFilesNothing is the whole point of the
// correction, end to end: a board built before filing became a decision still
// has logbook-on-entry: true in its done.md, and a move landing there used to
// sweep every finished ticket on the board — other people's included — into the
// logbook. It must now file nothing and leave the others standing.
func TestMoveIntoDoneOnAnOldBoardFilesNothing(t *testing.T) {
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
	// Set the board up, then make it an old one: the doorway back in done.md
	// and no corrections marker, which is the state such a board is in.
	if _, err := lane.Load(dir); err != nil {
		t.Fatal(err)
	}
	lanesDir := lane.ProjectLanesDir(dir)
	if err := os.WriteFile(filepath.Join(lanesDir, "done.md"), []byte(doneDoorway), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(lanesDir, "corrections")); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	stranger := mk(t, s, map[string]string{ticket.FieldStatus: "done"})
	mine := mk(t, s, map[string]string{ticket.FieldStatus: "signoff"})

	lanes, err := lane.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if done, _ := lanes.Get("done"); done == nil || done.LogbookOnEntry {
		t.Fatal("the correction did not reach this board")
	}
	warned := false
	for _, w := range lanes.Warnings {
		if strings.Contains(w, "logbook-on-entry") {
			warned = true
		}
	}
	if !warned {
		t.Errorf("the board was corrected without saying so: %v", lanes.Warnings)
	}

	all, _ := s.List()
	res, err := Move(s, gate.Env{Lanes: lanes, All: all}, mine.ID, Request{
		To: "done", Actor: "berk", Force: true, Folder: "bc-20260906",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Filed || len(res.Trimmed) != 0 {
		t.Fatalf("filed=%v trimmed=%d, want a move that files nothing", res.Filed, len(res.Trimmed))
	}
	if _, err := s.Load(stranger.ID); err != nil {
		t.Errorf("somebody else's finished ticket was swept off the board: %v", err)
	}
	if _, err := s.Load(mine.ID); err != nil {
		t.Errorf("the moved ticket left the board: %v", err)
	}
}
