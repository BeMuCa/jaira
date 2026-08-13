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
		"Merge rank:  50",
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

// TestLanesTableLabelsRankNotPosition asserts the human table column header
// no longer claims the number is a display position (task 10): the value is
// a merge rank, and column order follows the after: anchor, not this field.
func TestLanesTableLabelsRankNotPosition(t *testing.T) {
	lanesTestCatalogue(t)
	out, err := runLanes(t, t.TempDir())
	if err != nil {
		t.Fatalf("lanes: %v\n%s", err, out)
	}
	if strings.Contains(out, "PREC") {
		t.Errorf("lanes table = %q, must not label the column PREC (a position claim)", out)
	}
	if !strings.Contains(out, "RANK") {
		t.Errorf("lanes table = %q, want a RANK column labelling the merge rank", out)
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

// lanesTestProject initialises a board at a fresh directory, isolated from
// the real ~/.jaira, for tests that exercise use/publish/adopt against a
// project rather than only the catalogue.
func lanesTestProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestLanesUseWritesToProjectDir asserts 'lanes use <id>' copies the named
// lane, verbatim, into the project's own lane directory — the same call the
// lane settings screen's 'u' key makes.
func TestLanesUseWritesToProjectDir(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)

	out, err := runLanes(t, root, "use", "review")
	if err != nil {
		t.Fatalf("lanes use review: %v\n%s", err, out)
	}
	dst := filepath.Join(lane.ProjectLanesDir(root), "review.md")
	if _, statErr := os.Stat(dst); statErr != nil {
		t.Errorf("expected %s to exist: %v", dst, statErr)
	}
	if !strings.Contains(out, dst) {
		t.Errorf("lanes use output = %q, want it to name %s", out, dst)
	}
}

// TestLanesUseRefusesThenOverwritesWithForce asserts a second 'use' of the
// same lane refuses without --force and succeeds with it.
func TestLanesUseRefusesThenOverwritesWithForce(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)

	if _, err := runLanes(t, root, "use", "review"); err != nil {
		t.Fatalf("first use: %v", err)
	}
	_, err := runLanes(t, root, "use", "review")
	if err == nil {
		t.Fatal("expected a second 'use' of the same lane to refuse")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitValidation {
		t.Errorf("code = %d, want %d", ce.code, ExitValidation)
	}

	if _, err := runLanes(t, root, "use", "review", "--force"); err != nil {
		t.Fatalf("use --force: %v", err)
	}
}

// TestLanesUseJSONCarriesPath asserts --json reports the id and the path
// written, for an agent to consume without scraping a sentence.
func TestLanesUseJSONCarriesPath(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)

	out, err := runLanes(t, root, "use", "review", "--json")
	if err != nil {
		t.Fatalf("lanes use --json: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if payload["id"] != "review" {
		t.Errorf("id = %v, want %q", payload["id"], "review")
	}
	if payload["path"] != filepath.Join(lane.ProjectLanesDir(root), "review.md") {
		t.Errorf("path = %v", payload["path"])
	}
}

// TestLanesUseUnknownIsUsageError asserts an unrecognised id exits 2 rather
// than crashing.
func TestLanesUseUnknownIsUsageError(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)

	_, err := runLanes(t, root, "use", "nope")
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
}

// TestLanesPublishWritesUnderIdentitySlug asserts 'lanes publish <id>' writes
// into .jaira/shared/<slug>/, stamping the acting identity as creator — the
// same call the lane settings screen's 'p' key makes.
func TestLanesPublishWritesUnderIdentitySlug(t *testing.T) {
	t.Setenv("JAIRA_USER", "Alex Doe")
	catalogue := lanesTestCatalogue(t)
	root := lanesTestProject(t)
	// A custom lane with no creator: line, since the built-in lanes default to
	// creator "jaira" (see parse()) and would never show the stamp.
	if err := os.WriteFile(filepath.Join(catalogue, "custom.md"), []byte(
		"---\nid: custom\nname: Custom\nafter: human\nprecedence: 41\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runLanes(t, root, "publish", "custom")
	if err != nil {
		t.Fatalf("lanes publish custom: %v\n%s", err, out)
	}
	dst := filepath.Join(root, ".jaira", "shared", "alex-doe", "custom.md")
	got, statErr := os.ReadFile(dst)
	if statErr != nil {
		t.Fatalf("expected %s to exist: %v", dst, statErr)
	}
	if !strings.Contains(string(got), "creator: alex-doe") {
		t.Errorf("published file missing a stamped creator, got:\n%s", got)
	}
}

// TestLanesPublishRefusesWithoutForce mirrors the 'use' overwrite refusal.
func TestLanesPublishRefusesWithoutForce(t *testing.T) {
	t.Setenv("JAIRA_USER", "alex")
	lanesTestCatalogue(t)
	root := lanesTestProject(t)

	if _, err := runLanes(t, root, "publish", "review"); err != nil {
		t.Fatalf("first publish: %v", err)
	}
	_, err := runLanes(t, root, "publish", "review")
	if err == nil {
		t.Fatal("expected a second 'publish' of the same lane to refuse")
	}
	if _, err := runLanes(t, root, "publish", "review", "--force"); err != nil {
		t.Fatalf("publish --force: %v", err)
	}
}

// TestLanesAdoptCopiesSharedLaneIntoCatalogue asserts 'lanes adopt <path>'
// takes the path 'lanes shared' prints and copies it into the catalogue.
func TestLanesAdoptCopiesSharedLaneIntoCatalogue(t *testing.T) {
	catalogue := lanesTestCatalogue(t)
	root := lanesTestProject(t)
	sharedDir := filepath.Join(root, ".jaira", "shared", "sam")
	if err := os.MkdirAll(sharedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(sharedDir, "hitl.md")
	if err := os.WriteFile(src, []byte(
		"---\nid: hitl\nname: HITL\nafter: human\nprecedence: 41\ncreator: sam\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runLanes(t, root, "adopt", src)
	if err != nil {
		t.Fatalf("lanes adopt: %v\n%s", err, out)
	}
	dst := filepath.Join(catalogue, "hitl.md")
	if _, statErr := os.Stat(dst); statErr != nil {
		t.Errorf("expected %s to exist: %v", dst, statErr)
	}
}

// TestLanesAdoptRefusesThenOverwritesWithForce mirrors the collision handling
// the lane settings screen gives, minus the two-press confirmation (the CLI
// has --force instead).
func TestLanesAdoptRefusesThenOverwritesWithForce(t *testing.T) {
	catalogue := lanesTestCatalogue(t)
	root := lanesTestProject(t)
	if err := os.WriteFile(filepath.Join(catalogue, "hitl.md"), []byte(
		"---\nid: hitl\nname: Existing\nafter: human\nprecedence: 41\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sharedDir := filepath.Join(root, ".jaira", "shared", "sam")
	if err := os.MkdirAll(sharedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(sharedDir, "hitl.md")
	if err := os.WriteFile(src, []byte(
		"---\nid: hitl\nname: Sam's HITL\nafter: human\nprecedence: 41\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := runLanes(t, root, "adopt", src)
	if err == nil {
		t.Fatal("expected adopt to refuse an existing catalogue id without --force")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitValidation {
		t.Errorf("code = %d, want %d", ce.code, ExitValidation)
	}

	if _, err := runLanes(t, root, "adopt", src, "--force"); err != nil {
		t.Fatalf("adopt --force: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(catalogue, "hitl.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "Sam's HITL") {
		t.Errorf("adopt --force did not overwrite the catalogue file, got:\n%s", got)
	}
}

// TestLanesAdoptUnknownPathIsErrorNotCrash asserts a path that does not exist
// fails cleanly.
func TestLanesAdoptUnknownPathIsErrorNotCrash(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)
	_, err := runLanes(t, root, "adopt", filepath.Join(root, "nope.md"))
	if err == nil {
		t.Fatal("expected an error for a nonexistent path")
	}
}

// TestLanesDefaultPrintsBuiltinsWhenAbsent asserts reading an absent default
// board prints the built-in state rather than failing.
func TestLanesDefaultPrintsBuiltinsWhenAbsent(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))

	out, err := runLanes(t, t.TempDir(), "default")
	if err != nil {
		t.Fatalf("lanes default: %v\n%s", err, out)
	}
	if !strings.Contains(out, "backlog") {
		t.Errorf("lanes default = %q, want it to list the built-in lanes", out)
	}

	jout, err := runLanes(t, t.TempDir(), "default", "--json")
	if err != nil {
		t.Fatalf("lanes default --json: %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	if payload["using_builtins"] != true {
		t.Errorf("using_builtins = %v, want true", payload["using_builtins"])
	}
}

// TestLanesDefaultRoundTrips asserts writing then reading the default board
// carries the selection back.
func TestLanesDefaultRoundTrips(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))
	dir := t.TempDir()

	out, err := runLanes(t, dir, "default", "--lanes", "backlog,todo", "--options", "brainstorm")
	if err != nil {
		t.Fatalf("lanes default --lanes: %v\n%s", err, out)
	}

	jout, err := runLanes(t, dir, "default", "--json")
	if err != nil {
		t.Fatalf("lanes default --json: %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	if payload["using_builtins"] != false {
		t.Errorf("using_builtins = %v, want false after setting --lanes", payload["using_builtins"])
	}
	lanesGot, _ := payload["lanes"].([]any)
	if len(lanesGot) != 2 || lanesGot[0] != "backlog" || lanesGot[1] != "todo" {
		t.Errorf("lanes = %v, want [backlog todo]", lanesGot)
	}
	optsGot, _ := payload["options"].([]any)
	if len(optsGot) != 1 || optsGot[0] != "brainstorm" {
		t.Errorf("options = %v, want [brainstorm]", optsGot)
	}
}

// TestLanesDefaultRefusesUnknownLane asserts an id nothing installed claims
// is refused rather than written.
func TestLanesDefaultRefusesUnknownLane(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))

	_, err := runLanes(t, t.TempDir(), "default", "--lanes", "nosuchlane")
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
}

// TestLanesDefaultRefusesUnknownOption asserts an option name nothing
// installed requires is refused rather than written.
func TestLanesDefaultRefusesUnknownOption(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "default-board.md"))

	_, err := runLanes(t, t.TempDir(), "default", "--options", "nosuchoption")
	if err == nil {
		t.Fatal("expected an error for an unknown option name")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitUsage {
		t.Errorf("code = %d, want %d", ce.code, ExitUsage)
	}
}

// TestLanesDefaultClearRemovesFile asserts --clear removes a written default
// board and is idempotent when there is nothing to remove.
func TestLanesDefaultClearRemovesFile(t *testing.T) {
	lanesTestCatalogue(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	t.Setenv("JAIRA_DEFAULT_BOARD", board)

	if _, err := runLanes(t, t.TempDir(), "default", "--lanes", "backlog"); err != nil {
		t.Fatalf("lanes default --lanes: %v", err)
	}
	if _, err := os.Stat(board); err != nil {
		t.Fatalf("setup: expected %s to exist: %v", board, err)
	}

	if _, err := runLanes(t, t.TempDir(), "default", "--clear"); err != nil {
		t.Fatalf("lanes default --clear: %v", err)
	}
	if _, err := os.Stat(board); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed, stat err = %v", board, err)
	}

	// Idempotent: clearing an already-absent default board is a no-op
	// success, not an error, matching the concurrency-safe-retry convention.
	if _, err := runLanes(t, t.TempDir(), "default", "--clear"); err != nil {
		t.Fatalf("lanes default --clear (already absent): %v", err)
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

// TestLanesPathNamesDefaultBoard asserts 'lanes path' names the default
// board file in both output modes, in and out of a project — task 3 only
// required the two directories.
func TestLanesPathNamesDefaultBoard(t *testing.T) {
	lanesTestCatalogue(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	t.Setenv("JAIRA_DEFAULT_BOARD", board)

	outside := t.TempDir()
	out, err := runLanes(t, outside, "path")
	if err != nil {
		t.Fatalf("lanes path: %v\n%s", err, out)
	}
	if !strings.Contains(out, board) {
		t.Errorf("lanes path outside a project = %q, want it to name the default board file", out)
	}
	if !strings.Contains(out, "not written yet") {
		t.Errorf("lanes path = %q, want it to say the default board has not been written", out)
	}

	jout, err := runLanes(t, outside, "path", "--json")
	if err != nil {
		t.Fatalf("lanes path --json: %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	if payload["default_board"] != board {
		t.Errorf("lanes path --json default_board = %v, want %q", payload["default_board"], board)
	}
	if payload["default_board_exists"] != false {
		t.Errorf("lanes path --json default_board_exists = %v, want false", payload["default_board_exists"])
	}

	if err := os.WriteFile(board, []byte("---\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jout2, err := runLanes(t, outside, "path", "--json")
	if err != nil {
		t.Fatalf("lanes path --json (board exists): %v\n%s", err, jout2)
	}
	var payload2 map[string]any
	if err := json.Unmarshal([]byte(jout2), &payload2); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout2)
	}
	if payload2["default_board_exists"] != true {
		t.Errorf("lanes path --json default_board_exists = %v, want true once written", payload2["default_board_exists"])
	}
}

// TestLanesTemplateBoardParses asserts 'lanes template --board' prints a
// file LoadDefaultBoard accepts.
func TestLanesTemplateBoardParses(t *testing.T) {
	lanesTestCatalogue(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	t.Setenv("JAIRA_DEFAULT_BOARD", board)

	out, err := runLanes(t, t.TempDir(), "template", "--board")
	if err != nil {
		t.Fatalf("lanes template --board: %v\n%s", err, out)
	}
	if err := os.WriteFile(board, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	db, err := lane.LoadDefaultBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Warnings) != 0 {
		t.Errorf("template --board output produced warnings: %v", db.Warnings)
	}
	if len(db.Lanes) == 0 {
		t.Errorf("expected the template to name at least one lane, got none")
	}
}

// TestLanesWarnsOnBadDefaultBoardLaneAndOption asserts a default board naming
// an uninstalled lane and an unclaimed option are both warned about, on
// stderr in human mode and in the warnings array under --json, and the
// command still exits 0.
func TestLanesWarnsOnBadDefaultBoardLaneAndOption(t *testing.T) {
	lanesTestCatalogue(t)
	board := filepath.Join(t.TempDir(), "default-board.md")
	if err := os.WriteFile(board, []byte("---\nlanes: [backlog, nosuchlane]\noptions: [nosuchoption]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_DEFAULT_BOARD", board)

	// The warning line goes straight to os.Stderr (see captureStdio), not
	// through the cobra command's own writers, so it has to be captured at
	// the file-descriptor level rather than via runLanes' buffer.
	dir := t.TempDir()
	var out string
	var err error
	_, stderr := captureStdio(t, func() {
		out, err = runLanes(t, dir)
	})
	if err != nil {
		t.Fatalf("lanes must exit 0 even with a bad default board: %v\n%s", err, out)
	}
	if !strings.Contains(stderr, "nosuchlane") {
		t.Errorf("lanes stderr = %q, want a warning naming nosuchlane", stderr)
	}
	if !strings.Contains(stderr, "nosuchoption") {
		t.Errorf("lanes stderr = %q, want a warning naming nosuchoption", stderr)
	}

	jout, err := runLanes(t, t.TempDir(), "--json")
	if err != nil {
		t.Fatalf("lanes --json must exit 0 even with a bad default board: %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	warnings, _ := payload["warnings"].([]any)
	found := map[string]bool{}
	for _, w := range warnings {
		s, _ := w.(string)
		if strings.Contains(s, "nosuchlane") {
			found["lane"] = true
		}
		if strings.Contains(s, "nosuchoption") {
			found["option"] = true
		}
	}
	if !found["lane"] || !found["option"] {
		t.Errorf("lanes --json warnings = %v, want both the lane and option warnings", warnings)
	}
}

// TestLanesSharedListsPublishedLanes covers the CLI half of criterion 8: an
// agent has no way to press the TUI's adopt key, so 'lanes shared' is the
// read path.
func TestLanesSharedListsPublishedLanes(t *testing.T) {
	lanesTestCatalogue(t)
	root := t.TempDir()
	s, err := ticket.At(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	sharedDir := filepath.Join(root, ".jaira", "shared", "sam")
	if err := os.MkdirAll(sharedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sharedDir, "hitl.md"), []byte(
		"---\nid: hitl\nname: HITL\nafter: human\nprecedence: 41\ncreator: sam\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runLanes(t, root, "shared")
	if err != nil {
		t.Fatalf("lanes shared: %v\n%s", err, out)
	}
	if !strings.Contains(out, "hitl") || !strings.Contains(out, "sam") {
		t.Errorf("lanes shared = %q, want it to name the lane, folder and creator", out)
	}

	jout, err := runLanes(t, root, "shared", "--json")
	if err != nil {
		t.Fatalf("lanes shared --json: %v\n%s", err, jout)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	arr, ok := payload["shared"].([]any)
	if !ok || len(arr) != 1 {
		t.Fatalf("lanes shared --json = %v, want exactly one shared lane", payload)
	}
	entry := arr[0].(map[string]any)
	if entry["id"] != "hitl" || entry["folder"] != "sam" || entry["creator"] != "sam" {
		t.Errorf("lanes shared --json entry = %v", entry)
	}
}

// TestLanesSharedOutsideProjectSaysSo asserts the command does not crash or
// misreport when there is no project to look inside.
func TestLanesSharedOutsideProjectSaysSo(t *testing.T) {
	lanesTestCatalogue(t)
	out, err := runLanes(t, t.TempDir(), "shared")
	if err != nil {
		t.Fatalf("lanes shared (outside a project): %v\n%s", err, out)
	}
	if !strings.Contains(out, "not in a project directory") {
		t.Errorf("lanes shared outside a project = %q, want it to say there is no project", out)
	}
}
