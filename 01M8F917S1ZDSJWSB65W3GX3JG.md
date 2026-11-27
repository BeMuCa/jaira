---
id: 01M8F917S1ZDSJWSB65W3GX3JG
title: Worktree ohne .jaira findet das Board des Hauptcheckouts
status: todo
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
updated-at: 2026-11-27T19:52:41Z
assignee: Sascha
updated-by: Sascha
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

