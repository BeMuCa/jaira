---
id: 01M5GFANC07HBSWWZXF78RD0YP
title: "Skill jaira-new-task: ein Interview fuehrt eine neue Aufgabe durch das Projekt, bevor sie ein Ticket wird"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Wer eine neue Aufgabe oder ein Feature beschreibt, wird im Interview durch das eigene Projekt gefuehrt und endet mit Tickets, die wissen, wo im Repo und in der Architektur die Arbeit liegt."
context: |-
  Idee von berkc am 22.10.2026, woertlich: 'Skill idea: Jaira-new-task, not sure if I already talked about this or have it as a ticket, but a combination of the brainstorming skill in an interview format. User describes a new task or feature, and the agent checks how it fits into the project they are working on and leads them through it in an interview style: where in the repo the feature and code parts would fit, the architectural changes from containers, microservices, frontend, database etc. based on the project - meaning if' - der Satz bricht dort ab; berkc fragen, wie er weitergeht.
  Was es schon gibt:
  - /jaira-role-brainstorm (core/role/builtin/jaira-role-brainstorm/SKILL.md): macht aus einer rohen Idee Tickets, prueft gegen CLAUDE.md, aber ohne Interview und ohne Architektur-Abgleich.
  - PP24H1 (Backlog): Interview-Skill, der Board, Projekt und Milestones aufsetzt - das ist Einrichtung eines Vorhabens, nicht eine einzelne neue Aufgabe in einem bestehenden Projekt. Beide teilen sich wahrscheinlich den Interview-Teil.
  - Die Lane Feature Ideas (W0JX10) ist der natuerliche Startpunkt: eine Idee dort, jaira-new-task macht daraus Backlog-Tickets mit follows auf die Idee.
definition-of-done: "Entschieden, ob jaira-new-task ein eigener Skill ist oder ein Interview-Modus von /jaira-role-brainstorm (Abgleich mit PP24H1)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:13Z
updated-at: 2026-10-22T00:14:13Z
---

# Skill jaira-new-task: ein Interview fuehrt eine neue Aufgabe durch das Projekt, bevor sie ein Ticket wird

## Definition of Done

- [ ] Entschieden, ob jaira-new-task ein eigener Skill ist oder ein Interview-Modus von /jaira-role-brainstorm (Abgleich mit PP24H1)
- [ ] Der Skill liest Repo-Struktur und Architektur (Container, Services, Frontend, Datenbank) und fragt pro Bereich, ob und wie die Aufgabe ihn beruehrt
- [ ] Ergebnis sind Backlog-Tickets mit follows auf die Idee; ausprobiert an einer echten Idee aus der Lane Feature Ideas

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

