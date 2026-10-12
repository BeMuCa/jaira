---
id: 01M4QED6YSC1DXTB8TY5F3VNVC
title: move schreibt eine veraltete Ticket-Datei von der Platte ueber einen neueren Ref
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira move und jaira pull arbeiten nie mit einer Ticket-Datei, die aelter ist als der Ref des Tickets"
context: |-
  Was passiert ist (12.10.2026): 'jaira move K1W2JR --to done --force' auf master. master hing upstream hinterher; die Datei .jaira/tickets/01M263RFDQNPFBZVAYJNK1W2JR-*.md war vom 10.09. (status backlog). Der Ref refs/jaira/tickets/01M263RFDQNPFBZVAYJNK1W2JR stand auf signoff, mit Review, test-verdict, DoD-Proof und 9 Notes (Commit 3290438).
  move las die Datei von der Platte und schrieb sie als d84afdb auf den Ref und nach upstream: Review, Outcome, Proof, Notes weg, assignee zurueck auf BeMuCa. Von Hand wiederhergestellt per Forward-Commit 5ca86a2.
  Die 5 Gate-Verweigerungen (assignee, 3x outcome leer, DoD offen) kamen alle aus der alten Datei - --force hat genau das Symptom verdeckt.

  Die Branch ist nicht sicherer: critique, testing und review committen nach Board-Regel nichts, also hinkt die Ticket-Datei jeder Branch dem Ref um diese Lanes hinterher (K1W2JR: Branch critique, Ref signoff; QTHD8H: Branch in-progress, Ref signoff).
  'jaira pull <id>' hilft nicht: es ist ein No-op, sobald irgendeine Datei da ist (core/refsync/refsync.go:733), es liest den Ref dann gar nicht.
  Ohne Datei verweigert move ('on its ref and not on your disk; pull it first').
  Umgangen wurde es mit einem Skript, das die Datei mit dem Inhalt des Refs ueberschreibt und erst dann move/logbook ruft.
definition-of-done: "move auf eine Datei, die aelter ist als ihr Ref, schreibt nicht auf den Ref, sondern verweigert (auch mit --force) und nennt beide Staende; ein Test pinnt das"
tags:
  - cli
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T06:57:07Z
updated-at: 2026-10-12T07:34:34Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37581
claimed-at: 2026-10-12T07:24:00Z
outcome-what: "Jeder Schreibweg (move auch mit --force, note, dod, claim, set, logbook, archive) verweigert eine Ticket-Datei, deren updated-at vor dem lokalen Ref liegt, und nennt beide Staende; jaira pull bringt eine vorhandene aeltere Datei auf den Stand des Refs."
outcome-why: "move las die veraltete Datei und schrieb sie mit gueltigem Lease auf den Ref (K1W2JR: Review, Proofs, 9 Notes weg); pull war ein No-op, sobald eine Datei existierte, konnte den Zustand also nicht reparieren."
outcome-resolves: "F3VNVC: eine Datei aelter als ihr Ref kann den Ref nicht mehr ueberschreiben, und pull ist der Weg, sie nachzuziehen."
---

# move schreibt eine veraltete Ticket-Datei von der Platte ueber einen neueren Ref

## Definition of Done

- [x] move auf eine Datei, die aelter ist als ihr Ref, schreibt nicht auf den Ref, sondern verweigert (auch mit --force) und nennt beide Staende; ein Test pinnt das
  proof: core/refsync/refsync.go:727 CheckFresh, core/move/move.go:95 vor dem Gate, core/ticket/store.go:1300 in Mutate; TestMoveRefusesAFileOlderThanItsRefEvenWithForce
- [x] jaira pull bringt eine vorhandene, aeltere Datei auf den Stand des Refs, statt still nichts zu tun; ein Test pinnt das
  proof: core/refsync/refsync.go:787 refresh; TestPullBringsAStaleFileUpToItsRef, TestPullLeavesAStaleFileWithAnUnsentWriteAlone
- [x] Zeile in core/release/NOTES.md; go test ./... -race gruen
  proof: core/release/NOTES.md:17; go test ./... -race -count=1 exit 0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Syncer: lokalen Ref lesen, updated_at mit der Datei vergleichen, beide Staende zurueckgeben
- [x] Store.Mutate fragt den Recorder vor fn und verweigert mit ErrStale (beide Staende, Hinweis auf jaira pull)
- [x] Syncer.Pull: vorhandene aeltere Datei auffrischen, Outbox-Pending verweigert
- [x] Tests: move --force auf alter Datei verweigert, Ref unveraendert; pull frischt auf
- [x] NOTES.md-Zeile, go test ./... -race

## Progress
- **2026-10-12 07:24 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 07:25 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan: keine, die nur Alex treffen kann. Selbst entschieden, jeweils mit Grund:
1. 'Aelter' heisst: updated_at der Datei liegt strikt vor updated_at des lokalen Refs. Der lokale Ref ist genau der Stand, gegen den der Lease beim Push prueft (Syncer.Record), und kostet keinen Netzwerkweg. Byte-Vergleich mit der Ref-Historie verworfen: eine Datei, die nie so auf dem Ref stand (K1W2JR: Stand vom 10.09., Ref erst spaeter angelegt), faellt da durch.
2. Die Sperre sitzt in Store.Mutate, nicht nur in move: jeder Schreibweg (note, dod, claim, set, move) liest die Datei und schreibt sie auf den Ref, also verliert jeder denselben Stand. Mutate fragt den Recorder ueber ein optionales Interface; --force kommt dort gar nicht an, also greift die Sperre auch mit --force.
3. pull auf eine vorhandene Datei: fetch, dann Ref-updated_at vergleichen. Ref neuer -> derselbe Weg wie ein frischer pull (Ref uebernehmen, Datei mit den Ref-Bytes ueberschreiben; alter Dateiname weg, falls der Titel sich geaendert hat). Ref gleich alt oder aelter -> wie bisher AlreadyHere. Liegt fuer das Ticket ein ungesendeter Schreibvorgang im Outbox, verweigert pull das Ueberschreiben statt lokale Arbeit zu verwerfen.
- **2026-10-12 07:32 · Alexander Sacharov** — Umsetzung, was der Code nicht sagt:
- Pruefung auch in Store.Archive und Store.Logbook: beide schreiben per recordFiled den Dateistand auf den Ref, also derselbe Verlust. Delete bewusst nicht: Loeschen ist Absicht.
- move.Move prueft VOR dem Gate, nicht nur in Mutate: sonst listet das Gate die Luecken der alten Datei (bei K1W2JR die 5 Verweigerungen) und --force ueberstimmt sie, bevor Mutate verweigert.
- pull-Refresh schreibt den Ref nicht und aendert den Claim nicht; er laeuft auch offline gegen den lokalen Ref (fetch nur best effort), damit wiederholtes pull offline weiter geht.
- pull-Refresh nimmt keine Store-Sperre (Store.lock ist privat); bei gleichzeitigem Mutate desselben Tickets koennte es sich ueberschneiden. Fuer einen Handbefehl hingenommen.
- Ein Ref ohne updated-at gibt kein Urteil (Store.Create setzt es nicht von selbst, die CLI schon) - dann bleibt alles wie vorher.
- Test-Falle: runCLI ruft afterCommand nicht, die Ref-Schreibvorgaenge bleiben im Outbox; flushRefs() im Test aufrufen.
