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

	// jairaMarkerLocal opens a project's own area inside the otherwise
	// wholesale-regenerated block. Before this existed, a project's own rules
	// had nowhere to live but outside the managed block, where the next
	// regeneration could not touch them — but they ended up contradicting the
	// note instead, since the two were never written with each other in mind.
	// Everything from this marker to the end marker is now preserved verbatim
	// across regeneration, so project rules belong there instead of fighting
	// the block from outside it.
	jairaMarkerLocal = "<!-- jaira:local -->"
)

// agentNote is what an agent reads on entering the repository.
//
// The skill's own description already says to use it in a repository with a
// .jaira directory, but that relies on the model noticing. A line in the file
// the agent is given at the start of every session does not.
const agentNote = "## Task tracking: jaira\n" +
	"\n" +
	"This repository has a jaira board (`.jaira/`). Multi-step work is tracked\n" +
	"there as markdown tickets so it survives session boundaries.\n" +
	"\n" +
	"Capturing and picking work:\n" +
	"\n" +
	"- `jaira create <title> --goal <...> --context <...> --dod <...>` — one call files a\n" +
	"  complete ticket; without a goal, a definition of done, the context it came from\n" +
	"  and an assignee it cannot leave the backlog\n" +
	"- the context is the only record of why a ticket exists. Write it for someone who\n" +
	"  was not in this conversation and reads it weeks from now: what is wrong today,\n" +
	"  what triggered it, what is already known or ruled out. Write it as if that\n" +
	"  reader has mild ADHD and knows none of what you know — lead with what is wrong,\n" +
	"  short concrete lines with one point each, names and paths rather than\n" +
	"  adjectives, no jargon and no preamble. It may span several lines, but someone\n" +
	"  should be able to act after the first two. If acting on it would need a\n" +
	"  question answered first, it is not finished\n" +
	"- `jaira list --actionable --json` — everything that could be started right now\n" +
	"- `jaira next --json` — the single next actionable ticket\n" +
	"\n" +
	"Working a ticket:\n" +
	"\n" +
	"- `jaira claim <id>` — take it first; other sessions read this board too\n" +
	"- `jaira show <id> --for-lane <lane> --json` — the lane's prompt, the bounded input,\n" +
	"  the model tier, and the outputs the lane expects back\n" +
	"- `jaira dod <id> <n> --doing|--done` — mark checklist items as you go\n" +
	"- `jaira note <id> <text>` — at every pause, write down what the repository does\n" +
	"  not already say: dead ends, why this and not that, what you had to find out.\n" +
	"  Not what the checklist and git already record. A killed session never gets a\n" +
	"  turn to write anything down, so do not save it for the end\n" +
	"- `jaira move <id> --to <lane> --what <...> --why <...> --resolves <...>` — finish\n" +
	"  the step. jaira works out the commit list itself from git history — the union\n" +
	"  of the ticket file's own history and commits naming its id — so nothing needs\n" +
	"  to be typed here; it is written onto the ticket once, when the ticket leaves\n" +
	"  the board\n" +
	"- `jaira resume` — work left in progress, with everything recorded about it\n" +
	"- **the ticket rides in the same commit as the code.** Move the ticket first,\n" +
	"  then `git add` the changed file under `.jaira/tickets/` alongside your source\n" +
	"  changes and commit them together. A reviewer then sees the change and what it\n" +
	"  was for in one place, instead of a diff whose ticket is still in whatever\n" +
	"  state the last commit left it. Same for a ticket you create and hand to\n" +
	"  someone else: commit it, or nobody but you knows it exists — and now this is\n" +
	"  also what makes the commit list derivable at all: that shared commit is how\n" +
	"  git ties the ticket to the change\n" +
	"- `jaira sync <id>` — once a ticket reaches the terminal lane, stamps its\n" +
	"  commits and files it under `.jaira/sync/<you>-<date>/`, taking it off the\n" +
	"  board. `jaira restore <file>` brings it back\n" +
	"\n" +
	"`jaira <command> --help` for everything else.\n" +
	"\n" +
	"Do not edit files under `.jaira/tickets/` directly; the CLI is the write path.\n" +
	"The human review lane cannot be left by an agent — a person accepts the work there.\n" +
	"\n" +
	"Anything written after the `jaira:local` marker (see below) and before the end\n" +
	"marker survives the next time this section is regenerated, so project-specific\n" +
	"rules belong there rather than fighting this note from outside it."

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
		newBlock := strings.TrimSuffix(block, "\n")
		// The local marker only counts when it sits strictly inside the
		// managed block; one written before the start marker or after the end
		// marker is ordinary user text, not the boundary, and is left alone.
		if localRel := strings.Index(s[start:end], jairaMarkerLocal); localRel >= 0 {
			localIdx := start + localRel
			preserved := strings.TrimRight(s[localIdx+len(jairaMarkerLocal):end], "\n")
			newBlock = jairaMarkerStart + "\n" + agentNote + "\n\n" + jairaMarkerLocal + preserved + "\n" + jairaMarkerEnd
		}
		updated := s[:start] + newBlock + s[end+len(jairaMarkerEnd):]
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
