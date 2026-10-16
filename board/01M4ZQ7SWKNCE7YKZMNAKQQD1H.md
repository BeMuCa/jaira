---
id: 01M4ZQ7SWKNCE7YKZMNAKQQD1H
title: "Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Jede Swimlane der Milestone-Roadmap ist eine 'Road' mit Namen: ohne Namen Road-1 bis Road-n, und im TUI waehlt man eine Road mit den Pfeiltasten und benennt sie mit Enter um"
context: |-
  Was heute fehlt: Die Roadmap (L, aus D0SAHM) zeigt Swimlanes, also Ketten von Milestones, die ueber after: verbunden sind. Eine Swimlane hat keinen Namen.
  Der Anfang einer Swimlane ist ein Milestone ohne Vorgaenger ('jaira milestone create --after ""' macht ihn zum Anfang einer eigenen Swimlane).
  Berk am 15.10.: 'die einzelnen Swimlanes als Milestone Road bezeichnen, und die sollen jeweils auch Namen bekommen. Per Default zaehlen sie einfach hoch, Road-1 .. Road-n. Aber man kann auf sie wechseln per Pfeiltaste und mit Enter deren Bezeichnung dann aendern. Dafuer muss man den Milestones dann wahrscheinlich noch einen Road-Namen hinzufuegen?'
  Vorschlag: Den Namen traegt nur der erste Milestone der Road, als Feld 'road: <name>'. Die anderen erben ihn ueber die Kette, damit es keine doppelte Wahrheit gibt. Ein Milestone wird wie heute ueber seinen Ref geteilt (seit EBTX6X ueber die ID).
  Offen fuer brainstorm: Was passiert, wenn zwei Roads zusammenlaufen (ein Milestone mit after auf zwei Roads), und welche Pfeiltasten (in der Roadmap bewegen sich Pfeile heute zwischen Bubbles).
definition-of-done: "Jede Swimlane der Roadmap zeigt einen Road-Namen: 'road:' ihres ersten Milestones, sonst Road-1..Road-n in Anzeigereihenfolge; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-15T12:05:22Z
updated-at: 2026-10-15T12:05:22Z
---

# Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen

## Definition of Done

- [ ] Jede Swimlane der Roadmap zeigt einen Road-Namen: 'road:' ihres ersten Milestones, sonst Road-1..Road-n in Anzeigereihenfolge; Test
- [ ] In der Roadmap waehlt eine Pfeiltaste die Road, Enter oeffnet ein Namensfeld, Enter speichert 'road:' am ersten Milestone (geteilt wie jede Milestone-Aenderung), Esc bricht ab; Test
- [ ] 'jaira milestone set <name> --road <road>' setzt den Namen auch ueber die CLI; 'milestone ls --json' traegt road; Test
- [ ] README, docs/COMMANDS.md und eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

