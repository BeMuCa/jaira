---
id: 01M57ZXP95EN0S4B99GPNC323D
title: "Cockpit: Tickets aller Boards und der HITL-Lane in 'Waiting for you'"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Cockpit-Band 'Waiting for you' zeigt die Tickets, die auf einen Menschen warten, auf jedem Board des Repos und auch in der Lane human (HITL), nicht nur in signoff des Default-Boards."
context: |-
  Was falsch ist: Das Cockpit (core/mod/builtin/jaira-cockpit) sieht nur das Default-Board.
  Ursache 1: register.tsx refreshBoard ruft 'jaira lanes --json' und 'jaira list --json' ohne --board. Tickets auf .jaira/boards/<name> (z. B. video-Board, siehe S0EG0M/8NHD8Q) fehlen im Band.
  Ursache 1b: 'jaira lanes' kennt kein --board (internal/cli/tickets.go newLanesCmd); das Cockpit kann die Lanes eines anderen Boards gar nicht lesen.
  Ursache 2: board.ts isPersons liest nur human_exit. Die Lane human hat requires-question, nicht requires-human-exit, also human_exit=false in 'jaira lanes --json'. Tickets in HITL (z. B. 6TCTFV) erscheinen nie im Band.
  Der Kommentar in tickets.go bei human_exit behauptet '(signoff, human)' – stimmt nicht.
  Ausloeser: Alex, 2026-10-18: 'cockpit sieht Tickets anderer Boards nicht; HITL-Tickets sollen auch sichtbar sein, nicht nur signoff'.
  Ausser Scope: Overview-Pane und Workers-Tab bleiben beim Default-Board.
definition-of-done: "'jaira lanes --board <name> --json' liest die Lanes des Boards <name>; unbekanntes Board ist exit 2"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T17:11:04Z
updated-at: 2026-10-18T17:17:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1943142
claimed-at: 2026-10-18T17:11:16Z
outcome-what: jaira lanes kennt --board und gibt requires_question aus; das Cockpit-Band liest alle Boards aus jaira boards und zaehlt Lanes mit requires_question (human) als wartend
outcome-why: "Das Band las nur das Default-Board und erkannte HITL nicht, weil human nur requires-question traegt, nicht human_exit"
outcome-resolves: "Tickets anderer Boards und HITL-Tickets fehlten in 'Waiting for you'"
---

# Cockpit: Tickets aller Boards und der HITL-Lane in 'Waiting for you'

## Definition of Done

- [x] 'jaira lanes --board <name> --json' liest die Lanes des Boards <name>; unbekanntes Board ist exit 2
  proof: internal/cli/tickets.go newLanesCmd --board; TestLanesTakeABoard
- [x] 'jaira lanes --json' gibt pro Lane requires_question aus
  proof: internal/cli/tickets.go requires_question; TestLanesTakeABoard
- [x] Das Cockpit-Band zeigt wartende Tickets aller Boards aus 'jaira boards --json'
  proof: jaira-cockpit/hooks/register.tsx refreshBoard otherBoardsOf; board.test.ts 'tickets waiting on another board are listed too'
- [x] Das Cockpit-Band zeigt Tickets in einer Lane mit requires_question (human) als wartend, ohne Accept
  proof: jaira-cockpit/hooks/board.ts isPersons requires_question; board.test.ts 'the band lists the tickets in a person's lanes' (human lane now requires_question only)
- [x] Release-Note in core/release/notes/<ID>.md; Go- und Cockpit-Tests gruen
  proof: core/release/notes/NC323D.md; claude plugin test 160 pass; go test ./internal/cli ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-18 17:17 · Alexander Sacharov** — Ausser Scope geblieben: Overview-Pane und Workers-Tab lesen weiter nur das Default-Board (pro Board andere Lanes, braeuchte eine Board-Auswahl). Give back schiebt fest nach in-progress – auf einem Board ohne in-progress (video) scheitert das mit der Fehlermeldung von jaira move.
