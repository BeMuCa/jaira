package cli

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A lane that takes definition-of-done gets every item of the body checklist,
// not only the frontmatter scalar 'jaira create --dod' wrote: the review lane
// judges against it without the outcome, and one item of six is a fraction
// reported whole. Proofs stay out — they are the implementer's account.
func TestForLaneDeliversTheWholeDoDChecklist(t *testing.T) {
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
		ticket.FieldGoal: "g", ticket.FieldDoD: "first criterion",
		// Create writes the scalar only; the body checklist is what
		// 'jaira create' and 'jaira dod --add' build, so build it here.
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	h := ticket.Handle(tk.ID)
	if out, err := runCLI(t, dir, "dod", h, "--add", "first criterion", "--add", "second criterion"); err != nil {
		t.Fatalf("dod --add: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "dod", h, "1", "--done", "--proof", "the author's proof"); err != nil {
		t.Fatalf("dod --done: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "show", h, "--for-lane", "review", "--json")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	var payload struct {
		Input map[string]string `json:"input"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	dod := payload.Input[ticket.FieldDoD]
	for _, want := range []string{"1. [x] first criterion", "2. [ ] second criterion"} {
		if !strings.Contains(dod, want) {
			t.Errorf("definition-of-done input lacks %q:\n%s", want, dod)
		}
	}
	if strings.Contains(dod, "the author's proof") {
		t.Errorf("a proof reached the lane:\n%s", dod)
	}
}
