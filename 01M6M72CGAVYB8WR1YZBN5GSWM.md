---
id: 01M6M72CGAVYB8WR1YZBN5GSWM
title: Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau
status: todo
ready: true
creator: berkc
assignee: berkc
goal: Die sechsstellige Ticket-ID auf jeder Board-Karte ist auf dunklem und hellem Terminal auf einen Blick lesbar.
context: |-
  Was falsch ist: die ID (z.B. 6TCTFV) auf der Karte ist kaum zu erkennen. berkc am 04.11.2026 mit Screenshot: 'die ist schwer zu erkennen'.
  Ursache: internal/tui/view.go:63 styHandle = Foreground(colFaint), colFaint = Color("240") (view.go:45) - dasselbe Grau wie Rahmen und Lane-Zaehler.
  Der Name daneben (berkc) ist heller, die ID wirkt dadurch wie ausgegraut.
  Nicht gewollt: eine Farbe pro Ticket oder eine neue Einstellung - nur ein lesbarer, ruhiger Ton, der nicht mit Tag- und Lane-Farben konkurriert.
definition-of-done: "styHandle nutzt eine hellere Farbe (z.B. 245-250 oder eine adaptive Farbe hell/dunkel), die auf der ausgewaehlten und der normalen Karte lesbar ist; Test in internal/tui"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T21:22:33Z
updated-at: 2026-11-04T21:25:40Z
updated-by: berkc
claimed-by: berk-ca-1280093
claimed-at: 2026-11-04T21:24:26Z
---

# Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau

## Definition of Done

- [ ] styHandle nutzt eine hellere Farbe (z.B. 245-250 oder eine adaptive Farbe hell/dunkel), die auf der ausgewaehlten und der normalen Karte lesbar ist; Test in internal/tui
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

