---
id: 01M2E248SM9X1JRZBNTHC9V7ZV
title: "Rollen-Prompts im Binary ausliefern: jaira roles install"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira liefert die Rollen-Prompts selbst aus: 'jaira roles install' schreibt sie nach .claude/skills/ (--project) oder ~/.claude/skills/ (--global), ohne lokale Aenderungen zu ueberschreiben"
context: |-
  Heute liegen die Rollen-Prompts nur in ~/.claude/skills auf genau einem Rechner. Ein Teamkollege klont das Repo, startet jaira und hat die Lanes, aber niemanden, der sie faehrt.
  Die Rollen sind: teamlead (redet mit dem Menschen, priorisiert), dispatcher (faehrt ein Ticket Lane fuer Lane, delegiert), role-lane (genau eine Lane), role-tester (Tests, ohne zu reparieren), role-pr (PR oeffnen, nie mergen), role-brainstorm, role-research.
  Der Mechanismus existiert schon zweimal im Code und wird nicht neu erfunden: core/lane/lane.go:34 bettet builtin/*.md per go:embed ein; core/board/announce.go:295 schreibt einen verwalteten Block in CLAUDE.md/AGENTS.md und schuetzt Handarbeit mit dem jaira:local-Marker.
  Drei Entscheidungen sind schon gefallen:
  - Flacher Ordnername mit Praefix, nicht verschachtelt: .claude/skills/jaira-teamlead/SKILL.md. Verschachtelte Ordner sind in Claude Code kein Namensraum, nur Plugins haben einen - eine Rolle unter skills/jaira/teamlead/ wuerde stillschweigend nicht geladen.
  - Zielordner wie in announce.go: in die Skill-Ordner schreiben, die es gibt (.claude/, .codex/, .agents/), und wenn es keinen gibt, .claude/ anlegen.
  - Eine Datei, die vom eingebetteten Original abweicht, wird ohne --force nicht ueberschrieben. Wie bei 'lanes adopt'. Eine selbst angepasste Rolle ueberlebt 'jaira update'.
  Die sieben Prompt-Dateien sind bereits geschrieben und liegen unter ~/.claude/skills/ - sie sind der Inhalt von core/role/builtin/, nicht neu zu erfinden.
  Nicht Teil dieses Tickets: ein roles-Marktplatz analog zu lanes/ auf GitHub. Erst wenn jemand eine eigene Rolle teilen will.
definition-of-done: "'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T18:57:58Z
updated-at: 2026-09-13T18:57:58Z
---

# Rollen-Prompts im Binary ausliefern: jaira roles install

## Definition of Done

- [ ] 'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

