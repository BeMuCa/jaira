package lane

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeLane writes a minimal lane file and returns its path.
func writeLane(t *testing.T, dir, filename, body string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, filename)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func containsWarning(warnings []string, substr string) bool {
	for _, w := range warnings {
		if strings.Contains(w, substr) {
			return true
		}
	}
	return false
}

// TestLoadBuiltinsAlone asserts that with an empty catalogue and no project,
// Load("") returns exactly the shipped built-ins, in filename order.
func TestLoadBuiltinsAlone(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}

	want, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Lanes) != len(want) {
		t.Fatalf("got %d lanes, want %d built-ins: %v", len(set.Lanes), len(want), set.IDs())
	}
	for i, l := range want {
		if set.Lanes[i].ID != l.ID {
			t.Errorf("position %d: got %q, want %q (built-ins alone must render in shipped order)", i, set.Lanes[i].ID, l.ID)
		}
	}
	if len(set.Warnings) != 0 {
		t.Errorf("a clean install must produce zero warnings, got: %v", set.Warnings)
	}
}

// TestLoadCatalogue asserts a catalogue lane loads and orders by its anchor.
func TestLoadCatalogue(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)

	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
agentic: false
terminal: false
description: A custom step after HITL.
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := set.Get("triage"); !ok {
		t.Fatalf("catalogue lane did not load: %v", set.IDs())
	}
	humanIdx := set.Index("human")
	triageIdx := set.Index("triage")
	if triageIdx != humanIdx+1 {
		t.Errorf("triage should land immediately after human: ids=%v", set.IDs())
	}
}

// TestLoadProjectAuthoritative covers D-03: a non-empty project directory
// replaces the catalogue as this project's second lane source. A catalogue
// lane not copied into the project directory does not appear for that root,
// while it does appear for a root with no project directory at all.
func TestLoadProjectAuthoritative(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "catalogue-only.md", `---
id: catalogue-only
name: Catalogue Only
after: human
precedence: 41
---
`)

	root := t.TempDir()
	projDir := ProjectLanesDir(root)
	writeLane(t, projDir, "project-only.md", `---
id: project-only
name: Project Only
after: human
precedence: 41
---
`)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := set.Get("project-only"); !ok {
		t.Fatalf("project lane did not load: %v", set.IDs())
	}
	if _, ok := set.Get("catalogue-only"); ok {
		t.Fatalf("catalogue lane leaked into a project with its own lane directory: %v", set.IDs())
	}

	// A root with no project directory of its own still gets the catalogue.
	otherRoot := t.TempDir()
	otherSet, err := Load(otherRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := otherSet.Get("catalogue-only"); !ok {
		t.Fatalf("catalogue lane should appear for a root with no project directory: %v", otherSet.IDs())
	}
	if _, ok := otherSet.Get("project-only"); ok {
		t.Fatalf("another root's project lane must not leak: %v", otherSet.IDs())
	}
}

// TestLoadEmptyProjectDirWarns covers the D-03 mitigation: a project
// directory that exists but holds no lane files falls back to the catalogue
// (same as if it were absent) but must not do so silently.
func TestLoadEmptyProjectDirWarns(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "fallback.md", `---
id: fallback
name: Fallback
after: human
precedence: 41
---
`)

	root := t.TempDir()
	projDir := ProjectLanesDir(root)
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !containsWarning(set.Warnings, projDir) {
		t.Errorf("empty project directory must warn, naming the directory; got: %v", set.Warnings)
	}
	if _, ok := set.Get("fallback"); !ok {
		t.Errorf("an empty project directory must still fall back to the catalogue: %v", set.IDs())
	}
}

// TestLoadOverride covers the reversed collision rule: a custom lane whose id
// matches a built-in replaces it, prompt included, and the warning names the
// file, the id and the built-in it displaced.
func TestLoadOverride(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	src := writeLane(t, catalogue, "review.md", `---
id: review
name: My Review
after: human
precedence: 50
agentic: true
terminal: false
model-tier: strong
---
A completely different prompt.
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("review")
	if !ok {
		t.Fatalf("review lane missing after override: %v", set.IDs())
	}
	if l.Prompt != "A completely different prompt." {
		t.Errorf("override must replace the prompt, got %q", l.Prompt)
	}
	if l.Builtin {
		t.Errorf("an overriding lane must not be marked Builtin")
	}

	want := "lane " + src + ": id \"review\" overrides the built-in lane of the same name"
	if !containsList(set.Warnings, want) {
		t.Errorf("missing exact override warning %q, got: %v", want, set.Warnings)
	}
}

func containsList(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// TestLoadOverrideDropsProtection covers T-5-01: overriding a built-in that
// carries requires-human-exit is allowed, but dropping the field must produce
// a second, distinct warning naming it.
func TestLoadOverrideDropsProtection(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	src := writeLane(t, catalogue, "signoff.md", `---
id: signoff
name: Sign-off
after: review
precedence: 55
agentic: false
terminal: false
description: An override that silently drops the human-exit gate.
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("signoff")
	if !ok {
		t.Fatalf("signoff lane missing: %v", set.IDs())
	}
	if l.RequiresHumanExit {
		t.Fatalf("test lane fixture must not declare requires-human-exit")
	}

	want := "lane " + src + ": overriding \"signoff\" drops requires-human-exit — an agent could now accept its own work here undetected"
	if !containsList(set.Warnings, want) {
		t.Errorf("missing exact gate-weakening warning %q, got: %v", want, set.Warnings)
	}
	// The ordinary override warning must also still be present, as its own line.
	ordinary := "lane " + src + ": id \"signoff\" overrides the built-in lane of the same name"
	if !containsList(set.Warnings, ordinary) {
		t.Errorf("gate-weakening warning must be in addition to, not instead of, the ordinary override warning; got: %v", set.Warnings)
	}
}

// TestLoadOverrideDropsMultipleProtections covers a "done" override that
// drops terminal, requires-outcome and requires-nonmodel-signal at once: all
// three must be named in one warning.
func TestLoadOverrideDropsMultipleProtections(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	src := writeLane(t, catalogue, "done.md", `---
id: done
name: Done
after: signoff
precedence: 60
agentic: false
terminal: false
requires-outcome: false
requires-nonmodel-signal: false
description: An override that drops every protection Done carries.
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	want := "lane " + src + ": overriding \"done\" drops terminal, requires-outcome, requires-nonmodel-signal — an agent could now accept its own work here undetected"
	if !containsList(set.Warnings, want) {
		t.Errorf("missing exact multi-protection warning %q, got: %v", want, set.Warnings)
	}
}

// builtinBytes returns the shipped bytes of the named built-in lane, exactly
// what 'jaira lanes use' writes into a project's own lane directory.
func builtinBytes(t *testing.T, id string) []byte {
	t.Helper()
	builtins, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range builtins {
		if l.ID == id {
			raw, err := Bytes(l)
			if err != nil {
				t.Fatal(err)
			}
			return raw
		}
	}
	t.Fatalf("no built-in lane %q", id)
	return nil
}

// TestLoadOverrideEquivalentToBuiltinIsSilent covers the exact shape 'jaira
// lanes use' produces: a shadowing file whose bytes are the built-in's own,
// unchanged. That must warn about nothing and must not be marked Overrides —
// it is a copy, not an override in any sense the reader cares about.
func TestLoadOverrideEquivalentToBuiltinIsSilent(t *testing.T) {
	raw := builtinBytes(t, "review")

	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "review.md", string(raw))

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("review")
	if !ok {
		t.Fatalf("review lane missing: %v", set.IDs())
	}
	if l.Overrides != "" {
		t.Errorf("a byte-identical copy must not be marked as overriding, got %q", l.Overrides)
	}
	if containsWarning(set.Warnings, "overrides the built-in") {
		t.Errorf("a byte-identical copy must not warn about overriding anything: %v", set.Warnings)
	}
}

// TestLoadOverrideDifferingOnlyByCreatorIsSilent covers the field the
// comparison must deliberately ignore: creator is stamped by
// 'jaira lanes publish', so a copy differing only there is still a copy.
func TestLoadOverrideDifferingOnlyByCreatorIsSilent(t *testing.T) {
	raw := string(builtinBytes(t, "review"))
	withCreator := strings.Replace(raw, "---\n", "---\ncreator: alice\n", 1)
	if withCreator == raw {
		t.Fatal("test fixture did not actually add a creator line")
	}

	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "review.md", withCreator)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("review")
	if !ok {
		t.Fatalf("review lane missing: %v", set.IDs())
	}
	if l.Creator != "alice" {
		t.Fatalf("test fixture broken: creator not parsed, got %q", l.Creator)
	}
	if l.Overrides != "" {
		t.Errorf("a copy differing only by creator must not be marked as overriding, got %q", l.Overrides)
	}
	if containsWarning(set.Warnings, "overrides the built-in") {
		t.Errorf("a copy differing only by creator must not warn about overriding anything: %v", set.Warnings)
	}
}

// TestLoadOverrideChangedPromptStillWarns covers the other side: a shadowing
// file identical to the built-in except for its prompt is a real behaviour
// change and must warn exactly as before.
func TestLoadOverrideChangedPromptStillWarns(t *testing.T) {
	raw := string(builtinBytes(t, "review"))
	const oldPrompt = "Review this change against its definition of done."
	const newPrompt = "Review this change and reject it if you have any doubt at all."
	changed := strings.Replace(raw, oldPrompt, newPrompt, 1)
	if changed == raw {
		t.Fatal("test fixture did not actually change the prompt")
	}

	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	src := writeLane(t, catalogue, "review.md", changed)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("review")
	if !ok {
		t.Fatalf("review lane missing: %v", set.IDs())
	}
	if l.Overrides != "review" {
		t.Errorf("a changed prompt must still be marked as overriding, got %q", l.Overrides)
	}
	want := "lane " + src + ": id \"review\" overrides the built-in lane of the same name"
	if !containsList(set.Warnings, want) {
		t.Errorf("missing exact override warning %q, got: %v", want, set.Warnings)
	}
}

// TestLoadDuplicateInDirectory asserts that two files in one directory
// claiming the same non-built-in id still resolve deterministically to the
// first by sorted filename, with a warning — unchanged from before this task.
func TestLoadDuplicateInDirectory(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "a-first.md", `---
id: dup
name: First
after: human
precedence: 41
---
`)
	writeLane(t, catalogue, "b-second.md", `---
id: dup
name: Second
after: human
precedence: 41
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("dup")
	if !ok {
		t.Fatalf("dup lane missing: %v", set.IDs())
	}
	if l.Name != "First" {
		t.Errorf("duplicate id must resolve to the first by sorted filename, got %q", l.Name)
	}
	if !containsWarning(set.Warnings, "already defined by") {
		t.Errorf("duplicate id must warn, got: %v", set.Warnings)
	}
}

// TestLoadWarnsWhenNoRequiresSpecified covers the D-03 consequence: if an
// override takes requires-specified out of circulation entirely, the loader
// must say so, because that board can never move a ticket into work.
func TestLoadWarnsWhenNoRequiresSpecified(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "todo.md", `---
id: todo
name: Todo
after: backlog
precedence: 20
agentic: false
terminal: false
description: An override that drops requires-specified.
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if !containsWarning(set.Warnings, "requires-specified") {
		t.Errorf("dropping the only requires-specified lane must warn, got: %v", set.Warnings)
	}
}

// TestCreatorParsedFromFrontmatter asserts creator: lands on Lane.Creator.
func TestCreatorParsedFromFrontmatter(t *testing.T) {
	l, err := parse([]byte(`---
id: attributed
name: Attributed
creator: alex
---
`), "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if l.Creator != "alex" {
		t.Errorf("Creator = %q, want %q", l.Creator, "alex")
	}
}

// TestBuiltinDefaultsCreatorToJaira asserts a built-in with no creator: field
// reports "jaira" — the nine shipped files carry no such line, so this is the
// default doing the work rather than the files.
func TestBuiltinDefaultsCreatorToJaira(t *testing.T) {
	lanes, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lanes {
		if l.Creator != "jaira" {
			t.Errorf("built-in %q: Creator = %q, want %q", l.ID, l.Creator, "jaira")
		}
	}
}

// TestCustomLaneCreatorDefaultsEmpty asserts a custom lane with no creator:
// field reports empty, not "jaira" — absent provenance and "shipped by the
// tool" are different facts.
func TestCustomLaneCreatorDefaultsEmpty(t *testing.T) {
	l, err := parse([]byte(`---
id: mine
name: Mine
---
`), "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if l.Creator != "" {
		t.Errorf("Creator = %q, want empty for a custom lane with no creator: field", l.Creator)
	}
}

// TestLoadUnknownAnchorWarnsAndOrders covers the shared-lane-file case: a
// lane anchored to an id this installation does not have still loads, lands
// before the terminal lane, and warns naming the missing anchor.
func TestLoadUnknownAnchorWarnsAndOrders(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "orphan.md", `---
id: orphan
name: Orphan
after: not-installed
precedence: 999
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := set.Get("orphan"); !ok {
		t.Fatalf("orphan lane must still load: %v", set.IDs())
	}
	terminal := set.Terminal()
	if terminal == nil {
		t.Fatal("no terminal lane found")
	}
	if set.Index("orphan") >= set.Index(terminal.ID) {
		t.Errorf("a lane anchored to a missing lane must land before the terminal lane: ids=%v", set.IDs())
	}
	if !containsWarning(set.Warnings, `anchor "not-installed" is not installed`) {
		t.Errorf("missing-anchor warning must name the anchor, got: %v", set.Warnings)
	}
}

// TestOrderChainResolvesRegardlessOfFilenameOrder covers a chain of two custom
// lanes — B anchored to A, A anchored to a built-in — where the dependent
// lane's file sorts before its dependency's. This is what the fixed-point
// loop in order() is for: a single pass would place chained-b before
// chained-a even existed.
func TestOrderChainResolvesRegardlessOfFilenameOrder(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	// "a-" sorts before "b-", so chained-b (the dependent) is read first.
	writeLane(t, catalogue, "a-dependent.md", `---
id: chained-b
name: Chained B
after: chained-a
precedence: 42
---
`)
	writeLane(t, catalogue, "b-dependency.md", `---
id: chained-a
name: Chained A
after: human
precedence: 41
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	humanIdx := set.Index("human")
	aIdx := set.Index("chained-a")
	bIdx := set.Index("chained-b")
	if aIdx != humanIdx+1 {
		t.Fatalf("chained-a should land immediately after human: ids=%v", set.IDs())
	}
	if bIdx != aIdx+1 {
		t.Fatalf("chained-b should land immediately after chained-a: ids=%v", set.IDs())
	}
	if len(set.Warnings) != 0 {
		t.Errorf("a resolvable chain must not warn, got: %v", set.Warnings)
	}
}

// TestOrderCycleWarnsBothAppear covers two lanes anchored to each other: the
// cycle is reported, but a bad file must never remove a column, so both
// lanes still appear on the board.
func TestOrderCycleWarnsBothAppear(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "cycle1.md", `---
id: cycle1
name: Cycle 1
after: cycle2
precedence: 41
---
`)
	writeLane(t, catalogue, "cycle2.md", `---
id: cycle2
name: Cycle 2
after: cycle1
precedence: 42
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := set.Get("cycle1"); !ok {
		t.Errorf("cycle1 must still appear: %v", set.IDs())
	}
	if _, ok := set.Get("cycle2"); !ok {
		t.Errorf("cycle2 must still appear: %v", set.IDs())
	}
	if !containsWarning(set.Warnings, "forms a cycle") {
		t.Errorf("a cyclic anchor must warn, got: %v", set.Warnings)
	}
}

// TestOrderNoAnchorLandsBeforeTerminal covers a lane with no after: at all:
// it must park before the terminal lane, the same place an unresolvable
// anchor lands.
func TestOrderNoAnchorLandsBeforeTerminal(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "unanchored.md", `---
id: unanchored
name: Unanchored
precedence: 41
---
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	terminal := set.Terminal()
	if terminal == nil {
		t.Fatal("no terminal lane found")
	}
	if set.Index("unanchored") >= set.Index(terminal.ID) {
		t.Errorf("a lane with no anchor must land before the terminal lane: ids=%v", set.IDs())
	}
}

// TestAgenticLaneRequiresModelTier asserts parse rejects an agentic lane with
// no model-tier — already true, this proves it rather than asserting it in a
// summary.
func TestAgenticLaneRequiresModelTier(t *testing.T) {
	_, err := parse([]byte(`---
id: agentic-no-tier
name: Agentic No Tier
agentic: true
---
a prompt
`), "test", false)
	if err == nil {
		t.Fatal("expected an error for an agentic lane with no model-tier")
	}
	if !strings.Contains(err.Error(), "model-tier") {
		t.Errorf("error must name model-tier, got: %v", err)
	}
}

// TestModelTierRoundTripsExact asserts a lane's ModelTier is the exact alias
// string given, not resolved, rewritten or validated against a fixed list —
// a user is free to invent their own alias for their own runner.
func TestModelTierRoundTripsExact(t *testing.T) {
	l, err := parse([]byte(`---
id: custom-tier
name: Custom Tier
agentic: true
model-tier: my-local-nano
---
a prompt
`), "test", false)
	if err != nil {
		t.Fatal(err)
	}
	if l.ModelTier != "my-local-nano" {
		t.Errorf("ModelTier = %q, want the alias unchanged: %q", l.ModelTier, "my-local-nano")
	}
}

// TestModelTierNeverComparedToModelName is a guard test: it reads the
// non-test source tree rather than asserting an absence from memory. A
// shared lane file surviving a model rename depends entirely on nothing in
// this repo knowing what a model is called.
func TestModelTierNeverComparedToModelName(t *testing.T) {
	banned := []string{"claude", "gpt", "sonnet", "opus", "haiku"}
	err := filepath.WalkDir("../..", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "builtin":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		b, readErr := os.ReadFile(p)
		if readErr != nil {
			return readErr
		}
		for i, line := range strings.Split(string(b), "\n") {
			if !strings.Contains(line, "ModelTier") {
				continue
			}
			lower := strings.ToLower(line)
			for _, name := range banned {
				if strings.Contains(lower, name) {
					t.Errorf("%s:%d: ModelTier compared against a model name %q: %s", p, i+1, name, line)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestColumnsPassthroughForUnknownStatus asserts Set.Columns appends a
// read-only column for a status no lane claims.
func TestColumnsPassthroughForUnknownStatus(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	cols := set.Columns([]string{"backlog", "some-ghost-lane"})
	found := false
	for _, c := range cols {
		if c.ID == "some-ghost-lane" {
			found = true
			if !c.Unknown {
				t.Errorf("passthrough column must be marked Unknown")
			}
		}
	}
	if !found {
		t.Errorf("Columns must append a passthrough column for an unrecognized status: %v", idsOf(cols))
	}
}

// TestColumnsMultipleUnknownSorted asserts multiple unknown statuses are
// appended in a deterministic (alphabetical) order.
func TestColumnsMultipleUnknownSorted(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	cols := set.Columns([]string{"zebra-lane", "alpha-lane", "backlog"})
	var unknownIDs []string
	for _, c := range cols {
		if c.Unknown {
			unknownIDs = append(unknownIDs, c.ID)
		}
	}
	want := []string{"alpha-lane", "zebra-lane"}
	if len(unknownIDs) != len(want) {
		t.Fatalf("unknown columns = %v, want %v", unknownIDs, want)
	}
	for i := range want {
		if unknownIDs[i] != want[i] {
			t.Errorf("unknown columns = %v, want %v", unknownIDs, want)
		}
	}
}

func idsOf(lanes []*Lane) []string {
	out := make([]string, 0, len(lanes))
	for _, l := range lanes {
		out = append(out, l.ID)
	}
	return out
}

// TestBuiltinOrderIsUnchanged is the task 10 regression baseline: the ten
// built-ins render in exactly the order the shipped binary prints them in
// today, captured by running the binary rather than copied from a plan that
// may be stale. Never delete this test — task 10's whole point is that
// nothing here moved.
func TestBuiltinOrderIsUnchanged(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"backlog", "brainstorm", "todo", "pre-process", "in-progress",
		"human", "review", "signoff", "done", "blocked",
	}
	got := set.IDs()
	if len(got) != len(want) {
		t.Fatalf("got %d lanes, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("position %d: got %q, want %q (ids=%v)", i, got[i], want[i], got)
		}
	}
}

// TestBuiltinsProduceZeroWarnings is the bar the contract check must not
// drop: a clean install's ten built-ins never warn, including the new
// input-requires/output-produces check.
func TestBuiltinsProduceZeroWarnings(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Warnings) != 0 {
		t.Errorf("a clean install must produce zero warnings, got: %v", set.Warnings)
	}
}

// TestContractCheckWarnsWhenRequiredBeforeProducer asserts a custom lane
// requiring "plan", ordered before pre-process (which produces it), warns
// naming the field and both lanes.
func TestContractCheckWarnsWhenRequiredBeforeProducer(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "early.md", `---
id: early-consumer
name: Early Consumer
after: backlog
precedence: 6
agentic: true
model-tier: cheap
input-requires: [plan]
---
needs a plan
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if set.Index("early-consumer") >= set.Index("pre-process") {
		t.Fatalf("test fixture must land before pre-process: ids=%v", set.IDs())
	}
	if !containsWarning(set.Warnings, `lane early-consumer: requires "plan", but pre-process (which produces it) is ordered after it`) {
		t.Errorf("expected a warning naming the field and both lanes, got: %v", set.Warnings)
	}
}

// TestContractCheckSilentWhenRequiredAfterProducer asserts the same lane,
// anchored after pre-process instead, produces no warning.
func TestContractCheckSilentWhenRequiredAfterProducer(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "late.md", `---
id: late-consumer
name: Late Consumer
after: pre-process
precedence: 26
agentic: true
model-tier: cheap
input-requires: [plan]
---
needs a plan
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if containsWarning(set.Warnings, "late-consumer") {
		t.Errorf("a lane ordered after its producer must not warn, got: %v", set.Warnings)
	}
}

// TestContractCheckSilentForTicketSuppliedField asserts a lane requiring
// "goal" — which the ticket itself can already carry from creation — never
// warns, even positioned before anything that produces it.
func TestContractCheckSilentForTicketSuppliedField(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "wants-goal.md", `---
id: wants-goal
name: Wants Goal
after: backlog
precedence: 1
agentic: true
model-tier: cheap
input-requires: [goal]
---
needs a goal
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if containsWarning(set.Warnings, "wants-goal") {
		t.Errorf("a field the ticket itself supplies must never warn, got: %v", set.Warnings)
	}
}

// TestContractCheckWarnsWhenNothingProducesTheField asserts a lane requiring
// a field no installed lane produces at all warns, rather than being
// silently accepted.
func TestContractCheckWarnsWhenNothingProducesTheField(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "wants-nothing.md", `---
id: wants-mystery
name: Wants Mystery
after: backlog
precedence: 2
agentic: true
model-tier: cheap
input-requires: [some-field-nothing-produces]
---
needs the impossible
`)

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	want := `lane wants-mystery: requires "some-field-nothing-produces", which no installed lane produces`
	if !containsList(set.Warnings, want) {
		t.Errorf("missing exact warning %q, got: %v", want, set.Warnings)
	}
}

// TestPrecedenceUnknownReturnsNegOne asserts a merge never promotes a ticket
// into a lane this installation cannot reason about.
func TestPrecedenceUnknownReturnsNegOne(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if p := set.Precedence("no-such-lane"); p != -1 {
		t.Errorf("Precedence(unknown) = %d, want -1", p)
	}
}
