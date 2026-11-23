---
id: 01M8543DKBZA0P1CHFGKMMKENA
title: "refsync: der Fehler von Box.List beim Lesen des Outbox wird verschluckt"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Scheitert das Auflisten des Outbox, meldet refsync das als Fehler statt so zu tun, als waere der Outbox leer."
context: |-
  Was falsch ist: core/refsync/refsync.go:847 'entries, _ = y.Box.List()' wirft den Fehler weg und arbeitet mit einer leeren Liste weiter. Ein voruebergehender Lesefehler sieht dann aus wie 'nichts zu senden' - ausstehende Ticket-Aenderungen werden in diesem Lauf still nicht gepusht.
  Gleiche Fehlerart wie RGFJ0Q (Departed verschluckte den Fehler von Repo.List); gefunden im Review zu RGFJ0Q am 23.11.2026, gab es schon vorher.
  Erst nach RGFJ0Q starten: beide aendern core/refsync/refsync.go.
definition-of-done: "Ein Fehler von Box.List wird zurueckgegeben bzw. auf stderr gemeldet, statt als leerer Outbox zu gelten; Test mit einem scheiternden Box.List"
tags:
  - concurrency
blocked-by: []
related: []
follows: 01M851Y9VNEV64RE6FAXRGFJ0Q
commits: []
created-at: 2026-11-23T21:13:28Z
updated-at: 2026-11-23T21:13:28Z
---

# refsync: der Fehler von Box.List beim Lesen des Outbox wird verschluckt

## Definition of Done

- [ ] Ein Fehler von Box.List wird zurueckgegeben bzw. auf stderr gemeldet, statt als leerer Outbox zu gelten; Test mit einem scheiternden Box.List
- [ ] Release-Note in core/release/notes/<ID>.md, falls Nutzer eine neue Meldung sehen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

