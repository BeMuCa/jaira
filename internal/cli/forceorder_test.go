package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A forced move ends on its success line: the override bullets come first, so
// the last line never reads like a refusal — 'move --force | tail -1' once
// raised a false alarm that cost a day (BDV0HM).
func TestForcedMoveEndsOnItsSuccessLine(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	// Unspecified on purpose: the promotion gate refuses, --force overrides.
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "move", ticket.Handle(tk.ID), "--to", "todo", "--force")
	if err != nil {
		t.Fatalf("move --force: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Overrode") {
		t.Fatalf("no override happened — the fixture no longer exercises the order:\n%s", out)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	last := lines[len(lines)-1]
	if !strings.Contains(last, "→ todo") {
		t.Errorf("the last line is %q, want the success line ending the output:\n%s", last, out)
	}
	if strings.HasPrefix(strings.TrimSpace(last), "-") {
		t.Errorf("the last line is a refusal bullet:\n%s", out)
	}
}
