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
	// The cockpit mod tells a person's lane from an agent's by this field.
	if s := byID["signoff"]; s["human_exit"] != true {
		t.Errorf("signoff in lanes --json: human_exit=%v, want true", s["human_exit"])
	}
	if r := byID["review"]; r["human_exit"] != false {
		t.Errorf("review in lanes --json: human_exit=%v, want false", r["human_exit"])
	}

	table, err := runLanes(t, dir)
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, table)
	}
	if !strings.Contains(table, "CONTEXT") || !strings.Contains(table, "fresh*") || !strings.Contains(table, "* "+impliedContextNote) {
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
	setLaneContext(t, dir, "review", "Fresh")

	table, err := runLanes(t, dir)
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, table)
	}
	if strings.Contains(table, impliedContextNote) || strings.Contains(table, "fresh*") {
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

// A typo never takes freshness away. "context: Shared" on in-progress runs
// in-progress fresh, but it is no declaration the board can act on, so review
// keeps the default's fresh as well; 'jaira lanes' shows both marks with both
// footnotes and 'jaira validate' still fails.
func TestAnInvalidContextLeavesTheDefaultFresh(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "none.md"))
	dir, _ := newModeTicket(t)
	if out, err := runLanes(t, dir); err != nil {
		t.Fatalf("lanes: %v\n%s", err, out)
	}
	setLaneContext(t, dir, "in-progress", "Shared")

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
	if p := byID["in-progress"]; p["context"] != "fresh" || p["context_invalid"] != "Shared" {
		t.Errorf("in-progress: context=%v invalid=%v, want fresh/Shared", p["context"], p["context_invalid"])
	}
	if r := byID["review"]; r["context"] != "fresh" || r["context_implied"] != true {
		t.Errorf("review: context=%v implied=%v, want fresh/true: the typo must not switch the default off", r["context"], r["context_implied"])
	}

	table, err := runLanes(t, dir)
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, table)
	}
	for _, want := range []string{"fresh!", "fresh*", "! " + invalidContextNote, "* " + impliedContextNote} {
		if !strings.Contains(table, want) {
			t.Errorf("lanes table lacks %q:\n%s", want, table)
		}
	}

	out, err = runCLI(t, dir, "validate")
	var ce *codedError
	if !errors.As(err, &ce) || ce.code != ExitValidation {
		t.Fatalf("validate err = %v, want exit %d\n%s", err, ExitValidation, out)
	}
	if !strings.Contains(out, `context "Shared"`) {
		t.Errorf("validate does not name the bad value:\n%s", out)
	}
}

// setLaneContext writes a context: line into the frontmatter of one of the
// board's lane files.
func setLaneContext(t *testing.T, dir, id, value string) {
	t.Helper()
	p := filepath.Join(lane.ProjectLanesDir(dir), id+".md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(b), "---\n", "---\ncontext: "+value+"\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}
