---
id: 01M7WWYKK8CMJSM2S45WJDDCEZ
title: "Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Hat ein Repository .jaira/worktree-setup, kann es ein .jaira/worktree-teardown haben; reap.sh ruft es vor git worktree remove auf (z.B. Docker-Stack des Tickets stoppen und Volumes loeschen)."
context: |-
  Heute: spawn.sh ruft .jaira/worktree-setup <worktree> <slug> <port-offset> <repo-root> auf (~/.claude/skills/jaira-dispatcher/scripts/spawn.sh:47,170).
  reap.sh entfernt Worktree und Branch, kennt aber nichts, was das Setup gestartet hat: Container, Ports, Volumes laufen weiter.
  In einem Projekt liefen zwei Stacks auf demselben Port, einmal mit Ports und Projektname des Hauptstacks. Freie Ports pruefen ist Sache des Setup-Skripts im Projekt; das Aufraeumen braucht aber einen Haken in jaira.
  Alex 20.11.2026: "ok, aber man muss es auch schliessen".
definition-of-done: "reap.sh ruft ein ausfuehrbares .jaira/worktree-teardown mit denselben Argumenten wie das Setup auf, bevor es den Worktree entfernt (Test)"
tags: []
blocked-by:
  - 01M7R61GC6C8D7547N7CN7JW2Q
related:
  - 01M2MFEEE9VZAABPKKS27MG5GB
  - 01M6M72EN4ZZ2JBZJ1Q6DSDPZ4
commits: []
created-at: 2026-11-20T16:34:35Z
updated-at: 2026-11-20T16:50:10Z
updated-by: Alexander Sacharov
---

# Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab

## Definition of Done

- [ ] reap.sh ruft ein ausfuehrbares .jaira/worktree-teardown mit denselben Argumenten wie das Setup auf, bevor es den Worktree entfernt (Test)
- [ ] Scheitert das Teardown, behaelt reap.sh den Worktree und nennt den Grund ('kept <ticket>: teardown failed')
- [ ] spawn.sh --help und reap.sh --help beschreiben das Paar setup/teardown

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

