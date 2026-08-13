package lane

import (
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
