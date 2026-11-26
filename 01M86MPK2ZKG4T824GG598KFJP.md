---
id: 01M86MPK2ZKG4T824GG598KFJP
title: "gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…"
status: in-progress
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
updated-at: 2026-11-26T07:56:29Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-37950
claimed-at: 2026-11-26T07:54:51Z
goal: "Ein Push, den der Remote angenommen hat, gilt als gesendet, auch wenn danach das lokale update-ref fehlschlaegt: Flush meldet Sent und der Outbox-Eintrag verschwindet."
---

# gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…

## Definition of Done

- [ ] Ein Test haelt den lokalen Ref mit einem .lock fest: Write liefert nach angenommenem Push die SHA ohne Fehler, und der Remote traegt den Inhalt
- [ ] Ein Outbox-Test: Flush meldet in diesem Fall Sent und der Eintrag ist danach nicht mehr in der Outbox
- [ ] Release-Note in core/release/notes/98KFJP.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:55 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: publish gibt nach angenommenem Push nil zurueck statt den update-ref-Fehler - wie Delete es mit update-ref -d schon tut. Kein Retry: der lokale Ref ist nur ein Cache des Remote, der naechste Fetch zieht ihn nach. Offene Entscheidungen fuer die Person: keine.
