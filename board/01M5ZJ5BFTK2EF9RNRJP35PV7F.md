---
id: 01M5ZJ5BFTK2EF9RNRJP35PV7F
title: "jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht"
status: backlog
ready: false
creator: BeMuCa
goal: "Legt jemand einen Milestone ohne --after an, sagt die Ausgabe, dass er im Picker steht, nicht auf der Roadmap, und wie er dorthin kommt (jaira milestone set <name> --after ...)."
context: |-
  Was schief ist: ein Milestone ohne --after landet im M-Picker, nicht auf der Roadmap; der Mensch schaute auf die Roadmap und hielt ihn fuer verloren (O14, 27.10.2026).
  Heute steht das nur in jaira milestone create --help ('Without --after the milestone stays off the roadmap'), nicht in der Ausgabe des Befehls.
  Ein Satz in der Ausgabe reicht.
definition-of-done: "jaira milestone create <name> ohne --after druckt einen Satz mit dem set --after Befehl"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:21Z
updated-at: 2026-10-27T20:52:21Z
---

# jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht

## Definition of Done

- [ ] jaira milestone create <name> ohne --after druckt einen Satz mit dem set --after Befehl
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

