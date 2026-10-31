---
id: 01M67RYAXEZG7QD2XJ37EH6Z0N
title: "In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Steht ein Ticket in human oder signoff, zeigt die Ticketansicht seinen Diff gegen master Datei fuer Datei (wie die Merge-Ansicht von GitHub), und je Datei oeffnet ein Knopf die Datei in VS Code - so wird die Abnahme ein echter Code-Review."
context: |-
  Was fehlt: wer abnimmt, liest Outcome und Review-Text, aber nicht den Code. Den Diff muss man selbst in git oder auf GitHub suchen.
  Ausloeser: Maintainer 31.10.2026, Punkt 8, mit Screenshot (Image #7 - im Agenten-Kontext nicht sichtbar; wo genau die Ansicht hin soll, steht dort): 'Wenn wir im human review step sind, dann sollte man auch die Code changes sehen, wie bei git wenn man mergen will ... mit einem button kann man dann auch in vsc zu dem file wechseln.'
  Bekannt: der Lane-Payload kennt den Diff schon (input-requires: diff in .jaira/lanes/review.md). VS Code oeffnet eine Datei per 'code -g <pfad>:<zeile>'. Ort (Cockpit-Pane aus Q63EEV, TUI-Ticketansicht oder jaira web) im Brainstorm mit dem Maintainer klaeren.
definition-of-done: Brainstorm hat mit dem Maintainer den Ort festgelegt (Screenshot aus Punkt 8)
tags:
  - tui
blocked-by: []
related:
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-10-31T01:24:47Z
updated-at: 2026-10-31T01:24:47Z
---

# In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code

## Definition of Done

- [ ] Brainstorm hat mit dem Maintainer den Ort festgelegt (Screenshot aus Punkt 8)
- [ ] Ein Ticket in signoff zeigt seinen Diff je Datei; Test
- [ ] Ein Knopf/Taste oeffnet die Datei an der Zeile in VS Code
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

