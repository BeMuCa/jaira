---
id: 01M6HZQPZZFAF3MXY3AWQFWRXH
title: "Eine PR-Lane zeigt gebuendelt, welche Tickets in welchem PR liegen, mit Name, Branch und Pipeline-Farbe"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Im Board sieht man pro offenem PR eine gebuendelte Gruppe von Tickets (aufklappbar beim Hovern/Auswaehlen) mit PR-Name, Branch und einer Farbe fuer den CI-Stand"
context: |-
  Wunsch berkc 04.11.2026 (Punkt 11): 'PR Lane klingt eig fun, gebuendelt darstellen, wobei es sich ausbreitet wenn man drueber hovered. Es sollte den PR-Namen oder Branch-Namen halten, und Farbe nach State der Pipeline. Vielleicht Cron-Job, der stuendlich checkt?'
  Heute: PRs entstehen pro Milestone (6TCTFV): jaira-role-pr baut milestone/<name> aus origin/master und merged die Ticket-Branches mit --no-ff. Den PR-Stand kennt jaira nicht.
  Randbedingung: kein Daemon, kein Server in jaira. Kandidaten fuer den Stand: 'gh pr view --json' beim Auffrischen des Boards (nur wenn gh da ist), oder eine GitHub Action, die den Stand auf den Ticket-Ref schreibt.
  Offen fuer brainstorm: Lane oder nur Gruppierung im Board; Quelle und Haeufigkeit des PR-Stands; GitLab.
definition-of-done: Brainstorm hat Darstellung und Quelle des PR-Stands festgelegt
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T00:35:54Z
updated-at: 2026-11-04T21:27:04Z
updated-by: berkc
---

# Eine PR-Lane zeigt gebuendelt, welche Tickets in welchem PR liegen, mit Name, Branch und Pipeline-Farbe

## Definition of Done

- [ ] Brainstorm hat Darstellung und Quelle des PR-Stands festgelegt
- [ ] Tickets eines offenen PRs erscheinen gebuendelt mit PR-Name/Branch und Farbe nach CI-Stand; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

