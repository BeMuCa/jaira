package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
)

// The agent carrying a ticket learns which lane is not its own from two
// places — 'jaira lanes' for the board, 'show --for-lane --json' for the lane
// in hand. On a board whose lanes say nothing, both have to report the
// default: review is fresh, and says it was implied.
func TestLanesAndForLaneReportContext(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "none.md"))
	dir, h := newModeTicket(t)

	out, err := runLanes(t, dir, "--json")
	if err != nil {
		t.Fatalf("lanes --json: %v\n%s", err, out)
	}
	var payload struct {
		Lanes []map[string]any `json:"lanes"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	byID := map[string]map[string]any{}
	for _, l := range payload.Lanes {
		byID[l["id"].(string)] = l
	}
	if r := byID["review"]; r["context"] != "fresh" || r["context_implied"] != true {
		t.Errorf("review in lanes --json: context=%v implied=%v, want fresh/true", r["context"], r["context_implied"])
	}
	if p := byID["in-progress"]; p["context"] != "" || p["context_implied"] != false {
		t.Errorf("in-progress in lanes --json: context=%v implied=%v, want empty/false", p["context"], p["context_implied"])
	}

	table, err := runLanes(t, dir)
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, table)
	}
	if !strings.Contains(table, "CONTEXT") || !strings.Contains(table, "fresh*") || !strings.Contains(table, "* no lane on this board declares context:") {
		t.Errorf("lanes table does not show the implied fresh lane and why:\n%s", table)
	}

	if got := jsonCLI(t, dir, "show", h, "--for-lane", "review")["context"]; got != "fresh" {
		t.Errorf("show --for-lane review --json context = %v, want fresh", got)
	}
	got := jsonCLI(t, dir, "show", h, "--for-lane", "in-progress")
	if v, ok := got["context"]; !ok || v != "" {
		t.Errorf("show --for-lane in-progress --json context = %v (present %v), want an empty key", v, ok)
	}
}

// A typo in a lane's context: is loud everywhere a person looks: 'jaira
// lanes' marks the lane and does not claim the board declares nothing, and
// 'jaira validate' fails with the validation exit code.
func TestAnInvalidContextIsReportedNotDropped(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "none.md"))
	dir, _ := newModeTicket(t)
	if out, err := runLanes(t, dir); err != nil {
		t.Fatalf("lanes: %v\n%s", err, out)
	}
	p := filepath.Join(lane.ProjectLanesDir(dir), "review.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(b), "---\n", "---\ncontext: Fresh\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}

	table, err := runLanes(t, dir)
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, table)
	}
	if strings.Contains(table, "no lane on this board declares context:") || strings.Contains(table, "fresh*") {
		t.Errorf("lanes still claims the default decides:\n%s", table)
	}
	if !strings.Contains(table, "fresh!") || !strings.Contains(table, "! "+invalidContextNote) {
		t.Errorf("lanes does not mark review's bad context:\n%s", table)
	}

	out, err := runCLI(t, dir, "validate")
	var ce *codedError
	if !errors.As(err, &ce) || ce.code != ExitValidation {
		t.Fatalf("validate err = %v, want exit %d\n%s", err, ExitValidation, out)
	}
	if !strings.Contains(out, `context "Fresh"`) {
		t.Errorf("validate does not name the bad value:\n%s", out)
	}
}
