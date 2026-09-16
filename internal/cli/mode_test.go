package cli

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// newModeTicket is the shared fixture: a board with one ticket in in-progress.
func newModeTicket(t *testing.T) (dir, handle string) {
	t.Helper()
	t.Setenv("JAIRA_USER", "berk")
	dir = t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 16, 8, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "in-progress",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return dir, ticket.Handle(tk.ID)
}

// The mode is what survives a killed session: written once, it has to come
// back off disk on the next read rather than out of anybody's context.
func TestModeSurvivesRoundTrip(t *testing.T) {
	dir, h := newModeTicket(t)

	if out, err := runCLI(t, dir, "set", h, "mode="+ticket.ModeConversational); err != nil {
		t.Fatalf("set mode: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "show", h, "--json")
	if err != nil {
		t.Fatalf("show --json: %v\n%s", err, out)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("show --json is not json: %v\n%s", err, out)
	}
	if got["mode"] != ticket.ModeConversational {
		t.Errorf("show --json carries mode %v, want %q", got["mode"], ticket.ModeConversational)
	}

	// And it is clearable by hand: nothing clears it automatically, so the
	// empty write is the only way back out of the mode.
	if out, err := runCLI(t, dir, "set", h, "mode="); err != nil {
		t.Fatalf("clear mode: %v\n%s", err, out)
	}
	out, err = runCLI(t, dir, "show", h, "--json")
	if err != nil {
		t.Fatalf("show --json: %v\n%s", err, out)
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["mode"] != "" {
		t.Errorf("mode survived being cleared: %v", got["mode"])
	}
}

// The worker never sees the ticket file — it sees 'show --for-lane --json'. If
// the mode is not in that payload it may as well not exist.
func TestForLaneCarriesMode(t *testing.T) {
	dir, h := newModeTicket(t)

	out, err := runCLI(t, dir, "show", h, "--for-lane", "in-progress", "--json")
	if err != nil {
		t.Fatalf("show --for-lane: %v\n%s", err, out)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if _, ok := got["mode"]; !ok {
		t.Fatalf("the lane payload has no mode key at all:\n%s", out)
	}
	if got["mode"] != "" {
		t.Errorf("a fresh ticket is already in mode %v", got["mode"])
	}

	if out, err := runCLI(t, dir, "set", h, "mode="+ticket.ModeConversational); err != nil {
		t.Fatalf("set mode: %v\n%s", err, out)
	}
	out, err = runCLI(t, dir, "show", h, "--for-lane", "in-progress", "--json")
	if err != nil {
		t.Fatalf("show --for-lane: %v\n%s", err, out)
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["mode"] != ticket.ModeConversational {
		t.Errorf("the lane payload carries mode %v, want %q", got["mode"], ticket.ModeConversational)
	}
}

// A typo stored silently is the failure this mode exists to prevent: the
// person believes they are being asked before each increment, and the worker,
// comparing against exactly one word, commits as usual.
func TestSetRefusesUnknownMode(t *testing.T) {
	dir, h := newModeTicket(t)

	for _, bad := range []string{"chat", "conversation", "Conversational"} {
		out, err := runCLI(t, dir, "set", h, "mode="+bad)
		if err == nil {
			t.Fatalf("set mode=%s was accepted:\n%s", bad, out)
		}
		if !strings.Contains(out+err.Error(), "mode") {
			t.Errorf("set mode=%s failed without naming mode: %v\n%s", bad, err, out)
		}
	}

	out, err := runCLI(t, dir, "show", h, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["mode"] != "" {
		t.Errorf("a refused mode was written anyway: %v", got["mode"])
	}
}

// The gap critique found: the check trimmed, the write path did not, so
// "mode= conversational " came through the front door and was stored with its
// padding. The worker compares against exactly one word, so a padded value
// reads as no mode at all — the same silent failure TestSetRefusesUnknownMode
// covers, reached by a value that is not even a typo.
func TestSetStoresModeTrimmed(t *testing.T) {
	dir, h := newModeTicket(t)

	if out, err := runCLI(t, dir, "set", h, "mode=  "+ticket.ModeConversational+"  "); err != nil {
		t.Fatalf("set padded mode: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "show", h, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["mode"] != ticket.ModeConversational {
		t.Errorf("padded mode stored as %q, want %q", got["mode"], ticket.ModeConversational)
	}
}

// The mode is meant to survive a killed session, which it only does if the
// person who set it can see it is still on. Before this it was writable from
// the CLI and the TUI and readable only in --json.
func TestShowPrintsModeForPeople(t *testing.T) {
	dir, h := newModeTicket(t)

	out, err := runCLI(t, dir, "show", h)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "mode") {
		t.Errorf("a ticket without a mode spends a row on it:\n%s", out)
	}

	if out, err := runCLI(t, dir, "set", h, "mode="+ticket.ModeConversational); err != nil {
		t.Fatalf("set mode: %v\n%s", err, out)
	}
	out, err = runCLI(t, dir, "show", h)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "mode") || !strings.Contains(out, ticket.ModeConversational) {
		t.Errorf("'jaira show' does not print the mode:\n%s", out)
	}
}
