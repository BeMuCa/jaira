---
id: 01M6H737KDSDZK78K4KRPB664F
title: "jaira web 'Braucht dich': offene Fragen von Tickets in human und signoff im Browser beantworten"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM beantwortet die Fragen, an denen Agenten warten, im Browser - das ist der eigentliche Steuerhebel."
context: |-
  Was heute fehlt: steht ein Ticket mit offener Frage (Feld 'question') in einer Personen-Lane, sieht und beantwortet man sie nur im Terminal.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden: die Antwort wird eine Notiz am Ticket (wie 'jaira note') mit dem Namen des PM, und das Feld 'question' wird geleert. Die Lane bewegt sich dabei NICHT.
  Offen in CMN1YT, nicht hier: ob der PM signoff annehmen darf. Erst nach dieser Entscheidung ein eigenes Ticket.
  Voraussetzung: DVPAPE fuer den Namen an der Notiz.
definition-of-done: "Eine Liste zeigt jedes Ticket mit gesetztem 'question' in einer Personen-Lane, mit der Frage und ihren Optionen; Test"
tags:
  - web
blocked-by:
  - V9S8RN
  - DVPAPE
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:17Z
updated-at: 2026-11-03T17:25:17Z
---

# jaira web 'Braucht dich': offene Fragen von Tickets in human und signoff im Browser beantworten

## Definition of Done

- [ ] Eine Liste zeigt jedes Ticket mit gesetztem 'question' in einer Personen-Lane, mit der Frage und ihren Optionen; Test
- [ ] Antworten schreibt eine Notiz mit dem Namen des PM und leert 'question'; die Lane bleibt; Test
- [ ] jaira web --help, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

