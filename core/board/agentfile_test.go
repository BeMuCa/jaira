package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, root, name string) (string, bool) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, name))
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(b), true
}

func hasBlock(t *testing.T, root, name string) bool {
	s, _ := readFile(t, root, name)
	return strings.Contains(s, jairaMarkerStart)
}

// TestChooseAgentsDeletesAClaudeFileHoldingOnlyTheBlock: an empty CLAUDE.md
// would still hide AGENTS.md from Claude Code, so it goes.
func TestChooseAgentsDeletesAClaudeFileHoldingOnlyTheBlock(t *testing.T) {
	root := t.TempDir()
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ChooseAgentFiles(root, AgentFileAgents, nil); err != nil {
		t.Fatal(err)
	}
	if !hasBlock(t, root, "AGENTS.md") {
		t.Error("AGENTS.md lost the block")
	}
	if _, ok := readFile(t, root, "CLAUDE.md"); ok {
		t.Error("CLAUDE.md holding only the block should be deleted")
	}
}

// TestChooseAgentsImportsIntoAClaudeFileWithOwnContent: with a CLAUDE.md
// present, Claude Code reads AGENTS.md only through an import.
func TestChooseAgentsImportsIntoAClaudeFileWithOwnContent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "CLAUDE.md", "# Project\n\nOwn rules.\n")
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ChooseAgentFiles(root, AgentFileAgents, nil); err != nil {
		t.Fatal(err)
	}
	got, _ := readFile(t, root, "CLAUDE.md")
	want := claudeImport + "\n# Project\n\nOwn rules.\n"
	if got != want {
		t.Errorf("CLAUDE.md =\n%q\nwant\n%q", got, want)
	}
	// Going back to both takes jaira's import out again, or Claude Code
	// would read the block twice.
	if _, err := ChooseAgentFiles(root, AgentFileBoth, nil); err != nil {
		t.Fatal(err)
	}
	got, _ = readFile(t, root, "CLAUDE.md")
	if strings.Contains(got, "@AGENTS.md") || !strings.HasPrefix(got, "# Project\n\nOwn rules.\n") {
		t.Errorf("after both, CLAUDE.md =\n%s", got)
	}
}

// TestChooseAgentsLeavesTheUsersOwnImportAlone: an import somebody wrote is
// not jaira's to add twice or to take away.
func TestChooseAgentsLeavesTheUsersOwnImportAlone(t *testing.T) {
	root := t.TempDir()
	own := "@AGENTS.md\n\nOwn rules.\n"
	writeFile(t, root, "CLAUDE.md", own)
	if _, err := ChooseAgentFiles(root, AgentFileAgents, nil); err != nil {
		t.Fatal(err)
	}
	if got, _ := readFile(t, root, "CLAUDE.md"); got != own {
		t.Errorf("CLAUDE.md =\n%q\nwant it untouched", got)
	}
}

// TestChooseClaudeRemovesTheBlockFromAgentsAndKeepsTheRest.
func TestChooseClaudeRemovesTheBlockFromAgentsAndKeepsTheRest(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "AGENTS.md", "# Codex notes\n")
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ChooseAgentFiles(root, AgentFileClaude, nil); err != nil {
		t.Fatal(err)
	}
	if got, _ := readFile(t, root, "AGENTS.md"); got != "# Codex notes\n" {
		t.Errorf("AGENTS.md = %q, want the user's text back exactly", got)
	}
	if !hasBlock(t, root, "CLAUDE.md") {
		t.Error("CLAUDE.md lost the block")
	}
}

// TestTheChoiceHoldsForTheNextWrite: a later write without a choice — 'jaira
// update' with no flag, a lane change — keeps to the files carrying the block.
func TestTheChoiceHoldsForTheNextWrite(t *testing.T) {
	for _, tc := range []struct{ choice, in, out string }{
		{AgentFileAgents, "AGENTS.md", "CLAUDE.md"},
		{AgentFileClaude, "CLAUDE.md", "AGENTS.md"},
	} {
		root := t.TempDir()
		if _, err := ChooseAgentFiles(root, tc.choice, nil); err != nil {
			t.Fatal(err)
		}
		if _, err := AnnounceInAgentFiles(root, []LaneFact{{Name: "todo"}}); err != nil {
			t.Fatal(err)
		}
		if !hasBlock(t, root, tc.in) {
			t.Errorf("%s: %s lost the block", tc.choice, tc.in)
		}
		if hasBlock(t, root, tc.out) {
			t.Errorf("%s: %s got the block back", tc.choice, tc.out)
		}
	}
}

// TestNoBlockAnywhereWritesBoth is the old behaviour, for every board that
// never chose.
func TestNoBlockAnywhereWritesBoth(t *testing.T) {
	root := t.TempDir()
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	if !hasBlock(t, root, "AGENTS.md") || !hasBlock(t, root, "CLAUDE.md") {
		t.Error("a board that never chose should get the block in both files")
	}
}

func TestChooseRefusesAnUnknownChoice(t *testing.T) {
	if _, err := ChooseAgentFiles(t.TempDir(), "cursor", nil); err == nil {
		t.Error("want an error for an unknown choice")
	}
}

// TestChooseRefusesToDropAHandWrittenLocalArea: the local area inside the
// block is the user's own text, and removing the block would delete it.
func TestChooseRefusesToDropAHandWrittenLocalArea(t *testing.T) {
	root := t.TempDir()
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	s, _ := readFile(t, root, "CLAUDE.md")
	s = strings.Replace(s, jairaMarkerEnd, jairaMarkerLocal+"\nkeep me\n"+jairaMarkerEnd, 1)
	writeFile(t, root, "CLAUDE.md", s)
	if _, err := ChooseAgentFiles(root, AgentFileAgents, nil); err == nil {
		t.Fatal("want a refusal")
	}
	if got, _ := readFile(t, root, "CLAUDE.md"); !strings.Contains(got, "keep me") {
		t.Error("the local area was lost")
	}
}
