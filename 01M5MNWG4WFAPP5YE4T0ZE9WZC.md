---
id: 01M5MNWG4WFAPP5YE4T0ZE9WZC
title: "Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen"
status: backlog
ready: false
creator: BeMuCa
goal: "Worker und Subagenten fragen nur bei wirklich kritischen Befehlen (rm, push, reset --hard, kill-server) nach; alles, was eine Lane braucht, ist im Repo erlaubt, und die Doku sagt, wo Verbots- und Erlaubnislisten liegen und welche Settings ein Worktree liest."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 5) und am 23.10. nachgeprueft: das Repo hat keine .claude/settings.json mit permissions; die User-Settings erlauben im Auto-Modus nur sechs jaira-Unterbefehle (create, list, next, show, claim, note). go build/vet/test, gofmt, git add/commit/worktree, jaira move/dod/set/pull/logbook fragen deshalb nach.
  Eine Verbotsliste, die ein Teamlead in einen Worktree legte, wurde nie gelesen: Worktrees lesen die Settings des Haupt-Checkouts - das steht nirgends.
  Verwandt: spawn.sh-Ticket zum Berechtigungsmodus (gleiche Runde). Der Skill 'fewer-permission-prompts' kann die Nur-Lese-Haelfte aus Transkripten vorschlagen.
definition-of-done: "Eine committete .claude/settings.json erlaubt die Befehle der Lanes (go build/vet/test/gofmt, git add/commit/status/diff/log/show/worktree/branch/fetch/stash, alle jaira-Unterbefehle, mkdir/cp/ls/cat/grep/find, tmux ohne kill-server) und verbietet rm -rf ausserhalb des Scratchpads, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:47Z
updated-at: 2026-10-23T15:25:47Z
---

# Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen

## Definition of Done

- [ ] Eine committete .claude/settings.json erlaubt die Befehle der Lanes (go build/vet/test/gofmt, git add/commit/status/diff/log/show/worktree/branch/fetch/stash, alle jaira-Unterbefehle, mkdir/cp/ls/cat/grep/find, tmux ohne kill-server) und verbietet rm -rf ausserhalb des Scratchpads, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server
- [ ] docs/COMMANDS.md oder README sagt: Worktrees lesen die Settings des Haupt-Checkouts; eine Verbotsliste gehoert in dessen settings.local.json oder in die User-Settings
- [ ] Ein frischer Worker auf einem Scratch-Board faehrt eine Lane ohne eine einzige Rueckfrage durch; das Ergebnis steht als Notiz am Ticket

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

