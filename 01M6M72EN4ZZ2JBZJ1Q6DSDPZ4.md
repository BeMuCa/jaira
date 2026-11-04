---
id: 01M6M72EN4ZZ2JBZJ1Q6DSDPZ4
title: Ticket-Worktrees liegen gesammelt in einem Ordner jaira-worktrees/ neben dem Repository
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Alle Worktrees, die spawn.sh und die Worker fuer Tickets anlegen, liegen in <eltern>/<repo>-worktrees/<ID> statt einzeln neben dem Repository."
context: |-
  Was falsch ist: jeder Ticket-Worktree liegt direkt neben dem Repo (/home/berkc/code/X_General/jaira-KQQD1H, jaira-EH6Z0N, ...). Bei 10+ Tickets ist der Elternordner unuebersichtlich.
  Wunsch berkc 04.11.2026: 'Sollten wir die Worktrees ... immer in ein jaira-worktrees dir packen? dann haetten wirs gebuendelt'.
  Wo der Pfad entsteht: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (Worktree-Pfad), reap.sh (findet und loescht sie), jaira-worker/SKILL.md 'Where you work' (beside the repository and never inside it).
  Nicht ins Repo selbst legen: ein Worktree im Repo taucht in git status und in go ./... auf.
  Bestehende Worktrees nicht verschieben (git worktree move nur, wenn sauber); neue landen im neuen Ordner.
  Reihenfolge: nach 80MGXT, das dieselben Skripte (reap.sh) aendert.
definition-of-done: "spawn.sh legt neue Worktrees unter <eltern>/<repo>-worktrees/<id> an; reap.sh findet Worktrees in beiden Orten; Test fuer die Skripte"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T21:22:35Z
updated-at: 2026-11-04T21:22:35Z
---

# Ticket-Worktrees liegen gesammelt in einem Ordner jaira-worktrees/ neben dem Repository

## Definition of Done

- [ ] spawn.sh legt neue Worktrees unter <eltern>/<repo>-worktrees/<id> an; reap.sh findet Worktrees in beiden Orten; Test fuer die Skripte
- [ ] jaira-worker und jaira-dispatcher SKILL.md nennen den neuen Ort; shipped-Hashes ergaenzt
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

