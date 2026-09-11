package move

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

func fixture(t *testing.T) (*ticket.Store, gate.Env) {
	t.Helper()
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
	lanes, err := lane.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	all, _ := s.List()
	return s, gate.Env{Lanes: lanes, All: all}
}

func mk(t *testing.T, s *ticket.Store, fields map[string]string) *ticket.Ticket {
	t.Helper()
	base := map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 6, 8, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}
	for k, v := range fields {
		base[k] = v
	}
	tk, err := s.Create(base, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

// A refused move writes nothing; force writes and reports what it overrode.
func TestMoveRefusesThenForceOverrides(t *testing.T) {
	s, env := fixture(t)
	tk := mk(t, s, nil) // unspecified: the promotion gate refuses todo

	res, err := Move(s, env, tk.ID, Request{To: "todo", Actor: "berk"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Refused) == 0 {
		t.Fatal("an unspecified ticket entered todo")
	}
	if got, _ := s.Load(tk.ID); got.Status != "backlog" {
		t.Fatalf("a refused move wrote status %q", got.Status)
	}

	res, err = Move(s, env, tk.ID, Request{To: "todo", Actor: "berk", Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Overrode) == 0 || len(res.Refused) != 0 {
		t.Fatalf("force reported refused=%d overrode=%d", len(res.Refused), len(res.Overrode))
	}
	if got, _ := s.Load(tk.ID); got.Status != "todo" {
		t.Fatalf("forced move landed in %q", got.Status)
	}
}

// The pull claims: an unassigned ticket gains the mover, staged for the gate
// and written with the move. An existing assignee is never overwritten.
func TestMoveClaimsOnPull(t *testing.T) {
	s, env := fixture(t)
	tk := mk(t, s, map[string]string{
		ticket.FieldGoal: "g", ticket.FieldContext: "c", ticket.FieldDoD: "d",
	})
	res, err := Move(s, env, tk.ID, Request{To: "todo", Actor: "berk", ClaimOnPull: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Refused) > 0 {
		t.Fatalf("refused: %v", res.Refused.Err())
	}
	if !res.Claimed || res.Ticket.Assignee != "berk" {
		t.Fatalf("claimed=%v assignee=%q", res.Claimed, res.Ticket.Assignee)
	}
	res, err = Move(s, env, tk.ID, Request{To: "backlog", Actor: "somebody-else", ClaimOnPull: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Claimed || res.Ticket.Assignee != "berk" {
		t.Fatalf("an existing assignee was overwritten: claimed=%v assignee=%q", res.Claimed, res.Ticket.Assignee)
	}
}

// Stage persists before the gate — a refused move keeps the staged fields.
func TestMoveStagePersistsThroughARefusal(t *testing.T) {
	s, env := fixture(t)
	tk := mk(t, s, nil)
	res, err := Move(s, env, tk.ID, Request{
		To: "todo", Actor: "berk",
		Stage: func(t *ticket.Ticket) error {
			return t.Doc().SetScalar(ticket.FieldGoal, "staged goal")
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Refused) == 0 {
		t.Fatal("expected a refusal (dod/context still missing)")
	}
	if got, _ := s.Load(tk.ID); got.Goal != "staged goal" {
		t.Fatalf("the staged field did not survive the refusal: %q", got.Goal)
	}
}

// Landing in a doorway settles it: the move files the ticket to the logbook.
//
// The doorway is declared here rather than taken from the shipped lanes,
// because the lanes this tool ships no longer have one — filing is a decision
// somebody makes with 'jaira logbook', not something finishing a ticket does.
// The mechanism stays for a board that asks for it, and this is what pins it.
func TestMoveSettlesTheDoorway(t *testing.T) {
	s, env := fixture(t)
	if done, ok := env.Lanes.Get("done"); ok {
		done.LogbookOnEntry = true
	}
	tk := mk(t, s, map[string]string{ticket.FieldStatus: "signoff"})
	res, err := Move(s, env, tk.ID, Request{
		To: "done", Actor: "berk", Force: true, Folder: "bc-20260906",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Filed || len(res.Trimmed) != 1 {
		t.Fatalf("filed=%v trimmed=%d, want the doorway to file the arrival", res.Filed, len(res.Trimmed))
	}
	if _, err := os.Stat(res.Trimmed[0].Path); err != nil {
		t.Fatalf("filed ticket not in the logbook: %v", err)
	}
}

// And the shipped default is the opposite: landing in done files nothing, so a
// ticket waits there with everybody else's until somebody cuts.
func TestMoveIntoDoneFilesNothingByItself(t *testing.T) {
	s, env := fixture(t)
	tk := mk(t, s, map[string]string{ticket.FieldStatus: "signoff"})
	res, err := Move(s, env, tk.ID, Request{
		To: "done", Actor: "berk", Force: true, Folder: "bc-20260906",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Filed || len(res.Trimmed) != 0 {
		t.Fatalf("filed=%v trimmed=%d, want a move that files nothing", res.Filed, len(res.Trimmed))
	}
	after, err := s.Load(tk.ID)
	if err != nil {
		t.Fatalf("the ticket left the board anyway: %v", err)
	}
	if after.Status != "done" {
		t.Errorf("status is %q, want done", after.Status)
	}
}
