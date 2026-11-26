---
id: 01M86MPK2ZKG4T824GG598KFJP
title: "gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…"
status: testing
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
updated-at: 2026-11-26T08:00:03Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-37950
claimed-at: 2026-11-26T07:54:51Z
goal: "Ein Push, den der Remote angenommen hat, gilt als gesendet, auch wenn danach das lokale update-ref fehlschlaegt: Flush meldet Sent und der Outbox-Eintrag verschwindet."
outcome-what: "gitref publish meldet nach angenommenem Push keinen Fehler mehr, wenn das lokale update-ref scheitert; der lokale Ref wird beim naechsten Fetch nachgezogen."
outcome-why: "Der Remote ist die Wahrheit; ein Fehler liess die Outbox eine angekommene Schreibung mit alter Lease behalten und erneut senden, wo sie gegen sich selbst verlor. Delete und PushBranch behandeln update-ref schon so."
outcome-resolves: "DoD 1-3: gitref_test TestAnAcceptedPushIsAWriteEvenWhenTheLocalRefIsLocked, outbox_test TestAnAcceptedPushIsSentEvenWhenTheLocalRefIsLocked, Release-Note."
review-summary: none
review-gaps: "none: publish nutzt jetzt dasselbe r.run-Muster wie Delete/PushBranch; das Anlegen des .lock steht in gitref_test und outbox_test doppelt - zwei Pakete, ein gemeinsamer Helfer waere ein Testpaket fuer 6 Zeilen, gelassen"
---

# gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…

## Definition of Done

- [x] Ein Test haelt den lokalen Ref mit einem .lock fest: Write liefert nach angenommenem Push die SHA ohne Fehler, und der Remote traegt den Inhalt
  proof: core/gitref/gitref_test.go TestAnAcceptedPushIsAWriteEvenWhenTheLocalRefIsLocked (ticket/milestone/migration/epic; rot ohne Fix)
- [x] Ein Outbox-Test: Flush meldet in diesem Fall Sent und der Eintrag ist danach nicht mehr in der Outbox
  proof: core/outbox/outbox_test.go TestAnAcceptedPushIsSentEvenWhenTheLocalRefIsLocked (rot ohne Fix)
- [x] Release-Note in core/release/notes/98KFJP.md
  proof: core/release/notes/98KFJP.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:55 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: publish gibt nach angenommenem Push nil zurueck statt den update-ref-Fehler - wie Delete es mit update-ref -d schon tut. Kein Retry: der lokale Ref ist nur ein Cache des Remote, der naechste Fetch zieht ihn nach. Offene Entscheidungen fuer die Person: keine.
- **2026-11-26 07:58 · Sascha** — Geschwister geprueft: Delete (update-ref -d) und PushBranch ignorieren das lokale update-ref schon; publish war der einzige Schreibpfad, der es als Fehler meldete. Bekannte Grenze: bleibt der lokale Ref wegen des Locks stehen, liest eine Schreibung derselben Session vor dem naechsten Fetch die alte SHA als Lease und verliert - das war vorher genauso, nur kommt jetzt die erste Schreibung nicht mehr als Failed zurueck. Testweg: ein .lock-File neben dem Ref unter .git/ laesst update-ref scheitern, ohne den Push zu stoeren.
- **2026-11-26 07:58 · Sascha** — critique: nichts. Fix sitzt im Callee (publish), alle vier Schreibpfade laufen hindurch; Muster wie Delete/PushBranch in gitref.go. Kein einfacherer Schnitt.
- **2026-11-26 08:00 · Sascha** — testing: note: edge: lokaler Ref bleibt nach Lock stehen; Beispiel: Scratch-Board, .lock auf den Ticket-Ref, 'jaira set' -> exit 0, Remote neu, lokal alt; Lock weg, naechstes 'jaira set' -> 'was not sent: ada wrote it first' (gemeldet, Datei unveraendert, kein Datenverlust). Vorbestehend (vorher landete die erste Schreibung als Failed in der Outbox und verlor genauso). Fix waere: Lease nach einem Rejected gegen den Remote-Tip pruefen, ob der eigene Inhalt schon dort liegt - neue Regel, nicht dieses Ticket.
