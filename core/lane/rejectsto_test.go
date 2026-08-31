package lane

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// catalogueWith puts lane files in a temp catalogue and loads the set from it.
func catalogueWith(t *testing.T, files map[string]string) *Set {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", dir)
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	s, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// A lane that sends work back declares where to. Without it the loop lives only
// in the prompt's prose, and an agent reading the board cannot see it at all.
func TestRejectsToIsParsed(t *testing.T) {
	s := catalogueWith(t, map[string]string{"critique.md": `---
id: critique
name: Critique
after: in-progress
rejects-to: in-progress
agentic: true
model-tier: strong
description: Judges whether this is the right implementation.
---
# Prompt

Say what is wrong.
`})
	l, ok := s.Get("critique")
	if !ok {
		t.Fatal("lane did not load")
	}
	if !slices.Equal(l.RejectsTo, []string{"in-progress"}) {
		t.Errorf("RejectsTo = %q, want [in-progress]", l.RejectsTo)
	}
}

// A back-edge to a lane nobody installed is a dead end, and the board says so
// rather than letting an agent walk into it.
func TestRejectsToMustResolve(t *testing.T) {
	s := catalogueWith(t, map[string]string{"critique.md": `---
id: critique
name: Critique
after: in-progress
rejects-to: nowhere-at-all
description: x
---
`})
	var found bool
	for _, w := range s.Warnings {
		if strings.Contains(w, "rejects-to") && strings.Contains(w, "nowhere-at-all") {
			found = true
		}
	}
	if !found {
		t.Errorf("a dead back-edge was not reported:\n%s", strings.Join(s.Warnings, "\n"))
	}
}

// A lane may not send work back to itself: that is not a loop, it is a stall.
func TestRejectsToRefusesItself(t *testing.T) {
	s := catalogueWith(t, map[string]string{"critique.md": `---
id: critique
name: Critique
after: in-progress
rejects-to: critique
description: x
---
`})
	var found bool
	for _, w := range s.Warnings {
		if strings.Contains(w, "itself") {
			found = true
		}
	}
	if !found {
		t.Errorf("a self back-edge was not reported:\n%s", strings.Join(s.Warnings, "\n"))
	}
}

// Absent is the normal state: most lanes never send anything back.
func TestRejectsToDefaultsEmpty(t *testing.T) {
	bs, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range bs {
		if len(l.RejectsTo) > 0 {
			t.Errorf("built-in %q declares rejects-to %q; no built-in should", l.ID, l.RejectsTo)
		}
	}
}

// A project copy of a built-in must not start warning just because the field
// exists now: both sides are empty, so they are still equivalent.
func TestRejectsToDoesNotMakeCopiesDrift(t *testing.T) {
	bs, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range bs {
		copyOf := *l
		if !lanesEquivalent(l, &copyOf) {
			t.Errorf("an unmodified copy of %q now reads as an override", l.ID)
		}
	}
}

// A changed back-edge is a changed lane, so an override that moves the loop is
// reported like any other behaviour change.
func TestRejectsToCountsAsDrift(t *testing.T) {
	base := &Lane{ID: "critique", Name: "Critique"}
	moved := &Lane{ID: "critique", Name: "Critique", RejectsTo: []string{"todo"}}
	if lanesEquivalent(base, moved) {
		t.Error("moving the back-edge did not register as a change")
	}
}

// A rejection is not one kind of thing: a critique that found a flaw sends the
// work back to be implemented, one that found a decision hands it to a person.
// Both are the same lane's back edge, so a lane may name both.
func TestRejectsToAcceptsAList(t *testing.T) {
	s := catalogueWith(t, map[string]string{
		"in-progress.md": laneFile("in-progress", ""),
		"human.md":       laneFile("human", ""),
		"critique.md":    laneFile("critique", "rejects-to: [in-progress, human]"),
	})
	l, ok := s.Get("critique")
	if !ok {
		t.Fatal("lane did not load")
	}
	if !slices.Equal(l.RejectsTo, []string{"in-progress", "human"}) {
		t.Errorf("RejectsTo = %q, want [in-progress human]", l.RejectsTo)
	}
	for _, w := range s.Warnings {
		if strings.Contains(w, "rejects-to") {
			t.Errorf("a good back edge warned: %s", w)
		}
	}
}

// A list is validated target by target, so a lane that declares one good back
// edge and two bad ones is told which of them is bad rather than that something
// is.
func TestRejectsToWarnsPerTarget(t *testing.T) {
	s := catalogueWith(t, map[string]string{
		"in-progress.md": laneFile("in-progress", ""),
		"critique.md":    laneFile("critique", "rejects-to: [in-progress, nowhere-at-all, critique]"),
	})
	var dead, self bool
	for _, w := range s.Warnings {
		if strings.Contains(w, "rejects-to") && strings.Contains(w, "nowhere-at-all") {
			dead = true
		}
		if strings.Contains(w, "itself") {
			self = true
		}
	}
	if !dead || !self {
		t.Errorf("a list's bad targets were not both reported (dead=%v self=%v):\n%s",
			dead, self, strings.Join(s.Warnings, "\n"))
	}
	for _, w := range s.Warnings {
		if strings.Contains(w, "rejects-to") && strings.Contains(w, "in-progress") {
			t.Errorf("the good target was reported too: %s", w)
		}
	}
}

// laneFile writes a minimal valid lane, with one extra frontmatter line to test.
func laneFile(id, extra string) string {
	return "---\nid: " + id + "\nname: " + id + "\ndescription: x\n" + extra + "\n---\n"
}
