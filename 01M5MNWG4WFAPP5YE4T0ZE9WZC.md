---
id: 01M5MNWG4WFAPP5YE4T0ZE9WZC
title: "Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen"
status: in-progress
ready: true
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
updated-at: 2026-10-26T19:17:42Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-253301
claimed-at: 2026-10-26T19:10:55Z
---

# Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen

## Definition of Done

- [x] Eine committete .claude/settings.json erlaubt die Befehle der Lanes (go build/vet/test/gofmt, git add/commit/status/diff/log/show/worktree/branch/fetch/stash, alle jaira-Unterbefehle, mkdir/cp/ls/cat/grep/find, tmux ohne kill-server) und verbietet rm -rf ausserhalb des Scratchpads, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server
  proof: .claude/settings.json:3
- [ ] docs/COMMANDS.md oder README sagt: Worktrees lesen die Settings des Haupt-Checkouts; eine Verbotsliste gehoert in dessen settings.local.json oder in die User-Settings
- [ ] Ein frischer Worker auf einem Scratch-Board faehrt eine Lane ohne eine einzige Rueckfrage durch; das Ergebnis steht als Notiz am Ticket

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:12 · BeMuCa** — Driver: pre-process was passed by (next_lane went todo -> in-progress, the plan is empty). The lane worker gets its plan from the DoD plus the driver brief. Permission semantics are being checked against the Claude Code docs before writing; the user settings put their jaira rules under autoMode.allow, not permissions.allow. The dispatcher's spawn.sh drives herdr panes, not tmux, so the tmux rules only cover what a session types by hand.
- **2026-10-26 19:17 · BeMuCa** — Driver: the DoD as written cannot be met literally. Claude Code checks deny before allow and ignores how narrow a rule is, so 'deny rm -rf except in the scratchpad' cannot be written. Chosen instead: deny rm -rf / rm -fr everywhere and allow 'rm -r /tmp/*' for scratch boards. Edit(/**) was added beyond the DoD list, because default mode prompts on every file edit and DoD 3 could not pass without it. The user settings' autoMode.allow entries such as 'Bash(jaira create:*)' are prose for the auto-mode classifier, not tool patterns, and autoMode is never read from a project file. That is the likely reason the six jaira commands looked like the only ones allowed. Sources: code.claude.com/docs/en/permissions.md, settings.md, auto-mode-config.md.
- **2026-10-26 19:17 · BeMuCa** — Beyond the ticket list I added Bash(git ls-files *): core/role/builtin/jaira-role-lane/SKILL.md:108 tells a lane to run 'git ls-files --others --exclude-standard'. Also Bash(git config user.email) (exact, no wildcard): house rule 1 (maintainer's global instructions) requires it before the first commit; not in a lane prompt. No other lane prompt in .jaira/lanes/*.md names a shell command outside the rules.
