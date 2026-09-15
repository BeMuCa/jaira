---
id: 01M2HV0SV78FDFEZFSECJSRN8G
title: Dieselbe Regel ueber fertige Tickets steht dreimal von Hand im Code
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Was als 'in dieser Lane und diesem Klon zum Ablegen' gilt, steht an einer Stelle - und was ein abgelegtes Ticket im JSON heisst, heisst es in beiden Befehlen gleich."
context: |-
  Aus der review-Lane von 74VM40 am 2026-09-15, ausdruecklich als Folgearbeit benannt und dort nicht behoben.

  Erstens: die Bedingung 'Status == lane && !t.ReadOnly' steht dreimal von Hand geschrieben. 74VM40 hat zwei davon in filesInLane() zusammengezogen, die dritte in internal/tui/model.go:764 blieb - die Funktion ist nicht exportiert und die TUI zaehlt ohnehin ueber m.tickets statt ueber den Store. Drei Kopien einer Regel laufen auseinander, und diese hier entscheidet, was das Board als 'N to file' anzeigt: weicht sie von dem ab, was 'jaira logbook --all' wirklich ablegt, sagt das Board eine Zahl und der Befehl tut etwas anderes. Genau dieser Widerspruch ist in 74VM40 schon einmal aufgetreten.

  Zweitens: 'jaira logbook --all --json' traegt seit 74VM40 ein 'title' je Eintrag, weil es auf trimmedJSON gefaltet wurde - denselben Renderer, den 'jaira move' fuer ein weggeraeumtes Ticket benutzt. 'jaira logbook <id> --json' wurde nicht mitgezogen und beschreibt dasselbe abgelegte Ticket weiterhin ohne title. Ein Programm, das beide liest, bekommt fuer dieselbe Sache zwei Formen.

  Nicht Teil dieses Tickets: der Zaehler auf dem Board selbst und die Schwelle von zehn - beides ist in 74VM40 entschieden und belegt.
definition-of-done: "Die Regel, welche Tickets dieser Klon in einer Lane ablegen darf, steht an genau einer Stelle im Code, und die TUI liest sie von dort statt sie ein drittes Mal zu schreiben. Nachgestellt, indem die Regel an dieser einen Stelle geaendert wird und sich Board-Zaehler und Befehl gemeinsam mitaendern."
tags:
  - cli
blocked-by: []
related:
  - 01M28MHSDBABYVD8785A74VM40
commits: []
created-at: 2026-09-15T06:10:42Z
updated-at: 2026-09-15T06:10:42Z
---

# Dieselbe Regel ueber fertige Tickets steht dreimal von Hand im Code

## Definition of Done

- [ ] Die Regel, welche Tickets dieser Klon in einer Lane ablegen darf, steht an genau einer Stelle im Code, und die TUI liest sie von dort statt sie ein drittes Mal zu schreiben. Nachgestellt, indem die Regel an dieser einen Stelle geaendert wird und sich Board-Zaehler und Befehl gemeinsam mitaendern.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

