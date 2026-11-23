---
id: 01M81GE85YVMDSRN9YYNGPSAR2
title: Die Board-Dateien in master stehen auf jaira 0.6.1
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein frischer worktree hat nach jaira update keine uncommitteten Board-Dateien mehr.
context: |-
  Problem: die Board-Dateien in master stehen noch auf einem alten Stand; 'jaira update --lanes' auf 0.6.1 ändert sie in jedem Checkout.
  - Folge: in jedem worktree liegen dieselben uncommitteten Änderungen an .jaira/lanes/*.md, AGENTS.md, CLAUDE.md. Der signoff-Gate verweigert dann den Schritt (KRBGV9, 7MG5GB am 22.11.), und reap.sh behält worktrees (409XCV, PDZ444).
  - Lösung: die Ausgabe von 'jaira update --roles --lanes' 0.6.1 einmal in master committen.
definition-of-done: "master enthält .jaira/lanes/*.md, AGENTS.md und CLAUDE.md so, wie jaira update 0.6.1 sie schreibt"
tags:
  - release
blocked-by: []
related: []
commits:
  - 8476bc4791435252749a345c658e84d23bf574f5
created-at: 2026-11-22T11:32:08Z
updated-at: 2026-11-22T12:20:13Z
updated-by: Alexander Sacharov
---

# Die Board-Dateien in master stehen auf jaira 0.6.1

## Definition of Done

- [x] master enthält .jaira/lanes/*.md, AGENTS.md und CLAUDE.md so, wie jaira update 0.6.1 sie schreibt
  proof: 8476bc479: .jaira/lanes/*.md, AGENTS.md, CLAUDE.md aus jaira update 0.6.1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

