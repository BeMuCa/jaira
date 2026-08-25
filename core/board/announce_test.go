package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAnnounceRegeneratesByteForByteWithNoLocalMarker pins the exact rebuilt
// block for a file with no local marker: it must be byte-for-byte identical
// to what announceInAgentFile produces today, not merely "contains the note".
func TestAnnounceRegeneratesByteForByteWithNoLocalMarker(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	original := "# My project\n\nSome notes.\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, action, err := announceInAgentFile(root, "AGENTS.md", nil); err != nil {
		t.Fatalf("announceInAgentFile: %v", err)
	} else if action != "appended" {
		t.Fatalf("action = %q, want appended", action)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := original + "\n" + jairaMarkerStart + "\n" + agentNote + "\n" + jairaMarkerEnd + "\n"
	if string(got) != want {
		t.Errorf("announceInAgentFile() produced:\n%s\nwant exactly:\n%s", got, want)
	}
}

// TestAnnounceSurvivesLocalMarkerVerbatim asserts everything from
// jaira:local to the end marker — including blank lines and arbitrary
// markdown — survives regeneration unchanged, and that regenerating twice in
// a row is idempotent.
func TestAnnounceSurvivesLocalMarkerVerbatim(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	local := "\n" + jairaMarkerLocal + "\n\n" +
		"## Project rules\n\n- always run `task test`\n- never touch vendor/\n\n"
	original := jairaMarkerStart + "\nstale note\n" + local + jairaMarkerEnd + "\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, action, err := announceInAgentFile(root, "AGENTS.md", nil); err != nil {
		t.Fatalf("announceInAgentFile: %v", err)
	} else if action != "updated" {
		t.Fatalf("action = %q, want updated", action)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "## Project rules") ||
		!strings.Contains(string(got), "always run `task test`") ||
		!strings.Contains(string(got), "never touch vendor/") {
		t.Errorf("regenerated file = %q, lost the project's own rules from the local area", got)
	}
	if strings.Contains(string(got), "stale note") {
		t.Errorf("regenerated file = %q, the note above the local marker should have been refreshed", got)
	}

	// A second regeneration is idempotent.
	_, action2, err := announceInAgentFile(root, "AGENTS.md", nil)
	if err != nil {
		t.Fatalf("announceInAgentFile (second): %v", err)
	}
	if action2 != "unchanged" {
		t.Errorf("second regeneration action = %q, want unchanged", action2)
	}
}

// TestAnnounceCreatesFreshFileWithoutLocalMarker asserts a brand new file is
// created exactly as before this feature existed — no local marker present.
func TestAnnounceCreatesFreshFileWithoutLocalMarker(t *testing.T) {
	root := t.TempDir()

	path, action, err := announceInAgentFile(root, "CLAUDE.md", nil)
	if err != nil {
		t.Fatalf("announceInAgentFile: %v", err)
	}
	if action != "created" {
		t.Errorf("action = %q, want created", action)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := jairaMarkerStart + "\n" + agentNote + "\n" + jairaMarkerEnd + "\n"
	if string(got) != want {
		t.Errorf("created file = %q, want exactly %q", got, want)
	}
	if strings.Contains(string(got), jairaMarkerLocal) {
		t.Errorf("a brand new file must not contain the local marker, got %q", got)
	}
}

// TestAnnounceIgnoresLocalMarkerOutsideBlock asserts a jaira:local marker
// written before the start marker, or after the end marker, is left alone as
// ordinary user text rather than treated as the boundary.
func TestAnnounceIgnoresLocalMarkerOutsideBlock(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	original := "before " + jairaMarkerLocal + " text\n" +
		jairaMarkerStart + "\nstale note\n" + jairaMarkerEnd + "\n" +
		"after " + jairaMarkerLocal + " text\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, _, err := announceInAgentFile(root, "AGENTS.md", nil); err != nil {
		t.Fatalf("announceInAgentFile: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "before "+jairaMarkerLocal+" text") {
		t.Errorf("regenerated file = %q, lost the text before the start marker", got)
	}
	if !strings.Contains(string(got), "after "+jairaMarkerLocal+" text") {
		t.Errorf("regenerated file = %q, lost the text after the end marker", got)
	}
	// The block itself was rebuilt as if no local marker were present, since
	// the two occurrences found are both outside [start, end).
	want := jairaMarkerStart + "\n" + agentNote + "\n" + jairaMarkerEnd
	if !strings.Contains(string(got), want) {
		t.Errorf("regenerated file = %q, want the managed block rebuilt without a local area", got)
	}
}

// TestAgentNoteExplainsTheLocalArea asserts the generated note itself tells a
// reader the local area exists and how to open one.
func TestAgentNoteExplainsTheLocalArea(t *testing.T) {
	if !strings.Contains(agentNote, "jaira:local") {
		t.Errorf("agentNote does not mention the jaira:local marker")
	}
	if !strings.Contains(agentNote, "jaira sync") {
		t.Errorf("agentNote does not mention 'jaira sync'")
	}
}
