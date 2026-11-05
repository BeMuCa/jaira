---
id: 01M6HZQRR1QPZFEX7BBYRH10T7
title: "Ein Worker, der auf eine Antwort der Person wartet, laesst das Ticket nicht still in einer Agenten-Lane liegen"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wartet ein Worker oder Dispatcher auf eine Entscheidung der Person, sieht man das im Board als wartend auf die Person, nicht als Ticket in critique ohne Agent"
context: |-
  Beobachtet 03./04.11.2026: CZWKNY stand in critique, kein Agent aktiv; MG80DF vorher genauso in todo. Nicht das erste Mal (berkc).
  Ursache: jaira-worker/SKILL.md 'How you ask' - ein Worker ohne AskUserQuestion (Subagent) schreibt die Frage als Note und ins question-Feld und 'stops in its lane'. Der Subagent endet; das Ticket bleibt in der Agenten-Lane. Der Teamlead muss die Antwort zurueckbringen und den Worker neu starten - vergisst er das (hier: Antwort war eine Rueckfrage), liegt das Ticket still.
  Vorschlag: statt in der Lane zu bleiben, nach human mit --question moven (die Lane ist genau dafuer da, und das Board/Cockpit zeigt sie als 'wartet auf dich'); nach der Antwort zurueck in die Lane. Verwandt: SPVWWH (Tickets ohne lebenden Claim melden).
definition-of-done: "jaira-worker und jaira-dispatcher SKILL.md: eine offene Frage an die Person geht mit --question in die human-Lane, nach der Antwort zurueck in die Lane"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T00:35:56Z
updated-at: 2026-11-04T00:35:56Z
---

# Ein Worker, der auf eine Antwort der Person wartet, laesst das Ticket nicht still in einer Agenten-Lane liegen

## Definition of Done

- [ ] jaira-worker und jaira-dispatcher SKILL.md: eine offene Frage an die Person geht mit --question in die human-Lane, nach der Antwort zurueck in die Lane
- [ ] Shipped-Hashes angehaengt, Release-Notiz

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

