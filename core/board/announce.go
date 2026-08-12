package board

import (
	"os"
	"path/filepath"
	"strings"
)

// jairaMarker bounds the block jaira manages, so the note can be updated or
// removed later without touching anything the user wrote around it.
const (
	jairaMarkerStart = "<!-- jaira:start -->"
	jairaMarkerEnd   = "<!-- jaira:end -->"
)

// agentNote is what an agent reads on entering the repository.
//
// The skill's own description already says to use it in a repository with a
// .jaira directory, but that relies on the model noticing. A line in the file
// the agent is given at the start of every session does not.
const agentNote = `## Task tracking: jaira

This repository has a jaira board (` + "`.jaira/`" + `). Multi-step work is tracked
there as markdown tickets so it survives session boundaries.

- ` + "`jaira next --json`" + ` — the next actionable ticket
- ` + "`jaira create <title>`" + ` — capture work; only a title is needed
- ` + "`jaira dod <id> <n> --doing|--done`" + ` — mark checklist items as you go
- ` + "`jaira show <id> --for-lane <lane> --json`" + ` — the prompt and bounded input for a step

Do not edit files under ` + "`.jaira/tickets/`" + ` directly; the CLI is the write path.
The sign-off lane cannot be left by an agent — a person accepts the work there.`

// agentFiles are the instruction files coding agents read.
//
// There is no single convention. AGENTS.md is the closest thing to a cross-tool
// standard and is what Codex and several others read; CLAUDE.md is Claude Code's.
// Both are written, because a board nobody's agent knows about is a board that
// does not get used, and the cost of an extra markdown section is nothing.
var agentFiles = []string{"AGENTS.md", "CLAUDE.md"}

// AnnounceInAgentFiles writes the note into each agent instruction file.
func AnnounceInAgentFiles(root string) (written []string, err error) {
	for _, name := range agentFiles {
		path, action, ferr := announceInAgentFile(root, name)
		if ferr != nil {
			return written, ferr
		}
		if action != "unchanged" {
			written = append(written, filepath.Base(path)+" ("+action+")")
		}
	}
	return written, nil
}

// announceInAgentFile writes the note into one file, creating it if absent.
//
// It reports what it did rather than doing it quietly, because editing a file
// the user wrote is not something a tool should do invisibly.
func announceInAgentFile(root, name string) (path string, action string, err error) {
	path = filepath.Join(root, name)
	block := jairaMarkerStart + "\n" + agentNote + "\n" + jairaMarkerEnd + "\n"

	existing, readErr := os.ReadFile(path)
	if readErr != nil {
		if !os.IsNotExist(readErr) {
			return path, "", readErr
		}
		if err := os.WriteFile(path, []byte(block), 0o644); err != nil {
			return path, "", err
		}
		return path, "created", nil
	}

	s := string(existing)
	start := strings.Index(s, jairaMarkerStart)
	end := strings.Index(s, jairaMarkerEnd)
	if start >= 0 && end > start {
		updated := s[:start] + strings.TrimSuffix(block, "\n") + s[end+len(jairaMarkerEnd):]
		if updated == s {
			return path, "unchanged", nil
		}
		if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
			return path, "", err
		}
		return path, "updated", nil
	}

	sep := "\n"
	if !strings.HasSuffix(s, "\n") {
		sep = "\n\n"
	} else if !strings.HasSuffix(s, "\n\n") {
		sep = "\n"
	}
	if err := os.WriteFile(path, []byte(s+sep+block), 0o644); err != nil {
		return path, "", err
	}
	return path, "appended", nil
}
