---
id: 01M86MPK2ZKG4T824GG598KFJP
title: "gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…"
status: backlog
ready: true
creator: Sascha
context: |-
  Was falsch ist: core/gitref/gitref.go publish (Zeile 542-553) pusht erst, dann setzt es den lokalen Ref mit update-ref. Schlaegt update-ref fehl (z.B. ein fremdes .lock auf dem Ref), gibt publish (commit, err) zurueck.
  Folge: core/outbox/outbox.go Box.Flush sieht den Fehler, meldet Failed und laesst den Eintrag mit der alten Lease stehen, obwohl die Schreibung auf dem Remote ist.
  Der naechste Flush pusht mit der alten Lease gegen die eigene Schreibung, verliert (ErrRaceLost) und meldet Rejected: Ticket-Eintraege werden verworfen, Milestone/Epic in den 'Gewinner' gemergt - beides Laerm fuer eine erfolgreiche Schreibung.
  Vorbild: Delete (gitref.go:558) ignoriert ein fehlgeschlagenes lokales update-ref -d bereits; der Remote ist die Wahrheit, der naechste Fetch zieht den lokalen Ref nach.
  Gefunden: Review zu WV2J2W am 24.11.2026 (note: elsewhere). Vorbestehend.
tags: []
blocked-by: []
related: []
follows: 01M865XBRDFR99ENKNMKWV2J2W
commits: []
created-at: 2026-11-24T11:22:48Z
updated-at: 2026-11-26T07:55:42Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-37950
claimed-at: 2026-11-26T07:54:51Z
goal: "Ein Push, den der Remote angenommen hat, gilt als gesendet, auch wenn danach das lokale update-ref fehlschlaegt: Flush meldet Sent und der Outbox-Eintrag verschwindet."
---

# gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…

## Definition of Done

- [ ] <A checkable statement, readable by someone who was not here>

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:55 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
