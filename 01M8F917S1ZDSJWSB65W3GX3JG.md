---
id: 01M8F917S1ZDSJWSB65W3GX3JG
title: Worktree ohne .jaira findet das Board des Hauptcheckouts
status: in-progress
ready: true
creator: Sascha
goal: "jaira (und damit der cockpit) arbeitet in einem git-worktree auch dann, wenn dort kein .jaira liegt: es nimmt das Board des Hauptcheckouts."
context: |-
  Was falsch ist: Accept im cockpit schlägt in einem frischen Projekt fast immer fehl mit 'ticket: no .jaira directory found; run jaira init' (reason no_store).
  Warum: vor 'jaira share' steht .jaira nicht in git. Jeder worktree, den spawn.sh anlegt, hat darum kein .jaira. Dasselbe gilt für Branches, die vor share entstanden sind.
  Der cockpit ruft jaira im worktree des Tickets auf. ticket.Discover (core/ticket/store.go:268) sucht nur aufwärts nach .jaira/tickets und findet nichts.
  Heute hilft nur: im Hauptcheckout 'jaira move <id> --to done --force' von Hand.
  Schon vorhanden: core/project/worktree.go mainCheckout() findet den Hauptcheckout über die commondir-Datei.
  Idee: Discover fällt, wenn aufwärts nichts gefunden wird und dir in einem linked worktree liegt, auf mainCheckout(dir) zurück. Ein worktree MIT eigenem .jaira bleibt bei seinem eigenen.
  Offen für pre-process: ob Schreiben aus dem worktree ins Hauptcheckout-Board Sperren/Refs richtig trifft (Store.Root zeigt dann auf den Hauptcheckout).
definition-of-done: Im worktree ohne .jaira findet jaira list/show/move das Board des Hauptcheckouts (Test in core/ticket oder core/project)
tags:
  - cli
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-27T19:52:04Z
updated-at: 2026-11-27T20:34:47Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-68058
claimed-at: 2026-11-27T20:33:41Z
---

# Worktree ohne .jaira findet das Board des Hauptcheckouts

## Definition of Done

- [ ] Im worktree ohne .jaira findet jaira list/show/move das Board des Hauptcheckouts (Test in core/ticket oder core/project)
- [ ] Ein worktree mit eigenem .jaira benutzt weiter sein eigenes (Test)
- [ ] cockpit Accept auf ein Ticket, dessen worktree kein .jaira hat, geht durch (gesehen von einer Person)
- [ ] Release-Note core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-27 20:33 · Sascha** — parent session: uds:/run/user/1000/cc-socks/24078.sock
- **2026-11-27 20:34 · Sascha** — decisions (worker, none open for a person): 1) fallback lives in ticket.Discover, only after the upward walk finds nothing - a worktree with its own .jaira keeps it. 2) the git-file reading of project/worktree.go moves into a new leaf package core/gitwt (Top, MainCheckout, Linked) - project imports ticket, so ticket cannot call project; gitwt keeps one copy. 3) writes from the worktree land on the main checkout's board: locks and state are that board's, refs/jaira are shared by all worktrees, CommitsForTicket reads HEAD plus --branches, so branch commits are still found. Known limit: Branch() (create --follows from the branch name) reads the main checkout's branch when run from such a worktree. 4) DoD 3 is the person's check at signoff.
