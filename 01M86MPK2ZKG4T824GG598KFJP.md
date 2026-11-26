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
updated-at: 2026-11-26T08:02:06Z
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
test-verdict: "pass: go test -race ./... gruen (RC=0, Cache geleert), go vet ok; DoD 1-3 im Baum geprueft, beide Tests rot ohne Fix; Scratch-Board: .lock auf Ticket-Ref, 'jaira set' exit 0, Remote neu, keine Outbox-Leiche; Fehlerpfad = genau der Lock; zwei Schreiber: die naechste Schreibung derselben Session verliert gegen den veralteten lokalen Ref, wird gemeldet (Note edge, vorbestehend)"
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
- **2026-11-26 08:00 · Sascha** — waiting: bughunt (Subagent) und go test -race ./... since 09:00
- **2026-11-26 08:01 · Sascha** — bughunt: change: blocker: core/gitref/gitref.go:554 - class: after a push the remote accepted, a failed local update-ref is swallowed, so the local ref stays behind the remote while every lease reader treats it as 'what the remote last confirmed' (refsync.Record core/refsync/refsync.go:157, RecordDelete :176, RecordMilestone :222, core/refsync/epic.go:67, CheckFresh/isStale); the next write leased from it loses against this clone's own push, outbox.Flush (core/outbox/outbox.go:573-580) settles a ticket entry away as Rejected, and the ref keeps the older state. Same class, second example: a permanent update-ref failure (no write permission on .git/refs, full disk) is now silent on every write, where it used to surface. example: write v1 with the ref lock held -> Sent; lock released; next Record leases '' (local ref still empty); Flush fetches, pushes v2 -> 'the ticket moved since it was read', queue empty, ref carries v1. ran: go test -overlay (scratch test in core/outbox, Queue+Flush with a .lock file, then SHA/Queue/Fetch/Flush) -> second flush rejected, 0 queued, ref = v1. Marked change because line 554 is on the path and core/release/notes/98KFJP.md promises 'the next sync no longer reports it as a lost race against your own write'; the same test against the pre-change gitref.go also rejects v2, so this is the fix stopping one write short, not a regression
- **2026-11-26 08:02 · Sascha** — bughunt: pre-existing: core/gitref/gitref.go:998 - class: any server-side refusal whose stderr contains 'rejected' (pre-receive hook declined, protected ref, permission or quota refusal: '[remote rejected] ...') is classified as ErrRaceLost, and outbox.Flush (core/outbox/outbox.go:573-580) then settles a ticket entry away, so the write is dropped and reported as a lost race instead of kept as Failed; example: bare remote with a pre-receive hook that exits 1 -> outcome rejected, 'the ticket moved since it was read: remote: refs/jaira is protected', 0 entries left; blame: 03961b0ca4 / d3bc99ef64 (2026-09-10, before the ticket); ran: go test -overlay (scratch test in core/outbox with a refusing pre-receive hook)
