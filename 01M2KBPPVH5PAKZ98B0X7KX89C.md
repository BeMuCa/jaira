---
id: 01M2KBPPVH5PAKZ98B0X7KX89C
title: "Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jeder, der dieses Repository klont, hat den Dispatcher-Skill und sein spawn.sh mitsamt --no-worktree - ohne dass ihm jemand Dateien aus seinem Heimverzeichnis schickt."
context: |-
  Der Dispatcher-Skill liegt heute NUR in ~/.claude/skills/jaira-dispatcher/ auf Alex' Maschine. Wer das Repository klont, bekommt ihn nicht. Das Repository traegt aber schon Skills: .claude/skills/jaira/SKILL.md ist da und wird mitgeklont - der Ort ist also erprobt, nur der Dispatcher steht nicht drin.

  Ausgeloest am 2026-09-15: Alex wollte, dass ein Worker nicht immer in einem eigenen Worktree startet, sondern auch in der Hauptordner-Arbeitskopie. Der Schalter wurde gebaut - scripts/spawn.sh nimmt jetzt --no-worktree (und liest JAIRA_NO_WORKTREE=1), SKILL.md beschreibt ihn. Beides steht in ~/.claude und ist damit auf genau einer Maschine vorhanden. Alex' Satz dazu: 'сделай тикет под это чтобы все могли пользоваться'.

  Was --no-worktree tut: wt = repo-root statt ../.worktrees/<repo>-<slug>; kein Worktree, kein eigener Branch. Gedacht fuer eine einzelne Lane auf dem schon ausgecheckten Branch. Der Preis steht im Hilfetext: zwei Worker teilen sich dann ein Verzeichnis, also nur einer gleichzeitig.

  Zu klaeren und Sache der Plan-Lane, nicht hier entschieden:
  - ob der Skill ins Repository kopiert wird oder ob ~/.claude/skills/jaira-dispatcher ein Symlink dorthin wird, damit es nicht zwei driftende Kopien gibt
  - ob die anderen jaira-role-*-Skills denselben Weg gehen; sie haben dasselbe Problem, sind aber nicht das, was Alex verlangt hat
  - spawn.sh setzt HERDR_* voraus. Auf einer Maschine ohne Herdr muss der Skill das sagen und nicht scheitern

  Die aktuellen Dateien liegen zum Uebernehmen bereit unter ~/.claude/skills/jaira-dispatcher/ (SKILL.md und scripts/spawn.sh, Stand 2026-09-15); scripts/*.bak sind alte Sicherungen und gehoeren NICHT mit ins Repository.
definition-of-done: "Der Skill jaira-dispatcher liegt unter .claude/skills/ im Repository, neben dem schon vorhandenen jaira-Skill; scripts/spawn.sh faehrt mit und kennt --no-worktree; ein frisch geklontes Repository kann einen Worker starten, ohne dass etwas aus ~/.claude kopiert wird"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T20:21:31Z
updated-at: 2026-09-16T06:57:04Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30509
claimed-at: 2026-09-16T06:49:05Z
outcome-what: "Plan geschrieben: der Dispatcher wird in core/role/builtin/ nachgezogen, nicht nach .claude/skills/ kopiert"
outcome-why: "Die Rollen-Infrastruktur mit go:embed und 'jaira roles install' existiert schon; eine zweite Kopie waere genau die Drift, gegen die das Ticket existiert"
outcome-resolves: "DoD auf den echten Ort umgeschrieben, 7 Planschritte, Symlink- und Herdr-Frage aus dem Kontext beantwortet"
---

# Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree

## Definition of Done

- [ ] Der Rollen-Ordner core/role/builtin/jaira-dispatcher/ traegt SKILL.md und scripts/spawn.sh im Stand von ~/.claude (2026-09-15); spawn.sh kennt --no-worktree; 'jaira roles install' in ein leeres Verzeichnis liefert beide Dateien, spawn.sh ausfuehrbar, ohne dass etwas aus ~/.claude kopiert wird

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Hunk --no-worktree aus ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh nach core/role/builtin/jaira-dispatcher/scripts/spawn.sh uebernehmen: usage(), Flag-Schleife, wt-Zweig
- [x] Hunks 'dispatch'-Lane und proj-Variable mitnehmen; die vier .env-Zeilen (VITE_PORT_HOST, BACKEND_PORT_HOST, IMAGE_NS, COMPOSE_PROFILES) auslassen
- [x] Den --no-worktree-Absatz in core/role/builtin/jaira-dispatcher/SKILL.md nachziehen, wortgleich zum Stand in ~/.claude
- [x] Pruefen, dass ohne Herdr weder Skript noch Skill ins Leere laufen: HERDR_ENV-Check meldet sauber, SKILL.md nennt die Ersatzwege
- [x] Test in core/role: 'bash -n' ueber jede eingebettete .sh-Datei, und 'spawn.sh --help' endet mit 0 - heute faellt ein kaputtes Skript erst beim Worker auf
- [x] Zeile unter '## Unreleased' in core/release/NOTES.md: --no-worktree da, und 'jaira roles install --force' noetig, weil eine veraenderte Kopie uebersprungen wird
- [x] Abnahme: 'jaira roles install' in ein leeres Verzeichnis, beide Dateien da, spawn.sh ausfuehrbar, --no-worktree drin; 'go test ./core/role/...' gruen

## Progress
- **2026-09-16 06:52 · Alexander Sacharov** — Die Praemisse des Tickets ist ueberholt: der Dispatcher liegt laengst im Repository, nur an einem anderen Ort als die DoD sagte. core/role/builtin/jaira-dispatcher/ traegt SKILL.md und scripts/spawn.sh, sie sind per go:embed im Binary und 'jaira roles install' schreibt sie nach ~/.claude/skills. .claude/skills/ waere ein zweiter Ort und damit genau die driftende Zweitkopie, gegen die das Ticket existiert. DoD entsprechend umgeschrieben; Symlink-Frage damit erledigt, das Installationskommando ist der Weg.
- **2026-09-16 06:52 · Alexander Sacharov** — Der Diff ~/.claude gegen core/role/builtin ist groesser als --no-worktree. Vier Hunks: (1) --no-worktree samt usage() und Flag-Schleife - das ist der Auftrag; (2) Lane-Name 'dispatch' startet /jaira-dispatcher statt /jaira-role-lane - generisch, kommt mit; (3) proj-Variable statt Inline-Ausdruck, reines Refactoring, kommt mit; (4) vier .env-Zeilen VITE_PORT_HOST, BACKEND_PORT_HOST, IMAGE_NS, COMPOSE_PROFILES - die nennen 'task app:dev' eines anderen Projekts und gehoeren nicht in eine generische Builtin-Rolle. Bleiben draussen, Folgeticket.
- **2026-09-16 06:52 · Alexander Sacharov** — jaira roles install ueberschreibt eine veraenderte Datei nicht, es meldet sie als 'skipped' (core/role/install.go). Alex' ~/.claude-Kopie ist veraendert, also bleibt sie nach diesem Ticket die Abweichung, bis er einmal 'jaira roles install --force' laeuft. Das gehoert in die NOTES.md-Zeile, sonst merkt es niemand.
- **2026-09-16 06:52 · Alexander Sacharov** — Ohne Herdr scheitert spawn.sh nicht stumm: es prueft HERDR_ENV=1 und sagt 'not inside a Herdr pane', exit 1. SKILL.md nennt zwei Ersatzwege (Peer-Sessions, Subagents). Die dritte offene Frage aus dem Kontext ist damit beantwortet, es bleibt ein Pruefschritt und keine Aenderung.
