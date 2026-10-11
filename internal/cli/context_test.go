package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
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
