package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A lane may declare notes in input-requires: the bounded input then carries
// the ticket's Progress entries — the handover channel for loop findings. An
// empty journal is omitted, never reported missing.
func TestForLaneDeliversNotesWhenDeclared(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 5, 8, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "in-progress",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	h := ticket.Handle(tk.ID)

	// No notes yet: the field is absent, and NOT listed as missing.
	out, err := runCLI(t, dir, "show", h, "--for-lane", "in-progress", "--json")
	if err != nil {
		t.Fatalf("show --for-lane: %v\n%s", err, out)
	}
	if strings.Contains(out, `"notes"`) {
		t.Errorf("an empty journal produced a notes field:\n%s", out)
	}
	if strings.Contains(out, "notes (") {
		t.Errorf("an empty journal was reported missing:\n%s", out)
	}

	if _, err := runCLI(t, dir, "note", h, "testing: der writer buffert alles; fix: chunked schreiben"); err != nil {
		t.Fatal(err)
	}
	out, err = runCLI(t, dir, "show", h, "--for-lane", "in-progress", "--json")
	if err != nil {
		t.Fatalf("show --for-lane: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"notes"`) || !strings.Contains(out, "chunked schreiben") {
		t.Errorf("the declared notes input does not carry the note:\n%s", out)
	}

	// A builtin lane that does not declare notes never receives them.
	out, err = runCLI(t, dir, "show", h, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	if strings.Contains(out, "chunked schreiben") {
		t.Errorf("review got notes it never declared:\n%s", out)
	}
}
