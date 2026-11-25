---
id: 01M86NG5EY79JY4ERRXPE3MBCB
title: "Epic wird Mission mit eigenen Lanes: spec, tickets, knowledge - Spec-Bestaetigung durch den Menschen"
status: backlog
ready: true
creator: Sascha
assignee: Sascha
goal: "Eine Mission (heute Epic) laeuft durch eigene Lanes; in 'spec' schreibt ein Agent die Spec-Delta und ADRs einmal fuer die ganze Mission, ein Mensch bestaetigt sie (wie signoff), Tickets der Mission verweisen nur auf Requirement-IDs, und 'knowledge' prueft am Ende, dass Specs und ADRs aktuell sind."
context: |-
  Teil von SDD (siehe 7YQ1JC). Sascha 24.11.2026: Entwicklung nicht bremsen, sondern beschleunigen; Menschen schreiben keine Specs, sie bestaetigen nur - per Lane-Workflow. Epic -> Mission + Mission Control kam aus der Begriffsdiskussion vom 23.11. (Whiteboard).
  Vorschlag Lanes: idea -> spec (Agent schreibt Delta, Mensch bestaetigt) -> tickets (Tickets entstehen mit 'spec: <ID>') -> knowledge (jaira validate + eine billige Modellpruefung: Delta eingearbeitet, Proposed geschlossen) -> done.
  Gross: zuerst ein ADR (Proposed), den Sascha und Berk annehmen, dann Umsetzung. Berk arbeitet an Roads/Roadmap (KQQD1H) - abstimmen.
  Token: nur spec und knowledge lesen Specs (Index + eine Datei); normale Tickets wie heute.
definition-of-done: "ADR 'Epic wird Mission mit eigenen Lanes' (Proposed) liegt vor und ist von Sascha und Berk angenommen"
tags:
  - lanes
blocked-by: []
related: []
follows: 01M86NFX4NQF351XWAXZ7YQ1JC
commits: []
created-at: 2026-11-24T11:36:46Z
updated-at: 2026-11-24T11:36:46Z
---

# Epic wird Mission mit eigenen Lanes: spec, tickets, knowledge - Spec-Bestaetigung durch den Menschen

## Definition of Done

- [ ] ADR 'Epic wird Mission mit eigenen Lanes' (Proposed) liegt vor und ist von Sascha und Berk angenommen
- [ ] Mission-Lanes umgesetzt; in 'spec' sieht die Person Delta + ADR und bestaetigt per Accept; Test
- [ ] Tickets einer Mission tragen 'spec: <ID>'; 'knowledge' laesst die Mission erst nach jaira validate und Delta-Pruefung nach done; Test
- [ ] Release-Note in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

