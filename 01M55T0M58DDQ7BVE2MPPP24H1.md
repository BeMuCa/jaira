---
id: 01M55T0M58DDQ7BVE2MPPP24H1
title: "Ein Interview-Skill fragt den Menschen aus und baut daraus Board, Projekt und Milestones, die zu seinem Vorhaben passen"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein neuer Role-Skill (Arbeitsname /jaira-role-interview) fuehrt ein Interview mit dem Menschen, hilft ihm seine Absicht klar zu denken, und legt danach Board (Lanes), Projekt-Setup und Milestones genau passend an."
context: |-
  Ausloeser: Berk am 17.10.: 'Interview Mode skill, um Board-, Projekt- und Milestone-Erstellung zu unterstuetzen. Aehnlicher Prompting-Stil wie unser Brainstorming-Skill, aber kombiniert mit einem besseren Interview-Stil, der beim Nutzer Gedanken anstoesst und seine Absichten herausfindet. Dann soll der Agent ein Board anlegen, das genau seinen Beduerfnissen entspricht.'
  Heute: core/role/builtin/jaira-role-brainstorm/SKILL.md macht aus einer Idee Tickets (Fragen: was tut heute weh, was hat es ausgeloest, was ist ausgeschlossen, was muesste wahr sein damit es eine schlechte Idee ist). Fuer Boards, Lanes und Milestones gibt es keinen gefuehrten Weg: der Mensch muss 'jaira lanes add/template', 'jaira milestone create --after' usw. selbst kennen.
  Bausteine, die schon da sind: mehrere Boards pro Repo (S0EG0M, done), 'jaira lanes template --board', 'jaira lanes market', 'jaira milestone create <name> --after <a,b>', 'jaira epic create'.
  Offen fuer brainstorm: (1) eine Rolle mit drei Modi (board/projekt/milestone) oder drei Rollen; (2) was 'Projekt' hier heisst - ein Repo mit 'jaira init', oder ein Epic; (3) Interview-Technik: eine Frage pro Runde, Optionen anbieten, Zusammenfassen-und-bestaetigen-lassen vor jedem Schreiben; (4) wird mit dem Binary ausgeliefert (core/role/builtin) oder bleibt persoenlich.
definition-of-done: |-
  Der Skill liegt unter core/role/builtin/ und wird von 'jaira roles install' mitinstalliert
  Er stellt eine Frage pro Runde und fasst vor jedem Schreiben zusammen, was er anlegen wird, und wartet auf ein Ja
  Er legt am Ende Board-Lanes, Milestones und ggf. Epics nur ueber die jaira-CLI an
  Ein Probelauf in einem Temp-Repo ist im Ticket beschrieben: Fragen, Antworten, entstandenes Board
  Release-Note
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-17T20:49:20Z
updated-at: 2026-10-17T23:48:15Z
updated-by: berkc
---

# Ein Interview-Skill fragt den Menschen aus und baut daraus Board, Projekt und Milestones, die zu seinem Vorhaben passen

## Definition of Done

- [ ] Der Skill liegt unter core/role/builtin/ und wird von 'jaira roles install' mitinstalliert
Er stellt eine Frage pro Runde und fasst vor jedem Schreiben zusammen, was er anlegen wird, und wartet auf ein Ja
Er legt am Ende Board-Lanes, Milestones und ggf. Epics nur ueber die jaira-CLI an
Ein Probelauf in einem Temp-Repo ist im Ticket beschrieben: Fragen, Antworten, entstandenes Board
Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

