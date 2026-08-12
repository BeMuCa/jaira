package ticket

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The board lists tickets by reading only the first 16KB. A checklist beyond
// that point would be invisible to the card counts, which would then quietly
// disagree with the gate.
func TestListSeesChecklistsPastTheProbe(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}

	filler := strings.Repeat("Some long description line that pads the body out.\n", 500)
	body := "# big\n\n## Description\n\n" + filler +
		"\n## Definition of Done\n\n- [x] one\n- [ ] two\n"
	if len(body) < 16<<10 {
		t.Fatalf("test body is only %d bytes, not past the probe", len(body))
	}
	if _, err := s.Create(map[string]string{
		FieldID: NewID(timeNowForTest()), FieldTitle: "big", FieldStatus: "backlog",
	}, nil, body); err != nil {
		t.Fatal(err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("got %d tickets", len(list))
	}
	if n := len(list[0].DoDItems); n != 2 {
		t.Errorf("List saw %d checklist items, want 2 — the body was truncated at the probe", n)
	}
	_ = os.Stdout
}

func timeNowForTest() time.Time { return time.Now() }
