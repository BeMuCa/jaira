---
id: 01M8F917S1ZDSJWSB65W3GX3JG
title: Worktree ohne .jaira findet das Board des Hauptcheckouts
status: critique
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
updated-at: 2026-11-27T20:39:42Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-68058
claimed-at: 2026-11-27T20:33:41Z
outcome-what: "ticket.Discover falls back to the main checkout's board when a linked worktree has none; the git-file reading moved from core/project into a new leaf package core/gitwt shared by both"
outcome-why: "worktrees cut before jaira share have no .jaira, so every jaira call there - the cockpit's Accept included - failed with no_store"
outcome-resolves: jaira list/show/move and cockpit Accept work in a worktree without .jaira; a worktree with its own board keeps it
review-summary: "none"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Worktree ohne .jaira findet das Board des Hauptcheckouts

## Definition of Done

- [x] Im worktree ohne .jaira findet jaira list/show/move das Board des Hauptcheckouts (Test in core/ticket oder core/project)
  proof: core/ticket/store.go:284 fallback; internal/cli/worktreeboard_test.go TestWorktreeWithoutBoardUsesTheMainCheckouts (list/show/move); core/ticket/discover_test.go TestDiscoverInLinkedWorktree, TestDiscoverWorktreeWritesMainBoard
- [x] Ein worktree mit eigenem .jaira benutzt weiter sein eigenes (Test)
  proof: core/ticket/discover_test.go TestDiscoverInLinkedWorktree cases 'the worktree's own board wins' (+ from below)
- [ ] cockpit Accept auf ein Ticket, dessen worktree kein .jaira hat, geht durch (gesehen von einer Person)
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/3GX3JG.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-27 20:33 · Sascha** — parent session: uds:/run/user/1000/cc-socks/24078.sock
- **2026-11-27 20:34 · Sascha** — decisions (worker, none open for a person): 1) fallback lives in ticket.Discover, only after the upward walk finds nothing - a worktree with its own .jaira keeps it. 2) the git-file reading of project/worktree.go moves into a new leaf package core/gitwt (Top, MainCheckout, Linked) - project imports ticket, so ticket cannot call project; gitwt keeps one copy. 3) writes from the worktree land on the main checkout's board: locks and state are that board's, refs/jaira are shared by all worktrees, CommitsForTicket reads HEAD plus --branches, so branch commits are still found. Known limit: Branch() (create --follows from the branch name) reads the main checkout's branch when run from such a worktree. 4) DoD 3 is the person's check at signoff.
- **2026-11-27 20:39 · Sascha** — critique: core/project/worktree.go:50 Trees' second comment paragraph repeats gitwt.Linked's doc - drop it. Let stand: the fallback walks up from the main-checkout place past the repo top, same as a command run in the main checkout itself would - the worktree behaves as its main checkout; a bound to the repo top would make the two disagree.
