---
id: 01M2KBPPVH5PAKZ98B0X7KX89C
title: "Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree"
status: todo
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
updated-at: 2026-09-16T06:49:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30509
claimed-at: 2026-09-16T06:49:05Z
---

# Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree

## Definition of Done

- [ ] Der Skill jaira-dispatcher liegt unter .claude/skills/ im Repository, neben dem schon vorhandenen jaira-Skill; scripts/spawn.sh faehrt mit und kennt --no-worktree; ein frisch geklontes Repository kann einen Worker starten, ohne dass etwas aus ~/.claude kopiert wird

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

