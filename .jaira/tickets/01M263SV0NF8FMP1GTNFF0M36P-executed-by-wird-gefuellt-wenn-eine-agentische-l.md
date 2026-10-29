---
id: 01M263SV0NF8FMP1GTNFF0M36P
title: "executed-by wird gefuellt, wenn eine agentische Lane verlassen wird"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "An jedem Ticket, das eine agentische Lane durchlaufen hat, steht, welches Modell die Arbeit gemacht hat"
context: |-
  Auf dem Board eines anderen Projekts tragen 9 von 11 Tickets in review kein executed-by. Nur CNTYMW und ZTR2CW sagen 'sonnet'.
  Damit ist bei der Abnahme nicht erkennbar, ob ein Ticket von einem starken oder einem billigen Modell bearbeitet wurde - und die Lanes deklarieren genau das: in-progress hat 'model-tier: cheap', critique und optimize haben 'strong'.
  Das Feld existiert im Schema (core/ticket/schema.go, FieldExecutedBy) und 'jaira move' hat den Schalter --executed-by. Es wird nur von niemandem gesetzt und von nichts eingefordert.
  Kommentar im Schema zur Abgrenzung beachten: executed-by ist das Modell, die Verantwortung bleibt beim assignee. Es ersetzt also keinen Besitzer, es ergaenzt ihn.
  Haengt mit D28H7V zusammen: wenn eine agentische Lane uebersprungen wird, gibt es auch kein Modell, das sie ausgefuehrt haette - ein leeres executed-by ist dann korrekt und darf nicht erzwungen werden.
definition-of-done: "Ein Move aus einer Lane mit agentic:true fordert executed-by ein oder setzt es; ein Move aus einer nicht-agentischen Lane tut es nicht; die Abnahmeansicht zeigt das Feld; Tests decken beide Faelle ab; go test ./... -race gruen"
tags:
  - gates
blocked-by: []
commits: []
created-at: 2026-09-10T16:53:17Z
updated-at: 2026-09-10T16:53:18Z
follows: 01M263NEZ66D099THPSTD28H7V
updated-by: BeMuCa
---

# executed-by wird gefuellt, wenn eine agentische Lane verlassen wird

## Definition of Done

- [ ] Ein Move aus einer Lane mit agentic:true fordert executed-by ein oder setzt es; ein Move aus einer nicht-agentischen Lane tut es nicht; die Abnahmeansicht zeigt das Feld; Tests decken beide Faelle ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

