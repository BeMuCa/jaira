---
id: 01M67C18P3AB3GAPJ103DAQKZX
title: "Beim Anlegen eines Tickets entscheidet der Teamlead, ob brainstorm und planning noetig sind, oder fragt"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Die optionalen Lanes brainstorm und pre-process werden bewusst an- oder abgewaehlt, statt dass die Options-Haken unbeachtet leer bleiben."
context: |-
  Frage von berkc am 30.10.2026: 'Do we use all the lanes? What does pre processing do? do we ever use brainstorm? ... When working with the Teamlead for creating tickets, it should decide or ask us if they are relevant.'
  Heute: brainstorm (requires-option: brainstorm) und pre-process (requires-option: planning) laufen nur, wenn im Options-Abschnitt des Tickets der Haken gesetzt ist (core/lane/lane.go:120). Sonst ueberspringt der Weg sie.
  brainstorm klaert, WAS das Ticket sein soll (Ziel); pre-process schreibt den Plan, WIE es gebaut wird.
  Kein Skill fragt heute danach.
definition-of-done: "Teamlead-Skill: beim Anlegen setzt er brainstorm/planning nach Regel (unklares Ziel -> brainstorm; mehr als eine Datei/Entscheidung -> planning) oder fragt einmal"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T21:39:11Z
updated-at: 2026-10-30T21:39:11Z
---

# Beim Anlegen eines Tickets entscheidet der Teamlead, ob brainstorm und planning noetig sind, oder fragt

## Definition of Done

- [ ] Teamlead-Skill: beim Anlegen setzt er brainstorm/planning nach Regel (unklares Ziel -> brainstorm; mehr als eine Datei/Entscheidung -> planning) oder fragt einmal
- [ ] jaira create kann die Options setzen (z.B. --option planning); Test
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

