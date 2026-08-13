package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// lanesTestCatalogue isolates a test from the real user catalogue at
// ~/.jaira/lanes, the same way core/lane's own tests do.
func lanesTestCatalogue(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", dir)
	return dir
}

// runLanes drives the real 'lanes' command tree, the same path a user or
// agent invokes it through, and returns its combined stdout/stderr.
func runLanes(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "lanes"}, args...))
	err := root.Execute()
	return out.String(), err
}

// TestLanesJSONIncludesPromptAndCreator asserts every lane in 'lanes --json'
// carries its prompt and creator, not only the fields the table already shows.
func TestLanesJSONIncludesPromptAndCreator(t *testing.T) {
	lanesTestCatalogue(t)
	dir := t.TempDir()

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
	if len(payload.Lanes) == 0 {
		t.Fatal("no lanes in json output")
	}

	found := false
	for _, l := range payload.Lanes {
		if _, ok := l["creator"]; !ok {
			t.Errorf("lane %v missing creator field", l["id"])
		}
		if _, ok := l["prompt"]; !ok {
			t.Errorf("lane %v missing prompt field", l["id"])
		}
		if l["id"] == "review" {
			found = true
			prompt, _ := l["prompt"].(string)
			if !strings.Contains(prompt, "Review this change") {
				t.Errorf("review prompt = %q, want it to contain the built-in prompt text", prompt)
			}
			if l["creator"] != "jaira" {
				t.Errorf("review creator = %v, want %q", l["creator"], "jaira")
			}
		}
	}
	if !found {
		t.Fatal("review lane not present in --json output")
	}
}

// TestLanesShowPrintsFullContract asserts 'lanes show <id>' carries the id,
// name, anchor, precedence, tier, creator, source and the full prompt body
// with no ticket in hand.
func TestLanesShowPrintsFullContract(t *testing.T) {
	lanesTestCatalogue(t)
	dir := t.TempDir()

	out, err := runLanes(t, dir, "show", "review")
	if err != nil {
		t.Fatalf("lanes show review: %v\n%s", err, out)
	}
	for _, want := range []string{
		"ID:          review",
		"Name:        Review",
		"Anchor:      human",
		"Precedence:  50",
		"Tier:        strong",
		"Creator:     jaira",
		"Source:      built-in",
		"Review this change against its definition of done.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("lanes show review missing %q, got:\n%s", want, out)
		}
	}
}

// TestLanesShowUnknownExitsUsage asserts an id no lane claims exits 2 with a
// machine-readable reason, matching how 'move' already reports it.
func TestLanesShowUnknownExitsUsage(t *testing.T) {
	lanesTestCatalogue(t)
	dir := t.TempDir()

	_, err := runLanes(t, dir, "show", "nope")
	if err == nil {
		t.Fatal("expected an error for an unknown lane id")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitUsage {
		t.Errorf("code = %d, want %d", ce.code, ExitUsage)
	}
	if ce.reason != "no_such_lane" {
		t.Errorf("reason = %q, want %q", ce.reason, "no_such_lane")
	}
}

// TestLanesShowJSONCarriesFullContract asserts 'lanes show <id> --json' is
// valid JSON carrying the same fields the human output shows.
func TestLanesShowJSONCarriesFullContract(t *testing.T) {
	lanesTestCatalogue(t)
	dir := t.TempDir()

	out, err := runLanes(t, dir, "show", "review", "--json")
	if err != nil {
		t.Fatalf("lanes show review --json: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	for _, key := range []string{
		"id", "name", "after", "precedence", "model_tier",
		"input_requires", "output_produces", "creator", "source",
		"prompt", "description", "overrides",
	} {
		if _, ok := payload[key]; !ok {
			t.Errorf("lanes show --json missing field %q: %v", key, payload)
		}
	}
	if payload["id"] != "review" {
		t.Errorf("id = %v, want %q", payload["id"], "review")
	}
}

// TestLanesTypoFallsThroughToUsageError asserts an unrecognised subcommand
// name falls through to the parent's noArgs() check and exits 2, rather than
// being silently swallowed. This stays true across cobra upgrades only if it
// is asserted.
func TestLanesTypoFallsThroughToUsageError(t *testing.T) {
	lanesTestCatalogue(t)
	dir := t.TempDir()

	_, err := runLanes(t, dir, "typo")
	if err == nil {
		t.Fatal("expected 'jaira lanes typo' to fail")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitUsage {
		t.Errorf("code = %d, want %d", ce.code, ExitUsage)
	}
}

// TestLanesPathNamesBothDirectoriesAndMarksActive asserts 'lanes path' prints
// the catalogue and project directories and marks whichever is authoritative
// under D-03, both outside a project and once one exists.
func TestLanesPathNamesBothDirectoriesAndMarksActive(t *testing.T) {
	catalogue := lanesTestCatalogue(t)

	// Outside any project: only the catalogue is in force.
	outside := t.TempDir()
	out, err := runLanes(t, outside, "path")
	if err != nil {
		t.Fatalf("lanes path (outside a project): %v\n%s", err, out)
	}
	if !strings.Contains(out, catalogue) || !strings.Contains(out, "(active)") {
		t.Errorf("lanes path outside a project = %q, want the catalogue marked active", out)
	}
	if !strings.Contains(out, "not in a project directory") {
		t.Errorf("lanes path outside a project = %q, want it to say there is no project", out)
	}

	jout, err := runLanes(t, outside, "path", "--json")
	if err != nil {
		t.Fatalf("lanes path --json (outside a project): %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	if payload["active"] != "catalogue" || payload["in_project"] != false {
		t.Errorf("lanes path --json outside a project = %v", payload)
	}

	// Inside a project with no lane directory of its own: still the catalogue.
	root := t.TempDir()
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	out2, err := runLanes(t, root, "path")
	if err != nil {
		t.Fatalf("lanes path (empty project): %v\n%s", err, out2)
	}
	projDir := lane.ProjectLanesDir(root)
	if !strings.Contains(out2, projDir) {
		t.Errorf("lanes path (empty project) = %q, want it to name %q", out2, projDir)
	}
	catalogueLine, projectLine := pathLines(t, out2)
	if !strings.Contains(catalogueLine, "(active)") || strings.Contains(projectLine, "(active)") {
		t.Errorf("lanes path (empty project) = %q, want the catalogue marked active", out2)
	}

	// Once the project directory holds a lane file, it becomes authoritative.
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "todo.md"), []byte("---\nid: todo\nname: Todo\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out3, err := runLanes(t, root, "path")
	if err != nil {
		t.Fatalf("lanes path (active project): %v\n%s", err, out3)
	}
	catalogueLine3, projectLine3 := pathLines(t, out3)
	if strings.Contains(catalogueLine3, "(active)") || !strings.Contains(projectLine3, "(active)") {
		t.Errorf("lanes path (active project) = %q, want the project directory marked active", out3)
	}

	jout3, err := runLanes(t, root, "path", "--json")
	if err != nil {
		t.Fatalf("lanes path --json (active project): %v\n%s", err, jout3)
	}
	var payload3 map[string]any
	if err := json.Unmarshal([]byte(jout3), &payload3); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout3)
	}
	if payload3["active"] != "project" || payload3["in_project"] != true {
		t.Errorf("lanes path --json (active project) = %v", payload3)
	}
}

// pathLines splits 'lanes path' human output into its Catalogue and Project
// lines.
func pathLines(t *testing.T, out string) (catalogueLine, projectLine string) {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "Catalogue:"):
			catalogueLine = line
		case strings.HasPrefix(line, "Project:"):
			projectLine = line
		}
	}
	return catalogueLine, projectLine
}

// TestLanesTemplateParses asserts 'lanes template' prints a lane file the
// loader accepts without complaint, so the template cannot rot away from the
// parser it feeds.
func TestLanesTemplateParses(t *testing.T) {
	catalogue := lanesTestCatalogue(t)
	dir := t.TempDir()

	out, err := runLanes(t, dir, "template")
	if err != nil {
		t.Fatalf("lanes template: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(catalogue, "from-template.md"), []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}

	set, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := set.Get("my-lane"); !ok {
		t.Fatalf("template's own id %q did not load: ids=%v warnings=%v", "my-lane", set.IDs(), set.Warnings)
	}
	for _, w := range set.Warnings {
		if strings.Contains(w, "from-template.md") {
			t.Errorf("template produced a load warning: %s", w)
		}
	}
}
