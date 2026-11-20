---
id: 01M7WX0C09TC63TKPXT555CS2M
title: "Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Schlussmeldung eines Workers/Dispatchers hat ein festes Format, das ein Teamlead ohne Nachfragen und ohne das Ticket zu lesen auswerten kann."
context: |-
  Heute: Schlussmeldungen sind freier Text in drei Zeilen; der Teamlead liest das Ticket nach, um Branch, Commits und ob gepusht werden kann zu finden. Alex: "sonst liest er viel".
  Vorschlag des Teamleads: eine Zeile mit Branch, Commits (Spanne), Guard/Lint ok, offene Notes (Anzahl), Lane.
  Verwandt: MVR4PY (Cockpit liest Lane und Urteil aus einer festen Schlusszeile).
definition-of-done: "jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill"
tags: []
blocked-by: []
related:
  - 01M67RXSPRWC7WVBCCM7MVR4PY
commits: []
created-at: 2026-11-20T16:35:33Z
updated-at: 2026-11-20T16:35:33Z
---

# Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes

## Definition of Done

- [ ] jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill
- [ ] Das Cockpit liest die Felder aus dieser Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

